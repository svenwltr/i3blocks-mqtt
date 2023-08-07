package main

import (
	"github.com/rebuy-de/rebuy-go-sdk/v5/pkg/cmdutil"
	"github.com/sirupsen/logrus"
	"github.com/svenwltr/i3block-mqtt/cmd"
)

func main() {
	defer cmdutil.HandleExit()
	if err := cmd.NewRootCommand().Execute(); err != nil {
		logrus.Fatal(err)
	}
}
