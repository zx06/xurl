package main

import (
	"strconv"

	"github.com/spf13/cobra"
	"github.com/zx06/xurl/internal/errors"
	"github.com/zx06/xurl/internal/history"
	"github.com/zx06/xurl/internal/output"
)

func NewHistoryCommand(w output.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "history",
		Short: "View and re-run past requests",
	}
	cmd.AddCommand(newHistoryList(w))
	cmd.AddCommand(newHistoryRun(w))
	return cmd
}

func newHistoryList(w output.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List recent runs (newest first; index 1 = most recent)",
		RunE: func(cmd *cobra.Command, args []string) error {
			format, xe := parseOutputFormat(global.FormatStr)
			if xe != nil {
				return xe
			}
			store := history.NewStore(global.Resolved.HistoryPath)
			entries, xe := store.ListNewestFirst()
			if xe != nil {
				return xe
			}
			type row struct {
				Index       int    `json:"index"`
				ID          string `json:"id"`
				At          string `json:"at"`
				Path        string `json:"path"`
				RequestName string `json:"request_name,omitempty"`
				Profile     string `json:"profile"`
				OK          bool   `json:"ok"`
			}
			rows := make([]row, 0, len(entries))
			for i, e := range entries {
				rows = append(rows, row{
					Index:       i + 1,
					ID:          e.ID,
					At:          e.At.Format("2006-01-02T15:04:05Z07:00"),
					Path:        e.Path,
					RequestName: e.RequestName,
					Profile:     e.Profile,
					OK:          e.OK,
				})
			}
			return w.WriteOK(format, map[string]any{"entries": rows})
		},
	}
}

func newHistoryRun(w output.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "run <index>",
		Short: "Re-run a history entry by index (from 'history list')",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			format, xe := parseOutputFormat(global.FormatStr)
			if xe != nil {
				return xe
			}
			idx, err := strconv.Atoi(args[0])
			if err != nil || idx < 1 {
				return errors.New(errors.CodeCfgInvalid, "invalid history index", map[string]any{"index": args[0]})
			}
			store := history.NewStore(global.Resolved.HistoryPath)
			entry, xe := store.GetByIndex(idx)
			if xe != nil {
				return xe
			}
			profileName := entry.Profile
			prof, xe := profileForRun(entry.Profile)
			if xe != nil {
				return xe
			}
			if profileFlagChanged(cmd) && global.ProfileStr != "" {
				profileName = global.ProfileStr
				prof = global.Resolved.Profile
			}
			return executeRun(w, format, runParams{
				Target:      entry.Path,
				RequestName: entry.RequestName,
				ProfileName: profileName,
				Profile:     prof,
				RecordHist:  true,
			})
		},
	}
}

func profileFlagChanged(cmd *cobra.Command) bool {
	f := cmd.Root().PersistentFlags().Lookup("profile")
	return f != nil && f.Changed
}
