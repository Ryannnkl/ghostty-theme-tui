#!/usr/bin/env sh
set -eu

repo="${GHOSTTY_THEME_TUI_REPO:-Ryannnkl/ghostty-theme-tui}"
version="${GHOSTTY_THEME_TUI_VERSION:-latest}"
bindir="${GHOSTTY_THEME_TUI_BINDIR:-$HOME/.local/bin}"
binary="ghostty-theme-tui"

need() {
	if ! command -v "$1" >/dev/null 2>&1; then
		echo "error: missing required command: $1" >&2
		exit 1
	fi
}

detect_os() {
	case "$(uname -s)" in
		Linux) echo "linux" ;;
		Darwin) echo "darwin" ;;
		*) echo "error: unsupported OS: $(uname -s)" >&2; exit 1 ;;
	esac
}

detect_arch() {
	case "$(uname -m)" in
		x86_64 | amd64) echo "amd64" ;;
		arm64 | aarch64) echo "arm64" ;;
		*) echo "error: unsupported architecture: $(uname -m)" >&2; exit 1 ;;
	esac
}

download() {
	url="$1"
	output="$2"
	if command -v curl >/dev/null 2>&1; then
		curl -fsSL "$url" -o "$output"
	elif command -v wget >/dev/null 2>&1; then
		wget -q "$url" -O "$output"
	else
		echo "error: missing required command: curl or wget" >&2
		exit 1
	fi
}

need uname
need tar
need mktemp
need install

os="$(detect_os)"
arch="$(detect_arch)"
asset="${binary}_${os}_${arch}.tar.gz"

if [ "$version" = "latest" ]; then
	url="https://github.com/${repo}/releases/latest/download/${asset}"
else
	url="https://github.com/${repo}/releases/download/${version}/${asset}"
fi

tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT HUP INT TERM

echo "Downloading ${repo} ${version} for ${os}/${arch}..."
download "$url" "$tmpdir/$asset"

tar -xzf "$tmpdir/$asset" -C "$tmpdir"
mkdir -p "$bindir"
install -m 0755 "$tmpdir/$binary" "$bindir/$binary"

echo "Installed $binary to $bindir/$binary"
case ":$PATH:" in
	*":$bindir:"*) ;;
	*) echo "Add $bindir to your PATH to run $binary from anywhere." ;;
esac
