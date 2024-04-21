package policy

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	ErrKeyNotFound = fmt.Errorf("key not found")
)

func NewKeyNotFoundError(key string) error {
	return fmt.Errorf("%w: %s", ErrKeyNotFound, key)
}

type Policy struct {
	Value      any          `yaml:"value"`
	Conditions []*Condition `yaml:"conditions"`
}

type valueReturner interface {
	ValueReturnEnabled() bool
	EvaluateWithReturnValue(v map[string]any) (any, bool, error)
}

func (p *Policy) Evaluate(v map[string]any) (value any, hit bool, err error) {
	for _, c := range p.Conditions {
		if vr, ok := c.Spec.(valueReturner); ok {
			if !vr.ValueReturnEnabled() {
				if hit, err = c.Spec.Evaluate(v); err != nil {
					return
				}
				if hit {
					return p.Value, true, nil
				}
			} else {
				value, hit, err = vr.EvaluateWithReturnValue(v)
				if err != nil {
					return
				}
				if hit {
					return value, true, nil
				}
			}
		} else {
			if hit, err = c.Spec.Evaluate(v); err != nil {
				return
			}
			if hit {
				return p.Value, true, nil
			}
		}
	}

	return nil, false, nil
}

func (p *Policy) String() string {
	b := strings.Builder{}
	for i, c := range p.Conditions {
		b.WriteString(fmt.Sprintf("%s(%s)", strings.Repeat(" ", i), c))
	}

	return b.String()
}

type PolicyEngine struct {
	policies []*Policy `yaml:"-"`
}

type SpecMap map[string]func() ConditionSpec

func (pe *PolicyEngine) UnmarshalYAML(value *yaml.Node) error {
	if len(conditionsSpecMap) == 0 {
		return fmt.Errorf("no specs registered")
	}

	x := []yaml.Node{}

	if err := value.Decode(&x); err != nil {
		return err
	}

	pe.policies = make([]*Policy, len(x))
	for i, p := range x {
		policy := Policy{}
		if err := p.Decode(&policy); err != nil {
			return err
		}
		pe.policies[i] = &policy
	}

	return nil
}

func (pe *PolicyEngine) Evaluate(v map[string]any) (value any, hit bool, err error) {
	for _, p := range pe.policies {
		if value, hit, err = p.Evaluate(v); err != nil {
			return
		}
		if hit {
			return
		}
	}

	return nil, false, nil
}

func (pe *PolicyEngine) String() string {
	b := strings.Builder{}
	a := ""
	for _, p := range pe.policies {
		b.WriteString(fmt.Sprintf("%s(%s)", a, p))
		a = " OR "
	}
	return b.String()
}
