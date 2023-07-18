//go:build linux

package netlink

import (
	"github.com/vishvananda/netlink"
)

func NewRule() Rule {
	// defaults found from netlink.NewRule() for fields we use,
	// the rest of the defaults is set when converting from a `Rule`
	// to a `netlink.Rule`
	return Rule{
		Priority: -1,
		Mark:     -1,
	}
}

func (n *NetLink) RuleList(family int) (rules []Rule, err error) {
	netlinkRules, err := netlink.RuleList(family)
	if err != nil {
		return nil, err
	}

	rules = make([]Rule, len(netlinkRules))
	for i := range netlinkRules {
		rules[i] = netlinkRuleToRule(netlinkRules[i])
	}
	return rules, nil
}

func (n *NetLink) RuleAdd(rule Rule) error {
	netlinkRule := ruleToNetlinkRule(rule)
	return netlink.RuleAdd(&netlinkRule)
}

func (n *NetLink) RuleDel(rule Rule) error {
	netlinkRule := ruleToNetlinkRule(rule)
	return netlink.RuleDel(&netlinkRule)
}

func ruleToNetlinkRule(rule Rule) (netlinkRule netlink.Rule) {
	netlinkRule = *netlink.NewRule()
	netlinkRule.Priority = rule.Priority
	netlinkRule.Family = rule.Family
	netlinkRule.Table = rule.Table
	netlinkRule.Mark = rule.Mark
	netlinkRule.Src = netipPrefixToIPNet(rule.Src)
	netlinkRule.Dst = netipPrefixToIPNet(rule.Dst)
	netlinkRule.Invert = rule.Invert
	return netlinkRule
}

func netlinkRuleToRule(netlinkRule netlink.Rule) (rule Rule) {
	return Rule{
		Priority: netlinkRule.Priority,
		Family:   netlinkRule.Family,
		Table:    netlinkRule.Table,
		Mark:     netlinkRule.Mark,
		Src:      netIPNetToNetipPrefix(netlinkRule.Src),
		Dst:      netIPNetToNetipPrefix(netlinkRule.Dst),
		Invert:   netlinkRule.Invert,
	}
}
