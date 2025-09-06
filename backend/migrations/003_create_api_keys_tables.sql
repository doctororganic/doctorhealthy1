-- Migration: Create API Keys Tables
-- Version: 003
-- Description: Creates tables for API key management, usage tracking, and security

-- Create API keys table
CREATE TABLE IF NOT EXISTS api_keys (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    key_hash VARCHAR(64) NOT NULL UNIQUE, -- SHA-256 hash of the API key
    prefix VARCHAR(10) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    scopes JSON NOT NULL, -- Array of scopes as JSON
    rate_limit INTEGER NOT NULL DEFAULT 1000, -- Requests per minute
    expires_at TIMESTAMP NULL,
    last_used_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    metadata JSON NULL, -- Additional metadata as JSON
    
    -- Indexes
    INDEX idx_api_keys_user_id (user_id),
    INDEX idx_api_keys_key_hash (key_hash),
    INDEX idx_api_keys_status (status),
    INDEX idx_api_keys_expires_at (expires_at),
    INDEX idx_api_keys_created_at (created_at),
    
    -- Constraints
    CONSTRAINT chk_api_keys_status CHECK (status IN ('active', 'inactive', 'revoked', 'expired')),
    CONSTRAINT chk_api_keys_rate_limit CHECK (rate_limit > 0 AND rate_limit <= 10000),
    CONSTRAINT chk_api_keys_name_length CHECK (LENGTH(name) >= 3 AND LENGTH(name) <= 100)
);

-- Create API key usage tracking table
CREATE TABLE IF NOT EXISTS api_key_usage (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    api_key_id VARCHAR(255) NOT NULL,
    endpoint VARCHAR(500) NOT NULL,
    method VARCHAR(10) NOT NULL,
    status_code INTEGER NOT NULL,
    response_time BIGINT NOT NULL, -- Response time in milliseconds
    ip_address VARCHAR(45) NOT NULL, -- IPv4 or IPv6
    user_agent TEXT,
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Indexes
    INDEX idx_api_key_usage_api_key_id (api_key_id),
    INDEX idx_api_key_usage_timestamp (timestamp),
    INDEX idx_api_key_usage_endpoint (endpoint),
    INDEX idx_api_key_usage_status_code (status_code),
    INDEX idx_api_key_usage_composite (api_key_id, timestamp),
    
    -- Foreign key constraint
    FOREIGN KEY (api_key_id) REFERENCES api_keys(id) ON DELETE CASCADE
);

-- Create API key rate limiting table (for distributed rate limiting)
CREATE TABLE IF NOT EXISTS api_key_rate_limits (
    api_key_id VARCHAR(255) NOT NULL,
    window_start TIMESTAMP NOT NULL,
    request_count INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    PRIMARY KEY (api_key_id, window_start),
    
    -- Indexes
    INDEX idx_rate_limits_window_start (window_start),
    
    -- Foreign key constraint
    FOREIGN KEY (api_key_id) REFERENCES api_keys(id) ON DELETE CASCADE
);

-- Create API key security events table
CREATE TABLE IF NOT EXISTS api_key_security_events (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    api_key_id VARCHAR(255) NULL, -- Can be NULL for events without valid API key
    event_type VARCHAR(50) NOT NULL,
    severity VARCHAR(20) NOT NULL DEFAULT 'info',
    description TEXT NOT NULL,
    ip_address VARCHAR(45) NOT NULL,
    user_agent TEXT,
    endpoint VARCHAR(500),
    method VARCHAR(10),
    additional_data JSON NULL,
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Indexes
    INDEX idx_security_events_api_key_id (api_key_id),
    INDEX idx_security_events_event_type (event_type),
    INDEX idx_security_events_severity (severity),
    INDEX idx_security_events_timestamp (timestamp),
    INDEX idx_security_events_ip_address (ip_address),
    
    -- Constraints
    CONSTRAINT chk_security_events_severity CHECK (severity IN ('info', 'warning', 'error', 'critical')),
    
    -- Foreign key constraint (with NULL allowed)
    FOREIGN KEY (api_key_id) REFERENCES api_keys(id) ON DELETE SET NULL
);

-- Create view for API key statistics
CREATE OR REPLACE VIEW api_key_stats AS
SELECT 
    ak.id,
    ak.name,
    ak.user_id,
    ak.status,
    ak.rate_limit,
    ak.created_at,
    ak.last_used_at,
    ak.expires_at,
    COALESCE(usage_stats.total_requests, 0) as total_requests,
    COALESCE(usage_stats.requests_last_24h, 0) as requests_last_24h,
    COALESCE(usage_stats.requests_last_7d, 0) as requests_last_7d,
    COALESCE(usage_stats.requests_last_30d, 0) as requests_last_30d,
    COALESCE(usage_stats.avg_response_time, 0) as avg_response_time_ms,
    COALESCE(usage_stats.error_rate, 0) as error_rate_percent
FROM api_keys ak
LEFT JOIN (
    SELECT 
        api_key_id,
        COUNT(*) as total_requests,
        COUNT(CASE WHEN timestamp >= DATE_SUB(NOW(), INTERVAL 24 HOUR) THEN 1 END) as requests_last_24h,
        COUNT(CASE WHEN timestamp >= DATE_SUB(NOW(), INTERVAL 7 DAY) THEN 1 END) as requests_last_7d,
        COUNT(CASE WHEN timestamp >= DATE_SUB(NOW(), INTERVAL 30 DAY) THEN 1 END) as requests_last_30d,
        AVG(response_time) as avg_response_time,
        (COUNT(CASE WHEN status_code >= 400 THEN 1 END) * 100.0 / COUNT(*)) as error_rate
    FROM api_key_usage
    GROUP BY api_key_id
) usage_stats ON ak.id = usage_stats.api_key_id;

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_api_key_usage_date_range ON api_key_usage (api_key_id, timestamp, status_code);
CREATE INDEX IF NOT EXISTS idx_api_key_usage_daily_stats ON api_key_usage (api_key_id, DATE(timestamp));

-- Create triggers for automatic cleanup
DELIMITER //

-- Trigger to automatically set status to 'expired' for expired API keys
CREATE TRIGGER IF NOT EXISTS tr_api_keys_check_expiry
    BEFORE UPDATE ON api_keys
    FOR EACH ROW
BEGIN
    IF NEW.expires_at IS NOT NULL AND NEW.expires_at <= NOW() AND NEW.status = 'active' THEN
        SET NEW.status = 'expired';
        SET NEW.updated_at = NOW();
    END IF;
END//

-- Trigger to log security events for suspicious activity
CREATE TRIGGER IF NOT EXISTS tr_api_key_usage_security_check
    AFTER INSERT ON api_key_usage
    FOR EACH ROW
BEGIN
    DECLARE request_count INT DEFAULT 0;
    DECLARE rate_limit INT DEFAULT 0;
    
    -- Get the rate limit for this API key
    SELECT ak.rate_limit INTO rate_limit
    FROM api_keys ak
    WHERE ak.id = NEW.api_key_id;
    
    -- Count requests in the last minute
    SELECT COUNT(*) INTO request_count
    FROM api_key_usage
    WHERE api_key_id = NEW.api_key_id
    AND timestamp >= DATE_SUB(NOW(), INTERVAL 1 MINUTE);
    
    -- Log security event if rate limit is exceeded
    IF request_count > rate_limit THEN
        INSERT INTO api_key_security_events (
            api_key_id, event_type, severity, description, 
            ip_address, user_agent, endpoint, method
        ) VALUES (
            NEW.api_key_id, 'rate_limit_exceeded', 'warning',
            CONCAT('Rate limit exceeded: ', request_count, ' requests in 1 minute (limit: ', rate_limit, ')'),
            NEW.ip_address, NEW.user_agent, NEW.endpoint, NEW.method
        );
    END IF;
    
    -- Log security event for 4xx/5xx responses
    IF NEW.status_code >= 400 THEN
        INSERT INTO api_key_security_events (
            api_key_id, event_type, severity, description,
            ip_address, user_agent, endpoint, method,
            additional_data
        ) VALUES (
            NEW.api_key_id, 
            CASE 
                WHEN NEW.status_code = 401 THEN 'unauthorized_access'
                WHEN NEW.status_code = 403 THEN 'forbidden_access'
                WHEN NEW.status_code >= 500 THEN 'server_error'
                ELSE 'client_error'
            END,
            CASE 
                WHEN NEW.status_code IN (401, 403) THEN 'warning'
                WHEN NEW.status_code >= 500 THEN 'error'
                ELSE 'info'
            END,
            CONCAT('HTTP ', NEW.status_code, ' response for ', NEW.method, ' ', NEW.endpoint),
            NEW.ip_address, NEW.user_agent, NEW.endpoint, NEW.method,
            JSON_OBJECT('status_code', NEW.status_code, 'response_time', NEW.response_time)
        );
    END IF;
END//

DELIMITER ;

-- Create stored procedures for common operations
DELIMITER //

-- Procedure to clean up old usage data
CREATE PROCEDURE IF NOT EXISTS CleanupOldUsageData(IN days_to_keep INT)
BEGIN
    DECLARE rows_deleted INT DEFAULT 0;
    
    -- Delete usage data older than specified days
    DELETE FROM api_key_usage 
    WHERE timestamp < DATE_SUB(NOW(), INTERVAL days_to_keep DAY);
    
    SET rows_deleted = ROW_COUNT();
    
    -- Log the cleanup
    INSERT INTO api_key_security_events (
        event_type, severity, description, ip_address
    ) VALUES (
        'data_cleanup', 'info', 
        CONCAT('Cleaned up ', rows_deleted, ' old usage records (older than ', days_to_keep, ' days)'),
        '127.0.0.1'
    );
END//

-- Procedure to get API key usage summary
CREATE PROCEDURE IF NOT EXISTS GetAPIKeyUsageSummary(
    IN p_api_key_id VARCHAR(255),
    IN p_days INT
)
BEGIN
    SELECT 
        DATE(timestamp) as date,
        COUNT(*) as total_requests,
        COUNT(CASE WHEN status_code < 400 THEN 1 END) as successful_requests,
        COUNT(CASE WHEN status_code >= 400 THEN 1 END) as error_requests,
        AVG(response_time) as avg_response_time,
        MIN(response_time) as min_response_time,
        MAX(response_time) as max_response_time
    FROM api_key_usage
    WHERE api_key_id = p_api_key_id
    AND timestamp >= DATE_SUB(NOW(), INTERVAL p_days DAY)
    GROUP BY DATE(timestamp)
    ORDER BY date DESC;
END//

DELIMITER ;

-- Insert default scopes documentation
INSERT IGNORE INTO api_key_security_events (
    event_type, severity, description, ip_address
) VALUES (
    'system_initialization', 'info',
    'API key management system initialized with tables and procedures',
    '127.0.0.1'
);

-- Create event to automatically clean up old data
CREATE EVENT IF NOT EXISTS evt_cleanup_old_usage_data
ON SCHEDULE EVERY 1 DAY
STARTS CURRENT_TIMESTAMP
DO
  CALL CleanupOldUsageData(90); -- Keep 90 days of usage data

-- Create event to clean up old rate limit data
CREATE EVENT IF NOT EXISTS evt_cleanup_old_rate_limit_data
ON SCHEDULE EVERY 1 HOUR
STARTS CURRENT_TIMESTAMP
DO
  DELETE FROM api_key_rate_limits 
  WHERE window_start < DATE_SUB(NOW(), INTERVAL 2 HOUR);

-- Create event to clean up old security events
CREATE EVENT IF NOT EXISTS evt_cleanup_old_security_events
ON SCHEDULE EVERY 1 DAY
STARTS CURRENT_TIMESTAMP
DO
  DELETE FROM api_key_security_events 
  WHERE timestamp < DATE_SUB(NOW(), INTERVAL 180 DAY); -- Keep 6 months of security events

-- Enable event scheduler if not already enabled
SET GLOBAL event_scheduler = ON;

COMMIT;