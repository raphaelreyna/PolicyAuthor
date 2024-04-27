package conditions

import (
	"fmt"

	"github.com/raphaelreyna/policyauthor"
	"github.com/raphaelreyna/policyauthor/pkg/maputils"
)

type SubstringSpec struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

func (s *SubstringSpec) String() string {
	return fmt.Sprintf("[%s] SUBSTRING %+v", s.Key, s.Value)
}

func (s *SubstringSpec) Evaluate(v map[string]any) (bool, error) {
	if val, found := maputils.RecursiveGet(s.Key, v); found {
		if val, ok := val.(string); ok {
			return s.Value == val, nil
		}

		return false, fmt.Errorf("ContainSpec error: value at key %s is not a string, got %T", s.Key, val)
	}
	return false, policyauthor.NewKeyNotFoundError(s.Key)
}
