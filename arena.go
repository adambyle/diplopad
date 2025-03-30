package diplo

import (
	"errors"
	"maps"
	"slices"
)

type Outcome int

const (
	// OutcomeSuccess says the order works.
	OutcomeSuccess Outcome = iota
	// OutcomeMalformed says an order doesn't take a valid form for the phase.
	OutcomeMalformed
	// OutcomeRepeatUnit says the unit has already been given an order.
	OutcomeRepeatUnit
	// OutcomeEnemyUnit says the unit does not belong to the ordering country.
	OutcomeEnemyUnit
	// OutcomeMissingUnit says an order was given to a unit that doesn't exist.
	OutcomeMissingUnit
	// OutcomeBadTerrain says a unit tried to move onto terrain it couldn't occupy.
	OutcomeBadTerrain
	// OutcomeBadTarget says a unit tried to move somewhere it couldn't reach.
	OutcomeBadTarget
	// OutcomeCoastAmbiguous says a target coast was needed but not specified
	OutcomeCoastAmbiguous
	// OutcomeNoConvoy says an Army couldn't move since it wasn't convoyed.
	OutcomeNoConvoy
	// OutcomeBadRecipient says a unit cannot support the given unit.
	OutcomeBadRecipient
	// OutcomeMissingRecipient says a unit doesn't exist where support was given.
	OutcomeMissingRecipient
	// OutcomeDislodged says a support or convoy failed because the unit was dislodged.
	OutcomeDislodged
	// OutcomeCut says a support failed because it was cut by an attack.
	OutcomeCut
	// OutcomeWeak says a move failed due to insufficient power.
	OutcomeWeak
	// OutcomeOpposed says a move failed because other units opposed it.
	OutcomeStandoff
	// OutcomeOverpowered says a move failed because another unit moved in.
	OutcomeOverpowered
	// BUILDS ONLY:
	// OutcomeNoBuilds says a country is out of builds and cannot make this unit.
	OutcomeNoBuilds
	// OutcomeNoDisbands says a country has no excess of units and cannot disband any.
	OutcomeNoDisbands
	// OutcomeNotHome says a unit cannot be built on a supply center not home
	// to the building country.
	OutcomeNotHome
	// OutcomeUnowned says a unit cannot be built on a supply center not controlled
	// by the building country.
	OutcomeNotControlled
	// OutcomeOccupied says a unit cannot be builty where a unit already is.
	OutcomeOccupied
)

type build struct {
	unit  Unit
	coast string
}

// Arena is an interactive helper for resolving orders.
//
// Orders can be added and queried incrementally.
type Arena struct {
	game          *Game
	countryOrders map[string]map[Order]Outcome
	unitOrders    map[*Occupancy]Order
	builds        map[*Province]build
	buildCount    map[string]int
}

// Arena creates a new interactive set of unit orders.
func (g *Game) Arena() *Arena {
	a := &Arena{
		game:          g,
		countryOrders: make(map[string]map[Order]Outcome),
		unitOrders:    make(map[*Occupancy]Order),
		builds:        make(map[*Province]build),
	}
	if g.phase == Winter {
		a.buildCount = make(map[string]int)
		for _, country := range g.board.countries {
			a.buildCount[country] = g.CenterCount(country) - g.UnitCount(country)
		}
	}
	return a
}

// CivilDisorder gets an arena with the default orders for each country.
func (g *Game) CivilDisorder() *Arena {
	a := g.Arena()
	if g.phase == Winter {
		// Disband as needed.
		for _, c := range g.board.countries {
			centerCount := g.CenterCount(c)
			unitCount := g.UnitCount(c)
			if unitCount <= centerCount {
				continue
			}
			// Country has more units than centers; disband that many
			// in order of farthest from center first.
			units := a.game.FarthestUnits(c)
			for i := range unitCount - centerCount {
				a.Add(c, OrderHoldDisband(units[i].province))
			}
		}
	} else {
		// All units hold or disband.
		for u := range g.AllUnits() {
			a.Add(u.country, OrderHoldDisband(u.province))
		}
	}
	return a
}

func (a *Arena) doMovePhase(country string, order Order) (*Occupancy, Outcome) {
	if k := order.Kind(); k == InvalidOrder || k == Build {
		return nil, OutcomeMalformed
	}
	unit := a.game.Unit(order.Unit)
	if unit == nil {
		return nil, OutcomeMissingUnit
	}
	if unit.country != country {
		return nil, OutcomeEnemyUnit
	}
	if _, ok := a.unitOrders[unit]; ok {
		return nil, OutcomeRepeatUnit
	}
	// TODO
	return unit, OutcomeSuccess
}

func (a *Arena) doRetreatPhase(country string, order Order) (*Occupancy, Outcome) {
	if k := order.Kind(); k != HoldDisband && k != MoveRetreat {
		return nil, OutcomeMalformed
	}
	unit := a.game.DislodgedUnit(order.Unit)
	if unit == nil {
		return nil, OutcomeMissingUnit
	}
	if unit.country != country {
		return nil, OutcomeEnemyUnit
	}
	if _, ok := a.unitOrders[unit]; ok {
		return nil, OutcomeRepeatUnit
	}
	// TODO
	return unit, OutcomeSuccess
}

func (a *Arena) doBuildPhase(country string, order Order) (*Occupancy, Outcome) {
	switch order.Kind() {
	case HoldDisband:
		if a.buildCount[country] >= 0 {
			return nil, OutcomeNoDisbands
		}
		unit := a.game.Unit(order.Unit)
		if unit == nil {
			return nil, OutcomeMissingUnit
		}
		if unit.country != country {
			return nil, OutcomeEnemyUnit
		}
		if _, ok := a.unitOrders[unit]; ok {
			return nil, OutcomeRepeatUnit
		}
		return unit, OutcomeSuccess
	case Build:
		if a.buildCount[country] <= 0 {
			return nil, OutcomeNoBuilds
		}
		if order.Target.country != country {
			return nil, OutcomeNotHome
		}
		if c, _ := a.game.Center(order.Target); c != country {
			return nil, OutcomeNotControlled
		}
		if a.game.Unit(order.Target) != nil {
			return nil, OutcomeOccupied
		}
		if order.Target.terrain.Supports(order.Build) != nil {
			return nil, OutcomeBadTerrain
		}
		if c := order.Target.coasts; len(c) > 0 && !slices.Contains(c, order.TargetCoast) {
			return nil, OutcomeCoastAmbiguous
		}
		return nil, OutcomeSuccess
	default:
		return nil, OutcomeMalformed
	}
}

func (a *Arena) do(country string, order Order, add bool) Outcome {
	var (
		o Outcome
		u *Occupancy
	)
	switch {
	case a.game.phase.Move():
		u, o = a.doMovePhase(country, order)
	case a.game.phase.Retreat():
		u, o = a.doRetreatPhase(country, order)
	case a.game.phase == Winter:
		u, o = a.doBuildPhase(country, order)
		if add && o == OutcomeSuccess {
			if order.Kind() == Build {
				a.builds[order.Target] = build{order.Build, order.TargetCoast}
				a.buildCount[country]--
			} else {
				a.buildCount[country]++
			}
		}
	}
	if add {
		a.countryOrders[country][order] = o
		if u != nil {
			a.unitOrders[u] = order
		}
	}
	return o
}

func (a *Arena) undo(country string, order Order) {

}

// Orders is all orders given by a country.
func (a *Arena) Orders(country string) []Order {
	orders := make([]Order, 0, len(a.countryOrders[country]))
	for o := range a.countryOrders[country] {
		orders = append(orders, o)
	}
	return orders
}

// Outcomes gets the outcome of each order a country has given.
func (a *Arena) Outcomes(country string) map[Order]Outcome {
	return maps.Clone(a.countryOrders[country])
}

// Add processes and saves a country's order.
func (a *Arena) Add(country string, order Order) (Outcome, error) {
	if !slices.Contains(a.game.board.countries, country) {
		return 0, errors.New("invalid country")
		// TODO country parsing
	}
	return a.do(country, order, true), nil
}

// Remove undoes a country's order, if it has been made.
func (a *Arena) Remove(country string, order Order) {
	if _, ok := a.countryOrders[country][order]; ok {
		a.undo(country, order)
		delete(a.countryOrders[country], order)
	}
}

// Clear resets the orders given by a certain country.
func (a *Arena) Clear(country string) {
	for order := range a.countryOrders[country] {
		a.undo(country, order)
	}
	delete(a.countryOrders, country)
}

// Query gets the status of a country's order if it was added
// to the arena.
//
// (It does not get the status of an existing order.)
func (a *Arena) Query(country string, order Order) Outcome {
	return a.do(country, order, false)
}

func (a *Arena) Go() *Game {
	next := &Game{
		board:       a.game.board,
		year:        a.game.year,
		phase:       a.game.phase,
		units:       make(map[*Province]*Occupancy),
		centers:     maps.Clone(a.game.centers),
		coastParser: a.game.coastParser,
		dislodged:   make(map[*Province]*Occupancy),
		standoffs:   make(map[*Province]bool),
	}
	// Apply successful orders.
	switch {
	case a.game.phase.Move():
		// TODO
		// Don't forget to update dislodges.
	case a.game.phase.Retreat():
		// TODO
	case a.game.phase == Winter:
		// Civil disorder: disband units.
		cd := make(map[*Occupancy]bool)
		for _, c := range a.game.board.countries {
			bc := a.buildCount[c]
			if bc >= 0 {
				continue
			}
			// Country has disbands unaccounted for; disband farthest-first.
			units := a.game.FarthestUnits(c)
			for i := range -bc {
				cd[units[i]] = true
			}
		}
		// Exclude disbanded units.
		for u := range a.game.AllUnits() {
			// If unit removed by civil disorder...
			if cd[u] {
				continue
			}
			// If unit had a valid order assigned, it was disbanded.
			if _, ok := a.unitOrders[u]; ok {
				continue
			}
			next.SetUnit(u.province, u.coast, u.unit, u.country)
		}
		// Add built units.
		for p, b := range a.builds {
			next.SetUnit(p, b.coast, b.unit, p.country)
		}
	}
	// End of fall: add centers.
	if a.game.phase == FallRetreats {
		for c := range next.centers {
			if u := next.units[c]; u != nil {
				next.centers[c] = u.country
			}
		}
	}
	// TODO skip empty retreat and build phases. Skip() method?
	// Advance phase and year.
	next.phase++
	if next.phase > Winter {
		next.phase = Spring
		next.year++
	}
	return next
}
