// Package diplo allows simulation and analysis of Diplomacy games.
package diplo

import (
	"errors"
	"iter"
	"maps"
	"slices"
	"strings"
)

// StartYear is the year of the game's beginning.
const StartYear = 1901

// Phase is a stage of the game (i.e. moving, retreating, building)
// that determines what kind of orders players must make and what
// happens to units afterward.
type Phase int

const (
	// Spring has players move units.
	Spring Phase = iota
	// SpringRetreats has players retreat and disband dislodged units.
	SpringRetreats
	// Fall has players move units. Supply centers are capturable.
	Fall
	// FallRetreats has players retreat and disband dislodged units.
	// Supply centers are capturable.
	FallRetreats
	// Winter has players build units and disband them as required.
	Winter
)

// Move tells whether it is a move phase ([Spring] or [Fall]).
func (p Phase) Move() bool {
	return p == Spring || p == Fall
}

// Retreat tells whether it is a retreat-and-disband phase
// ([SpringRetreats] or [FallRetreats]).
func (p Phase) Retreat() bool {
	return p == SpringRetreats || p == FallRetreats
}

// Occupancy is a unit occupying a space on the board.
type Occupancy struct {
	province *Province
	coast    string
	unit     Unit
	country  string
}

// Province is which space the unit occupies.
func (o *Occupancy) Province() *Province {
	return o.province
}

// Coast is which specific coast the Fleet occupies.
//
// If the province is not coastal, or is coastal but has only
// one continuous coast Fleets may occupy, this returns false.
func (o *Occupancy) Coast() (string, bool) {
	if o.unit == Army || o.coast == "" {
		return "", false
	}
	return o.coast, true
}

// Unit is what kind of unit occupies the space.
func (o *Occupancy) Unit() Unit {
	return o.unit
}

// Country is who the unit belongs to.
func (o *Occupancy) Country() string {
	return o.country
}

// Game is a state of an ongoing game, in between resolving orders.
//
// The zero-value for Game is unusable, since a game board needs to be provided.
type Game struct {
	board   *Board
	year    int
	phase   Phase
	units   map[*Province]*Occupancy
	centers map[*Province]string
	// Used for retreats
	dislodged map[*Province]*Occupancy // dislodged units
	contests  map[*Province]bool       // cannot retreat here
	attackers map[*Occupancy]*Province // cannot retreat to attacker province
}

// NewGame creates a fresh game state from the specified board.
//
// Important: each country will automatically control its home supply centers,
// but the starting units must be added manually. See [StandardGameSetup].
//
// The game starts in [Spring] of [StartYear].
func NewGame(board *Board) *Game {
	game := &Game{
		board:   board,
		year:    StartYear,
		phase:   Spring,
		units:   make(map[*Province]*Occupancy),
		centers: make(map[*Province]string),
	}
	for p := range board.Centers() {
		game.centers[p] = p.country
	}
	game.resetRetreats()
	return game
}

func (g *Game) resetRetreats() {
	g.dislodged = make(map[*Province]*Occupancy)
	g.contests = make(map[*Province]bool)
}

// Board is the geographical layout the game uses.
func (g *Game) Board() *Board {
	return g.board
}

// Year is what round of phases the game is on.
func (g *Game) Year() int {
	return g.year
}

// Phase controls which stage of orders are needed.
func (g *Game) Phase() Phase {
	return g.phase
}

// AllUnits is every unit in a space on the board.
func (g *Game) AllUnits() iter.Seq[*Occupancy] {
	return maps.Values(g.units)
}

// Units gets every unit on the board for a country.
func (g *Game) Units(country string) iter.Seq[*Occupancy] {
	// TODO country parsing.
	return func(yield func(*Occupancy) bool) {
		for o := range g.AllUnits() {
			if strings.EqualFold(o.country, country) {
				if !yield(o) {
					return
				}
			}
		}
	}
}

// Unit gets what is occupying a space, or nil if nothing is there.
func (g *Game) Unit(province *Province) *Occupancy {
	return g.units[province]
}

// UnitCount is how many units a country has on the board.
func (g *Game) UnitCount(country string) int {
	return count(g.Units(country))
}

// AllCenters is every supply center and which country (empty string if none) controls it.
func (g *Game) AllCenters() iter.Seq2[*Province, string] {
	return maps.All(g.centers)
}

// Centers gets all controlled supply centers for a country.
func (g *Game) Centers(country string) iter.Seq[*Province] {
	return func(yield func(*Province) bool) {
		for p, c := range g.AllCenters() {
			if strings.EqualFold(c, country) {
				if !yield(p) {
					return
				}
			}
		}
	}
}

// Center gets which country, if any, controls a supply center.
//
// Will return empty string and false if not a supply center
// or not controlled by any country.
func (g *Game) Center(province *Province) (string, bool) {
	country := g.centers[province]
	if country == "" {
		return "", false
	} else {
		return country, true
	}
}

// CenterCount gets how many supply centers a country controls.
func (g *Game) CenterCount(country string) int {
	return count(g.Centers(country))
}

// OpenHomeCenters gets all of a country's home supply centers
// that they control and are not currently occupied by a unit.
func (g *Game) OpenHomeCenters(country string) iter.Seq[*Province] {
	return func(yield func(*Province) bool) {
		for p := range g.board.HomeCenters(country) {
			if g.centers[p] != country {
				continue
			}
			if _, ok := g.units[p]; ok {
				continue
			}
			if !yield(p) {
				return
			}
		}
	}
}

// OpenHomeCenterCount gets how many of a country's home supply centers
// are controlled by them and not currently occupied by a unit.
func (g *Game) OpenHomeCenterCount(country string) int {
	return count(g.OpenHomeCenters(country))
}

// Contested is which provinces had a standoff in the previous
// [Spring] or [Fall] phase; these may not be retreated to.
func (g *Game) Contested() iter.Seq[*Province] {
	if !g.phase.Retreat() {
		return nil
	}
	return maps.Keys(g.contests)
}

// DislodgedUnit gets a unit that has been dislodged from the province
// for this retreat phase.
func (g *Game) DislodgedUnit(province *Province) *Occupancy {
	return g.dislodged[province]
}

// Dislodged is all the units that were dislodged in a prior
// move phase and need to retreat or disband.
//
// The [Occupancy] values represent how things were in the prior phase;
// the unit must retreat to an adjacent province or disband. It does
// not mean a unit is actually present there for any other purpose.
func (g *Game) AllDislodged() iter.Seq[*Occupancy] {
	if !g.phase.Retreat() {
		return nil
	}
	return maps.Values(g.dislodged)
}

// CenterDistance finds the distance from a province to the nearest
// supply center controlled by a certain country. -1 is there are no centers
// or the province doesn't exist.
//
// This value is used when a country is in civil disorder. When units
// must be disbanded, the furthest from any controlled supply center are first.
func (g *Game) CenterDistance(province *Province, country string) int {
	if province == nil {
		return -1
	}
	distance := 0
	var (
		nodes   []*Province
		next    = []*Province{province}
		visited = map[*Province]bool{province: true}
	)
	for len(next) > 0 {
		nodes, next = next, nil
		for _, n := range nodes {
			if strings.EqualFold(g.centers[n], country) {
				return distance
			}
			for c := range g.board.ConnectionsFrom(n) {
				if visited[c.to] {
					continue
				}
				visited[c.to] = true
				next = append(next, c.to)
			}
		}
		distance += 1
	}
	return -1
}

// FarthestUnits gets a country's units in order of furthest from a controlled supply center.
func (g *Game) FarthestUnits(country string) []*Occupancy {
	var units []*Occupancy
	distance := make(map[*Occupancy]int)
	for u := range g.Units(country) {
		d := g.CenterDistance(u.province, country)
		units = append(units, u)
		distance[u] = d
	}
	slices.SortFunc(units, func(a, b *Occupancy) int {
		if cmp := distance[b] - distance[a]; cmp != 0 {
			return cmp
		} else {
			return strings.Compare(a.province.name, b.province.name)
		}
	})
	return units
}

func (g *Game) convoyChains(
	chains [][]*Province,
	base []*Province,
	dest *Province,
) [][]*Province {
	baseLength := len(base)
	last := base[baseLength-1]
	for c := range g.board.ConnectionsFrom(last) {
		to := c.to
		if to == dest {
			return append(chains, base[1:])
		}
		if to.terrain != Water {
			continue
		}
		if o, ok := g.units[to]; !ok || o.unit != Fleet {
			continue
		}
		if slices.Contains(base, to) {
			continue
		}
		next := make([]*Province, baseLength+1)
		copy(base, next[:baseLength])
		next[baseLength] = to
		chains = g.convoyChains(chains, next, dest)
	}
	return chains
}

// ConvoyChains finds all valid sequences of water provinces
// occupied by fleets from one coastal province to a
// destination coastal province.
func (g *Game) ConvoyChains(from, to *Province) [][]*Province {
	if from == nil || to == nil ||
		from.terrain != Coastal || to.terrain != Coastal {
		return nil
	}
	return g.convoyChains(nil, []*Province{from}, to)
}

// HasDestination determines whether a unit can legally travel to the
// specified province.
//
// It requires the province to be adjacent or for there to be a legal convoy route.
func (g *Game) HasDestination(unit *Occupancy, destination *Province) bool {
	if unit == nil || destination == nil {
		return false
	}
	if c := g.board.Connection(unit.province, destination); c != nil {
		return c.Traversable(unit.unit)
	}
	// Convoy routes.
	if unit.province.terrain != Coastal || destination.terrain != Coastal {
		return false
	}
	var (
		nodes   []*Province
		next    = []*Province{unit.province}
		visited = map[*Province]bool{unit.province: true}
	)
	for len(next) > 0 {
		nodes, next = next, nil
		for _, n := range nodes {
			for c := range g.board.ConnectionsFrom(n) {
				to := c.to
				if to == destination {
					return true
				}
				if visited[to] || to.terrain != Water {
					continue
				}
				visited[to] = true
				// Fleet must be there, to perform convoy.
				if g.Unit(to) == nil {
					continue
				}
				nodes = append(nodes, to)
				continue
			}
		}
	}
	return false
}

// Destinations gets which provinces a unit can travel to, including
// across oceans for Armies.
//
// [Game.Neighbors] is a subset of this.
func (g *Game) Destinations(unit *Occupancy) iter.Seq[*Province] {
	if unit == nil {
		return nil
	}
	return func(yield func(*Province) bool) {
		for p := range g.Neighbors(unit) {
			if !yield(p) {
				return
			}
		}
		// Follow Fleets to potential convoy destinations.
		if unit.province.terrain != Coastal {
			return
		}
		var (
			nodes   []*Province
			next    = []*Province{unit.province}
			visited = map[*Province]bool{unit.province: true}
		)
		for len(next) > 0 {
			nodes, next = next, nil
			for _, n := range nodes {
				for c := range g.board.ConnectionsFrom(n) {
					to := c.to
					if visited[to] {
						continue
					}
					visited[to] = true
					if to.terrain == Water {
						// Fleet must be there, to perform convoy.
						if g.Unit(to) == nil {
							continue
						}
						nodes = append(nodes, to)
						continue
					}
					// Coastal endpoint found.
					if !yield(to) {
						return
					}
				}
			}
		}
	}
}

// HasNeighbor determines whether a unit can travel to the adjancent destination.
func (g *Game) HasNeighbor(unit *Occupancy, destination *Province) bool {
	c := g.board.Connection(unit.province, destination)
	if c == nil {
		return false
	}
	return c.Traversable(unit.unit)
}

// Neighbors gets which adjacent provinces a unit can travel to.
func (g *Game) Neighbors(unit *Occupancy) iter.Seq[*Province] {
	if unit == nil {
		return nil
	}
	return func(yield func(*Province) bool) {
		for c := range g.board.ConnectionsFrom(unit.province) {
			if !c.Traversable(unit.unit) {
				continue
			}
			if !yield(c.to) {
				return
			}
		}
	}
}

func (g *Game) validSetUnit(p *Province, cs string, u Unit, cn string) (*Occupancy, error) {
	if err := g.board.validProvince(p); err != nil {
		return nil, err
	}
	if err := g.board.validCountry(cn); err != nil {
		return nil, err
	}
	if p.terrain == Coastal && u == Fleet {
		if err := p.validCoast(cs); err != nil {
			return nil, err
		}
	} else {
		cs = ""
	}
	if !p.terrain.Supports(u) {
		return nil, errors.New("unit cannot occupy terrain")
	}
	return &Occupancy{p, cs, u, cn}, nil
}

// SetUnit puts a unit on the board in the given space, removing any existing one.
//
// Coast must be set only if the unit is a Fleet and the province is coastal, having
// more than one distinct named coast that can be occupied. It is blank otherwise.
func (g *Game) SetUnit(province *Province, coast string, unit Unit, country string) error {
	occ, err := g.validSetUnit(province, coast, unit, country)
	if err != nil {
		return nil
	}
	g.units[province] = occ
	return nil
}

// RemoveUnit vacates a space on the game board.
func (g *Game) RemoveUnit(province *Province) {
	delete(g.units, province)
}

// TakeCenter gives control of a supply center to a country.
func (g *Game) TakeCenter(center *Province, country string) error {
	if err := g.board.validCenter(center); err != nil {
		return err
	}
	if err := g.board.validCountry(country); err != nil {
		return err
	}
	g.centers[center] = country
	return nil
}

// FreeCenter removes a supply center from control by any country.
func (g *Game) FreeCenter(center *Province) error {
	if err := g.board.validCenter(center); err != nil {
		return err
	}
	g.centers[center] = ""
	return nil
}

// SetYear changes the active year in the game.
func (g *Game) SetYear(year int) {
	g.year = max(year, StartYear)
}

// SetPhase changes the current phase of the game, without adjudicating orders.
//
// Unit counts are not enforced, i.e. when changing to or from [Winter].
// Retreating units are automatically disbanded; they disappear from the game.
func (g *Game) SetPhase(phase Phase) {
	g.resetRetreats()
	g.phase = phase
}

// BlockRetreat prevents a space from being retreated to during this retreat phase,
// as if a standoff had occured there.
func (g *Game) BlockRetreat(province *Province) {
	if !g.phase.Retreat() {
		return
	}
	g.contests[province] = true
}

// UnblockRetreat allows a space to be retreated to during this retreat phase,
// as if no standoff had occured there.
//
// The space still cannot be retreated to if it is occupied by a unit.
func (g *Game) UnblockRetreat(province *Province) {
	if !g.phase.Retreat() {
		return
	}
	delete(g.contests, province)
}

// AddDislodged creates a unit that must retreat from the given space or disband
// during this retreat phase.
//
// Creating a dislodged unit must follow the same rules as placing a unit there.
func (g *Game) AddDislodged(province *Province, coast string, unit Unit, country string, from *Province) error {
	if !g.phase.Retreat() {
		return errors.New("game not in retreat phase")
	}
	occ, err := g.validSetUnit(province, coast, unit, country)
	if err != nil {
		return err
	}
	g.dislodged[province] = occ
	g.attackers[occ] = from
	return nil
}
