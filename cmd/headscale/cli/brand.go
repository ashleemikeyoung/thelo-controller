// Copyright (C) 2026 Qavats Ranch, LLC
// SPDX-License-Identifier: AGPL-3.0-or-later

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
	brandCredit  = "Free software, AGPL-3.0-or-later (source: https://github.com/ashleemikeyoung/headscale). " +
		"Built on Headscale (https://github.com/juanfont/headscale) and " +
		"Tailscale (https://github.com/tailscale/tailscale), BSD-3-Clause."
	// checkUpstreamReleases is off: theloctl does not ask GitHub about Headscale
	// releases. Updates come from our own builds.
	checkUpstreamReleases = false
)

var (
	// "headscale nodes list", `headscale`, 'headscale ...': the command, not the product.
	commandMention = regexp.MustCompile("(^|[\\s`'\"(])headscale( |`|'|\"|$)")
	// "/etc/headscale/config.yaml" and the other default paths, in help text only.
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
// the Thelo Controller, a tailnet a thelonet. The upstream credit is left as it is,
// and so is any name inside a path, URL or identifier (/etc/headscale/,
// juanfont/headscale, tailnet_id), so error messages still point at real things.
func brandText(s string) string {
	if before, after, found := strings.Cut(s, brandCredit); found {
		return brandText(before) + brandCredit + brandText(after)
	}

	s = replaceStandalone(s, productAfterPreposition, func(m []string) string { return m[1] + " the " + brandProduct })
	s = replaceStandalone(s, productNoun, func([]string) string { return brandProduct })
	s = commandMention.ReplaceAllString(s, "${1}"+brandCommand+"${2}")
	s = replaceStandalone(s, productMention, func([]string) string { return brandProduct })
	s = replaceStandalone(s, netMention, func(m []string) string {
		if m[1] == "T" {
			return "T" + brandNet[1:] + m[2]
		}

		return brandNet + m[2]
	})

	return s
}

// brandHelp is brandText for help and flag usage, where the default paths are
// also ours.
func brandHelp(s string) string {
	return brandText(pathMention.Replace(s))
}

// replaceStandalone replaces the matches of re that stand on their own: not
// joined to a letter, digit or one of "/._-" on either side. A full stop ends a
// sentence unless a letter or digit follows it (headscale.sock).
func replaceStandalone(s string, re *regexp.Regexp, repl func(submatches []string) string) string {
	alnum := func(c byte) bool { return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' }
	joinedBefore := func(i int) bool { return i >= 0 && (alnum(s[i]) || strings.IndexByte("/._-", s[i]) >= 0) }
	joinedAfter := func(i int) bool {
		if i >= len(s) {
			return false
		}
		if s[i] == '.' {
			return i+1 < len(s) && alnum(s[i+1])
		}

		return alnum(s[i]) || strings.IndexByte("/_-", s[i]) >= 0
	}

	var b strings.Builder

	last := 0
	for _, loc := range re.FindAllStringSubmatchIndex(s, -1) {
		if joinedBefore(loc[0]-1) || joinedAfter(loc[1]) {
			continue
		}

		sub := make([]string, len(loc)/2)
		for i := range sub {
			if loc[2*i] >= 0 {
				sub[i] = s[loc[2*i]:loc[2*i+1]]
			}
		}

		b.WriteString(s[last:loc[0]])
		b.WriteString(repl(sub))
		last = loc[1]
	}

	b.WriteString(s[last:])

	return b.String()
}

// brandCommands applies brandHelp to every command's help and every flag's usage.
func brandCommands(c *cobra.Command) {
	c.Short = brandHelp(c.Short)
	c.Long = brandHelp(c.Long)
	c.Example = brandHelp(c.Example)
	c.Deprecated = brandHelp(c.Deprecated)

	for _, fs := range []*pflag.FlagSet{c.Flags(), c.PersistentFlags()} {
		fs.VisitAll(func(f *pflag.Flag) { f.Usage = brandHelp(f.Usage) })
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
