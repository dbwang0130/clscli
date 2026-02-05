package cmd

import (
	"context"
	"fmt"

	"github.com/clscli/clscli/internal/cls"
	"github.com/clscli/clscli/internal/output"
	"github.com/spf13/cobra"
)

var (
	contextTopic string
)

var contextCmd = &cobra.Command{
	Use:   "context [PkgId] [PkgLogId]",
	Short: "Get log context",
	Long:  "Retrieve context logs around a given log (PkgId and PkgLogId from SearchLog results).",
	Args:  cobra.ExactArgs(2),
	RunE:  runContext,
}

func init() {
	rootCmd.AddCommand(contextCmd)
	contextCmd.Flags().StringVarP(&contextTopic, "topic", "t", "", "Topic ID (required)")
	contextCmd.MarkFlagRequired("topic")
}

func runContext(cmd *cobra.Command, args []string) error {
	client, err := getCLSClient()
	if err != nil {
		return err
	}
	pkgID := args[0]
	pkgLogIDStr := args[1]
	pkgLogID, err := cls.PkgLogIdFromString(pkgLogIDStr)
	if err != nil {
		return fmt.Errorf("invalid PkgLogId %q: %w", pkgLogIDStr, err)
	}

	f, p := resolveOutput(cmd)
	writer, err := output.NewWriter(f, p)
	if err != nil {
		return err
	}
	defer writer.Close()

	in := cls.GetContextInput{
		TopicId:  contextTopic,
		PkgId:    pkgID,
		PkgLogId: pkgLogID,
	}
	logs, err := client.GetContext(context.Background(), in)
	if err != nil {
		return err
	}
	_ = writer.WriteTableHeaderLogContext()
	for _, log := range logs {
		if err := writer.WriteLogContextInfo(log); err != nil {
			return err
		}
	}
	return writer.Flush()
}
