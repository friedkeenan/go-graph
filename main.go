package main

import (
	"os"
	"log"
	"math"
	"image/color"
)

func main() {
	g := NewGraph(NewArea(-5, 5, 5, -5), 2048, 2048)

	g.DrawGrid()

	/* Draws a disk of radius 1 centered at the origin */
	g.DrawBoolExpression(func (c *Coord) bool {
		return math.Pow(c.Y, 2) + math.Pow(c.X, 2) <= 1
	})

	/* Draws a circle of radius 1 centered at the origin */
	g.DrawDiffExpressionWithColor(func (c *Coord) float64 {
		return (math.Pow(c.Y, 2) + math.Pow(c.X, 2)) - 1
	}, color.RGBA{0xFF, 0x00, 0xFF, 0xFF})

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