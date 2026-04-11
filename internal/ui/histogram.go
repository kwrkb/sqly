package ui

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

const maxHistogramBuckets = 20
const maxHistogramRows = 10_000

// detectNumericColumn returns true if the column type string suggests a numeric type.
func detectNumericColumn(colType string) bool {
	lower := strings.ToLower(colType)
	for _, kw := range []string{"int", "real", "float", "double", "decimal", "numeric", "number"} {
		if strings.Contains(lower, kw) {
			return true
		}
	}
	return false
}

// looksLikeNumeric checks the first few non-null values to see if they parse as float64.
func looksLikeNumeric(rows [][]string, colIdx int) bool {
	checked := 0
	for _, row := range rows {
		if colIdx >= len(row) {
			continue
		}
		val := row[colIdx]
		if val == "NULL" {
			continue
		}
		if _, err := strconv.ParseFloat(val, 64); err != nil {
			return false
		}
		checked++
		if checked >= 3 {
			break
		}
	}
	return checked > 0
}

// computeHistogram scans a numeric column, buckets values into equal-width bins,
// and returns a histogramData with rendered bars.
// Returns zero-value histogramData if the column has fewer than 2 parseable values
// or if row count exceeds maxHistogramRows.
func computeHistogram(rows [][]string, colIdx int) histogramData {
	if len(rows) > maxHistogramRows {
		return histogramData{Skipped: true}
	}

	values := make([]float64, 0, len(rows))
	nonNull := 0
	for _, row := range rows {
		if colIdx >= len(row) {
			continue
		}
		val := row[colIdx]
		if val == "NULL" {
			continue
		}
		nonNull++
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			continue
		}
		if math.IsNaN(f) || math.IsInf(f, 0) {
			continue
		}
		values = append(values, f)
	}

	// Suppress histogram if less than half of non-NULL values are numeric,
	// to avoid misleading distributions on mixed-type columns.
	if nonNull > 0 && len(values)*2 < nonNull {
		return histogramData{}
	}

	if len(values) < 2 {
		return histogramData{}
	}

	minVal, maxVal := values[0], values[0]
	for _, v := range values[1:] {
		if v < minVal {
			minVal = v
		}
		if v > maxVal {
			maxVal = v
		}
	}

	if minVal == maxVal {
		return histogramData{}
	}

	numBuckets := max(min(len(values), maxHistogramBuckets), 2)

	width := (maxVal - minVal) / float64(numBuckets)
	counts := make([]int, numBuckets)
	for _, v := range values {
		idx := int((v - minVal) / width)
		if idx >= numBuckets {
			idx = numBuckets - 1
		}
		counts[idx]++
	}

	bars := renderSparklineBars(counts)
	if bars == "" {
		return histogramData{}
	}

	label := formatHistogramLabel(minVal, maxVal)

	return histogramData{
		Bars:  bars,
		Label: label,
	}
}

// formatHistogramLabel returns a compact range label for the histogram.
// Uses integer format if both values are whole numbers, otherwise %.4g.
func formatHistogramLabel(minVal, maxVal float64) string {
	if minVal == math.Trunc(minVal) && maxVal == math.Trunc(maxVal) {
		return fmt.Sprintf("%.0f–%.0f", minVal, maxVal)
	}
	return fmt.Sprintf("%.4g–%.4g", minVal, maxVal)
}
