package service

import . "github.com/haorenfsa/milvus-ops/model"

type Health struct {
}

func NewHealth() *Health {
	return &Health{}
}

var defaultHealthCheck = HealthCheck{
	Healthy: true,
}

func (h Health) Get() *HealthCheck {
	return &defaultHealthCheck
}
