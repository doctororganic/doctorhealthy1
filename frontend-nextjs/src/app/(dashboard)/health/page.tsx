'use client';

import { useState } from 'react';

// Types
interface UserProfile {
  name: string;
  diseases: string[];
  medications: string[];
  complaint: string;
}

interface DiseaseInfo {
  name: string;
  description: string;
  symptoms: string[];
  dietaryRecommendations: string[];
  foodsToInclude: string[];
  foodsToAvoid: string[];
  lifestyleChanges: string[];
  medications: {
    name: string;
    purpose: string;
    sideEffects: string[];
    interactions: string[];
  }[];
}

// Disease data - in a real implementation, this would come from a database
const diseaseDatabase: Record<string, DiseaseInfo> = {
  'diabetes': {
    name: 'Diabetes',
    description: 'A chronic condition that affects how your body processes blood sugar (glucose).',
    symptoms: [
      'Increased thirst',
      'Frequent urination',
      'Extreme hunger',
      'Unexplained weight loss',
      'Fatigue',
      'Irritability',
      'Blurred vision',
      'Slow-healing sores',
      'Frequent infections'
    ],
    dietaryRecommendations: [
      'Eat regular meals with consistent carbohydrate intake',
      'Choose whole grains over refined grains',
      'Include lean proteins with each meal',
      'Eat plenty of non-starchy vegetables',
      'Limit sugary drinks and refined carbohydrates',
      'Choose healthy fats like nuts, seeds, avocado, and olive oil',
      'Monitor carbohydrate intake and adjust insulin or medication accordingly'
    ],
    foodsToInclude: [
      'Leafy greens (spinach, kale)',
      'Whole grains (oats, quinoa, brown rice)',
      'Lean proteins (chicken, fish, tofu)',
      'Beans and lentils',
      'Berries',
      'Nuts and seeds',
      'Greek yogurt',
      'Avocado',
      'Olive oil',
      'Cinnamon'
    ],
    foodsToAvoid: [
      'Sugary drinks (soda, sweetened teas, fruit juice)',
      'White bread, white rice, white pasta',
      'Processed snacks (chips, crackers, cookies)',
      'Fried foods',
      'Full-fat dairy products',
      'High-sugar fruits (mango, grapes, bananas in excess)',
      'Processed meats (bacon, sausage, hot dogs)'
    ],
    lifestyleChanges: [
      'Regular physical activity (30 minutes most days)',
      'Monitor blood sugar regularly',
      'Maintain a healthy weight',
      'Get enough sleep',
      'Manage stress',
      'Quit smoking',
      'Limit alcohol consumption'
    ],
    medications: [
      {
        name: 'Metformin',
        purpose: 'Reduce glucose production in the liver and improve insulin sensitivity',
        sideEffects: ['Nausea', 'Diarrhea', 'Stomach upset', 'Gas', 'Weakness'],
        interactions: ['Certain contrast agents', 'Excessive alcohol', 'Some antibiotics']
      },
      {
        name: 'Insulin',
        purpose: 'Help control blood sugar levels',
        sideEffects: ['Low blood sugar (hypoglycemia)', 'Weight gain', 'Allergic reactions', 'Injection site reactions'],
        interactions: ['Alcohol', 'Certain blood pressure medications', 'Steroids']
      }
    ]
  },
  'hypertension': {
    name: 'High Blood Pressure (Hypertension)',
    description: 'A condition in which the force of the blood against your artery walls is too high.',
    symptoms: [
      'Often no symptoms',
      'Headaches',
      'Shortness of breath',
      'Nosebleeds',
      'Vision problems',
      'Chest pain',
      'Irregular heartbeat',
      'Blood in urine',
      'Pounding in your chest, neck, or ears'
    ],
    dietaryRecommendations: [
      'Reduce sodium intake to less than 2,300mg per day',
      'Eat plenty of fruits and vegetables',
      'Choose whole grains over refined grains',
      'Include low-fat dairy products',
      'Limit saturated and trans fats',
      'Eat lean proteins (fish, poultry, beans, nuts)',
      'Limit added sugars and sugary drinks',
      'Follow the DASH eating plan (Dietary Approaches to Stop Hypertension)'
    ],
    foodsToInclude: [
      'Leafy greens (spinach, kale, collard greens)',
      'Berries (blueberries, strawberries)',
      'Beets',
      'Oats and barley',
      'Bananas',
      'Avocado',
      'Fatty fish (salmon, mackerel)',
      'Greek yogurt',
      'Seeds (pumpkin, flax, chia)',
      'Garlic',
      'Dark chocolate'
    ],
    foodsToAvoid: [
      'High-sodium foods (processed foods, canned soups, deli meats)',
      'Fried foods',
      'Fast food',
      'Processed meats (bacon, sausage, hot dogs)',
      'Full-fat dairy products',
      'Sugary drinks and desserts',
      'Excessive alcohol',
      'Caffeine in excess'
    ],
    lifestyleChanges: [
      'Regular physical activity (150 minutes of moderate exercise per week)',
      'Maintain a healthy weight',
      'Limit alcohol consumption',
      'Quit smoking',
      'Manage stress',
      'Get enough sleep',
      'Monitor blood pressure at home'
    ],
    medications: [
      {
        name: 'ACE Inhibitors (e.g., Lisinopril)',
        purpose: 'Relax blood vessels by blocking the formation of angiotensin II',
        sideEffects: ['Dry cough', 'Dizziness', 'Headache', 'Fatigue', 'Elevated potassium levels'],
        interactions: ['NSAIDs', 'Potassium supplements', 'Lithium']
      },
      {
        name: 'Calcium Channel Blockers (e.g., Amlodipine)',
        purpose: 'Relax blood vessels by preventing calcium from entering cells',
        sideEffects: ['Swelling in ankles and feet', 'Dizziness', 'Headache', 'Flushing', 'Palpitations'],
        interactions: ['Grapefruit juice', 'Some antibiotics', 'Antifungal medications']
      }
    ]
  },
  'heart_disease': {
    name: 'Heart Disease',
    description: 'A range of conditions that affect your heart, including coronary artery disease, heart rhythm problems, and heart defects.',
    symptoms: [
      'Chest pain, tightness, pressure, or discomfort',
      'Shortness of breath',
      'Pain in the neck, jaw, throat, upper abdomen or back',
      'Numbness or weakness in your arms or legs',
      'Extreme fatigue',
      'Irregular heartbeat',
      'Swelling in your legs, ankles, and feet'
    ],
    dietaryRecommendations: [
      'Follow a heart-healthy diet rich in fruits, vegetables, whole grains, and lean proteins',
      'Limit saturated and trans fats',
      'Reduce sodium intake',
      'Choose foods with omega-3 fatty acids',
      'Limit added sugars and refined grains',
      'Control portion sizes',
      'Limit alcohol consumption'
    ],
    foodsToInclude: [
      'Fatty fish (salmon, mackerel, sardines)',
      'Oats and barley',
      'Berries',
      'Dark leafy greens',
      'Avocado',
      'Nuts and seeds',
      'Legumes',
      'Olive oil',
      'Dark chocolate',
      'Green tea'
    ],
    foodsToAvoid: [
      'Processed meats (bacon, sausage, hot dogs)',
      'Refined carbohydrates (white bread, white rice, pasta)',
      'Sugary drinks and desserts',
      'Fried foods',
      'Full-fat dairy products',
      'Excessive alcohol',
      'High-sodium foods',
      'Trans fats (found in many packaged and fried foods)'
    ],
    lifestyleChanges: [
      'Regular physical activity (150 minutes of moderate exercise per week)',
      'Maintain a healthy weight',
      'Quit smoking',
      'Manage stress',
      'Get enough sleep',
      'Control blood pressure and cholesterol',
      'Limit alcohol consumption'
    ],
    medications: [
      {
        name: 'Statins (e.g., Atorvastatin)',
        purpose: 'Lower cholesterol levels in the blood',
        sideEffects: ['Muscle pain and weakness', 'Liver damage', 'Digestive problems', 'Increased blood sugar'],
        interactions: ['Grapefruit juice', 'Some antibiotics', 'Antifungal medications', 'Niacin supplements']
      },
      {
        name: 'Aspirin',
        purpose: 'Prevent blood clots that can cause heart attacks and strokes',
        sideEffects: ['Stomach irritation', 'Bleeding', 'Allergic reactions', 'Ringing in the ears'],
        interactions: ['Blood thinners', 'NSAIDs', 'Alcohol', 'Certain antidepressants']
      }
    ]
  },
  'obesity': {
    name: 'Obesity',
    description: 'A complex disease involving an excessive amount of body fat that increases the risk of health problems.',
    symptoms: [
      'Excess body fat, particularly around the waist',
      'Increased BMI (30 or higher)',
      'Breathing problems during sleep (sleep apnea)',
      'Joint pain',
      'Fatigue',
      'Excessive sweating',
      'Skin problems (folds infected)',
      'Shortness of breath with exertion'
    ],
    dietaryRecommendations: [
      'Focus on nutrient-dense, low-calorie foods',
      'Increase intake of fruits, vegetables, and fiber',
      'Choose lean proteins and whole grains',
      'Limit added sugars and refined carbohydrates',
      'Practice portion control',
      'Drink plenty of water',
      'Avoid liquid calories (sugary drinks, alcohol)',
      'Eat slowly and mindfully'
    ],
    foodsToInclude: [
      'Non-starchy vegetables (broccoli, spinach, cauliflower)',
      'Lean proteins (chicken breast, fish, tofu, beans, lentils)',
      'Whole grains (quinoa, oats, brown rice)',
      'Berries and other low-sugar fruits',
      'Nuts and seeds (in moderation)',
      'Greek yogurt',
      'Eggs',
      'Avocado',
      'Leafy greens',
      'Water and herbal tea'
    ],
    foodsToAvoid: [
      'Sugary drinks and desserts',
      'Fast food and processed meals',
      'White bread, white rice, white pasta',
      'Fried foods',
      'High-fat dairy products',
      'Processed meats',
      'Excessive oils and fats',
      'High-calorie coffee drinks',
      'Alcohol'
    ],
    lifestyleChanges: [
      'Regular physical activity (300 minutes of moderate exercise per week)',
      'Strength training at least twice a week',
      'Get enough sleep',
      'Manage stress',
      'Limit screen time',
      'Eat meals at regular times',
      'Keep a food journal'
    ],
    medications: [
      {
        name: 'Orlistat (Xenical)',
        purpose: 'Prevents absorption of some of the fat you eat',
        sideEffects: ['Oily spotting on underwear', 'Gas and oily discharge', 'More frequent bowel movements', 'Nausea'],
        interactions: ['Fat-soluble vitamins (A, D, E, K)', 'Blood thinners', 'Diabetes medications']
      },
      {
        name: 'Phentermine',
        purpose: 'Suppresses appetite',
        sideEffects: ['Increased blood pressure', 'Insomnia', 'Dry mouth', 'Constipation', 'Nervousness'],
        interactions: ['MAO inhibitors', 'Other weight loss medications', 'Alcohol', 'Antidepressants']
      }
    ]
  }
};

export default function HealthPage() {
  const [userProfile, setUserProfile] = useState<UserProfile>({
    name: '',
    diseases: [],
    medications: [],
    complaint: ''
  });

  const [selectedDisease, setSelectedDisease] = useState<string>('');
  const [diseaseInfo, setDiseaseInfo] = useState<DiseaseInfo | null>(null);
  const [showInfo, setShowInfo] = useState(false);
  const [loading, setLoading] = useState(false);

  // Disease options
  const diseaseOptions = [
    'diabetes',
    'hypertension',
    'heart_disease',
    'obesity',
    'asthma',
    'arthritis',
    'osteoporosis',
    'depression',
    'anxiety',
    'migraine',
    'eczema',
    'psoriasis',
    'ibs',
    'celiac_disease',
    'liver_disease',
    'kidney_disease'
  ];

  // Complaint options
  const complaintOptions = [
    'fatigue',
    'low_energy',
    'poor_sleep',
    'stress',
    'anxiety',
    'digestive_issues',
    'headaches',
    'joint_pain',
    'muscle_soreness',
    'weight_gain',
    'weight_loss_difficulty',
    'low_libido',
    'memory_issues',
    'mood_swings',
    'hormonal_imbalance'
  ];

  // Get disease information
  const getDiseaseInfo = (disease: string) => {
    setLoading(true);
    
    try {
      if (diseaseDatabase[disease]) {
        setDiseaseInfo(diseaseDatabase[disease]);
        setShowInfo(true);
      }
    } catch (error) {
      console.error('Failed to get disease information:', error);
    } finally {
      setLoading(false);
    }
  };

  // Handle disease selection
  const handleDiseaseSelect = (disease: string) => {
    setSelectedDisease(disease);
    setUserProfile(prev => ({
      ...prev,
      diseases: [...prev.diseases, disease]
    }));
    getDiseaseInfo(disease);
  };

  return (
    <div className="min-h-screen bg-gradient-to-b from-white to-yellow-50 p-4">
      <div className="max-w-4xl mx-auto">
        <div className="bg-white rounded-xl shadow-lg p-6 mb-6">
          <h1 className="text-2xl font-bold text-blue-600 mb-6">Diseases and Healthy-Lifestyle</h1>
          
          <form className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Your Name</label>
              <input
                type="text"
                value={userProfile.name}
                onChange={(e) => setUserProfile({...userProfile, name: e.target.value})}
                className="w-full p-2 border border-gray-300 rounded-md"
              />
            </div>
            
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Diseases</label>
              <select
                multiple
                value={userProfile.diseases}
                onChange={(e) => setUserProfile({...userProfile, diseases: Array.from(e.target.selectedOptions, option => option.value)})}
                className="w-full p-2 border border-gray-300 rounded-md"
              >
                {diseaseOptions.map(disease => (
                  <option key={disease} value={disease}>
                    {disease.replace('_', ' ').replace(/\b\w/g, l => l.toUpperCase())}
                  </option>
                ))}
              </select>
            </div>
            
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Medications</label>
              <textarea
                value={userProfile.medications.join(', ')}
                onChange={(e) => setUserProfile({...userProfile, medications: e.target.value.split(',').map(m => m.trim())})}
                className="w-full p-2 border border-gray-300 rounded-md"
                rows={2}
                placeholder="e.g., Metformin, Lisinopril, Atorvastatin"
              />
            </div>
            
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Current Complaint</label>
              <select
                value={userProfile.complaint}
                onChange={(e) => setUserProfile({...userProfile, complaint: e.target.value})}
                className="w-full p-2 border border-gray-300 rounded-md"
              >
                <option value="">Select a complaint</option>
                {complaintOptions.map(complaint => (
                  <option key={complaint} value={complaint}>
                    {complaint.replace('_', ' ').replace(/\b\w/g, l => l.toUpperCase())}
                  </option>
                ))}
              </select>
            </div>
            
            <button
              type="button"
              onClick={() => {
                if (userProfile.diseases.length > 0) {
                  getDiseaseInfo(userProfile.diseases[0]);
                }
              }}
              disabled={loading || userProfile.diseases.length === 0}
              className="w-full btn-secondary"
            >
              {loading ? 'Loading...' : 'Get Health Advice'}
            </button>
          </form>
        </div>
        
        {/* Disease Information */}
        {showInfo && diseaseInfo && (
          <div className="bg-white rounded-xl shadow-lg p-6">
            <h2 className="text-xl font-bold text-blue-600 mb-4">
              Health Information for: {diseaseInfo.name}
            </h2>
            
            <div className="mb-6">
              <p className="text-gray-600 mb-4">{diseaseInfo.description}</p>
              
              <div className="mb-4">
                <h3 className="font-semibold text-gray-800 mb-2">Common Symptoms</h3>
                <ul className="list-disc list-inside text-gray-600 ml-4">
                  {diseaseInfo.symptoms.map((symptom, index) => (
                    <li key={index}>{symptom}</li>
                  ))}
                </ul>
              </div>
            </div>
            
            <div className="mb-6">
              <h3 className="font-semibold text-gray-800 mb-2">Dietary Recommendations</h3>
              <ul className="list-disc list-inside text-gray-600 ml-4">
                {diseaseInfo.dietaryRecommendations.map((recommendation, index) => (
                  <li key={index}>{recommendation}</li>
                ))}
              </ul>
            </div>
            
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
              <div>
                <h3 className="font-semibold text-gray-800 mb-2">Foods to Include</h3>
                <ul className="list-disc list-inside text-gray-600 ml-4">
                  {diseaseInfo.foodsToInclude.map((food, index) => (
                    <li key={index}>{food}</li>
                  ))}
                </ul>
              </div>
              
              <div>
                <h3 className="font-semibold text-gray-800 mb-2">Foods to Avoid</h3>
                <ul className="list-disc list-inside text-gray-600 ml-4">
                  {diseaseInfo.foodsToAvoid.map((food, index) => (
                    <li key={index}>{food}</li>
                  ))}
                </ul>
              </div>
            </div>
            
            <div className="mb-6">
              <h3 className="font-semibold text-gray-800 mb-2">Lifestyle Changes</h3>
              <ul className="list-disc list-inside text-gray-600 ml-4">
                {diseaseInfo.lifestyleChanges.map((change, index) => (
                  <li key={index}>{change}</li>
                ))}
              </ul>
            </div>
            
            <div className="mb-6">
              <h3 className="font-semibold text-gray-800 mb-2">Medications Information</h3>
              <div className="space-y-4">
                {diseaseInfo.medications.map((med, index) => (
                  <div key={index} className="border-l-4 border-blue-500 pl-4">
                    <h4 className="font-medium text-gray-800">{med.name}</h4>
                    <p className="text-sm text-gray-600 mb-2">
                      <span className="font-medium">Purpose:</span> {med.purpose}
                    </p>
                    <div className="mb-2">
                      <span className="font-medium text-sm">Side Effects:</span>
                      <ul className="list-disc list-inside text-gray-600 ml-4 text-sm">
                        {med.sideEffects.map((effect, i) => (
                          <li key={i}>{effect}</li>
                        ))}
                      </ul>
                    </div>
                    <div>
                      <span className="font-medium text-sm">Interactions:</span>
                      <ul className="list-disc list-inside text-gray-600 ml-4 text-sm">
                        {med.interactions.map((interaction, i) => (
                          <li key={i}>{interaction}</li>
                        ))}
                      </ul>
                    </div>
                  </div>
                ))}
              </div>
            </div>
            
            {/* Disease-specific sections for multiple diseases */}
            {userProfile.diseases.length > 1 && (
              <div className="border-t pt-4">
                <h3 className="font-semibold text-gray-800 mb-2">Additional Conditions</h3>
                <p className="text-sm text-gray-600">
                  You have selected multiple conditions. It's important to consult with your healthcare provider for a comprehensive treatment plan that considers all your health conditions and medications.
                </p>
              </div>
            )}
          </div>
        )}
        
        {/* Medical Disclaimer */}
        <div className="disclaimer mt-8">
          <p>
            <strong>Disclaimer:</strong> This site is to update you with information and guide you to useful advice for the purpose of education and awareness and is not a substitute for a doctor's visit. Always consult with a qualified healthcare provider for diagnosis and treatment of medical conditions.
          </p>
          <p>
            <strong>Medication Safety:</strong> Never change your medication dosage or stop taking prescribed medications without consulting your healthcare provider. Inform your healthcare provider about all medications, supplements, and herbal remedies you are taking to avoid potential interactions.
          </p>
          <p>
            <strong>Emergency Information:</strong> If you are experiencing severe symptoms such as chest pain, difficulty breathing, severe headache, or other emergency symptoms, call emergency services immediately or go to the nearest emergency department.
          </p>
        </div>
      </div>
    </div>
  );
}