package model

// Stats for player
type Stats struct {
	Name             string  `json:"name"`
	Wins             int     `json:"wins"`
	Ties             int     `json:"ties"`
	Losses           int     `json:"losses"`
	Games            int     `json:"games"`
	WinPercentage    float64 `json:"winPercentage"`
	HighestWinStreak int     `json:"highestWinStreak"`
	CurrentWinStreak int     `json:"currentWinStreak"`
}

// StatsMap holds stats for all players
type StatsMap map[string]*Stats

// SortedStats is sorted
type SortedStats []*Stats

// Len is part of sort.Interface.
func (ss SortedStats) Len() int {
	return len(ss)
}

// Swap is part of sort.Interface.
func (ss SortedStats) Swap(i, j int) {
	ss[i], ss[j] = ss[j], ss[i]
}

// Less is part of sort.Interface. We use count as the value to sort by
func (ss SortedStats) Less(i, j int) bool {
	return ss[i].WinPercentage > ss[j].WinPercentage
}

// Check if key exist in map
func (sm StatsMap) hasKey(name string) bool {
	_, ok := sm[name]
	if ok {
		return true
	}
	return false
}

func (sm StatsMap) initPlayerStats(name string) {
	player := new(Stats)
	player.Wins = 0
	player.Ties = 0
	player.Losses = 0
	player.WinPercentage = 0
	player.CurrentWinStreak = 0
	player.HighestWinStreak = 0
	sm[name] = player
}

// Init stats to 0 if player doesnt exist
func (sm StatsMap) Init(winner string, loser string) {
	// Init if winner doesn't exist
	if !sm.hasKey(winner) {
		sm.initPlayerStats(winner)
	}

	// Init if loser doesn't exist
	if !sm.hasKey(loser) {
		sm.initPlayerStats(loser)
	}
}

// Increase stats for winner and loser
func (sm StatsMap) Increase(winner string, loser string, isTie bool) {
	if isTie {
		sm[winner].Ties++
		sm[loser].Ties++

		sm[winner].CurrentWinStreak = 0
		sm[loser].CurrentWinStreak = 0
	} else {
		sm[winner].Wins++
		sm[winner].CurrentWinStreak++
		if sm[winner].CurrentWinStreak > sm[winner].HighestWinStreak {
			sm[winner].HighestWinStreak = sm[winner].CurrentWinStreak
		}

		sm[loser].Losses++
		sm[loser].CurrentWinStreak = 0
	}
}

// CalculateComputed calc's computed stats
func (sm StatsMap) CalculateComputed() {
	for name := range sm {
		sm[name].Games = sm[name].Wins + sm[name].Ties + sm[name].Losses
		sm[name].WinPercentage = float64(sm[name].Wins) / float64(sm[name].Games)
	}
}

// StatsFromMatches ...
func StatsFromMatches(matches []Match) (SortedStats, error) {
	players := make(StatsMap)

	for _, match := range matches {
		// Set stats to 0 for new players
		players.Init(match.Winner, match.Loser)
		// Increase stats
		players.Increase(match.Winner, match.Loser, match.IsTie)
	}

	players.CalculateComputed()

	// return SortedStats
	ss := make(SortedStats, 0, len(players))

	for key, p := range players {
		p.Name = key
		ss = append(ss, p)
	}

	return ss, nil
}
