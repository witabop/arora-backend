package body

import "time"

type Creator struct {
	ID           *int64  `json:"id"`
	Name         *string `json:"name"`
	Type         *string `json:"type"`
	IsRNVAccount *bool   `json:"isRNVAccount"`
}

type Universe struct {
	ID                        *int64         `json:"id"`
	RootPlaceID               *int64         `json:"rootPlaceId"`
	Name                      *string        `json:"name"`
	Description               *string        `json:"description"`
	Creator                   *Creator       `json:"creator"`
	Price                     *interface{}   `json:"price"`
	AllowedGearGenres         *[]string      `json:"allowedGearGenres"`
	AllowedGearCategories     *[]interface{} `json:"allowedGearCategories"`
	IsGenreEnforced           *bool          `json:"isGenreEnforced"`
	CopyingAllowed            *bool          `json:"copyingAllowed"`
	Playing                   *int           `json:"playing"`
	Visits                    *int           `json:"visits"`
	MaxPlayers                *int           `json:"maxPlayers"`
	Created                   *time.Time     `json:"created"`
	Updated                   *time.Time     `json:"updated"`
	StudioAccessToApisAllowed *bool          `json:"studioAccessToApisAllowed"`
	CreateVipServersAllowed   *bool          `json:"createVipServersAllowed"`
	UniverseAvatarType        *string        `json:"universeAvatarType"`
	Genre                     *string        `json:"genre"`
	GenreL1                   *string        `json:"genre_l1"`
	GenreL2                   *string        `json:"genre_l2"`
	IsAllGenre                *bool          `json:"isAllGenre"`
	IsFavoritedByUser         *bool          `json:"isFavoritedByUser"`
	FavoritedCount            *int           `json:"favoritedCount"`
	LicenseDescription        *string        `json:"licenseDescription"`
}

type UResponse struct {
	Data *[]Universe `json:"data"`
}
