package models

import "time"

type KilledStep struct {
	BossStepID string    `bson:"bossStepID" json:"boss_step_id"`
	KilledAt   time.Time `bson:"killedAt" json:"killed_at"`
}

type Run struct {
	CustomID    string       `bson:"customID" json:"id"`
	DungeonID   string       `bson:"dungeonID" json:"dungeon_id"`
	PlayerID    string       `bson:"playerID" json:"player_id"`
	State       string       `bson:"state" json:"state"`
	CurrentStep int          `bson:"currentStep" json:"current_step"`
	KilledSteps []KilledStep `bson:"killedSteps" json:"killed_steps"`
	StartedAt   time.Time    `bson:"startedAt" json:"started_at"`
	EndedAt     *time.Time   `bson:"endedAt,omitempty" json:"ended_at,omitempty"`
}

func (r *Run) Collection() string {
	return "run"
}
