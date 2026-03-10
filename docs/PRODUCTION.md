# LiteLog Production Deployment Guide

This guide covers deploying LiteLog in production environments.

## Prerequisites

- Linux-based server (Ubuntu 20.04+, Debian 11+, CentOS 8+, or similar)
- Root or sudo access
- Basic understanding of systemd

---

## Systemd Service

### 1. Create dedicated user
\`\`\`bash
sudo useradd -r -s /bin/false litelog
\`\`\`

### 2. Create systemd service file

Create \`/etc/systemd/system/litelog.service\`:

\`\`\`ini
[Unit]
Description=LiteLog Log Ingestion Server
After=network.target

[Service]
Type=simple
User=litelog
Group=litelog
WorkingDirectory=/opt/litelog

ExecStart=/opt/litelog/litelog start \\
    --port 8080 \\
    --db /var/lib/litelog/litelog.db \\
    --retention 30d

Restart=always
RestartSec=5

LimitNOFILE=65536
LimitNPROC=4096

NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/var/lib/litelog /var/log/litelog

StandardOutput=journal
StandardError=journal
SyslogIdentifier=litelog

[Install]
WantedBy=multi-user.target
\`\`\`

### 3. Create directories and install binary

\`\`\`bash
sudo mkdir -p /opt/litelog /var/lib/litelog /var/log/litelog
sudo chown -R litelog:litelog /opt/litelog /var/lib/litelog /var/log/litelog

sudo cp /path/to/litelog /opt/litelog/litelog
sudo chmod +x /opt/litelog/litelog
sudo chown litelog:litelog /opt/litelog/litelog
\`\`\`

### 4. Start service

\`\`\`bash
sudo systemctl daemon-reload
sudo systemctl enable litelog
sudo systemctl start litelog
sudo systemctl status litelog
\`\`\`

---

## Docker/Docker Compose

### Dockerfile

\`\`\`dockerfile
FROM golang:1.24.2-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o litelog ./cmd/litelog

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/litelog /app/litelog

RUN addgroup -g 1000 litelog && \\
    adduser -D -u 1000 -G litelog litelog

USER litelog
EXPOSE 8080

CMD ["./litelog", "start", "--port", "8080", "--retention", "30d"]
\`\`\`

### docker-compose.yml

\`\`\`yaml
version: "3.8"

services:
  litelog:
    image: litelog:latest
    container_name: litelog
    restart: unless-stopped
    ports:
      - "8080:8080"
    volumes:
      - litelog-data:/var/lib/litelog
      - litelog-logs:/var/log/litelog
    environment:
      - LITELOG_PORT=8080
      - LITELOG_RETENTION=30d
    healthcheck:
      test: ["CMD", "wget", "--spider", "-q", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 10s
    deploy:
      resources:
        limits:
          cpus: "1"
          memory: 256M
        reservations:
          cpus: "0.25"
          memory: 64M

volumes:
  litelog-data:
  litelog-logs:
\`\`\`

### Start with Docker Compose

\`\`\`bash
docker-compose up -d
docker-compose logs -f litelog
\`\`\`

---

## NGINX Reverse Proxy with SSL

### 1. Install NGINX

\`\`\`bash
sudo apt update && sudo apt install -y nginx
\`\`\`

### 2. Generate SSL with Let's Encrypt

\`\`\`bash
sudo apt install -y certbot python3-certbot-nginx
sudo certbot --nginx -d litelog.yourdomain.com
\`\`\`

### 3. Create NGINX configuration

Create \`/etc/nginx/sites-available/litelog\`:

\`\`\`nginx
upstream litelog_backend {
    server 127.0.0.1:8080;
    keepalive 32;
}

server {
    listen 80;
    server_name litelog.yourdomain.com;
    return 301 https://\$server_name\$request_uri;
}

server {
    listen 443 ssl http2;
    server_name litelog.yourdomain.com;

    ssl_certificate /etc/letsencrypt/live/litelog.yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/litelog.yourdomain.com/privkey.pem;

    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;

    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    add_header X-Content-Type-Options nosniff always;
    add_header X-Frame-Options DENY always;
    add_header X-XSS-Protection "1; mode=block" always;

    access_log /var/log/nginx/litelog-access.log;
    error_log /var/log/nginx/litelog-error.log;

    location / {
        proxy_pass http://litelog_backend;
        proxy_http_version 1.1;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        proxy_set_header Connection "";
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
        proxy_buffering off;
        proxy_request_buffering off;
    }

    location /health {
        proxy_pass http://litelog_backend/health;
        access_log off;
    }
}
\`\`\`

### 4. Enable configuration

\`\`\`bash
sudo ln -s /etc/nginx/sites-available/litelog /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
\`\`\`

---

## OS Limits (ulimit)

### Check current limits

\`\`\`bash
ulimit -n   # Open files
ulimit -u   # User processes
\`\`\`

### Increase limits

\`\`\`bash
ulimit -n 65536
ulimit -u 4096
\`\`\`

### Systemd service limits

Already included in service file:
\`\`\`ini
LimitNOFILE=65536
LimitNPROC=4096
\`\`\`

---

## Security Best Practices

### 1. Run as non-root user

\`\`\`bash
sudo useradd -r -s /bin/false litelog
\`\`\`

### 2. Network isolation

Bind to localhost (127.0.0.1) and use reverse proxy. Use firewall rules to restrict access.

### 3. Database security

\`\`\`bash
sudo chmod 750 /var/lib/litelog
sudo chown -R litelog:litelog /var/lib/litelog
\`\`\`

### 4. Log retention

\`\`\`bash
litelog start --retention 30d
\`\`\`

### 5. Rate limiting (NGINX)

\`\`\`nginx
limit_req_zone \$binary_remote_addr zone=litelog_limit:10m rate=100r/s;
limit_req zone=litelog_limit burst=200 nodelay;
\`\`\`

---

## Performance Tuning

### SQLite optimization

- Connection pooling: SQLite handles concurrent reads efficiently in WAL mode
- Batch ingestion: Send logs in batches
- Indexing: Database schema optimized for timestamp-based queries

### System resources

\`\`\`ini
MemoryLimit=256M
CPUQuota=100%
\`\`\`

### Monitor metrics

- Ingestion rate (logs/sec)
- Database size growth
- Memory usage
- Disk I/O
- Response time for /ingest endpoint

---

## Monitoring and Maintenance

### Health checks

\`\`\`bash
sudo systemctl status litelog
sudo netstat -tlnp | grep 8080
sudo journalctl -u litelog -n 100 --no-pager
\`\`\`

### Log rotation

Create \`/etc/logrotate.d/litelog\`:

\`\`\`bash
/var/log/litelog/*.log {
    daily
    rotate 7
    compress
    delaycompress
    missingok
    notifempty
    create 0640 litelog litelog
}
\`\`\`

### Database maintenance

\`\`\`bash
sqlite3 /var/lib/litelog/litelog.db "VACUUM;"
\`\`\`

### Backup script

\`\`\`bash
#!/bin/bash
BACKUP_DIR="/backup/litelog"
DATE=\$(date +%Y%m%d_%H%M%S)
mkdir -p "\$BACKUP_DIR"
cp /var/lib/litelog/litelog.db "\$BACKUP_DIR/litelog_\$DATE.db"
find "\$BACKUP_DIR" -name "litelog_*.db" -mtime +7 -delete
echo "Backup completed: litelog_\$DATE.db"
\`\`\`

Add to crontab:

\`\`\`bash
0 2 * * * /opt/litelog/backup.sh >> /var/log/litelog/backup.log 2>&1
\`\`\`

---

## Troubleshooting

### Service won't start

\`\`\`bash
sudo journalctl -u litelog -n 50
ls -la /opt/litelog/litelog
ls -la /var/lib/litelog
\`\`\`

### Can't connect from remote host

- Check firewall rules
- Verify NGINX configuration
- Ensure LiteLog is listening on correct port

### Database locked

- Check for other processes accessing database
- Verify retention settings

### High memory usage

- Reduce retention period
- Check for unbounded query results
- Monitor ingestion rate

---

## Summary

This guide covers:

- ✅ Systemd service for automatic startup
- ✅ Docker/Docker Compose for containerized deployment
- ✅ NGINX reverse proxy with SSL/TLS
- ✅ OS limits for high-volume ingestion
- ✅ Security best practices
- ✅ Performance tuning
- ✅ Monitoring and maintenance

For additional help or questions, please:

- Open an issue on [GitHub](https://github.com/yashnaiduu/Litelog/issues)
- Review the [usage guide](GUIDE.md)
- Check the [contributing guide](../CONTRIBUTING.md)

---

<div align="center">
  <sub>LiteLog — centralized logging without infrastructure.</sub>
</div>
