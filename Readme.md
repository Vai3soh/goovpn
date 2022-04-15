Goovpn

This gui openvpn client for linux.
Program uses the following libraries:

| Package                                   | Changes commits 
| ----------------------------------------- | ----------------------------------------
| github.com/mysteriumnetwork/go-openvpn    | 7ec797ccb0654e1ecc5459b1199471afcf2e9554
| github.com/therecipe/qt                   |
| github.com/fangdingjun/go-log/v5          |

Build:

```git clone github.com/Vai3soh/goovpn```

Build binary with docker:
```make build_docker```

And then run build(deb/rpm package):
```make build_package```

Or build appimage:
```make build_appimage```

Download deb, rpm, appimage package in realese:

```github.com/Vai3soh/goovpn/releases```

Install package:
```sudo dpkg -i goovpn_1.0.0_amd64.deb or sudo dnf goovpn-1.0.0.x86_64.rpm```

After install run:
```sudo edit /etc/goovpn/config.yml```

and modify path to configs dir (configs_path: '~/ovpnconfigs/')
add path to current user aka: ```/home/user/ovpnconfigs```,
create dir ```mkdir /home/user/ovpnconfigs```,
move your openvpn configs files to this dir and run program:

1. From terminal: ```goovpn -config /etc/goovpn/config.yml```
2. From menu in DE 

If use Goovpn-x86_64.AppImage, config file is located ```~/.config/goovpn/config.yml```

DNS query:
1. If your distr with systemd, modify config directive ```use_systemd: false``` set ```use_systemd: true``` and install ```systemd-resolve```.

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
2. If use ```use_systemd: false``` install ```resolvconf```

Screenshot:


![Data_Label](https://raw.githubusercontent.com/Vai3soh/goovpn/master/goovpn_screen.png)


