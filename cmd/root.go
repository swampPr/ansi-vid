// Package cmd provides cmd  ->  The root command of "ansi-vid"
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/swampPr/ansi-vid/internal/render"
)

var rootCmd = &cobra.Command{
	Use:   "ansi-vid",
	Short: "Render your video as ANSI art.",
	RunE:  execAnsi,
}

func execAnsi(_ *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("please input a path to the source video")
	}

	absPath, err := filepath.Abs(args[0])
	if err != nil {
		return fmt.Errorf("an error occurred: %w", err)
	}

	err = render.Process(absPath)
	if err != nil {
		return fmt.Errorf("an error occurred: %w", err)
	}

	return nil
}

// Execute function  ->  Executes the command
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
