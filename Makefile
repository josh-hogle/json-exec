COMMAND_NAME ?= json-exec
TAG ?= $(shell git tag | sort -V | tail -1)
VERSION ?= $(shell if [ "${TAG}" != "" ]; then echo ${TAG} | cut -c2- ; else echo "0.0.0"; fi)
BUILD ?= $(shell git rev-parse HEAD | head -c 8)
TAG_DATE ?= $(shell git log --tags --simplify-by-decoration --pretty="format:%ai %d" | grep "tag: ${TAG}" | cut -d' ' -f1)
RELEASE_DATE ?= $(shell if [ "${TAG_DATE}" != "" ]; then echo "${TAG_DATE}"; else echo "Unreleased"; fi)
MAKEFILE_DIR := $(shell dirname $(firstword $(MAKEFILE_LIST)))
DIST_DIR ?= ${MAKEFILE_DIR}/dist
LDFLAGS_COMMON := \
	-X 'github.com/josh-hogle/json-exec/internal/app.CommandName=${COMMAND_NAME}' \
	-X 'github.com/josh-hogle/json-exec/internal/app.Version=${VERSION}' \
	-X 'github.com/josh-hogle/json-exec/internal/app.Build=${BUILD}' \
	-X 'github.com/josh-hogle/json-exec/internal/app.ReleaseDate=${RELEASE_DATE}'
LDFLAGS_RELEASE := \
	-X 'github.com/josh-hogle/json-exec/internal/app.IsDevelopment=false'
LDFLAGS_DEV := \
	-X 'github.com/josh-hogle/json-exec/internal/app.IsDevelopment=true'
SOURCES := $(shell find "${MAKEFILE_DIR}" -name "*.go")

# --- meta targets ---
all: release

release: build_release
linux: linux_amd64
darwin: darwin_amd64
windows: windows_amd64

dev: build_dev
dev_linux: dev_linux_amd64
dev_darwin: dev_darwin_amd64
dev_windows: dev_windows_amd64

clean:
	rm -rf "${DIST_DIR}"

build_release: \
	linux_amd64 \
	darwin_amd64 \
	windows_amd64

build_dev: \
	dev_linux_amd64 \
	dev_darwin_amd64 \
	dev_windows_amd64

# --- binary targets ---
dev_linux_amd64: export GOOS=linux
dev_linux_amd64: export GOARCH=amd64
dev_linux_amd64: ${SOURCES}
	mkdir -p "${DIST_DIR}/dev/${COMMAND_NAME}-${VERSION}-${GOOS}_${GOARCH}"
	go build -ldflags="${LDFLAGS_COMMON} ${LDFLAGS_DEV}" \
		-o "${DIST_DIR}/dev/${COMMAND_NAME}-${VERSION}-${GOOS}_${GOARCH}/${COMMAND_NAME}" \
		"${MAKEFILE_DIR}/cmd/${COMMAND_NAME}"

dev_darwin_amd64: export GOOS=darwin
dev_darwin_amd64: export GOARCH=amd64
dev_darwin_amd64: ${SOURCES}
	mkdir -p "${DIST_DIR}/dev/${COMMAND_NAME}-${VERSION}-${GOOS}_${GOARCH}"
	go build -ldflags="${LDFLAGS_COMMON} ${LDFLAGS_DEV}" \
		-o "${DIST_DIR}/dev/${COMMAND_NAME}-${VERSION}-${GOOS}_${GOARCH}/${COMMAND_NAME}" \
		"${MAKEFILE_DIR}/cmd/${COMMAND_NAME}"

dev_windows_amd64: export GOOS=windows
dev_windows_amd64: export GOARCH=amd64
dev_windows_amd64: ${SOURCES}
	mkdir -p "${DIST_DIR}/dev/${COMMAND_NAME}-${VERSION}-${GOOS}_${GOARCH}"
	go build -ldflags="${LDFLAGS_COMMON} ${LDFLAGS_DEV}" \
		-o "${DIST_DIR}/dev/${COMMAND_NAME}-${VERSION}-${GOOS}_${GOARCH}/${COMMAND_NAME}.exe" \
		"${MAKEFILE_DIR}/cmd/${COMMAND_NAME}"

linux_amd64: export GOOS=linux
linux_amd64: export GOARCH=amd64
linux_amd64: ${SOURCES}
	mkdir -p "${DIST_DIR}/release/${COMMAND_NAME}-${VERSION}-${GOOS}_${GOARCH}"
	go build -ldflags="${LDFLAGS_COMMON} ${LDFLAGS_RELEASE}" \
		-o "${DIST_DIR}/release/${COMMAND_NAME}-${VERSION}-${GOOS}_${GOARCH}/${COMMAND_NAME}" \
		"${MAKEFILE_DIR}/cmd/${COMMAND_NAME}"
	cp "${MAKEFILE_DIR}/README.md" "${DIST_DIR}/release/${COMMAND_NAME}-${VERSION}-${GOOS}_${GOARCH}/"
	cp "${MAKEFILE_DIR}/LICENSE" "${DIST_DIR}/release/${COMMAND_NAME}-${VERSION}-${GOOS}_${GOARCH}/"
	tar -zcvf "${DIST_DIR}/${COMMAND_NAME}-${VERSION}-${GOOS}_${GOARCH}.tar.gz" -C "${DIST_DIR}/release" \
		"./${COMMAND_NAME}-${VERSION}-${GOOS}_${GOARCH}"

darwin_amd64: export GOOS=darwin
darwin_amd64: export GOARCH=amd64
darwin_amd64: ${SOURCES}
	mkdir -p "${DIST_DIR}/release/${COMMAND_NAME}-${VERSION}-${GOOS}_${GOARCH}"
	go build -ldflags="${LDFLAGS_COMMON} ${LDFLAGS_RELEASE}" \
		-o "${DIST_DIR}/release/${COMMAND_NAME}-${VERSION}-${GOOS}_${GOARCH}/${COMMAND_NAME}" \
		"${MAKEFILE_DIR}/cmd/${COMMAND_NAME}"
	cp "${MAKEFILE_DIR}/README.md" "${DIST_DIR}/release/${COMMAND_NAME}-${VERSION}-${GOOS}_${GOARCH}/"
	cp "${MAKEFILE_DIR}/LICENSE" "${DIST_DIR}/release/${COMMAND_NAME}-${VERSION}-${GOOS}_${GOARCH}/"
	tar -zcvf "${DIST_DIR}/${COMMAND_NAME}-${VERSION}-${GOOS}_${GOARCH}.tar.gz" -C "${DIST_DIR}/release" \
		"./${COMMAND_NAME}-${VERSION}-${GOOS}_${GOARCH}"

windows_amd64: export GOOS=windows
windows_amd64: export GOARCH=amd64
windows_amd64: ${SOURCES}
	mkdir -p "${DIST_DIR}/release/${COMMAND_NAME}-${VERSION}-${GOOS}_${GOARCH}"
	go build -ldflags="${LDFLAGS_COMMON} ${LDFLAGS_RELEASE}" \
		-o "${DIST_DIR}/release/${COMMAND_NAME}-${VERSION}-${GOOS}_${GOARCH}/${COMMAND_NAME}.exe" \
		"${MAKEFILE_DIR}/cmd/${COMMAND_NAME}"
	cp "${MAKEFILE_DIR}/README.md" "${DIST_DIR}/release/${COMMAND_NAME}-${VERSION}-${GOOS}_${GOARCH}/"
	cp "${MAKEFILE_DIR}/LICENSE" "${DIST_DIR}/release/${COMMAND_NAME}-${VERSION}-${GOOS}_${GOARCH}/"
	cd "${DIST_DIR}/release" \
		&& zip -r "../${COMMAND_NAME}-${VERSION}-${GOOS}_${GOARCH}.zip" "${COMMAND_NAME}-${VERSION}-${GOOS}_${GOARCH}"
