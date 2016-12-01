# halite_reload
Unofficial Halite Engine / Environment / "Clone"


It's mostly for testing purposes. It allows you to load an HLT file **at any turn** and have your bots play on from there. Sample usage:

    reload.exe  some_file.hlt  100  "ruby mybot.rb"  "python my_other_bot.py"  "foobot.exe"

This loads the HLT file at turn 100 and plays on from there, using the specified bots to generate the moves.

It is written in [Go](https://golang.org/) using only the standard library, and so should compile everywhere. It does not currently save any output to file (but it prints to Stdout). The engine has been extensively tested for correctness. The source code for the combat/move resolution system is [here](https://github.com/fohristiwhirl/halite_reload/blob/master/gohalite_minimal/sim.go).

Yes, this is basically a complete engine, except that without worldgen code, it cannot start a new game from scratch, but must load from an HLT file (possibly at turn 0).
