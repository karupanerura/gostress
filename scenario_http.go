package gostress

import (
	"time"
)

type HttpScenario struct {
	Method     string
	Path       string
	Headers    map[string]string
	Content    interface{}
	BeforeRun  func(ScenarioState, *HttpScenario)
	OnComplete func(ScenarioState, *HttpResponse, time.Duration)
	OnError    func(ScenarioState, error)
}

func (scenario *HttpScenario) run(c *ScenarioContext) chan struct{} {
	ch := make(chan struct{}, 1)
	go func() {
		if cb := scenario.BeforeRun; cb != nil {
			cb(c.State, scenario)
		}
		startAt := time.Now()
		res, err := c.client.Request(scenario.Method, scenario.Path, scenario.Headers, scenario.Content)
		endAt := time.Now()
		if err == nil {
			if cb := scenario.OnComplete; cb != nil {
				duration := endAt.Sub(startAt)
				cb(c.State, res, duration)
			}
		} else {
			if cb := scenario.OnError; cb != nil {
				cb(c.State, err)
			} else {
				panic(err)
			}
		}
		ch <- struct{}{}
	}()
	return ch
}
