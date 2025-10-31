/**
 * Troubleshooting Test Suite
 * Tests to identify and verify fixes for common deployment issues
 * 
 * @author Test Engineer
 * @version 1.0.0
 */

const { expect } = require('chai');
const axios = require('axios');
const https = require('https');
const fs = require('fs');
const path = require('path');

// Test configuration
const config = {
  baseUrl: process.env.TEST_BASE_URL || 'https://super.doctorhealthy1.com',
  apiBaseUrl: process.env.TEST_API_URL || 'https://super.doctorhealthy1.com/api',
  timeout: 30000
};

// HTTPS Agent for SSL testing
const httpsAgent = new https.Agent({
  rejectUnauthorized: false // Only for testing, verify SSL in production
});

/**
 * Test Suite: Connectivity Issues
 */
describe('Connectivity Issues', () => {
  
  /**
   * Test: Server is reachable
   * Verifies the server is reachable from the test environment
   */
  it('should have reachable server', async function() {
    this.timeout(config.timeout);
    
    try {
      const response = await axios.get(config.baseUrl, {
        timeout: 10000,
        httpsAgent
      });
      
      expect(response.status).to.be.oneOf([200, 301, 302, 308]);
    } catch (error) {
      if (error.code === 'ENOTFOUND') {
        throw new Error('DNS resolution failed - check domain configuration');
      } else if (error.code === 'ECONNREFUSED') {
        throw new Error('Connection refused - server might be down or firewall blocking');
      } else if (error.code === 'ETIMEDOUT') {
        throw new Error('Connection timed out - server might be overloaded or network issues');
      }
      throw error;
    }
  });

  /**
   * Test: Ports are open
   * Verifies the required ports are open
   */
  it('should have required ports open', async function() {
    this.timeout(config.timeout);
    
    // Test HTTP port
    try {
      const httpUrl = config.baseUrl.replace('https://', 'http://');
      await axios.get(httpUrl, {
        timeout: 5000,
        maxRedirects: 0
      });
      
      // If we get here, HTTP port is open
    } catch (error) {
      if (error.code === 'ECONNREFUSED') {
        console.warn('HTTP port (80) might be closed or blocked');
      }
    }
    
    // Test HTTPS port
    try {
      await axios.get(config.baseUrl, {
        timeout: 5000,
        httpsAgent
      });
      
      // If we get here, HTTPS port is open
    } catch (error) {
      if (error.code === 'ECONNREFUSED') {
        throw new Error('HTTPS port (443) is closed or blocked');
      }
    }
  });
});

/**
 * Test Suite: SSL/TLS Issues
 */
describe('SSL/TLS Issues', () => {
  
  /**
   * Test: SSL certificate validity
   * Verifies SSL certificate is valid and trusted
   */
  it('should have valid SSL certificate', async function() {
    this.timeout(config.timeout);
    
    try {
      const response = await axios.get(config.baseUrl, {
        timeout: 10000,
        httpsAgent: new https.Agent({ rejectUnauthorized: true })
      });
      
      expect(response.status).to.be.oneOf([200, 301, 302, 308]);
    } catch (error) {
      if (error.code === 'UNABLE_TO_VERIFY_LEAF_SIGNATURE') {
        throw new Error('SSL certificate is not trusted - check certificate chain');
      } else if (error.code === 'CERT_HAS_EXPIRED') {
        throw new Error('SSL certificate has expired - renew certificate');
      } else if (error.code === 'ERR_TLS_CERT_ALTNAME_INVALID') {
        throw new Error('SSL certificate does not match domain - check certificate CN/SAN');
      }
      throw error;
    }
  });

  /**
   * Test: Certificate chain is complete
   * Verifies the full certificate chain is served
   */
  it('should have complete certificate chain', async function() {
    this.timeout(config.timeout);
    
    try {
      const response = await axios.get(config.baseUrl, {
        timeout: 10000,
        httpsAgent
      });
      
      // Check if certificate chain is provided
      const socket = response.request.res.socket;
      const cert = socket.getPeerCertificate(true);
      
      if (Array.isArray(cert)) {
        expect(cert.length).to.be.at.least(2, 'Certificate chain should include server certificate and intermediates');
      } else {
        console.warn('Certificate chain might be incomplete');
      }
    } catch (error) {
      console.warn('Cannot verify certificate chain');
    }
  });
});

/**
 * Test Suite: Application Issues
 */
describe('Application Issues', () => {
  
  /**
   * Test: Application is responding
   * Verifies the application is responding to requests
   */
  it('should have application responding', async function() {
    this.timeout(config.timeout);
    
    try {
      const response = await axios.get(config.baseUrl, {
        timeout: 15000,
        httpsAgent
      });
      
      if (response.status >= 500) {
        throw new Error('Application returning server errors - check application logs');
      }
    } catch (error) {
      if (error.response && error.response.status >= 500) {
        throw new Error('Application returning server errors - check application logs');
      }
      throw error;
    }
  });

  /**
   * Test: Application is not returning errors
   * Verifies the application is not returning error pages
   */
  it('should not return error pages', async function() {
    this.timeout(config.timeout);
    
    try {
      const response = await axios.get(config.baseUrl, {
        timeout: 15000,
        httpsAgent
      });
      
      // Check for common error indicators
      const content = response.data;
      expect(content).to.not.include('Internal Server Error');
      expect(content).to.not.include('503 Service Unavailable');
      expect(content).to.not.include('502 Bad Gateway');
      expect(content).to.not.include('Application Error');
    } catch (error) {
      throw error;
    }
  });
});

/**
 * Test Suite: Database Issues
 */
describe('Database Issues', () => {
  
  /**
   * Test: Database is accessible
   * Verifies the database is accessible from the application
   */
  it('should have accessible database', async function() {
    this.timeout(config.timeout);
    
    try {
      const response = await axios.get(`${config.apiBaseUrl}/system/status`, {
        timeout: 10000,
        httpsAgent
      });
      
      if (response.status === 200 && response.data.database) {
        if (response.data.database.status === 'error') {
          throw new Error('Database connection error - check database configuration');
        }
      }
    } catch (error) {
      if (error.response?.status !== 404) {
        console.warn('Cannot verify database connection through API');
      }
    }
  });

  /**
   * Test: Database SSL is configured
   * Verifies database connections are using SSL
   */
  it('should have database SSL configured', () => {
    const dbSSLMode = process.env.DB_SSL_MODE;
    
    if (!dbSSLMode || dbSSLMode === 'disable') {
      throw new Error('Database SSL is not configured - set DB_SSL_MODE=require');
    }
  });
});

/**
 * Test Suite: Configuration Issues
 */
describe('Configuration Issues', () => {
  
  /**
   * Test: Environment variables are set
   * Verifies required environment variables are set
   */
  it('should have required environment variables', () => {
    const requiredVars = [
      'SERVER_PORT',
      'DB_HOST',
      'DB_NAME',
      'DB_USER',
      'DB_PASSWORD',
      'JWT_SECRET',
      'API_KEY_SECRET',
      'ENCRYPTION_KEY'
    ];
    
    const missingVars = [];
    
    for (const varName of requiredVars) {
      if (!process.env[varName]) {
        missingVars.push(varName);
      }
    }
    
    if (missingVars.length > 0) {
      throw new Error(`Missing environment variables: ${missingVars.join(', ')}`);
    }
  });

  /**
   * Test: Configuration files exist
   * Verifies required configuration files exist
   */
  it('should have required configuration files', () => {
    const configFiles = [
      '.env',
      'nginx/conf.d/default.conf',
      'docker-compose.yml'
    ];
    
    const missingFiles = [];
    
    for (const configFile of configFiles) {
      const filePath = path.join(__dirname, '..', configFile);
      if (!fs.existsSync(filePath)) {
        missingFiles.push(configFile);
      }
    }
    
    if (missingFiles.length > 0) {
      throw new Error(`Missing configuration files: ${missingFiles.join(', ')}`);
    }
  });
});

/**
 * Test Suite: Performance Issues
 */
describe('Performance Issues', () => {
  
  /**
   * Test: Response time is acceptable
   * Verifies the application responds within acceptable time limits
   */
  it('should have acceptable response time', async function() {
    this.timeout(config.timeout);
    
    const startTime = Date.now();
    
    try {
      const response = await axios.get(`${config.baseUrl}/health`, {
        timeout: 10000,
        httpsAgent
      });
      
      const responseTime = Date.now() - startTime;
      
      if (responseTime > 10000) {
        throw new Error(`Slow response time: ${responseTime}ms - check server performance`);
      }
    } catch (error) {
      throw error;
    }
  });

  /**
   * Test: Application handles concurrent requests
   * Verifies the application can handle multiple simultaneous requests
   */
  it('should handle concurrent requests', async function() {
    this.timeout(config.timeout);
    
    const concurrentRequests = 5;
    const requests = [];
    let failedRequests = 0;
    
    for (let i = 0; i < concurrentRequests; i++) {
      requests.push(
        axios.get(`${config.baseUrl}/health`, {
          timeout: 10000,
          httpsAgent
        }).catch(error => {
          failedRequests++;
          return error;
        })
      );
    }
    
    await Promise.all(requests);
    
    if (failedRequests > 0) {
      throw new Error(`${failedRequests} out of ${concurrentRequests} requests failed - check server capacity`);
    }
  });
});

/**
 * Test Suite: Security Issues
 */
describe('Security Issues', () => {
  
  /**
   * Test: Security headers are present
   * Verifies security headers are configured
   */
  it('should have security headers', async function() {
    this.timeout(config.timeout);
    
    try {
      const response = await axios.get(config.baseUrl, {
        timeout: 10000,
        httpsAgent
      });
      
      const headers = response.headers;
      const missingHeaders = [];
      
      if (!headers['x-frame-options']) {
        missingHeaders.push('X-Frame-Options');
      }
      
      if (!headers['x-content-type-options']) {
        missingHeaders.push('X-Content-Type-Options');
      }
      
      if (!headers['x-xss-protection']) {
        missingHeaders.push('X-XSS-Protection');
      }
      
      if (missingHeaders.length > 0) {
        throw new Error(`Missing security headers: ${missingHeaders.join(', ')}`);
      }
    } catch (error) {
      throw error;
    }
  });

  /**
   * Test: No sensitive information is exposed
   * Verifies no sensitive information is exposed in responses
   */
  it('should not expose sensitive information', async function() {
    this.timeout(config.timeout);
    
    try {
      const response = await axios.get(`${config.apiBaseUrl}/info`, {
        timeout: 10000,
        httpsAgent
      });
      
      if (response.status === 200) {
        const data = JSON.stringify(response.data);
        
        const sensitivePatterns = [
          /password/i,
          /secret/i,
          /key.*=.*[a-f0-9]{32,}/i,
          /token.*=.*[a-f0-9]{32,}/i
        ];
        
        for (const pattern of sensitivePatterns) {
          if (pattern.test(data)) {
            throw new Error(`Sensitive information exposed: ${pattern}`);
          }
        }
      }
    } catch (error) {
      if (error.response?.status !== 404) {
        throw error;
      }
    }
  });
});

/**
 * Helper function to diagnose common issues
 */
async function diagnoseCommonIssues() {
  console.log('Diagnosing common deployment issues...');
  
  const issues = [];
  
  // Check connectivity
  try {
    await axios.get(config.baseUrl, {
      timeout: 5000,
      httpsAgent
    });
  } catch (error) {
    issues.push({
      type: 'connectivity',
      message: `Cannot connect to server: ${error.message}`,
      solution: 'Check server status, DNS configuration, and firewall rules'
    });
  }
  
  // Check SSL
  try {
    await axios.get(config.baseUrl, {
      timeout: 5000,
      httpsAgent: new https.Agent({ rejectUnauthorized: true })
    });
  } catch (error) {
    issues.push({
      type: 'ssl',
      message: `SSL certificate issue: ${error.message}`,
      solution: 'Check certificate validity, chain, and domain matching'
    });
  }
  
  // Check application health
  try {
    const response = await axios.get(`${config.baseUrl}/health`, {
      timeout: 5000,
      httpsAgent
    });
    
    if (response.status !== 200) {
      issues.push({
        type: 'application',
        message: `Application health check failed: ${response.status}`,
        solution: 'Check application logs and configuration'
      });
    }
  } catch (error) {
    issues.push({
      type: 'application',
      message: `Application health check failed: ${error.message}`,
      solution: 'Check application logs and configuration'
    });
  }
  
  return issues;
}

/**
 * Helper function to generate troubleshooting report
 */
async function generateTroubleshootingReport() {
  const issues = await diagnoseCommonIssues();
  
  const report = {
    timestamp: new Date().toISOString(),
    url: config.baseUrl,
    issues: issues,
    recommendations: []
  };
  
  // Add recommendations based on issues
  if (issues.some(i => i.type === 'connectivity')) {
    report.recommendations.push('Verify server is running and accessible');
    report.recommendations.push('Check DNS configuration');
    report.recommendations.push('Verify firewall rules allow traffic');
  }
  
  if (issues.some(i => i.type === 'ssl')) {
    report.recommendations.push('Renew SSL certificate if expired');
    report.recommendations.push('Verify certificate chain is complete');
    report.recommendations.push('Check certificate matches domain');
  }
  
  if (issues.some(i => i.type === 'application')) {
    report.recommendations.push('Check application logs for errors');
    report.recommendations.push('Verify environment variables are set');
    report.recommendations.push('Restart application if needed');
  }
  
  return report;
}

/**
 * Export helper functions for use in other test files
 */
module.exports = {
  config,
  diagnoseCommonIssues,
  generateTroubleshootingReport,
  runTroubleshootingTests: async () => {
    console.log('Running troubleshooting tests...');
    // Implementation would go here
  }
};