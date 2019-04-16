package requests

type RegisterBody struct {
	Sha256    string `json:"sha256" binding:"required"`
	Size      int    `json:"size" binding:"required"`
	ChunkSize int    `json:"chunk_size" binding:"required"`
}
