package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/wgyuuu/reqbot/app/controller"
	"github.com/wgyuuu/reqbot/common/dlog"
)

type Options struct {
	isDebug    bool
	configFile string
	flags      *pflag.FlagSet
}

func newCommand() *cobra.Command {
	opts := Options{}

	cmd := &cobra.Command{
		Use:           "reqbot [OPTIONS]",
		Short:         "robot for send request",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.flags = cmd.Flags()
			return runReqbot(opts)
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&opts.isDebug, "debug", "d", false, "look up send info")
	flags.StringVarP(&opts.configFile, "config", "f", "config.dt", "Info file for request message")

	return cmd
}

func main() {
	cmd := newCommand()
	if err := cmd.Execute(); err != nil {
		dlog.Errorf("%s\n", err.Error())
		os.Exit(1)
	}
}

func runReqbot(opts Options) error {
	dlog.Init(opts.isDebug)
	return controller.ProcCNF(opts.configFile)
}
