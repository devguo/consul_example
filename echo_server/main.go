package main

import (
	"github.com/devguo/consul_example/pkg/pb"
	"golang.org/x/net/context"
	"net"
	"fmt"
	"log"
	"google.golang.org/grpc"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type Server struct{
	Ip string
	Port int
	ConsulAddr string
}

func NewServer(ip string, port int,consulAddr string) *Server  {
	return &Server{
		Ip:ip,
		Port:port,
		ConsulAddr:consulAddr,
	}
}

func (s* Server) Echo(ctx context.Context,req *svc.EchoRequest) (*svc.EchoResponse, error) {
	rsp := &svc.EchoResponse{
		Msg:req.Msg,
	}
	return rsp,nil
}


//used by consul for health checking
type HealthServer struct{}

func (srv *HealthServer) Check(context.Context, *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

//register service to consul, hard code register info
func (s *Server) RegisterToConsul() error  {
	cfg := api.DefaultConfig()
	cfg.Address = s.ConsulAddr

	client, err := api.NewClient(cfg)
	if err != nil {
		return err
	}

	service := &api.AgentServiceRegistration{
		ID: "echo-001",
		Name: "echo",
		Address:s.Ip,
		Port:s.Port,
		Check:&api.AgentServiceCheck{
			Interval:"1s",
			GRPC:fmt.Sprintf("%v:%v/%v", s.Ip, s.Port, "echo"),
			DeregisterCriticalServiceAfter:"30s",
		},
	}

	if err := client.Agent().ServiceRegister(service) ; err != nil {
		return err
	}

	return nil
}

func main()  {
	//replace ip and port with yours
	srv := NewServer("192.168.2.209",13579,"192.168.2.251:8500")
	err := srv.RegisterToConsul()
	if err != nil {
		log.Printf("register to consul failed. %v",err)
		return
	}

	addr := fmt.Sprintf("%s:%d",srv.Ip,srv.Port)
	lis , err := net.Listen("tcp",addr)
	if err != nil {
		log.Printf("listen failed. addr=%s, err=%v",addr,err)
		return
	}

	gServer := grpc.NewServer()
	svc.RegisterEchoServer(gServer,srv)
	grpc_health_v1.RegisterHealthServer(gServer,&HealthServer{})
	err = gServer.Serve(lis)
	if err != nil {
		log.Printf("run server failed. err=%v",err)
	}

}