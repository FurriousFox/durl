package tester

import (
	"crypto/tls"
	"net"
	"net/url"
	"sync"
	"time"

	ui "argv.nl/durl/internal/app"
	. "argv.nl/durl/internal/test" //lint:ignore ST1001 intentional dot import
)

func Test(url *url.URL, ip net.IP, port string, model *ui.Model) {
	var state = model.State

	model.Mu.Lock()
	state[ip.String()]["test"] = Status{State: Pending}
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
		state[ip.String()]["tcp"] = Status{State: Failed, Msg: dial_err.Error()}

		state[ip.String()]["tls_10"] = Status{State: Failed}
		state[ip.String()]["tls_11"] = Status{State: Failed}
		state[ip.String()]["tls_12"] = Status{State: Failed}
		state[ip.String()]["tls_13"] = Status{State: Failed}

		state[ip.String()]["http_11"] = Status{State: Failed}
		state[ip.String()]["http_2"] = Status{State: Failed}
		state[ip.String()]["http_3"] = Status{State: Failed}

		state[ip.String()]["test"] = Status{State: Failed}
		model.Mu.Unlock()

		// skip tls/http, as tcp failed
		return
	} else {
		model.Mu.Lock()
		state[ip.String()]["tcp"] = Status{State: Success}
		model.Mu.Unlock()
		conn.Close()
	}

	var wg sync.WaitGroup

	// try tls 1.0
	wg.Go(func() {
		result := *Tls(url, address, tls.VersionTLS10)
		model.Mu.Lock()
		state[ip.String()]["tls_10"] = result
		model.Mu.Unlock()
		model.TriggerUpdate()
	})

	// tls 1.1
	wg.Go(func() {
		result := *Tls(url, address, tls.VersionTLS11)
		model.Mu.Lock()
		state[ip.String()]["tls_11"] = result
		model.Mu.Unlock()
		model.TriggerUpdate()
	})

	// tls 1.2
	wg.Go(func() {
		result := *Tls(url, address, tls.VersionTLS12)
		model.Mu.Lock()
		state[ip.String()]["tls_12"] = result
		model.Mu.Unlock()
		model.TriggerUpdate()
	})

	// tls 1.3
	wg.Go(func() {
		result := *Tls(url, address, tls.VersionTLS13)
		model.Mu.Lock()
		state[ip.String()]["tls_13"] = result
		model.Mu.Unlock()
		model.TriggerUpdate()
	})

	// try http 1.0

	// try http 1.1
	wg.Go(func() {
		result := *Http_11(url, address)
		model.Mu.Lock()
		state[ip.String()]["http_11"] = result
		model.Mu.Unlock()
		model.TriggerUpdate()
	})

	// try http 2
	wg.Go(func() {
		result := *Http_2(url, address)
		model.Mu.Lock()
		state[ip.String()]["http_2"] = result
		model.Mu.Unlock()
		model.TriggerUpdate()
	})

	// try http 3
	wg.Go(func() {
		result := *Http_3(url, address)
		model.Mu.Lock()
		state[ip.String()]["http_3"] = result
		model.Mu.Unlock()
		model.TriggerUpdate()
	})

	// certificate info

	wg.Wait()
	model.Mu.Lock()

	if model.State[ip.String()]["http_11"].State == Success || model.State[ip.String()]["http_2"].State == Success || model.State[ip.String()]["http_3"].State == Success {
		state[ip.String()]["test"] = Status{State: Success}
	} else {
		state[ip.String()]["test"] = Status{State: Failed}
	}

	model.Mu.Unlock()

	model.TriggerUpdate()
}
