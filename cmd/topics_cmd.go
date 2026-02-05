package cmd

import (
	"context"
	"fmt"

	"github.com/clscli/clscli/internal/cls"
	"github.com/clscli/clscli/internal/output"
	"github.com/spf13/cobra"
)

var (
	topicsTopicName  string
	topicsLogsetName string
	topicsLogsetId   string
	topicsOffset     int64
	topicsLimit      int64
)

var topicsCmd = &cobra.Command{
	Use:   "topics",
	Short: "List log topics",
	Long:  "List log topics in the given region. Use filters to find topic ID and region for query/context.",
	RunE:  runTopics,
}

func init() {
	rootCmd.AddCommand(topicsCmd)
	topicsCmd.Flags().StringVar(&topicsTopicName, "topic-name", "", "Filter by topic name (fuzzy match)")
	topicsCmd.Flags().StringVar(&topicsLogsetName, "logset-name", "", "Filter by logset name (fuzzy match)")
	topicsCmd.Flags().StringVar(&topicsLogsetId, "logset-id", "", "Filter by logset ID")
	topicsCmd.Flags().Int64Var(&topicsOffset, "offset", 0, "Pagination offset")
	topicsCmd.Flags().Int64Var(&topicsLimit, "limit", 20, "Max topics per request (max 100)")
}

func runTopics(cmd *cobra.Command, args []string) error {
	client, err := getCLSClient()
	if err != nil {
		return err
	}

	limit := topicsLimit
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	f, p := resolveOutput(cmd)
	writer, err := output.NewWriter(f, p)
	if err != nil {
		return err
	}
	defer writer.Close()

	in := cls.DescribeTopicsInput{
		TopicName:  topicsTopicName,
		LogsetName: topicsLogsetName,
		LogsetId:   topicsLogsetId,
		Offset:     topicsOffset,
		Limit:      limit,
		BizType:    0, // log topic
	}
	res, err := client.DescribeTopics(context.Background(), in)
	if err != nil {
		return fmt.Errorf("DescribeTopics: %w", err)
	}

	_ = writer.WriteTableHeaderTopics()
	for _, topic := range res.Topics {
		if err := writer.WriteTopicInfo(region, topic); err != nil {
			return err
		}
	}
	return writer.Flush()
}
