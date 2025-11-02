package tester

import (
	"crypto/tls"
	"net"
	"net/url"
	"sync"
	"time"

	ui "argv.nl/durl/internal/app"
	"argv.nl/durl/internal/test"
)

func Test(url *url.URL, ip net.IP, port string, model *ui.Model) {
	var state = model.State

	model.Mu.Lock()
	state[ip.String()]["test"] = test.Status{State: test.Pending}
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
		state[ip.String()]["tcp"] = test.Status{State: test.Failed, Msg: dial_err.Error()}

		state[ip.String()]["tls_10"] = test.Status{State: test.Failed}
		state[ip.String()]["tls_11"] = test.Status{State: test.Failed}
		state[ip.String()]["tls_12"] = test.Status{State: test.Failed}
		state[ip.String()]["tls_13"] = test.Status{State: test.Failed}

		state[ip.String()]["http_11"] = test.Status{State: test.Failed}
		state[ip.String()]["http_20"] = test.Status{State: test.Failed}
		state[ip.String()]["http_30"] = test.Status{State: test.Failed}

		state[ip.String()]["test"] = test.Status{State: test.Failed}
		model.Mu.Unlock()

		// skip tls/http, as tcp failed
		return
	} else {
		model.Mu.Lock()
		state[ip.String()]["tcp"] = test.Status{State: test.Success}
		model.Mu.Unlock()
		conn.Close()
	}

	var wg sync.WaitGroup

	// try tls 1.0
	wg.Go(func() {
		result := *test.Tls(url, address, tls.VersionTLS10)
		model.Mu.Lock()
		state[ip.String()]["tls_10"] = result
		model.Mu.Unlock()
		model.TriggerUpdate()
	})

	// tls 1.1
	wg.Go(func() {
		result := *test.Tls(url, address, tls.VersionTLS11)
		model.Mu.Lock()
		state[ip.String()]["tls_11"] = result
		model.Mu.Unlock()
		model.TriggerUpdate()
	})

	// tls 1.2
	wg.Go(func() {
		result := *test.Tls(url, address, tls.VersionTLS12)
		model.Mu.Lock()
		state[ip.String()]["tls_12"] = result
		model.Mu.Unlock()
		model.TriggerUpdate()
	})

	// tls 1.3
	wg.Go(func() {
		result := *test.Tls(url, address, tls.VersionTLS13)
		model.Mu.Lock()
		state[ip.String()]["tls_13"] = result
		model.Mu.Unlock()
		model.TriggerUpdate()
	})

	// try http 1.0

	// try http 1.1
	wg.Go(func() {
		result := *test.Http_11(url, address)
		model.Mu.Lock()
		state[ip.String()]["http_11"] = result
		model.Mu.Unlock()
		model.TriggerUpdate()
	})

	// try http 2
	wg.Go(func() {
		result := *test.Http_2(url, address)
		model.Mu.Lock()
		state[ip.String()]["http_20"] = result
		model.Mu.Unlock()
		model.TriggerUpdate()
	})

	// try http 3
	wg.Go(func() {
		result := *test.Http_3(url, address)
		model.Mu.Lock()
		state[ip.String()]["http_30"] = result
		model.Mu.Unlock()
		model.TriggerUpdate()
	})

	// certificate info

	wg.Wait()
	model.Mu.Lock()

	if model.State[ip.String()]["http_11"].State == test.Success || model.State[ip.String()]["http_20"].State == test.Success || model.State[ip.String()]["http_30"].State == test.Success {
		state[ip.String()]["test"] = test.Status{State: test.Success}
	} else {
		state[ip.String()]["test"] = test.Status{State: test.Failed}
	}

	model.Mu.Unlock()

	model.TriggerUpdate()
}
