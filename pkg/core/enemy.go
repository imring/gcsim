package core

type Enemy interface {
	Target
	Reactable

	// hp related
	MaxHP() float64
	HP() float64
	// hitlag related
	ApplyHitlag(factor, dur float64)
	QueueEnemyTask(f func(), delay int)
	// modifier related
	// add
	// AddStatus(key string, dur int, hitlag bool)
	// AddResistMod(mod ResistMod) // TODO: change after refactor
	// AddDefMod(mod DefMod)
	// delete
	// DeleteStatus(key string)
	// DeleteResistMod(key string)
	// DeleteDefMod(key string)
	// StatusExpiry(key string) int
}
