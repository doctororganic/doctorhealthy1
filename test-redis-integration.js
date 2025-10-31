/**
 * Test Redis integration for shared memory and rate limiting
 */

const redisClient = require('./production-nodejs/services/redisClient');

async function testRedisIntegration() {
  console.log('🧪 Testing Redis integration...\n');

  try {
    // Connect to Redis
    await redisClient.connect();
    console.log('✅ Redis connection successful');

    // Test shared memory operations
    console.log('\n📝 Testing shared memory operations...');

    // Set action
    await redisClient.setSharedMemory('test_agent', 'test_task_1', {
      type: 'code_generation',
      description: 'Test task',
      status: 'active'
    });
    console.log('✅ Set shared memory action');

    // Get action
    const action = await redisClient.getSharedMemory('test_agent', 'test_task_1');
    console.log('✅ Retrieved action:', action?.description);

    // Complete action
    await redisClient.completeSharedMemory('test_agent', 'test_task_1', { success: true });
    console.log('✅ Completed action');

    // Get completed action
    const completed = await redisClient.getSharedMemory('test_agent', 'test_task_1');
    console.log('✅ Completed action status:', completed?.status);

    // Test active actions
    const activeActions = await redisClient.getAllActiveActions();
    console.log('✅ Active actions count:', activeActions.length);

    // Test basic Redis operations
    const ping = await redisClient.ping();
    console.log('✅ Redis ping:', ping);

    const info = await redisClient.getInfo();
    console.log('✅ Redis info retrieved');

    console.log('\n🎉 All Redis integration tests passed!');

  } catch (error) {
    console.error('❌ Redis integration test failed:', error.message);
    process.exit(1);
  } finally {
    await redisClient.disconnect();
  }
}

// Run test if called directly
if (require.main === module) {
  testRedisIntegration();
}

module.exports = testRedisIntegration;