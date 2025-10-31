#!/usr/bin/env node

/**
 * AI Collaboration Automation Starter
 * Human-in-the-loop orchestration for AI agent coordination
 */

const sharedMemory = require('./ai-collaboration/shared-memory-service');
const orchestrationController = require('./ai-collaboration/orchestration-controller');
const roleAssignmentManager = require('./ai-collaboration/role-assignments');
const logger = require('./production-nodejs/utils/logger');

class AIAutomationStarter {
  constructor() {
    this.isInitialized = false;
    this.orchestrationActive = false;
  }

  /**
   * Initialize the AI automation system
   */
  async initialize() {
    try {
      logger.info('🚀 Initializing AI Collaboration Automation System');

      // Step 1: Initialize shared memory
      logger.info('Step 1: Initializing shared memory system');
      await sharedMemory.setAction('system', 'automation_init', {
        type: 'system_initialization',
        status: 'in_progress',
        timestamp: Date.now()
      });

      // Step 2: Check system requirements
      logger.info('Step 2: Checking system requirements');
      await this.checkSystemRequirements();

      // Step 3: Initialize monitoring
      logger.info('Step 3: Starting monitoring systems');
      await this.initializeMonitoring();

      // Step 4: Set up role assignments
      logger.info('Step 4: Setting up AI agent role assignments');
      await this.setupRoleAssignments();

      this.isInitialized = true;

      await sharedMemory.completeAction('system', 'automation_init', {
        success: true,
        initialized: true,
        timestamp: Date.now()
      });

      logger.info('✅ AI automation system initialized successfully');
      return true;
    } catch (error) {
      logger.error('Failed to initialize AI automation system', { error: error.message });
      throw error;
    }
  }

  /**
   * Check system requirements for automation
   */
  async checkSystemRequirements() {
    try {
      // Check if Redis is available
      try {
        await sharedMemory.redis.ping();
        logger.info('✅ Redis connection verified');
      } catch (error) {
        logger.warn('⚠️ Redis not available, using in-memory fallback');
      }

      // Check if monitoring services are running
      const monitoringStatus = await this.checkMonitoringServices();
      if (monitoringStatus.allHealthy) {
        logger.info('✅ All monitoring services healthy');
      } else {
        logger.warn('⚠️ Some monitoring services not available');
      }

      // Check if Node.js dependencies are installed
      const fs = require('fs');
      if (fs.existsSync('./production-nodejs/node_modules')) {
        logger.info('✅ Node.js dependencies installed');
      } else {
        logger.warn('⚠️ Node.js dependencies not installed');
      }

      return true;
    } catch (error) {
      logger.error('System requirements check failed', { error: error.message });
      throw error;
    }
  }

  /**
   * Check monitoring services health
   */
  async checkMonitoringServices() {
    const services = [
      { name: 'Grafana', url: 'http://localhost:3001/api/health', port: 3001 },
      { name: 'Prometheus', url: 'http://localhost:9090/-/healthy', port: 9090 },
      { name: 'Loki', url: 'http://localhost:3100/ready', port: 3100 }
    ];

    const results = {};

    for (const service of services) {
      try {
        const response = await fetch(service.url);
        results[service.name] = response.ok;
      } catch (error) {
        results[service.name] = false;
      }
    }

    const allHealthy = Object.values(results).every(status => status);

    return {
      services: results,
      allHealthy,
      healthyCount: Object.values(results).filter(status => status).length,
      totalCount: services.length
    };
  }

  /**
   * Initialize monitoring systems
   */
  async initializeMonitoring() {
    try {
      logger.info('Starting monitoring system initialization');

      // Start Loki, Grafana, Prometheus if not already running
      const { exec } = require('child_process');

      return new Promise((resolve, reject) => {
        exec('docker-compose -f docker-compose.monitoring.yml up -d', (error, stdout, stderr) => {
          if (error) {
            logger.warn('Monitoring services may already be running', { error: error.message });
            resolve(false);
          } else {
            logger.info('Monitoring services started', { stdout });
            resolve(true);
          }
        });
      });
    } catch (error) {
      logger.error('Failed to initialize monitoring', { error: error.message });
      throw error;
    }
  }

  /**
   * Set up AI agent role assignments
   */
  async setupRoleAssignments() {
    try {
      logger.info('Setting up AI agent role assignments');

      // Initialize role assignments
      await roleAssignmentManager.initializeRoleAssignments();

      // Set up initial workflow state
      await sharedMemory.setAction('workflow', 'initial_state', {
        type: 'workflow_initialization',
        status: 'ready',
        agents: {
          kilo: { role: 'frontend_development', status: 'ready' },
          roo: { role: 'backend_integration', status: 'ready' },
          codesupernova: { role: 'monitoring_systems', status: 'active' }
        },
        timestamp: Date.now()
      });

      logger.info('Role assignments completed');
      return true;
    } catch (error) {
      logger.error('Failed to setup role assignments', { error: error.message });
      throw error;
    }
  }

  /**
   * Start the automation workflow
   */
  async startAutomation() {
    try {
      if (!this.isInitialized) {
        await this.initialize();
      }

      logger.info('🎯 Starting AI automation workflow');

      // Start orchestration with human oversight
      const result = await orchestrationController.startCollaborationWorkflow();

      this.orchestrationActive = true;

      logger.info('✅ AI automation workflow started successfully');
      return result;
    } catch (error) {
      logger.error('Failed to start automation', { error: error.message });
      throw error;
    }
  }

  /**
   * Monitor automation progress
   */
  async monitorProgress() {
    try {
      const orchestrationStatus = await orchestrationController.getOrchestrationStatus();
      const workflowStatus = await roleAssignmentManager.getAssignmentStatus();

      const progressReport = {
        timestamp: Date.now(),
        orchestration: orchestrationStatus,
        workflow: workflowStatus,
        humanInterventions: orchestrationController.humanApprovalQueue.length,
        activeAgents: workflowStatus ? Object.keys(workflowStatus.assignments) : []
      };

      logger.info('Progress monitoring report', progressReport);
      return progressReport;
    } catch (error) {
      logger.error('Progress monitoring failed', { error: error.message });
      return null;
    }
  }

  /**
   * Interactive automation starter
   */
  async startInteractive() {
    try {
      console.log('🤖 AI Collaboration Automation System');
      console.log('=====================================');
      console.log('');

      // Check if already initialized
      if (this.isInitialized) {
        console.log('✅ System already initialized');
      } else {
        console.log('🔧 Initializing system...');
        await this.initialize();
        console.log('✅ System initialized');
      }

      console.log('');
      console.log('🎯 Available Actions:');
      console.log('1. 🚀 Start Full Automation Workflow');
      console.log('2. 🎨 Start Frontend Development (Kilo Code)');
      console.log('3. 🔧 Start Backend Integration (Roo Code)');
      console.log('4. 🧪 Start Testing & Validation');
      console.log('5. 📊 Monitor Progress');
      console.log('6. 🚨 Emergency Stop');
      console.log('');

      // For now, start the full workflow
      console.log('🚀 Starting full automation workflow...');
      await this.startAutomation();

      // Monitor progress
      setInterval(async () => {
        await this.monitorProgress();
      }, 30000); // Every 30 seconds

      return true;
    } catch (error) {
      console.error('❌ Interactive automation failed:', error.message);
      throw error;
    }
  }

  /**
   * Get system status
   */
  async getSystemStatus() {
    try {
      const orchestrationStatus = await orchestrationController.getOrchestrationStatus();
      const workflowStatus = await roleAssignmentManager.getAssignmentStatus();
      const monitoringStatus = await this.checkMonitoringServices();

      return {
        initialized: this.isInitialized,
        orchestrationActive: this.orchestrationActive,
        orchestration: orchestrationStatus,
        workflow: workflowStatus,
        monitoring: monitoringStatus,
        timestamp: Date.now()
      };
    } catch (error) {
      logger.error('Failed to get system status', { error: error.message });
      return null;
    }
  }
}

// CLI interface
async function main() {
  const automation = new AIAutomationStarter();

  try {
    if (process.argv[2] === '--status') {
      const status = await automation.getSystemStatus();
      console.log('📊 System Status:', JSON.stringify(status, null, 2));
      return;
    }

    if (process.argv[2] === '--monitor') {
      console.log('👀 Starting progress monitoring...');
      setInterval(async () => {
        const progress = await automation.monitorProgress();
        if (progress) {
          console.log('📈 Progress Update:', new Date(progress.timestamp).toISOString());
        }
      }, 10000); // Every 10 seconds
      return;
    }

    // Default: Start interactive mode
    await automation.startInteractive();

  } catch (error) {
    console.error('❌ Automation failed:', error.message);
    process.exit(1);
  }
}

// Export for use in other modules
module.exports = AIAutomationStarter;

// Run if called directly
if (require.main === module) {
  main().catch(console.error);
}