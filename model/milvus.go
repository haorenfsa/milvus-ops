package model

type Milvus struct {
	Namespace string `json:"namespace,omitempty"`
	Name      string `json:"name,omitempty"`
	Status    string `json:"status,omitempty"`
	// for crd, use image tag, for helm use app_version
	Version   string `json:"version,omitempty"`
	ManagedBy string `json:"managed_by,omitempty"`
	Type      string `json:"type"` // standalone/cluster
}

type ClassifiedPods struct {
	Components    []string            `json:"components"`
	ComponentPods map[string][]string `json:"component_pods"`
}

func NewClassifiedPods() *ClassifiedPods {
	return &ClassifiedPods{
		Components:    []string{},
		ComponentPods: make(map[string][]string),
	}
}
