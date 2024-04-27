package policyauthor_test

import (
	"testing"

	"github.com/raphaelreyna/policyauthor"
	"github.com/raphaelreyna/policyauthor/pkg/conditions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func Test_Fail_NoSpecsRegistered(t *testing.T) {
	conf := `
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

	p := struct {
		Policies *policyauthor.PolicyEngine `yaml:"policies"`
	}{
		Policies: &policyauthor.PolicyEngine{},
	}

	err := yaml.Unmarshal([]byte(conf), &p)
	require.Error(t, err)
}

type ContextTest struct {
	Map      map[string]any
	TestFunc func(t *testing.T, idx int, value any, hit bool, err error)
}

type TestConfig struct {
	Config       string
	ContextTests []ContextTest
}

func TestMain(t *testing.T) {
	policyauthor.RegisterConditions(conditions.AllConditionsMap())

	var tests = map[string]TestConfig{
		"basic": {
			Config: `
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
`,
			ContextTests: []ContextTest{
				{
					Map: map[string]any{
						"remote_addr": "foo",
						"header": map[string]any{
							"X-Forwarded-For": "12.1.0.1",
						},
					},
					TestFunc: func(t *testing.T, idx int, value any, hit bool, err error) {
						assert.NoError(t, err)
						assert.True(t, hit)
						assert.Equal(t, "latexmk", value)
					},
				},
				{
					Map: map[string]any{
						"remote_addr": "foo",
						"header": map[string]any{
							"X-Forwarded-For": "10.1.0.1",
						},
					},
					TestFunc: func(t *testing.T, idx int, value any, hit bool, err error) {
						assert.NoError(t, err)
						assert.False(t, hit)
						assert.Nil(t, value)
					},
				},
			},
		},
		"basic_return-values": {
			Config: `
policies:
- value: foo
  conditions:
    - type: or
      spec:
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
`,
			ContextTests: []ContextTest{
				{
					Map: map[string]any{
						"host":        "baz.example.com",
						"remote_addr": "qux",
					},
					TestFunc: func(t *testing.T, idx int, value any, hit bool, err error) {
						assert.NoError(t, err)
						assert.True(t, hit)
						assert.Equal(t, "baz", value)
					},
				},
				{
					Map: map[string]any{
						"host":        "123",
						"remote_addr": "bar",
					},
					TestFunc: func(t *testing.T, idx int, value any, hit bool, err error) {
						assert.NoError(t, err)
						assert.True(t, hit)
						assert.Equal(t, "foo", value)
					},
				},
				{
					Map: map[string]any{
						"host":        "123",
						"remote_addr": "123",
					},
					TestFunc: func(t *testing.T, idx int, value any, hit bool, err error) {
						assert.NoError(t, err)
						assert.False(t, hit)
						assert.Nil(t, value)
					},
				},
			},
		},
		"multiple_conditions": {
			Config: `
policies:
  - value: foo
    conditions:
      - type: equal
        spec:
          key: "remote_addr"
          value: "1"
  - value: bar
    conditions:
      - type: equal
        spec:
          key: "remote_addr"
          value: "2"
`,
			ContextTests: []ContextTest{
				{
					Map: map[string]any{
						"remote_addr": "1",
					},
					TestFunc: func(t *testing.T, idx int, value any, hit bool, err error) {
						assert.NoError(t, err)
						assert.True(t, hit)
						assert.Equal(t, "foo", value)
					},
				},
				{
					Map: map[string]any{
						"remote_addr": "2",
					},
					TestFunc: func(t *testing.T, idx int, value any, hit bool, err error) {
						assert.NoError(t, err)
						assert.True(t, hit)
						assert.Equal(t, "bar", value)
					},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			p := struct {
				Policies *policyauthor.PolicyEngine `yaml:"policies"`
			}{
				Policies: &policyauthor.PolicyEngine{},
			}

			err := yaml.Unmarshal([]byte(test.Config), &p)
			require.NoError(t, err)

			for ctxIdx, ctxTest := range test.ContextTests {
				value, hit, err := p.Policies.Evaluate(ctxTest.Map)
				ctxTest.TestFunc(t, ctxIdx, value, hit, err)
			}
		})
	}
}

func TestEdgeCases(t *testing.T) {
	policyauthor.RegisterConditions(conditions.AllConditionsMap())

	var edgeCaseTests = map[string]TestConfig{
		"empty_map": {
			Config: `
policies:
- value: empty_test
  conditions:
  - type: equal
    spec:
      key: "nonexistent"
      value: "none"
`,
			ContextTests: []ContextTest{
				{
					Map: map[string]any{},
					TestFunc: func(t *testing.T, idx int, value any, hit bool, err error) {
						assert.Error(t, err)
						assert.False(t, hit)
						assert.Nil(t, value)
					},
				},
			},
		},
		"malformed_config": {
			Config: `
policies:
 - value:`,
			ContextTests: []ContextTest{
				{
					Map:      map[string]any{"remote_addr": "foo"},
					TestFunc: func(t *testing.T, idx int, value any, hit bool, err error) {},
				},
			},
		},
	}

	for name, test := range edgeCaseTests {
		t.Run(name, func(t *testing.T) {
			p := struct {
				Policies *policyauthor.PolicyEngine `yaml:"policies"`
			}{
				Policies: &policyauthor.PolicyEngine{},
			}

			err := yaml.Unmarshal([]byte(test.Config), &p)
			if name == "malformed_config" {
				require.Error(t, err)
				return
			} else {
				require.NoError(t, err)
			}

			for ctxIdx, ctxTest := range test.ContextTests {
				value, hit, err := p.Policies.Evaluate(ctxTest.Map)
				ctxTest.TestFunc(t, ctxIdx, value, hit, err)
			}
		})
	}
}
