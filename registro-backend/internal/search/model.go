package search

type SearchResultItem struct {
	ID          string `json:"id"`
	Type        string `json:"type"` // 'student', 'teacher', 'grade', 'communication', 'lesson'
	Title       string `json:"title"`
	Subtitle    string `json:"subtitle,omitempty"`
	Description string `json:"description,omitempty"`
	URL         string `json:"url,omitempty"`
}

type SearchResponse struct {
	Query   string             `json:"query"`
	Total   int                `json:"total"`
	Results []SearchResultItem `json:"results"`
}
