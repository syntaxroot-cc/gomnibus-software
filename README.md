# gomnibus-software

A library of reusable software definitions for [gomnibus](https://github.com/syntaxroot-cc/gomnibus) — the Go-based Omnibus build system. Each YAML file describes how to fetch, build, and install an upstream software component into a gomnibus install root.

## Usage

Point your project's `gomnibus.yaml` at this library by adding it to `software_dirs`:

```yaml
# gomnibus.yaml
software_dirs:
  - /path/to/gomnibus-software/config/software
  - config/software   # your project's own definitions
```

Or pin a specific commit/tag by cloning it as a git submodule:

```bash
git submodule add https://github.com/syntaxroot-cc/gomnibus-software vendor/gomnibus-software
```

```yaml
# gomnibus.yaml
software_dirs:
  - vendor/gomnibus-software/config/software
  - config/software
```

Then declare the software in your project:

```yaml
# config/projects/myapp.yaml
dependencies:
  - zlib
  - openssl
  - myapp-software
```

## Available software

| Name | Default version | License |
|------|----------------|---------|
| [bzip2](config/software/bzip2.yaml) | 1.0.8 | BSD-2-Clause |
| [cacerts](config/software/cacerts.yaml) | 2022-10-11 | MPL-2.0 |
| [curl](config/software/curl.yaml) | 8.14.1 | MIT |
| [libffi](config/software/libffi.yaml) | 3.4.6 | MIT |
| [libyaml](config/software/libyaml.yaml) | 0.2.5 | MIT |
| [ncurses](config/software/ncurses.yaml) | 6.4 | MIT |
| [openssl](config/software/openssl.yaml) | 3.3.1 | OpenSSL |
| [zlib](config/software/zlib.yaml) | 1.3.1 | Zlib |

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md). Enable the pre-commit gofmt hook after cloning:

```bash
git config core.hooksPath .githooks
```

To add a new definition, create `config/software/<name>.yaml` following the [gomnibus software YAML schema](https://github.com/syntaxroot-cc/gomnibus/blob/main/internal/software/definition.go). The CI suite validates that all definitions parse correctly, all URL sources have checksums, and all declared dependencies exist in the library.
