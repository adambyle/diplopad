package orders

import (
	"github.com/adambyle/diplopad/diplomacy/geo"
	"github.com/adambyle/diplopad/diplomacy/unit"
)

// Hold leaves a unit in place.
type Hold struct {
	Unit *geo.Province // The unit to leave in place.
}

// Move transports a unit from one place to another.
type Move struct {
	Unit  *geo.Province // The source territory.
	Dest  *geo.Province // The destination territory.
	Coast geo.Coast     // The destination coast (optional).
}

// Support adds power to another unit (optionally to some place).
type Support struct {
	Unit   *geo.Province // The unit giving support.
	Target *geo.Province // The unit receiving support.
	Dest   *geo.Province // The destination of the target unit (may be nil for hold).
	Coast  geo.Coast     // The destination coast (optional).
}

// Convoy transports an Army across water provinces.
type Convoy struct {
	Unit   *geo.Province // The Fleet convoying.
	Target *geo.Province // The Army being moved.
	Dest   *geo.Province // The final destination of the Army.
}

// Orders collects all types of orders for one nation.
type Orders struct {
	Holds    []Hold
	Moves    []Move
	Supports []Support
	Convoys  []Convoy
}

// Build adds or removes units.
type Build struct {
	Unit    unit.Unit     // The type of unit to construct.
	Target  *geo.Province // The province to build/disband in.
	Coast   geo.Coast     // The coast to build on (for fleets when needed).
	Disband bool          // Disband here instead; ignores all but Target.
}

// Retreat directs units that have been dislodged.
type Retreat struct {
	Unit *geo.Province // The unit that was dislodged.
	Dest *geo.Province // Where to retreat (nil for disband).
}
