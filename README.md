# Treeship

*NOTE*:
Project Under Construction. 

Treeship is a Kubernetes agent and server system that performs actions on behalf of users in Kubernetes environments. It uses gRPC for efficient communication between server components and distributed agents.

## Overview
Treeship consists of two main components:

Server: Central coordination point that receives user requests and delegates actions
Agents: Distributed components that execute operations within Kubernetes clusters

## Features
- Secure delegation of Kubernetes operations
- Communication using gRPC bidirectional streaming
- Scalable architecture supporting multiple agents and clusters
- TODO: Add impersonate the client for K8s.
- TODO: Role-based access control for operation permissions

## Architecture
```
User Request → Treeship Server → gRPC → Treeship Agent → Kubernetes API
```
