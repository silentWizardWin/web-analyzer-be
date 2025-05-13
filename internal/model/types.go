package model

type AnalyzeRequest struct {
	URL string `json:"url"`
}

type AnalyzeResponse struct {
	HTMLVersion     string            `json:"html_version"`
	Title           string            `json:"title"`
	HeadingsCount   map[string]int    `json:"headings_count"`
	LoginFormExists bool              `json:"login_form_exists"`
	LinkAnalysis    LinkAnalysis      `json:"link_analysis"`
}

type LinkAnalysis struct {
	InternalLinks     int `json:"internal_links"`
	ExternalLinks     int `json:"external_links"`
	InaccessibleLinks int `json:"inaccessible_links"`
}
