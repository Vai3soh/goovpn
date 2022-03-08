FROM therecipe/qt:linux_static AS qt

RUN apt-get update && \
apt-get --no-install-recommends -y install git curl && \
rm -rf /usr/local/go/ && \
GO=go1.17.6.linux-amd64.tar.gz && \
curl -sL --retry 10 --retry-delay 60 -O https://go.dev/dl/$GO && \
tar -xzf $GO -C /usr/local

WORKDIR /home/user/work/src/github.com/Vai3soh/
COPY . /home/user/work/src/github.com/Vai3soh/

RUN go get github.com/therecipe/qt/internal/binding/files/docs/5.13.0@v0.0.0-20200904063919-c0c124a5770d

ENV QT_MXE_ARCH=amd64
WORKDIR /home/user/work/src/github.com/Vai3soh/cmd/app/
RUN go mod vendor && qtdeploy build desktop 
