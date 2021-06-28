package entities

type Ingress struct {
	Name string
}

func (i *Ingress) String() string {
	return i.Name
}

func (i *Ingress) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var output string
	if err := unmarshal(&output); err != nil {
		return err
	}
	i.Name = output
	return nil
}
