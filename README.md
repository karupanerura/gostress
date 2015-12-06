# gostress

HTTP/HTTPS stress test framework.

```go
package main

import (
	"github.com/karupanerura/gostress"
	"log"
	"math/rand"
	"runtime"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	rand.Seed(time.Now().UnixNano())

	client := gostress.NewHttpClient(
		gostress.HttpClientConfig{
			Server: gostress.ServerConfig{
				Hostname: "myhost.com",
				Secure: false,
			},
			Headers: map[string]string{},
			UserAgent: "Gostress/alpha",
			MaxIdleConnsPerHost: 1024,
			RequestEncoder: &gostress.JsonRequestEncoder{},
			ResponseDecoder: &gostress.JsonResponseDecoder{},
		},
	)
	state := map[string]string{}
	context := gostress.NewScenarioContext(client, state)
	scenario := makeScenario()
	context.Run(scenario)
}

func makeScenario() gostress.Scenario {
	scenarios := gostress.NewSeriesScenarioGroup(256)
	scenarios.MinInterval =   500 * time.Millisecond
	scenarios.MaxInterval = 10000 * time.Millisecond
	scenarios.Next(makeHTTPScenario("GET", "/", nil))
	scenarios.Next(
		&gostress.SleepScenario{Duration: 1 * time.Millisecond},
	)
	scenarios.Next(
		gostress.NewConcurrentScenarioGroup(3).Add(
			makeHTTPScenario("GET", "/api/foo", nil),
		).Add(
			makeHTTPScenario("GET", "/api/bar", nil),
		).Add(
			makeHTTPScenario("GET", "/api/baz", nil),
		),
	)
	scenarios.Next(
		&gostress.DelayScenario{
			Duration: 1 * time.Millisecond,
			Scenario: makeHTTPScenario("GET", "/api/hoge", nil),
		},
	)
	scenarios.Next(
		&gostress.DeferScenario{
			Defer: func (state gostress.ScenarioState) gostress.Scenario {
				return makeHTTPScenario("GET", "/api/fuga", nil)
			},
		},
	)
	scenarios.Next(
		&gostress.DeferScenario{
			Defer: func (state gostress.ScenarioState) gostress.Scenario {
				return &gostress.NoopScenario{}
			},
		},
	)
	return scenarios
}

func makeHTTPScenario(method, path string, content interface{}) gostress.Scenario {
	return &gostress.HttpScenario{
		Method:  method,
		Path:    path,
		Content: content,
		OnComplete: func(state gostress.ScenarioState, res *gostress.HttpResponse, duration time.Duration) {
			log.Printf("method:%s\tpath:%s\tstatus:%d\ttime:%f", method, path, res.StatusCode, duration.Seconds())
		},
		OnError: func(state gostress.ScenarioState, err error) {
			log.Printf("Error: %s", err)
		},
	}
}
```
