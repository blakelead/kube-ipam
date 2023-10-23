package main

import (
	"log"

	"github.com/blakelead/kube-ipam/internal/client"
	"github.com/blakelead/kube-ipam/internal/config"
	"github.com/blakelead/kube-ipam/internal/informer"
	"github.com/blakelead/kube-ipam/internal/ipam"
)

func main() {
	log.Println("kube-ipam: allocate external IPs to services of type Loadbalancer")
	log.Println("Usage: kube-ipam --cidr=192.168.0.0/24")

	c := config.ParseFlags()

	clientset, err := client.GetClientset(&c.Kubeconfig)
	if err != nil {
		log.Fatalf("Failed to create kube client: %v", err)
	}

	ipPool, err := ipam.NewIPPool(c.CIDR)
	if err != nil {
		log.Fatalf("Failed to create IP pool: %v", err)
	}

	informer.StartInformer(clientset, ipPool)
}
