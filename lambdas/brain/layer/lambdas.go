package layer

// universe search criteria
type SearchCriteria struct {
	Name           *string `json:"name"`
	Description    *string `json:"description"`
	Playing        *uint64 `json:"playing"`
	Visits         *uint64 `json:"visits"`
	FavoritedCount *uint64 `json:"favoritedCount"`
}

// lambda request struct
type RequestData struct {
	NumGames       uint8          `json:"numGames"`
	SearchCriteria SearchCriteria `json:"searchCriteria"`
}

// lambda response struct
type ResponseData struct {
	Data []Universe `json:"data"`
}
