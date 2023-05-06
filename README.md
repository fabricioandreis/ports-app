# Ports application
[![Continuous Delivery pipeline](https://github.com/fabricioandreis/ports-app/actions/workflows/continuous-delivery.yml/badge.svg)](https://github.com/fabricioandreis/ports-app/actions/workflows/continuous-delivery.yml)

This application reads Ports from an input JSON file and saves them into a Redis in-memory database.

## System design
The system contains two main components:
- A Docker container running a Go application that reads a JSON file mounted into it
- A Docker container running a Redis database to store the Ports

I created Docker Compose files to run the system locally on a development environment and also to run acceptance tests on a Continuous Delivery pipeline.

## Architecture
In my experience, one must pay attention to some important concepts and design choices when building software systems in order for them to implement the correct behavior, provide maintainability and be testable.

I tried to showcase the importance of SOLID Principles and Clean Architecture in my design. 

The application respects the Dependency Rule (a form of Dependency Inversion), where high level modules, those related to the policies (e.g.  domain entities, use cases), should not depend on low level modules, which are concerned with details (e.g reading files, writing to a database).

The way I achieved that was by injecting concrete implementations of contracts between packages (interfaces) into the higher level modules.

Moreover, I practiced Test Driven Development during the design of the solution. I prefer this approach because it guides my design, forcing me to understand the core of the problem first and then design the business rules' APIs from the outside in. A nice side effect of this is that there is a good (for the time available) unit test coverage of the high level modules of the application (package `usecase`).

### Decision Records (ADRs)
1. I decided to marshal data with Protocol Buffers when saving to the database in order to reduce the size of each object as compared to a JSON serialization.
2. The application reads the JSON file as an IO stream instead of loading all of the contents into main memory, because the input file could contain millions of rows, which probably would not fit into available memory.
3. The reading of the JSON file and the writing into the database are performed asynchronously in separate go routines that communicate exclusively via channels. I tried to adhere to the Go proverb "[don't communicate by sharing memory, share memory by communicating](https://www.youtube.com/watch?v=PAAkCSZUG1c&t=2m48s)".

4. The build of the statically linked binary executable is performed outside Docker. I choose a **distroless** Docker image to run the application for the following reasons:
   - Since the executable is statically linked, there are no dependencies on the Operating System, so no linux libraries are needed to run it
   - A small base Docker image generates much smaller application containers
   - Containers without unused dependencies are more secure because they reduce the attack surface
5. I used context cancellation to provide a graceful shutdown of the application


### Structure
- `cmd`: contains the entrypoint of the application
- `internal`: contains all the packages that should not be exposed to users of the program
  - `domain`: packages for high level policies. In this case, the Port entity handled by the program. In a regular application, the domain model should not be anemic (no behavior) like this one, because we would add more use cases to the application
  - `usecase`: packages to interact with several objects of the domain in order to accomplish a known use case
  - `contracts`: interfaces that abstract the behavior of low level modules like a database repository
  - `infra`: concrete implementations of the low level modules that interact with the application infrastructure
  - `tests`: acceptance tests (in this case similar to integration tests) that check the behavior of the application from the perspective of the users

## How to use
The following dependencies should be installed in your environment:
- Go 1.20+
- Docker
- Make

### Run
Run the application and its dependencies locally with:
```
make local
```

### Unit tests
Run the unit tests with:
```
make test
```

### Acceptance tests
Build the application and run acceptance tests locally with the following command:
```
make local-acceptance-tests
```
This will run a `docker compose run` command that:
- builds and runs the Docker container of the application
- runs a Redis in-memory database
- builds and runs another Docker container that performs the acceptance tests against the Redis database after the completion of the application


### Simulate SIGTERM
It is possible to simulate a SIGTERM event during the application execution with:
1. Run local instance of Redis on Docker: ```docker run -d --name redis -p 6379:6379 redis:7.0-alpine```
2. Run: ``` make run ```
3. Type: `Ctrl + C`

The application should print logs indicating a graceful shutdown.