ifeq ($(DESTDIR),)
	DESTDIR := /usr/local
endif

.DEFAULT_GOAL := help

.PHONY: help
help:
	@echo "Usage: make [target]"
	@cat Makefile | awk -F ":.*##"  '/##/ { printf "    %-17s %s\n", $$1, $$2 }' | grep -v  grep

.PHONY: format-c
format-c: ## formats all the C code
	cd testsuite/testdata && clang-format --style=file -i *.c *.h
	clang-format --style=file -i *.h

.PHONY: build
build: ## builds the Linux static libraries and leave them and a copy of the definitions in .build directory
	go build -ldflags="-s -w" -buildmode c-archive -o .build/uplink.a .
	mv .build/uplink.a .build/libuplink.a
	cp uplink_definitions.h .build/
	./uplinkc_builder.sh

.PHONY: build-gpl2
build-gpl2: ## builds the Linux static libraries GPL2 license compatible and leave them and a copy of the definitions in .build directory
	./check-licenses.sh
	go build -modfile=go-gpl2.mod -ldflags="-s -w" -buildmode c-archive -tags stdsha256 -o .build/uplink.a .
	mv .build/uplink.a .build/libuplink.a
	cp uplink_definitions.h .build/
	./uplinkc_builder.sh

.PHONY: build-gpl2-linux2win
build-gpl2-linux2win: ## cross-compiles the Windows static libraries GPL2 license compatible from Linux and leave them and a copy of the definitions in .build directory
	@if [ $(shell uname) = Linux ]; then \
		./check-licenses.sh; \
		GOOS="windows" GOARCH="amd64" CGO_ENABLED="1" CXX="x86_64-w64-mingw32-g++" CC="x86_64-w64-mingw32-gcc" go build -modfile=go-gpl2.mod -ldflags="-s -w" -buildmode c-archive -tags stdsha256 -o .build/uplink.a .; \
		mv .build/uplink.a .build/libuplink.a; \
		cp uplink_definitions.h .build/; \
		./uplinkc_builder.sh; \
	else \
		echo "The host OS is not Linux";\
	fi

.PHONY: bump-dependencies
bump-dependencies: ## bumps the dependencies
	go get storj.io/common@master storj.io/uplink@master
	go mod tidy
	cd testsuite;\
		go get storj.io/common@master storj.io/storj@master storj.io/uplink@master;\
		go mod tidy

.PHONY: install
install: ## installs the uplink-c build, header files and the pkg-config file into the destination or specified directory
	if [ $(shell uname) = Darwin ]; then \
		sed -i '' "s|prefix=@DESTDIR@|prefix=\"$(DESTDIR)\"|g" libuplinkc.pc; \
	else \
		sed -i "s|prefix=@DESTDIR@|prefix=\"$(DESTDIR)\"|g" libuplinkc.pc; \
	fi
	mkdir -p "$(DESTDIR)/include/storj/"
	mkdir -p "$(DESTDIR)/lib/pkgconfig/"
	cp .build/libuplink.a "$(DESTDIR)/lib"
	cp .build/uplink.h "$(DESTDIR)/include/storj/"
	cp .build/uplink_definitions.h "$(DESTDIR)/include/storj/"
	cp libuplinkc.pc  "$(DESTDIR)/lib/pkgconfig"
	if [ $(shell uname) = Darwin ]; then \
		sed -i '' "s|prefix=\"$(DESTDIR)\"|prefix=@DESTDIR@|g" libuplinkc.pc; \
	else \
		sed -i "s|prefix=\"$(DESTDIR)\"|prefix=@DESTDIR@|g" libuplinkc.pc; \
	fi
