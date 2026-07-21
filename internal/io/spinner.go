package io

import (
	"fmt"
	"sync"
	"time"

	"go.vervstack.ru/verv/internal/io/colors"
)

const (
	spinnerInterval = 80 * time.Millisecond
	clearLine       = "\r\033[K"
)

var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// Spinner renders an animated, single-line progress indicator (à la `docker pull`)
// on top of an IO printer, then collapses it into a static done/fail line.
type Spinner struct {
	printer IO

	mu      sync.Mutex
	label   string
	running bool
	stopCh  chan struct{}
	doneCh  chan struct{}
}

func NewSpinner(printer IO) *Spinner {
	return &Spinner{
		printer: printer,
	}
}

// Start begins animating the spinner in place with the given label.
// It is a no-op if the spinner is already running.
func (s *Spinner) Start(label string) {
	s.mu.Lock()

	if s.running {
		s.mu.Unlock()

		return
	}

	s.running = true
	s.label = label
	s.stopCh = make(chan struct{})
	s.doneCh = make(chan struct{})

	s.mu.Unlock()

	go s.animate()
}

// Stop halts the animation and prints a final, static line summarizing the step.
func (s *Spinner) Stop(success bool, finalMsg string) {
	s.mu.Lock()

	if !s.running {
		s.mu.Unlock()

		return
	}

	s.running = false

	s.mu.Unlock()

	close(s.stopCh)
	<-s.doneCh

	mark := "✅"
	color := colors.ColorGreen

	if !success {
		mark = "❌"
		color = colors.ColorRed
	}

	s.printer.Print(clearLine)
	s.printer.PrintlnColored(color, fmt.Sprintf("%s %s", mark, finalMsg))
}

func (s *Spinner) animate() {
	defer close(s.doneCh)

	ticker := time.NewTicker(spinnerInterval)
	defer ticker.Stop()

	frame := 0

	for {
		select {
		case <-s.stopCh:
			return
		case <-ticker.C:
			s.mu.Lock()

			label := s.label

			s.mu.Unlock()

			line := clearLine + colors.TerminalColor(colors.ColorCyan) + spinnerFrames[frame] +
				colors.TerminalColor(colors.ColorDefault) + " " + label
			s.printer.Print(line)

			frame = (frame + 1) % len(spinnerFrames)
		}
	}
}
