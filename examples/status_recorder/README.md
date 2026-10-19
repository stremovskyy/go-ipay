# Status with a GORM recorder

This example keeps its SQLite/GORM dependencies in a separate module so they are
not required to build the SDK. The SDK still supports Go 1.24; this example needs
Go 1.25 or newer for the patched `golang.org/x/text` dependency.

From the repository root:

```sh
go -C examples/status_recorder test -race ./...
go -C examples/status_recorder run .
```

The module replaces `github.com/stremovskyy/go-ipay` with the local repository.
Provide the existing `IPAY_*` environment variables before running the example.
For an environment file, set `IPAY_EXAMPLES_ENV_FILE` to its absolute path.

The test uses a temporary SQLite database and a local HTTP transport stub; it
does not contact iPay.
