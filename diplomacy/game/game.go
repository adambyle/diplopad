package game

import (
	"github.com/adambyle/diplopad/diplomacy/geo"
	"github.com/adambyle/diplopad/diplomacy/unit"
)

const StartYear = 1901

// Phase represents a stage of the game.
type Phase int

const (
	Spring Phase = iota
	SpringRetreats
	Fall
	FallRetreats
	Winter
)

// Placement represents a unit in a province.
//
// For Armies, Coast is meaningless.
// For Fleets, Coast is only meaningful on coastal provinces.
// It will be one of the province's specific coasts or UnnamedCoast otherwise.
type Placement struct {
	Unit   unit.Unit
	Nation geo.Nation
	Coast  geo.Coast
}

type retreat struct {
	standoff bool
	nation   geo.Nation
	unit     unit.Unit
}

// Game represents a state of the game.
type Game struct {
	board    *geo.Board
	year     int
	phase    Phase
	units    map[*geo.Province]Placement
	retreats map[*geo.Province]retreat
	centers  map[*geo.Province]geo.Nation
}

// New creates a fresh game state.
func New() *Game {
	const (
		u    = geo.UnnamedCoast
		a, f = unit.Army, unit.Fleet
	)
	startUnits := []struct {
		p string
		u unit.Unit
		n geo.Nation
		c geo.Coast
	}{
		{"Budapest", f, geo.Austria, u},
		{"Trieste", a, geo.Austria, u},
		{"Vienna", a, geo.Austria, u},
		{"Edinburgh", f, geo.England, u},
		{"London", f, geo.England, u},
		{"Liverpool", a, geo.England, u},
		{"Brest", f, geo.France, u},
		{"Marseilles", a, geo.France, u},
		{"Paris", a, geo.France, u},
		{"Berlin", a, geo.Germany, u},
		{"Kiel", f, geo.Germany, u},
		{"Munich", a, geo.Germany, u},
		{"Naples", f, geo.Italy, u},
		{"Rome", a, geo.Italy, u},
		{"Venice", a, geo.Italy, u},
		{"Moscow", a, geo.Russia, u},
		{"Sevastopol", f, geo.Russia, u},
		{"StPetersburg", f, geo.Russia, geo.SouthCoast},
		{"Warsaw", a, geo.Russia, u},
		{"Ankara", f, geo.Turkey, u},
		{"Constantinople", a, geo.Turkey, u},
		{"Smyrna", a, geo.Turkey, u},
	}
	g := &Game{
		board: geo.NewBoard(),
		year:  StartYear,
		units: make(map[*geo.Province]Placement),
	}
	for _, u := range startUnits {
		p := g.board.Province(u.p)
		g.setUnit(p, Placement{u.u, u.n, u.c})
		g.centers[p] = u.n
	}
	return g
}

func (g *Game) setUnit(prov *geo.Province, place Placement) {
	if !prov.Terrain().Occupiable(place.Unit) {
		panic("unit cannot occupy province")
	}
	if !prov.HasCoast(place.Coast) {
		panic("invalid coast")
	}
	if place.Nation == geo.NoNation {
		panic("nation needed")
	}
	g.units[prov] = place
}

// UnitCount counts the units for the given nation.
func (g *Game) UnitCount(n geo.Nation) int {
	count := 0
	for _, u := range g.units {
		if u.Nation == n {
			count++
		}
	}
	return count
}

// CenterCount counts the supply centers owned by the given nation.
func (g *Game) CenterCount(n geo.Nation) int {
	count := 0
	for p, u := range g.units {
		if c, _ := p.Center(); c && u.Nation == n {
			count++
		}
	}
	return count
}
