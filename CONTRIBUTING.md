# Contributing to gomnibus-software

## Dev setup

```bash
git clone https://github.com/syntaxroot-cc/gomnibus-software
cd gomnibus-software
git config core.hooksPath .githooks   # enable gofmt pre-commit hook
go test ./...
```

## Adding a definition

1. Create `config/software/<name>.yaml`. Required fields: `name`, `default_version`, `license`. If the source has a URL, include a `sha256` checksum.

2. If your definition depends on other software in the library, list them under `dependencies:`. The CI test will fail if a dependency is not found in the library.

3. If multiple versions are supported, add them under `versions:` with per-version `source` and `relative_path` blocks as needed.

4. Run `go test ./...` locally — it catches YAML parse errors, missing checksums, unknown dependencies, and malformed version blocks.

## Schema reference

The YAML schema mirrors gomnibus's `software.Definition` struct. Key fields:

```yaml
name: mylib
default_version: "1.2.3"
license: MIT
license_file: LICENSE        # optional: path inside source tree

source:
  url: "https://example.com/mylib-1.2.3.tar.gz"
  sha256: "abc123..."

relative_path: "mylib-1.2.3"   # extracted directory name (if different from name-version)

dependencies:
  - zlib                        # must exist in this library

versions:
  - version: "1.2.2"
    source:
      url: "https://example.com/mylib-1.2.2.tar.gz"
      sha256: "def456..."
    relative_path: "mylib-1.2.2"

build:
  - configure: []               # runs ./configure --prefix=${install_dir}
  - configure:                  # with extra args (no ${var} expansion in args)
      - "--disable-debug"
  - command: "..."              # shell command; ${install_dir}, ${src_dir}, ${build_dir} are expanded
  - make: []                    # runs make
  - make: [install]             # runs make install
  - mkdir: "${install_dir}/bin"
  - copy:
      src: "..."
      dst: "..."
```

**Note:** Variable substitution (`${install_dir}` etc.) only works in `command:` steps, not in `configure:` or `make:` argument lists. Use a `command:` step when configure args reference install paths.

## Code conventions

- No comments explaining *what* — only non-obvious *why*.
- `gofmt` is enforced by the pre-commit hook and CI.
