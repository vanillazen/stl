# Notes

This file serves as a dump for important notes regarding the project.

## Migrator
A migrator has been implemented, but it is still necessary to develop the corresponding tests and ensure its correct functionality before integrating it with the service. Currently, it is only available to work with SQLite, but a version for PostgreSQL is underway.

The migrator tool relies heavily on the new embedding feature from the 'embed' package, which allows files to be embedded directly into the Go binary. However, due to the design decision of requiring the 'go:embed' directive to be placed in a file within the root module directory, the migrator tool is compelled to move its 'main.go' file from its original location under 'cmd/stl/main.go' to the root of the project. This can be inconvenient, as it deviates from the common practice of organizing the main application file under a specific subdirectory. Fortunately, discussions are underway within the Go community to address this use case and provide better support this spread use case.

Rollback (n to all) still needs tweaking.

Finally, there is room to simplify the migrator operation as well as tidy it up: remove unnecessary error logging and error wrapping. some interface methods used are not really required, too.

