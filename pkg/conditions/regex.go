package conditions

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/raphaelreyna/policyauthor"
	"github.com/raphaelreyna/policyauthor/pkg/maputils"
	"gopkg.in/yaml.v3"
)

type RegexSpec struct {
	Key     string `yaml:"key"`
	Pattern string `yaml:"pattern"`
	Return  string `yaml:"return"`

	r *regexp.Regexp `yaml:"-"`
}

func (s *RegexSpec) String() string {
	return fmt.Sprintf("[%s] MATCHES REGEX %+v", s.Key, s.Pattern)
}

func (s *RegexSpec) UnmarshalYAML(value *yaml.Node) error {
	type S RegexSpec

	obj := S{}
	if err := value.Decode(&obj); err != nil {
		return err
	}
	*s = RegexSpec(obj)

	r, err := regexp.Compile(s.Pattern)
	if err != nil {
		return err
	}
	s.r = r

	return nil
}

func (s *RegexSpec) Evaluate(v map[string]interface{}) (bool, error) {
	if val, found := maputils.RecursiveGet(s.Key, v); found {
		if val, ok := val.(string); ok {
			return s.r.MatchString(val), nil
		}

		return false, fmt.Errorf("key %s is not a string", s.Key)
	}
	return false, policyauthor.NewKeyNotFoundError(s.Key)
}

func (s *RegexSpec) ValueReturnEnabled() bool {
	return s.Return != ""
}

func (s *RegexSpec) EvaluateWithReturnValue(v map[string]interface{}) (interface{}, bool, error) {
	if val, found := maputils.RecursiveGet(s.Key, v); found {
		val, ok := val.(string)
		if !ok {
			return nil, false, fmt.Errorf("key %s is not a string", s.Key)
		}

		if !s.r.MatchString(val) {
			return nil, false, nil
		}

		return formatWithRegex(s.r, val, s.Return), true, nil
	}
	return nil, false, policyauthor.NewKeyNotFoundError(s.Key)
}

// formatWithRegex takes a compiled regex, a target string, and a format string.
// It replaces placeholders like \1, \2, etc., in the format string with the corresponding matched groups.
func formatWithRegex(regex *regexp.Regexp, target, format string) string {
	matches := regex.FindStringSubmatch(target)
	if len(matches) == 0 {
		return "" // No match found
	}
	if matches[0] == "" {
		return "" // No match found
	}

	// Split the format string and replace placeholders with the matched groups.
	result := strings.Builder{}
	lastPos := 0
	for i := 0; i < len(format)-1; i++ {
		if format[i] == '\\' && '1' <= format[i+1] && format[i+1] <= '9' {
			// Append the part before the placeholder
			result.WriteString(format[lastPos:i])

			// Convert the character after '\' to an index
			index, _ := strconv.Atoi(string(format[i+1]))

			// Append the corresponding group if it exists
			if index < len(matches) {
				result.WriteString(matches[index])
			}

			// Skip over the number
			i++
			lastPos = i + 1
		}
	}

	// Append the remaining part of the format string
	if lastPos < len(format) {
		result.WriteString(format[lastPos:])
	}

	return result.String()
}
