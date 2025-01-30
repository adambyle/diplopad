// Package geo handles game geography.
package geo

import (
	"github.com/adambyle/diplopad/diplomacy/unit"
)

// Nation describes one of the seven Great Powers vying for victory.
type Nation string

const (
	NoNation = ""
	Austria  = "Austria"
	England  = "England"
	France   = "France"
	Germany  = "Germany"
	Italy    = "Italy"
	Russia   = "Russia"
	Turkey   = "Turkey"
)

// Nations gets a sequence of all seven nations.
func Nations() []Nation {
	return []Nation{
		Austria,
		England,
		France,
		Germany,
		Italy,
		Russia,
		Turkey,
	}
}

// Terrain describes which units may occupy a territory.
type Terrain byte

const (
	Inland  Terrain = iota // Army occupiable.
	Coastal                // Army and Fleet occupiable.
	Water                  // Fleet occupiable.
)

func (t Terrain) String() string {
	switch t {
	case Inland:
		return "Inland"
	case Water:
		return "Water"
	case Coastal:
		return "Coastal"
	default:
		return ""
	}
}

// Occupiable confirms a unit type can occupy this terrain.
func (t Terrain) Occupiable(u unit.Unit) bool {
	switch t {
	case Inland:
		return u == unit.Army
	case Water:
		return u == unit.Fleet
	default:
		return true
	}
}

// Coast describes one of three directional coasts on some provinces,
// or the "main coast" for provinces with just one coast.
type Coast byte

const (
	UnnamedCoast Coast = iota // Only coast on provinces with just one coast.
	EastCoast
	NorthCoast
	SouthCoast
	BothCoasts // Both of this province's special coasts are a valid destination.
)

// Specific confirms a coast value refers to a distinct coast (north, east, south).
func (c Coast) Specific() bool {
	return c == NorthCoast || c == EastCoast || c == SouthCoast
}

// Province describes a single territory on the board which may be occupied by one unit.
type Province struct {
	name    string
	abbrs   []string
	terrain Terrain
	coasts  [2]Coast
	nation  Nation
	center  bool
}

func newProvince(name string, terrain Terrain, abbrs ...string) *Province {
	if len(abbrs) == 0 {
		panic("province must have at least one abbreviation")
	}
	p := &Province{
		name:    name,
		abbrs:   abbrs,
		terrain: terrain,
	}
	return p
}

func (p *Province) makeCenter() *Province {
	if p.terrain == Water {
		panic("water province cannot be a supply center")
	}
	p.center = true
	return p
}

func (p *Province) makeHomeCenter(nation Nation) *Province {
	p.makeCenter()
	p.nation = nation
	return p
}

func (p *Province) setCoasts(coasts [2]Coast) *Province {
	if p.terrain != Coastal {
		panic("non-coastal province cannot have coasts")
	}
	for _, c := range coasts {
		if c.Specific() {
			panic("coasts must be specific")
		}
	}
	p.coasts = coasts
	return p
}

func (p *Province) hasNode(n Coast) bool {
	return p.HasCoast(n) || p.specificCoasts() && n == BothCoasts
}

func (p *Province) specificCoasts() bool {
	return p.coasts[0] != UnnamedCoast
}

func (p *Province) String() string {
	return p.abbrs[0]
}

// Name returns the unique full name of this province.
func (p *Province) Name() string {
	return p.name
}

// Abbrs returns valid abbreviations for this province, which
// may overlap with the abbreviations for other provinces and
// require disambiguation.
func (p *Province) Abbrs() []string {
	return p.abbrs
}

// Terrain specifies what kind of province this is, which in turn
// determines which units can occupy it.
func (p *Province) Terrain() Terrain {
	return p.terrain
}

// Center specifies whether this is a supply center, and if so,
// what nation it is a home supply center for.
func (p *Province) Center() (center bool, nation Nation) {
	return p.center, p.nation
}

// Coasts lists the province's specific coasts.
//
// For non-coastal provinces and coastal provinces with no specific coasts,
// ok is false.
func (p *Province) Coasts() (coasts [2]Coast, ok bool) {
	if p.terrain != Coastal || !p.specificCoasts() {
		return
	}
	return p.coasts, true
}

// HasCoast checks whether a coast is a valid target
// for movement for this province.
//
// For coastal provinces with more than one coast, this is true
// for either valid specific coast. In all other cases, only true
// for UnnamedCoast.
func (p *Province) HasCoast(c Coast) bool {
	return p.coasts[0] == c || p.coasts[1] == c
}
