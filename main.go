package main

import (
    "os"
    "fmt"
    "log"
    "image/color"
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

    for i := 3; i < len(os.Args); i++ {
        var col color.Color = g.RelationColor
        arg_swallowed := false

        if i < len(os.Args) - 1 && os.Args[i + 1][0] == '#' {
            var r, g, b uint8
            _, err = fmt.Sscanf(os.Args[i + 1], "#%02x%02x%02x", &r, &g, &b)

            if err == nil {
                col = color.RGBA{r, g, b, 0xFF}
                arg_swallowed = true
            }
        }

        expr, err := Eval(os.Args[i])
        if err != nil {
            log.Fatal(err)
        }

        switch expr.(type) {
            case Function:
                g.DrawFunctionWithColor(expr.(Function), col)

            case PolarFunction:
                g.DrawPolarFunctionWithColor(expr.(PolarFunction), col)

            case Relation:
                g.DrawRelationWithColor(expr.(Relation), col)
        }

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