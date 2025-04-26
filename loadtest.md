# Load Testing Report

## Test Environment

### Hardware Specifications

| Component | Specification                                  |
| --------- | ---------------------------------------------- |
| CPU       | Apple M4, 10 physical cores, 10 logical cores  |
| Memory    | 24 GB (25769803776 bytes)                      |
| Storage   | 460 GB SSD (10 GB used, 283 GB available)      |
| OS        | macOS Darwin Kernel 24.4.0, ARM64 architecture |

### Software Specifications

| Component  | Specification                                  |
| ---------- | ---------------------------------------------- |
| Go Version | go1.24.2 darwin/arm64                          |
| Database   | PostgreSQL 16 (Docker container, alpine image) |

### Parameters

| Key          | Value |
| ------------ | ----- |
| users        | 100   |
| transactions | 1000  |
| concurrency  | 10    |

## Test 1

### General Results

| Metric                    | Value                                     |
| ------------------------- | ----------------------------------------- |
| Total execution time      | 45.659215417s                             |
| Total requests            | 100100 (Users: 100, Transactions: 100000) |
| Successful requests       | 100100 (100.00%)                          |
| Failed requests           | 0 (0.00%)                                 |
| Requests per second (RPS) | 2192.33                                   |

### User Creation Request Metrics

| Metric                    | Value                                                       |
| ------------------------- | ----------------------------------------------------------- |
| Average response time     | 954.993µs                                                   |
| Minimum response time     | 401.833µs                                                   |
| Maximum response time     | 8.924125ms                                                  |
| Response time percentiles | P50=557.75µs, P90=1.274542ms, P95=3.23475ms, P99=8.924125ms |

### Transaction Creation Request Metrics

| Metric                    | Value                                                            |
| ------------------------- | ---------------------------------------------------------------- |
| Average response time     | 4.529362ms                                                       |
| Minimum response time     | 695.625µs                                                        |
| Maximum response time     | 110.621ms                                                        |
| Response time percentiles | P50=3.538459ms, P90=8.998666ms, P95=10.440875ms, P99=16.322917ms |

## Test 2

### Applied Performance Optimizations

- Fast JSON serialization with JsonIter
- Optimized database connection pool
- Middleware for JSON

### General Results

| Metric                    | Value                                     |
| ------------------------- | ----------------------------------------- |
| Total execution time      | 29.456336667s                             |
| Total requests            | 100100 (Users: 100, Transactions: 100000) |
| Successful requests       | 100100 (100.00%)                          |
| Failed requests           | 0 (0.00%)                                 |
| Requests per second (RPS) | 3398.25                                   |

### User Creation Request Metrics

| Metric                    | Value                                                      |
| ------------------------- | ---------------------------------------------------------- |
| Average response time     | 859.102µs                                                  |
| Minimum response time     | 556.209µs                                                  |
| Maximum response time     | 9.956042ms                                                 |
| Response time percentiles | P50=688.542µs, P90=1.143708ms, P95=1.314ms, P99=9.956042ms |

### Transaction Creation Request Metrics

| Metric                    | Value                                                       |
| ------------------------- | ----------------------------------------------------------- |
| Average response time     | 2.928914ms                                                  |
| Minimum response time     | 526.792µs                                                   |
| Maximum response time     | 63.350084ms                                                 |
| Response time percentiles | P50=2.811ms, P90=3.399958ms, P95=3.683333ms, P99=5.278459ms |

## Performance Improvement Summary

- **Total execution time reduction**: 35.5%
- **RPS increase**: 55.0%
- **Transaction response time improvement**: from 4.53ms to 2.93ms (35.3%)
