package test

import (
	"crypto/tls"
	"net"
	"net/url"
	"time"
)

func Tls(url *url.URL, address string, version uint16) *Status {
	conn, tls_err := tls.DialWithDialer(&net.Dialer{
		Timeout: 5 * time.Second,
	}, "tcp", address, &tls.Config{
		ServerName: url.Hostname(),
		MinVersion: version,
		MaxVersion: version,
	})
	if tls_err != nil {
		// fmt.Fprintln(os.Stderr, "tls dial err", tls_err)

		return &Status{
			State: Failed,
			Msg:   tls_err.Error(),
		}
	} else {
		conn.Close()

		return &Status{
			State: Success,
		}
	}
}
