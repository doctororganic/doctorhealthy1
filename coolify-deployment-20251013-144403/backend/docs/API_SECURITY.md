# API Security Documentation

## Overview

This document outlines the comprehensive security measures implemented for the Nutrition Platform API key management system. The implementation follows industry best practices and security standards to ensure robust protection against common vulnerabilities.

## Security Architecture

### 1. API Key Generation

#### Cryptographically Secure Random Generation
- Uses Go's `crypto/rand` package for cryptographically secure random number generation
- Minimum key length of 32 characters (recommended 64+ characters)
- Character set includes alphanumeric characters plus safe special characters (`-`, `_`)
- Prefixed keys for easy identification and categorization (`nk_`, `np_`, `nt_`, `na_`)

#### Key Format
```
[prefix]_[random_string]
Example: nk_8Kf9mN2pQ7rS3tU6vW9xY1zA4bC7dE0fG3hI5j
```

### 2. Key Storage and Hashing

#### Secure Storage
- API keys are **never stored in plaintext**
- SHA-256 hashing applied before database storage
- Original keys are only shown once during creation
- Database stores only the hash, metadata, and permissions

#### Hash Validation
- Constant-time comparison to prevent timing attacks
- Custom implementation of constant-time string comparison
- Protection against side-channel attacks

### 3. Authentication and Authorization

#### Multiple Authentication Methods
1. **X-API-Key Header**: `X-API-Key: nk_your_api_key_here`
2. **Authorization Header**: `Authorization: Bearer nk_your_api_key_here`
3. **Query Parameter**: `?api_key=nk_your_api_key_here` (not recommended for production)

#### Scope-Based Authorization
- Granular permission system with predefined scopes
- Available scopes:
  - `nutrition:read` - Read access to nutrition data
  - `nutrition:write` - Write access to nutrition data
  - `admin:read` - Read access to admin functions
  - `admin:write` - Write access to admin functions
  - `analytics:read` - Read access to analytics data

#### Endpoint Protection Levels
1. **Public Endpoints**: No authentication required
2. **Optional API Key**: Enhanced features with API key
3. **Required API Key**: Mandatory authentication
4. **Admin Only**: Requires admin scope
5. **Specific Scopes**: Requires particular permissions

### 4. Rate Limiting

#### Multi-Level Rate Limiting
- **Per-API-Key Limits**: Configurable per key (default: 1000 requests/hour)
- **Global Rate Limiting**: Overall system protection
- **Endpoint-Specific Limits**: Different limits for different endpoints
- **Burst Protection**: Prevents sudden traffic spikes

#### Rate Limit Headers
```http
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 999
X-RateLimit-Reset: 1640995200
X-RateLimit-Retry-After: 3600
```

#### Rate Limit Storage
- Redis-based distributed rate limiting
- Sliding window algorithm
- Automatic cleanup of expired entries

### 5. Security Monitoring and Logging

#### Comprehensive Logging
- All API key usage tracked with:
  - Timestamp
  - Endpoint accessed
  - HTTP method
  - Response time
  - Status code
  - IP address
  - User agent
  - Request size

#### Security Event Detection
- Automatic detection of:
  - Rate limit violations
  - Invalid key usage attempts
  - Suspicious access patterns
  - Multiple failed authentication attempts
  - Unusual geographic access patterns

#### Real-time Alerting
- Configurable thresholds for security events
- Email/webhook notifications for critical events
- Automatic key suspension for severe violations

### 6. Key Lifecycle Management

#### Key States
- `active` - Normal operation
- `suspended` - Temporarily disabled
- `revoked` - Permanently disabled
- `expired` - Past expiration date

#### Automatic Management
- Configurable expiration dates
- Automatic cleanup of expired keys
- Usage-based key rotation recommendations
- Inactive key detection and alerts

### 7. Validation and Security Checks

#### Comprehensive Key Validation
The system performs multiple security checks on API keys:

##### Length Validation
- Minimum: 32 characters
- Recommended: 64+ characters
- Maximum: 128 characters

##### Entropy Analysis
- Shannon entropy calculation
- Minimum threshold: 4.0 bits per character
- Pattern detection for weak keys

##### Security Pattern Detection
- Sequential character detection
- Common word detection (password, secret, admin, etc.)
- Keyboard pattern detection (qwerty, asdf, etc.)
- Date pattern detection
- Character distribution analysis

##### Security Scoring
- 0-100 point scoring system
- Security levels: Low, Medium, High, Critical
- Automatic recommendations for improvement

### 8. Error Handling and Security

#### Secure Error Messages
- Generic error messages to prevent information leakage
- Detailed logging for debugging (server-side only)
- No exposure of internal system details
- Consistent error format across all endpoints

#### Error Response Format
```json
{
  "error": {
    "code": "INVALID_API_KEY",
    "message": "Invalid API key provided",
    "timestamp": "2024-01-15T10:30:00Z",
    "request_id": "req_123456789"
  }
}
```

### 9. Security Headers

#### Automatic Security Headers
- `X-Content-Type-Options: nosniff`
- `X-Frame-Options: DENY`
- `X-XSS-Protection: 1; mode=block`
- `Strict-Transport-Security: max-age=31536000; includeSubDomains`
- `Content-Security-Policy: default-src 'self'`
- `Referrer-Policy: strict-origin-when-cross-origin`

### 10. Database Security

#### Table Structure
```sql
-- API Keys table with security considerations
CREATE TABLE api_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    key_hash VARCHAR(64) NOT NULL UNIQUE, -- SHA-256 hash
    scopes TEXT[] NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    rate_limit INTEGER NOT NULL DEFAULT 1000,
    user_id UUID REFERENCES users(id),
    expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_used_at TIMESTAMP WITH TIME ZONE,
    created_ip INET,
    last_used_ip INET
);

-- Indexes for performance and security
CREATE INDEX idx_api_keys_hash ON api_keys(key_hash);
CREATE INDEX idx_api_keys_status ON api_keys(status);
CREATE INDEX idx_api_keys_expires ON api_keys(expires_at);
```

#### Security Triggers
- Automatic status updates for expired keys
- Security event logging triggers
- Rate limit violation detection
- Suspicious activity pattern detection

### 11. Compliance and Standards

#### Security Standards Compliance
- **OWASP API Security Top 10** - Full compliance
- **NIST Cybersecurity Framework** - Aligned implementation
- **ISO 27001** - Security management practices
- **PCI DSS** - Where applicable for payment data

#### Specific OWASP API Security Measures
1. **API1: Broken Object Level Authorization** - Scope-based access control
2. **API2: Broken User Authentication** - Strong API key authentication
3. **API3: Excessive Data Exposure** - Minimal data in responses
4. **API4: Lack of Resources & Rate Limiting** - Comprehensive rate limiting
5. **API5: Broken Function Level Authorization** - Endpoint-specific permissions
6. **API6: Mass Assignment** - Input validation and sanitization
7. **API7: Security Misconfiguration** - Secure defaults and configuration
8. **API8: Injection** - Parameterized queries and input validation
9. **API9: Improper Assets Management** - API versioning and documentation
10. **API10: Insufficient Logging & Monitoring** - Comprehensive logging system

### 12. Security Testing

#### Automated Security Testing
- Unit tests for all security functions
- Integration tests for authentication flows
- Performance benchmarks for security operations
- Penetration testing scenarios

#### Test Coverage Areas
- API key generation and validation
- Authentication and authorization
- Rate limiting functionality
- Error handling security
- Timing attack resistance
- Input validation and sanitization

### 13. Deployment Security

#### Environment Configuration
```bash
# Required environment variables
API_KEY_ENCRYPTION_KEY=your-32-byte-encryption-key
DATABASE_URL=postgresql://user:pass@host:port/db?sslmode=require
REDIS_URL=redis://user:pass@host:port/db
JWT_SECRET=your-jwt-secret-key
CORS_ORIGINS=https://yourdomain.com
RATE_LIMIT_REDIS_URL=redis://host:port/1
```

#### Production Checklist
- [ ] HTTPS enforced for all endpoints
- [ ] Database connections use SSL/TLS
- [ ] Redis connections secured
- [ ] Environment variables properly set
- [ ] Logging configured and monitored
- [ ] Rate limiting properly configured
- [ ] Security headers enabled
- [ ] CORS properly configured
- [ ] API documentation secured
- [ ] Monitoring and alerting active

### 14. Incident Response

#### Security Incident Procedures
1. **Detection** - Automated monitoring and alerting
2. **Assessment** - Severity classification and impact analysis
3. **Containment** - Immediate threat mitigation
4. **Investigation** - Root cause analysis
5. **Recovery** - System restoration and validation
6. **Lessons Learned** - Process improvement

#### Automated Response Actions
- Automatic key suspension for severe violations
- Rate limit adjustment during attacks
- IP blocking for malicious actors
- Alert escalation for critical events

### 15. Security Maintenance

#### Regular Security Tasks
- **Daily**: Monitor security logs and alerts
- **Weekly**: Review API key usage patterns
- **Monthly**: Security audit and vulnerability assessment
- **Quarterly**: Penetration testing and security review
- **Annually**: Full security architecture review

#### Key Rotation Policy
- Recommended rotation: Every 90 days
- Forced rotation: After security incidents
- Automatic notifications: 30 days before expiration
- Grace period: 7 days for key transition

### 16. API Endpoints Security Summary

#### Public Endpoints (No Authentication)
- `GET /health` - Health check
- `GET /api/v1/nutrition/meals` - Basic meal data
- `GET /api/v1/nutrition/recipes` - Basic recipe data

#### API Key Protected Endpoints
- `GET /api/v1/nutrition/*` - Enhanced nutrition data
- `POST /api/v1/meals` - Create meals (requires write scope)
- `PUT /api/v1/meals/:id` - Update meals (requires write scope)
- `DELETE /api/v1/meals/:id` - Delete meals (requires write scope)

#### Admin Only Endpoints
- `POST /admin/api-keys` - Create API keys
- `GET /admin/api-keys/:id` - Get API key details
- `DELETE /admin/api-keys/:id` - Revoke API keys
- `GET /admin/api-keys/:id/stats` - API key statistics

### 17. Security Metrics and KPIs

#### Key Security Metrics
- API key validation success rate
- Rate limiting effectiveness
- Security event detection rate
- Response time for security operations
- Failed authentication attempts
- Suspicious activity detection

#### Performance Benchmarks
- API key validation: < 1ms
- Rate limit check: < 0.5ms
- Security logging: < 2ms
- Hash comparison: < 0.1ms

## Conclusion

This comprehensive security implementation provides enterprise-grade protection for the Nutrition Platform API. The multi-layered security approach ensures protection against common vulnerabilities while maintaining high performance and usability.

For questions or security concerns, please contact the security team or refer to the incident response procedures outlined above.

---

**Last Updated**: January 2024  
**Version**: 1.0  
**Classification**: Internal Use