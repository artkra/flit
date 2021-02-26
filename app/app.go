package main

import (
	"github.com/artkra/flit/lserver"
)

func main() {

	server := lserver.NewLServer(10, 1024)

	server.ListenAndServe()
}
