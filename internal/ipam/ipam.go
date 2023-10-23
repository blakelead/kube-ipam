package ipam

import (
	"fmt"
	"log"
	"net"
	"sync"
)

type IPPool struct {
	pool  map[string]bool
	mutex sync.Mutex
}

func NewIPPool(cidr string) (*IPPool, error) {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}
	pool := make(map[string]bool)
	for ip := ipNet.IP.Mask(ipNet.Mask); ipNet.Contains(ip); inc(ip) {
		pool[ip.String()] = false
	}
	delete(pool, ipNet.IP.String())
	delete(pool, lastIP(ipNet).String())
	return &IPPool{pool: pool}, nil
}

func (p *IPPool) Display() {
	for ip, status := range p.pool {
		log.Printf("IP: %s\tAllocated: %v\n", ip, status)
	}
}

func (p *IPPool) AllocateIP() (net.IP, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	for ip, allocated := range p.pool {
		if !allocated {
			p.pool[ip] = true
			return net.ParseIP(ip), nil
		}
	}
	return nil, fmt.Errorf("no more IPs")
}

func (p *IPPool) DeallocateIP(ip string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, exists := p.pool[ip]; exists {
		p.pool[ip] = false
	}
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func lastIP(n *net.IPNet) net.IP {
	var ip net.IP
	for i := 0; i < len(n.IP); i++ {
		ip = append(ip, n.IP[i]^n.Mask[i])
	}
	return ip
}

// func isIPinCIDR(ipStr, cidrStr string) bool {
// 	_, cidr, err := net.ParseCIDR(cidrStr)
// 	if err != nil {
// 		return false
// 	}
// 	ip := net.ParseIP(ipStr)
// 	if ip == nil {
// 		return false
// 	}
// 	return cidr.Contains(ip)
// }
