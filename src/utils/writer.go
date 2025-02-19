package utils

import (
	"errors"
	"fmt"

	"github.com/xuri/excelize/v2"
)

// Field constants to use when specifying which fields to update
const (
	FieldAuthToken = "AUTH_TOKEN"
	FieldProxy     = "PROXY"
	FieldUsername  = "USERNAME"
	FieldStatus    = "STATUS"
)

const AccountsFilePath = "data/accounts.xlsx"

// UpdateAccountInExcel updates specified fields for an account in the Excel file.
//
// Usage examples:
//
//	// Update single field
//	account := utils.Account{
//	    AuthToken: "token123",  // Required to identify the row
//	}
//	err := utils.UpdateAccountInExcel("./data/accounts.xlsx", account, map[string]string{
//	    utils.FieldStatus: "Active",
//	})
//
//	// Update multiple fields
//	updates := map[string]string{
//	    utils.FieldStatus: "Active",
//	    utils.FieldProxy: "new_proxy",
//	    utils.FieldUsername: "new_username",
//	}
//	err = utils.UpdateAccountInExcel("./data/accounts.xlsx", account, updates)
//
// The function is thread-safe and will only update the specified fields.
// AUTH_TOKEN is used to find the correct row to update.
func UpdateAccountInExcel(filePath string, account Account, fieldsToUpdate map[string]string) error {
	// Lock the mutex before accessing the file
	ExcelMutex.Lock()
	defer ExcelMutex.Unlock()

	// Open the Excel file
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Get the first sheet
	sheetName := f.GetSheetName(0)

	// Get all rows from the sheet
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return err
	}

	// Find the row with matching AUTH_TOKEN
	rowIndex := -1
	for i := 1; i < len(rows); i++ {
		if len(rows[i]) > 0 && rows[i][0] == account.AuthToken {
			rowIndex = i + 1 // Excel rows are 1-based
			break
		}
	}

	if rowIndex == -1 {
		return errors.New("account not found in Excel file")
	}

	// Update specified fields
	for field, newValue := range fieldsToUpdate {
		var colIndex string
		switch field {
		case FieldAuthToken:
			colIndex = "A"
		case FieldProxy:
			colIndex = "B"
		case FieldUsername:
			colIndex = "C"
		case FieldStatus:
			colIndex = "D"
		default:
			return errors.New("invalid field name: " + field)
		}

		// Update the cell
		cellRef := colIndex + string(rune('0'+rowIndex))
		if err := f.SetCellValue(sheetName, cellRef, newValue); err != nil {
			return err
		}
	}

	// Save the file
	if err := f.Save(); err != nil {
		return err
	}

	return nil
}

// ReportFailedAccount updates the status of an account in the Excel file
func UpdateAccountStatus(filePath string, account Account, status string) error {
	// Lock the mutex before accessing the file
	ExcelMutex.Lock()
	defer ExcelMutex.Unlock()

	// Open the Excel file
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to open Excel file: %w", err)
	}
	defer f.Close()

	// Get the first sheet
	sheetName := f.GetSheetName(0)

	// Get all rows from the sheet
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("failed to get rows: %w", err)
	}

	// Find the row with matching auth token
	rowIndex := -1
	for i, row := range rows {
		if len(row) > 0 && row[0] == account.AuthToken {
			rowIndex = i + 1 // Excel rows are 1-based
			break
		}
	}

	if rowIndex == -1 {
		return fmt.Errorf("account not found in Excel file")
	}

	// Create cell style based on status
	var style int
	switch status {
	case "VALID":
		style, err = f.NewStyle(&excelize.Style{
			Fill: excelize.Fill{Type: "pattern", Color: []string{"#90EE90"}, Pattern: 1}, // Light green
		})
	case "UNKNOWN", "AUTH_ERROR":
		style, err = f.NewStyle(&excelize.Style{
			Fill: excelize.Fill{Type: "pattern", Color: []string{"#FFB6B6"}, Pattern: 1}, // Light red
		})
	case "LOCKED":
		style, err = f.NewStyle(&excelize.Style{
			Fill: excelize.Fill{Type: "pattern", Color: []string{"#ADD8E6"}, Pattern: 1}, // Light blue
		})
	case "SUSPENDED":
		style, err = f.NewStyle(&excelize.Style{
			Fill: excelize.Fill{Type: "pattern", Color: []string{"#D3D3D3"}, Pattern: 1}, // Light gray
		})
	}
	if err != nil {
		return fmt.Errorf("failed to create cell style: %w", err)
	}

	// Update the status column (column D) with value and style
	cellRef := fmt.Sprintf("D%d", rowIndex)
	err = f.SetCellValue(sheetName, cellRef, status)
	if err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}
	err = f.SetCellStyle(sheetName, cellRef, cellRef, style)
	if err != nil {
		return fmt.Errorf("failed to set cell style: %w", err)
	}

	// Save the changes
	err = f.Save()
	if err != nil {
		return fmt.Errorf("failed to save Excel file: %w", err)
	}

	Logger{}.Success("%s | Updated account status to: %s", account.AuthToken, status)
	return nil
}
