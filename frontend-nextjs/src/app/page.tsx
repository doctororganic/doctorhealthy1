import Link from 'next/link';
import { MealsIcon } from '@/components/icons/MealsIcon';
import { WorkoutIcon } from '@/components/icons/WorkoutIcon';
import { RecipeIcon } from '@/components/icons/RecipeIcon';
import { DiseaseIcon } from '@/components/icons/DiseaseIcon';

export default function HomePage() {
  return (
    <div className="min-h-screen bg-gradient-to-b from-white to-yellow-50">
      {/* Header */}
      <header className="bg-white shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center py-6">
            <div className="flex items-center">
              <div className="text-2xl font-bold text-green-600">
                Dr. Pass Nutrition Platform
              </div>
            </div>
            <nav className="hidden md:flex space-x-10">
              <a href="#" className="text-gray-700 hover:text-green-600 transition-colors">
                Home
              </a>
              <a href="#" className="text-gray-700 hover:text-green-600 transition-colors">
                About
              </a>
              <a href="#" className="text-gray-700 hover:text-green-600 transition-colors">
                Contact
              </a>
            </nav>
          </div>
        </div>
      </header>

      {/* Hero Section */}
      <section className="py-12 px-4 sm:px-6 lg:px-8">
        <div className="max-w-7xl mx-auto text-center">
          <h1 className="text-4xl font-extrabold text-gray-900 sm:text-5xl md:text-6xl">
            Your Personalized
            <span className="block text-green-600"> Nutrition Journey</span>
          </h1>
          <p className="mt-3 max-w-md mx-auto text-base text-gray-500 sm:text-lg md:mt-5 md:text-xl md:max-w-3xl">
            Get customized meal plans, workout routines, recipes, and health advice tailored to your specific needs.
          </p>
        </div>
      </section>

      {/* Main Content - 4 Boxes */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
          
          {/* Box 1: Meals and Body Enhancing */}
          <div className="feature-card p-8">
            <div className="flex items-center mb-4">
              <div className="p-3 bg-green-100 rounded-full mr-4">
                <MealsIcon className="w-8 h-8 text-green-600" />
              </div>
              <h2 className="feature-title">Meals and Body Enhancing</h2>
            </div>
            <p className="text-gray-600 mb-6">
              Get personalized meal plans based on your body metrics, goals, and dietary preferences. Calculate exact calories, macros, and meal timing.
            </p>
            <Link href="/meals" className="btn-primary inline-block">
              Get Meal Plan
            </Link>
          </div>

          {/* Box 2: Workouts and Injuries */}
          <div className="feature-card p-8">
            <div className="flex items-center mb-4">
              <div className="p-3 bg-blue-100 rounded-full mr-4">
                <WorkoutIcon className="w-8 h-8 text-blue-600" />
              </div>
              <h2 className="feature-title">Workouts and Injuries</h2>
            </div>
            <p className="text-gray-600 mb-6">
              Customized workout routines that consider your fitness level, goals, and any injuries or physical limitations.
            </p>
            <Link href="/workouts" className="btn-secondary inline-block">
              Start Workout
            </Link>
          </div>

          {/* Box 3: Recipes and Review */}
          <div className="feature-card p-8">
            <div className="flex items-center mb-4">
              <div className="p-3 bg-green-100 rounded-full mr-4">
                <RecipeIcon className="w-8 h-8 text-green-600" />
              </div>
              <h2 className="feature-title">Recipes and Review</h2>
            </div>
            <p className="text-gray-600 mb-6">
              Explore recipes from different cuisines, all with halal options and nutritional information. Review and save your favorites.
            </p>
            <Link href="/recipes" className="btn-primary inline-block">
              Browse Recipes
            </Link>
          </div>

          {/* Box 4: Diseases and Healthy-Lifestyle */}
          <div className="feature-card p-8">
            <div className="flex items-center mb-4">
              <div className="p-3 bg-blue-100 rounded-full mr-4">
                <DiseaseIcon className="w-8 h-8 text-blue-600" />
              </div>
              <h2 className="feature-title">Diseases and Healthy-Lifestyle</h2>
            </div>
            <p className="text-gray-600 mb-6">
              Get nutritional advice tailored to your health conditions and medications, with proper disclaimers and professional guidance.
            </p>
            <Link href="/health" className="btn-secondary inline-block">
              Health Advice
            </Link>
          </div>
        </div>
      </main>

      {/* Footer */}
      <footer className="bg-white border-t border-gray-200">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          <div className="md:flex md:justify-between">
            <div className="mb-6 md:mb-0">
              <div className="flex items-center">
                <div className="text-xl font-bold text-green-600">
                  Dr. Pass Nutrition Platform
                </div>
              </div>
              <p className="mt-2 text-sm text-gray-600">
                Your health, our priority.
              </p>
            </div>
            <div className="grid grid-cols-2 gap-8 md:gap-20">
              <div>
                <h3 className="text-sm font-semibold text-gray-900 uppercase tracking-wider">
                  Features
                </h3>
                <ul className="mt-4 space-y-2 text-sm text-gray-600">
                  <li><a href="#" className="hover:text-green-600">Meal Plans</a></li>
                  <li><a href="#" className="hover:text-green-600">Workouts</a></li>
                  <li><a href="#" className="hover:text-green-600">Recipes</a></li>
                  <li><a href="#" className="hover:text-green-600">Health Advice</a></li>
                </ul>
              </div>
              <div>
                <h3 className="text-sm font-semibold text-gray-900 uppercase tracking-wider">
                  Support
                </h3>
                <ul className="mt-4 space-y-2 text-sm text-gray-600">
                  <li><a href="#" className="hover:text-green-600">Help Center</a></li>
                  <li><a href="#" className="hover:text-green-600">Contact Us</a></li>
                  <li><a href="#" className="hover:text-green-600">FAQ</a></li>
                </ul>
              </div>
            </div>
          </div>
          <div className="mt-8 border-t border-gray-200 pt-6">
            <p className="text-center text-sm text-gray-600">
              &copy; {new Date().getFullYear()} Dr. Pass Nutrition Platform. All rights reserved.
            </p>
          </div>
        </div>
      </footer>
    </div>
  );
}