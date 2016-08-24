package engine

import (
	"testing"

	"github.com/fsouza/go-dockerclient"
)

func TestRuleSet_AddRule(t *testing.T) {
	rs := NewRuleSet()
	if len(rs.GetRules()) != 0 {
		t.Error("Initial rulest lengh not equal to 0")
	}

	rule := Rule{
		"test_rule",
		docker.ListContainersOptions{All: true},
		"condition",
		"action",
	}

	rs.AddRule(rule)

	result := rs.GetRules()
	if result["test_rule"].target != "condition" {
		t.Error("Mismatch condition filed on the inserted rule")
	}
	if len(result) != 1 {
		t.Error("Error to add a new rule")
	}

}

func TestRuleSet_RemoveRule(t *testing.T) {
	rs := NewRuleSet()

	rule := Rule{
		"test_rule",
		docker.ListContainersOptions{All: true},
		"condition",
		"action",
	}

	rs.AddRule(rule)
	rs.RemoveRule("test_rule")
	if len(rs.GetRules()) != 0 {
		t.Error("Error to remove a rule")
	}
}
