package tags

type Tags map[string]int

func (t Tags) Get(tag string) int      { return t[tag] }
func (t Tags) Set(tag string, val int) { t[tag] = val }
func (t Tags) Remove(tag string)       { delete(t, tag) }
