package assert

import (
	"cmp"
)

// fieldDiff represents a single difference between two compared values.
// It stores the path to the differing field, along with the expected and actual values.
// This is used by the Match function to provide detailed information about differences.
type fieldDiff struct {
	Path     string      // The path to the field, using dot notation for nested fields
	Expected interface{} // The expected value at this path
	Actual   interface{} // The actual value at this path
}

// SimilarItem represents a similar item found
type SimilarItem struct {
	Value      interface{}
	Index      int
	Similarity float64
	DiffType   string // "typo", "case", "prefix", "suffix", "substring"
	Details    string // description of the difference
}

// ContainResult result of the contains search
type ContainResult struct {
	Found   bool
	Exact   bool
	Similar []SimilarItem
	Context []interface{}
	MaxShow int
	Total   int
}

// findInsertionInfo finds information about where a target would be inserted in a sorted collection.
// It returns:
// - found: boolean, true if the target is found in the collection.
// - insertIndex: the index where the target is or would be inserted in the sorted collection.
// - prev: the element just before the insertion point in the sorted collection (if any).
// - next: the element at the insertion point in the sorted collection (if any).
type insertionInfo[T Ordered] struct {
	found        bool
	insertIndex  int
	prev         *T
	next         *T
	sortedWindow string
}

type duplicateGroup struct {
	Value   any
	Indexes []int
}

// PanicInfo contains information about a panic that occurred.
type panicInfo struct {
	Panicked  bool
	Recovered any
	Stack     string
}

// Ordered is a type constraint for types that can be ordered.
// It includes all integer and floating-point types.
type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// Sortable is a type constraint for types that can be sorted.
// It uses Go's cmp.Ordered constraint for type-safe sorting operations.
type Sortable interface {
	cmp.Ordered
}

// MapContainResult represents the result of checking if a map contains a key or value
type MapContainResult struct {
	Found        bool
	Exact        bool
	MaxShow      int
	Total        int
	Context      []interface{}
	Similar      []SimilarItem
	CloseMatches []CloseMatch
}

// CloseMatch holds information about a value that is partially similar to the target.
type CloseMatch struct {
	Value       interface{}
	Differences []string
}

// sortViolation represents a single violation in sort order
type sortViolation struct {
	Index int
	Value interface{}
	Next  interface{}
}

// sortCheckResult contains the result of checking if a collection is sorted
type sortCheckResult struct {
	IsSorted   bool
	Violations []sortViolation
	Total      int
}
