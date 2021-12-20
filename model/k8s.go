package model

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tevino/log"
)

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

type LogOption struct {
	Since      *time.Time
	LimitBytes int64
}

func GetLogOption(ctx *gin.Context) LogOption {
	since := ctx.Query("since")
	size := ctx.Query("size_mb")
	var sizeMB float64
	var err error
	if size != "" {
		sizeMB, err = strconv.ParseFloat(size, 64)
		if err != nil {
			log.Infof("invalid size: %s", size)
		}
	}
	var sinceTime time.Time
	if since != "" {
		sinceTime, err = time.Parse(time.RFC3339, since)
		if err != nil {
			log.Infof("invalid since: %s", since)
		}
	}
	limitBytes := int64(sizeMB * 1024 * 1024)
	if limitBytes < 1 {
		limitBytes = 1
	}
	log.Info(limitBytes, sizeMB*1024*1024)
	return LogOption{
		Since:      &sinceTime,
		LimitBytes: limitBytes,
	}
}
