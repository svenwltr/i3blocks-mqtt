package cmd

import (
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func AutoEnv(cmd *cobra.Command, args []string) {
	cmd.PersistentFlags().VisitAll(func(flag *pflag.Flag) {
		if !flag.Changed {
			var (
				envName  = "BLOCKS_" + ConvertToValidEnvVarName(flag.Name)
				envValue = os.Getenv(envName)
			)

			if envValue == "" {
				return
			}

			flag.Value.Set(envValue)
		}
	})
}

func ConvertToValidEnvVarName(s string) string {
	// Remove non-alphanumeric characters and replace them with underscores
	re := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	s = re.ReplaceAllString(s, "_")

	// Remove leading and trailing underscores
	s = strings.Trim(s, "_")

	// Convert to uppercase
	s = strings.ToUpper(s)

	return s
}
