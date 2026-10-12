package output

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/zx06/xurl/internal/errors"
)

type Format string

const (
	FormatJSON  Format = "json"
	FormatTable Format = "table"
)

func ParseFormat(s string) (Format, *errors.XError) {
	switch s {
	case "json", "":
		return FormatJSON, nil
	case "table":
		return FormatTable, nil
	default:
		return "", errors.New(errors.CodeCfgInvalid, "invalid format; use json or table", map[string]any{"format": s})
	}
}

type Writer struct {
	Out io.Writer
	Err io.Writer
}

func New(out, err io.Writer) Writer {
	return Writer{Out: out, Err: err}
}

func (w Writer) WriteRunResult(format Format, data any) error {
	switch format {
	case FormatJSON:
		enc := json.NewEncoder(w.Out)
		enc.SetEscapeHTML(false)
		enc.SetIndent("", "  ")
		return enc.Encode(data)
	case FormatTable:
		return w.writeTable(data)
	default:
		return fmt.Errorf("unsupported format")
	}
}

func (w Writer) WriteOK(format Format, data any) error {
	env := map[string]any{"ok": true, "data": data}
	switch format {
	case FormatJSON:
		enc := json.NewEncoder(w.Out)
		enc.SetEscapeHTML(false)
		enc.SetIndent("", "  ")
		return enc.Encode(env)
	case FormatTable:
		b, _ := json.MarshalIndent(data, "", "  ")
		_, err := w.Out.Write(b)
		return err
	default:
		return fmt.Errorf("unsupported format")
	}
}

func (w Writer) WriteError(format Format, xe *errors.XError) error {
	env := map[string]any{
		"ok": false,
		"error": map[string]any{
			"code":    xe.Code,
			"message": xe.Message,
			"details": xe.Details,
		},
	}
	switch format {
	case FormatJSON:
		enc := json.NewEncoder(w.Out)
		enc.SetEscapeHTML(false)
		enc.SetIndent("", "  ")
		return enc.Encode(env)
	case FormatTable:
		_, err := fmt.Fprintf(w.Out, "ERROR %s: %s\n", xe.Code, xe.Message)
		return err
	default:
		return fmt.Errorf("unsupported format")
	}
}

func (w Writer) writeTable(data any) error {
	// Minimal table view for run results
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		_, _ = w.Out.Write(b)
		return nil
	}
	tw := tabwriter.NewWriter(w.Out, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "OK\t%v\n", m["ok"])
	if req, ok := m["request"].(map[string]any); ok {
		fmt.Fprintf(tw, "METHOD\t%v\n", req["method"])
		fmt.Fprintf(tw, "URL\t%v\n", req["url"])
	}
	if st, ok := m["status"]; ok {
		fmt.Fprintf(tw, "STATUS\t%v\n", st)
	}
	if body, ok := m["body"].(string); ok {
		fmt.Fprintf(tw, "BODY\t%s\n", body)
	}
	return tw.Flush()
}
