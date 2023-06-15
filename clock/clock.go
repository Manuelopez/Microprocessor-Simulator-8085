package clock

import (
	"fmt"
	"sync"
	"time"
)

const (
	State0To1   = "0 To 1"
	State1To0   = "1 To 0"
	StateStable = "In State"
	StateOff    = "OFF"
)

type Clock struct {
	Frequency       float64 `json:"frequency"`
	State           int     `json:"state"`
	TransitionState string  `json:"transitionState"`
	Off             bool
	sync.Mutex
}

func New(frequency float64) Clock {
    return Clock{Frequency: frequency, State: 0, TransitionState: StateStable, Off: true}
}

func (c Clock) GetState() (int, string) {
	c.Mutex.Lock()

	defer c.Mutex.Unlock()
	return c.State, c.TransitionState
}

func (c Clock) GetFrequency() float64 {

	c.Mutex.Lock()

	defer c.Mutex.Unlock()
	return c.Frequency
}

func (c *Clock) SetFrequency(frequency float64) {
	c.Mutex.Lock()

	defer c.Mutex.Unlock()
	c.Frequency = frequency
}

func (c *Clock) TurnOn() {


	c.Mutex.Lock()
    c.Off = false
    c.Mutex.Unlock()
	timeGap := (1 / c.Frequency) * float64(time.Second)
	for {
        c.Mutex.Lock()
        if c.Off{
            break
        }
		fmt.Println(" is running ")
		if c.State == 0 {
			c.State = 1
			c.TransitionState = State0To1
			time.Sleep(time.Duration(timeGap))
		} else {
			c.State = 0
			c.TransitionState = State1To0
			time.Sleep(time.Duration(timeGap))
		}
		if c.TransitionState == StateOff {

			c.Mutex.Unlock()
			break
		}

		c.Mutex.Unlock()
		time.Sleep(time.Duration(timeGap))
	}
}

func (c *Clock) TurnOff() {

	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	c.Off = true
}

func (c *Clock) Wait() {
	c.Lock()

	timeGap := (1 / c.Frequency) * float64(time.Second)
	time.Sleep(time.Duration(timeGap))
	c.Unlock()
}
