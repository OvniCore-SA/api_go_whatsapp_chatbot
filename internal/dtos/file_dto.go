package dtos

type FileDto struct {
	ID                       int64  `json:"id"`
	AssistantsID             int64  `json:"assistants_id"`
	OpenaiFilesID            string `json:"openai_files_id"`
	Filename                 string `json:"filename"`
	Purpose                  string `json:"purpose"`
	OpenaiVectorStoreIDs     string `json:"openai_vector_store_ids"`
	OpenaiVectorStoreFileIDs string `json:"openai_vector_store_file_ids"`
}
