package main

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zx06/xurl/internal/config"
	"github.com/zx06/xurl/internal/errors"
	"github.com/zx06/xurl/internal/output"
)

var (
	version = "dev"
)

type GlobalConfig struct {
	ConfigStr  string
	ProfileStr string
	FormatStr  string
	Resolved   config.Resolved
	Attrs      map[string]string
}

var global = &GlobalConfig{}

func NewRootCommand() *cobra.Command {
	var cliAttrs []string

	root := &cobra.Command{
		Use:           "xurl",
		Short:         "AI-first HTTP client CLI",
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       version,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			formatSet := cmd.Flags().Changed("format")
			profileSet := cmd.Flags().Changed("profile")
			configSet := cmd.Flags().Changed("config")
			if configSet && global.ConfigStr == "" {
				return errors.New(errors.CodeCfgInvalid, "config path is empty", nil)
			}

			r, xe := config.Resolve(config.Options{
				ConfigPath:    global.ConfigStr,
				CLIProfile:    global.ProfileStr,
				CLIProfileSet: profileSet,
				CLIFormat:     global.FormatStr,
				CLIFormatSet:  formatSet,
				EnvProfile:    os.Getenv("XURL_PROFILE"),
				EnvFormat:     os.Getenv("XURL_FORMAT"),
			})
			if xe != nil {
				return xe
			}
			global.Resolved = r
			global.FormatStr = r.Format
			global.ProfileStr = r.ProfileName
			global.Attrs = parseAttrs(cliAttrs)
			return nil
		},
	}

	root.PersistentFlags().StringVar(&global.ConfigStr, "config", "", "Config file path (default: ./xurl.yaml or ~/.config/xurl/xurl.yaml)")
	root.PersistentFlags().StringVarP(&global.ProfileStr, "profile", "p", "", "Profile name")
	root.PersistentFlags().StringVarP(&global.FormatStr, "format", "f", "json", "Output format: json|table")
	root.PersistentFlags().StringArrayVar(&cliAttrs, "attr", nil, "Variable key=value (repeatable)")

	out := output.New(os.Stdout, os.Stderr)
	root.AddCommand(NewRunCommand(out))
	root.AddCommand(NewProfileCommand(out))
	root.AddCommand(NewLibCommand(out))
	root.AddCommand(NewHistoryCommand(out))

	return root
}

func parseAttrs(pairs []string) map[string]string {
	out := map[string]string{}
	for _, p := range pairs {
		i := strings.Index(p, "=")
		if i <= 0 {
			continue
		}
		out[strings.TrimSpace(p[:i])] = strings.TrimSpace(p[i+1:])
	}
	return out
}

func parseOutputFormat(s string) (output.Format, *errors.XError) {
	return output.ParseFormat(s)
}
