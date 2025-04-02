package diplo

import (
	"errors"
	"iter"
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
	// OutcomeBadCoast says a fleet tried to move to a coast it couldn't reach.
	OutcomeBadCoast
	// OutcomeCoastAmbiguous says a target coast was needed but not specified
	OutcomeCoastAmbiguous
	// OutcomeNoConvoy says an Army couldn't move since it wasn't convoyed.
	OutcomeNoConvoy
	// OutcomeBadRecipient says a unit cannot support the given unit.
	OutcomeBadRecipient
	// OutcomeMissingRecipient says a unit doesn't exist where support was given.
	OutcomeMissingRecipient
	// OutcomeDislodged says a hold, support, or convoy failed because the unit was dislodged.
	OutcomeDislodged
	// OutcomeCut says a support failed because it was cut by an attack.
	OutcomeCut
	// OutcomeWeak says a move failed due to insufficient power.
	OutcomeWeak
	// OutcomeStandoff says a move failed because other units opposed it.
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

type unitOrder struct {
	order   Order
	outcome Outcome
}

// Arena is an interactive helper for resolving orders.
//
// Orders can be added and queried incrementally.
type Arena struct {
	game          *Game
	countryOrders map[string]map[Order]Outcome
	unitOrders    map[*Occupancy]unitOrder
	builds        map[*Province]build
	buildCount    map[string]int
}

// Arena creates a new interactive set of unit orders.
func (g *Game) Arena() *Arena {
	a := &Arena{
		game:          g,
		countryOrders: make(map[string]map[Order]Outcome),
		unitOrders:    make(map[*Occupancy]unitOrder),
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
	k := order.Kind()
	if k != HoldDisband && k != MoveRetreat {
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
	if k == HoldDisband {
		// Disbands always work.
		return unit, OutcomeSuccess
	}
	if !a.game.HasNeighbor(unit, order.Target) {
		if order.Target.terrain.Supports(unit.unit) {
			return unit, OutcomeBadTarget
		} else {
			return unit, OutcomeBadTerrain
		}
	}
	if a.game.standoffs[order.Target] {
		return unit, OutcomeStandoff
	}
	if unit.unit == Fleet {
		cs := a.game.board.Connection(unit.province, order.Target).toCoasts
		tc := order.TargetCoast
		if len(cs) > 1 && tc == "" {
			return unit, OutcomeCoastAmbiguous
		}
		if len(cs) > 0 && tc != "" && !slices.Contains(cs, tc) {
			return unit, OutcomeBadCoast
		}
	}
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
		if !order.Target.terrain.Supports(order.Build) {
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
			a.unitOrders[u] = unitOrder{order, o}
		}
	}
	return o
}

func outcomeAssigned(outcome Outcome) bool {
	switch outcome {
	case OutcomeMalformed, OutcomeRepeatUnit, OutcomeEnemyUnit, OutcomeMissingUnit:
		return false
	default:
		return true
	}
}

func (a *Arena) undo(country string, order Order, outcome Outcome) {
	switch {
	case a.game.phase.Move():
		if !outcomeAssigned(outcome) {
			return
		}
		unit := a.game.Unit(order.Unit)
		delete(a.unitOrders, unit)
	case a.game.phase.Retreat():
		if !outcomeAssigned(outcome) {
			return
		}
		unit := a.game.DislodgedUnit(order.Unit)
		delete(a.unitOrders, unit)
	case a.game.phase == Winter:
		if outcome != OutcomeSuccess {
			return
		}
		if order.Kind() == HoldDisband {
			unit := a.game.Unit(order.Unit)
			delete(a.unitOrders, unit)
			a.buildCount[country]++
		} else {
			delete(a.builds, order.Target)
			a.buildCount[country]--
		}
	}
}

// Query gets the status of a country's order.
//
// If the order does not exist, it gets the outcome of that order
// as if it had been added, without adding it.
//
// If the order already exists, it gets the outcome of that order.
func (a *Arena) Query(country string, order Order) Outcome {
	return a.do(country, order, false)
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

// Unit gets the order given to a certain unit.
func (a *Arena) Unit(unit *Occupancy) (Order, Outcome, bool) {
	if o, ok := a.unitOrders[unit]; ok {
		return o.order, o.outcome, true
	} else {
		return Order{}, 0, false
	}
}

// Unordered is all units that have not been given an order yet.
func (a *Arena) Unordered() iter.Seq[*Occupancy] {
	return func(yield func(*Occupancy) bool) {
		for u := range a.game.AllUnits() {
			if _, ok := a.unitOrders[u]; ok {
				continue
			}
			if !yield(u) {
				return
			}
		}
	}
}

// Retreats gets which units have been given retreat orders this phase.
func (a *Arena) Retreats() iter.Seq2[*Occupancy, Order] {
	if !a.game.phase.Retreat() {
		return nil
	}
	return func(yield func(*Occupancy, Order) bool) {
		for u, o := range a.unitOrders {
			if o.order.Kind() != MoveRetreat {
				continue
			}
			if !yield(u, o.order) {
				return
			}
		}
	}
}

// Disbandments gets which units have been ordered to be disbanded
// in a retreat or build phase.
//
// It does not include automatic disbandments, such as unordered
// units or civil-disorder disbandments in a build phase.
// See [Arena.Unordered] or [Game.FarthestUnits].
func (a *Arena) Disbandments() iter.Seq[*Occupancy] {
	if a.game.phase.Move() {
		return nil
	}
	return func(yield func(*Occupancy) bool) {
		for u, o := range a.unitOrders {
			if o.order.Kind() != HoldDisband {
				continue
			}
			if !yield(u) {
				return
			}
		}
	}
}

// Build gets what kind of unit is being built on a province this phase, if any.
func (a *Arena) Build(province *Province) (unit Unit, coast string, ok bool) {
	if b, ok := a.builds[province]; ok {
		return b.unit, b.coast, true
	} else {
		return 0, "", false
	}
}

// Builds gets all provinces being built in this phase.
func (a *Arena) Builds() iter.Seq[*Province] {
	return maps.Keys(a.builds)
}

// BuildCountLeft is the number of remaining builds a country has.
//
// It is simply the difference between controlled supply centers and existing units;
// it does not factor in how many home supply centers are available.
//
// If it is negative, the magnitude is how many more disbandments must be made this phase.
func (a *Arena) BuildCountLeft(country string) int {
	return a.buildCount[country]
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
	if outcome, ok := a.countryOrders[country][order]; ok {
		a.undo(country, order, outcome)
		delete(a.countryOrders[country], order)
	}
}

// Clear resets the orders given by a certain country.
func (a *Arena) Clear(country string) {
	for order, outcome := range a.countryOrders[country] {
		a.undo(country, order, outcome)
	}
	delete(a.countryOrders, country)
}

// FillIn gives the default orders to unordered units.
//
// It does not handle civil disorder disband conditions for Winter. The default
// order set in a Winter phase has no builds and no disbands.
func (a *Arena) FillIn() {
	if a.game.phase == Winter {
		return
	}
	// Unordered units hold or disband.
	var units iter.Seq[*Occupancy]
	if a.game.phase.Move() {
		units = a.game.AllUnits()
	} else {
		units = a.game.AllDislodged()
	}
	for u := range units {
		if _, ok := a.unitOrders[u]; ok {
			continue
		}
		a.Add(u.country, OrderHoldDisband(u.province))
	}
}

// Go creates a new game state following the adjudication
// of the orders added to the arena.
func (a *Arena) Go() *Game {
	a.FillIn()
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
		for u := range a.game.AllUnits() {
			uo := a.unitOrders[u]
			if uo.outcome == OutcomeDislodged {
				next.AddDislodged(u.province, u.coast, u.unit, u.country)
			} else if uo.order.Kind() == MoveRetreat && uo.outcome == OutcomeSuccess {
				coast := uo.order.TargetCoast
				cs := a.game.board.Connection(u.province, uo.order.Target).toCoasts
				if coast == "" && len(cs) > 0 {
					coast = cs[0]
				}
				if len(cs) == 0 {
					coast = ""
				}
				next.SetUnit(uo.order.Target, coast, u.unit, u.country)
			}
		}
	case a.game.phase.Retreat():
		for u := range a.game.AllDislodged() {
			uo := a.unitOrders[u]
			if uo.order.Kind() != MoveRetreat || uo.outcome != OutcomeSuccess {
				// Order failed or unit deliberately disbanded.
				continue
			}
			coast := uo.order.TargetCoast
			cs := a.game.board.Connection(u.province, uo.order.Target).toCoasts
			if coast == "" && len(cs) > 0 {
				coast = cs[0]
			}
			if len(cs) == 0 {
				coast = ""
			}
			next.SetUnit(uo.order.Target, coast, u.unit, u.country)
		}
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
