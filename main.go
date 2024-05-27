package main

import (
    "encoding/csv"
    "flag"
    "fmt"
    "log"
    "os"
    "strconv"
    "strings"

    "github.com/xuri/excelize/v2"
)

func main() {
    inputFile := flag.String("input", "", "Path to the input Excel file")
    outputFile := flag.String("output", "", "Path to the output CSV file")
    creditCardMode := flag.Bool("cc", false, "Enable credit card mode (amounts will be negative)")

    flag.Parse()

    if *inputFile == "" || *outputFile == "" {
        log.Fatalf("Usage: %s -input <input.xlsx> -output <output.csv> [-cc]\n", os.Args[0])
    }

    // Open the Excel file
    f, err := excelize.OpenFile(*inputFile)
    if err != nil {
        log.Fatalf("Failed to open Excel file: %s\n", err)
    }

    // Determine the format of the Excel file and process accordingly
    sheetName := f.GetSheetName(0)
    rows, err := f.GetRows(sheetName)
    if err != nil {
        log.Fatalf("Failed to read rows: %s\n", err)
    }

    // Open the CSV file for writing
    csvFile, err := os.Create(*outputFile)
    if err != nil {
        log.Fatalf("Failed to create CSV file: %s\n", err)
    }
    defer csvFile.Close()

    csvWriter := csv.NewWriter(csvFile)
    defer csvWriter.Flush()

    // Process the rows
    processRows(rows, csvWriter, *creditCardMode)

    fmt.Println("CSV file created successfully:", *outputFile)
}

func processRows(rows [][]string, csvWriter *csv.Writer, creditCardMode bool) {
    for i, row := range rows {
        if i == 0 {
            continue // skip the header row
        }
        if len(row) < 4 || isDividerRow(row) {
            continue // skip incomplete or divider rows
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

        // Remove commas and parse the amount
        amount = strings.ReplaceAll(amount, ",", "")
        amountFloat, err := strconv.ParseFloat(amount, 64)
        if err == nil {
            if creditCardMode {
                amountFloat = -amountFloat
            }
            amount = fmt.Sprintf("%.2f", amountFloat) // Ensure two decimal places
        }

        // Split the concept into PAYEE and DESC
        payee := ""
        desc := ""
        parts := strings.SplitN(concept, "/", 2)
        if len(parts) > 0 {
            payee = strings.TrimSpace(parts[0])
        }
        if len(parts) > 1 {
            desc = strings.TrimSpace(parts[1])
        }

        // Write to CSV
        if err := csvWriter.Write([]string{date, payee, desc, amount}); err != nil {
            log.Fatalf("Failed to write to CSV: %s\n", err)
        }
    }
}

func isDividerRow(row []string) bool {
    // A divider row typically has a unique pattern, such as a single non-empty cell
    if len(row) == 0 {
        return false
    }
    nonEmptyCells := 0
    for _, cell := range row {
        if strings.TrimSpace(cell) != "" {
            nonEmptyCells++
        }
    }
    return nonEmptyCells == 1
}
