package output

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"

	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"
)

// Writer writes search/context results in json or csv format to stdout or a file.
type Writer struct {
	w       io.Writer
	bw      *bufio.Writer
	csvw    *csv.Writer
	format  string
	close   func() error
}

// NewWriter creates a Writer. If path is empty, writes to stdout.
func NewWriter(format, path string) (*Writer, error) {
	var w io.Writer = os.Stdout
	var close func() error = func() error { return nil }
	if path != "" {
		f, err := os.Create(path)
		if err != nil {
			return nil, err
		}
		w = f
		close = f.Close
	}
	bw := bufio.NewWriter(w)
	csvw := csv.NewWriter(bw)
	return &Writer{w: w, bw: bw, csvw: csvw, format: format, close: close}, nil
}

// Flush flushes buffered output.
func (ow *Writer) Flush() error {
	if ow.format == "csv" && ow.csvw != nil {
		ow.csvw.Flush()
	}
	return ow.bw.Flush()
}

// Close flushes and closes the underlying writer if it is a file.
func (ow *Writer) Close() error {
	_ = ow.Flush()
	return ow.close()
}

// WriteLogInfo writes a single LogInfo in the configured format.
func (ow *Writer) WriteLogInfo(log *cls.LogInfo) error {
	if ow.format == "json" {
		b, err := json.Marshal(log)
		if err != nil {
			return err
		}
		_, err = ow.bw.Write(append(b, '\n'))
		return err
	}
	// csv: Time, TopicId, PkgId, PkgLogId, LogJson
	timeStr := ""
	if log.Time != nil {
		timeStr = fmt.Sprintf("%d", *log.Time)
	}
	topicID := ""
	if log.TopicId != nil {
		topicID = *log.TopicId
	}
	pkgID := ""
	if log.PkgId != nil {
		pkgID = *log.PkgId
	}
	pkgLogID := ""
	if log.PkgLogId != nil {
		pkgLogID = *log.PkgLogId
	}
	logJSON := ""
	if log.LogJson != nil {
		logJSON = *log.LogJson
	}
	return ow.csvw.Write([]string{timeStr, topicID, pkgID, pkgLogID, logJSON})
}

// WriteLogContextInfo writes a single LogContextInfo.
func (ow *Writer) WriteLogContextInfo(log *cls.LogContextInfo) error {
	if ow.format == "json" {
		b, err := json.Marshal(log)
		if err != nil {
			return err
		}
		_, err = ow.bw.Write(append(b, '\n'))
		return err
	}
	btime := ""
	if log.BTime != nil {
		btime = fmt.Sprintf("%d", *log.BTime)
	}
	pkgID := ""
	if log.PkgId != nil {
		pkgID = *log.PkgId
	}
	pkgLogID := ""
	if log.PkgLogId != nil {
		pkgLogID = fmt.Sprintf("%d", *log.PkgLogId)
	}
	content := ""
	if log.Content != nil {
		content = *log.Content
	}
	return ow.csvw.Write([]string{btime, pkgID, pkgLogID, content})
}

// WriteTableHeader writes CSV header for LogInfo.
func (ow *Writer) WriteTableHeaderLogInfo() error {
	if ow.format != "csv" {
		return nil
	}
	return ow.csvw.Write([]string{"Time", "TopicId", "PkgId", "PkgLogId", "LogJson"})
}

// WriteTableHeaderLogContext writes CSV header for LogContextInfo.
func (ow *Writer) WriteTableHeaderLogContext() error {
	if ow.format != "csv" {
		return nil
	}
	return ow.csvw.Write([]string{"BTime", "PkgId", "PkgLogId", "Content"})
}

// WriteAnalysisRecords writes SQL analysis result (columns + records).
func (ow *Writer) WriteAnalysisRecords(columns []*cls.Column, records []*string) error {
	if ow.format == "json" {
		for _, r := range records {
			if r != nil {
				ow.bw.WriteString(*r)
				ow.bw.WriteByte('\n')
			}
		}
		return nil
	}
	// csv: column names then rows (each record as one field)
	var colNames []string
	for _, col := range columns {
		if col != nil && col.Name != nil {
			colNames = append(colNames, *col.Name)
		}
	}
	if len(colNames) == 0 {
		colNames = []string{"result"}
	}
	if err := ow.csvw.Write(colNames); err != nil {
		return err
	}
	for _, r := range records {
		if r != nil {
			if err := ow.csvw.Write([]string{*r}); err != nil {
				return err
			}
		}
	}
	return nil
}
