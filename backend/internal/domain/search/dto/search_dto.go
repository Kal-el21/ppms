package dto

type SearchResult struct {
	EntityType string `json:"entity_type"`
	EntityID   uint64 `json:"entity_id"`
	Title      string `json:"title"`
	Snippet    string `json:"snippet"`
}

type SearchResponse struct {
	Projects []SearchResult `json:"projects"`
	Tasks    []SearchResult `json:"tasks"`
	Requests []SearchResult `json:"requests"`
}
