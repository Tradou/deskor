// see here https://stackoverflow.com/questions/75752696/custom-fyne-adaptive-grid-layout

package graphic

import (
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

func New(layout fyne.Layout, objects ...fyne.CanvasObject) *fyne.Container {
	return container.New(layout, objects...)
}

func NewAdaptiveGridWithRatios(ratios []float32, objects ...fyne.CanvasObject) *fyne.Container {
	return New(NewAdaptiveGridLayoutWithRatios(ratios), objects...)
}

// Declare conformity with Layout interface
var _ fyne.Layout = (*adaptiveGridLayoutWithRatios)(nil)

type adaptiveGridLayoutWithRatios struct {
	ratios          []float32
	adapt, vertical bool
}

func NewAdaptiveGridLayoutWithRatios(ratios []float32) fyne.Layout {
	return &adaptiveGridLayoutWithRatios{ratios: ratios, adapt: true}
}

func (g *adaptiveGridLayoutWithRatios) horizontal() bool {
	if g.adapt {
		return fyne.IsHorizontal(fyne.CurrentDevice().Orientation())
	}

	return !g.vertical
}

func (g *adaptiveGridLayoutWithRatios) countRows(objects []fyne.CanvasObject) int {
	count := 0
	for _, child := range objects {
		if child.Visible() {
			count++
		}
	}

	return int(math.Ceil(float64(count) / float64(len(g.ratios))))
}

// Layout is called to pack all child objects into a specified size.
// For a GridLayout this will pack objects into a table format with the number
// of columns specified in our constructor.
func (g *adaptiveGridLayoutWithRatios) Layout(objects []fyne.CanvasObject, size fyne.Size) {

	rows := g.countRows(objects)
	cols := len(g.ratios)

	padWidth := float32(cols-1) * theme.Padding()
	padHeight := float32(rows-1) * theme.Padding()
	tGap := float64(padWidth)
	tcellWidth := float64(size.Width) - tGap
	cellHeight := float64(size.Height-padHeight) / float64(rows)

	if !g.horizontal() {
		padWidth, padHeight = padHeight, padWidth
		tcellWidth = float64(size.Width-padWidth) - tGap
		cellHeight = float64(size.Height-padHeight) / float64(cols)
	}

	row, col := 0, 0
	i := 0
	var x1, x2, y1, y2 float32 = 0.0, 0.0, 0.0, 0.0
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		if i == 0 {
			x1 = 0
			y1 = 0
		} else {
			x1 = x2 + float32(theme.Padding())*float32(1)
			y1 = y2 - float32(cellHeight)
		} // float32(tGap/float64(col))
		//  (size + float64(theme.Padding())) * float64(offset)  float32(theme.Padding())*float32(1)
		x2 = x1 + float32(tcellWidth*float64(g.ratios[i]))
		y2 = float32(cellHeight)

		child.Move(fyne.NewPos(x1, y1))
		child.Resize(fyne.NewSize((x2 - x1), y2-y1))

		if g.horizontal() {
			if (i+1)%cols == 0 {
				row++
				col = 0
			} else {
				col++
			}
		} else {
			if (i+1)%cols == 0 {
				col++
				row = 0
			} else {
				row++
			}
		}
		i++
	}
}

func (g *adaptiveGridLayoutWithRatios) MinSize(objects []fyne.CanvasObject) fyne.Size {
	rows := g.countRows(objects)
	minSize := fyne.NewSize(0, 0)
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		minSize = minSize.Max(child.MinSize())
	}

	if g.horizontal() {
		minContentSize := fyne.NewSize(minSize.Width*float32(len(g.ratios)), minSize.Height*float32(rows))
		return minContentSize.Add(fyne.NewSize(theme.Padding()*fyne.Max(float32(len(g.ratios)-1), 0), theme.Padding()*fyne.Max(float32(rows-1), 0)))
	}

	minContentSize := fyne.NewSize(minSize.Width*float32(rows), minSize.Height*float32(len(g.ratios)))
	return minContentSize.Add(fyne.NewSize(theme.Padding()*fyne.Max(float32(rows-1), 0), theme.Padding()*fyne.Max(float32(len(g.ratios)-1), 0)))
}
