package lserver

import "sync"

// Dispatcher is a map, carrying info about routing
// Router holds a map of type receiver - sender. We extract receiver address from message, find such key in Router
// and write message into it's writer (conn.dispatcher.Router["somereceiver123"].Write(msg))
type Dispatcher struct {
	GuestBook *map[string]string
	Router    *map[string]*Conn
	mux       sync.Mutex
}

func Dispatch(dispatcher *Dispatcher, conn *Conn, msg []byte) {

}
