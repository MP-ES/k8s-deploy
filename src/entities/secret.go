package entities

type Secret struct {
	Name string
}

func (s *Secret) String() string {
	return s.Name
}

func (s *Secret) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var output string
	if err := unmarshal(&output); err != nil {
		return err
	}
	s.Name = output
	return nil
}
