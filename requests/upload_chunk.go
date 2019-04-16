package requests

type UploadChunkBody struct {
	ID   int    `json:"id" binding:"required"`
	Size int    `json:"size"`
	Data string `json:"data" binding:"required"`
}
