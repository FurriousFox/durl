package test

import (
	"crypto/tls"
	"net"
	"net/url"
	"sync"
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

		model.Mu.Lock()
		state[ip.String()]["test"] = []any{Failed}
		model.Mu.Unlock()

		// skip tls/http, as tcp failed
		return
	} else {
		model.Mu.Lock()
		state[ip.String()]["tcp"] = true
		model.Mu.Unlock()
		conn.Close()
	}

	var wg sync.WaitGroup

	// try tls 1.0
	wg.Go(func() {
		var status = Tls(url, address, tls.VersionTLS10)
		model.Mu.Lock()
		if status.State == true {
			state[ip.String()]["tls_10"] = true
		} else {
			state[ip.String()]["tls_10"] = []any{status.State, status.Msg}
		}
		model.Mu.Unlock()
		model.TriggerUpdate()
	})

	// tls 1.1
	wg.Go(func() {
		var status = Tls(url, address, tls.VersionTLS11)
		model.Mu.Lock()
		if status.State == true {
			state[ip.String()]["tls_11"] = true
		} else {
			state[ip.String()]["tls_11"] = []any{status.State, status.Msg}
		}
		model.Mu.Unlock()
		model.TriggerUpdate()
	})

	// tls 1.2
	wg.Go(func() {
		var status = Tls(url, address, tls.VersionTLS12)
		model.Mu.Lock()
		if status.State == true {
			state[ip.String()]["tls_12"] = true
		} else {
			state[ip.String()]["tls_12"] = []any{status.State, status.Msg}
		}
		model.Mu.Unlock()
		model.TriggerUpdate()
	})

	// tls 1.3
	wg.Go(func() {
		var status = Tls(url, address, tls.VersionTLS13)
		model.Mu.Lock()
		if status.State == true {
			state[ip.String()]["tls_13"] = true
		} else {
			state[ip.String()]["tls_13"] = []any{status.State, status.Msg}
		}
		model.Mu.Unlock()
		model.TriggerUpdate()
	})

	// try http 1.0

	// try http 1.1
	wg.Go(func() {
		var status = Http_11(url, address)
		model.Mu.Lock()
		state[ip.String()]["http_11"] = []any{status.State, status.Msg}
		model.Mu.Unlock()
		model.TriggerUpdate()
	})

	// try http 2
	wg.Go(func() {
		var status = Http_2(url, address)
		model.Mu.Lock()
		state[ip.String()]["http_20"] = []any{status.State, status.Msg}
		model.Mu.Unlock()
		model.TriggerUpdate()
	})

	// try http 3
	wg.Go(func() {
		var status = Http_3(url, address)
		model.Mu.Lock()
		state[ip.String()]["http_30"] = []any{status.State, status.Msg}
		model.Mu.Unlock()
		model.TriggerUpdate()
	})

	// certificate info

	wg.Wait()
	model.Mu.Lock()

	if value, ok := model.State[ip.String()]["http_11"].([]any); ok && value[0] == true {
		state[ip.String()]["test"] = []any{Success}
	} else if value, ok := model.State[ip.String()]["http_20"].([]any); ok && value[0] == true {
		state[ip.String()]["test"] = []any{Success}
	} else if value, ok := model.State[ip.String()]["http_30"].([]any); ok && value[0] == true {
		state[ip.String()]["test"] = []any{Success}
	} else {
		state[ip.String()]["test"] = []any{Failed}
	}

	model.Mu.Unlock()

	model.TriggerUpdate()
}
