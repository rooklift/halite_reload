package gohalite

import (
    "time"
)

const (
    STILL = 0
    NORTH = 1
    EAST = 2
    SOUTH = 3
    WEST = 4

    UP = NORTH
    RIGHT = EAST
    DOWN = SOUTH
    LEFT = WEST
)

type Neighbour struct {
    Index               int
    Dir                 int
}

type Game struct {

    // HLT file that the game was loaded from, if any (treat as read-only):
    HLT                 *HLT

    // Lookup table of neighbouring indices and directions:
    Neighbours          [][]Neighbour

    // Constant after game started:
    GameStart           time.Time
    Width               int
    Height              int
    Size                int
    Id                  int
    InitialPlayerCount  int
    Logfile             *Logfile

    // Single values that get set each turn:
    TurnStart           time.Time
    Turn                int

    // Slices handled by the main parsers:
    Production          []int
    Owner               []int
    Strength            []int

    // Other slices, updated each turn:
    Moves               []int           // Direction to move this turn?
}

func (g *Game) Copy() *Game {

    // There must be a better way of deep copying a struct?

    result := new(Game)

    result.HLT = g.HLT

    result.GameStart = g.GameStart
    result.Width = g.Width
    result.Height = g.Height
    result.Size = g.Size
    result.Id = g.Id
    result.InitialPlayerCount = g.InitialPlayerCount
    result.Logfile = g.Logfile

    result.TurnStart = g.TurnStart
    result.Turn = g.Turn

    result.Neighbours = g.Neighbours    // OK since the original is read-only in practice
    result.MakeSlices()

    copy(result.Production, g.Production)
    copy(result.Owner, g.Owner)
    copy(result.Strength, g.Strength)
    copy(result.Moves, g.Moves)
    // There's more to this function in the real library.

    return result
}

func (g *Game) CountPlayers() int {
    set := make(map[int]bool)
    for i := 0 ; i < g.Size ; i++ {
        if g.Owner[i] > 0 {
            set[g.Owner[i]] = true
        }
    }
    return len(set)
}

func (g *Game) TotalStrengths() []int {
    result := make([]int, g.InitialPlayerCount + 1, g.InitialPlayerCount + 1)
    for i := 0 ; i < g.Size ; i++ {
        result[g.Owner[i]] += g.Strength[i]
    }
    return result
}

func (g *Game) StrengthOfPlayer(id int) int {
    result := 0
    for i := 0 ; i < g.Size ; i++ {
        if g.Owner[i] == id {
            result += g.Strength[i]
        }
    }
    return result
}

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
    // There's more to this function in the real library.
}

func (g *Game) SetMove(index, direction int) {
    g.Moves[index] = direction
    // There's more to this function in the real library.
}

func (g *Game) I_to_XY(i int) (int, int) {
    x := i % g.Width
    y := i / g.Width
    return x, y
}

func (g *Game) XY_to_I(x, y int) (int) {

    if x < 0 {
        x += -(x / g.Width) * g.Width + g.Width      // Can make x == g.Width, so must still use % later
    }

    x %= g.Width

    if y < 0 {
        y += -(y / g.Height) * g.Height + g.Height   // Can make y == g.Height, so must still use % later
    }

    y %= g.Height

    return y * g.Width + x
}

func (g *Game) Movement_to_I(src, direction int) int {

    // Given a source square and a direction, what index do we land on?

    x, y := g.I_to_XY(src)

    switch direction {
    case RIGHT:
        return g.XY_to_I(x + 1, y)
    case LEFT:
        return g.XY_to_I(x - 1, y)
    case UP:
        return g.XY_to_I(x, y - 1)
    case DOWN:
        return g.XY_to_I(x, y + 1)
    default:
        return src
    }
}
