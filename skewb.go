package skewb

import (
	"errors"
	"fmt"
	"image/png"
	"os"
	"strings"

	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/backend/softwarebackend"
)

type Skewber interface {
	Drawer
	MovesApplier
	CenterDowner
	Equaler
	ExactEqualer
	Mirrorer
	CornerColorsGetter
	CenterColorGetter
}

type Drawer interface {
	Draw(fileName string) error
}

type MovesApplier interface {
	ApplyWCAMoves(wcaMoves string) error
	ApplyRubiskewbMoves(rubiskewbMoves string) error
}

type WCAMovesApplier interface {
	ApplyWCAMoves(wcaMoves string) error
}

type RubiskewbMovesApplier interface {
	ApplyRubiskewbMoves(rubiskewbMoves string) error
}

type CenterDowner interface {
	CenterDown(color string) error
}

type Equaler interface {
	Equal(other Skewber) bool
}

type ExactEqualer interface {
	ExactEqual(other Skewber) bool
}

type Mirrorer interface {
	OneLayerMirrorer
	FullMirrorer
}

type OneLayerMirrorer interface {
	OneLayerMirror(other Skewber, layerColor string) bool
}

type FullMirrorer interface {
	FullMirror(other Skewber) bool
}

type CornerColorsGetter interface {
	GetUFRCornerColors() [3]string
	GetURBCornerColors() [3]string
	GetULFCornerColors() [3]string
	GetUBLCornerColors() [3]string
	GetDRFCornerColors() [3]string
	GetDBRCornerColors() [3]string
	GetDFLCornerColors() [3]string
	GetDLBCornerColors() [3]string
}

type CenterColorGetter interface {
	GetUpCenterColor() string
	GetFrontCenterColor() string
	GetRightCenterColor() string
	GetBackCenterColor() string
	GetLeftCenterColor() string
	GetDownCenterColor() string
}

type Skewb struct {
	ufr corner
	urb corner
	ulf corner
	ubl corner
	drf corner
	dbr corner
	dfl corner
	dlb corner

	up    center
	front center
	right center
	back  center
	left  center
	down  center
}

type corner struct {
	colors          CornerColors
	firstPositions  cornerPositions
	secondPositions cornerPositions
	thirdPositions  cornerPositions
}

type CornerColors struct {
	first  string
	second string
	third  string
}

type cornerPositions struct {
	starting   [2]float64
	firstLine  [2]float64
	secondLine [2]float64
}

type center struct {
	color     string
	positions centerPositions
}

type centerPositions struct {
	starting   [2]float64
	firstLine  [2]float64
	secondLine [2]float64
	thirdLine  [2]float64
}

type Move string

var (
	U            Move = "U"
	UPrime       Move = "U'"
	R            Move = "R"
	RPrime       Move = "R'"
	LittleR      Move = "r"
	LittleRPrime Move = "r'"
	B            Move = "B"
	BPrime       Move = "B'"
	LittleB      Move = "b"
	LittleBPrime Move = "b'"
	L            Move = "L"
	LPrime       Move = "L'"
	LittleL      Move = "l"
	LittleLPrime Move = "l'"
	F            Move = "F"
	FPrime       Move = "F'"
	LittleF      Move = "f"
	LittleFPrime Move = "f'"

	X      Move = "x"
	XPrime Move = "x'"
	X2     Move = "x2"
	Y      Move = "y"
	YPrime Move = "y'"
	Y2     Move = "y2"
	Z      Move = "z"
	ZPrime Move = "z'"
	Z2     Move = "z2"

	ErrWCAMove = errors.New("wca move is not supported; valid types are: \"U\", \"U'\" \"R\", \"R'\" \"B\", \"B'\" \"L\", \"L'\", \"x\", \"x'\", \"x2\", \"y\", \"y'\", \"y2\", \"z\", \"z'\", \"z2\"")
	ErrColor   = errors.New("color is not part of Skewb")

	clockwise, itsY, equal             = true, true, true
	counterClockwise, itsntY, notEqual = false, false, false

	layerCornerColor  = 0
	firstCornerColor  = 1
	secondCornerColor = 2
	otherCornerColor  = -1
)

func New(upColor, frontColor, rightColor, backColor, leftColor, downColor string) Skewb {
	return Skewb{
		ufr: corner{
			colors: CornerColors{
				first:  upColor,
				second: frontColor,
				third:  rightColor,
			},
			firstPositions: cornerPositions{
				starting:   [2]float64{180, 90},
				firstLine:  [2]float64{300, 90},
				secondLine: [2]float64{240, 120},
			},
			secondPositions: cornerPositions{
				starting:   [2]float64{180, 90},
				firstLine:  [2]float64{240, 120},
				secondLine: [2]float64{240, 195},
			},
			thirdPositions: cornerPositions{
				starting:   [2]float64{240, 120},
				firstLine:  [2]float64{300, 90},
				secondLine: [2]float64{240, 195},
			},
		},
		urb: corner{
			colors: CornerColors{
				first:  upColor,
				second: rightColor,
				third:  backColor,
			},
			firstPositions: cornerPositions{
				starting:   [2]float64{300, 30},
				firstLine:  [2]float64{360, 60},
				secondLine: [2]float64{300, 90},
			},
			secondPositions: cornerPositions{
				starting:   [2]float64{300, 90},
				firstLine:  [2]float64{360, 60},
				secondLine: [2]float64{360, 135},
			},
			thirdPositions: cornerPositions{
				starting:   [2]float64{360, 60},
				firstLine:  [2]float64{420, 30},
				secondLine: [2]float64{360, 135},
			},
		},
		ulf: corner{
			colors: CornerColors{
				first:  upColor,
				second: leftColor,
				third:  frontColor,
			},
			firstPositions: cornerPositions{
				starting:   [2]float64{120, 60},
				firstLine:  [2]float64{180, 30},
				secondLine: [2]float64{180, 90},
			},
			secondPositions: cornerPositions{
				starting:   [2]float64{60, 30},
				firstLine:  [2]float64{120, 135},
				secondLine: [2]float64{120, 60},
			},
			thirdPositions: cornerPositions{
				starting:   [2]float64{120, 60},
				firstLine:  [2]float64{180, 90},
				secondLine: [2]float64{120, 135},
			},
		},
		ubl: corner{
			colors: CornerColors{
				first:  upColor,
				second: backColor,
				third:  leftColor,
			},
			firstPositions: cornerPositions{
				starting:   [2]float64{180, 30},
				firstLine:  [2]float64{240, 0},
				secondLine: [2]float64{300, 30},
			},
			secondPositions: cornerPositions{
				starting:   [2]float64{420, 30},
				firstLine:  [2]float64{480, 0},
				secondLine: [2]float64{480, 75},
			},
			thirdPositions: cornerPositions{
				starting:   [2]float64{0, 0},
				firstLine:  [2]float64{60, 30},
				secondLine: [2]float64{0, 75},
			},
		},
		drf: corner{
			colors: CornerColors{
				first:  downColor,
				second: rightColor,
				third:  frontColor,
			},
			firstPositions: cornerPositions{
				starting:   [2]float64{180, 240},
				firstLine:  [2]float64{240, 270},
				secondLine: [2]float64{240, 345},
			},
			secondPositions: cornerPositions{
				starting:   [2]float64{240, 195},
				firstLine:  [2]float64{300, 240},
				secondLine: [2]float64{240, 270},
			},
			thirdPositions: cornerPositions{
				starting:   [2]float64{240, 195},
				firstLine:  [2]float64{240, 270},
				secondLine: [2]float64{180, 240},
			},
		},
		dbr: corner{
			colors: CornerColors{
				first:  downColor,
				second: backColor,
				third:  rightColor,
			},
			firstPositions: cornerPositions{
				starting:   [2]float64{240, 345},
				firstLine:  [2]float64{240, 420},
				secondLine: [2]float64{180, 390},
			},
			secondPositions: cornerPositions{
				starting:   [2]float64{360, 135},
				firstLine:  [2]float64{420, 180},
				secondLine: [2]float64{360, 210},
			},
			thirdPositions: cornerPositions{
				starting:   [2]float64{300, 240},
				firstLine:  [2]float64{360, 135},
				secondLine: [2]float64{360, 210},
			},
		},
		dfl: corner{
			colors: CornerColors{
				first:  downColor,
				second: frontColor,
				third:  leftColor,
			},
			firstPositions: cornerPositions{
				starting:   [2]float64{120, 210},
				firstLine:  [2]float64{180, 240},
				secondLine: [2]float64{120, 285},
			},
			secondPositions: cornerPositions{
				starting:   [2]float64{120, 135},
				firstLine:  [2]float64{180, 240},
				secondLine: [2]float64{120, 210},
			},
			thirdPositions: cornerPositions{
				starting:   [2]float64{120, 135},
				firstLine:  [2]float64{120, 210},
				secondLine: [2]float64{60, 180},
			},
		},
		dlb: corner{
			colors: CornerColors{
				first:  downColor,
				second: leftColor,
				third:  backColor,
			},
			firstPositions: cornerPositions{
				starting:   [2]float64{120, 285},
				firstLine:  [2]float64{180, 390},
				secondLine: [2]float64{120, 360},
			},
			secondPositions: cornerPositions{
				starting:   [2]float64{0, 75},
				firstLine:  [2]float64{60, 180},
				secondLine: [2]float64{0, 150},
			},
			thirdPositions: cornerPositions{
				starting:   [2]float64{420, 180},
				firstLine:  [2]float64{480, 75},
				secondLine: [2]float64{480, 150},
			},
		},

		up: center{
			color: upColor,
			positions: centerPositions{
				starting:   [2]float64{180, 90},
				firstLine:  [2]float64{180, 30},
				secondLine: [2]float64{300, 30},
				thirdLine:  [2]float64{300, 90},
			},
		},
		front: center{
			color: frontColor,
			positions: centerPositions{
				starting:   [2]float64{120, 135},
				firstLine:  [2]float64{180, 90},
				secondLine: [2]float64{240, 195},
				thirdLine:  [2]float64{180, 240},
			},
		},
		right: center{
			color: rightColor,
			positions: centerPositions{
				starting:   [2]float64{240, 195},
				firstLine:  [2]float64{300, 90},
				secondLine: [2]float64{360, 135},
				thirdLine:  [2]float64{300, 240},
			},
		},
		back: center{
			color: backColor,
			positions: centerPositions{
				starting:   [2]float64{360, 135},
				firstLine:  [2]float64{420, 30},
				secondLine: [2]float64{480, 75},
				thirdLine:  [2]float64{420, 180},
			},
		},
		left: center{
			color: leftColor,
			positions: centerPositions{
				starting:   [2]float64{0, 75},
				firstLine:  [2]float64{60, 30},
				secondLine: [2]float64{120, 135},
				thirdLine:  [2]float64{60, 180},
			},
		},
		down: center{
			color: downColor,
			positions: centerPositions{
				starting:   [2]float64{180, 240},
				firstLine:  [2]float64{240, 345},
				secondLine: [2]float64{180, 390},
				thirdLine:  [2]float64{120, 285},
			},
		},
	}
}

func (s *Skewb) Draw(fileName string) error {
	backend := softwarebackend.New(490, 430)
	cv := canvas.New(backend)
	image := cv.GetImageData(0, 0, 490, 430)

	// Positions for drawing: https://github.com/AnnikaStein/SkewbPage/blob/7ced702e91ed90de86f3020403c0c17ce484f4ac/SkewbSkills/skewbskillsscripts.js#L1750
	cv.Translate(10, 10)
	cv.SetStrokeStyle("#000000FF")
	cv.SetLineWidth(3.0)
	cv.SetLineJoin(canvas.Round)
	cv.SetLineCap(canvas.Round)

	s.ufr.draw(cv)
	s.urb.draw(cv)
	s.ulf.draw(cv)
	s.ubl.draw(cv)
	s.drf.draw(cv)
	s.dbr.draw(cv)
	s.dfl.draw(cv)
	s.dlb.draw(cv)
	s.up.draw(cv)

	s.front.draw(cv)
	s.right.draw(cv)
	s.back.draw(cv)
	s.left.draw(cv)
	s.down.draw(cv)

	cv.Translate(-10, -10)

	file, err := os.Create(fmt.Sprintf("%v.png", fileName))

	if err != nil {
		return err
	}

	err = png.Encode(file, image)

	if err != nil {
		return err
	}

	return nil
}

func (c *corner) draw(cv *canvas.Canvas) {
	c.drawFirstLayer(cv)
	c.drawSecondLayer(cv)
	c.drawThirdLayer(cv)
}

func (c *corner) drawFirstLayer(cv *canvas.Canvas) {
	cv.SetFillStyle(c.colors.first)
	cv.BeginPath()
	cv.MoveTo(c.firstPositions.starting[0], c.firstPositions.starting[1])
	cv.LineTo(c.firstPositions.firstLine[0], c.firstPositions.firstLine[1])
	cv.LineTo(c.firstPositions.secondLine[0], c.firstPositions.secondLine[1])
	cv.LineTo(c.firstPositions.starting[0], c.firstPositions.starting[1])
	cv.ClosePath()
	cv.Fill()
	cv.Stroke()
}

func (c *corner) drawSecondLayer(cv *canvas.Canvas) {
	cv.SetFillStyle(c.colors.second)
	cv.BeginPath()
	cv.MoveTo(c.secondPositions.starting[0], c.secondPositions.starting[1])
	cv.LineTo(c.secondPositions.firstLine[0], c.secondPositions.firstLine[1])
	cv.LineTo(c.secondPositions.secondLine[0], c.secondPositions.secondLine[1])
	cv.LineTo(c.secondPositions.starting[0], c.secondPositions.starting[1])
	cv.ClosePath()
	cv.Fill()
	cv.Stroke()
}

func (c *corner) drawThirdLayer(cv *canvas.Canvas) {
	cv.SetFillStyle(c.colors.third)
	cv.BeginPath()
	cv.MoveTo(c.thirdPositions.starting[0], c.thirdPositions.starting[1])
	cv.LineTo(c.thirdPositions.firstLine[0], c.thirdPositions.firstLine[1])
	cv.LineTo(c.thirdPositions.secondLine[0], c.thirdPositions.secondLine[1])
	cv.LineTo(c.thirdPositions.starting[0], c.thirdPositions.starting[1])
	cv.ClosePath()
	cv.Fill()
	cv.Stroke()
}

func (c *center) draw(cv *canvas.Canvas) {
	cv.SetFillStyle(c.color)
	cv.BeginPath()
	cv.MoveTo(c.positions.starting[0], c.positions.starting[1])
	cv.LineTo(c.positions.firstLine[0], c.positions.firstLine[1])
	cv.LineTo(c.positions.secondLine[0], c.positions.secondLine[1])
	cv.LineTo(c.positions.thirdLine[0], c.positions.thirdLine[1])
	cv.LineTo(c.positions.starting[0], c.positions.starting[1])
	cv.ClosePath()
	cv.Fill()
	cv.Stroke()
}

func (s *Skewb) ApplyWCAMoves(wcaMoves string) error {
	for _, m := range strings.Split(wcaMoves, " ") {
		var isClockwise bool
		move := Move(m)

		switch move {
		case U, R, B, L, X, X2, Y, Y2, Z, Z2:
			isClockwise = clockwise
		case UPrime, RPrime, BPrime, LPrime, XPrime, YPrime, ZPrime:
			isClockwise = counterClockwise
		default:
			return fmt.Errorf("%v %v", move, ErrWCAMove)
		}

		switch move {
		case U, UPrime:
			applyMove(&s.ubl, &s.dlb, &s.urb, &s.ulf, &s.up, &s.left, &s.back, isClockwise)
		case R, RPrime:
			applyMove(&s.dbr, &s.drf, &s.urb, &s.dlb, &s.right, &s.back, &s.down, isClockwise)
		case B, BPrime:
			applyMove(&s.dlb, &s.dbr, &s.ubl, &s.dfl, &s.back, &s.left, &s.down, isClockwise)
		case L, LPrime:
			applyMove(&s.dfl, &s.dlb, &s.ulf, &s.drf, &s.front, &s.down, &s.left, isClockwise)
		case X, XPrime, X2:
			applyRotation(&s.ufr, &s.urb, &s.dbr, &s.drf, &s.ulf, &s.ubl, &s.dlb, &s.dfl, &s.front, &s.up, &s.back, &s.down, isClockwise, move == X2, itsntY)
		case Y, YPrime, Y2:
			applyRotation(&s.ufr, &s.ulf, &s.ubl, &s.urb, &s.drf, &s.dfl, &s.dlb, &s.dbr, &s.front, &s.left, &s.back, &s.right, isClockwise, move == Y2, itsY)
		case Z, ZPrime, Z2:
			applyRotation(&s.ulf, &s.ufr, &s.drf, &s.dfl, &s.ubl, &s.urb, &s.dbr, &s.dlb, &s.up, &s.right, &s.down, &s.left, isClockwise, move == Z2, itsntY)
		}
	}

	return nil
}

func (s *Skewb) ApplyRubiskewbMoves(rubiskewbMoves string) error {
	for _, m := range strings.Split(rubiskewbMoves, " ") {
		var isClockwise bool
		move := Move(m)

		switch move {
		case R, LittleR, B, LittleB, L, LittleL, F, LittleF, X, X2, Y, Y2, Z, Z2:
			isClockwise = clockwise
		case RPrime, LittleRPrime, BPrime, LittleBPrime, LPrime, LittleLPrime, FPrime, LittleFPrime, XPrime, YPrime, ZPrime:
			isClockwise = counterClockwise
		default:
			return fmt.Errorf("%v %v", move, ErrWCAMove)
		}

		switch move {
		case R, RPrime:
			applyMove(&s.urb, &s.ufr, &s.ubl, &s.dbr, &s.right, &s.up, &s.back, isClockwise)
		case LittleR, LittleRPrime:
			applyMove(&s.dbr, &s.drf, &s.urb, &s.dlb, &s.right, &s.back, &s.down, isClockwise)
		case B, BPrime:
			applyMove(&s.ubl, &s.dlb, &s.urb, &s.ulf, &s.up, &s.left, &s.back, isClockwise)
		case LittleB, LittleBPrime:
			applyMove(&s.dlb, &s.dbr, &s.ubl, &s.dfl, &s.back, &s.left, &s.down, isClockwise)
		case L, LPrime:
			applyMove(&s.ulf, &s.ubl, &s.ufr, &s.dfl, &s.up, &s.front, &s.left, isClockwise)
		case LittleL, LittleLPrime:
			applyMove(&s.dfl, &s.dlb, &s.ulf, &s.drf, &s.front, &s.down, &s.left, isClockwise)
		case F, FPrime:
			applyMove(&s.ufr, &s.ulf, &s.urb, &s.drf, &s.front, &s.up, &s.right, isClockwise)
		case LittleF, LittleFPrime:
			applyMove(&s.drf, &s.dfl, &s.ufr, &s.dbr, &s.front, &s.right, &s.down, isClockwise)
		case X, XPrime, X2:
			applyRotation(&s.ufr, &s.urb, &s.dbr, &s.drf, &s.ulf, &s.ubl, &s.dlb, &s.dfl, &s.front, &s.up, &s.back, &s.down, isClockwise, move == X2, itsntY)
		case Y, YPrime, Y2:
			applyRotation(&s.ufr, &s.ulf, &s.ubl, &s.urb, &s.drf, &s.dfl, &s.dlb, &s.dbr, &s.front, &s.left, &s.back, &s.right, isClockwise, move == Y2, itsY)
		case Z, ZPrime, Z2:
			applyRotation(&s.ulf, &s.ufr, &s.drf, &s.dfl, &s.ubl, &s.urb, &s.dbr, &s.dlb, &s.up, &s.right, &s.down, &s.left, isClockwise, move == Z2, itsntY)
		}
	}

	return nil
}

func applyMove(rotationCenter, firstCorner, secondCorner, thirdCorner *corner, firstCenter, secondCenter, thirdCenter *center, isClockwise bool) {
	rotationCenter.rotate(isClockwise)
	firstCorner.rotate(!isClockwise)
	secondCorner.rotate(!isClockwise)
	thirdCorner.rotate(!isClockwise)

	rotateThreeCenters(firstCenter, secondCenter, thirdCenter, isClockwise)

	rotateThreeCorners(firstCorner, secondCorner, thirdCorner, isClockwise)
}

func rotateThreeCenters(first, second, third *center, isClockwise bool) {
	if isClockwise == clockwise {
		first.color, second.color, third.color = third.color, first.color, second.color
	} else {
		first.color, second.color, third.color = second.color, third.color, first.color
	}
}

func rotateThreeCorners(first, second, third *corner, isClockwise bool) {
	if isClockwise == clockwise {
		first.colors, second.colors, third.colors = third.colors, first.colors, second.colors
	} else {
		first.colors, second.colors, third.colors = second.colors, third.colors, first.colors
	}
}

func (c *corner) rotate(isClockwise bool) {
	if isClockwise == clockwise {
		c.colors.first, c.colors.second, c.colors.third = c.colors.second, c.colors.third, c.colors.first
	} else {
		c.colors.first, c.colors.second, c.colors.third = c.colors.third, c.colors.first, c.colors.second
	}
}

func applyRotation(firstCorner, secondCorner, thirdCorner, fourthCorner, fifthCorner, sixthCorner, seventhCorner, eighthCorner *corner, firstCenter, secondCenter, thirdCenter, fourthCenter *center, isClockwise, isDouble, isY bool) {
	rotateFourCenters(firstCenter, secondCenter, thirdCenter, fourthCenter, isClockwise)

	rotateFourCorners(firstCorner, secondCorner, thirdCorner, fourthCorner, isClockwise)
	rotateFourCorners(fifthCorner, sixthCorner, seventhCorner, eighthCorner, isClockwise)

	if !isY {
		firstCorner.rotate(counterClockwise)
		secondCorner.rotate(clockwise)
		thirdCorner.rotate(counterClockwise)
		fourthCorner.rotate(clockwise)
		fifthCorner.rotate(clockwise)
		sixthCorner.rotate(counterClockwise)
		seventhCorner.rotate(clockwise)
		eighthCorner.rotate(counterClockwise)
	}

	if isDouble {
		rotateFourCenters(firstCenter, secondCenter, thirdCenter, fourthCenter, isClockwise)

		rotateFourCorners(firstCorner, secondCorner, thirdCorner, fourthCorner, isClockwise)
		rotateFourCorners(fifthCorner, sixthCorner, seventhCorner, eighthCorner, isClockwise)

		if !isY {
			firstCorner.rotate(counterClockwise)
			secondCorner.rotate(clockwise)
			thirdCorner.rotate(counterClockwise)
			fourthCorner.rotate(clockwise)
			fifthCorner.rotate(clockwise)
			sixthCorner.rotate(counterClockwise)
			seventhCorner.rotate(clockwise)
			eighthCorner.rotate(counterClockwise)
		}
	}
}

func rotateFourCenters(first, second, third, fourth *center, isClockwise bool) {
	if isClockwise == clockwise {
		first.color, second.color, third.color, fourth.color = fourth.color, first.color, second.color, third.color
	} else {
		first.color, second.color, third.color, fourth.color = second.color, third.color, fourth.color, first.color
	}
}

func rotateFourCorners(first, second, third, fourth *corner, isClockwise bool) {
	if isClockwise == clockwise {
		first.colors, second.colors, third.colors, fourth.colors = fourth.colors, first.colors, second.colors, third.colors
	} else {
		first.colors, second.colors, third.colors, fourth.colors = second.colors, third.colors, fourth.colors, first.colors
	}
}

func (s *Skewb) CenterDown(color string) error {
	switch {
	case s.up.color == color:
		return s.ApplyWCAMoves(fmt.Sprintf("%v", X2))
	case s.front.color == color:
		return s.ApplyWCAMoves(fmt.Sprintf("%v", XPrime))
	case s.right.color == color:
		return s.ApplyWCAMoves(fmt.Sprintf("%v", Z))
	case s.back.color == color:
		return s.ApplyWCAMoves(fmt.Sprintf("%v", X))
	case s.left.color == color:
		return s.ApplyWCAMoves(fmt.Sprintf("%v", ZPrime))
	case s.down.color == color:
		return nil
	default:
		return fmt.Errorf("%v %v", color, ErrColor)
	}
}

func (s *Skewb) Equal(other Skewber) bool {
	other.CenterDown(s.down.color)

	if s.equal(other) {
		return equal
	}

	for range 3 {
		other.ApplyWCAMoves(fmt.Sprintf("%v", Y))

		if s.equal(other) {
			return equal
		}
	}

	return notEqual
}

func (s *Skewb) ExactEqual(other Skewber) bool {
	return s.equal(other)
}

func (s *Skewb) equal(other Skewber) bool {
	switch {
	case [3]string{s.ufr.colors.first, s.ufr.colors.second, s.ufr.colors.third} != other.GetUFRCornerColors():
		return notEqual
	case [3]string{s.urb.colors.first, s.urb.colors.second, s.urb.colors.third} != other.GetURBCornerColors():
		return notEqual
	case [3]string{s.ulf.colors.first, s.ulf.colors.second, s.ulf.colors.third} != other.GetULFCornerColors():
		return notEqual
	case [3]string{s.ubl.colors.first, s.ubl.colors.second, s.ubl.colors.third} != other.GetUBLCornerColors():
		return notEqual
	case [3]string{s.drf.colors.first, s.drf.colors.second, s.drf.colors.third} != other.GetDRFCornerColors():
		return notEqual
	case [3]string{s.dbr.colors.first, s.dbr.colors.second, s.dbr.colors.third} != other.GetDBRCornerColors():
		return notEqual
	case [3]string{s.dfl.colors.first, s.dfl.colors.second, s.dfl.colors.third} != other.GetDFLCornerColors():
		return notEqual
	case [3]string{s.dlb.colors.first, s.dlb.colors.second, s.dlb.colors.third} != other.GetDLBCornerColors():
		return notEqual

	case s.up.color != other.GetUpCenterColor():
		return notEqual
	case s.front.color != other.GetFrontCenterColor():
		return notEqual
	case s.right.color != other.GetRightCenterColor():
		return notEqual
	case s.back.color != other.GetBackCenterColor():
		return notEqual
	case s.left.color != other.GetLeftCenterColor():
		return notEqual
	case s.down.color != other.GetDownCenterColor():
		return notEqual
	}

	return equal
}

func (s *Skewb) OneLayerMirror(other Skewber, layerColor string) bool {
	s.CenterDown(layerColor)
	other.CenterDown(s.down.color)

	if s.oneLayerMirror(other, layerColor) {
		return equal
	}

	rotations := []string{"y", "y'", "y2"}

	for _, rotation := range rotations {
		other.ApplyWCAMoves(rotation)

		if s.oneLayerMirror(other, layerColor) {
			reverse := createReverse(rotation)
			other.ApplyWCAMoves(reverse)

			return equal
		}

		reverse := createReverse(rotation)
		other.ApplyWCAMoves(reverse)
	}

	return notEqual
}

func (s *Skewb) oneLayerMirror(other Skewber, layerColor string) bool {
	sLayerCorners := [][3]string{}
	otherLayerCorners := [][3]string{}

	if id := index(s.GetUFRCornerColors(), layerColor); id != -1 {
		sLayerCorners = append(sLayerCorners, s.GetUFRCornerColors())
		otherLayerCorners = append(otherLayerCorners, other.GetUFRCornerColors())
	}

	if id := index(s.GetURBCornerColors(), layerColor); id != -1 {
		sLayerCorners = append(sLayerCorners, s.GetURBCornerColors())
		otherLayerCorners = append(otherLayerCorners, other.GetURBCornerColors())
	}

	if id := index(s.GetULFCornerColors(), layerColor); id != -1 {
		sLayerCorners = append(sLayerCorners, s.GetULFCornerColors())
		otherLayerCorners = append(otherLayerCorners, other.GetULFCornerColors())
	}

	if id := index(s.GetUBLCornerColors(), layerColor); id != -1 {
		sLayerCorners = append(sLayerCorners, s.GetUBLCornerColors())
		otherLayerCorners = append(otherLayerCorners, other.GetUBLCornerColors())
	}

	if id := index(s.GetDRFCornerColors(), layerColor); id != -1 {
		sLayerCorners = append(sLayerCorners, s.GetDRFCornerColors())
		otherLayerCorners = append(otherLayerCorners, other.GetDRFCornerColors())
	}

	if id := index(s.GetDBRCornerColors(), layerColor); id != -1 {
		sLayerCorners = append(sLayerCorners, s.GetDBRCornerColors())
		otherLayerCorners = append(otherLayerCorners, other.GetDBRCornerColors())
	}

	if id := index(s.GetDFLCornerColors(), layerColor); id != -1 {
		sLayerCorners = append(sLayerCorners, s.GetDFLCornerColors())
		otherLayerCorners = append(otherLayerCorners, other.GetDFLCornerColors())
	}

	if id := index(s.GetDLBCornerColors(), layerColor); id != -1 {
		sLayerCorners = append(sLayerCorners, s.GetDLBCornerColors())
		otherLayerCorners = append(otherLayerCorners, other.GetDLBCornerColors())
	}

	sLayerColors := getLayerColor(sLayerCorners, layerColor)
	otherLayerColors := getLayerColor(otherLayerCorners, layerColor)

	if sLayerColors == otherLayerColors {
		return equal
	}

	return notEqual
}

func getLayerColor(layerCorners [][3]string, layerColor string) [4][3]int {
	layerColors := [4][3]int{}

	switch {
	case layerCorners[0][0] == layerColor:
		layerColors[0][0] = layerCornerColor
		layerColors[0][1] = firstCornerColor
		layerColors[0][2] = secondCornerColor
	case layerCorners[0][1] == layerColor:
		layerColors[0][0] = firstCornerColor
		layerColors[0][1] = layerCornerColor
		layerColors[0][2] = secondCornerColor
	case layerCorners[0][2] == layerColor:
		layerColors[0][0] = firstCornerColor
		layerColors[0][1] = secondCornerColor
		layerColors[0][2] = layerCornerColor
	}

	for i, corner := range layerCorners[1:] {
		for j, color := range corner {
			switch {
			case layerCorners[0][0] == color:
				layerColors[i+1][j] = layerColors[0][0]
			case layerCorners[0][1] == color:
				layerColors[i+1][j] = layerColors[0][1]
			case layerCorners[0][2] == color:
				layerColors[i+1][j] = layerColors[0][2]
			default:
				layerColors[i+1][j] = otherCornerColor
			}
		}
	}

	return layerColors
}

func index(array [3]string, element string) int {
	for i, value := range array {
		if value == element {
			return i
		}
	}

	return -1
}

func (s *Skewb) FullMirror(other Skewber) bool {
	if s.fullMirror(other) {
		return equal
	}

	rotations := []string{"x", "x'", "x2", "y", "y'", "y2", "z", "z'", "z2", "x y", "x y'", "x y2", "x z", "x z'", "x z2", "x' y", "x' y'", "x' z", "x' z'", "x2 y", "x2 y'", "x2 z", "x2 z'"}

	for _, rotation := range rotations {
		other.ApplyWCAMoves(rotation)

		if s.fullMirror(other) {
			reverse := createReverse(rotation)
			other.ApplyWCAMoves(reverse)

			return equal
		}

		reverse := createReverse(rotation)
		other.ApplyWCAMoves(reverse)
	}

	return notEqual
}

func (s *Skewb) fullMirror(other Skewber) bool {
	sRelativeColors := getRelativeColors(s)
	otherRelativeColors := getRelativeColors(other)

	if sRelativeColors == otherRelativeColors {
		return equal
	}

	return notEqual
}

func getRelativeColors(s Skewber) [6][5]int {
	relativeColors := [6][5]int{}

	upColor := s.GetUpCenterColor()
	relativeColors[0][0] = 0
	frontColor := s.GetFrontCenterColor()
	relativeColors[1][0] = 1
	rightColor := s.GetRightCenterColor()
	relativeColors[2][0] = 2
	backColor := s.GetBackCenterColor()
	relativeColors[3][0] = 3
	leftColor := s.GetLeftCenterColor()
	relativeColors[4][0] = 4
	downColor := s.GetDownCenterColor()
	relativeColors[5][0] = 5

	ufrRelativeColors := getCornerRelativeColors(s.GetUFRCornerColors(), upColor, frontColor, rightColor, backColor, leftColor, downColor)
	relativeColors[0][1] = ufrRelativeColors[0]
	relativeColors[1][1] = ufrRelativeColors[1]
	relativeColors[2][1] = ufrRelativeColors[2]
	urbRelativeColors := getCornerRelativeColors(s.GetURBCornerColors(), upColor, frontColor, rightColor, backColor, leftColor, downColor)
	relativeColors[0][2] = urbRelativeColors[0]
	relativeColors[2][2] = urbRelativeColors[1]
	relativeColors[3][1] = urbRelativeColors[2]
	ulfRelativeColors := getCornerRelativeColors(s.GetULFCornerColors(), upColor, frontColor, rightColor, backColor, leftColor, downColor)
	relativeColors[0][3] = ulfRelativeColors[0]
	relativeColors[4][1] = ulfRelativeColors[1]
	relativeColors[1][2] = ulfRelativeColors[2]
	ublRelativeColors := getCornerRelativeColors(s.GetUBLCornerColors(), upColor, frontColor, rightColor, backColor, leftColor, downColor)
	relativeColors[0][4] = ublRelativeColors[0]
	relativeColors[3][2] = ublRelativeColors[1]
	relativeColors[4][2] = ublRelativeColors[2]
	drfRelativeColors := getCornerRelativeColors(s.GetDRFCornerColors(), upColor, frontColor, rightColor, backColor, leftColor, downColor)
	relativeColors[5][1] = drfRelativeColors[0]
	relativeColors[2][3] = drfRelativeColors[1]
	relativeColors[1][3] = drfRelativeColors[2]
	dbrRelativeColors := getCornerRelativeColors(s.GetDBRCornerColors(), upColor, frontColor, rightColor, backColor, leftColor, downColor)
	relativeColors[5][2] = dbrRelativeColors[0]
	relativeColors[3][3] = dbrRelativeColors[1]
	relativeColors[2][4] = dbrRelativeColors[2]
	dflRelativeColors := getCornerRelativeColors(s.GetDFLCornerColors(), upColor, frontColor, rightColor, backColor, leftColor, downColor)
	relativeColors[5][3] = dflRelativeColors[0]
	relativeColors[1][4] = dflRelativeColors[1]
	relativeColors[4][3] = dflRelativeColors[2]
	dlbRelativeColors := getCornerRelativeColors(s.GetDLBCornerColors(), upColor, frontColor, rightColor, backColor, leftColor, downColor)
	relativeColors[5][4] = dlbRelativeColors[0]
	relativeColors[4][4] = dlbRelativeColors[1]
	relativeColors[3][4] = dlbRelativeColors[2]

	return relativeColors
}

func getCornerRelativeColors(colors [3]string, upColor, frontColor, rightColor, backColor, leftColor, downColor string) [3]int {
	relativeColors := [3]int{}

	for i := range 3 {
		switch colors[i] {
		case upColor:
			relativeColors[i] = 0
		case frontColor:
			relativeColors[i] = 1
		case rightColor:
			relativeColors[i] = 2
		case backColor:
			relativeColors[i] = 3
		case leftColor:
			relativeColors[i] = 4
		case downColor:
			relativeColors[i] = 5
		}
	}

	return relativeColors
}

func createReverse(moves string) string {
	reverse := ""

	for _, move := range strings.Split(moves, " ") {
		if strings.Contains(move, "'") {
			move = strings.ReplaceAll(move, "'", "")
		} else {
			if !strings.Contains(move, "2") {
				move += "'"
			}
		}

		reverse = fmt.Sprintf("%v %v", move, reverse)
	}

	return reverse
}

func (s *Skewb) GetUFRCornerColors() [3]string {
	return [3]string{s.ufr.colors.first, s.ufr.colors.second, s.ufr.colors.third}
}

func (s *Skewb) GetURBCornerColors() [3]string {
	return [3]string{s.urb.colors.first, s.urb.colors.second, s.urb.colors.third}
}

func (s *Skewb) GetULFCornerColors() [3]string {
	return [3]string{s.ulf.colors.first, s.ulf.colors.second, s.ulf.colors.third}
}

func (s *Skewb) GetUBLCornerColors() [3]string {
	return [3]string{s.ubl.colors.first, s.ubl.colors.second, s.ubl.colors.third}
}

func (s *Skewb) GetDRFCornerColors() [3]string {
	return [3]string{s.drf.colors.first, s.drf.colors.second, s.drf.colors.third}
}

func (s *Skewb) GetDBRCornerColors() [3]string {
	return [3]string{s.dbr.colors.first, s.dbr.colors.second, s.dbr.colors.third}
}

func (s *Skewb) GetDFLCornerColors() [3]string {
	return [3]string{s.dfl.colors.first, s.dfl.colors.second, s.dfl.colors.third}
}

func (s *Skewb) GetDLBCornerColors() [3]string {
	return [3]string{s.dlb.colors.first, s.dlb.colors.second, s.dlb.colors.third}
}

func (s *Skewb) GetUpCenterColor() string {
	return s.up.color
}

func (s *Skewb) GetFrontCenterColor() string {
	return s.front.color
}

func (s *Skewb) GetRightCenterColor() string {
	return s.right.color
}

func (s *Skewb) GetBackCenterColor() string {
	return s.back.color
}

func (s *Skewb) GetLeftCenterColor() string {
	return s.left.color
}

func (s *Skewb) GetDownCenterColor() string {
	return s.down.color
}
