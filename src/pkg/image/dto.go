package image

// UploadResponse represents the response after a successful image upload
type UploadResponse struct {
	Message string     `json:"message"`
	Data    UploadData `json:"data"`
}

// UploadData contains the uploaded image URL
type UploadData struct {
	ImageURL string `json:"imageUrl"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error"`
}