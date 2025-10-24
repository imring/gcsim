package reactable

// func TestBurning(t *testing.T) {
// 	c := testCore()
// 	trg := addTargetToCore(c)

// 	c.Init()

// 	//TODO: write tests for burning (this is copypasted from quicken for now)
// 	trg.AttachOrRefill(&info.AttackEvent{
// 		Info: combat.AttackInfo{
// 			Element:    attributes.ElementGrass,
// 			Durability: 25,
// 		},
// 	})
// 	trg.React(&info.AttackEvent{
// 		Info: combat.AttackInfo{
// 			Element:    attributes.Electro,
// 			Durability: 25,
// 		},
// 	})
// 	// dendro electro gone; 20 quicken
// 	if !durApproxEqual(20, trg.Durability[attributes.ElementOverdose], 0.00001) {
// 		t.Errorf("expecting 20 cryo attached, got %v", trg.Durability[attributes.ElementOverdose])
// 	}
// 	if trg.AuraContains(attributes.ElementGrass, attributes.Electro) {
// 		t.Error("expecting target to not contain any remaining dendro or electro aura")
// 	}
// }
