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

func (s *AndSpec) ValueReturnEnabled() bool {
	for _, c := range s.Conditions {
		if vr, ok := c.Spec.(policy.ValueReturner); ok {
			if vr.ValueReturnEnabled() {
				return true
			}
		}
	}
	return false
}

func (s *AndSpec) EvaluateWithReturnValue(v map[string]any) (any, bool, error) {
	var (
		val      any
		foundVal bool
	)

	for _, c := range s.Conditions {
		if !foundVal {
			if vr, ok := c.Spec.(policy.ValueReturner); ok {
				if vr.ValueReturnEnabled() {
					var err error
					v, hit, err := vr.EvaluateWithReturnValue(v)
					if err != nil {
						return nil, false, err
					}
					if !hit {
						return nil, false, nil
					}
					val = v
					foundVal = true
				} else {
					hit, err := c.Spec.Evaluate(v)
					if err != nil {
						return nil, false, err
					}
					if !hit {
						return policy.ValueReturnerNil{}, false, nil
					}
				}
			} else {
				hit, err := c.Spec.Evaluate(v)
				if err != nil {
					return nil, false, err
				}
				if !hit {
					return nil, false, nil
				}
			}
		}
	}
	return val, true, nil
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

func (s *OrSpec) ValueReturnEnabled() bool {
	for _, c := range s.Conditions {
		if vr, ok := c.Spec.(policy.ValueReturner); ok {
			if vr.ValueReturnEnabled() {
				return true
			}
		}
	}
	return false
}

func (s *OrSpec) EvaluateWithReturnValue(v map[string]any) (any, bool, error) {
	for _, c := range s.Conditions {
		if vr, ok := c.Spec.(policy.ValueReturner); ok {
			if vr.ValueReturnEnabled() {
				var err error
				v, hit, err := vr.EvaluateWithReturnValue(v)
				if err != nil {
					return nil, false, err
				}
				if hit {
					return v, true, nil
				}
			} else {
				hit, err := c.Spec.Evaluate(v)
				if err != nil {
					return nil, false, err
				}
				if hit {
					return policy.ValueReturnerNil{}, true, nil
				}
			}
		} else {
			hit, err := c.Spec.Evaluate(v)
			if err != nil {
				return nil, false, err
			}
			if hit {
				return policy.ValueReturnerNil{}, true, nil
			}
		}
	}
	return nil, false, nil
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

func (s *NotSpec) ValueReturnEnabled() bool {
	if vr, ok := s.Condition.Spec.(policy.ValueReturner); ok {
		return vr.ValueReturnEnabled()
	}
	return false
}

func (s *NotSpec) EvaluateWithReturnValue(v map[string]any) (any, bool, error) {
	if vr, ok := s.Condition.Spec.(policy.ValueReturner); ok {
		if vr.ValueReturnEnabled() {
			v, hit, err := vr.EvaluateWithReturnValue(v)
			if err != nil {
				return nil, false, err
			}
			return v, !hit, nil
		} else {
			hit, err := s.Condition.Spec.Evaluate(v)
			if err != nil {
				return nil, false, err
			}
			return policy.ValueReturnerNil{}, !hit, nil
		}
	}
	hit, err := s.Condition.Spec.Evaluate(v)
	if err != nil {
		return nil, false, err
	}
	return nil, !hit, nil
}
