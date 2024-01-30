package utils

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	alphabet "github.com/rh2g17/md-brasilian-alphabet-sort"

	"github.com/OpenBanking-Brasil/conformance/tree/main/export_data/models"
)

func ImportData(url string) (models.Participants, error) {
	client := http.Client{}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, err
	}

	defer resp.Body.Close()

	var data models.Participants
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Println(err)
	}

	return data, nil
}

func Contains(arr []string, str string) bool {
	for _, s := range arr {
		if s == str {
			return true
		}
	}
	return false
}

func GetIndex(arr []string, str string) int {
	for i, s := range arr {
		if s == str {
			return i
		}
	}
	return -1
}

func MakeOrganisationsMap() map[string]string {
	organisations := make(map[string]string)

	// Requesting data from the API roles endpoint
	roles, err := ImportData("https://data.directory.openbankingbrasil.org.br/roles")
	if err != nil {
		log.Fatal("Failed to request data from the roles API:", err)
	}

	for _, role := range roles {
		organisations[role.RegistrationNumber] = role.RegisteredName
	}

	return organisations
}

func IsRightVersion(apiVersion string, targetVersion string) bool {
	return strings.Split(apiVersion, ".")[0] == targetVersion
}

// Copied from the main.go file in the root directory
func GenerateFromCsv(inputFile string, outputFile string, headers []string, separator rune) {
	f, err := os.Open(inputFile)
	if err != nil {
		log.Fatal("Failed to open file: ", err)
	}
	defer f.Close()

	// Read lines from file
	reader := csv.NewReader(f)
	reader.Comma = separator
	lines, err := reader.ReadAll()
	if err != nil {
		log.Fatal("Failed to read csv file: ", err)
	}
	lines = lines[1:]
	var sortLines []string

	for _, line := range lines {
		joinedLine := strings.Join(line, string(separator))
		sortLines = append(sortLines, joinedLine)
	}
	sortLines = alphabet.MergeSort(sortLines)
	lines = [][]string{}

	for _, item := range sortLines {
		split := strings.Split(item, string(separator))
		lines = append(lines, split)
	}

	// Set the table to output as a string
	tableOutput := &strings.Builder{}
	table := tablewriter.NewWriter(tableOutput)

	var indexHeaders []string

	for index := range headers {
		indexHeaders = append(indexHeaders, strconv.Itoa(index))
	}

	// Configure table
	table.SetHeader(indexHeaders)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetAutoWrapText(false)
	table.SetCenterSeparator("|")
	table.AppendBulk(lines)
	table.Render()

	// Open output file
	output, err := os.OpenFile(outputFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0775)
	if err != nil {
		log.Fatal("Failed to open file: ", err)
	}

	toWrite := tableOutput.String()

	// Replace headersxs
	for index, value := range headers {
		toWrite = strings.Replace(toWrite, strconv.Itoa(index), value, 1)
	}

	// Write result of table to file
	output.Write([]byte(toWrite))
	output.Close()
}

func FilterEntriesWithoutSpecificApi(inputFile string, separator rune, apiIndex int) {
	fileRead, err := os.Open(inputFile)
	if err != nil {
		log.Fatal("Failed to open file: ", err)
	}

	// Read lines from file
	reader := csv.NewReader(fileRead)
	reader.Comma = separator
	lines, err := reader.ReadAll()
	if err != nil {
		log.Fatal("Failed to read csv file: ", err)
	}
	fileRead.Close()

	// Reopen file
	fileWrite, err := os.Create(inputFile)
	if err != nil {
		log.Fatal("Failed to open file: ", err)
	}
	defer fileWrite.Close()

	// Create writer
	writer := csv.NewWriter(fileWrite)
	writer.Comma = separator
	defer writer.Flush()

	// Writing header to the file
	if err := writer.Write(lines[0]); err != nil {
		log.Fatal("Failed to write to file: ", err)
	}
	lines = lines[1:]

	// Writing the rest of the lines filtering entries without data in consents api column
	for _, line := range lines {
		if line[apiIndex + 2] != "" {
			if err := writer.Write(line); err != nil {
				log.Fatal("Failed to write to file: ", err)
			}
		}
	}
}

func FilterDuplicateEntries(inputFile string, separator rune) {
	fileRead, err := os.Open(inputFile)
	if err != nil {
		log.Fatal("Failed to open file: ", err)
	}

	// Read lines from file
	reader := csv.NewReader(fileRead)
	reader.Comma = separator
	lines, err := reader.ReadAll()
	if err != nil {
		log.Fatal("Failed to read csv file: ", err)
	}
	fileRead.Close()

	// Reopen file
	fileWrite, err := os.Create(inputFile)
	if err != nil {
		log.Fatal("Failed to open file: ", err)
	}
	defer fileWrite.Close()

	// Create writer
	writer := csv.NewWriter(fileWrite)
	writer.Comma = separator
	defer writer.Flush()

	// Writing header to the file
	if err := writer.Write(lines[0]); err != nil {
		log.Fatal("Failed to write to file: ", err)
	}
	lines = lines[1:]

	// Keep only one instance of exactly equal lines
	isWritten := make(map[string]bool)
	for _, line := range lines {
		if _, written := isWritten[strings.Join(line, " ")]; !written {
			isWritten[strings.Join(line, " ")] = true
			if err := writer.Write(line); err != nil {
				log.Fatal("Failed to write to file: ", err)
			}
		}
	}
}

func DateFromFileName(fileName string) string {
	re, err := regexp.Compile(`\d+-\D{3}-\d{4}`)
	if err != nil {
		log.Fatal("Could not create regular expression: ", err)
	}
	
	return re.FindString(fileName)
}

func ConvertDate(date string) string {
	t, err := time.Parse("02/01/2006", date)
	if err != nil {
		log.Fatalf("Failed to convert date (%s): %s", date, err)
	}
	
	return t.Format("02-Jan-2006")
}
