package cli

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func TestBrandText(t *testing.T) {
	for in, want := range map[string]string{
		"Manage the nodes of Headscale":                                        "Manage the nodes of the Thelo Controller",
		"Launches the headscale server":                                        "Launches the Thelo Controller",
		"use 'headscale auth register --auth-id <id>'":                         "use 'theloctl auth register --auth-id <id>'",
		"connecting to headscale: dial unix /var/run/headscale/headscale.sock": "connecting to the Thelo Controller: dial unix /var/run/headscale/headscale.sock",
		"loading policy: open /etc/headscale/policy.hujson: no such file":      "loading policy: open /etc/headscale/policy.hujson: no such file",
		"see https://github.com/juanfont/headscale and the tailnet_id field":   "see https://github.com/juanfont/headscale and the tailnet_id field",
		"Add a device to your tailnet; Tailnets are listed":                    "Add a device to your thelonet; Thelonets are listed",
		"The version of headscale.":                                            "The version of the Thelo Controller.",
		"Manage the Headscale ACL Policy":                                      "Manage the Thelo Controller ACL Policy",
		brandCredit:                                                            brandCredit,
	} {
		assert.Equal(t, want, brandText(in), in)
	}
}

func TestBrandHelpUsesOurPaths(t *testing.T) {
	assert.Equal(t, "config file (default is /etc/thelo-controller/config.yaml)",
		brandHelp("config file (default is /etc/headscale/config.yaml)"))
}

// Nothing in theloctl's help names Headscale or a tailnet, except the credit.
func TestHelpIsBranded(t *testing.T) {
	brandCommands(rootCmd)
	assert.Equal(t, "theloctl", rootCmd.Name())
	assert.Contains(t, rootCmd.Long, brandCredit)

	var walk func(c *cobra.Command)
	walk = func(c *cobra.Command) {
		texts := []string{strings.ReplaceAll(c.Long, brandCredit, ""), c.Short, c.Example, c.Deprecated}
		for _, fs := range []*pflag.FlagSet{c.Flags(), c.PersistentFlags()} {
			fs.VisitAll(func(f *pflag.Flag) { texts = append(texts, f.Usage) })
		}
		for _, s := range texts {
			assert.NotRegexp(t, `(?i)headscale|tailnet`, s, c.CommandPath())
		}
		for _, sub := range c.Commands() {
			walk(sub)
		}
	}
	walk(rootCmd)
}
