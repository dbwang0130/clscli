package cls

import (
	"context"

	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"
)

// SearchInput holds parameters for SearchLog.
type SearchInput struct {
	TopicId    string
	TopicIds   []string
	From       int64
	To         int64
	Query      string
	Limit      *int64
	Context    string
	Sort       string
	SyntaxRule *uint64
}

// SearchResult holds one page of SearchLog response.
type SearchResult struct {
	Context     string
	ListOver    bool
	Analysis    bool
	Results     []*cls.LogInfo
	ColNames    []*string
	Columns     []*cls.Column
	AnalysisRecords []*string
}

// Search runs SearchLog. If Context is set, it is used for pagination.
func (c *Client) Search(ctx context.Context, in SearchInput) (*SearchResult, error) {
	req := cls.NewSearchLogRequest()
	req.From = &in.From
	req.To = &in.To
	req.Query = &in.Query
	if in.SyntaxRule != nil {
		req.SyntaxRule = in.SyntaxRule
	} else {
		one := uint64(1)
		req.SyntaxRule = &one
	}
	useNew := true
	req.UseNewAnalysis = &useNew

	if in.TopicId != "" {
		req.TopicId = &in.TopicId
	}
	if len(in.TopicIds) > 0 {
		topics := make([]*cls.MultiTopicSearchInformation, 0, len(in.TopicIds))
		for _, id := range in.TopicIds {
			topics = append(topics, &cls.MultiTopicSearchInformation{TopicId: &id})
		}
		req.Topics = topics
	}
	if in.Limit != nil {
		req.Limit = in.Limit
	}
	if in.Context != "" {
		req.Context = &in.Context
	}
	if in.Sort != "" {
		req.Sort = &in.Sort
	}

	resp, err := c.api.SearchLogWithContext(ctx, req)
	if err != nil {
		return nil, err
	}
	p := resp.Response
	if p == nil {
		return &SearchResult{}, nil
	}
	out := &SearchResult{
		ListOver: p.ListOver != nil && *p.ListOver,
		Analysis: p.Analysis != nil && *p.Analysis,
	}
	if p.Context != nil {
		out.Context = *p.Context
	}
	if p.Results != nil {
		out.Results = p.Results
	}
	if p.ColNames != nil {
		out.ColNames = p.ColNames
	}
	if p.Columns != nil {
		out.Columns = p.Columns
	}
	if p.AnalysisRecords != nil {
		out.AnalysisRecords = p.AnalysisRecords
	}
	return out, nil
}
