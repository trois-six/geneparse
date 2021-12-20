package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Trois-Six/geneparse/pkg/geneanet"
	"github.com/spf13/cobra"
)

const (
	gedcomTimeout       = "60s"
	errInputDirNotExits = "input directory does not exist: %w"
	errFailedParse      = "failed to parse: %w"
)

var errInputDirectoryRequired = errors.New("input directory MUST be a directory")

type GedcomCmd struct{}

func (c *GedcomCmd) Command() *cobra.Command {
	var (
		inputDir string
		timeout  string
	)

	cmd := &cobra.Command{
		Use:   "gedcom",
		Short: "parse Geneanet bases and create a gedcom file",
		Long: `The gedcom command will parse Geneanet bases downloaded by the dlextr command ` +
			`and will create the corresponding gedcom file.`,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			i, err := cmd.Flags().GetString("inputdir")
			if err != nil {
				return fmt.Errorf(errParseInput, err)
			}

			ts, err := cmd.Flags().GetString("timeout")
			if err != nil {
				return fmt.Errorf(errParseInput, err)
			}

			t, err := time.ParseDuration(ts)
			if err != nil {
				return fmt.Errorf(errParseTimeout, err)
			}

			return c.Run(i, t)
		},
	}

	cmd.Flags().StringVarP(&inputDir, "inputdir", "i", "output", "Input directory for Geneanet bases")
	cmd.Flags().StringVarP(&timeout, "timeout", "t", gedcomTimeout, "Timeout to process bases")

	if err := cmd.MarkFlagRequired("inputdir"); err != nil {
		return nil
	}

	return cmd
}

func (c *GedcomCmd) Run(inputDir string, timeout time.Duration) error {
	info, err := os.Stat(inputDir)
	if err != nil {
		return fmt.Errorf(errInputDirNotExits, err)
	} else if !info.IsDir() {
		return errInputDirectoryRequired
	}

	ctx := context.Background()
	g := geneanet.New("", "", inputDir, timeout)

	if err := g.ParseInfoBase(ctx); err != nil {
		return fmt.Errorf(errFailedParse, err)
	}

	log.Printf("NbPersons: %d, Sosa: %d, Date: %s", g.GetNbPersons(), g.GetSosa(), time.Unix(g.GetTimestamp(), 0))

	return nil
}
