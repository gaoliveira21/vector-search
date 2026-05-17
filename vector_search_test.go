package main

import (
	"slices"
	"testing"
)

func TestRoundUp(t *testing.T) {
	tests := []struct {
		name      string
		val       float64
		precision int
		want      float64
	}{
		{"round up 0.063660", 0.063660, 3, 0.064},
		{"round down 0.063664", 0.063664, 3, 0.064},
		{"exact 0.5", 0.5, 2, 0.5},
		{"large value", 123.456, 1, 123.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := roundUp(tt.val, tt.precision)
			if got != tt.want {
				t.Errorf("roundUp(%f, %d) = %f, want %f", tt.val, tt.precision, got, tt.want)
			}
		})
	}
}

func TestClamp(t *testing.T) {
	tests := []struct {
		name string
		val  float64
		max  float64
		want float64
	}{
		{"under max", 125.0, 10000.0, 0.02},
		{"at max", 10000.0, 10000.0, 1.0},
		{"over max", 50000.0, 10000.0, 1.0},
		{"zero val", 0.0, 100.0, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := clamp(tt.val, tt.max)
			if got != tt.want {
				t.Errorf("clamp(%f, %f) = %f, want %f", tt.val, tt.max, got, tt.want)
			}
		})
	}
}

func FuzzClampValues(f *testing.F) {
	f.Add(12500.00, 10000.00, 1.0)
	f.Add(22.00, 23.00, 0.96)
	f.Add(4800.00, 5000.00, 0.96)

	f.Fuzz(func(t *testing.T, val float64, max float64, expected float64) {
		result := clamp(val, max)
		if result < 0 || result > 1 {
			t.Errorf("clamp result out of range [0,1]: %f", result)
		}
		_ = expected
	})
}

func TestNewVector(t *testing.T) {
	tests := []struct {
		name     string
		tx       Transaction
		expected Vector
	}{
		{
			name:     "max values",
			tx:       Transaction{Amount: 1250000, Hour: 22, CustomerAvgAmount: 480000},
			expected: Vector{1.0, 0.96, 0.96},
		},
		{
			name:     "zero values",
			tx:       Transaction{Amount: 0, Hour: 0, CustomerAvgAmount: 0},
			expected: Vector{0.0, 0.0, 0.0},
		},
		{
			name:     "mid values",
			tx:       Transaction{Amount: 500000, Hour: 12, CustomerAvgAmount: 250000},
			expected: Vector{0.5, 0.53, 0.5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewVector(&tt.tx)
			if !slices.Equal(got, tt.expected) {
				t.Errorf("NewVector() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestEuclideanDist(t *testing.T) {
	tests := []struct {
		name     string
		v        Vector
		ref      Vector
		expected float64
	}{
		{"same vector", Vector{1.0, 1.0, 1.0}, Vector{1.0, 1.0, 1.0}, 0.0},
		{"diff vector", Vector{1.0, 0.96, 0.96}, Vector{0.9708, 1.0000, 1.00}, 0.064},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EuclideanDist(tt.v, tt.ref)
			if got != tt.expected {
				t.Errorf("EuclideanDist() = %f, want %f", got, tt.expected)
			}
		})
	}
}

func TestManhattanDist(t *testing.T) {
	tests := []struct {
		name     string
		v        Vector
		ref      Vector
		expected float64
	}{
		{"same vector", Vector{1.0, 1.0, 1.0}, Vector{1.0, 1.0, 1.0}, 0.0},
		{"diff vector", Vector{1.0, 0.96, 0.96}, Vector{0.9708, 1.0000, 1.00}, 0.11},
		{"simple case", Vector{0.0, 0.0}, Vector{3.0, 4.0}, 7.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ManhattanDist(tt.v, tt.ref)
			if got != tt.expected {
				t.Errorf("ManhattanDist() = %f, want %f", got, tt.expected)
			}
		})
	}
}

func TestCosineSimilarity(t *testing.T) {
	tests := []struct {
		name     string
		v        Vector
		ref      Vector
		expected float64
	}{
		{"identical vectors", Vector{1.0, 0.0}, Vector{1.0, 0.0}, 1.0},
		{"opposite vectors", Vector{1.0, 0.0}, Vector{-1.0, 0.0}, -1.0},
		{"orthogonal vectors", Vector{1.0, 0.0}, Vector{0.0, 1.0}, 0.0},
		{"similar vectors", Vector{1.0, 1.0}, Vector{0.5, 0.5}, 1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CosineSimilarity(tt.v, tt.ref)
			if got != tt.expected {
				t.Errorf("CosineSimilarity() = %f, want %f", got, tt.expected)
			}
		})
	}
}

func TestDotProduct(t *testing.T) {
	tests := []struct {
		name     string
		v        Vector
		ref      Vector
		expected float64
	}{
		{"simple", Vector{1.0, 2.0, 3.0}, Vector{4.0, 5.0, 6.0}, 32.0},
		{"zero vector", Vector{1.0, 2.0}, Vector{0.0, 0.0}, 0.0},
		{"same vector", Vector{2.0, 3.0}, Vector{2.0, 3.0}, 13.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DotProduct(tt.v, tt.ref)
			if got != tt.expected {
				t.Errorf("DotProduct() = %f, want %f", got, tt.expected)
			}
		})
	}
}

func TestSearch(t *testing.T) {
	tr := &Transaction{
		Amount:            1250000,
		Hour:              22,
		CustomerAvgAmount: 480000,
	}

	result := Search(tr, EuclideanDist)

	if len(result) != len(ds) {
		t.Errorf("expected %d results, got %d", len(ds), len(result))
	}

	if result[0].Score != 0.064 {
		t.Errorf("expected first score 0.064, got %f", result[0].Score)
	}

	for i := 1; i < len(result); i++ {
		if result[i].Score < result[i-1].Score {
			t.Error("results are not sorted by score")
		}
	}

	expectedScores := []float64{0.064, 0.425, 0.595, 1.566, 1.606, 1.606}
	for i, expected := range expectedScores {
		if result[i].Score != expected {
			t.Errorf("result[%d] = %f, want %f", i, result[i].Score, expected)
		}
	}
}

func TestSearchResultsSorted(t *testing.T) {
	transactions := []Transaction{
		{Amount: 1250000, Hour: 22, CustomerAvgAmount: 480000},
		{Amount: 100, Hour: 1, CustomerAvgAmount: 100},
		{Amount: 500000, Hour: 12, CustomerAvgAmount: 250000},
	}

	for _, tx := range transactions {
		t.Run("", func(t *testing.T) {
			result := Search(&tx, EuclideanDist)
			for i := 1; i < len(result); i++ {
				if result[i].Score < result[i-1].Score {
					t.Errorf("results not sorted for tx %+v: score[%d]=%f < score[%d]=%f",
						tx, i, result[i].Score, i-1, result[i-1].Score)
				}
			}
		})
	}
}

func TestScore(t *testing.T) {
	tests := []struct {
		name      string
		results   []SearchResult
		k         int
		wantScore float64
	}{
		{
			name:      "all fraud",
			results:   []SearchResult{{Score: 0.1, fraud: true}, {Score: 0.2, fraud: true}, {Score: 0.3, fraud: true}},
			k:         3,
			wantScore: 1.0,
		},
		{
			name:      "all legit",
			results:   []SearchResult{{Score: 0.1, fraud: false}, {Score: 0.2, fraud: false}, {Score: 0.3, fraud: false}},
			k:         3,
			wantScore: 0.0,
		},
		{
			name:      "mixed - 2 fraud 1 legit",
			results:   []SearchResult{{Score: 0.1, fraud: true}, {Score: 0.2, fraud: true}, {Score: 0.3, fraud: false}},
			k:         3,
			wantScore: 2.0 / 3.0,
		},
		{
			name:      "mixed - 1 fraud 2 legit",
			results:   []SearchResult{{Score: 0.1, fraud: true}, {Score: 0.2, fraud: false}, {Score: 0.3, fraud: false}},
			k:         3,
			wantScore: 1.0 / 3.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Score(tt.results)
			if got != tt.wantScore {
				t.Errorf("Score() = %f, want %f", got, tt.wantScore)
			}
		})
	}
}

func TestScoreWithTransaction(t *testing.T) {
	tr := &Transaction{
		Amount:            1250000,
		Hour:              22,
		CustomerAvgAmount: 480000,
	}

	result := Search(tr, EuclideanDist)
	score := Score(result)

	if score != 1.0 {
		t.Errorf("expected score 1.0 (all top-3 are fraud), got %f", score)
	}
}
