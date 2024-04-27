package conditions

import (
	"fmt"
	"time"

	"github.com/raphaelreyna/policyauthor/pkg/maputils"
	"github.com/raphaelreyna/policyauthor/pkg/policy"
	"gopkg.in/yaml.v3"
)

type TimeSpec struct {
	Key    string `yaml:"key"`
	Layout string `yaml:"layout"`
	Before string `yaml:"before"`
	After  string `yaml:"after"`

	before time.Time `yaml:"-"`
	after  time.Time `yaml:"-"`
}

func (s *TimeSpec) UnmarshalYAML(value *yaml.Node) error {
	type T TimeSpec
	var t T
	err := value.Decode(&t)
	if err != nil {
		return err
	}
	*s = TimeSpec(t)

	layout := time.RFC3339
	if s.Layout != "" {
		layout = s.Layout
	}
	s.before, err = time.Parse(layout, s.Before)
	if err != nil {
		return fmt.Errorf("TimeSpec error: could not parse 'before' time: %s", err)
	}

	s.after, err = time.Parse(time.RFC3339, s.After)
	if err != nil {
		return fmt.Errorf("TimeSpec error: could not parse 'after' time: %s", err)
	}

	return nil
}

func (s *TimeSpec) String() string {
	switch {
	case s.Before != "" && s.After != "":
		return fmt.Sprintf("[%s] BETWEEN %s AND %s", s.Key, s.Before, s.After)
	case s.Before != "":
		return fmt.Sprintf("[%s] BEFORE %s", s.Key, s.Before)
	case s.After != "":
		return fmt.Sprintf("[%s] AFTER %s", s.Key, s.After)
	default:
		return fmt.Sprintf("[%s] BETWEEN %s AND %s", s.Key, s.Before, s.After)
	}
}

func (s *TimeSpec) Evaluate(v map[string]interface{}) (bool, error) {
	layout := time.RFC3339
	if s.Layout != "" {
		layout = s.Layout
	}

	if val, found := maputils.RecursiveGet(s.Key, v); found {
		if val, ok := val.(string); ok {
			t, err := time.Parse(time.RFC3339, val)
			if err != nil {
				return false, fmt.Errorf("TimeSpec error: value at key %s does not conform to the expected layout (%s): %s", s.Key, layout, err)
			}
			return s.before.Before(t) && s.after.After(t), nil
		}
	}
	return false, policy.NewKeyNotFoundError(s.Key)
}
