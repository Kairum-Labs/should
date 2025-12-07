package assert

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"slices"
	"sort"
	"strings"
	"time"
	"unicode/utf8"
)

// maxSimilarLen limits substring length for similarity checks.
// Longer substrings are skipped to improve performance.
// 20 was chosen as a practical balance between accuracy and speed.
const maxSimilarLen = 20

// similarityThreshold defines the minimum similarity difference to consider
// when removing substring matches. Values within this threshold are treated
// as effectively equal, preferring the more complete string.
const similarityThreshold = 0.05

// isSliceOrArray checks if the provided value is a slice or an array.
// It handles nil values by returning false.
func isSliceOrArray(v interface{}) bool {
	if v == nil {
		return false
	}
	kind := reflect.TypeOf(v).Kind()
	return kind == reflect.Slice || kind == reflect.Array
}

// isPrimitive checks if the provided reflect.Kind represents a primitive type.
func isPrimitive(kind reflect.Kind) bool {
	switch kind {
	case reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64, reflect.Bool:
		return true
	default:
		return false
	}
}

// formatSlice formats a slice or an array into a human-readable string.
// e.g., [1, 2, 3] or ["apple", "banana", "orange"].
func formatSlice(slice interface{}) string {
	if !isSliceOrArray(slice) {
		return fmt.Sprintf("<not a slice or array: %T>", slice)
	}
	return formatValueComparison(reflect.ValueOf(slice))
}

// formatComparisonValue formats an arbitrary value for human-readable comparison.
// It applies proper formatting based on the value's type.
func formatComparisonValue(obj interface{}) string {
	return formatValueComparison(reflect.ValueOf(obj))
}

// formatValueComparison handles the formatting logic for different reflect.Value types
// to provide consistent and readable output for comparison purposes.
func formatValueComparison(v reflect.Value) string {
	if !v.IsValid() {
		return "nil"
	}

	switch v.Kind() {
	case reflect.Struct:
		var parts []string
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i)
			if !field.IsExported() {
				continue
			}
			fieldValue := v.Field(i)
			parts = append(parts, fmt.Sprintf("%s: %s", field.Name, formatValueComparison(fieldValue)))
		}
		return fmt.Sprintf("{%s}", strings.Join(parts, ", "))

	case reflect.String:
		return fmt.Sprintf(`"%s"`, v.String())

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Uintptr, reflect.Float32, reflect.Float64,
		reflect.Bool:
		return fmt.Sprint(v.Interface())

	case reflect.Ptr:
		if v.IsNil() {
			return "nil"
		}
		return formatValueComparison(v.Elem())

	case reflect.Slice, reflect.Array:
		if v.IsNil() {
			return "nil"
		}
		var elements []string
		for i := 0; i < v.Len(); i++ {
			elements = append(elements, formatValueComparison(v.Index(i)))
		}
		return fmt.Sprintf("[%s]", strings.Join(elements, ", "))

	case reflect.Map:
		if v.IsNil() {
			return "nil"
		}
		var pairs []string
		for _, key := range v.MapKeys() {
			value := v.MapIndex(key)
			pairs = append(pairs, fmt.Sprintf("%s: %s", formatValueComparison(key), formatValueComparison(value)))
		}
		return fmt.Sprintf("map[%s]", strings.Join(pairs, ", "))

	default:
		if v.CanInterface() {
			return fmt.Sprint(v.Interface())
		}
		return fmt.Sprint(v)
	}
}

// formatDiffValue formats a value specifically for showing differences.
// It handles basic types differently than complex types for better readability.
func formatDiffValue(value interface{}) string {
	if value == nil {
		return "nil"
	}

	switch v := value.(type) {
	case string:
		return fmt.Sprintf("%q", v)
	case bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return fmt.Sprint(v)
	default:
		// for complex types, use our comparison formatter
		return formatComparisonValue(v)
	}
}

// findDifferences locates all differences between two values and returns them as a slice of fieldDiff.
// The function works recursively for nested structures.
func findDifferences(expected, actual interface{}) []fieldDiff {
	return compareExpectedActual(expected, actual, "")
}

// compareExpectedActual compares two values recursively and records any differences in the provided diffs slice.
// It handles complex structures like structs, maps, slices, and arrays.
func compareExpectedActual(expected, actual interface{}, path string) (diffs []fieldDiff) {
	expectedValue := reflect.ValueOf(expected)
	actualValue := reflect.ValueOf(actual)

	if !expectedValue.IsValid() || !actualValue.IsValid() {
		if expectedValue.IsValid() != actualValue.IsValid() {
			diffs = append(diffs, fieldDiff{
				Path:     path,
				Expected: expected,
				Actual:   actual,
			})
		}
		return
	}

	if expectedValue.Kind() != actualValue.Kind() {
		diffs = append(diffs, fieldDiff{
			Path:     path,
			Expected: expectedValue.Kind(),
			Actual:   actualValue.Kind(),
		})
		return
	}

	switch expectedValue.Kind() {
	case reflect.Struct:
		if expectedValue.Type() != actualValue.Type() {
			diffs = append(diffs, fieldDiff{
				Path:     path,
				Expected: expectedValue.Type(),
				Actual:   actualValue.Type(),
			})
			return
		}
		typeOfT := expectedValue.Type()
		for i := 0; i < expectedValue.NumField(); i++ {
			field := typeOfT.Field(i)
			if !field.IsExported() {
				continue
			}
			newPath := buildPath(path, field.Name)

			expectedField := expectedValue.Field(i).Interface()
			actualField := actualValue.Field(i).Interface()

			diffs = append(diffs, compareExpectedActual(expectedField, actualField, newPath)...)
		}

	case reflect.String:
		if expectedValue.String() != actualValue.String() {
			diffs = append(diffs, fieldDiff{
				Path:     path,
				Expected: expectedValue.String(),
				Actual:   actualValue.String(),
			})
		}

	case reflect.Bool:
		if expectedValue.Bool() != actualValue.Bool() {
			diffs = append(diffs, fieldDiff{
				Path:     path,
				Expected: expectedValue.Bool(),
				Actual:   actualValue.Bool(),
			})
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if expectedValue.Int() != actualValue.Int() {
			diffs = append(diffs, fieldDiff{
				Path:     path,
				Expected: expectedValue.Interface(),
				Actual:   actualValue.Interface(),
			})
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if expectedValue.Uint() != actualValue.Uint() {
			diffs = append(diffs, fieldDiff{
				Path:     path,
				Expected: expectedValue.Interface(),
				Actual:   actualValue.Interface(),
			})
		}

	case reflect.Float32, reflect.Float64:
		if expectedValue.Float() != actualValue.Float() {
			diffs = append(diffs, fieldDiff{
				Path:     path,
				Expected: expectedValue.Interface(),
				Actual:   actualValue.Interface(),
			})
		}

	case reflect.Ptr:
		if expectedValue.IsNil() != actualValue.IsNil() {
			diffs = append(diffs, fieldDiff{
				Path:     path,
				Expected: expectedValue.Interface(),
				Actual:   actualValue.Interface(),
			})
			return
		}
		if !expectedValue.IsNil() {
			diffs = append(diffs, compareExpectedActual(expectedValue.Elem().Interface(), actualValue.Elem().Interface(), path)...)
		}

	case reflect.Slice, reflect.Array:
		if expectedValue.IsNil() != actualValue.IsNil() {
			diffs = append(diffs, fieldDiff{
				Path:     path,
				Expected: expectedValue.Interface(),
				Actual:   actualValue.Interface(),
			})
			return
		}

		if expectedValue.Len() == 0 && actualValue.Len() == 0 {
			return
		}

		if expectedValue.Len() != actualValue.Len() {
			diffs = append(diffs, fieldDiff{
				Path:     path,
				Expected: nil,
				Actual:   nil,
				Message: fmt.Sprintf("length mismatch (expected: %d, actual: %d)",
					expectedValue.Len(), actualValue.Len()),
			})
			return
		}

		//compare elements one by one
		for i := 0; i < expectedValue.Len(); i++ {
			if !reflect.DeepEqual(expectedValue.Index(i).Interface(), actualValue.Index(i).Interface()) {
				elementPath := buildPath(path, fmt.Sprintf("[%d]", i))
				diffs = append(
					diffs,
					compareExpectedActual(
						expectedValue.Index(i).Interface(),
						actualValue.Index(i).Interface(),
						elementPath,
					)...,
				)
			}
		}

	case reflect.Map:
		if expectedValue.IsNil() != actualValue.IsNil() {
			diffs = append(diffs, fieldDiff{
				Path:     path,
				Expected: expectedValue.Interface(),
				Actual:   actualValue.Interface(),
			})
			return
		}

		if expectedValue.IsNil() {
			return
		}

		for _, key := range expectedValue.MapKeys() {
			actualVal := actualValue.MapIndex(key)
			keyStr := fmt.Sprint(key.Interface())
			keyPath := buildPath(path, fmt.Sprintf("[%s]", keyStr))

			if !actualVal.IsValid() {
				diffs = append(diffs, fieldDiff{
					Path:     keyPath,
					Expected: expectedValue.MapIndex(key).Interface(),
					Actual:   "<missing>",
				})
				continue
			}

			if !reflect.DeepEqual(expectedValue.MapIndex(key).Interface(), actualVal.Interface()) {
				diffs = append(diffs, compareExpectedActual(
					expectedValue.MapIndex(key).Interface(),
					actualVal.Interface(),
					keyPath,
				)...)
			}
		}

		for _, key := range actualValue.MapKeys() {
			expectedVal := expectedValue.MapIndex(key)
			if !expectedVal.IsValid() {
				keyStr := fmt.Sprint(key.Interface())
				keyPath := buildPath(path, fmt.Sprintf("[%s]", keyStr))

				diffs = append(diffs, fieldDiff{
					Path:     keyPath,
					Expected: "<missing>",
					Actual:   actualValue.MapIndex(key).Interface(),
				})
			}
		}
	}
	return
}

// buildPath creates a dotted path for nested fields to provide clear identification
// of where differences occur in complex structures.
func buildPath(parent, field string) string {
	if parent == "" {
		return field
	}
	return parent + "." + field
}

// formatMultilineString formats long strings into a readable multi-line layout
// for use in error messages. Strings shorter than 280 characters are returned as-is.
// Otherwise, it shows up to 5 initial lines (56 chars each), and if longer,
// appends the last 3 lines for context.
func formatMultilineString(s string) string {
	if len(s) < 280 {
		return s
	}

	builder := strings.Builder{}

	totalLine := len(s) / 56

	builder.WriteString("Length: ")
	builder.WriteString(fmt.Sprintf("%d", len(s)))
	builder.WriteString(" characters, ")
	builder.WriteString(fmt.Sprintf("%d lines", totalLine))
	builder.WriteString("\n")

	// max 5 lines and 56 characters per line
	for i := range 5 {
		builder.WriteString(fmt.Sprintf("%d. ", i+1))
		builder.WriteString(s[i*56 : min(i*56+56, len(s))])
		builder.WriteString("\n")
	}

	if totalLine > 5 {
		builder.WriteString("\n")
		builder.WriteString("Last lines:\n")

		// print last 3 lines
		for i := totalLine - 3; i < totalLine; i++ {
			builder.WriteString(fmt.Sprintf("%d. ", i+1))
			builder.WriteString(s[i*56 : min(i*56+56, len(s))])
			builder.WriteString("\n")
		}
	}

	return builder.String()
}

//  === THIS SECTION IS TO FIND SIMILAR STRINGS IN A SLICE ===

// findSimilarStrings finds similar strings in a slice
func findSimilarStrings(target string, collection []string, maxResults int) []similarItem {
	var results []similarItem

	for i, item := range collection {
		if item == target {
			continue // skip exact matches, they have been treated
		}

		similarity := calculateStringSimilarity(target, item)
		if similarity.Similarity >= 0.6 { // threshold de 60%
			similarity.Index = i
			results = append(results, similarity)
		}
	}

	// sort by similarity (highest first)
	for i := 0; i < len(results)-1; i++ {
		for j := i + 1; j < len(results); j++ {
			if results[i].Similarity < results[j].Similarity {
				results[i], results[j] = results[j], results[i]
			}
		}
	}

	// limit results
	if len(results) > maxResults {
		results = results[:maxResults]
	}

	return results
}

// calculateStringSimilarity calcula similaridade entre duas strings
func calculateStringSimilarity(target, candidate string) similarItem {
	item := similarItem{
		Value: candidate,
	}

	// 0. Check exact match first
	if target == candidate {
		item.Similarity = 1.0
		return item
	}

	// 1. Check case sensitivity
	if strings.EqualFold(target, candidate) {
		item.Similarity = 0.95
		item.DiffType = "case"
		item.Details = "case difference"
		return item
	}

	// 2. Check prefix/suffix
	if strings.HasPrefix(candidate, target) {
		item.Similarity = 0.9
		item.DiffType = "prefix"
		extra := candidate[len(target):]
		item.Details = fmt.Sprintf("extra '%s'", extra)
		return item
	}

	if strings.HasSuffix(candidate, target) {
		item.Similarity = 0.9
		item.DiffType = "suffix"
		extra := candidate[:len(candidate)-len(target)]
		item.Details = fmt.Sprintf("prefix '%s'", extra)
		return item
	}

	if strings.HasPrefix(target, candidate) {
		item.Similarity = 0.85
		item.DiffType = "prefix"
		missing := target[len(candidate):]
		item.Details = fmt.Sprintf("missing '%s'", missing)
		return item
	}

	if strings.HasSuffix(target, candidate) {
		item.Similarity = 0.85
		item.DiffType = "suffix"
		missing := target[:len(target)-len(candidate)]
		item.Details = fmt.Sprintf("missing prefix '%s'", missing)
		return item
	}

	// 3. Check substring
	if strings.Contains(candidate, target) {
		item.Similarity = 0.8
		item.DiffType = "substring"
		item.Details = "target is substring of candidate"
		return item
	}

	if strings.Contains(target, candidate) {
		item.Similarity = 0.75
		item.DiffType = "substring"
		item.Details = "candidate is substring of target"
		return item
	}

	// 4. Calculate Damerau-Levenshtein distance for typos
	distance := damerauLevenshteinDistance(target, candidate)
	maxLen := max(len(target), len(candidate))

	// NOTE: maxLen == 0 is impossible here because:
	// - If both empty: already returned in "exact match" (target == candidate)
	// - Otherwise: at least one string is non-empty
	// This guarantees no division by zero.

	similarity := 1.0 - float64(distance)/float64(maxLen)

	if similarity >= 0.6 {
		item.Similarity = similarity
		item.DiffType = "typo"
		item.Details = generateTypoDetails(target, candidate, distance)
	}

	return item
}

// levenshteinDistance calculates the Levenshtein distance between two strings
func levenshteinDistance(s1, s2 string) int {
	len1, len2 := len(s1), len(s2)

	matrix := make([][]int, len1+1)
	for i := range matrix {
		matrix[i] = make([]int, len2+1)
	}

	// Initialize first row and column
	for i := 0; i <= len1; i++ {
		matrix[i][0] = i
	}
	for j := 0; j <= len2; j++ {
		matrix[0][j] = j
	}

	// Fill matrix
	for i := 1; i <= len1; i++ {
		for j := 1; j <= len2; j++ {
			cost := 0
			if s1[i-1] != s2[j-1] {
				cost = 1
			}

			matrix[i][j] = min3(
				matrix[i-1][j]+1,      // deletion
				matrix[i][j-1]+1,      // insertion
				matrix[i-1][j-1]+cost, // substitution
			)
		}
	}

	return matrix[len1][len2]
}

// damerauLevenshteinDistance calculates the Damerau-Levenshtein distance between two strings.
// Unlike standard Levenshtein, it treats transposition of adjacent characters as a single operation.
// For example, "tets" -> "test" has distance 1 (transposition) instead of 2 (delete + insert).
func damerauLevenshteinDistance(s1, s2 string) int {
	len1, len2 := len(s1), len(s2)

	// Create a 2D matrix
	matrix := make([][]int, len1+1)
	for i := range matrix {
		matrix[i] = make([]int, len2+1)
	}

	// Initialize first row and column
	for i := 0; i <= len1; i++ {
		matrix[i][0] = i
	}
	for j := 0; j <= len2; j++ {
		matrix[0][j] = j
	}

	// Fill matrix
	for i := 1; i <= len1; i++ {
		for j := 1; j <= len2; j++ {
			cost := 0
			if s1[i-1] != s2[j-1] {
				cost = 1
			}

			matrix[i][j] = min3(
				matrix[i-1][j]+1,      // deletion
				matrix[i][j-1]+1,      // insertion
				matrix[i-1][j-1]+cost, // substitution
			)

			// Transposition
			if i > 1 && j > 1 && s1[i-1] == s2[j-2] && s1[i-2] == s2[j-1] {
				matrix[i][j] = min(matrix[i][j], matrix[i-2][j-2]+1)
			}
		}
	}

	return matrix[len1][len2]
}

// generateTypoDetails generates detailed description of the error type
func generateTypoDetails(target, candidate string, distance int) string {
	if distance == 1 {
		// Try to identify the specific type of error
		if len(target) == len(candidate) {
			// Check if it's a transposition (adjacent characters swapped)
			for i := 0; i < len(target)-1; i++ {
				if target[i] == candidate[i+1] && target[i+1] == candidate[i] {
					// It's a transposition
					return "1 character differs"
				}
			}

			// Otherwise it's a simple substitution
			for i := 0; i < len(target); i++ {
				if target[i] != candidate[i] {
					return fmt.Sprintf("'%c' ≠ '%c' at position %d", candidate[i], target[i], i+1)
				}
			}
		} else if len(candidate) == len(target)+1 {
			return "1 extra character"
		} else if len(target) == len(candidate)+1 {
			return "1 missing character"
		}
	}

	return fmt.Sprintf("%d characters differ", distance)
}

// auxiliary function for contains of string slices
func containsString(target string, collection []string) containResult {
	const maxShow = 5
	const maxSimilar = 3

	result := containResult{
		MaxShow: maxShow,
		Total:   len(collection),
	}

	// 1. Check exact match
	for _, item := range collection {
		if item == target {
			result.Found = true
			result.Exact = true
			return result
		}
	}

	result.Similar = findSimilarStrings(target, collection, maxSimilar)

	// 2. Prepare context (first elements to show)
	contextSize := maxShow
	if len(collection) > contextSize {
		result.Context = make([]interface{}, contextSize)
		for i := 0; i < contextSize; i++ {
			result.Context[i] = collection[i]
		}
	} else {
		result.Context = make([]interface{}, len(collection))
		for i, item := range collection {
			result.Context[i] = item
		}
	}

	return result
}

func min3(a, b, c int) int {
	if a <= b && a <= c {
		return a
	}
	if b <= c {
		return b
	}
	return c
}

// auxiliary function min/max
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func formatContainsError(target interface{}, result containResult) string {
	var msg strings.Builder

	msg.WriteString("Expected collection to contain element:\n")

	// Show context of the collection
	msg.WriteString(fmt.Sprintf("        Collection: %v", formatComparisonValue(result.Context)))
	if len(result.Context) < result.Total {
		msg.WriteString(fmt.Sprintf(" (showing %d of %d)", len(result.Context), result.Total))
	}
	msg.WriteString("\n")

	msg.WriteString(fmt.Sprintf("        Missing   : %v\n", target))

	// Show similar if found
	if len(result.Similar) > 0 {
		msg.WriteString("\n")
		if len(result.Similar) == 1 {
			similar := result.Similar[0]
			msg.WriteString(fmt.Sprintf("        Found similar: %v (at index %d) - %s\n",
				similar.Value, similar.Index, similar.Details))
		} else {
			msg.WriteString("        Hint: Similar elements found:\n")
			for _, similar := range result.Similar {
				msg.WriteString(fmt.Sprintf("          └─ %v (at index %d) - %s\n",
					similar.Value, similar.Index, similar.Details))
			}
		}
	}

	return msg.String()
}

//  === THIS SECTION IS TO FIND SIMILAR INT IN A SLICE ===

func findInsertionInfo[T Ordered](collection []T, target T) (insertionInfo[T], error) {
	info := insertionInfo[T]{}

	if len(collection) == 0 {
		info.insertIndex = -1
		return info, nil
	}

	if isFloat(target) {
		if math.IsNaN(float64(target)) {
			return info, fmt.Errorf("NaN values are not supported")
		}
	}

	sortedCollection := make([]T, len(collection))
	copy(sortedCollection, collection)
	slices.Sort(sortedCollection)

	if len(sortedCollection) > 0 && isFloat(sortedCollection[0]) {
		for _, v := range sortedCollection {
			if math.IsNaN(float64(v)) {
				return info, fmt.Errorf("collection contains NaN values")
			}
		}
	}

	insertIndex := sort.Search(len(sortedCollection), func(i int) bool {
		return sortedCollection[i] >= target
	})

	info.insertIndex = insertIndex

	if insertIndex < len(sortedCollection) && sortedCollection[insertIndex] == target {
		info.found = true
		return info, nil
	}

	if insertIndex > 0 {
		info.prev = &sortedCollection[insertIndex-1]
	}

	if insertIndex < len(sortedCollection) {
		info.next = &sortedCollection[insertIndex]
	}

	if len(collection) > 10 {
		windowSize := 4
		leftSide := windowSize / 2
		rightSide := windowSize / 2

		startIndex := max(0, insertIndex-leftSide)
		endIndex := min(len(sortedCollection), insertIndex+rightSide)

		// Adjust window boundaries to maximize element count within windowSize limit
		actualSize := endIndex - startIndex
		if actualSize < windowSize {
			if startIndex == 0 {
				// Already at the beginning, expand right to fill window
				endIndex = min(len(sortedCollection), windowSize)
			} else if endIndex == len(sortedCollection) {
				// Already at the end, expand left to fill window
				startIndex = max(0, len(sortedCollection)-windowSize)
			}
		}

		window := sortedCollection[startIndex:endIndex]

		var builder strings.Builder
		builder.WriteString("[")
		if startIndex > 0 {
			builder.WriteString("..., ")
		}

		var elements []string
		for _, val := range window {
			elements = append(elements, fmt.Sprintf("%v", val))
		}
		builder.WriteString(strings.Join(elements, ", "))

		if endIndex < len(sortedCollection) {
			builder.WriteString(", ...")
		}
		builder.WriteString("]")
		info.sortedWindow = builder.String()
	}

	return info, nil
}

func formatInsertionContext[T Ordered](collection []T, target T, info insertionInfo[T]) string {
	collectionLength := len(collection)
	builder := strings.Builder{}

	if collectionLength == 0 {
		builder.WriteString("Collection: []\n")
		builder.WriteString("Missing  : ")
		builder.WriteString(fmt.Sprint(target))
		return builder.String()
	}

	builder.WriteString("Collection: ")

	var elements []string
	if collectionLength <= 10 {
		// Show all elements
		for _, item := range collection {
			elements = append(elements, fmt.Sprintf("%v", item))
		}
		builder.WriteString(fmt.Sprintf("[%s]", strings.Join(elements, ", ")))
	} else {
		// Show first 5
		for i := 0; i < 5; i++ {
			elements = append(elements, fmt.Sprintf("%v", collection[i]))
		}
		// Show last 5
		lastElements := []string{}
		for i := collectionLength - 5; i < collectionLength; i++ {
			lastElements = append(lastElements, fmt.Sprintf("%v", collection[i]))
		}
		builder.WriteString(fmt.Sprintf("[%s, ..., %s]", strings.Join(elements, ", "), strings.Join(lastElements, ", ")))
		builder.WriteString(fmt.Sprintf(" (showing first 5 and last 5 of %d elements)", collectionLength))
	}

	builder.WriteString("\nMissing  : ")
	builder.WriteString(fmt.Sprint(target))

	if info.prev != nil || info.next != nil {
		builder.WriteString("\n\n")
		if info.prev != nil && info.next != nil {
			builder.WriteString(fmt.Sprintf("Element %v would fit between %v and %v in sorted order", target, *info.prev, *info.next))
		} else if info.prev != nil {
			builder.WriteString(fmt.Sprintf("Element %v would be after %v in sorted order", target, *info.prev))
		} else if info.next != nil {
			builder.WriteString(fmt.Sprintf("Element %v would be before %v in sorted order", target, *info.next))
		}
	}

	if info.sortedWindow != "" {
		if info.prev == nil && info.next == nil {
			builder.WriteString("\n")
		}
		builder.WriteString(fmt.Sprintf("\n└─ Sorted view: %s", info.sortedWindow))
	}

	return builder.String()
}

func isFloat[T Ordered](v T) bool {
	switch any(v).(type) {
	case float32, float64:
		return true
	default:
		return false
	}
}

// formatEmptyError formats a detailed error message for empty/not empty assertions
func formatEmptyError(value interface{}, expectedEmpty bool) string {
	var msg strings.Builder

	if expectedEmpty {
		msg.WriteString("Expected value to be empty, but it was not:\n")
	} else {
		msg.WriteString("Expected value to be not empty, but it was empty:\n")
	}

	actualValue := reflect.ValueOf(value)

	// Handle nil or zero value case
	if !actualValue.IsValid() {
		msg.WriteString("        Type    : nil\n")
		msg.WriteString("        Value   : nil\n")
		return msg.String()
	}

	switch actualValue.Kind() {
	case reflect.String:

		if len(actualValue.String()) > 180 {
			msg.WriteString(formatMultilineString(actualValue.String()))
			return msg.String()
		}

		str := actualValue.String()
		msg.WriteString("        Type    : string\n")
		msg.WriteString(fmt.Sprintf("        Length  : %d characters\n", len(str)))
		if expectedEmpty && len(str) > 0 {
			if len(str) <= 50 {
				msg.WriteString(fmt.Sprintf("        Content : %q\n", str))
			} else {
				msg.WriteString(fmt.Sprintf("        Content : %q... (truncated)\n", str[:47]))
			}
		}

	case reflect.Slice, reflect.Array:
		length := actualValue.Len()
		msg.WriteString(fmt.Sprintf("        Type    : %s\n", actualValue.Type()))
		msg.WriteString(fmt.Sprintf("        Length  : %d elements\n", length))
		if expectedEmpty && length > 0 {
			if length <= 5 {
				msg.WriteString(fmt.Sprintf("        Content : %s\n", formatComparisonValue(value)))
			} else {
				// Show first 3 elements
				elements := make([]string, 3)
				for i := 0; i < 3; i++ {
					elements[i] = formatValueComparison(actualValue.Index(i))
				}
				msg.WriteString(fmt.Sprintf("        Content : [%s, ...] (showing first 3 of %d)\n",
					strings.Join(elements, ", "), length))
			}
		}

	case reflect.Map:
		length := actualValue.Len()
		msg.WriteString(fmt.Sprintf("        Type    : %s\n", actualValue.Type()))
		msg.WriteString(fmt.Sprintf("        Length  : %d entries\n", length))
		if expectedEmpty && length > 0 {
			if length <= 3 {
				msg.WriteString(fmt.Sprintf("        Content : %s\n", formatComparisonValue(value)))
			} else {
				msg.WriteString(fmt.Sprintf("        Content : map[...] (showing %d entries)\n", length))
			}
		}

	case reflect.Chan:
		msg.WriteString(fmt.Sprintf("        Type    : %s\n", actualValue.Type()))
		msg.WriteString("        Note    : Channel length cannot be determined\n")

	default:
		msg.WriteString(fmt.Sprintf("        Type    : %s\n", actualValue.Type()))
		msg.WriteString(fmt.Sprintf("        Value   : %s\n", formatComparisonValue(value)))
	}

	return msg.String()
}

// formatNumericComparisonError formats a detailed error message for numeric comparisons
func formatNumericComparisonError(actual, expected interface{}, operation string) string {
	var msg strings.Builder

	actualV := reflect.ValueOf(actual)
	expectedV := reflect.ValueOf(expected)

	actualFloat, _ := toFloat64(actualV)
	expectedFloat, _ := toFloat64(expectedV)

	difference := actualFloat - expectedFloat

	switch operation {
	case "greater":
		msg.WriteString("Expected value to be greater than threshold:\n")
	case "less":
		msg.WriteString("Expected value to be less than threshold:\n")
	case "greaterOrEqual":
		msg.WriteString("Expected value to be greater than or equal to threshold:\n")
	case "lessOrEqual":
		msg.WriteString("Expected value to be less than or equal to threshold:\n")
	}

	msg.WriteString(fmt.Sprintf("        Value     : %v\n", actual))
	msg.WriteString(fmt.Sprintf("        Threshold : %v\n", expected))

	if difference > 0 {
		msg.WriteString(fmt.Sprintf("        Difference: +%v (value is %v greater)\n", difference, difference))
	} else if difference < 0 {
		msg.WriteString(fmt.Sprintf("        Difference: %v (value is %v smaller)\n", difference, -difference))
	} else {
		msg.WriteString("        Difference: 0 (values are equal)\n")
	}

	// Add contextual hint
	switch operation {
	case "greater":
		if difference <= 0 {
			msg.WriteString("        Hint      : Value should be larger than threshold\n")
		}
	case "less":
		if difference >= 0 {
			msg.WriteString("        Hint      : Value should be smaller than threshold\n")
		}
	case "greaterOrEqual":
		if difference < 0 {
			msg.WriteString("        Hint      : Value should be larger than or equal to threshold\n")
		}
	case "lessOrEqual":
		if difference > 0 {
			msg.WriteString("        Hint      : Value should be smaller than or equal to threshold\n")
		}
	}

	return msg.String()
}

// formatBeWithinError returns an error message when `actual` is outside `expected ± tolerance`.
func formatBeWithinError[T Float](actual, expected, tolerance T) string {
	diff := T(math.Abs(float64(actual - expected)))

	actualF := float64(actual)
	expectedF := float64(expected)
	diffF := float64(diff)
	toleranceF := float64(tolerance)

	format := chooseFormat(actualF, expectedF, diffF, toleranceF)

	var msg strings.Builder

	msg.WriteString(fmt.Sprintf(
		"Expected "+format+" to be within ±"+format+" of "+format+"\n",
		actualF, toleranceF, expectedF))
	msg.WriteString(fmt.Sprintf("Difference: "+format, diffF))

	if toleranceF > 0 {
		excess := ((diffF - toleranceF) / toleranceF)
		switch {
		case excess > 2:
			msg.WriteString(fmt.Sprintf(" (%.2f× greater than tolerance)", excess))
		case excess > 0:
			msg.WriteString(fmt.Sprintf(" (%.2f%% greater than tolerance)", 100*excess))
		}
	}

	return msg.String()
}

// chooseFormat selects the best format based on the magnitude of the numbers
func chooseFormat(values ...float64) string {
	maxAbs := 0.0
	minNonZero := math.Inf(1)

	for _, v := range values {
		abs := math.Abs(v)
		if abs > maxAbs {
			maxAbs = abs
		}
		if abs > 0 && abs < minNonZero {
			minNonZero = abs
		}
	}

	// Use scientific notation only for extremely large or small numbers
	if maxAbs >= 1e6 || (minNonZero != math.Inf(1) && minNonZero < 1e-6) {
		return "%.6e"
	}

	return "%.6f"
}

// formatLengthError formats a detailed error message for HaveLength assertions.
func formatLengthError(actual any, expected, actualLen int) string {
	var msg strings.Builder
	msg.WriteString("Expected collection to have specific length:\n")
	msg.WriteString(fmt.Sprintf("Type          : %T\n", actual))
	msg.WriteString(fmt.Sprintf("Expected Length: %d\n", expected))
	msg.WriteString(fmt.Sprintf("Actual Length : %d\n", actualLen))

	diff := actualLen - expected
	if diff > 0 {
		elementWord := "elements"
		if diff == 1 {
			elementWord = "element"
		}
		msg.WriteString(fmt.Sprintf("Difference    : +%d (%d %s extra)\n", diff, diff, elementWord))
	} else {
		elementWord := "elements"
		if -diff == 1 {
			elementWord = "element"
		}
		msg.WriteString(fmt.Sprintf("Difference    : %d (%d %s missing)\n", diff, -diff, elementWord))
	}

	return msg.String()
}

// formatTypeError formats a detailed error message for BeOfType assertions.
func formatTypeError(expectedType, actualType reflect.Type) string {
	var msg strings.Builder
	msg.WriteString("Expected value to be of specific type:\n")
	msg.WriteString(fmt.Sprintf("Expected Type: %v\n", expectedType))
	msg.WriteString(fmt.Sprintf("Actual Type  : %v\n", actualType))
	msg.WriteString("Difference   : Different concrete types\n")

	return msg.String()
}

// formatOneOfError formats a detailed error message for BeOneOf assertions.
func formatOneOfError[T any](actual T, options []T) string {
	var msg strings.Builder
	msg.WriteString("Expected value to be one of the allowed options:\n")
	msg.WriteString(fmt.Sprintf("Value   : %s\n", formatComparisonValue(actual)))

	// Truncate options if there are more than 4
	msg.WriteString("Options : ")
	if len(options) <= 4 {
		msg.WriteString(formatComparisonValue(options))
	} else {
		// Show first 4 options with truncation indicator
		truncatedOptions := make([]T, 4)
		copy(truncatedOptions, options[:4])
		baseStr := formatComparisonValue(truncatedOptions)
		// Remove the closing bracket and add truncation indicator
		if strings.HasSuffix(baseStr, "]") {
			msg.WriteString(baseStr[:len(baseStr)-1])
			msg.WriteString(fmt.Sprintf(", ...] (showing first 4 of %d)", len(options)))
		} else {
			msg.WriteString(baseStr)
		}
	}
	msg.WriteString("\n")

	msg.WriteString(fmt.Sprintf("Count   : 0 of %d options matched\n", len(options)))
	return msg.String()
}

func findUnhashableDuplicates(collection any) []duplicateGroup {
	rv := reflect.ValueOf(collection)
	length := rv.Len()
	visitedIndices := make([]bool, length)
	var duplicates []duplicateGroup

	for i := 0; i < length; i++ {
		if visitedIndices[i] {
			continue
		}

		item := rv.Index(i).Interface()
		var foundIndices []int

		for j := i + 1; j < length; j++ {
			if visitedIndices[j] {
				continue
			}

			candidate := rv.Index(j).Interface()
			if reflect.DeepEqual(item, candidate) {
				if len(foundIndices) == 0 {
					foundIndices = append(foundIndices, i)
					visitedIndices[i] = true
				}
				foundIndices = append(foundIndices, j)
				visitedIndices[j] = true
			}
		}

		if len(foundIndices) > 0 {
			duplicates = append(duplicates, duplicateGroup{Value: item, Indexes: foundIndices})
		}
	}
	return duplicates
}

// findDuplicates finds duplicate values in a collection.
// It uses a fast path for comparable types and a fallback for unhashable types.
// It returns a slice of duplicate groups, each containing the value and its indexes.
func findDuplicates(collection any) []duplicateGroup {
	rv := reflect.ValueOf(collection)

	// Check if the type is comparable to use the fast path with maps
	if rv.Type().Elem().Comparable() {
		return findComparableDuplicates(collection)
	}

	// Fallback to deep equality for unhashable types
	return findUnhashableDuplicates(collection)
}

func findComparableDuplicates(collection any) []duplicateGroup {
	rv := reflect.ValueOf(collection)
	indexes := make(map[any][]int)
	length := rv.Len()

	for i := 0; i < length; i++ {
		item := rv.Index(i).Interface()
		indexes[item] = append(indexes[item], i)
	}

	var duplicates []duplicateGroup
	for item, idxs := range indexes {
		if len(idxs) > 1 {
			duplicates = append(duplicates, duplicateGroup{Value: item, Indexes: idxs})
		}
	}
	return duplicates
}

func formatDuplicatesErrors(duplicates []duplicateGroup) string {
	var msg strings.Builder

	for _, group := range duplicates {
		if len(group.Indexes) > 4 {
			windowMsg := formatIndexesWindow(group.Indexes, 4)

			msg.WriteString(fmt.Sprintf(
				"\n└─ %s appears %d times at indexes %v",
				formatDuplicateItem(group.Value),
				len(group.Indexes),
				windowMsg,
			))
			continue
		}

		msg.WriteString(fmt.Sprintf("\n└─ %s appears %d times at indexes %v",
			formatDuplicateItem(group.Value), len(group.Indexes), formatComparisonValue(group.Indexes)))
	}

	return msg.String()
}

func formatDuplicateItem(item any) string {
	if item == nil {
		return "nil"
	}

	rv := reflect.ValueOf(item)
	rt := reflect.TypeOf(item)

	// For structs, use special formatting
	if rv.Kind() == reflect.Struct {
		return formatStructForDuplicates(rv, rt)
	}

	// For other types, use the existing formatComparisonValue
	return formatComparisonValue(item)
}

func formatStructForDuplicates(rv reflect.Value, rt reflect.Type) string {
	var parts []string
	charCount := 0
	maxChars := 80

	typeName := rt.Name()
	if typeName == "" {
		typeName = "struct"
	}

	result := typeName + "{"

	for i := 0; i < rv.NumField(); i++ {
		field := rt.Field(i)
		fieldValue := rv.Field(i)

		if !field.IsExported() {
			continue
		}

		fieldStr := fmt.Sprintf("%s: %v", field.Name, formatFieldForDuplicates(fieldValue))

		if charCount+len(fieldStr) > maxChars && len(parts) > 0 {
			parts = append(parts, "...")
			break
		}

		parts = append(parts, fieldStr)
		charCount += len(fieldStr)
	}

	result += strings.Join(parts, ", ") + "}"
	return result
}

func formatFieldForDuplicates(rv reflect.Value) string {
	switch rv.Kind() {
	case reflect.String:
		str := rv.String()
		if len(str) > 20 {
			return fmt.Sprintf("%q", str[:17]+"...")
		}
		return fmt.Sprintf("%q", str)
	case reflect.Struct:
		return fmt.Sprintf("%s{...}", rv.Type().Name())
	case reflect.Map, reflect.Slice, reflect.Array:
		return fmt.Sprintf("%s(...)", rv.Type().String())
	default:
		return fmt.Sprintf("%v", rv.Interface())
	}
}

func formatIndexesWindow(indexes []int, windowSize int) string {
	windowMsg := strings.Builder{}

	// If the number of indexes is less than or equal to the window size, return the indexes as is
	if len(indexes) <= windowSize {
		return formatComparisonValue(indexes)
	}

	/* windowMsg.WriteString(fmt.Sprintf("[%d, %d, %d, %d, ...]", indexes[0], indexes[1], indexes[2], indexes[3])) */

	windowMsg.WriteString("[")
	for i := range windowSize {
		windowMsg.WriteString(fmt.Sprintf("%d, ", indexes[i]))
	}

	windowMsg.WriteString("...")
	windowMsg.WriteString("]")

	return windowMsg.String()
}

/* func convertSliceToAny(slice interface{}) []any {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		return nil
	}

	result := make([]any, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		result[i] = rv.Index(i).Interface()
	}
	return result
} */

func addPrefixHighlight(msg *strings.Builder, actual, expected string) {
	prefixLength := len(expected)
	if len(actual) >= prefixLength {
		fmt.Fprintf(msg, "\n            %s", strings.Repeat("^", prefixLength))
		msg.WriteString("\n          (actual prefix)")
	}
}

func addPrefixHighlightToEnd(msg *strings.Builder, actual, expected string) {
	prefixLength := len(expected)
	if len(actual) >= prefixLength {
		blanksToAdd := len(actual) - prefixLength
		blanks := strings.Repeat(" ", blanksToAdd)
		fmt.Fprintf(msg, "\n            ")
		msg.WriteString(blanks)
		msg.WriteString(strings.Repeat("^", prefixLength))
		msg.WriteString("\n")
		msg.WriteString("            ")
		msg.WriteString(blanks)
		msg.WriteString("(actual suffix)")
	}
}

func formatStartsWithError(actual string, expected string, startWith string, noteMsg string, cfg *Config) string {
	var msg strings.Builder

	if cfg.IgnoreCase && strings.HasPrefix(strings.ToLower(actual), strings.ToLower(expected)) {
		msg.WriteString(fmt.Sprintf("Expected string to start with '%s', but it starts with '%s'", expected, startWith))
		msg.WriteString(fmt.Sprintf("\nExpected : '%s'", expected))
		msg.WriteString(fmt.Sprintf("\nActual   : '%s'", actual))
		addPrefixHighlight(&msg, actual, expected)
		msg.WriteString(noteMsg)
		return msg.String()
	}

	if !strings.HasPrefix(actual, expected) {
		msg.WriteString(fmt.Sprintf("Expected string to start with '%s', but it starts with '%s'", expected, startWith))
		msg.WriteString(fmt.Sprintf("\nExpected : '%s'", expected))
		msg.WriteString(fmt.Sprintf("\nActual   : '%s'", actual))
		addPrefixHighlight(&msg, actual, expected)
		msg.WriteString(noteMsg)
		return msg.String()
	}

	return ""
}

// formatEndsWithError formats a detailed error message for EndWith assertions.
func formatEndsWithError(actual string, expected string, actualEndSufix string, noteMsg string, cfg *Config) string {
	var msg strings.Builder
	if cfg.IgnoreCase && strings.HasSuffix(strings.ToLower(actualEndSufix), strings.ToLower(expected)) {
		msg.WriteString(fmt.Sprintf("Expected string to end with '%s', but it ends with '%s'", expected, actualEndSufix))
		msg.WriteString(fmt.Sprintf("\nExpected : '%s'", expected))
		msg.WriteString(fmt.Sprintf("\nActual   : '%s'", actual))
		addPrefixHighlight(&msg, actual, expected)
		msg.WriteString(noteMsg)
		return msg.String()
	}

	if !strings.HasSuffix(actualEndSufix, expected) {
		msg.WriteString(fmt.Sprintf("Expected string to end with '%s', but it ends with '%s'", expected, actualEndSufix))
		msg.WriteString(fmt.Sprintf("\nExpected : '%s'", expected))
		msg.WriteString(fmt.Sprintf("\nActual   : '%s'", actual))
		addPrefixHighlightToEnd(&msg, actual, expected)
		msg.WriteString(noteMsg)
		return msg.String()
	}

	return ""
}

// findExactCaseMismatch finds the exact case mismatch for a substring within a string
// Returns the position and the found substring if there's an exact case-only difference
func findExactCaseMismatch(actual, substring string) caseMismatchResult {
	if len(substring) == 0 {
		return caseMismatchResult{Found: false, Index: -1, Substring: ""}
	}

	actualRunes := []rune(actual)
	substringRunes := []rune(substring)

	actualLower := strings.ToLower(actual)
	substringLower := strings.ToLower(substring)

	// Check if there's a case-insensitive match
	byteIndex := strings.Index(actualLower, substringLower)
	if byteIndex == -1 {
		return caseMismatchResult{Found: false, Index: -1, Substring: ""}
	}

	// Convert byte index to rune index
	runeIndex := utf8.RuneCountInString(actual[:byteIndex])

	foundRunes := actualRunes[runeIndex : runeIndex+len(substringRunes)]
	foundSubstring := string(foundRunes)

	// Check if it's exactly the same except for case
	if strings.EqualFold(foundSubstring, substring) && foundSubstring != substring {
		return caseMismatchResult{Found: true, Index: runeIndex, Substring: foundSubstring}
	}

	return caseMismatchResult{Found: false, Index: -1, Substring: ""}
}

// formatSimpleCaseMismatchError formats a simplified error message for exact case mismatches
func formatSimpleCaseMismatchError(substring, foundSubstring string, position int) string {
	// Clean empty strings for display
	displayNeedle := substring
	displayFound := foundSubstring

	var b strings.Builder
	fmt.Fprintf(&b,
		"Expected string to contain %q, but found case difference\n"+
			"Substring: %q\n"+
			"Found    : %q at position %d\n"+
			"Note: Case mismatch detected (use should.WithIgnoreCase() if intended)",
		displayNeedle, displayNeedle, displayFound, position,
	)
	return b.String()
}

// formatContainSubstringError formats a detailed error message for ContainSubstring assertions.
func formatContainSubstringError(actual string, substring string, noteMsg string) string {
	var msg strings.Builder

	// Clean empty strings for display
	displayActual := actual
	displayNeedle := substring

	if strings.TrimSpace(actual) == "" {
		displayActual = "<empty>"
	}

	if strings.TrimSpace(substring) == "" {
		displayNeedle = "<empty>"
	}

	msg.WriteString(fmt.Sprintf("Expected string to contain %q, but it was not found", displayNeedle))
	msg.WriteString(fmt.Sprintf("\nSubstring   : %q", displayNeedle))

	// Handle very long strings with multiline formatting
	if len(actual) > 200 || strings.Contains(actual, "\n") {
		msg.WriteString(fmt.Sprintf("\nActual   : (length: %d)", len(actual)))
		msg.WriteString(fmt.Sprintf("\n%s", formatMultilineString(actual)))
	} else {
		msg.WriteString(fmt.Sprintf("\nActual   : %q", displayActual))
	}

	// Find similar substrings if substring is reasonable size
	if len(substring) > 0 && len(substring) <= maxSimilarLen && len(actual) > 0 {
		similarSubstrings := findSimilarSubstrings(actual, substring)

		// Filter out suggestions with empty details (similarity < 0.6)
		// These are too different to be helpful
		var goodSuggestions []similarItem
		for _, sim := range similarSubstrings {
			if sim.Details != "" {
				goodSuggestions = append(goodSuggestions, sim)
			}
		}

		if len(goodSuggestions) > 0 {
			sim := goodSuggestions[0]
			msg.WriteString("\n\nSimilar substring found:")
			msg.WriteString(fmt.Sprintf("\n  └─ '%s' at position %d - %s", sim.Value, sim.Index, sim.Details))
		}
	}

	if len(substring) > maxSimilarLen {
		msg.WriteString(fmt.Sprintf("\nNote: Substring is %d characters long (too large for similarity search)", len(substring)))
	}

	msg.WriteString(noteMsg)
	return msg.String()
}

// findSimilarSubstrings finds substrings in text that are similar to the target substring
// using a sliding window approach with Levenshtein distance. Limited to substrings <= 20 chars for performance.
// Returns at most 2 suggestions with:
// - Up to 2 characters difference
// - No duplicates
// - Not substrings of each other
// - At least 85% of the target length
func findSimilarSubstrings(text string, substring string) []similarItem {
	if len(substring) == 0 || len(substring) > maxSimilarLen || len(text) == 0 {
		return nil
	}

	var results []similarItem
	needleLen := len(substring)
	minLength := 0.85 * float64(needleLen) // 85% minimum size

	// Check if target has spaces - if not, filter out candidates with spaces
	targetHasSpaces := strings.Contains(substring, " ")

	// Use sliding window to extract all possible substrings of substring length
	for i := 0; i <= len(text)-needleLen; i++ {
		candidate := text[i : i+needleLen]

		if candidate == substring {
			continue // Skip exact matches
		}

		// Skip candidates with spaces if target doesn't have spaces
		if !targetHasSpaces && strings.Contains(candidate, " ") {
			continue
		}

		distance := damerauLevenshteinDistance(substring, candidate)
		if distance <= 2 { // Maximum 2 characters difference
			similarity := calculateStringSimilarity(substring, candidate)
			similarity.Index = i
			results = append(results, similarity)
		}
	}

	// Also check substrings of different lengths (±1, ±2 characters)
	for offset := -2; offset <= 2; offset++ {
		if offset == 0 {
			continue // Already checked exact length
		}

		substringLen := needleLen + offset
		if float64(substringLen) < minLength || substringLen > len(text) {
			continue // Skip if smaller than 85% of target size
		}

		for i := 0; i <= len(text)-substringLen; i++ {
			candidate := text[i : i+substringLen]

			// Skip candidates with spaces if target doesn't have spaces
			if !targetHasSpaces && strings.Contains(candidate, " ") {
				continue
			}

			distance := damerauLevenshteinDistance(substring, candidate)
			if distance <= 2 { // Maximum 2 characters difference
				similarity := calculateStringSimilarity(substring, candidate)
				similarity.Index = i
				results = append(results, similarity)
			}
		}
	}

	// Remove duplicates
	results = removeDuplicateSimilarItems(results)

	// Sort by similarity (highest first)
	for i := 0; i < len(results)-1; i++ {
		for j := i + 1; j < len(results); j++ {
			if results[i].Similarity < results[j].Similarity {
				results[i], results[j] = results[j], results[i]
			}
		}
	}

	// Keep only the best match (simplest and clearest for users)
	if len(results) > 1 {
		results = results[:1]
	}

	return results
}

// removeSubstringMatches removes items that are substrings of other items in the list
// Strategy: When one string is a substring of another, keep the longer one if it has
// equal or better similarity. This removes redundant suggestions like "ROR" when "ERROR" exists.
func removeSubstringMatches(items []similarItem) []similarItem {
	if len(items) <= 1 {
		return items
	}

	// Mark items to remove
	toRemove := make(map[int]bool)

	for i := range items {
		itemStr := strings.TrimSpace(fmt.Sprint(items[i].Value))

		for j := range items {
			if i == j || toRemove[j] {
				continue
			}

			otherStr := strings.TrimSpace(fmt.Sprint(items[j].Value))

			// Skip if they are the same string
			if itemStr == otherStr {
				continue
			}

			// Check if item is a substring of other
			if strings.Contains(otherStr, itemStr) && len(itemStr) < len(otherStr) {
				// item is shorter and contained in other (e.g., "ROR" vs "ERROR", or "test" vs "testing")
				// Remove item if:
				// 1. other has better or equal similarity, OR
				// 2. similarities are very close (within similarityThreshold) - prefer complete string
				if items[j].Similarity >= items[i].Similarity ||
					math.Abs(items[j].Similarity-items[i].Similarity) < similarityThreshold {
					toRemove[i] = true
					break
				}
			}

			// Check if other is a substring of item
			if strings.Contains(itemStr, otherStr) && len(otherStr) < len(itemStr) {
				// other is shorter and contained in item (e.g., "test" (other) vs "testing" (item))
				// Remove item if other has significantly better similarity (diff > similarityThreshold)
				if items[j].Similarity-items[i].Similarity > similarityThreshold {
					toRemove[i] = true
					break
				}
			}
		}
	}

	// Build filtered list
	var filtered []similarItem
	for i, item := range items {
		if !toRemove[i] {
			filtered = append(filtered, item)
		}
	}

	return filtered
}

// removeDuplicateSimilarItems removes duplicate similar items based on trimmed value
// Prefers items without leading/trailing spaces
func removeDuplicateSimilarItems(items []similarItem) []similarItem {
	if len(items) <= 1 {
		return items
	}

	// Group by trimmed value only (ignore index to eliminate duplicates)
	groups := make(map[string][]similarItem)

	for _, item := range items {
		trimmedValue := strings.TrimSpace(fmt.Sprint(item.Value))
		groups[trimmedValue] = append(groups[trimmedValue], item)
	}

	// For each group, choose the best item
	var unique []similarItem
	for _, group := range groups {
		if len(group) == 1 {
			unique = append(unique, group[0])
			continue
		}

		// Prefer item without spaces, then higher similarity, then earlier position
		best := group[0]
		bestStr := fmt.Sprint(best.Value)

		for _, item := range group[1:] {
			itemStr := fmt.Sprint(item.Value)

			// Prefer strings without leading/trailing spaces
			bestHasSpaces := strings.TrimSpace(bestStr) != bestStr
			itemHasSpaces := strings.TrimSpace(itemStr) != itemStr

			if !itemHasSpaces && bestHasSpaces {
				best = item
				bestStr = itemStr
			} else if itemHasSpaces == bestHasSpaces {
				// If both have or don't have spaces, prefer higher similarity
				if item.Similarity > best.Similarity {
					best = item
					bestStr = itemStr
				} else if item.Similarity == best.Similarity && item.Index < best.Index {
					// If similarity is equal, prefer earlier position
					best = item
					bestStr = itemStr
				}
			}
		}

		unique = append(unique, best)
	}

	return unique
}

// formatMapValuesList formats a slice of interface{} values for map error messages
// This function handles interface{} elements properly by getting their concrete values
func formatMapValuesList(values []interface{}) string {
	if values == nil {
		return "nil"
	}

	if len(values) == 0 {
		return "[]"
	}

	var elements []string
	for _, value := range values {
		// For strings, use single quotes to match existing test expectations
		if str, ok := value.(string); ok {
			elements = append(elements, fmt.Sprintf("'%s'", str))
		} else if value == nil {
			elements = append(elements, "nil")
		} else {
			// Handle interface{} values by getting their concrete type
			v := reflect.ValueOf(value)
			if v.Kind() == reflect.Interface && !v.IsNil() {
				v = v.Elem()
			}
			elements = append(elements, formatValueComparison(v))
		}
	}

	// Sort elements to ensure deterministic order for tests
	// But preserve original order when nil values are present
	hasNil := false
	for _, element := range elements {
		if element == "nil" {
			hasNil = true
			break
		}
	}

	if !hasNil {
		// Only sort when there are no nil values to preserve deterministic order
		sort.Strings(elements)
	}

	return fmt.Sprintf("[%s]", strings.Join(elements, ", "))
}

// containsMapKey checks if a map contains a specific key with similarity detection
func containsMapKey(mapValue interface{}, targetKey interface{}) mapContainResult {
	const maxShow = 5
	const maxSimilar = 3

	result := mapContainResult{
		MaxShow: maxShow,
	}

	v := reflect.ValueOf(mapValue)
	if v.Kind() != reflect.Map {
		return result
	}

	if v.IsNil() {
		result.Total = 0
		result.Context = nil
		return result
	}

	keys := v.MapKeys()
	result.Total = len(keys)

	// Extract all keys as interface{}
	allKeys := make([]interface{}, len(keys))
	for i, key := range keys {
		allKeys[i] = key.Interface()
	}

	// Check exact match
	targetVal := reflect.ValueOf(targetKey)
	for _, key := range keys {
		if reflect.DeepEqual(key.Interface(), targetKey) {
			result.Found = true
			result.Exact = true
			return result
		}
	}

	// Prepare context (keys to show)
	contextSize := maxShow
	if len(allKeys) > contextSize {
		result.Context = allKeys[:contextSize]
	} else {
		result.Context = allKeys
	}

	// Find similar keys based on type
	if targetVal.Kind() == reflect.String {
		// Handle string keys with similarity detection
		stringKeys := []string{}
		for _, key := range allKeys {
			if keyStr, ok := key.(string); ok {
				stringKeys = append(stringKeys, keyStr)
			}
		}
		if len(stringKeys) > 0 {
			targetStr := targetKey.(string)
			similarItems := findSimilarStrings(targetStr, stringKeys, maxSimilar)
			for _, item := range similarItems {
				result.Similar = append(result.Similar, similarItem{
					Value:      item.Value,
					Index:      item.Index,
					Similarity: item.Similarity,
					DiffType:   item.DiffType,
					Details:    item.Details,
				})
			}
		}
	} else if isNumericValue(targetKey) {
		// Handle numeric keys with numeric similarity
		result.Similar = findSimilarNumericKeys(allKeys, targetKey, maxSimilar)
	}

	return result
}

// containsMapValue checks if a map contains a specific value with similarity detection
func containsMapValue(mapValue interface{}, targetValue interface{}) mapContainResult {
	const maxShow = 5
	const maxSimilar = 3
	const maxCloseMatches = 2

	result := mapContainResult{
		MaxShow: maxShow,
	}

	v := reflect.ValueOf(mapValue)
	if v.Kind() != reflect.Map {
		return result
	}

	if v.IsNil() {
		result.Total = 0
		result.Context = nil
		return result
	}

	keys := v.MapKeys()
	result.Total = len(keys)

	// Extract all values as interface{}
	allValues := make([]interface{}, len(keys))
	for i, key := range keys {
		allValues[i] = v.MapIndex(key).Interface()
	}

	// Check exact match
	targetVal := reflect.ValueOf(targetValue)
	for _, value := range allValues {
		if reflect.DeepEqual(value, targetValue) {
			result.Found = true
			result.Exact = true
			return result
		}
	}

	// Prepare context (values to show)
	contextSize := maxShow
	if len(allValues) > contextSize {
		result.Context = allValues[:contextSize]
	} else {
		result.Context = allValues
	}

	isComplex := targetVal.Kind() == reflect.Struct || (targetVal.Kind() == reflect.Ptr && targetVal.Elem().Kind() == reflect.Struct)
	if isComplex {
		// Find close matches for structs
		var closeMatches []struct {
			match closeMatch
			diffs int
		}

		for _, val := range allValues {
			diffs := findDifferences(targetValue, val)
			if len(diffs) > 0 {
				var diffStrings []string
				for _, d := range diffs {
					diffStrings = append(
						diffStrings,
						fmt.Sprintf(
							"%s (%v ≠ %v)",
							d.Path,
							formatDiffValueConcise(d.Expected),
							formatDiffValueConcise(d.Actual),
						),
					)
				}
				closeMatches = append(closeMatches, struct {
					match closeMatch
					diffs int
				}{
					match: closeMatch{Value: val, Differences: diffStrings},
					diffs: len(diffs),
				})
			}
		}

		// Sort by number of differences (fewer is better)
		sort.Slice(closeMatches, func(i, j int) bool {
			return closeMatches[i].diffs < closeMatches[j].diffs
		})

		// Take the top N close matches
		for i := 0; i < len(closeMatches) && i < maxCloseMatches; i++ {
			result.CloseMatches = append(result.CloseMatches, closeMatches[i].match)
		}
		return result
	}

	// Find similar values based on type
	if targetVal.Kind() == reflect.String {
		// Handle string values with similarity detection
		stringValues := []string{}
		for _, value := range allValues {
			if valueStr, ok := value.(string); ok {
				stringValues = append(stringValues, valueStr)
			}
		}
		if len(stringValues) > 0 {
			targetStr := targetValue.(string)
			similarItems := findSimilarStrings(targetStr, stringValues, maxSimilar)
			for _, item := range similarItems {
				result.Similar = append(result.Similar, similarItem{
					Value:      item.Value,
					Index:      item.Index,
					Similarity: item.Similarity,
					DiffType:   item.DiffType,
					Details:    item.Details,
				})
			}
		}
	} else if isNumericValue(targetValue) {
		// Handle numeric values with numeric similarity
		result.Similar = findSimilarNumericKeys(allValues, targetValue, maxSimilar)
	}

	return result
}

// isNumericValue checks if a value is numeric
func isNumericValue(v interface{}) bool {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

// findSimilarNumericKeys finds numeric keys/values that are similar to the target
func findSimilarNumericKeys(items []interface{}, target interface{}, maxResults int) []similarItem {
	var results []similarItem

	targetVal := reflect.ValueOf(target)
	targetFloat, targetOk := toFloat64(targetVal)
	if !targetOk {
		return results
	}

	for i, item := range items {
		itemVal := reflect.ValueOf(item)
		itemFloat, itemOk := toFloat64(itemVal)
		if !itemOk {
			continue
		}

		if itemFloat == targetFloat {
			continue // Skip exact matches
		}

		diff := itemFloat - targetFloat
		absUltraDiff := diff
		if diff < 0 {
			absUltraDiff = -diff
		}

		// Consider numbers similar if they're within reasonable range
		var similarity float64
		var details string

		if absUltraDiff <= 1 {
			similarity = 0.9
			if diff > 0 {
				details = fmt.Sprintf("differs by %.0f", absUltraDiff)
			} else {
				details = fmt.Sprintf("differs by %.0f", absUltraDiff)
			}
		} else if absUltraDiff <= 10 {
			similarity = 0.8
			if diff > 0 {
				details = fmt.Sprintf("differs by %.0f", absUltraDiff)
			} else {
				details = fmt.Sprintf("differs by %.0f", absUltraDiff)
			}
		} else {
			// Check if target digits are contained in the item
			targetStr := fmt.Sprintf("%.0f", targetFloat)
			itemStr := fmt.Sprintf("%.0f", itemFloat)
			if strings.Contains(itemStr, targetStr) {
				similarity = 0.7
				details = "contains target digits"
			} else if strings.Contains(targetStr, itemStr) {
				similarity = 0.65
				details = "target contains these digits"
			} else {
				continue // Not similar enough
			}
		}

		if similarity >= 0.6 {
			results = append(results, similarItem{
				Value:      item,
				Index:      i,
				Similarity: similarity,
				DiffType:   "numeric",
				Details:    details,
			})
		}
	}

	// Sort by similarity (highest first)
	for i := 0; i < len(results)-1; i++ {
		for j := i + 1; j < len(results); j++ {
			if results[i].Similarity < results[j].Similarity {
				results[i], results[j] = results[j], results[i]
			}
		}
	}

	// Limit results
	if len(results) > maxResults {
		results = results[:maxResults]
	}

	return results
}

// formatMapContainKeyError formats error message for ContainKey assertion
func formatMapContainKeyError(target interface{}, result mapContainResult) string {
	var msg strings.Builder

	// Format target with single quotes for strings to match test expectations
	var targetStr string
	if str, ok := target.(string); ok {
		targetStr = fmt.Sprintf("'%s'", str)
	} else {
		targetStr = formatComparisonValue(target)
	}

	msg.WriteString(fmt.Sprintf("Expected map to contain key %s, but key was not found\n", targetStr))

	// Show available keys - use formatMapValuesList for better formatting
	msg.WriteString("Available keys: ")
	msg.WriteString(formatMapValuesList(result.Context))
	if len(result.Context) < result.Total {
		msg.WriteString(fmt.Sprintf(" (showing %d of %d)", len(result.Context), result.Total))
	}
	msg.WriteString("\n")

	msg.WriteString(fmt.Sprintf("Missing: %s\n", targetStr))

	// Show similar keys if found
	if len(result.Similar) > 0 {
		msg.WriteString("\n")
		if len(result.Similar) == 1 {
			similar := result.Similar[0]
			var similarStr string
			if str, ok := similar.Value.(string); ok {
				similarStr = fmt.Sprintf("'%s'", str)
			} else {
				similarStr = formatComparisonValue(similar.Value)
			}
			msg.WriteString("Similar key found:\n")
			msg.WriteString(fmt.Sprintf("  └─ %s - %s\n", similarStr, similar.Details))
		} else {
			msg.WriteString("Similar keys found:\n")
			for _, similar := range result.Similar {
				var similarStr string
				if str, ok := similar.Value.(string); ok {
					similarStr = fmt.Sprintf("'%s'", str)
				} else {
					similarStr = formatComparisonValue(similar.Value)
				}
				msg.WriteString(fmt.Sprintf("  └─ %s - %s\n", similarStr, similar.Details))
			}
		}
	}

	return msg.String()
}

// formatMapContainValueError formats error message for ContainValue assertion
func formatMapContainValueError(target interface{}, result mapContainResult) string {
	var msg strings.Builder

	targetV := reflect.ValueOf(target)
	isComplex := targetV.Kind() == reflect.Struct || (targetV.Kind() == reflect.Ptr && targetV.Elem().Kind() == reflect.Struct)

	// Use new formatting for complex types (structs)
	if isComplex {
		var typeName string
		if targetV.Kind() == reflect.Ptr {
			typeName = targetV.Elem().Type().Name()
		} else {
			typeName = targetV.Type().Name()
		}

		if typeName == "" {
			typeName = "struct"
		}

		msg.WriteString("Expected map to contain value, but it was not found:\n")
		msg.WriteString(fmt.Sprintf("Collection: %d values of type %s\n", result.Total, typeName))
		msg.WriteString(fmt.Sprintf("Missing   : %s\n", formatComplexType(target)))

		if len(result.Context) > 0 {
			msg.WriteString("\nAvailable values:\n")
			for i, v := range result.Context {
				prefix := "├─"
				if i == len(result.Context)-1 {
					prefix = "└─"
				}
				msg.WriteString(fmt.Sprintf("%s %s\n", prefix, formatComplexType(v)))
			}
		}

		if len(result.CloseMatches) > 0 {
			msg.WriteString("\nClose matches:\n")
			for i, match := range result.CloseMatches {
				prefix := "├─"
				if i == len(result.CloseMatches)-1 {
					prefix = "└─"
				}
				msg.WriteString(fmt.Sprintf("%s Match #%d: %s\n", prefix, i+1, formatComplexType(match.Value)))
				for _, diff := range match.Differences {
					msg.WriteString(fmt.Sprintf("│   └─ Differs in: %s\n", diff))
				}
			}
		}

		return strings.TrimSuffix(msg.String(), "\n")
	}

	// Format target with single quotes for strings to match test expectations
	var targetStr string
	if str, ok := target.(string); ok {
		targetStr = fmt.Sprintf("'%s'", str)
	} else {
		targetStr = formatComparisonValue(target)
	}

	msg.WriteString(fmt.Sprintf("Expected map to contain value %s, but value was not found\n", targetStr))

	// Show available values - use formatMapValuesList for better formatting
	msg.WriteString("Available values: ")
	msg.WriteString(formatMapValuesList(result.Context))
	if len(result.Context) < result.Total {
		msg.WriteString(fmt.Sprintf(" (showing %d of %d)", len(result.Context), result.Total))
	}
	msg.WriteString("\n")

	msg.WriteString(fmt.Sprintf("Missing: %s\n", targetStr))

	// Show similar values if found
	if len(result.Similar) > 0 {
		msg.WriteString("\n")
		if len(result.Similar) == 1 {
			similar := result.Similar[0]
			var similarStr string
			if str, ok := similar.Value.(string); ok {
				similarStr = fmt.Sprintf("'%s'", str)
			} else {
				similarStr = formatComparisonValue(similar.Value)
			}
			msg.WriteString("Similar value found:\n")
			msg.WriteString(fmt.Sprintf("  └─ %s - %s\n", similarStr, similar.Details))
		} else {
			msg.WriteString("Similar values found:\n")
			for _, similar := range result.Similar {
				var similarStr string
				if str, ok := similar.Value.(string); ok {
					similarStr = fmt.Sprintf("'%s'", str)
				} else {
					similarStr = formatComparisonValue(similar.Value)
				}
				msg.WriteString(fmt.Sprintf("  └─ %s - %s\n", similarStr, similar.Details))
			}
		}
	}

	return msg.String()
}

// formatComplexType formats a complex type (like a struct) with truncation for better readability.
func formatComplexType(item any) string {
	if item == nil {
		return "nil"
	}

	rv := reflect.ValueOf(item)

	// If it's a pointer, dereference it
	if rv.Kind() == reflect.Ptr && !rv.IsNil() {
		rv = rv.Elem()
	}

	rt := rv.Type()

	if rv.Kind() == reflect.Struct {
		return formatStructWithTruncation(rv, rt)
	}

	// Fallback for non-struct types
	return formatComparisonValue(item)
}

// formatStructWithTruncation creates a truncated string representation of a struct.
func formatStructWithTruncation(rv reflect.Value, rt reflect.Type) string {
	var parts []string
	charCount := 0
	maxChars := 80 // Same as formatStructForDuplicates

	typeName := rt.Name()
	if typeName == "" {
		typeName = "struct"
	}

	result := typeName + "{"

	for i := 0; i < rv.NumField(); i++ {
		field := rt.Field(i)
		fieldValue := rv.Field(i)

		if !field.IsExported() {
			continue
		}

		fieldStr := fmt.Sprintf("%s: %v", field.Name, formatFieldWithTruncation(fieldValue))

		if charCount+len(fieldStr) > maxChars && len(parts) > 0 {
			parts = append(parts, "...")
			break
		}

		parts = append(parts, fieldStr)
		charCount += len(fieldStr)
	}

	result += strings.Join(parts, ", ") + "}"
	return result
}

// formatFieldWithTruncation creates a truncated string representation of a field's value.
func formatFieldWithTruncation(rv reflect.Value) string {
	if !rv.IsValid() {
		return "nil"
	}

	switch rv.Kind() {
	case reflect.String:
		str := rv.String()
		if len(str) > 20 {
			return fmt.Sprintf("%q", str[:17]+"...")
		}
		return fmt.Sprintf("%q", str)
	case reflect.Ptr:
		if rv.IsNil() {
			return "nil"
		}
		return formatFieldWithTruncation(rv.Elem())
	case reflect.Struct:
		if !rv.IsValid() || rv.IsZero() {
			return fmt.Sprintf("%s{}", rv.Type().Name())
		}
		return fmt.Sprintf("%s{...}", rv.Type().Name())
	case reflect.Map, reflect.Slice, reflect.Array:
		if rv.IsNil() {
			return "nil"
		}
		if rv.Len() == 0 {
			return fmt.Sprintf("%s{}", rv.Type().String())
		}
		return fmt.Sprintf("%s(%d items)", rv.Type().String(), rv.Len())
	default:
		return fmt.Sprintf("%v", rv.Interface())
	}
}

// formatDiffValueConcise formats a value for difference display with truncation for readability.
func formatDiffValueConcise(value interface{}) string {
	if value == nil {
		return "nil"
	}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String:
		str := v.String()
		if len(str) > 30 {
			return fmt.Sprintf("%q", str[:27]+"...")
		}
		return fmt.Sprintf("%q", str)
	case reflect.Map:
		if v.IsNil() {
			return "nil"
		}
		if v.Len() == 0 {
			return "map[]"
		}
		if v.Len() == 1 {
			// Show single entry maps completely
			keys := v.MapKeys()
			key := keys[0]
			val := v.MapIndex(key)
			return fmt.Sprintf("map[%v: %v]", formatDiffValueConcise(key.Interface()), formatDiffValueConcise(val.Interface()))
		}
		// For maps with multiple entries, show count
		return fmt.Sprintf("map[%d entries]", v.Len())
	case reflect.Slice, reflect.Array:
		if v.IsNil() {
			return "nil"
		}
		if v.Len() == 0 {
			return "[]"
		}
		if v.Len() <= 3 {
			// Show small slices completely
			var elements []string
			for i := 0; i < v.Len(); i++ {
				elements = append(elements, formatDiffValueConcise(v.Index(i).Interface()))
			}
			return fmt.Sprintf("[%s]", strings.Join(elements, ", "))
		}
		// For larger slices, show count
		return fmt.Sprintf("[%d items]", v.Len())
	case reflect.Struct:
		// Use the existing complex type formatting
		return formatComplexType(value)
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return fmt.Sprint(value)
	default:
		// For other types, use a simple representation
		str := fmt.Sprint(value)
		if len(str) > 40 {
			return str[:37] + "..."
		}
		return str
	}
}

// formatMapNotContainKeyError formats error message for NotContainKey assertion
func formatMapNotContainKeyError(target interface{}, mapValue interface{}) string {
	var msg strings.Builder

	v := reflect.ValueOf(mapValue)
	mapType := v.Type().String()
	mapSize := v.Len()

	// Format target with single quotes for strings to match test expectations
	var targetStr string
	if str, ok := target.(string); ok {
		targetStr = fmt.Sprintf(`"%s"`, str)
	} else {
		targetStr = formatComparisonValue(target)
	}

	msg.WriteString("Expected map to NOT contain key, but key was found:\n")
	msg.WriteString(fmt.Sprintf("Map Type : %s\n", mapType))
	msg.WriteString(fmt.Sprintf("Map Size : %d entries\n", mapSize))
	msg.WriteString(fmt.Sprintf("Found Key: %s\n", targetStr))

	// Find the key and show its associated value
	keys := v.MapKeys()
	for _, key := range keys {
		if reflect.DeepEqual(key.Interface(), target) {
			value := v.MapIndex(key)
			var valueStr string
			if str, ok := value.Interface().(string); ok {
				valueStr = fmt.Sprintf(`"%s"`, str)
			} else {
				valueStr = formatComparisonValue(value.Interface())
			}
			msg.WriteString(fmt.Sprintf("Associated Value: %s\n", valueStr))
			break
		}
	}

	return strings.TrimSuffix(msg.String(), "\n")
}

// formatMapNotContainValueError formats error message for NotContainValue assertion
func formatMapNotContainValueError(target interface{}, mapValue interface{}) string {
	var msg strings.Builder

	v := reflect.ValueOf(mapValue)
	mapType := v.Type().String()
	mapSize := v.Len()

	msg.WriteString("Expected map to NOT contain value, but it was found:\n")
	msg.WriteString(fmt.Sprintf("Map Type : %s\n", mapType))
	msg.WriteString(fmt.Sprintf("Map Size : %d entries\n", mapSize))

	// Format the found value
	msg.WriteString(fmt.Sprintf("Found Value: %s\n", formatComplexType(target)))

	// Find which key(s) contain this value
	keys := v.MapKeys()
	var foundKeys []string
	for _, key := range keys {
		val := v.MapIndex(key)
		if reflect.DeepEqual(val.Interface(), target) {
			foundKeys = append(foundKeys, formatComparisonValue(key.Interface()))
		}
	}

	if len(foundKeys) == 1 {
		msg.WriteString(fmt.Sprintf("Found At: key %s", foundKeys[0]))
	} else if len(foundKeys) > 1 {
		if len(foundKeys) <= 3 {
			msg.WriteString(fmt.Sprintf("Found At: keys %s", strings.Join(foundKeys, ", ")))
		} else {
			msg.WriteString(fmt.Sprintf("Found At: %d keys (%s, ...)", len(foundKeys), strings.Join(foundKeys[:2], ", ")))
		}
	}

	return msg.String()
}

func formatRangeError[T Ordered](actual, minValue, maxValue T) string {
	if actual < minValue {
		return fmt.Sprintf("Expected value to be in range [%v, %v], but it was below:"+
			"\n        Value    : %v"+
			"\n        Range    : [%v, %v]"+
			"\n        Distance : %v below minimum (%v < %v)"+
			"\n        Hint     : Value should be >= %v",
			minValue, maxValue, actual, minValue, maxValue, minValue-actual, actual, minValue, minValue)
	}

	return fmt.Sprintf("Expected value to be in range [%v, %v], but it was above:"+
		"\n        Value    : %v"+
		"\n        Range    : [%v, %v]"+
		"\n        Distance : %v above maximum (%v > %v)"+
		"\n        Hint     : Value should be <= %v",
		minValue, maxValue, actual, minValue, maxValue, actual-maxValue, actual, maxValue, maxValue)
}

// formatNotPanicError formats a detailed error message for NotPanic assertions
func formatNotPanicError(panicInfo panicInfo, cfg *Config) string {
	var messageBuilder strings.Builder
	messageBuilder.WriteString("Expected for the function to not panic, but it panicked with: ")
	messageBuilder.WriteString(fmt.Sprintf("%v", panicInfo.Recovered))

	if cfg.StackTrace && panicInfo.Stack != "" {
		messageBuilder.WriteString("\nStack trace:\n")
		messageBuilder.WriteString(panicInfo.Stack)
	}

	return messageBuilder.String()
}

// checkIfSorted verifies if a slice is sorted in ascending order using generics
// Uses slices.IsSorted for fast primary check, then detailed analysis only if violations exist
func checkIfSorted[T Sortable](collection []T) sortCheckResult {
	length := len(collection)
	if length <= 1 {
		return sortCheckResult{
			IsSorted:   true,
			Violations: nil,
			Total:      length,
		}
	}

	if slices.IsSorted(collection) {
		return sortCheckResult{
			IsSorted:   true,
			Violations: nil,
			Total:      length,
		}
	}

	// Only do detailed analysis if we know there are violations
	var violations []sortViolation
	maxViolations := 6

	for i := 0; i < length-1; i++ {
		current := collection[i]
		next := collection[i+1]

		if current > next {
			violations = append(violations, sortViolation{
				Index: i,
				Value: current,
				Next:  next,
			})

			if len(violations) >= maxViolations {
				break
			}
		}
	}

	return sortCheckResult{
		IsSorted:   false,
		Violations: violations,
		Total:      length,
	}
}

// formatSortError creates a detailed error message for BeSorted failures using generics
func formatSortError(result sortCheckResult) string {
	var msg strings.Builder

	if len(result.Violations) == 0 {
		return ""
	}

	msg.WriteString("Expected collection to be in ascending order, but it is not:\n")

	collectionInfo := fmt.Sprintf("Collection: (total: %d elements)\n", result.Total)
	if result.Total > 100 {
		collectionInfo = fmt.Sprintf("Collection: [Large collection] (total: %d elements)\n", result.Total)
	}
	msg.WriteString(collectionInfo)

	violationCount := len(result.Violations)
	statusText := "Status    : 1 order violation found\n"
	if violationCount != 1 {
		statusText = fmt.Sprintf("Status    : %d order violations found\n", violationCount)
	}
	msg.WriteString(statusText)

	msg.WriteString("Problems  :\n")

	// Show up to 5 violations to avoid overwhelming output
	maxShow := min(violationCount, 5)
	for i := 0; i < maxShow; i++ {
		violation := result.Violations[i]
		msg.WriteString(fmt.Sprintf("  - Index %d: %v > %v\n",
			violation.Index, violation.Value, violation.Next))
	}

	remaining := violationCount - maxShow
	if remaining > 0 {
		remainingText := "  - ... and 1 more violation"
		if remaining != 1 {
			remainingText = fmt.Sprintf("  - ... and %d more violations", remaining)
		}
		msg.WriteString(remainingText)
	}

	return msg.String()
}

// formatBeSameTimeError builds a friendly error message for time equality comparisons.
func formatBeSameTimeError(expected time.Time, actual time.Time, diff time.Duration) string {
	var msg strings.Builder

	// Human readable difference text (e.g., 2.5s)
	human := humanizeDuration(diff)

	// Determine whether actual is later or earlier than expected
	relation := "later"
	if actual.Before(expected) {
		relation = "earlier"
	}

	msg.WriteString(fmt.Sprintf("Expected times to be the same, but difference is %s\n", human))
	msg.WriteString(fmt.Sprintf("Expected: %s\n", formatTimeForDisplay(expected)))
	msg.WriteString(fmt.Sprintf("Actual  : %s (%s %s)", formatTimeForDisplay(actual), human, relation))
	return msg.String()
}

// humanizeDuration returns a concise string representation of a duration,
// such as 2s, 2.5s, 150ms, 10m30s, or 1d1h.
//
// Negative durations are handled by returning their positive equivalent.
func humanizeDuration(duration time.Duration) string {
	if duration < 0 {
		duration = -duration
	}

	// Handle short durations with millisecond and second precision
	if duration < time.Millisecond {
		return fmt.Sprintf("%.3fms", float64(duration)/float64(time.Millisecond))
	}
	if duration < time.Second {
		ms := float64(duration) / float64(time.Millisecond)
		if ms-math.Floor(ms) == 0 {
			return fmt.Sprintf("%.0fms", ms)
		}
		return fmt.Sprintf("%.1fms", ms)
	}
	if duration < time.Minute {
		s := float64(duration) / float64(time.Second)
		if s-math.Floor(s) == 0 {
			return fmt.Sprintf("%.0fs", s)
		}
		return fmt.Sprintf("%.1fs", s)
	}

	// Handle longer durations with discrete units (minutes, hours, days)
	switch {
	case duration < time.Hour:
		minutes := duration / time.Minute
		seconds := (duration % time.Minute) / time.Second
		if seconds == 0 {
			return fmt.Sprintf("%dm", minutes)
		}
		return fmt.Sprintf("%dm%ds", minutes, seconds)
	case duration < 24*time.Hour:
		hours := duration / time.Hour
		minutes := (duration % time.Hour) / time.Minute
		if minutes == 0 {
			return fmt.Sprintf("%dh", hours)
		}
		return fmt.Sprintf("%dh%dm", hours, minutes)
	default:
		days := duration / (24 * time.Hour)
		hours := (duration % (24 * time.Hour)) / time.Hour
		if hours == 0 {
			return fmt.Sprintf("%dd", days)
		}
		return fmt.Sprintf("%dd%dh", days, hours)
	}
}

// formatTimeForDisplay formats a time.Time value for display, including
// fractional seconds up to nanosecond precision and the time zone name.
//
// Fractional seconds are trimmed of trailing zeros. The time zone is
// displayed as "UTC" for both UTC and the system's local time for
// consistent output, particularly in tests.
//
// Example output: "2006-01-02 15:04:05.123456789 UTC"
//
// Example output: "2006-01-02 15:04:05.5 EST"
func formatTimeForDisplay(t time.Time) string {
	// Use UTC for the base time to ensure consistent formatting.
	utcTime := t.UTC()

	// Format the base date and time string.
	baseFormat := "2006-01-02 15:04:05"
	formattedBase := utcTime.Format(baseFormat)

	// Build the fractional seconds part, if needed.
	fractionalPart := ""
	nanoseconds := utcTime.Nanosecond()
	if nanoseconds != 0 {
		// Format nanoseconds with leading zeros to 9 digits, then trim trailing zeros.
		fractionalStr := fmt.Sprintf("%09d", nanoseconds)
		fractionalStr = strings.TrimRight(fractionalStr, "0")
		fractionalPart = "." + fractionalStr
	}

	// Determine the time zone name. We use "UTC" for time.Local and time.UTC
	// to make test output consistent regardless of the machine's time zone.
	timeZoneName := t.Location().String()
	if t.Location() == time.UTC || t.Location() == time.Local {
		timeZoneName = "UTC"
	}

	return fmt.Sprintf("%s%s %s", formattedBase, fractionalPart, timeZoneName)
}

func formatBeErrorMessage(action string, err error, target interface{}) string {
	var msg strings.Builder

	var types []string
	unwrappedErr := err
	for unwrappedErr != nil {
		types = append(types, reflect.TypeOf(unwrappedErr).String())
		unwrappedErr = errors.Unwrap(unwrappedErr)
	}

	switch action {
	case "as":
		msg.WriteString(fmt.Sprintf("Expected error to be %T, but type not found in error chain\n", target))
	case "is":
		msg.WriteString(fmt.Sprintf("Expected error to be \"%s\", but not found in error chain\n", target))
	default:
		msg.WriteString("Assertion failed with an unknown type of error.\n")
	}

	msg.WriteString(fmt.Sprintf("Error: \"%s\"\n", err.Error()))
	msg.WriteString(fmt.Sprintf("Types  : [%s]", strings.Join(types, ", ")))

	return msg.String()
}

func formatNotBeErrorMessage(err error) string {
	var msg strings.Builder

	fmt.Fprintf(&msg, "Expected no error, but got an error\n")
	fmt.Fprintf(&msg, "Error: \"%s\"\n", err.Error())
	fmt.Fprintf(&msg, "Type: %T", err)

	return msg.String()
}
