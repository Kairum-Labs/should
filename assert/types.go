package assert

// Option is a functional option for configuring assertions.
type Option interface {
	Apply(config *Config)
}

// Config provides configuration options for assertions.
// It allows for custom error messages and future extensibility.
type Config struct {
	Message string
	/*
		 	Description    string
			DeepComparison bool
	*/
}

// message implements the Option interface for custom messages.
type message string

// Apply sets the custom message in the config.
func (m message) Apply(c *Config) {
	c.Message = string(m)
}

// WithMessage creates an option for setting a custom error message.
func WithMessage(msg string) Option {
	return message(msg)
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

// Ordered is a type constraint for types that can be ordered.
// It includes all integer and floating-point types.
type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}
