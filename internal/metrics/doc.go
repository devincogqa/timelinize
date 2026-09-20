// Package metrics provides small, dependency-free helpers for aggregating
// numeric samples that are collected while importing data sources.
//
// The helpers here are intentionally lightweight: they operate purely on
// in-memory float64 slices and do not touch the database or the import
// pipeline. They are meant to be used by callers that want quick summary
// statistics (means, percentiles, sliding-window rates) without pulling in
// a heavier statistics dependency.
package metrics
