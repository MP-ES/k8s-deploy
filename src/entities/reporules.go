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

func (r *RepositoryRules) IsImageEnabled(imageName string) bool {
	for _, i := range r.Images {
		if i.Name == imageName {
			return true
		}
	}
	return false
}

func (r *RepositoryRules) IsSecretEnabled(secretName string) bool {
	for _, s := range r.Secrets {
		if s.Name == secretName {
			return true
		}
	}
	return false
}

func (r *RepositoryRules) IsIngressEnabled(ingress string, kEnv K8sEnv) bool {
	if _, ok := (*r.Ingresses)[kEnv]; !ok {
		return false // K8S env not available
	}

	for _, i := range (*r.Ingresses)[kEnv] {
		if i.Name == ingress {
			return true
		}
	}
	return false
}
