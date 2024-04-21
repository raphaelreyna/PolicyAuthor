package conditions

import (
	"fmt"
	"net"

	"github.com/raphaelreyna/policyauthor/pkg/maputils"
	"github.com/raphaelreyna/policyauthor/pkg/policy"
	"gopkg.in/yaml.v3"
)

type CIDRSpec struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`

	cidrRange *net.IPNet `yaml:"-"`
}

func (s *CIDRSpec) UnmarshalYAML(value *yaml.Node) error {
	type T CIDRSpec
	var t T
	err := value.Decode(&t)
	if err != nil {
		return err
	}
	*s = CIDRSpec(t)

	_, s.cidrRange, err = net.ParseCIDR(s.Value)
	if err != nil {
		return err
	}

	return err
}

func (s *CIDRSpec) String() string {
	return fmt.Sprintf("[%s] IN CIDR RANGE %+v", s.Key, s.Value)
}

func (s *CIDRSpec) Evaluate(v map[string]any) (bool, error) {
	if val, found := maputils.RecursiveGet(s.Key, v); found {
		if val, ok := val.(string); ok {
			ip := net.ParseIP(val)
			if ip == nil {
				return false, fmt.Errorf("CIDRSpec error: value at key %s is not a valid IP address", s.Key)
			}
			return s.cidrRange.Contains(ip), nil
		}
	}
	return false, policy.NewKeyNotFoundError(s.Key)
}
