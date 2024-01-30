package models

// used https://mholt.github.io/json-to-go/ to generate the struct automatically
type Roles []struct {
	OrganisationID              string      `json:"OrganisationId"`
	Status                      string      `json:"Status"`
	RegistrationNumber          string      `json:"RegistrationNumber"`
	RegisteredName              string      `json:"RegisteredName"`
	ParentOrganisationReference interface{} `json:"ParentOrganisationReference"`
	OrgDomainClaims             []struct {
		AuthorisationDomainName string `json:"AuthorisationDomainName"`
		Status                  string `json:"Status"`
	} `json:"OrgDomainClaims"`
	OrgDomainRoleClaims []struct {
		AuthorisationDomainName     string `json:"AuthorisationDomainName"`
		AuthorisationDomainRoleName string `json:"AuthorisationDomainRoleName"`
		Status                      string `json:"Status"`
	} `json:"OrgDomainRoleClaims"`
}