#
# By default this will build the project on every non-mobile platform
# supported by the installed go environment.
#
# To limit a build to a single environment, you can force it to just a
# single platform by prefixing make with:
#
# PLATFORMS=linux:amd64: make clean all
#
# Just change the entry for your OS and CPU. These are listed in platforms.md
#
# Note: For 32 bit arm processors the 3rd parameter is important.
# e.g. use linux:arm:6 or linux:arm:7
#
# For all other processors, including arm64, leave the third field blank.
#
# To disable tests, you can prefix make with:
#
# GO_TEST="#" make clean all
#
# The quotes are important!
#
# You can combine the two as necessary.
#
# e.g. GO_TEST="#" PLATFORMS=linux:amd64: make clean all
#
# For a parallel builds you can use the -j parameter to make as usual.
#
# e.g.: make -j 8 clean all
#
# Pick a value suitable to the number of cores/thread your machine has.
# This is useful for a full build of all platforms as it will build all
# of the binaries in parallel speeding up the full build.
#

# The repository name/package prefix.
# This should match the value of module in go.mod
export PACKAGE_PREFIX = $(shell grep ^module go.mod | cut -f2 -d' ' | head -1)
export PACKAGE_NAME = $(shell basename $(PACKAGE_PREFIX))
export DIST_PREFIX = $(PACKAGE_NAME)_latest

# List of modules to build.
#
# Note: tools should be last as that generates executables and this
# allows the other modules to perform any tests first.
MODULES		= tools

# Where to place build artifacts. These must be subdirectories here and not
# a path elsewhere, otherwise it will break the build!
export BUILDS 	= builds
export DIST		= dist

# Tool names
export CP     	= @cp -p
export ECHO		= echo
export GO		= go
export MKDIR  	= mkdir -p
export PRINTF	= printf
export TAR		= tar

# Append -test.v to GO_TEST to show status of each test.
# Without it, only shows total time per module if they pass
export GO_TEST	?= $(GO) test

.PHONY: all build clean dist init test validate-go-version resolve-platforms platforms.md

# Used to separate commands in foreach in the all and dist targets.
# Using this allows each step to fail as if we entered them individually
# in the Makefile.
#
# In the foreach, we then use ${\n} to delimit the end of each command
# in the recipe.
#
# NOTE this MUST have 2 empty lines between define and endef for it to work!
define \n


endef

all: init
	@$(MKDIR) -pv $(BUILDS)
	$(foreach MODULE,$(MODULES),@$(MAKE) -C $(MODULE) all${\n})

clean:
	@$(GO) clean -testcache
	@$(RM) -r $(BUILDS) $(DIST)

dist: all platforms.md
	@$(MKDIR) -pv $(DIST)
	$(foreach MODULE,$(MODULES) tools,@$(MAKE) -C $(MODULE) dist${\n})

init: validate-go-version resolve-platforms
	@$(GO) mod download

# Validates the installed version of go against the version declared in go.mod
MINIMUM_SUPPORTED_GO_MAJOR_VERSION	= $(shell grep "^go" go.mod | cut -f2 -d' ' | cut -f1 -d'.')
MINIMUM_SUPPORTED_GO_MINOR_VERSION	= $(shell grep "^go" go.mod | cut -f2 -d' ' | cut -f2 -d'.')
GO_MAJOR_VERSION = $(shell $(GO) version | cut -f3 -d' ' | cut -c 3- | cut -f1 -d' ' | cut -f1 -d'.')
GO_MINOR_VERSION = $(shell $(GO) version | cut -f3 -d' ' | cut -c 3- | cut -f1 -d' ' | cut -f2 -d'.')
GO_VERSION_VALIDATION_ERR_MSG = Golang version $(GO_MAJOR_VERSION).$(GO_MINOR_VERSION) is not supported, please update to at least $(MINIMUM_SUPPORTED_GO_MAJOR_VERSION).$(MINIMUM_SUPPORTED_GO_MINOR_VERSION)
validate-go-version:
	@if [ $(GO_MAJOR_VERSION) -gt $(MINIMUM_SUPPORTED_GO_MAJOR_VERSION) ]; then \
		exit 0 ;\
	elif [ $(GO_MAJOR_VERSION) -lt $(MINIMUM_SUPPORTED_GO_MAJOR_VERSION) ]; then \
		$(ECHO) '$(GO_VERSION_VALIDATION_ERR_MSG)';\
		exit 1; \
	elif [ $(GO_MINOR_VERSION) -lt $(MINIMUM_SUPPORTED_GO_MINOR_VERSION) ] ; then \
		$(ECHO) '$(GO_VERSION_VALIDATION_ERR_MSG)';\
		exit 1; \
	fi

# This discovers all platforms supported by the locally installed go compiler.
# This will only expand then if the PLATFORMS environment variable was not set
# when invoking make
resolve-platforms:
ifeq ("$(PLATFORMS)","")
	$(eval DISC_PLATFORMS=)
	$(foreach DISC_PLATFORM,$(shell $(GO) tool dist list), \
		$(eval GOOS=$(word 1,$(subst /, ,$(DISC_PLATFORM)))) \
		$(if $(filter android,$(GOOS)),,\
			$(if $(filter ios,$(GOOS)),,\
				$(eval GOARCH=$(word 2,$(subst /, ,$(DISC_PLATFORM)))) \
				$(foreach GOARM, \
					$(if $(filter arm,$(GOARCH)),6 7,:), \
					$(eval DISC_PLATFORMS=$(DISC_PLATFORMS) $(GOOS):$(GOARCH):$(GOARM)) \
				) \
			)\
		)\
	)
	$(eval export PLATFORMS=$(DISC_PLATFORMS))
endif

# Generates platforms.md based on the local go installation.
# This does nothing other than keep that page in sync with what is currently
# supported by go and the build system.
platforms.md: resolve-platforms
	$(shell ( \
		echo "# Supported Platforms"; \
		echo; \
		echo "The following platforms are supported by virtue of how the build system works:"; \
		echo; \
		echo "| Operating System | CPU Architectures |"; \
		echo "| ---------------- | ----------------- |"; \
		$(foreach OS, $(shell ls $(BUILDS)), echo "| $(OS) | $(foreach ARCH,$(shell ls $(BUILDS)/$(OS)),$(ARCH)) |"; ) \
		echo; \
		echo "Operating Systems: $(shell ls $(BUILDS)|wc -l) CPU Architectures: $(shell ls -d $(BUILDS)/*/*| cut -f3 -d'/' | sort |uniq | wc -l)"; \
		echo; \
		echo "This is all non-mobile platforms supported by go version \`$(GO_MAJOR_VERSION).$(GO_MINOR_VERSION)\`" ;\
		echo; \
		echo "This page is automatically generated from the output of \`go tool dist list\`"; \
	  ) >$@ \
	)
