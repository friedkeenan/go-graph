package main

import (
    "os"
    "fmt"
    "log"
)

func main() {
    if len(os.Args) < 3 {
        log.Fatal("Too few arguments")
    }

    var x0, y0, x1, y1, scale float64
    _, err := fmt.Sscanf(os.Args[1], "{(%v, %v), (%v, %v)}", &x0, &y0, &x1, &y1)
    if err != nil {
        log.Fatal(err)
    }

    if x1 < x0 || y1 > y0 {
        log.Fatal("Invalid bounds")
    }

    _, err = fmt.Sscanf(os.Args[2], "%v", &scale)

    g := NewGraph(NewArea(x0, y0, x1, y1), scale)
    g.DrawGrid()

    g.DrawDifferentialExpression(func (c *Coord) float64 {
        return c.Y
    }, NewCoord(0, 1))

    for i := 3; i < len(os.Args); i++ {
        col := ExpressionColor
        arg_swallowed := false

        if i < len(os.Args) - 1 && os.Args[i + 1][0] == '#' {
            _, err = fmt.Sscanf(os.Args[i + 1], "#%02x%02x%02x", &col.R, &col.G, &col.B)

            if err == nil {
                arg_swallowed = true
            } else {
                col = ExpressionColor
            }
        }

        expr, err := EvalExpression(os.Args[i])
        if err != nil {
            log.Fatal(err)
        }

        g.DrawExpressionWithColor(expr, col)

        if arg_swallowed {
            i++
        }
    }

    f, err := os.Create("out.png")
    if err != nil {
        log.Fatal(err)
    }

    defer f.Close()

    err = g.SavePNG(f)
    if err != nil {
        log.Fatal(err)
    }
}