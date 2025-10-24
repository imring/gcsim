package character

import (
	"fmt"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/modifier"
	"github.com/genshinsim/gcsim/pkg/core/target"
	"github.com/genshinsim/gcsim/pkg/core/task"
	"github.com/genshinsim/gcsim/pkg/geometry"
	"github.com/genshinsim/gcsim/pkg/model"
	"github.com/genshinsim/gcsim/pkg/tags"
)

type Character struct {
	*target.Target

	core core.Core

	Index int

	// base
	Base    info.CharacterBase
	Weapon  info.WeaponProfile
	Talents info.TalentProfile
	Sets    info.Sets
	data    *model.AvatarData

	Equip    core.EquipInfo
	CharZone info.ZoneType
	CharBody info.BodyType

	NormalCon int
	SkillCon  int
	BurstCon  int

	// stats
	BaseProps attributes.Props
	Modifiers *modifier.Handler
	Tags      tags.Tags

	// current status
	ParticleDelay int // character custom particle delay
	energy        float64
	energyMax     float64

	// normal attack counter
	normalHitNum  int // how many hits in a normal combo
	normalCounter int

	// hp
	currentHPRatio float64
	currentHPDebt  float64

	// cd
	ActionCD               []int
	cdQueueWorkerStartedAt []int
	cdCurrentQueueWorker   []*func()
	cdQueue                [][]int
	AvailableCDCharge      []int
	additionalCDCharge     []int

	// dash cd: keeps track of remaining cd frames for off-field chars
	remainingDashCD int
	dashLockout     bool

	// hitlag stuff
	TimePassed   int // how many frames have passed since start of sim
	frozenFrames int // how many frames are we still frozen for
	queue        *task.Handler
}

type Opt struct {
	Core    core.Core
	Profile info.CharacterProfile
	Data    *model.AvatarData

	EnergyMax    float64
	NormalHitNum int
	NormalCon    int
	SkillCon     int
	BurstCon     int
}

func New(opt Opt) (*Character, error) {
	c := &Character{
		core:         opt.Core,
		Base:         opt.Profile.Base,
		Weapon:       opt.Profile.Weapon,
		Talents:      opt.Profile.Talents,
		data:         opt.Data,
		NormalCon:    opt.NormalCon,
		SkillCon:     opt.SkillCon,
		BurstCon:     opt.BurstCon,
		energyMax:    opt.EnergyMax,
		normalHitNum: opt.NormalHitNum,

		ActionCD:               make([]int, info.EndActionType),
		cdQueueWorkerStartedAt: make([]int, info.EndActionType),
		cdCurrentQueueWorker:   make([]*func(), info.EndActionType),
		cdQueue:                make([][]int, info.EndActionType),
		AvailableCDCharge:      make([]int, info.EndActionType),
		additionalCDCharge:     make([]int, info.EndActionType),
	}

	c.Equip.Sets = make(map[keys.Set]core.Set)
	copy(c.BaseProps[:], opt.Profile.Stats)

	if c.NormalCon <= 0 {
		c.NormalCon = -1
	}
	if c.SkillCon <= 0 {
		c.SkillCon = -1
	}
	if c.BurstCon <= 0 {
		c.BurstCon = -1
	}

	c.Target = target.NewTarget(opt.Core, geometry.Point{}, 0) // TODO: radius
	c.Modifiers = modifier.New(c.core, c.Target.Key())
	c.queue = task.New(&c.TimePassed)

	if c.Talents.Attack < 1 || c.Talents.Attack > 10 {
		return nil, fmt.Errorf("invalid talent lvl: attack - %v", c.Talents.Attack)
	}
	if c.Talents.Skill < 1 || c.Talents.Skill > 10 {
		return nil, fmt.Errorf("invalid talent lvl: skill - %v", c.Talents.Skill)
	}
	if c.Talents.Burst < 1 || c.Talents.Burst > 10 {
		return nil, fmt.Errorf("invalid talent lvl: burst - %v", c.Talents.Burst)
	}

	for i := 0; i < len(c.cdQueue); i++ {
		c.cdQueue[i] = make([]int, 0, 4)
		c.AvailableCDCharge[i] = 1
	}

	return c, nil
}

func (c *Character) Init() error        { return nil }
func (c *Character) GetIndex() int      { return c.Index }
func (c *Character) SetIndex(index int) { c.Index = index }

func (c *Character) GetBase() info.CharacterBase    { return c.Base }
func (c *Character) GetWeapon() info.WeaponProfile  { return c.Weapon }
func (c *Character) GetTalents() info.TalentProfile { return c.Talents }
func (c *Character) GetSets() info.Sets             { return c.Sets }
func (c *Character) Data() *model.AvatarData        { return c.data }

func (c *Character) GetEquip() core.EquipInfo { return c.Equip }

func (c *Character) Condition(fields []string) (any, error) {
	return false, fmt.Errorf("invalid character condition: .%v.%v", c.Base.Key.String(), strings.Join(fields, "."))
}

func (c *Character) consCheck() {
	consUnset := 0
	if c.NormalCon < 0 {
		consUnset++
	}
	if c.SkillCon < 0 {
		consUnset++
	}
	if c.BurstCon < 0 {
		consUnset++
	}
	if consUnset != 1 {
		panic(fmt.Sprintf("cons not set properly for %v, please set two out of three values:\nNormalCon: %v\nSkillCon: %v\nBurstCon: %v", c.Base.Key.String(), c.NormalCon, c.SkillCon, c.BurstCon))
	}
}

func (c *Character) TalentLvlAttack() int {
	c.consCheck()
	add := -1
	if c.Tag(keys.ChildePassive) > 0 {
		add++
	}
	if c.NormalCon > 0 && c.Base.Cons >= c.NormalCon {
		add += 3
	}
	if add >= 4 {
		add = 4
	}
	return c.Talents.Attack + add
}

func (c *Character) TalentLvlSkill() int {
	c.consCheck()
	add := -1
	if c.Tag(keys.SkirkPassive) > 0 {
		add++
	}
	if c.SkillCon > 0 && c.Base.Cons >= c.SkillCon {
		add += 3
	}
	if add >= 4 {
		add = 4
	}
	return c.Talents.Skill + add
}

func (c *Character) TalentLvlBurst() int {
	c.consCheck()
	add := -1
	if c.BurstCon > 0 && c.Base.Cons >= c.BurstCon {
		add += 3
	}
	if add >= 4 {
		add = 4
	}
	return c.Talents.Burst + add
}

func (c *Character) SetWeapon(w core.Weapon) {
	c.Equip.Weapon = w
}

func (c *Character) SetArtifactSet(key keys.Set, s core.Set) {
	c.Equip.Sets[key] = s
}

func (c *Character) SetDirectionToClosestEnemy() {
	src := c.Pos()
	// calculate direction towards closest enemy, or forward direction if none
	enemy := c.core.ClosestEnemy(src)
	if enemy == nil {
		c.SetDirection(geometry.DefaultDirection())
		return
	}
	c.SetDirection(enemy.Pos())
	c.core.Log().NewEvent("set target direction to closest enemy", glog.LogDebugEvent, -1).
		Write("enemy key", enemy.Key()).
		Write("direction", c.Direction())
}
