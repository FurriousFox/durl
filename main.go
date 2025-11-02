package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"os"
	"regexp"

	ui "argv.nl/durl/internal/app"
	"argv.nl/durl/internal/test"
	"argv.nl/durl/internal/tester"
)

func main() {
	log.SetOutput(io.Discard)

	if len(os.Args[1:]) != 1 {
		fmt.Fprintln(os.Stderr, "use 'durl <url>'")
		return
	}

	// parse url

	var url *url.URL
	var err error
	matched, err := regexp.MatchString(`^[a-zA-Z][a-zA-Z0-9+\-.]*://`, os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "internal exception:", err)
		return
	}
	if matched {
		// if strings.Contains(os.Args[1], ":") {
		url, err = url.Parse(os.Args[1])
	} else {
		url, err = url.Parse("https://" + os.Args[1])
	}
	if err != nil {
		var url2, err2 = url.Parse(os.Args[1])
		if err2 == nil {
			url = url2
		} else {
			fmt.Fprintln(os.Stderr, "invalid url:", err)
			return
		}
	}

	if url.Scheme != "https" {
		fmt.Fprintf(os.Stderr, "unsupported protocol '%s'\n", url.Scheme)
		return
	}

	if url.Hostname() == "" {
		fmt.Fprintln(os.Stderr, "hostname required")
		return
	}

	port := url.Port()
	if port == "" {
		port = "443"
	}

	// fmt.Println(url)
	// fmt.Println(url.Hostname())

	var hostname = url.Hostname()

	var ips, dns_err = net.LookupIP(hostname)
	if dns_err != nil {
		fmt.Fprintln(os.Stderr, "dns lookup error:", dns_err)
		return
	}

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

	ui.RunUI(model)
}
