package main

import (
	"crypto/tls"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"

	"github.com/mehrdadrad/ping"
	probing "github.com/prometheus-community/pro-bing"
	"github.com/valyala/fasthttp"
)

func endpointPing(endpoint string) (time.Duration, error) {
	if strings.HasPrefix(endpoint, "http") {
		return httpPing(endpoint)
	} else if strings.HasPrefix(endpoint, "tcp") {
		address := strings.TrimPrefix(endpoint, "tcp://")
		return tcpPing(address)
	}
	if args.AltPing {
		return icmpAltPing(endpoint)
	}
	return icmpPing(endpoint)
}

func icmpPing(endpoint string) (time.Duration, error) {
	pinger, err := probing.NewPinger(endpoint)
	if err != nil {
		return 0, err
	}
	defer pinger.Stop()

	pinger.Count = args.Count
	pinger.Timeout = args.Timeout

	resolve_count := 0
	for {
		err = pinger.Resolve()
		if err == nil || resolve_count > 3 {
			break
		}
		resolve_count++
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(200)))
	}

	err = pinger.Run()
	if err != nil {
		return 0, err
	}

	return pinger.Statistics().AvgRtt, nil
}

func icmpAltPing(endpoint string) (time.Duration, error) {
	pinger, err := ping.New(endpoint)
	if err != nil {
		return 0, err
	}

	pinger.SetCount(args.Count)
	pinger.SetTimeout(args.Timeout.String())

	r, err := pinger.Run()
	if err != nil {
		return 0, err
	}

	var rtt float64
	var count_success int
	for pr := range r {
		if pr.Err == nil && pr.RTT > 0 {
			rtt += pr.RTT
			count_success++
		}
	}
	return time.Duration((rtt / float64(count_success)) * float64(time.Millisecond)), nil
}

func httpPing(url string) (time.Duration, error) {
	var rtt time.Duration
	client := fasthttp.Client{
		ReadTimeout: args.Timeout,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	req := fasthttp.AcquireRequest()
	req.Header.SetMethod(fasthttp.MethodGet)
	defer fasthttp.ReleaseRequest(req)
	for i := 0; i < args.Count; i++ {
		stampedURL := fmt.Sprintf("%s?%d", url, time.Now().UnixNano()/int64(time.Millisecond))
		req.SetRequestURI(stampedURL)
		start := time.Now()
		resp := fasthttp.AcquireResponse()
		err := client.Do(req, resp)
		fasthttp.ReleaseResponse(resp)
		if err != nil {
			return 0, err
		}
		rtt += time.Since(start)
	}
	return rtt / time.Duration(args.Count), nil
}

func tcpPing(endpoint string) (time.Duration, error) {
	var rtt time.Duration
	for i := 0; i < args.Count; i++ {
		start := time.Now()
		conn, err := net.DialTimeout("tcp", endpoint, args.Timeout)
		if err != nil {
			return 0, err
		}
		conn.Close()
		rtt += time.Since(start)
	}
	return rtt / time.Duration(args.Count), nil
}
