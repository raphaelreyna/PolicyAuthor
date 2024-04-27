package conditions

import "github.com/raphaelreyna/policyauthor"

func AllConditionsMap() map[string]func() policyauthor.ConditionSpec {
	return map[string]func() policyauthor.ConditionSpec{
		"and":      func() policyauthor.ConditionSpec { return &AndSpec{} },
		"or":       func() policyauthor.ConditionSpec { return &OrSpec{} },
		"not":      func() policyauthor.ConditionSpec { return &NotSpec{} },
		"contains": func() policyauthor.ConditionSpec { return &SubstringSpec{} },
		"equal":    func() policyauthor.ConditionSpec { return &EqualSpec{} },
		"cidr":     func() policyauthor.ConditionSpec { return &CIDRSpec{} },
		"regex":    func() policyauthor.ConditionSpec { return &RegexSpec{} },
		"time":     func() policyauthor.ConditionSpec { return &TimeSpec{} },
		"range":    func() policyauthor.ConditionSpec { return &RangeSpec{} },
		"exists":   func() policyauthor.ConditionSpec { return &ExistsSpec{} },
	}
}
