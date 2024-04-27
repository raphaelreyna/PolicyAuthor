package policyauthor

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Condition struct {
	Type string        `yaml:"type"`
	Spec ConditionSpec `yaml:"-"`
}

func (c *Condition) UnmarshalYAML(value *yaml.Node) error {
	type C Condition
	type T struct {
		*C   `yaml:",inline"`
		Spec yaml.Node `yaml:"spec"`
	}

	obj := T{C: (*C)(c)}
	if err := value.Decode(&obj); err != nil {
		return err
	}

	specThunk, ok := conditionsSpecMap[c.Type]
	if !ok {
		return fmt.Errorf("unknown condition type: %s", c.Type)
	}
	c.Spec = specThunk()

	if err := obj.Spec.Decode(c.Spec); err != nil {
		return err
	}

	if c.Type == "" {
		return fmt.Errorf("condition type must be set")
	}

	if c.Spec == nil {
		return fmt.Errorf("condition spec must be set")
	}

	return nil
}

func (c *Condition) String() string {
	return c.Spec.String()
}
