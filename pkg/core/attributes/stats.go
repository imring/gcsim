package attributes

type ModifierChange struct {
	Props  Props
	Reason string
	Expiry int
}

type PropChange struct {
	PropMap PropMap
	Reason  string
	Expiry  int
}

type Stats struct {
	Props           Props
	ModifierChanges []ModifierChange
	Changes         []PropChange
}

func NewStats(base Props, modifierChanges []ModifierChange) Stats {
	props := Props{}
	props.Add(base)
	for _, m := range modifierChanges {
		props.Add(m.Props)
	}

	return Stats{
		Props:           props,
		ModifierChanges: modifierChanges,
		Changes:         make([]PropChange, 0, len(modifierChanges)),
	}
}

func (s *Stats) AddProperty(propChange PropChange) {
	for k, v := range propChange.PropMap {
		s.Props[k] += v
	}
	s.Changes = append(s.Changes, propChange)
}
