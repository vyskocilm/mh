# myhosts: manage /etc/hosts entries via CLI/HTTP API

> Note: `mh` is not yet released and provides no guarantee about command line
> or switches or the way it handles duplicates.

**Use case**: as a developer I need to create new DNS mappings for a while from
a script/program

## License

MIT, see `LICENSE` file


## Running a server

```sh
go get github.com/vyskocilm/mh
sudo mh server
# or via systemd to daemonize it
sudo systemd-run --name mh mh server
```

## Manipulate with entries
```
mh add ip name
mh del ip-or-name
mh list
```

## Groups

mh stores entries by the group allowing one to easilly drop all added entries from one group

```sh
mh add --group integration_project_1 ip name
mh add --group integration_project_1 ip2 name2
...
mh delgrp integration_project_1
```

if not specified, entries goes into `$default` group.

It is possible to export `MH_GROUP` variable to ensure `add` and `del` commands
stays indide group.


```sh
export MH_GROUP=integration_project_1

mh add ip1 name1
mh add ip2 name2
mh del ip1

...
mh delgrp
```

## Use with docker

Creates DNS mapping `api.test` to `apimoc_api_1` container

```sh
mh add $(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' apimock_api_1) api.test
```

## Garbage collect!

`mh` restores original `/etc/hosts` upon completion.

```
^Ctrl+C
# or systemctl stop mh
# check that ALL written entries are removed and /etc/hosts is restored
```

## Communication method

`mh` use unix socket `/vat/run/mh.sock` by default on unix systems and port
`:3003` on Windows by default. Can be changes via `-H/--host` command line
flag.

```sh
# listen on /tmp/mh.sock
mh --host unix:///tmp/mh.sock server
# listen on port 1234
mh --host 1234 server
```

If unix socket is specified and `mh` started as a root user, socket is assigned
under `docker` group. This is for convenience as most potential users would
have docker installed and configured.

## TODO
8. drop all other capabilities!!!
9. how to deal with duplicates? should it be smart?
10. drop data from memory on commit fail!
11. allow usage of different group for unix socket
12. remove socket begore server start
13. socket activation
14. HTTPS/TLS/CA/certificates support - browsers will treat http as insecure soonish
    so there must be a way to create TLS ready infrastructure

## Alternatives

Prior to writing `mh` I did experiment with a following methods, however did
not like any of those.

1. `sudo vim /etc/hosts`: slow, can't be automated, leaves garbage inside
2. `LD_PRELOAD` tricks: cumbersome, works only for code using `glibc`
3. `unshare` tricks: still involve root, hard to automate
4. `nss_wrapper`: same `LD_PRELOAD` trick
