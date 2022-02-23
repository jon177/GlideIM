package main

import (
	"go_im/im/client"
	"go_im/im/dao"
	"go_im/im/group"
	"go_im/pkg/db"
	"go_im/service"
	"go_im/service/api_service"
	"go_im/service/gateway_service"
	"go_im/service/group_messaging"
	"go_im/service/rpc"
)

func main() {
	db.Init()
	dao.Init()

	config, err := service.GetConfig()
	if err != nil {
		panic(err)
	}
	etcd := config.Etcd.Servers

	clientGateway, err := gateway_service.NewClient(&rpc.ClientOptions{
		Name:        config.Gateway.Client.Name,
		EtcdServers: etcd,
	})
	if err != nil {
		panic(err)
	}
	client.Manager = clientGateway

	groupManager, err := group_messaging.NewClient(&rpc.ClientOptions{
		Name:        config.GroupMessaging.Client.Name,
		EtcdServers: etcd,
	})
	if err != nil {
		panic(err)
	}
	group.Manager = groupManager

	server := api_service.NewServer(&rpc.ServerOptions{
		Name:        config.Api.Server.Name,
		Network:     config.Api.Server.Network,
		Addr:        config.Api.Server.Addr,
		Port:        config.Api.Server.Port,
		EtcdServers: etcd,
	})

	err = server.Run()

	if err != nil {
		panic(err)
	}
}