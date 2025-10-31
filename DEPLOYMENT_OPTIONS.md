# 🚀 Production Deployment Options

## Option 1: Direct Server Deployment (Recommended)

**Perfect if you have:** VPS, dedicated server, or cloud instance

### What You Need to Provide:
- ✅ Server IP address
- ✅ SSH username/password (or SSH key)
- ✅ Domain name pointing to your server
- ✅ Email address for SSL certificate

### What I'll Do Automatically:
- 🔧 Install all dependencies (Docker, PostgreSQL, Redis, Nginx)
- 🔒 Configure SSL certificates with Let's Encrypt
- 🚀 Deploy application with full monitoring stack
- 📊 Set up Prometheus, Grafana, and log aggregation
- 🔄 Configure automated backups and health checks
- 🛡️ Set up firewall and security hardening

### Command:
```bash
chmod +x deploy-to-server.sh
./deploy-to-server.sh --ip YOUR_SERVER_IP --user root --pass YOUR_PASSWORD --domain api.yourdomain.com --email your@email.com
```

### Server Requirements:
- **OS:** Ubuntu 20.04+ (preferred) or CentOS 8+
- **RAM:** Minimum 2GB, Recommended 4GB+
- **Storage:** Minimum 20GB, Recommended 50GB+
- **CPU:** 2+ cores recommended
- **Network:** Public IP with ports 80, 443, 22 accessible

---

## Option 2: Coolify Deployment

**Perfect if you have:** Coolify instance for container management

### What You Need to Provide:
- ✅ Coolify instance URL
- ✅ Coolify API access token
- ✅ Domain name
- ✅ Git repository access (if private)

### What I'll Do Automatically:
- 🏗️ Create project in Coolify
- 🗄️ Set up PostgreSQL and Redis databases
- 🚀 Deploy application with environment variables
- 🔒 Configure SSL certificate
- 📊 Set up monitoring and health checks
- 🔄 Configure automated deployments

### Command:
```bash
chmod +x deploy-with-coolify.sh
./deploy-with-coolify.sh --url https://coolify.yourdomain.com --token YOUR_COOLIFY_TOKEN --domain api.yourdomain.com
```

### Prerequisites:
- Coolify instance running and accessible
- API token with deployment permissions
- Domain DNS pointing to Coolify server

---

## Option 3: Manual Deployment Guide

**Perfect if you want:** Full control over the deployment process

### Step-by-Step Guide:

#### 1. Server Preparation
```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

#### 2. Application Setup
```bash
# Clone repository
git clone https://github.com/yourusername/nutrition-platform.git
cd nutrition-platform

# Configure environment
cp backend/.env.example backend/.env
nano backend/.env  # Edit with your values
```

#### 3. Deploy with Docker
```bash
# Build and start services
docker-compose -f docker-compose.production.yml up -d

# Run migrations
docker-compose -f docker-compose.production.yml exec backend go run cmd/migrate/main.go -direction up

# Seed database
docker-compose -f docker-compose.production.yml exec backend go run cmd/seed/main.go
```

#### 4. Configure Nginx and SSL
```bash
# Install Nginx
sudo apt install nginx certbot python3-certbot-nginx

# Configure domain
sudo nano /etc/nginx/sites-available/nutrition-platform
# Copy configuration from nginx/nginx.conf

# Enable site
sudo ln -s /etc/nginx/sites-available/nutrition-platform /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx

# Get SSL certificate
sudo certbot --nginx -d yourdomain.com
```

---

## 🔧 Post-Deployment Configuration

### 1. Create First API Key
```bash
# Access the application
curl https://yourdomain.com/health

# Create admin API key (you'll need to implement admin interface or use database directly)
```

### 2. Test API Endpoints
```bash
# Test health endpoint
curl https://yourdomain.com/health

# Test API info
curl https://yourdomain.com/api/info

# Test nutrition analysis (public endpoint)
curl -X POST https://yourdomain.com/api/nutrition/analyze \
  -H "Content-Type: application/json" \
  -d '{"food": "apple", "quantity": 100, "unit": "g", "checkHalal": true}'
```

### 3. Monitor Services
- **Application Health:** https://yourdomain.com/health
- **Metrics:** https://yourdomain.com/metrics (restricted access)
- **Grafana Dashboard:** https://yourdomain.com:3001 (if using Docker Compose)
- **Prometheus:** https://yourdomain.com:9091 (if using Docker Compose)

---

## 🆘 Troubleshooting

### Common Issues:

#### 1. **Connection Refused**
```bash
# Check if services are running
docker-compose -f docker-compose.production.yml ps

# Check logs
docker-compose -f docker-compose.production.yml logs backend
```

#### 2. **Database Connection Failed**
```bash
# Check database status
docker-compose -f docker-compose.production.yml logs postgres

# Test database connection
docker-compose -f docker-compose.production.yml exec postgres psql -U nutrition_user -d nutrition_platform
```

#### 3. **SSL Certificate Issues**
```bash
# Check certificate status
sudo certbot certificates

# Renew certificate
sudo certbot renew --dry-run
```

#### 4. **API Key Authentication Issues**
```bash
# Check API key format (should be: nk_[64_hex_characters])
# Verify API key in database
docker-compose -f docker-compose.production.yml exec postgres psql -U nutrition_user -d nutrition_platform -c "SELECT * FROM api_keys;"
```

---

## 📞 Support

### If You Need Help:
1. **Provide me with:**
   - Server IP and SSH credentials
   - Domain name
   - Email for SSL certificate
   - Any specific requirements

2. **I'll handle:**
   - Complete server setup and configuration
   - Application deployment and testing
   - SSL certificate configuration
   - Monitoring setup
   - Initial API key creation
   - Documentation and handover

### Contact Information:
- Provide server details and I'll deploy everything for you
- Or provide Coolify access and I'll use that platform
- Or follow the manual guide if you prefer hands-on approach

---

## 🎯 Recommended Approach

**For most users, I recommend Option 1 (Direct Server Deployment)** because:
- ✅ Complete automation
- ✅ Full monitoring stack included
- ✅ Automated backups and health checks
- ✅ Production-grade security configuration
- ✅ Easy maintenance and updates

**Just provide:** Server IP, SSH credentials, domain name, and email - I'll handle everything else!