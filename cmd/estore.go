package cmd

// eStore: code dealing with in-memory store of added DNS entries
// TODO: move out of cmd
import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sync"

	"github.com/spf13/cobra"
)

const (
	defaultGROUP = "$default"
	GROUP_ENV    = "MH_GROUP"
)

var (
	errIPOrNameNotFound = errors.New("IP or Name not found")
	errGroupNotFound    = errors.New("Group not found")
	groupVar            string
)

// add --group/-g flag to the cobra command
func addGroupFlag(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&groupVar, "group", "g", defaultGROUP, "add entry to specific group, enables fast removal")
}

// if --group/-g is not specified, read value from GROUP_ENV
func applyGroupEnv() {
	if groupVar == defaultGROUP {
		if value, ok := os.LookupEnv(GROUP_ENV); ok {
			groupVar = value
		}
	}
}

// e is entry struct
type e struct {
	IP   string
	Name string
}

// eStore - thread unsafe store of e
type eStore struct {
	entries map[string][]e
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
		entries: make(map[string][]e, 0),
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
	for group, entries := range es.entries {
		fmt.Fprintf(f, "### mh group: %s\n", group)
		for _, v := range entries {
			fmt.Fprintf(f, "%s\t%s\n", v.IP, v.Name)
		}
	}
	err = f.Close()
	if err != nil {
		return err
	}

	err = cpa(es.pathwrk, es.path)
	return err
}

// Add write new ip/name tuple
// FIXME: there are three!! strings in API
func (es *eStore) Add(ip, name, group string) {
	e := e{IP: ip, Name: name}
	es.entries[group] = append(es.entries[group], e)
}

// List returns all data from store
func (es *eStore) List() map[string][]e {
	return es.entries
}

// Del removes the first matching entry - either ip either name
func (es *eStore) Del(ipOrName, group string) error {

	idx := -1

	for i, e := range es.entries[group] {
		if e.IP == ipOrName || e.Name == ipOrName {
			idx = i
			break
		}
	}

	if idx != -1 {
		es.entries[group] = append(
			es.entries[group][:idx],
			es.entries[group][idx+1:]...,
		)
		if len(es.entries) == 0 {
			delete(es.entries, group)
		}
		return nil
	}
	return errIPOrNameNotFound
}

// DelGroup removes the group - returns errGroupNotFOund
func (es *eStore) DelGroup(group string) error {
	if _, ok := es.entries[group]; ok {
		delete(es.entries, group)
	}

	return errGroupNotFound
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
func (es *eStoreMx) Add(ip, name, group string) error {
	es.mx.Lock()
	defer es.mx.Unlock()
	es.es.Add(ip, name, group)

	err := es.es.Commit()
	return err
}

// List returns all data from store
func (es *eStoreMx) List() map[string][]e {
	es.mx.RLock()
	defer es.mx.RUnlock()
	return es.es.List()
}

// Del removes the first matching entry - either ip either name
func (es *eStoreMx) Del(ipOrName, group string) error {
	es.mx.Lock()
	defer es.mx.Unlock()

	err := es.es.Del(ipOrName, group)
	if err != nil {
		return err
	}
	err = es.es.Commit()
	if err != nil {
		return err
	}

	return nil
}

// DelGroup removes the group - returns errGroupNotFOund
func (es *eStoreMx) DelGroup(group string) error {
	es.mx.Lock()
	defer es.mx.Unlock()

	err := es.es.DelGroup(group)
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
