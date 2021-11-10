package main

import "testing"

func Benchmark11Node4Proc(b *testing.B) {
	g := parseGraph("graphs/Nodes_11_OutTree.dot")
	findOptimalSchedule(g, 4)
}
