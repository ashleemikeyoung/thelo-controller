// Copyright (C) 2026 Qavats Ranch, LLC
// SPDX-License-Identifier: AGPL-3.0-or-later

package templates

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHomeIsBrandedAndExplainsJoining(t *testing.T) {
	html := Home("https://pneuma.example.com").Render()

	assert.Contains(t, html, "<title>Thelo</title>")
	assert.Contains(t, html, "thelo login --login-server https://pneuma.example.com")
	assert.Contains(t, html, "Simlah")
	assert.Contains(t, html, "The Thelo Controller is operated by Qavats Ranch")
	assert.NotContains(t, html, "Pneuma<")
	assert.NotContains(t, html, "ElRoi account")
	// product mark in the header, company mark in the footer, each with a dark-mode variant
	for _, mark := range []string{"elroi", "qavats-it"} {
		assert.Contains(t, html, `src="/brand/`+mark+`-on-light.svg"`)
		assert.Contains(t, html, `srcset="/brand/`+mark+`-on-dark.svg"`)
	}
	// the licence and this version's source, which the AGPL asks for
	assert.Contains(t, html, "AGPL-3.0-or-later")
	assert.Contains(t, html, `href="`+SourceURL+`/tree/`)
	// the upstream credit the BSD-3-Clause licences ask for
	assert.Contains(t, html, "https://github.com/juanfont/headscale")
	assert.Contains(t, html, "https://github.com/tailscale/tailscale")
	assert.NotContains(t, html, "Powered by")
}
