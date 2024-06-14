package policyauthor

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
	ValueFrom  string       `yaml:"valueFrom"`
	Conditions []*Condition `yaml:"conditions"`
}

func (p *Policy) UnmarshalYAML(value *yaml.Node) error {
	type T Policy
	var t T
	err := value.Decode(&t)
	if err != nil {
		return err
	}
	*p = Policy(t)

	if p.ValueFrom != "" && p.Value != nil {
		return fmt.Errorf("cannot have both value and valueFrom")
	}

	return nil
}

func (p *Policy) Evaluate(evaluationContext map[string]any) (value any, hit bool, err error) {
	if len(evaluationContext) == 0 {
		return nil, false, fmt.Errorf("evaluation context is empty")
	}

	val := p.Value
	if p.ValueFrom != "" {
		val = evaluationContext[p.ValueFrom]
	}

	for _, c := range p.Conditions {
		if vr, ok := c.Spec.(ValueReturner); ok {
			if !vr.ValueReturnEnabled() {
				if hit, err = c.Spec.Evaluate(evaluationContext); err != nil {
					return
				}
				if hit {
					return val, true, nil
				}
			} else {
				value, hit, err = vr.EvaluateWithReturnValue(evaluationContext)
				if err != nil {
					return
				}
				if hit {
					if _, ok := value.(ValueReturnerNil); ok {
						return val, true, nil
					}
					return value, true, nil
				}
			}
		} else {
			if hit, err = c.Spec.Evaluate(evaluationContext); err != nil {
				return
			}
			if hit {
				return val, true, nil
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
