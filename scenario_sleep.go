package gostress

import (
	"time"
)

type SleepScenario struct {
	Duration   time.Duration
	OnComplete func(ScenarioState)
}

func (scenario *SleepScenario) run(c *ScenarioContext) chan struct{} {
	ch := make(chan struct{}, 1)
	go func() {
		time.Sleep(scenario.Duration)
		if cb := scenario.OnComplete; cb != nil {
			cb(c.State)
		}
		ch <- struct{}{}
	}()
	return ch
}
