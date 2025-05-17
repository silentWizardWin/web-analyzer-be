package service

import (
	"net/http"
	"web-analyzer-be/internal/model"
	"web-analyzer-be/internal/util"
)

func AnalyzeURL(url string) (*model.AnalyzeResponse, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := util.ParseHTML(resp.Body)
	if err != nil {
		return nil, err
	}

	return &model.AnalyzeResponse{
		HTMLVersion:     util.GetHTMLVersion(resp),
		Title:           util.GetTitle(doc),
		HeadingsCount:   util.CountHeadings(doc),
		LoginFormExists: util.ContainsLoginForm(doc),
		LinkAnalysis:    util.AnalyzeLinks(doc, url),
	}, nil
}
