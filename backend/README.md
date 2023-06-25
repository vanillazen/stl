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
go run ./cmd/stl/main.go
[INF] 2023/06/25 22:14:01 stl starting...
[INF] 2023/06/25 22:14:01 stl started!
[INF] 2023/06/25 22:14:01 list-repo start
[INF] 2023/06/25 22:14:01 http-server started listening at localhost:8080
[INF] 2023/06/25 22:14:01 sqlite-db database connected!
```

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

Notes
This project is intended to serve as a playground for experimentation and learning. It aims to develop a generator that simplifies the generation of similar projects based on the structure and patterns considered convenient within this implementation.

## License

This project is licensed under the MIT License. Feel free to use and modify it as per the terms of the license.
