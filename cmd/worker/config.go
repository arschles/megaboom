package main

import "time"

type config struct {
	Namespace     string        `envconfig:"NAMESPACE" required:"true"`
	TotalRequests int           `envconfig:"TOTAL_REQUESTS" required:"true"`
	Concurrency   int           `envconfig:"CONCURRENCY" required:"true"`
	MaxRunTime    time.Duration `envconfig:"MAX_RUN_TIME" required:"true"`
	Endpoint      string        `envconfig:"ENDPOINT" required:"true"`
}

func (c *config) loggerVals() []interface{} {
	return []interface{}{
		"totalRequests",
		c.TotalRequests,
		"concurrency",
		c.Concurrency,
	}
}
