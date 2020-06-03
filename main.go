package main

import (
    "os"
    "fmt"
    "strings"
    "log"
    "image/color"
)

func main() {
    g := NewGraph(NewArea(-5, 5, 5, -5), 1024, 1024)
    g.DrawGrid()

    for i := 1; i < len(os.Args); i++ {
        col := color.RGBA{0, 0, 0, 0xFF}
        arg_swallowed := false

        if i < len(os.Args) - 1 && os.Args[i + 1][0] == '#' {
            _, err := fmt.Sscanf(os.Args[i + 1], "#%02X%02X%02X", &col.R, &col.G, &col.B)

            if err == nil {
                arg_swallowed = true
            } else {
                log.Fatal(err)
                col = color.RGBA{0, 0, 0, 0xFF}
            }
        }

        if strings.Contains(os.Args[i], "==") {
            expr, err := EvalDiffExpression(os.Args[i])
            if err != nil {
                log.Fatal(err)
            }

            g.DrawDiffExpressionWithColor(expr, col)
        } else {
            expr, err := EvalBoolExpression(os.Args[i])
            if err != nil {
                log.Fatal(err)
            }

            g.DrawBoolExpressionWithColor(expr, col)
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