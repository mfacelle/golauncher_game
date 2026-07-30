package game

// holds into about the player, such as
// - number of levels completed
// - remaining projectiles
type PlayerInfo struct {
	NumParticles  int
	LevelsCleared int
}

func (p PlayerInfo) CalculateScore() int {
	if p.NumParticles == 0 {
		return p.LevelsCleared
	} else {
		return p.NumParticles * p.LevelsCleared
	}
}
