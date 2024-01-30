package utils

import (
	"fmt"
	"strings"

	"github.com/OpenBanking-Brasil/conformance/tree/main/conformance_table_generator/models"
)

func GenerateTable(apisWithVersions []string, phase string) {
	// import files
	repositoryUrl := "https://api.github.com/repos/OpenBanking-Brasil/conformance/git/trees/main?recursive=1"
	submissionFiles := importSubmittedFiles(repositoryUrl)

	// filter files by chosen APIs and versions
	filteredFiles := filterFilesByApisAndVersions(submissionFiles, apisWithVersions)

	// create table and dump
	tableHeaders := []string{"Organização", "Deployment"}
	apis := extractAPIsWithVersions(apisWithVersions)
	tableHeaders = append(tableHeaders, apis...)
	table := [][]string{tableHeaders}

	dumpHeaders := []string{"Id da Organização", "Deployment", "API", "Version", "Data"}
	dump := [][]string{dumpHeaders}

	organisationsMap := makeOrganisationsMap(false)

	for _, file := range filteredFiles {
		fileSplit := strings.Split(file, "/")
		fileName := fileSplit[len(fileSplit)-1]
		fileNameSplit := strings.Split(fileName, "_")
		orgId := fileNameSplit[0]
		orgName := organisationsMap[orgId]
		deploymentName := fileNameSplit[1]
		api := fileNameSplit[2]
		version := fileNameSplit[3]
		date := strings.Split(fileNameSplit[4], ".")[0]

		fileUrl := strings.Replace(repositoryUrl, "api.github.com/repos", "github.com", 1)
		fileUrl = strings.Replace(fileUrl, "git/trees/main?recursive=1", "blob/main/"+file, 1)
		fileUrl = strings.Replace(fileUrl, " ", "%20", -1)

		dump = append(dump, []string{
			orgId,
			deploymentName,
			api,
			version,
			date,
		})

		apiIndex := findAPIIndex(apisWithVersions, api, version)
		if ind := searchFileInTable(table, orgName, deploymentName); ind == -1 {
			newRow := make([]string, len(tableHeaders))
			newRow[0] = orgName
			newRow[1] = deploymentName
			newRow[apiIndex+2] = fmt.Sprintf("[%s](%s)", date, fileUrl)

			table = append(table, newRow)
		} else {
			table[ind][apiIndex+2] = fmt.Sprintf("[%s](%s)", date, fileUrl)
		}
	}

	// You should specify the 'phase' variable for file naming.
	dumpFileName := fmt.Sprintf("../results/%s/%s-conformance-dump.csv", phase, phase)
	exportTable(dump, dumpFileName)

	tableFileName := fmt.Sprintf("../results/%s/%s-conformance-table.csv", phase, phase)
	exportTable(table, tableFileName)
}

func filterFilesByApisAndVersions(submissionFiles models.GithubTree, apisWithVersions []string) []string {
	var filteredFiles []string

	for _, file := range submissionFiles {
		filePath := file.Path
		fileSplit := strings.Split(filePath, "/")
		fileApi := fileSplit[2]
		fileVersion := fileSplit[3]

		if fileApi == "payments-webhook" {
			fileApi = strings.Split(fileApi, "-")[1]
		}

		if findAPIIndex(apisWithVersions, fileApi, fileVersion) != -1 {
			filteredFiles = append(filteredFiles, filePath)
		}
	}

	return filteredFiles
}

func extractAPIsWithVersions(apisWithVersions []string) []string {
	var apisWithVersionsLabels []string
	for _, apiWithVersion := range apisWithVersions {
		parts := strings.Split(apiWithVersion, "_")
		if len(parts) == 2 {
			api := parts[0]
			version := parts[1]
			apisWithVersionsLabels = append(apisWithVersionsLabels, fmt.Sprintf("%s v%s", api, version))
		}
	}
	return apisWithVersionsLabels
}

func findAPIIndex(apisWithVersions []string, api string, version string) int {
	correctedVersion, _ := strings.CutSuffix(version, "-OL")
	correctedVersion, _ = strings.CutPrefix(correctedVersion, "v")

	if len(correctedVersion) == 1 {
		correctedVersion = correctedVersion + ".0"
	} else if strings.Count(correctedVersion, ".") == 2 {
		lastDotIndex := strings.LastIndex(correctedVersion, ".")
		correctedVersion = correctedVersion[:lastDotIndex]
	}

	for i, apiWithVersion := range apisWithVersions {
		parts := strings.Split(apiWithVersion, "_")
		if len(parts) == 2 && parts[0] == api && parts[1] == correctedVersion {
			return i
		}
	}
	return -1
}
