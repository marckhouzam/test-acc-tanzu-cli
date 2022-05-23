/*
Copyright 2021 VMware, Inc. All Rights Reserved.
*/

package commands

import (
	"context"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	cli "github.com/vmware-tanzu/apps-cli-plugin/pkg/cli-runtime"
)

type DocsOptions struct {
	Directory string
}

func DocsCommand(ctx context.Context, c *cli.Config) *cobra.Command {
	opts := &DocsOptions{}

	cmd := &cobra.Command{
		Use:     "docs",
		Short:   "generate docs in Markdown for this CLI",
		Example: fmt.Sprintf("%s docs", c.Name),
		Hidden:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := os.MkdirAll(opts.Directory, 0744); err != nil {
				return err
			}

			if noColorFlag := cmd.Root().Flag(cli.StripDash(cli.NoColorFlagName)); noColorFlag != nil {
				// force default to false for doc generation no matter the environment
				noColorFlag.DefValue = "false"
			}

			root := &cobra.Command{
				Use:               "tanzu",
				DisableAutoGenTag: true,
			}
			root.AddCommand(cmd.Root())

			// hack to rewrite the CommandPath content to add args
			cli.Visit(root, func(cmd *cobra.Command) error {
				if !cmd.HasSubCommands() {
					cmd.Use = cmd.Use + cli.FormatArgs(cmd)
				}
				return nil
			})

			if err := doc.GenMarkdownTree(root, opts.Directory); err != nil {
				return err
			}

			// remove synthetic root command
			if err := os.Remove(path.Join(opts.Directory, "tanzu.md")); err != nil {
				return err
			}
			return filepath.Walk(opts.Directory, func(path string, info fs.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if info.IsDir() {
					return nil
				}
				input, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}
				inlines := strings.Split(string(input), "\n")
				outlines := []string{}
				for _, line := range inlines {
					if !strings.HasPrefix(line, "* [tanzu](tanzu.md)") {
						outlines = append(outlines, line)
					}
				}
				return ioutil.WriteFile(path, []byte(strings.Join(outlines, "\n")), 0644)
			})
		},
	}

	cmd.Flags().StringVarP(&opts.Directory, "directory", "d", "docs", "the output `directory` for the docs")

	return cmd
}
