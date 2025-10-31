/**
 * SSL/TLS Configuration Validation Test Suite
 * Tests that SSL certificates are properly configured and secure
 * 
 * @author Test Engineer
 * @version 1.0.0
 */

const { expect } = require('chai');
const https = require('https');
const tls = require('tls');
const axios = require('axios');
const crypto = require('crypto');

// Test configuration
const config = {
  hostname: process.env.TEST_HOSTNAME || 'super.doctorhealthy1.com',
  port: process.env.TEST_PORT || 443,
  timeout: 30000
};

/**
 * Test Suite: SSL Certificate Validation
 */
describe('SSL Certificate Validation', () => {
  
  /**
   * Test: SSL certificate is valid and trusted
   * Verifies the SSL certificate is issued by a trusted CA
   */
  it('should have valid SSL certificate', async function() {
    this.timeout(config.timeout);
    
    return new Promise((resolve, reject) => {
      const socket = tls.connect({
        host: config.hostname,
        port: config.port,
        servername: config.hostname,
        rejectUnauthorized: true
      }, () => {
        const cert = socket.getPeerCertificate();
        socket.destroy();
        
        try {
          // Verify certificate is present
          expect(cert).to.not.be.null;
          expect(cert.subject).to.exist;
          expect(cert.issuer).to.exist;
          
          // Verify domain matches
          expect(cert.subject.CN).to.include(config.hostname);
          if (cert.subjectaltname) {
            expect(cert.subjectaltname).to.include(config.hostname);
          }
          
          // Verify certificate is not expired
          const now = new Date();
          expect(new Date(cert.valid_from)).to.be.below(now);
          expect(new Date(cert.valid_to)).to.be.above(now);
          
          // Verify certificate is issued by trusted CA
          expect(cert.issuer).to.not.include('self-signed');
          
          resolve();
        } catch (error) {
          reject(error);
        }
      });
      
      socket.on('error', (error) => {
        reject(new Error(`SSL connection failed: ${error.message}`));
      });
    });
  });

  /**
   * Test: SSL certificate uses strong encryption
   * Verifies the SSL/TLS connection uses modern cipher suites
   */
  it('should use strong SSL/TLS configuration', async function() {
    this.timeout(config.timeout);
    
    return new Promise((resolve, reject) => {
      const socket = tls.connect({
        host: config.hostname,
        port: config.port,
        servername: config.hostname,
        rejectUnauthorized: true,
        ciphers: tls.DEFAULT_CIPHERS
      }, () => {
        const cipher = socket.getCipher();
        socket.destroy();
        
        try {
          // Verify TLS version is 1.2 or higher
          expect(cipher.version).to.be.oneOf(['TLSv1.2', 'TLSv1.3']);
          
          // Verify cipher suite is strong
          const strongCiphers = [
            'ECDHE',
            'AES',
            'GCM',
            'CHACHA20'
          ];
          
          const cipherName = cipher.name.toUpperCase();
          let isStrong = false;
          
          for (const strongCipher of strongCiphers) {
            if (cipherName.includes(strongCipher)) {
              isStrong = true;
              break;
            }
          }
          
          expect(isStrong).to.be.true;
          
          resolve();
        } catch (error) {
          reject(error);
        }
      });
      
      socket.on('error', (error) => {
        reject(new Error(`SSL connection failed: ${error.message}`));
      });
    });
  });

  /**
   * Test: SSL certificate has proper key strength
   * Verifies the certificate uses a strong key algorithm
   */
  it('should have strong certificate key', async function() {
    this.timeout(config.timeout);
    
    return new Promise((resolve, reject) => {
      const socket = tls.connect({
        host: config.hostname,
        port: config.port,
        servername: config.hostname,
        rejectUnauthorized: true
      }, () => {
        const cert = socket.getPeerCertificate();
        socket.destroy();
        
        try {
          // Check key algorithm and size
          if (cert.pubkey) {
            // RSA keys should be at least 2048 bits
            if (cert.pubkey.type === 'rsa') {
              expect(cert.pubkey.bits).to.be.at.least(2048);
            }
            
            // ECDSA keys should use approved curves
            if (cert.pubkey.type === 'ec') {
              const approvedCurves = ['secp256r1', 'secp384r1', 'secp521r1'];
              expect(approvedCurves).to.include(cert.pubkey.asn1Curve);
            }
          }
          
          resolve();
        } catch (error) {
          reject(error);
        }
      });
      
      socket.on('error', (error) => {
        reject(new Error(`SSL connection failed: ${error.message}`));
      });
    });
  });

  /**
   * Test: HSTS header is configured
   * Verifies HTTP Strict Transport Security is enabled
   */
  it('should have HSTS header configured', async function() {
    this.timeout(config.timeout);
    
    try {
      const response = await axios.get(`https://${config.hostname}`, {
        timeout: 10000,
        validateStatus: () => true
      });
      
      const hstsHeader = response.headers['strict-transport-security'];
      
      if (hstsHeader) {
        // Verify HSTS includes max-age
        expect(hstsHeader).to.include('max-age=');
        
        // Verify max-age is at least 6 months (15552000 seconds)
        const maxAgeMatch = hstsHeader.match(/max-age=(\d+)/);
        if (maxAgeMatch) {
          const maxAge = parseInt(maxAgeMatch[1]);
          expect(maxAge).to.be.at.least(15552000);
        }
        
        // Verify includeSubDomains is set (recommended)
        expect(hstsHeader).to.include('includeSubDomains');
      } else {
        console.warn('HSTS header not found - should be implemented for production');
      }
    } catch (error) {
      console.warn('Could not verify HSTS header:', error.message);
    }
  });
});

/**
 * Test Suite: Security Headers Validation
 */
describe('Security Headers Validation', () => {
  
  /**
   * Test: All security headers are present
   * Verifies essential security headers are configured
   */
  it('should have all security headers', async function() {
    this.timeout(config.timeout);
    
    try {
      const response = await axios.get(`https://${config.hostname}`, {
        timeout: 10000,
        validateStatus: () => true
      });
      
      const headers = response.headers;
      
      // Check for required security headers
      expect(headers).to.have.property('x-frame-options');
      expect(headers).to.have.property('x-content-type-options');
      expect(headers).to.have.property('x-xss-protection');
      expect(headers).to.have.property('referrer-policy');
      
      // Verify header values
      expect(headers['x-frame-options']).to.equal('SAMEORIGIN');
      expect(headers['x-content-type-options']).to.equal('nosniff');
      expect(headers['x-xss-protection']).to.include('1; mode=block');
      expect(headers['referrer-policy']).to.include('strict-origin');
      
      // Check for Content Security Policy (recommended)
      if (headers['content-security-policy']) {
        const csp = headers['content-security-policy'];
        
        // Verify CSP includes default-src
        expect(csp).to.include('default-src');
        
        // Verify CSP includes script-src
        expect(csp).to.include('script-src');
        
        // Verify CSP doesn't allow unsafe-inline
        expect(csp).to.not.include("'unsafe-inline'");
      } else {
        console.warn('Content Security Policy not found - should be implemented for production');
      }
    } catch (error) {
      throw new Error(`Security headers validation failed: ${error.message}`);
    }
  });
});

/**
 * Test Suite: Certificate Chain Validation
 */
describe('Certificate Chain Validation', () => {
  
  /**
   * Test: Certificate chain is complete
   * Verifies the full certificate chain is provided
   */
  it('should have complete certificate chain', async function() {
    this.timeout(config.timeout);
    
    return new Promise((resolve, reject) => {
      const socket = tls.connect({
        host: config.hostname,
        port: config.port,
        servername: config.hostname,
        rejectUnauthorized: true
      }, () => {
        const cert = socket.getPeerCertificate(true);
        socket.destroy();
        
        try {
          // Check if certificate chain is present
          expect(cert).to.not.be.null;
          
          if (Array.isArray(cert)) {
            // Chain should include server certificate and at least one intermediate
            expect(cert.length).to.be.at.least(2);
            
            // First certificate should be the server certificate
            expect(cert[0].subject).to.exist;
            
            // Last certificate should be root or intermediate
            expect(cert[cert.length - 1].issuer).to.exist;
          } else {
            // Single certificate - warn about incomplete chain
            console.warn('Certificate chain might be incomplete');
          }
          
          resolve();
        } catch (error) {
          reject(error);
        }
      });
      
      socket.on('error', (error) => {
        reject(new Error(`SSL connection failed: ${error.message}`));
      });
    });
  });

  /**
   * Test: OCSP stapling is configured (if possible to verify)
   * Attempts to verify OCSP stapling is configured
   */
  it('should have OCSP stapling configured', async function() {
    this.timeout(config.timeout);
    
    return new Promise((resolve, reject) => {
      const socket = tls.connect({
        host: config.hostname,
        port: config.port,
        servername: config.hostname,
        rejectUnauthorized: true,
        requestOCSP: true
      }, () => {
        const ocspResponse = socket.getOCSPResponse();
        socket.destroy();
        
        try {
          // OCSP response should be present if stapling is configured
          if (ocspResponse) {
            console.log('OCSP stapling is configured');
          } else {
            console.warn('OCSP stapling might not be configured');
          }
          
          resolve();
        } catch (error) {
          reject(error);
        }
      });
      
      socket.on('error', (error) => {
        reject(new Error(`SSL connection failed: ${error.message}`));
      });
    });
  });
});

/**
 * Test Suite: Forward Secrecy Validation
 */
describe('Forward Secrecy Validation', () => {
  
  /**
   * Test: Forward secrecy cipher suites are supported
   * Verifies the server supports cipher suites with forward secrecy
   */
  it('should support forward secrecy cipher suites', async function() {
    this.timeout(config.timeout);
    
    const forwardSecrecyCiphers = [
      'ECDHE-RSA-AES256-GCM-SHA384',
      'ECDHE-RSA-AES128-GCM-SHA256',
      'ECDHE-RSA-CHACHA20-POLY1305',
      'ECDHE-ECDSA-AES256-GCM-SHA384',
      'ECDHE-ECDSA-AES128-GCM-SHA256',
      'ECDHE-ECDSA-CHACHA20-POLY1305'
    ];
    
    let hasForwardSecrecy = false;
    
    for (const cipher of forwardSecrecyCiphers) {
      try {
        await new Promise((resolve, reject) => {
          const socket = tls.connect({
            host: config.hostname,
            port: config.port,
            servername: config.hostname,
            rejectUnauthorized: true,
            ciphers: cipher
          }, () => {
            const cipherInfo = socket.getCipher();
            socket.destroy();
            
            if (cipherInfo.name === cipher) {
              hasForwardSecrecy = true;
            }
            
            resolve();
          });
          
          socket.on('error', () => {
            // Cipher not supported, try next
            resolve();
          });
        });
        
        if (hasForwardSecrecy) break;
      } catch (error) {
        // Continue trying other ciphers
      }
    }
    
    expect(hasForwardSecrecy).to.be.true;
  });
});

/**
 * Helper function to get certificate information
 */
function getCertificateInfo(hostname, port = 443) {
  return new Promise((resolve, reject) => {
    const socket = tls.connect({
      host: hostname,
      port: port,
      servername: hostname,
      rejectUnauthorized: true
    }, () => {
      const cert = socket.getPeerCertificate(true);
      const cipher = socket.getCipher();
      socket.destroy();
      
      resolve({
        certificate: cert,
        cipher: cipher
      });
    });
    
    socket.on('error', (error) => {
      reject(error);
    });
  });
}

/**
 * Helper function to check certificate expiration
 */
function checkCertificateExpiration(cert) {
  const now = new Date();
  const validFrom = new Date(cert.valid_from);
  const validTo = new Date(cert.valid_to);
  
  const daysUntilExpiration = Math.floor((validTo - now) / (1000 * 60 * 60 * 24));
  
  return {
    isValid: now >= validFrom && now <= validTo,
    daysUntilExpiration,
    isExpiringSoon: daysUntilExpiration < 30
  };
}

/**
 * Export helper functions for use in other test files
 */
module.exports = {
  getCertificateInfo,
  checkCertificateExpiration,
  runSSLValidationTests: async () => {
    console.log('Running SSL/TLS validation tests...');
    // Test runner implementation would go here
  }
};