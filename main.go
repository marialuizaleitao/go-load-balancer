package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type Server interface {
	Address() string
	IsAlive() bool
	Serve(rw http.ResponseWriter, r *http.Request)
}

type SimpleServer struct {
	address string
	proxy   *httputil.ReverseProxy
}

func (s *SimpleServer) Address() string {
	return s.address
}

func (s *SimpleServer) IsAlive() bool {
	return true
}

func (s *SimpleServer) Serve(rw http.ResponseWriter, r *http.Request) {
	s.proxy.ServeHTTP(rw, r)
}

type LoadBalancer struct {
	port            string
	roundRobinCount int
	servers         []Server
	mu              sync.Mutex
}

func NewLoadBalancer(port string, servers []Server) *LoadBalancer {
	return &LoadBalancer{
		port:            port,
		roundRobinCount: 0,
		servers:         servers,
	}
}

func NewSimpleServer(address string) *SimpleServer {
	serverURL, err := url.Parse(address)
	if err != nil {
		panic(fmt.Sprintf("Error parsing server URL %s: %v", address, err))
	}

	return &SimpleServer{
		address: address,
		proxy:   httputil.NewSingleHostReverseProxy(serverURL),
	}
}

func (lb *LoadBalancer) getNextAvailableServer() Server {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	for i := 0; i < len(lb.servers); i++ {
		server := lb.servers[lb.roundRobinCount%len(lb.servers)]
		lb.roundRobinCount++
		if server.IsAlive() {
			return server
		}
	}

	if len(lb.servers) > 0 {
		return lb.servers[0]
	}

	return nil
}

func (lb *LoadBalancer) serveProxy(rw http.ResponseWriter, r *http.Request) {
	targetServer := lb.getNextAvailableServer()
	if targetServer == nil {
		http.Error(rw, "No available servers", http.StatusServiceUnavailable)
		return
	}

	fmt.Printf("Forwarding request to address %s\n", targetServer.Address())
	targetServer.Serve(rw, r)
}

func main() {
	servers := []Server{
		NewSimpleServer("https://www.facebook.com"),
		NewSimpleServer("http://www.bing.com"),
		NewSimpleServer("https://www.google.com"),
	}
	lb := NewLoadBalancer("8000", servers)
	handleRedirect := func(rw http.ResponseWriter, r *http.Request) {
		lb.serveProxy(rw, r)
	}
	http.HandleFunc("/", handleRedirect)

	fmt.Printf("Serving requests at 'localhost:%v'\n", lb.port)
	http.ListenAndServe(":"+lb.port, nil)
}
