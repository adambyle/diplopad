// Package unit exports unit constants
package unit

type Unit byte

const (
	Army  Unit = iota // Can occupy land and coastal territory.
	Fleet             // Can occupy water and coastal territory, and convoy Armies.
)

func (u Unit) String() string {
	switch u {
	case Army:
		return "A"
	case Fleet:
		return "F"
	default:
		return ""
	}
}
