# TODO

- Fix missing chi router wiring in `cmd/store/main.go:26`; import `github.com/go-chi/chi/v5`, route through the router, or remove unused code so the binary builds.
- Correct the Go toolchain directive in `go.mod:3` to a supported form such as `go 1.21` so `go build ./...` works.
- Update `Taskfile.yml:31` to build from `./cmd/store` (current path `./main.go` fails).
- Initialise `lastPetID` in `api/petstore.go:88` based on the max seeded ID to avoid duplicate IDs when adding pets.
- Guard read paths in `api/petstore.go:34` and `api/petstore.go:120` with the same mutex used for writes (`petsLock`) to eliminate data races.
- Fix the first-party import prefix in `Taskfile.yml:58` to `github.com/neogan74/` so `gci` respects module thresholds.
