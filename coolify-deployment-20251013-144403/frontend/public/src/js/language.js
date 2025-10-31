// Language support for Arabic/English
const translations = {
    en: {
        app_name: 'Nutrition Platform',
        nav_home: 'Home',
        nav_meals: 'Meals',
        nav_recipes: 'Recipes',
        nav_workouts: 'Workouts',
        nav_products: 'Products',
        
        // Home page cards
        card_health_diet: 'Health & Diet',
        card_health_diet_desc: 'Personalized nutrition plans and healthy meal recommendations',
        card_workout_injuries: 'Workout & Injuries',
        card_workout_injuries_desc: 'Exercise routines and injury prevention guidance',
        card_sports_games: 'Sports & Games',
        card_sports_games_desc: 'Sports activities and recreational games information',
        card_food_reviews: 'Food Reviews',
        card_food_reviews_desc: 'Food and beauty product reviews and ratings',
        card_drug_doses: 'Drug Doses',
        card_drug_doses_desc: 'Medication dosage information and guidelines',
        
        // Footer
        footer_contact: 'Contact Us',
        footer_social: 'Follow Us',
        footer_newsletter: 'Newsletter',
        footer_newsletter_desc: 'Subscribe to get the latest health tips',
        footer_email_placeholder: 'Enter your email',
        footer_subscribe: 'Subscribe',
        footer_phone: 'Phone',
        footer_email: 'Email',
        footer_address: 'Address',
        footer_newsletter_text: 'Subscribe to get the latest health tips and updates',
        footer_email_placeholder: 'Enter your email',
        footer_subscribe: 'Subscribe',
        footer_privacy: 'Privacy Policy',
        footer_terms: 'Terms of Service',
        
        // Diet Planning
        diet_planning_title: 'Diet Planning System',
        diet_planning_subtitle: 'Personalized nutrition plans for healthy individuals and patients',
        healthy_people_tab: 'Healthy People',
        patients_tab: 'Patients',
        personal_info: 'Personal Information',
        name: 'Name',
        age: 'Age',
        height: 'Height (cm)',
        weight: 'Weight (kg)',
        gender: 'Gender',
        select_gender: 'Select Gender',
        male: 'Male',
        female: 'Female',
        activity_preferences: 'Activity & Preferences',
        activity_level: 'Activity Level',
        select_activity: 'Select Activity Level',
        sedentary: 'Sedentary (little/no exercise)',
        light_active: 'Lightly Active (light exercise 1-3 days/week)',
        moderate_active: 'Moderately Active (moderate exercise 3-5 days/week)',
        very_active: 'Very Active (hard exercise 6-7 days/week)',
        extra_active: 'Extra Active (very hard exercise, physical job)',
        diet_goal: 'Diet Goal',
        select_goal: 'Select Goal',
        maintain_weight: 'Maintain Weight',
        lose_weight: 'Lose Weight',
        gain_weight: 'Gain Weight',
        build_muscle: 'Build Muscle',
        cuisine_preference: 'Cuisine Preference',
        select_cuisine: 'Select Cuisine',
        mediterranean: 'Mediterranean',
        middle_eastern: 'Middle Eastern',
        asian: 'Asian',
        american: 'American',
        european: 'European',
        indian: 'Indian',
        dietary_preferences: 'Dietary Preferences',
        select_dietary: 'Select Preference',
        no_restrictions: 'No Restrictions',
        vegetarian: 'Vegetarian',
        vegan: 'Vegan',
        keto: 'Keto',
        paleo: 'Paleo',
        low_carb: 'Low Carb',
        allergies: 'Allergies & Food Restrictions',
        calculate_plan: 'Calculate My Plan',
        your_results: 'Your Nutritional Analysis',
        weekly_meal_plan: '7-Day Meal Plan',
        shopping_list: 'Shopping List',
        download_pdf: 'Download PDF',
        monthly_subscription: 'Monthly Meal Plan Subscription',
        subscription_description: 'Get personalized meal plans delivered monthly with shopping lists and nutrition tracking.',
        subscribe_online: 'Subscribe Online',
        contact_whatsapp: 'Contact via WhatsApp',
        medical_info: 'Medical Information',
        disease_condition: 'Disease/Condition',
        select_condition: 'Select Condition',
        diabetes_type1: 'Type 1 Diabetes',
        diabetes_type2: 'Type 2 Diabetes',
        hypertension: 'Hypertension',
        heart_disease: 'Heart Disease',
        kidney_disease: 'Kidney Disease',
        liver_disease: 'Liver Disease',
        obesity: 'Obesity',
        celiac: 'Celiac Disease',
        ibs: 'IBS',
        other: 'Other',
        medications: 'Current Medications',
        lab_results: 'Recent Lab Results',
        doctor_restrictions: 'Doctor\'s Dietary Restrictions',
        create_medical_plan: 'Create Medical Plan',
        medical_analysis: 'Medical Nutritional Analysis',
        medical_meal_plan: 'Medical 7-Day Meal Plan',
        medical_subscription: 'Medical Meal Plan Subscription',
        medical_subscription_description: 'Get condition-specific meal plans with medical follow-up and monitoring.',
        subscribe_medical: 'Subscribe to Medical Plan',
        contact_medical_whatsapp: 'Medical Consultation',
        
        // Card descriptions
        card_health_description: 'Personalized nutrition plans and healthy meal recommendations',
        card_workout_description: 'Exercise routines and injury prevention guidance',
        card_sports_description: 'Sports nutrition and performance optimization',
        card_reviews_description: 'Product reviews and beauty recommendations',
        card_drugs_description: 'Medication guidance and dosage information',
        login: 'Login',
        register: 'Register',
        logout: 'Logout',
        profile: 'Profile',
        settings: 'Settings',
        hero_title: 'Your Health Journey Starts Here',
        hero_subtitle: 'Discover personalized nutrition plans, healthy recipes, and effective workouts tailored just for you.',
        get_started: 'Get Started',
        learn_more: 'Learn More',
        features_title: 'Everything You Need for a Healthy Lifestyle',
        features_subtitle: 'Our comprehensive platform provides all the tools you need to achieve your health and fitness goals.',
        feature_meals: 'Meal Planning',
        feature_meals_desc: 'Personalized meal plans based on your dietary preferences and health goals.',
        feature_recipes: 'Healthy Recipes',
        feature_recipes_desc: 'Discover delicious and nutritious recipes that fit your lifestyle.',
        feature_workouts: 'Workout Plans',
        feature_workouts_desc: 'Customized exercise routines to help you stay fit and active.',
        feature_supplements: 'Supplements',
        feature_supplements_desc: 'Expert recommendations for supplements to support your health.',
        email: 'Email',
        password: 'Password',
        first_name: 'First Name',
        last_name: 'Last Name',
        confirm_password: 'Confirm Password',
        date_of_birth: 'Date of Birth',
        gender: 'Gender',
        select_gender: 'Select Gender',
        male: 'Male',
        female: 'Female',
        other: 'Other',
        preferred_language: 'Preferred Language',
        english: 'English',
        arabic: 'العربية',
        forgot_password: 'Forgot Password?',
        remember_me: 'Remember Me',
        already_have_account: 'Already have an account?',
        dont_have_account: "Don't have an account?",
        sign_up_here: 'Sign up here',
        sign_in_here: 'Sign in here',
        loading: 'Loading...',
        save: 'Save',
        cancel: 'Cancel',
        delete: 'Delete',
        edit: 'Edit',
        add: 'Add',
        search: 'Search',
        filter: 'Filter',
        sort: 'Sort',
        view_all: 'View All',
        no_results: 'No results found',
        error_occurred: 'An error occurred',
        success: 'Success',
        warning: 'Warning',
        info: 'Information',
        
        // Workout Generator translations
        'workout_generator_title': 'Workout Generator',
        'workout_tab': 'Workout Generator',
        'supplements_tab': 'Supplements',
        'product_upload_tab': 'Product Upload',
        
        // Personal Information
        'personal_info_workout': 'Personal Information',
        'age_workout': 'Age',
        'weight_workout': 'Weight (kg)',
        'height_workout': 'Height (cm)',
        'gender_workout': 'Gender',
        'male_workout': 'Male',
        'female_workout': 'Female',
        'experience_level': 'Experience Level',
        'beginner': 'Beginner',
        'intermediate': 'Intermediate',
        'advanced': 'Advanced',
        
        // Fitness Goals
        'fitness_goals': 'Fitness Goals',
        'primary_goal': 'Primary Goal',
        'weight_loss': 'Weight Loss',
        'muscle_gain': 'Muscle Gain',
        'endurance': 'Endurance',
        'strength': 'Strength',
        'cutting': 'Cutting',
        'bulking': 'Bulking',
        'secondary_goals': 'Secondary Goals (Optional)',
        'improve_flexibility': 'Improve Flexibility',
        'increase_stamina': 'Increase Stamina',
        'better_posture': 'Better Posture',
        'stress_relief': 'Stress Relief',
        
        // Equipment
        'available_equipment': 'Available Equipment',
        'gym_access': 'Full Gym Access',
        'home_gym': 'Home Gym',
        'bodyweight_only': 'Bodyweight Only',
        'dumbbells': 'Dumbbells',
        'resistance_bands': 'Resistance Bands',
        'pull_up_bar': 'Pull-up Bar',
        'kettlebells': 'Kettlebells',
        'barbell': 'Barbell',
        
        // Schedule
        'workout_schedule': 'Workout Schedule',
        'days_per_week': 'Days per Week',
        'session_duration': 'Session Duration (minutes)',
        
        // Injury Assessment
        'injury_assessment': 'Injury Assessment',
        'current_injuries': 'Current Injuries or Limitations',
        'injury_body_part': 'Body Part',
        'lower_back': 'Lower Back',
        'knees': 'Knees',
        'shoulders': 'Shoulders',
        'wrists': 'Wrists',
        'ankles': 'Ankles',
        'neck': 'Neck',
        'injury_severity': 'Severity',
        'mild': 'Mild',
        'moderate': 'Moderate',
        'severe': 'Severe',
        'injury_restrictions': 'Specific Restrictions',
        'add_injury': 'Add Another Injury',
        
        // Buttons
        'generate_workout': 'Generate Workout Plan',
        'download_pdf': 'Download PDF',
        'share_workout': 'Share Workout',
        'save_workout': 'Save Workout',
        
        // Supplements
        'basic_supplements': 'Basic Supplements',
        'advanced_supplements': 'Advanced Supplements',
        'supplement_stacks': 'Supplement Stacks',
        'benefits': 'Benefits',
        'dosage': 'Dosage',
        'timing': 'Timing',
        'side_effects': 'Side Effects',
        'interactions': 'Interactions',
        'food_sources': 'Food Sources',
        'more_details': 'More Details',
        'subscribe_now': 'Subscribe Now',
        'subscription_required': 'Subscription Required',
        
        // Product Upload
        'product_upload_title': 'Product Upload',
        'product_information': 'Product Information',
        'product_name': 'Product Name',
        'product_category': 'Category',
        'protein_powder': 'Protein Powder',
        'pre_workout': 'Pre-Workout',
        'post_workout': 'Post-Workout',
        'vitamins': 'Vitamins',
        'minerals': 'Minerals',
        'herbs': 'Herbs',
        'other': 'Other',
        'product_description': 'Description',
        'product_price': 'Price',
        'product_currency': 'Currency',
        'product_brand': 'Brand',
        'product_ingredients': 'Ingredients',
        'product_benefits': 'Benefits',
        'product_dosage': 'Dosage Instructions',
        'product_warnings': 'Warnings & Side Effects',
        'product_images': 'Product Images',
        'contact_information': 'Contact Information',
        'contact_name': 'Full Name',
        'contact_email': 'Email',
        'contact_phone': 'Phone Number',
        'contact_company': 'Company Name',
        'submit_product': 'Submit for Review',
        
        // Results
        'workout_results': 'Workout Results',
        'your_workout_plan': 'Your Workout Plan',
        'nutrition_advice': 'Nutrition Advice',
        'recovery_tips': 'Recovery Tips',
        'progressive_overload': 'Progressive Overload Guidelines',
        'pre_workout_nutrition': 'Pre-Workout',
        'post_workout_nutrition': 'Post-Workout',
        'daily_nutrition': 'Daily',
        'hydration': 'Hydration',
        'sleep': 'Sleep',
        'active_recovery': 'Active Recovery',
        'exercises': 'exercises',
        'minutes': 'minutes',
        'rest_day': 'Rest Day',
        'relax_recover': 'Relax and Recover',
        'workout_generator': 'Workout Generator',
        'supplements': 'Supplements',
        'workout_generator_title': 'Workout Generator',
        'workout_generator_subtitle': 'Create personalized workout plans based on your goals, experience, and injury considerations',
        'workout_form_title': 'Workout Assessment Form',
        'fitness_goal': 'Fitness Goal',
        'weight_loss': 'Weight Loss',
        'muscle_gain': 'Muscle Gain',
        'endurance': 'Endurance',
        'cutting': 'Cutting',
        'bulking': 'Bulking',
        'general_fitness': 'General Fitness',
        'experience_level': 'Experience Level',
        'beginner': 'Beginner',
        'beginner_desc': '0-6 months experience',
        'intermediate': 'Intermediate',
        'intermediate_desc': '6 months - 2 years',
        'advanced': 'Advanced',
        'advanced_desc': '2+ years experience',
        'available_equipment': 'Available Equipment',
        'bodyweight': 'Bodyweight Only',
        'dumbbells': 'Dumbbells',
        'barbell': 'Barbell',
        'resistance_bands': 'Resistance Bands',
        'gym_access': 'Full Gym Access',
        'cardio_equipment': 'Cardio Equipment',
        'injury_assessment': 'Injury Assessment',
        'injury_disclaimer': 'Please select any current injuries or areas of concern. This will help us modify your workout plan accordingly.',
        'no_injuries': 'No Current Injuries',
        'lower_back': 'Lower Back',
        'knee': 'Knee',
        'shoulder': 'Shoulder',
        'wrist': 'Wrist',
        'ankle': 'Ankle',
        'additional_notes': 'Additional Notes',
        'additional_notes_placeholder': 'Any additional information about your fitness goals, preferences, or limitations...',
        'generate_workout': 'Generate My Workout Plan',
        'your_workout_plan': 'Your 7-Day Workout Plan',
        'generating_workout': 'Generating Your Personalized Workout...',
        'please_wait': 'This may take a few moments',
        
        // Nutrition Advice
        'weight_loss_calories': 'Maintain a caloric deficit of 300-500 calories',
        'weight_loss_protein': '1.2-1.6g per kg body weight',
        'weight_loss_carbs': 'Moderate carbs, focus on complex carbs',
        'weight_loss_fats': '20-30% of total calories',
        'weight_loss_timing': 'Eat smaller, frequent meals',
        'weight_loss_pre_workout': 'Light snack 30-60 minutes before',
        'weight_loss_post_workout': 'Protein within 30 minutes',
        
        'muscle_gain_calories': 'Maintain a caloric surplus of 200-400 calories',
        'muscle_gain_protein': '1.6-2.2g per kg body weight',
        'muscle_gain_carbs': 'Higher carbs for energy and recovery',
        'muscle_gain_fats': '25-35% of total calories',
        'muscle_gain_timing': 'Eat every 3-4 hours',
        'muscle_gain_pre_workout': 'Carbs and protein 1-2 hours before',
        'muscle_gain_post_workout': 'High protein meal within 2 hours',
        
        'endurance_calories': 'Match energy expenditure with intake',
        'endurance_protein': '1.2-1.4g per kg body weight',
        'endurance_carbs': 'High carbs (6-10g per kg body weight)',
        'endurance_fats': '20-25% of total calories',
        'endurance_timing': 'Fuel before, during, and after long sessions',
        'endurance_pre_workout': 'Carbs 1-4 hours before exercise',
        'endurance_post_workout': 'Carbs and protein 3:1 ratio',
        
        'cutting_calories': 'Aggressive caloric deficit of 500-750 calories',
        'cutting_protein': '2.0-2.5g per kg body weight',
        'cutting_carbs': 'Lower carbs, time around workouts',
        'cutting_fats': '15-25% of total calories',
        'cutting_timing': 'Consider intermittent fasting',
        'cutting_pre_workout': 'Minimal carbs if needed',
        'cutting_post_workout': 'Lean protein focus',
        
        'bulking_calories': 'Large caloric surplus of 500-1000 calories',
        'bulking_protein': '2.0-2.5g per kg body weight',
        'bulking_carbs': 'Very high carbs for maximum growth',
        'bulking_fats': '25-35% of total calories',
        'bulking_timing': 'Frequent large meals',
        'bulking_pre_workout': 'Large carb and protein meal',
        'bulking_post_workout': 'Immediate high-calorie meal',
        
        // Recovery Tips
        'recovery_sleep': '7-9 hours of quality sleep nightly',
        'recovery_stretching': '10-15 minutes post-workout stretching',
        'recovery_rest_days': 'Take at least 1-2 rest days per week',
        'recovery_active': 'Light walking or yoga on rest days',
        'recovery_stress': 'Manage stress through meditation',
        'recovery_massage': 'Self-massage or foam rolling',
        
        // Hydration Tips
        'hydration_daily': '2-3 liters of water daily',
        'hydration_pre_workout': '500ml 2 hours before exercise',
        'hydration_during_workout': '150-250ml every 15-20 minutes',
        'hydration_post_workout': '150% of fluid lost through sweat',
        'hydration_signs': 'Monitor urine color for hydration status',
        
        // Workout Types
        'rest_day': 'Rest Day',
        'active_recovery': 'Active Recovery',
        'light_stretching': 'Light Stretching',
        'hydration_focus': 'Hydration Focus',
        'full_body_workout': 'Full Body Workout',
        'upper_body_workout': 'Upper Body Workout',
        'lower_body_workout': 'Lower Body Workout',
        'push_workout': 'Push Workout',
        'pull_workout': 'Pull Workout',
        'leg_workout': 'Leg Workout',
        'cardio_workout': 'Cardio Workout',
        'workout': 'Workout'
    },
    ar: {
        app_name: 'منصة التغذية',
        nav_home: 'الرئيسية',
        nav_meals: 'الوجبات',
        nav_recipes: 'الوصفات',
        nav_workouts: 'التمارين',
        nav_products: 'المنتجات',
        
        // Home page cards
        card_health_diet: 'صحتك ووجباتك',
        card_health_diet_desc: 'خطط تغذية مخصصة وتوصيات وجبات صحية',
        card_workout_injuries: 'الجيم والإصابات',
        card_workout_injuries_desc: 'روتين التمارين وإرشادات الوقاية من الإصابات',
        card_sports_games: 'الرياضات والألعاب',
        card_sports_games_desc: 'معلومات الأنشطة الرياضية والألعاب الترفيهية',
        card_food_reviews: 'تقييمات الأكل ومنتجات التجميل',
        card_food_reviews_desc: 'مراجعات وتقييمات الطعام ومنتجات التجميل',
        card_drug_doses: 'جرعات الدواء',
        card_drug_doses_desc: 'معلومات جرعات الأدوية والإرشادات',
        
        // Footer
        footer_contact: 'اتصل بنا',
        footer_social: 'تابعنا',
        footer_newsletter: 'النشرة الإخبارية',
        footer_newsletter_desc: 'اشترك للحصول على أحدث النصائح الصحية',
        footer_email_placeholder: 'أدخل بريدك الإلكتروني',
        footer_subscribe: 'اشترك',
        footer_phone: 'الهاتف',
        footer_email: 'البريد الإلكتروني',
        footer_address: 'العنوان',
        footer_newsletter_text: 'اشترك للحصول على أحدث النصائح الصحية والتحديثات',
        footer_email_placeholder: 'أدخل بريدك الإلكتروني',
        footer_subscribe: 'اشتراك',
        footer_privacy: 'سياسة الخصوصية',
        footer_terms: 'شروط الخدمة',
        
        // Diet Planning
        diet_planning_title: 'نظام التخطيط الغذائي',
        diet_planning_subtitle: 'خطط تغذية شخصية للأشخاص الأصحاء والمرضى',
        healthy_people_tab: 'الأشخاص الأصحاء',
        patients_tab: 'المرضى',
        personal_info: 'المعلومات الشخصية',
        name: 'الاسم',
        age: 'العمر',
        height: 'الطول (سم)',
        weight: 'الوزن (كجم)',
        gender: 'الجنس',
        select_gender: 'اختر الجنس',
        male: 'ذكر',
        female: 'أنثى',
        activity_preferences: 'النشاط والتفضيلات',
        activity_level: 'مستوى النشاط',
        select_activity: 'اختر مستوى النشاط',
        sedentary: 'خامل (قليل أو بدون تمرين)',
        light_active: 'نشط قليلاً (تمرين خفيف 1-3 أيام/أسبوع)',
        moderate_active: 'نشط متوسط (تمرين متوسط 3-5 أيام/أسبوع)',
        very_active: 'نشط جداً (تمرين شاق 6-7 أيام/أسبوع)',
        extra_active: 'نشط للغاية (تمرين شاق جداً، عمل بدني)',
        diet_goal: 'هدف النظام الغذائي',
        select_goal: 'اختر الهدف',
        maintain_weight: 'الحفاظ على الوزن',
        lose_weight: 'فقدان الوزن',
        gain_weight: 'زيادة الوزن',
        build_muscle: 'بناء العضلات',
        cuisine_preference: 'تفضيل المطبخ',
        select_cuisine: 'اختر المطبخ',
        mediterranean: 'متوسطي',
        middle_eastern: 'شرق أوسطي',
        asian: 'آسيوي',
        american: 'أمريكي',
        european: 'أوروبي',
        indian: 'هندي',
        dietary_preferences: 'التفضيلات الغذائية',
        select_dietary: 'اختر التفضيل',
        no_restrictions: 'بدون قيود',
        vegetarian: 'نباتي',
        vegan: 'نباتي صرف',
        keto: 'كيتو',
        paleo: 'باليو',
        low_carb: 'قليل الكربوهيدرات',
        allergies: 'الحساسية وقيود الطعام',
        calculate_plan: 'احسب خطتي',
        your_results: 'تحليلك الغذائي',
        weekly_meal_plan: 'خطة الوجبات لـ 7 أيام',
        shopping_list: 'قائمة التسوق',
        download_pdf: 'تحميل PDF',
        monthly_subscription: 'اشتراك خطة الوجبات الشهرية',
        subscription_description: 'احصل على خطط وجبات شخصية شهرياً مع قوائم التسوق وتتبع التغذية.',
        subscribe_online: 'اشترك عبر الإنترنت',
        contact_whatsapp: 'تواصل عبر واتساب',
        medical_info: 'المعلومات الطبية',
        disease_condition: 'المرض/الحالة',
        select_condition: 'اختر الحالة',
        diabetes_type1: 'السكري النوع الأول',
        diabetes_type2: 'السكري النوع الثاني',
        hypertension: 'ارتفاع ضغط الدم',
        heart_disease: 'أمراض القلب',
        kidney_disease: 'أمراض الكلى',
        liver_disease: 'أمراض الكبد',
        obesity: 'السمنة',
        celiac: 'مرض السيلياك',
        ibs: 'متلازمة القولون العصبي',
        other: 'أخرى',
        medications: 'الأدوية الحالية',
        lab_results: 'نتائج المختبر الحديثة',
        doctor_restrictions: 'قيود الطبيب الغذائية',
        create_medical_plan: 'إنشاء خطة طبية',
        medical_analysis: 'التحليل الغذائي الطبي',
        medical_meal_plan: 'خطة الوجبات الطبية لـ 7 أيام',
        medical_subscription: 'اشتراك خطة الوجبات الطبية',
        medical_subscription_description: 'احصل على خطط وجبات خاصة بالحالة مع متابعة ومراقبة طبية.',
        subscribe_medical: 'اشترك في الخطة الطبية',
        contact_medical_whatsapp: 'استشارة طبية',
        
        // Card descriptions
        card_health_description: 'خطط تغذية شخصية وتوصيات وجبات صحية',
        card_workout_description: 'روتين التمارين وإرشادات الوقاية من الإصابات',
        card_sports_description: 'تغذية رياضية وتحسين الأداء',
        card_reviews_description: 'مراجعات المنتجات وتوصيات التجميل',
        card_drugs_description: 'إرشادات الأدوية ومعلومات الجرعات',
        login: 'تسجيل الدخول',
        register: 'إنشاء حساب',
        logout: 'تسجيل الخروج',
        profile: 'الملف الشخصي',
        settings: 'الإعدادات',
        hero_title: 'رحلتك الصحية تبدأ من هنا',
        hero_subtitle: 'اكتشف خطط التغذية المخصصة والوصفات الصحية والتمارين الفعالة المصممة خصيصاً لك.',
        get_started: 'ابدأ الآن',
        learn_more: 'اعرف المزيد',
        features_title: 'كل ما تحتاجه لنمط حياة صحي',
        features_subtitle: 'منصتنا الشاملة توفر جميع الأدوات التي تحتاجها لتحقيق أهدافك الصحية واللياقة البدنية.',
        feature_meals: 'تخطيط الوجبات',
        feature_meals_desc: 'خطط وجبات مخصصة بناءً على تفضيلاتك الغذائية وأهدافك الصحية.',
        feature_recipes: 'وصفات صحية',
        feature_recipes_desc: 'اكتشف وصفات لذيذة ومغذية تناسب نمط حياتك.',
        feature_workouts: 'خطط التمارين',
        feature_workouts_desc: 'روتين تمارين مخصص لمساعدتك على البقاء بصحة جيدة ونشاط.',
        feature_supplements: 'المكملات الغذائية',
        feature_supplements_desc: 'توصيات الخبراء للمكملات الغذائية لدعم صحتك.',
        email: 'البريد الإلكتروني',
        password: 'كلمة المرور',
        first_name: 'الاسم الأول',
        last_name: 'اسم العائلة',
        confirm_password: 'تأكيد كلمة المرور',
        date_of_birth: 'تاريخ الميلاد',
        gender: 'الجنس',
        select_gender: 'اختر الجنس',
        male: 'ذكر',
        female: 'أنثى',
        other: 'آخر',
        preferred_language: 'اللغة المفضلة',
        english: 'English',
        arabic: 'العربية',
        forgot_password: 'نسيت كلمة المرور؟',
        remember_me: 'تذكرني',
        already_have_account: 'لديك حساب بالفعل؟',
        dont_have_account: 'ليس لديك حساب؟',
        sign_up_here: 'سجل هنا',
        sign_in_here: 'سجل دخولك هنا',
        loading: 'جاري التحميل...',
        save: 'حفظ',
        cancel: 'إلغاء',
        delete: 'حذف',
        edit: 'تعديل',
        add: 'إضافة',
        search: 'بحث',
        filter: 'تصفية',
        sort: 'ترتيب',
        view_all: 'عرض الكل',
        no_results: 'لم يتم العثور على نتائج',
        error_occurred: 'حدث خطأ',
        success: 'نجح',
        warning: 'تحذير',
        info: 'معلومات',
        
        // Workout Generator Arabic translations
        'workout_generator_title': 'مولد التمارين',
        'workout_tab': 'مولد التمارين',
        'supplements_tab': 'المكملات الغذائية',
        'product_upload_tab': 'رفع المنتجات',
        
        // Personal Information
        'personal_info_workout': 'المعلومات الشخصية',
        'age_workout': 'العمر',
        'weight_workout': 'الوزن (كيلو)',
        'height_workout': 'الطول (سم)',
        'gender_workout': 'الجنس',
        'male_workout': 'ذكر',
        'female_workout': 'أنثى',
        'experience_level': 'مستوى الخبرة',
        'beginner': 'مبتدئ',
        'intermediate': 'متوسط',
        'advanced': 'متقدم',
        
        // Fitness Goals
        'fitness_goals': 'أهداف اللياقة',
        'primary_goal': 'الهدف الأساسي',
        'weight_loss': 'فقدان الوزن',
        'muscle_gain': 'بناء العضلات',
        'endurance': 'التحمل',
        'strength': 'القوة',
        'cutting': 'التنشيف',
        'bulking': 'الضخامة',
        'secondary_goals': 'أهداف ثانوية (اختيارية)',
        'improve_flexibility': 'تحسين المرونة',
        'increase_stamina': 'زيادة القدرة على التحمل',
        'better_posture': 'تحسين الوضعية',
        'stress_relief': 'تخفيف التوتر',
        
        // Equipment
        'available_equipment': 'المعدات المتاحة',
        'gym_access': 'صالة رياضية كاملة',
        'home_gym': 'صالة منزلية',
        'bodyweight_only': 'وزن الجسم فقط',
        'dumbbells': 'دمبل',
        'resistance_bands': 'أحزمة المقاومة',
        'pull_up_bar': 'عقلة',
        'kettlebells': 'كيتل بيل',
        'barbell': 'بار حديد',
        
        // Schedule
        'workout_schedule': 'جدول التمارين',
        'days_per_week': 'أيام في الأسبوع',
        'session_duration': 'مدة الجلسة (دقائق)',
        
        // Injury Assessment
        'injury_assessment': 'تقييم الإصابات',
        'current_injuries': 'الإصابات أو القيود الحالية',
        'injury_body_part': 'جزء الجسم',
        'lower_back': 'أسفل الظهر',
        'knees': 'الركبتين',
        'shoulders': 'الأكتاف',
        'wrists': 'المعصمين',
        'ankles': 'الكاحلين',
        'neck': 'الرقبة',
        'injury_severity': 'شدة الإصابة',
        'mild': 'خفيفة',
        'moderate': 'متوسطة',
        'severe': 'شديدة',
        'injury_restrictions': 'قيود محددة',
        'add_injury': 'إضافة إصابة أخرى',
        
        // Buttons
        'generate_workout': 'إنشاء خطة التمرين',
        'download_pdf': 'تحميل PDF',
        'share_workout': 'مشاركة التمرين',
        'save_workout': 'حفظ التمرين',
        
        // Supplements
        'basic_supplements': 'المكملات الأساسية',
        'advanced_supplements': 'المكملات المتقدمة',
        'supplement_stacks': 'مجموعات المكملات',
        'benefits': 'الفوائد',
        'dosage': 'الجرعة',
        'timing': 'التوقيت',
        'side_effects': 'الآثار الجانبية',
        'interactions': 'التفاعلات',
        'food_sources': 'المصادر الغذائية',
        'more_details': 'المزيد من التفاصيل',
        'subscribe_now': 'اشترك الآن',
        'subscription_required': 'يتطلب اشتراك',
        
        // Product Upload
        'product_upload_title': 'رفع المنتجات',
        'product_information': 'معلومات المنتج',
        'product_name': 'اسم المنتج',
        'product_category': 'الفئة',
        'protein_powder': 'بودرة البروتين',
        'pre_workout': 'ما قبل التمرين',
        'post_workout': 'ما بعد التمرين',
        'vitamins': 'الفيتامينات',
        'minerals': 'المعادن',
        'herbs': 'الأعشاب',
        'other': 'أخرى',
        'product_description': 'الوصف',
        'product_price': 'السعر',
        'product_currency': 'العملة',
        'product_brand': 'العلامة التجارية',
        'product_ingredients': 'المكونات',
        'product_benefits': 'الفوائد',
        'product_dosage': 'تعليمات الجرعة',
        'product_warnings': 'التحذيرات والآثار الجانبية',
        'product_images': 'صور المنتج',
        'contact_information': 'معلومات الاتصال',
        'contact_name': 'الاسم الكامل',
        'contact_email': 'البريد الإلكتروني',
        'contact_phone': 'رقم الهاتف',
        'contact_company': 'اسم الشركة',
        'submit_product': 'إرسال للمراجعة',
        
        // Results
        'workout_results': 'نتائج التمرين',
        'your_workout_plan': 'خطة التمرين الخاصة بك',
        'nutrition_advice': 'نصائح التغذية',
        'recovery_tips': 'نصائح التعافي',
        'progressive_overload': 'إرشادات التحميل التدريجي',
        'pre_workout_nutrition': 'قبل التمرين',
        'post_workout_nutrition': 'بعد التمرين',
        'daily_nutrition': 'يومياً',
        'hydration': 'الترطيب',
        'sleep': 'النوم',
        'active_recovery': 'التعافي النشط',
        'exercises': 'تمارين',
        'minutes': 'دقائق',
        'rest_day': 'يوم راحة',
        'relax_recover': 'استرخ وتعافى',
        'workout_generator': 'مولد التمارين',
        'supplements': 'المكملات الغذائية',
        'workout_generator_title': 'مولد التمارين',
        'workout_generator_subtitle': 'إنشاء خطط تمارين شخصية بناءً على أهدافك وخبرتك واعتبارات الإصابة',
        'workout_form_title': 'نموذج تقييم التمرين',
        'fitness_goal': 'هدف اللياقة',
        'weight_loss': 'فقدان الوزن',
        'muscle_gain': 'بناء العضلات',
        'endurance': 'التحمل',
        'cutting': 'التنشيف',
        'bulking': 'الضخامة',
        'general_fitness': 'اللياقة العامة',
        'experience_level': 'مستوى الخبرة',
        'beginner': 'مبتدئ',
        'beginner_desc': 'خبرة 0-6 أشهر',
        'intermediate': 'متوسط',
        'intermediate_desc': '6 أشهر - سنتان',
        'advanced': 'متقدم',
        'advanced_desc': 'خبرة أكثر من سنتين',
        'available_equipment': 'المعدات المتاحة',
        'bodyweight': 'وزن الجسم فقط',
        'dumbbells': 'دمبل',
        'barbell': 'بار حديد',
        'resistance_bands': 'أحزمة المقاومة',
        'gym_access': 'صالة رياضية كاملة',
        'cardio_equipment': 'معدات الكارديو',
        'injury_assessment': 'تقييم الإصابات',
        'injury_disclaimer': 'يرجى تحديد أي إصابات حالية أو مناطق قلق. سيساعدنا هذا في تعديل خطة التمرين وفقاً لذلك.',
        'no_injuries': 'لا توجد إصابات حالية',
        'lower_back': 'أسفل الظهر',
        'knee': 'الركبة',
        'shoulder': 'الكتف',
        'wrist': 'المعصم',
        'ankle': 'الكاحل',
        'additional_notes': 'ملاحظات إضافية',
        'additional_notes_placeholder': 'أي معلومات إضافية حول أهداف اللياقة أو التفضيلات أو القيود...',
        'generate_workout': 'إنشاء خطة التمرين الخاصة بي',
        'your_workout_plan': 'خطة التمرين لـ 7 أيام',
        'generating_workout': 'إنشاء التمرين الشخصي الخاص بك...',
        'please_wait': 'قد يستغرق هذا بضع لحظات',
        
        // Nutrition Advice
        'weight_loss_calories': 'حافظ على عجز في السعرات الحرارية 300-500 سعرة',
        'weight_loss_protein': '1.2-1.6 جرام لكل كيلو من وزن الجسم',
        'weight_loss_carbs': 'كربوهيدرات معتدلة، ركز على الكربوهيدرات المعقدة',
        'weight_loss_fats': '20-30% من إجمالي السعرات الحرارية',
        'weight_loss_timing': 'تناول وجبات صغيرة ومتكررة',
        'weight_loss_pre_workout': 'وجبة خفيفة قبل 30-60 دقيقة',
        'weight_loss_post_workout': 'بروتين خلال 30 دقيقة',
        
        'muscle_gain_calories': 'حافظ على فائض في السعرات الحرارية 200-400 سعرة',
        'muscle_gain_protein': '1.6-2.2 جرام لكل كيلو من وزن الجسم',
        'muscle_gain_carbs': 'كربوهيدرات أعلى للطاقة والتعافي',
        'muscle_gain_fats': '25-35% من إجمالي السعرات الحرارية',
        'muscle_gain_timing': 'تناول الطعام كل 3-4 ساعات',
        'muscle_gain_pre_workout': 'كربوهيدرات وبروتين قبل 1-2 ساعة',
        'muscle_gain_post_workout': 'وجبة عالية البروتين خلال ساعتين',
        
        'endurance_calories': 'طابق استهلاك الطاقة مع المدخول',
        'endurance_protein': '1.2-1.4 جرام لكل كيلو من وزن الجسم',
        'endurance_carbs': 'كربوهيدرات عالية (6-10 جرام لكل كيلو)',
        'endurance_fats': '20-25% من إجمالي السعرات الحرارية',
        'endurance_timing': 'تزود بالوقود قبل وأثناء وبعد الجلسات الطويلة',
        'endurance_pre_workout': 'كربوهيدرات قبل 1-4 ساعات من التمرين',
        'endurance_post_workout': 'كربوهيدرات وبروتين بنسبة 3:1',
        
        'cutting_calories': 'عجز قوي في السعرات الحرارية 500-750 سعرة',
        'cutting_protein': '2.0-2.5 جرام لكل كيلو من وزن الجسم',
        'cutting_carbs': 'كربوهيدرات أقل، توقيتها حول التمارين',
        'cutting_fats': '15-25% من إجمالي السعرات الحرارية',
        'cutting_timing': 'فكر في الصيام المتقطع',
        'cutting_pre_workout': 'كربوهيدرات قليلة إذا لزم الأمر',
        'cutting_post_workout': 'تركيز على البروتين الخالي من الدهون',
        
        'bulking_calories': 'فائض كبير في السعرات الحرارية 500-1000 سعرة',
        'bulking_protein': '2.0-2.5 جرام لكل كيلو من وزن الجسم',
        'bulking_carbs': 'كربوهيدرات عالية جداً للنمو الأقصى',
        'bulking_fats': '25-35% من إجمالي السعرات الحرارية',
        'bulking_timing': 'وجبات كبيرة ومتكررة',
        'bulking_pre_workout': 'وجبة كبيرة من الكربوهيدرات والبروتين',
        'bulking_post_workout': 'وجبة عالية السعرات فورية',
        
        // Recovery Tips
        'recovery_sleep': '7-9 ساعات من النوم الجيد ليلياً',
        'recovery_stretching': '10-15 دقيقة من التمدد بعد التمرين',
        'recovery_rest_days': 'خذ على الأقل 1-2 يوم راحة أسبوعياً',
        'recovery_active': 'مشي خفيف أو يوغا في أيام الراحة',
        'recovery_stress': 'إدارة التوتر من خلال التأمل',
        'recovery_massage': 'تدليك ذاتي أو استخدام الفوم رولر',
        
        // Hydration Tips
        'hydration_daily': '2-3 لتر من الماء يومياً',
        'hydration_pre_workout': '500 مل قبل ساعتين من التمرين',
        'hydration_during_workout': '150-250 مل كل 15-20 دقيقة',
        'hydration_post_workout': '150% من السوائل المفقودة عبر العرق',
        'hydration_signs': 'راقب لون البول لحالة الترطيب',
        
        // Workout Types
        'rest_day': 'يوم راحة',
        'active_recovery': 'تعافي نشط',
        'light_stretching': 'تمدد خفيف',
        'hydration_focus': 'تركيز على الترطيب',
        'full_body_workout': 'تمرين الجسم كامل',
        'upper_body_workout': 'تمرين الجزء العلوي',
        'lower_body_workout': 'تمرين الجزء السفلي',
        'push_workout': 'تمرين الدفع',
        'pull_workout': 'تمرين السحب',
        'leg_workout': 'تمرين الأرجل',
        'cardio_workout': 'تمرين الكارديو',
        'workout': 'تمرين'
    }
};

let currentLanguage = 'en';

function setLanguage(lang) {
    currentLanguage = lang;
    localStorage.setItem('language', lang);
    
    // Update HTML attributes
    const html = document.documentElement;
    const body = document.body;
    
    if (lang === 'ar') {
        html.setAttribute('lang', 'ar');
        html.setAttribute('dir', 'rtl');
        body.classList.add('rtl');
        
        // Enable RTL Bootstrap
        const bootstrapRTL = document.getElementById('bootstrap-rtl');
        if (bootstrapRTL) {
            bootstrapRTL.disabled = false;
        }
        
        // Disable LTR Bootstrap
        const bootstrapLTR = document.querySelector('link[href*="bootstrap.min.css"]:not([id="bootstrap-rtl"])');
        if (bootstrapLTR) {
            bootstrapLTR.disabled = true;
        }
    } else {
        html.setAttribute('lang', 'en');
        html.setAttribute('dir', 'ltr');
        body.classList.remove('rtl');
        
        // Disable RTL Bootstrap
        const bootstrapRTL = document.getElementById('bootstrap-rtl');
        if (bootstrapRTL) {
            bootstrapRTL.disabled = true;
        }
        
        // Enable LTR Bootstrap
        const bootstrapLTR = document.querySelector('link[href*="bootstrap.min.css"]:not([id="bootstrap-rtl"])');
        if (bootstrapLTR) {
            bootstrapLTR.disabled = false;
        }
    }
    
    // Update all translatable elements
    updateTranslations();
    
    // Update language buttons
    updateLanguageButtons();
    
    // Update workout generator language if it exists
    if (window.workoutGenerator) {
        window.workoutGenerator.setLanguage(lang);
    }
}

function updateTranslations() {
    const elements = document.querySelectorAll('[data-translate]');
    
    elements.forEach(element => {
        const key = element.getAttribute('data-translate');
        const translation = translations[currentLanguage][key];
        
        if (translation) {
            if (element.tagName === 'INPUT' && (element.type === 'submit' || element.type === 'button')) {
                element.value = translation;
            } else if (element.tagName === 'INPUT' && element.placeholder !== undefined) {
                element.placeholder = translation;
            } else {
                element.textContent = translation;
            }
        }
    });
}

function updateLanguageButtons() {
    const enButton = document.querySelector('button[onclick="setLanguage(\'en\')"]');
    const arButton = document.querySelector('button[onclick="setLanguage(\'ar\')"]');
    
    if (enButton && arButton) {
        // Remove active class from both
        enButton.classList.remove('btn-secondary');
        enButton.classList.add('btn-outline-secondary');
        arButton.classList.remove('btn-secondary');
        arButton.classList.add('btn-outline-secondary');
        
        // Add active class to current language
        if (currentLanguage === 'en') {
            enButton.classList.remove('btn-outline-secondary');
            enButton.classList.add('btn-secondary');
        } else {
            arButton.classList.remove('btn-outline-secondary');
            arButton.classList.add('btn-secondary');
        }
    }
}

function getCurrentLanguage() {
    return currentLanguage;
}

function translate(key) {
    return translations[currentLanguage][key] || key;
}

// Language switching function for global use
function switchLanguage(lang) {
    setLanguage(lang);
    document.documentElement.lang = lang;
    document.documentElement.dir = lang === 'ar' ? 'rtl' : 'ltr';
    
    // Update all translations on the page
    updateTranslations();
    updateLanguageButtons();
    
    console.log('Language switched to:', lang);
}

// Initialize language on page load
document.addEventListener('DOMContentLoaded', () => {
    const savedLanguage = localStorage.getItem('language') || 'en';
    setLanguage(savedLanguage);
});

// Export functions for global use
window.setLanguage = setLanguage;
window.getCurrentLanguage = getCurrentLanguage;
window.translate = translate;
window.switchLanguage = switchLanguage;