package main

import (
	"github.com/spf13/cobra"
	"github.com/zx06/xurl/internal/errors"
	"github.com/zx06/xurl/internal/history"
	"github.com/zx06/xurl/internal/output"
)

func NewRunCommand(w output.Writer) *cobra.Command {
	var runLast bool

	cmd := &cobra.Command{
		Use:   "run <target> [request-name]",
		Short: "Execute HTTP request(s) from a .http file (path or library name)",
		Args: func(cmd *cobra.Command, args []string) error {
			last, _ := cmd.Flags().GetBool("last")
			if last {
				return nil
			}
			if len(args) < 1 {
				return errors.New(errors.CodeHTTPFile, "target required (or use --last)", nil)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			format, xe := parseOutputFormat(global.FormatStr)
			if xe != nil {
				return xe
			}

			if runLast {
				store := history.NewStore(global.Resolved.HistoryPath)
				entry, xe := store.Latest()
				if xe != nil {
					return xe
				}
				prof, xe := profileForRun(entry.Profile)
				if xe != nil {
					return xe
				}
				return executeRun(w, format, runParams{
					Target:      entry.Path,
					RequestName: entry.RequestName,
					ProfileName: entry.Profile,
					Profile:     prof,
					RecordHist:  true,
				})
			}

			if global.ProfileStr == "" {
				return errors.New(errors.CodeCfgInvalid, "profile is required; use -p/--profile or XURL_PROFILE", nil)
			}

			reqName := ""
			if len(args) >= 2 {
				reqName = args[1]
			}

			return executeRun(w, format, runParams{
				Target:      args[0],
				RequestName: reqName,
				ProfileName: global.ProfileStr,
				Profile:     global.Resolved.Profile,
				RecordHist:  true,
			})
		},
	}

	cmd.Flags().BoolVar(&runLast, "last", false, "Re-run the most recent history entry (uses stored profile)")
	return cmd
}
