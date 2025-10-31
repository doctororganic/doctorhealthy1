/**
 * AI Collaboration - Redis-Based Shared Memory Service
 * Enables seamless coordination between AI agents (Kilo Code, Roo Code, CodeSupernova)
 */

const Redis = require('ioredis');
const logger = require('../production-nodejs/utils/logger');

class SharedMemoryService {
  constructor() {
    this.redis = new Redis({
      host: process.env.REDIS_HOST || 'localhost',
      port: process.env.REDIS_PORT || 6379,
      password: process.env.REDIS_PASSWORD || undefined,
      retryDelayOnFailover: 100,
      enableReadyCheck: true,
      maxRetriesPerRequest: 3,
      lazyConnect: true
    });

    this.namespace = 'ai_collaboration';
    this.agentRoles = {
      'kilo': 'frontend_development',
      'roo': 'backend_integration',
      'codesupernova': 'monitoring_systems'
    };

    this.setupRedisEventHandlers();
  }

  /**
   * Set up Redis event handlers for connection management
   */
  setupRedisEventHandlers() {
    this.redis.on('connect', () => {
      logger.info('Shared memory Redis connected', { service: 'ai_collaboration' });
    });

    this.redis.on('error', (error) => {
      logger.error('Shared memory Redis error', { error: error.message });
    });

    this.redis.on('close', () => {
      logger.warn('Shared memory Redis connection closed');
    });
  }

  /**
   * Store action/state for an AI agent
   */
  async setAction(agentId, actionId, data) {
    try {
      await this.redis.connect();

      const key = `${this.namespace}:${agentId}:${actionId}`;
      const actionData = {
        ...data,
        timestamp: Date.now(),
        agentId,
        status: 'active',
        agentRole: this.agentRoles[agentId] || 'unknown'
      };

      await this.redis.setex(key, 3600, JSON.stringify(actionData));

      logger.info('Action stored in shared memory', {
        agentId,
        actionId,
        status: 'active'
      });

      return actionData;
    } catch (error) {
      logger.error('Failed to set action in shared memory', {
        error: error.message,
        agentId,
        actionId
      });
      throw error;
    }
  }

  /**
   * Get action/state for an AI agent
   */
  async getAction(agentId, actionId) {
    try {
      await this.redis.connect();

      const key = `${this.namespace}:${agentId}:${actionId}`;
      const data = await this.redis.get(key);

      if (data) {
        const parsed = JSON.parse(data);
        logger.debug('Action retrieved from shared memory', {
          agentId,
          actionId,
          status: parsed.status
        });
        return parsed;
      }

      return null;
    } catch (error) {
      logger.error('Failed to get action from shared memory', {
        error: error.message,
        agentId,
        actionId
      });
      return null;
    }
  }

  /**
   * Mark action as completed
   */
  async completeAction(agentId, actionId, result) {
    try {
      await this.redis.connect();

      const key = `${this.namespace}:${agentId}:${actionId}`;
      const current = await this.getAction(agentId, actionId);

      if (current) {
        const completedData = {
          ...current,
          status: 'completed',
          result,
          completedAt: Date.now()
        };

        await this.redis.setex(key, 3600, JSON.stringify(completedData));

        logger.info('Action completed in shared memory', {
          agentId,
          actionId,
          status: 'completed'
        });

        return completedData;
      }

      return null;
    } catch (error) {
      logger.error('Failed to complete action in shared memory', {
        error: error.message,
        agentId,
        actionId
      });
      throw error;
    }
  }

  /**
   * Wait for dependency completion
   */
  async waitForDependency(agentId, dependencyActionId, timeout = 30000) {
    const startTime = Date.now();

    logger.info('Waiting for dependency', {
      agentId,
      dependencyActionId,
      timeout
    });

    while (Date.now() - startTime < timeout) {
      const dependency = await this.getAction(agentId, dependencyActionId);

      if (dependency && dependency.status === 'completed') {
        logger.info('Dependency completed', {
          agentId,
          dependencyActionId,
          waitTime: Date.now() - startTime
        });
        return dependency.result;
      }

      // Wait 1 second before checking again
      await new Promise(resolve => setTimeout(resolve, 1000));
    }

    const error = new Error(`Dependency ${dependencyActionId} not completed within ${timeout}ms timeout`);
    logger.error('Dependency timeout', {
      agentId,
      dependencyActionId,
      timeout,
      error: error.message
    });
    throw error;
  }

  /**
   * Get all active actions across all agents
   */
  async getActiveActions() {
    try {
      await this.redis.connect();

      const pattern = `${this.namespace}:*:*:*`;
      const keys = await this.redis.keys(pattern);
      const actions = [];

      for (const key of keys) {
        const data = await this.redis.get(key);
        if (data) {
          const parsed = JSON.parse(data);
          if (parsed.status === 'active') {
            actions.push(parsed);
          }
        }
      }

      return actions;
    } catch (error) {
      logger.error('Failed to get active actions', { error: error.message });
      return [];
    }
  }

  /**
   * Get actions by specific agent
   */
  async getAgentActions(agentId) {
    try {
      await this.redis.connect();

      const pattern = `${this.namespace}:${agentId}:*`;
      const keys = await this.redis.keys(pattern);
      const actions = [];

      for (const key of keys) {
        const data = await this.redis.get(key);
        if (data) {
          actions.push(JSON.parse(data));
        }
      }

      return actions;
    } catch (error) {
      logger.error('Failed to get agent actions', {
        error: error.message,
        agentId
      });
      return [];
    }
  }

  /**
   * Clean up old completed actions
   */
  async cleanupOldActions(maxAge = 24 * 60 * 60 * 1000) { // 24 hours default
    try {
      await this.redis.connect();

      const pattern = `${this.namespace}:*:*:*`;
      const keys = await this.redis.keys(pattern);
      const cutoffTime = Date.now() - maxAge;
      let cleanedCount = 0;

      for (const key of keys) {
        const data = await this.redis.get(key);
        if (data) {
          const parsed = JSON.parse(data);
          if (parsed.timestamp && parsed.timestamp < cutoffTime) {
            await this.redis.del(key);
            cleanedCount++;
          }
        }
      }

      logger.info('Cleaned up old actions', { cleanedCount });
      return cleanedCount;
    } catch (error) {
      logger.error('Failed to cleanup old actions', { error: error.message });
      return 0;
    }
  }

  /**
   * Get collaboration status summary
   */
  async getCollaborationStatus() {
    try {
      const activeActions = await this.getActiveActions();
      const agentStats = {};

      // Group by agent
      for (const action of activeActions) {
        if (!agentStats[action.agentId]) {
          agentStats[action.agentId] = {
            agentId: action.agentId,
            role: action.agentRole,
            activeActions: 0,
            completedActions: 0
          };
        }
        agentStats[action.agentId].activeActions++;
      }

      // Get completed actions count for each agent
      for (const agentId of Object.keys(agentStats)) {
        const allActions = await this.getAgentActions(agentId);
        agentStats[agentId].completedActions = allActions.filter(a => a.status === 'completed').length;
      }

      return {
        totalActiveActions: activeActions.length,
        agentStats,
        timestamp: Date.now()
      };
    } catch (error) {
      logger.error('Failed to get collaboration status', { error: error.message });
      return null;
    }
  }

  /**
   * Close Redis connection
   */
  async close() {
    try {
      await this.redis.quit();
      logger.info('Shared memory Redis connection closed');
    } catch (error) {
      logger.error('Error closing shared memory Redis connection', { error: error.message });
    }
  }
}

// Create singleton instance
const sharedMemory = new SharedMemoryService();

module.exports = sharedMemory;

// Graceful shutdown
process.on('SIGTERM', async () => {
  await sharedMemory.close();
  process.exit(0);
});

process.on('SIGINT', async () => {
  await sharedMemory.close();
  process.exit(0);
});