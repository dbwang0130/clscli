package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/clscli/clscli/internal/cls"
	"github.com/clscli/clscli/internal/output"
	"github.com/spf13/cobra"
)

var (
	queryQuery   string
	queryTopic   string
	queryTopics  []string
	queryLast    string
	queryFrom    int64
	queryTo      int64
	queryLimit   int64
	queryMax     int64
	querySort    string
	queryContext string
)

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Search and analyze logs",
	Long:  "Query CLS with CQL or SQL. Use -q for query string, -t or --topics for topic(s), --last or --from/--to for time range.",
	RunE:  runQuery,
}

func init() {
	rootCmd.AddCommand(queryCmd)
	queryCmd.Flags().StringVarP(&queryQuery, "query", "q", "", "Search condition or SQL (e.g. level:ERROR or * | select count(*) as cnt)")
	queryCmd.MarkFlagRequired("query")
	queryCmd.Flags().StringVarP(&queryTopic, "topic", "t", "", "Single topic ID")
	queryCmd.Flags().StringSliceVar(&queryTopics, "topics", nil, "Comma-separated topic IDs (max 50); cannot use with -t")
	queryCmd.Flags().StringVar(&queryLast, "last", "", "Time range from now (e.g. 1h, 30m, 24h)")
	queryCmd.Flags().Int64Var(&queryFrom, "from", 0, "Start time Unix ms (or use --last)")
	queryCmd.Flags().Int64Var(&queryTo, "to", 0, "End time Unix ms (or use --last)")
	queryCmd.Flags().Int64Var(&queryLimit, "limit", 100, "Logs per request (max 1000)")
	queryCmd.Flags().Int64Var(&queryMax, "max", 0, "Max total logs (auto-paginate until this or ListOver); 0 = single request")
	queryCmd.Flags().StringVar(&querySort, "sort", "desc", "Sort: asc or desc")
	queryCmd.Flags().StringVar(&queryContext, "context", "", "Context cursor from previous response (manual pagination)")
}

func runQuery(cmd *cobra.Command, args []string) error {
	client, err := getCLSClient()
	if err != nil {
		return err
	}
	from, to, err := resolveTimeRange()
	if err != nil {
		return err
	}
	topicID, topicIDs, err := resolveTopics()
	if err != nil {
		return err
	}

	f, p := resolveOutput(cmd)
	writer, err := output.NewWriter(f, p)
	if err != nil {
		return err
	}
	defer writer.Close()

	limit := queryLimit
	if limit <= 0 {
		limit = 100
	}
	if limit > 1000 {
		limit = 1000
	}
	syntaxRule := uint64(1)
	in := cls.SearchInput{
		TopicId:    topicID,
		TopicIds:   topicIDs,
		From:       from,
		To:         to,
		Query:      queryQuery,
		Limit:      &limit,
		Sort:       querySort,
		SyntaxRule: &syntaxRule,
	}
	if queryContext != "" {
		in.Context = queryContext
	}

	ctx := context.Background()
	total := int64(0)
	maxTotal := queryMax
	first := true

	for {
		res, err := client.Search(ctx, in)
		if err != nil {
			return err
		}
		if res.Analysis {
			if first {
				_ = writer.WriteAnalysisRecords(res.Columns, res.AnalysisRecords)
				_ = writer.Flush()
			}
			break
		}
		if first {
			_ = writer.WriteTableHeaderLogInfo()
		}
		for _, log := range res.Results {
			if err := writer.WriteLogInfo(log); err != nil {
				return err
			}
			total++
			if maxTotal > 0 && total >= maxTotal {
				_ = writer.Flush()
				return nil
			}
		}
		_ = writer.Flush()
		if res.ListOver || res.Context == "" {
			break
		}
		in.Context = res.Context
		first = false
		if maxTotal > 0 && total >= maxTotal {
			break
		}
	}
	return nil
}

func resolveTimeRange() (from, to int64, err error) {
	if queryLast != "" {
		d, err := time.ParseDuration(queryLast)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid --last %q: %w", queryLast, err)
		}
		now := time.Now()
		end := now.UnixMilli()
		start := now.Add(-d).UnixMilli()
		return start, end, nil
	}
	if queryFrom > 0 && queryTo > 0 {
		return queryFrom, queryTo, nil
	}
	return 0, 0, fmt.Errorf("set --last (e.g. 1h) or both --from and --to (Unix ms)")
}

func resolveTopics() (topicID string, topicIDs []string, err error) {
	if queryTopic != "" && len(queryTopics) > 0 {
		return "", nil, fmt.Errorf("use either -t/--topic or --topics, not both")
	}
	if queryTopic != "" {
		return queryTopic, nil, nil
	}
	if len(queryTopics) > 0 {
		if len(queryTopics) > 50 {
			return "", nil, fmt.Errorf("--topics: max 50 topics")
		}
		return "", queryTopics, nil
	}
	return "", nil, fmt.Errorf("set -t/--topic or --topics")
}

// resolveOutput returns format and output path from --format, -o/--out, --output.
// When writing to a file, format is inferred from extension: .json -> json, .csv -> csv.
// If explicit format (--format or --output=json/csv) conflicts with extension, a warning is printed and extension wins.
func resolveOutput(cmd *cobra.Command) (formatOut string, outPath string) {
	formatOut = "csv"
	if format != "" {
		formatOut = format
	}
	outPath = outFile
	if outputFlag != "" {
		switch outputFlag {
		case "json", "csv":
			formatOut = outputFlag
		default:
			outPath = outputFlag
		}
	}
	if outPath != "" {
		if inferred := formatFromExt(outPath); inferred != "" {
			if inferred != formatOut {
				fmt.Fprintf(os.Stderr, "warning: output format %q overridden by file extension to %q\n", formatOut, inferred)
			}
			formatOut = inferred
		}
	}
	return formatOut, outPath
}

// formatFromExt returns format inferred from file extension: .json -> json, .csv -> csv.
func formatFromExt(path string) string {
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(path), "."))
	switch ext {
	case "json":
		return "json"
	case "csv":
		return "csv"
	default:
		return ""
	}
}
