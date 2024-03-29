package timer

import (
	"time"
)

const (
	DoorOpenTime     = 3 * time.Second
	decisionDeadline = 250 * time.Millisecond
	pollRate         = 20 * time.Millisecond
	PokeRate         = 500 * time.Millisecond
)

type Timer struct {
	isActive      bool
	endTime       time.Time
	timerDuration time.Duration
}

var (
	DoorTimer             = Timer{timerDuration: DoorOpenTime}
	DecisionDeadlineTimer = Timer{timerDuration: decisionDeadline}
	PokeCabTimer          = Timer{timerDuration: PokeRate}
)

func (timer *Timer) PollTimerOut(receiver chan<- bool) {
	prev := false
	for {
		time.Sleep(pollRate)
		v := timer.TimedOut()
		if v != prev && v {
			receiver <- v
		}
		prev = v
	}
}

func (timer *Timer) TimerStart() {
	timer.endTime = time.Now().Add(timer.timerDuration)
	timer.isActive = true
}

func (timer *Timer) TimerStop() {
	timer.isActive = false
}

func (timer *Timer) TimedOut() bool {
	return timer.isActive && time.Now().After(timer.endTime)
}
