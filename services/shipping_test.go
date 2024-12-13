package services

import "testing"

func TestCalculatePack(t *testing.T) {
	testcases := []struct {
		name           string
		order          int
		expectedResult map[int]int
	}{
		{
			name:  "orders of 1",
			order: 1,
			expectedResult: map[int]int{
				250: 1,
			},
		},
		{
			name:  "orders of 250",
			order: 250,
			expectedResult: map[int]int{
				250: 1,
			},
		},
		{
			name:  "orders of 251",
			order: 251,
			expectedResult: map[int]int{
				500: 1,
			},
		},
		{
			name:  "orders of 501",
			order: 501,
			expectedResult: map[int]int{
				500: 1,
				250: 1,
			},
		},
		{
			name:  "orders of 12001",
			order: 12001,
			expectedResult: map[int]int{
				5000: 2,
				2000: 1,
				250:  1,
			},
		},
		{
			name:  "orders of 13001",
			order: 13001,
			expectedResult: map[int]int{
				5000: 2, // 10,000
				2000: 1, // 2,000
				1000: 1, // 1,000
				250:  1, // 250
			},
		},
		{
			name:  "orders of 2390",
			order: 2390,
			expectedResult: map[int]int{
				2000: 1, // 20,000
				500:  1, // 500
			},
		},
	}

	var packSizes = []int{5000, 2000, 1000, 500, 250}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			result := findOptimalPacks(packSizes, tc.order)
			if len(result) != len(tc.expectedResult) {
				t.Errorf("number of entries should match, expected: %v; got %v", len(tc.expectedResult), len(result))
			}
			for k, v := range tc.expectedResult {
				if v != result[k] {
					t.Errorf("expected number of %d pack to be %d but got %v", k, v, result)
				}
			}
		})
	}
}
