# VANTUN Network Condition Test Scenarios

This document outlines the network condition test scenarios for evaluating the performance and stability of the VANTUN protocol using Charles Proxy.

## Test Environment

- **Local Machine**: macOS with Charles Proxy installed
- **VANTUN Server**: Running on localhost:4242
- **VANTUN Client**: Running on the same machine, connecting through Charles Proxy
- **Test Data**: 1024 bytes per request
- **Test Duration**: 1 minute per scenario
- **Request Interval**: 1 second

## Network Condition Scenarios

### Scenario 1: Good Network Conditions
- **Latency**: 10ms
- **Bandwidth**: 100 Mbps
- **Packet Loss**: 0%
- **Description**: Simulates a high-quality network connection

### Scenario 2: Typical WiFi Network
- **Latency**: 50ms
- **Bandwidth**: 50 Mbps
- **Packet Loss**: 0.1%
- **Description**: Simulates a typical home or office WiFi connection

### Scenario 3: Mobile Network (4G)
- **Latency**: 100ms
- **Bandwidth**: 20 Mbps
- **Packet Loss**: 0.5%
- **Description**: Simulates a 4G mobile network connection

### Scenario 4: Poor Network Conditions
- **Latency**: 300ms
- **Bandwidth**: 5 Mbps
- **Packet Loss**: 2%
- **Description**: Simulates a poor network connection, such as edge or congested WiFi

### Scenario 5: Very Poor Network Conditions
- **Latency**: 500ms
- **Bandwidth**: 1 Mbps
- **Packet Loss**: 5%
- **Description**: Simulates a very poor network connection, such as edge or heavily congested network

## Test Execution Plan

1. Start VANTUN server on localhost:4242
2. Configure Charles Proxy with the network conditions for the current scenario
3. Run the VANTUN client through Charles Proxy
4. Execute the stability test for 1 minute
5. Record the results (requests sent, errors, latency, throughput)
6. Repeat for all scenarios

## Metrics to Collect

- Total requests sent
- Successful responses received
- Error rate
- Average latency per request
- Throughput (bytes/second)
- Connection stability

## Comparison with TCP

For each scenario, we will also run a similar test using raw TCP to compare performance:
- Same test data size and frequency
- Same network conditions applied via Charles Proxy
- Record and compare metrics with VANTUN