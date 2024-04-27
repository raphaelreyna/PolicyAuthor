# PolicyAuthor

The PolicyAuthor package is a versatile library designed to evaluate policies defined in YAML against Go data structures. This package allows users to define complex policy rules using conditions such as equality checks, CIDR matches, regular expressions, etc. and apply these policies to dynamically evaluate data.

## Example config made possible by PolicyAuthor

```yaml
server:
  addr: :8080
  tls: {}
  routing:
    # PolicyAuthor YAML starts here
    - value: "https://mozilla.mysite.com"
      conditions:
        - type: regex
          spec:
            key: "headers.User-Agent"
            pattern: "^Mozilla/5.0.*"
    - value: "https://auth.mysite.com"
      conditions:
        - type: and
          spec:
            conditions:
              - type: equal
                spec:
                  key: "remote_addr"
                  value: "127.0.0.1"
              - type: exists
                spec:
                  key: "headers.X-My-Auth"
```

## Features

- _Flexible Policy Definitions_: Define your policies in YAML with support for multiple condition types.
- _Dynamic Data Evaluation_: Evaluate policies against Go data structures to determine compliance with defined rules.
- _Extensible_: Easily register new conditions to expand the functionality.

### Built-in conditions

- Logical conditions: and, or, not
- equal
- exists
- range
- regex
- substring
- cidr
- time

## Dev Example: Implementing Access Control

Here’s how you can use PolicyAuthor to enforce access control based on user location and request properties:

```go
package main

import (
    "github.com/raphaelreyna/policyauthor/pkg/conditions"
    "github.com/raphaelreyna/policyauthor/pkg/policy"
    "gopkg.in/yaml.v3"
    "log"
)

func main() {
    policy.RegisterConditions(conditions.AllConditionsMap())

    var config = `
    policies:
    - value: "AllowAccess"
      conditions:
      - type: and
        spec:
          conditions:
          - type: equal
            spec:
              key: "user_role"
              value: "admin"
          - type: not
            spec:
              condition:
                type: cidr
                spec:
                  key: "remote_addr"
                  value: "192.168.1.0/24"
    `

    var p struct {
        Policies *policy.PolicyEngine `yaml:"policies"`
    }

    err := yaml.Unmarshal([]byte(config), &p)
    if err != nil {
        log.Fatalf("error: %v", err)
    }

    userData := map[string]any{
        "user_role": "admin",
        "remote_addr": "192.168.2.1",
    }

    value, ok, err := p.Policies.Evaluate(userData)
    if err != nil {
        log.Fatalf("evaluation error: %v", err)
    }
    if ok && value == "AllowAccess" {
        log.Println("Access granted")
    } else {
        log.Println("Access denied")
    }
}
```