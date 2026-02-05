.PHONY: build run tidy clean pack

BINARY := clscli
ifeq ($(OS),Windows_NT)
	BINARY := clscli.exe
endif

DIST_DIR := dist/clscli
PACK_ARCHIVE := dist/clscli-skill.zip

build:
	go build -o $(BINARY) .

run: build
	./$(BINARY) $(ARGS)

tidy:
	go mod tidy

clean:
	rm -f $(BINARY)
	rm -rf dist

# Pack skill: build binary, gen _meta.json (version=tag, publishedAt=now), copy SKILL.md + binary to dist/clscli/, then zip.
pack: build
	@mkdir -p $(DIST_DIR)
	@V=$$(git describe --tags --always 2>/dev/null | sed 's/^v//') || V="1.0.0"; \
	TS=$$(date +%s)000; \
	echo '{"ownerId":"dbwang0130","slug":"clscli","version":"'$$V'","publishedAt":'$$TS'}' > $(DIST_DIR)/_meta.json
	cp SKILL.md $(DIST_DIR)/
	cp $(BINARY) $(DIST_DIR)/
	@(cd dist && zip -r clscli-skill.zip clscli)
	@echo "pack: $(PACK_ARCHIVE)"
