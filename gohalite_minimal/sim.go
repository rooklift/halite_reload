package gohalite

// This is the main simulation code for resolving a turn.
// It's slightly ugly but does have the virtue of working perfectly.

type Simulator struct {
    G                       *Game
    placements              [][]int
    damage                  [][]int
    present                 [][]bool
    zero_strength_winner    []int
}

func NewSimulator(g *Game) *Simulator {
    s := new(Simulator)
    s.G = g.Copy()

    s.placements = make([][]int, g.InitialPlayerCount + 1)
    s.damage = make([][]int, g.InitialPlayerCount + 1)
    s.present = make([][]bool, g.InitialPlayerCount + 1)

    for n := 0 ; n <= g.InitialPlayerCount ; n++ {
        s.placements[n] = make([]int, g.Size)
        s.damage[n] = make([]int, g.Size)
        s.present[n] = make([]bool, g.Size)
    }

    s.zero_strength_winner = make([]int, g.Size)

    return s
}

func (s *Simulator) Simulate() {

    // Each frame, each user "places" a certain amount of stuff on each square.
    // There are 3 sources:
    //
    //      Incoming pieces
    //      Stationary pieces
    //      Production, if that square didn't move
    //
    // The list of "placements" is almost all that is needed to generate the next frame.
    // But because of strength-0 combat being real, we also need a list of mere presences.

    g := s.G

    for i := 0 ; i < g.Size ; i++ {
        for n := 0 ; n <= g.InitialPlayerCount ; n++ {
            s.placements[n][i] = 0
            s.damage[n][i] = 0
            s.present[n][i] = false
        }
        s.zero_strength_winner[i] = 0
    }

    // Set placements from movement...

    for i := 0 ; i < g.Size ; i++ {
        if g.Owner[i] == 0 {
            g.Moves[i] = STILL
        }
        target := g.Movement_to_I(i, g.Moves[i])
        s.placements[g.Owner[i]][target] += g.Strength[i]

        // The player can be considered present at both source and dest:

        s.present[g.Owner[i]][target] = true
        s.present[g.Owner[i]][i] = true
    }

    // Add production from cells that didn't move...

    for i := 0 ; i < g.Size ; i++ {
        if g.Owner[i] != 0 {
            if g.Moves[i] == STILL {
                s.placements[g.Owner[i]][i] += g.Production[i]
            }
        }
    }

    // Cells with a presence, no strength, and no combat do live...

    for i := 0 ; i < g.Size ; i++ {

        players_with_presence := 0
        players_with_strength := 0
        zero_strength_player := 0
        for n := 1 ; n <= g.InitialPlayerCount ; n++ {
            if s.present[n][i] {
                players_with_presence++
                zero_strength_player = n
            }
            if s.placements[n][i] > 0 {
                players_with_strength++
            }
        }

        if players_with_presence == 1 && players_with_strength == 0 && (g.Strength[i] == 0 || g.Owner[i] != 0) {
            combat_flag := false
            for _, neighbour := range g.Neighbours[i] {
                if combat_flag {
                    break
                }
                for n := 1 ; n <= g.InitialPlayerCount ; n++ {
                    if n != zero_strength_player {
                        if s.present[n][neighbour.Index] {
                            combat_flag = true
                        }
                    }
                }
            }
            if combat_flag == false {
                 s.zero_strength_winner[i] = zero_strength_player
            }
        }
    }

    // Cap at 255...

    for i := 0 ; i < g.Size ; i++ {
        for n := 0 ; n <= g.InitialPlayerCount ; n++ {
            if s.placements[n][i] > 255 {
                s.placements[n][i] = 255
            }
        }
    }

    // Note how much damage each placement will do...

    for i := 0 ; i < g.Size ; i++ {
        for n := 0 ; n <= g.InitialPlayerCount ; n++ {
            s.damage[n][i] = s.placements[n][i]
        }
    }

    // Damage from coincidence...

    for i := 0 ; i < g.Size ; i++ {
        for n := 0 ; n <= g.InitialPlayerCount ; n++ {
            if s.damage[n][i] > 0 {
                for t := 0 ; t <= g.InitialPlayerCount ; t++ {
                    if t != n {
                        s.placements[t][i] -= s.damage[n][i]
                    }
                }
            }
        }
    }

    // Damage from adjacency...

    for i := 0 ; i < g.Size ; i++ {
        for n := 1 ; n <= g.InitialPlayerCount ; n++ {              // Note the n := 1, not 0
            if s.damage[n][i] > 0 {
                for _, neighbour := range g.Neighbours[i] {
                    for t := 1 ; t <= g.InitialPlayerCount ; t++ {  // Note the t := 1, not 0
                        if t != n {
                            s.placements[t][neighbour.Index] -= s.damage[n][i]
                        }
                    }
                }
            }
        }
    }

    // Place winner, if any...

    for i := 0 ; i < g.Size ; i++ {
        g.Owner[i] = 0                                          // Neutral by default
        g.Strength[i] = 0
        for n := 0 ; n <= g.InitialPlayerCount ; n++ {
            if s.placements[n][i] > 0 {                         // Should only be true once
                g.Owner[i] = n
                g.Strength[i] = s.placements[n][i]
                break
            }
        }
    }

    // Fix zero strength presences that survived...

    for i := 0 ; i < g.Size ; i++ {
        if s.zero_strength_winner[i] != 0 {
            g.Owner[i] = s.zero_strength_winner[i]
            g.Strength[i] = 0
        }
    }

    g.SetExtraState()   // Includes g.Turn += 1 and resets various slices
    return
}
