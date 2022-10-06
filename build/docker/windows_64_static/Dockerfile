#docker build -t therecipe/qt:builder -f Dockerfile .
#docker load --input builder_win_static.tar

FROM ubuntu:20.04 as base
LABEL maintainer therecipe

ENV USER user
ENV HOME /home/$USER
ENV GOPATH $HOME/work

COPY --from=therecipe/qt:linux $GOPATH/bin $GOPATH/bin
COPY --from=therecipe/qt:linux /usr/local/go /usr/local/go
COPY --from=therecipe/qt:linux $GOPATH/src/github.com/therecipe/qt \
$GOPATH/src/github.com/therecipe/qt

COPY --from=therecipe/qt:windows_64_static_base_posix /usr/lib/mxe/ \
/usr/lib/mxe/

RUN tar -cvzpf /tmp/all.tar.gz /usr/local/go $GOPATH/bin \
$GOPATH/src/github.com/therecipe/qt /usr/lib/mxe/

FROM ubuntu:20.04
LABEL maintainer therecipe

ENV USER=user \
	HOME=/home/user \
	GOPATH=/home/user/work \
	PATH=/usr/lib/mxe/usr/bin:/usr/local/go/bin:$PATH \
	QT_DOCKER=true \
	QT_MXE=true \
	QT_MXE_ARCH=amd64 \
	QT_MXE_STATIC=true \
	QT_PROJ_GIT_PATH=/home/user/work/src/github.com/therecipe/qt

COPY --from=base /tmp/all.tar.gz /tmp/all.tar.gz

RUN tar -xvzpf /tmp/all.tar.gz && rm -rf /tmp/all.tar.gz && \
apt-get -qq update && DEBIAN_FRONTEND=noninteractive apt-get \
--no-install-recommends -qq -y install \
ca-certificates git pkg-config 

RUN cd $QT_PROJ_GIT_PATH && \
$GOPATH/bin/qtsetup prep && \
$GOPATH/bin/qtsetup check windows && \
$GOPATH/bin/qtsetup generate windows && \
$GOPATH/bin/qtsetup install windows && \
cd $GOPATH/src/github.com/therecipe/qt/internal/examples/widgets/line_edits && \
$GOPATH/bin/qtdeploy build windows && rm -rf ./deploy
