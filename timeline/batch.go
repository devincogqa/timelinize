/*
	Timelinize
	Copyright (c) 2013 Matthew Holt

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU Affero General Public License as published
	by the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU Affero General Public License for more details.

	You should have received a copy of the GNU Affero General Public License
	along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package timeline

import (
	"sync"
)

// BatchProcessor collects items and processes them in batches
// for more efficient database operations during imports.
type BatchProcessor[T any] struct {
	mu        sync.Mutex
	items     []T
	batchSize int
	processFn func([]T) error
}

// NewBatchProcessor creates a new BatchProcessor with the given batch size
// and processing function.
func NewBatchProcessor[T any](batchSize int, processFn func([]T) error) *BatchProcessor[T] {
	return &BatchProcessor[T]{
		items:     make([]T, 0, batchSize),
		batchSize: batchSize,
		processFn: processFn,
	}
}

// Add adds an item to the batch. If the batch reaches the configured
// size, it automatically triggers processing.
func (bp *BatchProcessor[T]) Add(item T) error {
	bp.mu.Lock()
	defer bp.mu.Unlock()

	bp.items = append(bp.items, item)

	if len(bp.items) >= bp.batchSize {
		return bp.processCurrentBatch()
	}

	return nil
}

// Flush processes any remaining items in the batch.
func (bp *BatchProcessor[T]) Flush() error {
	bp.mu.Lock()
	defer bp.mu.Unlock()

	if len(bp.items) == 0 {
		return nil
	}

	return bp.processCurrentBatch()
}

// Count returns the number of items currently waiting in the batch.
func (bp *BatchProcessor[T]) Count() int {
	bp.mu.Lock()
	defer bp.mu.Unlock()
	return len(bp.items)
}

// processCurrentBatch sends current items to the processing function
// and resets the internal buffer.
func (bp *BatchProcessor[T]) processCurrentBatch() error {
	err := bp.processFn(bp.items)
	bp.items = bp.items[:0]
	return err
}
