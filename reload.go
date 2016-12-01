package main

import (
    "bufio"
    "fmt"
    "os"
    "os/exec"
    "strings"
    "strconv"
    "time"

    hal "./gohalite_minimal"        // A stripped-down version of the library I use for bots
)


const SHOW_PROGRESS_IN_CONSOLE = true
const SLOW = false


func bot_handler(cmd string, id int, g *hal.Game, nudges chan bool) {

    have_warned_eof := false

    cmd_split := strings.Fields(cmd)
    exec_command := exec.Command(cmd_split[0], cmd_split[1:]...)

    i_pipe, err := exec_command.StdinPipe()
    if err != nil {
        fmt.Printf("%v\n", err)
        os.Exit(1)
    }

    o_pipe, err := exec_command.StdoutPipe()
    if err != nil {
        fmt.Printf("%v\n", err)
        os.Exit(1)
    }
/*
    e_pipe, err := exec_command.StderrPipe()
    if err != nil {
        fmt.Printf("%v\n", err)
        os.Exit(1)
    }
*/
    err = exec_command.Start()
    if err != nil {
        fmt.Printf("Failed to start bot %d (%s)\n", id, cmd)
        os.Exit(1)
    }

    fmt.Fprintf(i_pipe, "%d\n", id)
    fmt.Fprintf(i_pipe, "%d %d\n", g.Width, g.Height)
    fmt.Fprintf(i_pipe, "%s\n", g.ProductionMapString())
    fmt.Fprintf(i_pipe, "%s\n", g.GameMapString())

    scanner := bufio.NewScanner(o_pipe)
    if scanner.Scan() == false {
        if have_warned_eof == false {
            fmt.Printf("Turn %d: bot %d output reached EOF\n", g.Turn, id)
            have_warned_eof = true
        }
    }
    botname := scanner.Text()

    for {
        <- nudges                                               // Hub tells us to act
        fmt.Fprintf(i_pipe, "%s\n", g.GameMapString())          // Send the map

        if scanner.Scan() == false {
            if have_warned_eof == false {
                fmt.Printf("Turn %d: bot %d output reached EOF\n", g.Turn, id)
                have_warned_eof = true
            }
        }
        fields := strings.Fields(scanner.Text())

        for n := 0 ; n < len(fields) - 2 ; n += 3 {
            x, err := strconv.Atoi(fields[n])
            if err != nil {
                fmt.Printf("Turn %d: %s sent some unfathomable outputs\n", g.Turn, botname)
                break
            }
            y, err := strconv.Atoi(fields[n + 1])
            if err != nil {
                fmt.Printf("Turn %d: %s sent some unfathomable outputs\n", g.Turn, botname)
                break
            }
            dir, err := strconv.Atoi(fields[n + 2])
            if err != nil {
                fmt.Printf("Turn %d: %s sent some unfathomable outputs\n", g.Turn, botname)
                break
            }

            i := g.XY_to_I(x,y)

            if g.Owner[i] == id {
                if dir >= 0 && dir <= 4 {
                    g.Moves[i] = dir
                }
            }
        }

        nudges <- true                                  // Tell Hub we're done. Using the same channel is OK in this simple case.
    }
}

func main() {

    // I'm making goroutines as bot handlers, though I may still
    // talk to the bots consecutively (rather than concurrently).
    // It may be safer (the bot handler can edit the game struct).

    if len(os.Args) < 4 {
        fmt.Printf("Usage: %s filename start_turn \"bot command one\" ...\n", os.Args[0])
        os.Exit(1)
    }

    hlt_file := os.Args[1]

    initial_turn, err := strconv.Atoi(os.Args[2])
    if err != nil {
        fmt.Printf("Usage: %s  <filename>  <start_turn>  \"bot command one\" ...\n", os.Args[0])
        os.Exit(1)
    }

    botlist := os.Args[3:]

    //

    channels := make([]chan bool, len(botlist))

    for n := 0 ; n < len(botlist) ; n++ {
        channels[n] = make(chan bool)
    }

    tmp := new(hal.Game)
    err = tmp.Load(hlt_file, initial_turn, 0)
    if err != nil {
        fmt.Printf("%v\n", err)
        os.Exit(1)
    }

    sim := hal.NewSimulator(tmp)

    for n := 0 ; n < len(botlist) ; n++ {
        go bot_handler(botlist[n], n + 1, sim.G, channels[n])   // Bot IDs are offset by 1, since ID 0 == neutral
    }

    for {
        for n := 0 ; n < len(botlist) ; n++ {   // This is consecutive, not concurrent
            channels[n] <- true
            <- channels[n]
        }

        sim.Simulate()

        if SHOW_PROGRESS_IN_CONSOLE && sim.G.Turn % 20 == 0 {
            print_map(sim.G)
            if SLOW {
                time.Sleep(1 * time.Second)
            }
        }

        if sim.G.Turn >= 500 {
            fmt.Printf("Turn %d reached\n", sim.G.Turn)
            break
        }

        if len(botlist) > 1 {
            if sim.G.CountPlayers() == 1 {
                if SHOW_PROGRESS_IN_CONSOLE {
                    print_map(sim.G)
                }
                all_strengths := sim.G.TotalStrengths()
                for p := 1 ; p <= sim.G.InitialPlayerCount ; p++ {
                    if all_strengths[p] > 0 {
                        fmt.Printf("Turn %d: Victory for player %d (%s)\n", sim.G.Turn, p, botlist[p - 1])
                    }
                }
                break
            }
        }
    }
}

func print_map(g *hal.Game) {

    translate := ".XRGBVM"

    for y := 0 ; y < g.Height ; y++ {
        for x := 0 ; x < g.Width ; x++ {
            i := g.XY_to_I(x,y)
            fmt.Printf("%c ", translate[g.Owner[i]])
        }
        fmt.Printf("\n")
    }
    fmt.Printf("\n")
}
