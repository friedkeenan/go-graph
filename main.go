package main

import (
    "os"
    "log"
)

func main() {
    g := NewGraph(NewArea(-5, 5, 5, -5), 1024, 1024)
    g.DrawGrid()

    for i := 1; i < len(os.Args); i++ {
        expr, err := EvalDiffExpression(os.Args[i])
        if err != nil {
            log.Fatal(err)
        }

        g.DrawDiffExpression(expr)
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