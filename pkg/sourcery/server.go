package sourcery

import (
	"net"
	"net/http"
	"time"
)

type contentServer struct {
	server   *http.Server
	listener net.Listener
	baseUrl  string
}

func createContentServer(handler http.Handler) (contentServer, error) {
	server := &http.Server{
		// TODO: Needs optimising.
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       20 * time.Second,
		MaxHeaderBytes:    16 << 10, // 16 KiB
	}

	// Create listener on a random loopback port.
	listener, e := net.Listen("tcp4", "127.0.0.1:0")
	if e != nil {
		return contentServer{}, e
	}

	baseUrl := "http://" + listener.Addr().(*net.TCPAddr).String()

	cs := contentServer{
		server:   server,
		listener: listener,
		baseUrl:  baseUrl,
	}

	return cs, nil
}

func (cs *contentServer) serve() error {
	return cs.server.Serve(cs.listener)
}

func (cs *contentServer) close() error {
	return cs.server.Close()
}
