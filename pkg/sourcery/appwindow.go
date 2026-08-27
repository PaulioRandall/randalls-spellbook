/*
	Source commit: https://github.com/crgimenes/glaze/blob/1860ae968677452cc1c30bcfd9306d31e39083ed/appwindow.go

	I want to bind functions to WebView but
	AppWindow(opts AppOptions), which is implemented by
	appwindow.go in Glaze's source, manages the entire
	lifecycle of the WebView.

	A copied and modified appwindow.go from Glaze and added
	'OnWebViewReady' option. The function accepts a Glaze
	WebView instance and returns an error. If the returned
	error is not nil then it will immediately be returned
	from the AppWindow function without starting the WebView.
*/

/*
	MIT License

	Copyright (c) 2025 Adam Bouqdib (original go-webview base)
	Copyright (c) 2026 crgimenes

	Permission is hereby granted, free of charge, to any person obtaining a copy
	of this software and associated documentation files (the "Software"), to deal
	in the Software without restriction, including without limitation the rights
	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
	copies of the Software, and to permit persons to whom the Software is
	furnished to do so, subject to the following conditions:

	The above copyright notice and this permission notice shall be included in all
	copies or substantial portions of the Software.

	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
	SOFTWARE.
*/

package sourcery

import (
	"github.com/crgimenes/glaze"

	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"runtime"
	"time"
)

// AppTransport selects how AppWindow serves HTTP to the embedded browser.
type AppTransport string

const (
	// AppTransportAuto chooses the recommended platform default.
	// - macOS/Linux: unix backend socket with loopback HTTP gateway.
	// - Windows: loopback TCP.
	AppTransportAuto AppTransport = "auto"

	// AppTransportTCP serves directly over loopback TCP.
	AppTransportTCP AppTransport = "tcp"

	// AppTransportUnix serves the application handler over a Unix domain socket.
	// A lightweight loopback HTTP gateway is created so the embedded browser can
	// still navigate with a standard http:// URL.
	AppTransportUnix AppTransport = "unix"
)

// AppReadyInfo contains transport details once AppWindow listeners are ready.
type AppReadyInfo struct {
	// URL is the navigable URL used by the embedded browser.
	URL string

	// Transport is the resolved backend transport in use.
	Transport AppTransport

	// Backend is the backend listener endpoint.
	// - tcp: "ip:port"
	// - unix: "/path/to/socket"
	Backend string

	// Gateway is the loopback gateway endpoint when unix transport is used.
	// For tcp transport this matches Backend.
	Gateway string
}

// AppOptions configures an AppWindow.
type AppOptions struct {
	// Title is the window title.
	Title string

	// Width and Height set the initial window dimensions.
	Width  int
	Height int

	// Hint controls window resize behaviour (HintNone, HintMin, HintMax, HintFixed).
	Hint glaze.Hint

	// Debug enables the browser developer tools.
	Debug bool

	// Transport selects the backend transport.
	// Defaults to AppTransportAuto.
	Transport AppTransport

	// Addr is the listen address for the local HTTP server.
	// Used by AppTransportTCP and defaults to "127.0.0.1:0".
	Addr string

	// UnixSocketPath is an optional socket path used when Transport is unix.
	// If empty, a temporary socket path is generated automatically.
	UnixSocketPath string

	// Handler is the HTTP handler to serve (typically an http.ServeMux).
	Handler http.Handler

	// OnWebViewReady is called just before the WebView
	// instance is run, allowing function binding etc.
	OnWebViewReady func(webview glaze.WebView) error

	// OnReady is called once listeners are up, with the navigable base URL.
	// Use it to log the address or perform additional setup.
	OnReady func(addr string)

	// OnReadyInfo is called once listeners are up, with transport details.
	// This is useful to inspect whether backend transport is tcp or unix.
	OnReadyInfo func(info AppReadyInfo)
}

// AppWindow creates a native window backed by a local HTTP server.
//
// It starts the server on a random loopback port (or the address specified
// in opts.Addr), opens a webview pointing to it, and runs the UI event loop.
// When the user closes the window, the server is shut down and AppWindow returns.
//
// This is the recommended way to wrap a full devengine application as a
// desktop app — pass the configured http.ServeMux as opts.Handler and
// everything (templates, assets, routes) works unmodified.
func AppWindow(opts AppOptions) error {
	if opts.Handler == nil {
		return fmt.Errorf("webview: AppOptions.Handler must not be nil")
	}
	if opts.Width <= 0 {
		opts.Width = 1024
	}
	if opts.Height <= 0 {
		opts.Height = 768
	}
	if opts.Title == "" {
		opts.Title = "App"
	}

	setup, err := setupAppTransport(opts)
	if err != nil {
		return err
	}
	defer func() {
		if setup.close != nil {
			_ = setup.close()
		}
	}()

	// Start extra transport components (for example, Unix loopback gateway).
	setup.start()

	// Start the application HTTP server in the background.
	srv := &http.Server{Handler: opts.Handler, ReadHeaderTimeout: 10 * time.Second}
	defer func() { _ = srv.Close() }()
	go func() { _ = srv.Serve(setup.listener) }()

	if opts.OnReady != nil {
		opts.OnReady(setup.baseURL)
	}
	if opts.OnReadyInfo != nil {
		opts.OnReadyInfo(AppReadyInfo{
			URL:       setup.baseURL,
			Transport: setup.transport,
			Backend:   setup.backend,
			Gateway:   setup.gateway,
		})
	}

	// Create the webview window.
	w, err := glaze.New(opts.Debug)
	if err != nil {
		return fmt.Errorf("webview: %w", err)
	}

	defer w.Destroy()

	w.SetTitle(opts.Title)
	w.SetSize(opts.Width, opts.Height, opts.Hint)
	w.Navigate(setup.baseURL)

	if opts.OnWebViewReady != nil {
		err = opts.OnWebViewReady(w)
		if err != nil {
			return err
		}
	}

	w.Run()
	return nil
}

type appTransportSetup struct {
	listener  net.Listener
	baseURL   string
	transport AppTransport
	backend   string
	gateway   string
	start     func()
	close     func() error
}

func setupAppTransport(opts AppOptions) (appTransportSetup, error) {
	transport, err := resolveAppTransport(opts.Transport, runtime.GOOS)
	if err != nil {
		return appTransportSetup{}, err
	}

	switch transport {
	case AppTransportTCP:
		return setupTCPTransport(opts.Addr)
	case AppTransportUnix:
		return setupUnixTransport(opts.UnixSocketPath)
	default:
		return appTransportSetup{}, fmt.Errorf("webview: unsupported transport %q", transport)
	}
}

func resolveAppTransport(requested AppTransport, goos string) (AppTransport, error) {
	switch requested {
	case "", AppTransportAuto:
		if goos == "windows" {
			return AppTransportTCP, nil
		}
		return AppTransportUnix, nil
	case AppTransportTCP:
		return AppTransportTCP, nil
	case AppTransportUnix:
		if goos == "windows" {
			return "", errors.New("webview: unix transport is not supported on windows")
		}
		return AppTransportUnix, nil
	default:
		return "", fmt.Errorf("webview: invalid transport %q", requested)
	}
}

func setupTCPTransport(addr string) (appTransportSetup, error) {
	if addr == "" {
		addr = "127.0.0.1:0"
	}

	// Validate that the requested address resolves to loopback only.
	// Desktop app HTTP handlers must not be exposed on external interfaces.
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return appTransportSetup{}, fmt.Errorf("webview: invalid listen address %q: %w", addr, err)
	}
	ip := net.ParseIP(host)
	if ip != nil && !ip.IsLoopback() {
		return appTransportSetup{}, fmt.Errorf("webview: refusing to listen on non-loopback address %q; use 127.0.0.1 or [::1]", addr)
	}
	// Also reject wildcard addresses like "" or "0.0.0.0" or "::".
	if ip == nil || ip.IsUnspecified() {
		return appTransportSetup{}, fmt.Errorf("webview: refusing to listen on wildcard address %q; use 127.0.0.1 or [::1]", addr)
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return appTransportSetup{}, fmt.Errorf("webview: listen %s: %w", addr, err)
	}

	tcpAddr, ok := ln.Addr().(*net.TCPAddr)
	if !ok {
		_ = ln.Close()
		return appTransportSetup{}, errors.New("webview: failed to read tcp listen address")
	}

	return appTransportSetup{
		listener: ln,
		// Build the URL from the address actually bound (tcpAddr.String wraps an
		// IPv6 host in brackets), so an "[::1]:0" listener is reached at
		// http://[::1]:port rather than a non-existent 127.0.0.1 listener.
		baseURL:   "http://" + tcpAddr.String(),
		transport: AppTransportTCP,
		backend:   tcpAddr.String(),
		gateway:   tcpAddr.String(),
		start:     func() {},
		close:     nil,
	}, nil
}

func setupUnixTransport(socketPath string) (appTransportSetup, error) {
	path, err := prepareUnixSocketPath(socketPath)
	if err != nil {
		return appTransportSetup{}, err
	}

	unixListener, err := net.Listen("unix", path)
	if err != nil {
		_ = removeUnixSocket(path)
		return appTransportSetup{}, fmt.Errorf("webview: listen unix %s: %w", path, err)
	}

	proxyListener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		_ = unixListener.Close()
		_ = removeUnixSocket(path)
		return appTransportSetup{}, fmt.Errorf("webview: listen tcp gateway: %w", err)
	}

	proxyURL := &url.URL{Scheme: "http", Host: "unix"}
	proxy := httputil.NewSingleHostReverseProxy(proxyURL)
	proxy.Transport = &http.Transport{
		DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
			var dialer net.Dialer
			return dialer.DialContext(ctx, "unix", path)
		},
	}
	proxyServer := &http.Server{Handler: proxy, ReadHeaderTimeout: 10 * time.Second}

	tcpAddr, ok := proxyListener.Addr().(*net.TCPAddr)
	if !ok {
		_ = proxyServer.Close()
		_ = proxyListener.Close()
		_ = unixListener.Close()
		_ = removeUnixSocket(path)
		return appTransportSetup{}, errors.New("webview: failed to read tcp gateway address")
	}

	return appTransportSetup{
		listener:  unixListener,
		baseURL:   fmt.Sprintf("http://127.0.0.1:%d", tcpAddr.Port),
		transport: AppTransportUnix,
		backend:   path,
		gateway:   tcpAddr.String(),
		start: func() {
			go func() { _ = proxyServer.Serve(proxyListener) }()
		},
		close: func() error {
			_ = proxyServer.Close()
			_ = proxyListener.Close()
			return removeUnixSocket(path)
		},
	}, nil
}

func prepareUnixSocketPath(socketPath string) (string, error) {
	path := socketPath
	if path == "" {
		tmpFile, err := os.CreateTemp("", "glaze-*.sock")
		if err != nil {
			return "", fmt.Errorf("webview: create temp socket path: %w", err)
		}
		path = tmpFile.Name()
		err = tmpFile.Close()
		if err != nil {
			return "", fmt.Errorf("webview: close temp file: %w", err)
		}
		err = os.Remove(path)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("webview: remove temp file %s: %w", path, err)
		}
		return path, nil
	}

	err := removeUnixSocket(path)
	if err != nil {
		return "", err
	}
	return path, nil
}

func removeUnixSocket(path string) error {
	if path == "" {
		return nil
	}
	info, err := os.Lstat(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("webview: stat unix socket %s: %w", path, err)
	}
	if info.Mode()&os.ModeSocket == 0 {
		return fmt.Errorf("webview: %s exists and is not a unix socket", path)
	}
	err = os.Remove(path)
	if err != nil {
		return fmt.Errorf("webview: remove unix socket %s: %w", path, err)
	}
	return nil
}
