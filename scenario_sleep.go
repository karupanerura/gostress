package gostress

import (
	"time"
)

type SleepScenario struct {
	Duration   time.Duration
	OnComplete func(ScenarioState)
}

func (scenario *SleepScenario) run(c *ScenarioContext) <-chan done {
	ch := make(chan done, 1)
	go func() {
		time.Sleep(scenario.Duration)
		if cb := scenario.OnComplete; cb != nil {
			cb(c.State)
		}
		ch <- done{}
		close(ch)
	}()
	return ch
}
