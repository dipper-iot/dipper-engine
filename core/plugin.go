package core

type Plugin interface {
	Forms()
	Rule() Rule
	Get()
	Save()
}
