package interactive

type Requirement interface {
	Satisfied() bool
}

type CommonRequirement struct {
	requirement func() bool
}

type Element struct {
	action       func()
	requirements []Requirement
}

func NewCommonRequirement(requirement func() bool) CommonRequirement {
	return CommonRequirement{
		requirement: requirement,
	}
}

func NewElement(action func(), requirements ...Requirement) Element {
	return Element{
		action:       action,
		requirements: requirements,
	}
}

func (cr CommonRequirement) Satisfied() bool {
	return cr.requirement()
}

func (e Element) available() bool {
	for _, requirement := range e.requirements {
		if !requirement.Satisfied() {
			return false
		}
	}

	return true
}

func (e Element) RunAction() {
	if !e.available() {
		return
	}

	e.action()
}
