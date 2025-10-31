# Vultr VPS Deployment Guide

## Nutrition Platform Deployment to love.doctorhealthy1.com

### Server Information
- **Domain**: love.doctorhealthy1.com
- **IPv6 Address**: 2001:19f0:1000:6d30:5400:05ff:fe9f:f721
- **Vultr Account**: https://my.vultr.com/

### Prerequisites
1. Vultr VPS with Ubuntu 20.04+ or Debian 11+
2. Root access to the VPS
3. Domain DNS configured to point to VPS IP

### Deployment Steps

#### 1. Prepare Local Files
The deployment script has created the following files:
- `nutrition-platform.tar` - Docker image
- `deploy-on-vps.sh` - VPS setup script
- `docker-compose.vultr.yml` - Docker Compose configuration
- `nginx/vultr.conf` - Nginx configuration

#### 2. Upload Files to VPS
```bash
# Upload deployment files
scp nutrition-platform.tar deploy-on-vps.sh root@YOUR_VPS_IP:/root/

# SSH into VPS
ssh root@YOUR_VPS_IP
```

#### 3. Run Deployment Script
```bash
# Make script executable
chmod +x deploy-on-vps.sh

# Run deployment
./deploy-on-vps.sh
```

#### 4. Configure DNS
In your domain registrar (Namecheap), add these DNS records:
```
Type: A
Name: love
Value: YOUR_VPS_IPv4_ADDRESS
TTL: 300

Type: AAAA
Name: love
Value: 2001:19f0:1000:6d30:5400:05ff:fe9f:f721
TTL: 300
```

#### 5. Verify Deployment
```bash
# Check application status
docker-compose ps

# Check logs
docker-compose logs -f

# Test health endpoint
curl http://localhost:8080/health

# Test external access
curl https://love.doctorhealthy1.com/health
```

### What the Deployment Script Does

1. **System Updates**: Updates Ubuntu/Debian packages
2. **Docker Installation**: Installs Docker and Docker Compose
3. **Application Setup**: 
   - Creates `/opt/nutrition-platform` directory
   - Loads Docker image
   - Creates docker-compose.yml
   - Sets up systemd service
4. **Nginx Configuration**:
   - Installs Nginx
   - Configures reverse proxy
   - Sets up SSL with Let's Encrypt
5. **SSL Certificate**: Automatically obtains SSL certificate for HTTPS

### Manual Configuration (if needed)

#### Environment Variables
Edit `/opt/nutrition-platform/docker-compose.yml` to add:
```yaml
environment:
  - DB_HOST=your_db_host
  - DB_PASSWORD=your_db_password
  - JWT_SECRET=your_jwt_secret
  - API_KEY_SECRET=your_api_secret
```

#### Database Setup
If using external database:
```bash
# Install PostgreSQL client
apt install -y postgresql-client

# Connect to database
psql -h your_db_host -U nutrition_user -d nutrition_platform
```

### Monitoring and Maintenance

#### Check Application Status
```bash
# Service status
systemctl status nutrition-platform

# Container status
docker ps

# Application logs
docker logs nutrition-platform

# Nginx logs
tail -f /var/log/nginx/love.doctorhealthy1.com.access.log
tail -f /var/log/nginx/love.doctorhealthy1.com.error.log
```

#### Update Application
```bash
# Stop services
systemctl stop nutrition-platform

# Update Docker image
docker load < new-nutrition-platform.tar

# Start services
systemctl start nutrition-platform
```

#### SSL Certificate Renewal
Certificates auto-renew via cron job. Manual renewal:
```bash
certbot renew --nginx
```

### Troubleshooting

#### Application Won't Start
```bash
# Check Docker logs
docker logs nutrition-platform

# Check system resources
free -h
df -h

# Restart services
systemctl restart nutrition-platform
```

#### SSL Issues
```bash
# Check certificate status
certbot certificates

# Test SSL configuration
nginx -t

# Reload Nginx
systemctl reload nginx
```

#### Database Connection Issues
```bash
# Test database connectivity
telnet your_db_host 5432

# Check environment variables
docker exec nutrition-platform env | grep DB
```

### Security Considerations

1. **Firewall**: Configure UFW to allow only necessary ports
```bash
ufw allow 22/tcp
ufw allow 80/tcp
ufw allow 443/tcp
ufw enable
```

2. **SSH Security**: Disable password authentication
3. **Regular Updates**: Keep system and Docker images updated
4. **Backup**: Regular database and file backups

### Performance Optimization

1. **Resource Monitoring**: Use `htop`, `iotop` for monitoring
2. **Log Rotation**: Configure logrotate for application logs
3. **Caching**: Redis can be added for improved performance
4. **CDN**: Consider using Cloudflare for static assets

### Support

For issues:
1. Check application logs: `docker logs nutrition-platform`
2. Check Nginx logs: `/var/log/nginx/`
3. Verify DNS propagation: `dig love.doctorhealthy1.com`
4. Test SSL: `openssl s_client -connect love.doctorhealthy1.com:443`

### Expected URLs
- **Main Application**: https://love.doctorhealthy1.com
- **Health Check**: https://love.doctorhealthy1.com/health
- **API Endpoints**: https://love.doctorhealthy1.com/api/*