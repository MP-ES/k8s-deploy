package entities

type ResourcesQuotas struct {
	LimitsCpu    string `yaml:"limits.cpu"`
	LimitsMemory string `yaml:"limits.memory"`
}
