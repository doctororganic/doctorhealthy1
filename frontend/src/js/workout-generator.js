// Workout Generator JavaScript
class WorkoutGenerator {
    constructor() {
        this.exercises = null;
        this.supplements = null;
        this.currentLanguage = localStorage.getItem('language') || 'en';
        this.init();
    }

    async init() {
        await this.loadData();
        this.setupEventListeners();
        this.updateLanguage();
    }

    async loadData() {
        try {
            // Load exercises database
            const exercisesResponse = await fetch('../data/exercises.json');
            this.exercises = await exercisesResponse.json();
            
            // Load supplements database
            const supplementsResponse = await fetch('../data/supplements.json');
            this.supplements = await supplementsResponse.json();
            
            console.log('Data loaded successfully');
        } catch (error) {
            console.error('Error loading data:', error);
        }
    }

    setupEventListeners() {
        // Tab switching
        document.querySelectorAll('.tab-btn').forEach(btn => {
            btn.addEventListener('click', (e) => this.switchTab(e.target.dataset.tab));
        });

        // Workout form submission
        const workoutForm = document.getElementById('workoutForm');
        if (workoutForm) {
            workoutForm.addEventListener('submit', (e) => this.handleWorkoutSubmission(e));
        }

        // Product upload form
        const uploadForm = document.getElementById('productUploadForm');
        if (uploadForm) {
            uploadForm.addEventListener('submit', (e) => this.handleProductUpload(e));
        }

        // Language change listener
        window.addEventListener('languageChanged', () => {
            this.currentLanguage = localStorage.getItem('language') || 'en';
            this.updateLanguage();
        });

        // Initialize supplements display
        this.displaySupplements();
    }

    switchTab(tabName) {
        // Hide all tab contents
        document.querySelectorAll('.tab-content').forEach(content => {
            content.classList.remove('active');
        });

        // Remove active class from all buttons
        document.querySelectorAll('.tab-btn').forEach(btn => {
            btn.classList.remove('active');
        });

        // Show selected tab content
        const selectedContent = document.getElementById(tabName);
        const selectedBtn = document.querySelector(`[data-tab="${tabName}"]`);
        
        if (selectedContent && selectedBtn) {
            selectedContent.classList.add('active');
            selectedBtn.classList.add('active');
        }
    }

    async handleWorkoutSubmission(e) {
        e.preventDefault();
        
        // Validate workout form using ValidationSystem
        const validationRules = [
            { field: 'age', type: 'age', required: true, name: 'Age', min: 13, max: 100 },
            { field: 'weight', type: 'weight', required: true, name: 'Weight' },
            { field: 'height', type: 'height', required: true, name: 'Height' },
            { field: 'gender', required: true, name: 'Gender' },
            { field: 'experience', required: true, name: 'Experience Level' },
            { field: 'goal', required: true, name: 'Primary Goal' },
            { field: 'days-per-week', type: 'numeric', required: true, name: 'Days per Week', min: 1, max: 7 },
            { field: 'session-duration', type: 'numeric', required: true, name: 'Session Duration', min: 15, max: 180 }
        ];
        
        if (!validation.validateForm('workoutForm', validationRules)) {
            return;
        }
        
        const formData = new FormData(e.target);
        const workoutData = {
            personalInfo: {
                age: parseInt(formData.get('age')),
                weight: parseFloat(formData.get('weight')),
                height: parseFloat(formData.get('height')),
                gender: formData.get('gender'),
                experience: formData.get('experience')
            },
            goals: {
                primary: formData.get('goal'),
                secondary: formData.getAll('secondary-goals')
            },
            equipment: formData.getAll('equipment'),
            schedule: {
                daysPerWeek: parseInt(formData.get('days-per-week')),
                sessionDuration: parseInt(formData.get('session-duration'))
            },
            injuries: this.getInjuryData(formData)
        };

        try {
            const workout = await this.generateWorkout(workoutData);
            this.displayWorkoutPlan(workout);
            this.saveWorkoutPlan(workout);
            this.showSuccessMessage('Workout plan generated successfully!');
            validation.clearErrors();
        } catch (error) {
            console.error('Error generating workout:', error);
            this.showError('Failed to generate workout plan. Please try again.');
        }
    }

    getInjuryData(formData) {
        const injuries = [];
        const injuryParts = formData.getAll('injury-part');
        const injurySeverities = formData.getAll('injury-severity');
        const injuryRestrictions = formData.getAll('injury-restrictions');

        for (let i = 0; i < injuryParts.length; i++) {
            if (injuryParts[i]) {
                injuries.push({
                    bodyPart: injuryParts[i],
                    severity: injurySeverities[i] || 'mild',
                    restrictions: injuryRestrictions[i] || ''
                });
            }
        }
        return injuries;
    }

    async generateWorkout(data) {
        if (!this.exercises) {
            throw new Error('Exercise database not loaded');
        }

        const { personalInfo, goals, equipment, schedule, injuries } = data;
        
        // Get appropriate template based on goal
        const template = this.exercises.workout_templates[goals.primary] || this.exercises.workout_templates.muscle_gain;
        
        // Filter exercises based on equipment and injuries
        const availableExercises = this.filterExercises(equipment, injuries);
        
        // Generate 7-day plan
        const weeklyPlan = this.generateWeeklyPlan(template, availableExercises, schedule, personalInfo);
        
        // Add nutrition and recovery advice
        const nutritionAdvice = this.getNutritionAdvice(goals.primary, personalInfo);
        const recoveryTips = this.getRecoveryTips();
        
        return {
            id: `workout-${Date.now()}`,
            timestamp: new Date().toISOString(),
            userInfo: personalInfo,
            goals: goals,
            weeklyPlan: weeklyPlan,
            nutritionAdvice: nutritionAdvice,
            recoveryTips: recoveryTips,
            progressiveOverload: this.getProgressiveOverloadGuidelines()
        };
    }

    filterExercises(equipment, injuries) {
        const filtered = {};
        const injuryRestrictions = injuries.map(injury => 
            this.exercises.injury_modifications[injury.bodyPart]?.avoid || []
        ).flat();

        Object.keys(this.exercises.exercises).forEach(bodyPart => {
            filtered[bodyPart] = this.exercises.exercises[bodyPart].filter(exercise => {
                // Check equipment availability
                const hasEquipment = exercise.equipment.some(eq => 
                    equipment.includes(eq) || eq === 'bodyweight'
                );
                
                // Check injury restrictions
                const isRestricted = injuryRestrictions.includes(exercise.id);
                
                // Additional filtering for exercise difficulty
                const difficultyMatch = this.isDifficultyAppropriate(exercise.difficulty, equipment);
                
                return hasEquipment && !isRestricted && difficultyMatch;
            });
        });

        return filtered;
    }

    isDifficultyAppropriate(exerciseDifficulty, userEquipment) {
        // Beginners should avoid advanced exercises without proper equipment
        if (exerciseDifficulty === 'advanced' && !userEquipment.includes('full_gym')) {
            return userEquipment.includes('dumbbells') || userEquipment.includes('barbell');
        }
        return true;
    }

    generateWeeklyPlan(template, exercises, schedule, personalInfo) {
        const plan = [];
        const daysPerWeek = schedule.daysPerWeek;
        const experienceLevel = personalInfo.experience;
        const sessionDuration = schedule.sessionDuration;
        
        for (let day = 1; day <= 7; day++) {
            if (day <= daysPerWeek) {
                const dayTemplate = template.weekly_split[`day_${day}`] || template.weekly_split.day_1;
                const dayPlan = this.generateDayPlan(dayTemplate, exercises, experienceLevel, sessionDuration);
                plan.push({
                    day: day,
                    name: dayTemplate.name,
                    name_ar: dayTemplate.name_ar,
                    exercises: dayPlan,
                    estimatedDuration: this.calculateDuration(dayPlan),
                    targetDuration: sessionDuration,
                    muscleGroups: dayTemplate.muscle_groups
                });
            } else {
                plan.push({
                    day: day,
                    name: 'Rest Day',
                    name_ar: 'يوم راحة',
                    exercises: [],
                    restDay: true,
                    recoveryActivities: this.getRestDayActivities()
                });
            }
        }
        
        return plan;
    }

    getRestDayActivities() {
        return [
            { activity: 'Light walking', duration: '20-30 minutes', intensity: 'low' },
            { activity: 'Stretching', duration: '10-15 minutes', intensity: 'low' },
            { activity: 'Foam rolling', duration: '10-15 minutes', intensity: 'low' },
            { activity: 'Meditation', duration: '10-20 minutes', intensity: 'low' }
        ];
    }

    generateDayPlan(dayTemplate, exercises, experienceLevel, targetDuration) {
        const dayPlan = [];
        let estimatedTime = 0;
        
        dayTemplate.muscle_groups.forEach(muscleGroup => {
            const availableExercises = exercises[muscleGroup] || [];
            let exerciseCount = this.getExerciseCount(muscleGroup, experienceLevel);
            
            // Adjust exercise count based on target duration
            if (targetDuration < 45) {
                exerciseCount = Math.max(1, exerciseCount - 1);
            } else if (targetDuration > 90) {
                exerciseCount = exerciseCount + 1;
            }
            
            // Select exercises for this muscle group
            const selectedExercises = this.selectExercises(availableExercises, exerciseCount);
            
            selectedExercises.forEach(exercise => {
                const sets = this.getSetsAndReps(exercise, experienceLevel);
                const exerciseTime = this.estimateExerciseTime(sets.sets, sets.rest);
                
                if (estimatedTime + exerciseTime <= targetDuration + 10) { // 10 min buffer
                    dayPlan.push({
                        ...exercise,
                        sets: sets.sets,
                        reps: sets.reps,
                        rest: sets.rest,
                        notes: this.getExerciseNotes(exercise, experienceLevel),
                        estimatedTime: exerciseTime
                    });
                    estimatedTime += exerciseTime;
                }
            });
        });
        
        return dayPlan;
    }

    estimateExerciseTime(sets, restPeriod) {
        const restMinutes = parseInt(restPeriod.split('-')[0]) / 60; // Convert seconds to minutes
        const setTime = 1.5; // Assume 1.5 minutes per set
        return (sets * setTime) + ((sets - 1) * restMinutes);
    }

    selectExercises(exercises, count) {
        if (exercises.length <= count) return exercises;
        
        // Prioritize compound movements and different difficulty levels
        const sorted = exercises.sort((a, b) => {
            const aCompound = a.muscle_groups.length;
            const bCompound = b.muscle_groups.length;
            return bCompound - aCompound;
        });
        
        return sorted.slice(0, count);
    }

    getExerciseCount(muscleGroup, experienceLevel) {
        const counts = {
            beginner: { chest: 2, back: 2, legs: 3, shoulders: 2, arms: 2, core: 2, cardio: 1 },
            intermediate: { chest: 3, back: 3, legs: 4, shoulders: 3, arms: 3, core: 2, cardio: 1 },
            advanced: { chest: 4, back: 4, legs: 5, shoulders: 4, arms: 4, core: 3, cardio: 1 }
        };
        
        return counts[experienceLevel]?.[muscleGroup] || 2;
    }

    getSetsAndReps(exercise, experienceLevel) {
        const level = exercise.difficulty;
        const userLevel = experienceLevel;
        
        const setReps = {
            beginner: { sets: 2, reps: '8-12', rest: '60-90s' },
            intermediate: { sets: 3, reps: '8-12', rest: '90-120s' },
            advanced: { sets: 4, reps: '6-12', rest: '120-180s' }
        };
        
        return setReps[userLevel] || setReps.intermediate;
    }

    getExerciseNotes(exercise, experienceLevel) {
        const notes = [];
        
        if (exercise.form_cues && exercise.form_cues.length > 0) {
            notes.push(this.currentLanguage === 'ar' ? 
                exercise.form_cues_ar?.[0] || exercise.form_cues[0] : 
                exercise.form_cues[0]
            );
        }
        
        if (exercise.alternatives && exercise.alternatives.length > 0) {
            notes.push(`Alternative: ${exercise.alternatives[0]}`);
        }
        
        return notes;
    }

    calculateDuration(dayPlan) {
        // Estimate 3-4 minutes per set including rest
        const totalSets = dayPlan.reduce((sum, exercise) => sum + exercise.sets, 0);
        return Math.round(totalSets * 3.5); // minutes
    }

    getNutritionAdvice(goal, personalInfo) {
        const baseCalories = this.calculateBaseCalories(personalInfo);
        
        const advice = {
            weight_loss: {
                pre_workout: 'Light snack 30-60 minutes before: banana with almond butter',
                post_workout: 'Protein + carbs within 30 minutes: protein shake with fruit',
                daily: 'Maintain caloric deficit, prioritize protein (1.6-2.2g/kg body weight)',
                hydration: '2-3L water daily, extra 500ml per hour of exercise',
                calories: `Target: ${Math.round(baseCalories * 0.85)} calories/day (15% deficit)`,
                macros: {
                    protein: `${Math.round(personalInfo.weight * 1.8)}g`,
                    carbs: `${Math.round(personalInfo.weight * 3)}g`,
                    fats: `${Math.round(personalInfo.weight * 0.8)}g`
                }
            },
            muscle_gain: {
                pre_workout: 'Balanced meal 1-2 hours before: oats with protein powder',
                post_workout: 'High protein meal within 2 hours: chicken with rice and vegetables',
                daily: 'Caloric surplus, high protein (2.2-3g/kg), complex carbs',
                hydration: '3-4L water daily, monitor urine color',
                calories: `Target: ${Math.round(baseCalories * 1.15)} calories/day (15% surplus)`,
                macros: {
                    protein: `${Math.round(personalInfo.weight * 2.2)}g`,
                    carbs: `${Math.round(personalInfo.weight * 5)}g`,
                    fats: `${Math.round(personalInfo.weight * 1)}g`
                }
            },
            endurance: {
                pre_workout: 'Carb-rich meal 2-3 hours before: pasta with lean protein',
                post_workout: 'Carbs + protein 3:1 ratio: chocolate milk or recovery drink',
                daily: 'High carbohydrate intake (6-10g/kg), moderate protein',
                hydration: '4-5L water daily, electrolyte replacement during long sessions',
                calories: `Target: ${Math.round(baseCalories * 1.1)} calories/day`,
                macros: {
                    protein: `${Math.round(personalInfo.weight * 1.6)}g`,
                    carbs: `${Math.round(personalInfo.weight * 7)}g`,
                    fats: `${Math.round(personalInfo.weight * 0.8)}g`
                }
            },
            strength: {
                pre_workout: 'Moderate carbs + protein: Greek yogurt with berries',
                post_workout: 'High protein meal: lean meat with sweet potato',
                daily: 'Adequate calories, high protein (2-2.5g/kg), strategic carb timing',
                hydration: '3L water daily, creatine supplementation beneficial',
                calories: `Target: ${Math.round(baseCalories * 1.05)} calories/day`,
                macros: {
                    protein: `${Math.round(personalInfo.weight * 2.2)}g`,
                    carbs: `${Math.round(personalInfo.weight * 4)}g`,
                    fats: `${Math.round(personalInfo.weight * 1)}g`
                }
            }
        };
        
        return advice[goal] || advice.muscle_gain;
    }

    calculateBaseCalories(personalInfo) {
        // Mifflin-St Jeor Equation
        let bmr;
        if (personalInfo.gender === 'male') {
            bmr = (10 * personalInfo.weight) + (6.25 * personalInfo.height) - (5 * personalInfo.age) + 5;
        } else {
            bmr = (10 * personalInfo.weight) + (6.25 * personalInfo.height) - (5 * personalInfo.age) - 161;
        }
        
        // Activity factor based on experience level
        const activityFactors = {
            beginner: 1.4,
            intermediate: 1.55,
            advanced: 1.7
        };
        
        return bmr * (activityFactors[personalInfo.experience] || 1.55);
    }

    getRecoveryTips() {
        return {
            sleep: {
                duration: '7-9 hours per night',
                quality: 'Dark, cool room (65-68°F), consistent schedule',
                tips: 'No screens 1 hour before bed, magnesium supplementation'
            },
            active_recovery: {
                activities: 'Light walking, yoga, swimming, stretching',
                frequency: '2-3 times per week on rest days',
                duration: '20-30 minutes low intensity'
            },
            stress_management: {
                techniques: 'Meditation, deep breathing, journaling',
                importance: 'Chronic stress elevates cortisol, impairs recovery',
                recommendation: '10-15 minutes daily stress reduction'
            },
            hydration: {
                daily: '35ml per kg body weight minimum',
                exercise: 'Additional 500-750ml per hour of training',
                indicators: 'Pale yellow urine, minimal thirst'
            }
        };
    }

    getProgressiveOverloadGuidelines() {
        return {
            principles: [
                'Gradually increase weight by 2.5-5% when you can complete all sets with perfect form',
                'Add 1-2 reps per set before increasing weight',
                'Increase training volume by 10-20% per week maximum',
                'Focus on one variable at a time (weight, reps, or sets)'
            ],
            weekly_progression: {
                week_1: 'Establish baseline with proper form',
                week_2: 'Add 1-2 reps or 2.5-5% weight',
                week_3: 'Continue progression if form remains good',
                week_4: 'Deload week - reduce volume by 40-50%'
            },
            warning_signs: [
                'Form breakdown',
                'Joint pain',
                'Excessive fatigue',
                'Plateau lasting 2+ weeks'
            ]
        };
    }

    displayWorkoutPlan(workout) {
        const resultsDiv = document.getElementById('workoutResults');
        if (!resultsDiv) return;

        const html = `
            <div class="workout-plan">
                <div class="plan-header">
                    <h3>${this.currentLanguage === 'ar' ? 'خطة التمرين الخاصة بك' : 'Your Workout Plan'}</h3>
                    <p class="plan-id">Plan ID: ${workout.id}</p>
                </div>
                
                <div class="weekly-overview">
                    ${workout.weeklyPlan.map(day => `
                        <div class="day-card ${day.restDay ? 'rest-day' : ''}">
                            <h4>${this.currentLanguage === 'ar' ? day.name_ar : day.name}</h4>
                            ${day.restDay ? 
                                `<p>${this.currentLanguage === 'ar' ? 'يوم راحة - استرخ وتعافى' : 'Rest Day - Relax and Recover'}</p>` :
                                `
                                <div class="exercises-summary">
                                    <p>${day.exercises.length} ${this.currentLanguage === 'ar' ? 'تمارين' : 'exercises'}</p>
                                    <p>${day.estimatedDuration} ${this.currentLanguage === 'ar' ? 'دقيقة' : 'minutes'}</p>
                                </div>
                                <div class="exercise-list">
                                    ${day.exercises.map(ex => `
                                        <div class="exercise-item">
                                            <span class="exercise-name">${this.currentLanguage === 'ar' ? ex.name_ar : ex.name}</span>
                                            <span class="exercise-details">${ex.sets} × ${ex.reps}</span>
                                        </div>
                                    `).join('')}
                                </div>
                                `
                            }
                        </div>
                    `).join('')}
                </div>
                
                <div class="nutrition-recovery">
                    <div class="nutrition-advice">
                        <h4>${this.currentLanguage === 'ar' ? 'نصائح التغذية' : 'Nutrition Advice'}</h4>
                        <div class="advice-grid">
                            <div class="advice-item">
                                <strong>${this.currentLanguage === 'ar' ? 'قبل التمرين:' : 'Pre-workout:'}</strong>
                                <p>${workout.nutritionAdvice.pre_workout}</p>
                            </div>
                            <div class="advice-item">
                                <strong>${this.currentLanguage === 'ar' ? 'بعد التمرين:' : 'Post-workout:'}</strong>
                                <p>${workout.nutritionAdvice.post_workout}</p>
                            </div>
                            <div class="advice-item">
                                <strong>${this.currentLanguage === 'ar' ? 'يومياً:' : 'Daily:'}</strong>
                                <p>${workout.nutritionAdvice.daily}</p>
                            </div>
                            <div class="advice-item">
                                <strong>${this.currentLanguage === 'ar' ? 'الترطيب:' : 'Hydration:'}</strong>
                                <p>${workout.nutritionAdvice.hydration}</p>
                            </div>
                        </div>
                    </div>
                    
                    <div class="recovery-tips">
                        <h4>${this.currentLanguage === 'ar' ? 'نصائح التعافي' : 'Recovery Tips'}</h4>
                        <div class="recovery-grid">
                            <div class="recovery-item">
                                <strong>${this.currentLanguage === 'ar' ? 'النوم:' : 'Sleep:'}</strong>
                                <p>${workout.recoveryTips.sleep.duration}</p>
                            </div>
                            <div class="recovery-item">
                                <strong>${this.currentLanguage === 'ar' ? 'التعافي النشط:' : 'Active Recovery:'}</strong>
                                <p>${workout.recoveryTips.active_recovery.activities}</p>
                            </div>
                        </div>
                    </div>
                </div>
                
                <div class="action-buttons">
                    <button onclick="workoutGenerator.downloadWorkoutPDF('${workout.id}')" class="btn btn-primary">
                        ${this.currentLanguage === 'ar' ? 'تحميل PDF' : 'Download PDF'}
                    </button>
                    <button onclick="workoutGenerator.shareWorkout('${workout.id}')" class="btn btn-secondary">
                        ${this.currentLanguage === 'ar' ? 'مشاركة' : 'Share'}
                    </button>
                </div>
            </div>
        `;
        
        resultsDiv.innerHTML = html;
        resultsDiv.style.display = 'block';
    }

    saveWorkoutPlan(workout) {
        try {
            // Save to localStorage
            const savedWorkouts = JSON.parse(localStorage.getItem('savedWorkouts') || '[]');
            savedWorkouts.push(workout);
            localStorage.setItem('savedWorkouts', JSON.stringify(savedWorkouts));
            
            // Save to JSON file
            const filename = `user-workout-${workout.timestamp.split('T')[0]}.json`;
            this.downloadJSON(workout, filename);
            
            console.log('Workout plan saved successfully');
        } catch (error) {
            console.error('Error saving workout plan:', error);
        }
    }

    downloadJSON(data, filename) {
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

    downloadWorkoutPDF(workoutId) {
        // This would integrate with a PDF generation library
        alert(this.currentLanguage === 'ar' ? 
            'ميزة تحميل PDF قيد التطوير' : 
            'PDF download feature coming soon'
        );
    }

    shareWorkout(workoutId) {
        const shareText = this.currentLanguage === 'ar' ? 
            'شاهد خطة التمرين الخاصة بي!' : 
            'Check out my workout plan!';
        
        if (navigator.share) {
            navigator.share({
                title: shareText,
                url: window.location.href
            });
        } else {
            // Fallback to copy to clipboard
            navigator.clipboard.writeText(window.location.href);
            alert(this.currentLanguage === 'ar' ? 
                'تم نسخ الرابط' : 
                'Link copied to clipboard'
            );
        }
    }

    displaySupplements() {
        if (!this.supplements) return;
        
        const supplementsDiv = document.getElementById('supplementsContent');
        if (!supplementsDiv) return;
        
        const isSubscribed = localStorage.getItem('isSubscribed') === 'true';
        
        let html = `
            <div class="supplements-section">
                <div class="basic-supplements">
                    <h3>${this.currentLanguage === 'ar' ? 'المكملات الأساسية' : 'Basic Supplements'}</h3>
                    <div class="supplements-grid">
        `;
        
        // Display basic supplements
        Object.entries(this.supplements.basic_supplements).forEach(([key, supplement]) => {
            html += this.createSupplementCard(supplement, false);
        });
        
        html += `
                    </div>
                </div>
        `;
        
        // Display advanced supplements (subscription required)
        if (isSubscribed) {
            html += `
                <div class="advanced-supplements">
                    <h3>${this.currentLanguage === 'ar' ? 'المكملات المتقدمة' : 'Advanced Supplements'}</h3>
                    <div class="supplements-grid">
            `;
            
            Object.entries(this.supplements.advanced_supplements).forEach(([key, supplement]) => {
                html += this.createSupplementCard(supplement, true);
            });
            
            html += `
                    </div>
                </div>
            `;
        } else {
            html += `
                <div class="subscription-prompt">
                    <h3>${this.currentLanguage === 'ar' ? 'المكملات المتقدمة' : 'Advanced Supplements'}</h3>
                    <p>${this.currentLanguage === 'ar' ? 
                        'اشترك للوصول إلى المكملات المتقدمة والتوصيات المخصصة' : 
                        'Subscribe to access advanced supplements and personalized recommendations'
                    }</p>
                    <button class="btn btn-primary" onclick="workoutGenerator.showSubscriptionModal()">
                        ${this.currentLanguage === 'ar' ? 'اشترك الآن' : 'Subscribe Now'}
                    </button>
                </div>
            `;
        }
        
        // Display supplement stacks
        html += `
                <div class="supplement-stacks">
                    <h3>${this.currentLanguage === 'ar' ? 'مجموعات المكملات' : 'Supplement Stacks'}</h3>
                    <div class="stacks-grid">
        `;
        
        Object.entries(this.supplements.supplement_stacks).forEach(([key, stack]) => {
            if (!stack.subscription_required || isSubscribed) {
                html += this.createStackCard(stack);
            }
        });
        
        html += `
                    </div>
                </div>
            </div>
        `;
        
        supplementsDiv.innerHTML = html;
    }

    createSupplementCard(supplement, isAdvanced) {
        const name = this.currentLanguage === 'ar' ? supplement.name_ar : supplement.name;
        const benefits = this.currentLanguage === 'ar' ? supplement.benefits_ar : supplement.benefits;
        
        return `
            <div class="supplement-card ${isAdvanced ? 'advanced' : 'basic'}">
                <div class="supplement-header">
                    <h4>${name}</h4>
                    <span class="category">${supplement.category}</span>
                </div>
                <div class="supplement-benefits">
                    <h5>${this.currentLanguage === 'ar' ? 'الفوائد:' : 'Benefits:'}</h5>
                    <ul>
                        ${benefits.slice(0, 3).map(benefit => `<li>${benefit}</li>`).join('')}
                    </ul>
                </div>
                <div class="supplement-dosage">
                    <h5>${this.currentLanguage === 'ar' ? 'الجرعة:' : 'Dosage:'}</h5>
                    <p>${this.currentLanguage === 'ar' ? 
                        supplement.dosage_ar?.general || supplement.dosage.general : 
                        supplement.dosage.general
                    }</p>
                </div>
                ${supplement.fda_guidelines ? `
                    <div class="fda-guidelines">
                        <small><strong>FDA:</strong> ${supplement.fda_guidelines}</small>
                    </div>
                ` : ''}
                <button class="btn btn-outline" onclick="workoutGenerator.showSupplementDetails('${supplement.name}')">
                    ${this.currentLanguage === 'ar' ? 'المزيد من التفاصيل' : 'More Details'}
                </button>
            </div>
        `;
    }

    createStackCard(stack) {
        const name = this.currentLanguage === 'ar' ? stack.name_ar : stack.name;
        
        return `
            <div class="stack-card">
                <h4>${name}</h4>
                <p>${stack.description}</p>
                <div class="stack-supplements">
                    <strong>${this.currentLanguage === 'ar' ? 'يشمل:' : 'Includes:'}</strong>
                    <ul>
                        ${stack.supplements.map(sup => `<li>${sup.replace('_', ' ')}</li>`).join('')}
                    </ul>
                </div>
                <div class="stack-cost">
                    <strong>${this.currentLanguage === 'ar' ? 'التكلفة:' : 'Cost:'}</strong> ${stack.total_cost}
                </div>
            </div>
        `;
    }

    showSupplementDetails(supplementName) {
        // This would show a modal with detailed supplement information
        alert(`Detailed information for ${supplementName} coming soon!`);
    }

    showSubscriptionModal() {
        // This would show a subscription modal
        alert(this.currentLanguage === 'ar' ? 
            'نموذج الاشتراك قيد التطوير' : 
            'Subscription modal coming soon'
        );
    }

    async handleProductUpload(e) {
        e.preventDefault();
        
        // Validate product upload form using ValidationSystem
        const validationRules = [
            { field: 'product-name', required: true, name: 'Product Name' },
            { field: 'product-category', required: true, name: 'Category' },
            { field: 'product-description', required: true, name: 'Description' },
            { field: 'product-price', type: 'numeric', required: true, name: 'Price', min: 0 },
            { field: 'product-currency', required: true, name: 'Currency' },
            { field: 'product-brand', required: true, name: 'Brand' },
            { field: 'contact-name', type: 'name', required: true, name: 'Contact Name' },
            { field: 'contact-email', type: 'email', required: true, name: 'Contact Email' },
            { field: 'contact-phone', type: 'phone', required: true, name: 'Contact Phone' }
        ];
        
        if (!validation.validateForm('productUploadForm', validationRules)) {
            return;
        }
        
        const formData = new FormData(e.target);
        
        // Validate product images
        const imageFiles = formData.getAll('product-images');
        for (const file of imageFiles) {
            if (file.size > 0) {
                if (!validation.validateFile(file, 'product-images', ['image/jpeg', 'image/png', 'image/gif', 'image/webp'], 5 * 1024 * 1024)) {
                    validation.displayErrors();
                    return;
                }
            }
        }
        
        const productData = {
            id: `product-${Date.now()}`,
            timestamp: new Date().toISOString(),
            name: formData.get('product-name'),
            name_ar: formData.get('product-name-ar'),
            category: formData.get('product-category'),
            description: formData.get('product-description'),
            description_ar: formData.get('product-description-ar'),
            price: parseFloat(formData.get('product-price')),
            currency: formData.get('product-currency'),
            brand: formData.get('product-brand'),
            ingredients: formData.get('product-ingredients'),
            benefits: formData.get('product-benefits'),
            dosage: formData.get('product-dosage'),
            warnings: formData.get('product-warnings'),
            contact: {
                name: formData.get('contact-name'),
                email: formData.get('contact-email'),
                phone: formData.get('contact-phone'),
                company: formData.get('contact-company')
            },
            status: 'pending',
            images: []
        };
        
        // Handle image uploads
        for (const file of imageFiles) {
            if (file.size > 0) {
                try {
                    const imageData = await this.processImage(file);
                    productData.images.push(imageData);
                } catch (error) {
                    console.error('Error processing image:', error);
                }
            }
        }
        
        try {
            await this.savePendingProduct(productData);
            this.showSuccessMessage('Product submitted for review!');
            e.target.reset();
            validation.clearErrors();
        } catch (error) {
            console.error('Error submitting product:', error);
            this.showError('Failed to submit product. Please try again.');
        }
    }

    async processImage(file) {
        return new Promise((resolve, reject) => {
            const reader = new FileReader();
            reader.onload = (e) => {
                resolve({
                    name: file.name,
                    size: file.size,
                    type: file.type,
                    data: e.target.result,
                    timestamp: new Date().toISOString()
                });
            };
            reader.onerror = reject;
            reader.readAsDataURL(file);
        });
    }

    async savePendingProduct(productData) {
        // Save to localStorage (in real app, this would go to a server)
        const pendingProducts = JSON.parse(localStorage.getItem('pendingProducts') || '[]');
        pendingProducts.push(productData);
        localStorage.setItem('pendingProducts', JSON.stringify(pendingProducts));
        
        // Save to JSON file for admin review
        const filename = `pending-product-${productData.id}.json`;
        this.downloadJSON(productData, filename);
        
        // Simulate Google Forms integration
        this.submitToGoogleForms(productData);
    }

    submitToGoogleForms(productData) {
        // This would integrate with Google Forms API
        console.log('Submitting to Google Forms:', productData);
        
        // For now, just log the data that would be submitted
        const formData = {
            'Product Name': productData.name,
            'Category': productData.category,
            'Price': `${productData.price} ${productData.currency}`,
            'Contact Email': productData.contact.email,
            'Submission Date': productData.timestamp
        };
        
        console.log('Google Forms data:', formData);
    }

    showSuccessMessage(message) {
        const alertDiv = document.createElement('div');
        alertDiv.className = 'alert alert-success';
        alertDiv.textContent = message;
        document.body.appendChild(alertDiv);
        
        setTimeout(() => {
            document.body.removeChild(alertDiv);
        }, 3000);
    }

    showError(message) {
        const alertDiv = document.createElement('div');
        alertDiv.className = 'alert alert-error';
        alertDiv.textContent = message;
        document.body.appendChild(alertDiv);
        
        setTimeout(() => {
            document.body.removeChild(alertDiv);
        }, 5000);
    }

    updateLanguage() {
        // Update all translatable elements
        document.querySelectorAll('[data-translate]').forEach(element => {
            const key = element.getAttribute('data-translate');
            if (translations[this.currentLanguage] && translations[this.currentLanguage][key]) {
                element.textContent = translations[this.currentLanguage][key];
            }
        });
        
        // Update placeholders
        document.querySelectorAll('[data-translate-placeholder]').forEach(element => {
            const key = element.getAttribute('data-translate-placeholder');
            if (translations[this.currentLanguage] && translations[this.currentLanguage][key]) {
                element.placeholder = translations[this.currentLanguage][key];
            }
        });
        
        // Update RTL/LTR direction
        document.documentElement.dir = this.currentLanguage === 'ar' ? 'rtl' : 'ltr';
        
        // Refresh supplements display
        this.displaySupplements();
    }

    translate(key) {
        // Simple translation helper method
        const translations = {
            en: {
                rest_day: 'Rest Day',
                full_body_workout: 'Full Body Workout',
                upper_body_workout: 'Upper Body Workout',
                lower_body_workout: 'Lower Body Workout',
                push_workout: 'Push Workout',
                pull_workout: 'Pull Workout',
                leg_workout: 'Leg Workout',
                cardio_workout: 'Cardio Workout',
                workout: 'Workout'
            },
            ar: {
                rest_day: 'يوم راحة',
                full_body_workout: 'تمرين الجسم كامل',
                upper_body_workout: 'تمرين الجزء العلوي',
                lower_body_workout: 'تمرين الجزء السفلي',
                push_workout: 'تمرين الدفع',
                pull_workout: 'تمرين السحب',
                leg_workout: 'تمرين الأرجل',
                cardio_workout: 'تمرين الكارديو',
                workout: 'تمرين'
            }
        };
        
        return translations[this.currentLanguage]?.[key] || key;
    }
}

// Initialize the workout generator when the page loads
let workoutGenerator;
document.addEventListener('DOMContentLoaded', () => {
    workoutGenerator = new WorkoutGenerator();
});

// Export for global access
window.workoutGenerator = workoutGenerator;