package femtioelva

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
	"strings"
)

type Grid struct {
	box      Box
	cellSize float64
	rows     int
	cols     int
	data     []int // row major
}

// NewGrid covering box with nxn cells.
func NewGrid(box Box, n int) Grid {
	return Grid{
		box:      box,
		cellSize: box.sideLength / float64(n),
		rows:     n,
		cols:     n,
		data:     make([]int, n*n),
	}
}

func (g Grid) Width() int {
	return g.cols
}

func (g Grid) Height() int {
	return g.rows
}

func (g Grid) String() string {
	var s strings.Builder
	for i, v := range g.data {
		if i > 0 && i%g.cols == 0 {
			fmt.Fprintf(&s, "\n")
		}
		if v > 0 {
			fmt.Fprintf(&s, "#")
		} else {
			fmt.Fprintf(&s, ".")
		}
	}
	return s.String()
}

func (g Grid) At(row, col int) int {
	index := row*g.cols + col
	return g.data[index]
}

func (g Grid) Incr(row, col int) {
	index := row*g.cols + col
	g.data[index]++
}

func (g Grid) Set(row, col, val int) {
	index := row*g.cols + col
	g.data[index] = val
}

func (g Grid) Max() int {
	max := math.MinInt
	for _, v := range g.data {
		if v > max {
			max = v
		}
	}
	return max
}

func (g Grid) IncrUTM(east, north float64) {
	// Shift to box coordinates
	minEast, minNorth := LatLong2UTM(g.box.LowLat, g.box.LowLong)
	east -= minEast
	north -= minNorth

	col := int(east / g.cellSize)
	if col < 0 || g.cols < col {
		return // Ignore out of bounds.
	}
	row := g.rows - int(north/g.cellSize) // rows increase south
	if row < 0 || g.rows <= row {
		return // Ignore out of bounds.
	}
	g.Incr(row, col)
}

// Draw grid with cell and margin pixel size as given.
func (g Grid) Draw(cell, margin int) *image.Gray {
	img := image.NewGray(image.Rect(0, 0,
		g.Width()*(cell+margin)+margin,
		g.Height()*(cell+margin)+margin))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.ZP, draw.Src)

	// Make the most popular spots completely dark
	scale := float64(g.Max()) / 255.0

	for row := 0; row < g.Height(); row++ {
		for col := 0; col < g.Width(); col++ {
			if g.At(row, col) == 0 {
				continue
			}
			x0 := margin + row*(margin+cell)
			y0 := margin + col*(margin+cell)
			r := image.Rect(x0, y0, x0+cell, y0+cell)
			val := 255 - uint8(float64(g.At(row, col))/scale) // Higher count -> darker
			draw.Draw(img, r, &image.Uniform{color.Gray{val}}, image.ZP, draw.Src)
		}
	}
	return img
}
