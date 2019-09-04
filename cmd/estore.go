package cmd

// eStore: code dealing with in-memory store of added DNS entries
// TODO: move out of cmd
import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sync"
)

var (
	errIPOrNameNotFound = errors.New("IP or Name not found")
)

// e is entry struct
type e struct {
	IP   string
	Name string
}

// eStore - thread unsafe store of e
type eStore struct {
	entries []e
	path    string // original hosts
	pathmh  string // mh backup
	pathwrk string // mh working file
}

func cpa(src, dst string) error {
	//FIXME: this is inherently unportable, but works well on Linux/Mac
	cpCmd := exec.Command("cp", "-a", src, dst)
	err := cpCmd.Run()
	return err
}

func newEStore(path string) (eStore, error) {
	pathmh := fmt.Sprintf("%s.mh", path)
	pathwrk := fmt.Sprintf("%s.wrk", pathmh)

	// I gave up trying to check for permissions via Stat and Stat_t
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return eStore{}, err
	}
	f.Close()

	err = cpa(path, pathmh)
	if err != nil {
		return eStore{}, err
	}

	return eStore{
		entries: make([]e, 0),
		path:    path,
		pathmh:  pathmh,
		pathwrk: pathwrk,
	}, nil
}

// Close restores the state back
func (es *eStore) Close() error {
	err := cpa(es.pathmh, es.path)
	os.Remove(es.pathwrk)
	os.Remove(es.pathmh)
	return err
}

// Commit writes all the changes to the disk
// 1. copy backup from start to work file
// 2. append all data there
// 3. copy the file
func (es *eStore) Commit() error {
	err := cpa(es.pathmh, es.pathwrk)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(es.pathwrk, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	fmt.Fprintf(f, "### added by github.com/vyskocilm/mh, entries will be dropped when server stops ###\n")
	for _, v := range es.entries {
		fmt.Fprintf(f, "%s\t%s\n", v.IP, v.Name)
	}
	err = f.Close()
	if err != nil {
		return err
	}

	err = cpa(es.pathwrk, es.path)
	return err
}

// Add write new ip/name tuple
func (es *eStore) Add(ip, name string) {
	e := e{IP: ip, Name: name}
	es.entries = append(es.entries, e)
}

// List returns all data from store
func (es *eStore) List() []e {
	return es.entries
}

// Del removes the first matching entry - either ip either name
func (es *eStore) Del(ipOrName string) error {

	idx := -1

	for i, e := range es.entries {
		if e.IP == ipOrName || e.Name == ipOrName {
			idx = i
			break
		}
	}

	if idx != -1 {
		es.entries = append(
			es.entries[:idx],
			es.entries[idx+1:]...,
		)
		return nil
	}
	return errIPOrNameNotFound
}

// Adds two new options for eStore
// 1. all operations are now guarded via RWMutex
// 2. all write operations are commited by default
//
// newEStoreMx and Close methods are NOT thread-safe by design
type eStoreMx struct {
	mx sync.RWMutex
	es eStore
}

func newEStoreMx(path string) (eStoreMx, error) {
	estore, err := newEStore(path)
	if err != nil {
		return eStoreMx{}, err
	}

	return eStoreMx{es: estore}, nil
}

// Add write new ip/name tuple
func (es *eStoreMx) Add(ip, name string) error {
	es.mx.Lock()
	defer es.mx.Unlock()
	es.es.Add(ip, name)

	err := es.es.Commit()
	return err
}

// List returns all data from store
func (es *eStoreMx) List() []e {
	es.mx.RLock()
	defer es.mx.RUnlock()
	return es.es.List()
}

// Del removes the first matching entry - either ip either name
func (es *eStoreMx) Del(ipOrName string) error {
	es.mx.Lock()
	defer es.mx.Unlock()

	err := es.es.Del(ipOrName)
	if err != nil {
		return err
	}
	err = es.es.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (es *eStoreMx) Close() error {
	return es.es.Close()
}
