BINARY := ghostty-theme-tui
PREFIX ?= $(HOME)/.local
BINDIR ?= $(PREFIX)/bin
GO ?= go

.PHONY: build install uninstall test clean

build:
	$(GO) build -o $(BINARY) ./cmd/ghostty-theme-tui

install:
	mkdir -p "$(BINDIR)"
	$(GO) build -o "$(BINDIR)/$(BINARY)" ./cmd/ghostty-theme-tui
	chmod +x "$(BINDIR)/$(BINARY)"
	@echo "Installed $(BINARY) to $(BINDIR)"
	@echo "Make sure $(BINDIR) is in your PATH."

uninstall:
	rm -f "$(BINDIR)/$(BINARY)"
	@echo "Removed $(BINDIR)/$(BINARY)"

test:
	$(GO) test ./...

clean:
	rm -rf dist "$(BINARY)"

