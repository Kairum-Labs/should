#!/bin/bash

set -e

if [ $# -ne 3 ]; then
    echo "Usage: $0 <baseline_file> <current_file> <threshold_percent>"
    echo "Example: $0 benchmark.baseline benchmark.current 2"
    exit 1
fi

BASELINE_FILE="$1"
CURRENT_FILE="$2"
THRESHOLD="$3"

if [ ! -f "$BASELINE_FILE" ]; then
    echo "Error: Baseline file '$BASELINE_FILE' not found"
    exit 1
fi

if [ ! -f "$CURRENT_FILE" ]; then
    echo "Error: Current file '$CURRENT_FILE' not found"
    exit 1
fi

# Function to extract benchmark results
extract_benchmarks() {
    local file="$1"
    grep "^Benchmark" "$file" | grep "ns/op" | while read -r line; do
        name=$(echo "$line" | awk '{print $1}')
        ns_per_op=$(echo "$line" | awk '{for(i=1;i<=NF;i++) if($i ~ /^[0-9.]+$/ && $(i+1) == "ns/op") print $i}')
        if [ -n "$ns_per_op" ]; then
            echo "$name:$ns_per_op"
        fi
    done
}

# Extract benchmarks from both files
baseline_results=$(extract_benchmarks "$BASELINE_FILE")
current_results=$(extract_benchmarks "$CURRENT_FILE")

if [ -z "$baseline_results" ]; then
    echo "Warning: No benchmark results found in baseline file"
    exit 0
fi

if [ -z "$current_results" ]; then
    echo "Warning: No benchmark results found in current file"
    exit 0
fi

# Compare results
echo "Benchmark Comparison Results:"
echo "=============================="

# Use a temporary file to track regressions
regression_file=$(mktemp)
echo "0" > "$regression_file"

echo "$current_results" | while IFS=: read -r current_name current_ns; do
    baseline_ns=$(echo "$baseline_results" | grep "^$current_name:" | cut -d: -f2)
    
    if [ -n "$baseline_ns" ] && [ "$baseline_ns" != "0" ]; then
        increase=$(awk "BEGIN {printf \"%.2f\", (($current_ns - $baseline_ns) / $baseline_ns) * 100}")
        
        if awk "BEGIN {exit !($increase > $THRESHOLD)}"; then
            echo "REGRESSION: $current_name"
            echo "   Baseline: ${baseline_ns} ns/op"
            echo "   Current:  ${current_ns} ns/op"
            echo "   Increase: ${increase}% (threshold: ${THRESHOLD}%)"
            echo ""
            echo "1" > "$regression_file"
        else
            echo "OK: $current_name (${increase}% change)"
        fi
    else
        echo "NEW: $current_name (${current_ns} ns/op)"
    fi
done

# Check final result and exit accordingly  
regressions_found=$(cat "$regression_file")
rm -f "$regression_file"

if [ "$regressions_found" -eq 1 ]; then
    echo ""
    echo "Performance regressions detected! Build should fail."
    exit 1
else
    echo ""
    echo "No performance regressions detected."
    exit 0
fi