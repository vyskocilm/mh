// +build !windows

package cmd

import(
    "os/user"
    "strconv"

    "golang.org/x/sys/unix"
)

// socketGroup assigns docker group for root user on unix socket
func socketGroup(cfg *mhCfg) error {

    // do nothing for non root
    if unix.Getuid() != 0 {
        return nil
    }
    g, err := user.LookupGroup("docker")
    if err != nil {
        return err
    }
    gid, err := strconv.Atoi(g.Gid)
    if err != nil {
        panic(err)          // this should not happen - group id is supposed to be number
    }
    err = unix.Chown(cfg.unix, 0, gid)
    if err != nil {
        return err
    }

    err = unix.Chmod(cfg.unix, 0666)
    return err
}
