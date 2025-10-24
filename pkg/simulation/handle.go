package simulation

func (c *Core) handleEnergy() {
	// energy once interval=300 amount=1 #once at frame 300
	// if s.cfg.EnergySettings.Active && s.cfg.EnergySettings.Once {
	// 	f := s.cfg.EnergySettings.Start
	// 	s.cfg.EnergySettings.Active = false
	// 	s.C.Tasks.Add(func() {
	// 		s.C.Player.DistributeParticle(character.Particle{
	// 			Source: "enemy",
	// 			Num:    float64(s.cfg.EnergySettings.Amount),
	// 			Ele:    attributes.NoElement,
	// 		})
	// 	}, f)
	// 	s.C.Log.NewEventBuildMsg(glog.LogEnergyEvent, -1, "energy queued (once)").
	// 		Write("last", s.cfg.EnergySettings.LastEnergyDrop).
	// 		Write("cfg", s.cfg.EnergySettings).
	// 		Write("amt", s.cfg.EnergySettings.Amount).
	// 		Write("energy_frame", s.C.F+f)
	// }
	// // energy every interval=300,600 amount=1 #randomly every 300 to 600 frames
	// if s.cfg.EnergySettings.Active && s.C.F-s.cfg.EnergySettings.LastEnergyDrop >= s.cfg.EnergySettings.Start {
	// 	f := s.C.Rand.Intn(s.cfg.EnergySettings.End - s.cfg.EnergySettings.Start)
	// 	s.cfg.EnergySettings.LastEnergyDrop = s.C.F + f
	// 	s.C.Tasks.Add(func() {
	// 		s.C.Player.DistributeParticle(character.Particle{
	// 			Source: "drop",
	// 			Num:    float64(s.cfg.EnergySettings.Amount),
	// 			Ele:    attributes.NoElement,
	// 		})
	// 	}, f)
	// 	s.C.Log.NewEventBuildMsg(glog.LogEnergyEvent, -1, "energy queued").
	// 		Write("last", s.cfg.EnergySettings.LastEnergyDrop).
	// 		Write("cfg", s.cfg.EnergySettings).
	// 		Write("amt", s.cfg.EnergySettings.Amount).
	// 		Write("energy_frame", s.C.F+f)
	// }
}

func (c *Core) handleHurt() {
	// hurt once interval=300 amount=1,300 element=physical #once at frame 300 (or nearest)
	// if s.cfg.HurtSettings.Active && s.cfg.HurtSettings.Once {
	// 	f := s.cfg.HurtSettings.Start
	// 	amt := s.cfg.HurtSettings.Min + s.C.Rand.Float64()*(s.cfg.HurtSettings.Max-s.cfg.HurtSettings.Min)
	// 	s.cfg.HurtSettings.Active = false

	// 	s.C.Tasks.Add(func() {
	// 		ai := info.AttackInfo{
	// 			ActorIndex:       s.C.Player.Active(),
	// 			Abil:             "Hurt",
	// 			AttackTag:        attacks.AttackTagNone,
	// 			ICDTag:           attacks.ICDTagNone,
	// 			ICDGroup:         attacks.ICDGroupDefault,
	// 			StrikeType:       attacks.StrikeTypeDefault,
	// 			Durability:       0,
	// 			Element:          s.cfg.HurtSettings.Element,
	// 			FlatDmg:          amt,
	// 			IgnoreDefPercent: 1,
	// 		}
	// 		ap := combat.NewSingleTargetHit(s.C.Combat.Player().Key())
	// 		ap.SkipTargets[info.TargettablePlayer] = false
	// 		ap.SkipTargets[info.TargettableEnemy] = true
	// 		ap.SkipTargets[info.TargettableGadget] = true
	// 		s.C.QueueAttack(ai, ap, -1, 0) // -1 to avoid snapshot
	// 	}, f)

	// 	s.C.Log.NewEventBuildMsg(glog.LogHurtEvent, -1, "hurt queued (once)").
	// 		Write("last", s.cfg.HurtSettings.LastHurt).
	// 		Write("cfg", s.cfg.HurtSettings).
	// 		Write("amt", amt).
	// 		Write("hurt_frame", s.C.F+f)
	// }
	// // hurt every interval=480,720 amount=1,300 element=physical #randomly 1 to 300 dmg every 480 to 720 frames
	// if s.cfg.HurtSettings.Active && s.C.F-s.cfg.HurtSettings.LastHurt >= s.cfg.HurtSettings.Start {
	// 	f := s.C.Rand.Intn(s.cfg.HurtSettings.End - s.cfg.HurtSettings.Start)
	// 	amt := s.cfg.HurtSettings.Min + s.C.Rand.Float64()*(s.cfg.HurtSettings.Max-s.cfg.HurtSettings.Min)
	// 	s.cfg.HurtSettings.LastHurt = s.C.F + f

	// 	s.C.Tasks.Add(func() {
	// 		ai := info.AttackInfo{
	// 			ActorIndex:       s.C.Player.Active(),
	// 			Abil:             "Hurt",
	// 			AttackTag:        attacks.AttackTagNone,
	// 			ICDTag:           attacks.ICDTagNone,
	// 			ICDGroup:         attacks.ICDGroupDefault,
	// 			StrikeType:       attacks.StrikeTypeDefault,
	// 			Durability:       0,
	// 			Element:          s.cfg.HurtSettings.Element,
	// 			FlatDmg:          amt,
	// 			IgnoreDefPercent: 1,
	// 		}
	// 		ap := combat.NewSingleTargetHit(s.C.Combat.Player().Key())
	// 		ap.SkipTargets[info.TargettablePlayer] = false
	// 		ap.SkipTargets[info.TargettableEnemy] = true
	// 		ap.SkipTargets[info.TargettableGadget] = true
	// 		s.C.QueueAttack(ai, ap, -1, 0) // -1 to avoid snapshot
	// 	}, f)

	// 	s.C.Log.NewEventBuildMsg(glog.LogHurtEvent, -1, "hurt queued").
	// 		Write("last", s.cfg.HurtSettings.LastHurt).
	// 		Write("cfg", s.cfg.HurtSettings).
	// 		Write("amt", amt).
	// 		Write("hurt_frame", s.C.F+f)
	// }
}
