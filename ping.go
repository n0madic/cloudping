package main

import (
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"

	"github.com/mehrdadrad/ping"
	probing "github.com/prometheus-community/pro-bing"
	"github.com/valyala/fasthttp"
)

func endpointPing(template string, region *region) {
	defer wg.Done()

	if region.endpoint == "" {
		code := region.name
		if region.code != "" {
			code = region.code
		}
		region.endpoint = fmt.Sprintf(template, code)
	}

	var err error
	if strings.HasPrefix(region.endpoint, "http") {
		region.rtt, err = httpPing(region.endpoint, args.Count, args.Timeout)
		if err != nil {
			region.err = err
		}
	} else if strings.HasPrefix(region.endpoint, "tcp") {
		address := strings.TrimPrefix(region.endpoint, "tcp://")
		region.rtt, err = tcpPing(address, args.Count, args.Timeout)
		if err != nil {
			region.err = err
		}
	} else {
		icmpPingFn := icmpPing
		if args.AltPing {
			icmpPingFn = icmpAltPing
		}
		region.rtt, err = icmpPingFn(region.endpoint, args.Count, args.Timeout)
		if err != nil {
			region.err = err
		}
	}

	if region.rtt == 0 && region.err == nil {
		region.err = fmt.Errorf("timeout")
	}
}

func icmpPing(endpoint string, count int, timeout time.Duration) (time.Duration, error) {
	pinger, err := probing.NewPinger(endpoint)
	if err != nil {
		return 0, err
	}
	defer pinger.Stop()

	pinger.Count = count
	pinger.Timeout = timeout

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

func icmpAltPing(endpoint string, count int, timeout time.Duration) (time.Duration, error) {
	pinger, err := ping.New(endpoint)
	if err != nil {
		return 0, err
	}

	pinger.SetCount(count)
	pinger.SetTimeout(timeout.String())

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

func httpPing(url string, count int, timeout time.Duration) (time.Duration, error) {
	var rtt time.Duration
	client := fasthttp.Client{
		ReadTimeout: timeout,
	}
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	req.Header.SetMethod(fasthttp.MethodGet)
	defer fasthttp.ReleaseRequest(req)
	for i := 0; i < count; i++ {
		start := time.Now()
		resp := fasthttp.AcquireResponse()
		err := client.Do(req, resp)
		fasthttp.ReleaseResponse(resp)
		if err != nil {
			return 0, err
		}
		rtt += time.Since(start)
	}
	return rtt / time.Duration(count), nil
}

func tcpPing(endpoint string, count int, timeout time.Duration) (time.Duration, error) {
	var rtt time.Duration
	for i := 0; i < count; i++ {
		start := time.Now()
		conn, err := net.DialTimeout("tcp", endpoint, timeout)
		if err != nil {
			return 0, err
		}
		conn.Close()
		rtt += time.Since(start)
	}
	return rtt / time.Duration(count), nil
}
