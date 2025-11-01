package test

import (
	"crypto/tls"
	"net"
	"net/url"
	"time"

	ui "argv.nl/durl/internal/app"
)

func Test(url *url.URL, ip net.IP, port string, model *ui.Model) {
	var state = model.State

	model.Mu.Lock()
	state[ip.String()]["test"] = []any{Active}
	model.Mu.Unlock()

	var ipstring string

	if ip.To4() != nil {
		ipstring = ip.String()
	} else {
		ipstring = "[" + ip.String() + "]"
	}
	var address = ipstring + ":" + port

	// try tcp
	var conn, dial_err = net.DialTimeout("tcp", address, 5*time.Second)
	if dial_err != nil {
		// fmt.Fprintln(os.Stderr, "tcp dial err", dial_err)

		model.Mu.Lock()
		state[ip.String()]["tcp"] = false

		state[ip.String()]["tls_10"] = []any{false, dial_err.Error()}
		state[ip.String()]["tls_11"] = []any{false, dial_err.Error()}
		state[ip.String()]["tls_12"] = []any{false, dial_err.Error()}
		state[ip.String()]["tls_13"] = []any{false, dial_err.Error()}

		state[ip.String()]["http_11"] = []any{false, dial_err.Error()}
		state[ip.String()]["http_20"] = []any{false, dial_err.Error()}
		state[ip.String()]["http_30"] = []any{false, dial_err.Error()}

		model.Mu.Unlock()
		// skip tls/http, as tcp failed

		model.Mu.Lock()
		state[ip.String()]["test"] = []any{Failed}
		model.Mu.Unlock()

		return
	} else {
		model.Mu.Lock()
		state[ip.String()]["tcp"] = true
		model.Mu.Unlock()
		conn.Close()
	}

	// try tls 1.0
	conn, tls_10_err := tls.DialWithDialer(&net.Dialer{
		Timeout: 5 * time.Second,
	}, "tcp", address, &tls.Config{
		ServerName: url.Hostname(),
		MinVersion: tls.VersionTLS10,
		MaxVersion: tls.VersionTLS10,
	})
	if tls_10_err != nil {
		// fmt.Fprintln(os.Stderr, "tls dial err", tls_10_err)

		model.Mu.Lock()
		state[ip.String()]["tls_10"] = []any{false, tls_10_err.Error()}
		model.Mu.Unlock()
	} else {
		model.Mu.Lock()
		state[ip.String()]["tls_10"] = true
		model.Mu.Unlock()
		conn.Close()
	}

	// tls 1.1
	conn, tls_11_err := tls.DialWithDialer(&net.Dialer{
		Timeout: 5 * time.Second,
	}, "tcp", address, &tls.Config{
		ServerName: url.Hostname(),
		MinVersion: tls.VersionTLS11,
		MaxVersion: tls.VersionTLS11,
	})
	if tls_11_err != nil {
		// fmt.Fprintln(os.Stderr, "tls dial err", tls_11_err)

		model.Mu.Lock()
		state[ip.String()]["tls_11"] = []any{false, tls_11_err.Error()}
		model.Mu.Unlock()
	} else {
		model.Mu.Lock()
		state[ip.String()]["tls_11"] = true
		model.Mu.Unlock()
		conn.Close()
	}

	// tls 1.2
	conn, tls_12_err := tls.DialWithDialer(&net.Dialer{
		Timeout: 5 * time.Second,
	}, "tcp", address, &tls.Config{
		ServerName: url.Hostname(),
		MinVersion: tls.VersionTLS12,
		MaxVersion: tls.VersionTLS12,
	})
	if tls_12_err != nil {
		// fmt.Fprintln(os.Stderr, "tls dial err", tls_12_err)

		model.Mu.Lock()
		state[ip.String()]["tls_12"] = []any{false, tls_12_err.Error()}
		model.Mu.Unlock()
	} else {
		model.Mu.Lock()
		state[ip.String()]["tls_12"] = true
		model.Mu.Unlock()
		conn.Close()
	}

	// tls 1.3
	conn, tls_13_err := tls.DialWithDialer(&net.Dialer{
		Timeout: 5 * time.Second,
	}, "tcp", address, &tls.Config{
		ServerName: url.Hostname(),
		MinVersion: tls.VersionTLS13,
		MaxVersion: tls.VersionTLS13,
	})
	if tls_13_err != nil {
		// fmt.Fprintln(os.Stderr, "tls dial err", tls_13_err)

		model.Mu.Lock()
		state[ip.String()]["tls_13"] = []any{false, tls_13_err.Error()}
		model.Mu.Unlock()
	} else {
		model.Mu.Lock()
		state[ip.String()]["tls_13"] = true
		model.Mu.Unlock()
		conn.Close()
	}

	// try http 1.0

	// try http 1.1
	{
		var status = Http_11(url, address)
		model.Mu.Lock()
		state[ip.String()]["http_11"] = []any{status.State, status.Msg}
		model.Mu.Unlock()
	}

	// try http 2
	{
		var status = Http_2(url, address)
		model.Mu.Lock()
		state[ip.String()]["http_20"] = []any{status.State, status.Msg}
		model.Mu.Unlock()
	}

	// try http 3
	{
		var status = Http_3(url, address)
		model.Mu.Lock()
		state[ip.String()]["http_30"] = []any{status.State, status.Msg}
		model.Mu.Unlock()
	}
	// certificate info

	model.Mu.Lock()
	state[ip.String()]["test"] = []any{Success}
	model.Mu.Unlock()
}
