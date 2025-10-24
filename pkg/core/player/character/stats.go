package character

import (
	"strconv"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

func (c *Character) Stats() attributes.Stats {
	return attributes.NewStats(c.BaseProps, c.Modifiers.Stats())
}

func (c *Character) Stat(attr attributes.Prop) float64 {
	// TODO: optimize with caching and OnPropChange event?
	return c.Stats().Props[attr]
}

func (c *Character) Snapshot(attackInfo *info.Attack) info.Snapshot {
	var sb strings.Builder
	var debug []any
	var evt glog.Event

	if attackInfo != nil {
		evt = c.core.Log().NewEvent(attackInfo.Abil, glog.LogSnapshotEvent, c.Index).
			Write("abil", attackInfo.Abil).
			Write("mult", attackInfo.Mult).
			Write("ele", attackInfo.Element.String()).
			Write("durability", float64(attackInfo.Durability)).
			Write("icd_tag", attackInfo.ICDTag).
			Write("icd_group", attackInfo.ICDGroup)
	}

	// snapshot the stats
	s := info.Snapshot{
		Stats:       c.Stats(),
		Level:       c.Base.Level,
		SourceFrame: c.core.F(),
	}

	// TODO: infusion

	// logs
	for _, v := range s.Stats.ModifierChanges {
		sb.WriteString(v.Reason)
		modStatus := make([]string, 0, 2)
		modStatus = append(modStatus,
			"status: added",
			"expiry_frame: "+strconv.Itoa(v.Expiry),
		)
		modStatus = append(
			modStatus,
			attributes.PrettyPrintStatsSlice(v.Props[:])...,
		)
		debug = append(debug, sb.String(), modStatus)
		sb.Reset()
	}

	if attackInfo != nil {
		evt.WriteBuildMsg(debug...)
		evt.Write("final_stats", attributes.PrettyPrintStatsSlice(s.Stats.Props[:]))
		// if inf != attributes.ElementNone {
		// 	evt.Write("infused_ele", inf.String())
		// }
	}
	s.Logs = debug

	return s
}

func (c *Character) StatusDuration(status string) int {
	return c.Modifiers.GetDuration(status)
}

func (c *Character) StatusIsActive(status string) bool {
	return c.Modifiers.HasModifier(status)
}

func (c *Character) Tag(tag string) int {
	return c.Tags.Get(tag)
}

func (c *Character) SetTag(tag string, val int) {
	c.Tags.Set(tag, val)
}

func (c *Character) AddModifier(m info.Modifier) {
	c.Modifiers.Add(m)
}

func (c *Character) RemoveModifier(key string) bool {
	return c.Modifiers.Remove(key)
}
