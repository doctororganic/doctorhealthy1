'use client'

import Link from 'next/link'
import { useState } from 'react'

export default function Home() {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const [message, setMessage] = useState('')

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault()
    setIsLoading(true)
    setMessage('')

    try {
      const response = await fetch('http://localhost:8080/api/v1/auth/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email, password }),
      })

      const data = await response.json()
      
      if (response.ok) {
        setMessage('Login successful!')
        localStorage.setItem('token', data.token)
        window.location.href = '/dashboard'
      } else {
        setMessage(data.error || 'Login failed')
      }
    } catch (error) {
      setMessage('Network error - please try again')
    } finally {
      setIsLoading(false)
    }
  }

  const testApi = async () => {
    try {
      const response = await fetch('http://localhost:8080/health')
      const data = await response.json()
      setMessage(`API Status: ${data.status}`)
    } catch (error) {
      setMessage('API connection failed')
    }
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-green-50 to-blue-50">
      <div className="container mx-auto px-4 py-8">
        <header className="text-center mb-8">
          <h1 className="text-4xl font-bold text-gray-800 mb-2">
            Nutrition Platform
          </h1>
          <p className="text-gray-600">
            Your comprehensive health and fitness companion
          </p>
        </header>

        <div className="max-w-md mx-auto bg-white rounded-lg shadow-md p-6">
          <h2 className="text-2xl font-semibold text-center mb-6">
            Sign In
          </h2>
          
          <form onSubmit={handleLogin} className="space-y-4">
            <div>
              <label htmlFor="email" className="block text-sm font-medium text-gray-700 mb-1">
                Email
              </label>
              <input
                type="email"
                id="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-green-500"
                placeholder="Enter your email"
                required
              />
            </div>

            <div>
              <label htmlFor="password" className="block text-sm font-medium text-gray-700 mb-1">
                Password
              </label>
              <input
                type="password"
                id="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-green-500"
                placeholder="Enter your password"
                required
              />
            </div>

            <button
              type="submit"
              disabled={isLoading}
              className="w-full bg-green-600 text-white py-2 px-4 rounded-md hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-green-500 disabled:opacity-50"
            >
              {isLoading ? 'Signing in...' : 'Sign In'}
            </button>
          </form>

          {message && (
            <div className="mt-4 p-3 rounded-md bg-blue-50 text-blue-800 text-sm">
              {message}
            </div>
          )}

          <div className="mt-6 text-center">
            <button
              onClick={testApi}
              className="text-sm text-blue-600 hover:text-blue-800"
            >
              Test API Connection
            </button>
          </div>

          <div className="mt-6 text-center text-sm text-gray-600">
            Don't have an account?{' '}
            <Link href="/register" className="text-green-600 hover:text-green-800">
              Sign up
            </Link>
          </div>
        </div>

        <div className="mt-8 max-w-4xl mx-auto">
          <div className="grid md:grid-cols-3 gap-6">
            <div className="bg-white p-6 rounded-lg shadow-md">
              <h3 className="text-xl font-semibold text-gray-800 mb-2">
                🥗 Nutrition Tracking
              </h3>
              <p className="text-gray-600">
                Track your meals, monitor calories, and get personalized nutrition plans.
              </p>
            </div>

            <div className="bg-white p-6 rounded-lg shadow-md">
              <h3 className="text-xl font-semibold text-gray-800 mb-2">
                💪 Fitness Planning
              </h3>
              <p className="text-gray-600">
                Create workout routines, track progress, and achieve your fitness goals.
              </p>
            </div>

            <div className="bg-white p-6 rounded-lg shadow-md">
              <h3 className="text-xl font-semibold text-gray-800 mb-2">
                📊 Health Analytics
              </h3>
              <p className="text-gray-600">
                Monitor your health metrics and get insights into your wellness journey.
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
