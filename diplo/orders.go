package diplo

import (
	"fmt"
	"slices"
	"strings"
)

// OrderKind represents the structure of an [Order] object.
type OrderKind int

const (
	// InvalidOrder signifies a malformed [Order] object. See docs of
	// that type for valid forms.
	InvalidOrder OrderKind = iota
	// HoldDisband keeps a unit in place or disbands it (depending on phase).
	HoldDisband
	// MoveRetreat directs a unit to a neighboring province, or retreats there
	// (depending on phase).
	MoveRetreat
	// SupportHold aids another unit's hold, support, or convoy.
	SupportHold
	// SupportMove aids another unit's move.
	SupportMove
	// Convoy facilitates an Army's move across water.
	Convoy
	// Build is the construction of a new unit.
	Build
)

// Order is an instruction for a unit.
//
// An Order object takes different forms depending on what fields are set.
//
//  1. Hold order: Unit is the holding unit. Functions as a disband
//     during retreat and build phases.
//  2. Move order: Unit is the moving unit; Target is the destination
//     province. Functions as a retreat order during retreat phases.
//  3. Support hold order: Unit is the supporting unit. Recipient
//     is the holding unit being supported.
//  4. Support move order: Unit is the supporting unit. Recipient
//     is the moving unit being supported. Target is the destination province.
//  5. Convoy order: Unit is the Fleet making the convoy. Recipient
//     is the moving Army being convoyed. Target is the destination coastal province.
//  6. Build order: Unit is nil. Target is the supply center to build on.
//     Build is the unit type to add.
type Order struct {
	// Unit is the unit making the order. Must always be set, except for build orders.
	Unit *Province
	// Recipient is a unit being supported; nil if not a support.
	Recipient *Province
	// Target is a province to move to; nil if a hold.
	Target *Province
	// TargetCoast is the coast to move to.
	TargetCoast string
	// Convoy is true if a Fleet is making a convoy order;
	// if true, Other and Target must be set.
	Convoy bool
	// The kind of unit to build. Ignored outside of build phases.
	Build Unit
}

// OrderHoldDisband creates a hold or disband order.
func OrderHoldDisband(unit *Province) Order {
	return Order{
		Unit: unit,
	}
}

// OrderMoveRetreat creates a move or retreat order. Coast empty when not needed.
func OrderMoveRetreat(unit, destination *Province, coast string) Order {
	return Order{
		Unit:        unit,
		Target:      destination,
		TargetCoast: coast,
	}
}

// OrderSupportHold creates a support-hold order.
func OrderSupportHold(supporter, holder *Province) Order {
	return Order{
		Unit:      supporter,
		Recipient: holder,
	}
}

// OrderSupportMove creates a support-move order. Coast empty when not needed.
func OrderSupportMove(supporter, mover, destination *Province, coast string) Order {
	return Order{
		Unit:        supporter,
		Recipient:   mover,
		Target:      destination,
		TargetCoast: coast,
	}
}

// OrderConvoy creates a convoy order.
func OrderConvoy(fleet, army, destination *Province) Order {
	return Order{
		Unit:      fleet,
		Recipient: army,
		Target:    destination,
		Convoy:    true,
	}
}

// OrderBuild creates a build order.
func OrderBuild(province *Province, unit Unit) Order {
	return Order{
		Target: province,
		Build:  unit,
	}
}

// Kind determines the form of the order (see list in docs for [Order])
// based on which fields are set. It does not validate the order.
func (o *Order) Kind() OrderKind {
	var (
		u = o.Unit != nil
		r = o.Recipient != nil
		t = o.Target != nil
		c = o.Convoy
	)
	switch {
	case u && !r && !t && !c:
		return HoldDisband
	case u && !r && t && !c:
		return MoveRetreat
	case u && r && !t && !c:
		return SupportHold
	case u && r && t && !c:
		return SupportMove
	case u && r && t && c:
		return Convoy
	case !u && !r && t && !c:
		return Build
	default:
		return InvalidOrder
	}
}

func (g *Game) validParse(name string, valid []*Province) (*Province, error) {
	ps := g.board.ParseProvince(name)
	if len(ps) == 0 {
		return nil, fmt.Errorf("no province %s", name)
	}
	if len(ps) > 1 {
		var matches []*Province
		for _, p := range ps {
			if slices.Contains(valid, p) {
				matches = append(matches, p)
			}
		}
		if len(matches) == 1 {
			return matches[0], nil
		}
		return nil, fmt.Errorf("ambiguous province %s", name)
	}
	return ps[0], nil
}

// ParseOrder can interpret a simple string representation of a unit order.
//
// If coerce is a string, ParseOrder will resolve ambiguous province abbreviations
// by accepting only the ones which the unit could travel to, using the string value
// as the nation making the order. Coast names are always coerced (during resolution).
//
// ParseOrder cannot parse build or disband orders during a [Winter] phase.
//
// Parse order supports the order format described in the rulebook:
// A/F Unit [S/C Other] - Target
//
// The unit prefix A/F may be omitted, and an arrow (-> or -->) may replace the hyphen.
func (g *Game) ParseOrder(order string, coerce string) (*Order, error) {
	order = strings.ToLower(order)
	// Put space around hyphens and arrows for easier processing.
	// Sometimes coast designations use parentheses.
	order = strings.ReplaceAll(order, "-", " - ")
	order = strings.ReplaceAll(order, "--", " - ")
	order = strings.ReplaceAll(order, "->", " - ")
	order = strings.ReplaceAll(order, "-->", " - ")
	order = strings.ReplaceAll(order, "(", " ")
	order = strings.ReplaceAll(order, ")", " ")
	// Process parts.
	var (
		unit, recipient, target, coast = "", "", "", ""
		convoy                         = false
		mode                           = 0
	)
	for i, p := range strings.Fields(order) {
		// Ignore unit prefix.
		if i == 0 && (p == "a" || p == "f") {
			continue
		}
		if c, ok := g.board.ParseCoast(p); ok {
			// Coast only needed for target.
			if mode == 2 {
				coast = c
			}
			continue
		}
		switch mode {
		case 0:
			if unit == "" {
				unit = p
			} else if p == "s" {
				mode = 1
			} else if p == "c" {
				mode = 1
				convoy = true
			} else {
				unit += " " + p
			}
		case 1:
			if recipient == "" {
				recipient = p
			} else if p == "-" {
				mode = 2
			} else {
				recipient += " " + p
			}
		case 2:
			if target == "" {
				target = p
			} else {
				target += " " + p
			}
		}
	}
	o := &Order{
		TargetCoast: coast,
		Convoy:      convoy,
	}
	var err error
	// Unit is required since the only order which does not require a unit, a build order,
	// cannot be parsed.
	var unitV []*Province
	for u := range g.Units(coerce) {
		unitV = append(unitV, u.province)
	}
	if o.Unit, err = g.validParse(unit, unitV); err != nil {
		return nil, err
	}
	if target != "" {
		var targetV []*Province
		if coerce == "" {
			targetV = nil
		} else if recipient == "" {
			// Move order; target can be across water.
			targetV = slices.Collect(g.Destinations(g.Unit(o.Unit)))
		} else {
			// Support-move order; target must be neighbor.
			targetV = slices.Collect(g.Neighbors(g.Unit(o.Unit)))
		}
		if o.Target, err = g.validParse(target, targetV); err != nil {
			return nil, err
		}
	}
	if recipient != "" {
		var recipientV []*Province
		if coerce == "" {
			recipientV = nil
		} else if target == "" {
			// Support-hold order; recipient must be a neighbor
			// of the supporting unit.
			u := g.Unit(o.Unit)
			if u != nil {
				for n := range g.Neighbors(u) {
					if g.Unit(n) != nil {
						recipientV = append(recipientV, n)
					}
				}
			}
		} else {
			// Support-move order; recipient must have target
			// as a destination.
			u := g.Unit(o.Unit)
			for ru := range g.AllUnits() {
				if ru == u {
					continue
				}
				if g.HasDestination(ru, o.Target) {
					recipientV = append(recipientV, ru.province)
				}
			}
		}
		if o.Recipient, err = g.validParse(recipient, recipientV); err != nil {
			return nil, err
		}
	}
	return o, nil
}
