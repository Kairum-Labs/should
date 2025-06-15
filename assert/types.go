package assert

// Assertion is a struct that contains the value to be asserted.
type Assertion[T any] struct {
	value T // The value to be asserted.
}

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

// Ordered is a type constraint for types that can be ordered.
// It includes all integer and floating-point types.
type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}
