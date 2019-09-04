package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run mh server",
	Long:  `Start mh as HTTP server ready to accept conections`,
	Run: func(cmd *cobra.Command, args []string) {
		hosts := ""
		if len(args) > 0 {
			hosts = args[0]
		}

		// autodetect the location
		if hosts == "" {
			switch runtime.GOOS {
			case "darwin":
				hosts = "/etc/hosts"
			case "dragonfly":
				hosts = "/etc/hosts"
			case "freebsd":
				hosts = "/etc/hosts"
			case "linux":
				hosts = "/etc/hosts"
			case "netbsd":
				hosts = "/etc/hosts"
			case "openbsd":
				hosts = "/etc/hosts"
			case "solaris":
				hosts = "/etc/hosts"
			case "windows":
				hosts = `C:\Windows\System32\drivers\etc`
			default:
				exitOnErr(
					fmt.Errorf("Unknown operating system: %s, please pass hosts file location on a command line", runtime.GOOS))
			}
		}
		startServer(hosts)
	},
}

type mhSvc struct {
	estore *eStoreMx
}

// listEntries returns a list of entries added by server itself
func (s *mhSvc) listEntries(w http.ResponseWriter, r *http.Request) {

	list := s.estore.List()
	json.NewEncoder(w).Encode(list)

}

func (s *mhSvc) addEntry(w http.ResponseWriter, r *http.Request) {
	//TODO: validate hostname
	vars := mux.Vars(r)
	err := s.estore.Add(
		vars["ip"],
		vars["name"],
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "{\"error\": \"%s\"}", err)
		return
	}

	// TODO: more meaingful reply?
	w.WriteHeader(http.StatusNoContent)
}

func (s *mhSvc) delEntry(w http.ResponseWriter, r *http.Request) {
	//TODO: validate hostname
	vars := mux.Vars(r)
	err := s.estore.Del(
		vars["ipOrName"],
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "{\"error\": \"%s\"}", err)
		return
	}

	// TODO: more meaingful reply?
	w.WriteHeader(http.StatusNoContent)
}

// startServer starts HTTP server prividing an interface between commands and hosts file
func startServer(hosts string) {

	estore, err := newEStoreMx(hosts)
	if err != nil {
		exitOnErr(err)
	}

	svc := &mhSvc{
		estore: &estore,
	}

	r := mux.NewRouter()
	apiV1 := r.PathPrefix("/v1").Subrouter()
	apiV1.HandleFunc("/e", svc.listEntries).Methods("GET")
	apiV1.HandleFunc("/e/{ip}/{name}", svc.addEntry).Methods("PUT")
	apiV1.HandleFunc("/e/{ipOrName}", svc.delEntry).Methods("DELETE")

	srv := &http.Server{
		Addr: ":1234",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			fmt.Println(err.Error())
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// cleanup the state, remove temporary files
	estore.Close()

	wait := 15 * time.Second
	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	return

}
