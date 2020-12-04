package main

import "testing"

func TestMain(t *testing.T) {
	main()
}

func BenchmarkMain(t *testing.B) {
	main()
}
