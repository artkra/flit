package lserver

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net"
	"strings"
	"time"
)

const (
	beginToken = "+++idspispopd_____"
)

var serviceTokens = [...]string{
	"2hGXENUIwShosts1D3Z0W7l5KV4FqgYo",
	"oSeEtBW7MRsyncT4PD9Fg6idIbYXUfCn",
	"aoZJnr4pSEpingqhL6Vvb0dgXFM1sWwY",
}

type Server struct {
	Addr          string
	IdleTimeout   time.Duration
	MaxReadBuffer int64
	dispatcher    *Dispatcher
}

type Conn struct {
	net.Conn
	IdleTimeout   time.Duration
	MaxReadBuffer int64
}

func (c *Conn) Write(p []byte) (int, error) {
	c.UpdateDeadline()
	return c.Conn.Write(p)
}

func (c *Conn) Read(b []byte) (int, error) {
	c.UpdateDeadline()
	r := io.LimitReader(c.Conn, c.MaxReadBuffer)
	return r.Read(b)
}

func (c *Conn) UpdateDeadline() {
	idleDeadline := time.Now().Add(c.IdleTimeout)
	c.Conn.SetDeadline(idleDeadline)
}

func NewLServer(t time.Duration, m int64) *Server {
	guestBook := make(map[string]string)
	router := make(map[string]*Conn)

	dispatcher := Dispatcher{
		GuestBook: &guestBook,
		Router:    &router,
	}
	return &Server{
		Addr:          ":15070",
		IdleTimeout:   t * time.Second,
		MaxReadBuffer: m,
		dispatcher:    &dispatcher,
	}
}

func (srv Server) ListenAndServe() error {
	addr := srv.Addr

	log.Printf("Starting server on %v\n", addr)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	defer listener.Close()

	for {
		newConn, err := listener.Accept()
		if err != nil {
			log.Printf("--- error accepting connection %v\n", err)
			continue
		}
		conn := &Conn{
			Conn:          newConn,
			IdleTimeout:   srv.IdleTimeout,
			MaxReadBuffer: srv.MaxReadBuffer,
		}
		conn.SetDeadline(time.Now().Add(conn.IdleTimeout))
		log.Printf("+++ accepted connection from %v\n", conn.RemoteAddr())
		go handle(conn, srv.dispatcher)
	}
}

func lSplit(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	var payloadSize = len(data)

	if payloadSize >= 50 {
		for j := 0; j < payloadSize; j++ {
			if bytes.Equal(data[j:j+18], []byte(beginToken)) {
				return j + 1, data[j+18:], nil
			}
		}
		return payloadSize, nil, nil
	}

	return 0, nil, nil
}

func handle(conn *Conn, dispatcher *Dispatcher) error {
	// parse message, add record to dispatcher if not already there, check for receiver and send it to him
	defer func() {
		log.Printf("~ closing connection from %v\n", conn.RemoteAddr())
		conn.Close()
	}()

	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)

	scanr := bufio.NewScanner(r)
	scanr.Split(lSplit)

	for {
		scanned := scanr.Scan()

		if !scanned {
			if err := scanr.Err(); err != nil {
				log.Printf("%v(%v)", err, conn.RemoteAddr())
				return err
			}
			break
		}

		Dispatch(dispatcher, conn, scanr.Bytes())

		w.WriteString(strings.ToUpper(scanr.Text()) + "\n")
		w.Flush()
	}
	return nil
}
