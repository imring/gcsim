package target

import (
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

func (t *Target) WillApplyEle(tag info.ICDTag, grp info.ICDGroup, char int) float64 {
	// no icd if no tag
	if tag == info.ICDTagNone {
		return 1
	}

	// check if we need to start timer
	x := t.icdTagOnTimer[char][tag]
	if !t.icdTagOnTimer[char][tag] {
		t.icdTagOnTimer[char][tag] = true
		t.ResetTagCounterAfterDelay(tag, grp, char)
	}

	val := t.icdTagCounter[char][tag]
	t.icdTagCounter[char][tag]++

	// if counter > length, then use 0 for group seq
	groupSeq := info.ICDGroupEleApplicationSequence[grp][len(info.ICDGroupEleApplicationSequence[grp])-1]
	if val < len(info.ICDGroupEleApplicationSequence[grp]) {
		groupSeq = info.ICDGroupEleApplicationSequence[grp][val]
	}

	t.core.Log().NewEvent("ele icd check", glog.LogICDEvent, char).
		Write("grp", grp).
		Write("target", t.key).
		Write("tag", tag).
		Write("counter", val).
		Write("val", groupSeq).
		Write("group on timer", x)

	return groupSeq
}

func (t *Target) GroupTagDamageMult(tag info.ICDTag, grp info.ICDGroup, char int) float64 {
	// check if we need to start timer
	if !t.icdDamageTagOnTimer[char][tag] {
		t.icdDamageTagOnTimer[char][tag] = true
		t.ResetDamageCounterAfterDelay(tag, grp, char)
	}

	val := t.icdDamageTagCounter[char][tag]
	t.icdDamageTagCounter[char][tag]++

	// if counter > length, then use 0 for group seq
	groupSeq := info.ICDGroupDamageSequence[grp][len(info.ICDGroupDamageSequence[grp])-1]
	if val < len(info.ICDGroupDamageSequence[grp]) {
		groupSeq = info.ICDGroupDamageSequence[grp][val]
	}

	return groupSeq
}

func (t *Target) ResetDamageCounterAfterDelay(tag info.ICDTag, grp info.ICDGroup, char int) {
	t.core.Tasks().Add(func() {
		// set the counter back to 0
		t.icdDamageTagCounter[char][tag] = 0
		t.icdDamageTagOnTimer[char][tag] = false
		t.core.Log().NewEvent("damage counter reset", glog.LogICDEvent, char).
			Write("tag", tag).
			Write("grp", grp)
	}, info.ICDGroupResetTimer[grp]-1)
	t.core.Log().NewEvent("damage reset timer set", glog.LogICDEvent, char).
		Write("tag", tag).
		Write("grp", grp).
		Write("reset", t.core.F()+info.ICDGroupResetTimer[grp]-1)
}

func (t *Target) ResetTagCounterAfterDelay(tag info.ICDTag, grp info.ICDGroup, char int) {
	t.core.Tasks().Add(func() {
		// set the counter back to 0
		t.icdTagCounter[char][tag] = 0
		t.icdTagOnTimer[char][tag] = false
		t.core.Log().NewEvent("ele app counter reset", glog.LogICDEvent, char).
			Write("tag", tag).
			Write("grp", grp)
	}, info.ICDGroupResetTimer[grp]-1)
	t.core.Log().NewEvent("ele app reset timer set", glog.LogICDEvent, char).
		Write("tag", tag).
		Write("grp", grp).
		Write("reset", t.core.F()+info.ICDGroupResetTimer[grp]-1)
}
