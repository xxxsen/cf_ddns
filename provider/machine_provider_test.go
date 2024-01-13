package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMachineGet(t *testing.T) {
	{
		p, err := createMachineProvider(map[string]interface{}{
			"interface":     "br0",
			"type":          "ipv4",
			"allow_private": true,
		})
		assert.NoError(t, err)
		addr, err := p.Get()
		assert.NoError(t, err)
		t.Logf("v4:%s", addr)
	}
	{
		p, err := createMachineProvider(map[string]interface{}{
			"interface": "br0",
			"type":      "ipv6",
		})
		assert.NoError(t, err)
		addr, err := p.Get()
		assert.NoError(t, err)
		t.Logf("v6:%s", addr)
	}
}
