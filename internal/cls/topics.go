package cls

import (
	"context"

	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"
)

// DescribeTopicsInput holds parameters for DescribeTopics.
type DescribeTopicsInput struct {
	TopicName  string
	LogsetName string
	LogsetId   string
	Offset     int64
	Limit      int64
	BizType    uint64 // 0: log topic (default), 1: metric topic
}

// DescribeTopicsResult holds DescribeTopics response.
type DescribeTopicsResult struct {
	Topics     []*cls.TopicInfo
	TotalCount int64
}

// DescribeTopics returns log/metric topic list for the client's region.
func (c *Client) DescribeTopics(ctx context.Context, in DescribeTopicsInput) (*DescribeTopicsResult, error) {
	req := cls.NewDescribeTopicsRequest()
	if in.Limit <= 0 {
		in.Limit = 20
	}
	if in.Limit > 100 {
		in.Limit = 100
	}
	req.Offset = &in.Offset
	req.Limit = &in.Limit
	bizType := in.BizType
	req.BizType = &bizType

	var filters []*cls.Filter
	if in.TopicName != "" {
		filters = append(filters, &cls.Filter{
			Key:    stringPtr("topicName"),
			Values: []*string{&in.TopicName},
		})
	}
	if in.LogsetName != "" {
		filters = append(filters, &cls.Filter{
			Key:    stringPtr("logsetName"),
			Values: []*string{&in.LogsetName},
		})
	}
	if in.LogsetId != "" {
		filters = append(filters, &cls.Filter{
			Key:    stringPtr("logsetId"),
			Values: []*string{&in.LogsetId},
		})
	}
	if len(filters) > 0 {
		req.Filters = filters
	}

	resp, err := c.api.DescribeTopicsWithContext(ctx, req)
	if err != nil {
		return nil, err
	}
	p := resp.Response
	if p == nil {
		return &DescribeTopicsResult{}, nil
	}
	out := &DescribeTopicsResult{}
	if p.Topics != nil {
		out.Topics = p.Topics
	}
	if p.TotalCount != nil {
		out.TotalCount = *p.TotalCount
	}
	return out, nil
}

func stringPtr(s string) *string { return &s }
