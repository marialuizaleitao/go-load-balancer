# Simple Go Load Balancer

This project implements a simple and thread-safe load balancer in Go, capable of distributing incoming HTTP requests among multiple backend servers using a round-robin algorithm. Each backend server is represented as an interface, allowing flexibility in defining server behavior and health checks.

## Features

- **Round-Robin Load Balancing**: Requests are evenly distributed across available servers.
- **Thread-Safe Implementation**: Uses `sync.Mutex` to ensure safe access to shared resources.
- **Dynamic Server Addition**: Easily extendable to add or remove backend servers.

## How It Works

The load balancer listens on a specified port and forwards incoming requests to backend servers based on a round-robin selection mechanism. It ensures that only healthy servers receive requests by checking their health status through the `Server` interface.

## Usage

### Requirements

- Go 1.16 or later installed on your machine.

### API Endpoints

- **Root Endpoint `/`**: All incoming requests are load balanced among the configured backend servers.
