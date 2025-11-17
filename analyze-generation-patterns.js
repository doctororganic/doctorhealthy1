#!/usr/bin/env node

/**
 * Generation Pattern Analysis Tool
 * Analyzes API responses for repetitive patterns, variety, and generation issues
 */

const axios = require('axios');
const fs = require('fs');

const API_BASE_URL = 'http://localhost:8080';

class GenerationAnalyzer {
    constructor() {
        this.results = {
            endpoints: {},
            patterns: [],
            issues: [],
            recommendations: []
        };
        this.sampleSize = 10; // Number of requests per endpoint
    }

    async analyzeEndpoint(endpoint, description, params = {}) {
        console.log(`\n🔍 Analyzing: ${description} (${endpoint})`);
        
        const responses = [];
        const patterns = {
            identical: 0,
            similar: 0,
            unique: 0,
            avgLength: 0,
            variance: 0,
            commonWords: new Map(),
            structure: new Map()
        };

        // Make multiple requests to the same endpoint
        for (let i = 0; i < this.sampleSize; i++) {
            try {
                const response = await axios.get(`${API_BASE_URL}${endpoint}`, {
                    params: params,
                    timeout: 5000
                });
                
                if (response.status === 200) {
                    const data = JSON.stringify(response.data);
                    responses.push({
                        data: response.data,
                        string: data,
                        length: data.length,
                        structure: this.getStructure(response.data)
                    });
                }
            } catch (error) {
                console.log(`⚠️  Request ${i + 1} failed: ${error.message}`);
            }
        }

        if (responses.length === 0) {
            console.log(`❌ No successful responses for ${endpoint}`);
            return;
        }

        // Analyze patterns
        this.analyzeResponsePatterns(responses, patterns);
        
        // Store results
        this.results.endpoints[endpoint] = {
            description,
            responses: responses.length,
            patterns,
            sampleResponses: responses.slice(0, 3) // Keep first 3 for analysis
        };

        console.log(`✅ Analyzed ${responses.length} responses`);
        console.log(`   Identical: ${patterns.identical}, Similar: ${patterns.similar}, Unique: ${patterns.unique}`);
    }

    getStructure(data) {
        const structure = [];
        
        function traverse(obj, path = '') {
            if (obj === null || obj === undefined) return;
            
            if (Array.isArray(obj)) {
                structure.push(`${path}[]`);
                if (obj.length > 0) {
                    traverse(obj[0], path + '[0]');
                }
            } else if (typeof obj === 'object') {
                Object.keys(obj).forEach(key => {
                    structure.push(`${path}.${key}`);
                    traverse(obj[key], path + '.' + key);
                });
            } else {
                structure.push(`${path}(${typeof obj})`);
            }
        }
        
        traverse(data);
        return structure.sort().join('|');
    }

    analyzeResponsePatterns(responses, patterns) {
        const lengths = responses.map(r => r.length);
        patterns.avgLength = lengths.reduce((a, b) => a + b, 0) / lengths.length;
        patterns.variance = this.calculateVariance(lengths);
        
        // Check for identical responses
        const responseStrings = responses.map(r => r.string);
        const uniqueStrings = [...new Set(responseStrings)];
        
        patterns.identical = responses.length - uniqueStrings.length;
        patterns.unique = uniqueStrings.length;
        
        // Check for similar responses (structure similarity)
        const structures = responses.map(r => r.structure);
        const uniqueStructures = [...new Set(structures)];
        
        if (uniqueStructures.length === 1) {
            patterns.similar = responses.length;
        } else {
            patterns.similar = Math.max(...uniqueStructures.map(s => 
                structures.filter(st => st === s).length
            ));
        }
        
        // Analyze common words/phrases
        const allText = responseStrings.join(' ').toLowerCase();
        const words = allText.match(/\b[a-z]+\b/g) || [];
        
        words.forEach(word => {
            if (word.length > 3) { // Only consider words longer than 3 chars
                patterns.commonWords.set(word, (patterns.commonWords.get(word) || 0) + 1);
            }
        });
        
        // Sort by frequency
        patterns.commonWords = new Map([...patterns.commonWords.entries()]
            .sort((a, b) => b[1] - a[1])
            .slice(0, 10));
    }

    calculateVariance(numbers) {
        const mean = numbers.reduce((a, b) => a + b, 0) / numbers.length;
        const squaredDiffs = numbers.map(num => Math.pow(num - mean, 2));
        return squaredDiffs.reduce((a, b) => a + b, 0) / numbers.length;
    }

    async analyzeAllEndpoints() {
        console.log('🚀 Starting Generation Pattern Analysis...\n');

        // Define endpoints to analyze
        const endpoints = [
            // Health endpoints
            { path: '/api/v1/health/conditions', desc: 'Health Conditions' },
            { path: '/api/v1/health/tips', desc: 'Health Tips' },
            
            // Nutrition data endpoints
            { path: '/api/v1/metabolism', desc: 'Metabolism Guide' },
            { path: '/api/v1/meal-plans', desc: 'Meal Plans' },
            { path: '/api/v1/vitamins-minerals', desc: 'Vitamins & Minerals' },
            { path: '/api/v1/workout-techniques', desc: 'Workout Techniques' },
            { path: '/api/v1/calories', desc: 'Calories Data' },
            { path: '/api/v1/skills', desc: 'Skills Data' },
            { path: '/api/v1/diseases', desc: 'Disease Data' },
            { path: '/api/v1/type-plans', desc: 'Type Plans' },
            
            // Dynamic endpoints with different parameters
            { path: '/api/v1/health/conditions', desc: 'Health Conditions (Filtered)', params: { category: 'chronic' } },
            { path: '/api/v1/calories/weight-loss', desc: 'Calories - Weight Loss' },
            { path: '/api/v1/skills/beginner', desc: 'Skills - Beginner' },
            { path: '/api/v1/type-plans/muscle', desc: 'Type Plans - Muscle' },
            
            // Recipe endpoint with different parameters
            { path: '/api/v1/recipes', desc: 'Recipes - Breakfast', params: { meal_type: 'Breakfast' } },
            { path: '/api/v1/recipes', desc: 'Recipes - Lunch', params: { meal_type: 'Lunch' } },
            { path: '/api/v1/recipes', desc: 'Recipes - Dinner', params: { meal_type: 'Dinner' } },
            { path: '/api/v1/recipes', desc: 'Recipes - Medical', params: { medical_conditions: 'diabetes' } },
            
            // Workout endpoint with different parameters
            { path: '/api/v1/workouts', desc: 'Workouts - Weight Loss', params: { goal: 'weight_loss' } },
            { path: '/api/v1/workouts', desc: 'Workouts - Muscle Gain', params: { goal: 'muscle_gain' } },
            { path: '/api/v1/workouts', desc: 'Workouts - Injury Modified', params: { injury_location: 'shoulder', injury_status: 'recovering' } }
        ];

        // Analyze each endpoint
        for (const endpoint of endpoints) {
            await this.analyzeEndpoint(endpoint.path, endpoint.desc, endpoint.params || {});
            
            // Add small delay between requests to avoid overwhelming the server
            await new Promise(resolve => setTimeout(resolve, 100));
        }

        // Generate analysis report
        this.generateAnalysisReport();
    }

    generateAnalysisReport() {
        console.log('\n📊 GENERATING ANALYSIS REPORT...\n');

        // Identify patterns and issues
        this.identifyPatterns();
        this.identifyIssues();
        this.generateRecommendations();

        // Generate comprehensive report
        const report = {
            timestamp: new Date().toISOString(),
            summary: {
                totalEndpoints: Object.keys(this.results.endpoints).length,
                totalRequests: Object.values(this.results.endpoints)
                    .reduce((sum, ep) => sum + ep.responses, 0),
                analysisDate: new Date().toLocaleDateString()
            },
            endpoints: this.results.endpoints,
            patterns: this.patterns,
            issues: this.issues,
            recommendations: this.recommendations
        };

        // Save report
        fs.writeFileSync('generation-analysis-report.json', JSON.stringify(report, null, 2));
        console.log('📄 Analysis report saved to: generation-analysis-report.json');

        // Print summary
        this.printSummary();
    }

    identifyPatterns() {
        console.log('🔍 Identifying Patterns...');
        
        const endpoints = this.results.endpoints;
        
        // Check for identical responses across endpoints
        Object.entries(endpoints).forEach(([endpoint, data]) => {
            if (data.patterns.identical > data.responses * 0.8) {
                this.patterns.push({
                    type: 'IDENTICAL_RESPONSES',
                    severity: 'HIGH',
                    endpoint,
                    description: `80%+ responses are identical`,
                    impact: 'No variety in responses'
                });
            }
            
            if (data.patterns.similar > data.responses * 0.9) {
                this.patterns.push({
                    type: 'SIMILAR_STRUCTURE',
                    severity: 'MEDIUM',
                    endpoint,
                    description: `90%+ responses have similar structure`,
                    impact: 'Limited structural variety'
                });
            }
            
            if (data.patterns.variance < 100) {
                this.patterns.push({
                    type: 'LOW_VARIANCE',
                    severity: 'MEDIUM',
                    endpoint,
                    description: `Very low response length variance (${data.patterns.variance.toFixed(2)})`,
                    impact: 'Responses might be too uniform'
                });
            }
        });

        // Check for common repetitive words
        const allCommonWords = new Map();
        Object.values(endpoints).forEach(data => {
            data.patterns.commonWords.forEach((count, word) => {
                allCommonWords.set(word, (allCommonWords.get(word) || 0) + count);
            });
        });

        // Find overly common words
        allCommonWords.forEach((count, word) => {
            if (count > 50) {
                this.patterns.push({
                    type: 'REPETITIVE_WORDS',
                    severity: 'LOW',
                    word,
                    count,
                    description: `Word "${word}" appears ${count} times across all responses`
                });
            }
        });
    }

    identifyIssues() {
        console.log('⚠️  Identifying Issues...');
        
        Object.entries(this.results.endpoints).forEach(([endpoint, data]) => {
            // Check for empty responses
            if (data.responses === 0) {
                this.issues.push({
                    type: 'NO_RESPONSES',
                    severity: 'HIGH',
                    endpoint,
                    description: 'No successful responses received'
                });
            }
            
            // Check for very short responses
            if (data.patterns.avgLength < 100) {
                this.issues.push({
                    type: 'SHORT_RESPONSES',
                    severity: 'MEDIUM',
                    endpoint,
                    description: `Average response length too short: ${data.patterns.avgLength.toFixed(0)} chars`
                });
            }
            
            // Check for completely static responses
            if (data.patterns.identical === data.responses && data.responses > 1) {
                this.issues.push({
                    type: 'STATIC_RESPONSES',
                    severity: 'HIGH',
                    endpoint,
                    description: 'All responses are identical - no dynamic content'
                });
            }
        });

        // Check for missing variety across different parameters
        const recipeEndpoints = Object.entries(endpoints).filter(([ep]) => ep.includes('recipes'));
        if (recipeEndpoints.length > 1) {
            const recipeStructures = recipeEndpoints.map(([_, data]) => data.patterns.structure);
            const uniqueStructures = [...new Set(recipeStructures)];
            
            if (uniqueStructures.length === 1) {
                this.issues.push({
                    type: 'NO_PARAMETER_VARIETY',
                    severity: 'HIGH',
                    endpoint: '/api/v1/recipes',
                    description: 'Recipe endpoint returns same structure regardless of parameters'
                });
            }
        }
    }

    generateRecommendations() {
        console.log('💡 Generating Recommendations...');
        
        // Recommendations for repetitive patterns
        const identicalCount = this.patterns.filter(p => p.type === 'IDENTICAL_RESPONSES').length;
        if (identicalCount > 0) {
            this.recommendations.push({
                priority: 'HIGH',
                category: 'CONTENT_DIVERSITY',
                title: 'Implement Response Randomization',
                description: `Found ${identicalCount} endpoints with identical responses. Implement randomized content generation to increase variety.`,
                implementation: [
                    'Add randomization functions for content selection',
                    'Implement multiple response templates',
                    'Use dynamic content pools with random selection'
                ]
            });
        }

        // Recommendations for static responses
        const staticCount = this.issues.filter(i => i.type === 'STATIC_RESPONSES').length;
        if (staticCount > 0) {
            this.recommendations.push({
                priority: 'HIGH',
                category: 'DYNAMIC_CONTENT',
                title: 'Add Dynamic Content Generation',
                description: `Found ${staticCount} endpoints with completely static responses. Implement dynamic content based on parameters and context.`,
                implementation: [
                    'Create content variation algorithms',
                    'Implement parameter-based content filtering',
                    'Add contextual content adaptation'
                ]
            });
        }

        // Recommendations for parameter handling
        const noVarietyIssues = this.issues.filter(i => i.type === 'NO_PARAMETER_VARIETY');
        if (noVarietyIssues.length > 0) {
            this.recommendations.push({
                priority: 'HIGH',
                category: 'PARAMETER_PROCESSING',
                title: 'Improve Parameter-Based Response Generation',
                description: 'Endpoints are not properly utilizing request parameters to generate varied responses.',
                implementation: [
                    'Implement proper parameter validation and processing',
                    'Create parameter-specific content pools',
                    'Add parameter combination logic for unique responses'
                ]
            });
        }

        // Recommendations for content pools
        this.recommendations.push({
            priority: 'MEDIUM',
            category: 'CONTENT_MANAGEMENT',
            title: 'Expand Content Pools',
            description: 'Increase the size and variety of content pools to prevent repetition.',
            implementation: [
                'Create larger datasets for each content type',
                'Implement content categorization and tagging',
                'Add seasonal or time-based content variation'
            ]
        });

        // Recommendations for response variation
        this.recommendations.push({
            priority: 'MEDIUM',
            category: 'RESPONSE_LOGIC',
            title: 'Implement Response Variation Algorithms',
            description: 'Add algorithms to ensure responses vary even with similar inputs.',
            implementation: [
                'Add random seed generation',
                'Implement content shuffling',
                'Create response combination logic'
            ]
        });

        // Recommendations for quality assurance
        this.recommendations.push({
            priority: 'LOW',
            category: 'QUALITY_ASSURANCE',
            title: 'Implement Content Diversity Testing',
            description: 'Set up automated testing to monitor content diversity over time.',
            implementation: [
                'Create diversity metrics and monitoring',
                'Implement content similarity detection',
                    'Add automated variation testing in CI/CD'
                ]
            });
    }

    printSummary() {
        console.log('\n📋 ANALYSIS SUMMARY');
        console.log('==================');
        console.log(`Total Endpoints Analyzed: ${Object.keys(this.results.endpoints).length}`);
        console.log(`Total Requests Made: ${Object.values(this.results.endpoints).reduce((sum, ep) => sum + ep.responses, 0)}`);
        console.log(`Patterns Identified: ${this.patterns.length}`);
        console.log(`Issues Found: ${this.issues.length}`);
        console.log(`Recommendations Generated: ${this.recommendations.length}`);

        console.log('\n🚨 CRITICAL ISSUES:');
        this.issues
            .filter(issue => issue.severity === 'HIGH')
            .forEach(issue => {
                console.log(`  • ${issue.endpoint}: ${issue.description}`);
            });

        console.log('\n💡 TOP RECOMMENDATIONS:');
        this.recommendations
            .filter(rec => rec.priority === 'HIGH')
            .forEach(rec => {
                console.log(`  • ${rec.title}`);
                console.log(`    ${rec.description}`);
            });

        console.log('\n📊 DIVERSITY SCORES:');
        Object.entries(this.results.endpoints).forEach(([endpoint, data]) => {
            const diversityScore = ((data.patterns.unique / Math.max(data.responses, 1)) * 100).toFixed(1);
            console.log(`  ${endpoint}: ${diversityScore}% unique responses`);
        });
    }
}

// Main execution
async function main() {
    const analyzer = new GenerationAnalyzer();
    
    try {
        await analyzer.analyzeAllEndpoints();
    } catch (error) {
        console.error('💥 Analysis failed:', error.message);
        process.exit(1);
    }
}

// Handle graceful shutdown
process.on('SIGINT', () => {
    console.log('\n🛑 Analysis interrupted');
    process.exit(1);
});

// Run analysis if this file is executed directly
if (require.main === module) {
    main().catch(console.error);
}

module.exports = GenerationAnalyzer;
