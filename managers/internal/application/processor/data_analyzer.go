package processor

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
	"middleman/managers/internal/application/services"
)

// DataAnalyzer performs actual data analysis
type DataAnalyzer struct{}

// NewDataAnalyzer creates a new data analyzer
func NewDataAnalyzer() *DataAnalyzer {
	return &DataAnalyzer{}
}

// AnalyzeCSV analyzes CSV data
func (a *DataAnalyzer) AnalyzeCSV(data []byte) (*services.DataProcessingResult, error) {
	reader := csv.NewReader(bytes.NewReader(data))

	// Read header
	headers, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV headers: %w", err)
	}

	// Initialize column analysis
	columns := make([]ColumnAnalysis, len(headers))
	for i, header := range headers {
		columns[i] = ColumnAnalysis{
			Name:   header,
			Values: make(map[string]int),
			Types:  make(map[string]int),
		}
	}

	// Read and analyze rows
	rowCount := 0
	sampleData := []map[string]interface{}{}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading CSV row %d: %w", rowCount+1, err)
		}

		rowCount++

		// Analyze each column value
		rowData := make(map[string]interface{})
		for i, value := range record {
			if i < len(columns) {
				columns[i].AnalyzeValue(value)
				rowData[headers[i]] = value
			}
		}

		// Collect sample data (first 10 rows)
		if rowCount <= 10 {
			sampleData = append(sampleData, rowData)
		}
	}

	// Build result
	result := &services.DataProcessingResult{
		Summary:     fmt.Sprintf("CSV file with %d columns and %d rows", len(headers), rowCount),
		RowCount:    rowCount,
		ColumnCount: len(headers),
		SampleData:  sampleData,
		Metadata: map[string]interface{}{
			"columns":     a.buildColumnInfo(columns),
			"file_format": "csv",
			"headers":     headers,
		},
	}

	return result, nil
}

// AnalyzeJSON analyzes JSON data
func (a *DataAnalyzer) AnalyzeJSON(data []byte) (*services.DataProcessingResult, error) {
	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Analyze structure
	analysis := a.analyzeJSONStructure(jsonData)

	// Build result
	result := &services.DataProcessingResult{
		Summary:     analysis.Summary,
		RowCount:    analysis.RecordCount,
		ColumnCount: analysis.FieldCount,
		SampleData:  analysis.SampleData,
		Metadata: map[string]interface{}{
			"structure":   analysis.Structure,
			"file_format": "json",
			"data_type":   analysis.DataType,
		},
	}

	return result, nil
}

// AnalyzeXML analyzes XML data
func (a *DataAnalyzer) AnalyzeXML(data []byte) (*services.DataProcessingResult, error) {
	// For production, we would use encoding/xml to properly parse
	// For now, return basic analysis
	lines := strings.Split(string(data), "\n")

	result := &services.DataProcessingResult{
		Summary:     fmt.Sprintf("XML document with %d lines", len(lines)),
		RowCount:    -1, // Not applicable for XML
		ColumnCount: -1, // Not applicable for XML
		SampleData:  []map[string]interface{}{},
		Metadata: map[string]interface{}{
			"file_format": "xml",
			"size_bytes":  len(data),
		},
	}

	return result, nil
}

// AnalyzeYAML analyzes YAML data
func (a *DataAnalyzer) AnalyzeYAML(data []byte) (*services.DataProcessingResult, error) {
	var yamlData interface{}
	if err := yaml.Unmarshal(data, &yamlData); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Analyze structure (similar to JSON)
	analysis := a.analyzeJSONStructure(yamlData)

	result := &services.DataProcessingResult{
		Summary:     analysis.Summary,
		RowCount:    analysis.RecordCount,
		ColumnCount: analysis.FieldCount,
		SampleData:  analysis.SampleData,
		Metadata: map[string]interface{}{
			"structure":   analysis.Structure,
			"file_format": "yaml",
			"data_type":   analysis.DataType,
		},
	}

	return result, nil
}

// ValidateDataQuality performs quality checks on data
func (a *DataAnalyzer) ValidateDataQuality(data []byte, format string) (*services.DataQualityResult, error) {
	switch format {
	case "csv":
		return a.validateCSVQuality(data)
	case "json":
		return a.validateJSONQuality(data)
	default:
		return &services.DataQualityResult{
			QualityScore: 0.5,
			Issues: []services.DataQualityIssue{
				{
					Type:        "unsupported_format",
					Severity:    "warning",
					Count:       1,
					Description: fmt.Sprintf("Quality validation not fully implemented for format: %s", format),
				},
			},
			Metadata: map[string]interface{}{
				"format": format,
			},
		}, nil
	}
}

// Helper structures and methods

type ColumnAnalysis struct {
	Name         string
	Values       map[string]int
	Types        map[string]int
	NullCount    int
	TotalCount   int
	NumericCount int
	Min          float64
	Max          float64
	Sum          float64
}

func (c *ColumnAnalysis) AnalyzeValue(value string) {
	c.TotalCount++

	// Check for null/empty
	if value == "" || value == "null" || value == "NULL" {
		c.NullCount++
		return
	}

	// Track unique values (up to 1000)
	if len(c.Values) < 1000 {
		c.Values[value]++
	}

	// Detect type
	if num, err := strconv.ParseFloat(value, 64); err == nil {
		c.Types["numeric"]++
		c.NumericCount++
		if c.NumericCount == 1 {
			c.Min = num
			c.Max = num
		} else {
			c.Min = math.Min(c.Min, num)
			c.Max = math.Max(c.Max, num)
		}
		c.Sum += num
	} else if _, err := time.Parse(time.RFC3339, value); err == nil {
		c.Types["datetime"]++
	} else if value == "true" || value == "false" {
		c.Types["boolean"]++
	} else {
		c.Types["string"]++
	}
}

func (a *DataAnalyzer) buildColumnInfo(columns []ColumnAnalysis) []services.ColumnInfo {
	result := make([]services.ColumnInfo, len(columns))

	for i, col := range columns {
		// Determine primary data type
		dataType := "string"
		maxTypeCount := 0
		for typeName, count := range col.Types {
			if count > maxTypeCount {
				maxTypeCount = count
				dataType = typeName
			}
		}

		// Get sample values
		samples := []interface{}{}
		for value := range col.Values {
			samples = append(samples, value)
			if len(samples) >= 5 {
				break
			}
		}

		// Build statistics
		stats := map[string]interface{}{
			"null_rate":    float64(col.NullCount) / float64(col.TotalCount),
			"unique_count": len(col.Values),
		}

		if col.NumericCount > 0 {
			stats["min"] = col.Min
			stats["max"] = col.Max
			stats["mean"] = col.Sum / float64(col.NumericCount)
		}

		result[i] = services.ColumnInfo{
			Name:         col.Name,
			DataType:     dataType,
			SampleValues: samples,
			NullCount:    col.NullCount,
			UniqueCount:  len(col.Values),
			Statistics:   stats,
		}
	}

	return result
}

type JSONAnalysis struct {
	Summary     string
	RecordCount int
	FieldCount  int
	DataType    string
	Structure   map[string]interface{}
	SampleData  []map[string]interface{}
}

func (a *DataAnalyzer) analyzeJSONStructure(data interface{}) *JSONAnalysis {
	analysis := &JSONAnalysis{
		Structure:  make(map[string]interface{}),
		SampleData: []map[string]interface{}{},
	}

	switch v := data.(type) {
	case []interface{}:
		analysis.DataType = "array"
		analysis.RecordCount = len(v)
		analysis.Summary = fmt.Sprintf("JSON array with %d records", len(v))

		// Analyze first few records
		for i, item := range v {
			if i >= 10 {
				break
			}
			if m, ok := item.(map[string]interface{}); ok {
				analysis.SampleData = append(analysis.SampleData, m)
				if i == 0 {
					analysis.FieldCount = len(m)
					for key := range m {
						analysis.Structure[key] = detectJSONType(m[key])
					}
				}
			}
		}

	case map[string]interface{}:
		analysis.DataType = "object"
		analysis.RecordCount = 1
		analysis.FieldCount = len(v)
		analysis.Summary = fmt.Sprintf("JSON object with %d fields", len(v))
		analysis.SampleData = append(analysis.SampleData, v)

		for key, val := range v {
			analysis.Structure[key] = detectJSONType(val)
		}

	default:
		analysis.DataType = "primitive"
		analysis.Summary = "JSON primitive value"
	}

	return analysis
}

func detectJSONType(v interface{}) string {
	switch v.(type) {
	case string:
		return "string"
	case float64:
		return "number"
	case bool:
		return "boolean"
	case nil:
		return "null"
	case []interface{}:
		return "array"
	case map[string]interface{}:
		return "object"
	default:
		return "unknown"
	}
}

func (a *DataAnalyzer) validateCSVQuality(data []byte) (*services.DataQualityResult, error) {
	reader := csv.NewReader(bytes.NewReader(data))

	issues := []services.DataQualityIssue{}
	missingValues := make(map[string]int)
	rowLengths := make(map[int]int)
	duplicates := 0
	rowHashes := make(map[string]int)

	headers, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read headers: %w", err)
	}

	expectedCols := len(headers)
	rowNum := 0

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			issues = append(issues, services.DataQualityIssue{
				Type:        "parse_error",
				Severity:    "error",
				Description: fmt.Sprintf("Error parsing row %d: %v", rowNum+1, err),
				Count:       1,
			})
			continue
		}

		rowNum++

		// Check column count
		if len(record) != expectedCols {
			rowLengths[len(record)]++
		}

		// Check for empty values
		for i, value := range record {
			if value == "" && i < len(headers) {
				missingValues[headers[i]]++
			}
		}

		// Check for duplicate rows
		rowHash := strings.Join(record, "|")
		if rowHashes[rowHash] > 0 {
			duplicates++
		}
		rowHashes[rowHash]++
	}

	// Calculate quality score
	qualityScore := 1.0

	// Deduct for inconsistent row lengths
	if len(rowLengths) > 1 {
		issues = append(issues, services.DataQualityIssue{
			Type:        "inconsistent_columns",
			Severity:    "error",
			Description: "Rows have inconsistent number of columns",
			Count:       len(rowLengths) - 1,
		})
		qualityScore -= 0.2
	}

	// Deduct for missing values
	totalMissing := 0
	for col, count := range missingValues {
		totalMissing += count
		missRate := float64(count) / float64(rowNum)
		if missRate > 0.1 {
			issues = append(issues, services.DataQualityIssue{
				Type:        "missing_values",
				Severity:    "warning",
				Column:      col,
				Count:       count,
				Description: fmt.Sprintf("Column '%s' has %.1f%% missing values", col, missRate*100),
			})
		}
	}
	if totalMissing > 0 {
		qualityScore -= math.Min(0.3, float64(totalMissing)/float64(rowNum*len(headers)))
	}

	// Deduct for duplicates
	if duplicates > 0 {
		issues = append(issues, services.DataQualityIssue{
			Type:        "duplicate_rows",
			Severity:    "warning",
			Count:       duplicates,
			Description: fmt.Sprintf("Found %d duplicate rows", duplicates),
		})
		qualityScore -= math.Min(0.2, float64(duplicates)/float64(rowNum))
	}

	// Build recommendations
	recommendations := []string{}
	if len(missingValues) > 0 {
		recommendations = append(recommendations, "Consider handling missing values through imputation or removal")
	}
	if duplicates > 0 {
		recommendations = append(recommendations, "Remove duplicate rows to ensure data integrity")
	}
	if len(rowLengths) > 1 {
		recommendations = append(recommendations, "Fix rows with inconsistent column counts")
	}

	// Build data types map
	dataTypes := make(map[string]string)
	for _, header := range headers {
		dataTypes[header] = "string" // Would be enhanced with actual type detection
	}

	return &services.DataQualityResult{
		QualityScore:    math.Max(0, qualityScore),
		Issues:          issues,
		MissingValues:   missingValues,
		DuplicateRows:   duplicates,
		DataTypes:       dataTypes,
		ValidityChecks:  []services.ValidityCheck{},
		Recommendations: recommendations,
		Metadata: map[string]interface{}{
			"row_count":    rowNum,
			"column_count": len(headers),
		},
	}, nil
}

func (a *DataAnalyzer) validateJSONQuality(data []byte) (*services.DataQualityResult, error) {
	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return &services.DataQualityResult{
			QualityScore: 0,
			Issues: []services.DataQualityIssue{
				{
					Type:        "parse_error",
					Severity:    "error",
					Count:       1,
					Description: fmt.Sprintf("Invalid JSON: %v", err),
				},
			},
		}, nil
	}

	issues := []services.DataQualityIssue{}
	missingValues := make(map[string]int)
	qualityScore := 1.0

	// Analyze based on JSON type
	switch v := jsonData.(type) {
	case []interface{}:
		// Validate array consistency
		if len(v) == 0 {
			issues = append(issues, services.DataQualityIssue{
				Type:        "empty_data",
				Severity:    "warning",
				Count:       1,
				Description: "JSON array is empty",
			})
			qualityScore -= 0.5
		} else {
			// Check consistency of array elements
			fieldCounts := make(map[int]int)
			for i, item := range v {
				if m, ok := item.(map[string]interface{}); ok {
					fieldCounts[len(m)]++
					// Check for null/missing values
					for key, val := range m {
						if val == nil {
							missingValues[key]++
						}
					}
				} else {
					issues = append(issues, services.DataQualityIssue{
						Type:        "inconsistent_structure",
						Severity:    "warning",
						Count:       1,
						Description: fmt.Sprintf("Array element %d is not an object", i),
					})
				}
			}

			if len(fieldCounts) > 1 {
				issues = append(issues, services.DataQualityIssue{
					Type:        "inconsistent_schema",
					Severity:    "warning",
					Count:       len(fieldCounts),
					Description: "Objects in array have varying number of fields",
				})
				qualityScore -= 0.2
			}
		}

	case map[string]interface{}:
		// Check for null values in object
		for key, val := range v {
			if val == nil {
				missingValues[key]++
			}
		}
	}

	// Adjust score for missing values
	if len(missingValues) > 0 {
		qualityScore -= 0.1
	}

	recommendations := []string{}
	if len(issues) > 0 {
		recommendations = append(recommendations, "Fix structural inconsistencies in the JSON data")
	}
	if len(missingValues) > 0 {
		recommendations = append(recommendations, "Handle null values appropriately")
	}

	return &services.DataQualityResult{
		QualityScore:    math.Max(0, qualityScore),
		Issues:          issues,
		MissingValues:   missingValues,
		DuplicateRows:   0, // Not easily determined for JSON
		DataTypes:       map[string]string{},
		ValidityChecks:  []services.ValidityCheck{},
		Recommendations: recommendations,
		Metadata: map[string]interface{}{
			"structure_type": detectJSONType(jsonData),
		},
	}, nil
}
