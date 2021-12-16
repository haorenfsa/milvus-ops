package model

import "github.com/gin-gonic/gin"

type ContainerLocation struct {
	Cluster   string
	Namespace string
	Pod       string
	Container string
}

type MilvusLocateOption struct {
	Cluster   string
	Namespace string
	Milvus    string
	Component string
	Pod       string
	Container string
	ManagedBy string // helm / operator
}

func GetMilvusLocateOption(ctx *gin.Context) MilvusLocateOption {
	cluster := ctx.Param("cluster")
	if cluster == "_default" {
		cluster = ""
	}
	namespace := ctx.Param("namespace")
	milvus := ctx.Param("milvus")
	component := ctx.Query("component")
	pod := ctx.Query("pod")
	container := ctx.Query("container")
	managedBy := ctx.Query("by")
	return MilvusLocateOption{
		Cluster:   cluster,
		Namespace: namespace,
		Milvus:    milvus,
		Component: component,
		Pod:       pod,
		Container: container,
		ManagedBy: managedBy,
	}
}
