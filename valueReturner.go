package policyauthor

type ValueReturner interface {
	ValueReturnEnabled() bool
	EvaluateWithReturnValue(v map[string]any) (any, bool, error)
}

type ValueReturnerNil struct{}
