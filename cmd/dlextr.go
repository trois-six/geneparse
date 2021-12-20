package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Trois-Six/geneparse/pkg/geneanet"
	"github.com/spf13/cobra"
)

const (
	loginTimeout             = "10s"
	errParseInput            = "could not parse input: %w"
	errParseTimeout          = "could not parse timeout: %w"
	errAccessOutputDir       = "could not access filesystem: %w"
	errCreateOutputDir       = "could not create output directory: %w"
	errFailedLogin           = "failed to log in: %w"
	errFailedGetAccountInfos = "failed get account infos: %w"
	errFailedSetLogged       = "failed set as logged: %w"
	errFailedDownload        = "failed to download the Geneanet bases: %w"
)

var errOutputDirectoryRequired = errors.New("output directory MUST be a directory")

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
				return fmt.Errorf(errParseInput, err)
			}

			p, err := cmd.Flags().GetString("password")
			if err != nil {
				return fmt.Errorf(errParseInput, err)
			}

			o, err := cmd.Flags().GetString("outputdir")
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
			return fmt.Errorf(errAccessOutputDir, err)
		}

		if err = os.MkdirAll(outputDir, os.ModePerm|os.ModeDir); err != nil {
			return fmt.Errorf(errCreateOutputDir, err)
		}
	} else if !info.IsDir() {
		return errOutputDirectoryRequired
	}

	ctx := context.Background()
	g := geneanet.New(username, password, outputDir, timeout)

	if err := g.Login(ctx); err != nil {
		return fmt.Errorf(errFailedLogin, err)
	}

	if err := g.GetAccountInfos(ctx); err != nil {
		return fmt.Errorf(errFailedGetAccountInfos, err)
	}

	if err := g.SetLogged(ctx); err != nil {
		return fmt.Errorf(errFailedSetLogged, err)
	}

	if err := g.GetBase(ctx); err != nil {
		return fmt.Errorf(errFailedDownload, err)
	}

	return nil
}
