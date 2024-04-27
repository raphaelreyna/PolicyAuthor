package conditions

import (
	"fmt"

	"github.com/raphaelreyna/policyauthor/pkg/maputils"
)

type ExistsSpec struct {
	Key string `yaml:"key"`
}

func (s *ExistsSpec) String() string {
	return fmt.Sprintf("[%s] EXISTS", s.Key)
}

func (s *ExistsSpec) Evaluate(v map[string]interface{}) (bool, error) {
	_, found := maputils.RecursiveGet(s.Key, v)
	return found, nil
}
