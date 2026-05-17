# Vector Search - Fraud Detection Implementation Guide

## What is Vector Search?

Vector search finds similar items by comparing their vector embeddings (numerical representations) using distance metrics. This implementation applies it to fraud detection by finding similar past transactions.

---

## Implementation Overview

### Tech Stack

- **Language**: Go 1.26.2
- **Approach**: Brute-force k-nearest neighbors search
- **Use Case**: Fraud detection via transaction similarity

### Core Components

1. **Transaction** → **Vector**: Normalize transaction data into 3D vectors
2. **Distance Functions**: Euclidean, Manhattan, Cosine Similarity, Dot Product
3. **Search**: Find k nearest neighbors in reference dataset
4. **Scoring**: Calculate fraud probability from neighbor labels

---

## Data Flow

```
Transaction (Amount, Hour, CustomerAvgAmount)
    ↓ NewVector()
Vector [0.01, 0.08, 0.05]  (normalized, clamped to [0,1])
    ↓ Search() with distance function
[]SearchResult {Score, fraud}
    ↓ Score()
fraud_probability (0.0 to 1.0)
```

---

## Distance Metrics

| Metric | Formula | Best For |
|--------|---------|----------|
| Euclidean | √Σ(Ai-Bi)² | General similarity |
| Manhattan | Σ\|Ai-Bi\| | Grid-like paths |
| Cosine Similarity | (A·B)/(\|\|A\|\|×\|\|B\|\|) | Direction similarity |
| Dot Product | Σ(Ai×Bi) | Projection magnitude |

---

## Code Examples (Go)

### Vector Creation

```go
type Transaction struct {
    Amount            int
    Hour              int
    CustomerAvgAmount int
}

type Vector []float64

func NewVector(t *Transaction) Vector {
    return Vector{
        clamp(float64(t.Amount)/100, 10000.0),
        clamp(float64(t.Hour), 23.0),
        clamp(float64(t.CustomerAvgAmount)/100, 5000.0),
    }
}
```

### Distance Functions

```go
func EuclideanDist(v Vector, ref Vector) float64 {
    coordDist := 0.0
    for i := range v {
        coordDist += math.Pow(v[i]-ref[i], 2)
    }
    return roundUp(math.Sqrt(coordDist), 3)
}

func CosineSimilarity(v Vector, ref Vector) float64 {
    dot, normV, normRef := 0.0, 0.0, 0.0
    for i := range v {
        dot += v[i] * ref[i]
        normV += v[i] * v[i]
        normRef += ref[i] * ref[i]
    }
    return roundUp(dot/(math.Sqrt(normV)*math.Sqrt(normRef)), 3)
}
```

### Search & Scoring

```go
func Search(tr *Transaction, dist func(Vector, Vector) float64) []SearchResult {
    vec := NewVector(tr)
    result := []SearchResult{}
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
    fraud := 0
    for _, v := range sr[:3] {
        if v.fraud {
            fraud++
        }
    }
    return float64(fraud) / 3.0
}
```

---

## Running the Code

```bash
go run vector_search.go
go test -v
```

---

## Extending the Implementation

### Phase 2: KD-Tree

- [ ] Build tree by splitting on median (alternating dimensions)
- [ ] Implement insert and search
- [ ] Handle approximate nearest neighbors
- [ ] Compare performance vs brute force

### Phase 3: HNSW (Production-Ready)

- [ ] Understand skip lists concept
- [ ] Implement layered graph structure
- [ ] Build index (insertion with probabilistic layering)
- [ ] Implement search (greedy graph traversal)

---

## When to Use Each Approach

- **Brute Force (current)**: Dataset < 10,000 vectors, need exact results
- **KD-Tree**: Dataset < 1M, dimensions < 50, exact results
- **HNSW**: Millions of vectors, dimensions > 50, approximate results (production)

---

## Next Steps

1. **Add more reference data** to improve detection accuracy
2. **Implement KD-Tree** for faster searches with larger datasets
3. **Add more distance metrics** (Mahalanobis, Jaccard)
4. **Explore dimensionality reduction** for higher-dimensional data