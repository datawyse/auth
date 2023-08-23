package main

import (
	"log"

	"auth/cmd"

	"github.com/spf13/cobra"
)

type AuthServiceCli struct {
	RootCmd *cobra.Command
}

func ProjectService() *AuthServiceCli {
	cli := &AuthServiceCli{
		RootCmd: &cobra.Command{
			Use:   "project-service",
			Short: "Project Service CLI",
			FParseErrWhitelist: cobra.FParseErrWhitelist{
				UnknownFlags: true,
			},
			// no need to provide the default cobra completion command
			CompletionOptions: cobra.CompletionOptions{
				DisableDefaultCmd: true,
			},
		},
	}

	// hide the default help command (allow only `--help` flag)
	cli.RootCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	// register system commands
	cli.RootCmd.AddCommand(cmd.Serve())

	return cli
}

func (cli *AuthServiceCli) Start() error {
	if err := cli.RootCmd.Execute(); err != nil {
		return err
	}
	return nil
}

func main() {
	app := ProjectService()
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
