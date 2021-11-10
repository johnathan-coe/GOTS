package main

import "testing"

func BenchmarkFindOptimalSchedule(b *testing.B) {
	b.StopTimer()
	g := parseGraph("graphs/16Nodes4Processors.dot")
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		FindOptimalSchedule(g, 4)
	}
}
