package main

import (
	"github.com/hashicorp/consul/api"
	"log"
	"strings"
	"time"
)

func PrintServices(cli *api.Client) error {
	services, _, err := cli.Health().Service("echo", "", true, nil)
	if err != nil {
		log.Printf("get service failed. err=%v", err)
		return err
	}

	var list []string
	for _, s := range services {
		list = append(list, s.Service.ID)
	}

	log.Printf("get service list %s", strings.Join(list, ","))
	return nil
}

func main() {
	cfg := api.DefaultConfig()
	cfg.Address = "192.168.2.209:8500"
	cli, err := api.NewClient(cfg)
	if err != nil {
		log.Printf("init consul client failed. %v", err)
		return
	}

	ticker := time.NewTicker(5 * time.Second)
	for t := range ticker.C {
		err = PrintServices(cli)
		if err != nil {
			log.Printf("%s client run with error %v", t, err)
		}
	}
}
