package main

// mh: myhosts - small daemon managing /etc/hosts file allowing easy and dynamic change of DNS

import (
	"github.com/vyskocilm/mh/cmd"
)

func main() {
	cmd.Execute()
}
