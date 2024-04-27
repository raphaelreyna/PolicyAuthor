package conditions

import (
	"fmt"
	"reflect"

	"github.com/raphaelreyna/policyauthor"
	"github.com/raphaelreyna/policyauthor/pkg/maputils"
)

type EqualSpec struct {
	Key   string `yaml:"key"`
	Value any    `yaml:"value"`
}

func (s *EqualSpec) String() string {
	return fmt.Sprintf("[%s] EQUALS %+v", s.Key, s.Value)
}

func (s *EqualSpec) Evaluate(v map[string]any) (bool, error) {
	vv, found := maputils.RecursiveGet(s.Key, v)
	if !found {
		return false, policyauthor.NewKeyNotFoundError(s.Key)
	}

	// TODO(raphaelreyna): performance could probably be improved here

	return reflect.DeepEqual(vv, s.Value), nil
}
