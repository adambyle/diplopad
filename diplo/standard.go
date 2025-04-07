package diplo

// StandardBoard is used for a standard game of Diplomacy.
var StandardBoard *Board

func init() {
	standardBuilder := Builder{
		Countries: []string{
			"Austria",
			"England",
			"France",
			"Germany",
			"Italy",
			"Russia",
			"Turkey",
		},
		Provinces: []BuilderProvince{
			{
				Name:          "Adriatic Sea",
				Abbreviations: []string{"ADR"},
				Terrain:       Water,
			},
			{
				Name:          "Aegean Sea",
				Abbreviations: []string{"AEG"},
				Terrain:       Water,
			},
			{
				Name:          "Albania",
				Abbreviations: []string{"Alb"},
				Terrain:       Coastal,
			},
			{
				Name:          "Ankara",
				Abbreviations: []string{"Ank"},
				Terrain:       Coastal,
				Country:       "Turkey",
			},
			{
				Name:          "Apulia",
				Abbreviations: []string{"Apu"},
				Terrain:       Coastal,
			},
			{
				Name:          "Armenia",
				Abbreviations: []string{"Arm"},
				Terrain:       Coastal,
			},
			{
				Name:          "Baltic Sea",
				Abbreviations: []string{"BAL"},
				Terrain:       Water,
			},
			{
				Name:          "Barents Sea",
				Abbreviations: []string{"BAR"},
				Terrain:       Water,
			},
			{
				Name:          "Belgium",
				Abbreviations: []string{"Bel"},
				Terrain:       Coastal,
				Center:        true,
			},
			{
				Name:          "Berlin",
				Abbreviations: []string{"Ber"},
				Terrain:       Coastal,
				Country:       "Germany",
			},
			{
				Name:          "Black Sea",
				Abbreviations: []string{"BLA"},
				Terrain:       "Water",
			},
			{
				Name:          "Bohemia",
				Abbreviations: []string{"Boh"},
				Terrain:       Inland,
			},
			{
				Name:          "Brest",
				Abbreviations: []string{"Bre"},
				Terrain:       Coastal,
				Country:       "France",
			},
			{
				Name:          "Budapest",
				Abbreviations: []string{"Bud"},
				Terrain:       Inland,
				Country:       "Austria",
			},
			{
				Name:          "Bulgaria",
				Abbreviations: []string{"Bul"},
				Terrain:       Coastal,
				Coasts:        []string{"EC", "SC"},
				Center:        true,
			},
			{
				Name:          "Burgundy",
				Abbreviations: []string{"Bur"},
				Terrain:       Inland,
			},
			{
				Name:          "Clyde",
				Abbreviations: []string{"Cly"},
				Terrain:       Coastal,
			},
			{
				Name:          "Constantinople",
				Abbreviations: []string{"Con"},
				Terrain:       Coastal,
				Country:       "Turkey",
			},
			{
				Name:          "Denmark",
				Abbreviations: []string{"Den"},
				Terrain:       Coastal,
				Center:        true,
			},
			{
				Name:          "Eastern Mediterranean",
				Abbreviations: []string{"EAS", "EMS"},
				Terrain:       Water,
			},
			{
				Name:          "Edinburgh",
				Abbreviations: []string{"Edi"},
				Terrain:       Coastal,
				Country:       "England",
			},
			{
				Name:          "English Channel",
				Abbreviations: []string{"ENG"},
				Terrain:       Water,
			},
			{
				Name:          "Finland",
				Abbreviations: []string{"Fin"},
				Terrain:       Coastal,
			},
			{
				Name:          "Galicia",
				Abbreviations: []string{"Gal"},
				Terrain:       Inland,
			},
			{
				Name:          "Gascony",
				Abbreviations: []string{"Gas"},
				Terrain:       Coastal,
			},
			{
				Name:          "Greece",
				Abbreviations: []string{"Gre"},
				Terrain:       Coastal,
				Center:        true,
			},
			{
				Name:          "Gulf of Bothnia",
				Abbreviations: []string{"BOT"},
				Terrain:       Water,
			},
			{
				Name:          "Gulf of Lyon",
				Abbreviations: []string{"LYO", "GOL"},
				Terrain:       Water,
			},
			{
				Name:          "Helgoland Bight",
				Abbreviations: []string{"HEL"},
				Terrain:       Water,
			},
			{
				Name:          "Holland",
				Abbreviations: []string{"Hol"},
				Terrain:       Coastal,
				Center:        true,
			},
			{
				Name:          "Ionian Sea",
				Abbreviations: []string{"ION"},
				Terrain:       Water,
			},
			{
				Name:          "Irish Sea",
				Abbreviations: []string{"IRI"},
				Terrain:       Water,
			},
			{
				Name:          "Kiel",
				Abbreviations: []string{"Kie"},
				Terrain:       Coastal,
				Country:       "Germany",
			},
			{
				Name:          "Liverpool",
				Abbreviations: []string{"Lvp"},
				Terrain:       Coastal,
				Country:       "England",
			},
			{
				Name:          "Livonia",
				Abbreviations: []string{"Lvn"},
				Terrain:       Coastal,
			},
			{
				Name:          "London",
				Abbreviations: []string{"Lon"},
				Terrain:       Coastal,
				Country:       "England",
			},
			{
				Name:          "Marseilles",
				Abbreviations: []string{"Mar"},
				Terrain:       Coastal,
				Country:       "France",
			},
			{
				Name:          "Mid-Atlantic Ocean",
				Abbreviations: []string{"MAO", "MID"},
				Terrain:       Water,
			},
			{
				Name:          "Moscow",
				Abbreviations: []string{"Mos"},
				Terrain:       Inland,
				Country:       "Russia",
			},
			{
				Name:          "Munich",
				Abbreviations: []string{"Mun"},
				Terrain:       Inland,
				Country:       "Germany",
			},
			{
				Name:          "Naples",
				Abbreviations: []string{"Nap"},
				Terrain:       Coastal,
				Country:       "Italy",
			},
			{
				Name:          "North Africa",
				Abbreviations: []string{"NAf"},
				Terrain:       Coastal,
			},
			{
				Name:          "North Atlantic Ocean",
				Abbreviations: []string{"NAO", "NAT"},
				Terrain:       Water,
			},
			{
				Name:          "North Sea",
				Abbreviations: []string{"NTH"},
				Terrain:       Water,
			},
			{
				Name:          "Norway",
				Abbreviations: []string{"Nwy"},
				Terrain:       Coastal,
				Center:        true,
			},
			{
				Name:          "Norwegian Sea",
				Abbreviations: []string{"NWG", "NRG"},
				Terrain:       Water,
			},
			{
				Name:          "Paris",
				Abbreviations: []string{"Par"},
				Terrain:       Inland,
				Country:       "France",
			},
			{
				Name:          "Picardy",
				Abbreviations: []string{"Pic"},
				Terrain:       Coastal,
			},
			{
				Name:          "Piedmont",
				Abbreviations: []string{"Pie"},
				Terrain:       Coastal,
			},
			{
				Name:          "Portugal",
				Abbreviations: []string{"Por"},
				Terrain:       Coastal,
				Center:        true,
			},
			{
				Name:          "Prussia",
				Abbreviations: []string{"Pru"},
				Terrain:       Coastal,
			},
			{
				Name:          "Rome",
				Abbreviations: []string{"Rom"},
				Terrain:       Coastal,
				Country:       "Italy",
			},
			{
				Name:          "Ruhr",
				Abbreviations: []string{"Ruh"},
				Terrain:       Inland,
			},
			{
				Name:          "Rumania",
				Abbreviations: []string{"Rum"},
				Terrain:       Coastal,
				Center:        true,
			},
			{
				Name:          "Serbia",
				Abbreviations: []string{"Ser"},
				Terrain:       Inland,
				Center:        true,
			},
			{
				Name:          "Sevastopol",
				Abbreviations: []string{"Sev"},
				Terrain:       Coastal,
				Country:       "Russia",
			},
			{
				Name:          "Silesia",
				Abbreviations: []string{"Sil"},
				Terrain:       Inland,
			},
			{
				Name:          "Skagerrak",
				Abbreviations: []string{"SKA"},
				Terrain:       Water,
			},
			{
				Name:          "Smyrna",
				Abbreviations: []string{"Smy"},
				Terrain:       Coastal,
				Country:       "Turkey",
			},
			{
				Name:          "Spain",
				Abbreviations: []string{"Spa"},
				Terrain:       Coastal,
				Coasts:        []string{"NC", "SC"},
				Center:        true,
			},
			{
				Name:          "St. Petersburg",
				Abbreviations: []string{"StP"},
				Terrain:       Coastal,
				Coasts:        []string{"NC", "SC"},
				Country:       "Russia",
			},
			{
				Name:          "Sweden",
				Abbreviations: []string{"Swe"},
				Terrain:       Coastal,
				Center:        true,
			},
			{
				Name:          "Syria",
				Abbreviations: []string{"Syr"},
				Terrain:       Coastal,
			},
			{
				Name:          "Trieste",
				Abbreviations: []string{"Tri"},
				Terrain:       Coastal,
				Country:       "Austria",
			},
			{
				Name:          "Tunis",
				Abbreviations: []string{"Tun"},
				Terrain:       Coastal,
				Center:        true,
			},
			{
				Name:          "Tuscany",
				Abbreviations: []string{"Tus"},
				Terrain:       Coastal,
			},
			{
				Name:          "Tyrolia",
				Abbreviations: []string{"Tyr"},
				Terrain:       Inland,
			},
			{
				Name:          "Tyrrhenian Sea",
				Abbreviations: []string{"TYS", "TYN"},
				Terrain:       Water,
			},
			{
				Name:          "Ukraine",
				Abbreviations: []string{"Ukr"},
				Terrain:       Inland,
			},
			{
				Name:          "Venice",
				Abbreviations: []string{"Ven"},
				Terrain:       Coastal,
				Country:       "Italy",
			},
			{
				Name:          "Vienna",
				Abbreviations: []string{"Vie"},
				Terrain:       Inland,
				Country:       "Austria",
			},
			{
				Name:          "Wales",
				Abbreviations: []string{"Wal"},
				Terrain:       Coastal,
			},
			{
				Name:          "Warsaw",
				Abbreviations: []string{"War"},
				Terrain:       Inland,
				Country:       "Russia",
			},
			{
				Name:          "Western Mediterranean",
				Abbreviations: []string{"WES", "WMS"},
				Terrain:       Water,
			},
			{
				Name:          "Yorkshire",
				Abbreviations: []string{"Yor"},
				Terrain:       Coastal,
			},
		},
		Connections: []BuilderConnection{
			{
				From:  "ADR",
				ToAll: []string{"Alb", "Apu", "ION", "Tri", "Ven"},
			},
			{
				From:  "AEG",
				ToAll: []string{"Con", "EAS", "Gre", "ION", "Smy"},
			},
			{
				From:     "AEG",
				To:       "Bul",
				ToCoasts: []string{"SC"},
			},
			{
				From:  "Alb",
				ToAll: []string{"ION", "Ser"},
			},
			{
				From:    "Alb",
				ToAll:   []string{"Gre", "Tri"},
				Coastal: true,
			},
			{
				From:  "Ank",
				ToAll: []string{"BLA", "Smy"},
			},
			{
				From:    "Ank",
				ToAll:   []string{"Con", "Arm"},
				Coastal: true,
			},
			{
				From:  "Apu",
				ToAll: []string{"ION", "Rom"},
			},
			{
				From:    "Apu",
				ToAll:   []string{"Nap", "Ven"},
				Coastal: true,
			},
			{
				From:  "Arm",
				ToAll: []string{"BLA", "Smy", "Syr"},
			},
			{
				From:    "Arm",
				To:      "Sev",
				Coastal: true,
			},
			{
				From:  "BAL",
				ToAll: []string{"Ber", "BOT", "Den", "Kie", "Lvn", "Pru", "Swe"},
			},
			{
				From:  "BAR",
				ToAll: []string{"NWG", "Nwy"},
			},
			{
				From:     "BAR",
				To:       "StP",
				ToCoasts: []string{"NC"},
			},
			{
				From:  "Bel",
				ToAll: []string{"Bur", "ENG", "NTH", "Ruh"},
			},
			{
				From:    "Bel",
				ToAll:   []string{"Hol", "Pic"},
				Coastal: true,
			},
			{
				From:  "Ber",
				ToAll: []string{"Mun", "Sil"},
			},
			{
				From:    "Ber",
				ToAll:   []string{"Kie", "Pru"},
				Coastal: true,
			},
			{
				From:  "BLA",
				ToAll: []string{"Con", "Rum", "Sev"},
			},
			{
				From:     "BLA",
				To:       "Bul",
				ToCoasts: []string{"EC"},
			},
			{
				From:  "Boh",
				ToAll: []string{"Gal", "Mun", "Sil", "Tyr", "Vie"},
			},
			{
				From:  "Bre",
				ToAll: []string{"ENG", "MAO", "Par"},
			},
			{
				From:    "Bre",
				ToAll:   []string{"Gas", "Pic"},
				Coastal: true,
			},
			{
				From:  "Bud",
				ToAll: []string{"Gal", "Rum", "Ser", "Tri", "Vie"},
			},
			{
				From:       "Bul",
				To:         "Con",
				FromCoasts: []string{"EC", "SC"},
				Coastal:    true,
			},
			{
				From:       "Bul",
				To:         "Gre",
				FromCoasts: []string{"SC"},
				Coastal:    true,
			},
			{
				From:       "Bul",
				To:         "Rum",
				FromCoasts: []string{"EC"},
				Coastal:    true,
			},
			{
				From: "Bul",
				To:   "Ser",
			},
			{
				From:  "Bur",
				ToAll: []string{"Gas", "Mar", "Mun", "Par", "Pic", "Ruh"},
			},
			{
				From:  "Cly",
				ToAll: []string{"NAO", "NWG"},
			},
			{
				From:    "Cly",
				ToAll:   []string{"Edi", "Lvp"},
				Coastal: true,
			},
			{
				From:    "Con",
				To:      "Smy",
				Coastal: true,
			},
			{
				From:  "Den",
				ToAll: []string{"HEL", "NTH", "SKA"},
			},
			{
				From:    "Den",
				ToAll:   []string{"Kie", "Swe"},
				Coastal: true,
			},
			{
				From:  "EAS",
				ToAll: []string{"ION", "Smy", "Syr"},
			},
			{
				From:  "Edi",
				ToAll: []string{"Lvp", "NTH", "NWG"},
			},
			{
				From:    "Edi",
				To:      "Yor",
				Coastal: true,
			},
			{
				From:  "ENG",
				ToAll: []string{"IRI", "Lon", "MAO", "NTH", "Pic", "Wal"},
			},
			{
				From:  "Fin",
				ToAll: []string{"Nwy", "BOT"},
			},
			{
				From:    "Fin",
				To:      "Swe",
				Coastal: true,
			},
			{
				From:     "Fin",
				To:       "StP",
				ToCoasts: []string{"SC"},
				Coastal:  true,
			},
			{
				From:  "Gal",
				ToAll: []string{"Rum", "Sil", "Ukr", "Vie", "War"},
			},
			{
				From:  "Gas",
				ToAll: []string{"MAO", "Mar", "Par"},
			},
			{
				From:     "Gas",
				To:       "Spa",
				ToCoasts: []string{"NC"},
				Coastal:  true,
			},
			{
				From:  "Gre",
				ToAll: []string{"ION", "Ser"},
			},
			{
				From:  "BOT",
				ToAll: []string{"Lvn", "Swe"},
			},
			{
				From:     "BOT",
				To:       "StP",
				ToCoasts: []string{"SC"},
			},
			{
				From:  "LYO",
				ToAll: []string{"Mar", "Pie", "Tus", "TYS", "WES"},
			},
			{
				From:     "LYO",
				To:       "Spa",
				ToCoasts: []string{"SC"},
			},
			{
				From:  "HEL",
				ToAll: []string{"Hol", "Kie", "NTH"},
			},
			{
				From:  "Hol",
				ToAll: []string{"NTH", "Ruh"},
			},
			{
				From:    "Hol",
				To:      "Kie",
				Coastal: true,
			},
			{
				From:  "ION",
				ToAll: []string{"Nap", "Tun", "TYS"},
			},
			{
				From:  "IRI",
				ToAll: []string{"Lvp", "MAO", "NAO", "Wal"},
			},
			{
				From:  "Kie",
				ToAll: []string{"Mun", "Ruh"},
			},
			{
				From:  "Lvp",
				ToAll: []string{"NAO", "Yor"},
			},
			{
				From:    "Lvp",
				To:      "Wal",
				Coastal: true,
			},
			{
				From:  "Lvn",
				ToAll: []string{"Mos", "War"},
			},
			{
				From:    "Lvn",
				To:      "Pru",
				Coastal: true,
			},
			{
				From:     "Lvn",
				To:       "StP",
				ToCoasts: []string{"SC"},
				Coastal:  true,
			},
			{
				From: "Lon",
				To:   "NTH",
			},
			{
				From:    "Lon",
				ToAll:   []string{"Wal", "Yor"},
				Coastal: true,
			},
			{
				From:    "Mar",
				To:      "Pie",
				Coastal: true,
			},
			{
				From:     "Mar",
				To:       "Spa",
				ToCoasts: []string{"SC"},
				Coastal:  true,
			},
			{
				From:  "MAO",
				ToAll: []string{"NAf", "NAO", "Por", "WES"},
			},
			{
				From:     "MAO",
				To:       "Spa",
				ToCoasts: []string{"NC", "SC"},
			},
			{
				From:  "Mos",
				ToAll: []string{"Sev", "StP", "Ukr", "War"},
			},
			{
				From:  "Mun",
				ToAll: []string{"Ruh", "Sil", "Tyr"},
			},
			{
				From: "Nap",
				To:   "TYS",
			},
			{
				From:    "Nap",
				To:      "Rom",
				Coastal: true,
			},
			{
				From: "NAf",
				To:   "WES",
			},
			{
				From:    "NAf",
				To:      "Tun",
				Coastal: true,
			},
			{
				From: "NAO",
				To:   "NWG",
			},
			{
				From:  "NTH",
				ToAll: []string{"NWG", "Nwy", "SKA", "Yor"},
			},
			{
				From:  "Nwy",
				ToAll: []string{"NWG", "SKA"},
			},
			{
				From:    "Nwy",
				To:      "Swe",
				Coastal: true,
			},
			{
				From:     "Nwy",
				To:       "StP",
				ToCoasts: []string{"NC"},
				Coastal:  true,
			},
			{
				From: "Par",
				To:   "Pic",
			},
			{
				From:  "Pie",
				ToAll: []string{"Tyr", "Ven"},
			},
			{
				From:    "Pie",
				To:      "Tus",
				Coastal: true,
			},
			{
				From:     "Por",
				To:       "Spa",
				ToCoasts: []string{"NC", "SC"},
				Coastal:  true,
			},
			{
				From:  "Pru",
				ToAll: []string{"Sil", "War"},
			},
			{
				From:  "Rom",
				ToAll: []string{"TYS", "Ven"},
			},
			{
				From:    "Rom",
				To:      "Tus",
				Coastal: true,
			},
			{
				From:  "Rum",
				ToAll: []string{"Ser", "Ukr"},
			},
			{
				From:    "Rum",
				To:      "Sev",
				Coastal: true,
			},
			{
				From: "Ser",
				To:   "Tri",
			},
			{
				From: "Sev",
				To:   "Ukr",
			},
			{
				From: "Sil",
				To:   "War",
			},
			{
				From: "SKA",
				To:   "Swe",
			},
			{
				From:    "Smy",
				To:      "Syr",
				Coastal: true,
			},
			{
				From:       "Spa",
				To:         "WES",
				FromCoasts: []string{"SC"},
			},
			{
				From:  "Tri",
				ToAll: []string{"Tyr", "Vie"},
			},
			{
				From:    "Tri",
				To:      "Ven",
				Coastal: true,
			},
			{
				From:  "Tun",
				ToAll: []string{"TYS", "WES"},
			},
			{
				From:  "Tus",
				ToAll: []string{"TYS", "Ven"},
			},
			{
				From:  "Tyr",
				ToAll: []string{"Ven", "Vie"},
			},
			{
				From: "TYS",
				To:   "WES",
			},
			{
				From: "Ukr",
				To:   "War",
			},
			{
				From: "Wal",
				To:   "Yor",
			},
		},
	}
	b, err := standardBuilder.Build()
	if err != nil {
		panic(err)
	}
	StandardBoard = b
}

// StandardGameSetup sets up the board for a standard game
// by placing the units on the board.
//
// It should be called directly after creating a game object
// with [NewGame].
//
// Prefer to use [StandardGame].
//
// Panics if not used on a game with the standard board.
func StandardGameSetup(g *Game) {
	if g.board != StandardBoard {
		panic("standard board not used")
	}
	g.SetUnit(g.board.Province("London"), "", Fleet, "England")
	// TODO
}

// StandardGame returns the Spring 1901 game state for a
// standard game of Diplomacy.
func StandardGame() *Game {
	g := NewGame(StandardBoard)
	StandardGameSetup(g)
	return g
}
