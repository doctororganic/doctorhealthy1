# 🚀 SETUP SERVER AND DEPLOY NUTRITION PLATFORM

## 📋 Overview

This guide will help you:
1. Set up a server for deployment
2. Configure the server for Coolify
3. Deploy the nutrition platform

## 🔗 Coolify Dashboard Link
**Access URL:** https://api.doctorhealthy1.com

## 🖥️ Step 1: Set Up a Server

### Option 1: Create New VPS (Recommended)

#### Vultr (Recommended)
1. Go to [Vultr](https://www.vultr.com/)
2. Create an account
3. Deploy a new server:
   - **Type:** Ubuntu 22.04 LTS
   - **Plan:** Regular Performance ($6/month minimum)
   - **Region:** Choose closest to your users
   - **Enable IPv6:** Yes
   - **Server Hostname:** nutrition-platform

#### DigitalOcean
1. Go to [DigitalOcean](https://www.digitalocean.com/)
2. Create an account
3. Create a Droplet:
   - **Image:** Ubuntu 22.04 LTS
   - **Plan:** Basic ($5/month minimum)
   - **Region:** Choose closest to your users
   - **Hostname:** nutrition-platform

#### AWS EC2
1. Go to [AWS EC2](https://aws.amazon.com/ec2/)
2. Launch Instance:
   - **AMI:** Ubuntu Server 22.04 LTS
   - **Instance Type:** t2.micro (Free Tier)
   - **Region:** Choose closest to your users

### Option 2: Use Existing Server

If you already have a server:
1. Ensure it's running Ubuntu 20.04+ or CentOS 8+
2. Have root access or sudo privileges
3. Docker installed (or install with commands below)

## 🔧 Step 2: Configure Server

### SSH into Your Server
```bash
ssh root@YOUR_SERVER_IP
```

### Update System
```bash
apt update && apt upgrade -y
```

### Install Docker
```bash
# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh | sh
sh get-docker.sh

# Start Docker
systemctl start docker
systemctl enable docker

# Install Docker Compose
curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose

# Test Docker
docker run hello-world
```

### Install Required Tools
```bash
# Install Node.js (if needed)
curl -fsSL https://deb.nodesource.com/setup_18.x | bash -
apt-get install -y nodejs

# Install Git
apt install -y git

# Install Nginx (for reverse proxy)
apt install -y nginx
systemctl start nginx
systemctl enable nginx
```

### Configure Firewall
```bash
# Configure UFW firewall
ufw allow ssh
ufw allow 80
ufw allow 443
ufw allow 22
ufw --force enable
```

## 🔗 Step 3: Connect Server to Coolify

### Option 1: Connect via Coolify Dashboard
1. Go to **https://api.doctorhealthy1.com**
2. Click **"Servers"** in left sidebar
3. Click **"Add Server"**
4. Choose **"Connect existing server"**
5. Follow the instructions to add your server

### Option 2: Connect via SSH Key
1. Generate SSH key if you don't have one:
```bash
ssh-keygen -t ed25519 -C "coolify"
```

2. Copy the public key:
```bash
cat ~/.ssh/id_ed25519.pub
```

3. Add the key to Coolify:
   - In Coolify dashboard, go to **"Servers"**
   - Click **"Add Server"**
   - Choose **"Connect via SSH key"**
   - Paste the public key

## 📦 Step 4: Deploy Application

### Option 1: Deploy via Coolify (Recommended)
1. In Coolify dashboard, click **"Applications"**
2. Click **"Add Application"**
3. Upload the ZIP file: `nutrition-platform-coolify-20251013-164858.zip`
4. Configure as shown in MANUAL-COOLIFY-DEPLOYMENT-GUIDE.md

### Option 2: Deploy Directly on Server

#### Create Deployment Directory
```bash
mkdir -p /opt/nutrition-platform
cd /opt/nutrition-platform
```

#### Upload and Extract Application
```bash
# Upload ZIP file to server (using scp)
scp nutrition-platform-coolify-20251013-164858.zip root@YOUR_SERVER_IP:/opt/nutrition-platform/

# Extract ZIP file
unzip nutrition-platform-coolify-20251013-164858.zip
```

#### Create Environment File
```bash
cat > .env << 'EOF'
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=nutrition_platform
DB_USER=nutrition_user
DB_PASSWORD=ac287cc0e30f54afad53c6dc7e02fd0cccad979d62b75d75d97b1ede12daf8d5
DB_SSL_MODE=require

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=f606b2d16d6697e666ce78a8685574d042df15484ca8f18f39f2e67bf38dc09a

# Security Configuration
JWT_SECRET=9a00511e8e23764f8f4524c02f1db9eccc1923208c02fb36cb758d874d8d569bce9ea1b24ac18a958334abe15ef89e09d6010fe64a1d1ffc02a45b07898b2473
API_KEY_SECRET=5d2763e839f7e71b90ff88bef12f690a41802635aa131f6bc7160056ef0aeb7dc9caaeb07dbe0028128e617529a48903f8d01c6cc64ce61419eb7f309fdfc8bc
ENCRYPTION_KEY=cc1574e486b2f5abd69d86537079ba928974cc463e36ff410647b15b15533d23
SESSION_SECRET=f40776484ee20b35e4f754909fb3067cef2a186d0da7c4c24f1bcd54870d9fba

# Server Configuration
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

# CORS Configuration
CORS_ALLOWED_ORIGINS=https://super.doctorhealthy1.com,https://my.doctorhealthy1.com

# Features
RELIGIOUS_FILTER_ENABLED=true
FILTER_ALCOOL=true
FILTER_PORK=true
DEFAULT_LANGUAGE=en
SUPPORTED_LANGUAGES=en,ar
EOF
```

#### Create Docker Compose File
```bash
cat > docker-compose.yml << 'EOF'
version: "3.8"

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - NODE_ENV=production
    volumes:
      - ./.env:/app/.env
    depends_on:
      - postgres
      - redis
    restart: unless-stopped

  postgres:
    image: postgres:15
    environment:
      POSTGRES_DB: nutrition_platform
      POSTGRES_USER: nutrition_user
      POSTGRES_PASSWORD: ac287cc0e30f54afad53c6dc7e02fd0cccad979d62b75d75d97b1ede12daf8d5
    volumes:
      - postgres_data:/var/lib/postgresql/data
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    command: redis-server --requirepass f606b2d16d6697e666ce78a8685574d042df15484ca8f18f39f2e67bf38dc09a
    volumes:
      - redis_data:/data
    restart: unless-stopped

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/nginx/ssl
    depends_on:
      - app
    restart: unless-stopped

volumes:
  postgres_data:
  redis_data:
EOF
```

#### Create Nginx Configuration
```bash
cat > nginx.conf << 'EOF'
events {
    worker_connections 1024;
}

http {
    upstream app {
        server app:8080;
    }

    server {
        listen 80;
        server_name super.doctorhealthy1.com;
        return 301 https://$server_name$request_uri;
    }

    server {
        listen 443 ssl http2;
        server_name super.doctorhealthy1.com;

        ssl_certificate /etc/nginx/ssl/cert.pem;
        ssl_certificate_key /etc/nginx/ssl/key.pem;

        location / {
            proxy_pass http://app;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }
    }
}
EOF
```

#### Deploy Application
```bash
# Build and start services
docker-compose up -d --build

# Check status
docker-compose ps

# View logs
docker-compose logs -f app
```

## 🔐 Step 5: Configure SSL Certificate

### Option 1: Let's Encrypt (Recommended)
```bash
# Install Certbot
apt install certbot python3-certbot-nginx

# Get SSL certificate
certbot --nginx -d super.doctorhealthy1.com

# Set up auto-renewal
echo "0 12 * * * /usr/bin/certbot renew --quiet" | crontab -
```

### Option 2: Self-Signed Certificate (Testing Only)
```bash
# Create SSL directory
mkdir -p /etc/nginx/ssl

# Generate self-signed certificate
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout /etc/nginx/ssl/key.pem \
  -out /etc/nginx/ssl/cert.pem \
  -subj "/C=US/ST=State/L=City/O=Organization/CN=super.doctorhealthy1.com"

# Restart Nginx
systemctl restart nginx
```

## 🔍 Step 6: Verify Deployment

### Test Application
```bash
# Test locally on server
curl http://localhost:8080/health

# Test from external
curl https://super.doctorhealthy1.com/health
```

### Test API
```bash
curl -X POST https://super.doctorhealthy1.com/api/v1/nutrition/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "food": "chicken breast",
    "quantity": 100,
    "unit": "grams",
    "checkHalal": true,
    "language": "en"
  }'
```

## 📊 Monitoring

### Server Monitoring
```bash
# Check system resources
htop

# Check Docker containers
docker stats

# Check disk usage
df -h
```

### Application Monitoring
```bash
# View application logs
docker-compose logs -f app

# Check database logs
docker-compose logs postgres

# Check Redis logs
docker-compose logs redis
```

## 🚨 Troubleshooting

### Application Won't Start
1. Check logs: `docker-compose logs app`
2. Verify environment variables
3. Check if ports are available

### Database Connection Issues
1. Check PostgreSQL logs: `docker-compose logs postgres`
2. Verify database credentials
3. Check if database is running

### SSL Certificate Issues
1. Check Nginx error logs: `tail -f /var/log/nginx/error.log`
2. Verify certificate files exist
3. Check certificate expiration

## 📋 Post-Deployment Checklist

- [ ] Server is running and accessible
- [ ] Docker is installed and running
- [ ] Application is deployed and accessible
- [ ] SSL certificate is valid
- [ ] Database connections are working
- [ ] All endpoints are responding correctly
- [ ] Monitoring is configured
- [ ] Backups are configured

## 🎯 Success Criteria

Your deployment is successful when:
- ✅ Application loads at https://super.doctorhealthy1.com
- ✅ Health check returns 200 OK
- ✅ API endpoints respond correctly
- ✅ SSL certificate is valid and trusted
- ✅ Database is connected and accessible
- ✅ All features are working as expected

## 📞 Support

For any issues:
1. Check application logs
2. Verify server configuration
3. Check network connectivity
4. Review documentation

## 🎊 Congratulations!

Once deployed, your AI-powered nutrition platform will be live with:
- ✅ Real-time nutrition analysis
- ✅ 10 evidence-based diet plans
- ✅ Recipe management system
- ✅ Health tracking and analytics
- ✅ Medication management
- ✅ Workout programs
- ✅ Multi-language support (EN/AR)
- ✅ Religious dietary filtering
- ✅ SSL secured with HTTPS

---
**Last Updated:** October 13, 2025  
**Deployment Status:** ✅ READY FOR DEPLOYMENT  
**Security Level:** 🔒 PRODUCTION READY