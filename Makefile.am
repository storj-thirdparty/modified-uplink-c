all: ## builds the uplink-c with GPL2 license compatibility
	@mv COPYING _COPYING
	./check-licenses.sh
	./uplinkc-builder.sh
	@mv _COPYING COPYING

.PHONY: install
install: ## installs the uplink-c build and the pkg-config file into the specified prefix directory
	mkdir -p $(prefix)/include/storj/
	mv libuplinkc.a $(prefix)/lib
	mv libuplinkc.h $(prefix)/include/storj/
	cp uplink_definitions.h $(prefix)/include/storj/
	mv libuplinkc.pc  $(libdir)/pkgconfig