package gui

import (
	"fmt"
	"runtime"
	"strconv"

	"github.com/Vai3soh/goovpn/entity"
)

func getBool(v *string) bool {
	if v != nil {
		return *v != ""
	}
	return false
}

func getUint(v string) (*uint, error) {
	u64, err := strconv.ParseUint(v, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("convert string to uint failed: [%s]", err)

	}
	r := uint(u64)
	return &r, nil
}

func (g *Gui) setSessionParamsOtherOpt(data []entity.Message) error {
	for _, e := range data {
		switch e.AtrId {
		case `with_conn_timeout`:
			t, _ := strconv.Atoi(e.Value)
			g.sessionService.SetConnTimeout(t)
		case `with_allow_unused_addr_families`:
			g.sessionService.SetAllowUnusedAddrFamilies(e.Value)
		case `with_clock_tick_ms`:
			v, err := getUint(e.Value)
			if err != nil {
				return err
			}
			g.sessionService.SetClockTickMS(*v)
		case `with_external_pki_alias`:
			g.sessionService.SetExternalPkiAlias(e.Value)
		case `with_gremlin_config`:
			g.sessionService.SetGremlinConfig(e.Value)
		case `with_gui_version`:
			g.sessionService.SetGuiVersion(e.Value)
		case `with_hwaddr_override`:
			g.sessionService.SetHwAddrOverride(e.Value)
		case `with_platform_version`:
			g.sessionService.SetPlatformVersion(e.Value)
		case `with_port_override`:
			g.sessionService.SetPortOverride(e.Value)
		case `with_private_key_password`:
			g.sessionService.SetPrivateKeyPassword(e.Value)
		case `with_proto_override`:
			g.sessionService.SetProtoOverride(e.Value)
		case `with_proto_version_override`:
			t, _ := strconv.Atoi(e.Value)
			g.sessionService.SetProtoVersionOverride(t)
		case `with_proxy_allow_clear_text_auth`:
			g.sessionService.SetProxyAllowCleartextAuth(getBool(&e.Value))
		case `with_proxy_host`:
			g.sessionService.SetProxyHost(e.Value)
		case `with_proxy_username`:
			g.sessionService.SetProxyUsername(e.Value)
		case `with_proxy_password`:
			g.sessionService.SetProxyPassword(e.Value)
		case `with_proxy_port`:
			g.sessionService.SetProxyPort(e.Value)
		case `with_server_override`:
			g.sessionService.SetServerOverride(e.Value)
		case `with_sso_methods`:
			g.sessionService.SetSsoMethods(e.Value)
		case `with_tls_cert_profile_override`:
			g.sessionService.SetTlsCertProfileOverride(e.Value)
		case `with_tls_cipher_list`:
			g.sessionService.SetTlsCipherList(e.Value)
		case `with_tls_cipher_suites_list`:
			g.sessionService.SetTlsCiphersuitesList(e.Value)
		case `with_tls_version_min_override`:
			g.sessionService.SetTlsVersionMinOverride(e.Value)
		case `with_default_key_direction`:
			t, _ := strconv.Atoi(e.Value)
			g.sessionService.SetDefaultKeyDirection(t)
		}

	}
	return nil
}

func (g *Gui) setSessionParamsGeneralOvpnLib(data []entity.Message) error {
	for _, e := range data {
		switch e.AtrId {
		case `tun_persist`:
			g.sessionService.SetTunPersist(getBool(&e.Value))
		case `legacy_algo`:
			g.sessionService.SetEnableLegacyAlgorithms(getBool(&e.Value))
			if runtime.GOOS == "windows" {
				g.sessionService.SetEnableLegacyAlgorithms(false)
			}
		case `preferred_dc_algo`:
			g.sessionService.SetEnableNonPreferredDCAlgorithms(getBool(&e.Value))
		case `disable_client_cert`:
			g.sessionService.SetDisableClientCert(getBool(&e.Value))
		case `retry_on_failed`:
			g.sessionService.SetRetryOnAuthFailed(getBool(&e.Value))
		case `allow_local_dns_resolvers`:
			g.sessionService.SetAllowLocalDnsResolvers(getBool(&e.Value))
		case `allow_local_lan_access`:
			g.sessionService.SetAllowLocalLanAccess(getBool(&e.Value))
		case `alt_proxy`:
			g.sessionService.SetAltProxy(getBool(&e.Value))
		case `auto_login_sessions`:
			g.sessionService.SetAutologinSessions(getBool(&e.Value))
		case `use_dco`:
			g.sessionService.SetDco(getBool(&e.Value))
		case `with_echo`:
			g.sessionService.SetEcho(getBool(&e.Value))
		case `tun_builder_cupture_event`:
			g.sessionService.SetGenerate_tun_builder_capture_event(getBool(&e.Value))
		case `google_dns_fallback`:
			g.sessionService.SetGoogleDnsFallback(getBool(&e.Value))
		case `info`:
			g.sessionService.SetInfo(getBool(&e.Value))
		case `proxy_allow_clear_text_auth`:
			g.sessionService.SetProxyAllowCleartextAuth(getBool(&e.Value))
		case `synch_dns_lookup`:
			g.sessionService.SetSynchronousDnsLookup(getBool(&e.Value))
		case `enable_win_tun`:
			g.sessionService.SetWintun(getBool(&e.Value))

		}
	}
	return nil
}

func (g *Gui) setSessionParamsDebugSslAndCompression(data []entity.Message) {
	for _, e := range data {
		switch e.AtrId {
		case `#ssl`:
			v, _ := strconv.Atoi(e.Value)
			g.sessionService.SetSslDebugLevel(v)
		case `#cmp`:
			g.sessionService.SetCompressionMode(e.Value)
		}
	}
}
