# Zerodha-X: Distributed Trading Engine

A highly concurrent, event-driven trading platform capable of handling 10,000+ mock transactions per second with <50ms p99 latency.

## Technology Stack
- **Languages:** Go (Golang)
- **Databases:** PostgreSQL, Redis
- **Messaging:** Apache Kafka
- **Infrastructure:** Docker, Kubernetes (K8s)
- **Observability:** Prometheus, Grafana
- **Testing:** K6

## Concepts Learned & Implemented
| Concept | One-liner Description | Where it's used |
|---------|-----------------------|-----------------|
| Domain-Driven Design (DDD) | Focusing on the core domain and domain logic | Microservice boundaries |
| Microservices | Architecture that structures an application as a collection of services | System Architecture |
| gRPC vs REST | Choosing between request-response (REST) and high-performance RPC (gRPC) | Service Communication |
| FIFO Matching Algorithm | Price-Time Priority matching for fair execution | Matching Engine |
| Concurrency in Go | Using Mutexes and Goroutines for high-performance safety | Order Book Management |
| Time Complexity | Optimizing data structures for O(log N) or O(1) lookups | Matching Engine |

## Project Structure
- `services/api-gateway`: REST entry point for clients.
- `services/order-service`: Manages order lifecycle.
- `services/matching-engine`: High-performance FIFO order matching.
- `services/wallet-service`: Manages user balances and ledger.
