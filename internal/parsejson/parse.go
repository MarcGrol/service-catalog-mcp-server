package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/ghetzel/go-stockutil/maputil"
)

//go:embed crash_logs.json
var crashesJson []byte

type crashlogReport map[string]string

func main() {
	reports, err := parseCrashlogs()
	if err != nil {
		log.Fatalf("Failed to parse crashlogs")
	}

	columnNames := collectColumnNames(reports)

	csvBlob := printCSV(reports, columnNames)

	fmt.Fprint(os.Stdout, csvBlob)
}

func parseCrashlogs() ([]crashlogReport, error) {
	reports := []crashlogReport{}
	err := json.Unmarshal(crashesJson, &reports)
	if err != nil {
		return reports, fmt.Errorf("error unmarshalling crashes json: %v", err)
	}

	for idx, report := range reports {
		report, err = parseSubLog(report)
		if err != nil {
			return reports, fmt.Errorf("error parsing sub-log for line %d: %s", idx+1, err)
		}
		reports[idx] = report
	}
	return reports, nil
}

func parseSubLog(report crashlogReport) (crashlogReport, error) {
	data := map[string]any{}
	err := json.Unmarshal([]byte(report["data"]), &data)
	if err != nil {
		return report, fmt.Errorf("Error unmarshalling sub-json: %v", err)
	}

	flattened, err := maputil.CoalesceMap(data, ".")
	if err != nil {
		return report, fmt.Errorf("Error flattening sub json: %v", err)
	}
	for key, value := range flattened {
		if !strings.HasPrefix(key, "panic.panicInfo.panicMessage.") {
			stringValue := fmt.Sprintf("%v", value)
			report[key] = strings.ReplaceAll(stringValue, "\n", " ")
		}
	}
	delete(report, "data")

	return report, nil
}

func collectColumnNames(reports []crashlogReport) []string {
	keyMap := map[string]bool{}
	for _, report := range reports {
		for field, _ := range report {
			keyMap[field] = true
		}
	}

	keySlice := []string{}
	for field := range keyMap {
		keySlice = append(keySlice, field)
	}
	sort.Strings(keySlice)

	return keySlice
}

func printCSV(reports []crashlogReport, keys []string) string {
	buf := &bytes.Buffer{}
	printCSVHeader(keys, buf)

	for _, report := range reports {
		printCSVLine(keys, report, buf)
	}
	return buf.String()
}

func printCSVLine(keys []string, report crashlogReport, writer io.Writer) {
	for _, key := range keys {
		fmt.Fprintf(writer, "\"%s\",", report[key])
	}
	fmt.Fprintf(writer, "\n")
}

func printCSVHeader(keys []string, writer io.Writer) {
	for _, key := range keys {
		fmt.Fprintf(writer, "%s,", key)
	}
	fmt.Fprintf(writer, "\n")
}
