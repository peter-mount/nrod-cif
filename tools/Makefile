
MODULES	= $(shell ls -d */ | sed "s|/||")

DEPS = $(foreach PLATFORM,$(PLATFORMS), \
		$(eval GOOS=$(word 1,$(subst :, ,$(PLATFORM)))) \
		$(eval GOARCH=$(word 2,$(subst :, ,$(PLATFORM)))) \
		$(eval GOARM=$(word 3,$(subst :, ,$(PLATFORM)))) \
		$(eval BUILD=../$(BUILDS)/$(GOOS)/$(GOARCH)$(GOARM)) \
		$(foreach MODULE,$(MODULES),$(BUILD)/$(MODULE) )\
	)

DISTRIBUTIONS = $(foreach PLATFORM,$(PLATFORMS), \
		$(eval GOOS=$(word 1,$(subst :, ,$(PLATFORM)))) \
		$(eval GOARCH=$(word 2,$(subst :, ,$(PLATFORM)))) \
		$(eval GOARM=$(word 3,$(subst :, ,$(PLATFORM)))) \
		../$(DIST)/$(DIST_PREFIX)-$(GOOS)-$(GOARCH)$(GOARM).tgz)

.PHONY:	all dist

all: $(DEPS)

dist: $(DISTRIBUTIONS)

# This builds an individual tool based on it's relative path
../$(BUILDS)/%:
	$(eval REL=$(subst /, ,$@))
	$(eval GOOS=$(word 3,$(REL)))
	$(eval GOARCH=$(word 4,$(REL)))
	$(eval MODULE=$(word 5,$(REL)))
	$(eval GOARM=$(shell echo -n $(filter arm%,$(GOARCH)) | sed -E -e "s/arm([567])$$/\\1/g" -e "s/arm64$$//g"))
	$(eval GOARCH=$(shell echo -n $(GOARCH) | sed -E "s/arm([567])$$/arm/g"))
	@$(MKDIR) $(shell dirname $@)
	@$(PRINTF) "Compiling %s %s%s %s\n" "$(GOOS)" "$(GOARCH)" "$(GOARM)" "$(MODULE)"
	@CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) $(GO) build -o $@ $(MODULE)/bin/main.go

../$(DIST)/%.tgz:
	$(eval REL=$(subst _, ,$(basename $@)))
	$(eval GOOS=$(word 3,$(REL)))
	$(eval GOARCH=$(word 4,$(REL)))
	$(eval GOARM=$(shell echo -n $(filter arm%,$(GOARCH)) | sed -E -e "s/arm([567])$$/\\1/g" -e "s/arm64$$//g"))
	$(eval GOARCH=$(shell echo -n $(GOARCH) | sed -E "s/arm([567])$$/arm/g"))
	$(eval BUILD=../$(BUILDS)/$(GOOS)/$(GOARCH)$(GOARM))
	@$(MKDIR) $(shell dirname $@)
	@$(PRINTF) "Generating %s\n" $$(basename $@)
	@$(TAR) -P --transform "s|^../$(BUILDS)|$(DIST_PREFIX)|" -czpf $@ $(BUILD)
