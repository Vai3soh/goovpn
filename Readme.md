Goovpn

This gui openvpn client for linux.
Program uses the following libraries:

| Package                                   | Changes commits 
| ----------------------------------------- | ----------------------------------------
| github.com/mysteriumnetwork/go-openvpn    | 7ec797ccb0654e1ecc5459b1199471afcf2e9554
| github.com/therecipe/qt                   |

Build:

git clone github.com/Vai3soh/goovpn

Build binary with docker:
make build_docker

And then run build(deb/rpm package):
make build_package

Or build appimage:
make build_appimage

Download deb, rpm package in realese:

Install package:
sudo dpkg -i goovpn-1.0.0.x86_64.deb/sudo dnf goovpn-1.0.0.x86_64.deb

After install run:
sudo edit /etc/goovpn/config.yml

and modify path to configs dir (configs_path: '~/ovpnconfigs/')
add path to current user aka: /home/user/ovpnconfigs
move your openvpn configs files to this dir and run program:

1. From terminal: goovpn -config /etc/goovpn/config.yml
2. From menu in DE 

If use Goovpn-x86_64.AppImage, config file is located ~/.config/goovpn/config.yml

Screenshot:
![Data_Label](https://raw.githubusercontent.com/Vai3soh/goovpn/master/goovpn_screen.png)


