package gohalite

import (
    "crypto/sha1"
    "fmt"
    "os"
    "strings"
)

type Logfile struct {
    outfile         *os.File
    outfilename     string
    enabled         bool
    logged_once     map[string]bool
}

func NewLog(outfilename string, enabled bool) *Logfile {
    logged_once := make(map[string]bool)
    return &Logfile{nil, outfilename, enabled, logged_once}
}

func (log *Logfile) Dump(format_string string, args ...interface{}) {

    if log == nil {
        return
    }

    if log.enabled == false {
        return
    }

    if log.outfile == nil {

        var err error

        if _, tmp_err := os.Stat(log.outfilename); tmp_err == nil {
            // File exists
            log.outfile, err = os.OpenFile(log.outfilename, os.O_APPEND|os.O_WRONLY, 0666)
        } else {
            // File needs creating
            log.outfile, err = os.Create(log.outfilename)
        }

        if err != nil {
            log.enabled = false
            return
        }
    }

    fmt.Fprintf(log.outfile, format_string, args...)
    fmt.Fprintf(log.outfile, "\r\n")                    // Because I use Windows...
}

func (g *Game) Log(format_string string, args ...interface{}) {
    g.Logfile.Dump(format_string, args...)
}

func (g *Game) LogOnce(format_string string, args ...interface{}) bool {
    if g.Logfile.logged_once[format_string] == false {
        g.Logfile.logged_once[format_string] = true         // Note that it's format_string that is checked / saved
        g.Logfile.Dump(format_string, args...)
        return true
    }
    return false
}

func (g *Game) LogValueMap(value_map []int, translate map[int]string) {

    max_value := 0
    for i := 0 ; i < g.Size ; i++ {
        if value_map[i] > max_value {
            max_value = value_map[i]
        }
    }

    for y := 0 ; y < g.Height ; y++ {
        s := ""
        for x := 0 ; x < g.Width ; x++ {

            var add string

            if translate != nil {
                add = translate[value_map[g.XY_to_I(x, y)]]
            } else {
                add = fmt.Sprintf("%d", value_map[g.XY_to_I(x, y)])
            }

            if max_value > 99 {
                s += fmt.Sprintf(" %3s", add)
            } else if max_value > 9 {
                s += fmt.Sprintf(" %2s", add)
            } else {
                s += fmt.Sprintf(" %1s", add)
            }
        }
        g.Log(s)
    }
    g.Log("")
}

func (g *Game) LogOwner() {
    g.LogValueMap(g.Owner, nil)
}

func (g *Game) LogStrength() {
    g.LogValueMap(g.Strength, nil)
}

func (g *Game) LogProduction() {
    g.LogValueMap(g.Production, nil)
}

func (g *Game) LogMoves() {
    translate := make(map[int]string)
    translate[0] = "."
    translate[1] = "^"
    translate[2] = ">"
    translate[3] = "v"
    translate[4] = "<"
    g.LogValueMap(g.Moves, translate)
}

func hash_from_string(datastring string) string {
    data := []byte(datastring)
    sum := sha1.Sum(data)
    return fmt.Sprintf("%x", sum)
}

func (g *Game) MovesHash() string {
    var components []string
    for i := 0 ; i < g.Size ; i++ {
        if g.Owner[i] == g.Id && g.Moves[i] != STILL {
            components = append(components, fmt.Sprintf("%d %d", i, g.Moves[i]))
        }
    }
    fullstring := strings.Join(components, " ")
    return hash_from_string(fullstring)
}

func (g *Game) BoardHash() string {
    var components []string

    for i := 0 ; i < g.Size ; i++ {
        components = append(components, fmt.Sprintf("%d %d %d", g.Owner[i], g.Production[i], g.Strength[i]))
    }

    fullstring := strings.Join(components, " ")
    return hash_from_string(fullstring)
}

func (g *Game) LogMovesHash() {
    g.Log("Moves SHA1: %s", g.MovesHash())
}

func (g *Game) LogBoardHash() {
    g.Log("Board SHA1: %s", g.BoardHash())
}
