package main

import "github.com/nic758/bdtp-golang/bdtp"

func main() {
	//Be sure that env is set.
	bdtp.NewServer("4444")
}
