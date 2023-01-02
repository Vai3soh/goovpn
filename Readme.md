Goovpn

This gui openvpn client for linux, windows.
Program uses the following libraries:

| Package                                   | Changes commits 
| ----------------------------------------- | ----------------------------------------
| github.com/Vai3soh/ovpncli (core)         |
| github.com/wailsapp/wails                 |
| github.com/fangdingjun/go-log/v5          |


Build:

```git clone github.com/Vai3soh/goovpn```

Build binary:
	linux:
		```make build_bin_linux```
	windows:
		```make build_bin_windows```

And then run build(deb/rpm package):
```make build_package```

Or build appimage:
```make build_appimage```

Download deb, rpm, appimage package in realese:

```github.com/Vai3soh/goovpn/releases```

Install package:
```sudo dpkg -i goovpn_1.0.3_amd64.deb or sudo dnf goovpn-1.0.3.x86_64.rpm```

After install run:

Run soft ```goovpn``` and modify path to configs dir (option ```Configs dir path```)
add path to current user aka: ```/home/user/ovpnconfigs```,
create dir ```mkdir /home/user/ovpnconfigs```,
move your openvpn configs files to this dir.

DNS query:
1. If your distr with systemd, enable the option ```Use systemd``` and install ```systemd-resolve```.

After restart unit systemd-resolved, your ```/etc/resolv.conf```
```
# This file is managed by man:systemd-resolved(8). Do not edit.
# ....

nameserver 127.0.0.53
options edns0 trust-ad
search .
```
and after connect (resolvectl status) example:
 ```resolvectl status
Global
         Protocols: +LLMNR +mDNS -DNSOverTLS DNSSEC=no/unsupported
  resolv.conf mode: stub
Current DNS Server: 192.168.1.1
       DNS Servers: 192.168.1.1

Link 1 (enp1s0f1)
Current Scopes: LLMNR/IPv4 LLMNR/IPv6
     Protocols: -DefaultRoute +LLMNR -mDNS -DNSOverTLS DNSSEC=no/unsupported

Link 93 (tun)
    Current Scopes: DNS LLMNR/IPv4 LLMNR/IPv6
         Protocols: +DefaultRoute +LLMNR -mDNS -DNSOverTLS DNSSEC=no/unsupported
Current DNS Server: 10.211.254.254
       DNS Servers: 10.211.254.254 8.8.8.8
        DNS Domain: ~.
```
2. If system not use systemd disable ```Use systemd``` install ```resolvconf```

For windows OS (dependencies):
	install webview2: https://developer.microsoft.com/en-us/microsoft-edge/webview2/
	install tap or wintun driver:
		go to ```https://swupdate.openvpn.org/community/releases/OpenVPN-2.5.8-I604-amd64.msi``` download installer
		run and custom install only tap and wintun driver.
	 
Screenshot:


![Data_Label](https://raw.githubusercontent.com/Vai3soh/goovpn/master/goovpn_screen1.png)
![Data_Label](https://raw.githubusercontent.com/Vai3soh/goovpn/master/goovpn_screen2.png)

