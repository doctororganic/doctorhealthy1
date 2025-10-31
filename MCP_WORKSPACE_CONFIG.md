# MCP Workspace Configuration Guide

This document provides all required variables, data, and configuration for integrating MCP (Model Context Protocol) servers across the entire workspace.

## 🔧 Global Environment Variables

### Required for All MCP Servers

```bash
# Redis Configuration (Shared Memory)
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=secure_redis_password_2025

# Workspace Paths
WORKSPACE_ROOT=/Users/khaledahmedmohamed/Desktop/trae new healthy1
NUTRITION_PLATFORM_PATH=${WORKSPACE_ROOT}/nutrition-platform
MCP_SERVERS_PATH=${NUTRITION_PLATFORM_PATH}/mcp-servers

# Node.js Configuration
NODE_ENV=development
NODE_PATH=${MCP_SERVERS_PATH}/node_modules

# Python Configuration (if needed)
PYTHONPATH=${WORKSPACE_ROOT}
PIP_CACHE_DIR=${WORKSPACE_ROOT}/.cache/pip

# Browser Automation
BROWSER_HEADLESS=false
BROWSER_WIDTH=1920
BROWSER_HEIGHT=1080
BROWSER_TIMEOUT=30000

# OCR Configuration
OCR_LANGUAGE=eng
OCR_CACHE_DIR=${WORKSPACE_ROOT}/.cache/ocr

# Document Processing
DOCX_CACHE_DIR=${WORKSPACE_ROOT}/.cache/docx
PDF_CACHE_DIR=${WORKSPACE_ROOT}/.cache/pdf

# Desktop Automation
DESKTOP_DISPLAY=:0
DESKTOP_SCALE_FACTOR=1.0

# Logging
LOG_LEVEL=info
LOG_DIR=${WORKSPACE_ROOT}/logs
LOG_MAX_SIZE=10m
LOG_MAX_FILES=5
```

## 🏗️ MCP Server Configurations

### 1. Roo Automation MCP Server (Basic)

**File:** `nutrition-platform/mcp-servers/index.js`

**Purpose:** Basic system automation, file operations, shell commands

**Configuration:**
```json
{
  "mcpServers": {
    "roo-automation": {
      "command": "node",
      "args": ["/Users/khaledahmedmohamed/Desktop/trae new healthy1/nutrition-platform/mcp-servers/index.js"],
      "env": {
        "REDIS_HOST": "localhost",
        "REDIS_PORT": "6379",
        "REDIS_PASSWORD": "secure_redis_password_2025",
        "NODE_ENV": "development",
        "LOG_LEVEL": "info"
      },
      "timeout": 30000
    }
  }
}
```

**Required Data:**
- Redis connection string
- File system permissions
- Node.js runtime
- Shell environment

### 2. Advanced Automation MCP Server

**File:** `nutrition-platform/mcp-servers/advanced-automation-server.js`

**Purpose:** Comprehensive automation suite

**Configuration:**
```json
{
  "mcpServers": {
    "advanced-automation": {
      "command": "node",
      "args": ["/Users/khaledahmedmohamed/Desktop/trae new healthy1/nutrition-platform/mcp-servers/advanced-automation-server.js"],
      "env": {
        "REDIS_HOST": "localhost",
        "REDIS_PORT": "6379",
        "REDIS_PASSWORD": "secure_redis_password_2025",
        "BROWSER_HEADLESS": "false",
        "BROWSER_WIDTH": "1920",
        "BROWSER_HEIGHT": "1080",
        "BROWSER_TIMEOUT": "30000",
        "OCR_LANGUAGE": "eng",
        "DESKTOP_DISPLAY": ":0",
        "LOG_LEVEL": "info",
        "NODE_ENV": "development",
        "OCR_CACHE_DIR": "/Users/khaledahmedmohamed/Desktop/trae new healthy1/.cache/ocr",
        "DOCX_CACHE_DIR": "/Users/khaledahmedmohamed/Desktop/trae new healthy1/.cache/docx",
        "PDF_CACHE_DIR": "/Users/khaledahmedmohamed/Desktop/trae new healthy1/.cache/pdf"
      },
      "timeout": 60000,
      "capabilities": {
        "browser": true,
        "desktop": true,
        "documents": true,
        "ocr": true,
        "workflows": true,
        "ai_coordination": true
      }
    }
  }
}
```

**Required Data:**
- Redis connection details
- Browser automation settings
- Desktop display configuration
- OCR language packs
- Document processing libraries
- AI agent registry

### 3. Browser Automation MCP Server

**File:** `nutrition-platform/mcp-servers/browser-automation.js`

**Purpose:** Specialized browser automation with Puppeteer/Playwright

**Configuration:**
```json
{
  "mcpServers": {
    "browser-automation": {
      "command": "node",
      "args": ["/Users/khaledahmedmohamed/Desktop/trae new healthy1/nutrition-platform/mcp-servers/browser-automation.js"],
      "env": {
        "BROWSER_ENGINE": "puppeteer",
        "BROWSER_HEADLESS": "false",
        "BROWSER_WIDTH": "1920",
        "BROWSER_HEIGHT": "1080",
        "BROWSER_TIMEOUT": "30000",
        "BROWSER_USER_AGENT": "MCP Browser Automation/1.0",
        "BROWSER_PROXY": "",
        "BROWSER_DOWNLOAD_DIR": "/Users/khaledahmedmohamed/Desktop/trae new healthy1/downloads",
        "SCREENSHOT_DIR": "/Users/khaledahmedmohamed/Desktop/trae new healthy1/screenshots",
        "VIDEO_DIR": "/Users/khaledahmedmohamed/Desktop/trae new healthy1/videos",
        "REDIS_HOST": "localhost",
        "REDIS_PORT": "6379",
        "REDIS_PASSWORD": "secure_redis_password_2025"
      },
      "timeout": 120000
    }
  }
}
```

**Required Data:**
- Browser binary paths
- User agent strings
- Proxy configurations
- Download directories
- Screenshot/video storage paths

### 4. Desktop Automation MCP Server

**File:** `nutrition-platform/mcp-servers/desktop-automation.js`

**Purpose:** System-wide desktop automation

**Configuration:**
```json
{
  "mcpServers": {
    "desktop-automation": {
      "command": "node",
      "args": ["/Users/khaledahmedmohamed/Desktop/trae new healthy1/nutrition-platform/mcp-servers/desktop-automation.js"],
      "env": {
        "DISPLAY": ":0",
        "DESKTOP_SCALE_FACTOR": "1.0",
        "MOUSE_SPEED": "100",
        "KEYBOARD_LAYOUT": "us",
        "SCREENSHOT_DIR": "/Users/khaledahmedmohamed/Desktop/trae new healthy1/desktop-screenshots",
        "AUTOMATION_LOG_DIR": "/Users/khaledahmedmohamed/Desktop/trae new healthy1/logs/desktop",
        "SAFE_MODE": "true",
        "REDIS_HOST": "localhost",
        "REDIS_PORT": "6379",
        "REDIS_PASSWORD": "secure_redis_password_2025"
      },
      "timeout": 60000,
      "security": {
        "require_confirmation": true,
        "max_actions_per_minute": 60,
        "allowed_applications": ["*"],
        "blocked_applications": []
      }
    }
  }
}
```

**Required Data:**
- Display server information
- Screen resolution and scaling
- Keyboard layout settings
- Mouse sensitivity settings
- Security policies

### 5. Document Processing MCP Server

**File:** `nutrition-platform/mcp-servers/document-processing.js`

**Purpose:** Document analysis and processing

**Configuration:**
```json
{
  "mcpServers": {
    "document-processing": {
      "command": "node",
      "args": ["/Users/khaledahmedmohamed/Desktop/trae new healthy1/nutrition-platform/mcp-servers/document-processing.js"],
      "env": {
        "OCR_LANGUAGE": "eng",
        "OCR_CACHE_DIR": "/Users/khaledahmedmohamed/Desktop/trae new healthy1/.cache/ocr",
        "PDF_CACHE_DIR": "/Users/khaledahmedmohamed/Desktop/trae new healthy1/.cache/pdf",
        "DOCX_CACHE_DIR": "/Users/khaledahmedmohamed/Desktop/trae new healthy1/.cache/docx",
        "TEMP_DIR": "/tmp/mcp-documents",
        "MAX_FILE_SIZE": "100MB",
        "SUPPORTED_FORMATS": "pdf,docx,txt,rtf,odt",
        "OCR_ENGINES": "tesseract,easyocr",
        "PDF_EXTRACT_IMAGES": "true",
        "DOCX_EXTRACT_IMAGES": "true",
        "REDIS_HOST": "localhost",
        "REDIS_PORT": "6379",
        "REDIS_PASSWORD": "secure_redis_password_2025"
      },
      "timeout": 300000,
      "processing": {
        "max_concurrent_jobs": 3,
        "queue_size": 100,
        "cleanup_temp_files": true
      }
    }
  }
}
```

**Required Data:**
- OCR model files and language packs
- Document template libraries
- Font collections
- Image processing libraries
- Text analysis models

### 6. Workflow Orchestration MCP Server

**File:** `nutrition-platform/mcp-servers/workflow-orchestration.js`

**Purpose:** Complex multi-step workflow management

**Configuration:**
```json
{
  "mcpServers": {
    "workflow-orchestration": {
      "command": "node",
      "args": ["/Users/khaledahmedmohamed/Desktop/trae new healthy1/nutrition-platform/mcp-servers/workflow-orchestration.js"],
      "env": {
        "MAX_WORKFLOW_STEPS": "50",
        "WORKFLOW_TIMEOUT": "3600000",
        "PARALLEL_EXECUTION": "true",
        "MAX_PARALLEL_JOBS": "5",
        "WORKFLOW_STORAGE_DIR": "/Users/khaledahmedmohamed/Desktop/trae new healthy1/workflows",
        "WORKFLOW_LOG_DIR": "/Users/khaledahmedmohamed/Desktop/trae new healthy1/logs/workflows",
        "RETRY_FAILED_STEPS": "true",
        "MAX_RETRIES": "3",
        "REDIS_HOST": "localhost",
        "REDIS_PORT": "6379",
        "REDIS_PASSWORD": "secure_redis_password_2025",
        "NOTIFICATION_WEBHOOK": "",
        "SLACK_WEBHOOK": ""
      },
      "timeout": 3600000,
      "workflows": {
        "enabled_types": ["sequential", "parallel", "conditional", "loop"],
        "max_nested_workflows": 3,
        "state_persistence": true,
        "error_recovery": true
      }
    }
  }
}
```

**Required Data:**
- Workflow template library
- Step execution engines
- State persistence schemas
- Error handling strategies
- Notification configurations

### 7. AI Agent Coordination MCP Server

**File:** `nutrition-platform/mcp-servers/ai-coordination.js`

**Purpose:** Multi-agent communication and coordination

**Configuration:**
```json
{
  "mcpServers": {
    "ai-coordination": {
      "command": "node",
      "args": ["/Users/khaledahmedmohamed/Desktop/trae new healthy1/nutrition-platform/mcp-servers/ai-coordination.js"],
      "env": {
        "MAX_AGENTS": "10",
        "MESSAGE_TTL": "3600",
        "COORDINATION_TIMEOUT": "30000",
        "AGENT_REGISTRY_FILE": "/Users/khaledahmedmohamed/Desktop/trae new healthy1/agent-registry.json",
        "COORDINATION_LOG_DIR": "/Users/khaledahmedmohamed/Desktop/trae new healthy1/logs/coordination",
        "CONFLICT_RESOLUTION": "priority_based",
        "LOAD_BALANCING": "round_robin",
        "HEARTBEAT_INTERVAL": "30000",
        "REDIS_HOST": "localhost",
        "REDIS_PORT": "6379",
        "REDIS_PASSWORD": "secure_redis_password_2025",
        "API_ENDPOINT": "http://localhost:8080",
        "WEBSOCKET_PORT": "8081"
      },
      "timeout": 60000,
      "coordination": {
        "enable_websockets": true,
        "enable_rest_api": true,
        "enable_message_queue": true,
        "max_message_size": "1MB",
        "encryption_enabled": false
      }
    }
  }
}
```

**Required Data:**
- Agent capability registry
- Communication protocols
- Conflict resolution algorithms
- Load balancing strategies
- Message encryption keys

## 📁 Directory Structure Requirements

```
/Users/khaledahmedmohamed/Desktop/trae new healthy1/
├── nutrition-platform/
│   └── mcp-servers/
│       ├── index.js                          # Basic automation
│       ├── advanced-automation-server.js    # Comprehensive suite
│       ├── browser-automation.js            # Browser tools
│       ├── desktop-automation.js            # Desktop tools
│       ├── document-processing.js           # Document tools
│       ├── workflow-orchestration.js        # Workflow tools
│       └── ai-coordination.js               # Agent coordination
├── .cache/
│   ├── ocr/
│   ├── docx/
│   ├── pdf/
│   └── pip/
├── logs/
│   ├── desktop/
│   ├── workflows/
│   └── coordination/
├── screenshots/
├── videos/
├── downloads/
├── workflows/
└── agent-registry.json
```

## 🔐 Security Configuration

### Access Control
```json
{
  "security": {
    "allowed_commands": ["*"],
    "blocked_commands": ["rm -rf /", "sudo", "passwd"],
    "max_execution_time": 300000,
    "require_user_confirmation": true,
    "log_all_actions": true,
    "audit_trail_enabled": true
  }
}
```

### Network Security
```json
{
  "network": {
    "allowed_hosts": ["localhost", "127.0.0.1"],
    "blocked_ports": [22, 23, 3389],
    "max_connections_per_host": 10,
    "connection_timeout": 10000,
    "ssl_verification": true
  }
}
```

## 📊 Monitoring and Metrics

### Required Metrics Collection
```json
{
  "monitoring": {
    "enable_prometheus": true,
    "metrics_port": 9090,
    "collect_system_metrics": true,
    "collect_performance_metrics": true,
    "log_level": "info",
    "alert_webhooks": [],
    "health_check_interval": 30000
  }
}
```

## 🚀 Deployment Configuration

### Docker Compose Setup
```yaml
version: '3.8'
services:
  mcp-coordinator:
    image: node:18-alpine
    volumes:
      - ./nutrition-platform:/app
      - ./.cache:/cache
      - ./logs:/logs
    environment:
      - REDIS_HOST=redis
      - REDIS_PORT=6379
      - REDIS_PASSWORD=secure_redis_password_2025
    depends_on:
      - redis
    command: node mcp-servers/advanced-automation-server.js

  redis:
    image: redis:7-alpine
    command: redis-server --requirepass secure_redis_password_2025
    volumes:
      - redis_data:/data

volumes:
  redis_data:
```

## 🧪 Testing Configuration

### Test Environment Variables
```bash
MCP_TEST_MODE=true
MCP_MOCK_EXTERNAL_APIS=true
MCP_DISABLE_SECURITY_CHECKS=false
MCP_LOG_TEST_RESULTS=true
```

### Integration Test Data
```json
{
  "test_documents": {
    "pdf_path": "/test/sample.pdf",
    "docx_path": "/test/sample.docx",
    "image_path": "/test/sample.png"
  },
  "test_urls": {
    "sample_website": "https://example.com",
    "api_endpoint": "https://jsonplaceholder.typicode.com"
  },
  "test_workflows": {
    "basic_workflow": "workflow_basic_test.json",
    "complex_workflow": "workflow_complex_test.json"
  }
}
```

## 🔧 Troubleshooting Data

### Common Issues and Solutions
```json
{
  "troubleshooting": {
    "redis_connection_failed": {
      "symptoms": ["ECONNREFUSED", "authentication failed"],
      "solutions": ["Check Redis service", "Verify password", "Check network connectivity"]
    },
    "browser_launch_failed": {
      "symptoms": ["WebDriver error", "headless mode issues"],
      "solutions": ["Install browser binaries", "Check display settings", "Update browser drivers"]
    },
    "ocr_processing_failed": {
      "symptoms": ["tesseract not found", "language pack missing"],
      "solutions": ["Install tesseract", "Download language packs", "Check PATH"]
    }
  }
}
```

This configuration guide provides all the variables, data structures, and settings needed to fully integrate MCP automation capabilities across your entire workspace.