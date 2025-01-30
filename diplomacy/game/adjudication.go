package game

import (
	"fmt"
	"maps"

	"github.com/adambyle/diplopad/diplomacy/geo"
	"github.com/adambyle/diplopad/diplomacy/orders"
)

func (g *Game) Submit(ords map[geo.Nation]orders.Orders) (
	result *Game,
	outcomes map[geo.Nation]Outcomes,
) {
	return
}

func (g *Game) SubmitBuilds(ords map[geo.Nation][]orders.Build) (
	result *Game,
	outcomes map[geo.Nation]map[orders.Build]Outcome,
) {
	result = &Game{
		board:   g.board,
		year:    g.year + 1,
		phase:   Spring,
		units:   maps.Clone(g.units),
		centers: maps.Clone(g.centers),
	}
	for nation, ords := range ords {
		buildsLeft := g.CenterCount(nation) - g.UnitCount(nation)
		outcomes[nation] = make(map[orders.Build]Outcome)
		for _, ord := range ords {
			var (
				fail = func(r string) {
					outcomes[nation][ord] = failure(r)
				}
				succeed = func() {
					outcomes[nation][ord] = success()
				}
			)
			if ord.Disband {
				if buildsLeft >= 0 {
					fail("No disbands needed.")
					continue
				}
				if ord.Target == nil {
					fail("Unit to disband unspecified.")
					continue
				}
				place, ok := g.units[ord.Target]
				if !ok {
					fail(fmt.Sprintf("No unit to disband in %v.", ord.Target))
					continue
				}
				if place.Nation != nation {
					fail(fmt.Sprintf("Unit in %v belongs to %v.", ord.Target, place.Nation))
					continue
				}
				delete(result.units, ord.Target)
				buildsLeft++
				succeed()
			} else {
				if buildsLeft <= 0 {
					fail("No builds left.")
					continue
				}
				if ord.Target == nil {
					fail("Province to build in unspecified.")
					continue
				}
				if _, ok := g.units[ord.Target]; ok {
					fail(fmt.Sprintf("%v occupied.", ord.Target))
					continue
				}
				if !ord.Target.Terrain().Occupiable(ord.Unit) {
					fail(fmt.Sprintf("Unit %v cannot occupy %v.", ord.Unit, ord.Target))
					continue
				}
				// Coast handling: coast specification in order is ignored
				// if target province has no specific coast or is not coastal.
				var coast geo.Coast
				if _, ok := ord.Target.Coasts(); ok {
					coast = ord.Coast
					if !ord.Target.HasCoast(coast) {
						fail(fmt.Sprintf("Invalid coast for %v.", ord.Target))
						continue
					}
				} else {
					coast = geo.UnnamedCoast
				}
				result.setUnit(ord.Target, Placement{ord.Unit, nation, coast})
				buildsLeft--
				succeed()
			}
		}
		// TODO handle remaining negative builds.
	}
	return
}

func (g *Game) SubmitRetreats(ords map[geo.Nation][]orders.Retreat) (
	result *Game,
	outcomes map[geo.Nation]map[orders.Retreat]Outcome,
) {
	result = &Game{
		board:   g.board,
		year:    g.year,
		phase:   g.phase + 1,
		units:   maps.Clone(g.units),
		centers: maps.Clone(g.centers),
	}
	// TODO handle retreat orders to same location.
	for nation, ords := range ords {
		outcomes[nation] = make(map[orders.Retreat]Outcome)
		for _, ord := range ords {
			var (
				fail = func(r string) {
					outcomes[nation][ord] = failure(r)
				}
				succeed = func() {
					outcomes[nation][ord] = success()
				}
			)
			if ord.Unit == nil {
				fail("Unit to retreat or disband unspecified.")
				continue
			}
			disband := ord.Dest == nil
			ret, ok := g.retreats[ord.Unit]
			if !ok || ret.standoff {
				fail(fmt.Sprintf("No unit in %v was dislodged.", ord.Unit))
				continue
			}
			if ret.nation != nation {
				fail(fmt.Sprintf("Unit dislodged from %v belongs to %v.", ord.Unit, ret.nation))
				continue
			}
			if disband {
				// Unit is already removed, just not saved from "retreat limbo".
				succeed()
			} else {
				if _, ok := g.units[ord.Dest]; ok {
					fail(fmt.Sprintf("%v occupied.", ord.Dest))
					continue
				}
				if ret, ok := g.retreats[ord.Dest]; ok && ret.standoff {
					fail(fmt.Sprintf("A standoff occured in %v.", ord.Dest))
					continue
				}
				// TODO.
				succeed()
			}
		}
	}
	return
}

// Outcome summarizes the result of an order.
type Outcome struct {
	Success bool   // Whether the order had an effect on the game.
	Reason  string // Reason why the order failed.
}

func failure(reason string) Outcome {
	return Outcome{false, reason}
}

func success() Outcome {
	return Outcome{true, "Success."}
}

type Outcomes struct {
	Holds    map[orders.Hold]Outcome
	Moves    map[orders.Move]Outcome
	Supports map[orders.Support]Outcome
	Convoys  map[orders.Convoy]Outcome
}
