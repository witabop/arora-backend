package layer

// lambda request struct
type RequestData struct {
	MaxID *int64 `json:"maxID"`
}

// lambda response struct
type ResponseData struct {
	Success  int8    `json:"success"`
	ValidIDs []int64 `json:"validIDs"`
}
