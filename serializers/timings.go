package serializers

import "github.com/erraroo/erraroo/models"

type Timing struct {
	*models.Timing
}

type Timings struct {
	Timings []Timing
}

func NewTimings(ps []*models.Timing) Timings {
	timings := Timings{}
	timings.Timings = make([]Timing, len(ps))

	for i, p := range ps {
		timings.Timings[i] = Timing{p}
	}

	return timings
}
