package rules

import (
	"k8s.io/client-go/pkg/api/v1"
	"errors"
	"strings"
)

type Rule interface {
	load() (bool, error)
	ShouldReap(pod v1.Pod) (bool, string)
}

type Rules struct {
	LoadedRules []Rule
}

func LoadRules() (Rules, error) {
	// load all possible rules
	rules := []Rule{
		&duration{},
		&containerStatus{},
		&chaos{},
	}
	// return only the active rules
	loadedRules := []Rule{}
	for _, rule := range rules {
		load, err := rule.load()
		if err != nil {
			return Rules{LoadedRules:loadedRules}, err
		} else if load {
			loadedRules = append(loadedRules, rule)
		}
	}
	// return an err if no rules where loaded
	if len(loadedRules) == 0 {
		return Rules{LoadedRules:loadedRules}, errors.New("no rules were loaded")
	}
	return Rules{LoadedRules:loadedRules}, nil
}

func (rules Rules) ShouldReap(pod v1.Pod) (bool, string) {
	reasons := []string{}
	for _, rule := range rules.LoadedRules {
		reap, reason := rule.ShouldReap(pod)
		if !reap {
			return false, ""
		}
		reasons = append(reasons, reason)
	}
	return true, strings.Join(reasons, " AND ")
}
