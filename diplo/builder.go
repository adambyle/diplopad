package diplo

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

type BuilderProvince struct {
	// Name is the province's full name, capitailzed and formatted as you want it to be displayed.
	//
	// Names of provinces may contain spaces and punctuation, like "St. Petersburg". The parser will
	// be able to handle inputs that lack that space and punctuation. All provinces should have unique names.
	// Do not include hyphens, which are used to separate unit from destination in order formatting.
	Name string
	// Abbreviations is a list of valid abbreviations unique to the province, capitalized and formatted
	// as you want them to be displayed.
	//
	// There must be at least one. No two provinces may share any abbreviations.
	// Abbreviations should consist only of uppercase and lowercase characters. (Examples: WES, StP, not Wes-Med, not St.P)
	Abbreviations []string
	// Terrain is the province's terrain type.
	Terrain Terrain
	// Coasts is the symbols shown for distinct coasts a coastal province may have.
	//
	// This is ignored for non-coastal provinces. It is an error for only one coast name
	// to be present here. Coasts should only be named if there are multiple, distinct coasts.
	// For example, Portugal has one continuous, unnamed coast, but Spain has two coasts, NC and SC.
	//
	// List the names of the coasts here as you would like them to be displayed in-game; prefer
	// short abbreviations like NC, SC, EC. Be consistent across provinces, including with
	// capitalization. [Builder.CoastParser] should have a way to return each of the coasts
	// specified here. Coast names with spaces in them will be unparseable using the built-in parser;
	// avoid names like "East Coast".
	Coasts []string `json:",omitempty"`
	// Center is true if the province is a supply center. This is automatically set to true
	// when [BuilderProvince.Country] is set.
	Center bool `json:",omitempty"`
	// Country is the name of the country for which the province is a home supply center. When this is
	// set, [BuilderProvince.Center] is set to true.
	Country string `json:",omitempty"`
}

// BuilderConnection connection from one province to another
// (or many others if ToAll is set).
//
// If ToAll is set, the coast fields are ignored.
type BuilderConnection struct {
	// From is one endpoint of the connection.
	From string
	// To is the other endpoint of the connection.
	To string `json:",omitempty"`
	// ToAll is a list of provinces connected to From.
	ToAll []string `json:",omitempty"`
	// FromCoasts is a list of the coasts that are valid on From
	// in this connection (relevant only when from has more than one
	// distinct coast).
	FromCoasts []string `json:",omitempty"`
	// ToCoasts is a list of the coasts that are valid on To
	// in this connection (relevant only when from has more than one
	// distinct coast).
	ToCoasts []string `json:",omitempty"`
	// Coastal is true when From and To are coastal and they are connected
	// along the coast (so that a Fleet may travel between them). When false,
	// if From and To are coastal, only Armies may travel between them.
	Coastal bool `json:",omitempty"`
}

// Builder allows for the construction of custom boards. The zero-value is an empty builder
// ready to use; you can add countries, provinces, and connections directly.
//
// It can be deserialized from JSON.
type Builder struct {
	// Countries is a list of the names of countries on the board, capitalized and formatted
	// as you want them to be displayed in-game.
	Countries []string
	// Provinces are the spaces on the board.
	Provinces []BuilderProvince
	// Connections are the adjancencies between spaces.
	//
	// Connections should only be provided once per pair of provinces.
	Connections []BuilderConnection
	// CoastParser interprets string representations of coast names.
	//
	// If unset, the default coast parser can parse NC, EC, SC, and WC; you will
	// need to implement a custom parser if your map uses other coast names.
	CoastParser func(string) (string, bool) `json:"-"`
	// CountryParser interprets string representations of country names.
	//
	// If unset, the default country parser can parse the names of the standard
	// Diplomacy game's countries; you will need to implement a custom parser
	// if your map uses custom countries.
	CountryParser func(string) (string, bool) `json:"-"`
}

func validEndpoint(p, other *Province, cs []string) error {
	if p == nil {
		return errors.New("nil province")
	}
	if p.terrain == Coastal {
		var (
			pCoasts = len(p.coasts)
			cCoasts = len(cs)
		)
		if pCoasts == 0 && cCoasts > 0 {
			return fmt.Errorf(
				"%s has no named coasts but coasts is non-empty", p.name,
			)
		} else if pCoasts > 0 {
			// Province has named coasts; coasts must be specified in connection.
			if cCoasts == 0 && other.terrain != Inland {
				return fmt.Errorf(
					"%s has named coasts but coasts is empty", p.name,
				)
			}
			// Coasts specified in connection must be valid on province.
			for _, c := range cs {
				if !slices.Contains(p.coasts, c) {
					return fmt.Errorf(
						"no coast %s found on %s", c, p.name,
					)
				}
			}
		}
	} else {
		// Coasts may not be specified for non-coastal province.
		if len(cs) > 0 {
			return fmt.Errorf(
				"%s is not coastal but coasts is non-empty", p.name,
			)
		}
	}
	return nil
}

func (c *Connection) valid() error {
	if c.from == c.to {
		return fmt.Errorf("connection between same provinces: %s", c.from.name)
	}
	if c.from.terrain == Inland && c.to.terrain == Water ||
		c.from.terrain == Water && c.to.terrain == Inland {
		return fmt.Errorf(
			"connection between land and water: %s - %s",
			c.from.name,
			c.to.name,
		)
	}
	if c.coastal && (c.from.terrain != Coastal || c.to.terrain != Coastal) {
		return fmt.Errorf(
			"coastal connection involving non-coastal provinces: %s - %s",
			c.from.name,
			c.to.name,
		)
	}
	// Check each endpoint.
	var err error
	if err = validEndpoint(c.from, c.to, c.fromCoasts); err != nil {
		return err
	}
	if err = validEndpoint(c.to, c.from, c.toCoasts); err != nil {
		return err
	}
	return nil
}

func (b *Builder) Build() (*Board, error) {
	board := &Board{
		coastParser:   b.CoastParser,
		countryParser: b.CountryParser,
		connections:   make(map[endpoints]*Connection),
	}
	for _, c := range b.Countries {
		var ok bool
		c, ok = board.ParseCountry(c)
		if !ok {
			return nil, fmt.Errorf("invalid country %s (need to add custom parser?)", c)
		}
		if !slices.Contains(board.countries, c) {
			board.countries = append(board.countries, c)
		}
	}
	if len(board.countries) == 0 {
		return nil, errors.New("no valid countries")
	}
	for _, p := range b.Provinces {
		name := strings.TrimSpace(p.Name)
		if name == "" {
			return nil, errors.New("empty province name")
		}
		if board.Province(name) != nil {
			return nil, fmt.Errorf("duplicate name %s", name)
		}
		if len(p.Abbreviations) == 0 {
			return nil, fmt.Errorf("no abbreviations given for %s", name)
		}
		for i, abbr := range p.Abbreviations {
			abbr = strings.TrimSpace(abbr)
			if abbr == "" {
				return nil, fmt.Errorf("empty abbreviation for %s", name)
			}
			for _, bp := range board.provinces {
				if hasStringFold(bp.abbrs, abbr) {
					return nil, fmt.Errorf("duplicate abbreviation %s", abbr)
				}
			}
			p.Abbreviations[i] = abbr
		}
		if p.Country != "" {
			p.Center = true
			country, ok := board.ParseCountry(p.Country)
			if !ok {
				return nil, fmt.Errorf("unknown country %s for %s", p.Country, name)
			}
			p.Country = country
		}
		if p.Terrain != Coastal {
			p.Coasts = nil
		}
		if len(p.Coasts) == 1 {
			return nil, fmt.Errorf("province %s cannot have just one named coast", name)
		}
		for i, c := range p.Coasts {
			coast, ok := board.ParseCoast(c)
			if !ok {
				return nil, fmt.Errorf("unknown coast name %s on %s", c, name)
			}
			p.Coasts[i] = coast
		}
		board.provinces = append(board.provinces, &Province{
			name:    name,
			abbrs:   p.Abbreviations,
			terrain: p.Terrain,
			coasts:  p.Coasts,
			center:  p.Center,
			country: p.Country,
		})
	}
	for i := 0; i < len(b.Connections); i++ {
		c := b.Connections[i]
		if c.ToAll != nil {
			for _, t := range c.ToAll {
				b.Connections = append(b.Connections, BuilderConnection{
					From:    c.From,
					To:      t,
					Coastal: c.Coastal,
				})
			}
			continue
		}
		froms, tos := board.ParseProvince(c.From), board.ParseProvince(c.To)
		if len(froms) != 1 || len(tos) != 1 {
			return nil, fmt.Errorf(
				"unknown provinces in connection %s - %s",
				c.From,
				c.To,
			)
		}
		from, to := froms[0], tos[0]
		e := endpoints{from, to}
		if board.connections[e] != nil || board.connections[endpoints{to, from}] != nil {
			return nil, fmt.Errorf("duplicate connection %s - %s", from.name, to.name)
		}
		for i, coast := range c.FromCoasts {
			var ok bool
			c.FromCoasts[i], ok = board.ParseCoast(coast)
			if !ok {
				return nil, fmt.Errorf("unknown coast %s in connection %s - %s", coast, c.From, c.To)
			}
		}
		for i, coast := range c.ToCoasts {
			var ok bool
			c.ToCoasts[i], ok = board.ParseCoast(coast)
			if !ok {
				return nil, fmt.Errorf("unknown coast %s in connection %s - %s", coast, c.From, c.To)
			}
		}
		connection := &Connection{from, to, c.FromCoasts, c.ToCoasts, c.Coastal}
		if err := connection.valid(); err != nil {
			return nil, err
		}
		board.connections[e] = connection
	}
	return board, nil
}
