package main

// This tests the simulator against an actual HLT game record.
// Note that the displayed hashes are not file hashes, but rather hashes of a string of the state.

import (
    "bufio"
    "fmt"
    "os"

    hal "./gohalite_minimal"
)

func print_comparison(sim *hal.Game, real *hal.Game) {

    for i := 0 ; i < sim.Size ; i++ {
        x, y := sim.I_to_XY(i)
        if sim.Owner[i] != real.Owner[i] {
            fmt.Printf("   i %d [%d,%d]: Sim owner: %d ... HLT owner: %d\n", i, x, y, sim.Owner[i], real.Owner[i])
        }
        if sim.Strength[i] != real.Strength[i] {
            fmt.Printf("   i %d [%d,%d]: Sim strength: %d ... HLT strength: %d\n", i, x, y, sim.Strength[i], real.Strength[i])
        }
    }
}


func main() {

    if len(os.Args) == 1 {
        fmt.Printf("Usage: %s <filename>\n", os.Args[0])
        os.Exit(1)
    }

    MAINLOOP:
    for _, filename := range os.Args[1:] {

        fmt.Printf("%s\n", filename)

        tmp := new(hal.Game)
        tmp.Load(filename, 0, 1)
        s := hal.NewSimulator(tmp)

        comp := new(hal.Game)
        comp.Load(filename, 0, 1)

        for turn := 0 ; turn < len(s.G.HLT.Frames) - 1 ; turn++ {

            s.G.SetMovesFromHLT()
            s.Simulate()

            comp.Turn += 1
            comp.SetBoardFromHLT()

            if s.G.BoardHash() != comp.BoardHash() {
                fmt.Printf("\n")
                fmt.Printf("   Turn %3d: %s (simulated)\n", s.G.Turn, s.G.BoardHash())
                fmt.Printf("   Turn %3d: %s (HLT file)\n", comp.Turn, comp.BoardHash())
                fmt.Printf("\n")
                print_comparison(s.G, comp)
                fmt.Printf("\n")
                continue MAINLOOP
            }
        }

        fmt.Printf("Turn %3d: %s (simulated)\n", s.G.Turn, s.G.BoardHash())
        fmt.Printf("Turn %3d: %s (HLT file)\n", comp.Turn, comp.BoardHash())
    }

    var scanner *bufio.Scanner = bufio.NewScanner(os.Stdin)
    fmt.Printf("Done. Type any input to quit.\n")
    scanner.Scan()
}
