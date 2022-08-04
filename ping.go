package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/go-ping/ping"
	"github.com/valyala/fasthttp"
)

func endpointPing(template string, region *region) {
	defer wg.Done()

	if region.host == "" {
		code := region.name
		if region.code != "" {
			code = region.code
		}
		region.host = fmt.Sprintf(template, code)
	}

	if strings.HasPrefix(region.host, "http") {
		var err error
		region.rtt, err = httpPing(region.host, args.Count, args.Timeout)
		if err != nil {
			region.err = err
		}
	} else {
		pinger, err := ping.NewPinger(region.host)
		if err != nil {
			region.err = err
			return
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
			region.err = err
			return
		}

		region.rtt = pinger.Statistics().AvgRtt
	}

	if region.rtt == 0 && region.err == nil {
		region.err = fmt.Errorf("timeout")
	}
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
