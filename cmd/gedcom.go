package cmd

import (
	"fmt"
	"os"

	"github.com/Trois-Six/geneparse/pkg/geneanet"
	"github.com/Trois-Six/geneparse/pkg/geneanet/utils"
	"github.com/spf13/cobra"
)

type GedcomCmd struct{}

func (c *GedcomCmd) Command() *cobra.Command {
	var inputDir string

	cmd := &cobra.Command{
		Use:   "gedcom",
		Short: "parse Geneanet bases and create a gedcom file",
		Long: `The gedcom command will parse Geneanet bases downloaded by the dlextr command ` +
			`and will create the corresponding gedcom file.`,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			i, err := cmd.Flags().GetString("inputdir")
			if err != nil {
				return fmt.Errorf("could not parse input: %w", err)
			}

			return c.Run(i)
		},
	}

	cmd.Flags().StringVarP(&inputDir, "inputdir", "i", "output", "Input directory for Geneanet bases")

	if err := cmd.MarkFlagRequired("inputdir"); err != nil {
		return nil
	}

	return cmd
}

func (c *GedcomCmd) Run(inputDir string) error {
	info, err := os.Stat(inputDir)
	if err != nil {
		return fmt.Errorf("input directory does not exist: %w", err)
	} else if !info.IsDir() {
		return fmt.Errorf("%w: %s", utils.ErrDirMustBeADir, inputDir)
	}

	g, err := geneanet.New(inputDir)
	if err != nil {
		return fmt.Errorf("failed to initialize Geneanet: %w", err)
	}

	if err = g.Parse(); err != nil {
		return fmt.Errorf("failed to parse Geneanet: %w", err)
	}

	return nil
}
