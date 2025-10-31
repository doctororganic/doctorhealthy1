// Diet Planning System JavaScript

// BMR Calculation Functions
function calculateBMR(weight, height, age, gender, formula = 'mifflin') {
    let bmr;
    
    if (formula === 'mifflin') {
        // Mifflin-St Jeor Equation
        if (gender === 'male') {
            bmr = (10 * weight) + (6.25 * height) - (5 * age) + 5;
        } else {
            bmr = (10 * weight) + (6.25 * height) - (5 * age) - 161;
        }
    } else {
        // Harris-Benedict Equation
        if (gender === 'male') {
            bmr = 88.362 + (13.397 * weight) + (4.799 * height) - (5.677 * age);
        } else {
            bmr = 447.593 + (9.247 * weight) + (3.098 * height) - (4.330 * age);
        }
    }
    
    return Math.round(bmr);
}

// TDEE Calculation
function calculateTDEE(bmr, activityLevel) {
    const activityMultipliers = {
        sedentary: 1.2,
        light: 1.375,
        moderate: 1.55,
        very: 1.725,
        extra: 1.9
    };
    
    return Math.round(bmr * activityMultipliers[activityLevel]);
}

// Macro Breakdown Calculation
function calculateMacros(tdee, goal, weight) {
    let calories = tdee;
    let protein, carbs, fat;
    
    // Adjust calories based on goal
    switch (goal) {
        case 'lose':
            calories = tdee - 500; // 500 calorie deficit
            break;
        case 'gain':
            calories = tdee + 500; // 500 calorie surplus
            break;
        case 'muscle':
            calories = tdee + 300; // 300 calorie surplus for muscle building
            break;
        default:
            calories = tdee; // maintain
    }
    
    // Calculate macros
    if (goal === 'muscle') {
        protein = weight * 2.2; // Higher protein for muscle building
        fat = weight * 1.0;
        carbs = (calories - (protein * 4) - (fat * 9)) / 4;
    } else if (goal === 'lose') {
        protein = weight * 2.0; // Higher protein to preserve muscle
        fat = weight * 0.8;
        carbs = (calories - (protein * 4) - (fat * 9)) / 4;
    } else {
        protein = weight * 1.6; // Standard protein
        fat = weight * 1.0;
        carbs = (calories - (protein * 4) - (fat * 9)) / 4;
    }
    
    return {
        calories: Math.round(calories),
        protein: Math.round(protein),
        carbs: Math.round(Math.max(carbs, 100)), // Minimum 100g carbs
        fat: Math.round(fat)
    };
}

// Diet Type Recommendation
function recommendDietType(goal, dietary, health = null) {
    if (health) {
        // Medical recommendations
        switch (health) {
            case 'diabetes_type1':
            case 'diabetes_type2':
                return 'Low Glycemic Index Diet';
            case 'hypertension':
                return 'DASH Diet';
            case 'heart_disease':
                return 'Mediterranean Diet';
            case 'kidney_disease':
                return 'Low Protein Diet';
            case 'liver_disease':
                return 'Low Sodium Diet';
            case 'obesity':
                return 'Calorie-Controlled Diet';
            case 'celiac':
                return 'Gluten-Free Diet';
            case 'ibs':
                return 'Low FODMAP Diet';
            default:
                return 'Balanced Medical Diet';
        }
    }
    
    // Regular recommendations
    if (dietary === 'keto') return 'Ketogenic Diet';
    if (dietary === 'paleo') return 'Paleo Diet';
    if (dietary === 'vegetarian') return 'Vegetarian Mediterranean Diet';
    if (dietary === 'vegan') return 'Plant-Based Diet';
    if (dietary === 'low_carb') return 'Low Carbohydrate Diet';
    
    switch (goal) {
        case 'lose':
            return 'Balanced Weight Loss Diet';
        case 'gain':
            return 'High Calorie Balanced Diet';
        case 'muscle':
            return 'High Protein Diet';
        default:
            return 'Balanced Mediterranean Diet';
    }
}

// Sample Meal Plans
const mealPlans = {
    healthy: {
        day1: {
            breakfast: {
                name: 'Oatmeal with Berries',
                ingredients: ['1 cup oats', '1 cup milk', '1/2 cup mixed berries', '1 tbsp honey', '1 tbsp almonds'],
                calories: 350,
                protein: 12,
                carbs: 58,
                fat: 8,
                preparation: 'Cook oats with milk, top with berries, honey, and almonds'
            },
            snack1: {
                name: 'Greek Yogurt with Nuts',
                ingredients: ['1 cup Greek yogurt', '1 oz mixed nuts'],
                calories: 200,
                protein: 15,
                carbs: 8,
                fat: 12,
                preparation: 'Mix yogurt with nuts'
            },
            lunch: {
                name: 'Grilled Chicken Salad',
                ingredients: ['4 oz grilled chicken', '2 cups mixed greens', '1/2 avocado', '1 tbsp olive oil', '1 tbsp lemon juice'],
                calories: 400,
                protein: 35,
                carbs: 12,
                fat: 25,
                preparation: 'Grill chicken, mix with greens, avocado, and dressing'
            },
            snack2: {
                name: 'Apple with Peanut Butter',
                ingredients: ['1 medium apple', '2 tbsp peanut butter'],
                calories: 280,
                protein: 8,
                carbs: 25,
                fat: 16,
                preparation: 'Slice apple and serve with peanut butter'
            },
            dinner: {
                name: 'Baked Salmon with Quinoa',
                ingredients: ['5 oz salmon', '1 cup cooked quinoa', '1 cup steamed broccoli', '1 tbsp olive oil'],
                calories: 550,
                protein: 40,
                carbs: 45,
                fat: 22,
                preparation: 'Bake salmon, serve with quinoa and steamed broccoli'
            }
        }
        // Add more days...
    },
    medical: {
        diabetes: {
            day1: {
                breakfast: {
                    name: 'Low-GI Breakfast Bowl',
                    ingredients: ['1/2 cup steel-cut oats', '1 tbsp chia seeds', '1/4 cup blueberries', '1 tbsp almond butter'],
                    calories: 300,
                    protein: 10,
                    carbs: 35,
                    fat: 12,
                    preparation: 'Cook oats, add toppings',
                    notes: 'Low glycemic index, high fiber'
                }
                // Add more meals...
            }
        }
    }
};

// Weekly Meal Plan Generation Functions
async function generateWeeklyMealPlan(userType) {
    const weekSelector = document.getElementById(`${userType}WeekSelector`);
    const generateBtn = document.getElementById(`${userType}GenerateBtn`);
    const loadingSpinner = document.getElementById(`${userType}LoadingSpinner`);
    const mealPlanDisplay = document.getElementById(`${userType}MealPlanDisplay`);
    const weeklyPlanSection = document.getElementById(`${userType}WeeklyPlan`);
    
    const selectedWeek = weekSelector.value;
    
    // Show loading state
    generateBtn.disabled = true;
    loadingSpinner.style.display = 'block';
    mealPlanDisplay.style.display = 'none';
    
    try {
        // Collect user data
        const userData = collectUserData(userType);
        
        // Call API to generate meal plan
        const response = await fetch('/api/generate-meal-plan', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${getApiKey()}`
            },
            body: JSON.stringify({
                ...userData,
                week: selectedWeek,
                language: getCurrentLanguage()
            })
        });
        
        if (!response.ok) {
            throw new Error('Failed to generate meal plan');
        }
        
        const mealPlan = await response.json();
        
        // Display the generated meal plan
        displayWeeklyMealPlan(userType, mealPlan);
        
        // Show the meal plan display
        loadingSpinner.style.display = 'none';
        mealPlanDisplay.style.display = 'block';
        
    } catch (error) {
        console.error('Error generating meal plan:', error);
        alert('Failed to generate meal plan. Please try again.');
        loadingSpinner.style.display = 'none';
    } finally {
        generateBtn.disabled = false;
    }
}

function collectUserData(userType) {
    const prefix = userType === 'healthy' ? 'healthy' : 'patient';
    
    const userData = {
        name: document.getElementById(`${prefix}Name`).value,
        age: parseInt(document.getElementById(`${prefix}Age`).value),
        height: parseInt(document.getElementById(`${prefix}Height`).value),
        weight: parseInt(document.getElementById(`${prefix}Weight`).value),
        gender: document.getElementById(`${prefix}Gender`).value,
        activityLevel: document.getElementById(`${prefix}Activity`).value,
        goal: document.getElementById(`${prefix}Goal`).value,
        cuisine: document.getElementById(`${prefix}Cuisine`).value,
        dietary: document.getElementById(`${prefix}Dietary`).value,
        allergies: document.getElementById(`${prefix}Allergies`).value
    };
    
    // Add medical information for patients
    if (userType === 'patients') {
        userData.disease = document.getElementById('patientDisease').value;
        userData.medications = document.getElementById('patientMedications').value;
        userData.symptoms = document.getElementById('patientSymptoms').value;
    }
    
    return userData;
}

function displayWeeklyMealPlan(userType, mealPlan) {
    const gridContainer = document.getElementById(`${userType}WeeklyPlanGrid`);
    gridContainer.innerHTML = '';
    
    const days = ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday'];
    
    days.forEach((day, index) => {
        const dayPlan = mealPlan.days[index];
        const dayCard = createDayCard(day, dayPlan);
        gridContainer.appendChild(dayCard);
    });
}

function createDayCard(day, dayPlan) {
    const dayCard = document.createElement('div');
    dayCard.className = 'day-plan-card';
    
    dayCard.innerHTML = `
        <div class="day-header">
            <h5 data-translate="${day.toLowerCase()}">${day}</h5>
            <span class="total-calories">${dayPlan.totalCalories} kcal</span>
        </div>
        <div class="meals-container">
            ${createMealCard('Breakfast', dayPlan.breakfast)}
            ${createMealCard('Snack 1', dayPlan.snack1)}
            ${createMealCard('Lunch', dayPlan.lunch)}
            ${createMealCard('Snack 2', dayPlan.snack2)}
            ${createMealCard('Dinner', dayPlan.dinner)}
        </div>
        <div class="day-macros">
            <div class="macro-item">
                <span class="macro-label" data-translate="protein">Protein:</span>
                <span class="macro-value">${dayPlan.totalProtein}g</span>
            </div>
            <div class="macro-item">
                <span class="macro-label" data-translate="carbs">Carbs:</span>
                <span class="macro-value">${dayPlan.totalCarbs}g</span>
            </div>
            <div class="macro-item">
                <span class="macro-label" data-translate="fat">Fat:</span>
                <span class="macro-value">${dayPlan.totalFat}g</span>
            </div>
        </div>
    `;
    
    return dayCard;
}

function createMealCard(mealType, meal) {
    if (!meal) return '';
    
    return `
        <div class="meal-card">
            <div class="meal-header">
                <h6 class="meal-type" data-translate="${mealType.toLowerCase().replace(' ', '_')}">${mealType}</h6>
                <span class="meal-calories">${meal.calories} kcal</span>
            </div>
            <div class="meal-name">${meal.name}</div>
            <div class="meal-ingredients">
                <strong data-translate="ingredients">Ingredients:</strong>
                <ul>
                    ${meal.ingredients.map(ingredient => `<li>${ingredient}</li>`).join('')}
                </ul>
            </div>
            <div class="meal-preparation">
                <strong data-translate="preparation">Preparation:</strong>
                <p>${meal.preparation}</p>
            </div>
            ${meal.alternatives ? `
                <div class="meal-alternatives">
                    <strong data-translate="alternatives">Alternatives:</strong>
                    <ul>
                        ${meal.alternatives.map(alt => `<li>${alt}</li>`).join('')}
                    </ul>
                </div>
            ` : ''}
            ${meal.notes ? `
                <div class="meal-notes">
                    <strong data-translate="notes">Notes:</strong>
                    <p>${meal.notes}</p>
                </div>
            ` : ''}
        </div>
    `;
}

async function downloadMealPlanPDF(userType) {
    const weekSelector = document.getElementById(`${userType}WeekSelector`);
    const selectedWeek = weekSelector.value;
    
    try {
        const userData = collectUserData(userType);
        
        const response = await fetch('/api/generate-meal-plan-pdf', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${getApiKey()}`
            },
            body: JSON.stringify({
                ...userData,
                week: selectedWeek,
                language: getCurrentLanguage()
            })
        });
        
        if (!response.ok) {
            throw new Error('Failed to generate PDF');
        }
        
        const blob = await response.blob();
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.style.display = 'none';
        a.href = url;
        a.download = `meal-plan-week-${selectedWeek}.pdf`;
        document.body.appendChild(a);
        a.click();
        window.URL.revokeObjectURL(url);
        document.body.removeChild(a);
        
    } catch (error) {
        console.error('Error downloading PDF:', error);
        alert('Failed to download PDF. Please try again.');
    }
}

function generateShoppingList(userType) {
    // Implementation for shopping list generation
    console.log('Generating shopping list for', userType);
}

function getApiKey() {
    // For now, return a placeholder since the backend doesn't require API keys for meal generation
    // In a production environment, this would retrieve a real API key
    return 'demo-api-key';
}

function getCurrentLanguage() {
    return document.documentElement.lang || 'en';
}

// Form Handlers
document.addEventListener('DOMContentLoaded', function() {
    // Healthy People Form
    const healthyForm = document.getElementById('healthyForm');
    if (healthyForm) {
        healthyForm.addEventListener('submit', function(e) {
            e.preventDefault();
            handleHealthyFormSubmit();
        });
    }
    
    // Patients Form
    const patientsForm = document.getElementById('patientsForm');
    if (patientsForm) {
        patientsForm.addEventListener('submit', function(e) {
            e.preventDefault();
            handlePatientsFormSubmit();
        });
    }
});

async function handleHealthyFormSubmit() {
    // Get form data
    const formData = {
        name: document.getElementById('healthyName').value,
        age: parseInt(document.getElementById('healthyAge').value),
        height: parseInt(document.getElementById('healthyHeight').value),
        weight: parseInt(document.getElementById('healthyWeight').value),
        gender: document.getElementById('healthyGender').value,
        activity: document.getElementById('healthyActivity').value,
        goal: document.getElementById('healthyGoal').value,
        cuisine: document.getElementById('healthyCuisine').value,
        dietary: document.getElementById('healthyDietary').value,
        allergies: document.getElementById('healthyAllergies').value
    };
    
    // Validate form
    if (!validateHealthyForm(formData)) {
        return;
    }
    
    // Calculate results
    const results = calculateHealthyResults(formData);
    
    // Display results
    displayHealthyResults(results, formData);
    
    // Show weekly plan section
    const weeklyPlanSection = document.getElementById('healthyWeeklyPlan');
    if (weeklyPlanSection) {
        weeklyPlanSection.style.display = 'block';
    }
    
    // Generate meal plan (legacy support)
    await generateHealthyMealPlan(results, formData);
    
    // Store in localStorage
    localStorage.setItem('lastHealthyPlan', JSON.stringify({
        formData,
        results,
        timestamp: new Date().toISOString()
    }));
}

async function handlePatientsFormSubmit() {
    // Get form data
    const formData = {
        name: document.getElementById('patientName').value,
        age: parseInt(document.getElementById('patientAge').value),
        height: parseInt(document.getElementById('patientHeight').value),
        weight: parseInt(document.getElementById('patientWeight').value),
        gender: document.getElementById('patientGender').value,
        activity: document.getElementById('patientActivity').value,
        disease: document.getElementById('patientDisease').value,
        medications: document.getElementById('patientMedications').value,
        labResults: document.getElementById('patientLabResults').value,
        restrictions: document.getElementById('patientRestrictions').value
    };
    
    // Validate form
    if (!validatePatientsForm(formData)) {
        return;
    }
    
    // Calculate results
    const results = calculateMedicalResults(formData);
    
    // Display results
    displayMedicalResults(results, formData);
    
    // Generate meal plan
    await generateMedicalMealPlan(results, formData);
    
    // Store in localStorage
    localStorage.setItem('lastMedicalPlan', JSON.stringify({
        formData,
        results,
        timestamp: new Date().toISOString()
    }));
}

function validateHealthyForm(data) {
    validation.clearErrors();
    
    // Validate using ValidationSystem
    const validationRules = [
        { field: 'healthyName', type: 'name', required: true, name: 'Name' },
        { field: 'healthyAge', type: 'age', required: true, name: 'Age', min: 13, max: 100 },
        { field: 'healthyHeight', type: 'height', required: true, name: 'Height' },
        { field: 'healthyWeight', type: 'weight', required: true, name: 'Weight' },
        { field: 'healthyGender', required: true, name: 'Gender' },
        { field: 'healthyActivity', required: true, name: 'Activity Level' },
        { field: 'healthyGoal', required: true, name: 'Goal' }
    ];
    
    return validation.validateForm('healthyForm', validationRules);
}

function validatePatientsForm(data) {
    validation.clearErrors();
    
    // Validate using ValidationSystem
    const validationRules = [
        { field: 'patientName', type: 'name', required: true, name: 'Name' },
        { field: 'patientAge', type: 'age', required: true, name: 'Age', min: 13, max: 100 },
        { field: 'patientHeight', type: 'height', required: true, name: 'Height' },
        { field: 'patientWeight', type: 'weight', required: true, name: 'Weight' },
        { field: 'patientGender', required: true, name: 'Gender' },
        { field: 'patientActivity', required: true, name: 'Activity Level' },
        { field: 'patientDisease', required: true, name: 'Medical Condition' }
    ];
    
    return validation.validateForm('patientsForm', validationRules);
}

function calculateHealthyResults(data) {
    const bmrMifflin = calculateBMR(data.weight, data.height, data.age, data.gender, 'mifflin');
    const bmrHarris = calculateBMR(data.weight, data.height, data.age, data.gender, 'harris');
    const tdee = calculateTDEE(bmrMifflin, data.activity);
    const macros = calculateMacros(tdee, data.goal, data.weight);
    const dietType = recommendDietType(data.goal, data.dietary);
    
    return {
        bmrMifflin,
        bmrHarris,
        tdee,
        macros,
        dietType
    };
}

function calculateMedicalResults(data) {
    const bmrMifflin = calculateBMR(data.weight, data.height, data.age, data.gender, 'mifflin');
    const bmrHarris = calculateBMR(data.weight, data.height, data.age, data.gender, 'harris');
    const tdee = calculateTDEE(bmrMifflin, data.activity);
    
    // Adjust for medical condition
    let adjustedTdee = tdee;
    if (data.disease === 'diabetes_type1' || data.disease === 'diabetes_type2') {
        adjustedTdee = tdee * 0.95; // Slightly lower for diabetes
    } else if (data.disease === 'obesity') {
        adjustedTdee = tdee - 500; // Deficit for weight loss
    }
    
    const macros = calculateMacros(adjustedTdee, 'maintain', data.weight);
    const dietType = recommendDietType('maintain', null, data.disease);
    
    return {
        bmrMifflin,
        bmrHarris,
        tdee: adjustedTdee,
        macros,
        dietType,
        medicalNotes: getMedicalNotes(data.disease)
    };
}

function getMedicalNotes(disease) {
    const notes = {
        diabetes_type1: 'Monitor blood glucose levels closely. Coordinate carbohydrate intake with insulin.',
        diabetes_type2: 'Focus on low glycemic index foods. Monitor portion sizes.',
        hypertension: 'Limit sodium intake to less than 2300mg daily. Increase potassium-rich foods.',
        heart_disease: 'Limit saturated fats. Increase omega-3 fatty acids.',
        kidney_disease: 'Monitor protein and phosphorus intake. Limit sodium.',
        liver_disease: 'Avoid alcohol completely. Monitor protein intake.',
        obesity: 'Focus on portion control and regular physical activity.',
        celiac: 'Strictly avoid all gluten-containing foods.',
        ibs: 'Follow low FODMAP diet. Identify trigger foods.'
    };
    
    return notes[disease] || 'Follow general healthy eating guidelines.';
}

function displayHealthyResults(results, formData) {
    const resultsDiv = document.getElementById('healthyResults');
    const calculationsDiv = document.getElementById('healthyCalculations');
    
    calculationsDiv.innerHTML = `
        <div class="result-card">
            <h4>Basal Metabolic Rate (BMR)</h4>
            <p><strong>Mifflin-St Jeor Formula:</strong> ${results.bmrMifflin} calories/day</p>
            <p><strong>Harris-Benedict Formula:</strong> ${results.bmrHarris} calories/day</p>
            <small>BMR is the number of calories your body needs at rest.</small>
        </div>
        
        <div class="result-card">
            <h4>Total Daily Energy Expenditure (TDEE)</h4>
            <p><strong>${results.tdee} calories/day</strong></p>
            <small>TDEE includes your activity level and is your maintenance calories.</small>
        </div>
        
        <div class="result-card">
            <h4>Recommended Daily Intake</h4>
            <div class="row">
                <div class="col-md-3">
                    <strong>Calories:</strong><br>
                    ${results.macros.calories}
                </div>
                <div class="col-md-3">
                    <strong>Protein:</strong><br>
                    ${results.macros.protein}g
                </div>
                <div class="col-md-3">
                    <strong>Carbs:</strong><br>
                    ${results.macros.carbs}g
                </div>
                <div class="col-md-3">
                    <strong>Fat:</strong><br>
                    ${results.macros.fat}g
                </div>
            </div>
        </div>
        
        <div class="result-card">
            <h4>Recommended Diet Type</h4>
            <p><strong>${results.dietType}</strong></p>
            <small>Based on your goals and preferences.</small>
        </div>
    `;
    
    resultsDiv.style.display = 'block';
}

function displayMedicalResults(results, formData) {
    const resultsDiv = document.getElementById('patientResults');
    const calculationsDiv = document.getElementById('patientCalculations');
    
    calculationsDiv.innerHTML = `
        <div class="result-card">
            <h4>Medical Nutritional Analysis</h4>
            <p><strong>Condition:</strong> ${formData.disease.replace('_', ' ').toUpperCase()}</p>
            <p><strong>Recommended Diet:</strong> ${results.dietType}</p>
        </div>
        
        <div class="result-card">
            <h4>Caloric Requirements</h4>
            <p><strong>BMR (Mifflin-St Jeor):</strong> ${results.bmrMifflin} calories/day</p>
            <p><strong>Adjusted TDEE:</strong> ${results.tdee} calories/day</p>
        </div>
        
        <div class="result-card">
            <h4>Macronutrient Targets</h4>
            <div class="row">
                <div class="col-md-3">
                    <strong>Calories:</strong><br>
                    ${results.macros.calories}
                </div>
                <div class="col-md-3">
                    <strong>Protein:</strong><br>
                    ${results.macros.protein}g
                </div>
                <div class="col-md-3">
                    <strong>Carbs:</strong><br>
                    ${results.macros.carbs}g
                </div>
                <div class="col-md-3">
                    <strong>Fat:</strong><br>
                    ${results.macros.fat}g
                </div>
            </div>
        </div>
        
        <div class="result-card">
            <h4>Medical Notes</h4>
            <p>${results.medicalNotes}</p>
            <small><em>Always consult with your healthcare provider before making dietary changes.</em></small>
        </div>
    `;
    
    resultsDiv.style.display = 'block';
}

async function generateHealthyMealPlan(results, formData) {
    const mealPlanDiv = document.getElementById('healthyMealPlan');
    const contentDiv = document.getElementById('healthyMealPlanContent');
    
    let mealPlanHTML = '';
    
    for (let day = 1; day <= 7; day++) {
        const dayMeals = await generateDayMeals(day, results.macros, formData);
        mealPlanHTML += `
            <div class="day-card">
                <h4>Day ${day}</h4>
                <div class="row">
                    ${dayMeals}
                </div>
            </div>
        `;
    }
    
    contentDiv.innerHTML = mealPlanHTML;
    mealPlanDiv.style.display = 'block';
}

async function generateMedicalMealPlan(results, formData) {
    const mealPlanDiv = document.getElementById('patientMealPlan');
    const contentDiv = document.getElementById('patientMealPlanContent');
    
    let mealPlanHTML = '';
    
    for (let day = 1; day <= 7; day++) {
        const dayMeals = await generateMedicalDayMeals(day, results.macros, formData);
        mealPlanHTML += `
            <div class="day-card">
                <h4>Day ${day} - Medical Plan</h4>
                <div class="row">
                    ${dayMeals}
                </div>
            </div>
        `;
    }
    
    contentDiv.innerHTML = mealPlanHTML;
    mealPlanDiv.style.display = 'block';
}

async function generateDayMeals(day, macros, formData) {
    const meals = ['Breakfast', 'Snack 1', 'Lunch', 'Snack 2', 'Dinner'];
    const caloriesPerMeal = [0.25, 0.1, 0.3, 0.1, 0.25];
    
    let mealsHTML = '';
    
    for (let index = 0; index < meals.length; index++) {
        const meal = meals[index];
        const mealCalories = Math.round(macros.calories * caloriesPerMeal[index]);
        const sample = await getSampleMeal(meal, mealCalories, formData);
        mealsHTML += `
            <div class="col-md-4 mb-3">
                <div class="meal-item">
                    <h5>${meal}</h5>
                    <p><strong>Sample Meal:</strong> ${sample}</p>
                    <p><strong>Calories:</strong> ~${mealCalories}</p>
                    <small>Adjust portions based on your needs</small>
                </div>
            </div>
        `;
    }
    
    return mealsHTML;
}

async function generateMedicalDayMeals(day, macros, formData) {
    const meals = ['Breakfast', 'Snack 1', 'Lunch', 'Snack 2', 'Dinner'];
    const caloriesPerMeal = [0.25, 0.1, 0.3, 0.1, 0.25];
    
    let mealsHTML = '';
    
    for (let index = 0; index < meals.length; index++) {
        const meal = meals[index];
        const mealCalories = Math.round(macros.calories * caloriesPerMeal[index]);
        const sample = await getMedicalSampleMeal(meal, mealCalories, formData);
        mealsHTML += `
            <div class="col-md-4 mb-3">
                <div class="meal-item">
                    <h5>${meal}</h5>
                    <p><strong>Medical Meal:</strong> ${sample}</p>
                    <p><strong>Calories:</strong> ~${mealCalories}</p>
                    <small><em>Designed for ${formData.disease.replace('_', ' ')}</em></small>
                </div>
            </div>
        `;
    }
    
    return mealsHTML;
}

async function getSampleMeal(mealType, calories, formData) {
    // Fetch sample meal from backend API
    try {
        const response = await app.apiCall('/recipes?meal_type=' + mealType, 'GET');
        if (response.ok) {
            const recipes = await response.json();
            if (recipes.recipes && recipes.recipes.length > 0) {
                const randomIndex = Math.floor(Math.random() * recipes.recipes.length);
                return recipes.recipes[randomIndex].name;
            } else {
                return 'No recipe found for ' + mealType;
            }
        } else {
            return 'Error fetching recipes';
        }
    } catch (e) {
        console.error(e);
        return 'Network error while fetching recipes';
    }
}

async function getMedicalSampleMeal(mealType, calories, formData) {
    // Fetch medical sample meal from backend API, filtered by disease
    try {
        const response = await app.apiCall('/recipes?meal_type=' + mealType + '&medical_conditions=' + formData.disease, 'GET');
        if (response.ok) {
            const recipes = await response.json();
            if (recipes.recipes && recipes.recipes.length > 0) {
                const randomIndex = Math.floor(Math.random() * recipes.recipes.length);
                return recipes.recipes[randomIndex].name;
            } else {
                return 'No medical recipe found for ' + mealType;
            }
        } else {
            return 'Error fetching medical recipes';
        }
    } catch (e) {
        console.error(e);
        return 'Network error';
    }
}

// Utility Functions
function generateShoppingList(type = 'healthy') {
    const ingredients = [
        'Oats', 'Greek yogurt', 'Mixed berries', 'Almonds', 'Chicken breast',
        'Salmon', 'Quinoa', 'Brown rice', 'Mixed greens', 'Broccoli',
        'Sweet potato', 'Avocado', 'Olive oil', 'Eggs', 'Lentils'
    ];
    
    const shoppingList = ingredients.join('\n');
    
    // Create and download shopping list
    const blob = new Blob([shoppingList], { type: 'text/plain' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `shopping-list-${type}.txt`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
}

function downloadPDF(type) {
    // Create PDF content
    const content = document.getElementById(type + 'Results').innerHTML;
    const mealPlan = document.getElementById(type + 'MealPlanContent').innerHTML;
    
    const pdfContent = `
        <html>
        <head>
            <title>Nutrition Plan - ${type}</title>
            <style>
                body { font-family: Arial, sans-serif; margin: 20px; }
                .results-section { margin-bottom: 30px; }
                .meal-plan { margin-top: 30px; }
                .day-card { margin-bottom: 20px; border: 1px solid #ddd; padding: 15px; }
                .meal-item { margin-bottom: 10px; padding: 10px; background: #f9f9f9; }
            </style>
        </head>
        <body>
            <h1>Nutrition Plan</h1>
            <div class="results-section">${content}</div>
            <div class="meal-plan">${mealPlan}</div>
        </body>
        </html>
    `;
    
    // Create and download PDF
    const blob = new Blob([pdfContent], { type: 'text/html' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    const filename = `nutrition-plan-${type}-${new Date().toISOString().split('T')[0]}.html`;
    a.href = url;
    a.download = filename;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
}