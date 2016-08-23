package engine

import (
	"github.com/fsouza/go-dockerclient"
)

type Rule struct {
	filter    docker.ListContainersOptions
	condition string
	action    string
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

func (rs *RuleSet) AddRule(name string, r Rule) {
	rs.rules[name] = r
}

func (rs *RuleSet) RemoveRule(name string) {
	delete(rs.rules, name)
}
