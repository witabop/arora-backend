package layer

// universe search criteria
type SearchCriteria struct {
	Name           *string `json:"name"`
	Description    *string `json:"description"`
	Playing        *int    `json:"playing"`
	Visits         *int    `json:"visits"`
	FavoritedCount *int    `json:"favoritedCount"`
}

// lambda request struct
type RequestData struct {
	MaxID          int64          `json:"maxID"`
	SearchCriteria SearchCriteria `json:"searchCriteria"`
}

// lambda response struct
type ResponseData struct {
	ValidUniverses []Universe `json:"validUniverses"`
}
