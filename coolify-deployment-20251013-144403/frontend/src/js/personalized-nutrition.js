// Personalized Nutrition System JavaScript

// Global variables
let currentClientData = {};
let selectedCountry = null;
let generatedPlan = null;

// Initialize the system
document.addEventListener('DOMContentLoaded', function() {
    initializeSystem();
    setupEventListeners();
    loadCountries();
});

function initializeSystem() {
    console.log('Personalized Nutrition System initialized');
}

function setupEventListeners() {
    // Form submission
    document.getElementById('clientDataForm').addEventListener('submit', handleFormSubmission);
    
    // Real-time BMI calculation
    ['clientWeight', 'clientHeight'].forEach(id => {
        document.getElementById(id).addEventListener('input', calculateBMIRealTime);
    });
}

// Handle form submission
function handleFormSubmission(event) {
    event.preventDefault();
    
    // Collect client data
    currentClientData = collectClientData();
    
    // Validate data
    if (!validateClientData(currentClientData)) {
        return;
    }
    
    // Calculate nutrition requirements
    const nutritionCalculations = calculateNutritionRequirements(currentClientData);
    
    // Display calculations
    displayCalculationResults(nutritionCalculations);
    
    // Generate meal plan
    generateMealPlan(currentClientData, nutritionCalculations);
    
    // Show results
    document.getElementById('calculationResults').style.display = 'block';
    document.getElementById('weeklyPlan').style.display = 'block';
    
    // Scroll to results
    document.getElementById('calculationResults').scrollIntoView({ behavior: 'smooth' });
}

// Collect client data from form
function collectClientData() {
    return {
        name: document.getElementById('clientName').value,
        age: parseInt(document.getElementById('clientAge').value),
        weight: parseFloat(document.getElementById('clientWeight').value),
        height: parseInt(document.getElementById('clientHeight').value),
        gender: document.getElementById('clientGender').value,
        activityLevel: document.getElementById('activityLevel').value,
        metabolicRate: document.getElementById('metabolicRate').value,
        mainGoal: document.getElementById('mainGoal').value,
        foodRestrictions: document.getElementById('foodRestrictions').value,
        medicalConditions: document.getElementById('medicalConditions').value,
        medications: document.getElementById('medications').value
    };
}

// Enhanced validation using error handler
function validateClientData(data) {
    // Use enhanced error handler validation if available
    if (window.errorHandler && typeof window.errorHandler.validateClientData === 'function') {
        const errors = window.errorHandler.validateClientData({
            name: data.name,
            age: data.age,
            weight: data.weight,
            height: data.height,
            gender: data.gender,
            activityLevel: data.activityLevel,
            goal: data.mainGoal
        });
        
        if (errors.length > 0) {
            // Show validation errors using enhanced error handler
            if (typeof window.errorHandler.showValidationErrors === 'function') {
                window.errorHandler.showValidationErrors(errors, 'nutritionValidationErrors');
            } else {
                // Fallback to alert
                alert('يرجى تصحيح الأخطاء التالية:\n' + errors.join('\n'));
            }
            return false;
        }
        
        // Clear any existing validation errors
        const errorContainer = document.getElementById('nutritionValidationErrors');
        if (errorContainer) {
            errorContainer.style.display = 'none';
        }
        
        return true;
    }
    
    // Fallback validation if error handler is not available
    const requiredFields = ['name', 'age', 'weight', 'height', 'gender', 'activityLevel', 'metabolicRate', 'mainGoal'];
    
    for (let field of requiredFields) {
        if (!data[field]) {
            alert(`الرجاء ملء جميع الحقول المطلوبة: ${field}`);
            return false;
        }
    }
    
    // Validate ranges
    if (data.age < 13 || data.age > 100) {
        alert('العمر يجب أن يكون بين 13 و 100 سنة');
        return false;
    }
    
    if (data.weight < 30 || data.weight > 300) {
        alert('الوزن يجب أن يكون بين 30 و 300 كيلو');
        return false;
    }
    
    if (data.height < 100 || data.height > 250) {
        alert('الطول يجب أن يكون بين 100 و 250 سم');
        return false;
    }
    
    return true;
}

// Calculate BMI in real-time
function calculateBMIRealTime() {
    const weight = parseFloat(document.getElementById('clientWeight').value);
    const height = parseInt(document.getElementById('clientHeight').value);
    
    if (weight && height) {
        const bmi = weight / Math.pow(height / 100, 2);
        
        // You can display BMI somewhere if needed
        console.log('Current BMI:', bmi.toFixed(1));
    }
}

// Calculate nutrition requirements based on client data
function calculateNutritionRequirements(data) {
    const bmi = data.weight / Math.pow(data.height / 100, 2);
    
    // Calculate BMR using Mifflin-St Jeor equation
    let bmr;
    if (data.gender === 'male') {
        bmr = (10 * data.weight) + (6.25 * data.height) - (5 * data.age) + 5;
    } else {
        bmr = (10 * data.weight) + (6.25 * data.height) - (5 * data.age) - 161;
    }
    
    // Activity multipliers
    const activityMultipliers = {
        'sedentary': 1.2,
        'light': 1.375,
        'moderate': 1.55,
        'very': 1.725,
        'extra': 1.9
    };
    
    let tdee = bmr * (activityMultipliers[data.activityLevel] || 1.55);
    
    // Determine calories per kg based on BMI and goals
    let caloriesPerKg;
    
    if (bmi >= 18 && bmi <= 30) {
        // Normal to overweight range
        if (data.mainGoal === 'weight_loss' || data.mainGoal === 'maintain_weight') {
            caloriesPerKg = 20;
        } else {
            caloriesPerKg = 25;
        }
    } else if (bmi >= 15 && bmi < 18) {
        // Underweight range
        if (data.mainGoal === 'weight_gain' || data.mainGoal === 'maintain_weight') {
            caloriesPerKg = 25;
        } else {
            caloriesPerKg = 20;
        }
    } else {
        // Very high BMI or muscle building goals
        caloriesPerKg = 30;
    }
    
    // Adjust for metabolic rate
    if (data.metabolicRate === 'high' || data.mainGoal === 'muscle_strength') {
        caloriesPerKg = 30;
    }
    
    // Calculate base calories using the formula
    let baseCalories = data.weight * caloriesPerKg;
    
    // Adjust based on goal
    let goalMultiplier = 1;
    switch (data.mainGoal) {
        case 'weight_loss':
            goalMultiplier = 0.85;
            break;
        case 'weight_gain':
        case 'muscle_strength':
            goalMultiplier = 1.15;
            caloriesPerKg = Math.max(caloriesPerKg, 25);
            baseCalories = data.weight * caloriesPerKg;
            break;
        case 'body_recomposition':
            goalMultiplier = 1.05;
            break;
    }
    
    const totalCalories = Math.round(baseCalories * goalMultiplier);
    
    // Calculate protein requirements based on activity and goals
    let proteinPerKg;
    if (data.activityLevel === 'very' || data.activityLevel === 'extra' || data.mainGoal === 'muscle_strength') {
        proteinPerKg = 1.7; // High activity or muscle building (1.5-1.7g per kg)
    } else if (data.activityLevel === 'moderate') {
        proteinPerKg = 1.5;
    } else {
        proteinPerKg = 1.2; // Low activity (1-1.5g per kg)
    }
    
    const proteinGrams = Math.round(data.weight * proteinPerKg);
    const proteinCalories = proteinGrams * 4;
    
    // Calculate fat (25-35% of total calories)
    const fatPercentage = 0.30;
    const fatCalories = Math.round(totalCalories * fatPercentage);
    const fatGrams = Math.round(fatCalories / 9);
    
    // Calculate carbs (remaining calories)
    const carbCalories = totalCalories - proteinCalories - fatCalories;
    const carbGrams = Math.round(carbCalories / 4);
    
    // Determine meal distribution based on goals
    let mealDistribution;
    if (data.mainGoal === 'weight_gain' || data.mainGoal === 'muscle_strength') {
        // 4 meals + 2 non-consecutive intermittent fasting days
        mealDistribution = {
            type: 'weight_gain',
            regularDays: {
                meals: 4,
                snacks: 2,
                distribution: [0.25, 0.10, 0.30, 0.10, 0.20, 0.05] // breakfast, snack1, lunch, snack2, dinner, snack3
            },
            fastingDays: {
                meals: 2,
                snacks: 1,
                distribution: [0.40, 0.15, 0.45] // meal1, snack, meal2
            }
        };
    } else {
        // Standard 3 meals + snacks
        mealDistribution = {
            type: 'standard',
            regularDays: {
                meals: 3,
                snacks: 2,
                distribution: [0.25, 0.10, 0.35, 0.15, 0.15] // breakfast, snack1, lunch, snack2, dinner
            }
        };
    }
    
    return {
        bmi: bmi,
        totalCalories: totalCalories,
        protein: { grams: proteinGrams, calories: proteinCalories },
        fat: { grams: fatGrams, calories: fatCalories },
        carbs: { grams: carbGrams, calories: carbCalories },
        caloriesPerKg: caloriesPerKg,
        mealDistribution: mealDistribution,
        recommendedDietType: recommendDietType(data, bmi),
        bmr: Math.round(bmr),
        tdee: Math.round(tdee)
    };
}

// Recommend diet type based on client data
function recommendDietType(data, bmi) {
    // Check medical conditions first
    const conditions = data.medicalConditions.toLowerCase();
    
    if (conditions.includes('diabetes') || conditions.includes('سكري')) {
        return 'low-carb';
    }
    if (conditions.includes('hypertension') || conditions.includes('ضغط')) {
        return 'dash';
    }
    if (conditions.includes('heart') || conditions.includes('قلب')) {
        return 'mediterranean';
    }
    if (conditions.includes('kidney') || conditions.includes('كلى')) {
        return 'low-protein';
    }
    
    // Based on goals
    switch (data.mainGoal) {
        case 'weight_loss':
            return bmi > 25 ? 'low-carb' : 'balanced';
        case 'weight_gain':
        case 'muscle_strength':
            return 'high-carb';
        case 'body_recomposition':
            return 'balanced';
        default:
            return 'mediterranean';
    }
}

// Display calculation results
function displayCalculationResults(calculations) {
    const resultsDiv = document.getElementById('calculationDetails');
    
    const bmiStatus = getBMIStatus(calculations.bmi);
    
    resultsDiv.innerHTML = `
        <div class="row">
            <div class="col-md-6">
                <div class="card border-primary mb-3">
                    <div class="card-header bg-primary text-white">
                        <h5 class="mb-0"><i class="fas fa-calculator me-2"></i>المؤشرات الأساسية</h5>
                    </div>
                    <div class="card-body">
                        <p><strong>مؤشر كتلة الجسم (BMI):</strong> ${calculations.bmi.toFixed(1)} - ${bmiStatus}</p>
                        <p><strong>السعرات المستخدمة:</strong> ${calculations.caloriesPerKg} سعرة لكل كيلو</p>
                        <p><strong>نوع النظام المناسب:</strong> ${getDietTypeArabic(calculations.recommendedDietType)}</p>
                        <div class="alert alert-info">
                            <small><i class="fas fa-info-circle me-1"></i>تم استخدام معادلات علمية لحساب احتياجاتك</small>
                        </div>
                    </div>
                </div>
            </div>
            <div class="col-md-6">
                <div class="card border-success mb-3">
                    <div class="card-header bg-success text-white">
                        <h5 class="mb-0"><i class="fas fa-utensils me-2"></i>الاحتياجات اليومية</h5>
                    </div>
                    <div class="card-body">
                        <div class="nutrition-summary">
                            <div class="d-flex justify-content-between mb-2">
                                <span><strong>السعرات الحرارية:</strong></span>
                                <span class="text-primary fw-bold">${calculations.totalCalories} سعرة</span>
                            </div>
                            <div class="d-flex justify-content-between mb-2">
                                <span><strong>البروتين:</strong></span>
                                <span class="text-success fw-bold">${calculations.protein.grams}g</span>
                            </div>
                            <div class="d-flex justify-content-between mb-2">
                                <span><strong>الكربوهيدرات:</strong></span>
                                <span class="text-warning fw-bold">${calculations.carbs.grams}g</span>
                            </div>
                            <div class="d-flex justify-content-between">
                                <span><strong>الدهون:</strong></span>
                                <span class="text-info fw-bold">${calculations.fat.grams}g</span>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        
        ${generateMedicalAdvice(currentClientData)}
        ${generateMedicalDisclaimer()}
    `;
}

// Get BMI status in Arabic
function getBMIStatus(bmi) {
    if (bmi < 18.5) return 'نقص في الوزن';
    if (bmi < 25) return 'وزن طبيعي';
    if (bmi < 30) return 'زيادة في الوزن';
    return 'سمنة';
}

// Get diet type in Arabic
function getDietTypeArabic(type) {
    const types = {
        'low-carb': 'قليل الكربوهيدرات',
        'high-carb': 'عالي الكربوهيدرات',
        'balanced': 'متوازن',
        'mediterranean': 'البحر الأبيض المتوسط',
        'dash': 'داش (DASH)',
        'keto': 'كيتو',
        'low-protein': 'قليل البروتين'
    };
    return types[type] || 'متوازن';
}

// Generate medical disclaimer
function generateMedicalDisclaimer() {
    return `
        <div class="col-12 mt-3">
            <div class="card border-danger">
                <div class="card-header bg-danger text-white">
                    <h6 class="mb-0"><i class="fas fa-exclamation-triangle me-2"></i>إخلاء مسؤولية طبية</h6>
                </div>
                <div class="card-body">
                    <div class="alert alert-warning mb-0">
                        <strong>تنبيه مهم:</strong> هذا النظام الغذائي مخصص للأشخاص الأصحاء فقط. إذا كنت تعاني من أي حالة طبية أو تتناول أدوية، يجب استشارة طبيب مختص قبل اتباع أي نظام غذائي جديد. هذه المعلومات للإرشاد العام فقط ولا تغني عن الاستشارة الطبية المتخصصة.
                    </div>
                </div>
            </div>
        </div>
    `;
}

// Generate medical advice based on conditions
function generateMedicalAdvice(data) {
    if (!data.medicalConditions && !data.medications) {
        return '';
    }
    
    let advice = `
        <div class="col-12">
            <div class="card border-warning">
                <div class="card-header bg-warning text-dark">
                    <h5 class="mb-0"><i class="fas fa-exclamation-triangle me-2"></i>نصائح طبية مهمة</h5>
                </div>
                <div class="card-body">
    `;
    
    if (data.medicalConditions) {
        const conditions = data.medicalConditions.toLowerCase();
        
        if (conditions.includes('diabetes') || conditions.includes('سكري')) {
            advice += `
                <div class="alert alert-info">
                    <strong>نصائح لمرضى السكري:</strong>
                    <ul class="mb-0 mt-2">
                        <li>راقب مستوى السكر في الدم قبل وبعد الوجبات</li>
                        <li>تناول الكربوهيدرات بكميات ثابتة في كل وجبة</li>
                        <li>اختر الأطعمة ذات المؤشر الجلايسيمي المنخفض</li>
                        <li>تجنب السكريات المكررة والمشروبات الغازية</li>
                        <li>تناول البروتينات الخالية من الدهون والدهون الصحية</li>
                        <li>نسق تناول الكربوهيدرات مع توقيت الأنسولين</li>
                    </ul>
                </div>
            `;
        }
        
        if (conditions.includes('hypertension') || conditions.includes('ضغط')) {
            advice += `
                <div class="alert alert-info">
                    <strong>نصائح لمرضى ضغط الدم:</strong>
                    <ul class="mb-0 mt-2">
                        <li>قلل من تناول الصوديوم إلى أقل من 2300 ملغ يومياً</li>
                        <li>أكثر من تناول الأطعمة الغنية بالبوتاسيوم</li>
                        <li>اختر منتجات الألبان قليلة الدسم</li>
                        <li>تجنب الأطعمة المصنعة عالية الصوديوم</li>
                        <li>قلل من استهلاك الكحول</li>
                        <li>اتبع نظام DASH الغذائي</li>
                    </ul>
                </div>
            `;
        }
        
        if (conditions.includes('heart') || conditions.includes('قلب')) {
            advice += `
                <div class="alert alert-info">
                    <strong>نصائح لمرضى القلب:</strong>
                    <ul class="mb-0 mt-2">
                        <li>قلل من الدهون المشبعة إلى أقل من 7% من السعرات</li>
                        <li>أكثر من تناول أحماض أوميغا-3 الدهنية</li>
                        <li>تجنب الدهون المتحولة تماماً</li>
                        <li>اختر البروتينات الخالية من الدهون</li>
                        <li>أكثر من تناول الفواكه والخضروات</li>
                        <li>اتبع النظام الغذائي المتوسطي</li>
                    </ul>
                </div>
            `;
        }
        
        if (conditions.includes('kidney') || conditions.includes('كلى')) {
            advice += `
                <div class="alert alert-info">
                    <strong>نصائح لمرضى الكلى:</strong>
                    <ul class="mb-0 mt-2">
                        <li>راقب تناول البروتين حسب مرحلة المرض</li>
                        <li>قلل من الفوسفور والبوتاسيوم</li>
                        <li>تحكم في تناول الصوديوم</li>
                        <li>حافظ على السعرات الحرارية الكافية</li>
                        <li>راقب تناول السوائل إذا لزم الأمر</li>
                    </ul>
                </div>
            `;
        }
        
        if (conditions.includes('liver') || conditions.includes('كبد')) {
            advice += `
                <div class="alert alert-info">
                    <strong>نصائح لمرضى الكبد:</strong>
                    <ul class="mb-0 mt-2">
                        <li>تجنب الكحول تماماً</li>
                        <li>قلل من تناول الصوديوم</li>
                        <li>راقب تناول البروتين حسب الحالة</li>
                        <li>حافظ على وزن صحي</li>
                        <li>تجنب المواد السامة للكبد</li>
                        <li>تجنب المحار النيء أو غير المطبوخ جيداً</li>
                    </ul>
                </div>
            `;
        }
        
        if (conditions.includes('celiac') || conditions.includes('جلوتين')) {
            advice += `
                <div class="alert alert-info">
                    <strong>نصائح لمرضى السيلياك:</strong>
                    <ul class="mb-0 mt-2">
                        <li>تجنب جميع الحبوب المحتوية على الجلوتين</li>
                        <li>اقرأ الملصقات بعناية</li>
                        <li>امنع التلوث المتقاطع</li>
                        <li>اختر الأطعمة الخالية من الجلوتين طبيعياً</li>
                        <li>تأكد من الحصول على الألياف وفيتامينات ب الكافية</li>
                    </ul>
                </div>
            `;
        }
    }
    
    if (data.medications) {
        advice += `
            <div class="alert alert-warning">
                <strong><i class="fas fa-pills me-1"></i>معلومات عن الأدوية:</strong>
                <p class="mb-2">بعض الأدوية قد تتفاعل مع الطعام. استشر طبيبك أو الصيدلي حول:</p>
                <ul class="mb-0">
                    <li>الأطعمة التي يجب تجنبها مع أدويتك</li>
                    <li>أفضل أوقات تناول الدواء بالنسبة للوجبات</li>
                    <li>المكملات الغذائية التي قد تتفاعل مع أدويتك</li>
                </ul>
            </div>
        `;
    }
    
    advice += `
                    <div class="alert alert-danger">
                        <strong><i class="fas fa-exclamation-triangle me-1"></i>تنويه مهم:</strong>
                        هذه التوصيات عامة ولا تغني عن استشارة طبيب مختص. يجب مراجعة طبيبك قبل اتباع أي نظام غذائي جديد، خاصة إذا كنت تعاني من حالات طبية أو تتناول أدوية.
                    </div>
                </div>
            </div>
        </div>
    `;
    
    return advice;
}

// Generate meal plan
function generateMealPlan(clientData, calculations) {
    // Fetch recipes from API
    fetch('/api/recipes')
        .then(response => response.json())
        .then(recipes => {
            // Filter recipes based on restrictions
            const filteredRecipes = filterRecipesByRestrictions(recipes, clientData);
            
            // Generate weekly meal plan using the new system
            generateWeeklyMealPlan(clientData, filteredRecipes);
        })
        .catch(error => {
            console.error('Error fetching recipes:', error);
            // Fallback to old system if API fails
            generateFallbackMealPlan(clientData, calculations);
        });
}

// Fallback meal plan generation (old system)
function generateFallbackMealPlan(clientData, calculations) {
    const mealPlanGrid = document.getElementById('mealPlanGrid');
    
    // Generate 7-day meal plan
    const weekDays = ['الأحد', 'الاثنين', 'الثلاثاء', 'الأربعاء', 'الخميس', 'الجمعة', 'السبت'];
    const mealTypes = ['الإفطار', 'سناك صباحي', 'الغداء', 'سناك مسائي', 'العشاء'];
    
    let planHTML = '';
    
    weekDays.forEach((day, dayIndex) => {
        planHTML += `
            <div class="col-12 mb-4">
                <h4 class="text-primary mb-3">
                    <i class="fas fa-calendar-day me-2"></i>${day}
                </h4>
                <div class="row">
        `;
        
        mealTypes.forEach((mealType, mealIndex) => {
            const mealCalories = Math.round(calculations.totalCalories * calculations.mealDistribution.regularDays.distribution[mealIndex]);
            const meal = generateMealForType(mealType, mealCalories, clientData, calculations.recommendedDietType);
            
            planHTML += `
                <div class="col-md-6 col-lg-4 mb-3">
                    <div class="meal-box" onclick="showMealDetails('${day}', '${mealType}', ${mealIndex})">
                        <h5>${mealType}</h5>
                        <p class="meal-name">${meal.name}</p>
                        <p class="meal-ingredients">${meal.ingredients}</p>
                        
                        <div class="nutrition-info">
                            <div class="nutrition-item">
                                <div class="value">${mealCalories}</div>
                                <div class="label">سعرة</div>
                            </div>
                            <div class="nutrition-item">
                                <div class="value">${meal.protein}g</div>
                                <div class="label">بروتين</div>
                            </div>
                            <div class="nutrition-item">
                                <div class="value">${meal.carbs}g</div>
                                <div class="label">كارب</div>
                            </div>
                            <div class="nutrition-item">
                                <div class="value">${meal.fat}g</div>
                                <div class="label">دهون</div>
                            </div>
                        </div>
                        
                        <div class="mt-3">
                            <button class="btn btn-sm btn-outline-primary" onclick="event.stopPropagation(); showAlternative('${mealType}', ${mealCalories})">
                                <i class="fas fa-exchange-alt me-1"></i>بديل
                            </button>
                        </div>
                    </div>
                </div>
            `;
        });
        
        planHTML += `
                </div>
            </div>
        `;
    });
    
    mealPlanGrid.innerHTML = planHTML;
}

// Generate weekly meal plan with detailed meal boxes
function generateWeeklyMealPlan(clientData, recipes) {
    const mealPlanContainer = document.getElementById('mealPlanContainer');
    if (!mealPlanContainer) return;
    
    const dailyCalories = clientData.nutritionRequirements.calories;
    const dailyProtein = clientData.nutritionRequirements.protein;
    const dailyCarbs = clientData.nutritionRequirements.carbs;
    const dailyFat = clientData.nutritionRequirements.fat;
    
    // Determine meal count based on goal
    const mealsPerDay = getMealsPerDay(clientData.goal);
    const isIntermittentFasting = clientData.goal === 'muscle_building' || clientData.goal === 'weight_gain';
    
    let weeklyPlan = '<div class="weekly-meal-plan">';
    weeklyPlan += '<h4 class="text-center mb-4"><i class="fas fa-calendar-week"></i> الخطة الأسبوعية للوجبات</h4>';
    
    for (let day = 1; day <= 7; day++) {
        const isIFDay = isIntermittentFasting && (day === 2 || day === 5); // Non-consecutive IF days
        
        weeklyPlan += `
            <div class="day-container mb-4">
                <div class="day-header">
                    <h5><i class="fas fa-calendar-day"></i> اليوم ${day} ${isIFDay ? '(صيام متقطع)' : ''}</h5>
                </div>
                <div class="meals-grid">
        `;
        
        if (isIFDay) {
            weeklyPlan += generateIntermittentFastingDay(recipes, dailyCalories * 0.8, dailyProtein, dailyCarbs, dailyFat);
        } else {
            weeklyPlan += generateRegularDay(recipes, mealsPerDay, dailyCalories, dailyProtein, dailyCarbs, dailyFat);
        }
        
        weeklyPlan += `
                </div>
            </div>
        `;
    }
    
    weeklyPlan += '</div>';
    mealPlanContainer.innerHTML = weeklyPlan;
    
    // Add click handlers for meal boxes
    addMealBoxClickHandlers();
}

// Get meals per day based on goal
function getMealsPerDay(goal) {
    if (goal === 'weight_gain' || goal === 'muscle_building') {
        return 4; // 3 main meals + 1 snack or 4 meals
    }
    return 3; // 3 main meals
}

// Generate regular day meals
function generateRegularDay(recipes, mealsPerDay, dailyCalories, dailyProtein, dailyCarbs, dailyFat) {
    let dayMeals = '';
    const mealTypes = mealsPerDay === 4 ? ['الإفطار', 'الغداء', 'العشاء', 'وجبة خفيفة'] : ['الإفطار', 'الغداء', 'العشاء'];
    const caloriesPerMeal = Math.round(dailyCalories / mealsPerDay);
    const proteinPerMeal = Math.round(dailyProtein / mealsPerDay);
    const carbsPerMeal = Math.round(dailyCarbs / mealsPerDay);
    const fatPerMeal = Math.round(dailyFat / mealsPerDay);
    
    mealTypes.forEach((mealType, index) => {
        const meal = selectMealFromRecipes(recipes, mealType, caloriesPerMeal);
        const alternative = selectAlternativeMeal(recipes, mealType, caloriesPerMeal, meal);
        
        dayMeals += createMealBox(meal, alternative, mealType, caloriesPerMeal, proteinPerMeal, carbsPerMeal, fatPerMeal);
    });
    
    return dayMeals;
}

// Generate intermittent fasting day
function generateIntermittentFastingDay(recipes, dailyCalories, dailyProtein, dailyCarbs, dailyFat) {
    let dayMeals = '';
    const mealTypes = ['الغداء', 'العشاء']; // Only 2 meals during eating window
    const caloriesPerMeal = Math.round(dailyCalories / 2);
    const proteinPerMeal = Math.round(dailyProtein / 2);
    const carbsPerMeal = Math.round(dailyCarbs / 2);
    const fatPerMeal = Math.round(dailyFat / 2);
    
    mealTypes.forEach((mealType, index) => {
        const meal = selectMealFromRecipes(recipes, mealType, caloriesPerMeal);
        const alternative = selectAlternativeMeal(recipes, mealType, caloriesPerMeal, meal);
        
        dayMeals += createMealBox(meal, alternative, mealType, caloriesPerMeal, proteinPerMeal, carbsPerMeal, fatPerMeal);
    });
    
    // Add fasting note
    dayMeals += `
        <div class="col-12">
            <div class="alert alert-info">
                <i class="fas fa-clock"></i> <strong>نافذة الأكل:</strong> 12:00 ظهراً - 8:00 مساءً (صيام 16 ساعة)
            </div>
        </div>
    `;
    
    return dayMeals;
}

// Create meal box with recipe details
function createMealBox(meal, alternative, mealType, targetCalories, targetProtein, targetCarbs, targetFat) {
    const mealId = `meal_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
    const altId = `alt_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
    
    return `
        <div class="col-md-6 col-lg-4 mb-3">
            <div class="meal-box" data-meal-id="${mealId}">
                <div class="meal-header">
                    <h6><i class="fas fa-utensils"></i> ${mealType}</h6>
                </div>
                <div class="meal-content">
                    <h6 class="meal-name">${meal.name}</h6>
                    <div class="nutrition-info">
                        <div class="nutrition-item">
                            <span class="label">السعرات:</span>
                            <span class="value">${meal.nutrition?.calories || targetCalories}</span>
                        </div>
                        <div class="nutrition-item">
                            <span class="label">البروتين:</span>
                            <span class="value">${meal.nutrition?.protein || targetProtein}g</span>
                        </div>
                        <div class="nutrition-item">
                            <span class="label">الكارب:</span>
                            <span class="value">${meal.nutrition?.carbs || targetCarbs}g</span>
                        </div>
                        <div class="nutrition-item">
                            <span class="label">الدهون:</span>
                            <span class="value">${meal.nutrition?.fat || targetFat}g</span>
                        </div>
                    </div>
                    <div class="ingredients-preview">
                        <small><strong>المكونات:</strong> ${meal.ingredients?.slice(0, 3).join(', ') || 'مكونات متنوعة'}...</small>
                    </div>
                    <div class="meal-actions mt-2">
                        <button class="btn btn-sm btn-outline-primary" onclick="showMealDetails('${mealId}', ${JSON.stringify(meal).replace(/"/g, '&quot;')})">
                            <i class="fas fa-eye"></i> خطوات التحضير
                        </button>
                        <button class="btn btn-sm btn-outline-secondary" onclick="showAlternative('${altId}', ${JSON.stringify(alternative).replace(/"/g, '&quot;')})">
                            <i class="fas fa-exchange-alt"></i> البديل
                        </button>
                    </div>
                </div>
            </div>
        </div>
    `;
}

// Select meal from recipes based on type and calories
function selectMealFromRecipes(recipes, mealType, targetCalories) {
    // Filter recipes suitable for meal type
    let suitableRecipes = recipes.filter(recipe => {
        const calories = recipe.nutrition?.calories || 300;
        return calories >= (targetCalories * 0.8) && calories <= (targetCalories * 1.2);
    });
    
    if (suitableRecipes.length === 0) {
        suitableRecipes = recipes; // Fallback to all recipes
    }
    
    // Select random recipe
    const selectedRecipe = suitableRecipes[Math.floor(Math.random() * suitableRecipes.length)];
    
    // Adjust portions if needed
    return adjustRecipePortions(selectedRecipe, targetCalories);
}

// Select alternative meal
function selectAlternativeMeal(recipes, mealType, targetCalories, mainMeal) {
    let alternatives = recipes.filter(recipe => 
        recipe.id !== mainMeal.id && 
        Math.abs((recipe.nutrition?.calories || 300) - targetCalories) <= 100
    );
    
    if (alternatives.length === 0) {
        alternatives = recipes.filter(recipe => recipe.id !== mainMeal.id);
    }
    
    const alternative = alternatives[Math.floor(Math.random() * alternatives.length)];
    return adjustRecipePortions(alternative, targetCalories);
}

// Adjust recipe portions to match target calories
function adjustRecipePortions(recipe, targetCalories) {
    const originalCalories = recipe.nutrition?.calories || 300;
    const ratio = targetCalories / originalCalories;
    
    return {
        ...recipe,
        portion_ratio: ratio,
        adjusted_nutrition: {
            calories: Math.round(originalCalories * ratio),
            protein: Math.round((recipe.nutrition?.protein || 20) * ratio),
            carbs: Math.round((recipe.nutrition?.carbs || 30) * ratio),
            fat: Math.round((recipe.nutrition?.fat || 10) * ratio)
        }
    };
}

// Add click handlers for meal boxes
function addMealBoxClickHandlers() {
    document.querySelectorAll('.meal-box').forEach(box => {
        box.addEventListener('click', function(e) {
            if (!e.target.closest('button')) {
                const mealId = this.dataset.mealId;
                this.classList.toggle('expanded');
            }
        });
    });
}

// Generate meal for specific type
function generateMealForType(mealType, calories, clientData, dietType) {
    // Mock meal data - in real implementation, this would fetch from your APIs
    const meals = {
        'الإفطار': [
            {
                name: 'شوفان بالفواكه والمكسرات',
                ingredients: 'شوفان، حليب، موز، توت، لوز',
                protein: Math.round(calories * 0.2 / 4),
                carbs: Math.round(calories * 0.5 / 4),
                fat: Math.round(calories * 0.3 / 9),
                preparation: 'اخلط الشوفان مع الحليب، أضف الفواكه والمكسرات'
            },
            {
                name: 'بيض مسلوق مع خبز أسمر',
                ingredients: 'بيض، خبز أسمر، أفوكادو، طماطم',
                protein: Math.round(calories * 0.25 / 4),
                carbs: Math.round(calories * 0.45 / 4),
                fat: Math.round(calories * 0.3 / 9),
                preparation: 'اسلقي البيض، قطعي الأفوكادو والطماطم'
            }
        ],
        'الغداء': [
            {
                name: 'دجاج مشوي مع أرز بني وخضار',
                ingredients: 'صدر دجاج، أرز بني، بروكلي، جزر',
                protein: Math.round(calories * 0.3 / 4),
                carbs: Math.round(calories * 0.4 / 4),
                fat: Math.round(calories * 0.3 / 9),
                preparation: 'اشوي الدجاج، اطبخي الأرز، اسلقي الخضار'
            },
            {
                name: 'سمك مع كينوا وسلطة',
                ingredients: 'فيليه سمك، كينوا، خس، خيار، طماطم',
                protein: Math.round(calories * 0.35 / 4),
                carbs: Math.round(calories * 0.35 / 4),
                fat: Math.round(calories * 0.3 / 9),
                preparation: 'اشوي السمك، اطبخي الكينوا، حضري السلطة'
            }
        ],
        'العشاء': [
            {
                name: 'سلطة البروتين',
                ingredients: 'دجاج، خس، طماطم، خيار، جبن قليل الدسم',
                protein: Math.round(calories * 0.4 / 4),
                carbs: Math.round(calories * 0.3 / 4),
                fat: Math.round(calories * 0.3 / 9),
                preparation: 'قطعي الخضار، أضيفي الدجاج والجبن'
            }
        ],
        'سناك صباحي': [
            {
                name: 'زبادي بالمكسرات',
                ingredients: 'زبادي يوناني، لوز، عسل',
                protein: Math.round(calories * 0.3 / 4),
                carbs: Math.round(calories * 0.4 / 4),
                fat: Math.round(calories * 0.3 / 9),
                preparation: 'اخلطي الزبادي مع المكسرات والعسل'
            }
        ],
        'سناك مسائي': [
            {
                name: 'تفاح مع زبدة اللوز',
                ingredients: 'تفاح، زبدة لوز طبيعية',
                protein: Math.round(calories * 0.15 / 4),
                carbs: Math.round(calories * 0.55 / 4),
                fat: Math.round(calories * 0.3 / 9),
                preparation: 'قطعي التفاح وادهنيه بزبدة اللوز'
            }
        ]
    };
    
    const mealOptions = meals[mealType] || meals['الإفطار'];
    return mealOptions[Math.floor(Math.random() * mealOptions.length)];
}

// Show meal details modal
function showMealDetails(day, mealType, mealIndex) {
    // Create modal content
    const meal = generateMealForType(mealType, 400, currentClientData, 'balanced'); // Mock data
    
    const modalHTML = `
        <div class="modal fade" id="mealModal" tabindex="-1">
            <div class="modal-dialog modal-lg">
                <div class="modal-content">
                    <div class="modal-header">
                        <h5 class="modal-title">${mealType} - ${day}</h5>
                        <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
                    </div>
                    <div class="modal-body">
                        <h6>المكونات:</h6>
                        <p>${meal.ingredients}</p>
                        
                        <h6>خطوات التحضير:</h6>
                        <p>${meal.preparation}</p>
                        
                        <h6>القيم الغذائية:</h6>
                        <div class="row">
                            <div class="col-3 text-center">
                                <div class="fw-bold text-primary">${meal.protein}g</div>
                                <small>بروتين</small>
                            </div>
                            <div class="col-3 text-center">
                                <div class="fw-bold text-warning">${meal.carbs}g</div>
                                <small>كارب</small>
                            </div>
                            <div class="col-3 text-center">
                                <div class="fw-bold text-info">${meal.fat}g</div>
                                <small>دهون</small>
                            </div>
                            <div class="col-3 text-center">
                                <div class="fw-bold text-success">400</div>
                                <small>سعرة</small>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    `;
    
    // Remove existing modal
    const existingModal = document.getElementById('mealModal');
    if (existingModal) {
        existingModal.remove();
    }
    
    // Add new modal
    document.body.insertAdjacentHTML('beforeend', modalHTML);
    
    // Show modal
    const modal = new bootstrap.Modal(document.getElementById('mealModal'));
    modal.show();
}

// Show meal preparation details
function showMealDetails(mealId, meal) {
    const modal = document.createElement('div');
    modal.className = 'modal fade';
    modal.innerHTML = `
        <div class="modal-dialog modal-lg">
            <div class="modal-content">
                <div class="modal-header">
                    <h5 class="modal-title"><i class="fas fa-utensils"></i> ${meal.name}</h5>
                    <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
                </div>
                <div class="modal-body">
                    <div class="row">
                        <div class="col-md-6">
                            <h6><i class="fas fa-list"></i> المكونات:</h6>
                            <ul class="list-group list-group-flush">
                                ${meal.ingredients?.map(ingredient => `<li class="list-group-item">${ingredient}</li>`).join('') || '<li class="list-group-item">مكونات متنوعة</li>'}
                            </ul>
                        </div>
                        <div class="col-md-6">
                            <h6><i class="fas fa-chart-bar"></i> القيم الغذائية:</h6>
                            <div class="nutrition-details">
                                <div class="nutrition-item">السعرات: ${meal.adjusted_nutrition?.calories || meal.nutrition?.calories || 300}</div>
                                <div class="nutrition-item">البروتين: ${meal.adjusted_nutrition?.protein || meal.nutrition?.protein || 20}g</div>
                                <div class="nutrition-item">الكربوهيدرات: ${meal.adjusted_nutrition?.carbs || meal.nutrition?.carbs || 30}g</div>
                                <div class="nutrition-item">الدهون: ${meal.adjusted_nutrition?.fat || meal.nutrition?.fat || 10}g</div>
                            </div>
                        </div>
                    </div>
                    <hr>
                    <h6><i class="fas fa-clipboard-list"></i> خطوات التحضير:</h6>
                    <ol class="preparation-steps">
                        ${meal.instructions?.map(step => `<li>${step}</li>`).join('') || generateDefaultInstructions(meal)}
                    </ol>
                    ${meal.portion_ratio && meal.portion_ratio !== 1 ? `
                        <div class="alert alert-info mt-3">
                            <i class="fas fa-info-circle"></i> <strong>ملاحظة:</strong> تم تعديل الكميات بنسبة ${Math.round(meal.portion_ratio * 100)}% لتناسب احتياجاتك الغذائية.
                        </div>
                    ` : ''}
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">إغلاق</button>
                    <button type="button" class="btn btn-primary" onclick="addToShoppingList('${meal.name}', ${JSON.stringify(meal.ingredients || []).replace(/"/g, '&quot;')})">
                        <i class="fas fa-cart-plus"></i> إضافة للتسوق
                    </button>
                </div>
            </div>
        </div>
    `;
    
    document.body.appendChild(modal);
    const bsModal = new bootstrap.Modal(modal);
    bsModal.show();
    
    modal.addEventListener('hidden.bs.modal', () => {
        document.body.removeChild(modal);
    });
}

// Generate default instructions if not provided
function generateDefaultInstructions(meal) {
    return `
        <li>حضر جميع المكونات المطلوبة</li>
        <li>اتبع طريقة الطبخ المناسبة للوصفة</li>
        <li>تأكد من نضج جميع المكونات</li>
        <li>قدم الطبق ساخناً</li>
    `;
}

// Show alternative meal
function showAlternative(mealType, calories) {
    const alternative = generateMealForType(mealType, calories, currentClientData, 'balanced');
    
    alert(`بديل لـ ${mealType}:\n\n${alternative.name}\nالمكونات: ${alternative.ingredients}\nالتحضير: ${alternative.preparation}`);
}

// Load countries for cuisine selection
async function loadCountries() {
    try {
        // Fetch countries from recipes API
        const response = await fetch('/api/v1/recipes');
        const data = await response.json();
        
        let countries = [];
        if (data.countries && data.countries.length > 0) {
            countries = data.countries;
        } else {
            // Fallback data if API fails
            countries = [
                { name: 'المطبخ العربي', code: 'arab', flag: '🇸🇦' },
                { name: 'المطبخ التركي', code: 'turkish', flag: '🇹🇷' },
                { name: 'المطبخ الإيطالي', code: 'italian', flag: '🇮🇹' },
                { name: 'المطبخ الهندي', code: 'indian', flag: '🇮🇳' },
                { name: 'المطبخ الصيني', code: 'chinese', flag: '🇨🇳' },
                { name: 'المطبخ المكسيكي', code: 'mexican', flag: '🇲🇽' },
                { name: 'المطبخ اليوناني', code: 'greek', flag: '🇬🇷' },
                { name: 'المطبخ الفرنسي', code: 'french', flag: '🇫🇷' }
            ];
        }
        
        const countryGrid = document.getElementById('countryGrid');
        
        countryGrid.innerHTML = countries.map(country => `
            <div class="country-card" onclick="selectCountry('${country.code}', '${country.name}')">
                <div style="font-size: 2rem; margin-bottom: 10px;">${country.flag}</div>
                <div>${country.name}</div>
            </div>
        `).join('');
    } catch (error) {
        console.error('Error loading countries:', error);
        // Use fallback data on error
        const countries = [
            { name: 'المطبخ العربي', code: 'arab', flag: '🇸🇦' },
            { name: 'المطبخ التركي', code: 'turkish', flag: '🇹🇷' },
            { name: 'المطبخ الإيطالي', code: 'italian', flag: '🇮🇹' },
            { name: 'المطبخ الهندي', code: 'indian', flag: '🇮🇳' },
            { name: 'المطبخ الصيني', code: 'chinese', flag: '🇨🇳' },
            { name: 'المطبخ المكسيكي', code: 'mexican', flag: '🇲🇽' },
            { name: 'المطبخ اليوناني', code: 'greek', flag: '🇬🇷' },
            { name: 'المطبخ الفرنسي', code: 'french', flag: '🇫🇷' }
        ];
        
        const countryGrid = document.getElementById('countryGrid');
        countryGrid.innerHTML = countries.map(country => `
            <div class="country-card" onclick="selectCountry('${country.code}', '${country.name}')">
                <div style="font-size: 2rem; margin-bottom: 10px;">${country.flag}</div>
                <div>${country.name}</div>
            </div>
        `).join('');
    }
}

// Select country for cuisine-based plan
function selectCountry(code, name) {
    selectedCountry = { code, name };
    
    // Update UI
    document.querySelectorAll('.country-card').forEach(card => {
        card.classList.remove('selected');
    });
    
    event.target.closest('.country-card').classList.add('selected');
}

// Show cuisine selection
function showCuisineSelection() {
    document.getElementById('countrySelection').style.display = 'block';
    document.getElementById('countrySelection').scrollIntoView({ behavior: 'smooth' });
}

// Generate cuisine-based plan
async function generateCuisineBasedPlan() {
    if (!selectedCountry) {
        alert('الرجاء اختيار مطبخ أولاً');
        return;
    }
    
    if (!currentClientData.name) {
        alert('الرجاء ملء بيانات العميل أولاً');
        return;
    }
    
    try {
        // Show loading message
        alert(`جاري توليد خطة غذائية مخصصة لـ ${selectedCountry.name}...`);
        
        // Fetch recipes for the selected country
        const response = await fetch('/api/v1/recipes');
        const data = await response.json();
        
        let countryRecipes = [];
        if (data.recipes && data.recipes.length > 0) {
            countryRecipes = data.recipes.filter(recipe => recipe.country === selectedCountry.code);
        }
        
        if (countryRecipes.length === 0) {
            alert(`لا توجد وصفات متاحة لـ ${selectedCountry.name} حالياً. سيتم استخدام وصفات عامة.`);
            countryRecipes = data.recipes || [];
        }
        
        // Generate plan based on selected cuisine
        const nutritionCalculations = calculateNutritionRequirements(currentClientData);
        
        // Filter meals based on selected cuisine and dietary restrictions
        const filteredRecipes = filterRecipesByRestrictions(countryRecipes, currentClientData);
        
        // Generate cuisine-specific meal plan
        generateCuisineMealPlan(selectedCountry, currentClientData, nutritionCalculations, filteredRecipes);
        
        alert(`تم توليد خطة غذائية بناءً على ${selectedCountry.name} بنجاح!`);
        
    } catch (error) {
        console.error('Error generating cuisine plan:', error);
        // Fallback to regular meal plan
        const nutritionCalculations = calculateNutritionRequirements(currentClientData);
        generateMealPlan(currentClientData, nutritionCalculations);
        alert(`تم توليد خطة غذائية عامة بدلاً من ${selectedCountry.name}`);
    }
}

// Check if recipe is safe for medical condition
function isRecipeSafeForMedicalCondition(recipe, medicalConditions) {
    const conditions = medicalConditions.toLowerCase();
    const recipeName = recipe.name.toLowerCase();
    const ingredients = recipe.ingredients ? recipe.ingredients.join(' ').toLowerCase() : '';
    
    // Diabetes restrictions
    if (conditions.includes('diabetes') || conditions.includes('سكري')) {
        const diabeticAvoidList = ['sugar', 'honey', 'syrup', 'candy', 'cake', 'cookie', 'soda', 'juice', 'white bread', 'white rice', 'pasta', 'سكر', 'عسل', 'حلوى', 'كيك', 'بسكويت', 'صودا', 'عصير', 'خبز أبيض', 'أرز أبيض', 'معكرونة'];
        if (diabeticAvoidList.some(item => recipeName.includes(item) || ingredients.includes(item))) {
            return false;
        }
        // High carb recipes should be limited
        if (recipe.nutrition && recipe.nutrition.carbs > 45) {
            return false;
        }
    }
    
    // Hypertension restrictions
    if (conditions.includes('hypertension') || conditions.includes('ضغط') || conditions.includes('pressure')) {
        const hypertensionAvoidList = ['salt', 'sodium', 'pickled', 'canned', 'processed', 'deli', 'bacon', 'sausage', 'ملح', 'صوديوم', 'مخلل', 'معلب', 'مصنع', 'بيكون', 'نقانق'];
        if (hypertensionAvoidList.some(item => recipeName.includes(item) || ingredients.includes(item))) {
            return false;
        }
    }
    
    // Heart disease restrictions
    if (conditions.includes('heart') || conditions.includes('قلب') || conditions.includes('cardiac')) {
        const heartAvoidList = ['fried', 'butter', 'cream', 'cheese', 'red meat', 'trans fat', 'مقلي', 'زبدة', 'كريمة', 'جبنة', 'لحم أحمر'];
        if (heartAvoidList.some(item => recipeName.includes(item) || ingredients.includes(item))) {
            return false;
        }
        // High saturated fat recipes should be limited
        if (recipe.nutrition && recipe.nutrition.fat > 20) {
            return false;
        }
    }
    
    // Kidney disease restrictions
    if (conditions.includes('kidney') || conditions.includes('كلى') || conditions.includes('renal')) {
        const kidneyAvoidList = ['banana', 'orange', 'potato', 'tomato', 'nuts', 'dairy', 'chocolate', 'موز', 'برتقال', 'بطاطس', 'طماطم', 'مكسرات', 'ألبان', 'شوكولاتة'];
        if (kidneyAvoidList.some(item => recipeName.includes(item) || ingredients.includes(item))) {
            return false;
        }
        // High protein recipes should be limited
        if (recipe.nutrition && recipe.nutrition.protein > 25) {
            return false;
        }
    }
    
    // Liver disease restrictions
    if (conditions.includes('liver') || conditions.includes('كبد') || conditions.includes('hepatic')) {
        const liverAvoidList = ['alcohol', 'wine', 'beer', 'raw', 'shellfish', 'high fat', 'كحول', 'نبيذ', 'بيرة', 'نيء', 'محار', 'دهون عالية'];
        if (liverAvoidList.some(item => recipeName.includes(item) || ingredients.includes(item))) {
            return false;
        }
    }
    
    // Celiac disease restrictions
    if (conditions.includes('celiac') || conditions.includes('gluten') || conditions.includes('سيلياك') || conditions.includes('جلوتين')) {
        const glutenAvoidList = ['wheat', 'barley', 'rye', 'bread', 'pasta', 'flour', 'قمح', 'شعير', 'جاودار', 'خبز', 'معكرونة', 'دقيق'];
        if (glutenAvoidList.some(item => recipeName.includes(item) || ingredients.includes(item))) {
            return false;
        }
    }
    
    return true;
}

// Filter recipes based on dietary restrictions and allergies
function filterRecipesByRestrictions(recipes, clientData) {
    if (!clientData.foodRestrictions && !clientData.medicalConditions) {
        return recipes;
    }
    
    let filteredRecipes = recipes.filter(recipe => {
        // Check medical conditions first
        if (clientData.medicalConditions && clientData.medicalConditions.trim() !== '') {
            if (!isRecipeSafeForMedicalCondition(recipe, clientData.medicalConditions)) {
                return false;
            }
        }
        
        // Check for allergens and food restrictions
        if (clientData.foodRestrictions && clientData.foodRestrictions.trim() !== '') {
            const excludedFoods = clientData.foodRestrictions.toLowerCase().split(',').map(item => item.trim());
            
            // Check for allergens
            if (recipe.allergens && recipe.allergens.length > 0) {
                const hasAllergen = recipe.allergens.some(allergen => 
                    excludedFoods.some(excluded => 
                        allergen.toLowerCase().includes(excluded) || 
                        excluded.includes(allergen.toLowerCase())
                    )
                );
                if (hasAllergen) return false;
            }
            
            // Check ingredients for excluded foods with safe data access
            if (recipe.ingredients && recipe.ingredients.length > 0) {
                const hasExcludedIngredient = recipe.ingredients.some(ingredient => {
                    // Use safe data access if error handler is available
                    let ingredientName;
                    if (window.errorHandler && typeof window.errorHandler.safeGet === 'function') {
                        ingredientName = window.errorHandler.safeGet(ingredient, 'item', '') || 
                                       window.errorHandler.safeGet(ingredient, 'name', '') || 
                                       (typeof ingredient === 'string' ? ingredient : '');
                    } else {
                        // Fallback safe access
                        ingredientName = ingredient.item || ingredient.name || ingredient || '';
                    }
                    
                    if (!ingredientName || typeof ingredientName !== 'string') {
                        return false;
                    }
                    
                    return excludedFoods.some(excluded => {
                        try {
                            return ingredientName.toLowerCase().includes(excluded) || 
                                   excluded.includes(ingredientName.toLowerCase());
                        } catch (error) {
                            if (window.errorHandler) {
                                window.errorHandler.logError({
                                    type: 'Ingredient Filtering Error',
                                    message: `Error filtering ingredient: ${error.message}`,
                                    data: { ingredient, excluded },
                                    timestamp: new Date().toISOString()
                                });
                            }
                            return false;
                        }
                    });
                });
                if (hasExcludedIngredient) return false;
            }
        }
        
        return true;
    });
    
    // Apply halal filtering if requested
    if (clientData.foodRestrictions && 
        (clientData.foodRestrictions.toLowerCase().includes('halal') || 
         clientData.foodRestrictions.toLowerCase().includes('حلال'))) {
        
        // First try to get naturally halal recipes
        const halalRecipes = window.HalalFilter ? 
            window.HalalFilter.filterHalalRecipes(filteredRecipes) : 
            filteredRecipes;
        
        // If we have halal recipes, use them
        if (halalRecipes.length > 0) {
            filteredRecipes = halalRecipes;
        } else if (window.HalalFilter) {
            // If no naturally halal recipes, try with alternatives
            filteredRecipes = filteredRecipes.map(recipe => 
                window.HalalFilter.replaceWithHalalAlternatives(recipe)
            ).filter(recipe => 
                window.HalalFilter.filterHalalRecipes([recipe]).length > 0
            );
        }
        
        // Display halal compliance info if container exists
        if (window.HalalFilter && document.getElementById('halalComplianceInfo')) {
            window.HalalFilter.displayHalalComplianceInfo(filteredRecipes, 'halalComplianceInfo');
        }
    }
    
    return filteredRecipes;
}

// Generate cuisine-specific meal plan
function generateCuisineMealPlan(cuisine, clientData, calculations, recipes = []) {
    console.log('Generating cuisine-based meal plan for:', cuisine.name);
    
    // Add cuisine information header
    const mealPlanContainer = document.getElementById('mealPlanContainer');
    if (mealPlanContainer) {
        const cuisineHeader = `
            <div class="cuisine-header mb-4">
                <div class="alert alert-info">
                    <h5><i class="fas fa-globe me-2"></i>خطة غذائية مخصصة - ${cuisine.name}</h5>
                    <p class="mb-0">تم إنشاء هذه الخطة بناءً على المطبخ ${cuisine.name} مع مراعاة احتياجاتك الغذائية والصحية</p>
                </div>
            </div>
        `;
        mealPlanContainer.insertAdjacentHTML('afterbegin', cuisineHeader);
    }
    
    // Filter recipes for this cuisine
    const cuisineRecipes = recipes.filter(recipe => recipe.country === cuisine.code);
    const filteredRecipes = filterRecipesByRestrictions(cuisineRecipes, clientData);
    
    // Generate weekly meal plan using the new system
    generateWeeklyMealPlan(clientData, filteredRecipes);
}

// Fallback cuisine meal plan generation (old system)
function generateFallbackCuisineMealPlan(cuisine, clientData, calculations, recipes = []) {
    console.log('Generating fallback cuisine-based meal plan for:', cuisine.name);
    
    const mealPlanGrid = document.getElementById('mealPlanGrid');
    
    // Generate 7-day meal plan
    const weekDays = ['الأحد', 'الاثنين', 'الثلاثاء', 'الأربعاء', 'الخميس', 'الجمعة', 'السبت'];
    const mealTypes = ['الإفطار', 'سناك صباحي', 'الغداء', 'سناك مسائي', 'العشاء'];
    
    let planHTML = '';
    
    weekDays.forEach((day, dayIndex) => {
        planHTML += `
            <div class="col-12 mb-4">
                <h4 class="text-primary mb-3">
                    <i class="fas fa-calendar-day me-2"></i>${day} - ${cuisine.name}
                </h4>
                <div class="row">
        `;
        
        mealTypes.forEach((mealType, mealIndex) => {
            const mealCalories = Math.round(calculations.totalCalories * calculations.mealDistribution.regularDays.distribution[mealIndex]);
            const meal = generateCuisineMealForType(mealType, mealCalories, clientData, calculations.recommendedDietType, recipes, cuisine);
            
            planHTML += `
                <div class="col-md-6 col-lg-4 mb-3">
                    <div class="meal-box" onclick="showMealDetails('${day}', '${mealType}', ${mealIndex})">
                        <h5>${mealType}</h5>
                        <p class="meal-name">${meal.name}</p>
                        <p class="meal-ingredients">${meal.ingredients}</p>
                        <p class="cuisine-label"><small><i class="fas fa-globe"></i> ${cuisine.name}</small></p>
                        
                        <div class="nutrition-info">
                            <div class="nutrition-item">
                                <div class="value">${mealCalories}</div>
                                <div class="label">سعرة</div>
                            </div>
                            <div class="nutrition-item">
                                <div class="value">${meal.protein}g</div>
                                <div class="label">بروتين</div>
                            </div>
                            <div class="nutrition-item">
                                <div class="value">${meal.carbs}g</div>
                                <div class="label">كارب</div>
                            </div>
                            <div class="nutrition-item">
                                <div class="value">${meal.fat}g</div>
                                <div class="label">دهون</div>
                            </div>
                        </div>
                        
                        <div class="mt-3">
                            <button class="btn btn-sm btn-outline-primary" onclick="event.stopPropagation(); showAlternative('${mealType}', ${mealCalories})">
                                <i class="fas fa-exchange-alt me-1"></i>بديل
                            </button>
                        </div>
                    </div>
                </div>
            `;
        });
        
        planHTML += `
                </div>
            </div>
        `;
    });
    
    mealPlanGrid.innerHTML = planHTML;
}

// Generate cuisine-specific meal for specific type
function generateCuisineMealForType(mealType, calories, clientData, dietType, recipes, cuisine) {
    // Try to find cuisine-specific recipes first
    let suitableRecipes = recipes.filter(recipe => 
        recipe.category === mealType && recipe.country === cuisine.code
    );
    
    // If no cuisine-specific recipes, use general recipes
    if (suitableRecipes.length === 0) {
        suitableRecipes = recipes.filter(recipe => recipe.category === mealType);
    }
    
    // If still no recipes, use fallback
    if (suitableRecipes.length === 0) {
        return generateMealForType(mealType, calories, clientData, dietType);
    }
    
    const selectedRecipe = suitableRecipes[Math.floor(Math.random() * suitableRecipes.length)];
    
    return {
        name: selectedRecipe.name || `وجبة ${mealType}`,
        ingredients: selectedRecipe.ingredients ? 
            selectedRecipe.ingredients.map(ing => `${ing.name} - ${ing.amount}`).join(', ') : 
            'مكونات متنوعة',
        protein: Math.round(calories * 0.2 / 4),
        carbs: Math.round(calories * 0.5 / 4),
        fat: Math.round(calories * 0.3 / 9),
        preparation: selectedRecipe.instructions ? 
            selectedRecipe.instructions.join('. ') : 
            'تعليمات التحضير متاحة في التفاصيل',
        cuisine: cuisine.name
    };
}

// Generate shopping list
function generateShoppingList() {
    if (!generatedPlan) {
        alert('الرجاء توليد خطة غذائية أولاً');
        return;
    }
    
    // Mock shopping list
    const shoppingList = [
        'دجاج (2 كيلو)',
        'أرز بني (1 كيلو)',
        'خضار متنوعة',
        'فواكه طازجة',
        'بيض (12 حبة)',
        'حليب قليل الدسم',
        'زبادي يوناني',
        'مكسرات متنوعة'
    ];
    
    alert('قائمة التسوق:\n\n' + shoppingList.join('\n'));
}

// Download PDF
function downloadPDF() {
    if (!generatedPlan) {
        alert('الرجاء توليد خطة غذائية أولاً');
        return;
    }
    
    alert('سيتم تطوير وظيفة تحميل PDF قريباً');
}

// Language switching function
function switchLanguage(lang) {
    document.documentElement.lang = lang;
    document.documentElement.dir = lang === 'ar' ? 'rtl' : 'ltr';
    
    // Update text content based on language
    // This would integrate with your existing language system
    console.log('Language switched to:', lang);
}