# atlantis-configtest

[![CI](https://github.com/winebarrel/atlantis-configtest/actions/workflows/ci.yml/badge.svg)](https://github.com/winebarrel/atlantis-configtest/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/winebarrel/atlantis-configtest/graph/badge.svg)](https://codecov.io/gh/winebarrel/atlantis-configtest)
[![AI Generated](https://img.shields.io/badge/AI%20Generated-Claude-orange?logo=anthropic)](https://claude.ai/claude-code)

Check the structure of Atlantis configuration files without starting a server.

Atlantis parses `repos.yaml` only when the server boots, and `atlantis.yaml` only when it runs a plan for a pull request. A typo in either one shows up as a failed deploy or as a failed comment on someone's pull request. This command runs the same parser ahead of time, so a CI job can catch it first.

## Installation

Download a binary from the [releases page](https://github.com/winebarrel/atlantis-configtest/releases), or build it yourself:

```sh
go install github.com/winebarrel/atlantis-configtest/cmd/atlantis-configtest@latest
```

## Usage

```
Usage: atlantis-configtest <command> [flags]

Validate the structure of Atlantis configuration files.

Flags:
  -h, --help       Show context-sensitive help.
      --version

Commands:
  server-repo <file> ... [flags]
    Validate server-side repo config files (repos.yaml).

  repo <file> ... [flags]
    Validate repo-level config files (atlantis.yaml).

Run "atlantis-configtest <command> --help" for more information on a command.
```

Every file is reported on stderr, and the exit status is 1 if any of them fails:

```sh
$ atlantis-configtest server-repo repos.yaml
repos.yaml: ok
$ echo $?
0

$ atlantis-configtest server-repo repos.yaml broken.yaml
repos.yaml: ok
broken.yaml: repos: (0: (allowed_overrides: "nonexistent_key" is not a valid override, only "plan_requirements", ... are supported.).).
atlantis-configtest: error: 1 of 2 file(s) invalid
$ echo $?
1
```

Every file is checked even after one of them fails, so a single broken file does not hide the rest.
