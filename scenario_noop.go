package gostress

type NoopScenario struct {
}

func (scenario *NoopScenario) run(_ *ScenarioContext) chan struct{} {
	ch := make(chan struct{}, 1)
	ch <- struct{}{}
	return ch
}
