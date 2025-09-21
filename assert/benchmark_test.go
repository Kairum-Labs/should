package assert

import (
	"strings"
	"testing"
)

// Test data for benchmarks
var (
	// Simple types for BeEqual
	intValue1    = 42
	intValue2    = 42
	intValue3    = 43
	stringValue1 = "hello world"
	stringValue2 = "hello world"
	stringValue3 = "hello universe"

	// Complex struct for BeEqual
	complexStruct1 = struct {
		ID       int
		Name     string
		Active   bool
		Tags     []string
		Settings map[string]interface{}
	}{
		ID:     1,
		Name:   "Test Project",
		Active: true,
		Tags:   []string{"go", "testing", "benchmark"},
		Settings: map[string]interface{}{
			"debug":   true,
			"level":   2,
			"timeout": 30,
		},
	}

	complexStruct2 = struct {
		ID       int
		Name     string
		Active   bool
		Tags     []string
		Settings map[string]interface{}
	}{
		ID:     1,
		Name:   "Test Project",
		Active: true,
		Tags:   []string{"go", "testing", "benchmark"},
		Settings: map[string]interface{}{
			"debug":   true,
			"level":   2,
			"timeout": 30,
		},
	}

	// Strings for ContainSubstring benchmarks
	shortText    = "Hello, World! This is a short test string for benchmarking."
	mediumText   = strings.Repeat("Lorem ipsum dolor sit amet, consectetur adipiscing elit. ", 10)
	longText     = strings.Repeat("This is a very long text used for performance testing of substring operations. ", 100)
	veryLongText = strings.Repeat("Performance testing with ContainSubstring function using very long strings "+
		"to measure efficiency. ", 1000)
	needleShort    = "World"
	needleMedium   = "consectetur"
	needleLong     = "substring operations"
	needleVeryLong = "ContainSubstring function"
	needleNotFound = "nonexistent"
	needleCaseDiff = "WORLD"
)

// BeEqual Benchmarks - Critical function #1
func BenchmarkBeEqual_Primitives_Int_Same(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		BeEqual(b, intValue1, intValue2)
	}
}

func BenchmarkBeEqual_Primitives_Int_Different(b *testing.B) {
	b.ReportAllocs()
	mt := &mockT{}
	for i := 0; i < b.N; i++ {
		BeEqual(mt, intValue1, intValue3)
	}
}

func BenchmarkBeEqual_Primitives_String_Same(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		BeEqual(b, stringValue1, stringValue2)
	}
}

func BenchmarkBeEqual_Primitives_String_Different(b *testing.B) {
	b.ReportAllocs()
	mt := &mockT{}
	for i := 0; i < b.N; i++ {
		BeEqual(mt, stringValue1, stringValue3)
	}
}

func BenchmarkBeEqual_ComplexStruct_Same(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		BeEqual(b, complexStruct1, complexStruct2)
	}
}

func BenchmarkBeEqual_Slice_Small(b *testing.B) {
	slice1 := []string{"a", "b", "c"}
	slice2 := []string{"a", "b", "c"}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		BeEqual(b, slice1, slice2)
	}
}

func BenchmarkBeEqual_Slice_Large(b *testing.B) {
	slice1 := make([]int, 1000)
	slice2 := make([]int, 1000)
	for j := 0; j < 1000; j++ {
		slice1[j] = j
		slice2[j] = j
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		BeEqual(b, slice1, slice2)
	}
}

func BenchmarkBeEqual_Map_Small(b *testing.B) {
	map1 := map[string]int{"a": 1, "b": 2, "c": 3}
	map2 := map[string]int{"a": 1, "b": 2, "c": 3}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		BeEqual(b, map1, map2)
	}
}

func BenchmarkBeEqual_Map_Large(b *testing.B) {
	map1 := make(map[string]int, 100)
	map2 := make(map[string]int, 100)
	for j := 0; j < 100; j++ {
		key := string(rune('a' + j%26))
		map1[key] = j
		map2[key] = j
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		BeEqual(b, map1, map2)
	}
}

// ContainSubstring Benchmarks - Critical function #2
func BenchmarkContainSubstring_Short_Found_Beginning(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ContainSubstring(b, shortText, "Hello")
	}
}

func BenchmarkContainSubstring_Short_Found_Middle(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ContainSubstring(b, shortText, needleShort)
	}
}

func BenchmarkContainSubstring_Short_Found_End(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ContainSubstring(b, shortText, "benchmarking.")
	}
}

func BenchmarkContainSubstring_Short_NotFound(b *testing.B) {
	b.ReportAllocs()
	mt := &mockT{}
	for i := 0; i < b.N; i++ {
		ContainSubstring(mt, shortText, needleNotFound)
	}
}

func BenchmarkContainSubstring_Medium_Found(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ContainSubstring(b, mediumText, needleMedium)
	}
}

func BenchmarkContainSubstring_Medium_NotFound(b *testing.B) {
	b.ReportAllocs()
	mt := &mockT{}
	for i := 0; i < b.N; i++ {
		ContainSubstring(mt, mediumText, needleNotFound)
	}
}

func BenchmarkContainSubstring_Long_Found(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ContainSubstring(b, longText, needleLong)
	}
}

func BenchmarkContainSubstring_Long_NotFound(b *testing.B) {
	b.ReportAllocs()
	mt := &mockT{}
	for i := 0; i < b.N; i++ {
		ContainSubstring(mt, longText, needleNotFound)
	}
}

func BenchmarkContainSubstring_VeryLong_Found(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ContainSubstring(b, veryLongText, needleVeryLong)
	}
}

func BenchmarkContainSubstring_VeryLong_NotFound(b *testing.B) {
	b.ReportAllocs()
	mt := &mockT{}
	for i := 0; i < b.N; i++ {
		ContainSubstring(mt, veryLongText, needleNotFound)
	}
}

func BenchmarkContainSubstring_CaseDifference(b *testing.B) {
	b.ReportAllocs()
	mt := &mockT{}
	for i := 0; i < b.N; i++ {
		ContainSubstring(mt, shortText, needleCaseDiff)
	}
}

func BenchmarkContainSubstring_WithIgnoreCase(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ContainSubstring(b, shortText, needleCaseDiff, WithIgnoreCase())
	}
}
