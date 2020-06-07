package main

import (
    "os"
    "fmt"
    "log"
    "image/color"
    "github.com/friedkeenan/gograph"
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

    area, err := gograph.NewArea(x0, y0, x1, y1)
    if err != nil {
        log.Fatal(err)
    }

    _, err = fmt.Sscanf(os.Args[2], "%v", &scale)

    g, err := gograph.NewGraph(area, scale)
    if err != nil {
        log.Fatal(err)
    }

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

        expr, err := gograph.Eval(os.Args[i])
        if err != nil {
            log.Fatal(err)
        }

        switch expr.(type) {
            case gograph.Function:
                g.DrawFunctionWithColor(expr.(gograph.Function), col)

            case gograph.PolarFunction:
                g.DrawPolarFunctionWithColor(expr.(gograph.PolarFunction), col)

            case gograph.Relation:
                g.DrawRelationWithColor(expr.(gograph.Relation), col)
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