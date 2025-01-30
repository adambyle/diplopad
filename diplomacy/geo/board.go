package geo

import (
	"iter"
	"slices"
	"strings"

	"github.com/adambyle/diplopad/diplomacy/unit"
)

// Coasts describes how two provinces may be connected, specifically
// with regards to start and destination coasts.
//
// For non-coastal provinces, both are UnnamedCoasts (zero-value).
// For coastal provinces with no specific coasts, UnnamedCoasts is used.
// BothCoasts is lega in the context of a map connection, denoting that either
// of the province's specific coasts may be an endpoint of the connection.
type Coasts struct {
	From, To Coast
}

type connection struct {
	from, to *Province
}

func (c Coasts) reverse() Coasts {
	return Coasts{c.To, c.From}
}

type Board struct {
	provinces []*Province
	connects  map[connection]Coasts
}

func NewBoard() *Board {
	m := &Board{
		connects: make(map[connection]Coasts),
	}
	return m
}

func (m *Board) connect(from, to *Province, c Coasts) {
	if !from.hasNode(c.From) || !to.hasNode(c.To) {
		panic("invalid connection coasts")
	}
	if from.terrain == Inland && to.terrain == Water ||
		from.terrain == Water && to.terrain == Inland {

		panic("inland province cannot connect to water province")
	}
	m.connects[connection{from, to}] = c
	m.connects[connection{to, from}] = c.reverse()
}

func (m *Board) connectAll(connections []struct {
	from, to *Province
	c        Coasts
}) {
	for _, c := range connections {
		m.connect(c.from, c.to, c.c)
	}
}

// Provinces returns a sequence of the provinces on the board.
func (m *Board) Provinces() iter.Seq[*Province] {
	return slices.Values(m.provinces)
}

// Province finds information about a province by name (case-insensitive).
// Returns nil if not found.
func (m *Board) Province(name string) *Province {
	for _, p := range m.provinces {
		if strings.EqualFold(p.name, name) {
			return p
		}
	}
	return nil
}

// ParseProvince finds province(s) that match the given abbreviation.
//
// If no abbreviation matches are found, searches by start of name.
func (m *Board) ParseProvince(id string) []*Province {
	var results []*Province
	for _, p := range m.provinces {
		if slices.ContainsFunc(p.abbrs, func(a string) bool {
			return strings.EqualFold(a, id)
		}) {
			results = append(results, p)
		}
	}
	if len(results) > 0 {
		return results
	}
	// We check names only if no results have been found.
	for _, p := range m.provinces {
		if strings.EqualFold(p.name[:len(id)], id) {
			results = append(results, p)
		}
	}
	return results
}

// Connections retrieves all outgoing connections from a province.
//
// Results include pairs of destination province and coast connection info.
func (m *Board) Connections(from *Province) iter.Seq2[*Province, Coasts] {
	return func(yield func(*Province, Coasts) bool) {
		for provs, c := range m.connects {
			if provs.from != from {
				continue
			}
			if !yield(provs.to, c) {
				return
			}
		}
	}
}

// CanMove checks whether a move by a unit between two provinces is valid.
//
// Coasts must be specified for fleet movements. Water provinces always have UnnamedCoast.
// BothCoasts is invalid anywhere; must name a specific coast if applicable.
func (m *Board) CanMove(u unit.Unit, from, to *Province, c Coasts) bool {
	coasts, ok := m.connects[connection{from, to}]
	if !ok || !to.terrain.Occupiable(u) {
		return false
	}
	if u == unit.Army {
		// Armies don't need to specify coast.
		return true
	}
	if c.From == BothCoasts || c.To == BothCoasts {
		// It is nonsense to move from or to "both coasts".
		return false
	}
	var (
		fromValid = coasts.From == c.From || coasts.From == BothCoasts && from.HasCoast(c.From)
		toValid   = coasts.To == c.To || coasts.To == BothCoasts && to.HasCoast(c.To)
	)
	return fromValid && toValid
}
