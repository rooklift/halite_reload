package gohalite

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
