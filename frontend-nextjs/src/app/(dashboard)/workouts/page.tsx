'use client';

import { useState } from 'react';

// Types
interface UserProfile {
  name: string;
  weight: number;
  height: number;
  gender: 'male' | 'female';
  activityLevel: 'sedentary' | 'light' | 'moderate' | 'active' | 'very_active';
  workoutGoal: string;
  workoutLocation: 'gym' | 'home';
  injuries: string[];
  complaints: string[];
}

interface Exercise {
  name: string;
  type: string;
  sets: number;
  reps: number;
  rest: number;
  description: string;
  commonMistakes: string[];
  alternative: {
    name: string;
    type: string;
    sets: number;
    reps: number;
    rest: number;
    description: string;
    commonMistakes: string[];
  };
}

interface InjuryAdvice {
  injury: string;
  advice: string;
  exercises: string[];
  treatments: string[];
  restrictions: string[];
}

interface ComplaintSolution {
  complaint: string;
  advice: string;
  nutritionalRecommendations: string[];
  supplements: {
    name: string;
    dosage: string;
    frequency: string;
  }[];
}

export default function WorkoutsPage() {
  const [userProfile, setUserProfile] = useState<UserProfile>({
    name: '',
    weight: 70,
    height: 170,
    gender: 'male',
    activityLevel: 'moderate',
    workoutGoal: 'muscle_gain',
    workoutLocation: 'gym',
    injuries: [],
    complaints: []
  });

  const [workoutPlan, setWorkoutPlan] = useState<Exercise[]>([]);
  const [injuryAdvice, setInjuryAdvice] = useState<InjuryAdvice | null>(null);
  const [complaintSolution, setComplaintSolution] = useState<ComplaintSolution | null>(null);
  const [showPlan, setShowPlan] = useState(false);
  const [loading, setLoading] = useState(false);

  // Workout types
  const workoutTypes = [
    'strength',
    'cardio',
    'flexibility',
    'balance',
    'functional',
    'plyometric',
    'endurance',
    'hiit'
  ];

  // Complaints from metabolism
  const complaintsList = [
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

  // Injury types
  const injuryTypes = [
    'back_pain',
    'knee_pain',
    'shoulder_pain',
    'ankle_sprain',
    'wrist_pain',
    'neck_pain',
    'hip_pain',
    'elbow_pain',
    'muscle_strain',
    'tendonitis',
    'plantar_fasciitis',
    'rotator_cuff_injury',
    'herniated_disc'
  ];

  // Workout goals
  const workoutGoals = [
    'muscle_gain',
    'fat_loss',
    'endurance',
    'strength',
    'flexibility',
    'athletic_performance',
    'rehabilitation',
    'general_fitness',
    'weight_management',
    'stress_relief'
  ];

  // Generate workout plan
  const generateWorkoutPlan = () => {
    setLoading(true);
    
    try {
      // Generate exercises based on user profile
      const exercises = generateExercises();
      setWorkoutPlan(exercises);
      setShowPlan(true);
    } catch (error) {
      console.error('Failed to generate workout plan:', error);
    } finally {
      setLoading(false);
    }
  };

  // Generate exercises based on user profile
  const generateExercises = (): Exercise[] => {
    const { workoutGoal, workoutLocation, injuries, gender } = userProfile;
    
    // Sample exercises - in a real implementation, this would come from a database
    const exercisePool: Exercise[] = [
      {
        name: 'Bench Press',
        type: 'strength',
        sets: 4,
        reps: 10,
        rest: 90,
        description: 'Lie on a bench and press a barbell from chest to full arm extension.',
        commonMistakes: ['Arching your back', 'Bouncing the bar off your chest', 'Not using full range of motion'],
        alternative: {
          name: 'Dumbbell Press',
          type: 'strength',
          sets: 3,
          reps: 12,
          rest: 60,
          description: 'Lie on a bench and press dumbbells from chest to full arm extension.',
          commonMistakes: ['Uneven pressing motion', 'Dropping dumbbells too low', 'Flaring elbows out too wide']
        }
      },
      {
        name: 'Squats',
        type: 'strength',
        sets: 4,
        reps: 12,
        rest: 90,
        description: 'Stand with feet shoulder-width apart and lower your body until thighs are parallel to the floor.',
        commonMistakes: ['Knees caving inward', 'Heels lifting off the floor', 'Leaning too far forward'],
        alternative: {
          name: 'Goblet Squats',
          type: 'strength',
          sets: 3,
          reps: 15,
          rest: 60,
          description: 'Hold a dumbbell vertically against your chest and perform squats.',
          commonMistakes: ['Elbows flaring out', 'Not keeping chest up', 'Depth inconsistent']
        }
      },
      {
        name: 'Deadlifts',
        type: 'strength',
        sets: 4,
        reps: 8,
        rest: 120,
        description: 'Lift a barbell from the floor to a standing position, keeping your back straight.',
        commonMistakes: ['Rounding your back', 'Using momentum instead of strength', 'Not engaging your core'],
        alternative: {
          name: 'Romanian Deadlifts',
          type: 'strength',
          sets: 3,
          reps: 12,
          rest: 90,
          description: 'Hold a barbell in front of your thighs and hinge at the hips, lowering the bar while keeping your back straight.',
          commonMistakes: ['Rounding your back', 'Bending knees too much', 'Going too low']
        }
      },
      {
        name: 'Pull-ups',
        type: 'strength',
        sets: 3,
        reps: 8,
        rest: 90,
        description: 'Hang from a bar and pull your body up until your chin is over the bar.',
        commonMistakes: ['Using momentum', 'Not going through full range of motion', 'Kipping too early'],
        alternative: {
          name: 'Lat Pulldowns',
          type: 'strength',
          sets: 3,
          reps: 12,
          rest: 60,
          description: 'Sit at a lat pulldown machine and pull the bar down to your chest.',
          commonMistakes: ['Leaning back too far', 'Using body weight', 'Not controlling the negative']
        }
      },
      {
        name: 'Running',
        type: 'cardio',
        sets: 1,
        reps: 30,
        rest: 0,
        description: 'Run at a steady pace for 30 minutes.',
        commonMistakes: ['Overstriding', 'Poor running form', 'Not breathing properly'],
        alternative: {
          name: 'Cycling',
          type: 'cardio',
          sets: 1,
          reps: 45,
          rest: 0,
          description: 'Cycle at a moderate intensity for 45 minutes.',
          commonMistakes: ['Resistance too low', 'Poor posture', 'Not adjusting seat height']
        }
      },
      {
        name: 'Push-ups',
        type: 'strength',
        sets: 3,
        reps: 15,
        rest: 60,
        description: 'Start in a plank position and lower your body until your chest nearly touches the floor.',
        commonMistakes: ['Sagging hips', 'Not lowering far enough', 'Elbows flaring out'],
        alternative: {
          name: 'Knee Push-ups',
          type: 'strength',
          sets: 3,
          reps: 20,
          rest: 60,
          description: 'Start in a plank position with knees on the floor and lower your body.',
          commonMistakes: ['Hips too high', 'Not lowering far enough', 'Breaking straight line from head to knees']
        }
      }
    ];

    // Filter exercises based on injuries
    const filteredExercises = exercisePool.filter(exercise => {
      // Filter out exercises that would aggravate injuries
      if (injuries.includes('knee_pain') && (exercise.name === 'Squats' || exercise.name === 'Deadlifts')) {
        return false;
      }
      if (injuries.includes('back_pain') && (exercise.name === 'Deadlifts' || exercise.name === 'Bench Press')) {
        return false;
      }
      if (injuries.includes('shoulder_pain') && exercise.name === 'Bench Press') {
        return false;
      }
      return true;
    });

    // Select exercises based on goal
    let selectedExercises: Exercise[] = [];
    
    if (workoutGoal === 'muscle_gain' || workoutGoal === 'strength') {
      selectedExercises = filteredExercises.filter(ex => ex.type === 'strength').slice(0, 4);
    } else if (workoutGoal === 'fat_loss' || workoutGoal === 'endurance') {
      selectedExercises = [
        filteredExercises.find(ex => ex.type === 'cardio') || filteredExercises[4],
        ...filteredExercises.filter(ex => ex.type === 'strength').slice(0, 3)
      ];
    } else {
      selectedExercises = filteredExercises.slice(0, 4);
    }

    // Adjust for home vs gym workout
    if (workoutLocation === 'home') {
      selectedExercises = selectedExercises.map(exercise => {
        // Use bodyweight alternatives for home workouts
        if (exercise.name === 'Bench Press') {
          return {
            ...exercise,
            name: 'Push-ups',
            sets: 3,
            reps: 15,
            description: exercise.alternative.description,
            commonMistakes: exercise.alternative.commonMistakes
          };
        }
        return exercise;
      });
    }

    return selectedExercises;
  };

  // Get injury advice
  const getInjuryAdvice = () => {
    const { injuries } = userProfile;
    
    if (injuries.length === 0) return;
    
    // Sample injury advice - in a real implementation, this would come from a database
    const injuryAdvicePool: Record<string, InjuryAdvice> = {
      'knee_pain': {
        injury: 'Knee Pain',
        advice: 'Avoid high-impact activities and exercises that put stress on the knees. Focus on low-impact cardio and strengthening exercises for the muscles supporting the knees.',
        exercises: ['Swimming', 'Cycling', 'Leg Press', 'Hamstring Curls', 'Quadriceps Extensions (with limited range)'],
        treatments: ['RICE (Rest, Ice, Compression, Elevation)', 'Physical therapy', 'Anti-inflammatory medication'],
        restrictions: ['Avoid deep squats', 'No high-impact activities', 'Avoid exercises that cause pain']
      },
      'back_pain': {
        injury: 'Back Pain',
        advice: 'Focus on core strengthening exercises and proper posture. Avoid exercises that put excessive strain on the spine.',
        exercises: ['Plank', 'Bird Dog', 'Cat-Cow Stretch', 'Glute Bridges', 'Partial Range of Motion Deadlifts'],
        treatments: ['Heat therapy', 'Gentle stretching', 'Core strengthening exercises', 'Massage therapy'],
        restrictions: ['Avoid heavy lifting', 'No exercises that round the back', 'Avoid high-impact activities']
      },
      'shoulder_pain': {
        injury: 'Shoulder Pain',
        advice: 'Focus on exercises that strengthen the rotator cuff and improve shoulder mobility. Avoid overhead pressing movements.',
        exercises: ['Band Pull-Aparts', 'External Rotation', 'Scaption', 'Wall Slides', 'Pendulum Swings'],
        treatments: ['Physical therapy', 'Anti-inflammatory medication', 'Heat or ice therapy', 'Rest from aggravating activities'],
        restrictions: ['Avoid overhead presses', 'No heavy lifting with painful arm', 'Avoid exercises that cause pain']
      }
    };

    // Get advice for the first injury
    const firstInjury = injuries[0];
    if (injuryAdvicePool[firstInjury]) {
      setInjuryAdvice(injuryAdvicePool[firstInjury]);
    }
  };

  // Get complaint solution
  const getComplaintSolution = () => {
    const { complaints, weight, gender } = userProfile;
    
    if (complaints.length === 0) return;
    
    // Sample complaint solutions - in a real implementation, this would come from a database
    const complaintSolutionsPool: Record<string, ComplaintSolution> = {
      'fatigue': {
        complaint: 'Fatigue',
        advice: 'Focus on balanced nutrition, adequate sleep, and moderate exercise. Avoid excessive caffeine and processed foods.',
        nutritionalRecommendations: [
          'Eat whole foods rich in iron and B vitamins',
          'Stay hydrated with water throughout the day',
          'Include complex carbohydrates for sustained energy',
          'Ensure adequate protein intake'
        ],
        supplements: [
          { name: 'Iron', dosage: '18mg', frequency: 'daily' },
          { name: 'Vitamin B12', dosage: '1000mcg', frequency: 'daily' },
          { name: 'Vitamin D', dosage: '1000-2000 IU', frequency: 'daily' },
          { name: 'Magnesium', dosage: '200-400mg', frequency: 'daily' }
        ]
      },
      'stress': {
        complaint: 'Stress',
        advice: 'Practice stress management techniques like meditation, deep breathing, and regular exercise. Limit caffeine and alcohol intake.',
        nutritionalRecommendations: [
          'Eat foods rich in magnesium and B vitamins',
          'Include omega-3 fatty acids in your diet',
          'Avoid excessive caffeine and sugar',
          'Stay hydrated and eat regular meals'
        ],
        supplements: [
          { name: 'Ashwagandha', dosage: '300-600mg', frequency: 'twice daily' },
          { name: 'L-Theanine', dosage: '100-200mg', frequency: 'twice daily' },
          { name: 'Magnesium', dosage: '200-400mg', frequency: 'daily' },
          { name: 'B-Complex', dosage: '50mg', frequency: 'daily' }
        ]
      },
      'weight_gain': {
        complaint: 'Weight Gain',
        advice: 'Focus on portion control, whole foods, and regular physical activity. Limit processed foods, sugar, and unhealthy fats.',
        nutritionalRecommendations: [
          'Eat more fiber-rich foods to promote satiety',
          'Include lean protein sources with each meal',
          'Choose whole grains over refined grains',
          'Eat plenty of vegetables and limit sugary drinks'
        ],
        supplements: [
          { name: 'Green Tea Extract', dosage: '500mg', frequency: 'daily' },
          { name: 'Probiotics', dosage: '10-20 billion CFU', frequency: 'daily' },
          { name: 'Fiber Supplement', dosage: '5-10g', frequency: 'daily' },
          { name: 'Chromium Picolinate', dosage: '200-400mcg', frequency: 'daily' }
        ]
      }
    };

    // Get solution for the first complaint
    const firstComplaint = complaints[0];
    if (complaintSolutionsPool[firstComplaint]) {
      setComplaintSolution(complaintSolutionsPool[firstComplaint]);
    }
  };

  // Handle form submission
  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    generateWorkoutPlan();
    getInjuryAdvice();
    getComplaintSolution();
  };

  return (
    <div className="min-h-screen bg-gradient-to-b from-white to-yellow-50 p-4">
      <div className="max-w-4xl mx-auto">
        <div className="bg-white rounded-xl shadow-lg p-6 mb-6">
          <h1 className="text-2xl font-bold text-blue-600 mb-6">Workouts and Injuries</h1>
          
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Name</label>
                <input
                  type="text"
                  value={userProfile.name}
                  onChange={(e) => setUserProfile({...userProfile, name: e.target.value})}
                  className="w-full p-2 border border-gray-300 rounded-md"
                />
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Workout Location</label>
                <select
                  value={userProfile.workoutLocation}
                  onChange={(e) => setUserProfile({...userProfile, workoutLocation: e.target.value as any})}
                  className="w-full p-2 border border-gray-300 rounded-md"
                >
                  <option value="gym">Gym</option>
                  <option value="home">Home</option>
                </select>
              </div>
              
              <div className="grid grid-cols-3 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Weight (kg)</label>
                  <input
                    type="number"
                    value={userProfile.weight}
                    onChange={(e) => {
                      const num = parseFloat(e.target.value);
                      setUserProfile({...userProfile, weight: isNaN(num) ? userProfile.weight : Math.max(1, num)});
                    }}
                    className="w-full p-2 border border-gray-300 rounded-md"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Height (cm)</label>
                  <input
                    type="number"
                    value={userProfile.height}
                    onChange={(e) => {
                      const num = parseInt(e.target.value);
                      setUserProfile({...userProfile, height: isNaN(num) ? userProfile.height : Math.max(1, num)});
                    }}
                    className="w-full p-2 border border-gray-300 rounded-md"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Gender</label>
                  <select
                    value={userProfile.gender}
                    onChange={(e) => setUserProfile({...userProfile, gender: e.target.value as any})}
                    className="w-full p-2 border border-gray-300 rounded-md"
                  >
                    <option value="male">Male</option>
                    <option value="female">Female</option>
                  </select>
                </div>
              </div>
            </div>
            
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Workout Goal</label>
              <select
                value={userProfile.workoutGoal}
                onChange={(e) => setUserProfile({...userProfile, workoutGoal: e.target.value})}
                className="w-full p-2 border border-gray-300 rounded-md"
              >
                {workoutGoals.map(goal => (
                  <option key={goal} value={goal}>
                    {goal.replace('_', ' ').replace(/\b\w/g, l => l.toUpperCase())}
                  </option>
                ))}
              </select>
            </div>
            
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Injuries</label>
              <select
                multiple
                value={userProfile.injuries}
                onChange={(e) => setUserProfile({...userProfile, injuries: Array.from(e.target.selectedOptions, option => option.value)})}
                className="w-full p-2 border border-gray-300 rounded-md"
              >
                {injuryTypes.map(injury => (
                  <option key={injury} value={injury}>
                    {injury.replace('_', ' ').replace(/\b\w/g, l => l.toUpperCase())}
                  </option>
                ))}
              </select>
            </div>
            
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Complaints</label>
              <select
                multiple
                value={userProfile.complaints}
                onChange={(e) => setUserProfile({...userProfile, complaints: Array.from(e.target.selectedOptions, option => option.value)})}
                className="w-full p-2 border border-gray-300 rounded-md"
              >
                {complaintsList.map(complaint => (
                  <option key={complaint} value={complaint}>
                    {complaint.replace('_', ' ').replace(/\b\w/g, l => l.toUpperCase())}
                  </option>
                ))}
              </select>
            </div>
            
            <button
              type="submit"
              disabled={loading}
              className="w-full btn-secondary"
            >
              {loading ? 'Generating...' : 'Generate Workout Plan'}
            </button>
          </form>
        </div>
        
        {/* Workout Plan Results */}
        {showPlan && (
          <div className="space-y-6">
            {/* Workout Plan */}
            <div className="bg-white rounded-xl shadow-lg p-6">
              <h2 className="text-xl font-bold text-blue-600 mb-4">Your Workout Plan</h2>
              <div className="space-y-4">
                {workoutPlan.map((exercise, index) => (
                  <div key={index} className="workout-box">
                    <div className="flex justify-between items-start mb-2">
                      <h3 className="font-semibold text-gray-800">{exercise.name}</h3>
                      <div className="text-right text-sm">
                        <div className="text-gray-600">{exercise.type}</div>
                        <div>{exercise.sets} sets × {exercise.reps} reps</div>
                        <div>Rest: {exercise.rest}s</div>
                      </div>
                    </div>
                    <p className="text-gray-600 mb-2">{exercise.description}</p>
                    
                    <div className="mb-2">
                      <h4 className="font-medium text-gray-700">Common Mistakes to Avoid:</h4>
                      <ul className="list-disc list-inside text-sm text-gray-600 ml-4">
                        {exercise.commonMistakes.map((mistake, i) => (
                          <li key={i}>{mistake}</li>
                        ))}
                      </ul>
                    </div>
                    
                    <div className="border-t pt-2">
                      <h4 className="font-medium text-gray-700">Alternative Exercise</h4>
                      <div className="flex justify-between items-start mb-2">
                        <h5 className="font-medium text-gray-800">{exercise.alternative.name}</h5>
                        <div className="text-right text-sm">
                          <div className="text-gray-600">{exercise.alternative.type}</div>
                          <div>{exercise.alternative.sets} sets × {exercise.alternative.reps} reps</div>
                          <div>Rest: {exercise.alternative.rest}s</div>
                        </div>
                      </div>
                      <p className="text-gray-600">{exercise.alternative.description}</p>
                      
                      <h4 className="font-medium text-gray-700">Common Mistakes to Avoid:</h4>
                      <ul className="list-disc list-inside text-sm text-gray-600 ml-4">
                        {exercise.alternative.commonMistakes.map((mistake, i) => (
                          <li key={i}>{mistake}</li>
                        ))}
                      </ul>
                    </div>
                  </div>
                ))}
              </div>
            </div>
            
            {/* Injury Advice */}
            {injuryAdvice && (
              <div className="bg-white rounded-xl shadow-lg p-6">
                <h2 className="text-xl font-bold text-blue-600 mb-4">Injury Advice: {injuryAdvice.injury}</h2>
                <p className="text-gray-600 mb-4">{injuryAdvice.advice}</p>
                
                <div className="mb-4">
                  <h3 className="font-semibold text-gray-800 mb-2">Recommended Exercises</h3>
                  <ul className="list-disc list-inside text-gray-600 ml-4">
                    {injuryAdvice.exercises.map((exercise, i) => (
                      <li key={i}>{exercise}</li>
                    ))}
                  </ul>
                </div>
                
                <div className="mb-4">
                  <h3 className="font-semibold text-gray-800 mb-2">Treatments</h3>
                  <ul className="list-disc list-inside text-gray-600 ml-4">
                    {injuryAdvice.treatments.map((treatment, i) => (
                      <li key={i}>{treatment}</li>
                    ))}
                  </ul>
                </div>
                
                <div>
                  <h3 className="font-semibold text-gray-800 mb-2">Restrictions</h3>
                  <ul className="list-disc list-inside text-gray-600 ml-4">
                    {injuryAdvice.restrictions.map((restriction, i) => (
                      <li key={i}>{restriction}</li>
                    ))}
                  </ul>
                </div>
              </div>
            )}
            
            {/* Complaint Solution */}
            {complaintSolution && (
              <div className="bg-white rounded-xl shadow-lg p-6">
                <h2 className="text-xl font-bold text-blue-600 mb-4">Solution for: {complaintSolution.complaint}</h2>
                <p className="text-gray-600 mb-4">{complaintSolution.advice}</p>
                
                <div className="mb-4">
                  <h3 className="font-semibold text-gray-800 mb-2">Nutritional Recommendations</h3>
                  <ul className="list-disc list-inside text-gray-600 ml-4">
                    {complaintSolution.nutritionalRecommendations.map((recommendation, i) => (
                      <li key={i}>{recommendation}</li>
                    ))}
                  </ul>
                </div>
                
                <div>
                  <h3 className="font-semibold text-gray-800 mb-2">Recommended Supplements</h3>
                  <div className="overflow-x-auto">
                    <table className="min-w-full divide-y divide-gray-200">
                      <thead className="bg-gray-50">
                        <tr>
                          <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Supplement</th>
                          <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Dosage</th>
                          <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Frequency</th>
                        </tr>
                      </thead>
                      <tbody className="bg-white divide-y divide-gray-200">
                        {complaintSolution.supplements.map((supplement, i) => (
                          <tr key={i}>
                            <td className="px-4 py-2 text-sm text-gray-900">{supplement.name}</td>
                            <td className="px-4 py-2 text-sm text-gray-900">{supplement.dosage}</td>
                            <td className="px-4 py-2 text-sm text-gray-900">{supplement.frequency}</td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                </div>
              </div>
            )}
          </div>
        )}
        
        {/* Medical Disclaimer */}
        <div className="disclaimer mt-8">
          <p>
            <strong>Disclaimer:</strong> This site is to update you with information and guide you to useful advice for the purpose of education and awareness and is not a substitute for a doctor's visit. Always consult with a healthcare professional before starting any new exercise program or supplement regimen, especially if you have pre-existing injuries or medical conditions.
          </p>
        </div>
      </div>
    </div>
  );
}