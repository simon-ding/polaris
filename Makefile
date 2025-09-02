.PHONY: windows

VERSION=$(shell git describe --tags --long)

windows:
	@echo "Building for Windows..."
	go build -tags lib -ldflags="-X polaris/db.Version=$(git describe --tags --long)" -buildmode=c-shared -o ui/windows/libpolaris.dll ./cmd/binding
	cd ui && flutter build windows

polaris-web:
	@echo "Building..."
	CGO_ENABLED=0 go build -o polaris -ldflags="-X polaris/db.Version=$(VERSION) -X polaris/db.DefaultTmdbApiKey=$(TMDB_API_KEY)"  ./cmd/polaris