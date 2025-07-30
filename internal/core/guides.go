package core

import (
	"image/color"
	"math"
)

// Smart guides and alignment system

// Guide represents a guide line
type Guide struct {
	Type        GuideType
	Position    float64
	Orientation GuideOrientation
	Visible     bool
	Color       color.Color
}

// GuideType represents the type of guide
type GuideType int

const (
	GuideTypeManual GuideType = iota
	GuideTypeGrid
	GuideTypeCenter
	GuideTypeBaseline
	GuideTypeMargin
)

// GuideOrientation represents guide orientation
type GuideOrientation int

const (
	GuideHorizontal GuideOrientation = iota
	GuideVertical
)

// GuideManager manages guides and alignment
type GuideManager struct {
	Guides          []*Guide
	GridSize        float64
	GridVisible     bool
	GridColor       color.Color
	SnapToGrid      bool
	SnapToGuides    bool
	SnapDistance    float64
	ShowCenterLines bool
	ShowBaselines   bool
	Margins         Margins
}

// Margins represents canvas margins
type Margins struct {
	Top, Right, Bottom, Left float64
}

// AlignmentTarget represents something that can be aligned
type AlignmentTarget struct {
	X, Y, Width, Height float64
	ID                  string
}

// NewGuideManager creates a new guide manager
func NewGuideManager() *GuideManager {
	return &GuideManager{
		Guides:          make([]*Guide, 0),
		GridSize:        20,
		GridVisible:     false,
		GridColor:       color.RGBA{200, 200, 200, 100},
		SnapToGrid:      false,
		SnapToGuides:    true,
		SnapDistance:    5,
		ShowCenterLines: false,
		ShowBaselines:   false,
		Margins:         Margins{20, 20, 20, 20},
	}
}

// AddGuide adds a manual guide
func (gm *GuideManager) AddGuide(position float64, orientation GuideOrientation) *Guide {
	guide := &Guide{
		Type:        GuideTypeManual,
		Position:    position,
		Orientation: orientation,
		Visible:     true,
		Color:       color.RGBA{0, 150, 255, 200},
	}
	gm.Guides = append(gm.Guides, guide)
	return guide
}

// RemoveGuide removes a guide
func (gm *GuideManager) RemoveGuide(guide *Guide) bool {
	for i, g := range gm.Guides {
		if g == guide {
			gm.Guides = append(gm.Guides[:i], gm.Guides[i+1:]...)
			return true
		}
	}
	return false
}

// ClearGuides removes all manual guides
func (gm *GuideManager) ClearGuides() {
	newGuides := make([]*Guide, 0)
	for _, guide := range gm.Guides {
		if guide.Type != GuideTypeManual {
			newGuides = append(newGuides, guide)
		}
	}
	gm.Guides = newGuides
}

// SnapPoint snaps a point to guides or grid
func (gm *GuideManager) SnapPoint(x, y float64) (float64, float64) {
	snappedX, snappedY := x, y

	// Snap to grid
	if gm.SnapToGrid && gm.GridSize > 0 {
		snappedX = math.Round(x/gm.GridSize) * gm.GridSize
		snappedY = math.Round(y/gm.GridSize) * gm.GridSize
	}

	// Snap to guides
	if gm.SnapToGuides {
		for _, guide := range gm.Guides {
			if !guide.Visible {
				continue
			}

			switch guide.Orientation {
			case GuideVertical:
				if math.Abs(x-guide.Position) <= gm.SnapDistance {
					snappedX = guide.Position
				}
			case GuideHorizontal:
				if math.Abs(y-guide.Position) <= gm.SnapDistance {
					snappedY = guide.Position
				}
			}
		}
	}

	return snappedX, snappedY
}

// SnapRectangle snaps a rectangle to guides or grid
func (gm *GuideManager) SnapRectangle(x, y, width, height float64) (float64, float64, float64, float64) {
	// Snap top-left corner
	snappedX, snappedY := gm.SnapPoint(x, y)

	// Snap bottom-right corner
	rightX, bottomY := gm.SnapPoint(x+width, y+height)

	// Adjust width and height based on snapped corners
	snappedWidth := rightX - snappedX
	snappedHeight := bottomY - snappedY

	return snappedX, snappedY, snappedWidth, snappedHeight
}

// GetNearestGuides returns guides near a point
func (gm *GuideManager) GetNearestGuides(x, y float64, maxDistance float64) []*Guide {
	var nearGuides []*Guide

	for _, guide := range gm.Guides {
		if !guide.Visible {
			continue
		}

		var distance float64
		switch guide.Orientation {
		case GuideVertical:
			distance = math.Abs(x - guide.Position)
		case GuideHorizontal:
			distance = math.Abs(y - guide.Position)
		}

		if distance <= maxDistance {
			nearGuides = append(nearGuides, guide)
		}
	}

	return nearGuides
}

// GenerateCenterGuides generates center guides for the canvas
func (gm *GuideManager) GenerateCenterGuides(width, height float64) {
	// Remove existing center guides
	newGuides := make([]*Guide, 0)
	for _, guide := range gm.Guides {
		if guide.Type != GuideTypeCenter {
			newGuides = append(newGuides, guide)
		}
	}
	gm.Guides = newGuides

	if gm.ShowCenterLines {
		// Vertical center
		centerV := &Guide{
			Type:        GuideTypeCenter,
			Position:    width / 2,
			Orientation: GuideVertical,
			Visible:     true,
			Color:       color.RGBA{255, 100, 100, 150},
		}
		gm.Guides = append(gm.Guides, centerV)

		// Horizontal center
		centerH := &Guide{
			Type:        GuideTypeCenter,
			Position:    height / 2,
			Orientation: GuideHorizontal,
			Visible:     true,
			Color:       color.RGBA{255, 100, 100, 150},
		}
		gm.Guides = append(gm.Guides, centerH)
	}
}

// GenerateMarginGuides generates margin guides
func (gm *GuideManager) GenerateMarginGuides(width, height float64) {
	// Remove existing margin guides
	newGuides := make([]*Guide, 0)
	for _, guide := range gm.Guides {
		if guide.Type != GuideTypeMargin {
			newGuides = append(newGuides, guide)
		}
	}
	gm.Guides = newGuides

	// Add margin guides
	margins := []struct {
		position    float64
		orientation GuideOrientation
	}{
		{gm.Margins.Left, GuideVertical},
		{width - gm.Margins.Right, GuideVertical},
		{gm.Margins.Top, GuideHorizontal},
		{height - gm.Margins.Bottom, GuideHorizontal},
	}

	for _, margin := range margins {
		guide := &Guide{
			Type:        GuideTypeMargin,
			Position:    margin.position,
			Orientation: margin.orientation,
			Visible:     true,
			Color:       color.RGBA{100, 255, 100, 120},
		}
		gm.Guides = append(gm.Guides, guide)
	}
}

// Alignment functions

// AlignTargetsLeft aligns targets to the leftmost position
func AlignTargetsLeft(targets []AlignmentTarget) []AlignmentTarget {
	if len(targets) == 0 {
		return targets
	}

	minX := targets[0].X
	for _, target := range targets {
		if target.X < minX {
			minX = target.X
		}
	}

	for i := range targets {
		targets[i].X = minX
	}

	return targets
}

// AlignTargetsRight aligns targets to the rightmost position
func AlignTargetsRight(targets []AlignmentTarget) []AlignmentTarget {
	if len(targets) == 0 {
		return targets
	}

	maxX := targets[0].X + targets[0].Width
	for _, target := range targets {
		if target.X+target.Width > maxX {
			maxX = target.X + target.Width
		}
	}

	for i := range targets {
		targets[i].X = maxX - targets[i].Width
	}

	return targets
}

// AlignTargetsTop aligns targets to the topmost position
func AlignTargetsTop(targets []AlignmentTarget) []AlignmentTarget {
	if len(targets) == 0 {
		return targets
	}

	minY := targets[0].Y
	for _, target := range targets {
		if target.Y < minY {
			minY = target.Y
		}
	}

	for i := range targets {
		targets[i].Y = minY
	}

	return targets
}

// AlignTargetsBottom aligns targets to the bottommost position
func AlignTargetsBottom(targets []AlignmentTarget) []AlignmentTarget {
	if len(targets) == 0 {
		return targets
	}

	maxY := targets[0].Y + targets[0].Height
	for _, target := range targets {
		if target.Y+target.Height > maxY {
			maxY = target.Y + target.Height
		}
	}

	for i := range targets {
		targets[i].Y = maxY - targets[i].Height
	}

	return targets
}

// AlignTargetsCenterHorizontal aligns targets horizontally to center
func AlignTargetsCenterHorizontal(targets []AlignmentTarget) []AlignmentTarget {
	if len(targets) == 0 {
		return targets
	}

	// Find the center of all targets
	var totalCenterX float64
	for _, target := range targets {
		totalCenterX += target.X + target.Width/2
	}
	centerX := totalCenterX / float64(len(targets))

	for i := range targets {
		targets[i].X = centerX - targets[i].Width/2
	}

	return targets
}

// AlignTargetsCenterVertical aligns targets vertically to center
func AlignTargetsCenterVertical(targets []AlignmentTarget) []AlignmentTarget {
	if len(targets) == 0 {
		return targets
	}

	// Find the center of all targets
	var totalCenterY float64
	for _, target := range targets {
		totalCenterY += target.Y + target.Height/2
	}
	centerY := totalCenterY / float64(len(targets))

	for i := range targets {
		targets[i].Y = centerY - targets[i].Height/2
	}

	return targets
}

// DistributeHorizontally distributes targets evenly horizontally
func DistributeHorizontally(targets []AlignmentTarget) []AlignmentTarget {
	if len(targets) < 3 {
		return targets
	}

	// Sort by X position
	for i := 0; i < len(targets)-1; i++ {
		for j := i + 1; j < len(targets); j++ {
			if targets[i].X > targets[j].X {
				targets[i], targets[j] = targets[j], targets[i]
			}
		}
	}

	// Calculate spacing
	leftmost := targets[0].X
	rightmost := targets[len(targets)-1].X + targets[len(targets)-1].Width
	totalWidth := rightmost - leftmost

	var objectsWidth float64
	for _, target := range targets {
		objectsWidth += target.Width
	}

	spacing := (totalWidth - objectsWidth) / float64(len(targets)-1)

	// Distribute
	currentX := leftmost
	for i := range targets {
		targets[i].X = currentX
		currentX += targets[i].Width + spacing
	}

	return targets
}

// DistributeVertically distributes targets evenly vertically
func DistributeVertically(targets []AlignmentTarget) []AlignmentTarget {
	if len(targets) < 3 {
		return targets
	}

	// Sort by Y position
	for i := 0; i < len(targets)-1; i++ {
		for j := i + 1; j < len(targets); j++ {
			if targets[i].Y > targets[j].Y {
				targets[i], targets[j] = targets[j], targets[i]
			}
		}
	}

	// Calculate spacing
	topmost := targets[0].Y
	bottommost := targets[len(targets)-1].Y + targets[len(targets)-1].Height
	totalHeight := bottommost - topmost

	var objectsHeight float64
	for _, target := range targets {
		objectsHeight += target.Height
	}

	spacing := (totalHeight - objectsHeight) / float64(len(targets)-1)

	// Distribute
	currentY := topmost
	for i := range targets {
		targets[i].Y = currentY
		currentY += targets[i].Height + spacing
	}

	return targets
}
