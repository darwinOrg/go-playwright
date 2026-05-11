package extpw

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	dgctx "github.com/darwinOrg/go-common/context"
	dglogger "github.com/darwinOrg/go-logger"
)

const (
	minMouseSteps = 8
	maxMouseSteps = 50
)

// HumanLikeMouseOption configures human-like mouse movement behavior.
type HumanLikeMouseOption func(*humanLikeMouseOptions)

type humanLikeMouseOptions struct {
	speed float64
}

func defaultHumanLikeMouseOptions() *humanLikeMouseOptions {
	return &humanLikeMouseOptions{
		speed: 0.9,
	}
}

// WithMouseSpeed sets the movement speed (0.0~1.0, higher = faster).
// Default is 0.9.
func WithMouseSpeed(speed float64) HumanLikeMouseOption {
	return func(opts *humanLikeMouseOptions) {
		if speed > 0 && speed <= 1.0 {
			opts.speed = speed
		}
	}
}

// HumanLikeMoveToSelector moves the mouse to the element specified by selector
// with a natural human-like trajectory (Bezier curve + easing + micro-jitter).
// It does NOT click — use HumanLikeClick if a click is also needed.
func (p *ExtPage) HumanLikeMoveToSelector(ctx *dgctx.DgContext, selector string, options ...HumanLikeMouseOption) error {
	defer func() {
		if err := recover(); err != nil {
			dglogger.Errorf(ctx, "HumanLikeMoveToSelector[%s] panic: %v", selector, err)
		}
	}()

	p.CheckSuspend(ctx)

	opts := defaultHumanLikeMouseOptions()
	for _, o := range options {
		o(opts)
	}

	locator := p.ExtLocator(selector)
	if !locator.Exists(ctx) {
		err := fmt.Errorf("selector[%s] not found", selector)
		dglogger.Errorf(ctx, "HumanLikeMoveToSelector: %v", err)
		return err
	}

	bbox, err := locator.BoundingBox()
	if err != nil {
		dglogger.Errorf(ctx, "HumanLikeMoveToSelector: get bounding box error: %v", err)
		return err
	}
	if bbox == nil {
		err := fmt.Errorf("bounding box is nil for selector[%s]", selector)
		dglogger.Errorf(ctx, "HumanLikeMoveToSelector: %v", err)
		return err
	}

	// Target: random point within the central 60% of the element
	targetX := bbox.X + bbox.Width*0.2 + bbox.Width*0.6*rand.Float64()
	targetY := bbox.Y + bbox.Height*0.2 + bbox.Height*0.6*rand.Float64()

	// Determine starting position
	startX, startY, err := p.getMouseStartPosition()
	if err != nil {
		// Fallback: viewport bottom-right with randomness
		viewportSize := p.ViewportSize()
		if viewportSize != nil {
			startX = float64(viewportSize.Width) * 0.85
			startY = float64(viewportSize.Height) * 0.9
		} else {
			startX = 800
			startY = 600
		}
	}

	// Add jitter to start position so it doesn't always begin from the exact same spot
	startX += float64(rand.Intn(80) - 40)
	startY += float64(rand.Intn(50) - 25)

	return p.animateHumanLikeMouse(startX, startY, targetX, targetY, opts.speed)
}

// HumanLikeClick moves the mouse naturally to the element first, then clicks.
// This is a drop-in replacement for a "stealth" click that bypasses bot detection.
func (p *ExtPage) HumanLikeClick(ctx *dgctx.DgContext, selector string, options ...HumanLikeMouseOption) error {
	defer func() {
		if err := recover(); err != nil {
			dglogger.Errorf(ctx, "HumanLikeClick[%s] panic: %v", selector, err)
		}
	}()

	p.CheckSuspend(ctx)

	if err := p.HumanLikeMoveToSelector(ctx, selector, options...); err != nil {
		return err
	}

	// Tiny hesitation before pressing (humans naturally pause before clicking)
	time.Sleep(time.Duration(rand.Intn(120)+40) * time.Millisecond)

	bbox, err := p.ExtLocator(selector).BoundingBox()
	if err != nil || bbox == nil {
		dglogger.Errorf(ctx, "HumanLikeClick: re-get bounding box failed, fallback to regular click: %v", err)
		return p.Click(ctx, selector)
	}

	clickX := bbox.X + bbox.Width*0.2 + bbox.Width*0.6*rand.Float64()
	clickY := bbox.Y + bbox.Height*0.2 + bbox.Height*0.6*rand.Float64()

	if err := p.Mouse().Click(clickX, clickY); err != nil {
		dglogger.Errorf(ctx, "HumanLikeClick: mouse click error: %v", err)
		return err
	}

	return nil
}

// ---- internal helpers ----

// getMouseStartPosition reads the last-known mouse position stored in the page,
// or falls back to a reasonable default (viewport bottom-right).
func (p *ExtPage) getMouseStartPosition() (float64, float64, error) {
	result, err := p.Evaluate(`() => {
		const x = window._qwenLastMouseX;
		const y = window._qwenLastMouseY;
		if (x != null && y != null) return { x, y };
		return { x: window.innerWidth * 0.85, y: window.innerHeight * 0.9 };
	}`)
	if err != nil {
		return 0, 0, err
	}

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		return 0, 0, fmt.Errorf("unexpected mouse position result type: %T", result)
	}

	x, _ := resultMap["x"].(float64)
	y, _ := resultMap["y"].(float64)
	return x, y, nil
}

// animateHumanLikeMouse moves the mouse from (startX,startY) to (targetX,targetY)
// along a Bezier curve path with varying speed and micro-jitter.
func (p *ExtPage) animateHumanLikeMouse(startX, startY, targetX, targetY, speed float64) error {
	dist := math.Sqrt((targetX-startX)*(targetX-startX) + (targetY-startY)*(targetY-startY))

	// More steps = slower & smoother; less steps = faster & jerkier
	numSteps := int(dist / 18.0 * (1.0 + (1.0-speed)*0.6))
	switch {
	case numSteps < minMouseSteps:
		numSteps = minMouseSteps
	case numSteps > maxMouseSteps:
		numSteps = maxMouseSteps
	}

	path := generateBezierPath(startX, startY, targetX, targetY, numSteps)

	// Simulate human reaction delay before starting to move
	time.Sleep(time.Duration(rand.Intn(250)+50) * time.Millisecond)

	lastX, lastY := startX, startY
	for i, pt := range path {
		if err := p.Mouse().Move(pt[0], pt[1]); err != nil {
			return err
		}
		lastX, lastY = pt[0], pt[1]

		if i < len(path)-1 {
			progress := float64(i) / float64(len(path)-1)
			var delayMs int
			switch {
			case progress < 0.12:
				delayMs = rand.Intn(18) + 12 // start: slow (acceleration phase)
			case progress > 0.88:
				delayMs = rand.Intn(18) + 12 // end: slow (deceleration phase)
			default:
				delayMs = rand.Intn(8) + 3 // middle: fast cruising
			}
			time.Sleep(time.Duration(delayMs) * time.Millisecond)
		}
	}

	// Store last position in the page for the next movement
	_, _ = p.Evaluate(fmt.Sprintf(
		`() => { window._qwenLastMouseX = %f; window._qwenLastMouseY = %f; }`,
		lastX, lastY,
	))

	return nil
}

// ---- Bezier curve & easing ----

// cubicBezierPoint evaluates a cubic Bezier curve at parameter t ∈ [0,1].
func cubicBezierPoint(t, x0, y0, x1, y1, x2, y2, x3, y3 float64) (float64, float64) {
	u := 1 - t
	uu := u * u
	uuu := uu * u
	tt := t * t
	ttt := tt * t

	x := uuu*x0 + 3*uu*t*x1 + 3*u*tt*x2 + ttt*x3
	y := uuu*y0 + 3*uu*t*y1 + 3*u*tt*y2 + ttt*y3
	return x, y
}

// generateBezierPath generates waypoints along a cubic Bezier curve.
// Control points are perturbed randomly to create a natural-looking curved trajectory.
func generateBezierPath(startX, startY, endX, endY float64, numSteps int) [][2]float64 {
	dist := math.Sqrt((endX-startX)*(endX-startX) + (endY-startY)*(endY-startY))

	// Base control points at 25% and 75% along the direct line
	cp1x := startX + (endX-startX)*0.25
	cp1y := startY + (endY-startY)*0.25
	cp2x := startX + (endX-startX)*0.75
	cp2y := startY + (endY-startY)*0.75

	// Add perpendicular offsets for a curved trajectory
	if dist > 50 {
		offsetMag := dist * 0.12 * (rand.Float64()*0.6 + 0.2)
		dx := endX - startX
		dy := endY - startY
		lenD := math.Sqrt(dx*dx + dy*dy)
		if lenD > 0 {
			nx := -dy / lenD * offsetMag
			ny := dx / lenD * offsetMag
			// Randomly apply offset to one of the two control points
			if rand.Float64() < 0.5 {
				cp1x += nx
				cp1y += ny
			} else {
				cp2x += nx
				cp2y += ny
			}
		}
	}

	// Jitter control points themselves
	jitterAmt := dist * 0.06
	cp1x += (rand.Float64()*2 - 1) * jitterAmt
	cp1y += (rand.Float64()*2 - 1) * jitterAmt
	cp2x += (rand.Float64()*2 - 1) * jitterAmt
	cp2y += (rand.Float64()*2 - 1) * jitterAmt

	path := make([][2]float64, numSteps)
	for i := 0; i < numSteps; i++ {
		t := float64(i) / float64(numSteps-1)
		// Easing applied to the parameter t for more natural speed profile
		t = easeInOutQuad(t)

		x, y := cubicBezierPoint(t, startX, startY, cp1x, cp1y, cp2x, cp2y, endX, endY)

		// Micro-jitter: slight hand tremble at each step
		if dist > 20 {
			jitter := math.Max(dist*0.004, 0.8)
			x += (rand.Float64()*2 - 1) * jitter
			y += (rand.Float64()*2 - 1) * jitter
		}

		path[i] = [2]float64{x, y}
	}
	return path
}

// easeInOutQuad applies quadratic ease-in-out: slow → fast → slow.
func easeInOutQuad(t float64) float64 {
	if t < 0.5 {
		return 2 * t * t
	}
	return 1 - math.Pow(-2*t+2, 2)/2
}
