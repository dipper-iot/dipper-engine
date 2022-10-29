package core

type SessionControl interface {
	ListSession() []uint64
	StopSession(id uint64)
	InfoSession(id uint64) map[string]interface{}
}

func (d *DipperEngine) ListControl() []string {
	list := make([]string, 0)

	for name, _ := range d.mapSessionControl {
		list = append(list, name)
	}
	return list
}

func (d *DipperEngine) ControlSession(ruleName string) []uint64 {

	rule, ok := d.mapSessionControl[ruleName]
	if ok {
		return rule.ListSession()
	}
	return []uint64{}
}

func (d *DipperEngine) ControlGetRule(session uint64) []string {
	listRule := make([]string, 0)
	for ruleId, control := range d.mapSessionControl {
		listSession := control.ListSession()
		for _, name := range listSession {
			if name == session {
				listRule = append(listRule, ruleId)
				break
			}
		}
	}

	return listRule
}

func (d *DipperEngine) ControlStopSession(ruleName string, session uint64) {

	rule, ok := d.mapSessionControl[ruleName]
	if ok {
		rule.StopSession(session)
	}
	return
}
