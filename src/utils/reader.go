package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"encoding/base64"

	"github.com/xuri/excelize/v2"
)

// ReadAccountsFromExcel reads account data from an Excel file and returns a slice of Account structs.
// It uses the provided start and end indices to determine which accounts to process.
//
// Usage example:
//
//	accounts, err := utils.ReadAccountsFromExcel("./data/accounts.xlsx", 1, 10)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, account := range accounts {
//	    fmt.Printf("Account: %+v\n", account)
//	}
//
// The function is thread-safe and will read all accounts until it finds an empty AUTH_TOKEN.
// It skips the header row and handles missing columns by setting empty strings.
func ReadAccountsFromExcel(filePath string, startIndex, endIndex int) ([]Account, error) {
	// Lock the mutex before accessing the file
	ExcelMutex.Lock()
	defer ExcelMutex.Unlock()

	// Open the Excel file
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Get the first sheet
	sheetName := f.GetSheetName(0)

	// Get all rows from the sheet
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	var accounts []Account

	// Skip the first row (header) and process remaining rows
	for i := 1; i < len(rows); i++ {
		row := rows[i]

		// Check if we have enough columns
		if len(row) < 1 {
			continue
		}

		// Stop if AUTH_TOKEN is empty
		if row[0] == "" {
			break
		}

		// Create account struct, handling possible missing columns
		account := Account{
			AuthToken: row[0],
			Proxy:     "",
			Username:  "",
			Status:    "",
		}

		// Safely assign other fields if they exist
		if len(row) > 1 {
			account.Proxy = row[1]
		}
		if len(row) > 2 {
			account.Username = row[2]
		}
		if len(row) > 3 {
			account.Status = row[3]
		}

		accounts = append(accounts, account)
	}

	// If both are 1 and 0 or equal, process all accounts
	if (startIndex == 1 && endIndex == 0) || startIndex == endIndex {
		return accounts, nil
	}

	// Validate indices
	if startIndex < 1 {
		startIndex = 1
	}
	if endIndex == 0 || endIndex > len(accounts) {
		endIndex = len(accounts)
	}

	// Return the slice of accounts within the range
	return accounts[startIndex-1 : endIndex], nil
}

func UserInputToSlice(textToShow string) []string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("\n" + textToShow + ": ")
	userInput, _ := reader.ReadString('\n')
	userInput = strings.TrimSpace(userInput)
	userInput = strings.ReplaceAll(userInput, ",", " ")

	return strings.Fields(userInput)
}

func UserInputInteger(textToShow string) (int, bool) {
	var input int

	fmt.Print(textToShow + ": ")
	reader := bufio.NewReader(os.Stdin)
	userInput, _ := reader.ReadString('\n')
	userInput = strings.TrimSpace(userInput)

	input, err := strconv.Atoi(userInput)
	if err != nil {
		Logger{}.Error("Wrong input.")
		return input, false
	}
	return input, true
}

func EncodeFileToBase64(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

func ReadTxtFile(fileName string, filePath string) []string {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
		return []string{}
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	var items []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		items = append(items, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Failed to read file: %v", err)
		return []string{}
	}

	Logger{}.Info("Successfully loaded %d %s.", len(items), fileName)
	return items
}

func ReadPictures(path string) ([]string, bool) {
	var encodedImages []string

	files, err := os.ReadDir(path)
	if err != nil {
		Logger{}.Error("Failed to read pictures dir: %s", err)
		return encodedImages, false
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".png") || strings.HasSuffix(file.Name(), ".jpg") || strings.HasSuffix(file.Name(), ".jpeg") {
			filePath := filepath.Join(path, file.Name())

			imageFile, err := os.ReadFile(filePath)
			if err != nil {
				Logger{}.Error("Failed to read the picture: %s", err)
			}

			encodedImage := base64.StdEncoding.EncodeToString(imageFile)
			encodedImages = append(encodedImages, encodedImage)
		}
	}

	return encodedImages, true
}
