# TSMall E-commerce System

## Introduction

TSMall is a modern e-commerce microservice system built with the [Hertz](https://github.com/cloudwego/hertz/) framework. This project aims to provide a scalable, high-performance online shopping platform.

### Key Features

- Microservice Architecture
  - User Service: Authentication, authorization, and user management
  - Product Service: Product catalog and inventory management
  - Order Service: Order processing and management
  - Payment Service: Payment processing and integration
  - Cart Service: Shopping cart management
- Technical Features
  - High-performance HTTP API service based on Hertz framework
  - Nacos integration for service discovery and configuration management
  - Multi-environment deployment support
  - Unified error handling and response format
  - Comprehensive logging and monitoring

## Project Status

Currently in initial development phase:
- ✅ Basic framework setup
- ✅ Infrastructure configuration
- ✅ Development environment setup
- 🚧 Core services implementation (In Progress)
- 📅 API gateway (Planned)
- 📅 Business logic implementation (Planned)

## Directory Structure

```plaintext
TSMall/
├── biz/                    # Business logic directory
│   ├── bizcontext/        # Business context
│   ├── config/            # Configuration management
│   ├── constant/          # Constants definition
│   ├── errno/             # Error codes
│   ├── global_init/       # Global initialization
│   ├── handler/           # Request handlers
│   ├── router/            # Route registration
│   ├── service/           # Business service layer
│   └── utils/             # Utility functions
├── cmd/                    # Entry point
│   └── main.go
├── configs/               # Configuration files
│   ├── dev/              # Development environment
│   ├── test/             # Testing environment
│   └── prod/             # Production environment
├── hertz_gen/            # Hertz generated code
├── idl/                  # Interface definition files
├── kit/                  # Toolkit
├── script/               # Scripts
└── template/             # Templates
```

## Quick Start
### Prerequisites
- Go 1.21+
- Nacos (Optional, for configuration center and service discovery)
### Installation
1. Clone the repository
```bash
git clone https://github.com/dgdts/TSMall.git
cd TSMall
```

2. Install dependencies
```bash
go mod tidy
```

### Create New Service
Use the provided script to create a new service:

```bash
./template/new_project.sh <module_name> <service_name>
```

### Build and Run
1. Build the project
```bash
./build_hz.sh
```

2. Run the service
```bash
./script/run.sh
```

## Configuration
The project supports multi-environment configuration in the configs directory:
- dev/ : Development environment
- test/ : Testing environment
- prod/ : Production environment

Main configuration sections:
- Global settings (environment type, service name, etc.)
- Hertz service configuration
- Logging configuration
- Registry center configuration
- Configuration center settings
- Monitoring settings

## API Documentation
API definitions are located in proto files under the idl/ directory.

## Development Guide
### Error Handling
The system defines unified error codes and handling mechanisms in the biz/errno/ directory.

### Response Format
Standard response format:
```json
{
    "code": 0,
    "msg": "success",
    "data": {}
}
```
### Middleware
Integrated middleware components:

- CORS handling
- Recovery mechanism
- Access logging
- Gzip compression
- Pprof profiling

## Deployment
The project includes a Dockerfile for containerized deployment with multi-environment support.

## Health Check
Use script/healthcheck.sh for service health monitoring:

- Port listening check
- Process status check
## Contributing
Issues and Pull Requests are welcome.

