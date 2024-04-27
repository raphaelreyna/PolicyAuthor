package policyauthor

type ConditionSpec interface {
	String() string
	Evaluate(v map[string]any) (bool, error)
}

var conditionsSpecMap = map[string]func() ConditionSpec{}

func RegisterCondition(name string, f func() ConditionSpec) {
	conditionsSpecMap[name] = f
}

func RegisterConditions(m map[string]func() ConditionSpec) {
	for k, v := range m {
		conditionsSpecMap[k] = v
	}
}
