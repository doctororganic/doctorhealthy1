/**
 * File-Based Shared Memory Service (Alternative to Redis)
 * Simple, reliable shared memory using file system
 */

const fs = require('fs').promises;
const path = require('path');
const logger = require('../production-nodejs/utils/logger');

class FileSharedMemory {
  constructor(basePath = './ai-collaboration/memory') {
    this.basePath = basePath;
    this.namespace = 'ai_collaboration';
  }

  /**
   * Ensure memory directory exists
   */
  async ensureDirectory() {
    try {
      await fs.mkdir(this.basePath, { recursive: true });
    } catch (error) {
      if (error.code !== 'EEXIST') {
        throw error;
      }
    }
  }

  /**
   * Store action/state for an AI agent
   */
  async setAction(agentId, actionId, data) {
    try {
      await this.ensureDirectory();

      const dir = path.join(this.basePath, agentId);
      await fs.mkdir(dir, { recursive: true });

      const filePath = path.join(dir, `${actionId}.json`);
      const actionData = {
        ...data,
        timestamp: Date.now(),
        agentId,
        status: 'active'
      };

      await fs.writeFile(filePath, JSON.stringify(actionData, null, 2));

      logger.info('Action stored in file-based shared memory', {
        agentId,
        actionId,
        status: 'active',
        filePath
      });

      return actionData;
    } catch (error) {
      logger.error('Failed to set action in file-based shared memory', {
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
      const filePath = path.join(this.basePath, agentId, `${actionId}.json`);

      try {
        const data = await fs.readFile(filePath, 'utf8');
        const parsed = JSON.parse(data);

        logger.debug('Action retrieved from file-based shared memory', {
          agentId,
          actionId,
          status: parsed.status
        });

        return parsed;
      } catch (error) {
        if (error.code === 'ENOENT') {
          return null; // File doesn't exist
        }
        throw error;
      }
    } catch (error) {
      logger.error('Failed to get action from file-based shared memory', {
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
      const current = await this.getAction(agentId, actionId);

      if (current) {
        const completedData = {
          ...current,
          status: 'completed',
          result,
          completedAt: Date.now()
        };

        const filePath = path.join(this.basePath, agentId, `${actionId}.json`);
        await fs.writeFile(filePath, JSON.stringify(completedData, null, 2));

        logger.info('Action completed in file-based shared memory', {
          agentId,
          actionId,
          status: 'completed'
        });

        return completedData;
      }

      return null;
    } catch (error) {
      logger.error('Failed to complete action in file-based shared memory', {
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

    logger.info('Waiting for dependency in file-based shared memory', {
      agentId,
      dependencyActionId,
      timeout
    });

    while (Date.now() - startTime < timeout) {
      const dependency = await this.getAction(agentId, dependencyActionId);

      if (dependency && dependency.status === 'completed') {
        logger.info('Dependency completed in file-based shared memory', {
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
    logger.error('Dependency timeout in file-based shared memory', {
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
      await this.ensureDirectory();

      const agents = ['kilo', 'roo', 'codesupernova'];
      const actions = [];

      for (const agentId of agents) {
        const agentDir = path.join(this.basePath, agentId);

        try {
          const files = await fs.readdir(agentDir);

          for (const file of files) {
            if (file.endsWith('.json')) {
              const filePath = path.join(agentDir, file);
              const data = await fs.readFile(filePath, 'utf8');
              const parsed = JSON.parse(data);

              if (parsed.status === 'active') {
                actions.push(parsed);
              }
            }
          }
        } catch (error) {
          // Agent directory doesn't exist yet, skip
          continue;
        }
      }

      return actions;
    } catch (error) {
      logger.error('Failed to get active actions from file-based shared memory', {
        error: error.message
      });
      return [];
    }
  }

  /**
   * Get actions by specific agent
   */
  async getAgentActions(agentId) {
    try {
      const agentDir = path.join(this.basePath, agentId);
      const actions = [];

      try {
        const files = await fs.readdir(agentDir);

        for (const file of files) {
          if (file.endsWith('.json')) {
            const filePath = path.join(agentDir, file);
            const data = await fs.readFile(filePath, 'utf8');
            actions.push(JSON.parse(data));
          }
        }
      } catch (error) {
        // Agent directory doesn't exist yet
        return [];
      }

      return actions;
    } catch (error) {
      logger.error('Failed to get agent actions from file-based shared memory', {
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
      await this.ensureDirectory();

      const agents = ['kilo', 'roo', 'codesupernova'];
      const cutoffTime = Date.now() - maxAge;
      let cleanedCount = 0;

      for (const agentId of agents) {
        const agentDir = path.join(this.basePath, agentId);

        try {
          const files = await fs.readdir(agentDir);

          for (const file of files) {
            if (file.endsWith('.json')) {
              const filePath = path.join(agentDir, file);
              const data = await fs.readFile(filePath, 'utf8');
              const parsed = JSON.parse(data);

              if (parsed.timestamp && parsed.timestamp < cutoffTime) {
                await fs.unlink(filePath);
                cleanedCount++;
              }
            }
          }
        } catch (error) {
          // Agent directory doesn't exist yet, skip
          continue;
        }
      }

      logger.info('Cleaned up old actions from file-based shared memory', { cleanedCount });
      return cleanedCount;
    } catch (error) {
      logger.error('Failed to cleanup old actions from file-based shared memory', {
        error: error.message
      });
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
            activeActions: 0,
            completedActions: 0
          };
        }
        agentStats[action.agentId].activeActions++;
      }

      // Get completed actions count for each agent
      const agents = ['kilo', 'roo', 'codesupernova'];
      for (const agentId of agents) {
        const allActions = await this.getAgentActions(agentId);
        agentStats[agentId].completedActions = allActions.filter(a => a.status === 'completed').length;
      }

      return {
        totalActiveActions: activeActions.length,
        agentStats,
        timestamp: Date.now()
      };
    } catch (error) {
      logger.error('Failed to get collaboration status from file-based shared memory', {
        error: error.message
      });
      return null;
    }
  }
}

// Create singleton instance
const fileSharedMemory = new FileSharedMemory();

module.exports = fileSharedMemory;