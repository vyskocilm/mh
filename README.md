# myhosts: manage /etc/hosts entries via HTTP API

Use case: as a developer I need to create new DNS mappings for a test and clean them up

Possibilities

1. `sudo vim /etc/hosts`: slow, can't be automated, leaves garbage inside
2. `LD_PRELOAD` tricks: cumbersome, works only for code using `glibc`
3. `unshare` tricks: still involve root, hard to automate

## myhosts

```sh
go get github.com/vyskocilm/mh
sudo mh server
```

## add mapping
```
mh add ip name
mh del ip-or-name
mh list
```

## garbage collect!
```
^Ctrl+C
# check that ALL written entries are removed and /etc/hosts is restored
```

## TODO:
7. investigate the way to use Unix sockets like Docker and so does
8. drop all other capabilities!!!
9. how to deal with duplicates? should it be smart?
10. drop data from memory on commit fail!

## concept: transactions (?)

mark more edits and allow removal of more entries atomically

```
mh add -tx 11 ip name
mh del -tx 11
```
