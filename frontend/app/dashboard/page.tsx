'use client'

import { useState, useEffect } from 'react'

export default function Dashboard() {
  const [user, setUser] = useState<any>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [message, setMessage] = useState('')

  useEffect(() => {
    const token = localStorage.getItem('token')
    if (!token) {
      window.location.href = '/'
      return
    }

    const fetchUserData = async () => {
      try {
        const response = await fetch('http://localhost:8080/api/v1/auth/me', {
          headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json',
          },
        })

        if (response.ok) {
          const userData = await response.json()
          setUser(userData.user)
        } else {
          setMessage('Failed to load user data')
          localStorage.removeItem('token')
          setTimeout(() => {
            window.location.href = '/'
          }, 2000)
        }
      } catch (error) {
        setMessage('Network error')
        localStorage.removeItem('token')
        setTimeout(() => {
          window.location.href = '/'
        }, 2000)
      } finally {
        setIsLoading(false)
      }
    }

    fetchUserData()
  }, [])

  const handleLogout = () => {
    localStorage.removeItem('token')
    window.location.href = '/'
  }

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-green-600 mx-auto"></div>
          <p className="mt-4 text-gray-600">Loading...</p>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <nav className="bg-white shadow-sm border-b">
        <div className="container mx-auto px-4 py-4">
          <div className="flex justify-between items-center">
            <h1 className="text-2xl font-bold text-gray-800">Nutrition Platform</h1>
            <div className="flex items-center space-x-4">
              <span className="text-gray-600">Welcome, {user?.name}</span>
              <button
                onClick={handleLogout}
                className="bg-red-600 text-white px-4 py-2 rounded-md hover:bg-red-700"
              >
                Logout
              </button>
            </div>
          </div>
        </div>
      </nav>

      <div className="container mx-auto px-4 py-8">
        {message && (
          <div className="mb-4 p-3 rounded-md bg-red-50 text-red-800 text-sm">
            {message}
          </div>
        )}

        <div className="mb-8">
          <h2 className="text-3xl font-bold text-gray-800 mb-2">Dashboard</h2>
          <p className="text-gray-600">Welcome back! Here's your nutrition journey overview.</p>
        </div>

        <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
          <div className="bg-white p-6 rounded-lg shadow-sm border">
            <h3 className="text-lg font-semibold text-gray-800 mb-2">Daily Calories</h3>
            <p className="text-3xl font-bold text-green-600">2,150</p>
            <p className="text-sm text-gray-500">of 2,500 goal</p>
          </div>

          <div className="bg-white p-6 rounded-lg shadow-sm border">
            <h3 className="text-lg font-semibold text-gray-800 mb-2">Water Intake</h3>
            <p className="text-3xl font-bold text-blue-600">6</p>
            <p className="text-sm text-gray-500">of 8 glasses</p>
          </div>

          <div className="bg-white p-6 rounded-lg shadow-sm border">
            <h3 className="text-lg font-semibold text-gray-800 mb-2">Workouts</h3>
            <p className="text-3xl font-bold text-purple-600">3</p>
            <p className="text-sm text-gray-500">this week</p>
          </div>

          <div className="bg-white p-6 rounded-lg shadow-sm border">
            <h3 className="text-lg font-semibold text-gray-800 mb-2">Weight</h3>
            <p className="text-3xl font-bold text-orange-600">165</p>
            <p className="text-sm text-gray-500">lbs (-2 lbs)</p>
          </div>
        </div>

        <div className="grid md:grid-cols-2 gap-8">
          <div className="bg-white p-6 rounded-lg shadow-sm border">
            <h3 className="text-xl font-semibold text-gray-800 mb-4">Recent Meals</h3>
            <div className="space-y-3">
              <div className="flex justify-between items-center p-3 bg-gray-50 rounded">
                <div>
                  <p className="font-medium">Breakfast</p>
                  <p className="text-sm text-gray-500">8:00 AM</p>
                </div>
                <p className="text-green-600 font-medium">450 cal</p>
              </div>
              <div className="flex justify-between items-center p-3 bg-gray-50 rounded">
                <div>
                  <p className="font-medium">Lunch</p>
                  <p className="text-sm text-gray-500">12:30 PM</p>
                </div>
                <p className="text-green-600 font-medium">680 cal</p>
              </div>
              <div className="flex justify-between items-center p-3 bg-gray-50 rounded">
                <div>
                  <p className="font-medium">Snack</p>
                  <p className="text-sm text-gray-500">3:00 PM</p>
                </div>
                <p className="text-green-600 font-medium">150 cal</p>
              </div>
            </div>
          </div>

          <div className="bg-white p-6 rounded-lg shadow-sm border">
            <h3 className="text-xl font-semibold text-gray-800 mb-4">Quick Actions</h3>
            <div className="grid grid-cols-2 gap-4">
              <button className="p-4 bg-green-50 text-green-700 rounded-lg hover:bg-green-100 transition">
                <div className="text-2xl mb-2">🍎</div>
                <p className="font-medium">Log Meal</p>
              </button>
              <button className="p-4 bg-blue-50 text-blue-700 rounded-lg hover:bg-blue-100 transition">
                <div className="text-2xl mb-2">💪</div>
                <p className="font-medium">Log Workout</p>
              </button>
              <button className="p-4 bg-purple-50 text-purple-700 rounded-lg hover:bg-purple-100 transition">
                <div className="text-2xl mb-2">📊</div>
                <p className="font-medium">View Progress</p>
              </button>
              <button className="p-4 bg-orange-50 text-orange-700 rounded-lg hover:bg-orange-100 transition">
                <div className="text-2xl mb-2">🎯</div>
                <p className="font-medium">Set Goals</p>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
