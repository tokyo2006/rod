package rod

import (
	"sync"

	"github.com/ysmood/kit"
	"github.com/ysmood/rod/lib/cdp"
	"github.com/ysmood/rod/lib/input"
)

const defaultMouseButton = "left"

// Mouse represents the mouse on a page, it's always related the main frame
type Mouse struct {
	page *Page
	sync.Mutex

	x int64
	y int64

	// the buttons is currently beening pressed, reflects the press order
	buttons []string
}

// MoveE ...
func (m *Mouse) MoveE(x, y int64, steps int) error {
	if steps < 1 {
		steps = 1
	}

	m.Lock()
	defer m.Unlock()

	stepX := (x - m.x) / int64(steps)
	stepY := (y - m.y) / int64(steps)

	button, buttons := input.EncodeMouseButton(m.buttons)

	for i := 0; i < steps; i++ {
		toX := m.x + stepX
		toY := m.y + stepY

		_, err := m.page.Call("Input.dispatchMouseEvent", cdp.Object{
			"type":      "mouseMoved",
			"x":         toX,
			"y":         toY,
			"button":    button,
			"buttons":   buttons,
			"modifiers": m.page.Keyboard.modifiers,
		})
		if err != nil {
			return err
		}

		// to make sure set only when call is successful
		m.x = toX
		m.y = toY
	}

	return nil
}

// Move to the location
func (m *Mouse) Move(x, y int64) {
	kit.E(m.MoveE(x, y, 0))
}

// DownE ...
func (m *Mouse) DownE(button string, clicks int) error {
	m.Lock()
	defer m.Unlock()

	toButtons := append(m.buttons, button)

	_, buttons := input.EncodeMouseButton(toButtons)

	_, err := m.page.Call("Input.dispatchMouseEvent", cdp.Object{
		"type":       "mousePressed",
		"button":     button,
		"buttons":    buttons,
		"clickCount": clicks,
		"modifiers":  m.page.Keyboard.modifiers,
		"x":          m.x,
		"y":          m.y,
	})
	if err != nil {
		return err
	}
	m.buttons = toButtons
	return nil
}

// Down button: none, left, middle, right, back, forward
func (m *Mouse) Down(button string) {
	kit.E(m.DownE(button, 1))
}

// UpE ...
func (m *Mouse) UpE(button string, clicks int) error {
	m.Lock()
	defer m.Unlock()

	toButtons := []string{}
	for _, btn := range m.buttons {
		if btn == button {
			continue
		}
		toButtons = append(toButtons, btn)
	}

	_, buttons := input.EncodeMouseButton(toButtons)

	_, err := m.page.Call("Input.dispatchMouseEvent", cdp.Object{
		"type":       "mouseReleased",
		"button":     button,
		"buttons":    buttons,
		"clickCount": clicks,
		"x":          m.x,
		"y":          m.y,
	})
	if err != nil {
		return err
	}
	m.buttons = toButtons
	return nil
}

// Up button: none, left, middle, right, back, forward
func (m *Mouse) Up(button string) {
	kit.E(m.UpE(button, 1))
}

// ClickE ...
func (m *Mouse) ClickE(button string) error {
	if button == "" {
		button = defaultMouseButton
	}

	err := m.DownE(button, 1)
	if err != nil {
		return err
	}

	return m.UpE(button, 1)
}

// Click button: none, left, middle, right, back, forward
func (m *Mouse) Click(button string) {
	kit.E(m.ClickE(button))
}