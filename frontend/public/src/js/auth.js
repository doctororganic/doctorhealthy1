// Authentication management
class AuthManager {
    constructor() {
        this.baseURL = 'http://localhost:8080/api/v1';
        this.token = localStorage.getItem('access_token');
        this.refreshToken = localStorage.getItem('refresh_token');
        this.user = JSON.parse(localStorage.getItem('user') || 'null');
    }

    // Check if user is authenticated
    isAuthenticated() {
        return !!this.token && !!this.user;
    }

    // Check if user is admin
    isAdmin() {
        return this.user && this.user.role === 'admin';
    }

    // Get current user
    getCurrentUser() {
        return this.user;
    }

    // Set authentication data
    setAuth(accessToken, refreshToken, user) {
        this.token = accessToken;
        this.refreshToken = refreshToken;
        this.user = user;
        
        localStorage.setItem('access_token', accessToken);
        localStorage.setItem('refresh_token', refreshToken);
        localStorage.setItem('user', JSON.stringify(user));
    }

    // Clear authentication data
    clearAuth() {
        this.token = null;
        this.refreshToken = null;
        this.user = null;
        
        localStorage.removeItem('access_token');
        localStorage.removeItem('refresh_token');
        localStorage.removeItem('user');
    }

    // Make authenticated API request
    async apiRequest(endpoint, options = {}) {
        const url = `${this.baseURL}${endpoint}`;
        const headers = {
            'Content-Type': 'application/json',
            ...options.headers
        };

        if (this.token) {
            headers['Authorization'] = `Bearer ${this.token}`;
        }

        try {
            const response = await fetch(url, {
                ...options,
                headers
            });

            // Handle token expiration
            if (response.status === 401 && this.refreshToken) {
                const refreshed = await this.refreshAccessToken();
                if (refreshed) {
                    // Retry the original request with new token
                    headers['Authorization'] = `Bearer ${this.token}`;
                    return await fetch(url, {
                        ...options,
                        headers
                    });
                }
            }

            return response;
        } catch (error) {
            console.error('API request failed:', error);
            throw error;
        }
    }

    // Refresh access token
    async refreshAccessToken() {
        try {
            const response = await fetch(`${this.baseURL}/auth/refresh`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    refresh_token: this.refreshToken
                })
            });

            if (response.ok) {
                const data = await response.json();
                this.token = data.access_token;
                localStorage.setItem('access_token', data.access_token);
                return true;
            } else {
                // Refresh token is invalid, logout user
                this.logout();
                return false;
            }
        } catch (error) {
            console.error('Token refresh failed:', error);
            this.logout();
            return false;
        }
    }

    // Register new user
    async register(userData) {
        try {
            const response = await fetch(`${this.baseURL}/auth/register`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(userData)
            });

            const data = await response.json();

            if (response.ok) {
                this.setAuth(data.access_token, data.refresh_token, data.user);
                return { success: true, data };
            } else {
                return { success: false, error: data.message || 'Registration failed' };
            }
        } catch (error) {
            console.error('Registration error:', error);
            return { success: false, error: 'Network error occurred' };
        }
    }

    // Login user
    async login(email, password) {
        try {
            const response = await fetch(`${this.baseURL}/auth/login`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ email, password })
            });

            const data = await response.json();

            if (response.ok) {
                this.setAuth(data.access_token, data.refresh_token, data.user);
                return { success: true, data };
            } else {
                return { success: false, error: data.message || 'Login failed' };
            }
        } catch (error) {
            console.error('Login error:', error);
            return { success: false, error: 'Network error occurred' };
        }
    }

    // Logout user
    async logout() {
        try {
            if (this.token) {
                await this.apiRequest('/auth/logout', {
                    method: 'POST'
                });
            }
        } catch (error) {
            console.error('Logout error:', error);
        } finally {
            this.clearAuth();
            window.location.reload();
        }
    }

    // Get user profile
    async getProfile() {
        try {
            const response = await this.apiRequest('/auth/profile');
            if (response.ok) {
                const data = await response.json();
                this.user = data.user;
                localStorage.setItem('user', JSON.stringify(data.user));
                return { success: true, data: data.user };
            } else {
                const error = await response.json();
                return { success: false, error: error.error || 'Failed to get profile' };
            }
        } catch (error) {
            console.error('Get profile error:', error);
            return { success: false, error: 'Network error occurred' };
        }
    }

    // Update user profile
    async updateProfile(profileData) {
        try {
            const response = await this.apiRequest('/auth/profile', {
                method: 'PUT',
                body: JSON.stringify(profileData)
            });

            if (response.ok) {
                const data = await response.json();
                this.user = data.user;
                localStorage.setItem('user', JSON.stringify(data.user));
                return { success: true, data: data.user };
            } else {
                const error = await response.json();
                return { success: false, error: error.error || 'Failed to update profile' };
            }
        } catch (error) {
            console.error('Update profile error:', error);
            return { success: false, error: 'Network error occurred' };
        }
    }

    // Get user sessions
    async getSessions() {
        try {
            const response = await this.apiRequest('/auth/sessions');
            if (response.ok) {
                const data = await response.json();
                return { success: true, data: data.sessions };
            } else {
                const error = await response.json();
                return { success: false, error: error.error || 'Failed to get sessions' };
            }
        } catch (error) {
            console.error('Get sessions error:', error);
            return { success: false, error: 'Network error occurred' };
        }
    }

    // Logout from all devices
    async logoutAll() {
        try {
            const response = await this.apiRequest('/auth/logout-all', {
                method: 'POST'
            });
            if (response.ok) {
                this.clearAuth();
                window.location.reload();
                return { success: true };
            } else {
                const error = await response.json();
                return { success: false, error: error.error || 'Failed to logout from all devices' };
            }
        } catch (error) {
            console.error('Logout all error:', error);
            return { success: false, error: 'Network error occurred' };
        }
    }

    // Invalidate specific session
    async invalidateSession(sessionId) {
        try {
            const response = await this.apiRequest(`/auth/sessions/${sessionId}`, {
                method: 'DELETE'
            });
            if (response.ok) {
                return { success: true };
            } else {
                const error = await response.json();
                return { success: false, error: error.error || 'Failed to invalidate session' };
            }
        } catch (error) {
            console.error('Invalidate session error:', error);
            return { success: false, error: 'Network error occurred' };
        }
    }
}

// Create global auth manager instance
const authManager = new AuthManager();

// Form handlers
function handleLoginForm(event) {
    event.preventDefault();
    const form = event.target;
    const formData = new FormData(form);
    const email = formData.get('email');
    const password = formData.get('password');
    
    // Validate form using ValidationSystem
    const validationRules = [
        { field: 'email', type: 'email', required: true, name: 'Email' },
        { field: 'password', type: 'password', required: true, name: 'Password' }
    ];
    
    if (!validation.validateForm('loginForm', validationRules)) {
        return;
    }
    
    showLoading(true);
    
    authManager.login(email, password)
        .then(result => {
            if (result.success) {
                showToast('Login successful!', 'success');
                closeModal('loginModal');
                updateAuthUI();
                form.reset();
                validation.clearErrors();
            } else {
                showToast(result.error, 'error');
            }
        })
        .finally(() => {
            showLoading(false);
        });
}

function handleRegisterForm(event) {
    event.preventDefault();
    const form = event.target;
    const formData = new FormData(form);
    
    const password = formData.get('password');
    const confirmPassword = formData.get('confirmPassword');
    
    // Validate form using ValidationSystem
    const validationRules = [
        { field: 'email', type: 'email', required: true, name: 'Email' },
        { field: 'password', type: 'password', required: true, name: 'Password' },
        { field: 'firstName', type: 'name', required: true, name: 'First Name' },
        { field: 'lastName', type: 'name', required: true, name: 'Last Name' },
        { field: 'dateOfBirth', type: 'date', required: true, name: 'Date of Birth' }
    ];
    
    if (!validation.validateForm('registerForm', validationRules)) {
        return;
    }
    
    // Additional password confirmation validation
    if (!validation.validateConfirmPassword(password, confirmPassword, 'confirmPassword')) {
        validation.displayErrors();
        return;
    }
    
    const userData = {
        email: formData.get('email'),
        password: password,
        confirm_password: confirmPassword,
        first_name: formData.get('firstName'),
        last_name: formData.get('lastName'),
        date_of_birth: formData.get('dateOfBirth'),
        gender: formData.get('gender'),
        language: formData.get('language') || 'en'
    };
    
    showLoading(true);
    
    authManager.register(userData)
        .then(result => {
            if (result.success) {
                showToast('Registration successful!', 'success');
                closeModal('registerModal');
                updateAuthUI();
                form.reset();
                validation.clearErrors();
            } else {
                showToast(result.error, 'error');
            }
        })
        .finally(() => {
            showLoading(false);
        });
}

function handleLogout() {
    if (confirm(translate('confirm_logout') || 'Are you sure you want to logout?')) {
        authManager.logout();
    }
}

// UI update functions
function updateAuthUI() {
    const loginBtn = document.getElementById('loginBtn');
    const registerBtn = document.getElementById('registerBtn');
    const userDropdown = document.getElementById('userDropdown');
    const userName = document.getElementById('userName');
    
    if (authManager.isAuthenticated()) {
        // Hide login/register buttons
        if (loginBtn) loginBtn.style.display = 'none';
        if (registerBtn) registerBtn.style.display = 'none';
        
        // Show user dropdown
        if (userDropdown) {
            userDropdown.style.display = 'block';
            if (userName) {
                const user = authManager.getCurrentUser();
                userName.textContent = `${user.first_name} ${user.last_name}`;
            }
        }
    } else {
        // Show login/register buttons
        if (loginBtn) loginBtn.style.display = 'inline-block';
        if (registerBtn) registerBtn.style.display = 'inline-block';
        
        // Hide user dropdown
        if (userDropdown) userDropdown.style.display = 'none';
    }
}

// Initialize auth UI on page load
document.addEventListener('DOMContentLoaded', () => {
    updateAuthUI();
    
    // Add form event listeners
    const loginForm = document.getElementById('loginForm');
    const registerForm = document.getElementById('registerForm');
    
    if (loginForm) {
        loginForm.addEventListener('submit', handleLoginForm);
    }
    
    if (registerForm) {
        registerForm.addEventListener('submit', handleRegisterForm);
    }
});

// Export for global use
window.authManager = authManager;
window.handleLogout = handleLogout;
window.updateAuthUI = updateAuthUI;