#!/bin/sh
# cans installer — Apple Silicon macOS. Drops a binary on PATH, then: cans doctor
set -eu

REPO="veronica-agent/cans"
DEST="${CANS_BIN:-}"

die() {
	echo "cans-install: $*" >&2
	exit 1
}

os=$(uname -s)
arch=$(uname -m)
[ "$os" = "Darwin" ] || die "macOS only"
[ "$arch" = "arm64" ] || die "Apple Silicon only (got $arch)"

if [ -z "$DEST" ]; then
	if [ -d /opt/homebrew/bin ] && [ -w /opt/homebrew/bin ]; then
		DEST=/opt/homebrew/bin
	elif [ -d /usr/local/bin ] && [ -w /usr/local/bin ]; then
		DEST=/usr/local/bin
	else
		DEST="${HOME}/.local/bin"
	fi
fi

mkdir -p "$DEST"

github_curl() {
	url=$1
	dest=$2
	hdr=""
	if [ -n "${GITHUB_TOKEN:-}" ]; then
		hdr="Authorization: Bearer ${GITHUB_TOKEN}"
	elif command -v gh >/dev/null 2>&1; then
		tok=$(gh auth token 2>/dev/null || true)
		if [ -n "$tok" ]; then
			hdr="Authorization: Bearer ${tok}"
		fi
	fi
	if [ -n "$hdr" ]; then
		curl -fsSL -H "$hdr" -o "$dest" "$url"
	else
		curl -fsSL -o "$dest" "$url"
	fi
}

pick_archive() {
	dir=$1
	for f in "$dir"/*.tar.gz "$dir"/*.tgz; do
		if [ -f "$f" ]; then
			echo "$f"
			return 0
		fi
	done
	return 1
}

download_release() {
	tmp=$(mktemp -d -t cans-install)
	# shellcheck disable=SC2064
	trap 'rm -rf "$tmp"' EXIT

	if command -v gh >/dev/null 2>&1; then
		gh release download -R "$REPO" -p "*darwin_arm64.tar.gz" -D "$tmp"
	else
		api="$tmp/latest.json"
		github_curl "https://api.github.com/repos/${REPO}/releases/latest" "$api"
		tag=$(sed -n 's/.*"tag_name": *"\([^"]*\)".*/\1/p' "$api" | head -n 1)
		[ -n "$tag" ] || return 1
		ver=${tag#v}
		github_curl "https://github.com/${REPO}/releases/download/${tag}/cans_${ver}_darwin_arm64.tar.gz" "$tmp/cans.tar.gz"
	fi

	archive=$(pick_archive "$tmp") || return 1
	tar -xzf "$archive" -C "$tmp"
	bin=""
	if [ -f "$tmp/cans" ]; then
		bin="$tmp/cans"
	else
		bin=$(find "$tmp" -type f -name cans | head -n 1)
	fi
	[ -n "$bin" ] && [ -f "$bin" ] || return 1
	chmod 755 "$bin"
	mv "$bin" "${DEST}/cans"
	trap - EXIT
	rm -rf "$tmp"
}

install_with_go() {
	command -v go >/dev/null 2>&1 || return 1
	go install "github.com/${REPO}/cmd/cans@main"
	gobin=$(go env GOBIN)
	[ -n "$gobin" ] || gobin="$(go env GOPATH)/bin"
	[ -f "${gobin}/cans" ] || return 1
	if [ "$gobin" != "$DEST" ]; then
		cp "${gobin}/cans" "${DEST}/cans"
		chmod 755 "${DEST}/cans"
	fi
}

if download_release; then
	:
elif install_with_go; then
	echo "cans-install: installed from source (no GitHub release yet)" >&2
else
	die "need a GitHub release (gh auth) or Go on PATH"
fi

echo "installed ${DEST}/cans"
case ":${PATH}:" in
*":${DEST}:"*) ;;
*)
	echo "add ${DEST} to PATH" >&2
	;;
esac

echo "next: unpack qwen3-tts-native into ~/.cans/native, then cans doctor"
