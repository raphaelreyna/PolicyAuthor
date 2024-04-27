package conditions

import "github.com/raphaelreyna/policyauthor/pkg/policy"

func AllConditionsMap() map[string]func() policy.ConditionSpec {
	return map[string]func() policy.ConditionSpec{
		"and":      func() policy.ConditionSpec { return &AndSpec{} },
		"or":       func() policy.ConditionSpec { return &OrSpec{} },
		"not":      func() policy.ConditionSpec { return &NotSpec{} },
		"contains": func() policy.ConditionSpec { return &SubstringSpec{} },
		"equal":    func() policy.ConditionSpec { return &EqualSpec{} },
		"cidr":     func() policy.ConditionSpec { return &CIDRSpec{} },
		"regex":    func() policy.ConditionSpec { return &RegexSpec{} },
		"time":     func() policy.ConditionSpec { return &TimeSpec{} },
		"range":    func() policy.ConditionSpec { return &RangeSpec{} },
		"exists":   func() policy.ConditionSpec { return &ExistsSpec{} },
	}
}
