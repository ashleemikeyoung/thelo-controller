// Copyright (C) 2026 Qavats Ranch, LLC
// SPDX-License-Identifier: AGPL-3.0-or-later

package apiv2

import (
	"encoding/json"
	"net/netip"
	"testing"

	"github.com/juanfont/headscale/hscontrol/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDeviceClientConnectivity asserts the node's reported endpoints appear as
// clientConnectivity.endpoints for fields=all, as [] when there are none, and
// not at all in the default response, matching Tailscale.
func TestDeviceClientConnectivity(t *testing.T) {
	node := &types.Node{
		ID:       7,
		Hostname: "thevault",
		Endpoints: []netip.AddrPort{
			netip.MustParseAddrPort("203.0.113.9:41641"),
			netip.MustParseAddrPort("192.168.4.91:41641"),
		},
	}

	all := deviceFromView(node.View(), true, "")
	require.NotNil(t, all.ClientConnectivity)
	assert.Equal(t, []string{"203.0.113.9:41641", "192.168.4.91:41641"}, all.ClientConnectivity.Endpoints)

	def := deviceFromView(node.View(), false, "")
	assert.Nil(t, def.ClientConnectivity)

	body, err := json.Marshal(def)
	require.NoError(t, err)
	assert.NotContains(t, string(body), "clientConnectivity")

	none := deviceFromView((&types.Node{ID: 8}).View(), true, "")
	body, err = json.Marshal(none)
	require.NoError(t, err)
	assert.Contains(t, string(body), `"clientConnectivity":{"endpoints":[]}`)
}
