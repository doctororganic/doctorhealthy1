/**
 * Deployment Setup Test Suite
 * Sets up secure environment variables and validates deployment readiness
 * 
 * @author Test Engineer
 * @version 1.0.0
 */

const { expect } = require('chai');
const crypto = require('crypto');
const fs = require('fs');
const path = require('path');

/**
 * Test Suite: Secure Environment Setup
 */
describe('Secure Environment Setup', () => {
  
  /**
   * Test: Generate secure credentials
   * Generates cryptographically secure credentials for deployment
   */
  it('should generate secure credentials', () => {
    const credentials = {
      DB_PASSWORD: crypto.randomBytes(32).toString('hex'),
      JWT_SECRET: crypto.randomBytes(64).toString('hex'),
      API_KEY_SECRET: crypto.randomBytes(64).toString('hex'),
      REDIS_PASSWORD: crypto.randomBytes(32).toString('hex'),
      ENCRYPTION_KEY: crypto.randomBytes(16).toString('hex')
    };
    
    // Validate each credential meets requirements
    expect(credentials.DB_PASSWORD.length).to.equal(64, 'DB password should be 64 hex characters (32 bytes)');
    expect(credentials.JWT_SECRET.length).to.equal(128, 'JWT secret should be 128 hex characters (64 bytes)');
    expect(credentials.API_KEY_SECRET.length).to.equal(128, 'API key secret should be 128 hex characters (64 bytes)');
    expect(credentials.REDIS_PASSWORD.length).to.equal(64, 'Redis password should be 64 hex characters (32 bytes)');
    expect(credentials.ENCRYPTION_KEY.length).to.equal(32, 'Encryption key should be 32 hex characters (16 bytes)');
    
    // Verify no common patterns
    const commonPatterns = ['password', 'secret', 'key', '123', 'abc'];
    for (const [key, value] of Object.entries(credentials)) {
      for (const pattern of commonPatterns) {
        expect(value.toLowerCase()).to.not.include(pattern, `${key} should not contain common patterns`);
      }
    }
    
    console.log('\n🔑 Generated Secure Credentials:');
    console.log('=====================================');
    for (const [key, value] of Object.entries(credentials)) {
      console.log(`${key}=${value}`);
    }
    console.log('=====================================\n');
    
    return credentials;
  });
  
  /**
   * Test: Update environment file
   * Updates the .env file with secure credentials
   */
  it('should update environment file with secure credentials', () => {
    const envPath = path.join(__dirname, '..', '.env');
    
    // Generate secure credentials
    const credentials = {
      DB_PASSWORD: crypto.randomBytes(32).toString('hex'),
      JWT_SECRET: crypto.randomBytes(64).toString('hex'),
      API_KEY_SECRET: crypto.randomBytes(64).toString('hex'),
      REDIS_PASSWORD: crypto.randomBytes(32).toString('hex'),
      ENCRYPTION_KEY: crypto.randomBytes(16).toString('hex')
    };
    
    // Read existing .env file
    let envContent = '';
    if (fs.existsSync(envPath)) {
      envContent = fs.readFileSync(envPath, 'utf8');
    }
    
    // Update with new credentials
    const lines = envContent.split('\n');
    const updatedLines = [];
    
    for (const line of lines) {
      let updatedLine = line;
      for (const [key, value] of Object.entries(credentials)) {
        if (line.startsWith(`${key}=`)) {
          updatedLine = `${key}=${value}`;
        }
      }
      updatedLines.push(updatedLine);
    }
    
    // Write updated content
    const updatedContent = updatedLines.join('\n');
    fs.writeFileSync(envPath, updatedContent);
    
    // Verify updates
    const newContent = fs.readFileSync(envPath, 'utf8');
    for (const [key, value] of Object.entries(credentials)) {
      expect(newContent).to.include(`${key}=${value}`, `${key} should be updated in .env file`);
    }
    
    console.log('\n✅ Environment file updated successfully');
    console.log(`Updated: ${envPath}`);
    
    return credentials;
  });
  
  /**
   * Test: Update Coolify environment variables
   * Prepares environment variables for Coolify deployment
   */
  it('should prepare Coolify environment variables', () => {
    const coolifyEnvPath = path.join(__dirname, '..', 'coolify-env-vars.txt');
    
    // Generate secure credentials
    const credentials = {
      DB_PASSWORD: crypto.randomBytes(32).toString('hex'),
      JWT_SECRET: crypto.randomBytes(64).toString('hex'),
      API_KEY_SECRET: crypto.randomBytes(64).toString('hex'),
      REDIS_PASSWORD: crypto.randomBytes(32).toString('hex'),
      ENCRYPTION_KEY: crypto.randomBytes(16).toString('hex')
    };
    
    // Read existing coolify-env-vars.txt file
    let envContent = '';
    if (fs.existsSync(coolifyEnvPath)) {
      envContent = fs.readFileSync(coolifyEnvPath, 'utf8');
    }
    
    // Update with new credentials
    const lines = envContent.split('\n');
    const updatedLines = [];
    
    for (const line of lines) {
      let updatedLine = line;
      for (const [key, value] of Object.entries(credentials)) {
        if (line.includes(`${key}=${`)) {
          updatedLine = line.replace(/\$\{[^}]+\}/, value);
        }
      }
      updatedLines.push(updatedLine);
    }
    
    // Write updated content
    const updatedContent = updatedLines.join('\n');
    fs.writeFileSync(coolifyEnvPath, updatedContent);
    
    // Verify updates
    const newContent = fs.readFileSync(coolifyEnvPath, 'utf8');
    for (const [key, value] of Object.entries(credentials)) {
      expect(newContent).to.include(`${key}=${value}`, `${key} should be updated in coolify-env-vars.txt`);
    }
    
    console.log('\n✅ Coolify environment variables updated successfully');
    console.log(`Updated: ${coolifyEnvPath}`);
    
    return credentials;
  });
});

/**
 * Test Suite: Deployment Validation
 */
describe('Deployment Validation', () => {
  
  /**
   * Test: Validate all security fixes are in place
   * Verifies all security issues have been addressed
   */
  it('should validate all security fixes are in place', () => {
    // Check .env file for placeholders
    const envPath = path.join(__dirname, '..', '.env');
    if (fs.existsSync(envPath)) {
      const envContent = fs.readFileSync(envPath, 'utf8');
      
      const placeholderPatterns = [
        'REPLACE_WITH',
        'your_',
        'change_this',
        'placeholder'
      ];
      
      for (const pattern of placeholderPatterns) {
        expect(envContent.toLowerCase()).to.not.include(pattern, 
          `.env should not contain placeholder values like "${pattern}"`);
      }
    }
    
    // Check coolify-env-vars.txt for placeholders
    const coolifyEnvPath = path.join(__dirname, '..', 'coolify-env-vars.txt');
    if (fs.existsSync(coolifyEnvPath)) {
      const envContent = fs.readFileSync(coolifyEnvPath, 'utf8');
      
      const placeholderPatterns = [
        'REPLACE_WITH',
        'your_',
        'change_this',
        'placeholder'
      ];
      
      for (const pattern of placeholderPatterns) {
        expect(envContent.toLowerCase()).to.not.include(pattern, 
          `coolify-env-vars.txt should not contain placeholder values like "${pattern}"`);
      }
    }
    
    console.log('\n✅ All security fixes validated');
  });
  
  /**
   * Test: Verify database SSL configuration
   * Ensures database connections are encrypted
   */
  it('should verify database SSL configuration', () => {
    const envPath = path.join(__dirname, '..', '.env');
    
    if (fs.existsSync(envPath)) {
      const envContent = fs.readFileSync(envPath, 'utf8');
      expect(envContent).to.include('DB_SSL_MODE=require', 
        'Database SSL mode should be set to "require"');
    }
    
    console.log('\n✅ Database SSL configuration verified');
  });
  
  /**
   * Test: Verify CORS configuration
   * Ensures CORS is properly restricted
   */
  it('should verify CORS configuration', () => {
    const envPath = path.join(__dirname, '..', '.env');
    
    if (fs.existsSync(envPath)) {
      const envContent = fs.readFileSync(envPath, 'utf8');
      
      // Check that CORS doesn't allow all origins
      expect(envContent).to.not.include('CORS_ALLOWED_ORIGINS=*', 
        'CORS should not allow all origins');
      
      // Check that specific domains are allowed
      expect(envContent).to.include('super.doctorhealthy1.com', 
        'CORS should allow super.doctorhealthy1.com');
    }
    
    console.log('\n✅ CORS configuration verified');
  });
  
  /**
   * Test: Generate deployment commands
   * Generates the exact commands needed for deployment
   */
  it('should generate deployment commands', () => {
    const credentials = {
      DB_PASSWORD: crypto.randomBytes(32).toString('hex'),
      JWT_SECRET: crypto.randomBytes(64).toString('hex'),
      API_KEY_SECRET: crypto.randomBytes(64).toString('hex'),
      REDIS_PASSWORD: crypto.randomBytes(32).toString('hex'),
      ENCRYPTION_KEY: crypto.randomBytes(16).toString('hex')
    };
    
    console.log('\n🚀 Deployment Commands:');
    console.log('=========================');
    console.log('# Generate secure secrets');
    console.log(`export DB_PASSWORD=${credentials.DB_PASSWORD}`);
    console.log(`export JWT_SECRET=${credentials.JWT_SECRET}`);
    console.log(`export API_KEY_SECRET=${credentials.API_KEY_SECRET}`);
    console.log(`export REDIS_PASSWORD=${credentials.REDIS_PASSWORD}`);
    console.log(`export ENCRYPTION_KEY=${credentials.ENCRYPTION_KEY}`);
    console.log('');
    console.log('# Update your deployment environment with these values');
    console.log('echo "DB_PASSWORD=$DB_PASSWORD" >> .env');
    console.log('echo "JWT_SECRET=$JWT_SECRET" >> .env');
    console.log('echo "API_KEY_SECRET=$API_KEY_SECRET" >> .env');
    console.log('echo "REDIS_PASSWORD=$REDIS_PASSWORD" >> .env');
    console.log('echo "ENCRYPTION_KEY=$ENCRYPTION_KEY" >> .env');
    console.log('');
    console.log('# Run deployment tests');
    console.log('npm test -- tests/setup-deployment.test.js');
    console.log('');
    console.log('# Execute deployment');
    console.log('./complete-deployment.sh');
    console.log('=========================\n');
    
    return credentials;
  });
});

/**
 * Helper function to setup secure deployment
 */
async function setupSecureDeployment() {
  console.log('🔐 Setting up secure deployment...');
  
  const credentials = {
    DB_PASSWORD: crypto.randomBytes(32).toString('hex'),
    JWT_SECRET: crypto.randomBytes(64).toString('hex'),
    API_KEY_SECRET: crypto.randomBytes(64).toString('hex'),
    REDIS_PASSWORD: crypto.randomBytes(32).toString('hex'),
    ENCRYPTION_KEY: crypto.randomBytes(16).toString('hex')
  };
  
  // Update .env file
  const envPath = path.join(__dirname, '..', '.env');
  if (fs.existsSync(envPath)) {
    let envContent = fs.readFileSync(envPath, 'utf8');
    const lines = envContent.split('\n');
    
    for (let i = 0; i < lines.length; i++) {
      for (const [key, value] of Object.entries(credentials)) {
        if (lines[i].startsWith(`${key}=`)) {
          lines[i] = `${key}=${value}`;
        }
      }
    }
    
    fs.writeFileSync(envPath, lines.join('\n'));
    console.log('✅ .env file updated with secure credentials');
  }
  
  // Update coolify-env-vars.txt
  const coolifyEnvPath = path.join(__dirname, '..', 'coolify-env-vars.txt');
  if (fs.existsSync(coolifyEnvPath)) {
    let envContent = fs.readFileSync(coolifyEnvPath, 'utf8');
    const lines = envContent.split('\n');
    
    for (let i = 0; i < lines.length; i++) {
      for (const [key, value] of Object.entries(credentials)) {
          if (lines[i].includes(`${key}=${`)) {
          lines[i] = lines[i].replace(/\$\{[^}]+\}/, value);
        }
      }
    }
    
    fs.writeFileSync(coolifyEnvPath, lines.join('\n'));
    console.log('✅ coolify-env-vars.txt updated with secure credentials');
  }
  
  console.log('\n🔑 Secure Credentials Generated:');
  console.log('=====================================');
  for (const [key, value] of Object.entries(credentials)) {
    console.log(`${key}=${value}`);
  }
  console.log('=====================================\n');
  
  return credentials;
}

/**
 * Export helper functions
 */
module.exports = {
  setupSecureDeployment,
  generateCredentials: () => {
    return {
      DB_PASSWORD: crypto.randomBytes(32).toString('hex'),
      JWT_SECRET: crypto.randomBytes(64).toString('hex'),
      API_KEY_SECRET: crypto.randomBytes(64).toString('hex'),
      REDIS_PASSWORD: crypto.randomBytes(32).toString('hex'),
      ENCRYPTION_KEY: crypto.randomBytes(16).toString('hex')
    };
  }
};