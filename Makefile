.PHONY: windows

windows:
	@echo "Building for Windows..."
	go build -tags lib -ldflags="-X polaris/db.Version=$(git describe --tags --long)" -buildmode=c-shared -o ui/windows/libpolaris.dll ./cmd/binding
	cd ui && flutter build windows