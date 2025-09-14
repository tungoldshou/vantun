# Change Log

This document records all significant changes to the VANTUN project, following [Semantic Versioning](https://semver.org/).

## [Unreleased]

### Added Features
- Enhanced documentation system with detailed usage guides and configuration instructions
- Improved demo guide with clearer usage examples

### Fixed
- Fixed inconsistent wording in documentation

### Changed
- Reorganized README document structure to highlight core project advantages
- Optimized technical architecture description to better showcase project highlights

## [1.0.0] - 2025-09-14

### Added Features
- High-performance tunnel protocol implementation based on QUIC for exceptional network transmission performance
- Secure handshake and session negotiation mechanism for connection security
- Multiple logical stream types: interactive, bulk, telemetry, optimized for different business scenarios
- Forward Error Correction (FEC) functionality for enhanced data transmission robustness
- Multipath support for improved performance and reliability using multiple network paths
- Hybrid congestion control algorithm combining underlying QUIC CC with upper-layer token bucket rate limiting
- Pluggable obfuscation module to reduce the risk of being misidentified as illegal traffic
- Telemetry data collection system for comprehensive performance monitoring
- Command-line client/server programs for rapid deployment and ease of use

### Fixed
- Fixed duplicate code and logic errors in main.go to improve code quality
- Improved multipath controller implementation for enhanced multipath transmission stability
- Implemented real telemetry data collection for accurate performance monitoring
- Improved token bucket controller integration with telemetry system for optimized congestion control

### Changed
- Refactored project structure to improve modularity for easier maintenance and extension
- Enhanced test coverage to ensure system stability and reliability
- Improved error handling and resource management for enhanced system robustness

Initial project version with complete implementation of all core features.