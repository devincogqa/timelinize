package metrics

// Sum returns the total of all values in the slice.
func Sum(values []float64) float64 {
	var total float64
	for _, v := range values {
		total += v
	}
	return total
}

// Chunk splits values into consecutive slices of at most size elements. The
// final chunk may be shorter. A size of zero or less returns the input as a
// single chunk.
func Chunk(values []float64, size int) [][]float64 {
	if size <= 0 {
		return [][]float64{values}
	}

	chunks := make([][]float64, 0, (len(values)+size-1)/size)
	for start := 0; start < len(values); start += size {
		end := start + size
		if end > len(values) {
			end = len(values)
		}
		chunks = append(chunks, values[start:end])
	}
	return chunks
}
