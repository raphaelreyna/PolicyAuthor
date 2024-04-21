package conditions

import "fmt"

type ExistsSpec struct {
	Key string `yaml:"key"`
}

func (s *ExistsSpec) String() string {
	return fmt.Sprintf("[%s] EXISTS", s.Key)
}

func (s *ExistsSpec) Evaluate(v map[string]interface{}) (bool, error) {
	_, ok := v[s.Key]
	return ok, nil
}
