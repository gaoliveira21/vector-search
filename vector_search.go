package main

import (
	"math"
	"sort"
)

const MaxAmount = 10000.00
const MaxHour = 23.00
const MaxAvg = 5000.00

type Transaction struct {
	Amount            int
	Hour              int
	CustomerAvgAmount int
}

type Vector []float64

type RefData struct {
	Vector Vector
	Label  string
}

type Dataset []RefData

type SearchResult struct {
	Score float64
	fraud bool
}

var ds = Dataset{
	RefData{Vector: Vector{0.0100, 0.0833, 0.05}, Label: "legit"},
	RefData{Vector: Vector{0.5796, 0.9167, 1.00}, Label: "fraud"},
	RefData{Vector: Vector{0.0035, 0.1667, 0.05}, Label: "legit"},
	RefData{Vector: Vector{0.9708, 1.0000, 1.00}, Label: "fraud"},
	RefData{Vector: Vector{0.4082, 1.0000, 1.00}, Label: "fraud"},
	RefData{Vector: Vector{0.0092, 0.0833, 0.05}, Label: "legit"},
}

func roundUp(val float64, precision int) float64 {
	pow := math.Pow(10, float64(precision))
	return math.Ceil(val*pow) / pow
}

func clamp(val float64, max float64) float64 {
	return min(roundUp(val/max, 2), 1.0)
}

func NewVector(t *Transaction) Vector {
	return Vector{
		clamp(float64(t.Amount)/float64(100), MaxAmount),
		clamp(float64(t.Hour), MaxHour),
		clamp(float64(t.CustomerAvgAmount)/float64(100), MaxAvg),
	}
}

func EuclideanDist(v Vector, ref Vector) float64 {
	coordDist := 0.0

	for i := range v {
		coordDist += math.Pow(v[i]-ref[i], 2)
	}
	score := math.Sqrt(coordDist)

	return roundUp(score, 3)
}

func ManhattanDist(v Vector, ref Vector) float64 {
	coordDist := 0.0

	for i := range v {
		coordDist += math.Abs(v[i] - ref[i])
	}

	return roundUp(coordDist, 3)
}

func CosineSimilarity(v Vector, ref Vector) float64 {
	dot := 0.0
	normV := 0.0
	normRef := 0.0

	for i := range v {
		dot += v[i] * ref[i]
		normV += v[i] * v[i]
		normRef += ref[i] * ref[i]
	}

	return roundUp(dot/(math.Sqrt(normV)*math.Sqrt(normRef)), 3)
}

func DotProduct(v Vector, ref Vector) float64 {
	result := 0.0

	for i := range v {
		result += v[i] * ref[i]
	}

	return roundUp(result, 3)
}

func Search(tr *Transaction, dist func(Vector, Vector) float64) []SearchResult {
	result := []SearchResult{}
	vec := NewVector(tr)

	for _, v := range ds {
		score := dist(vec, v.Vector)
		result = append(result, SearchResult{Score: score, fraud: v.Label == "fraud"})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Score < result[j].Score
	})

	return result
}

func Score(sr []SearchResult) float64 {
	legit := 0
	fraud := 0
	k := 3

	sr = sr[:k]

	for _, v := range sr {
		if v.fraud {
			fraud++
		} else {
			legit++
		}
	}

	return float64(fraud) / float64(k)
}
