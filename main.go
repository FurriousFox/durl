package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"

	ui "argv.nl/durl/internal/app"
	"argv.nl/durl/internal/test"
	"argv.nl/durl/internal/tester"
	"argv.nl/durl/internal/util"
)

func main() {
	log.SetOutput(io.Discard)

	if len(os.Args[1:]) != 1 {
		fmt.Fprintln(os.Stderr, "use 'durl <url>'")
		return
	}

	// parse url
	url, err := util.ParseURL(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	port := url.Port()
	if port == "" {
		port = "443"
	}

	// dns
	var hostname = url.Hostname()

	var ips, dns_err = net.LookupIP(hostname)
	if dns_err != nil {
		fmt.Fprintln(os.Stderr, "dns lookup error:", dns_err)
		return
	}

	// run tests
	state := map[string]map[string]test.Status{}
	model := &ui.Model{State: state}

	go func() {
		model.Mu.Lock()
		for _, ip := range ips {
			state[ip.String()] = map[string]test.Status{}
		}
		model.Mu.Unlock()

		for _, ip := range ips {
			go tester.Test(url, ip, port, model)
		}

		// model.Mu.RLock()
		// jsonBytes, err := json.Marshal(state)
		// model.Mu.RUnlock()
		// if err != nil {
		// 	panic(err)
		// }
		// fmt.Println(string(jsonBytes))
	}()

	// ui
	ui.RunUI(model)
}
