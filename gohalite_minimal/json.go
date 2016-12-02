package gohalite

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "time"
)

type Site struct {
    Owner       int
    Strength    int
}

type HLT struct {
    Version     int                 `json:"version"`
    Width       int                 `json:"width"`
    Height      int                 `json:"height"`
    NumPlayers  int                 `json:"num_players"`
    NumFrames   int                 `json:"num_frames"`
    PlayerNames []string            `json:"player_names"`
    Productions [][]int             `json:"productions"`
    Frames      [][][]Site          `json:"frames"`
    Moves       [][][]int           `json:"moves"`
}

func (n *Site) UnmarshalJSON(buf []byte) error {

    // Convert a len-2 json array into a Site.
    // Stolen from http://eagain.net/articles/go-json-array-to-struct/

    tmp := []interface{}{&n.Owner, &n.Strength}
    wantLen := len(tmp)
    if err := json.Unmarshal(buf, &tmp); err != nil {
        return err
    }
    if g, e := len(tmp), wantLen; g != e {
        return fmt.Errorf("wrong number of fields in Site: %d != %d", g, e)
    }
    return nil
}

func (n *Site) MarshalJSON() ([]byte, error) {
    s := fmt.Sprintf("[%d,%d]", n.Owner, n.Strength)
    return []byte(s), nil
}

func LoadHLT(filename string) (*HLT, error) {

    file, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }

    hlt := new(HLT)
    err = json.Unmarshal(file, hlt)

    if err != nil {
        return nil, err
    }

    return hlt, nil
}

func (g *Game) SetBoardFromHLT(hlt *HLT, turn int, id int) error {

    if len(hlt.Frames) <= turn {
        return fmt.Errorf("SetBoardFromHLT: wanted turn %d but file only had %d frames", turn, len(hlt.Frames))
    }

    if g.Width != hlt.Width || g.Height != hlt.Height {
        g.Width = hlt.Width
        g.Height = hlt.Height
        g.Size = g.Width * g.Height
        g.MakeLookupTable()
        g.MakeSlices()
    }

    g.Id = id
    g.Turn = turn

    g.InitialPlayerCount = hlt.NumPlayers

    for y := 0 ; y < g.Height ; y++ {
        for x := 0 ; x < g.Width ; x++ {
            i := g.XY_to_I(x, y)
            g.Production[i] = hlt.Productions[y][x]                 // note y,x format in source
            g.Owner[i] = hlt.Frames[g.Turn][y][x].Owner             // note y,x format in source
            g.Strength[i] = hlt.Frames[g.Turn][y][x].Strength       // note y,x format in source
        }
    }

    g.TurnStart = time.Now()

    return nil
}

func (g *Game) SetMovesFromHLT(hlt *HLT) error {

    if hlt.Width != g.Width || hlt.Height != g.Height {
        return fmt.Errorf("SetMovesFromHLT: HLT dimensions didn't match game")
    }

    if len(hlt.Moves) <= g.Turn {
        return fmt.Errorf("SetMovesFromHLT: wanted turn %d but file only had %d movelists", g.Turn, len(hlt.Moves))
    }

    for y := 0 ; y < g.Height ; y++ {
        for x := 0 ; x < g.Width ; x++ {
            i := g.XY_to_I(x, y)
            g.Moves[i] = hlt.Moves[g.Turn][y][x]                    // note y,x format in source
        }
    }

    return nil
}

// Note that the HLT file stored in the game object (if any) is what we loaded the game from.
// These functions that follow are to save to some other HLT file, not that one.

func (h *HLT) AddFrame(g *Game) error {

    if h.Width != g.Width || h.Height != g.Height {
        return fmt.Errorf("AddFrame: HLT dimensions didn't match game")
    }

    h.Frames = append(h.Frames, make([][]Site, g.Height))
    for y := 0 ; y < g.Height ; y++ {
        h.Frames[len(h.Frames) - 1][y] = make([]Site, g.Width)
        for x := 0 ; x < g.Width ; x++ {
            site := Site{}
            site.Owner = g.Owner[g.XY_to_I(x,y)]
            site.Strength = g.Strength[g.XY_to_I(x,y)]
            h.Frames[len(h.Frames) - 1][y][x] = site
        }
    }

    h.NumFrames++

    return nil
}

func (h *HLT) AddMoves(g *Game) error {

    if h.Width != g.Width || h.Height != g.Height {
        return fmt.Errorf("AddMoves: HLT dimensions didn't match game")
    }

    h.Moves = append(h.Moves, make([][]int, g.Height))
    for y := 0 ; y < g.Height ; y++ {
        h.Moves[len(h.Moves) - 1][y] = make([]int, g.Width)
        for x := 0 ; x < g.Width ; x++ {
            h.Moves[len(h.Moves) - 1][y][x] = g.Moves[g.XY_to_I(x,y)]
        }
    }

    return nil
}

func (h *HLT) SetProductions(g *Game) error {

    if h.Width != g.Width || h.Height != g.Height {
        return fmt.Errorf("SetProductions: HLT dimensions didn't match game")
    }

    h.Productions = nil
    for y := 0 ; y < g.Height ; y++ {
        h.Productions = append(h.Productions, nil)
        for x := 0 ; x < g.Width ; x++ {
            h.Productions[y] = append(h.Productions[y], g.Production[g.XY_to_I(x,y)])
        }
    }

    return nil
}
