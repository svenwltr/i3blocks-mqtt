package cmd

import (
	"context"

	"github.com/rebuy-de/rebuy-go-sdk/v5/pkg/cmdutil"
	"github.com/spf13/cobra"
)

type RunnerFunc func(ctx context.Context) error

func (r RunnerFunc) Bind(cmd *cobra.Command) error {
	return nil
}

func (r RunnerFunc) Run(ctx context.Context) error {
	return r(ctx)
}

func NewRootCommand() *cobra.Command {
	return cmdutil.New(
		"i3blocks-mqtt", "Tools for adding MQTT support to i3blocks",
		cmdutil.WithLogVerboseFlag(),
		cmdutil.WithVersionCommand(),

		cmdutil.WithSubCommand(cmdutil.New(
			"subscribe", "Subscribes to a single topic and output the data to stdout",
			cmdutil.WithRunner(new(SubscribeRunner)),
		)),
	)
}
