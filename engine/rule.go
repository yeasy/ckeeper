package engine

import (
	"github.com/fsouza/go-dockerclient"
)

type Rule struct {
	name   string
	option docker.ListContainersOptions
	target string
	action string
}

type RuleSet struct {
	rules map[string]Rule
}

func NewRuleSet() *RuleSet {
	rs := RuleSet{}
	rs.rules = make(map[string]Rule)
	return &rs
}

func (rs *RuleSet) Init() {
}

func (rs *RuleSet) GetRules() map[string]Rule {
	return rs.rules
}

func (rs *RuleSet) AddRule(r Rule) {
	rs.rules[r.name] = r
}

func (rs *RuleSet) RemoveRule(name string) {
	delete(rs.rules, name)
}
