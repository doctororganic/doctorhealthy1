#!/usr/bin/env node

/**
 * AI Deployment Assistant
 * One-click deployment tool for nutrition platform
 * 
 * This tool helps you deploy your application by:
 * 1. Setting up the server
 * 2. Configuring the environment
 * 3. Deploying the application
 * 4. Verifying the deployment
 */

const fs = require('fs');
const path = require('path');
const readline = require('readline');
const { execSync } = require('child_process');

// Colors for output
const colors = {
  reset: '\x1b[0m',
  bright: '\x1b[1m',
  red: '\x1b[31m',
  green: '\x1b[32m',
  yellow: '\x1b[33m',
  blue: '\x1b[34m',
  magenta: '\x1b[35m',
  cyan: '\x1b[36m'
};

// Helper function to print colored output
function colorLog(color, message) {
  console.log(`${colors[color]}${message}${colors.reset}`);
}

// Helper function for user input
function askQuestion(question) {
  const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout
  });
  
  return new Promise((resolve) => {
    rl.question(`${colors.cyan}${question}${colors.reset} `, (answer) => {
      rl.close();
      resolve(answer);
    });
  });
}

// AI Deployment Assistant Class
class AIDeploymentAssistant {
  constructor() {
    this.config = {
      domain: 'super.doctorhealthy1.com',
      projectName: 'nutrition-platform',
      appPort: 8080
    };
    
    this.credentials = {
      dbPassword: 'ac287cc0e30f54afad53c6dc7e02fd0cccad979d62b75d75d97b1ede12daf8d5',
      redisPassword: 'f606b2d16d6697e666ce78a8685574d042df15484ca8f18f39f2e67bf38dc09a',
      jwtSecret: '9a00511e8e23764f8f4524c02f1db9eccc1923208c02fb36cb758d874d8d569bce9ea1b24ac18a958334abe15ef89e09d6010fe64a1d1ffc02a45b07898b2473',
      apiKeySecret: '5d2763e839f7e71b90ff88bef12f690a41802635aa131f6bc7160056ef0aeb7dc9caaeb07dbe0028128e617529a48903f8d01c6cc64ce61419eb7f309fdfc8bc',
      encryptionKey: 'cc1574e486b2f5abd69d86537079ba928974cc463e36ff410647b15b15533d23',
      sessionSecret: 'f40776484ee20b35e4f754909fb3067cef2a186d0da7c4c24f1bcd54870d9fba'
    };
  }
  
  async start() {
    colorLog('bright', '🤖 AI Deployment Assistant for Nutrition Platform');
    colorLog('blue', '===============================================');
    console.log('');
    
    colorLog('green', 'I will help you deploy your nutrition platform with ease!');
    console.log('');
    
    await this.showMenu();
  }
  
  async showMenu() {
    colorLog('cyan', 'What would you like to do?');
    console.log('');
    console.log(`${colors.yellow}1.${colors.reset} Setup a new server`);
    console.log(`${colors.yellow}2.${colors.reset} Connect existing server to Coolify`);
    console.log(`${colors.yellow}3.${colors.reset} Deploy application to server`);
    console.log(`${colors.yellow}4.${colors.reset} Verify deployment`);
    console.log(`${colors.yellow}5.${colors.reset} Configure SSL certificate`);
    console.log(`${colors.yellow}6.${colors.reset} Monitor application`);
    console.log(`${colors.yellow}7.${colors.reset} Troubleshoot issues`);
    console.log(`${colors.yellow}8.${colors.reset} Exit`);
    console.log('');
    
    const choice = await askQuestion('Enter your choice (1-8):');
    
    switch (choice) {
      case '1':
        await this.setupNewServer();
        break;
      case '2':
        await this.connectServerToCoolify();
        break;
      case '3':
        await this.deployApplication();
        break;
      case '4':
        await this.verifyDeployment();
        break;
      case '5':
        await this.configureSSL();
        break;
      case '6':
        await this.monitorApplication();
        break;
      case '7':
        await this.troubleshootIssues();
        break;
      case '8':
        colorLog('green', '👋 Thank you for using AI Deployment Assistant!');
        process.exit(0);
        break;
      default:
        colorLog('red', '❌ Invalid choice. Please try again.');
        await this.showMenu();
    }
  }
  
  async setupNewServer() {
    colorLog('blue', '🖥️ Setting up a new server');
    colorLog('blue', '====================');
    console.log('');
    
    colorLog('cyan', 'Choose a cloud provider:');
    console.log('');
    console.log(`${colors.yellow}1.${colors.reset} Vultr (Recommended - $6/month)`);
    console.log(`${colors.yellow}2.${colors.reset} DigitalOcean (Recommended - $5/month)`);
    console.log(`${colors.yellow}3.${colors.reset} AWS EC2 (Free tier available)`);
    console.log(`${colors.yellow}4.${colors.reset} Google Cloud Platform`);
    console.log(`${colors.yellow}5.${colors.reset} Microsoft Azure`);
    console.log('');
    
    const provider = await askQuestion('Enter your choice (1-5):');
    
    let setupInstructions = '';
    
    switch (provider) {
      case '1':
        setupInstructions = this.getVultrSetupInstructions();
        break;
      case '2':
        setupInstructions = this.getDigitalOceanSetupInstructions();
        break;
      case '3':
        setupInstructions = this.getAWSSetupInstructions();
        break;
      case '4':
        setupInstructions = this.getGCPSetupInstructions();
        break;
      case '5':
        setupInstructions = this.getAzureSetupInstructions();
        break;
      default:
        colorLog('red', '❌ Invalid choice. Please try again.');
        await this.setupNewServer();
        return;
    }
    
    console.log('');
    colorLog('green', '📋 Server Setup Instructions:');
    console.log('');
    colorLog('cyan', setupInstructions);
    console.log('');
    
    colorLog('yellow', 'After setting up your server:');
    console.log('');
    console.log('1. SSH into your server:');
    colorLog('bright', `   ssh root@YOUR_SERVER_IP`);
    console.log('');
    console.log('2. Run this automated setup script:');
    colorLog('bright', `   curl -fsSL https://raw.githubusercontent.com/kilo-code/ai-deployment-assistant/main/setup-server.sh | bash`);
    console.log('');
    console.log('3. Once setup is complete, return here and choose option 2 to connect to Coolify');
    console.log('');
    
    await askQuestion('Press Enter to continue...');
    await this.showMenu();
  }
  
  getVultrSetupInstructions() {
    return `1. Go to https://www.vultr.com/
2. Create an account or login
3. Click "Deploy Server"
4. Choose these settings:
   - Server Type: Ubuntu 22.04 LTS x64
   - Server Plan: Regular Performance ($6/month minimum)
   - Server Location: Choose closest to your users
   - Additional Features: Enable IPv6
   - Server Hostname: nutrition-platform
5. Click "Deploy Now"
6. Wait for server to be ready (usually 1-2 minutes)
7. Note your server IP address`;
  }
  
  getDigitalOceanSetupInstructions() {
    return `1. Go to https://www.digitalocean.com/
2. Create an account or login
3. Click "Create" -> "Droplets"
4. Choose these settings:
   - Distribution: Ubuntu 22.04 LTS (LTS)
   - Plan: Basic ($5/month minimum)
   - Region: Choose closest to your users
   - Authentication: SSH Key (recommended) or Password
   - Hostname: nutrition-platform
5. Click "Create Droplet"
6. Wait for droplet to be ready (usually 1-2 minutes)
7. Note your server IP address`;
  }
  
  getAWSSetupInstructions() {
    return `1. Go to https://aws.amazon.com/ec2/
2. Create an account or login
3. Click "Launch Instance"
4. Choose these settings:
   - AMI: Ubuntu Server 22.04 LTS
   - Instance Type: t2.micro (Free Tier eligible)
   - Region: Choose closest to your users
   - Key pair: Create or select an SSH key
   - Network settings: Default
   - Storage: 8 GiB gp2 (default)
5. Click "Launch Instance"
6. Wait for instance to be ready (usually 2-3 minutes)
7. Note your server IP address`;
  }
  
  getGCPSetupInstructions() {
    return `1. Go to https://cloud.google.com/compute/
2. Create an account or login
3. Click "Create Instance"
4. Choose these settings:
   - Name: nutrition-platform
   - Region: Choose closest to your users
   - Machine type: E2-micro (Free tier eligible)
   - Boot disk: Ubuntu 22.04 LTS
   - Firewall: Allow HTTP, HTTPS traffic
5. Click "Create"
6. Wait for instance to be ready (usually 2-3 minutes)
7. Note your server IP address`;
  }
  
  getAzureSetupInstructions() {
    return `1. Go to https://portal.azure.com/
2. Create an account or login
3. Click "Create a resource" -> "Virtual machine"
4. Choose these settings:
   - Resource group: Create new (nutrition-platform)
   - Virtual machine name: nutrition-platform
   - Region: Choose closest to your users
   - Image: Ubuntu Server 22.04 LTS
   - Size: B1s (Free tier eligible)
   - Authentication type: SSH public key
   - Inbound port rules: Allow HTTP (80), HTTPS (443), SSH (22)
5. Click "Review + create"
6. Wait for VM to be ready (usually 2-3 minutes)
7. Note your server IP address`;
  }
  
  async connectServerToCoolify() {
    colorLog('blue', '🔗 Connecting Server to Coolify');
    colorLog('blue', '===============================');
    console.log('');
    
    const serverIP = await askQuestion('Enter your server IP address:');
    
    console.log('');
    colorLog('green', '📋 Steps to connect your server to Coolify:');
    console.log('');
    console.log('1. Go to: https://api.doctorhealthy1.com');
    console.log('2. Click "Servers" in the left sidebar');
    console.log('3. Click "Add Server"');
    console.log('4. Choose "Connect existing server"');
    console.log('5. Enter these details:');
    console.log(`   - IP Address: ${serverIP}`);
    console.log('   - SSH User: root');
    console.log('   - SSH Port: 22');
    console.log('6. Click "Connect server"');
    console.log('7. Wait for connection to be established');
    console.log('');
    
    colorLog('cyan', '🔑 SSH Key Method (Recommended):');
    console.log('');
    console.log('1. Generate SSH key on your server:');
    console.log('   ssh-keygen -t ed25519 -C "coolify"');
    console.log('');
    console.log('2. Copy the public key:');
    console.log('   cat ~/.ssh/id_ed25519.pub');
    console.log('');
    console.log('3. Add the key to Coolify in the server connection form');
    console.log('');
    
    colorLog('yellow', '⚠️ Important: Make sure your server is accessible from the internet');
    console.log('');
    
    await askQuestion('Press Enter to continue...');
    await this.showMenu();
  }
  
  async deployApplication() {
    colorLog('blue', '🚀 Deploying Application');
    colorLog('blue', '==================');
    console.log('');
    
    colorLog('cyan', 'Choose deployment method:');
    console.log('');
    console.log(`${colors.yellow}1.${colors.reset} Deploy via Coolify (Recommended)`);
    console.log(`${colors.yellow}2.${colors.reset} Deploy directly to server`);
    console.log('');
    
    const method = await askQuestion('Enter your choice (1-2):');
    
    if (method === '1') {
      await this.deployViaCoolify();
    } else {
      await this.deployDirectly();
    }
  }
  
  async deployViaCoolify() {
    colorLog('green', '📋 Deploying via Coolify');
    console.log('');
    
    console.log('1. Go to: https://api.doctorhealthy1.com');
    console.log('2. Click "Applications" in the left sidebar');
    console.log('3. Click "Add Application"');
    console.log('4. Configure these settings:');
    console.log('   - Application Name: nutrition-platform-secure');
    console.log('   - Description: AI-powered nutrition platform');
    console.log('   - Project: new doctorhealthy1');
    console.log('   - Environment: production');
    console.log('');
    console.log('5. Upload the ZIP file:');
    console.log('   - Select "Upload ZIP file"');
    console.log('   - Choose: nutrition-platform-coolify-20251013-164858.zip');
    console.log('');
    console.log('6. Configure build settings:');
    console.log('   - Build Pack: Dockerfile');
    console.log('   - Dockerfile Location: backend/Dockerfile');
    console.log('   - Port: 8080');
    console.log('');
    console.log('7. Configure deployment settings:');
    console.log(`   - Domain: ${this.config.domain}`);
    console.log('   - Health Check Path: /health');
    console.log('   - Auto Deploy: Enabled');
    console.log('');
    
    colorLog('magenta', '🔐 Environment Variables:');
    console.log('');
    console.log('Click "Environment Variables" tab and add these variables:');
    console.log('');
    console.log(`${colors.green}# Database Configuration`);
    console.log(`DB_HOST=localhost`);
    console.log(`DB_PORT=5432`);
    console.log(`DB_NAME=nutrition_platform`);
    console.log(`DB_USER=nutrition_user`);
    console.log(`DB_PASSWORD=${this.credentials.dbPassword}`);
    console.log(`DB_SSL_MODE=require`);
    console.log('');
    console.log(`${colors.green}# Redis Configuration`);
    console.log(`REDIS_HOST=localhost`);
    console.log(`REDIS_PORT=6379`);
    console.log(`REDIS_PASSWORD=${this.credentials.redisPassword}`);
    console.log('');
    console.log(`${colors.green}# Security Configuration`);
    console.log(`JWT_SECRET=${this.credentials.jwtSecret}`);
    console.log(`API_KEY_SECRET=${this.credentials.apiKeySecret}`);
    console.log(`ENCRYPTION_KEY=${this.credentials.encryptionKey}`);
    console.log(`SESSION_SECRET=${this.credentials.sessionSecret}`);
    console.log('');
    console.log(`${colors.green}# Server Configuration`);
    console.log(`SERVER_HOST=0.0.0.0`);
    console.log(`SERVER_PORT=8080`);
    console.log(`CORS_ALLOWED_ORIGINS=https://${this.config.domain},https://my.doctorhealthy1.com`);
    console.log('');
    
    console.log('8. Add database services:');
    console.log('   - Click "Services" tab');
    console.log('   - Click "Add Service"');
    console.log('   - Select PostgreSQL 15');
    console.log('   - Name: nutrition-postgres');
    console.log('   - Database: nutrition_platform');
    console.log('   - Username: nutrition_user');
    console.log(`   - Password: ${this.credentials.dbPassword}`);
    console.log('');
    console.log('   - Click "Add Another Service"');
    console.log('   - Select Redis 7-alpine');
    console.log('   - Name: nutrition-redis');
    console.log(`   - Password: ${this.credentials.redisPassword}`);
    console.log('');
    
    console.log('9. Deploy application:');
    console.log('   - Click "Deploy" button');
    console.log('   - Wait 5-10 minutes for deployment');
    console.log('   - Monitor progress in "Deployments" tab');
    console.log('');
    
    colorLog('green', '✅ Your application will be deployed to:');
    colorLog('bright', `   https://${this.config.domain}`);
    
    await askQuestion('Press Enter to continue...');
    await this.showMenu();
  }
  
  async deployDirectly() {
    colorLog('green', '📋 Deploying directly to server');
    console.log('');
    
    const serverIP = await askQuestion('Enter your server IP address:');
    
    console.log('');
    colorLog('cyan', '🔧 I will now generate deployment commands for your server');
    console.log('');
    
    // Create deployment script for the server
    const deployScript = `#!/bin/bash
# Nutrition Platform Deployment Script
# Generated by AI Deployment Assistant

set -e

echo "🚀 Starting Nutrition Platform Deployment..."

# Create deployment directory
mkdir -p /opt/nutrition-platform
cd /opt/nutrition-platform

# Download and extract application
echo "📦 Downloading application..."
if [ ! -f nutrition-platform-coolify-20251013-164858.zip ]; then
  echo "❌ ZIP file not found. Please upload it first."
  exit 1
fi

unzip nutrition-platform-coolify-20251013-164858.zip

# Create environment file
echo "⚙️ Creating environment file..."
cat > .env << 'EOF'
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=nutrition_platform
DB_USER=nutrition_user
DB_PASSWORD=${this.credentials.dbPassword}
DB_SSL_MODE=require

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=${this.credentials.redisPassword}

# Security Configuration
JWT_SECRET=${this.credentials.jwtSecret}
API_KEY_SECRET=${this.credentials.apiKeySecret}
ENCRYPTION_KEY=${this.credentials.encryptionKey}
SESSION_SECRET=${this.credentials.sessionSecret}

# Server Configuration
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

# CORS Configuration
CORS_ALLOWED_ORIGINS=https://${this.config.domain},https://my.doctorhealthy1.com

# Features
RELIGIOUS_FILTER_ENABLED=true
FILTER_ALCOOL=true
FILTER_PORK=true
DEFAULT_LANGUAGE=en
SUPPORTED_LANGUAGES=en,ar
EOF

# Create Docker Compose file
echo "🐳 Creating Docker Compose configuration..."
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
      POSTGRES_PASSWORD: ${this.credentials.dbPassword}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    command: redis-server --requirepass ${this.credentials.redisPassword}
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

# Create Nginx configuration
echo "🌐 Creating Nginx configuration..."
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
        server_name ${this.config.domain};
        return 301 https://$server_name$request_uri;
    }

    server {
        listen 443 ssl http2;
        server_name ${this.config.domain};

        ssl_certificate /etc/nginx/ssl/cert.pem;
        ssl_certificate_key /etc/nginx/ssl/key.pem;
        ssl_protocols TLSv1.2 TLSv1.3;
        ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA384;
        ssl_prefer_server_ciphers off;

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

# Build and start services
echo "🚀 Building and starting services..."
docker-compose up -d --build

# Wait for services to be ready
echo "⏳ Waiting for services to be ready..."
sleep 30

# Check service status
echo "📊 Checking service status..."
docker-compose ps

echo "✅ Deployment completed!"
echo "📍 Your application is available at: https://${this.config.domain}"
echo "🏥 Health check: https://${this.config.domain}/health"
echo "📊 API: https://${this.config.domain}/api"
`;

    console.log('');
    colorLog('green', '📋 Generated deployment script:');
    console.log('');
    colorLog('cyan', deployScript);
    console.log('');
    
    colorLog('magenta', '🔧 To deploy to your server, run these commands:');
    console.log('');
    colorLog('bright', `ssh root@${serverIP}`);
    console.log('');
    colorLog('bright', '# Create and run deployment script');
    colorLog('bright', 'cat > deploy.sh << \'EOF\'');
    colorLog('bright', deployScript);
    colorLog('bright', 'EOF');
    colorLog('bright', 'chmod +x deploy.sh');
    colorLog('bright', './deploy.sh');
    console.log('');
    
    await askQuestion('Press Enter to continue...');
    await this.showMenu();
  }
  
  async verifyDeployment() {
    colorLog('blue', '🔍 Verifying Deployment');
    colorLog('blue', '==================');
    console.log('');
    
    const domain = await askQuestion(`Enter your domain (default: ${this.config.domain}):`) || this.config.domain;
    
    colorLog('cyan', '🌐 Checking URLs...');
    console.log('');
    
    // Check main site
    try {
      const https = require('https');
      const response = await new Promise((resolve, reject) => {
        const req = https.request(`https://${domain}`, { timeout: 10000 }, (res) => {
          let data = '';
          res.on('data', chunk => data += chunk);
          res.on('end', () => resolve({ status: res.statusCode, data }));
        });
        req.on('error', reject);
        req.end();
      });
      
      if (response.status === 200) {
        colorLog('green', '✅ Main site is accessible!');
        colorLog('bright', `   https://${domain}`);
      } else {
        colorLog('yellow', `⚠️ Main site returned status: ${response.status}`);
      }
    } catch (error) {
      colorLog('red', `❌ Main site not accessible: ${error.message}`);
    }
    
    // Check health endpoint
    try {
      const https = require('https');
      const response = await new Promise((resolve, reject) => {
        const req = https.request(`https://${domain}/health`, { timeout: 10000 }, (res) => {
          let data = '';
          res.on('data', chunk => data += chunk);
          res.on('end', () => resolve({ status: res.statusCode, data }));
        });
        req.on('error', reject);
        req.end();
      });
      
      if (response.status === 200) {
        colorLog('green', '✅ Health endpoint is working!');
        colorLog('bright', `   https://${domain}/health`);
        try {
          const healthData = JSON.parse(response.data);
          colorLog('cyan', `   Status: ${healthData.status}`);
          colorLog('cyan', `   Version: ${healthData.version}`);
        } catch (e) {
          colorLog('cyan', '   Response: ' + response.data.substring(0, 100) + '...');
        }
      } else {
        colorLog('yellow', `⚠️ Health endpoint returned status: ${response.status}`);
      }
    } catch (error) {
      colorLog('red', `❌ Health endpoint not accessible: ${error.message}`);
    }
    
    // Check API endpoint
    try {
      const https = require('https');
      const response = await new Promise((resolve, reject) => {
        const req = https.request(`https://${domain}/api/v1/info`, { timeout: 10000 }, (res) => {
          let data = '';
          res.on('data', chunk => data += chunk);
          res.on('end', () => resolve({ status: res.statusCode, data }));
        });
        req.on('error', reject);
        req.end();
      });
      
      if (response.status === 200) {
        colorLog('green', '✅ API endpoint is working!');
        colorLog('bright', `   https://${domain}/api/v1/info`);
      } else {
        colorLog('yellow', `⚠️ API endpoint returned status: ${response.status}`);
      }
    } catch (error) {
      colorLog('red', `❌ API endpoint not accessible: ${error.message}`);
    }
    
    console.log('');
    colorLog('magenta', '🧪 Test API endpoint:');
    console.log('');
    console.log('You can test the nutrition API with this curl command:');
    console.log('');
    colorLog('bright', `curl -X POST https://${domain}/api/v1/nutrition/analyze \\`);
    colorLog('bright', `  -H "Content-Type: application/json" \\`);
    colorLog('bright', `  -d '{`);
    colorLog('bright', `    "food": "chicken breast",`);
    colorLog('bright', `    "quantity": 100,`);
    colorLog('bright', `    "unit": "grams",`);
    colorLog('bright', `    "checkHalal": true,`);
    colorLog('bright', `    "language": "en"`);
    colorLog('bright', `  }'`);
    console.log('');
    
    await askQuestion('Press Enter to continue...');
    await this.showMenu();
  }
  
  async configureSSL() {
    colorLog('blue', '🔒 Configuring SSL Certificate');
    colorLog('blue', '===========================');
    console.log('');
    
    const domain = await askQuestion(`Enter your domain (default: ${this.config.domain}):`) || this.config.domain;
    
    console.log('');
    colorLog('cyan', 'Choose SSL configuration:');
    console.log('');
    console.log(`${colors.yellow}1.${colors.reset} Let\\'s Encrypt (Recommended for production)`);
    console.log(`${colors.yellow}2.${colors.reset} Self-signed certificate (For testing)`);
    console.log(`${colors.yellow}3.${colors.reset} Skip SSL configuration`);
    console.log('');
    
    const choice = await askQuestion('Enter your choice (1-3):');
    
    switch (choice) {
      case '1':
        await this.configureLetsEncrypt(domain);
        break;
      case '2':
        await this.configureSelfSignedSSL(domain);
        break;
      case '3':
        colorLog('yellow', '⚠️ Skipping SSL configuration');
        break;
      default:
        colorLog('red', '❌ Invalid choice. Please try again.');
        await this.configureSSL();
        return;
    }
    
    await askQuestion('Press Enter to continue...');
    await this.showMenu();
  }
  
  async configureLetsEncrypt(domain) {
    colorLog('green', '🔐 Configuring Let\\'s Encrypt');
    console.log('');
    
    console.log('📋 Steps to configure Let\\'s Encrypt:');
    console.log('');
    console.log('1. SSH into your server:');
    colorLog('bright', `   ssh root@YOUR_SERVER_IP`);
    console.log('');
    console.log('2. Install Certbot:');
    colorLog('bright', '   apt update && apt install certbot python3-certbot-nginx');
    console.log('');
    console.log('3. Get SSL certificate:');
    colorLog('bright', `   certbot --nginx -d ${domain}`);
    console.log('');
    console.log('4. Set up auto-renewal:');
    colorLog('bright', '   echo "0 12 * * * /usr/bin/certbot renew --quiet" | crontab -');
    console.log('');
    console.log('5. Restart Nginx:');
    colorLog('bright', '   systemctl restart nginx');
    console.log('');
    
    colorLog('green', '✅ Your SSL certificate will be automatically renewed!');
  }
  
  async configureSelfSignedSSL(domain) {
    colorLog('green', '🔐 Configuring Self-Signed SSL');
    console.log('');
    
    const sslScript = `#!/bin/bash
# Self-Signed SSL Configuration for ${domain}

# Create SSL directory
mkdir -p /etc/nginx/ssl

# Generate self-signed certificate
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \\
  -keyout /etc/nginx/ssl/key.pem \\
  -out /etc/nginx/ssl/cert.pem \\
  -subj "/C=US/ST=State/L=City/O=Organization/CN=${domain}"

# Update Nginx configuration
sed -i 's/ssl_certificate/#ssl_certificate/g' /etc/nginx/nginx.conf
sed -i 's/ssl_certificate_key/#ssl_certificate_key/g' /etc/nginx/nginx.conf
sed -i 's|#ssl_certificate|ssl_certificate /etc/nginx/ssl/cert.pem|g' /etc/nginx/nginx.conf
sed -i 's|#ssl_certificate_key|ssl_certificate_key /etc/nginx/ssl/key.pem|g' /etc/nginx/nginx.conf

# Restart Nginx
systemctl restart nginx

echo "✅ Self-signed SSL certificate configured!"
echo "📍 Certificate files:"
echo "   Certificate: /etc/nginx/ssl/cert.pem"
echo "   Key: /etc/nginx/ssl/key.pem"
echo ""
echo "⚠️ Note: Self-signed certificates will show security warnings in browsers."
echo "   For production, use Let's Encrypt instead.";
`;
    
    console.log('');
    colorLog('cyan', '📋 Generated SSL configuration script:');
    console.log('');
    colorLog('bright', sslScript);
    console.log('');
    
    colorLog('magenta', '🔧 To configure SSL on your server, run these commands:');
    console.log('');
    colorLog('bright', 'ssh root@YOUR_SERVER_IP');
    console.log('');
    colorLog('bright', '# Create and run SSL script');
    colorLog('bright', 'cat > setup-ssl.sh << \'EOF\'');
    colorLog('bright', sslScript);
    colorLog('bright', 'EOF');
    colorLog('bright', 'chmod +x setup-ssl.sh');
    colorLog('bright', './setup-ssl.sh');
    console.log('');
  }
  
  async monitorApplication() {
    colorLog('blue', '📊 Monitoring Application');
    colorLog('blue', '==================');
    console.log('');
    
    const domain = await askQuestion(`Enter your domain (default: ${this.config.domain}):`) || this.config.domain;
    
    console.log('');
    colorLog('cyan', '🔍 Real-time monitoring commands:');
    console.log('');
    colorLog('bright', '1. Check application status:');
    console.log(`   curl -f https://${domain}/health`);
    console.log('');
    colorLog('bright', '2. Check system resources:');
    console.log('   ssh root@YOUR_SERVER_IP "htop"');
    console.log('');
    colorLog('bright', '3. Check Docker containers:');
    console.log('   ssh root@YOUR_SERVER_IP "docker stats"');
    console.log('');
    colorLog('bright', '4. Check application logs:');
    console.log('   ssh root@YOUR_SERVER_IP "docker-compose logs -f app"');
    console.log('');
    colorLog('bright', '5. Check database logs:');
    console.log('   ssh root@YOUR_SERVER_IP "docker-compose logs -f postgres"');
    console.log('');
    colorLog('bright', '6. Check Redis logs:');
    console.log('   ssh root@YOUR_SERVER_IP "docker-compose logs -f redis"');
    console.log('');
    
    console.log('');
    colorLog('magenta', '📋 Monitoring Dashboard Setup:');
    console.log('');
    console.log('To set up a monitoring dashboard, run these commands:');
    console.log('');
    colorLog('bright', 'ssh root@YOUR_SERVER_IP');
    console.log('');
    colorLog('bright', '# Install monitoring tools');
    colorLog('bright', 'curl -fsSL https://raw.githubusercontent.com/kilo-code/ai-deployment-assistant/main/setup-monitoring.sh | bash');
    console.log('');
    colorLog('bright', '# Access monitoring dashboard');
    colorLog('bright', 'http://YOUR_SERVER_IP:3000');
    console.log('');
    
    await askQuestion('Press Enter to continue...');
    await this.showMenu();
  }
  
  async troubleshootIssues() {
    colorLog('blue', '🔧 Troubleshooting Issues');
    colorLog('blue', '=====================');
    console.log('');
    
    console.log('🤔 Choose the issue you\'re experiencing:');
    console.log('');
    console.log(`${colors.yellow}1.${colors.reset} Application not starting`);
    console.log(`${colors.yellow}2.${colors.reset} Database connection issues`);
    console.log(`${colors.yellow}3.${colors.reset} SSL certificate problems`);
    console.log(`${colors.yellow}4.${colors.reset} API endpoints not working`);
    console.log(`${colors.yellow}5.${colors.reset} Performance issues`);
    console.log(`${colors.yellow}6.${colors.reset} Security concerns`);
    console.log('');
    
    const issue = await askQuestion('Enter your choice (1-6):');
    
    switch (issue) {
      case '1':
        await this.troubleshootAppNotStarting();
        break;
      case '2':
        await this.troubleshootDatabaseIssues();
        break;
      case '3':
        await this.troubleshootSSLCertificate();
        break;
      case '4':
        await this.troubleshootAPIEndpoints();
        break;
      case '5':
        await this.troubleshootPerformance();
        break;
      case '6':
        await this.troubleshootSecurity();
        break;
      default:
        colorLog('red', '❌ Invalid choice. Please try again.');
        await this.troubleshootIssues();
        return;
    }
    
    await askQuestion('Press Enter to continue...');
    await this.showMenu();
  }
  
  async troubleshootAppNotStarting() {
    colorLog('red', '🚨 Application Not Starting');
    console.log('');
    
    console.log('🔍 Common causes and solutions:');
    console.log('');
    console.log('1. Check application logs:');
    console.log('   ssh root@YOUR_SERVER_IP "docker-compose logs app"');
    console.log('');
    console.log('2. Check if port is available:');
    console.log('   ssh root@YOUR_SERVER_IP "netstat -tulpn | grep 8080"');
    console.log('');
    console.log('3. Check Docker status:');
    console.log('   ssh root@YOUR_SERVER_IP "systemctl status docker"');
    console.log('');
    console.log('4. Restart services:');
    console.log('   ssh root@YOUR_SERVER_IP "docker-compose down && docker-compose up -d"');
    console.log('');
    console.log('5. Check disk space:');
    console.log('   ssh root@YOUR_SERVER_IP "df -h"');
    console.log('');
    
    console.log('🔧 Auto-fix script:');
    console.log('');
    console.log('Run this command on your server:');
    console.log('');
    colorLog('bright', 'curl -fsSL https://raw.githubusercontent.com/kilo-code/ai-deployment-assistant/main/fix-app-not-starting.sh | bash');
    console.log('');
  }
  
  async troubleshootDatabaseIssues() {
    colorLog('red', '🚨 Database Connection Issues');
    console.log('');
    
    console.log('🔍 Common causes and solutions:');
    console.log('');
    console.log('1. Check database logs:');
    console.log('   ssh root@YOUR_SERVER_IP "docker-compose logs postgres"');
    console.log('');
    console.log('2. Test database connection:');
    console.log('   ssh root@YOUR_SERVER_IP "docker exec -it nutrition-platform_postgres psql -U nutrition_user -d nutrition_platform -c \\\"SELECT 1;\\\""');
    console.log('');
    console.log('3. Check Redis connection:');
    console.log('   ssh root@YOUR_SERVER_IP "docker exec -it nutrition-platform_redis redis-cli -a YOUR_REDIS_PASSWORD ping"');
    console.log('');
    console.log('4. Verify environment variables:');
    console.log('   ssh root@YOUR_SERVER_IP "docker exec -it nutrition-platform_app env | grep DB_"');
    console.log('');
    
    console.log('🔧 Auto-fix script:');
    console.log('');
    console.log('Run this command on your server:');
    console.log('');
    colorLog('bright', 'curl -fsSL https://raw.githubusercontent.com/kilo-code/ai-deployment-assistant/main/fix-database-issues.sh | bash');
    console.log('');
  }
  
  async troubleshootSSLCertificate() {
    colorLog('red', '🚨 SSL Certificate Problems');
    console.log('');
    
    console.log('🔍 Common causes and solutions:');
    console.log('');
    console.log('1. Check certificate validity:');
    console.log(`   openssl s_client -connect ${this.config.domain}:443 -servername ${this.config.domain} < /dev/null`);
    console.log('');
    console.log('2. Check certificate files:');
    console.log('   ssh root@YOUR_SERVER_IP "ls -la /etc/nginx/ssl/"');
    console.log('');
    console.log('3. Check Nginx configuration:');
    console.log('   ssh root@YOUR_SERVER_IP "nginx -t"');
    console.log('');
    console.log('4. Check Nginx error logs:');
    console.log('   ssh root@YOUR_SERVER_IP "tail -f /var/log/nginx/error.log"');
    console.log('');
    
    console.log('🔧 Auto-fix script:');
    console.log('');
    console.log('Run this command on your server:');
    console.log('');
    colorLog('bright', 'curl -fsSL https://raw.githubusercontent.com/kilo-code/ai-deployment-assistant/main/fix-ssl-issues.sh | bash');
    console.log('');
  }
  
  async troubleshootAPIEndpoints() {
    colorLog('red', '🚨 API Endpoints Not Working');
    console.log('');
    
    console.log('🔍 Common causes and solutions:');
    console.log('');
    console.log('1. Check if application is running:');
    console.log('   ssh root@YOUR_SERVER_IP "docker-compose ps"');
    console.log('');
    console.log('2. Check application logs:');
    console.log('   ssh root@YOUR_SERVER_IP "docker-compose logs app"');
    console.log('');
    console.log('3. Test API endpoints locally:');
    console.log('   ssh root@YOUR_SERVER_IP "curl -f http://localhost:8080/health"');
    console.log('');
    console.log('4. Check CORS configuration:');
    console.log('   ssh root@YOUR_SERVER_IP "docker exec -it nutrition-platform_app env | grep CORS"');
    console.log('');
    
    console.log('🔧 Auto-fix script:');
    console.log('');
    console.log('Run this command on your server:');
    console.log('');
    colorLog('bright', 'curl -fsSL https://raw.githubusercontent.com/kilo-code/ai-deployment-assistant/main/fix-api-endpoints.sh | bash');
    console.log('');
  }
  
  async troubleshootPerformance() {
    colorLog('red', '🚨 Performance Issues');
    console.log('');
    
    console.log('🔍 Common causes and solutions:');
    console.log('');
    console.log('1. Check system resources:');
    console.log('   ssh root@YOUR_SERVER_IP "htop"');
    console.log('');
    console.log('2. Check Docker resource usage:');
    console.log('   ssh root@YOUR_SERVER_IP "docker stats"');
    console.log('');
    console.log('3. Check database performance:');
    console.log('   ssh root@YOUR_SERVER_IP "docker exec -it nutrition-platform_postgres psql -U nutrition_user -d nutrition_platform -c \\\"SELECT count(*) FROM pg_stat_activity;\\\""');
    console.log('');
    console.log('4. Check Redis performance:');
    console.log('   ssh root@YOUR_SERVER_IP "docker exec -it nutrition-platform_redis redis-cli info"');
    console.log('');
    
    console.log('🔧 Auto-fix script:');
    console.log('');
    console.log('Run this command on your server:');
    console.log('');
    colorLog('bright', 'curl -fsSL https://raw.githubusercontent.com/kilo-code/ai-deployment-assistant/main/fix-performance-issues.sh | bash');
    console.log('');
  }
  
  async troubleshootSecurity() {
    colorLog('red', '🚨 Security Concerns');
    console.log('');
    
    console.log('🔍 Security checklist:');
    console.log('');
    console.log('1. Check for open ports:');
    console.log('   ssh root@YOUR_SERVER_IP "netstat -tulpn"');
    console.log('');
    console.log('2. Check firewall rules:');
    console.log('   ssh root@YOUR_SERVER_IP "ufw status"');
    console.log('');
    console.log('3. Check for failed SSH attempts:');
    console.log('   ssh root@YOUR_SERVER_IP "grep \"Failed password\" /var/log/auth.log | tail -10"');
    console.log('');
    console.log('4. Check application vulnerabilities:');
    console.log('   ssh root@YOUR_SERVER_IP "docker run --rm -v /opt/nutrition-platform npm audit"');
    console.log('');
    
    console.log('🔧 Security hardening script:');
    console.log('');
    console.log('Run this command on your server:');
    console.log('');
    colorLog('bright', 'curl -fsSL https://raw.githubusercontent.com/kilo-code/ai-deployment-assistant/main/harden-security.sh | bash');
    console.log('');
  }
}

// Start the AI Deployment Assistant
const assistant = new AIDeploymentAssistant();
assistant.start().catch(console.error);