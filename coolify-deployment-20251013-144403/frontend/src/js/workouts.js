// Global variables
let selectedWorkoutType = '';
let selectedCategories = [];
let selectedComplaints = [];
let selectedInjuries = [];
let clientData = {};
let workoutData = {};
let metabolismData = {};
let injuryData = {};

// Load data on page load
document.addEventListener('DOMContentLoaded', function() {
    loadWorkoutData();
    loadMetabolismData();
    loadInjuryData();
});

// Enhanced data loading with error handling
async function loadWorkoutData() {
    try {
        const response = await fetch('/api/exercises');
        
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        // Parse JSON directly from backend response
        workoutData = await response.json();
        
        populateWorkoutCategories();
        
        // Log successful loading
        if (window.errorHandler) {
            window.errorHandler.logInfo('Workout data loaded successfully');
        }
        
    } catch (error) {
        console.error('Error loading workout data:', error);
        
        // Log error if error handler is available
        if (window.errorHandler) {
            window.errorHandler.logError(error, 'loadWorkoutData');
        }
        
        // Enhanced fallback data
        workoutData = {
            en: {
                categories: {
                    'upper_body': 'Upper Body',
                    'lower_body': 'Lower Body',
                    'cardio': 'Cardio',
                    'strength': 'Strength Training',
                    'flexibility': 'Flexibility',
                    'core': 'Core Training',
                    'functional': 'Functional Training',
                    'rehabilitation': 'Rehabilitation'
                }
            },
            ar: {
                categories: {
                    'upper_body': 'تمارين الجزء العلوي',
                    'lower_body': 'تمارين الجزء السفلي',
                    'cardio': 'تمارين القلب',
                    'strength': 'تمارين القوة',
                    'flexibility': 'تمارين المرونة',
                    'core': 'تمارين البطن والجذع',
                    'functional': 'التمارين الوظيفية',
                    'rehabilitation': 'تمارين إعادة التأهيل'
                }
            }
        };
        
        populateWorkoutCategories();
        
        // Show user-friendly error message
        if (window.errorHandler) {
            window.errorHandler.showMessage('تم تحميل البيانات الاحتياطية للتمارين - Fallback workout data loaded', 'warning');
        }
    }
}

// Load metabolism data for complaints
// TODO: Integrate with backend API for metabolism data
async function loadMetabolismData() {
    try {
        const response = await fetch('../metabolism.js');
        const text = await response.text();
        
        // Extract JSON from JS file - look for the main object
        const jsonMatch = text.match(/\{[\s\S]*\}/);
        if (jsonMatch) {
            metabolismData = JSON.parse(jsonMatch[0]);
            console.log('Loaded metabolism data:', metabolismData);
            populateComplaints();
        }
    } catch (error) {
        console.error('Error loading metabolism data:', error);
        // Enhanced fallback complaints data based on metabolism sections
        metabolismData = {
            metabolism_guide: {
                title: {
                    en: "Comprehensive Guide to Metabolism: Nutrition, Exercise, and Physiological Factors",
                    ar: "دليل شامل لعملية الأيض: التغذية والتمارين والعوامل الفسيولوجية"
                },
                sections: [
                    { 
                        section_id: 'metabolism_during_eating', 
                        title: { ar: 'عملية الأيض أثناء الأكل', en: 'Metabolism During Eating' },
                        content: {
                            ar: {
                                important_notes: ['تأثير الطعام الحراري يمثل 10٪ من إجمالي الطاقة اليومية'],
                                practice_and_experiments: ['الوجبات المختلطة تحسن الاستجابة الأيضية'],
                                analysis_rules: ['البروتين هو العامل الرئيسي لتحديد حجم TEF']
                            }
                        }
                    },
                    { 
                        section_id: 'metabolism_fasting', 
                        title: { ar: 'عملية الأيض أثناء الصيام', en: 'Metabolism During Fasting' },
                        content: {
                            ar: {
                                important_notes: ['الصيام يؤدي إلى تحول أيضي من الجلوكوز إلى الدهون'],
                                practice_and_experiments: ['الصيام المتقطع يزيد معدل الأيض بنسبة 3-14٪'],
                                analysis_rules: ['الصيام المطول يمكن أن يقلل من معدل الأيض']
                            }
                        }
                    },
                    { 
                        section_id: 'metabolism_stress', 
                        title: { ar: 'عملية الأيض والتوتر', en: 'Metabolism and Stress' },
                        content: {
                            ar: {
                                important_notes: ['التوتر المزمن يرفع الكورتيزول ويؤثر على الأيض'],
                                practice_and_experiments: ['تقنيات إدارة التوتر تحسن الصحة الأيضية'],
                                analysis_rules: ['جودة النوم تؤثر بشكل كبير على أيض التوتر']
                            }
                        }
                    },
                    { 
                        section_id: 'metabolism_hormones', 
                        title: { ar: 'عملية الأيض والتوازن الهرموني', en: 'Metabolism and Hormonal Balance' },
                        content: {
                            ar: {
                                important_notes: ['الهرمونات تنظم جميع جوانب عملية الأيض'],
                                practice_and_experiments: ['النظام الغذائي المتوازن يدعم التوازن الهرموني'],
                                analysis_rules: ['الأنسولين والجلوكاجون لهما تأثيرات متضادة على الأيض']
                            }
                        }
                    }
                ]
            }
        };
        populateComplaints();
    }
}

// Load injury data
// TODO: Integrate with backend API for injury data
async function loadInjuryData() {
    try {
        const response = await fetch('../injury easy trae json/injury 1 easy.js');
        const text = await response.text();
        
        // Extract all injury JSON objects from the file
        const injuryMatches = text.match(/```json\s*\n\s*\{[\s\S]*?\n\s*\}/g);
        
        if (injuryMatches && injuryMatches.length > 0) {
            const injuries = [];
            
            injuryMatches.forEach((match, index) => {
                try {
                    // Clean the match and extract JSON
                    const cleanMatch = match.replace(/```json\s*\n/, '').replace(/\n\s*\}\s*$/, '}');
                    const injuryObj = JSON.parse(cleanMatch);
                    
                    // Create injury entry with proper structure
                    const injury = {
                        id: `injury_${index + 1}`,
                        name: {
                            ar: injuryObj.title?.arabic || `إصابة ${index + 1}`,
                            en: injuryObj.title?.english || `Injury ${index + 1}`
                        },
                        description: {
                            ar: injuryObj.description?.arabic || '',
                            en: injuryObj.description?.english || ''
                        },
                        management_plan: {
                            ar: injuryObj.management_plan?.arabic || '',
                            en: injuryObj.management_plan?.english || ''
                        },
                        supplements: injuryObj.supplements || [],
                        medications: injuryObj.medications || [],
                        plants_and_herbs: injuryObj.plants_and_herbs || [],
                        recipes: injuryObj.recipes || [],
                        gym_tips: {
                            ar: injuryObj.gym_tips?.arabic || injuryObj.gym_tips || '',
                            en: injuryObj.gym_tips?.english || injuryObj.gym_tips || ''
                        },
                        monthly_workout_plan: injuryObj.monthly_workout_plan || [],
                        disclaimer: {
                            ar: injuryObj.disclaimer?.arabic || '',
                            en: injuryObj.disclaimer?.english || ''
                        }
                    };
                    
                    injuries.push(injury);
                } catch (parseError) {
                    console.warn(`Error parsing injury ${index + 1}:`, parseError);
                }
            });
            
            injuryData = { injuries };
            console.log('Loaded injuries:', injuries.length);
            populateInjuries();
        } else {
            throw new Error('No injury data found in file');
        }
    } catch (error) {
        console.error('Error loading injury data:', error);
        // Enhanced fallback injury data
        injuryData = {
            injuries: [
                { 
                    id: 'gym_sprain', 
                    name: { ar: 'التواء الجيم', en: 'Gym Sprain' },
                    description: { ar: 'إصابة شائعة في الجيم', en: 'Common gym injury' }
                },
                { 
                    id: 'back_pain', 
                    name: { ar: 'آلام الظهر', en: 'Back Pain' },
                    description: { ar: 'آلام في منطقة الظهر', en: 'Pain in the back area' }
                },
                { 
                    id: 'knee_pain', 
                    name: { ar: 'آلام الركبة', en: 'Knee Pain' },
                    description: { ar: 'آلام في منطقة الركبة', en: 'Pain in the knee area' }
                },
                { 
                    id: 'shoulder_pain', 
                    name: { ar: 'آلام الكتف', en: 'Shoulder Pain' },
                    description: { ar: 'آلام في منطقة الكتف', en: 'Pain in the shoulder area' }
                },
                { 
                    id: 'neck_pain', 
                    name: { ar: 'آلام الرقبة', en: 'Neck Pain' },
                    description: { ar: 'آلام في منطقة الرقبة', en: 'Pain in the neck area' }
                },
                { 
                    id: 'ankle_sprain', 
                    name: { ar: 'التواء الكاحل', en: 'Ankle Sprain' },
                    description: { ar: 'التواء في منطقة الكاحل', en: 'Sprain in the ankle area' }
                },
                { 
                    id: 'wrist_pain', 
                    name: { ar: 'آلام المعصم', en: 'Wrist Pain' },
                    description: { ar: 'آلام في منطقة المعصم', en: 'Pain in the wrist area' }
                },
                { 
                    id: 'muscle_strain', 
                    name: { ar: 'شد عضلي', en: 'Muscle Strain' },
                    description: { ar: 'شد في العضلات', en: 'Strain in muscles' }
                }
            ]
        };
        populateInjuries();
    }
}

// Select workout type with enhanced animation
function selectWorkoutType(type) {
    selectedWorkoutType = type;
    
    // Update UI with animation
    document.querySelectorAll('.icon-option').forEach(option => {
        option.classList.remove('selected');
    });
    
    const selectedOption = document.querySelector(`[data-type="${type}"]`);
    selectedOption.classList.add('selected');
    
    // Add feedback animation
    selectedOption.style.animation = 'none';
    setTimeout(() => {
        selectedOption.style.animation = 'selectedPulse 2s infinite';
    }, 10);
    
    // Update workout categories based on type
    updateWorkoutCategoriesForType(type);
    
    // Show success message
    showMessage(`تم اختيار ${type === 'gym' ? 'تمارين الجيم' : 'التمارين المنزلية'}`, 'success');
    
    console.log('Selected workout type:', type);
}

// Update workout categories based on selected type
function updateWorkoutCategoriesForType(type) {
    const container = document.getElementById('workoutCategories');
    if (!container) return;
    
    // Enhanced categories with type-specific filtering
    const allCategories = [
        { id: 'upper_body', name: 'تمارين الجزء العلوي', icon: 'fas fa-hand-rock', gym: true, home: true },
        { id: 'lower_body', name: 'تمارين الجزء السفلي', icon: 'fas fa-running', gym: true, home: true },
        { id: 'cardio', name: 'تمارين القلب', icon: 'fas fa-heartbeat', gym: true, home: true },
        { id: 'strength', name: 'تمارين القوة', icon: 'fas fa-dumbbell', gym: true, home: false },
        { id: 'flexibility', name: 'تمارين المرونة', icon: 'fas fa-child', gym: true, home: true },
        { id: 'core', name: 'تمارين البطن والجذع', icon: 'fas fa-circle', gym: true, home: true },
        { id: 'functional', name: 'التمارين الوظيفية', icon: 'fas fa-cogs', gym: true, home: true },
        { id: 'rehabilitation', name: 'تمارين إعادة التأهيل', icon: 'fas fa-medkit', gym: false, home: true },
        { id: 'bodyweight', name: 'تمارين وزن الجسم', icon: 'fas fa-male', gym: false, home: true },
        { id: 'weights', name: 'تمارين الأوزان', icon: 'fas fa-weight-hanging', gym: true, home: false }
    ];
    
    // Filter categories based on workout type
    const filteredCategories = allCategories.filter(category => {
        return type === 'gym' ? category.gym : category.home;
    });
    
    // Update container with filtered categories
    container.innerHTML = filteredCategories.map(category => `
        <div class="col-md-3 mb-3">
            <div class="workout-card p-3 text-center" onclick="toggleCategory('${category.id}')" data-category="${category.id}">
                <i class="${category.icon} fa-2x mb-2 text-primary"></i>
                <h6>${category.name}</h6>
            </div>
        </div>
    `).join('');
    
    // Clear previously selected categories that are no longer available
    selectedCategories = selectedCategories.filter(categoryId => 
        filteredCategories.some(cat => cat.id === categoryId)
    );
}

// Populate workout categories
function populateWorkoutCategories() {
    const container = document.getElementById('workoutCategories');
    const categories = [
        { id: 'upper_body', name: 'تمارين الجزء العلوي', icon: 'fas fa-hand-rock' },
        { id: 'lower_body', name: 'تمارين الجزء السفلي', icon: 'fas fa-running' },
        { id: 'cardio', name: 'تمارين القلب', icon: 'fas fa-heartbeat' },
        { id: 'strength', name: 'تمارين القوة', icon: 'fas fa-dumbbell' },
        { id: 'flexibility', name: 'تمارين المرونة', icon: 'fas fa-child' },
        { id: 'core', name: 'تمارين البطن والجذع', icon: 'fas fa-circle' },
        { id: 'functional', name: 'التمارين الوظيفية', icon: 'fas fa-cogs' },
        { id: 'rehabilitation', name: 'تمارين إعادة التأهيل', icon: 'fas fa-medkit' }
    ];
    
    container.innerHTML = categories.map(category => `
        <div class="col-md-3 mb-3">
            <div class="workout-card p-3 text-center" onclick="toggleCategory('${category.id}')" data-category="${category.id}">
                <i class="${category.icon} fa-2x mb-2 text-primary"></i>
                <h6>${category.name}</h6>
            </div>
        </div>
    `).join('');
}

// Toggle category selection
function toggleCategory(categoryId) {
    const card = document.querySelector(`[data-category="${categoryId}"]`);
    
    if (selectedCategories.includes(categoryId)) {
        selectedCategories = selectedCategories.filter(id => id !== categoryId);
        card.classList.remove('selected');
    } else {
        selectedCategories.push(categoryId);
        card.classList.add('selected');
    }
}

// Populate complaints
function populateComplaints() {
    const container = document.getElementById('complaintsSection');
    if (!container) return;
    
    // Use loaded metabolism data or fallback
    let complaints = [];
    
    if (metabolismData?.metabolism_guide?.sections) {
        complaints = metabolismData.metabolism_guide.sections.map(section => ({
            id: section.section_id,
            name: section.title?.ar || section.title?.en || section.section_id,
            description: section.content?.ar?.important_notes?.[0] || section.content?.en?.important_notes?.[0] || ''
        }));
    } else {
        // Fallback complaints
        complaints = [
            { id: 'metabolism_during_eating', name: 'مشاكل الأيض أثناء الأكل', description: 'صعوبة في هضم الطعام وبطء الأيض' },
            { id: 'metabolism_fasting', name: 'مشاكل الأيض أثناء الصيام', description: 'انخفاض الطاقة وبطء حرق الدهون' },
            { id: 'metabolism_stress', name: 'تأثير التوتر على الأيض', description: 'التوتر المزمن يؤثر على عملية الأيض' },
            { id: 'metabolism_hormones', name: 'خلل الهرمونات والأيض', description: 'اضطرابات هرمونية تؤثر على الأيض' },
            { id: 'metabolism_sleep', name: 'مشاكل النوم والأيض', description: 'قلة النوم تؤثر على معدل الأيض' },
            { id: 'metabolism_aging', name: 'بطء الأيض مع التقدم في العمر', description: 'انخفاض معدل الأيض مع التقدم في السن' },
            { id: 'metabolism_thyroid', name: 'مشاكل الغدة الدرقية', description: 'اضطرابات الغدة الدرقية وتأثيرها على الأيض' },
            { id: 'metabolism_diabetes', name: 'مشاكل السكري والأيض', description: 'اضطرابات السكر وتأثيرها على الأيض' }
        ];
    }
    
    console.log('Populating complaints:', complaints.length);
    
    container.innerHTML = `
        <div class="row">
            ${complaints.map(complaint => `
                <div class="col-md-4 mb-2">
                    <div class="complaint-item" onclick="toggleComplaint('${complaint.id}')" data-complaint="${complaint.id}" title="${complaint.description || ''}">
                        <i class="fas fa-exclamation-circle me-2"></i>
                        <span class="complaint-name">${complaint.name}</span>
                        ${complaint.description ? `<small class="d-block text-muted mt-1">${complaint.description.substring(0, 50)}${complaint.description.length > 50 ? '...' : ''}</small>` : ''}
                    </div>
                </div>
            `).join('')}
        </div>
    `;
}

// Toggle complaint selection
function toggleComplaint(complaintId) {
    const item = document.querySelector(`[data-complaint="${complaintId}"]`);
    
    if (selectedComplaints.includes(complaintId)) {
        selectedComplaints = selectedComplaints.filter(id => id !== complaintId);
        item.classList.remove('selected');
    } else {
        selectedComplaints.push(complaintId);
        item.classList.add('selected');
    }
}

// Populate injuries
function populateInjuries() {
    const container = document.getElementById('injuriesSection');
    if (!container) return;
    
    // Use loaded injury data or fallback
    const injuries = injuryData?.injuries || [
        { id: 'back_pain', name: { ar: 'آلام الظهر' } },
        { id: 'knee_pain', name: { ar: 'آلام الركبة' } },
        { id: 'shoulder_pain', name: { ar: 'آلام الكتف' } }
    ];
    
    container.innerHTML = `
        <div class="row">
            ${injuries.map(injury => `
                <div class="col-md-4 mb-2">
                    <div class="injury-item" onclick="toggleInjury('${injury.id}')" data-injury="${injury.id}" title="${injury.description?.ar || injury.description || ''}">
                        <i class="fas fa-band-aid me-2"></i>
                        <span class="injury-name">${injury.name?.ar || injury.name}</span>
                        ${injury.description?.ar ? `<small class="d-block text-muted mt-1">${injury.description.ar.substring(0, 50)}...</small>` : ''}
                    </div>
                </div>
            `).join('')}
        </div>
    `;
    
    console.log('Populated injuries:', injuries.length);
}

// Toggle injury selection
function toggleInjury(injuryId) {
    const item = document.querySelector(`[data-injury="${injuryId}"]`);
    
    if (selectedInjuries.includes(injuryId)) {
        selectedInjuries = selectedInjuries.filter(id => id !== injuryId);
        item.classList.remove('selected');
    } else {
        selectedInjuries.push(injuryId);
        item.classList.add('selected');
    }
}

// Collect client data
function collectClientData() {
    clientData = {
        name: document.getElementById('clientName').value,
        gender: document.getElementById('clientGender').value,
        weight: parseFloat(document.getElementById('clientWeight').value),
        height: parseFloat(document.getElementById('clientHeight').value),
        age: parseInt(document.getElementById('clientAge').value),
        activityLevel: document.getElementById('activityLevel').value,
        workoutGoal: document.getElementById('workoutGoal').value
    };
    
    return clientData;
}

// Enhanced validation using error handler
function validateFormData() {
    const data = collectClientData();
    
    // Use enhanced error handler validation if available
    if (window.errorHandler && typeof window.errorHandler.validateClientData === 'function') {
        const errors = window.errorHandler.validateClientData({
            name: data.name,
            age: data.age,
            weight: data.weight,
            height: data.height,
            gender: data.gender,
            activityLevel: data.activityLevel,
            workoutGoal: data.workoutGoal
        });
        
        // Add workout-specific validations
        if (!selectedWorkoutType) {
            errors.push('الرجاء اختيار نوع التمارين (جيم أو منزلية) - Please select workout type (gym or home)');
        }
        
        if (selectedCategories.length === 0) {
            errors.push('الرجاء اختيار نوع واحد على الأقل من التمارين - Please select at least one workout category');
        }
        
        // Show validation errors using enhanced error handler
        if (errors.length > 0) {
            if (typeof window.errorHandler.showValidationErrors === 'function') {
                window.errorHandler.showValidationErrors(errors, 'workoutValidationErrors');
            } else {
                alert(errors.join('\n'));
            }
            return false;
        }
        
        // Clear any existing validation errors
        const errorContainer = document.getElementById('workoutValidationErrors');
        if (errorContainer) {
            errorContainer.innerHTML = '';
            errorContainer.style.display = 'none';
        }
        
        return true;
    }
    
    // Fallback validation if error handler is not available
    if (!data.name || !data.gender || !data.weight || !data.height || !data.age || !data.activityLevel || !data.workoutGoal) {
        alert('الرجاء ملء جميع البيانات المطلوبة');
        return false;
    }
    
    if (!selectedWorkoutType) {
        alert('الرجاء اختيار نوع التمارين (جيم أو منزلية)');
        return false;
    }
    
    if (selectedCategories.length === 0) {
        alert('الرجاء اختيار نوع واحد على الأقل من التمارين');
        return false;
    }
    
    return true;
}

// Generate workout plan
function generateWorkoutPlan() {
    if (!validateFormData()) {
        return;
    }
    
    const data = collectClientData();
    
    // Show loading
    const resultsContainer = document.getElementById('workoutPlanResults');
    resultsContainer.style.display = 'block';
    resultsContainer.innerHTML = `
        <div class="text-center py-5">
            <div class="spinner-border text-primary" role="status">
                <span class="visually-hidden">جاري التحميل...</span>
            </div>
            <p class="mt-3">جاري توليد خطة التمارين الشخصية...</p>
        </div>
    `;
    
    // Simulate processing time
    setTimeout(() => {
        const workoutPlan = createAdvancedWorkoutPlan(data);
        displayAdvancedResults(workoutPlan);
    }, 2000);
}

// Create advanced workout plan with detailed exercises
function createAdvancedWorkoutPlan(clientData) {
    const plan = {
        clientInfo: clientData,
        workoutType: selectedWorkoutType,
        categories: selectedCategories,
        complaints: selectedComplaints,
        injuries: selectedInjuries,
        weeklyPlan: generateWeeklyWorkoutPlan(clientData),
        exercises: [],
        nutritionAdvice: [],
        supplements: [],
        injuryTreatments: [],
        commonMistakes: generateCommonMistakes(selectedCategories),
        alternatives: generateExerciseAlternatives(selectedCategories, selectedInjuries),
        medicalDisclaimer: "هذا الموقع لتحديثك بالمعلومات وتوجيهك للنصائح المفيدة بهدف التعليم والوعي ولا يغني عن زيارة الطبيب"
    };
    
    // Generate exercises based on selected categories and workout type
    selectedCategories.forEach(category => {
        const categoryExercises = getAdvancedExercisesForCategory(category, clientData);
        plan.exercises.push(...categoryExercises);
    });
    
    // Add injury-specific treatments
    selectedInjuries.forEach(injuryId => {
        const treatment = generateInjuryTreatment(injuryId);
        if (treatment) {
            plan.injuryTreatments.push(treatment);
        }
    });
    
    // Add complaint-specific nutrition and supplements
    selectedComplaints.forEach(complaintId => {
        const nutritionAdvice = generateNutritionAdviceForComplaint(complaintId, clientData);
        if (nutritionAdvice) {
            plan.nutritionAdvice.push(nutritionAdvice);
        }
    });
    
    // Generate supplements
    const supplements = generateSupplementAdvice(selectedComplaints, clientData);
    plan.supplements = supplements;
    
    return plan;
}

// Get advanced exercises for specific category
function getAdvancedExercisesForCategory(category, clientData) {
    const exercises = generateExercisesForCategory(category, clientData.workoutGoal, selectedWorkoutType);
    
    // Enhance exercises with detailed information
    return exercises.map(exercise => ({
        ...exercise,
        category: category,
        sets: calculateAdvancedSets(exercise, clientData),
        reps: calculateAdvancedReps(exercise, clientData),
        rest: calculateAdvancedRest(exercise, clientData),
        progressionTips: generateProgressionTips(exercise, clientData),
        safetyNotes: generateSafetyNotes(exercise, selectedInjuries)
    }));
}

// Calculate advanced sets based on client data and goals
function calculateAdvancedSets(exercise, clientData) {
    let baseSets = exercise.sets || 3;
    
    // Adjust based on experience level (derived from activity level)
    if (clientData.activityLevel === 'sedentary') {
        baseSets = Math.max(2, baseSets - 1);
    } else if (clientData.activityLevel === 'very_active') {
        baseSets = baseSets + 1;
    }
    
    // Adjust based on goal
    switch (clientData.workoutGoal) {
        case 'strength':
            return Math.max(4, baseSets + 1);
        case 'muscle_gain':
            return Math.max(3, baseSets);
        case 'weight_loss':
        case 'endurance':
            return Math.max(2, baseSets);
        default:
            return baseSets;
    }
}

// Calculate advanced reps based on client data and goals
function calculateAdvancedReps(exercise, clientData) {
    const baseReps = exercise.reps || '10-12';
    
    switch (clientData.workoutGoal) {
        case 'strength':
            return '4-6';
        case 'muscle_gain':
            return '8-12';
        case 'weight_loss':
        case 'endurance':
            return '12-15';
        case 'general_fitness':
            return '10-12';
        default:
            return baseReps;
    }
}

// Calculate advanced rest time based on exercise and goals
function calculateAdvancedRest(exercise, clientData) {
    switch (clientData.workoutGoal) {
        case 'strength':
            return '2-3 دقائق';
        case 'muscle_gain':
            return '90-120 ثانية';
        case 'weight_loss':
        case 'endurance':
            return '30-60 ثانية';
        default:
            return exercise.rest || '60-90 ثانية';
    }
}

// Generate progression tips
function generateProgressionTips(exercise, clientData) {
    const tips = [];
    
    switch (clientData.workoutGoal) {
        case 'strength':
            tips.push('زد الوزن تدريجياً كل أسبوع بنسبة 2.5-5%');
            tips.push('ركز على الشكل الصحيح قبل زيادة الوزن');
            break;
        case 'muscle_gain':
            tips.push('زد التكرارات أو الوزن عندما تصبح التمارين سهلة');
            tips.push('تأكد من الشعور بالتعب في آخر تكرارين');
            break;
        case 'weight_loss':
            tips.push('قلل فترات الراحة تدريجياً لزيادة حرق السعرات');
            tips.push('أضف تمارين مركبة لحرق المزيد من السعرات');
            break;
    }
    
    return tips;
}

// Generate safety notes based on injuries
function generateSafetyNotes(exercise, injuries) {
    const notes = [];
    
    if (injuries.includes('back_pain')) {
        notes.push('تجنب الانحناء المفرط للظهر');
        notes.push('حافظ على استقامة العمود الفقري');
    }
    
    if (injuries.includes('knee_pain')) {
        notes.push('تجنب النزول العميق إذا شعرت بألم في الركبة');
        notes.push('تأكد من عدم تجاوز الركبتين لأصابع القدمين');
    }
    
    if (injuries.includes('shoulder_pain')) {
        notes.push('تجنب الحركات فوق مستوى الرأس');
        notes.push('ابدأ بأوزان خفيفة وزد تدريجياً');
    }
    
    return notes;
}

// Generate weekly workout plan
function generateWeeklyWorkoutPlan(clientData) {
    const daysOfWeek = ['الأحد', 'الاثنين', 'الثلاثاء', 'الأربعاء', 'الخميس', 'الجمعة', 'السبت'];
    const weeklyPlan = [];
    
    // Determine workout frequency based on activity level and goal
    let workoutDays = 3; // Default
    if (clientData.activityLevel === 'very_active' || clientData.workoutGoal === 'muscle_gain') {
        workoutDays = 5;
    } else if (clientData.activityLevel === 'moderately_active') {
        workoutDays = 4;
    }
    
    // Create workout schedule
    const workoutSchedule = generateWorkoutSchedule(workoutDays, selectedCategories);
    
    daysOfWeek.forEach((day, index) => {
        const dayPlan = {
            day: day,
            dayNumber: index + 1,
            isWorkoutDay: workoutSchedule[index] !== null,
            focus: workoutSchedule[index],
            exercises: [],
            duration: '45-60 دقيقة',
            intensity: getIntensityForDay(index, clientData.workoutGoal)
        };
        
        if (dayPlan.isWorkoutDay) {
            dayPlan.exercises = generateDayExercises(dayPlan.focus, clientData);
        } else {
            dayPlan.restActivity = getRestDayActivity(clientData.activityLevel);
        }
        
        weeklyPlan.push(dayPlan);
    });
    
    return weeklyPlan;
}

// Generate workout schedule for the week
function generateWorkoutSchedule(workoutDays, categories) {
    const schedule = [null, null, null, null, null, null, null]; // 7 days
    
    if (workoutDays === 3) {
        schedule[0] = categories[0] || 'upper_body'; // Sunday
        schedule[2] = categories[1] || 'lower_body'; // Tuesday
        schedule[4] = categories[2] || 'full_body'; // Thursday
    } else if (workoutDays === 4) {
        schedule[0] = 'upper_body'; // Sunday
        schedule[1] = 'lower_body'; // Monday
        schedule[3] = 'upper_body'; // Wednesday
        schedule[5] = 'lower_body'; // Friday
    } else if (workoutDays === 5) {
        schedule[0] = 'chest_triceps'; // Sunday
        schedule[1] = 'back_biceps'; // Monday
        schedule[2] = 'legs'; // Tuesday
        schedule[4] = 'shoulders'; // Thursday
        schedule[5] = 'full_body'; // Friday
    }
    
    return schedule;
}

// Generate exercises for specific day
function generateDayExercises(focus, clientData) {
    const exercises = [];
    
    switch (focus) {
        case 'upper_body':
        case 'chest_triceps':
            exercises.push(
                { name: 'تمرين الضغط', sets: '3-4', reps: '8-12', rest: '60-90 ثانية', muscle: 'الصدر والذراعين' },
                { name: 'تمرين العقلة', sets: '3', reps: '5-10', rest: '90 ثانية', muscle: 'الظهر والذراعين' },
                { name: 'تمرين الديبس', sets: '3', reps: '8-12', rest: '60 ثانية', muscle: 'الصدر والذراعين' }
            );
            break;
        case 'lower_body':
        case 'legs':
            exercises.push(
                { name: 'القرفصاء', sets: '4', reps: '10-15', rest: '90 ثانية', muscle: 'الفخذين والمؤخرة' },
                { name: 'الطعنات', sets: '3', reps: '10 لكل رجل', rest: '60 ثانية', muscle: 'الفخذين والمؤخرة' },
                { name: 'رفع الساق الخلفية', sets: '3', reps: '12-15', rest: '45 ثانية', muscle: 'المؤخرة' }
            );
            break;
        case 'full_body':
            exercises.push(
                { name: 'البيربي', sets: '3', reps: '8-12', rest: '90 ثانية', muscle: 'الجسم كامل' },
                { name: 'تمرين الجبل', sets: '3', reps: '20 ثانية', rest: '60 ثانية', muscle: 'الجسم كامل' },
                { name: 'القفز مع فتح الذراعين', sets: '3', reps: '15-20', rest: '45 ثانية', muscle: 'القلب والأوعية' }
            );
            break;
    }
    
    return exercises.map(ex => ({
        ...ex,
        sets: calculateAdvancedSets(ex, clientData),
        reps: calculateAdvancedReps(ex, clientData),
        rest: calculateAdvancedRest(ex, clientData)
    }));
}

// Get intensity for specific day
function getIntensityForDay(dayIndex, goal) {
    const intensities = ['متوسطة', 'عالية', 'متوسطة', 'عالية', 'متوسطة', 'منخفضة', 'راحة'];
    
    if (goal === 'strength') {
        return ['عالية', 'متوسطة', 'عالية', 'راحة', 'عالية', 'متوسطة', 'راحة'][dayIndex];
    }
    
    return intensities[dayIndex];
}

// Get rest day activity
function getRestDayActivity(activityLevel) {
    const activities = {
        'sedentary': 'مشي خفيف لمدة 20-30 دقيقة',
        'lightly_active': 'يوجا أو تمدد لمدة 30 دقيقة',
        'moderately_active': 'مشي سريع أو سباحة خفيفة',
        'very_active': 'نشاط خفيف أو تمدد عميق'
    };
    
    return activities[activityLevel] || 'راحة تامة';
}

// Generate common mistakes for selected categories
function generateCommonMistakes(categories) {
    const mistakes = {
        'chest': [
            {
                mistake: 'عدم النزول بالكامل في تمرين الضغط',
                correction: 'انزل حتى يلامس صدرك الأرض تقريباً للحصول على أقصى استفادة',
                consequence: 'تقليل فعالية التمرين وعدم تطوير القوة بشكل كامل'
            },
            {
                mistake: 'رفع الوزن بسرعة مفرطة',
                correction: 'تحكم في الحركة واجعلها بطيئة ومنضبطة',
                consequence: 'زيادة خطر الإصابة وتقليل تأثير التمرين'
            }
        ],
        'back': [
            {
                mistake: 'استخدام الزخم بدلاً من العضلات',
                correction: 'ركز على سحب الوزن بعضلات الظهر وليس بالزخم',
                consequence: 'عدم تطوير عضلات الظهر بشكل صحيح'
            },
            {
                mistake: 'عدم ضغط لوحي الكتف معاً',
                correction: 'اضغط لوحي الكتف معاً في نهاية الحركة',
                consequence: 'تقليل تفعيل عضلات الظهر الوسطى'
            }
        ],
        'legs': [
            {
                mistake: 'عدم النزول بعمق كافٍ في القرفصاء',
                correction: 'انزل حتى تصبح فخذيك موازية للأرض على الأقل',
                consequence: 'عدم تطوير عضلات الفخذ والمؤخرة بشكل كامل'
            },
            {
                mistake: 'تجاوز الركبتين لأصابع القدمين',
                correction: 'حافظ على الركبتين خلف أصابع القدمين',
                consequence: 'زيادة الضغط على مفاصل الركبة وخطر الإصابة'
            }
        ],
        'shoulders': [
            {
                mistake: 'رفع الأوزان فوق الرأس بشكل خاطئ',
                correction: 'حافظ على استقامة الظهر وتجنب التقوس المفرط',
                consequence: 'إجهاد أسفل الظهر وزيادة خطر الإصابة'
            }
        ]
    };
    
    const selectedMistakes = [];
    categories.forEach(category => {
        if (mistakes[category]) {
            selectedMistakes.push(...mistakes[category]);
        }
    });
    
    // Add general mistakes if no specific ones found
    if (selectedMistakes.length === 0) {
        selectedMistakes.push(
            {
                mistake: 'عدم الإحماء قبل التمرين',
                correction: 'قم بإحماء لمدة 5-10 دقائق قبل بدء التمارين',
                consequence: 'زيادة خطر الإصابة وتقليل الأداء'
            },
            {
                mistake: 'إهمال فترات الراحة',
                correction: 'خذ راحة كافية بين المجموعات (60-120 ثانية)',
                consequence: 'تقليل الأداء في المجموعات التالية'
            }
        );
    }
    
    return selectedMistakes;
}

// Generate exercise alternatives based on categories and injuries
function generateExerciseAlternatives(categories, injuries) {
    const alternatives = {
        'chest': {
            original: 'تمرين الضغط العادي',
            alternatives: [
                {
                    name: 'تمرين الضغط على الركبتين',
                    difficulty: 'أسهل',
                    reason: 'للمبتدئين أو من لديهم ضعف في القوة'
                },
                {
                    name: 'تمرين الضغط المائل',
                    difficulty: 'أسهل',
                    reason: 'تقليل الحمل على الجزء العلوي من الجسم'
                },
                {
                    name: 'تمرين الضغط بقدم واحدة',
                    difficulty: 'أصعب',
                    reason: 'لزيادة التحدي وتطوير التوازن'
                }
            ]
        },
        'back': {
            original: 'تمرين العقلة',
            alternatives: [
                {
                    name: 'تمرين العقلة بالمساعدة',
                    difficulty: 'أسهل',
                    reason: 'للمبتدئين الذين لا يستطيعون رفع وزن الجسم'
                },
                {
                    name: 'تمرين السحب الأفقي',
                    difficulty: 'أسهل',
                    reason: 'بديل آمن لمن لديهم مشاكل في الكتف'
                }
            ]
        },
        'legs': {
            original: 'القرفصاء العميق',
            alternatives: [
                {
                    name: 'القرفصاء الجزئي',
                    difficulty: 'أسهل',
                    reason: 'لمن لديهم مشاكل في الركبة أو مرونة محدودة'
                },
                {
                    name: 'القرفصاء على كرسي',
                    difficulty: 'أسهل',
                    reason: 'للمبتدئين لتعلم الحركة الصحيحة'
                },
                {
                    name: 'القرفصاء بقدم واحدة',
                    difficulty: 'أصعب',
                    reason: 'لتطوير التوازن والقوة الوظيفية'
                }
            ]
        }
    };
    
    const selectedAlternatives = [];
    
    categories.forEach(category => {
        if (alternatives[category]) {
            let categoryAlternatives = alternatives[category];
            
            // Filter alternatives based on injuries
            if (injuries.includes('knee_pain') && category === 'legs') {
                categoryAlternatives.alternatives = categoryAlternatives.alternatives.filter(
                    alt => alt.difficulty === 'أسهل'
                );
            }
            
            if (injuries.includes('shoulder_pain') && (category === 'chest' || category === 'back')) {
                categoryAlternatives.alternatives = categoryAlternatives.alternatives.filter(
                    alt => alt.reason.includes('آمن') || alt.difficulty === 'أسهل'
                );
            }
            
            selectedAlternatives.push(categoryAlternatives);
        }
    });
    
    return selectedAlternatives;
}

// Generate injury treatment plan
function generateInjuryTreatment(injuryId) {
    const injuryTreatments = {
        'back_pain': {
            name: 'آلام الظهر',
            description: 'خطة علاج شاملة لآلام الظهر',
            exercises: [
                {
                    name: 'تمدد القطة والجمل',
                    sets: '2-3',
                    reps: '10-15',
                    instructions: 'على اليدين والركبتين، قوس الظهر لأعلى ثم لأسفل ببطء'
                },
                {
                    name: 'تمرين الجسر',
                    sets: '2-3',
                    reps: '10-12',
                    instructions: 'استلق على الظهر، ارفع الوركين مع شد عضلات المؤخرة'
                }
            ],
            supplements: [
                {
                    name: 'الكركم',
                    dosage: '500-1000 مجم يومياً',
                    benefits: 'مضاد للالتهاب طبيعي'
                },
                {
                    name: 'أوميجا 3',
                    dosage: '1000-2000 مجم يومياً',
                    benefits: 'يقلل الالتهاب ويحسن الشفاء'
                }
            ],
            tips: [
                'تطبيق الثلج لمدة 15-20 دقيقة كل ساعتين في أول 48 ساعة',
                'تجنب الجلوس لفترات طويلة',
                'النوم على جانبك مع وسادة بين الركبتين'
            ]
        },
        'knee_pain': {
            name: 'آلام الركبة',
            description: 'برنامج إعادة تأهيل الركبة',
            exercises: [
                {
                    name: 'تقوية العضلة الرباعية',
                    sets: '2-3',
                    reps: '10-15',
                    instructions: 'اجلس وامدد الساق ببطء مع شد عضلة الفخذ الأمامية'
                },
                {
                    name: 'تمدد أوتار الركبة',
                    sets: '2-3',
                    reps: '30 ثانية',
                    instructions: 'استلق وارفع الساق مع سحبها نحو الصدر'
                }
            ],
            supplements: [
                {
                    name: 'الجلوكوزامين',
                    dosage: '1500 مجم يومياً',
                    benefits: 'يدعم صحة الغضاريف'
                },
                {
                    name: 'الكوندرويتين',
                    dosage: '1200 مجم يومياً',
                    benefits: 'يحافظ على مرونة المفاصل'
                }
            ],
            tips: [
                'تجنب الأنشطة عالية التأثير',
                'استخدم دعامة الركبة عند الحاجة',
                'حافظ على وزن صحي لتقليل الضغط على الركبة'
            ]
        },
        'shoulder_pain': {
            name: 'آلام الكتف',
            description: 'برنامج علاج آلام الكتف',
            exercises: [
                {
                    name: 'دوران الكتف',
                    sets: '2-3',
                    reps: '10 في كل اتجاه',
                    instructions: 'حرك الكتفين في دوائر بطيئة للأمام ثم للخلف'
                },
                {
                    name: 'تمدد الكتف المتقاطع',
                    sets: '2-3',
                    reps: '30 ثانية',
                    instructions: 'اسحب الذراع عبر الصدر مع الضغط بلطف'
                }
            ],
            supplements: [
                {
                    name: 'المغنيسيوم',
                    dosage: '400-600 مجم يومياً',
                    benefits: 'يساعد على استرخاء العضلات'
                },
                {
                    name: 'فيتامين د',
                    dosage: '1000-2000 وحدة دولية يومياً',
                    benefits: 'يدعم صحة العظام والعضلات'
                }
            ],
            tips: [
                'تجنب النوم على الكتف المصاب',
                'استخدم كمادات دافئة قبل التمرين',
                'تجنب الحركات المفاجئة فوق مستوى الرأس'
            ]
        }
    };
    
    return injuryTreatments[injuryId] || null;
}

// Generate nutrition advice for complaints
function generateNutritionAdviceForComplaint(complaintId, clientData) {
    // Use loaded metabolism data if available
    if (metabolismData?.metabolism_guide?.sections) {
        const section = metabolismData.metabolism_guide.sections.find(s => s.section_id === complaintId);
        if (section) {
            const content = section.content?.ar || section.content?.en || {};
            return {
                complaint: section.title?.ar || section.title?.en || complaintId,
                foods: content.practice_and_experiments || [
                    'تناول الأطعمة الطبيعية والمتوازنة',
                    'الإكثار من الخضروات والفواكه',
                    'شرب كمية كافية من الماء'
                ],
                avoid: [
                    'الأطعمة المصنعة والمعالجة',
                    'السكريات المكررة',
                    'الدهون المتحولة',
                    'المشروبات الغازية'
                ],
                tips: content.important_notes || [
                    'اتباع نظام غذائي متوازن',
                    'ممارسة الرياضة بانتظام',
                    'الحصول على قسط كافٍ من النوم'
                ],
                analysis_rules: content.analysis_rules || [],
                references: content.references || []
            };
        }
    }
    
    // Fallback nutrition advice for common complaints
    const fallbackAdvice = {
        'metabolism_during_eating': {
            complaint: 'مشاكل الأيض أثناء الأكل',
            foods: [
                'البروتين الخالي من الدهون لتحفيز الأيض',
                'الأطعمة الغنية بالألياف لتحسين الهضم',
                'الشاي الأخضر لتسريع الأيض',
                'التوابل الحارة مثل الفلفل والزنجبيل'
            ],
            avoid: [
                'الوجبات الكبيرة والثقيلة',
                'السكريات المكررة',
                'الأطعمة المصنعة',
                'تناول الطعام بسرعة'
            ],
            tips: [
                'تناول وجبات صغيرة ومتكررة',
                'امضغ الطعام ببطء وجيداً',
                'اشرب الماء قبل الوجبات بـ 30 دقيقة'
            ]
        },
        'metabolism_fasting': {
            complaint: 'مشاكل الأيض أثناء الصيام',
            foods: [
                'البروتين عالي الجودة في وجبة السحور',
                'الكربوهيدرات المعقدة للطاقة المستدامة',
                'الدهون الصحية مثل المكسرات والأفوكادو',
                'الأطعمة الغنية بالماء والألياف'
            ],
            avoid: [
                'السكريات البسيطة في السحور',
                'الأطعمة المالحة المفرطة',
                'الكافيين المفرط',
                'الإفراط في الطعام عند الإفطار'
            ],
            tips: [
                'تناول السحور متأخراً قدر الإمكان',
                'ابدأ الإفطار بالتمر والماء',
                'تجنب الأنشطة المجهدة أثناء الصيام'
            ]
        },
        'metabolism_stress': {
            complaint: 'تأثير التوتر على الأيض',
            foods: [
                'الأطعمة الغنية بالمغنيسيوم مثل المكسرات',
                'الأسماك الدهنية الغنية بأوميغا 3',
                'الشاي الأخضر والبابونج للاسترخاء',
                'الفواكه والخضروات الملونة'
            ],
            avoid: [
                'الكافيين المفرط',
                'السكريات المكررة',
                'الأطعمة المصنعة',
                'الوجبات السريعة'
            ],
            tips: [
                'تناول وجبات منتظمة لتجنب انخفاض السكر',
                'مارس تقنيات الاسترخاء قبل الأكل',
                'تجنب الأكل العاطفي'
            ]
        }
    };
    
    return fallbackAdvice[complaintId] || {
        complaint: 'نصائح عامة للأيض',
        foods: [
            'البروتين الخالي من الدهون',
            'الخضروات الورقية الخضراء',
            'الفواكه الطازجة',
            'الحبوب الكاملة'
        ],
        avoid: [
            'الأطعمة المصنعة',
            'السكريات المضافة',
            'الدهون المتحولة',
            'المشروبات السكرية'
        ],
        tips: [
            'اتباع نظام غذائي متوازن',
            'ممارسة الرياضة بانتظام',
            'شرب كمية كافية من الماء',
            'الحصول على نوم جيد'
        ]
    };
}

// Generate supplement advice based on complaints and client data
function generateSupplementAdvice(complaints, clientData) {
    const supplements = [];
    const weight = parseFloat(clientData.weight) || 70;
    const age = parseInt(clientData.age) || 30;
    
    // Process each complaint and add specific supplements
    complaints.forEach(complaintId => {
        // Try to get supplement recommendations from metabolism data
        if (metabolismData?.metabolism_guide?.sections) {
            const section = metabolismData.metabolism_guide.sections.find(s => s.section_id === complaintId);
            if (section) {
                const content = section.content?.ar || section.content?.en || {};
                
                // Extract supplement recommendations from practice_and_experiments or analysis_rules
                const practiceItems = content.practice_and_experiments || [];
                const analysisRules = content.analysis_rules || [];
                
                practiceItems.concat(analysisRules).forEach((item, index) => {
                    if (typeof item === 'string' && (item.includes('مكمل') || item.includes('فيتامين') || item.includes('معدن'))) {
                        supplements.push({
                            name: `مكمل متخصص للـ${section.title?.ar || section.title?.en}`,
                            dosage: 'حسب التوجيهات الطبية',
                            benefits: item,
                            timing: 'مع الوجبة',
                            source: 'metabolism_data'
                        });
                    }
                });
            }
        }
        
        // Fallback supplements for specific complaints
        const complaintSupplements = {
            'metabolism_during_eating': [
                {
                    name: 'إنزيمات الهضم',
                    dosage: 'قرص واحد مع كل وجبة رئيسية',
                    benefits: 'يحسن هضم الطعام ويزيد امتصاص العناصر الغذائية',
                    timing: 'مع بداية كل وجبة'
                },
                {
                    name: 'الكروم بيكولينات',
                    dosage: '200 ميكروجرام يومياً',
                    benefits: 'يحسن استقلاب الكربوهيدرات ويقلل الرغبة في السكريات',
                    timing: 'قبل الوجبة الرئيسية بـ 30 دقيقة'
                }
            ],
            'metabolism_fasting': [
                {
                    name: 'الأحماض الأمينية المتشعبة (BCAA)',
                    dosage: '10-15 جرام أثناء فترة الأكل',
                    benefits: 'يحافظ على الكتلة العضلية أثناء الصيام',
                    timing: 'مع أول وجبة بعد الصيام'
                },
                {
                    name: 'المغنيسيوم',
                    dosage: '400-600 مجم يومياً',
                    benefits: 'يقلل التعب ويحسن جودة النوم أثناء الصيام',
                    timing: 'قبل النوم مع السحور'
                }
            ],
            'metabolism_stress': [
                {
                    name: 'الأشواغاندا',
                    dosage: '300-500 مجم مرتين يومياً',
                    benefits: 'يقلل هرمون الكورتيزول ويحسن استجابة الجسم للتوتر',
                    timing: 'صباحاً ومساءً مع الطعام'
                },
                {
                    name: 'فيتامين ب المركب',
                    dosage: 'قرص واحد يومياً',
                    benefits: 'يدعم الجهاز العصبي ويحسن إنتاج الطاقة',
                    timing: 'صباحاً مع الإفطار'
                }
            ],
            'metabolism_hormones': [
                {
                    name: 'فيتامين د3',
                    dosage: '2000-4000 وحدة دولية يومياً',
                    benefits: 'يدعم التوازن الهرموني ويحسن الأيض',
                    timing: 'مع وجبة تحتوي على دهون'
                },
                {
                    name: 'الزنك',
                    dosage: '15-30 مجم يومياً',
                    benefits: 'يدعم إنتاج الهرمونات ويحسن الأيض',
                    timing: 'على معدة فارغة أو مع وجبة خفيفة'
                }
            ],
            'metabolism_sleep': [
                {
                    name: 'الميلاتونين',
                    dosage: '1-3 مجم قبل النوم بـ 30 دقيقة',
                    benefits: 'يحسن جودة النوم ويدعم إيقاع الساعة البيولوجية',
                    timing: 'قبل النوم بـ 30-60 دقيقة'
                },
                {
                    name: 'الجليسين',
                    dosage: '3 جرام قبل النوم',
                    benefits: 'يحسن جودة النوم ويقلل درجة حرارة الجسم',
                    timing: 'قبل النوم بـ 30 دقيقة'
                }
            ]
        };
        
        if (complaintSupplements[complaintId]) {
            supplements.push(...complaintSupplements[complaintId]);
        }
    });
    
    // Add general supplements for all clients based on age and weight
    const generalSupplements = [
        {
            name: 'أوميجا 3',
            dosage: `${Math.round(weight * 20)} مجم يومياً`,
            benefits: 'يقلل الالتهاب ويدعم صحة القلب والدماغ والأيض',
            timing: 'مع الوجبة الرئيسية'
        },
        {
            name: 'فيتامين د3',
            dosage: age > 50 ? '3000-4000 وحدة دولية' : '2000-3000 وحدة دولية',
            benefits: 'يدعم صحة العظام والمناعة والتوازن الهرموني',
            timing: 'مع وجبة تحتوي على دهون'
        },
        {
            name: 'مالتي فيتامين عالي الجودة',
            dosage: 'قرص واحد يومياً',
            benefits: 'يغطي النقص في الفيتامينات والمعادن الأساسية',
            timing: 'مع وجبة الإفطار'
        }
    ];
    
    supplements.push(...generalSupplements);
    
    // Remove duplicates based on supplement name
    const uniqueSupplements = supplements.filter((supplement, index, self) => 
        index === self.findIndex(s => s.name === supplement.name)
    );
    
    return uniqueSupplements;
}

// Generate personalized workout plan
function generatePersonalizedWorkoutPlan(data) {
    const plan = {
        clientInfo: data,
        workoutType: selectedWorkoutType,
        categories: selectedCategories,
        complaints: selectedComplaints,
        injuries: selectedInjuries,
        exercises: [],
        nutritionAdvice: [],
        supplementAdvice: [],
        injuryTreatment: []
    };
    
    // Generate exercises based on categories and goals
    selectedCategories.forEach(category => {
        const exercises = generateExercisesForCategory(category, data.workoutGoal, selectedWorkoutType);
        plan.exercises.push(...exercises);
    });
    
    // Generate nutrition advice for complaints
    selectedComplaints.forEach(complaint => {
        const advice = generateNutritionAdviceForComplaint(complaint, data);
        plan.nutritionAdvice.push(advice);
    });
    
    // Generate supplement advice
    const supplements = generateSupplementAdvice(selectedComplaints, data);
    plan.supplementAdvice = supplements;
    
    // Generate injury treatment
    selectedInjuries.forEach(injury => {
        const treatment = generateInjuryTreatment(injury);
        plan.injuryTreatment.push(treatment);
    });
    
    return plan;
}

// Generate exercises for category
function generateExercisesForCategory(category, goal, workoutType) {
    const exercises = [];
    
    // Exercise database based on category and type
    const exerciseDatabase = {
        upper_body: {
            gym: [
                {
                    name: 'تمرين البنش برس',
                    sets: 3,
                    reps: '8-12',
                    rest: '90 ثانية',
                    instructions: 'استلق على البنش واحمل البار بقبضة متوسطة العرض. أنزل البار ببطء إلى الصدر ثم ادفعه لأعلى.',
                    commonMistakes: ['عدم التحكم في النزول', 'رفع الوركين عن البنش', 'عدم لمس الصدر'],
                    alternative: {
                        name: 'تمرين الدمبل فلاي',
                        sets: 3,
                        reps: '10-15',
                        rest: '60 ثانية',
                        instructions: 'استلق على البنش واحمل دمبل في كل يد. افتح الذراعين على الجانبين ثم اجمعهما فوق الصدر.'
                    }
                },
                {
                    name: 'تمرين العقلة',
                    sets: 3,
                    reps: '5-10',
                    rest: '120 ثانية',
                    instructions: 'تعلق من البار واسحب جسمك لأعلى حتى يصل الذقن فوق البار.',
                    commonMistakes: ['استخدام الزخم', 'عدم النزول بالكامل', 'تأرجح الجسم'],
                    alternative: {
                        name: 'تمرين السحب بالكابل',
                        sets: 3,
                        reps: '8-12',
                        rest: '90 ثانية',
                        instructions: 'اجلس أمام جهاز الكابل واسحب البار إلى الصدر مع الضغط على لوحي الكتف.'
                    }
                }
            ],
            home: [
                {
                    name: 'تمرين الضغط',
                    sets: 3,
                    reps: '10-15',
                    rest: '60 ثانية',
                    instructions: 'ابدأ في وضع البلانك واخفض صدرك نحو الأرض ثم ادفع جسمك لأعلى.',
                    commonMistakes: ['ترهل الوركين', 'عدم النزول بالكامل', 'وضع اليدين خطأ'],
                    alternative: {
                        name: 'تمرين الضغط على الركبتين',
                        sets: 3,
                        reps: '12-20',
                        rest: '45 ثانية',
                        instructions: 'نفس تمرين الضغط لكن مع الارتكاز على الركبتين بدلاً من أصابع القدمين.'
                    }
                }
            ]
        },
        lower_body: {
            gym: [
                {
                    name: 'تمرين السكوات بالبار',
                    sets: 3,
                    reps: '8-12',
                    rest: '120 ثانية',
                    instructions: 'ضع البار على الكتفين واخفض جسمك كأنك تجلس على كرسي ثم قف مرة أخرى.',
                    commonMistakes: ['الركبتين تتجهان للداخل', 'عدم النزول بالكامل', 'انحناء الظهر'],
                    alternative: {
                        name: 'تمرين السكوات بالدمبل',
                        sets: 3,
                        reps: '10-15',
                        rest: '90 ثانية',
                        instructions: 'احمل دمبل أمام الصدر وقم بنفس حركة السكوات.'
                    }
                }
            ],
            home: [
                {
                    name: 'تمرين السكوات بوزن الجسم',
                    sets: 3,
                    reps: '15-20',
                    rest: '60 ثانية',
                    instructions: 'قف مع فتح القدمين بعرض الكتفين واخفض جسمك كأنك تجلس على كرسي.',
                    commonMistakes: ['الركبتين تتجهان للداخل', 'عدم النزول بالكامل', 'رفع الكعبين'],
                    alternative: {
                        name: 'تمرين الطعنات',
                        sets: 3,
                        reps: '10-12 لكل رجل',
                        rest: '60 ثانية',
                        instructions: 'خذ خطوة كبيرة للأمام واخفض الركبة الخلفية نحو الأرض.'
                    }
                }
            ]
        },
        cardio: {
            gym: [
                {
                    name: 'الجري على السير',
                    sets: 1,
                    reps: '20-30 دقيقة',
                    rest: 'حسب الحاجة',
                    instructions: 'ابدأ بسرعة معتدلة وزد تدريجياً. حافظ على وضعية جسم مستقيمة.',
                    commonMistakes: ['البدء بسرعة عالية', 'عدم الإحماء', 'وضعية جسم خاطئة'],
                    alternative: {
                        name: 'الدراجة الثابتة',
                        sets: 1,
                        reps: '25-35 دقيقة',
                        rest: 'حسب الحاجة',
                        instructions: 'اضبط المقاومة حسب مستواك وحافظ على إيقاع ثابت.'
                    }
                }
            ],
            home: [
                {
                    name: 'الجري في المكان',
                    sets: 3,
                    reps: '2-3 دقائق',
                    rest: '30 ثانية',
                    instructions: 'اجر في مكانك مع رفع الركبتين عالياً وتحريك الذراعين.',
                    commonMistakes: ['عدم رفع الركبتين', 'الهبوط بقوة', 'عدم تحريك الذراعين'],
                    alternative: {
                        name: 'تمرين الجامبينغ جاكس',
                        sets: 3,
                        reps: '30-45 ثانية',
                        rest: '15 ثانية',
                        instructions: 'اقفز مع فتح الرجلين ورفع الذراعين فوق الرأس ثم العودة للوضع الأصلي.'
                    }
                }
            ]
        }
    };
    
    const categoryExercises = exerciseDatabase[category]?.[workoutType] || [];
    
    // Adjust exercises based on goal
    categoryExercises.forEach(exercise => {
        const adjustedExercise = adjustExerciseForGoal(exercise, goal);
        exercises.push(adjustedExercise);
    });
    
    return exercises;
}

// Adjust exercise for specific goal
function adjustExerciseForGoal(exercise, goal) {
    const adjusted = { ...exercise };
    
    switch (goal) {
        case 'weight_loss':
            adjusted.sets = Math.max(3, adjusted.sets);
            adjusted.reps = adjusted.reps.includes('-') ? 
                adjusted.reps.split('-')[1] + '-' + (parseInt(adjusted.reps.split('-')[1]) + 5) :
                '15-20';
            adjusted.rest = '45-60 ثانية';
            break;
        case 'muscle_gain':
            adjusted.sets = Math.max(4, adjusted.sets + 1);
            adjusted.reps = '6-10';
            adjusted.rest = '90-120 ثانية';
            break;
        case 'strength':
            adjusted.sets = Math.max(4, adjusted.sets + 1);
            adjusted.reps = '3-6';
            adjusted.rest = '120-180 ثانية';
            break;
        case 'endurance':
            adjusted.sets = Math.max(3, adjusted.sets);
            adjusted.reps = '15-25';
            adjusted.rest = '30-45 ثانية';
            break;
    }
    
    return adjusted;
}

// Generate nutrition advice for complaint
function generateNutritionAdviceForComplaint(complaint, clientData) {
    const adviceDatabase = {
        fatigue: {
            title: 'نصائح غذائية للتعب والإرهاق',
            advice: [
                'تناول وجبات متوازنة تحتوي على البروتين والكربوهيدرات المعقدة',
                'شرب كمية كافية من الماء (8-10 أكواب يومياً)',
                'تجنب السكريات البسيطة والكافيين الزائد',
                'تناول الأطعمة الغنية بالحديد مثل السبانخ واللحوم الحمراء',
                'إضافة فيتامين B12 والمغنيسيوم للنظام الغذائي'
            ]
        },
        slow_metabolism: {
            title: 'نصائح غذائية لتسريع الأيض',
            advice: [
                'تناول وجبات صغيرة ومتكررة (5-6 وجبات يومياً)',
                'زيادة تناول البروتين (1.2-1.6 جم لكل كيلو من وزن الجسم)',
                'شرب الشاي الأخضر والقهوة باعتدال',
                'تناول الأطعمة الحارة التي تحتوي على الكابسيسين',
                'عدم تفويت وجبة الإفطار'
            ]
        },
        weight_gain: {
            title: 'نصائح غذائية للتحكم في الوزن',
            advice: [
                'تقليل السعرات الحرارية بنسبة 10-20%',
                'زيادة تناول الألياف والخضروات',
                'تجنب الأطعمة المصنعة والسكريات المضافة',
                'تناول البروتين في كل وجبة',
                'شرب الماء قبل الوجبات'
            ]
        },
        digestive_issues: {
            title: 'نصائح غذائية لمشاكل الهضم',
            advice: [
                'تناول الأطعمة الغنية بالبروبيوتيك مثل الزبادي',
                'زيادة تناول الألياف تدريجياً',
                'تجنب الأطعمة المقلية والدهنية',
                'مضغ الطعام جيداً وتناوله ببطء',
                'شرب شاي الزنجبيل والنعناع'
            ]
        }
    };
    
    return adviceDatabase[complaint] || {
        title: 'نصائح غذائية عامة',
        advice: ['تناول نظام غذائي متوازن ومتنوع']
    };
}

// Generate supplement advice
function generateSupplementAdvice(complaints, clientData) {
    const supplements = [];
    
    // Base supplements for everyone
    supplements.push({
        name: 'فيتامين د3',
        dose: calculateVitaminDDose(clientData.weight),
        timing: 'مع وجبة تحتوي على دهون',
        benefits: 'يدعم صحة العظام والمناعة'
    });
    
    supplements.push({
        name: 'أوميغا 3',
        dose: '1000-2000 مجم يومياً',
        timing: 'مع الوجبات',
        benefits: 'يقلل الالتهابات ويدعم صحة القلب'
    });
    
    // Complaint-specific supplements
    if (complaints.includes('fatigue')) {
        supplements.push({
            name: 'فيتامين B12',
            dose: '1000 مكجم يومياً',
            timing: 'صباحاً على معدة فارغة',
            benefits: 'يحسن الطاقة ووظائف الأعصاب'
        });
        
        supplements.push({
            name: 'الحديد',
            dose: calculateIronDose(clientData.gender, clientData.weight),
            timing: 'على معدة فارغة مع فيتامين C',
            benefits: 'يمنع فقر الدم ويحسن الطاقة'
        });
    }
    
    if (complaints.includes('slow_metabolism')) {
        supplements.push({
            name: 'الكافيين الطبيعي',
            dose: '100-200 مجم قبل التمرين',
            timing: 'قبل التمرين بـ 30 دقيقة',
            benefits: 'يزيد معدل الأيض والطاقة'
        });
    }
    
    if (complaints.includes('digestive_issues')) {
        supplements.push({
            name: 'البروبيوتيك',
            dose: '10-50 مليار وحدة يومياً',
            timing: 'مع الطعام أو بعده',
            benefits: 'يحسن صحة الأمعاء والهضم'
        });
    }
    
    return supplements;
}

// Calculate vitamin D dose based on weight
function calculateVitaminDDose(weight) {
    const basedose = 1000; // IU
    const additionalPerKg = 10;
    return `${basedose + (weight * additionalPerKg)} وحدة دولية يومياً`;
}

// Calculate iron dose
function calculateIronDose(gender, weight) {
    const baseDose = gender === 'female' ? 18 : 8; // mg
    return `${baseDose} مجم يومياً`;
}

// Generate injury treatment
function generateInjuryTreatment(injuryId) {
    // Find injury data from loaded injury data
    const injury = injuryData?.injuries?.find(inj => inj.id === injuryId);
    
    if (injury) {
        return {
            title: injury.name?.ar || injury.name || 'علاج الإصابة',
            description: injury.description?.ar || injury.description || '',
            management_plan: injury.management_plan?.ar || injury.management_plan || '',
            exercises: injury.monthly_workout_plan ? formatWorkoutPlan(injury.monthly_workout_plan) : [],
            advice: injury.gym_tips?.ar || injury.gym_tips || [],
            supplements: injury.supplements || [],
            medications: injury.medications || [],
            plants_and_herbs: injury.plants_and_herbs || [],
            recipes: injury.recipes || [],
            disclaimer: injury.disclaimer?.ar || injury.disclaimer || ''
        };
    }
    
    // Fallback treatment data
    return {
        title: 'علاج عام للإصابة',
        description: 'يرجى استشارة طبيب مختص للحصول على تشخيص دقيق وخطة علاج مناسبة.',
        exercises: [],
        advice: ['استشر طبيب مختص', 'تطبيق الراحة والثلج', 'تجنب الأنشطة المؤلمة'],
        supplements: [],
        medications: [],
        plants_and_herbs: [],
        recipes: [],
        disclaimer: 'هذه المعلومات للأغراض التعليمية فقط ولا تغني عن استشارة طبية متخصصة.'
    };
}

// Helper function to format workout plan from injury data
function formatWorkoutPlan(monthlyPlan) {
    const exercises = [];
    
    if (monthlyPlan && typeof monthlyPlan === 'object') {
        Object.keys(monthlyPlan).forEach(month => {
            const monthData = monthlyPlan[month];
            if (monthData && monthData.exercises) {
                monthData.exercises.forEach(exercise => {
                    exercises.push({
                        name: exercise.name || 'تمرين علاجي',
                        sets: exercise.sets || '2-3',
                        reps: exercise.reps || '10-15',
                        instructions: exercise.notes || exercise.instructions || '',
                        month: month
                    });
                });
            }
        });
    }
    
    return exercises;
}

// Display advanced workout plan
function displayAdvancedResults(plan) {
    const container = document.getElementById('workoutPlanResults');
    container.style.display = 'block';
    
    let html = `
        <div class="workout-plan-header">
            <h2><i class="fas fa-star"></i> خطة التمارين الشخصية لـ ${plan.clientInfo.name}</h2>
        </div>
        
        <div class="client-info-section">
            <h3><i class="fas fa-user"></i> معلومات العميل</h3>
            <div class="row">
                <div class="col-md-6">
                    <ul class="client-details">
                        <li><strong>الجنس:</strong> ${plan.clientInfo.gender === 'male' ? 'ذكر' : 'أنثى'}</li>
                        <li><strong>العمر:</strong> ${plan.clientInfo.age} سنة</li>
                        <li><strong>الوزن:</strong> ${plan.clientInfo.weight} كجم</li>
                        <li><strong>الطول:</strong> ${plan.clientInfo.height} سم</li>
                    </ul>
                </div>
                <div class="col-md-6">
                    <ul class="client-details">
                        <li><strong>مستوى النشاط:</strong> ${getActivityLevelName(plan.clientInfo.activityLevel)}</li>
                        <li><strong>الهدف:</strong> ${getGoalName(plan.clientInfo.workoutGoal)}</li>
                        <li><strong>نوع التمارين:</strong> ${plan.workoutType === 'gym' ? 'تمارين الجيم' : 'تمارين منزلية'}</li>
                    </ul>
                </div>
            </div>
        </div>
    `;
    
    // Display weekly plan
    if (plan.weeklyPlan && plan.weeklyPlan.length > 0) {
        html += `
            <div class="weekly-plan-section">
                <h3><i class="fas fa-calendar-week"></i> الخطة الأسبوعية</h3>
                <div class="weekly-grid">
        `;
        
        plan.weeklyPlan.forEach(day => {
            html += `
                <div class="day-card ${day.isWorkoutDay ? 'workout-day' : 'rest-day'}">
                    <div class="day-header">
                        <h4>${day.day}</h4>
                        <span class="day-type">${day.isWorkoutDay ? 'يوم تمرين' : 'يوم راحة'}</span>
                    </div>
                    <div class="day-content">
                        ${day.isWorkoutDay ? `
                            <div class="workout-info">
                                <p><strong>التركيز:</strong> ${getCategoryName(day.focus)}</p>
                                <p><strong>المدة:</strong> ${day.duration}</p>
                                <p><strong>الشدة:</strong> ${day.intensity}</p>
                            </div>
                            <div class="day-exercises">
                                <h5>التمارين:</h5>
                                <ul>
                                    ${day.exercises.map(ex => `
                                        <li>
                                            <strong>${ex.name}</strong><br>
                                            ${ex.sets} مجموعات × ${ex.reps} تكرار
                                            <small>(راحة: ${ex.rest})</small>
                                        </li>
                                    `).join('')}
                                </ul>
                            </div>
                        ` : `
                            <div class="rest-info">
                                <p><strong>النشاط المقترح:</strong></p>
                                <p>${day.restActivity}</p>
                            </div>
                        `}
                    </div>
                </div>
            `;
        });
        
        html += `
                </div>
            </div>
        `;
    }
    
    // Display exercises
    if (plan.exercises && plan.exercises.length > 0) {
        html += `
            <div class="exercises-section">
                <h3><i class="fas fa-dumbbell"></i> التمارين المخصصة</h3>
                <div class="exercises-grid">
        `;
        
        plan.exercises.forEach((exercise, index) => {
            html += `
                <div class="exercise-card">
                    <div class="exercise-header">
                        <h4>${exercise.name || `تمرين ${index + 1}`}</h4>
                        <span class="exercise-category">${getCategoryName(exercise.category)}</span>
                    </div>
                    <div class="exercise-details">
                        <div class="exercise-specs">
                            <span class="spec"><i class="fas fa-repeat"></i> ${exercise.sets} مجموعات</span>
                            <span class="spec"><i class="fas fa-hashtag"></i> ${exercise.reps} تكرار</span>
                            <span class="spec"><i class="fas fa-clock"></i> راحة ${exercise.rest}</span>
                        </div>
                        
                        ${exercise.instructions ? `
                            <div class="exercise-instructions">
                                <h5>التعليمات:</h5>
                                <p>${exercise.instructions}</p>
                            </div>
                        ` : ''}
                        
                        ${exercise.progressionTips && exercise.progressionTips.length > 0 ? `
                            <div class="progression-tips">
                                <h5>نصائح التطوير:</h5>
                                <ul>
                                    ${exercise.progressionTips.map(tip => `<li>${tip}</li>`).join('')}
                                </ul>
                            </div>
                        ` : ''}
                        
                        ${exercise.safetyNotes && exercise.safetyNotes.length > 0 ? `
                            <div class="safety-notes">
                                <h5>ملاحظات الأمان:</h5>
                                <ul>
                                    ${exercise.safetyNotes.map(note => `<li>${note}</li>`).join('')}
                                </ul>
                            </div>
                        ` : ''}
                        
                        ${exercise.alternative ? `
                            <div class="alternative-exercise">
                                <h5>التمرين البديل:</h5>
                                <p><strong>${exercise.alternative.name}</strong></p>
                                <p>${exercise.alternative.instructions || 'تعليمات مشابهة للتمرين الأساسي'}</p>
                            </div>
                        ` : ''}
                    </div>
                </div>
            `;
        });
        
        html += `
                </div>
            </div>
        `;
    }
    
    // Display injury treatments
    if (plan.injuryTreatments && plan.injuryTreatments.length > 0) {
        html += `
            <div class="injury-treatments-section">
                <h3><i class="fas fa-band-aid"></i> علاج الإصابات</h3>
        `;
        
        plan.injuryTreatments.forEach(treatment => {
            html += `
                <div class="treatment-card">
                    <div class="treatment-header">
                        <h4><i class="fas fa-heartbeat"></i> ${treatment.title}</h4>
                        ${treatment.description ? `<p class="treatment-description">${treatment.description}</p>` : ''}
                    </div>
                    
                    ${treatment.management_plan ? `
                        <div class="management-plan">
                            <h5><i class="fas fa-clipboard-list"></i> خطة العلاج:</h5>
                            <p>${treatment.management_plan}</p>
                        </div>
                    ` : ''}
                    
                    ${treatment.exercises && treatment.exercises.length > 0 ? `
                        <div class="treatment-exercises">
                            <h5><i class="fas fa-dumbbell"></i> التمارين العلاجية:</h5>
                            ${treatment.exercises.map(exercise => `
                                <div class="therapeutic-exercise">
                                    <h6>${exercise.name}</h6>
                                    <div class="exercise-specs">
                                        <span><i class="fas fa-repeat"></i> ${exercise.sets} مجموعات</span>
                                        <span><i class="fas fa-hashtag"></i> ${exercise.reps} تكرار</span>
                                        ${exercise.month ? `<span><i class="fas fa-calendar"></i> ${exercise.month}</span>` : ''}
                                    </div>
                                    ${exercise.instructions ? `<p class="exercise-instructions">${exercise.instructions}</p>` : ''}
                                </div>
                            `).join('')}
                        </div>
                    ` : ''}
                    
                    ${treatment.supplements && treatment.supplements.length > 0 ? `
                        <div class="treatment-supplements">
                            <h5><i class="fas fa-pills"></i> المكملات الغذائية:</h5>
                            ${treatment.supplements.map(supplement => `
                                <div class="supplement-item">
                                    <div class="supplement-header">
                                        <strong>${supplement.name?.ar || supplement.name}</strong>
                                        ${supplement.dose?.ar || supplement.dose ? `<span class="dose">${supplement.dose?.ar || supplement.dose}</span>` : ''}
                                    </div>
                                    ${supplement.benefits?.ar || supplement.benefits ? `<p class="supplement-benefits">${supplement.benefits?.ar || supplement.benefits}</p>` : ''}
                                    ${supplement.administration?.ar || supplement.administration ? `<small class="administration">${supplement.administration?.ar || supplement.administration}</small>` : ''}
                                </div>
                            `).join('')}
                        </div>
                    ` : ''}
                    
                    ${treatment.medications && treatment.medications.length > 0 ? `
                        <div class="treatment-medications">
                            <h5><i class="fas fa-prescription-bottle-alt"></i> الأدوية:</h5>
                            ${treatment.medications.map(medication => `
                                <div class="medication-item">
                                    <strong>${medication.name?.ar || medication.name}</strong>
                                    ${medication.dose?.ar || medication.dose ? `<span class="dose">${medication.dose?.ar || medication.dose}</span>` : ''}
                                    ${medication.benefits?.ar || medication.benefits ? `<p>${medication.benefits?.ar || medication.benefits}</p>` : ''}
                                </div>
                            `).join('')}
                        </div>
                    ` : ''}
                    
                    ${treatment.plants_and_herbs && treatment.plants_and_herbs.length > 0 ? `
                        <div class="treatment-herbs">
                            <h5><i class="fas fa-leaf"></i> الأعشاب والنباتات الطبية:</h5>
                            ${treatment.plants_and_herbs.map(herb => `
                                <div class="herb-item">
                                    <strong>${herb.name?.ar || herb.name}</strong>
                                    ${herb.dose?.ar || herb.dose ? `<span class="dose">${herb.dose?.ar || herb.dose}</span>` : ''}
                                    ${herb.benefits?.ar || herb.benefits ? `<p>${herb.benefits?.ar || herb.benefits}</p>` : ''}
                                </div>
                            `).join('')}
                        </div>
                    ` : ''}
                    
                    ${treatment.recipes && treatment.recipes.length > 0 ? `
                        <div class="treatment-recipes">
                            <h5><i class="fas fa-utensils"></i> وصفات مساعدة:</h5>
                            ${treatment.recipes.map(recipe => `
                                <div class="recipe-item">
                                    <h6>${recipe.name?.ar || recipe.name}</h6>
                                    ${recipe.ingredients?.ar || recipe.ingredients ? `<p><strong>المكونات:</strong> ${recipe.ingredients?.ar || recipe.ingredients}</p>` : ''}
                                    ${recipe.instructions?.ar || recipe.instructions ? `<p><strong>التحضير:</strong> ${recipe.instructions?.ar || recipe.instructions}</p>` : ''}
                                </div>
                            `).join('')}
                        </div>
                    ` : ''}
                    
                    ${treatment.advice && treatment.advice.length > 0 ? `
                        <div class="treatment-tips">
                            <h5><i class="fas fa-lightbulb"></i> نصائح إضافية:</h5>
                            <ul>
                                ${treatment.advice.map(tip => `<li>${tip}</li>`).join('')}
                            </ul>
                        </div>
                    ` : ''}
                    
                    ${treatment.disclaimer ? `
                        <div class="treatment-disclaimer">
                            <div class="alert alert-warning">
                                <i class="fas fa-exclamation-triangle"></i>
                                <strong>تنبيه طبي:</strong> ${treatment.disclaimer}
                            </div>
                        </div>
                    ` : ''}
                </div>
            `;
        });
        
        html += `</div>`;
    }
    
    // Display nutrition advice
    if (plan.nutritionAdvice && plan.nutritionAdvice.length > 0) {
        html += `
            <div class="nutrition-advice-section">
                <h3><i class="fas fa-apple-alt"></i> النصائح الغذائية</h3>
        `;
        
        plan.nutritionAdvice.forEach(advice => {
            html += `
                <div class="nutrition-card">
                    <h4>نصائح لـ ${advice.complaint}</h4>
                    
                    <div class="recommended-foods">
                        <h5>الأطعمة المُوصى بها:</h5>
                        <ul>
                            ${advice.foods.map(food => `<li>${food}</li>`).join('')}
                        </ul>
                    </div>
                    
                    <div class="foods-to-avoid">
                        <h5>الأطعمة التي يجب تجنبها:</h5>
                        <ul>
                            ${advice.avoid.map(food => `<li>${food}</li>`).join('')}
                        </ul>
                    </div>
                    
                    <div class="nutrition-tips">
                        <h5>نصائح إضافية:</h5>
                        <ul>
                            ${advice.tips.map(tip => `<li>${tip}</li>`).join('')}
                        </ul>
                    </div>
                </div>
            `;
        });
        
        html += `</div>`;
    }
    
    // Display supplements
    if (plan.supplements && plan.supplements.length > 0) {
        html += `
            <div class="supplements-section">
                <h3><i class="fas fa-pills"></i> المكملات الغذائية</h3>
                <div class="supplements-grid">
        `;
        
        plan.supplements.forEach(supplement => {
            html += `
                <div class="supplement-card">
                    <h4>${supplement.name}</h4>
                    <div class="supplement-details">
                        <p><strong>الجرعة:</strong> ${supplement.dosage}</p>
                        <p><strong>التوقيت:</strong> ${supplement.timing}</p>
                        <p><strong>الفوائد:</strong> ${supplement.benefits}</p>
                        ${supplement.note ? `<p class="supplement-note"><strong>ملاحظة:</strong> ${supplement.note}</p>` : ''}
                    </div>
                </div>
            `;
        });
        
        html += `
                </div>
            </div>
        `;
    }
    
    // Display common mistakes
    if (plan.commonMistakes && plan.commonMistakes.length > 0) {
        html += `
            <div class="common-mistakes-section">
                <h3><i class="fas fa-exclamation-circle"></i> الأخطاء الشائعة وكيفية تجنبها</h3>
                <div class="mistakes-grid">
        `;
        
        plan.commonMistakes.forEach(mistake => {
            html += `
                <div class="mistake-card">
                    <div class="mistake-header">
                        <h4><i class="fas fa-times-circle text-danger"></i> ${mistake.mistake}</h4>
                    </div>
                    <div class="mistake-content">
                        <div class="correction">
                            <h5><i class="fas fa-check-circle text-success"></i> الطريقة الصحيحة:</h5>
                            <p>${mistake.correction}</p>
                        </div>
                        ${mistake.tips && mistake.tips.length > 0 ? `
                            <div class="mistake-tips">
                                <h6><i class="fas fa-lightbulb"></i> نصائح إضافية:</h6>
                                <ul>
                                    ${mistake.tips.map(tip => `<li>${tip}</li>`).join('')}
                                </ul>
                            </div>
                        ` : ''}
                    </div>
                </div>
            `;
        });
        
        html += `
                </div>
            </div>
        `;
    }
    
    // Display exercise alternatives
    if (plan.alternatives && plan.alternatives.length > 0) {
        html += `
            <div class="alternatives-section">
                <h3><i class="fas fa-exchange-alt"></i> التمارين البديلة</h3>
                <div class="alternatives-grid">
        `;
        
        plan.alternatives.forEach(alternative => {
            html += `
                <div class="alternative-card">
                    <div class="alternative-header">
                        <h4><i class="fas fa-dumbbell"></i> ${alternative.name}</h4>
                        <span class="alternative-category">${getCategoryName(alternative.category)}</span>
                    </div>
                    <div class="alternative-content">
                        <div class="alternative-specs">
                            <span class="spec"><i class="fas fa-repeat"></i> ${alternative.sets} مجموعات</span>
                            <span class="spec"><i class="fas fa-hashtag"></i> ${alternative.reps} تكرار</span>
                            <span class="spec"><i class="fas fa-clock"></i> راحة ${alternative.rest}</span>
                        </div>
                        
                        ${alternative.instructions ? `
                            <div class="alternative-instructions">
                                <h5>التعليمات:</h5>
                                <p>${alternative.instructions}</p>
                            </div>
                        ` : ''}
                        
                        ${alternative.benefits && alternative.benefits.length > 0 ? `
                            <div class="alternative-benefits">
                                <h5><i class="fas fa-star"></i> الفوائد:</h5>
                                <ul>
                                    ${alternative.benefits.map(benefit => `<li>${benefit}</li>`).join('')}
                                </ul>
                            </div>
                        ` : ''}
                        
                        ${alternative.suitableFor && alternative.suitableFor.length > 0 ? `
                            <div class="suitable-for">
                                <h6><i class="fas fa-user-check"></i> مناسب لـ:</h6>
                                <div class="suitable-tags">
                                    ${alternative.suitableFor.map(condition => `<span class="suitable-tag">${condition}</span>`).join('')}
                                </div>
                            </div>
                        ` : ''}
                    </div>
                </div>
            `;
        });
        
        html += `
                </div>
            </div>
        `;
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

// Helper function to get goal name in Arabic
function getGoalName(goal) {
    const goals = {
        'weight_loss': 'فقدان الوزن',
        'muscle_gain': 'زيادة العضلات',
        'strength': 'زيادة القوة',
        'endurance': 'تحسين التحمل',
        'general_fitness': 'لياقة عامة',
        'rehabilitation': 'إعادة تأهيل'
    };
    return goals[goal] || goal;
}

// Helper function to get category name in Arabic
function getCategoryName(category) {
    const categories = {
        'upper_body': 'تمارين الجزء العلوي',
        'lower_body': 'تمارين الجزء السفلي',
        'cardio': 'تمارين القلب',
        'strength': 'تمارين القوة',
        'flexibility': 'تمارين المرونة',
        'core': 'تمارين البطن والجذع',
        'functional': 'التمارين الوظيفية',
        'rehabilitation': 'تمارين إعادة التأهيل'
    };
    return categories[category] || category;
}

// Show message function
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
    if (mainContainer) {
        mainContainer.insertBefore(messageDiv, mainContainer.firstChild);
    }
    
    // Auto remove after 5 seconds
    setTimeout(() => {
        if (messageDiv.parentNode) {
            messageDiv.remove();
        }
    }, 5000);
}

// Display workout plan
function displayWorkoutPlan(plan) {
    const container = document.getElementById('workoutPlanResults');
    
    let html = `
        <div class="workout-plan">
            <h2 class="text-center mb-4">
                <i class="fas fa-star me-2"></i>خطة التمارين الشخصية لـ ${plan.clientInfo.name}
            </h2>
            
            <div class="row mb-4">
                <div class="col-md-6">
                    <h5><i class="fas fa-user me-2"></i>معلومات العميل</h5>
                    <ul class="list-unstyled">
                        <li><strong>الجنس:</strong> ${plan.clientInfo.gender === 'male' ? 'ذكر' : 'أنثى'}</li>
                        <li><strong>العمر:</strong> ${plan.clientInfo.age} سنة</li>
                        <li><strong>الوزن:</strong> ${plan.clientInfo.weight} كجم</li>
                        <li><strong>الطول:</strong> ${plan.clientInfo.height} سم</li>
                        <li><strong>الهدف:</strong> ${getGoalName(plan.clientInfo.workoutGoal)}</li>
                    </ul>
                </div>
                <div class="col-md-6">
                    <h5><i class="fas fa-dumbbell me-2"></i>نوع التمارين</h5>
                    <p><strong>${plan.workoutType === 'gym' ? 'تمارين الجيم' : 'تمارين منزلية'}</strong></p>
                    <h5><i class="fas fa-list me-2"></i>الفئات المختارة</h5>
                    <ul class="list-unstyled">
                        ${plan.categories.map(cat => `<li>• ${getCategoryName(cat)}</li>`).join('')}
                    </ul>
                </div>
            </div>
        </div>
    `;
    
    // Add exercises
    if (plan.exercises.length > 0) {
        html += `
            <div class="card mb-4">
                <div class="card-header bg-success text-white">
                    <h4><i class="fas fa-dumbbell me-2"></i>التمارين المقترحة</h4>
                </div>
                <div class="card-body">
                    ${plan.exercises.map(exercise => `
                        <div class="exercise-card">
                            <h5 class="text-primary">${exercise.name}</h5>
                            <div class="row">
                                <div class="col-md-8">
                                    <p><strong>المجموعات:</strong> ${exercise.sets} | <strong>التكرارات:</strong> ${exercise.reps} | <strong>الراحة:</strong> ${exercise.rest}</p>
                                    <p><strong>التعليمات:</strong> ${exercise.instructions}</p>
                                    <div class="alert alert-warning">
                                        <strong>الأخطاء الشائعة:</strong>
                                        <ul class="mb-0">
                                            ${exercise.commonMistakes.map(mistake => `<li>${mistake}</li>`).join('')}
                                        </ul>
                                    </div>
                                </div>
                                <div class="col-md-4">
                                    ${exercise.alternative ? `
                                        <div class="alternative-exercise">
                                            <h6 class="text-success"><i class="fas fa-exchange-alt me-1"></i>البديل</h6>
                                            <p><strong>${exercise.alternative.name}</strong></p>
                                            <small>${exercise.alternative.sets} مجموعات × ${exercise.alternative.reps}</small>
                                            <p class="small mt-2">${exercise.alternative.instructions}</p>
                                        </div>
                                    ` : ''}
                                </div>
                            </div>
                        </div>
                    `).join('')}
                </div>
            </div>
        `;
    }
    
    // Add nutrition advice
    if (plan.nutritionAdvice.length > 0) {
        html += `
            <div class="card mb-4">
                <div class="card-header bg-info text-white">
                    <h4><i class="fas fa-apple-alt me-2"></i>النصائح الغذائية</h4>
                </div>
                <div class="card-body">
                    ${plan.nutritionAdvice.map(advice => `
                        <div class="mb-4">
                            <h5 class="text-info">${advice.title}</h5>
                            <ul>
                                ${advice.advice.map(tip => `<li>${tip}</li>`).join('')}
                            </ul>
                        </div>
                    `).join('')}
                </div>
            </div>
        `;
    }
    
    // Add supplement advice
    if (plan.supplementAdvice.length > 0) {
        html += `
            <div class="card mb-4">
                <div class="card-header bg-warning text-white">
                    <h4><i class="fas fa-pills me-2"></i>المكملات الغذائية المقترحة</h4>
                </div>
                <div class="card-body">
                    <div class="row">
                        ${plan.supplementAdvice.map(supplement => `
                            <div class="col-md-6 mb-3">
                                <div class="border rounded p-3">
                                    <h6 class="text-warning">${supplement.name}</h6>
                                    <p><strong>الجرعة:</strong> ${supplement.dose}</p>
                                    <p><strong>التوقيت:</strong> ${supplement.timing}</p>
                                    <p><strong>الفوائد:</strong> ${supplement.benefits}</p>
                                </div>
                            </div>
                        `).join('')}
                    </div>
                </div>
            </div>
        `;
    }
    
    // Add injury treatment
    if (plan.injuryTreatment.length > 0) {
        html += `
            <div class="card mb-4">
                <div class="card-header bg-danger text-white">
                    <h4><i class="fas fa-medkit me-2"></i>علاج الإصابات</h4>
                </div>
                <div class="card-body">
                    ${plan.injuryTreatment.map(treatment => `
                        <div class="mb-4">
                            <h5 class="text-danger">${treatment.title}</h5>
                            
                            ${treatment.exercises.length > 0 ? `
                                <h6 class="text-success mt-3">تمارين العلاج:</h6>
                                ${treatment.exercises.map(exercise => `
                                    <div class="border-start border-success ps-3 mb-2">
                                        <strong>${exercise.name}</strong> - ${exercise.sets} مجموعات × ${exercise.reps}<br>
                                        <small>${exercise.instructions}</small>
                                    </div>
                                `).join('')}
                            ` : ''}
                            
                            <h6 class="text-info mt-3">نصائح العلاج:</h6>
                            <ul>
                                ${treatment.advice.map(tip => `<li>${tip}</li>`).join('')}
                            </ul>
                            
                            ${treatment.supplements.length > 0 ? `
                                <h6 class="text-warning mt-3">المكملات المساعدة:</h6>
                                ${treatment.supplements.map(supplement => `
                                    <div class="border-start border-warning ps-3 mb-2">
                                        <strong>${supplement.name}</strong> - ${supplement.dose}<br>
                                        <small>${supplement.benefits}</small>
                                    </div>
                                `).join('')}
                            ` : ''}
                        </div>
                    `).join('')}
                </div>
            </div>
        `;
    }
    
    // Add disclaimer
    html += `
        <div class="disclaimer">
            <h6><i class="fas fa-info-circle me-2"></i>تنويه طبي</h6>
            <p class="mb-0">هذا الموقع لتحديثك بالمعلومات وتوجيهك للنصائح المفيدة بهدف التعليم والوعي ولا يغني عن زيارة الطبيب</p>
        </div>
    `;
    
    container.innerHTML = html;
    
    // Scroll to results
    container.scrollIntoView({ behavior: 'smooth' });
}

// Helper functions
function getGoalName(goal) {
    const goals = {
        'weight_loss': 'فقدان الوزن',
        'muscle_gain': 'زيادة العضلات',
        'strength': 'زيادة القوة',
        'endurance': 'تحسين التحمل',
        'general_fitness': 'لياقة عامة',
        'rehabilitation': 'إعادة تأهيل'
    };
    return goals[goal] || goal;
}

function getCategoryName(category) {
    const categories = {
        'upper_body': 'تمارين الجزء العلوي',
        'lower_body': 'تمارين الجزء السفلي',
        'cardio': 'تمارين القلب',
        'strength': 'تمارين القوة',
        'flexibility': 'تمارين المرونة',
        'core': 'تمارين البطن والجذع',
        'functional': 'التمارين الوظيفية',
        'rehabilitation': 'تمارين إعادة التأهيل'
    };
    return categories[category] || category;
}