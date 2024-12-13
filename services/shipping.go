package services

import (
	"math"
)

type dpEntry struct {
	count int
	prev  int
	pack  int
}

// FindOptimalPacks takes a sorted ascending list of available pack sizes
func FindOptimalPacks(packSizes []int, N int) map[int]int { //nolint:cyclop
	if len(packSizes) < 1 {
		return nil
	}

	maxCheck := N + packSizes[0]

	dp := make([]dpEntry, maxCheck+1)
	for i := range dp {
		dp[i].count = -1
	}
	dp[0].count = 0

	for x := 0; x <= maxCheck; x++ {
		if dp[x].count == -1 {
			continue
		}
		for _, p := range packSizes {
			next := x + p
			if next <= maxCheck {
				if dp[next].count == -1 || dp[next].count > dp[x].count+1 {
					dp[next].count = dp[x].count + 1
					dp[next].prev = x
					dp[next].pack = p
				}
			}
		}
	}

	minLeftover := math.MaxInt
	minPacks := math.MaxInt
	bestSum := -1

	for x := N; x <= maxCheck; x++ {
		if dp[x].count != -1 {
			leftover := x - N
			if leftover < minLeftover {
				minLeftover = leftover
				minPacks = dp[x].count
				bestSum = x
			} else if leftover == minLeftover && dp[x].count < minPacks {
				// same leftover but fewer packs
				minPacks = dp[x].count
				bestSum = x
			}
		}
	}

	if bestSum == -1 {
		return nil
	}

	// Reconstruct solution
	packCount := make(map[int]int)
	current := bestSum
	for current > 0 {
		entry := dp[current]
		packCount[entry.pack]++
		current = entry.prev
	}

	return packCount
}
