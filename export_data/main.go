package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/OpenBanking-Brasil/conformance/tree/main/export_data/utils"
)

// go run main.go -t <phaseNo> -v <versionNo>
// example
// go run main.go -t phase2 -v 2

var (
	Target            string
	Version           string
)

func init() {
	flag.StringVar(&Target, "t", "phase2", "Target Table")
	flag.StringVar(&Version, "v", "2", "API Version")
	flag.Parse()
}

func main() {
	// Versions that are allowed
	allowedVersions := []string {"1", "2"}
	if !utils.Contains(allowedVersions, Version) {
		log.Fatal("Version chosen is not allowed: ", Version)
	}

	resultsPathCsv := fmt.Sprintf("../results/%s/v%s/%s-v%s-data.csv", Target, Version, Target, Version)
	resultsPathMd  := fmt.Sprintf("../results/%s/v%s/%s-v%s-data.md" , Target, Version, Target, Version)
	// It's recommended to use semicolon as separator as some organisations might have comma in their names
	separator   := ';'
	var apiFamilyTypes []string
	var apiHeaderNames []string
	headers := []string {
		"Conglomerado",
		"Marca",
	}

	switch Target {
	case "phase1":
		apiFamilyTypes = []string{"acceptance-and-branches-abroad", "financial-risk", "responsibility"}

		apiHeaderNames = []string{""}
	case "phase2":
		apiFamilyTypes = []string {
			"consents",
			"customers-personal",
			"customers-business",
			"resources",
			"accounts",
			"credit-cards-accounts",
			"loans",
			"financings",
			"unarranged-accounts-overdraft",
			"invoice-financings",
		}
	
		apiHeaderNames = []string {
			"Consentimento API",
			"Dados Cadastrais (PF) API",
			"Dados Cadastrais (PJ) API",
			"Resources API",
			"Contas API",
			"Cartão de Crédito API",
			"Operações de Crédito - Empréstimos API",
			"Operações de Crédito - Financiamentos API",
			"Operações de Crédito - Adiantamento a Depositantes API",
			"Operações de Crédito - Direitos Creditórios Descontados API",
		}
	case "phase3":
		apiFamilyTypes = []string {
			"payments-pix",
		}

		apiHeaderNames = []string {
			"Pagamentos API",
		}
	// case "phase4":
	// 	apiFamilyTypes = []string {
	// 		"opendata-investments_funds",
	// 		"opendata-investments_bank-fixed-incomes",
	// 		"opendata-investments_credit-fixed-incomes",
	// 		"opendata-investments_variable-incomes",
	// 		"opendata-investments_treasure-titles",
	// 		"opendata-capitalization_bonds",
	// 		"opendata-exchange_online-rates",
	// 		"opendata-exchange_vet-values",
	// 		"opendata-acquiring-services_personals",
	// 		"opendata-acquiring-services_businesses",
	// 		"opendata-pension_risk-coverages",
	// 		"opendata-pension_survival-coverages",
	// 		"opendata-insurance_automotives",
	// 		"opendata-insurance_homes",
	// 		"opendata-insurance_personals",
	// 	}
	default:
		log.Fatalf("Invalid phase entered: %s", Target)
	}

	exportData(apiFamilyTypes, apiHeaderNames, Version, resultsPathCsv, separator)
	utils.FilterDuplicateEntries(resultsPathCsv, separator)
	// Specifically for phase 2, we should filter out entries that do not have certification for consents API
	// And for phase 3, we should filter out entries that do not have certifications for payments API
	if Target == "phase2" {
		utils.FilterEntriesWithoutSpecificApi(resultsPathCsv, separator, utils.GetIndex(apiFamilyTypes, "consents"))
	} else if Target == "phase3" {
		utils.FilterEntriesWithoutSpecificApi(resultsPathCsv, separator, utils.GetIndex(apiFamilyTypes, "payments-pix"))
	}
	headers = append(headers, apiHeaderNames...)
	utils.GenerateFromCsv(resultsPathCsv, resultsPathMd, headers, separator)
}

func exportData(apiFamilyTypes []string, apiHeaderNames []string, version string, fileName string, separator rune) {
	// Creating the header for the table
	tableHeader := []string{"Conglomerado", "Marca"}
	tableHeader = append(tableHeader, apiHeaderNames...)

	// Requesting data from the API participants endpoint
	participants, err := utils.ImportData("https://data.directory.openbankingbrasil.org.br/participants")
	if err != nil {
		log.Fatal("Failed to request data from the participants API:", err)
	}

	// Creating the map from registration number to registered name
	organisations := utils.MakeOrganisationsMap()

	// Creating the csv file
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal("Failed to open file: ", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = separator
	defer writer.Flush()

	// Writing header to the file
	if err := writer.Write(tableHeader); err != nil {
		log.Fatal("Failed to write to file:", err)
	}

	// Writing data to the file
	for _, participant := range participants {
		for _, server := range participant.AuthorisationServers {
			row_elements := make(map[string]string)
			// We look for the parent organisation in order to get the conglomerate
			if participant.ParentOrganisationReference != "" {
				row_elements["Conglomerado"] = organisations[participant.ParentOrganisationReference]
			} else {
				row_elements["Conglomerado"] = participant.RegisteredName
			}
			row_elements["Marca"] = server.CustomerFriendlyName

			// Iterate through all servers
			for _, resource := range server.APIResources {
				// The family type must be in apiFamilyTypes and there must be an APICertificationURI
				if utils.Contains(apiFamilyTypes, resource.APIFamilyType) && resource.APICertificationURI != nil && utils.IsRightVersion(resource.APIVersion, version) {
					// Search for the date in the file containing the certification
					certDate := utils.DateFromFileName(fmt.Sprintf("%v", resource.APICertificationURI))
					if certDate == "" {
						certDate = "No date"
					}
					// If the date is available in the endpoint, it should overwrite the one from the file
					if resource.CertificationStartDate != nil {
						certDate = utils.ConvertDate(fmt.Sprintf("%v", resource.CertificationStartDate))
					}
					row_elements[resource.APIFamilyType] = fmt.Sprintf(
						"[%s](%s)",
						certDate,
						resource.APICertificationURI,
					)
				}
			}

			row := make([]string, len(apiFamilyTypes) + 2)
			row[0] = row_elements["Conglomerado"]
			row[1] = row_elements["Marca"]
			for i, familyType := range apiFamilyTypes {
				row[i + 2] = row_elements[familyType]
			}

			if err := writer.Write(row); err != nil {
				log.Fatal("Failed to write to file:", err)
			}
		}
	}
}