# VANTUN Docker Deployment Guide

## 🐳 Quick Start with Docker

Docker provides the easiest way to deploy VANTUN with consistent performance across different platforms.

### Prerequisites
- Docker 20.10+
- Docker Compose 1.29+ (optional)
- 64MB+ available RAM
- 50MB+ disk space

## 📦 Docker Images

### Official Images
```bash
# Latest stable release
docker pull tungoldshou/vantun:latest

# Specific version
docker pull tungun/vantun:v1.0.0

# Development build
docker pull tungoldshou/vantun:dev
```

### Image Variants
- `latest`: Stable release (recommended)
- `alpine`: Lightweight Alpine-based image
- `distroless`: Minimal security-focused image
- `dev`: Latest development build

## 🚀 Basic Deployment

### 1. Simple Server Setup

```bash
# Create configuration directory
mkdir -p /etc/vantun

# Server configuration file (/etc/vantun/server.json)
cat > /etc/vantun/server.json << 'EOF'
{
  "server": true,
  "address": "0.0.0.0:4242",
  "log_level": "info",
  "multipath": true,
  "obfs": true,
  "fec_data": 10,
  "fec_parity": 3
}
EOF

# Run server
docker run -d \
  --name vantun-server \
  --restart unless-stopped \
  --network host \
  -v /etc/vantun:/etc/vantun \
  tungoldshou/vantun:latest \
  server --config /etc/vantun/server.json
```

### 2. Simple Client Setup

```bash
# Create configuration directory
mkdir -p /etc/vantun

# Client configuration file (/etc/vantun/client.json)
cat > /etc/vantun/client.json << 'EOF'
{
  "server": false,
  "address": "your-server.com:4242",
  "log_level": "info",
  "multipath": true,
  "obfs": true,
  "fec_data": 10,
  "fec_parity": 3,
  "local_addr": "127.0.0.1:1080",
  "socks5": true
}
EOF

# Run client
docker run -d \
  --name vantun-client \
  --restart unless-stopped \
  -p 1080:1080 \
  -v /etc/vantun:/etc/vantun \
  tungoldshou/vantun:latest \
  client --config /etc/vantun/client.json
```

## 🔧 Advanced Configuration

### Docker Compose Setup

Create `docker-compose.yml`:

```yaml
version: '3.8'

services:
  vantun-server:
    image: tungoldshou/vantun:latest
    container_name: vantun-server
    restart: unless-stopped
    ports:
      - "4242:4242/udp"
      - "4242:4242/tcp"
    volumes:
      - ./config/server.json:/etc/vantun/server.json:ro
      - ./certs:/etc/vantun/certs:ro
    command: server --config /etc/vantun/server.json
    networks:
      - vantun-network
    sysctls:
      - net.core.rmem_max=134217728
      - net.core.wmem_max=134217728
    ulimits:
      nofile:
        soft: 1048576
        hard: 1048576

  vantun-client:
    image: tungoldshou/vantun:latest
    container_name: vantun-client
    restart: unless-stopped
    ports:
      - "1080:1080"
    volumes:
      - ./config/client.json:/etc/vantun/client.json:ro
    command: client --config /etc/vantun/client.json
    networks:
      - vantun-network
    depends_on:
      - vantun-server
    sysctls:
      - net.core.rmem_max=134217728
      - net.core.wmem_max=134217728

networks:
  vantun-network:
    driver: bridge
    ipam:
      config:
        - subnet: 172.20.0.0/16
```

Deploy with:
```bash
docker-compose up -d
```

### Production Setup with TLS

Create production configuration:

```bash
# Create directories
mkdir -p /opt/vantun/{config,certs,logs}

# Generate certificates (or use Let's Encrypt)
openssl req -x509 -newkey rsa:4096 -keyout /opt/vantun/certs/server.key -out /opt/vantun/certs/server.crt -days 365 -nodes -subj "/CN=your-domain.com"

# Server configuration
cat > /opt/vantun/config/server.json << 'EOF'
{
  "server": true,
  "address": "0.0.0.0:4242",
  "log_level": "info",
  "log_file": "/var/log/vantun/server.log",
  "multipath": true,
  "obfs": true,
  "fec_data": 10,
  "fec_parity": 3,
  "tls": {
    "cert": "/etc/vantun/certs/server.crt",
    "key": "/etc/vantun/certs/server.key"
  },
  "performance": {
    "workers": 4,
    "buffer_size": 2097152
  }
}
EOF

# Production Docker command
docker run -d \
  --name vantun-server-prod \
  --restart unless-stopped \
  -p 4242:4242 \
  -p 4242:4242/udp \
  -v /opt/vantun/config:/etc/vantun:ro \
  -v /opt/vantun/certs:/etc/vantun/certs:ro \
  -v /opt/vantun/logs:/var/log/vantun \
  --log-driver json-file \
  --log-opt max-size=10m \
  --log-opt max-file=3 \
  tungoldshou/vantun:latest \
  server --config /etc/vantun/server.json
```

## 🛡️ Security Best Practices

### 1. Non-root Container
```dockerfile
FROM tungoldshou/vantun:latest

# Create non-root user
RUN adduser -D -H -s /sbin/nologin vantun

# Switch to non-root user
USER vantun

ENTRYPOINT ["vantun"]
```

### 2. Network Isolation
```yaml
version: '3.8'

services:
  vantun-server:
    image: tungoldshou/vantun:latest
    cap_drop:
      - ALL
    cap_add:
      - NET_BIND_SERVICE
    security_opt:
      - no-new-privileges:true
    read_only: true
    tmpfs:
      - /tmp:noexec,nosuid,size=100m
```

### 3. Secret Management
```bash
# Use Docker secrets
echo "your-secret-key" | docker secret create vantun_key -

# In docker-compose.yml
secrets:
  vantun_key:
    external: true

services:
  vantun:
    secrets:
      - vantun_key
```

## 📊 Monitoring and Logging

### Basic Logging
```bash
# View logs
docker logs vantun-server

# Follow logs in real-time
docker logs -f vantun-server

# Export logs
docker logs vantun-server > vantun-server.log 2>&1
```

### Advanced Monitoring with Prometheus

Create `docker-compose.monitoring.yml`:

```yaml
version: '3.8'

services:
  vantun:
    image: tungoldshou/vantun:latest
    ports:
      - "4242:4242"
      - "8080:8080"  # Metrics port
    environment:
      - VANTUN_METRICS_ENABLED=true
      - VANTUN_METRICS_PORT=8080
    command: server --config /etc/vantun/server.json --metrics :8080

  prometheus:
    image: prom/prometheus:latest
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml:ro

  grafana:
    image: grafana/grafana:latest
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
```

### Health Checks
```dockerfile
# Custom Dockerfile with health check
FROM tungoldshou/vantun:latest

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD nc -z localhost 4242 || exit 1
```

## 🔧 Performance Optimization

### Kernel Parameters
```bash
# Create optimization script
cat > optimize-docker.sh << 'EOF'
#!/bin/bash
# Docker host optimization for VANTUN

echo "Optimizing kernel parameters..."

# Network performance
echo "net.core.rmem_max = 134217728" >> /etc/sysctl.conf
echo "net.core.wmem_max = 134217728" >> /etc/sysctl.conf
echo "net.core.netdev_max_backlog = 5000" >> /etc/sysctl.conf
echo "net.ipv4.tcp_congestion_control = bbr" >> /etc/sysctl.conf

# UDP performance (for QUIC)
echo "net.ipv4.udp_mem = 65536 131072 262144" >> /etc/sysctl.conf
echo "net.ipv4.udp_rmem_min = 16384" >> /etc/sysctl.conf
echo "net.ipv4.udp_wmem_min = 16384" >> /etc/sysctl.conf

# Apply changes
sysctl -p

echo "Optimization completed!"
EOF

chmod +x optimize-docker.sh
./optimize-docker.sh
```

### Resource Limits
```yaml
# docker-compose.yml with resource limits
version: '3.8'

services:
  vantun:
    image: tungoldshou/vantun:latest
    deploy:
      resources:
        limits:
          cpus: '2.0'
          memory: 512M
        reservations:
          cpus: '0.5'
          memory: 128M
```

## 🔄 Update Management

### Automated Updates
```bash
# Create update script
cat > update-vantun.sh << 'EOF'
#!/bin/bash
# VANTUN Docker update script

CONTAINER_NAME="vantun-server"
IMAGE_NAME="tungoldshou/vantun:latest"

# Pull latest image
docker pull $IMAGE_NAME

# Stop and remove old container
docker stop $CONTAINER_NAME
docker rm $CONTAINER_NAME

# Run new container
docker run -d \
  --name $CONTAINER_NAME \
  --restart unless-stopped \
  -p 4242:4242 \
  -v /etc/vantun:/etc/vantun:ro \
  $IMAGE_NAME \
  server --config /etc/vantun/server.json

echo "Update completed!"
EOF

chmod +x update-vantun.sh

# Add to crontab for automatic updates
# 0 2 * * * /path/to/update-vantun.sh >> /var/log/vantun-update.log 2>&1
```

### Rolling Updates with Docker Swarm
```bash
# Initialize Docker Swarm
docker swarm init

# Deploy stack
docker stack deploy -c docker-compose.yml vantun

# Rolling update
docker service update --image tungoldshou/vantun:latest vantun_vantun
```

## 🐛 Troubleshooting

### Common Issues

#### 1. Container Won't Start
```bash
# Check logs
docker logs vantun-server

# Check configuration
docker run --rm -v /etc/vantun:/etc/vantun:ro tungoldshou/vantun:latest vantun --config /etc/vantun/server.json --check
```

#### 2. Network Connectivity Issues
```bash
# Test network connectivity
docker run --rm appropriate/curl curl -v your-server:4242

# Check port binding
docker port vantun-server

# Test from inside container
docker exec -it vantun-server nc -v localhost 4242
```

#### 3. Performance Issues
```bash
# Monitor resource usage
docker stats vantun-server

# Check kernel parameters
docker exec vantun-server sysctl -a | grep -E "(rmem|wmem|congestion)"

# Test UDP performance
docker exec vantun-server iperf3 -s -p 4242
```

### Debug Mode
```bash
# Run with debug logging
docker run -d \
  --name vantun-debug \
  -p 4242:4242 \
  -v /etc/vantun:/etc/vantun:ro \
  tungoldshou/vantun:latest \
  server --config /etc/vantun/server.json --log-level debug

# Follow debug logs
docker logs -f vantun-debug
```

## 📚 Reference Configurations

### High-Performance Server
```json
{
  "server": true,
  "address": "0.0.0.0:4242",
  "log_level": "info",
  "multipath": true,
  "obfs": true,
  "fec_data": 10,
  "fec_parity": 3,
  "performance": {
    "workers": 8,
    "buffer_size": 4194304,
    "read_buffer": 2097152,
    "write_buffer": 2097152
  }
}
```

### Mobile-Optimized Client
```json
{
  "server": false,
  "address": "server.example.com:4242",
  "log_level": "info",
  "multipath": true,
  "obfs": true,
  "fec_data": 8,
  "fec_parity": 4,
  "mobile": {
    "network_detection": true,
    "adaptive_fec": true,
    "power_optimization": true
  }
}
```

---

## 🚀 Quick Commands Reference

```bash
# One-liner server deployment
docker run -d --name vantun-server --restart unless-stopped -p 4242:4242 -v /etc/vantun:/etc/vantun tungoldshou/vantun:latest server --config /etc/vantun/server.json

# One-liner client deployment
docker run -d --name vantun-client --restart unless-stopped -p 1080:1080 -v /etc/vantun:/etc/vantun tungoldshou/vantun:latest client --config /etc/vantun/client.json

# Check status
docker ps | grep vantun

# View logs
docker logs -f vantun-server

# Update container
docker pull tungoldshou/vantun:latest && docker restart vantun-server
```

For more information, visit our [GitHub repository](https://github.com/tungoldshou/vantun) or join our [Telegram group](https://t.me/vantun01).