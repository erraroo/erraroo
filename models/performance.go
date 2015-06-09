package models

type windowPerformance struct {
	Memory struct {
		JsHeapSizeLimit int64
		TotalJSHeapSize int64
		UsedJSHeapSize  int64
	}
	Navigation struct {
		RedirectCount int64
		Type          int
	}
	Timing Timing
}
