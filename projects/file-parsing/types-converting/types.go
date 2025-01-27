package filesConverting

type Score struct {
	Name      string
	HighScore int32 `json:"high_score"`
}

func Want() []Score {
	want := []Score{
		{Name: "Aya", HighScore: 10},
		{Name: "Prisha", HighScore: 30},
		{Name: "Charlie", HighScore: -1},
		{Name: "Margot", HighScore: 25},
	}
	return want
}
