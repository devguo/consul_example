package main

import (
	"google.golang.org/grpc"
	"fmt"
	"github.com/devguo/consul_example/pkg/pb"
	"context"
	"time"
	"log"
	"github.com/hashicorp/consul/api"
	)

type Client struct {
	Ip string
	Port int
	conn *grpc.ClientConn
}

func NewClient(ip string, port int) *Client  {
	return &Client {
		Ip:ip,
		Port:port,
	}
}

func (c *Client) Init() error  {
	addr := fmt.Sprintf("%s:%d",c.Ip, c.Port)
	conn, err := grpc.Dial(addr,grpc.WithInsecure())
	if err != nil{
		return err
	}

	c.conn = conn
	return  nil
}

func (c *Client) Run() error  {
	req := &svc.EchoRequest{
		Msg:"Hello Hello Hello",
	}

	echoCli := svc.NewEchoClient(c.conn)

	ctx , cancel := context.WithTimeout(context.Background(),time.Second)
	defer cancel()

	rsp , err := echoCli.Echo(ctx, req)
	if err != nil {
		return err
	}

	log.Printf("get dummy respose. %s", rsp.GetMsg())
	return nil
}

func main()  {
	//init consul client
	cfg := api.DefaultConfig()
	cfg.Address = "192.168.2.220:8500"
	cli ,err := api.NewClient(cfg)
	if err != nil {
		log.Printf("init consul client failed. %v",err)
		return
	}

	//services,err := cli.Agent().Services()
	services,_,err := cli.Health().Service("echo","",true,nil)
	if err != nil {
		log.Printf("get service failed. err=%v",err)
		return
	}


	//find service by id, can find by other info, eg, Tag, Service
	echo := services[0].Service


	c := NewClient(echo.Address,echo.Port)
	err = c.Init()
	if err != nil {
		log.Printf("establish connection failed. %v",err)
		return
	}

	ticker := time.NewTicker(5 * time.Second)
	for t := range ticker.C {
		err = c.Run()
		if err != nil {
			log.Printf("%s client run with error %v",t,err)
		}
	}
}
