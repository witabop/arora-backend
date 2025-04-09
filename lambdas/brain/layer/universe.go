package layer

// universe struct
type Universe struct {
	ID             *int64   `json:"id"`
	RootPlaceID    *int64   `json:"rootPlaceId"`
	Name           *string  `json:"name"`
	Description    *string  `json:"description"`
	Playing        *int     `json:"playing"`
	Visits         *int     `json:"visits"`
	FavoritedCount *int     `json:"favoritedCount"`
	PercentMatch   *float64 `json:"percentMatch"`
}

// roblox api response
type UniverseResponse struct {
	ValidUniverses *[]Universe `json:"validUniverses"`
}
