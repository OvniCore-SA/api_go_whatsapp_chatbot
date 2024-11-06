package openaivectorfiles

// Id de un file en OpenAI
type AddFileRequest struct {
	FileID string `json:"file_id"`
}

// Id de un file en OpenAI
type RequestVectorStorage struct {
	FileID        string `json:"file_id"`
	VectorStoreID string `json:"vector_storage_id"`
}

// Se utiliza para las operaciones CRUDs con vector files en OpenAI
type VectorStoreFile struct {
	ID            string `json:"id"`
	Object        string `json:"object"`
	CreatedAt     int64  `json:"created_at"`
	UsageBytes    int    `json:"usage_bytes"`
	VectorStoreID string `json:"vector_store_id"`
	Status        string `json:"status"`
	LastError     string `json:"last_error"`
}
