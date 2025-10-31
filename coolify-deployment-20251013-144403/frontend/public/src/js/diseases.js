// Disease Management System
let diseaseData = [];
let oldDiseaseData = [];
let metabolismData = [];
let skillsData = [];

// Initialize the disease management system
document.addEventListener('DOMContentLoaded', function() {
    loadAllData();
    setupEventListeners();
});

// Load all required data files
async function loadAllData() {
    try {
        await Promise.all([
            loadDiseaseData(),
            loadOldDiseaseData(),
            loadMetabolismData(),
            loadSkillsData()
        ]);
        
        populateDiseaseOptions();
        populateComplaintOptions();
        
        console.log('All disease management data loaded successfully');
    } catch (error) {
        console.error('Error loading disease data:', error);
        showMessage('حدث خطأ في تحميل البيانات', 'error');
    }
}

// Load disease data
async function loadDiseaseData() {
    try {
        const response = await fetch('../backend/data/disease.json');
        if (!response.ok) throw new Error('Failed to load disease data');
        diseaseData = await response.json();
    } catch (error) {
        console.error('Error loading disease data:', error);
        diseaseData = [];
    }
}

// Load old disease data
async function loadOldDiseaseData() {
    try {
        const response = await fetch('../backend/data/old disease.json');
        if (!response.ok) throw new Error('Failed to load old disease data');
        oldDiseaseData = await response.json();
    } catch (error) {
        console.error('Error loading old disease data:', error);
        oldDiseaseData = [];
    }
}

// Load metabolism data
async function loadMetabolismData() {
    try {
        const response = await fetch('../backend/data/metabolism.js');
        if (!response.ok) throw new Error('Failed to load metabolism data');
        const text = await response.text();
        
        // Extract JSON from JavaScript file
        const jsonMatch = text.match(/const\s+metabolism\s*=\s*(\{[\s\S]*?\});?\s*(?:\/\/|$)/m);
        if (jsonMatch) {
            metabolismData = JSON.parse(jsonMatch[1]);
        }
    } catch (error) {
        console.error('Error loading metabolism data:', error);
        metabolismData = {};
    }
}

// Load skills data
async function loadSkillsData() {
    try {
        const response = await fetch('../backend/data/skills.js');
        if (!response.ok) throw new Error('Failed to load skills data');
        const text = await response.text();
        
        // Extract JSON from JavaScript file
        const jsonMatch = text.match(/const\s+skills\s*=\s*(\{[\s\S]*?\});?\s*(?:\/\/|$)/m);
        if (jsonMatch) {
            skillsData = JSON.parse(jsonMatch[1]);
        }
    } catch (error) {
        console.error('Error loading skills data:', error);
        skillsData = {};
    }
}

// Populate disease options
function populateDiseaseOptions() {
    const diseaseSelect = document.getElementById('diseaseSelect');
    diseaseSelect.innerHTML = '<option value="">اختر المرض</option>';
    
    // Combine both disease datasets
    const allDiseases = [...diseaseData, ...oldDiseaseData];
    
    allDiseases.forEach(disease => {
        const option = document.createElement('option');
        option.value = disease.name || disease.title;
        option.textContent = disease.name || disease.title;
        diseaseSelect.appendChild(option);
    });
}

// Populate complaint options from metabolism data
function populateComplaintOptions() {
    const complaintSelect = document.getElementById('complaintSelect');
    complaintSelect.innerHTML = '<option value="">اختر الشكوى (اختياري)</option>';
    
    if (metabolismData.complaints) {
        metabolismData.complaints.forEach(complaint => {
            const option = document.createElement('option');
            option.value = complaint.name || complaint.title;
            option.textContent = complaint.name || complaint.title;
            complaintSelect.appendChild(option);
        });
    }
}

// Setup event listeners
function setupEventListeners() {
    const form = document.getElementById('diseaseForm');
    form.addEventListener('submit', handleFormSubmit);
}

// Handle form submission
async function handleFormSubmit(event) {
    event.preventDefault();
    
    if (!validateForm()) {
        return;
    }
    
    const clientData = collectClientData();
    
    // Show loading modal
    const loadingModal = new bootstrap.Modal(document.getElementById('loadingModal'));
    loadingModal.show();
    
    try {
        const treatmentPlan = await generateTreatmentPlan(clientData);
        displayTreatmentResults(treatmentPlan);
        loadingModal.hide();
        showMessage('تم إنشاء خطة العلاج بنجاح!', 'success');
    } catch (error) {
        console.error('Error generating treatment plan:', error);
        loadingModal.hide();
        showMessage('حدث خطأ في إنشاء خطة العلاج', 'error');
    }
}

// Enhanced validation using error handler
function validateForm() {
    const clientData = collectClientData();
    
    // Use enhanced error handler validation if available
    if (window.errorHandler && typeof window.errorHandler.validateClientData === 'function') {
        const errors = window.errorHandler.validateClientData({
            name: clientData.name,
            age: clientData.age,
            weight: clientData.weight,
            height: clientData.height,
            gender: clientData.gender,
            activityLevel: clientData.activityLevel
        });
        
        // Add disease-specific validations
        if (!clientData.disease) {
            errors.push('الرجاء اختيار المرض - Please select a disease');
        }
        
        // Show validation errors using enhanced error handler
        if (errors.length > 0) {
            if (typeof window.errorHandler.showValidationErrors === 'function') {
                window.errorHandler.showValidationErrors(errors, 'diseaseValidationErrors');
            } else {
                showMessage(errors.join('\n'), 'error');
            }
            return false;
        }
        
        // Clear any existing validation errors
        const errorContainer = document.getElementById('diseaseValidationErrors');
        if (errorContainer) {
            errorContainer.innerHTML = '';
            errorContainer.style.display = 'none';
        }
        
        return true;
    }
    
    // Fallback validation if error handler is not available
    const requiredFields = [
        'clientName', 'clientAge', 'clientGender', 
        'clientWeight', 'clientHeight', 'activityLevel', 'diseaseSelect'
    ];
    
    for (const fieldId of requiredFields) {
        const field = document.getElementById(fieldId);
        if (!field.value.trim()) {
            field.focus();
            showMessage(`يرجى ملء حقل ${field.previousElementSibling.textContent}`, 'error');
            return false;
        }
    }
    
    return true;
}

// Collect client data from form
function collectClientData() {
    return {
        name: document.getElementById('clientName').value.trim(),
        age: parseInt(document.getElementById('clientAge').value),
        gender: document.getElementById('clientGender').value,
        weight: parseFloat(document.getElementById('clientWeight').value),
        height: parseInt(document.getElementById('clientHeight').value),
        activityLevel: document.getElementById('activityLevel').value,
        disease: document.getElementById('diseaseSelect').value,
        complaint: document.getElementById('complaintSelect').value,
        currentMedication: document.getElementById('currentMedication').value.trim(),
        medicalHistory: document.getElementById('medicalHistory').value.trim(),
        symptoms: document.getElementById('symptoms').value.trim()
    };
}

// Generate comprehensive treatment plan
async function generateTreatmentPlan(clientData) {
    const plan = {
        clientInfo: clientData,
        diseaseInfo: findDiseaseInfo(clientData.disease),
        treatmentPlan: [],
        nutritionAdvice: [],
        lifestyle: [],
        medications: [],
        followUp: [],
        cookingTips: [],
        medicalDisclaimer: "هذا الموقع لتحديثك بالمعلومات وتوجيهك للنصائح المفيدة بهدف التعليم والوعي ولا يغني عن زيارة الطبيب. يجب استشارة طبيب مختص قبل تطبيق أي نصائح أو تغيير في العلاج."
    };
    
    // Generate treatment based on disease
    if (plan.diseaseInfo) {
        plan.treatmentPlan = generateDiseaseSpecificTreatment(plan.diseaseInfo, clientData);
        plan.nutritionAdvice = generateNutritionAdvice(plan.diseaseInfo, clientData);
        plan.lifestyle = generateLifestyleAdvice(plan.diseaseInfo, clientData);
        plan.medications = generateMedicationAdvice(plan.diseaseInfo, clientData);
        plan.followUp = generateFollowUpPlan(plan.diseaseInfo, clientData);
    }
    
    // Add complaint-specific advice if available
    if (clientData.complaint) {
        const complaintAdvice = generateComplaintAdvice(clientData.complaint, clientData);
        if (complaintAdvice) {
            plan.nutritionAdvice.push(...complaintAdvice.nutrition);
            plan.lifestyle.push(...complaintAdvice.lifestyle);
        }
    }
    
    // Add cooking tips from skills data
    plan.cookingTips = generateCookingTips(plan.nutritionAdvice);
    
    return plan;
}

// Find disease information
function findDiseaseInfo(diseaseName) {
    // Search in disease data first
    let disease = diseaseData.find(d => d.name === diseaseName || d.title === diseaseName);
    
    // If not found, search in old disease data
    if (!disease) {
        disease = oldDiseaseData.find(d => d.name === diseaseName || d.title === diseaseName);
    }
    
    return disease;
}

// Generate disease-specific treatment
function generateDiseaseSpecificTreatment(diseaseInfo, clientData) {
    const treatments = [];
    
    if (diseaseInfo.treatment) {
        treatments.push({
            category: "العلاج الأساسي",
            description: diseaseInfo.treatment,
            priority: "عالية"
        });
    }
    
    if (diseaseInfo.management) {
        treatments.push({
            category: "إدارة المرض",
            description: diseaseInfo.management,
            priority: "عالية"
        });
    }
    
    if (diseaseInfo.prevention) {
        treatments.push({
            category: "الوقاية",
            description: diseaseInfo.prevention,
            priority: "متوسطة"
        });
    }
    
    return treatments;
}

// Generate nutrition advice
function generateNutritionAdvice(diseaseInfo, clientData) {
    const advice = [];
    
    if (diseaseInfo.diet) {
        advice.push({
            category: "النظام الغذائي",
            recommendations: Array.isArray(diseaseInfo.diet) ? diseaseInfo.diet : [diseaseInfo.diet],
            importance: "عالية"
        });
    }
    
    if (diseaseInfo.foods_to_avoid) {
        advice.push({
            category: "الأطعمة التي يجب تجنبها",
            recommendations: Array.isArray(diseaseInfo.foods_to_avoid) ? diseaseInfo.foods_to_avoid : [diseaseInfo.foods_to_avoid],
            importance: "عالية"
        });
    }
    
    if (diseaseInfo.recommended_foods) {
        advice.push({
            category: "الأطعمة المُوصى بها",
            recommendations: Array.isArray(diseaseInfo.recommended_foods) ? diseaseInfo.recommended_foods : [diseaseInfo.recommended_foods],
            importance: "عالية"
        });
    }
    
    // Add general nutrition advice based on BMI
    const bmi = calculateBMI(clientData.weight, clientData.height);
    advice.push(generateBMIBasedAdvice(bmi, clientData));
    
    return advice;
}

// Generate lifestyle advice
function generateLifestyleAdvice(diseaseInfo, clientData) {
    const advice = [];
    
    if (diseaseInfo.lifestyle) {
        advice.push({
            category: "نمط الحياة",
            recommendations: Array.isArray(diseaseInfo.lifestyle) ? diseaseInfo.lifestyle : [diseaseInfo.lifestyle]
        });
    }
    
    if (diseaseInfo.exercise) {
        advice.push({
            category: "التمارين الرياضية",
            recommendations: Array.isArray(diseaseInfo.exercise) ? diseaseInfo.exercise : [diseaseInfo.exercise]
        });
    }
    
    // Add general lifestyle advice
    advice.push({
        category: "نصائح عامة",
        recommendations: [
            "الحصول على نوم كافٍ (7-8 ساعات يومياً)",
            "إدارة التوتر والضغط النفسي",
            "شرب كمية كافية من الماء (8-10 أكواب يومياً)",
            "تجنب التدخين والكحول",
            "المتابعة الدورية مع الطبيب المختص"
        ]
    });
    
    return advice;
}

// Generate medication advice
function generateMedicationAdvice(diseaseInfo, clientData) {
    const medications = [];
    
    if (diseaseInfo.medications) {
        medications.push({
            category: "الأدوية الموصوفة",
            details: Array.isArray(diseaseInfo.medications) ? diseaseInfo.medications : [diseaseInfo.medications],
            note: "يجب استشارة الطبيب قبل تناول أي دواء"
        });
    }
    
    if (diseaseInfo.supplements) {
        medications.push({
            category: "المكملات الغذائية",
            details: Array.isArray(diseaseInfo.supplements) ? diseaseInfo.supplements : [diseaseInfo.supplements],
            note: "يُفضل استشارة الطبيب أو الصيدلي"
        });
    }
    
    return medications;
}

// Generate follow-up plan
function generateFollowUpPlan(diseaseInfo, clientData) {
    const followUp = [];
    
    if (diseaseInfo.monitoring) {
        followUp.push({
            category: "المتابعة والمراقبة",
            schedule: Array.isArray(diseaseInfo.monitoring) ? diseaseInfo.monitoring : [diseaseInfo.monitoring]
        });
    }
    
    // Add general follow-up recommendations
    followUp.push({
        category: "المتابعة العامة",
        schedule: [
            "فحص دوري كل 3-6 أشهر",
            "مراقبة الأعراض يومياً",
            "قياس الوزن أسبوعياً",
            "تسجيل أي تغييرات في الحالة الصحية"
        ]
    });
    
    return followUp;
}

// Generate complaint-specific advice
function generateComplaintAdvice(complaint, clientData) {
    if (!metabolismData.complaints) return null;
    
    const complaintInfo = metabolismData.complaints.find(c => 
        c.name === complaint || c.title === complaint
    );
    
    if (!complaintInfo) return null;
    
    return {
        nutrition: [{
            category: `نصائح غذائية لـ ${complaint}`,
            recommendations: complaintInfo.nutrition || [],
            importance: "متوسطة"
        }],
        lifestyle: [{
            category: `نصائح نمط الحياة لـ ${complaint}`,
            recommendations: complaintInfo.lifestyle || []
        }]
    };
}

// Generate cooking tips
function generateCookingTips(nutritionAdvice) {
    const tips = [];
    
    if (skillsData.cooking_tips) {
        // Add general cooking tips
        tips.push({
            category: "نصائح الطبخ العامة",
            tips: skillsData.cooking_tips.general || []
        });
        
        // Add healthy cooking methods
        if (skillsData.cooking_tips.healthy_methods) {
            tips.push({
                category: "طرق الطبخ الصحية",
                tips: skillsData.cooking_tips.healthy_methods
            });
        }
    }
    
    return tips;
}

// Calculate BMI
function calculateBMI(weight, height) {
    const heightInMeters = height / 100;
    return weight / (heightInMeters * heightInMeters);
}

// Generate BMI-based advice
function generateBMIBasedAdvice(bmi, clientData) {
    let category, recommendations;
    
    if (bmi < 18.5) {
        category = "نقص الوزن";
        recommendations = [
            "زيادة السعرات الحرارية بطريقة صحية",
            "تناول وجبات صغيرة ومتكررة",
            "إضافة البروتينات والدهون الصحية",
            "ممارسة تمارين القوة لبناء العضلات"
        ];
    } else if (bmi >= 18.5 && bmi < 25) {
        category = "الوزن الطبيعي";
        recommendations = [
            "الحفاظ على النظام الغذائي المتوازن",
            "ممارسة الرياضة بانتظام",
            "تناول الفواكه والخضروات يومياً",
            "شرب كمية كافية من الماء"
        ];
    } else if (bmi >= 25 && bmi < 30) {
        category = "زيادة الوزن";
        recommendations = [
            "تقليل السعرات الحرارية تدريجياً",
            "زيادة النشاط البدني",
            "تجنب الأطعمة المصنعة والسكريات",
            "تناول البروتين في كل وجبة"
        ];
    } else {
        category = "السمنة";
        recommendations = [
            "استشارة أخصائي تغذية",
            "وضع خطة لفقدان الوزن تدريجياً",
            "ممارسة الرياضة بإشراف طبي",
            "مراقبة السعرات الحرارية يومياً"
        ];
    }
    
    return {
        category: `نصائح للوزن (BMI: ${bmi.toFixed(1)} - ${category})`,
        recommendations: recommendations,
        importance: "متوسطة"
    };
}

// Display treatment results
function displayTreatmentResults(plan) {
    const container = document.getElementById('treatmentPlanResults');
    container.style.display = 'block';
    
    let html = `
        <div class="treatment-plan-header">
            <h2><i class="fas fa-heartbeat"></i> خطة العلاج الشخصية لـ ${plan.clientInfo.name}</h2>
        </div>
        
        <div class="client-summary">
            <h3><i class="fas fa-user-md"></i> ملخص الحالة</h3>
            <div class="row">
                <div class="col-md-6">
                    <ul class="client-details">
                        <li><strong>العمر:</strong> ${plan.clientInfo.age} سنة</li>
                        <li><strong>الجنس:</strong> ${plan.clientInfo.gender === 'male' ? 'ذكر' : 'أنثى'}</li>
                        <li><strong>الوزن:</strong> ${plan.clientInfo.weight} كجم</li>
                        <li><strong>الطول:</strong> ${plan.clientInfo.height} سم</li>
                    </ul>
                </div>
                <div class="col-md-6">
                    <ul class="client-details">
                        <li><strong>المرض:</strong> ${plan.clientInfo.disease}</li>
                        <li><strong>مستوى النشاط:</strong> ${getActivityLevelName(plan.clientInfo.activityLevel)}</li>
                        ${plan.clientInfo.complaint ? `<li><strong>الشكوى:</strong> ${plan.clientInfo.complaint}</li>` : ''}
                        <li><strong>BMI:</strong> ${calculateBMI(plan.clientInfo.weight, plan.clientInfo.height).toFixed(1)}</li>
                    </ul>
                </div>
            </div>
        </div>
    `;
    
    // Display treatment plan
    if (plan.treatmentPlan && plan.treatmentPlan.length > 0) {
        html += `
            <div class="treatment-section">
                <h3><i class="fas fa-stethoscope"></i> خطة العلاج</h3>
        `;
        
        plan.treatmentPlan.forEach(treatment => {
            html += `
                <div class="treatment-card">
                    <h4>${treatment.category}</h4>
                    <div class="priority-badge priority-${treatment.priority}">${treatment.priority}</div>
                    <p>${treatment.description}</p>
                </div>
            `;
        });
        
        html += `</div>`;
    }
    
    // Display nutrition advice
    if (plan.nutritionAdvice && plan.nutritionAdvice.length > 0) {
        html += `
            <div class="nutrition-section">
                <h3><i class="fas fa-apple-alt"></i> النصائح الغذائية</h3>
        `;
        
        plan.nutritionAdvice.forEach(advice => {
            html += `
                <div class="nutrition-card">
                    <h4>${advice.category}</h4>
                    <div class="importance-badge importance-${advice.importance}">${advice.importance}</div>
                    <ul>
                        ${advice.recommendations.map(rec => `<li>${rec}</li>`).join('')}
                    </ul>
                </div>
            `;
        });
        
        html += `</div>`;
    }
    
    // Display lifestyle advice
    if (plan.lifestyle && plan.lifestyle.length > 0) {
        html += `
            <div class="lifestyle-section">
                <h3><i class="fas fa-heart"></i> نصائح نمط الحياة</h3>
        `;
        
        plan.lifestyle.forEach(lifestyle => {
            html += `
                <div class="lifestyle-card">
                    <h4>${lifestyle.category}</h4>
                    <ul>
                        ${lifestyle.recommendations.map(rec => `<li>${rec}</li>`).join('')}
                    </ul>
                </div>
            `;
        });
        
        html += `</div>`;
    }
    
    // Display medications
    if (plan.medications && plan.medications.length > 0) {
        html += `
            <div class="medications-section">
                <h3><i class="fas fa-pills"></i> الأدوية والمكملات</h3>
        `;
        
        plan.medications.forEach(medication => {
            html += `
                <div class="medication-card">
                    <h4>${medication.category}</h4>
                    <ul>
                        ${medication.details.map(detail => `<li>${detail}</li>`).join('')}
                    </ul>
                    <div class="medication-note">
                        <i class="fas fa-exclamation-triangle"></i>
                        ${medication.note}
                    </div>
                </div>
            `;
        });
        
        html += `</div>`;
    }
    
    // Display cooking tips
    if (plan.cookingTips && plan.cookingTips.length > 0) {
        html += `
            <div class="cooking-tips-section">
                <h3><i class="fas fa-utensils"></i> نصائح الطبخ</h3>
        `;
        
        plan.cookingTips.forEach(tipCategory => {
            html += `
                <div class="cooking-tips-card">
                    <h4>${tipCategory.category}</h4>
                    <ul>
                        ${tipCategory.tips.map(tip => `<li>${tip}</li>`).join('')}
                    </ul>
                </div>
            `;
        });
        
        html += `</div>`;
    }
    
    // Display follow-up plan
    if (plan.followUp && plan.followUp.length > 0) {
        html += `
            <div class="followup-section">
                <h3><i class="fas fa-calendar-check"></i> خطة المتابعة</h3>
        `;
        
        plan.followUp.forEach(followup => {
            html += `
                <div class="followup-card">
                    <h4>${followup.category}</h4>
                    <ul>
                        ${followup.schedule.map(item => `<li>${item}</li>`).join('')}
                    </ul>
                </div>
            `;
        });
        
        html += `</div>`;
    }
    
    // Add medical disclaimer
    html += `
        <div class="medical-disclaimer">
            <h4><i class="fas fa-exclamation-triangle"></i> تنويه طبي</h4>
            <p>${plan.medicalDisclaimer}</p>
        </div>
    `;
    
    container.innerHTML = html;
    container.scrollIntoView({ behavior: 'smooth' });
}

// Helper function to get activity level name in Arabic
function getActivityLevelName(level) {
    const levels = {
        'sedentary': 'قليل النشاط',
        'light': 'نشاط خفيف',
        'moderate': 'نشاط متوسط',
        'active': 'نشط',
        'very_active': 'نشط جداً'
    };
    return levels[level] || level;
}

// Show message to user
function showMessage(message, type = 'info') {
    // Create message element
    const messageDiv = document.createElement('div');
    messageDiv.className = `alert alert-${type === 'error' ? 'danger' : type === 'success' ? 'success' : 'info'} alert-dismissible fade show message-alert`;
    messageDiv.innerHTML = `
        ${message}
        <button type="button" class="btn-close" data-bs-dismiss="alert"></button>
    `;
    
    // Insert at top of main container
    const mainContainer = document.querySelector('.main-container');
    mainContainer.insertBefore(messageDiv, mainContainer.firstChild);
    
    // Auto remove after 5 seconds
    setTimeout(() => {
        if (messageDiv.parentNode) {
            messageDiv.remove();
        }
    }, 5000);
}