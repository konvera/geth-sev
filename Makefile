GETH_BRANCH ?= main
GETH_REPO ?= https://github.com/flashbots/builder

GORUN = env GO111MODULE=on go run
SRCDIR = go-ethereum

GOMODCACHE = $(shell go env GOMODCACHE)

.PHONY: all
all: geth

############################## GETH EXECUTABLE ###############################

# Clone Geth and fetch dependencies
$(SRCDIR)/Makefile:
	git clone -b $(GETH_BRANCH) $(GETH_REPO) $(SRCDIR)
	cd $(SRCDIR) && go mod download

# patch Geth
$(SRCDIR)/PATCHED: $(SRCDIR)/Makefile
ifeq ($(TLS),1)
	patch -d $(SRCDIR) -p1 < geth-patches/0003-go-ethereum-tls.patch
endif
ifeq ($(PROTECT),1)
	patch -d $(SRCDIR) -p1 < geth-patches/0004-protect.patch
endif
	touch $(SRCDIR)/PATCHED

# Build Geth
$(SRCDIR)/build/bin/geth:  $(SRCDIR)/PATCHED
	cd $(SRCDIR) && \
		go mod tidy && \
		go build -ldflags "-extldflags '-Wl,-z,stack-size=0x800000,-fuse-ld=gold'" -tags urfave_cli_no_docs -trimpath -v -o $(PWD)/$(SRCDIR)/build/bin/geth ./cmd/geth

########################### COPIES OF EXECUTABLES #############################

# Geth build process creates the final executable as build/bin/geth. For
# simplicity, copy it into our root directory.

geth: $(SRCDIR)/build/bin/geth
	cp $< $@

################################## CLEANUP ####################################

.PHONY: distclean
distclean:
	$(RM) -rf $(SRCDIR) geth
