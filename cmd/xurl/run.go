package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/zx06/xurl/internal/errors"
	"github.com/zx06/xurl/internal/httpfile"
	"github.com/zx06/xurl/internal/output"
	"github.com/zx06/xurl/internal/runner"
)

func NewRunCommand(w output.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "run <file.http> [request-name]",
		Short: "Execute HTTP request(s) from a .http file",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			format, xe := parseOutputFormat(global.FormatStr)
			if xe != nil {
				return xe
			}
			if global.ProfileStr == "" {
				return errors.New(errors.CodeCfgInvalid, "profile is required; use -p/--profile or XURL_PROFILE", nil)
			}

			path := args[0]
			reqName := ""
			if len(args) == 2 {
				reqName = args[1]
			}

			b, err := os.ReadFile(path)
			if err != nil {
				return errors.Wrap(errors.CodeHTTPFile, "failed to read HTTP file", map[string]any{"path": path}, err)
			}

			requests, xe := httpfile.Parse(b)
			if xe != nil {
				return xe
			}

			req, ok := httpfile.FindByName(requests, reqName)
			if !ok {
				if reqName == "" {
					return errors.New(errors.CodeHTTPFile, "file contains multiple requests; specify a request name", map[string]any{"count": len(requests)})
				}
				return errors.New(errors.CodeHTTPFile, "request not found", map[string]any{"name": reqName})
			}

			result := runner.Run(req, runner.Options{
				Profile: global.Resolved.Profile,
				Attrs:   global.Attrs,
			})

			return w.WriteRunResult(format, result)
		},
	}
}
