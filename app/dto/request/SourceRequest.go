package request

type UpdateSourceRequest struct {
	SourceName string `json:"source_name"`
	SourceDesc string `json:"source_desc"`
	SourceLink string `json:"source_link"`
	SourceType string `json:"source_type"`
	SourceExp  string `json:"source_exp"`
}

type CreateSourceRequest struct {
	SourceName string `json:"source_name"`
	SourceDesc string `json:"source_desc"`
	SourceLink string `json:"source_link"`
	SourceType int    `json:"source_type"`
	SourceExp  string `json:"source_exp"`
}
