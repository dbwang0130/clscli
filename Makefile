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

# Pack skill: build binary, copy SKILL.md + binary to dist/clscli/, then zip.
pack: build
	@mkdir -p $(DIST_DIR)
	cp SKILL.md $(DIST_DIR)/
	cp $(BINARY) $(DIST_DIR)/
	@(cd dist && zip -r clscli-skill.zip clscli)
	@echo "pack: $(PACK_ARCHIVE)"
