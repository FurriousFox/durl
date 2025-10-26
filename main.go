package main

import (
	"fmt"
	"net"
	"net/url"
	"os"
)

func main() {
	if len(os.Args[1:]) != 1 {
		fmt.Fprintln(os.Stderr, "use 'durl <url>'")
		return
	}

	// parse url
	var url, error = url.Parse(os.Args[1])
	if error != nil {
		var url2, error2 = url.Parse("https://" + os.Args[1])
		if error2 == nil {
			url = url2
		} else {
			fmt.Fprintln(os.Stderr, "idk error", error)
			return
		}
	}

	if url.Scheme == "" {
		var url2, error2 = url.Parse("https://" + os.Args[1])
		if error2 == nil {
			url = url2
		} else {
			fmt.Fprintln(os.Stderr, "idk error", error)
			return
		}
	} else if url.Scheme != "https" {
		fmt.Fprintf(os.Stderr, "unsupported protocol '%s'\n", url.Scheme)
		return
	}

	if url.Hostname() == "" {
		fmt.Fprintln(os.Stderr, "hostname required")
		return
	}

	fmt.Println(url)
	fmt.Println(url.Hostname())

	var hostname = url.Hostname()

	var ips, dns_error = net.LookupIP(hostname)
	if dns_error != nil {
		fmt.Fprintln(os.Stderr, "idk error", error)
		return
	}

	for _, ip := range ips {
		fmt.Println(ip)
	}

	// try tcp

	// try tls 1.*

	// try http 1
	// try http 2
	// try http 3

	// certificate info
}
