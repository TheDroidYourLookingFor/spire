package expansions

var kunark = Expansion{
	ExpansionNumber: 1,
	ExpansionName:   "Ruins of Kunark",
	ShortName:       "RoK",
	MaxLevel:        60,
	ContentFlags:    []ContentFlag{},
	Rules: []Rule{
		{
			Name:    "Expansion:CurrentExpansion",
			Value:   "1",
			Comment: "Current Expansion",
		},
		{
			Name:    "World:ExpansionSettings",
			Value:   "1",
			Comment: "Kunark Client-Based Expansion Setting",
		},
		{
			Name:    "World:CharacterSelectExpansionSettings",
			Value:   "1",
			Comment: "Kunark Client-Based Expansion Setting",
		},
		{
			Name:    "Character:MaxExpLevel",
			Value:   "60",
			Comment: "Level 60 cap until PoP",
		},
		{
			Name:    "Character:MaxLevel",
			Value:   "60",
			Comment: "Level 60 cap until PoP",
		},
	},
}
