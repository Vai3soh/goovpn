package memory

import (
	"testing"

	"github.com/Vai3soh/goovpn/entity"
	"github.com/stretchr/testify/require"
)

var test_config = `
    dev tun0
	proto tcp
	remote 8.8.8.8 443
	cipher AES-128-CBC
	auth SHA1
	resolv-retry infinite
	nobind
	persist-key
	persist-tun
	auth-user-pass auth.txt
	client
	verb 3
`

func TestStore(t *testing.T) {
	pr := make(map[string]entity.Profile)
	pr["/etc/openvpn/test.ovpn"] = entity.Profile{Body: test_config}
	db := NewDb(WithMapMemory(pr))
	profile := entity.Profile{
		Body: `remote localhost`,
		Path: "/etc/openvpn/localhost.ovpn",
	}
	db.Store(profile)
	key := db.memory["/etc/openvpn/localhost.ovpn"]
	require.Equal(t, key, profile)
}
