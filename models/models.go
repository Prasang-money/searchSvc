package models

//

type Country struct {
	Name       Name                  `json:"name"`
	Population int                   `json:"population"`
	Capital    []string              `json:"capital"`
	Currencies map[string]Currencies `json:"currencies"`
}

type Name struct {
	Common string `json:"common"`
}

type Currencies struct {
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

type CountryMetadata struct {
	Name       string `json:"name"`
	Population int    `json:"population"`
	Capital    string `json:"capital"`
	Currency   string `json:"currency"`
}
