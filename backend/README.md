# STL - Simple Todo List Backend

STL is a simple example application that serves as a reference implementation for creating Go REST applications. It is designed to focus on microservices for CRUD operations on single resources, providing a clear example of a RESTful microservice implemented in an elegant and straightforward way.

## Key Features

- **CRUD Operations:** The backend module of STL provides a robust implementation for performing CRUD operations on a single resource, specifically a todo list. It handles creating, reading, updating, and deleting todo list items efficiently.
- **RESTful API:** The backend module exposes a RESTful API that follows best practices and conventions for building web services. It ensures that the API endpoints are intuitive and adhere to the principles of REST architecture.
- **Scalable and Maintainable:** The backend module is designed to be scalable and maintainable. It follows modular architecture principles, allowing easy integration with other microservices in the same repository. The codebase is organized and well-structured, promoting code reusability and maintainability.

## Usage

To use the STL backend module, follow these steps:

```shell
$ make run
go run ./cmd/stl/main.go
[INF] 2023/06/28 11:32:19 stl database connected
[INF] 2023/06/28 11:32:19 stl starting...
[INF] 2023/06/28 11:32:19 list-repo started
[INF] 2023/06/28 11:32:19 stl started
[INF] 2023/06/28 11:32:19 Migration 'users' already applied
[INF] 2023/06/28 11:32:19 Migration 'lists' already applied
[INF] 2023/06/28 11:32:19 Migration 'tasks' already applied
[INF] 2023/06/28 11:32:19 list-service start
[INF] 2023/06/28 11:32:19 stl started
[INF] 2023/06/28 11:32:19 http-server listening at localhost:8080
```

```shell
$ make api/get-list 
./scripts/curl/get-list.sh -h localhost -p 8080 -l cdc7a443-3c6a-431b-b45a-b14735953a19 --pretty
Calling API endpoint: GET http://localhost:8080/api/v1/lists/cdc7a443-3c6a-431b-b45a-b14735953a19
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100   394  100   394    0     0   237k      0 --:--:-- --:--:-- --:--:--  384k
{
  "count": 1,
  "pages": 1,
  "data": {
    "UserID": "",
    "Name": "List 1",
    "Description": "List 1 Description",
    "CreatedAt": "0001-01-01T00:00:00Z",
    "UpdatedAt": "0001-01-01T00:00:00Z",
    "Tasks": [
      {
        "Name": "Task 1",
        "Description": "Task 1 Description",
        "Category": [
          "Category 1"
        ],
        "Tags": [
          "Tag 1",
          "Tag 2"
        ],
        "Location": [
          "Location 1"
        ],
        "CreatedAt": "0001-01-01T00:00:00Z",
        "UpdatedAt": "0001-01-01T00:00:00Z"
      }
    ]
  },
  "error": {
    "OK": true
  }
}
```


## API 
* OpenAPI: [http://localhost:8080/api/v1/docs](http://localhost:8080/api/v1/docs)

## Other Modules

In addition to the backend module, the STL project consists of the following sibling modules:

- **API Gateway:** The API Gateway module acts as a central entry point for all incoming requests. It handles routing and request forwarding to the appropriate microservices. This module, along with other sibling modules, is versioned and part of the same repository.
- **BFF (Backend for Frontend):** The BFF module is responsible for aggregating data from multiple microservices and providing a unified API specifically tailored for the frontend client. It is also part of the same repository and versioned alongside other modules.
- **Isomorphic Web Manager:** The isomorphic web manager module enables rendering the web application on both the server-side and client-side, providing better performance and a seamless user experience. This module is part of the same repository and versioned together with other modules.

These modules will be designed to work together as a cohesive system, providing a complete solution for building microservice-based applications.

## Dependency Management and Go Stdlib

In the STL project, we prioritize using Go's standard library as much as possible and avoid relying on external dependencies. We strive to develop our own implementations for networking, database access, and OpenAPI, leveraging the capabilities provided by the Go stdlib.

While we aim to avoid external dependencies, we acknowledge that certain complex functionalities such as GRPC, telemetry, monitoring, and tracing may require the use of proven solutions. In such cases, we will carefully evaluate the trade-offs and opt for established libraries or tools when necessary.

By minimizing external dependencies and maximizing the utilization of Go's standard library, we maintain better control over the project's codebase and reduce potential compatibility issues.

While abstractions are encouraged, the project prioritizes simplicity over complex architectures. The use of interfaces instead of directly using structs is welcomed as it allows for easy implementation changes and facilitates mocking during testing.

## Tests
```shell
$ make test
make -f makefile.test test-selected
make[1]: Entering directory '/home/adrian/Projects/labs/vanillazen/stl/backend'
clear
make -f makefile.test test-migrator
make[2]: Entering directory '/home/adrian/Projects/labs/vanillazen/stl/backend'
go test -v -run TestMigrator -count=1 -timeout=10s internal/infra/migration/sqlite/*.go
=== RUN   TestMigrator
=== RUN   TestMigrator/TestMigrateBase
=== NAME  TestMigrator
    migrator_test.go:170: TestMigrateBase: OK
=== RUN   TestMigrator/TestMigrateAndAgain
=== NAME  TestMigrator
    migrator_test.go:170: TestMigrateAndAgain: OK
=== RUN   TestMigrator/TestRollback1
=== NAME  TestMigrator
    migrator_test.go:170: TestRollback1: OK
=== RUN   TestMigrator/TestRollback2
=== NAME  TestMigrator
    migrator_test.go:170: TestRollback2: OK
=== RUN   TestMigrator/TestRollbackAll
=== NAME  TestMigrator
    migrator_test.go:170: TestRollbackAll: OK
--- PASS: TestMigrator (0.55s)
    --- PASS: TestMigrator/TestMigrateBase (0.15s)
    --- PASS: TestMigrator/TestMigrateAndAgain (0.02s)
    --- PASS: TestMigrator/TestRollback1 (0.11s)
    --- PASS: TestMigrator/TestRollback2 (0.14s)
    --- PASS: TestMigrator/TestRollbackAll (0.13s)
PASS
ok  	command-line-arguments	0.558s
make[2]: Leaving directory '/home/adrian/Projects/labs/vanillazen/stl/backend'
make -f makefile.test test-http
make[2]: Entering directory '/home/adrian/Projects/labs/vanillazen/stl/backend'
go test -v -run TestGetResourceInfo -count=1 -timeout=10s internal/infra/http/*.go
=== RUN   TestGetResourceInfo
=== RUN   TestGetResourceInfo/Valid_URL_parts
=== RUN   TestGetResourceInfo/Invalid_URL_parts_count
=== RUN   TestGetResourceInfo/Invalid_URL_ID
--- PASS: TestGetResourceInfo (0.00s)
    --- PASS: TestGetResourceInfo/Valid_URL_parts (0.00s)
    --- PASS: TestGetResourceInfo/Invalid_URL_parts_count (0.00s)
    --- PASS: TestGetResourceInfo/Invalid_URL_ID (0.00s)
PASS
ok  	command-line-arguments	0.003s
make[2]: Leaving directory '/home/adrian/Projects/labs/vanillazen/stl/backend'
make[1]: Leaving directory '/home/adrian/Projects/labs/vanillazen/stl/backend'
```

## Notes
This project is intended to serve as a playground for experimentation and learning. It aims to develop a generator that simplifies the generation of similar projects based on the structure and patterns considered convenient within this implementation.

Additional updates [here](docs/notes.md)

## License

This project is licensed under the MIT License. Feel free to use and modify it as per the terms of the license.
