package gohalite

import (
    "fmt"
    "strings"
)

func (g *Game) ProductionMapString() string {
    var components []string
    for y := 0 ; y < g.Height ; y++ {
        for x := 0 ; x < g.Width ; x++ {
            components = append(components, fmt.Sprintf("%d", g.Production[g.XY_to_I(x,y)]))
        }
    }
    return strings.Join(components, " ")
}

func (g *Game) GameMapString() string {
    var components []string

    // FIXME: we're supposed to send a sort of RLE thing

    for y := 0 ; y < g.Height ; y++ {
        for x := 0 ; x < g.Width ; x++ {
            components = append(components, fmt.Sprintf("1 %d", g.Owner[g.XY_to_I(x,y)]))
        }
    }

    for y := 0 ; y < g.Height ; y++ {
        for x := 0 ; x < g.Width ; x++ {
            components = append(components, fmt.Sprintf("%d", g.Strength[g.XY_to_I(x,y)]))
        }
    }
    return strings.Join(components, " ")
}
