package main

import (
	"hash/fnv"
)

func fnvhash(b []byte) uint64 {
	hash := fnv.New64a()
	hash.Write(b)
	return hash.Sum64()
}

func volumeMultiplier(level int) float64 {
	switch level {
	case 1:
		return 0.25
	case 2:
		return 0.55
	case 3:
		return 0.8
	case 4:
		return 1
	default:
		return 0
	}
}
