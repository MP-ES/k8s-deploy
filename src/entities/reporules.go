package entities

type RepositoryRules struct {
	Name            string                 `yaml:"name"`
	K8sEnvs         []*K8sEnv              `yaml:"k8s-envs,flow"`
	Images          []*Image               `yaml:"images,flow"`
	Secrets         []*Secret              `yaml:"secrets,flow"`
	ResourcesQuotas *ResourcesQuotas       `yaml:"resources-quotas"`
	Ingresses       *map[K8sEnv][]*Ingress `yaml:"ingresses"`
}

func (r *RepositoryRules) IsK8sEnvEnabled(kEnv *K8sEnv) bool {
	for _, k := range r.K8sEnvs {
		if *k == *kEnv {
			return true
		}
	}
	return false
}
