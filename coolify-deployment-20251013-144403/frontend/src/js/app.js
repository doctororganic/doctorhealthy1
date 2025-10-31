// Main application JavaScript
class NutritionApp {
    constructor() {
        this.apiBaseUrl = 'http://localhost:8080/api/v1';
        this.currentUser = null;
        this.init();
    }

    init() {
        // Initialize app
        this.checkAuthStatus();
        this.setupEventListeners();
        this.loadLanguage();
        
        // Initialize tooltips and popovers
        this.initializeBootstrapComponents();
        
        // Initialize PWA features
        this.initializePWA();
        
        // Setup offline detection
        this.setupOfflineDetection();
    }

    initializeBootstrapComponents() {
        // Initialize Bootstrap tooltips
        const tooltipTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="tooltip"]'));
        tooltipTriggerList.map(function (tooltipTriggerEl) {
            return new bootstrap.Tooltip(tooltipTriggerEl);
        });

        // Initialize Bootstrap popovers
        const popoverTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="popover"]'));
        popoverTriggerList.map(function (popoverTriggerEl) {
            return new bootstrap.Popover(popoverTriggerEl);
        });
    }

    setupEventListeners() {
        // Login form
        const loginForm = document.getElementById('loginForm');
        if (loginForm) {
            loginForm.addEventListener('submit', (e) => this.handleLogin(e));
        }

        // Register form
        const registerForm = document.getElementById('registerForm');
        if (registerForm) {
            registerForm.addEventListener('submit', (e) => this.handleRegister(e));
        }

        // Navigation links
        document.addEventListener('click', (e) => {
            if (e.target.matches('[data-page]')) {
                e.preventDefault();
                this.loadPage(e.target.dataset.page);
            }
        });
    }

    async checkAuthStatus() {
        const token = localStorage.getItem('auth_token');
        if (token) {
            try {
                const response = await this.apiCall('/users/profile', 'GET');
                if (response.ok) {
                    const user = await response.json();
                    this.setCurrentUser(user);
                } else {
                    this.logout();
                }
            } catch (error) {
                console.error('Auth check failed:', error);
                this.logout();
            }
        }
    }

    async handleLogin(e) {
        e.preventDefault();
        
        const email = document.getElementById('loginEmail').value;
        const password = document.getElementById('loginPassword').value;

        this.showLoading(true);

        try {
            const response = await this.apiCall('/auth/login', 'POST', {
                email,
                password
            });

            if (response.ok) {
                const data = await response.json();
                localStorage.setItem('auth_token', data.token);
                localStorage.setItem('refresh_token', data.refresh_token);
                
                this.setCurrentUser(data.user);
                this.hideModal('loginModal');
                this.showSuccess('Login successful!');
            } else {
                const error = await response.json();
                this.showError(error.error || 'Login failed');
            }
        } catch (error) {
            console.error('Login error:', error);
            this.showError('Network error. Please try again.');
        } finally {
            this.showLoading(false);
        }
    }

    async handleRegister(e) {
        e.preventDefault();
        
        const firstName = document.getElementById('firstName').value;
        const lastName = document.getElementById('lastName').value;
        const email = document.getElementById('registerEmail').value;
        const password = document.getElementById('registerPassword').value;

        this.showLoading(true);

        try {
            const response = await this.apiCall('/auth/register', 'POST', {
                first_name: firstName,
                last_name: lastName,
                email,
                password
            });

            if (response.ok) {
                const data = await response.json();
                localStorage.setItem('auth_token', data.token);
                localStorage.setItem('refresh_token', data.refresh_token);
                
                this.setCurrentUser(data.user);
                this.hideModal('registerModal');
                this.showSuccess('Registration successful!');
            } else {
                const error = await response.json();
                this.showError(error.error || 'Registration failed');
            }
        } catch (error) {
            console.error('Registration error:', error);
            this.showError('Network error. Please try again.');
        } finally {
            this.showLoading(false);
        }
    }

    async apiCall(endpoint, method = 'GET', data = null) {
        const url = this.apiBaseUrl + endpoint;
        const options = {
            method,
            headers: {
                'Content-Type': 'application/json',
            },
        };

        const token = localStorage.getItem('auth_token');
        if (token) {
            options.headers['Authorization'] = `Bearer ${token}`;
        }

        if (data) {
            options.body = JSON.stringify(data);
        }

        return fetch(url, options);
    }

    setCurrentUser(user) {
        this.currentUser = user;
        this.updateUI();
    }

    updateUI() {
        const authButtons = document.getElementById('auth-buttons');
        const userMenu = document.getElementById('user-menu');
        const userName = document.getElementById('user-name');

        if (this.currentUser) {
            authButtons.style.display = 'none';
            userMenu.style.display = 'block';
            userName.textContent = `${this.currentUser.first_name} ${this.currentUser.last_name}`;
        } else {
            authButtons.style.display = 'block';
            userMenu.style.display = 'none';
        }
    }

    logout() {
        localStorage.removeItem('auth_token');
        localStorage.removeItem('refresh_token');
        this.currentUser = null;
        this.updateUI();
        this.showSuccess('Logged out successfully');
    }

    showLoading(show) {
        const loading = document.getElementById('loading');
        if (show) {
            loading.classList.add('show');
        } else {
            loading.classList.remove('show');
        }
    }

    showModal(modalId) {
        const modal = new bootstrap.Modal(document.getElementById(modalId));
        modal.show();
    }

    hideModal(modalId) {
        const modal = bootstrap.Modal.getInstance(document.getElementById(modalId));
        if (modal) {
            modal.hide();
        }
    }

    showSuccess(message) {
        this.showToast(message, 'success');
    }

    showError(message) {
        this.showToast(message, 'danger');
    }

    showToast(message, type = 'info') {
        // Create toast element
        const toast = document.createElement('div');
        toast.className = `toast align-items-center text-white bg-${type} border-0`;
        toast.setAttribute('role', 'alert');
        toast.innerHTML = `
            <div class="d-flex">
                <div class="toast-body">
                    ${message}
                </div>
                <button type="button" class="btn-close btn-close-white me-2 m-auto" data-bs-dismiss="toast"></button>
            </div>
        `;

        // Add to toast container or create one
        let toastContainer = document.querySelector('.toast-container');
        if (!toastContainer) {
            toastContainer = document.createElement('div');
            toastContainer.className = 'toast-container position-fixed bottom-0 end-0 p-3';
            document.body.appendChild(toastContainer);
        }

        toastContainer.appendChild(toast);

        // Show toast
        const bsToast = new bootstrap.Toast(toast);
        bsToast.show();

        // Remove toast after it's hidden
        toast.addEventListener('hidden.bs.toast', () => {
            toast.remove();
        });
    }

    loadLanguage() {
        const savedLanguage = localStorage.getItem('language') || 'en';
        this.changeLanguage(savedLanguage);
    }

    // Language Management
    changeLanguage(lang) {
        localStorage.setItem('language', lang);
        
        // Update HTML attributes
        document.documentElement.lang = lang;
        document.documentElement.dir = lang === 'ar' ? 'rtl' : 'ltr';
        
        // Switch Bootstrap CSS for RTL/LTR
        this.switchBootstrapCSS(lang);
        
        // Update body class for RTL support
        document.body.classList.toggle('rtl', lang === 'ar');
        
        // Update flag button states
        this.updateFlagButtons(lang);
        
        // Update all translatable elements
        if (typeof updateTranslations === 'function') {
            updateTranslations();
        }
    }

    // Switch Bootstrap CSS based on language direction
    switchBootstrapCSS(lang) {
        const bootstrapCSS = document.getElementById('bootstrap-css');
        const bootstrapRTLCSS = document.getElementById('bootstrap-rtl-css');
        
        if (lang === 'ar') {
            bootstrapCSS.disabled = true;
            bootstrapRTLCSS.disabled = false;
        } else {
            bootstrapCSS.disabled = false;
            bootstrapRTLCSS.disabled = true;
        }
    }

    // Update flag button active states
    updateFlagButtons(lang) {
        const enFlag = document.getElementById('en-flag');
        const arFlag = document.getElementById('ar-flag');
        
        if (enFlag && arFlag) {
            enFlag.classList.toggle('active', lang === 'en');
            arFlag.classList.toggle('active', lang === 'ar');
        }
    }

    loadPage(page) {
        // Implement page loading logic
        console.log('Loading page:', page);
        // This would typically load different content based on the page
    }

    // PWA Installation
    initializePWA() {
        let deferredPrompt;
        
        // Listen for beforeinstallprompt event
        window.addEventListener('beforeinstallprompt', (e) => {
            console.log('PWA install prompt available');
            e.preventDefault();
            deferredPrompt = e;
            this.showPWAInstallPrompt(deferredPrompt);
        });

        // Listen for app installed event
        window.addEventListener('appinstalled', () => {
            console.log('PWA was installed');
            this.hidePWAInstallPrompt();
            this.showSuccess('App installed successfully!');
        });

        // Register service worker
        if ('serviceWorker' in navigator) {
            window.addEventListener('load', () => {
                navigator.serviceWorker.register('/sw.js')
                    .then((registration) => {
                        console.log('SW registered: ', registration);
                    })
                    .catch((registrationError) => {
                        console.log('SW registration failed: ', registrationError);
                    });
            });
        }
    }

    showPWAInstallPrompt(deferredPrompt) {
        const promptHTML = `
            <div class="pwa-install-prompt" id="pwaInstallPrompt">
                <div class="d-flex justify-content-between align-items-center">
                    <div>
                        <h6 class="mb-1">Install Nutrition Platform</h6>
                        <small>Get the full app experience</small>
                    </div>
                    <div>
                        <button class="btn btn-light btn-sm me-2" onclick="app.installPWA()">Install</button>
                        <button class="btn btn-outline-light btn-sm" onclick="app.hidePWAInstallPrompt()">×</button>
                    </div>
                </div>
            </div>
        `;
        
        document.body.insertAdjacentHTML('beforeend', promptHTML);
        setTimeout(() => {
            document.getElementById('pwaInstallPrompt').classList.add('show');
        }, 2000);
        
        this.deferredPrompt = deferredPrompt;
    }

    async installPWA() {
        if (this.deferredPrompt) {
            this.deferredPrompt.prompt();
            const { outcome } = await this.deferredPrompt.userChoice;
            console.log(`User response to the install prompt: ${outcome}`);
            this.deferredPrompt = null;
            this.hidePWAInstallPrompt();
        }
    }

    hidePWAInstallPrompt() {
        const prompt = document.getElementById('pwaInstallPrompt');
        if (prompt) {
            prompt.remove();
        }
    }

    // Offline Detection
    setupOfflineDetection() {
        const offlineHTML = `
            <div class="offline-indicator" id="offlineIndicator">
                <i class="fas fa-wifi-slash me-2"></i>
                You are currently offline. Some features may not be available.
            </div>
        `;
        
        document.body.insertAdjacentHTML('afterbegin', offlineHTML);
        
        window.addEventListener('online', () => {
            document.getElementById('offlineIndicator').classList.remove('show');
            this.showSuccess('Connection restored!');
        });
        
        window.addEventListener('offline', () => {
            document.getElementById('offlineIndicator').classList.add('show');
            this.showError('You are now offline');
        });
        
        // Check initial connection status
        if (!navigator.onLine) {
            document.getElementById('offlineIndicator').classList.add('show');
        }
    }
}

// Global functions
function showLogin() {
    app.showModal('loginModal');
}

function showRegister() {
    app.showModal('registerModal');
}

function logout() {
    app.logout();
}

// Navigation functions for home page cards
function navigateToSection(section) {
    console.log(`Navigating to ${section} section`);
    const pageMap = {
        'health': 'diet-planning.html',
        'workout': 'workout-generator.html',
        'sports': 'sports-games.html',
        'reviews': 'food-reviews.html',
        'drugs': 'drug-doses.html'
    };
    const page = pageMap[section];
    if (page) {
        window.location.href = page;
    } else {
        alert('Section not implemented yet!');
    }
}

function navigateToHealthDiet() {
    console.log('Navigating to Health & Diet section');
    window.location.href = 'diet-planning.html';
}

function navigateToWorkoutInjuries() {
    console.log('Navigating to Workout & Injuries section');
    window.location.href = 'workout-generator.html';
}

function navigateToWorkout() {
    // Navigate to workout generator page
    window.location.href = 'workout-generator.html';
}

function navigateToSportsGames() {
    console.log('Navigating to Sports & Games section');
    window.location.href = 'sports-games.html';
}

function navigateToFoodReviews() {
    console.log('Navigating to Food Reviews section');
    window.location.href = 'food-reviews.html';
}

function navigateToDrugDoses() {
    console.log('Navigating to Drug Doses section');
    window.location.href = 'drug-doses.html';
}

// Newsletter subscription
function subscribeNewsletter(event) {
    event.preventDefault();
    const email = event.target.querySelector('input[type="email"]').value;
    
    if (email) {
        console.log('Newsletter subscription:', email);
        app.showSuccess('Newsletter subscription successful!');
        event.target.reset();
    }
}

// Language change function
function changeLanguage(lang) {
    if (app) {
        app.changeLanguage(lang);
    }
}

// Initialize app when DOM is loaded
let app;
document.addEventListener('DOMContentLoaded', () => {
    app = new NutritionApp();
});