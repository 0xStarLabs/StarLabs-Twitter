package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

const (
	SOFTWARE_VERSION = "2.0.0"
)

// ShowLogo displays the STARLABS ASCII art logo
func ShowLogo() {
	// Create ASCII art with a different font
	myFigure := figure.NewFigure("STAR LABS", "standard", true)
	logo := myFigure.String()

	// Create color objects for blue gradient effect
	blue := color.New(color.FgHiBlue, color.Bold)

	// Split logo into lines
	lines := strings.Split(logo, "\n")

	fmt.Println() // Empty line before logo

	// Print each line in bold blue
	for _, line := range lines {
		if len(strings.TrimSpace(line)) > 0 {
			blue.Println(line)
		}
	}
	fmt.Println() // Empty line after logo
}

// ShowDevInfo displays development and version information
func ShowDevInfo() {
	// Create blue color
	blue := color.New(color.FgHiBlue)

	// Box drawing characters
	const (
		topLeft     = "╔"
		topRight    = "╗"
		bottomLeft  = "╚"
		bottomRight = "╝"
		horizontal  = "═"
		vertical    = "║"
		separator   = "----------------------------------------"
	)

	// Create the info box content
	width := 42
	horizontalLine := strings.Repeat(horizontal, width-2)

	info := []string{
		topLeft + horizontalLine + topRight,
		vertical + centerText("Twitter Bot "+SOFTWARE_VERSION, width-2) + vertical,
		vertical + separator + vertical,
		vertical + centerText("", width-2) + vertical,
		vertical + centerText("GitHub: https://github.com/StarLabs", width-2) + vertical,
		vertical + centerText("", width-2) + vertical,
		vertical + centerText("Developer: https://t.me/StarLabsTech", width-2) + vertical,
		vertical + centerText("Chat: https://t.me/StarLabsChat", width-2) + vertical,
		vertical + centerText("", width-2) + vertical,
		bottomLeft + horizontalLine + bottomRight,
	}

	fmt.Println() // Empty line before info box

	// Print all lines in the same blue color
	for _, line := range info {
		blue.Println(line)
	}
	fmt.Println() // Empty line after info box
}

// centerText centers the text within the given width
func centerText(text string, width int) string {
	padding := width - len(text)
	if padding <= 0 {
		return text
	}
	leftPad := padding / 2
	rightPad := padding - leftPad
	return strings.Repeat(" ", leftPad) + text + strings.Repeat(" ", rightPad)
}

// PrintStatistics prints a beautiful table with processing statistics
func PrintStatistics(total, processed, success, failed, locked, suspended int) {
	// Create new table
	table := tablewriter.NewWriter(os.Stdout)

	// Configure table style
	table.SetHeader([]string{"Statistics", "Value"})
	table.SetBorder(true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderAlignment(tablewriter.ALIGN_CENTER)

	// Make table bigger
	table.SetColumnSeparator("│")
	table.SetRowSeparator("─")
	table.SetCenterSeparator("┼")
	table.SetHeaderLine(true)
	table.SetColWidth(30)
	table.SetAutoWrapText(false)

	// Set header colors (cyan)
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.FgHiCyanColor, tablewriter.Bold},
		tablewriter.Colors{tablewriter.FgHiCyanColor, tablewriter.Bold},
	)

	// Set different colors for each row
	table.SetColumnColor(
		tablewriter.Colors{tablewriter.FgHiWhiteColor},
		tablewriter.Colors{}, // Will be set per row
	)

	// Add data rows with different colors for values
	data := []struct {
		label string
		value int
		color int
	}{
		{"Total", total, tablewriter.FgHiWhiteColor},
		{"Processed", processed, tablewriter.FgHiWhiteColor},
		{"Success", success, tablewriter.FgHiGreenColor},
		{"Failed", failed, tablewriter.FgHiYellowColor},
		{"Locked", locked, tablewriter.FgHiRedColor},
		{"Suspended", suspended, tablewriter.FgHiRedColor},
	}

	for _, row := range data {
		table.Rich(
			[]string{row.label, fmt.Sprintf("%d", row.value)},
			[]tablewriter.Colors{
				{tablewriter.FgHiWhiteColor},
				{row.color},
			},
		)
	}

	fmt.Println()
	table.Render()
	fmt.Println()
}

func printStatRow(borderColor *color.Color, valueColor *color.Color, label string, value int, vertical string) {
	borderColor.Printf(vertical)
	fmt.Printf(" %-16s ", label)
	borderColor.Printf(vertical)
	valueColor.Printf(" %8d ", value)
	borderColor.Println(vertical)
}
