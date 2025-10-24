package info

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/geometry"
)

type AttackTag int // attacktag is used instead of actions etc..

const (
	AttackTagNone AttackTag = iota
	AttackTagNormal
	AttackTagExtra
	AttackTagPlunge
	AttackTagElementalArt
	AttackTagElementalArtHold
	AttackTagElementalBurst
	AttackTagWeaponSkill
	AttackTagMonaBubbleBreak
	AttackTagNoneStat
	ReactionAttackDelim
	AttackTagOverloadDamage
	AttackTagSuperconductDamage
	AttackTagECDamage
	AttackTagShatter
	AttackTagSwirlPyro
	AttackTagSwirlHydro
	AttackTagSwirlCryo
	AttackTagSwirlElectro
	AttackTagBurningDamage
	AttackTagBloom
	AttackTagBountifulCore // special tag for nilou
	AttackTagBurgeon
	AttackTagHyperbloom
	AttackTagLength
)

type StrikeType int

const (
	StrikeTypeDefault StrikeType = iota
	StrikeTypePierce
	StrikeTypeBlunt
	StrikeTypeSlash
	StrikeTypeSpear
)

type AdditionalTag int

const (
	AdditionalTagNone AdditionalTag = iota
	AdditionalTagNightsoul
	AdditionalTagKinichCannon
)

type Attack struct {
	ActorIndex       int         // character this attack belongs to
	DamageSrc        keys.Target // source of this attack; should be a unique key identifying the target
	Abil             string      // name of ability triggering the damage
	AttackTag        AttackTag
	AdditionalTags   []AdditionalTag
	PoiseDMG         float64 // only needed on blunt attacks for frozen consumption before shatter for now
	ICDTag           ICDTag
	ICDGroup         ICDGroup
	Element          attributes.ElementType // element of ability
	Durability       Durability             // durability of aura, 0 if nothing applied
	NoImpulse        bool
	HitWeakPoint     bool
	Mult             float64 // ability multiplier. could set to 0 from initial Mona dmg
	StrikeType       StrikeType
	UseDef           bool    // we use this instead of flatdmg to make sure stat snapshotting works properly
	UseHP            bool    // we use this instead of flatdmg to make sure stat snapshotting works properly
	FlatDmg          float64 // flat dmg;
	IgnoreDefPercent float64 // by default this value is 0; if = 1 then the attack will ignore defense; raiden c2 should be set to 0.6 (i.e. ignore 60%)
	IgnoreInfusion   bool
	// amp info
	Amped   bool         // new flag used by new reaction system
	AmpMult float64      // amplier
	AmpType ReactionType // melt or vape i guess
	// catalyze info
	Catalyzed     bool
	CatalyzedType ReactionType
	// special flag for sim generated attack
	SourceIsSim bool
	DoNotLog    bool
	// hitlag stuff
	HitlagHaltFrames     float64 // this is the number of frames to pause by
	HitlagFactor         float64 // this is factor to slow clock by
	CanBeDefenseHalted   bool    // for whacking ruin gaurds
	IsDeployable         bool    // if this is true, then hitlag does not affect owner
	HitlagOnHeadshotOnly bool    // if this is true, will only apply if HitWeakpoint is also true
}

type AttackPattern struct {
	Shape       geometry.Shape
	SkipTargets [TargettableTypeCount]bool
	IgnoredKeys []keys.Target
}

type AttackCB struct {
	Target      keys.Target
	AttackEvent *AttackEvent
	Damage      float64
	IsCrit      bool
}

type AttackCBFunc func(AttackCB)

type AttackEvent struct {
	Info        Attack
	Pattern     AttackPattern
	Attacker    Snapshot
	Target      Snapshot
	SourceFrame int  // source frame
	Reacted     bool // true if a reaction already took place - for purpose of attach/refill

	Callbacks []AttackCBFunc `json:"-"`
}

type QueueAttack struct {
	Info    Attack
	Pattern AttackPattern

	SnapshotDelay int // ignored if snapshot is not nil
	DmgDelay      int

	Callbacks []AttackCBFunc
}
