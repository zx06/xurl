package main

import (
	"github.com/spf13/cobra"
	"github.com/zx06/xurl/internal/library"
	"github.com/zx06/xurl/internal/output"
)

func NewLibCommand(w output.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lib",
		Short: "Manage the default .http file library",
	}
	cmd.AddCommand(newLibPath(w))
	cmd.AddCommand(newLibList(w))
	cmd.AddCommand(newLibAdd(w))
	return cmd
}

func newLibPath(w output.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "path",
		Short: "Print the library directory path",
		RunE: func(cmd *cobra.Command, args []string) error {
			format, xe := parseOutputFormat(global.FormatStr)
			if xe != nil {
				return xe
			}
			dir := global.Resolved.LibraryDir
			if xe := library.EnsureDir(dir); xe != nil {
				return xe
			}
			return w.WriteOK(format, map[string]any{"library_dir": dir})
		},
	}
}

func newLibList(w output.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List .http files in the library",
		RunE: func(cmd *cobra.Command, args []string) error {
			format, xe := parseOutputFormat(global.FormatStr)
			if xe != nil {
				return xe
			}
			dir := global.Resolved.LibraryDir
			files, xe := library.List(dir)
			if xe != nil {
				return xe
			}
			return w.WriteOK(format, map[string]any{"library_dir": dir, "files": files})
		},
	}
}

func newLibAdd(w output.Writer) *cobra.Command {
	var symlink bool
	cmd := &cobra.Command{
		Use:   "add <file.http>",
		Short: "Copy (or symlink) a .http file into the library",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			format, xe := parseOutputFormat(global.FormatStr)
			if xe != nil {
				return xe
			}
			dest, xe := library.Add(args[0], global.Resolved.LibraryDir, symlink)
			if xe != nil {
				return xe
			}
			return w.WriteOK(format, map[string]any{"library_dir": global.Resolved.LibraryDir, "path": dest})
		},
	}
	cmd.Flags().BoolVar(&symlink, "symlink", false, "Create a symlink instead of copying")
	return cmd
}
