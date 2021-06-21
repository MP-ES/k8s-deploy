package entities

type RepositorySchema struct {
	Name            string                  `yaml:"name"`
	K8sEnvs         []*K8sEnv               `yaml:"k8s-envs,flow"`
	Images          []*Image                `yaml:"images,flow"`
	Secrets         []*Secret               `yaml:"secrets,flow"`
	ResourcesQuotas *ResourcesQuotas        `yaml:"resources-quotas"`
	Ingresses       *map[*K8sEnv][]*Ingress `yaml:"ingresses"`
}
