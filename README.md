# Vector Search

A Go implementation of vector-based similarity search for fraud detection.

## Overview

This project converts transaction data into normalized vectors and uses distance metrics to find similar historical transactions, enabling fraud detection through k-nearest neighbor classification.

## Usage

```go
tr := &Transaction{
    Amount:            1250000,
    Hour:              22,
    CustomerAvgAmount: 480000,
}

results := Search(tr, EuclideanDist)
score := Score(results) // fraud probability (0.0 - 1.0)
```

## Distance Metrics

- **EuclideanDist** - Straight-line distance between vectors
- **ManhattanDist** - Sum of absolute coordinate differences
- **CosineSimilarity** - Angle between vectors (-1 to 1)
- **DotProduct** - Vector projection magnitude

## API

### `NewVector(t *Transaction) Vector`
Normalizes a transaction into a 3D vector [amount, hour, customerAvgAmount].

### `Search(tr *Transaction, dist func(Vector, Vector) float64) []SearchResult`
Finds all reference vectors sorted by distance to the query vector.

### `Score(sr []SearchResult) float64`
Calculates fraud probability based on the k=3 nearest neighbors.

## Running Tests

```bash
go test -v ./...
go test -fuzz=FuzzClampValues -fuzztime=10s ./...
```