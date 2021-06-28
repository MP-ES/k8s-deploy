package entities

type Image struct {
	Name string
}

func (i *Image) String() string {
	return i.Name
}

func (i *Image) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var output string
	if err := unmarshal(&output); err != nil {
		return err
	}
	i.Name = output
	return nil
}
