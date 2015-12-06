package gostress

import (
	"time"
)

type DelayScenario struct {
	Duration   time.Duration
	Scenario   Scenario
	OnComplete func(ScenarioState)
}

func (scenario *DelayScenario) run(c *ScenarioContext) <-chan done {
	ch := make(chan done, 1)
	ch <- done{}
	close(ch)

	c.wg.Add(1)
	go func() {
		time.Sleep(scenario.Duration)
		ch := scenario.Scenario.run(c)
		<-ch
		if cb := scenario.OnComplete; cb != nil {
			cb(c.State)
		}
		c.wg.Done()
	}()
	return ch
}
