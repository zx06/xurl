package main

import (
	"github.com/spf13/cobra"
	"github.com/zx06/xurl/internal/config"
	"github.com/zx06/xurl/internal/errors"
	"github.com/zx06/xurl/internal/output"
)

func NewProfileCommand(w output.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profile",
		Short: "Manage profiles",
	}
	cmd.AddCommand(newProfileList(w))
	cmd.AddCommand(newProfileShow(w))
	return cmd
}

func newProfileList(w output.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List configured profiles",
		RunE: func(cmd *cobra.Command, args []string) error {
			format, xe := parseOutputFormat(global.FormatStr)
			if xe != nil {
				return xe
			}
			profiles := global.Resolved.AllProfiles
			if len(profiles) == 0 {
				cfg, _, xe := config.LoadConfig(config.Options{ConfigPath: global.ConfigStr})
				if xe != nil && xe.Code != errors.CodeCfgNotFound {
					return xe
				}
				profiles = cfg.Profiles
			}
			return w.WriteOK(format, map[string]any{"profiles": config.ListProfiles(profiles)})
		},
	}
}

func newProfileShow(w output.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "show <name>",
		Short: "Show profile details (secrets masked)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			format, xe := parseOutputFormat(global.FormatStr)
			if xe != nil {
				return xe
			}
			profiles := global.Resolved.AllProfiles
			p, ok := profiles[name]
			if !ok {
				cfg, _, xe := config.LoadConfig(config.Options{ConfigPath: global.ConfigStr})
				if xe != nil && xe.Code != errors.CodeCfgNotFound {
					return xe
				}
				p, ok = cfg.Profiles[name]
				if !ok {
					return errors.New(errors.CodeCfgInvalid, "profile not found", map[string]any{"profile": name})
				}
			}
			return w.WriteOK(format, config.MaskProfile(name, p))
		},
	}
}
