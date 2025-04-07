package diplo

import (
	"errors"
	"fmt"
	"iter"
	"maps"
	"slices"
	"strings"
)

func DefaultCoastParser(coast string) (string, bool) {
	coast = strings.TrimSpace(strings.ToLower(coast))
	switch coast {
	case "east", "ec", "e":
		return "EC", true
	case "north", "nc", "n":
		return "NC", true
	case "south", "sc", "s":
		return "SC", true
	case "west", "wc", "w":
		return "WC", true
	default:
		return "", false
	}
}

func DefaultCountryParser(country string) (string, bool) {
	country = strings.TrimSpace(strings.ToLower(country))
	switch country {
	case "a", "aus", "austria", "hungary", "austria-hungary":
		return "Austria", true
	case "e", "eng", "england":
		return "England", true
	case "f", "fra", "france":
		return "France", true
	case "g", "ger", "germany":
		return "Germany", true
	case "i", "ita", "italy":
		return "Italy", true
	case "r", "rus", "russia":
		return "Russia", true
	case "t", "tur", "turkey":
		return "Turkey", true
	default:
		return "", false
	}
}

// Terrain is a property of a province controlling
// which units may occupy it.
type Terrain string

const (
	// Inland provinces may be occupied by Armies.
	Inland Terrain = "Inland"
	// Coastal provinces may be occupied by Armies or Fleets.
	Coastal Terrain = "Coastal"
	// Water provinces may be occupied by Fleets.
	Water Terrain = "Water"
)

// Supports tells whether a unit type can occupy the terrain.
func (t Terrain) Supports(u Unit) bool {
	var (
		badArmy  = u == Army && t == Water
		badFleet = u == Fleet && t == Inland
	)
	return !badArmy && !badFleet
}

// Province is a space is on the game board that a unit can occupy.
type Province struct {
	name    string   // Full name
	abbrs   []string // Unique abbreviations
	terrain Terrain
	coasts  []string // Named coasts, ignored if not coastal
	center  bool     // Is supply center
	country string   // Supply center home
}

func (p *Province) validCoast(coast string) error {
	if !hasStringFold(p.coasts, coast) {
		return fmt.Errorf("province %s has no coast %s", p.name, coast)
	}
	return nil
}

// Name is the full name of the province.
func (p *Province) Name() string {
	return p.name
}

// Abbreviations is all valid, unique abbreviations for the province.
func (p *Province) Abbreviations() []string {
	return slices.Clone(p.abbrs)
}

// Terrain controls which units can occupy the province.
func (p *Province) Terrain() Terrain {
	return p.terrain
}

// Coasts is the names of the distinct coasts the province has.
// This is empty for non-coastal provinces and for coastal provinces
// that have one continuous coast occupiable by Fleets.
func (p *Province) Coasts() []string {
	if p.terrain != Coastal {
		return nil
	}
	return p.coasts
}

// Center tells whether the province is a supply center.
func (p *Province) Center() bool {
	return p.center
}

// Country tells which country the province is a home supply center for, if any.
// Returns false if not a home supply center.
func (p *Province) Country() (string, bool) {
	if p.center && p.country != "" {
		return p.country, true
	} else {
		return "", false
	}
}

type endpoints struct {
	from, to *Province
}

// Connection is an adjacency between two provinces.
//
// Connections are symmetrical, so the From and To methods
// are only meaningful in certain contexts (such as when outbound
// connections are requested from a certain province).
type Connection struct {
	from, to             *Province
	fromCoasts, toCoasts []string
	coastal              bool
}

// From is the start province in a connection.
func (c *Connection) From() *Province {
	return c.from
}

// To is the destination province in a connection.
func (c *Connection) To() *Province {
	return c.to
}

// FromCoasts are the coasts on the start province
// that a Fleet may travel from to get to the destination province.
func (c *Connection) FromCoasts() []string {
	return slices.Clone(c.fromCoasts)
}

// ToCoasts are the coasts on the destination province
// that a Fleet may travel to from the start province.
func (c *Connection) ToCoasts() []string {
	return slices.Clone(c.toCoasts)
}

// Coastal tells whether the provinces are connected
// along a shared coast (so Fleets may move between them).
func (c *Connection) Coastal() bool {
	return c.coastal
}

// Traversable tells whether a unit kind can cross this connection.
func (c *Connection) Traversable(unit Unit) bool {
	// Unit must be able to occupy provinces.
	if !c.from.terrain.Supports(unit) || !c.to.terrain.Supports(unit) {
		return false
	}
	if unit == Fleet {
		// Water provinces are fleet-traversable; otherwise the provinces are
		// coastal and must be coastally connected.
		return c.from.terrain == Water || c.to.terrain == Water || c.coastal
	}
	return true
}

// Reverse flips the "from" and "to" provinces, returning
// an equivalent connection in the opposite direction.
func (c *Connection) Reverse() *Connection {
	return &Connection{
		from:       c.to,
		to:         c.from,
		fromCoasts: c.toCoasts,
		toCoasts:   c.fromCoasts,
		coastal:    c.coastal,
	}
}

// Board is a game map with countries, provinces, and the connections
// between provinces.
type Board struct {
	countries     []string
	provinces     []*Province
	connections   map[endpoints]*Connection
	coastParser   func(string) (string, bool)
	countryParser func(string) (string, bool)
}

func (b *Board) validCountry(country string) error {
	if !hasStringFold(b.countries, country) {
		return fmt.Errorf("board does not have country '%s'", country)
	}
	return nil
}

func (b *Board) validCenter(p *Province) error {
	if err := b.validProvince(p); err != nil {
		return err
	}
	if !p.center {
		return fmt.Errorf("province %s is not a supply center", p.name)
	}
	return nil
}

func (b *Board) validProvince(p *Province) error {
	if p == nil {
		return errors.New("nil province")
	}
	if !slices.Contains(b.provinces, p) {
		return fmt.Errorf("board does not have province '%s'", p.name)
	}
	return nil
}

// Countries are the Great Powers active on this board.
func (b *Board) Countries() []string {
	return slices.Clone(b.countries)
}

// Provinces is all provinces on the board.
func (b *Board) Provinces() iter.Seq[*Province] {
	return slices.Values(b.provinces)
}

// Province gets the province on the board with the given name.
// Returns nil if it doesn't exist.
func (b *Board) Province(name string) *Province {
	for _, p := range b.provinces {
		if strings.EqualFold(p.name, name) {
			return p
		}
	}
	return nil
}

func (b *Board) ParseProvince(id string) []*Province {
	id = simplify(id)
	var results []*Province
	for _, p := range b.provinces {
		if hasStringFold(p.abbrs, id) {
			results = append(results, p)
		}
	}
	if len(results) > 0 {
		return results
	}
	for _, p := range b.provinces {
		if strings.HasPrefix(simplify(p.name), id) {
			results = append(results, p)
		}
	}
	return results
}

// Centers is all supply centers on the board.
func (b *Board) Centers() iter.Seq[*Province] {
	return func(yield func(*Province) bool) {
		for _, p := range b.provinces {
			if p.center {
				if !yield(p) {
					return
				}
			}
		}
	}
}

// AllHomeCenters is all home supply centers for every country.
func (b *Board) AllHomeCenters() iter.Seq[*Province] {
	return func(yield func(*Province) bool) {
		for p := range b.Centers() {
			if p.country != "" {
				if !yield(p) {
					return
				}
			}
		}
	}
}

// HomeCenters gets the home supply centers for a country.
func (b *Board) HomeCenters(country string) iter.Seq[*Province] {
	return func(yield func(*Province) bool) {
		for p := range b.Centers() {
			if p.country == country {
				if !yield(p) {
					return
				}
			}
		}
	}
}

// Connections is all adjacencies between provinces on the board.
//
// Duplicates are not present; all connections will be in one direction, and the
// reverse counterpart will not be present.
func (b *Board) Connections() iter.Seq[*Connection] {
	return maps.Values(b.connections)
}

// Connection gets the adjacency between two provinces.
// Returns nil if no connection exists.
func (b *Board) Connection(from, to *Province) *Connection {
	if c, ok := b.connections[endpoints{from, to}]; ok {
		return c
	}
	if c, ok := b.connections[endpoints{to, from}]; ok {
		return c.Reverse()
	}
	return nil
}

// Connects tests whether one province connects to another through
// the given coasts.
//
// Use [Board.Connects] instead to validate Army movements.
//
// fromCoast and toCoast may be the empty string when no distinct
// coasts exist on the respective province, or when only one would be valid
// anyway.
func (b *Board) Connects(from, to *Province, fromCoast, toCoast string) bool {
	c := b.Connection(from, to)
	if c == nil {
		return false
	}
	var (
		// If to-coast unspecified, accept if only one coast valid.
		coerceTo = len(c.toCoasts) == 1 && toCoast == ""
		// From-coast and to-coast are acceptable if...
		fromValid = len(c.fromCoasts) == 0 || hasStringFold(c.fromCoasts, fromCoast)
		toValid   = coerceTo || len(c.toCoasts) == 0 || hasStringFold(c.toCoasts, toCoast)
	)
	return fromValid && toValid
}

// ConnectionsFrom gets all outbound connections from the province.
func (b *Board) ConnectionsFrom(province *Province) iter.Seq[*Connection] {
	return func(yield func(*Connection) bool) {
		for e, c := range b.connections {
			if e.from == province {
				if !yield(c) {
					return
				}
				continue
			}
			if e.to == province {
				if !yield(c.Reverse()) {
					return
				}
			}
		}
	}
}

// ConnectionsTo gets all inbound connections to the province.
func (b *Board) ConnectionsTo(province *Province) iter.Seq[*Connection] {
	return func(yield func(*Connection) bool) {
		for e, c := range b.connections {
			if e.to == province {
				if !yield(c) {
					return
				}
				continue
			}
			if e.from == province {
				if !yield(c.Reverse()) {
					return
				}
			}
		}
	}
}

func (b *Board) ParseCoast(coast string) (string, bool) {
	if p := b.coastParser; p != nil {
		return p(coast)
	} else {
		return DefaultCoastParser(coast)
	}
}

func (b *Board) ParseCountry(country string) (string, bool) {
	if p := b.countryParser; p != nil {
		return p(country)
	} else {
		return DefaultCountryParser(country)
	}
}
