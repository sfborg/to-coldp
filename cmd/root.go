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

	"github.com/sfborg/sflib/io/archio"
	"github.com/sfborg/sflib/io/dbio"
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

		// initiate sfga file
		sfga, err := archio.New(sfgaPath, cfg.CacheDir)
		if err != nil {
			slog.Error("Cannot initialize SFGA archive", "error", err)
			os.Exit(1)
		}

		err = sfga.Extract()
		if err != nil {
			slog.Error("Cannot extract SFGA archive", "error", err)
			os.Exit(1)
		}

		// initiate sfga db instance
		sfdb := dbio.New(cfg.CacheDbDir)
		db, err := sfdb.Connect()
		if err != nil {
			slog.Error("Cannot connect to SFGA database", "error", err)
			os.Exit(1)
		}

		isComp := sfdb.IsCompatible(cfg.MinVersionSFGA)
		if !isComp {
			ver := sfdb.Version()
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

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.to-coldp.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
