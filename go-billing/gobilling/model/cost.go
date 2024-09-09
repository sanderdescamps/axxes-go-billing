package model

import (
	"fmt"
	"time"
)

type Cost struct {
	Fixed      float64 `json:"cost_fixed"`
	CostPerSec float64 `json:"cost_per_sec"`
}

func NewCost() *Cost {
	return &Cost{
		Fixed:      0.0,
		CostPerSec: 0.0,
	}
}

func (c *Cost) Add(o Cost) *Cost {
	c.CostPerSec = c.CostPerSec + o.CostPerSec
	c.Fixed = c.Fixed + o.Fixed
	return c
}

func (c *Cost) Multiply(a float64) *Cost {
	c.CostPerSec = c.CostPerSec * a
	c.Fixed = c.Fixed * a
	return c
}

func (c *Cost) Equal(o Cost) bool {
	return c.CostPerSec == o.CostPerSec && c.Fixed == o.Fixed
}

func (c Cost) CostForDuration(duration time.Duration) float64 {
	return c.Fixed + c.CostPerSec*duration.Seconds()
}

func (c Cost) ToString() string {
	return fmt.Sprintf("fixed: %.2f, cost_per_second: %.2f", c.Fixed, c.CostPerSec)
}
