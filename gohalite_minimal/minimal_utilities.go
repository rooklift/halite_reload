package gohalite

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
