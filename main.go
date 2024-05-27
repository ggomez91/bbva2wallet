package main

import (
    "encoding/csv"
    "fmt"
    "log"
    "os"
    "github.com/xuri/excelize/v2"
)

func main() {
    if len(os.Args) < 3 {
        log.Fatalf("Usage: %s <input.xlsx> <output.csv>\n", os.Args[0])
    }

    inputFile := os.Args[1]
    outputFile := os.Args[2]

    // Open the Excel file
    f, err := excelize.OpenFile(inputFile)
    if err != nil {
        log.Fatalf("Failed to open Excel file: %s\n", err)
    }

    // Read the rows from the first sheet
    rows, err := f.GetRows(f.GetSheetName(0))
    if err != nil {
        log.Fatalf("Failed to read rows: %s\n", err)
    }

    // Open the CSV file for writing
    csvFile, err := os.Create(outputFile)
    if err != nil {
        log.Fatalf("Failed to create CSV file: %s\n", err)
    }
    defer csvFile.Close()

    csvWriter := csv.NewWriter(csvFile)
    defer csvWriter.Flush()

    // Process and write the rows
    for i, row := range rows {
		if i < 3 {
            continue // skip the first 3 rows
        }	
        if len(row) < 4 {
            continue // skip incomplete rows
        }
        date := row[0]
        concept := row[1]
        cargo := row[2]
        abono := row[3]

        // Merge CARGO and ABONO into one column
        amount := cargo
        if cargo == "" {
            amount = abono
        }

        // Write to CSV
        if err := csvWriter.Write([]string{date, concept, amount}); err != nil {
            log.Fatalf("Failed to write to CSV: %s\n", err)
        }
    }

    fmt.Println("CSV file created successfully:", outputFile)
}
