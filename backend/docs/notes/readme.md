# Notes

This file serves as a dump for important notes regarding the project.

## Migrator
A migrator has been implemented, but it is still necessary to develop the corresponding tests and ensure its correct functionality. 

Currently, it is only available to work with SQLite, but a version for PostgreSQL is underway.

The migrator tool relies heavily on the new embedding feature from the 'embed' package, which allows files to be embedded directly into the Go binary. However, due to the design decision of requiring the 'go:embed' directive to be placed in a file within the root module directory, the migrator tool is compelled to move its 'main.go' file from its original location under 'cmd/stl/main.go' to the root of the project. This can be inconvenient, as it deviates from the common practice of organizing the main application file under a specific subdirectory. Fortunately, discussions are underway within the Go community to address this use case and provide better support this spread use case.

Some error messages can be improved. It is confusing the use of `rollback` in the same sentence when referring to two different issues (Transaction rollback vs migration rollback).

There is also an opportunity to streamline and improve the migrator operation by eliminating unnecessary error logging and wrapping, and by creating a standardized interface that enables its use across various implementations in a more versatile manner.

A command line tool that allows managing migrations and rollbacks would also be desirable.
