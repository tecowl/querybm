METADATA_YAML=.project.yaml
$(METADATA_YAML): metadata-gen

METADATA_LINTERS=$(strip $(shell $(MAKE) linters-enabled --no-print-directory 2>/dev/null | grep . | wc -l))
.PHONY: metadata-gen
metadata-gen: 
	@echo "linters: $(METADATA_LINTERS)" > $(METADATA_YAML)
