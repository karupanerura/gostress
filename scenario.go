package gostress

import "sync"

type ScenarioState interface{}

type Scenario interface {
	run(*ScenarioContext) chan struct{}
}

type ScenarioContext struct {
	client *HttpClient
	wg     *sync.WaitGroup
	State  ScenarioState
}

func NewScenarioContext(client *HttpClient, state ScenarioState) *ScenarioContext {
	return &ScenarioContext{
		client: client,
		wg:     &sync.WaitGroup{},
		State:  state,
	}
}

func (c *ScenarioContext) Run(scenario Scenario) {
	done := scenario.run(c)
	<-done
	close(done)
	c.wg.Wait()
}
