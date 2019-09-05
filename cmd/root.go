package cmd

import (
	"fmt"
	"os"
    "context"
    "errors"
    "strings"
    "strconv"
	"runtime"
    "net/http"
    "net"

	"github.com/spf13/cobra"
)

var (
    host string             // --hosts -H flag
    cfg mhCfg               // parsed configuration
    errUnsupportedHost = errors.New("Unsupported --host format")
)

type mhCfg struct {
    hostsPath    string
    unix        string
    port        string
}

// HttpClient returns new HTTP client based on command line arguments
// solves the TCP vs Unix socket problem
func (c *mhCfg) HTTPClient() http.Client {
    if (c.unix != "") {
        return http.Client{
            Transport: &http.Transport{
                DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
                    return net.Dial("unix", c.unix)
                },
            },
        }
    }
    return http.Client{}
}

func init() {

    cfg = mhCfg{
        hostsPath: "",
        unix: "",
        port: "",
    }

    rootCmd.PersistentFlags().StringVarP(&host, "host", "H", "", `The host to bind/connect to
Supported formats are port number or unix:///path/to/socket`)
}

func parseArgs(cmd *cobra.Command, args []string) error {

    // values on cmdline
    chosts := ""
    cunix := ""
    cport := ""

    // ignore /etc/hosts search for non-server command
    if len(args) > 0 {
        chosts = args[0]
    }

    if strings.HasPrefix(host, "unix://") {
        cunix = host[7:len(host)]
    } else if host != "" {
        _, err := strconv.Atoi(host)
        if err != nil {
            return errUnsupportedHost
        }
        cport = host
    }


    // autodetected values
    ahosts := ""
    aunix := ""
    aport := ""
    switch runtime.GOOS {
    case "darwin":
        fallthrough
    case "dragonfly":
        fallthrough
    case "freebsd":
        fallthrough
    case "linux":
        fallthrough
    case "netbsd":
        fallthrough
    case "openbsd":
        fallthrough
    case "solaris":
        ahosts = "/etc/hosts"
        aunix = "/var/run/mh.sock"
    case "windows":
        ahosts = `C:\Windows\System32\drivers\etc\hosts`
        aport = ":3003"
    default:
        if chosts == "" && (cport == "" || cunix == "") {
            exitOnErr(
                fmt.Errorf("Unknown operating system: %s, please pass hosts file location and --host on a command line", runtime.GOOS))
            }
        }

    if chosts != "" {
        cfg.hostsPath = chosts
    } else {
        cfg.hostsPath = ahosts
    }

    if cport == "" {
        if cunix != "" {
            cfg.unix = cunix
        } else {
            cfg.unix = aunix
        }
    } else if cport != "" {
        cfg.port = cport
    } else {
        if aunix != "" {
            cfg.unix = aunix
        } else {
            cfg.port = aport
        }
    }

    // final assert
    if cfg.unix == "" && cfg.port == "" {
        err := fmt.Errorf("either unix socket part either HTTP port are empty, this should not happen")
        return err
    }

    return nil
}

var rootCmd = &cobra.Command{
	Use:   "mh",
	Short: "mh let you dynamically manage /etc/hosts",
	Long: `Easy to use way of managing /etc/hosts
        with a simple and clean command line interface and garbage collection`,
    Args: func(cmd *cobra.Command, args []string) error {
        return parseArgs(cmd, args)
    },
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

// Execute run mh command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
