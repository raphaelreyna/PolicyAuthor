package policy_test

import (
	"testing"

	"github.com/raphaelreyna/policyauthor/pkg/conditions"
	"github.com/raphaelreyna/policyauthor/pkg/policy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

var config1 = `
policies:
- value: latexmk
  conditions:
  - type: and
    spec:
      conditions:
      - type: equal
        spec:
          key: "remote_addr"
          value: "foo"
      - type: not
        spec:
          condition:
            type: cidr
            spec:
              key: "header.X-Forwarded-For"
              value: "10.1.0.1/24"
`

func TestBasic(t *testing.T) {
	policy.RegisterConditions(conditions.AllConditionsMap())

	p := struct {
		Policies *policy.PolicyEngine `yaml:"policies"`
	}{
		Policies: &policy.PolicyEngine{},
	}

	err := yaml.Unmarshal([]byte(config1), &p)
	require.NoError(t, err)

	m1 := map[string]any{
		"remote_addr": "foo",
		"header": map[string]any{
			"X-Forwarded-For": "12.1.0.1",
		},
	}

	value, ok, err := p.Policies.Evaluate(m1)
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, "latexmk", value)

	m2 := map[string]any{
		"remote_addr": "foo",
		"header": map[string]any{
			"X-Forwarded-For": "10.1.0.1",
		},
	}

	value, ok, err = p.Policies.Evaluate(m2)
	assert.NoError(t, err)
	assert.False(t, ok)
	assert.Nil(t, value)
}

func Test_Fail_NoSpecsRegistered(t *testing.T) {
	p := struct {
		Policies *policy.PolicyEngine `yaml:"policies"`
	}{
		Policies: &policy.PolicyEngine{},
	}

	err := yaml.Unmarshal([]byte(config1), &p)
	require.Error(t, err)
}

var config2 = `
policies:
- value: foo
  conditions:
    - type: regex
      spec:
        key: "host"
        pattern: "(.*)\\.example\\.com"
        return: \1
    - type: equal
      spec:
        key: "remote_addr"
        value: "bar"
`

func TestReturnValues(t *testing.T) {
	policy.RegisterConditions(conditions.AllConditionsMap())

	p := struct {
		Policies *policy.PolicyEngine `yaml:"policies"`
	}{
		Policies: &policy.PolicyEngine{},
	}

	err := yaml.Unmarshal([]byte(config2), &p)
	require.NoError(t, err)

	m1 := map[string]any{
		"host":        "baz.example.com",
		"remote_addr": "qux",
	}

	value, ok, err := p.Policies.Evaluate(m1)
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, "baz", value)

	m2 := map[string]any{
		"host":        "123",
		"remote_addr": "bar",
	}
	value, ok, err = p.Policies.Evaluate(m2)
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, "foo", value)

	m3 := map[string]any{
		"host":        "123",
		"remote_addr": "123",
	}
	value, ok, err = p.Policies.Evaluate(m3)
	assert.NoError(t, err)
	assert.False(t, ok)
	assert.Nil(t, value)
}
