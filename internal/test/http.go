package test

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"time"

	quic "github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
)

func _http(url *url.URL, address string, http2 bool) *Status {
	var protocols = &http.Protocols{}
	protocols.SetHTTP1(!http2)
	protocols.SetHTTP2(http2)
	protocols.SetUnencryptedHTTP2(false)

	var client = &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives: true,
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return (&net.Dialer{Timeout: 5 * time.Second}).DialContext(ctx, network, address)
			},
			TLSClientConfig: &tls.Config{
				ServerName: url.Hostname(),
			},
			Protocols: protocols,
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error { // don't follow redirects
			return http.ErrUseLastResponse
		},
	}
	var req, err = http.NewRequest("GET", url.String(), nil)
	if err != nil {
		panic(err)
	}
	req.Close = true

	var resp, err2 = client.Do(req)
	if resp != nil {
		resp.Body.Close()
	}

	if err2 != nil {
		return &Status{
			State: Failed,
			Msg:   err2.Error(),
		}
	} else {
		return &Status{
			State: Success,
			Msg:   resp.Status,
		}
	}
}

func Http_11(url *url.URL, address string) *Status {
	return _http(url, address, false)
}

func Http_2(url *url.URL, address string) *Status {
	return _http(url, address, true)
}

func Http_3(url *url.URL, address string) *Status {
	udpConn, err3 := net.ListenUDP("udp", &net.UDPAddr{})
	if err3 != nil {
		return &Status{
			State: Failed,
			Msg:   err3.Error(),
		}
	} else {
		quic_tr := &quic.Transport{Conn: udpConn}

		tr := &http3.Transport{
			TLSClientConfig: &tls.Config{
				NextProtos: []string{http3.NextProtoH3},
				ServerName: url.Hostname(),
			},
			QUICConfig: &quic.Config{},
			Dial: func(ctx context.Context, addr string, tlsCfg *tls.Config, cfg *quic.Config) (*quic.Conn, error) {
				var a, err = net.ResolveUDPAddr("udp", address)
				if err != nil {
					panic(err)
				}
				return quic_tr.Dial(ctx, a, tlsCfg, cfg)
			},
		}

		client := &http.Client{
			Timeout:   5 * time.Second,
			Transport: tr,
			CheckRedirect: func(req *http.Request, via []*http.Request) error { // don't follow redirects
				return http.ErrUseLastResponse
			},
		}

		var req, err = http.NewRequest("GET", url.String(), nil)
		if err != nil {
			panic(err)
		}
		req.Close = true

		var resp, err2 = client.Do(req)
		if resp != nil {
			resp.Body.Close()
		}
		tr.Close()
		quic_tr.Close()
		udpConn.Close()

		if err2 != nil {
			return &Status{
				State: Failed,
				Msg:   err2.Error(),
			}
		} else {
			return &Status{
				State: Success,
				Msg:   resp.Status,
			}
		}
	}
}
