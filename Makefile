CONTAINER_PATH := /home/user/work/src/github.com/Vai3soh/cmd/app/deploy/linux/app
project_name = $(notdir $(shell pwd))
main_path = ./cmd/app
submodule_path := third_party/go-openvpn

generate: mock-gen

mock-gen:
	@rm -rf ./test/mocks/packages
	@go generate ./...

build_docker:
	git submodule update --init
	cd $(submodule_path) && git checkout goovpn_dev && \
	scripts/xgo_run.sh scripts/build-bridge.sh && cd -
	cp $(submodule_path)/openvpn3/bridge/libopenvpn3_linux_amd64.a pkg/openvpn3/bridge/
	@docker pull therecipe/qt:linux_static
	@docker build -t goovpn:latest -f Dockerfile .
	@docker run --name goovpn goovpn:latest
	@docker cp goovpn:$(CONTAINER_PATH) $(main_path)/$(project_name)
	@docker stop goovpn 
	@docker rm goovpn
	@docker image rm goovpn:latest
	@docker image rm therecipe/qt:linux_static
	rm -rf $(submodule_path)

build_package:
	@nfpm package -t ./build/package -p deb
	@nfpm package -t ./build/package -p rpm

build_appimage: 
	./scripts/build_appimage.sh
	
fmt:
	gofmt -s -w .
	
clean:
	@rm -rf $(main_path)/$(project_name)

changelog_update:
	git log v1.0.0...v1.0.2 --oneline --decorate > CHANGELOG.md
