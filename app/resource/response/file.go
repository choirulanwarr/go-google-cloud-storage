package response

type UploadFileResponse struct {
	Path string `json:"path"`
}

type GetFileResponse struct {
	Name      string `json:"name"`
	URL       string `json:"url"`
	Size      string `json:"size"`
	Type      string `json:"type"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}
