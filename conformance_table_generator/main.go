package main

import (
	"flag"
	"log"
	"strings"

	"github.com/OpenBanking-Brasil/conformance/tree/main/conformance_table_generator/utils"
)

var (
	PhaseName string
	Apis string
)

func init() {
	flag.StringVar(&PhaseName, "phase", "", "The name of the phase (e.g., 'phase2')")
	flag.StringVar(&Apis, "apis", "", "APIs and Versions (e.g., 'accounts_2.1 credit-card_1.0 credit-card_2.1 resources_3.0')")
	flag.Parse()
}

func main() {
	if Apis == "" {
		log.Fatal("APIs and versions not provided. Please use the -apis flag to specify the APIs and versions.")
	}

	if PhaseName == "" {
		log.Fatal("Phase name not provided. Please use the -phase flag to specify the phase name.")
	}

	apisWithVersions := strings.Fields(Apis)
	if len(apisWithVersions) == 0 {
		log.Fatal("No APIs and versions provided.")
	}

	utils.GenerateTable(apisWithVersions, PhaseName)
}
