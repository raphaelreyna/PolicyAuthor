package policyauthor

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

type PolicyEngine struct {
	policies []*Policy `yaml:"-"`
}

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
