project_name = $(notdir $(shell pwd))
main_path = ./cmd/app

generate: mock-gen

mock-gen:
	@rm -rf ./test/mocks/packages
	@go generate ./...

build_bin_linux:
	cd $(main_path) && wails build -platform linux/amd64 

build_bin_windows:
	cd $(main_path) && CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++-posix wails build -platform windows/amd64 -skipbindings -ldflags "-linkmode external -extldflags -static"  
		
build_debug_race:
	cd $(main_path) && wails build -debug -race

build_package:
	@nfpm package -t ./build/package -p deb
	@nfpm package -t ./build/package -p rpm

build_appimage: 
	./scripts/build_appimage.sh
	
fmt:
	gofmt -s -w .
	
clean:
	@rm -rf ./build/bin/goovpn ./build/bin/goovpn.exe
	
changelog_update:
	rm -rf changelog.yml
	chglog init
	chglog format --template repo > CHANGELOG.md
