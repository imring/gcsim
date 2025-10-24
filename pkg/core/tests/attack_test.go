package tests

import (
	"math/rand"
	"testing"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/target"
	"github.com/genshinsim/gcsim/pkg/core/target/enemy"
	"github.com/genshinsim/gcsim/pkg/core/task"
	"github.com/genshinsim/gcsim/tests/mock"
	"go.uber.org/mock/gomock"
)

func NewMockCore(ctrl *gomock.Controller, f *int) *mock.MockCore {
	core := mock.NewMockCore(ctrl)
	events := &event.System{}
	tasks := task.New(f)
	core.EXPECT().F().Return(0).AnyTimes()
	core.EXPECT().Events().Return(events).AnyTimes()
	core.EXPECT().IsDamageMode().Return(true).AnyTimes()
	core.EXPECT().Log().Return(&glog.NilLogger{}).AnyTimes()
	core.EXPECT().Rand().Return(rand.New(rand.NewSource(0))).AnyTimes()
	core.EXPECT().Tasks().Return(tasks).AnyTimes()

	return core
}

func NewCharacter(core *mock.MockCore) core.Character {
	char, err := character.New(character.Opt{
		Core:    core,
		Profile: info.CharacterProfile{},
	})
	if err != nil {
		panic("failed to create character")
	}
	return char
}

func NewPlayerHandler(core *mock.MockCore) *player.Handler {
	h := player.New(player.Opt{})
	h.AddChar(NewCharacter(core))
	return h
}

func NewEnemy(core *mock.MockCore) core.Target {
	return enemy.New(core, info.EnemyProfile{})
}

func NewTarget(core *mock.MockCore) *target.Handler {
	return target.New()
}

func NewCombatHandler(core *mock.MockCore, f *int, targetHandler *target.Handler) *combat.Handler {
	return combat.New(combat.Opt{
		F:      f,
		Events: core.Events(),
		Player: NewPlayerHandler(core),
		Target: targetHandler,
	})
}

func TestAttack(t *testing.T) {
	ctrl := gomock.NewController(t)

	f := 0
	core := NewMockCore(ctrl, &f)

	targetHandler := NewTarget(core)
	enemy := NewEnemy(core)
	targetHandler.Add(enemy)

	count := 0
	combatHandler := NewCombatHandler(core, &f, targetHandler)
	core.Events().EnemyHit.Subscribe(func(event.TargetHitEvent) {
		count++
	})

	combatHandler.QueueAttack(info.QueueAttack{
		Info:    info.Attack{},
		Pattern: combat.NewSingleTargetHit(enemy.Key()),
	})

	if count != 1 {
		t.Errorf("expected 1 enemy hit, got %d", count)
	}
}
