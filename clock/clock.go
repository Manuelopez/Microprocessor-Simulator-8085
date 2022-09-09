package clock

import (
	"time"
)

const (
	State0To1 = "0 To 1"
	State1To0 = "1 To 0"
	StateStable     = "In State"
)

type Clock struct {
	Frequency       float64
	State           int
	TransitionState string
  Stream chan string
}

func New(frequency float64) Clock {
  return Clock{Frequency: frequency, State: 0, TransitionState: StateStable, Stream: make(chan string)}
}

func (c Clock) GetState() (int, string) {
	return c.State, c.TransitionState
}

func (c Clock) GetFrequency() float64 {
	return c.Frequency
}

func (c *Clock) SetFrequency(frequency float64) {
	c.Frequency = frequency
}

func (c *Clock) TurnOn() {
	timeGap := (1 / c.Frequency) * float64(time.Second)
	for {
		if c.State == 0 {
			c.State = 1
			c.TransitionState = State0To1
			c.Stream <- c.TransitionState
			time.Sleep(time.Duration(timeGap))
		} else {
			c.State = 0
			c.TransitionState = State1To0
			c.Stream <- c.TransitionState
			time.Sleep(time.Duration(timeGap))
		}

		time.Sleep(time.Duration(timeGap))
		c.TransitionState = StateStable
		c.Stream <- c.TransitionState
	}
}

func (c *Clock) TurnOff() {
	close(c.Stream)
}

func (c *Clock) Wait() {
	for state := range c.Stream {
		if state == State0To1{
			return
		}
	}
}
