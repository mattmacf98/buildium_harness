# Buildium Test Harness

Go test harness helpers for creating interactive tutorials with **Buildium**.

This package provides utilities to streamline the creation of step-by-step tutorials using Go's testing framework. It enables you to build structured, verifiable learning experiences where each step can be tested and validated.

## Installation

```bash
go get github.com/builium-org/buildium_harness
```

## Package Overview

| Package | Description |
|---------|-------------|
| `logger` | Colorized logging with step tracking and log collection |
| `meta` | Project metadata parsing from `meta.json` |
| `supabase` | Supabase client for authentication and run reporting |
| `testcli` | Test runner for CLI-based tutorials |
| `testserver` | Test runner for server-based tutorials |

## Usage

### CLI Tutorials

For tutorials that test command-line applications:

```go
package main

import "github.com/builium-org/buildium_harness/testcli"

func main() {
    testcli.RunCliTest([]func(config *testcli.CliTestConfig) error{
        Step1_BasicOutput,
        Step2_FlagParsing,
        Step3_FileHandling,
    })
}

func Step1_BasicOutput(config *testcli.CliTestConfig) error {
    config.Logger.LogTitle("Basic Output")
    config.Logger.LogInfo("Testing that the CLI outputs 'Hello, World!'")
    
    // Run the user's executable and validate output
    // config.Executable contains the path to the user's compiled binary
    
    return nil // or return an error if the test fails
}
```

### Server Tutorials

For tutorials that test HTTP servers or long-running processes:

```go
package main

import "github.com/builium-org/buildium_harness/testserver"

func main() {
    testserver.RunServerTest([]func(config *testserver.ServerTestConfig) error{
        Step1_ServerStarts,
        Step2_GetEndpoint,
        Step3_PostEndpoint,
    })
}

func Step1_ServerStarts(config *testserver.ServerTestConfig) error {
    config.Logger.LogTitle("Server Starts")
    config.Logger.LogInfo("Testing that the server starts on port 8080")
    
    // config.Server is automatically started before each test
    // and stopped after each test
    
    return nil
}
```

## Project Configuration

Each tutorial project requires a `meta.json` file:

```json
{
  "stage": 2,
  "entrypoint": "app",
  "projectId": "abc123-def456"
}
```

| Field | Description |
|-------|-------------|
| `stage` | Current stage the user is on (0-indexed). Tests only run up to this stage. |
| `entrypoint` | Name of the compiled executable to test |
| `projectId` | Unique identifier for tracking progress in Supabase |

## Environment Variables

| Variable | Description |
|----------|-------------|
| `BUILDIUM_EMAIL` | User's Buildium account email |
| `BUILDIUM_PASSWORD` | User's Buildium account password |
| `ENVIRONMENT` | Set to `PROD` for production, `BUILDING` to skip reporting, or leave empty for local development |
| `SERVER_STARTUP_TIME` | Milliseconds to wait for server to start (default: 500) |

## Logger API

The `Logger` provides structured logging with colored terminal output:

```go
logger.LogTitle("Test Name")      // Header for test section
logger.LogSuccess("Passed!")      // Green success message
logger.LogInfo("Checking...")     // Blue info message
logger.LogError("Failed!")        // Red error message
logger.LogClientCode(output)      // Yellow output from user's code
```

All logs are collected and can be retrieved with `logger.GetAllLogs()` for reporting.

## Running Tests

```bash
# Run with a path to the user's project
go run main.go -path=/path/to/user/project

# The path should contain:
# - meta.json (configuration)
# - The compiled executable (specified in meta.json entrypoint)
```

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    Test Harness                         │
├─────────────────────────────────────────────────────────┤
│  testcli / testserver                                   │
│  ├── Reads meta.json for configuration                  │
│  ├── Runs test steps sequentially up to current stage   │
│  ├── Logs results with colorized output                 │
│  └── Reports progress to Supabase                       │
├─────────────────────────────────────────────────────────┤
│  logger                                                 │
│  ├── Step tracking                                      │
│  ├── Colorized console output                           │
│  └── Log collection for reporting                       │
├─────────────────────────────────────────────────────────┤
│  supabase                                               │
│  ├── User authentication                                │
│  └── Project run submission                             │
└─────────────────────────────────────────────────────────┘
```
