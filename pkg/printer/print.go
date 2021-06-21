package printer

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/jkandasa/jenkinsctl/pkg/utils"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v2"
)

// output types
const (
	OutputConsole = "console"
	OutputYAML    = "yaml"
	OutputJSON    = "json"
)

func Print(out io.Writer, headers []string, data interface{}, hideHeader bool, output string, pretty bool) {
	switch output {
	case OutputConsole:
		dataConsole, ok := data.([]interface{})
		if !ok {
			fmt.Fprintln(out, "data not in table format")
			return
		}
		PrintConsole(out, headers, dataConsole, hideHeader)
		return

	case OutputJSON:
		var jsonBytes []byte
		var err error
		if pretty {
			jsonBytes, err = json.MarshalIndent(data, "", " ")
		} else {
			jsonBytes, err = json.Marshal(data)
		}
		if err != nil {
			fmt.Println("error on converting to json", err)
			return
		}
		fmt.Fprint(out, string(jsonBytes))

	case OutputYAML:
		bytes, err := yaml.Marshal(data)
		if err != nil {
			fmt.Println("error on converting to yaml", err)
			return
		}
		fmt.Fprint(out, string(bytes))
	}
}

func PrintConsole(out io.Writer, headers []string, data []interface{}, hideHeader bool) {
	// convert the data
	rows := make([][]string, 0)
	for index := range data {
		structData := data[index]
		mapData := utils.StructToMap(structData)
		row := make([]string, 0)
		for _, key := range headers {
			if value, ok := mapData[key]; ok {
				row = append(row, fmt.Sprintf("%v", value))
			} else {
				row = append(row, "")
			}
		}
		rows = append(rows, row)
	}

	table := tablewriter.NewWriter(os.Stdout)
	if !hideHeader {
		table.SetHeader(headers)
		table.SetAutoFormatHeaders(true)
		table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	}
	table.SetAutoWrapText(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t") // pad with tabs
	table.SetNoWhiteSpace(true)
	table.AppendBulk(rows) // Add Bulk Data
	table.Render()
}
