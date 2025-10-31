const http = require('http');
const url = require('url');

const nutritionDB = {
  apple: {calories: 52, protein: 0.3, carbs: 14, fat: 0.2, fiber: 2.4, sugar: 10.4},
  banana: {calories: 89, protein: 1.1, carbs: 23, fat: 0.3, fiber: 2.6, sugar: 12.2},
  chicken: {calories: 165, protein: 31, carbs: 0, fat: 3.6, fiber: 0, sugar: 0},
  rice: {calories: 130, protein: 2.7, carbs: 28, fat: 0.3, fiber: 0.4, sugar: 0.1},
  bread: {calories: 265, protein: 9, carbs: 49, fat: 3.2, fiber: 2.7, sugar: 5.7},
  egg: {calories: 155, protein: 13, carbs: 1.1, fat: 11, fiber: 0, sugar: 1.1},
  milk: {calories: 42, protein: 3.4, carbs: 5, fat: 1, fiber: 0, sugar: 5},
  orange: {calories: 47, protein: 0.9, carbs: 12, fat: 0.1, fiber: 2.4, sugar: 9.4}
};

const halalFoods = {
  apple: true, banana: true, orange: true, rice: true,
  bread: true, egg: true, milk: true, chicken: true
};

const server = http.createServer((req, res) => {
  const parsedUrl = url.parse(req.url, true);
  const path = parsedUrl.pathname;
  const method = req.method;

  res.setHeader('Access-Control-Allow-Origin', '*');
  res.setHeader('Access-Control-Allow-Methods', 'GET, POST, OPTIONS');
  res.setHeader('Access-Control-Allow-Headers', 'Content-Type');

  if (method === 'OPTIONS') {
    res.writeHead(200);
    res.end();
    return;
  }

  if (path === '/' && method === 'GET') {
    res.writeHead(200, {'Content-Type': 'text/html'});
    res.end(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Trae New Healthy1 - AI Nutrition Platform</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh; color: #333;
        }
        .container { max-width: 1200px; margin: 0 auto; padding: 20px; }
        .header { text-align: center; color: white; margin-bottom: 40px; padding: 40px 0; }
        .header h1 { font-size: 3em; margin-bottom: 10px; text-shadow: 2px 2px 4px rgba(0,0,0,0.3); }
        .main-content { background: white; border-radius: 20px; padding: 40px; box-shadow: 0 20px 40px rgba(0,0,0,0.1); }
        .status-badge { background: #27ae60; color: white; padding: 15px 30px; border-radius: 50px; text-align: center; font-size: 1.1em; font-weight: bold; margin-bottom: 30px; }
        .test-section { background: #ecf0f1; padding: 30px; border-radius: 15px; margin: 30px 0; }
        .test-button { background: #3498db; color: white; padding: 15px 30px; border: none; border-radius: 8px; font-size: 16px; cursor: pointer; }
        .result { margin-top: 20px; padding: 20px; background: white; border-radius: 8px; display: none; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🍎 Trae New Healthy1</h1>
            <p>AI-Powered Nutrition & Health Management Platform</p>
        </div>
        <div class="main-content">
            <div class="status-badge">✅ Platform is LIVE and Ready!</div>
            <div class="test-section">
                <h3>🧪 Test Nutrition Analysis</h3>
                <button class="test-button" onclick="testHealth()">Test Health Endpoint</button>
                <div id="result" class="result"></div>
            </div>
        </div>
    </div>
    <script>
        async function testHealth() {
            const resultDiv = document.getElementById('result');
            resultDiv.style.display = 'block';
            resultDiv.innerHTML = '<p>🔄 Testing health endpoint...</p>';

            try {
                const response = await fetch('/health');
                const data = await response.json();
                resultDiv.innerHTML = '<h4>✅ Health Check Passed</h4><pre>' + JSON.stringify(data, null, 2) + '</pre>';
            } catch (error) {
                resultDiv.innerHTML = '<p>❌ Error: ' + error.message + '</p>';
            }
        }
    </script>
</body>
</html>`);
    return;
  }

  if (path === '/health' && method === 'GET') {
    res.writeHead(200, {'Content-Type': 'application/json'});
    res.end(JSON.stringify({
      status: 'healthy',
      timestamp: new Date().toISOString(),
      uptime: process.uptime(),
      message: 'Trae New Healthy1 is running',
      version: '1.0.0'
    }));
    return;
  }

  if (path === '/api/info' && method === 'GET') {
    res.writeHead(200, {'Content-Type': 'application/json'});
    res.end(JSON.stringify({
      name: 'Trae New Healthy1',
      description: 'AI-powered nutrition platform',
      version: '1.0.0',
      status: 'active',
      endpoints: {
        health: '/health',
        nutrition: '/api/nutrition/analyze',
        info: '/api/info'
      }
    }));
    return;
  }

  if (path === '/api/nutrition/analyze' && method === 'POST') {
    let body = '';
    req.on('data', chunk => { body += chunk.toString(); });
    req.on('end', () => {
      try {
        const data = JSON.parse(body);
        const { food, quantity, unit, checkHalal } = data;

        const foodData = nutritionDB[food] || {
          calories: 100, protein: 5, carbs: 15, fat: 2, fiber: 1, sugar: 5
        };

        const multiplier = quantity / 100;
        const isHalal = halalFoods[food] || false;

        res.writeHead(200, {'Content-Type': 'application/json'});
        res.end(JSON.stringify({
          food, quantity, unit,
          calories: foodData.calories * multiplier,
          protein: foodData.protein * multiplier,
          carbs: foodData.carbs * multiplier,
          fat: foodData.fat * multiplier,
          fiber: foodData.fiber * multiplier,
          sugar: foodData.sugar * multiplier,
          isHalal,
          status: 'success',
          message: 'Analysis completed',
          timestamp: new Date().toISOString()
        }));
      } catch (error) {
        res.writeHead(400, {'Content-Type': 'application/json'});
        res.end(JSON.stringify({error: 'Invalid JSON'}));
      }
    });
    return;
  }

  res.writeHead(404, {'Content-Type': 'application/json'});
  res.end(JSON.stringify({error: 'Not found'}));
});

const port = process.env.PORT || 8080;
server.listen(port, '0.0.0.0', () => {
  console.log(`🚀 Server running on port ${port}`);
});