package cmd

import (
	"fmt"
	"os"

	"github.com/clscli/clscli/internal/cls"
	"github.com/spf13/cobra"
)

var (
	region string
	format string
	outFile string
)

var rootCmd = &cobra.Command{
	Use:   "clscli",
	Short: "CLI for Tencent Cloud CLS (Cloud Log Service)",
	Long:  "Retrieve and analyze CLS logs: query (-q), context (-c).",
}

func init() {
	rootCmd.PersistentFlags().StringVar(&format, "format", "csv", "Output format: json or csv")
	rootCmd.PersistentFlags().StringVarP(&outFile, "out", "o", "", "Write output to file (default stdout)")
	rootCmd.PersistentFlags().StringVar(&region, "region", "", "CLS region (e.g. ap-guangzhou)")
	rootCmd.PersistentFlags().StringVar(&outputFlag, "output", "", "Output: json, csv, or file path (overrides format/out)")
}

var outputFlag string

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func getCLSClient() (*cls.Client, error) {
	if region == "" {
		return nil, fmt.Errorf("region required: use --region (e.g. ap-guangzhou)")
	}
	secretID := os.Getenv("TENCENTCLOUD_SECRET_ID")
	secretKey := os.Getenv("TENCENTCLOUD_SECRET_KEY")
	if secretID == "" || secretKey == "" {
		return nil, fmt.Errorf("credentials required: set TENCENTCLOUD_SECRET_ID and TENCENTCLOUD_SECRET_KEY")
	}
	return cls.NewClient(secretID, secretKey, region)
}
