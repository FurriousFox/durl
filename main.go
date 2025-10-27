package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	quic "github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"argv.nl/durl/internal/app"
)

func main() {
	if len(os.Args[1:]) != 1 {
		fmt.Fprintln(os.Stderr, "use 'durl <url>'")
		return
	}

	// parse url
	var url, err = url.Parse(os.Args[1])
	if err != nil {
		var url2, err2 = url.Parse("https://" + os.Args[1])
		if err2 == nil {
			url = url2
		} else {
			fmt.Fprintln(os.Stderr, "idk err", err)
			return
		}
	}

	if url.Scheme == "" {
		var url2, err2 = url.Parse("https://" + os.Args[1])
		if err2 == nil {
			url = url2
		} else {
			fmt.Fprintln(os.Stderr, "idk err", err)
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

	port := url.Port()
	if port == "" {
		port = "443"
	}

	fmt.Println(url)
	fmt.Println(url.Hostname())

	var hostname = url.Hostname()

	var ips, dns_err = net.LookupIP(hostname)
	if dns_err != nil {
		fmt.Fprintln(os.Stderr, "dns lookup err", dns_err)
		return
	}

	state := map[string]map[string]any{}

	for _, ip := range ips {
		var ipstring string

		if ip.To4() != nil {
			ipstring = ip.String()
		} else {
			ipstring = "[" + ip.String() + "]"
		}
		var host = ipstring + ":" + port

		state[ip.String()] = map[string]any{}

		fmt.Println(host)

		// try tcp
		var conn, dial_err = net.DialTimeout("tcp", host, 5*time.Second)
		if dial_err != nil {
			fmt.Fprintln(os.Stderr, "tcp dial err", dial_err)

			state[ip.String()]["tcp"] = false
			// skip tls/http, as tcp failed
			continue
		} else {
			state[ip.String()]["tcp"] = true
			conn.Close()
		}

		// try tls 1.0
		conn, tls_10_err := tls.DialWithDialer(&net.Dialer{
			Timeout: 5 * time.Second,
		}, "tcp", host, &tls.Config{
			ServerName: url.Hostname(),
			MinVersion: tls.VersionTLS10,
			MaxVersion: tls.VersionTLS10,
		})
		if tls_10_err != nil {
			fmt.Fprintln(os.Stderr, "tls dial err", tls_10_err)

			state[ip.String()]["tls_10"] = []any{false, tls_10_err.Error()}
		} else {
			state[ip.String()]["tls_10"] = true
			conn.Close()
		}

		// tls 1.1
		conn, tls_11_err := tls.DialWithDialer(&net.Dialer{
			Timeout: 5 * time.Second,
		}, "tcp", host, &tls.Config{
			ServerName: url.Hostname(),
			MinVersion: tls.VersionTLS11,
			MaxVersion: tls.VersionTLS11,
		})
		if tls_11_err != nil {
			fmt.Fprintln(os.Stderr, "tls dial err", tls_11_err)

			state[ip.String()]["tls_11"] = []any{false, tls_11_err.Error()}
		} else {
			state[ip.String()]["tls_11"] = true
			conn.Close()
		}

		// tls 1.2
		conn, tls_12_err := tls.DialWithDialer(&net.Dialer{
			Timeout: 5 * time.Second,
		}, "tcp", host, &tls.Config{
			ServerName: url.Hostname(),
			MinVersion: tls.VersionTLS12,
			MaxVersion: tls.VersionTLS12,
		})
		if tls_12_err != nil {
			fmt.Fprintln(os.Stderr, "tls dial err", tls_12_err)

			state[ip.String()]["tls_12"] = []any{false, tls_12_err.Error()}
		} else {
			state[ip.String()]["tls_12"] = true
			conn.Close()
		}

		// tls 1.3
		conn, tls_13_err := tls.DialWithDialer(&net.Dialer{
			Timeout: 5 * time.Second,
		}, "tcp", host, &tls.Config{
			ServerName: url.Hostname(),
			MinVersion: tls.VersionTLS13,
			MaxVersion: tls.VersionTLS13,
		})
		if tls_13_err != nil {
			fmt.Fprintln(os.Stderr, "tls dial err", tls_13_err)

			state[ip.String()]["tls_13"] = []any{false, tls_13_err.Error()}
		} else {
			state[ip.String()]["tls_13"] = true
			conn.Close()
		}

		// try http 1.0
		// var client = &http.Client{
		// 	Timeout: 5 * time.Second,
		// 	Transport: &http.Transport{
		// 		DisableKeepAlives: true,
		// 		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
		// 			return (&net.Dialer{Timeout: 5 * time.Second}).DialContext(ctx, network, host)
		// 		},
		// 		TLSClientConfig: &tls.Config{
		// 			ServerName: url.Hostname(),
		// 		},
		// 	}}
		// var req, err = http.NewRequest("GET", url.String(), nil)
		// if err != nil {
		// 	panic(err)
		// }
		// req.Proto = "HTTP/1.0"
		// req.ProtoMajor = 1
		// req.ProtoMinor = 0
		// req.Close = true

		// var resp, err2 = client.Do(req)
		// if err2 != nil {
		// 	panic(err2)
		// }
		// resp.Body.Close()
		// fmt.Println("Response status:", resp.Status)

		// try http 1.1
		{
			var client = &http.Client{
				Timeout: 5 * time.Second,
				Transport: &http.Transport{
					DisableKeepAlives: true,
					DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
						return (&net.Dialer{Timeout: 5 * time.Second}).DialContext(ctx, network, host)
					},
					TLSClientConfig: &tls.Config{
						ServerName: url.Hostname(),
					},
				}}
			var req, err = http.NewRequest("GET", url.String(), nil)
			if err != nil {
				panic(err)
			}
			req.Close = true

			var resp, err2 = client.Do(req)
			if err2 != nil {
				state[ip.String()]["http_11"] = []any{false, err2.Error()}
			} else {
				state[ip.String()]["http_11"] = []any{true, resp.Status, resp.StatusCode}
			}
			if resp != nil {
				// scanner := bufio.NewScanner(resp.Body)
				// for i := 0; scanner.Scan() && i < 5; i++ {
				// 	fmt.Println(scanner.Text())
				// }

				resp.Body.Close()
				fmt.Println("Response status:", resp.Status)
			}
		}

		// try http 2
		{
			var protocols = &http.Protocols{}
			protocols.SetHTTP1(false)
			protocols.SetHTTP2(true)

			var client = &http.Client{
				Timeout: 5 * time.Second,
				Transport: &http.Transport{
					DisableKeepAlives: true,
					DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
						return (&net.Dialer{Timeout: 5 * time.Second}).DialContext(ctx, network, host)
					},
					TLSClientConfig: &tls.Config{
						ServerName: url.Hostname(),
					},
					Protocols: protocols,
				}}
			var req, err = http.NewRequest("GET", url.String(), nil)
			if err != nil {
				panic(err)
			}
			req.Close = true

			var resp, err2 = client.Do(req)
			if err2 != nil {
				state[ip.String()]["http_20"] = []any{false, err2.Error()}
			} else {
				state[ip.String()]["http_20"] = []any{true, resp.Status, resp.StatusCode}
			}
			if resp != nil {
				resp.Body.Close()
				fmt.Println("Response status:", resp.Status)
			}
		}

		// try http 3
		{
			udpConn, err3 := net.ListenUDP("udp", &net.UDPAddr{})
			if err3 != nil {
				state[ip.String()]["http_30"] = []any{false, err3.Error()}
			} else {
				quic_tr := &quic.Transport{Conn: udpConn}

				tr := &http3.Transport{
					TLSClientConfig: &tls.Config{
						NextProtos: []string{http3.NextProtoH3},
						ServerName: url.Hostname(),
					},
					QUICConfig: &quic.Config{},
					Dial: func(ctx context.Context, addr string, tlsCfg *tls.Config, cfg *quic.Config) (*quic.Conn, error) {
						var a, err = net.ResolveUDPAddr("udp", host)
						if err != nil {
							panic(err)
						}
						return quic_tr.Dial(ctx, a, tlsCfg, cfg)
					},
				}

				client := &http.Client{
					Timeout:   5 * time.Second,
					Transport: tr,
				}

				var req, err = http.NewRequest("GET", url.String(), nil)
				if err != nil {
					panic(err)
				}
				req.Close = true

				var resp, err2 = client.Do(req)
				if err2 != nil {
					state[ip.String()]["http_30"] = []any{false, err2.Error()}
				} else {
					state[ip.String()]["http_30"] = []any{true, resp.Status, resp.StatusCode}
				}
				if resp != nil {
					resp.Body.Close()
					fmt.Println("Response status:", resp.Status)
				}

				tr.Close()
				quic_tr.Close()
				udpConn.Close()
			}
		}
		// certificate info
		// something to add later ig
	}

	jsonBytes, err := json.Marshal(state)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonBytes))

	maxWidth := len("IP")
	for ip := range state {
		if l := len(ip); l > maxWidth {
			maxWidth = l
		}
	}

	columns := []table.Column{
		{Title: "IP", Width: maxWidth},
	}

	rows := []table.Row{}
	for ip := range state {
		rows = append(rows, table.Row{ip})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	p := tea.NewProgram(app.Model{Focus: 0, Table: t, IpLen: maxWidth, State: state}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

