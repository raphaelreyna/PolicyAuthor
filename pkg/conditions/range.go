package conditions

import (
	"fmt"

	"github.com/raphaelreyna/policyauthor/pkg/policy"
	"gopkg.in/yaml.v3"
)

type RangeSpec struct {
	Key   string   `yaml:"key"`
	Lower *float64 `yaml:"lower,omitempty"`
	Upper *float64 `yaml:"upper,omitempty"`
}

func (s *RangeSpec) String() string {
	return fmt.Sprintf("[%s] BETWEEN %d AND %d", s.Key, s.Lower, s.Upper)
}

func (s *RangeSpec) UnmarshalYAML(value *yaml.Node) error {
	type S RangeSpec
	var ss S
	if err := value.Decode(&ss); err != nil {
		return err
	}
	*s = RangeSpec(ss)

	if s.Lower == nil && s.Upper == nil {
		return fmt.Errorf("RangeSpec error: both lower and upper bounds are nil")
	}

	return nil
}

func (s *RangeSpec) Evaluate(v map[string]interface{}) (bool, error) {
	val, ok := v[s.Key]
	if !ok {
		return false, policy.NewKeyNotFoundError(s.Key)
	}

	switch {
	case s.Lower != nil && s.Upper != nil:
		switch v := val.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			x := float64(v.(int))
			return x >= *s.Lower && x <= *s.Upper, nil
		case float32, float64:
			x := v.(float64)
			return x >= *s.Lower && x <= *s.Upper, nil
		default:
			return false, fmt.Errorf("key %s is not a number", s.Key)
		}
	case s.Lower != nil:
		switch v := val.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			x := float64(v.(int))
			return x >= *s.Lower, nil
		case float32, float64:
			x := v.(float64)
			return x >= *s.Lower, nil
		default:
			return false, fmt.Errorf("key %s is not a number", s.Key)
		}
	case s.Upper != nil:
		switch v := val.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			x := float64(v.(int))
			return x <= *s.Upper, nil
		case float32, float64:
			x := v.(float64)
			return x <= *s.Upper, nil
		default:
			return false, fmt.Errorf("key %s is not a number", s.Key)
		}
	default:
		return false, fmt.Errorf("RangeSpec error: both lower and upper bounds are nil")
	}

}
