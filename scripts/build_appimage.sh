#!/usr/bin/env bash

set -eux
version=`grep -r 'version' nfpm.yaml | awk '{ print $NF }' | tr -d \"v`
dir_build='/tmp/AppDir'

[ -f ./build/bin/goovpn ] || exit 0

mkdir -p ${dir_build}/goovpn/
mkdir -p ${dir_build}/goovpn/usr/bin/
mkdir -p ${dir_build}/goovpn/usr/sbin/
mkdir -p ${dir_build}/goovpn/usr/lib/
mkdir -p ${dir_build}/goovpn/usr/share/icons/hicolor/128x128/apps/
mkdir -p ${dir_build}/goovpn/etc/goovpn/

cp ./build/bin/goovpn ${dir_build}/goovpn/usr/sbin/goovpn
chmod +x ${dir_build}/goovpn/usr/sbin/goovpn

cat <<EOF > ${dir_build}/goovpn/usr/bin/goovpn
#!/usr/bin/env bash

dir=\$(dirname \$PWD)
cmd=\${dir}/usr/sbin/goovpn

mkdir -p /tmp/appimage_root/  \$HOME/.config/goovpn/
cp \$cmd /tmp/appimage_root/

pkexec --disable-internal-agent \
env DISPLAY=\$DISPLAY SUDO_USER=\$USER \
DB_GOOVPN=\$HOME/.config/goovpn/goovpn.db \
XAUTHORITY=\$XAUTHORITY /tmp/appimage_root/goovpn 
EOF

chmod +x ${dir_build}/goovpn/usr/bin/goovpn


cat <<EOF > ${dir_build}/goovpn/goovpn.desktop
[Desktop Entry]
Name=Goovpn
Exec=goovpn
Icon=goovpn
Terminal=false
Type=Application
Categories=Network;Qt;
EOF

sed -i 's/HOME/\\$HOME/'g ${dir_build}/goovpn/goovpn.desktop

cp ./embedfile/assets/app.png ${dir_build}/goovpn/goovpn.png
cp ./embedfile/assets/app.png  ${dir_build}/goovpn/usr/share/icons/hicolor/128x128/apps/goovpn.png
ldd ./build/bin/goovpn | grep "=> /" | awk '{print $3}' | \
xargs -I '{}' cp -v '{}' ${dir_build}/goovpn/usr/lib/

url='https://github.com/AppImage/AppImageKit/releases/download/continuous/AppRun-x86_64'
wget -O ${dir_build}/goovpn/AppRun ${url}
chmod +x ${dir_build}/goovpn/AppRun

url='https://github.com/AppImage/AppImageKit/releases/download/continuous/appimagetool-x86_64.AppImage'
wget -O ${dir_build}/appimagetool-x86_64.AppImage ${url}
chmod +x ${dir_build}/appimagetool-x86_64.AppImage
cd ${dir_build}/ && ARCH=x86_64 ./appimagetool-x86_64.AppImage -v -n --comp xz goovpn
cd - || exit

mv ${dir_build}/Goovpn-x86_64.AppImage \
./build/package/Goovpn-${version}_x86_64.AppImage
rm -rf ${dir_build}/

