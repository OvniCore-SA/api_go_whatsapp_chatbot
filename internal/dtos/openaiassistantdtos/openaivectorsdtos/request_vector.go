package openaivectorsdtos

type CreateVectorStoreRequest struct {
	Name string `json:"name"`
}

type FileCounts struct {
	InProgress int `json:"in_progress"`
	Completed  int `json:"completed"`
	Failed     int `json:"failed"`
	Cancelled  int `json:"cancelled"`
	Total      int `json:"total"`
}

type VectorStore struct {
	ID         string     `json:"id"`
	Object     string     `json:"object"`
	CreatedAt  int64      `json:"created_at"`
	Name       string     `json:"name"`
	Bytes      int        `json:"bytes"`
	FileCounts FileCounts `json:"file_counts"`
}
