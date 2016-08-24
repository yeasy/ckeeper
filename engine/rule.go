package engine

import (
	"github.com/fsouza/go-dockerclient"
)

// Rule represents one rule
type Rule struct {
	name   string
	option docker.ListContainersOptions
	target string
	action string
}

// RuleSet represents a group of rules
type RuleSet struct {
	rules map[string]Rule
}

// NewRuleSet will return a RuleSet object
func NewRuleSet() *RuleSet {
	rs := RuleSet{}
	rs.rules = make(map[string]Rule)
	return &rs
}

// GetRules will return all rules in the RuleSet
func (rs *RuleSet) GetRules() map[string]Rule {
	return rs.rules
}

// AddRule will add one rule into the RuleSet
func (rs *RuleSet) AddRule(r Rule) {
	rs.rules[r.name] = r
}

// RemoveRule will remove one rule from the RuleSet
func (rs *RuleSet) RemoveRule(name string) {
	delete(rs.rules, name)
}
