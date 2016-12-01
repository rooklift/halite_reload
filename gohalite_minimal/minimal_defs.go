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
    HasOrders           []bool          // Has this piece explicitly been given orders (even if STILL)?
    Allocation          []int           // How much strength are we sending to this spot?
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
    copy(result.HasOrders, g.HasOrders)
    copy(result.Allocation, g.Allocation)

    return result
}

type Neighbour struct {
    Index               int
    Dir                 int
}
