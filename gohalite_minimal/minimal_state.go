package gohalite

func (g *Game) SetExtraState() {
    g.Turn += 1
    for i := 0 ; i < g.Size ; i++ {
        g.Moves[i] = STILL
    }
    // There's more to this function in the real library.
}

func (g *Game) MakeLookupTable() {
    g.Neighbours = make([][]Neighbour, g.Size, g.Size)
    for i := 0 ; i < g.Size ; i++ {

        g.Neighbours[i] = make([]Neighbour, 4, 4)

        x, y := g.I_to_XY(i)

        // Do these in a swirl pattern? Has some subtle ramifications...

        g.Neighbours[i][0] = Neighbour{g.XY_to_I(x, y - 1), UP}
        g.Neighbours[i][1] = Neighbour{g.XY_to_I(x + 1, y), RIGHT}
        g.Neighbours[i][2] = Neighbour{g.XY_to_I(x, y + 1), DOWN}
        g.Neighbours[i][3] = Neighbour{g.XY_to_I(x - 1, y), LEFT}
    }
}

func (g *Game) MakeSlices() {
    g.Production = make([]int, g.Size, g.Size)
    g.Owner = make([]int, g.Size, g.Size)
    g.Strength = make([]int, g.Size, g.Size)

    g.Moves = make([]int, g.Size, g.Size)
    g.HasOrders = make([]bool, g.Size, g.Size)
    g.Allocation = make([]int, g.Size, g.Size)
}

func (g *Game) SetMove(index, direction int) {
    g.Moves[index] = direction
    // There's more to this function in the real library.
}
