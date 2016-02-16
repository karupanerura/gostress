package gostress

type DeferScenario struct {
	Defer func(ScenarioState) Scenario
}

func (scenario *DeferScenario) run(c *ScenarioContext) chan struct{} {
	next := scenario.Defer(c.State)
	return next.run(c)
}
