package combat

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/geometry"
)

func (h *Handler) Enemy(i int) core.Enemy {
	if i < 0 || i > len(h.enemies) {
		return nil
	}
	return h.enemies[i]
}

func (h *Handler) SetEnemyPos(i int, p geometry.Point) bool {
	if i < 0 || i > len(h.enemies)-1 {
		return false
	}

	h.enemies[i].SetPos(p)
	return true
}

func (h *Handler) KillEnemy(i int) {
	h.enemies[i].Kill(nil)
	h.log.NewEvent("enemy dead", glog.LogSimEvent, -1).Write("index", i)
}

func (h *Handler) AddEnemy(t core.Enemy) {
	h.enemies = append(h.enemies, t)
	h.target.Add(t)
}

func (h *Handler) Enemies() []core.Enemy {
	return h.enemies
}

func (h *Handler) EnemyCount() int {
	return len(h.enemies)
}

func (h *Handler) PrimaryTarget() core.Target {
	for _, v := range h.enemies {
		if v.Key() == h.defaultTarget {
			if !v.IsAlive() {
				h.log.NewEvent("default target is dead", glog.LogWarnings, -1)
			}
			return v
		}
	}
	panic("default target does not exist?!")
}
