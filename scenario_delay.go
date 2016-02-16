package gostress

import (
	"time"
)

type DelayScenario struct {
	Duration   time.Duration
	Scenario   Scenario
	OnComplete func(ScenarioState)
}

func (scenario *DelayScenario) run(c *ScenarioContext) chan struct{} {
	ch := make(chan struct{}, 1)
	ch <- struct{}{}

	c.wg.Add(1)
	go func() {
		time.Sleep(scenario.Duration)
		ch := scenario.Scenario.run(c)
		<-ch
		close(ch)
		if cb := scenario.OnComplete; cb != nil {
			cb(c.State)
		}
		c.wg.Done()
	}()
	return ch
}
