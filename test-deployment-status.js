/**
 * Simple deployment status test
 * Checks if the application is deployed and accessible
 */

const https = require('https');
const http = require('http');

const domains = [
  'https://super.doctorhealthy1.com',
  'https://my.doctorhealthy1.com',
  'http://super.doctorhealthy1.com',
  'http://my.doctorhealthy1.com'
];

async function checkDeployment() {
  console.log('🔍 Checking deployment status...\n');
  
  for (const domain of domains) {
    try {
      console.log(`Testing: ${domain}`);
      
      const response = await new Promise((resolve, reject) => {
        const protocol = domain.startsWith('https') ? https : http;
        const req = protocol.request(domain, { timeout: 10000 }, (res) => {
          let data = '';
          res.on('data', chunk => data += chunk);
          res.on('end', () => resolve({ status: res.statusCode, data }));
        });
        
        req.on('error', reject);
        req.on('timeout', () => reject(new Error('Request timeout')));
        req.end();
      });
      
      console.log(`✅ Status: ${response.status}`);
      
      if (response.status === 200) {
        console.log('🎉 Application is LIVE and accessible!');
        console.log(`📍 URL: ${domain}`);
        
        // Check for health endpoint
        try {
          const healthResponse = await new Promise((resolve, reject) => {
            const protocol = domain.startsWith('https') ? https : http;
            const req = protocol.request(`${domain}/health`, { timeout: 5000 }, (res) => {
              let data = '';
              res.on('data', chunk => data += chunk);
              res.on('end', () => resolve({ status: res.statusCode, data }));
            });
            
            req.on('error', reject);
            req.on('timeout', () => reject(new Error('Request timeout')));
            req.end();
          });
          
          if (healthResponse.status === 200) {
            console.log('✅ Health endpoint is working');
            console.log(`📊 Response: ${healthResponse.data.trim()}`);
          }
        } catch (e) {
          console.log('⚠️ Health endpoint not accessible');
        }
      }
      
      console.log('---');
    } catch (error) {
      console.log(`❌ Error: ${error.message}`);
      console.log('---');
    }
  }
}

checkDeployment().catch(console.error);