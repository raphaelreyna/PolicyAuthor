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

type ValueReturner interface {
	ValueReturnEnabled() bool
	EvaluateWithReturnValue(v map[string]any) (any, bool, error)
}

type ValueReturnerNil struct{}

func (p *Policy) Evaluate(evaluationContext map[string]any) (value any, hit bool, err error) {
	if len(evaluationContext) == 0 {
		return nil, false, fmt.Errorf("evaluation context is empty")
	}

	for _, c := range p.Conditions {
		if vr, ok := c.Spec.(ValueReturner); ok {
			if !vr.ValueReturnEnabled() {
				if hit, err = c.Spec.Evaluate(evaluationContext); err != nil {
					return
				}
				if hit {
					return p.Value, true, nil
				}
			} else {
				value, hit, err = vr.EvaluateWithReturnValue(evaluationContext)
				if err != nil {
					return
				}
				if hit {
					if _, ok := value.(ValueReturnerNil); ok {
						return p.Value, true, nil
					}
					return value, true, nil
				}
			}
		} else {
			if hit, err = c.Spec.Evaluate(evaluationContext); err != nil {
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

		if len(policy.Conditions) == 0 {
			return fmt.Errorf("no conditions found in policy %d", i)
		}

		pe.policies[i] = &policy
	}

	if len(pe.policies) == 0 {
		return fmt.Errorf("no policies found")
	}

	return nil
}

func (pe *PolicyEngine) Evaluate(evaluationContext map[string]any) (value any, hit bool, err error) {
	if len(evaluationContext) == 0 {
		return nil, false, fmt.Errorf("evaluation context is empty")
	}

	for _, p := range pe.policies {
		if value, hit, err = p.Evaluate(evaluationContext); err != nil {
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
