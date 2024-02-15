package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/trois-six/geneparse/pkg/geneanet/dlextr"
	"github.com/trois-six/geneparse/pkg/geneanet/utils"
	"github.com/spf13/cobra"
)

const loginTimeout = "10s"

type DownloadAndExtractCmd struct{}

func (c *DownloadAndExtractCmd) Command() *cobra.Command {
	var (
		username  string
		password  string
		outputDir string
		timeout   string
	)

	cmd := &cobra.Command{
		Use:   "dlextr",
		Short: "download and extract Geneanet bases",
		Long: `The dlextr command will connect to Geneanet as if it was the Geneanet Android app, ` +
			`and will download the Geneanet bases. These bases use the Geneweb format.`,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			u, err := cmd.Flags().GetString("username")
			if err != nil {
				return fmt.Errorf(utils.ErrParseInput, err)
			}

			p, err := cmd.Flags().GetString("password")
			if err != nil {
				return fmt.Errorf(utils.ErrParseInput, err)
			}

			o, err := cmd.Flags().GetString("outputdir")
			if err != nil {
				return fmt.Errorf(utils.ErrParseInput, err)
			}

			ts, err := cmd.Flags().GetString("timeout")
			if err != nil {
				return fmt.Errorf(utils.ErrParseInput, err)
			}

			t, err := time.ParseDuration(ts)
			if err != nil {
				return fmt.Errorf("could not parse timeout: %w", err)
			}

			return c.Run(u, p, o, t)
		},
	}

	cmd.Flags().StringVarP(&username, "username", "u", "", "Username or email address to log in to Geneanet (required)")
	cmd.Flags().StringVarP(&password, "password", "p", "", "Password to log in to Geneanet (required)")
	cmd.Flags().StringVarP(&outputDir, "outputdir", "o", "output", "Output directory for Geneanet bases")
	cmd.Flags().StringVarP(&timeout, "timeout", "t", loginTimeout, "Connection timeout for requests to Geneanet")

	if err := cmd.MarkFlagRequired("username"); err != nil {
		return nil
	}

	if err := cmd.MarkFlagRequired("password"); err != nil {
		return nil
	}

	return cmd
}

func (c *DownloadAndExtractCmd) Run(username, password, outputDir string, timeout time.Duration) error {
	info, err := os.Stat(outputDir)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("could not access filesystem: %w", err)
		}

		if err = os.MkdirAll(outputDir, os.ModePerm|os.ModeDir); err != nil {
			return fmt.Errorf("could not create output directory: %w", err)
		}
	} else if !info.IsDir() {
		return fmt.Errorf("%w: %s", utils.ErrDirMustBeADir, outputDir)
	}

	d := dlextr.New(username, password, outputDir, timeout)

	if err := d.Login(); err != nil {
		return fmt.Errorf("failed to log in: %w", err)
	}

	if err := d.GetAccountInfos(); err != nil {
		return fmt.Errorf("failed get account infos: %w", err)
	}

	if err := d.SetLogged(); err != nil {
		return fmt.Errorf("failed set as logged: %w", err)
	}

	if err := d.GetBase(); err != nil {
		return fmt.Errorf("failed to download the Geneanet bases: %w", err)
	}

	if err := d.Unzip(); err != nil {
		return fmt.Errorf("failed to extract the Geneanet bases: %w", err)
	}

	return nil
}
