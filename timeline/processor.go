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
	"context"
	"encoding/json"
	"runtime/debug"
	"sync/atomic"

	"github.com/ringsaturn/tzf"
	"go.uber.org/zap"
)

// processor orchestrates the import pipeline for a single data source.
// It manages batching, checkpointing, and coordinating the flow of items
// from a data source into the timeline database.
type processor struct {
	// estimatedCount tracks the estimated number of items to process;
	// accessed atomically and 64-bit aligned. Also used for estimating
	// total size before running.
	estimatedCount *int64

	// outerLoopIdx and innerLoopIdx track the current position for
	// checkpointing, allowing imports to be resumed from where they
	// left off.
	outerLoopIdx, innerLoopIdx int

	// ij is the parent import job that this processor belongs to.
	ij *ImportJob

	// ds is the data source being imported from.
	ds DataSource

	// dsRowID is the database row ID of the data source.
	dsRowID uint64

	// dsOpt holds data-source-specific configuration options.
	dsOpt any

	// tl is the timeline that items are being imported into.
	tl *Timeline

	// log is the structured logger for this processor.
	log *zap.Logger

	// batch accumulates graphs for bulk insertion, which greatly
	// increases import speed.
	batch []*Graph

	// batchSize is at least len(batch), but edges on a graph can
	// add to it, so it may exceed the batch length.
	batchSize int

	// rootGraphCount is a counter used for periodic database optimization
	// during an import.
	// TODO: this should ideally be global per timeline, even if multiple jobs run simultaneously
	rootGraphCount int

	// tzFinder is used for guessing time zones based on geographic coordinates.
	tzFinder tzf.F
}

// process runs the import pipeline for a single directory entry. If estimatedCount
// is non-nil, it performs a size estimation pass instead of a full import. Otherwise,
// it runs the complete import: creating a file importer, sending items through the
// processing pipeline, and waiting for all workers to finish.
func (p processor) process(ctx context.Context, dirEntry DirEntry, dsCheckpoint json.RawMessage) error {
	defer p.recoverFromPanic()

	ctx = p.contextWithOwner(ctx)

	importParams := p.buildImportParams(dirEntry, dsCheckpoint)

	if p.estimatedCount != nil {
		return p.estimateSize(ctx, dirEntry, importParams)
	}

	return p.executeImport(ctx, dirEntry, importParams)
}

// recoverFromPanic recovers from any panic during import processing,
// preventing a misbehaving data source from crashing the entire application.
func (p processor) recoverFromPanic() {
	if r := recover(); r != nil {
		p.log.Error("panic",
			zap.Any("error", r),
			zap.String("stack", string(debug.Stack())))
	}
}

// contextWithOwner adds the timeline owner entity to the context, making it
// available to data sources that need an item owner when one is otherwise unknown.
func (p processor) contextWithOwner(ctx context.Context) context.Context {
	owner, err := p.tl.LoadEntity(ownerEntityID)
	if err == nil {
		ctx = context.WithValue(ctx, RepoOwnerCtxKey, owner)
	}
	return ctx
}

// buildImportParams constructs the ImportParams used by file importers,
// populating it with the logger, continuation function, timeframe filter,
// checkpoint data, and data-source-specific options.
func (p processor) buildImportParams(dirEntry DirEntry, dsCheckpoint json.RawMessage) ImportParams {
	return ImportParams{
		Log:               p.log.With(zap.String("filename", dirEntry.FullPath())),
		Continue:          p.ij.job.Continue,
		Timeframe:         p.ij.ProcessingOptions.Timeframe,
		Checkpoint:        dsCheckpoint,
		DataSourceOptions: p.dsOpt,
	}
}

// estimateSize runs a pre-import pass to estimate the number of items that will
// be processed. If the data source implements SizeEstimator, it uses the optimized
// estimation path. Otherwise, it falls back to running a full import in count-only
// mode and tallying the incoming graphs.
func (p processor) estimateSize(ctx context.Context, dirEntry DirEntry, params ImportParams) error {
	fileImporter := p.ds.NewFileImporter()

	if estimator, ok := fileImporter.(SizeEstimator); ok {
		return p.estimateSizeOptimized(ctx, dirEntry, params, estimator)
	}

	return p.estimateSizeByCountingGraphs(ctx, dirEntry, params, fileImporter)
}

// estimateSizeOptimized uses a data source's native SizeEstimator implementation
// for faster (but probably still slow) size estimation.
func (p processor) estimateSizeOptimized(ctx context.Context, dirEntry DirEntry, params ImportParams, estimator SizeEstimator) error {
	totalSize, err := estimator.EstimateSize(ctx, dirEntry, params)
	if err != nil {
		params.Log.Error("could not estimate import size", zap.Error(err))
	}
	// don't call SetTotal() yet -- wait until we're done counting,
	// so progress bars don't think we are done with the estimate
	atomic.AddInt64(p.estimatedCount, int64(totalSize))
	return nil
}

// estimateSizeByCountingGraphs falls back to running the import in count-only mode
// when the data source lacks an optimized SizeEstimator. It counts the graph sizes
// that come through the pipeline to produce an estimate.
func (p processor) estimateSizeByCountingGraphs(ctx context.Context, dirEntry DirEntry, params ImportParams, fileImporter FileImporter) error {
	done := make(chan struct{})
	wg, graphsCh := p.beginProcessing(ctx, p.ij.ProcessingOptions, true, done)
	params.Pipeline = graphsCh

	err := fileImporter.FileImport(ctx, dirEntry, params)
	if err != nil {
		params.Log.Error("failed estimating size before import", zap.Error(err))
	}

	// sending on the pipeline must be complete by now; signal to workers to exit
	close(done)

	// wait for all processing workers to complete so we have an accurate count
	wg.Wait()

	return nil
}

// executeImport performs the full import pipeline: it starts processing workers,
// runs the file importer to send items through the pipeline, then waits for all
// workers to finish before returning.
func (p processor) executeImport(ctx context.Context, dirEntry DirEntry, params ImportParams) error {
	done := make(chan struct{})
	wg, graphsCh := p.beginProcessing(ctx, p.ij.ProcessingOptions, false, done)
	params.Pipeline = graphsCh

	// use a fresh file importer to avoid any potentially reused state (TODO: necessary?)
	importErr := p.ds.NewFileImporter().FileImport(ctx, dirEntry, params)

	// sending on the pipeline must be complete by now; signal to workers to exit
	close(done)

	params.Log.Info("importer done sending items; waiting for processing to finish", zap.Error(importErr))

	// wait for all processing workers to complete
	wg.Wait()

	return importErr
}

// ProcessingOptions configures how item processing is carried out.
type ProcessingOptions struct {
	// Whether to perform integrity checks
	Integrity bool `json:"integrity,omitempty"`

	// Constrain processed items to within a timeframe
	Timeframe Timeframe `json:"timeframe,omitempty"`

	// If true, items with manual modifications may be updated, overwriting local changes.
	OverwriteLocalChanges bool `json:"overwrite_local_changes,omitempty"`

	// Names of columns in the items table to check for sameness when loading an item
	// that doesn't have data_source+original_id. The field/column is the same if the
	// values are identical or if one of the values is NULL. If the map value is true,
	// however, strict NULL comparison is applied, where NULL=NULL only. In other words,
	// with a false value, it is as if the field is not in the unique constraints at
	// all when applied to a field that is null.
	ItemUniqueConstraints map[string]bool `json:"item_unique_constraints,omitempty"`

	// How to update existing items, specified per-field.
	// If not set, items that already exist will simply be skipped.
	ItemUpdatePreferences []FieldUpdatePreference `json:"item_update_preferences,omitempty"`

	// TODO: WIP (Should this be in importjob or processingoptions?)
	Interactive *InteractiveImport `json:"interactive,omitempty"`

	// How many items to process in one database transaction. This is a minimum value,
	// not a maximum, due to the recursive nature of graphs. (Hopefully data sources
	// do not build out graphs that are too big for available memory. Most graphs
	// just have 1-2 things on it). This setting is sensitive. Smaller batches can
	// reduce performance, slowing down imports. Larger batches can be faster, allowing
	// for more throughput especially when data files are large, but reduces the
	// granularity of the live progress updates, and the performance benefit gets
	// smaller as the database gets larger and the items get smaller.
	BatchSize int `json:"batch_size,omitempty"`

	// When enabled, items that are received for processing which have coordinates and
	// a timestamp but lack a time zone may have the time zone augmented by doing a
	// geopoly lookup. The data set uses more memory but can help make timestamps more
	// correct/complete. Since Go does not actually support time.Time values without a
	// time zone, by "lack a time zone" we mean the location is time.Local (i.e. "wall
	// time"), which essentially means, "time zone unknown".
	InferTimeZone bool `json:"infer_time_zone,omitempty"`

	// Generate thumbnails of images and videos as they are processed, rather than waiting
	// until the thumbnailing job after the import finishes. This uses more memory and
	// slows down the import job during media items.
	Thumbnails bool `json:"thumbnails,omitempty"`
}

type ctxKey string

// RepoOwnerCtxKey is the context key used to store the repository owner entity.
var RepoOwnerCtxKey ctxKey = "repo_owner"

// ownerEntityID is the fixed database ID for the timeline owner entity.
const ownerEntityID = 1
