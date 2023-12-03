package sample

import (
	ffmt "fmt"
	"time"
)

type Sample struct {
	Name string
}

func NewSample(t time.Time) *Sample {
	return &Sample{
		Name: ffmt.Sprintf("sample_%s", t.Format(time.DateOnly)),
	}
}

func SampleString() string {
	return "sampleString"
}
