package main

import "testing"

func BenchmarkFindOptimalSchedule(b *testing.B) {
	g := parseGraph("graphs/16Nodes4Processors.dot")

	for i := 0; i < b.N; i++ {
		FindOptimalSchedule(g, 4)
	}
}
