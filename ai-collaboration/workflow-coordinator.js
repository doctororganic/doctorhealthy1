/**
 * AI Collaboration - Workflow Coordinator
 * Manages sequential and parallel workflows between AI agents
 */

const sharedMemory = require('./file-shared-memory');
const logger = require('../production-nodejs/utils/logger');

// Disable Redis connections in workflow coordinator
process.env.REDIS_HOST = 'disabled';

class WorkflowCoordinator {
  constructor() {
    this.agentRoles = {
      'kilo': {
        name: 'Kilo Code',
        primaryRole: 'frontend_development',
        secondaryRole: 'testing_and_fixations',
        capabilities: ['react', 'vue', 'javascript', 'testing', 'ui/ux']
      },
      'roo': {
        name: 'Roo Code',
        primaryRole: 'backend_integration',
        secondaryRole: 'validation_and_reviewing',
        capabilities: ['nodejs', 'api', 'database', 'validation', 'review']
      },
      'codesupernova': {
        name: 'CodeSupernova',
        primaryRole: 'monitoring_systems',
        secondaryRole: 'system_architecture',
        capabilities: ['monitoring', 'observability', 'docker', 'devops']
      }
    };

    this.workflowStages = {
      'frontend_development': {
        assignedTo: 'kilo',
        dependencies: ['monitoring_system'],
        deliverables: ['modern_web_interface', 'real_time_dashboard', 'mobile_responsive_design']
      },
      'backend_integration': {
        assignedTo: 'roo',
        dependencies: ['frontend_development'],
        deliverables: ['api_enhancement', 'database_integration', 'mobile_sdk']
      },
      'testing_and_fixations': {
        assignedTo: 'kilo',
        dependencies: ['backend_integration'],
        deliverables: ['comprehensive_tests', 'bug_fixes', 'performance_optimization']
      },
      'validation_and_reviewing': {
        assignedTo: 'roo',
        dependencies: ['testing_and_fixations'],
        deliverables: ['code_review', 'security_audit', 'final_validation']
      }
    };
  }

  /**
   * Assign task to specific AI agent
   */
  async assignTask(agentId, taskId, taskData) {
    try {
      const agentRole = this.agentRoles[agentId];
      if (!agentRole) {
        throw new Error(`Unknown agent: ${agentId}`);
      }

      const assignment = {
        ...taskData,
        assignedTo: agentId,
        assignedBy: 'workflow_coordinator',
        assignmentTime: Date.now(),
        agentRole: agentRole.primaryRole,
        status: 'assigned'
      };

      await sharedMemory.setAction(agentId, taskId, assignment);

      logger.info('Task assigned to AI agent', {
        agentId,
        taskId,
        agentRole: agentRole.name,
        taskType: taskData.type
      });

      return assignment;
    } catch (error) {
      logger.error('Failed to assign task', {
        error: error.message,
        agentId,
        taskId
      });
      throw error;
    }
  }

  /**
   * Start workflow stage
   */
  async startWorkflowStage(stageName) {
    try {
      const stage = this.workflowStages[stageName];
      if (!stage) {
        throw new Error(`Unknown workflow stage: ${stageName}`);
      }

      // Check dependencies only if not monitoring_system (which should be pre-completed)
      for (const dependency of stage.dependencies) {
        if (dependency !== 'monitoring_system') {
          await this.checkDependency(stage.assignedTo, dependency);
        }
      }

      // Create stage task
      const stageTask = {
        type: 'workflow_stage',
        stageName,
        assignedTo: stage.assignedTo,
        deliverables: stage.deliverables,
        startTime: Date.now(),
        status: 'in_progress'
      };

      await sharedMemory.setAction(stage.assignedTo, `stage_${stageName}`, stageTask);

      logger.info('Workflow stage started', {
        stageName,
        assignedTo: stage.assignedTo,
        deliverables: stage.deliverables
      });

      return stageTask;
    } catch (error) {
      logger.error('Failed to start workflow stage', {
        error: error.message,
        stageName
      });
      throw error;
    }
  }

  /**
   * Check if dependency is completed
   */
  async checkDependency(agentId, dependencyStage) {
    try {
      const dependency = await sharedMemory.getAction(agentId, `stage_${dependencyStage}`);

      if (!dependency || dependency.status !== 'completed') {
        throw new Error(`Dependency not completed: ${dependencyStage}`);
      }

      return dependency.result;
    } catch (error) {
      logger.error('Dependency check failed', {
        error: error.message,
        agentId,
        dependencyStage
      });
      throw error;
    }
  }

  /**
   * Complete workflow stage
   */
  async completeWorkflowStage(agentId, stageName, result) {
    try {
      const stage = this.workflowStages[stageName];
      if (!stage) {
        throw new Error(`Unknown workflow stage: ${stageName}`);
      }

      // Verify agent is authorized for this stage
      if (stage.assignedTo !== agentId) {
        throw new Error(`Agent ${agentId} not authorized for stage ${stageName}`);
      }

      // Complete the stage
      await sharedMemory.completeAction(agentId, `stage_${stageName}`, result);

      logger.info('Workflow stage completed', {
        agentId,
        stageName,
        result: result.success
      });

      return result;
    } catch (error) {
      logger.error('Failed to complete workflow stage', {
        error: error.message,
        agentId,
        stageName
      });
      throw error;
    }
  }

  /**
   * Get workflow status
   */
  async getWorkflowStatus() {
    try {
      const activeActions = await sharedMemory.getActiveActions();
      const status = {
        totalActiveTasks: activeActions.length,
        stages: {},
        agents: {},
        timestamp: Date.now()
      };

      // Group by stage
      for (const action of activeActions) {
        if (action.type === 'workflow_stage') {
          status.stages[action.stageName] = {
            assignedTo: action.assignedTo,
            status: action.status,
            progress: action.progress || 0,
            startTime: action.startTime
          };
        }

        // Group by agent
        if (!status.agents[action.agentId]) {
          status.agents[action.agentId] = {
            agentId: action.agentId,
            role: action.agentRole,
            activeTasks: 0,
            completedTasks: 0
          };
        }
        status.agents[action.agentId].activeTasks++;
      }

      // Get completed actions count for each agent
      for (const agentId of Object.keys(status.agents)) {
        try {
          const allActions = await sharedMemory.getAgentActions(agentId);
          if (allActions && Array.isArray(allActions)) {
            status.agents[agentId].completedTasks = allActions.filter(a => a.status === 'completed').length;
          } else {
            status.agents[agentId].completedTasks = 0;
          }
        } catch (error) {
          // Skip if unable to get agent actions
          status.agents[agentId].completedTasks = 0;
        }
      }

      return status;
    } catch (error) {
      logger.error('Failed to get workflow status', { error: error.message });
      return null;
    }
  }

  /**
   * Coordinate parallel tasks
   */
  async coordinateParallelTasks(tasks) {
    try {
      const results = [];
      const promises = [];

      // Start all tasks in parallel
      for (const task of tasks) {
        const promise = this.assignTask(task.agentId, task.taskId, task.taskData)
          .then(() => task.taskId);
        promises.push(promise);
      }

      // Wait for all tasks to be assigned
      await Promise.all(promises);

      logger.info('Parallel tasks coordinated', {
        taskCount: tasks.length,
        agents: tasks.map(t => t.agentId)
      });

      return tasks.map(t => t.taskId);
    } catch (error) {
      logger.error('Failed to coordinate parallel tasks', {
        error: error.message,
        taskCount: tasks.length
      });
      throw error;
    }
  }

  /**
   * Sequential workflow execution
   */
  async executeSequentialWorkflow(stages) {
    try {
      const results = [];

      for (const stageName of stages) {
        logger.info('Starting sequential stage', { stageName });

        // Start the stage
        await this.startWorkflowStage(stageName);

        // Wait for completion
        const stage = this.workflowStages[stageName];
        const result = await sharedMemory.waitForDependency(
          stage.assignedTo,
          `stage_${stageName}`
        );

        results.push({ stageName, result });

        logger.info('Sequential stage completed', {
          stageName,
          result: result.success
        });
      }

      return results;
    } catch (error) {
      logger.error('Sequential workflow execution failed', {
        error: error.message,
        completedStages: results.length
      });
      throw error;
    }
  }

  /**
   * Get agent capabilities and status
   */
  async getAgentStatus(agentId) {
    try {
      const agentRole = this.agentRoles[agentId];
      if (!agentRole) {
        throw new Error(`Unknown agent: ${agentId}`);
      }

      const actions = await sharedMemory.getAgentActions(agentId);
      const activeActions = actions.filter(a => a.status === 'active');
      const completedActions = actions.filter(a => a.status === 'completed');

      return {
        agentId,
        name: agentRole.name,
        primaryRole: agentRole.primaryRole,
        secondaryRole: agentRole.secondaryRole,
        capabilities: agentRole.capabilities,
        activeTasks: activeActions.length,
        completedTasks: completedActions.length,
        recentActions: actions.slice(-5) // Last 5 actions
      };
    } catch (error) {
      logger.error('Failed to get agent status', {
        error: error.message,
        agentId
      });
      return null;
    }
  }

  /**
   * Create task handoff between agents
   */
  async createTaskHandoff(fromAgent, toAgent, taskId, context) {
    try {
      const handoff = {
        type: 'task_handoff',
        fromAgent,
        toAgent,
        taskId,
        context,
        handoffTime: Date.now(),
        status: 'pending'
      };

      await sharedMemory.setAction(toAgent, `handoff_${taskId}`, handoff);

      logger.info('Task handoff created', {
        fromAgent,
        toAgent,
        taskId
      });

      return handoff;
    } catch (error) {
      logger.error('Failed to create task handoff', {
        error: error.message,
        fromAgent,
        toAgent,
        taskId
      });
      throw error;
    }
  }

  /**
   * Accept task handoff
   */
  async acceptTaskHandoff(agentId, taskId) {
    try {
      const handoff = await sharedMemory.getAction(agentId, `handoff_${taskId}`);

      if (!handoff) {
        throw new Error(`No handoff found for task: ${taskId}`);
      }

      handoff.status = 'accepted';
      handoff.acceptedTime = Date.now();

      await sharedMemory.setAction(agentId, `handoff_${taskId}`, handoff);

      logger.info('Task handoff accepted', {
        agentId,
        taskId,
        fromAgent: handoff.fromAgent
      });

      return handoff;
    } catch (error) {
      logger.error('Failed to accept task handoff', {
        error: error.message,
        agentId,
        taskId
      });
      throw error;
    }
  }
}

// Create singleton instance
const workflowCoordinator = new WorkflowCoordinator();

module.exports = workflowCoordinator;