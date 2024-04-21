package conditions

import (
	"fmt"
	"strings"

	"github.com/raphaelreyna/policyauthor/pkg/policy"
)

type AndSpec struct {
	Conditions []*policy.Condition `yaml:"conditions"`
}

func (s *AndSpec) String() string {
	b := strings.Builder{}
	a := ""
	for _, c := range s.Conditions {
		b.WriteString(fmt.Sprintf("%s(%s)", a, c))
		a = " AND "
	}
	return b.String()
}

func (s *AndSpec) Evaluate(v map[string]any) (bool, error) {
	for _, c := range s.Conditions {
		hit, err := c.Spec.Evaluate(v)
		if err != nil {
			return false, err
		}
		if !hit {
			return false, nil
		}
	}
	return true, nil
}

type OrSpec struct {
	Conditions []*policy.Condition `yaml:"conditions"`
}

func (s *OrSpec) String() string {
	b := strings.Builder{}
	a := ""
	for _, c := range s.Conditions {
		b.WriteString(fmt.Sprintf("%s(%s)", a, c))
		a = " OR "
	}
	return b.String()
}

func (s *OrSpec) Evaluate(v map[string]any) (bool, error) {
	for _, c := range s.Conditions {
		hit, err := c.Spec.Evaluate(v)
		if err != nil {
			return false, err
		}
		if hit {
			return true, nil
		}
	}
	return false, nil
}

type NotSpec struct {
	Condition policy.Condition `yaml:"condition"`
}

func (s *NotSpec) String() string {
	return fmt.Sprintf("NOT (%s)", s.Condition)
}

func (s *NotSpec) Evaluate(v map[string]any) (bool, error) {
	hit, err := s.Condition.Spec.Evaluate(v)
	if err != nil {
		return false, err
	}
	return !hit, nil
}
