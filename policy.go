package policyauthor

import (
	"fmt"
	"strings"
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
