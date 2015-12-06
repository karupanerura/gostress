package gostress

type NoopScenario struct {
}

func (scenario *NoopScenario) run(_ *ScenarioContext) <-chan done {
	ch := make(chan done, 1)
	ch <- done{}
	close(ch)
	return ch
}
