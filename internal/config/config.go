package config

import (
	"flag"
	"os"
	"path/filepath"
)

type Config struct {
	Kubeconfig string
	CIDR       string
}

func ParseFlags() Config {
	var config Config
	flag.StringVar(&config.Kubeconfig, "kubeconfig", filepath.Join(os.Getenv("HOME"), ".kube", "config"), "absolute path to the kubeconfig file")
	flag.StringVar(&config.CIDR, "cidr", "192.168.0.0/24", "CIDR block for LoadBancer IPs allocation")
	flag.Parse()
	return config
}
