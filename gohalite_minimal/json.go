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
    s := fmt.Sprintf("[%d, %d]", n.Owner, n.Strength)
    return []byte(s), nil
}

func (g *Game) Load(filename string, turn int, id int) error {

    file, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

    var hlt HLT
	err = json.Unmarshal(file, &hlt)

	if err != nil {
		return err
	}

    if len(hlt.Frames) <= turn {
        return fmt.Errorf("wanted turn %d but file only had %d frames", turn, len(hlt.Frames))
    }

    // ----------------------------------------------------------

    g.HLT = &hlt

    g.Id = id
    g.Turn = turn

    g.Width = g.HLT.Width
    g.Height = g.HLT.Height
    g.Size = g.Width * g.Height

    g.InitialPlayerCount = g.HLT.NumPlayers

    g.MakeLookupTable()
    g.MakeSlices()

    g.SetBoardFromHLT()

    g.GameStart = time.Now()

	return nil
}

func (g *Game) SetBoardFromHLT() error {

    if g.HLT == nil {
        return fmt.Errorf("no HLT in Game object")
    }

    if len(g.HLT.Frames) <= g.Turn {
        return fmt.Errorf("wanted turn %d but file only had %d frames", g.Turn, len(g.HLT.Frames))
    }

    for y := 0 ; y < g.Height ; y++ {
        for x := 0 ; x < g.Width ; x++ {
            i := g.XY_to_I(x, y)
            g.Production[i] = g.HLT.Productions[y][x]               // note y,x format in source
        }
    }

    for y := 0 ; y < g.Height ; y++ {
        for x := 0 ; x < g.Width ; x++ {
            i := g.XY_to_I(x, y)
            g.Owner[i] = g.HLT.Frames[g.Turn][y][x].Owner             // note y,x format in source
            g.Strength[i] = g.HLT.Frames[g.Turn][y][x].Strength       // note y,x format in source
        }
    }

    g.TurnStart = time.Now()

    return nil
}

func (g *Game) SetMovesFromHLT() error {

    if g.HLT == nil {
        return fmt.Errorf("no HLT in Game object")
    }

    if len(g.HLT.Moves) <= g.Turn {
        return fmt.Errorf("wanted turn %d but file only had %d movelists", g.Turn, len(g.HLT.Moves))
    }

    for y := 0 ; y < g.Height ; y++ {
        for x := 0 ; x < g.Width ; x++ {
            i := g.XY_to_I(x, y)
            g.Moves[i] = g.HLT.Moves[g.Turn][y][x]                 // note y,x format in source
        }
    }

    return nil
}




func (h *HLT) AddFrame(g *Game) {

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
}


func (h *HLT) AddMoves(g *Game) {
    h.Moves = append(h.Moves, make([][]int, g.Height))
    for y := 0 ; y < g.Height ; y++ {
        h.Moves[len(h.Moves) - 1][y] = make([]int, g.Width)
        for x := 0 ; x < g.Width ; x++ {
            h.Moves[len(h.Moves) - 1][y][x] = g.Moves[g.XY_to_I(x,y)]
        }
    }
}

func (h *HLT) SetProductions(g *Game) {
    h.Productions = nil
    for y := 0 ; y < g.Height ; y++ {
        h.Productions = append(h.Productions, nil)
        for x := 0 ; x < g.Width ; x++ {
            h.Productions[y] = append(h.Productions[y], g.Production[g.XY_to_I(x,y)])
        }
    }
}
