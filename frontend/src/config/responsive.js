// Responsive configuration for nutrition platform
// Supports RTL/LTR, accessibility, offline UX, and PWA features

// Breakpoints for responsive design
export const breakpoints = {
  xs: '320px',
  sm: '576px',
  md: '768px',
  lg: '992px',
  xl: '1200px',
  xxl: '1400px'
};

// Media queries
export const mediaQueries = {
  mobile: `(max-width: ${breakpoints.md})`,
  tablet: `(min-width: ${breakpoints.md}) and (max-width: ${breakpoints.lg})`,
  desktop: `(min-width: ${breakpoints.lg})`,
  largeDesktop: `(min-width: ${breakpoints.xl})`,
  touch: '(hover: none) and (pointer: coarse)',
  hover: '(hover: hover) and (pointer: fine)',
  reducedMotion: '(prefers-reduced-motion: reduce)',
  darkMode: '(prefers-color-scheme: dark)',
  lightMode: '(prefers-color-scheme: light)',
  highContrast: '(prefers-contrast: high)'
};

// RTL/LTR configuration
export const directionConfig = {
  // Supported languages with their text direction
  languages: {
    en: { dir: 'ltr', name: 'English', flag: '🇺🇸' },
    ar: { dir: 'rtl', name: 'العربية', flag: '🇸🇦' },
    he: { dir: 'rtl', name: 'עברית', flag: '🇮🇱' },
    fa: { dir: 'rtl', name: 'فارسی', flag: '🇮🇷' },
    ur: { dir: 'rtl', name: 'اردو', flag: '🇵🇰' },
    es: { dir: 'ltr', name: 'Español', flag: '🇪🇸' },
    fr: { dir: 'ltr', name: 'Français', flag: '🇫🇷' },
    de: { dir: 'ltr', name: 'Deutsch', flag: '🇩🇪' },
    it: { dir: 'ltr', name: 'Italiano', flag: '🇮🇹' },
    pt: { dir: 'ltr', name: 'Português', flag: '🇵🇹' },
    ru: { dir: 'ltr', name: 'Русский', flag: '🇷🇺' },
    zh: { dir: 'ltr', name: '中文', flag: '🇨🇳' },
    ja: { dir: 'ltr', name: '日本語', flag: '🇯🇵' },
    ko: { dir: 'ltr', name: '한국어', flag: '🇰🇷' }
  },
  
  // Default language
  defaultLanguage: 'en',
  
  // RTL-specific CSS properties mapping
  rtlProperties: {
    'margin-left': 'margin-right',
    'margin-right': 'margin-left',
    'padding-left': 'padding-right',
    'padding-right': 'padding-left',
    'border-left': 'border-right',
    'border-right': 'border-left',
    'left': 'right',
    'right': 'left',
    'text-align': {
      'left': 'right',
      'right': 'left'
    },
    'float': {
      'left': 'right',
      'right': 'left'
    },
    'transform': {
      'translateX': (value) => `translateX(${-parseFloat(value)}${value.replace(/[0-9.-]/g, '')})`
    }
  }
};

// Accessibility configuration
export const accessibilityConfig = {
  // ARIA labels and descriptions
  ariaLabels: {
    navigation: 'Main navigation',
    search: 'Search nutrition information',
    menu: 'Menu',
    close: 'Close',
    loading: 'Loading content',
    error: 'Error message',
    success: 'Success message',
    warning: 'Warning message',
    info: 'Information message',
    skipToContent: 'Skip to main content',
    languageSelector: 'Select language',
    themeToggle: 'Toggle dark/light theme',
    userMenu: 'User account menu',
    notifications: 'Notifications',
    settings: 'Settings'
  },
  
  // Focus management
  focusConfig: {
    trapFocus: true,
    restoreFocus: true,
    initialFocus: '[data-autofocus]',
    focusableSelectors: [
      'a[href]',
      'button:not([disabled])',
      'input:not([disabled])',
      'select:not([disabled])',
      'textarea:not([disabled])',
      '[tabindex]:not([tabindex="-1"])',
      '[contenteditable="true"]'
    ].join(', ')
  },
  
  // Color contrast ratios (WCAG AA compliance)
  contrastRatios: {
    normal: 4.5,
    large: 3,
    enhanced: 7
  },
  
  // Font size scaling
  fontScaling: {
    min: 0.8,
    max: 2.0,
    step: 0.1,
    default: 1.0
  },
  
  // Animation preferences
  animations: {
    respectReducedMotion: true,
    defaultDuration: 300,
    reducedDuration: 0,
    easing: 'cubic-bezier(0.4, 0, 0.2, 1)'
  }
};

// PWA configuration
export const pwaConfig = {
  // Service Worker settings
  serviceWorker: {
    scope: '/',
    updateViaCache: 'none',
    skipWaiting: true,
    clientsClaim: true
  },
  
  // Cache strategies
  cacheStrategies: {
    // Static assets (CSS, JS, images)
    static: {
      strategy: 'CacheFirst',
      cacheName: 'static-cache',
      maxEntries: 100,
      maxAgeSeconds: 30 * 24 * 60 * 60 // 30 days
    },
    
    // API responses
    api: {
      strategy: 'NetworkFirst',
      cacheName: 'api-cache',
      maxEntries: 50,
      maxAgeSeconds: 5 * 60, // 5 minutes
      networkTimeoutSeconds: 3
    },
    
    // Images
    images: {
      strategy: 'CacheFirst',
      cacheName: 'image-cache',
      maxEntries: 200,
      maxAgeSeconds: 7 * 24 * 60 * 60 // 7 days
    },
    
    // Fonts
    fonts: {
      strategy: 'CacheFirst',
      cacheName: 'font-cache',
      maxEntries: 30,
      maxAgeSeconds: 365 * 24 * 60 * 60 // 1 year
    }
  },
  
  // Offline fallbacks
  offlineFallbacks: {
    document: '/offline.html',
    image: '/images/offline-placeholder.svg',
    audio: '/audio/offline-notification.mp3'
  },
  
  // Background sync
  backgroundSync: {
    enabled: true,
    tagPrefix: 'nutrition-sync-',
    maxRetryTime: 24 * 60 * 60 * 1000 // 24 hours
  },
  
  // Push notifications
  pushNotifications: {
    enabled: true,
    vapidPublicKey: process.env.REACT_APP_VAPID_PUBLIC_KEY,
    applicationServerKey: process.env.REACT_APP_APPLICATION_SERVER_KEY
  }
};

// Offline UX configuration
export const offlineConfig = {
  // Detection settings
  detection: {
    checkInterval: 5000, // Check every 5 seconds
    timeout: 3000, // 3 second timeout for network requests
    endpoints: ['/api/health', '/ping']
  },
  
  // UI feedback
  ui: {
    showOfflineBanner: true,
    showOfflineToast: true,
    offlineIndicator: true,
    syncIndicator: true
  },
  
  // Data synchronization
  sync: {
    autoSync: true,
    syncOnReconnect: true,
    conflictResolution: 'client-wins', // 'client-wins', 'server-wins', 'manual'
    maxPendingActions: 100
  },
  
  // Local storage
  storage: {
    quota: 50 * 1024 * 1024, // 50MB
    persistentStorage: true,
    compressionEnabled: true
  }
};

// Performance configuration
export const performanceConfig = {
  // Lazy loading
  lazyLoading: {
    enabled: true,
    rootMargin: '50px',
    threshold: 0.1,
    placeholderColor: '#f0f0f0'
  },
  
  // Image optimization
  images: {
    formats: ['webp', 'avif', 'jpg', 'png'],
    sizes: [320, 640, 960, 1280, 1920],
    quality: 85,
    progressive: true
  },
  
  // Code splitting
  codeSplitting: {
    enabled: true,
    chunkSize: 244 * 1024, // 244KB
    preloadCritical: true
  },
  
  // Resource hints
  resourceHints: {
    preconnect: [
      'https://fonts.googleapis.com',
      'https://fonts.gstatic.com',
      'https://api.nutrition-platform.com'
    ],
    prefetch: [
      '/api/user/profile',
      '/api/nutrition/favorites'
    ]
  }
};

// Theme configuration
export const themeConfig = {
  // Color schemes
  colorSchemes: {
    light: {
      primary: '#2563eb',
      secondary: '#64748b',
      accent: '#f59e0b',
      background: '#ffffff',
      surface: '#f8fafc',
      text: '#1e293b',
      textSecondary: '#64748b',
      border: '#e2e8f0',
      error: '#dc2626',
      warning: '#f59e0b',
      success: '#16a34a',
      info: '#0ea5e9'
    },
    dark: {
      primary: '#3b82f6',
      secondary: '#94a3b8',
      accent: '#fbbf24',
      background: '#0f172a',
      surface: '#1e293b',
      text: '#f1f5f9',
      textSecondary: '#94a3b8',
      border: '#334155',
      error: '#ef4444',
      warning: '#fbbf24',
      success: '#22c55e',
      info: '#38bdf8'
    },
    highContrast: {
      primary: '#000000',
      secondary: '#666666',
      accent: '#ff6600',
      background: '#ffffff',
      surface: '#f5f5f5',
      text: '#000000',
      textSecondary: '#333333',
      border: '#000000',
      error: '#cc0000',
      warning: '#ff6600',
      success: '#008800',
      info: '#0066cc'
    }
  },
  
  // Typography
  typography: {
    fontFamilies: {
      sans: ['Inter', 'system-ui', 'sans-serif'],
      serif: ['Georgia', 'serif'],
      mono: ['Fira Code', 'monospace'],
      arabic: ['Noto Sans Arabic', 'Arial', 'sans-serif'],
      chinese: ['Noto Sans SC', 'PingFang SC', 'sans-serif'],
      japanese: ['Noto Sans JP', 'Hiragino Sans', 'sans-serif']
    },
    
    fontSizes: {
      xs: '0.75rem',
      sm: '0.875rem',
      base: '1rem',
      lg: '1.125rem',
      xl: '1.25rem',
      '2xl': '1.5rem',
      '3xl': '1.875rem',
      '4xl': '2.25rem',
      '5xl': '3rem'
    },
    
    lineHeights: {
      tight: 1.25,
      normal: 1.5,
      relaxed: 1.75
    },
    
    fontWeights: {
      light: 300,
      normal: 400,
      medium: 500,
      semibold: 600,
      bold: 700
    }
  },
  
  // Spacing
  spacing: {
    xs: '0.25rem',
    sm: '0.5rem',
    md: '1rem',
    lg: '1.5rem',
    xl: '2rem',
    '2xl': '3rem',
    '3xl': '4rem'
  },
  
  // Border radius
  borderRadius: {
    none: '0',
    sm: '0.125rem',
    md: '0.375rem',
    lg: '0.5rem',
    xl: '0.75rem',
    full: '9999px'
  },
  
  // Shadows
  shadows: {
    sm: '0 1px 2px 0 rgba(0, 0, 0, 0.05)',
    md: '0 4px 6px -1px rgba(0, 0, 0, 0.1)',
    lg: '0 10px 15px -3px rgba(0, 0, 0, 0.1)',
    xl: '0 20px 25px -5px rgba(0, 0, 0, 0.1)'
  }
};

// Validation configuration
export const validationConfig = {
  // Form validation rules
  rules: {
    required: {
      message: 'This field is required',
      test: (value) => value !== null && value !== undefined && value !== ''
    },
    email: {
      message: 'Please enter a valid email address',
      test: (value) => /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value)
    },
    password: {
      message: 'Password must be at least 8 characters with uppercase, lowercase, number, and special character',
      test: (value) => /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$/.test(value)
    },
    phone: {
      message: 'Please enter a valid phone number',
      test: (value) => /^[+]?[1-9]?[0-9]{7,15}$/.test(value.replace(/\s/g, ''))
    },
    url: {
      message: 'Please enter a valid URL',
      test: (value) => /^https?:\/\/.+/.test(value)
    },
    number: {
      message: 'Please enter a valid number',
      test: (value) => !isNaN(parseFloat(value)) && isFinite(value)
    },
    positiveNumber: {
      message: 'Please enter a positive number',
      test: (value) => !isNaN(parseFloat(value)) && isFinite(value) && parseFloat(value) > 0
    }
  },
  
  // Real-time validation settings
  realTime: {
    enabled: true,
    debounceMs: 300,
    validateOnBlur: true,
    validateOnChange: true
  },
  
  // Error display
  errorDisplay: {
    showInline: true,
    showSummary: true,
    focusFirstError: true,
    scrollToError: true
  }
};

// Animation configuration
export const animationConfig = {
  // Transition durations
  durations: {
    fast: 150,
    normal: 300,
    slow: 500
  },
  
  // Easing functions
  easings: {
    linear: 'linear',
    easeIn: 'cubic-bezier(0.4, 0, 1, 1)',
    easeOut: 'cubic-bezier(0, 0, 0.2, 1)',
    easeInOut: 'cubic-bezier(0.4, 0, 0.2, 1)',
    bounce: 'cubic-bezier(0.68, -0.55, 0.265, 1.55)'
  },
  
  // Common animations
  animations: {
    fadeIn: {
      from: { opacity: 0 },
      to: { opacity: 1 }
    },
    slideUp: {
      from: { transform: 'translateY(20px)', opacity: 0 },
      to: { transform: 'translateY(0)', opacity: 1 }
    },
    slideDown: {
      from: { transform: 'translateY(-20px)', opacity: 0 },
      to: { transform: 'translateY(0)', opacity: 1 }
    },
    scaleIn: {
      from: { transform: 'scale(0.9)', opacity: 0 },
      to: { transform: 'scale(1)', opacity: 1 }
    }
  }
};

// Export utility functions
export const utils = {
  // Check if device is mobile
  isMobile: () => window.matchMedia(mediaQueries.mobile).matches,
  
  // Check if device supports touch
  isTouch: () => window.matchMedia(mediaQueries.touch).matches,
  
  // Check if user prefers reduced motion
  prefersReducedMotion: () => window.matchMedia(mediaQueries.reducedMotion).matches,
  
  // Check if user prefers dark mode
  prefersDarkMode: () => window.matchMedia(mediaQueries.darkMode).matches,
  
  // Get current language direction
  getDirection: (language = directionConfig.defaultLanguage) => {
    return directionConfig.languages[language]?.dir || 'ltr';
  },
  
  // Check if language is RTL
  isRTL: (language = directionConfig.defaultLanguage) => {
    return utils.getDirection(language) === 'rtl';
  },
  
  // Convert CSS property for RTL
  convertRTLProperty: (property, value, isRTL = false) => {
    if (!isRTL) return { [property]: value };
    
    const rtlProperty = directionConfig.rtlProperties[property];
    if (!rtlProperty) return { [property]: value };
    
    if (typeof rtlProperty === 'string') {
      return { [rtlProperty]: value };
    }
    
    if (typeof rtlProperty === 'object') {
      const convertedValue = rtlProperty[value] || value;
      return { [property]: convertedValue };
    }
    
    return { [property]: value };
  },
  
  // Generate responsive CSS
  responsive: (styles) => {
    const css = {};
    
    Object.entries(styles).forEach(([breakpoint, style]) => {
      if (breakpoint === 'base') {
        Object.assign(css, style);
      } else if (breakpoints[breakpoint]) {
        css[`@media (min-width: ${breakpoints[breakpoint]})`] = style;
      }
    });
    
    return css;
  },
  
  // Calculate contrast ratio
  getContrastRatio: (color1, color2) => {
    const getLuminance = (color) => {
      const rgb = parseInt(color.slice(1), 16);
      const r = (rgb >> 16) & 0xff;
      const g = (rgb >> 8) & 0xff;
      const b = (rgb >> 0) & 0xff;
      
      const [rs, gs, bs] = [r, g, b].map(c => {
        c = c / 255;
        return c <= 0.03928 ? c / 12.92 : Math.pow((c + 0.055) / 1.055, 2.4);
      });
      
      return 0.2126 * rs + 0.7152 * gs + 0.0722 * bs;
    };
    
    const l1 = getLuminance(color1);
    const l2 = getLuminance(color2);
    
    return (Math.max(l1, l2) + 0.05) / (Math.min(l1, l2) + 0.05);
  }
};

// Default export with all configurations
export default {
  breakpoints,
  mediaQueries,
  directionConfig,
  accessibilityConfig,
  pwaConfig,
  offlineConfig,
  performanceConfig,
  themeConfig,
  validationConfig,
  animationConfig,
  utils
};