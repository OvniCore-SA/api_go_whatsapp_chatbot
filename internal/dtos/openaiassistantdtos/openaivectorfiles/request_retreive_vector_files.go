package openaivectorfiles

type RequestVectorStoreFile struct {
	ID            string  `json:"id"`
	Object        string  `json:"object"`
	CreatedAt     int64   `json:"created_at"`
	VectorStoreID string  `json:"vector_store_id"`
	Status        string  `json:"status"`
	LastError     *string `json:"last_error,omitempty"`
}

type RequestVectorStoreFiles struct {
	FileID        string `json:"file_id"`
	VectorStoreID string `json:"vector_store_id"`
}
