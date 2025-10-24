package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/genshinsim/gcsim/pkg/gcs/eval"
	"github.com/genshinsim/gcsim/pkg/gcs/parser"
	"github.com/genshinsim/gcsim/pkg/simulation"
)

const cfg = `
bennett char lvl=90/90 cons=6 talent=9,9,9;
bennett add weapon="mist" refine=5 lvl=90/90;
bennett add set="cwof" count=4;
bennett add stats hp=4780 atk=311 atk%=0.583 pyro%=0.466 cr=0.311 ; # main
bennett add stats def=39.36 def%=0.124 hp=507.88 hp%=0.0992 atk=33.08 atk%=0.496 er=0.1102 em=39.64 cr=0.2979 cd=0.4634 ;	

#enemy config

options swap_delay=12 iteration=1000;
target lvl=100 resist=0.1 radius=2 pos=0,2.4; 
energy every interval=480,720 amount=1;
active bennett;

#action list 

for { 
	switch {
	#case .bennett.burst.ready && .bennett.energy == .bennett.energymax:
	#	bennett burst;
	#case .bennett.skill.ready:
	#	bennett skill;
	case .bennett.normal > 4:
		bennett dash;
	default:
		bennett attack;
	}
}
`

func TestMain(t *testing.T) {
	// parse cfg
	p := parser.New(cfg)
	cfg, gcsl, err := p.Parse()
	if err != nil {
		panic(err)
	}

	fmt.Println(gcsl.String())

	// create new eval
	eval, err := eval.NewEvaluator(gcsl)
	if err != nil {
		panic(err)
	}

	// create simulation
	sim, err := simulation.New(simulation.Opt{
		Config:            cfg,
		Eval:              eval,
		Seed:              0,
		Debug:             true,
		DefHalt:           cfg.Settings.DefHalt,
		EnableHitlag:      cfg.Settings.EnableHitlag,
		DamageMode:        cfg.Settings.DamageMode,
		IgnoreBurstEnergy: cfg.Settings.IgnoreBurstEnergy,
		Delays:            cfg.Settings.Delays,
	})
	if err != nil {
		panic(err)
	}

	// run simulation
	_, err = sim.Run()
	if err != nil {
		panic(err)
	}

	// fmt.Println(res)

	logs, err := sim.Log().Dump()
	if err != nil {
		panic(err)
	}
	// fmt.Println(string(logs))

	os.Remove("logs.json")
	os.WriteFile("logs.json", logs, 0o600)
}
