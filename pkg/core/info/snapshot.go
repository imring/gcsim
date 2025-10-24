package info

import "github.com/genshinsim/gcsim/pkg/core/attributes"

type Snapshot struct {
	Stats attributes.Stats // total character stats including from artifact, bonuses, etc...
	Level int

	SourceFrame int           // frame snapshot was generated at
	Logs        []interface{} // logs for the snapshot
}
