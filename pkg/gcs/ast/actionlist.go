package ast

import (
	"encoding/json"
	"log"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character/profile"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

type ActionList struct {
	Targets     []enemy.EnemyProfile       `json:"targets"`
	PlayerPos   core.Coord                 `json:"player_initial_pos"`
	Characters  []profile.CharacterProfile `json:"characters"`
	InitialChar keys.Char                  `json:"initial"`
	Program     *BlockStmt                 `json:"-"`
	Energy      EnergySettings             `json:"energy_settings"`
	Settings    SimulatorSettings          `json:"settings"`
	Errors      []error                    `json:"-"` //These represents errors preventing ActionList from being executed
	ErrorMsgs   []string                   `json:"errors"`
}

type EnergySettings struct {
	Active         bool
	Once           bool //how often
	Start          int
	End            int
	Amount         int
	LastEnergyDrop int
}

type SimulatorSettings struct {
	Duration     float64
	DamageMode   bool
	EnableHitlag bool
	DefHalt      bool // for hitlag
	//other stuff
	NumberOfWorkers int // how many workers to run the simulation
	Iterations      int // how many iterations to run
	Delays          player.Delays
}

type Delays struct {
	Skill  int
	Burst  int
	Attack int
	Charge int
	Aim    int
	Dash   int
	Jump   int
	Swap   int
}

func (c *ActionList) Copy() *ActionList {

	r := *c

	r.Targets = make([]enemy.EnemyProfile, len(c.Targets))
	for i, v := range c.Targets {
		r.Targets[i] = v.Clone()
	}

	r.Characters = make([]profile.CharacterProfile, len(c.Characters))
	for i, v := range c.Characters {
		r.Characters[i] = v.Clone()
	}

	r.Program = c.Program.CopyBlock()
	return &r
}

func (a *ActionList) PrettyPrint() string {
	prettyJson, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	return string(prettyJson)
}
