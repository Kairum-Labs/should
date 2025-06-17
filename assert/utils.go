package assert

import (
	"fmt"
	"math"
	"reflect"
	"slices"
	"sort"
	"strings"
)

// isSliceOrArray checks if the provided value is a slice or an array.
// It handles nil values by returning false.
func isSliceOrArray(v interface{}) bool {
	if v == nil {
		return false
	}
	kind := reflect.TypeOf(v).Kind()
	return kind == reflect.Slice || kind == reflect.Array
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
		reflect.Float32, reflect.Float64,
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

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
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
				Expected: expectedValue.Interface(),
				Actual:   actualValue.Interface(),
			})
			return
		}

		//compare elements one by one
		for i := 0; i < expectedValue.Len(); i++ {
			if !reflect.DeepEqual(expectedValue.Index(i).Interface(), actualValue.Index(i).Interface()) {
				elementPath := buildPath(path, fmt.Sprintf("[%d]", i))
				diffs = append(diffs, compareExpectedActual(expectedValue.Index(i).Interface(), actualValue.Index(i).Interface(), elementPath)...)
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
func findSimilarStrings(target string, collection []string, maxResults int) []SimilarItem {
	var results []SimilarItem

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
func calculateStringSimilarity(target, candidate string) SimilarItem {
	item := SimilarItem{
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

	// 4. Calculate Levenshtein distance
	distance := levenshteinDistance(target, candidate)
	maxLen := max(len(target), len(candidate))

	if maxLen == 0 {
		item.Similarity = 1.0
		return item
	}

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

// generateTypoDetails generates detailed description of the error type
func generateTypoDetails(target, candidate string, distance int) string {
	if distance == 1 {
		// Try to identify the specific type of error
		if len(target) == len(candidate) {
			// Substitution of a character
			for i := 0; i < len(target); i++ {
				if target[i] != candidate[i] {
					return fmt.Sprintf("'%c' ≠ '%c' at position %d", candidate[i], target[i], i+1)
				}
			}
		} else if len(candidate) == len(target)+1 {
			return "1 extra char"
		} else if len(target) == len(candidate)+1 {
			return "1 missing char"
		}
	}

	return fmt.Sprintf("%d char diff", distance)
}

// auxiliary function for contains of string slices
func containsString(target string, collection []string) ContainResult {
	const maxShow = 5
	const maxSimilar = 3

	result := ContainResult{
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

func formatContainsError(target interface{}, result ContainResult) string {
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
			msg.WriteString("        Hint: Possible typo detected")
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

// findInsertionContext returns the context for inserting a target into a sorted version of the collection.
//
// If the target exists: returns ("", index).
// If not: returns (window string, insertion index). The window shows up to 4 nearby sorted values.
//
// Example:
//
//	findInsertionContext([]int{10, 20, 40}, 30) => ("[..., 10, 20, 40]", 2)
func findInsertionContext[T Ordered](collection []T, target T) (string, int) {
	if len(collection) == 0 {
		return "", -1
	}

	if isFloat(target) {
		if math.IsNaN(float64(target)) {
			return "error: NaN values are not supported", -1
		}
	}

	sortedCollection := make([]T, len(collection))
	copy(sortedCollection, collection)
	slices.Sort(sortedCollection)

	if len(sortedCollection) > 0 && isFloat(sortedCollection[0]) {
		for _, v := range sortedCollection {
			if math.IsNaN(float64(v)) {
				return "error: collection contains NaN values", -1
			}
		}
	}

	//we need to find the index where the target should be inserted in the sorted collection
	insertIndex := sort.Search(len(sortedCollection), func(i int) bool {
		return sortedCollection[i] >= target
	})

	if insertIndex < len(sortedCollection) && sortedCollection[insertIndex] == target {
		return "", insertIndex
	}

	windowSize := 4

	leftSide := windowSize / 2
	rightSide := windowSize / 2

	startIndex := max(0, insertIndex-leftSide)
	endIndex := min(len(sortedCollection), insertIndex+rightSide)

	// Adjust window boundaries to maximize element count within windowSize limit
	if len(sortedCollection) <= windowSize {
		// For small collections, show all elements
		startIndex = 0
		endIndex = len(sortedCollection)
	} else {
		// For larger collections, ensure we have windowSize elements when possible
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
	}

	// extracting elements from the collection
	window := sortedCollection[startIndex:endIndex]

	builder := strings.Builder{}

	builder.WriteString("[")

	if startIndex > 0 {
		builder.WriteString("..., ")
	}

	for i, val := range window {
		builder.WriteString(fmt.Sprintf("%v", val))
		if i < len(window)-1 {
			builder.WriteString(", ")
		}
	}

	if endIndex < len(sortedCollection) {
		builder.WriteString(", ...")
	}

	builder.WriteString("]")

	return builder.String(), insertIndex
}

func formatInsertionContext[T Ordered](collection []T, target T, window string) string {
	collectionLength := len(collection)
	builder := strings.Builder{}

	if collectionLength == 0 {
		builder.WriteString("\nCollection: []")
		builder.WriteString("\nMissing  : ")
		builder.WriteString(fmt.Sprint(target))
		return builder.String()
	}

	windowSize := 4
	var elementsShown int

	if collectionLength <= windowSize {
		// Small collections show all elements
		elementsShown = collectionLength
	} else {
		// Large collections show up to windowSize elements
		elementsShown = windowSize
	}

	builder.WriteString("\nCollection: ")
	builder.WriteString(window)
	if collectionLength > elementsShown {
		builder.WriteString(fmt.Sprintf(" (showing %d of %d elements)", elementsShown, collectionLength))
	}
	builder.WriteString("\nMissing  : ")
	builder.WriteString(fmt.Sprint(target))

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
			return formatMultilineString(actualValue.String())
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
