const http = require('http');

const tests = [
  {
    name: 'Health Check',
    method: 'GET',
    path: '/health',
    expected: 'healthy'
  },
  {
    name: 'API Info',
    method: 'GET',
    path: '/api/info',
    expected: 'Trae New Healthy1'
  },
  {
    name: 'Nutrition Analysis - Apple',
    method: 'POST',
    path: '/api/nutrition/analyze',
    body: {food: 'apple', quantity: 100, unit: 'g', checkHalal: true},
    expected: 'success'
  },
  {
    name: 'Nutrition Analysis - Chicken',
    method: 'POST',
    path: '/api/nutrition/analyze',
    body: {food: 'chicken', quantity: 200, unit: 'g', checkHalal: true},
    expected: 'success'
  }
];

async function runTests() {
  console.log('🧪 Running API Tests...\n');

  for (const test of tests) {
    try {
      const result = await makeRequest(test);
      console.log(`✅ ${test.name}: PASSED`);
    } catch (error) {
      console.log(`❌ ${test.name}: FAILED - ${error.message}`);
    }
  }
}

function makeRequest(test) {
  return new Promise((resolve, reject) => {
    const options = {
      hostname: 'localhost',
      port: 8080,
      path: test.path,
      method: test.method,
      headers: {'Content-Type': 'application/json'}
    };

    const req = http.request(options, (res) => {
      let data = '';
      res.on('data', chunk => data += chunk);
      res.on('end', () => {
        if (res.statusCode === 200 && data.includes(test.expected)) {
          resolve(data);
        } else {
          reject(new Error(`Status ${res.statusCode}, expected "${test.expected}"`));
        }
      });
    });

    req.on('error', reject);
    if (test.body) req.write(JSON.stringify(test.body));
    req.end();
  });
}

runTests();