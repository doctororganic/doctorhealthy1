/**
 * Monitor deployment progress
 * Continuously checks deployment status
 */

const https = require('https');
const http = require('http');

const domains = [
  { url: 'https://super.doctorhealthy1.com', name: 'Super DoctorHealthy1 (HTTPS)' },
  { url: 'https://my.doctorhealthy1.com', name: 'My DoctorHealthy1 (HTTPS)' },
  { url: 'http://super.doctorhealthy1.com', name: 'Super DoctorHealthy1 (HTTP)' },
  { url: 'http://my.doctorhealthy1.com', name: 'My DoctorHealthy1 (HTTP)' }
];

async function checkSite(url, rejectUnauthorized = false) {
  return new Promise((resolve, reject) => {
    const protocol = url.startsWith('https') ? https : http;
    const req = protocol.request(url, { 
      timeout: 10000,
      rejectUnauthorized: !rejectUnauthorized
    }, (res) => {
      let data = '';
      res.on('data', chunk => data += chunk);
      res.on('end', () => resolve({ 
        status: res.statusCode, 
        headers: res.headers,
        data: data.trim()
      }));
    });
    
    req.on('error', reject);
    req.on('timeout', () => reject(new Error('Request timeout')));
    req.end();
  });
}

async function monitorDeployment() {
  console.log('🚀 Monitoring Nutrition Platform Deployment...\n');
  
  let deploymentComplete = false;
  let attempts = 0;
  const maxAttempts = 30;
  
  while (!deploymentComplete && attempts < maxAttempts) {
    attempts++;
    console.log(`\n📊 Check attempt ${attempts}/${maxAttempts}`);
    console.log('=====================================');
    
    for (const site of domains) {
      try {
        console.log(`\n🔍 Testing ${site.name}`);
        
        // Test main URL
        const response = await checkSite(site.url, true);
        console.log(`✅ Status: ${response.status}`);
        
        if (response.status === 200) {
          console.log('🎉 Application is LIVE!');
          console.log(`📍 URL: ${site.url}`);
          
          // Check content type
          if (response.headers['content-type']) {
            console.log(`📄 Content-Type: ${response.headers['content-type']}`);
          }
          
          // Check for health endpoint
          try {
            const healthResponse = await checkSite(`${site.url}/health`, true);
            if (healthResponse.status === 200) {
              console.log('✅ Health endpoint working');
              console.log(`📊 Health: ${healthResponse.data}`);
            }
          } catch (e) {
            console.log('⚠️ Health endpoint not accessible');
          }
          
          // Check for API endpoint
          try {
            const apiResponse = await checkSite(`${site.url}/api/info`, true);
            if (apiResponse.status === 200) {
              console.log('✅ API endpoint working');
            }
          } catch (e) {
            console.log('⚠️ API endpoint not accessible');
          }
          
          deploymentComplete = true;
        } else if (response.status === 404) {
          console.log('⚠️ Server responding but application not yet deployed');
        } else if (response.status === 502 || response.status === 503) {
          console.log('⚠️ Server temporarily unavailable');
        } else {
          console.log(`⚠️ Unexpected status: ${response.status}`);
        }
      } catch (error) {
        if (error.message.includes('self-signed certificate')) {
          console.log('🔒 SSL configured (self-signed certificate)');
        } else if (error.message.includes('ECONNREFUSED')) {
          console.log('❌ Connection refused (server not ready)');
        } else {
          console.log(`❌ Error: ${error.message}`);
        }
      }
    }
    
    if (!deploymentComplete) {
      console.log('\n⏳ Waiting 30 seconds before next check...');
      await new Promise(resolve => setTimeout(resolve, 30000));
    }
  }
  
  console.log('\n=====================================');
  if (deploymentComplete) {
    console.log('🎉 DEPLOYMENT SUCCESSFUL!');
    console.log('\n📋 Next Steps:');
    console.log('1. Test all application features');
    console.log('2. Verify SSL certificate (may take 5-15 minutes for Let\'s Encrypt)');
    console.log('3. Monitor application performance');
    console.log('4. Set up monitoring and alerts');
  } else {
    console.log('⚠️ Deployment still in progress');
    console.log('Please check the deployment logs for more information');
  }
}

// Run monitoring
monitorDeployment().catch(console.error);