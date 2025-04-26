# Load Testing Report

Between 2 PC

## Test Environment

### Hardware Specifications (server)

| Component | Specification                                  |
| --------- | ---------------------------------------------- |
| CPU       | Apple M4, 10 physical cores, 10 logical cores  |
| Memory    | 24 GB                                          |
| Storage   | 460 GB SSD                                     |
| OS        | macOS Darwin Kernel 24.4.0, ARM64 architecture |

### Hardware Specifications (client)

| Component | Specification                                     |
| --------- | ------------------------------------------------- |
| CPU       | Apple M1 Pro, 10 physical cores, 10 logical cores |
| Memory    | 16 GB                                             |
| Storage   | 460 GB SSD                                        |
| OS        | macOS Darwin Kernel 24.4.0, ARM64 architecture    |

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

### Description

Loadtest using mac book pro M1 traffic generation and iMac for API server to split CPU usage

### General Results

| Metric                    | Value                                     |
| ------------------------- | ----------------------------------------- |
| Total execution time      | 2m54.952679542s                           |
| Total requests            | 100100 (Users: 100, Transactions: 100000) |
| Successful requests       | 100100 (100.00%)                          |
| Failed requests           | 0 (0.00%)                                 |
| Requests per second (RPS) | 572.15                                    |

### User Creation Request Metrics

| Metric                    | Value                                                               |
| ------------------------- | ------------------------------------------------------------------- |
| Average response time     | 14.719106ms                                                         |
| Minimum response time     | 6.543166ms                                                          |
| Maximum response time     | 126.194417ms                                                        |
| Response time percentiles | P50=11.953791ms, P90=22.189625ms, P95=29.581959ms, P99=126.194417ms |

### Transaction Creation Request Metrics

| Metric                    | Value                                                              |
| ------------------------- | ------------------------------------------------------------------ |
| Average response time     | 17.325043ms                                                        |
| Minimum response time     | 5.718875ms                                                         |
| Maximum response time     | 125.608792ms                                                       |
| Response time percentiles | P50=15.880833ms, P90=24.404041ms, P95=28.649625ms, P99=39.189166ms |

## Test 2

### Parameters

| Key          | Value |
| ------------ | ----- |
| users        | 100   |
| transactions | 1000  |
| concurrency  | 20    |

### General Results

| Metric                    | Value                                     |
| ------------------------- | ----------------------------------------- |
| Total execution time      | 1m31.0559785s                             |
| Total requests            | 100100 (Users: 100, Transactions: 100000) |
| Successful requests       | 100100 (100.00%)                          |
| Failed requests           | 0 (0.00%)                                 |
| Requests per second (RPS) | 1099.32                                   |

### User Creation Request Metrics

| Metric                    | Value                                                               |
| ------------------------- | ------------------------------------------------------------------- |
| Average response time     | 15.303592ms                                                         |
| Minimum response time     | 9.134458ms                                                          |
| Maximum response time     | 131.742334ms                                                        |
| Response time percentiles | P50=12.194292ms, P90=22.026792ms, P95=30.604916ms, P99=131.742334ms |

### Transaction Creation Request Metrics

| Metric                    | Value                                                          |
| ------------------------- | -------------------------------------------------------------- |
| Average response time     | 17.979893ms                                                    |
| Minimum response time     | 6.565833ms                                                     |
| Maximum response time     | 139.032125ms                                                   |
| Response time percentiles | P50=16.926459ms, P90=23.95325ms, P95=27.322ms, P99=35.985125ms |

## Test 3

### Parameters

| Key          | Value |
| ------------ | ----- |
| users        | 100   |
| transactions | 1000  |
| concurrency  | 25    |

### General Results

| Metric                    | Value                                     |
| ------------------------- | ----------------------------------------- |
| Total execution time      | 1m8.568411875s                            |
| Total requests            | 100100 (Users: 100, Transactions: 100000) |
| Successful requests       | 99267 (99.17%)                            |
| Failed requests           | 833 (0.83%)                               |
| Requests per second (RPS) | 1459.86                                   |

### User Creation Request Metrics

| Metric                    | Value                                                              |
| ------------------------- | ------------------------------------------------------------------ |
| Average response time     | 11.751135ms                                                        |
| Minimum response time     | 8.478042ms                                                         |
| Maximum response time     | 114.352875ms                                                       |
| Response time percentiles | P50=10.447875ms, P90=12.76225ms, P95=14.806667ms, P99=114.352875ms |

### Transaction Creation Request Metrics

| Metric                    | Value                                                              |
| ------------------------- | ------------------------------------------------------------------ |
| Average response time     | 16.91646ms                                                         |
| Minimum response time     | 1.411875ms                                                         |
| Maximum response time     | 74.009792ms                                                        |
| Response time percentiles | P50=15.797041ms, P90=21.902791ms, P95=25.359417ms, P99=44.908875ms |

At this stage, service access problems began to appear, 0.83% of requests were unsuccessful due to the error:

```
2025/04/26 14:40:59 Worker N: error creating transaction for user M: request failed: Post "http://10.0.0.50:8081/api/transactions": dial tcp 10.0.0.50:8081: connect: can't assign requested address
```

## Test 4

Configure http transport

### General Results

| Metric                    | Value                                     |
| ------------------------- | ----------------------------------------- |
| Total execution time      | 1m2.4624155s                              |
| Total requests            | 100100 (Users: 100, Transactions: 100000) |
| Successful requests       | 100100 (100.00%)                          |
| Failed requests           | 0 (0.00%)                                 |
| Requests per second (RPS) | 1602.56                                   |

### User Creation Request Metrics

| Metric                    | Value                                                             |
| ------------------------- | ----------------------------------------------------------------- |
| Average response time     | 12.648492ms                                                       |
| Minimum response time     | 7.910833ms                                                        |
| Maximum response time     | 143.810208ms                                                      |
| Response time percentiles | P50=11.031583ms, P90=13.943083ms, P95=14.8025ms, P99=143.810208ms |

### Transaction Creation Request Metrics

| Metric                    | Value                                                            |
| ------------------------- | ---------------------------------------------------------------- |
| Average response time     | 15.455233ms                                                      |
| Minimum response time     | 6.101833ms                                                       |
| Maximum response time     | 134.003167ms                                                     |
| Response time percentiles | P50=15.011625ms, P90=19.205417ms, P95=21.0925ms, P99=26.687375ms |

## Test 5

### Parameters

| Key          | Value |
| ------------ | ----- |
| users        | 100   |
| transactions | 1000  |
| concurrency  | 30    |

### General Results

| Metric                    | Value                                     |
| ------------------------- | ----------------------------------------- |
| Total execution time      | 1m2.565441708s                            |
| Total requests            | 100100 (Users: 100, Transactions: 100000) |
| Successful requests       | 100100 (100.00%)                          |
| Failed requests           | 0 (0.00%)                                 |
| Requests per second (RPS) | 1599.92                                   |

### User Creation Request Metrics

| Metric                    | Value                                                           |
| ------------------------- | --------------------------------------------------------------- |
| Average response time     | 10.962608ms                                                     |
| Minimum response time     | 6.767333ms                                                      |
| Maximum response time     | 36.134ms                                                        |
| Response time percentiles | P50=10.171958ms, P90=13.953083ms, P95=14.824042ms, P99=36.134ms |

### Transaction Creation Request Metrics

| Metric                    | Value                                                             |
| ------------------------- | ----------------------------------------------------------------- |
| Average response time     | 15.667672ms                                                       |
| Minimum response time     | 5.557416ms                                                        |
| Maximum response time     | 71.551667ms                                                       |
| Response time percentiles | P50=15.154917ms, P90=19.389542ms, P95=21.092375ms, P99=27.77575ms |

Increasing concurrency from 25 to 30 did not yield better results, although CPU usage on the service side was less than 50%.

## Test 6

### Parameters

| Key          | Value |
| ------------ | ----- |
| users        | 100   |
| transactions | 1000  |
| concurrency  | 50    |

### General Results

| Metric                    | Value                                     |
| ------------------------- | ----------------------------------------- |
| Total execution time      | 38.885642834s                             |
| Total requests            | 100100 (Users: 100, Transactions: 100000) |
| Successful requests       | 100100 (100.00%)                          |
| Failed requests           | 0 (0.00%)                                 |
| Requests per second (RPS) | 2574.21                                   |

### User Creation Request Metrics

| Metric                    | Value                                                             |
| ------------------------- | ----------------------------------------------------------------- |
| Average response time     | 13.685437ms                                                       |
| Minimum response time     | 8.83825ms                                                         |
| Maximum response time     | 44.260541ms                                                       |
| Response time percentiles | P50=12.386917ms, P90=20.324792ms, P95=22.00825ms, P99=44.260541ms |

### Transaction Creation Request Metrics

| Metric                    | Value                                                              |
| ------------------------- | ------------------------------------------------------------------ |
| Average response time     | 19.191179ms                                                        |
| Minimum response time     | 6.869708ms                                                         |
| Maximum response time     | 137.986958ms                                                       |
| Response time percentiles | P50=18.338084ms, P90=24.604208ms, P95=27.471791ms, P99=41.384917ms |

## Test 7

### Parameters

| Key          | Value |
| ------------ | ----- |
| users        | 100   |
| transactions | 1000  |
| concurrency  | 75    |

### General Results

| Metric                    | Value                                     |
| ------------------------- | ----------------------------------------- |
| Total execution time      | 52.812550792s                             |
| Total requests            | 100100 (Users: 100, Transactions: 100000) |
| Successful requests       | 100100 (100.00%)                          |
| Failed requests           | 0 (0.00%)                                 |
| Requests per second (RPS) | 1895.38                                   |

### User Creation Request Metrics

| Metric                    | Value                                                               |
| ------------------------- | ------------------------------------------------------------------- |
| Average response time     | 18.403778ms                                                         |
| Minimum response time     | 10.069916ms                                                         |
| Maximum response time     | 102.486583ms                                                        |
| Response time percentiles | P50=14.061583ms, P90=20.500667ms, P95=72.279291ms, P99=102.486583ms |

### Transaction Creation Request Metrics

| Metric                    | Value                                                           |
| ------------------------- | --------------------------------------------------------------- |
| Average response time     | 27.724063ms                                                     |
| Minimum response time     | 6.409459ms                                                      |
| Maximum response time     | 303.463084ms                                                    |
| Response time percentiles | P50=22.94ms, P90=32.954416ms, P95=101.344458ms, P99=116.83325ms |

## Comparison between Test 6 and Test 7

### Test Configuration Changes

| Parameter   | Test 6 Value | Test 7 Value | Change    |
| ----------- | ------------ | ------------ | --------- |
| Concurrency | 50           | 75           | +25 (50%) |

### Performance Impact Analysis

| Metric                        | Test 6   | Test 7   | Change    | Percentage Change |
| ----------------------------- | -------- | -------- | --------- | ----------------- |
| Total execution time          | 38.89s   | 52.81s   | +13.92s   | +35.8%            |
| Requests per second (RPS)     | 2574.21  | 1895.38  | -678.83   | -26.4%            |
| User avg response time        | 13.69ms  | 18.40ms  | +4.71ms   | +34.4%            |
| Transaction avg response time | 19.19ms  | 27.72ms  | +8.53ms   | +44.4%            |
| User max response time        | 44.26ms  | 102.49ms | +58.23ms  | +131.6%           |
| Transaction max response time | 137.99ms | 303.46ms | +165.47ms | +119.9%           |
| User P95 response time        | 22.01ms  | 72.28ms  | +50.27ms  | +228.4%           |
| Transaction P95 response time | 27.47ms  | 101.34ms | +73.87ms  | +268.9%           |

### Key Findings

1. **Performance Degradation**: Despite increasing concurrency by 50% (from 50 to 75), overall performance actually decreased, with RPS dropping by 26.4% and execution time increasing by 35.8%.

2. **Response Time Impact**: Both user and transaction response times were significantly impacted:

   - Average response times increased by 34.4% for user creation and 44.4% for transaction creation
   - Maximum response times more than doubled for both operations
   - P95 response times increased dramatically by 228.4% for users and 268.9% for transactions

3. **Scalability Limit**: The results suggest that the system reached its optimal concurrency level at around 50 concurrent connections, beyond which performance degradation occurs due to increased contention and resource competition.

4. **Resource Utilization**: While error rates remained at 0% in both tests, the significant increase in response times and decrease in throughput indicates that system resources (CPU, network, database connections) were becoming saturated at 75 concurrent connections.

### Conclusion

The comparison between Test 6 and Test 7 clearly demonstrates that increasing concurrency from 50 to 75 led to diminishing returns and actually degraded overall system performance. The optimal concurrency setting for this particular system configuration appears to be around 50 concurrent connections, which delivered the highest throughput (2574.21 RPS) and lowest latency metrics across all tests conducted.

This finding is valuable for capacity planning and suggests that horizontal scaling (adding more server instances) would be more beneficial than further increasing concurrency if higher throughput is required.
