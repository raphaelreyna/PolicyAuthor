package conditions

import (
	"fmt"

	"github.com/raphaelreyna/policyauthor/pkg/maputils"
	"github.com/raphaelreyna/policyauthor/pkg/policy"
)

type ContainsSpec struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

func (s *ContainsSpec) String() string {
	return fmt.Sprintf("[%s] CONTAINS %+v", s.Key, s.Value)
}

func (s *ContainsSpec) Evaluate(v map[string]any) (bool, error) {
	if val, found := maputils.RecursiveGet(s.Key, v); found {
		if val, ok := val.(string); ok {
			return s.Value == val, nil
		}

		return false, fmt.Errorf("ContainSpec error: value at key %s is not a string, got %T", s.Key, val)
	}
	return false, policy.NewKeyNotFoundError(s.Key)
}
