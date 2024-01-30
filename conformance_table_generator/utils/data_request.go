package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/OpenBanking-Brasil/conformance/tree/main/conformance_table_generator/models"
	"github.com/joho/godotenv"
)

func importSubmittedFiles(repositoryUrl string) models.GithubTree {
	// import every file and folder from github repository
	client := http.Client{}

	req, err := http.NewRequest(http.MethodGet, repositoryUrl, nil)
	if err != nil {
		log.Fatal("Failed to make a request to github api: ", err)
	}

	envErr := godotenv.Load(".env")
	if envErr != nil{
		log.Print("No env file loaded")
	}
	token := "Bearer " + os.Getenv("PAT")
	req.Header.Add("Authorization", token)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Failed to obtain a response from github api: ", err)
	} else if resp.StatusCode != http.StatusOK {
		log.Fatal("Failed to obtain a 201 response from github api: ", resp.StatusCode)
	}

	defer resp.Body.Close()

	var respBody models.GithubResponse
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		log.Fatal("Failed to decode response from github api: ", err)
	}

	// filter response to only contain the files from submissions folder
	var filteredResults models.GithubTree
	for _, file := range respBody.Tree {
		filePath := file.Path
		filePathSplit := strings.Split(filePath, "/")
		
		isInSubmissionsFolder := filePathSplit[0] == "submissions"
		isTheCorrectLenght := len(filePathSplit) == 5
		isTheCorrectExtension := filePath[len(filePath) - 3:] == "zip" || filePath[len(filePath) - 4:] == "json"
		if isInSubmissionsFolder && isTheCorrectLenght  && isTheCorrectExtension {
			filteredResults = append(filteredResults, file)
		}
	}

	return filteredResults
}

func makeOrganisationsMap(useParentOrg bool) map[string]string {
	organisations := make(map[string]string)

	client := http.Client{}

	req, err := http.NewRequest(http.MethodGet, "https://data.directory.openbankingbrasil.org.br/roles", nil)
	if err != nil {
		log.Fatal("Failed to make a request to roles endpoint: ", err)
	}

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Fatal("Failed to obtain a response from roles endpoint: ", err)
	}

	defer resp.Body.Close()

	var roles models.Roles
	if err := json.NewDecoder(resp.Body).Decode(&roles); err != nil {
		log.Fatal("Failed to decode response from roles endpoint: ", err)
	}

	if useParentOrg {
		orgsWithParent := make(map[string]string)
		for _, role := range roles {
			if role.ParentOrganisationReference != nil && role.ParentOrganisationReference != role.RegistrationNumber {
				orgsWithParent[role.RegistrationNumber[:8]] = fmt.Sprintf("%v", role.ParentOrganisationReference)[:8]
			} else {
				organisations[role.RegistrationNumber[:8]] = role.RegisteredName
			}
		}
		for child, parent := range orgsWithParent {
			organisations[child] = organisations[parent]
		}
	} else {
		for _, role := range roles {
			organisations[role.RegistrationNumber[:8]] = role.RegisteredName
		}
	}

	return organisations
}
