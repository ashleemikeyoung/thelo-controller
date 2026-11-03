package cli

import (
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// The names people see for this control server. Headscale is the code underneath;
// it appears only in the credit, never as the product (BRANDING.md in
// identity_authority_stack).
const (
	brandProduct = "Thelo Controller"
	brandCommand = "theloctl"
	brandNet     = "thelonet"
	brandCredit  = "Built on Headscale (https://github.com/juanfont/headscale) and " +
		"Tailscale (https://github.com/tailscale/tailscale), BSD-3-Clause."
	// checkUpstreamReleases is off: theloctl does not ask GitHub about Headscale
	// releases. Updates come from our own builds.
	checkUpstreamReleases = false
)

var (
	// "headscale nodes list", `headscale`, 'headscale ...': the command, not the product.
	commandMention = regexp.MustCompile("(^|[\\s`'\"(])headscale( |`|'|\"|$)")
	// "/etc/headscale/config.yaml" and the other default paths.
	pathMention = strings.NewReplacer(
		"/etc/headscale", "/etc/thelo-controller",
		"/var/lib/headscale", "/var/lib/thelo-controller",
		"/var/run/headscale/headscale.sock", "/var/run/thelo-controller/theloctl.sock",
	)
	// "of Headscale" reads "of the Thelo Controller"; "the headscale server" is the product too.
	productAfterPreposition = regexp.MustCompile(`\b(of|in|for|with|to) [Hh]eadscale( server| instance)?\b`)
	productNoun             = regexp.MustCompile(`\b[Hh]eadscale (server|instance)\b`)
	productMention          = regexp.MustCompile(`\b[Hh]eadscale\b`)
	netMention              = regexp.MustCompile(`\b([Tt])ailnet(s?)\b`)
)

// brandText rewrites text people read: the command becomes theloctl, the product
// the Thelo Controller, a tailnet a thelonet. The upstream credit is left as it is.
func brandText(s string) string {
	if before, after, found := strings.Cut(s, brandCredit); found {
		return brandText(before) + brandCredit + brandText(after)
	}

	s = pathMention.Replace(s)
	s = productAfterPreposition.ReplaceAllString(s, "$1 the "+brandProduct)
	s = productNoun.ReplaceAllString(s, brandProduct)
	s = commandMention.ReplaceAllString(s, "${1}"+brandCommand+"${2}")
	s = productMention.ReplaceAllString(s, brandProduct)
	s = netMention.ReplaceAllStringFunc(s, func(m string) string {
		out := brandNet + strings.TrimPrefix(strings.TrimPrefix(m, "tailnet"), "Tailnet")
		if m[0] == 'T' {
			return "T" + out[1:]
		}

		return out
	})

	return s
}

// brandCommands applies brandText to every command's help and every flag's usage.
func brandCommands(c *cobra.Command) {
	c.Short = brandText(c.Short)
	c.Long = brandText(c.Long)
	c.Example = brandText(c.Example)
	c.Deprecated = brandText(c.Deprecated)

	for _, fs := range []*pflag.FlagSet{c.Flags(), c.PersistentFlags()} {
		fs.VisitAll(func(f *pflag.Flag) { f.Usage = brandText(f.Usage) })
	}

	for _, sub := range c.Commands() {
		brandCommands(sub)
	}
}

func init() {
	rootCmd.Use = brandCommand
	rootCmd.Short = brandCommand + " - the " + brandProduct
	rootCmd.Long = "\n" + brandCommand + " runs and manages the " + brandProduct +
		", the control server of a " + brandNet + ".\n\n" + brandCredit
}
