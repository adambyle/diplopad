package diplo

type Arena struct {
	game    *Game
	defense map[*Occupancy]map[*Occupancy]bool // set of defenders per unit
	offense map[*Province]map[*Occupancy]bool  // set of offenders per space
}
