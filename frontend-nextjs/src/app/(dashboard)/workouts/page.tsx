'use client';

import { useState } from 'react';
import { useWorkouts } from '@/hooks/useNutritionData';
import { LoadingSkeleton } from '@/components/ui/LoadingSkeleton';
import { ErrorDisplay } from '@/components/ui/ErrorDisplay';
import { EmptyState } from '@/components/ui/EmptyState';

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
  const [currentPage, setCurrentPage] = useState(1);
  const { data: workoutData, loading: workoutsLoading, error: workoutsError, refetch: refetchWorkouts, pagination: workoutsPagination } = useWorkouts({ 
    goal: userProfile.workoutGoal,
    level: userProfile.activityLevel,
    page: currentPage, 
    limit: 20 
  });
  const [injuryAdvice, setInjuryAdvice] = useState<InjuryAdvice | null>(null);
  const [complaintSolution, setComplaintSolution] = useState<ComplaintSolution | null>(null);
  const [showPlan, setShowPlan] = useState(false);
  const [loading, setLoading] = useState(false);

  // Workout types
  // Mock data removed - workoutTypes now come from API

  // Complaint options
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
    'weight_loss_difficulty'
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
    setShowPlan(true);
    refetchWorkouts(); // Load real data from API
  };

  // Use API data from hook
  const workouts = (workoutData?.items || []) as any[];

  // Get injury advice
  const getInjuryAdvice = () => {
    const { injuries } = userProfile;
    
    if (injuries.length === 0) return;
    
    // TODO: Load injury advice from API using injuries hook
    // For now, set a placeholder since mock data was removed
    const firstInjury = injuries[0];
    setInjuryAdvice({
      injury: firstInjury,
      advice: `Please consult with a healthcare professional for advice regarding ${firstInjury}.`,
      exercises: [],
      treatments: [],
      restrictions: []
    });
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
          {/* Live fetched workouts preview */}
          <div className="mb-4 p-3 rounded border bg-blue-50">
            <div className="font-medium text-blue-800">Fetched Workouts Preview</div>
            {workoutsLoading && <div className="text-blue-700">Loading workouts...</div>}
            {workoutsError && <div className="text-red-700">{String(workoutsError)}</div>}
            {workoutData && (
              <div className="text-sm text-blue-900">{workoutData.items?.length || 0} workouts available</div>
            )}
            <button type="button" onClick={() => refetchWorkouts()} className="mt-2 btn-secondary">Refresh Workouts</button>
          </div>
          
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
              {workoutsLoading && <div className="text-center py-8">Loading workouts...</div>}
              {workoutsError && <div className="text-red-600 py-4">Error loading workouts: {workoutsError}</div>}
              {!workoutsLoading && !workoutsError && workouts.length === 0 && (
                <div className="text-center py-8 text-gray-500">No workouts found. Try adjusting your filters.</div>
              )}
              {!workoutsLoading && !workoutsError && workouts.length > 0 && (
                <div className="space-y-4">
                  {workouts.map((workout: any, index: number) => (
                    <div key={workout.id || index} className="workout-box">
                      <h3 className="font-semibold text-gray-800">{workout.name || workout.title || `Workout ${index + 1}`}</h3>
                      {workout.description && <p className="text-gray-600 mt-2">{workout.description}</p>}
                      {workout.exercises && (
                        <div className="mt-4">
                          <h4 className="font-medium text-gray-700 mb-2">Exercises:</h4>
                          <ul className="list-disc list-inside text-sm text-gray-600 ml-4">
                            {Array.isArray(workout.exercises) ? workout.exercises.map((ex: any, i: number) => (
                              <li key={i}>{ex.name || ex}</li>
                            )) : <li>{workout.exercises}</li>}
                          </ul>
                        </div>
                      )}
                    </div>
                  ))}
                </div>
              )}
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

      {/* Pagination Controls */}
      {workoutsPagination && (
        <div className="flex justify-between items-center mt-6 px-4">
          <button 
            onClick={() => setCurrentPage(prev => Math.max(prev - 1, 1))}
            disabled={currentPage <= 1 || workoutsLoading}
            className="px-4 py-2 bg-gray-300 text-gray-700 rounded-md disabled:opacity-50 hover:bg-gray-400"
          >
            ← Previous
          </button>
          <span className="text-sm text-gray-600">
              Page {workoutsPagination.page} of {workoutsPagination.totalPages}
            {workoutsPagination.total && ` (${workoutsPagination.total} total workouts)`}
          </span>
          <button 
            onClick={() => setCurrentPage(prev => prev + 1)}
            disabled={currentPage >= workoutsPagination.totalPages || workoutsLoading}
            className="px-4 py-2 bg-blue-600 text-white rounded-md disabled:opacity-50 hover:bg-blue-700"
          >
            Next →
          </button>
        </div>
      )}
    </div>
  );
}
