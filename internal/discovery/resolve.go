package discovery

import (
	"context"
	"github.com/grandcat/zeroconf"
	"math"
	"net"
	"os"
	"time"
)

func StartResolver(service string, domain string, waitTime int, updateTime int) {
	for range time.Tick(time.Second * time.Duration(updateTime)) {
		Resolve(service, domain, int(math.Max(float64(waitTime), 1)))
	}
}

func Resolve(service string, domain string, waitTime int) {
	resolver, err := zeroconf.NewResolver(zeroconf.SelectIPTraffic(zeroconf.IPv4))
	if err != nil {
		Log.Error(err, err.Error())
		os.Exit(1)
	}

	entries := make(chan *zeroconf.ServiceEntry)
	go func(results <-chan *zeroconf.ServiceEntry) {
		var res []*TxtData
		for entry := range results {
			if isForeign(entry.AddrIPv4) {
				if txtData, err := Decode(entry.Text[0]); err == nil {
					res = append(res, txtData)
				}
			}
		}
		UpdateForeign(res)
	}(entries)

	var ctx context.Context = nil
	var cancel context.CancelFunc = nil
	if waitTime > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), time.Second*time.Duration(waitTime))
		defer cancel()
	} else {
		ctx = context.Background()
	}

	err = resolver.Browse(ctx, service, domain, entries)
	if err != nil {
		Log.Error(err, err.Error())
		os.Exit(1)
	}

	<-ctx.Done()
}

func getIPs() map[string]bool {
	myIps := map[string]bool{}
	ifaces, err := net.Interfaces()
	if err != nil {
		Log.Error(err, err.Error())
		os.Exit(1)
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			ip := getIP(addr)
			if ip != nil {
				myIps[ip.String()] = true
			}
		}
	}
	return myIps
}

func isForeign(foreignIps []net.IP) bool {
	myIps := getIPs()
	for _, fIp := range foreignIps {
		Log.Info("Received packet from " + fIp.String())
		if myIps[fIp.String()] {
			return false
		}
	}
	return true
}

func getIP(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	return ip
}