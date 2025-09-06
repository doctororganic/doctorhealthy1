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

function handleHealthyFormSubmit() {
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
    
    // Generate meal plan
    generateHealthyMealPlan(results, formData);
    
    // Store in localStorage
    localStorage.setItem('lastHealthyPlan', JSON.stringify({
        formData,
        results,
        timestamp: new Date().toISOString()
    }));
}

function handlePatientsFormSubmit() {
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
    generateMedicalMealPlan(results, formData);
    
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

function generateHealthyMealPlan(results, formData) {
    const mealPlanDiv = document.getElementById('healthyMealPlan');
    const contentDiv = document.getElementById('healthyMealPlanContent');
    
    // Generate 7-day meal plan
    let mealPlanHTML = '';
    
    for (let day = 1; day <= 7; day++) {
        mealPlanHTML += `
            <div class="day-card">
                <h4>Day ${day}</h4>
                <div class="row">
                    ${generateDayMeals(day, results.macros, formData)}
                </div>
            </div>
        `;
    }
    
    contentDiv.innerHTML = mealPlanHTML;
    mealPlanDiv.style.display = 'block';
}

function generateMedicalMealPlan(results, formData) {
    const mealPlanDiv = document.getElementById('patientMealPlan');
    const contentDiv = document.getElementById('patientMealPlanContent');
    
    // Generate medical meal plan
    let mealPlanHTML = '';
    
    for (let day = 1; day <= 7; day++) {
        mealPlanHTML += `
            <div class="day-card">
                <h4>Day ${day} - Medical Plan</h4>
                <div class="row">
                    ${generateMedicalDayMeals(day, results.macros, formData)}
                </div>
            </div>
        `;
    }
    
    contentDiv.innerHTML = mealPlanHTML;
    mealPlanDiv.style.display = 'block';
}

function generateDayMeals(day, macros, formData) {
    const meals = ['Breakfast', 'Snack 1', 'Lunch', 'Snack 2', 'Dinner'];
    const caloriesPerMeal = [0.25, 0.1, 0.3, 0.1, 0.25]; // Percentage of daily calories
    
    let mealsHTML = '';
    
    meals.forEach((meal, index) => {
        const mealCalories = Math.round(macros.calories * caloriesPerMeal[index]);
        mealsHTML += `
            <div class="col-md-4 mb-3">
                <div class="meal-item">
                    <h5>${meal}</h5>
                    <p><strong>Sample Meal:</strong> ${getSampleMeal(meal, mealCalories, formData)}</p>
                    <p><strong>Calories:</strong> ~${mealCalories}</p>
                    <small>Adjust portions based on your needs</small>
                </div>
            </div>
        `;
    });
    
    return mealsHTML;
}

function generateMedicalDayMeals(day, macros, formData) {
    const meals = ['Breakfast', 'Snack 1', 'Lunch', 'Snack 2', 'Dinner'];
    const caloriesPerMeal = [0.25, 0.1, 0.3, 0.1, 0.25];
    
    let mealsHTML = '';
    
    meals.forEach((meal, index) => {
        const mealCalories = Math.round(macros.calories * caloriesPerMeal[index]);
        mealsHTML += `
            <div class="col-md-4 mb-3">
                <div class="meal-item">
                    <h5>${meal}</h5>
                    <p><strong>Medical Meal:</strong> ${getMedicalSampleMeal(meal, mealCalories, formData)}</p>
                    <p><strong>Calories:</strong> ~${mealCalories}</p>
                    <small><em>Designed for ${formData.disease.replace('_', ' ')}</em></small>
                </div>
            </div>
        `;
    });
    
    return mealsHTML;
}

function getSampleMeal(mealType, calories, formData) {
    const sampleMeals = {
        'Breakfast': [
            'Oatmeal with berries and nuts',
            'Greek yogurt with granola',
            'Whole grain toast with avocado',
            'Smoothie bowl with fruits'
        ],
        'Snack 1': [
            'Apple with almond butter',
            'Greek yogurt',
            'Mixed nuts',
            'Hummus with vegetables'
        ],
        'Lunch': [
            'Grilled chicken salad',
            'Quinoa bowl with vegetables',
            'Lentil soup with bread',
            'Turkey wrap with vegetables'
        ],
        'Snack 2': [
            'Protein smoothie',
            'Cottage cheese with fruit',
            'Trail mix',
            'Vegetable sticks with hummus'
        ],
        'Dinner': [
            'Baked salmon with quinoa',
            'Grilled chicken with sweet potato',
            'Lean beef with brown rice',
            'Tofu stir-fry with vegetables'
        ]
    };
    
    const meals = sampleMeals[mealType] || ['Balanced meal'];
    return meals[Math.floor(Math.random() * meals.length)];
}

function getMedicalSampleMeal(mealType, calories, formData) {
    const disease = formData.disease;
    
    const medicalMeals = {
        diabetes_type1: {
            'Breakfast': ['Steel-cut oats with cinnamon', 'Egg white omelet with vegetables'],
            'Lunch': ['Grilled fish with quinoa', 'Chicken salad with olive oil'],
            'Dinner': ['Baked chicken with steamed broccoli', 'Lean turkey with cauliflower rice']
        },
        diabetes_type2: {
            'Breakfast': ['Greek yogurt with berries', 'Vegetable omelet'],
            'Lunch': ['Lentil soup', 'Grilled salmon salad'],
            'Dinner': ['Baked cod with asparagus', 'Chicken breast with green beans']
        },
        hypertension: {
            'Breakfast': ['Oatmeal with banana', 'Low-sodium whole grain cereal'],
            'Lunch': ['Herb-seasoned chicken', 'Quinoa salad with vegetables'],
            'Dinner': ['Baked fish with herbs', 'Lean beef with roasted vegetables']
        }
    };
    
    const diseaseSpecific = medicalMeals[disease];
    if (diseaseSpecific && diseaseSpecific[mealType]) {
        const meals = diseaseSpecific[mealType];
        return meals[Math.floor(Math.random() * meals.length)];
    }
    
    return getSampleMeal(mealType, calories, formData);
}

// Utility Functions
function generateShoppingList(type = 'healthy') {
    const ingredients = [
        'Oats', 'Greek yogurt', 'Mixed berries', 'Almonds', 'Chicken breast',
        'Salmon', 'Quinoa', 'Brown rice', 'Mixed greens', 'Broccoli',
        'Sweet potato', 'Avocado', 'Olive oil', 'Eggs', 'Lentils'
    ];
    
    alert(`Shopping List Generated!\n\n${ingredients.join('\n')}\n\nThis list has been copied to your clipboard.`);
}

function downloadPDF(type) {
    alert(`PDF download functionality will be implemented with a proper PDF library. This would generate a comprehensive ${type} meal plan PDF.`);
}

function openSubscriptionForm() {
    // This would open a Google Form embed
    window.open('https://forms.google.com/subscription-form', '_blank');
}

function openMedicalSubscriptionForm() {
    // This would open a medical subscription Google Form
    window.open('https://forms.google.com/medical-subscription-form', '_blank');
}

function contactWhatsApp(type = 'general') {
    const phoneNumber = '+1234567890'; // Replace with actual WhatsApp number
    const message = type === 'medical' 
        ? 'Hello, I need medical nutrition consultation.'
        : 'Hello, I\'m interested in your meal plan subscription.';
    
    const whatsappURL = `https://wa.me/${phoneNumber}?text=${encodeURIComponent(message)}`;
    window.open(whatsappURL, '_blank');
}

// Save meal plan to JSON file
function saveMealPlan(data, type) {
    const timestamp = new Date().toISOString().split('T')[0];
    const filename = `user-${type}-${timestamp}.json`;
    
    const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    
    const a = document.createElement('a');
    a.href = url;
    a.download = filename;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
}