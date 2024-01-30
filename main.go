package main

import (
	"encoding/csv"
	"flag"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	alphabet "github.com/rh2g17/md-brasilian-alphabet-sort"
)

var (
	Target string
)

/*
* To generate all tables - go run main.go
* To generate phase 2 table - go run main.go -t phase2
* To generate phase 2 v2 table - go run main.go -t phase2v2
* To generate phase 3 table - go run main.go -t phase3
* To generate phase 4 table - go run main.go -t phase4
 */

func init() {
	flag.StringVar(&Target, "t", "all", "Target Table")
	flag.Parse()
}

func main() {
	if Target == "phase2" || Target == "all" {
		generateFromCsv("./phase2-data.csv", "./phase2-output.txt", []string{"Organisation", "Deployment", "Consentimento API", "Dados Cadastrais (PF) API", "Dados Cadastrais (PJ) API", "Resources API", "Contas API", "Cartão de Crédito API", "Operações de Crédito - Empréstimos API", "Operações de Crédito - Financiamentos API", "Operações de Crédito - Adiantamento a Depositantes API", "Operações de Crédito - Direitos Creditórios Descontados API"})
	}
	if Target == "phase2v2" || Target == "all" {
		generateFromCsv("./phase2v2-data.csv", "./phase2v2-output.txt", []string{"Organisation", "Deployment", "Consentimento API", "Dados Cadastrais (PF) API", "Dados Cadastrais (PJ) API", "Resources API", "Contas API", "Cartão de Crédito API", "Operações de Crédito - Empréstimos API", "Operações de Crédito - Financiamentos API", "Operações de Crédito - Adiantamento a Depositantes API", "Operações de Crédito - Direitos Creditórios Descontados API"})
	}
	if Target == "phase3" || Target == "all" {
		generateFromCsv("./phase3-data.csv", "./phase3-output.txt", []string{"Organisation", "Deployment", "MANU/DICT/INIC - T0", "MANU/DICT/INIC/QRES/QRDN - T2"})
	}
	if Target == "phase4" || Target == "all" {
		// comment when phase 4 comes out
		//generateFromCsv("./phase4-data.csv", "./phase4-output.txt", []string{"Brand Name", "Accounts", "Admin", "Channels", "Consents", "Credit Cards Accounts", "Customers Business", "Customers Personal", "Discovery", "Financings", "Invoice Financings", "Loans", "Payments Consents", "Payments Pix", "Products Services", "Resources", "Unarranged Accounts Overdraft"})

		// uncomment when phase 4 comes out
		// generateFromCsv("./phase4-data.csv", "./phase4-output.txt", []string{"Brand Name", "Open Data Investments - Funds", "Open Data Investments - Bank Fixed Incomes", "Open Data Investments - Credit Fixed Incomes", "Open Data Investments - Variable Incomes", "Open Data Investments - Treasure Titles", "Open Data Capitalization - Bonds", "Open Data Exchange - Online Rates", "Open Data Exchange - Vet Values", "Open Data Acquiring Services - Personals", "Open Data Acquiring Services - Businesses", "Open Data Pension - Risk Coverages", "Open Data Pension - Survival Coverages", "Open Data Insurance - Automotives", "Open Data Insurance - Homes", "Open Data Insurance - Personals"})
	}
}

func generateFromCsv(inputFile string, outputFile string, headers []string) {
	f, _ := os.Open(inputFile)
	defer f.Close()

	//Read lines from file
	lines, _ := csv.NewReader(f).ReadAll()
	lines = lines[1:]
	var sortLines []string

	for _, line := range lines {
		joinedLine := strings.Join(line, ",")
		sortLines = append(sortLines, joinedLine)
	}
	sortLines = alphabet.MergeSort(sortLines)
	lines = [][]string{}

	for _, item := range sortLines {
		split := strings.Split(item, ",")
		lines = append(lines, split)
	}

	//Set the table to output as a string
	tableOutput := &strings.Builder{}
	table := tablewriter.NewWriter(tableOutput)

	var indexHeaders []string

	for index := range headers {
		indexHeaders = append(indexHeaders, strconv.Itoa(index))
	}

	//Configure table
	table.SetHeader(indexHeaders)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetAutoWrapText(false)
	table.SetCenterSeparator("|")
	table.AppendBulk(lines)
	table.Render()

	//Open output file
	output, _ := os.OpenFile(outputFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0775)

	toWrite := tableOutput.String()

	//Replace headersxs
	for index, value := range headers {
		toWrite = strings.Replace(toWrite, strconv.Itoa(index), value, 1)
	}

	//Write result of table to file
	output.Write([]byte(toWrite))
	output.Close()
}
