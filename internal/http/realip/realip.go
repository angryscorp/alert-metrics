package realip

import (
	"net"
	"net/http"
)

type Transport struct {
	transport http.RoundTripper
	realIP    string
}

func New(transport http.RoundTripper) *Transport {
	if transport == nil {
		transport = http.DefaultTransport
	}

	realIP := getLocalIP()

	return &Transport{
		transport: transport,
		realIP:    realIP,
	}
}

func (rt *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if rt.realIP != "" {
		req.Header.Set("X-Real-IP", rt.realIP)
	}
	return rt.transport.RoundTrip(req)
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			return ipNet.IP.String()
		}
	}

	return ""
}
