package main

import (
	"fmt"
	"slices"

	diplo "github.com/adambyle/diplopad"
)

func main() {
	board := diplo.StandardBoard
	for p := range board.Provinces() {
		cs := slices.Collect(board.ConnectionsFrom(p))
		fmt.Println(p.Name(), len(cs))
		for _, c := range cs {
			if c.Coastal() {
				fmt.Println(" ", c.FromCoasts(), "->", c.To().Name(), c.ToCoasts(), "(Coastal)")
			} else {
				fmt.Println(" ", c.FromCoasts(), "->", c.To().Name(), c.ToCoasts())
			}
		}
	}
}
