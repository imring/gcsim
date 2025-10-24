package event

type System struct {
	// combat
	ApplyAttack ApplyAttackEventHandler
	EnemyHit    TargetHitEventHandler
	EnemyDamage EnemyDamageEventHandler
	GadgetHit   TargetHitEventHandler

	// construct
	ConstructSpawned ConstructSpawnedEventHandler

	// player
	EnergyChange    EnergyChangeEventHandler
	HPDebt          HPDebtEventHandler
	CharacterAction CharacterActionEventHandler
	StamUse         StamUseEventHandler
	StateChange     StateChangeEventHandler
	ActionFailed    ActionFailedEventHandler
	CharacterSwap   CharacterSwapEventHandler
	ActionExec      ActionExecEventHandler
	Dash            ActionEventHandler
	Skill           ActionEventHandler
	Burst           ActionEventHandler
	Attack          ActionEventHandler
	ChargeAttack    ActionEventHandler
	Plunge          ActionEventHandler
	AimShoot        ActionEventHandler

	// reaction
	Bloom                  ReactionEventHandler
	DendroCoreSpawned      DendroCoreSpawnedEventHandler
	Hyperbloom             ReactionEventHandler
	Burgeon                ReactionEventHandler
	Burning                ReactionEventHandler
	Aggravate              ReactionEventHandler
	Spread                 ReactionEventHandler
	Quicken                ReactionEventHandler
	CrystallizeElectro     ReactionEventHandler
	CrystallizeCryo        ReactionEventHandler
	CrystallizeHydro       ReactionEventHandler
	CrystallizePyro        ReactionEventHandler
	ElectroCharged         ReactionEventHandler
	Frozen                 ReactionEventHandler
	Shatter                ReactionEventHandler
	Melt                   ReactionEventHandler
	Overload               ReactionEventHandler
	Superconduct           ReactionEventHandler
	SwirlElectro           ReactionEventHandler
	SwirlCryo              ReactionEventHandler
	SwirlHydro             ReactionEventHandler
	SwirlPyro              ReactionEventHandler
	Vaporize               ReactionEventHandler
	AuraDurabilityAdded    AuraDurabilityAddedEventHandler
	AuraDurabilityDepleted AuraDurabilityDepletedEventHandler

	// shield
	Shielded    ShieldedEventHandler
	ShieldBreak ShieldBreakEventHandler

	// target
	TargetDied  TargetDiedEventHandler
	TargetMoved TargetMovedEventHandler

	// simulation
	Tick                 TickEventHandler
	SimEndedSuccessfully SimEndedSuccessfullyHandler
}
