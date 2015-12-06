package gostress

import (
	"sync"
	"time"
	"math/rand"
)

type ScenarioGroup struct {
	scenarios  []Scenario
	OnComplete func(ScenarioState, time.Duration)
}

type ConcurrentScenarioGroup struct {
	ScenarioGroup
}

type SeriesScenarioGroup struct {
	ScenarioGroup
	MaxInterval time.Duration
	MinInterval time.Duration
}

const ZERO_SEC = 0 * time.Second

func NewConcurrentScenarioGroup(size int) *ConcurrentScenarioGroup {
	return &ConcurrentScenarioGroup{
		ScenarioGroup{
			scenarios:  make([]Scenario, 0, size),
			OnComplete: nil,
		},
	}
}

func NewSeriesScenarioGroup(size int) *SeriesScenarioGroup {
	return &SeriesScenarioGroup{
		ScenarioGroup: ScenarioGroup{
			scenarios:  make([]Scenario, 0, size),
			OnComplete: nil,
		},
		MaxInterval: ZERO_SEC,
		MinInterval: ZERO_SEC,
	}
}

func (c *ConcurrentScenarioGroup) Add(scenario Scenario) *ConcurrentScenarioGroup {
	c.scenarios = append(c.scenarios, scenario)
	return c
}

func (c *ConcurrentScenarioGroup) AddNth(count uint, scenario Scenario) *ConcurrentScenarioGroup {
	for i := (uint)(0); i < count; i++ {
		c.Add(scenario)
	}
	return c
}

func (c *SeriesScenarioGroup) Next(scenario Scenario) *SeriesScenarioGroup {
	c.scenarios = append(c.scenarios, scenario)
	return c
}

func (group *ConcurrentScenarioGroup) run(c *ScenarioContext) <-chan done {
	ch := make(chan done, 1)
	wg := &sync.WaitGroup{}

	for _, scenario := range group.scenarios {
		scenario := scenario // redeclare c for the closure
		wg.Add(1)
		go func() {
			done := scenario.run(c)
			<-done
			wg.Done()
		}()
	}

	go func() {
		startAt := time.Now()
		wg.Wait()
		endAt := time.Now()
		if cb := group.OnComplete; cb != nil {
			duration := endAt.Sub(startAt)
			cb(c.State, duration)
		}
		ch <- done{}
		close(ch)
	}()

	return ch
}

func (group *SeriesScenarioGroup) run(c *ScenarioContext) <-chan done {
	ch := make(chan done, 1)

	go func() {
		startAt := time.Now()
		for _, scenario := range group.scenarios {
			if group.MaxInterval > ZERO_SEC {
				baseInterval := group.MinInterval
				rangedInterval := group.MaxInterval - group.MinInterval
				interval := baseInterval + time.Duration(rand.Int63n(rangedInterval.Nanoseconds()))
				time.Sleep(interval)
			}
			done := scenario.run(c)
			<-done
		}
		endAt := time.Now()
		if cb := group.OnComplete; cb != nil {
			duration := endAt.Sub(startAt)
			cb(c.State, duration)
		}
		ch <- done{}
		close(ch)
	}()

	return ch
}
