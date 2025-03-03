/*
Copyright © 2025 Dmitry Mozzherin <dmozzherin@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"log/slog"
	"os"

	"github.com/gnames/gnsys"
	"github.com/sfborg/sflib/io/sfgaio"
	"github.com/sfborg/to-coldp/internal/io/cfio"
	tocoldp "github.com/sfborg/to-coldp/pkg"
	"github.com/sfborg/to-coldp/pkg/config"
	"github.com/spf13/cobra"
)

var (
	opts []config.Option
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "to-coldp",
	Short: "Converts SFGA archive to CoLDP",
	Long: `Takes a path to SFGArchive and outputs CoLDP file into provided
second path. SFGArchive can be compressed (Zip, Tar), be in binary or
SQL dump format.
`,
	Run: func(cmd *cobra.Command, args []string) {
		versionFlag(cmd)
		flags := []flagFunc{
			nameUsageFlag,
		}
		for _, v := range flags {
			v(cmd)
		}

		cfg := config.New(opts...)

		if len(args) != 2 {
			cmd.Help()
			os.Exit(0)
		}

		sfgaPath := args[0]
		coldpPath := args[1]

		slog.Info("Extracting SFGA data", "path", sfgaPath)

		err := prepareFileStructure()
		if err != nil {
			slog.Error("Cannot setup cache file structure", "error", err)
			os.Exit(1)
		}

		// initiate sfga file
		sfga := sfgaio.New()

		err = sfga.Import(sfgaPath, cfg.CacheSfgaDir)
		if err != nil {
			slog.Error("Cannot extract SFGA archive", "error", err)
			os.Exit(1)
		}

		db, err := sfga.Connect()
		if err != nil {
			slog.Error("Cannot connect to SFGA database", "error", err)
			os.Exit(1)
		}

		isComp := sfga.IsCompatible(cfg.MinVersionSFGA)
		if !isComp {
			ver := sfga.Version()
			slog.Error("Incompatible SFGA version",
				"want", cfg.MinVersionSFGA, "got", ver,
			)
			os.Exit(1)
		}

		clf, err := cfio.New(db, cfg.CacheColdpDir)
		if err != nil {
			slog.Error("Cannot set CoLDP", "error", err)
			os.Exit(1)
		}

		tcdp := tocoldp.New(cfg, clf)

		slog.Info("Exporting SFGA data to CoLDP")
		err = tcdp.Export(coldpPath)
		if err != nil {
			slog.Error("Cannot export SFGA to CoLDP", "error", err)
			os.Exit(1)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func prepareFileStructure() error {
	var err error
	cfg := config.New()
	root := cfg.CacheDir
	err = gnsys.MakeDir(root)
	if err != nil {
		return err
	}
	err = gnsys.CleanDir(root)
	if err != nil {
		return err
	}
	dirs := []string{
		cfg.CacheColdpDir,
		cfg.CacheSfgaDir,
	}
	for _, v := range dirs {
		err = gnsys.MakeDir(v)
		if err != nil {
			return err
		}
	}

	return nil
}

func init() {
}
