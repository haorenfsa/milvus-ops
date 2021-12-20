package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/haorenfsa/milvus-ops/ctrl"
	"github.com/haorenfsa/milvus-ops/helm"
	"github.com/haorenfsa/milvus-ops/k8s"
	"github.com/haorenfsa/milvus-ops/service"
	"github.com/haorenfsa/milvus-ops/storage/file"

	"github.com/haorenfsa/milvus-ops/server"
)

func main() {
	var staticPath string
	var kubeconfigPath string
	var port int
	flag.StringVar(&staticPath, "s", "./web/build", "static file path to serve, not serve when empty")
	flag.IntVar(&port, "p", 8080, "server port")
	flag.StringVar(&kubeconfigPath, "k", "", "kubeconfig path")
	flag.Parse()
	log.Print(staticPath)

	theServer := server.NewHTTPServer()

	ctx := context.Background()

	kubeStorage := file.NewStorage(kubeconfigPath)
	k8sCliGetter := k8s.NewK8sClientGetter(kubeStorage)
	clis, err := k8sCliGetter.ListClients(ctx)
	if err != nil {
		log.Fatalf("list k8s client failed: %v", err)
	}
	helmCli := helm.NewClients(clis)
	healthService := service.NewHealth()
	healthCtrl := ctrl.NewHealthController(healthService)

	milvusService := service.NewMilvusService(k8sCliGetter, helmCli)
	milvusCtrl := ctrl.NewMilvusController(milvusService)

	theServer.UseControllers([]server.Controller{
		healthCtrl,
		milvusCtrl,
	})
	theServer.ServeStaticPath(staticPath)
	theServer.Run(fmt.Sprintf(":%d", port))
}
