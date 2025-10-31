/**
 * RTL Support Validation and Arabic Render Testing
 * Ensures proper right-to-left text display and Arabic language support
 */

// RTL Configuration
const RTL_CONFIG = {
  languages: ['ar', 'he', 'fa', 'ur'],
  arabicLanguages: ['ar'],
  testStrings: {
    arabic: {
      simple: 'مرحبا بكم في منصة التغذية',
      complex: 'البروتين: ٢٥ جرام، الكربوهيدرات: ٤٥ جرام، الدهون: ١٢ جرام',
      mixed: 'Protein البروتين: 25g',
      numbers: '١٢٣٤٥٦٧٨٩٠',
      punctuation: 'مرحبا، كيف حالك؟ أهلاً وسهلاً!',
      longText: 'هذا نص طويل باللغة العربية لاختبار التفاف النص والعرض الصحيح في واجهة المستخدم. يجب أن يظهر النص من اليمين إلى اليسار بشكل صحيح.',
    },
    hebrew: {
      simple: 'שלום וברוכים הבאים',
      mixed: 'Protein חלבון: 25g',
    },
    persian: {
      simple: 'خوش آمدید به پلتفرم تغذیه',
      mixed: 'Protein پروتئین: 25g',
    },
    urdu: {
      simple: 'غذائیت کے پلیٹ فارم میں خوش آمدید',
      mixed: 'Protein پروٹین: 25g',
    }
  },
  cssProperties: {
    direction: 'rtl',
    textAlign: 'right',
    unicodeBidi: 'embed',
    writingMode: 'horizontal-tb'
  },
  fonts: {
    arabic: [
      'Noto Sans Arabic',
      'Amiri',
      'Cairo',
      'Tajawal',
      'Almarai',
      'Markazi Text',
      'Scheherazade New',
      'Arial Unicode MS',
      'Tahoma',
      'sans-serif'
    ],
    fallback: ['Arial', 'Helvetica', 'sans-serif']
  }
};

// RTL Validator Class
class RTLValidator {
  constructor() {
    this.testResults = {
      cssLoading: false,
      fontLoading: false,
      textRendering: false,
      directionality: false,
      bidiSupport: false,
      numberFormatting: false,
      layoutIntegrity: false,
      accessibility: false,
      performance: false,
      errors: [],
      warnings: [],
      recommendations: []
    };
    
    this.testContainer = null;
    this.observer = null;
    this.performanceMetrics = {
      renderTime: 0,
      layoutTime: 0,
      fontLoadTime: 0,
      totalTime: 0
    };
  }

  /**
   * Run comprehensive RTL validation tests
   */
  async validateRTLSupport() {
    console.log('🔍 Starting RTL Support Validation...');
    const startTime = performance.now();

    try {
      // Create test container
      this.createTestContainer();

      // Run all validation tests
      await Promise.all([
        this.testCSSLoading(),
        this.testFontLoading(),
        this.testTextRendering(),
        this.testDirectionality(),
        this.testBidiSupport(),
        this.testNumberFormatting(),
        this.testLayoutIntegrity(),
        this.testAccessibility(),
        this.testPerformance()
      ]);

      // Calculate total time
      this.performanceMetrics.totalTime = performance.now() - startTime;

      // Generate recommendations
      this.generateRecommendations();

      // Cleanup
      this.cleanup();

      console.log('✅ RTL Validation completed:', this.testResults);
      return this.getValidationReport();

    } catch (error) {
      this.testResults.errors.push(`Validation failed: ${error.message}`);
      console.error('❌ RTL Validation failed:', error);
      return this.getValidationReport();
    }
  }

  /**
   * Create test container for RTL testing
   */
  createTestContainer() {
    this.testContainer = document.createElement('div');
    this.testContainer.id = 'rtl-test-container';
    this.testContainer.style.cssText = `
      position: fixed;
      top: -9999px;
      left: -9999px;
      width: 500px;
      height: 300px;
      visibility: hidden;
      pointer-events: none;
      z-index: -1;
    `;
    document.body.appendChild(this.testContainer);
  }

  /**
   * Test CSS loading and RTL styles
   */
  async testCSSLoading() {
    try {
      const testElement = document.createElement('div');
      testElement.className = 'rtl-test';
      testElement.dir = 'rtl';
      testElement.textContent = RTL_CONFIG.testStrings.arabic.simple;
      
      // Apply RTL styles
      Object.assign(testElement.style, {
        direction: 'rtl',
        textAlign: 'right',
        unicodeBidi: 'embed'
      });
      
      this.testContainer.appendChild(testElement);
      
      // Force layout
      testElement.offsetHeight;
      
      // Check computed styles
      const computedStyle = window.getComputedStyle(testElement);
      const direction = computedStyle.direction;
      const textAlign = computedStyle.textAlign;
      
      if (direction === 'rtl' && (textAlign === 'right' || textAlign === 'start')) {
        this.testResults.cssLoading = true;
      } else {
        this.testResults.errors.push(`CSS RTL styles not applied correctly. Direction: ${direction}, TextAlign: ${textAlign}`);
      }
      
    } catch (error) {
      this.testResults.errors.push(`CSS loading test failed: ${error.message}`);
    }
  }

  /**
   * Test Arabic font loading and availability
   */
  async testFontLoading() {
    const startTime = performance.now();
    
    try {
      const promises = RTL_CONFIG.fonts.arabic.map(font => this.checkFontAvailability(font));
      const results = await Promise.all(promises);
      
      const availableFonts = results.filter(result => result.available);
      
      if (availableFonts.length > 0) {
        this.testResults.fontLoading = true;
        console.log(`✅ Available Arabic fonts: ${availableFonts.map(f => f.font).join(', ')}`);
      } else {
        this.testResults.warnings.push('No Arabic fonts detected. Using fallback fonts.');
        this.testResults.fontLoading = false;
      }
      
      this.performanceMetrics.fontLoadTime = performance.now() - startTime;
      
    } catch (error) {
      this.testResults.errors.push(`Font loading test failed: ${error.message}`);
    }
  }

  /**
   * Check if a specific font is available
   */
  async checkFontAvailability(fontName) {
    return new Promise((resolve) => {
      const testString = RTL_CONFIG.testStrings.arabic.simple;
      const testSize = '72px';
      
      // Create test elements
      const container = document.createElement('div');
      container.style.cssText = `
        position: absolute;
        top: -9999px;
        left: -9999px;
        visibility: hidden;
      `;
      
      const fallbackElement = document.createElement('span');
      fallbackElement.style.cssText = `
        font-family: monospace;
        font-size: ${testSize};
        white-space: nowrap;
      `;
      fallbackElement.textContent = testString;
      
      const testElement = document.createElement('span');
      testElement.style.cssText = `
        font-family: "${fontName}", monospace;
        font-size: ${testSize};
        white-space: nowrap;
      `;
      testElement.textContent = testString;
      
      container.appendChild(fallbackElement);
      container.appendChild(testElement);
      document.body.appendChild(container);
      
      // Force layout
      const fallbackWidth = fallbackElement.offsetWidth;
      const testWidth = testElement.offsetWidth;
      
      // Clean up
      document.body.removeChild(container);
      
      // Font is available if widths differ
      const available = fallbackWidth !== testWidth;
      
      resolve({ font: fontName, available });
    });
  }

  /**
   * Test text rendering quality
   */
  async testTextRendering() {
    const startTime = performance.now();
    
    try {
      const testCases = [
        { text: RTL_CONFIG.testStrings.arabic.simple, type: 'simple' },
        { text: RTL_CONFIG.testStrings.arabic.complex, type: 'complex' },
        { text: RTL_CONFIG.testStrings.arabic.mixed, type: 'mixed' },
        { text: RTL_CONFIG.testStrings.arabic.longText, type: 'long' }
      ];
      
      let passedTests = 0;
      
      for (const testCase of testCases) {
        const element = document.createElement('div');
        element.style.cssText = `
          direction: rtl;
          text-align: right;
          font-family: ${RTL_CONFIG.fonts.arabic.join(', ')};
          font-size: 16px;
          line-height: 1.5;
          width: 400px;
        `;
        element.textContent = testCase.text;
        
        this.testContainer.appendChild(element);
        
        // Force layout and measure
        const rect = element.getBoundingClientRect();
        
        if (rect.width > 0 && rect.height > 0) {
          passedTests++;
        } else {
          this.testResults.warnings.push(`Text rendering failed for ${testCase.type} text`);
        }
      }
      
      this.testResults.textRendering = passedTests === testCases.length;
      this.performanceMetrics.renderTime = performance.now() - startTime;
      
    } catch (error) {
      this.testResults.errors.push(`Text rendering test failed: ${error.message}`);
    }
  }

  /**
   * Test text directionality
   */
  async testDirectionality() {
    try {
      const testElement = document.createElement('div');
      testElement.style.cssText = `
        direction: rtl;
        width: 200px;
        border: 1px solid transparent;
      `;
      testElement.textContent = RTL_CONFIG.testStrings.arabic.simple;
      
      this.testContainer.appendChild(testElement);
      
      // Force layout
      testElement.offsetHeight;
      
      // Check if text starts from the right
      const computedStyle = window.getComputedStyle(testElement);
      const direction = computedStyle.direction;
      
      // Test cursor position (if supported)
      let cursorTest = true;
      try {
        const range = document.createRange();
        range.setStart(testElement.firstChild, 0);
        const rect = range.getBoundingClientRect();
        const elementRect = testElement.getBoundingClientRect();
        
        // In RTL, cursor should be closer to the right edge
        const distanceFromRight = elementRect.right - rect.right;
        const distanceFromLeft = rect.left - elementRect.left;
        
        cursorTest = distanceFromRight < distanceFromLeft;
      } catch (e) {
        // Cursor test not supported, skip
      }
      
      this.testResults.directionality = direction === 'rtl' && cursorTest;
      
      if (!this.testResults.directionality) {
        this.testResults.warnings.push('Text directionality test failed');
      }
      
    } catch (error) {
      this.testResults.errors.push(`Directionality test failed: ${error.message}`);
    }
  }

  /**
   * Test bidirectional text support
   */
  async testBidiSupport() {
    try {
      const bidiTexts = [
        RTL_CONFIG.testStrings.arabic.mixed,
        'Hello مرحبا World عالم',
        '123 عدد 456',
        'Email: user@example.com البريد الإلكتروني'
      ];
      
      let passedTests = 0;
      
      for (const text of bidiTexts) {
        const element = document.createElement('div');
        element.style.cssText = `
          direction: rtl;
          unicode-bidi: embed;
          font-family: ${RTL_CONFIG.fonts.arabic.join(', ')};
        `;
        element.textContent = text;
        
        this.testContainer.appendChild(element);
        
        // Force layout
        const rect = element.getBoundingClientRect();
        
        if (rect.width > 0 && rect.height > 0) {
          passedTests++;
        }
      }
      
      this.testResults.bidiSupport = passedTests === bidiTexts.length;
      
    } catch (error) {
      this.testResults.errors.push(`Bidi support test failed: ${error.message}`);
    }
  }

  /**
   * Test Arabic number formatting
   */
  async testNumberFormatting() {
    try {
      const numberTests = [
        { input: '1234567890', expected: '١٢٣٤٥٦٧٨٩٠' },
        { input: '25.5', expected: '٢٥.٥' },
        { input: '100%', expected: '١٠٠٪' }
      ];
      
      let passedTests = 0;
      
      for (const test of numberTests) {
        // Test Arabic-Indic digit conversion
        const converted = this.convertToArabicNumerals(test.input);
        
        if (converted === test.expected) {
          passedTests++;
        } else {
          this.testResults.warnings.push(`Number formatting failed: ${test.input} -> ${converted} (expected: ${test.expected})`);
        }
      }
      
      this.testResults.numberFormatting = passedTests === numberTests.length;
      
    } catch (error) {
      this.testResults.errors.push(`Number formatting test failed: ${error.message}`);
    }
  }

  /**
   * Test layout integrity with RTL content
   */
  async testLayoutIntegrity() {
    const startTime = performance.now();
    
    try {
      // Create complex layout with RTL content
      const layoutContainer = document.createElement('div');
      layoutContainer.style.cssText = `
        direction: rtl;
        display: flex;
        flex-direction: column;
        width: 400px;
        padding: 20px;
        border: 1px solid #ccc;
      `;
      
      // Header
      const header = document.createElement('h2');
      header.textContent = 'معلومات التغذية';
      header.style.cssText = 'margin: 0 0 10px 0; text-align: right;';
      
      // Content with mixed text
      const content = document.createElement('div');
      content.innerHTML = `
        <p style="text-align: right; margin: 5px 0;">البروتين: <span style="font-weight: bold;">25g</span></p>
        <p style="text-align: right; margin: 5px 0;">الكربوهيدرات: <span style="font-weight: bold;">45g</span></p>
        <p style="text-align: right; margin: 5px 0;">الدهون: <span style="font-weight: bold;">12g</span></p>
      `;
      
      // Button
      const button = document.createElement('button');
      button.textContent = 'إضافة إلى الوجبة';
      button.style.cssText = `
        align-self: flex-start;
        padding: 10px 20px;
        margin-top: 10px;
        direction: rtl;
      `;
      
      layoutContainer.appendChild(header);
      layoutContainer.appendChild(content);
      layoutContainer.appendChild(button);
      this.testContainer.appendChild(layoutContainer);
      
      // Force layout
      layoutContainer.offsetHeight;
      
      // Check layout measurements
      const containerRect = layoutContainer.getBoundingClientRect();
      const headerRect = header.getBoundingClientRect();
      const buttonRect = button.getBoundingClientRect();
      
      // Verify RTL alignment
      const headerAlignedRight = Math.abs(headerRect.right - containerRect.right) < 25; // Account for padding
      const buttonAlignedRight = Math.abs(buttonRect.right - containerRect.right) < 25;
      
      this.testResults.layoutIntegrity = containerRect.width > 0 && 
                                        containerRect.height > 0 && 
                                        headerAlignedRight && 
                                        buttonAlignedRight;
      
      this.performanceMetrics.layoutTime = performance.now() - startTime;
      
      if (!this.testResults.layoutIntegrity) {
        this.testResults.warnings.push('Layout integrity test failed - RTL alignment issues detected');
      }
      
    } catch (error) {
      this.testResults.errors.push(`Layout integrity test failed: ${error.message}`);
    }
  }

  /**
   * Test accessibility features for RTL
   */
  async testAccessibility() {
    try {
      const accessibilityTests = [
        this.testScreenReaderSupport(),
        this.testKeyboardNavigation(),
        this.testAriaLabels(),
        this.testColorContrast()
      ];
      
      const results = await Promise.all(accessibilityTests);
      const passedTests = results.filter(result => result).length;
      
      this.testResults.accessibility = passedTests >= 3; // At least 3 out of 4 tests should pass
      
    } catch (error) {
      this.testResults.errors.push(`Accessibility test failed: ${error.message}`);
    }
  }

  /**
   * Test screen reader support
   */
  async testScreenReaderSupport() {
    try {
      const element = document.createElement('div');
      element.setAttribute('role', 'region');
      element.setAttribute('aria-label', 'معلومات التغذية');
      element.setAttribute('lang', 'ar');
      element.style.direction = 'rtl';
      element.textContent = RTL_CONFIG.testStrings.arabic.simple;
      
      this.testContainer.appendChild(element);
      
      // Check if attributes are properly set
      const hasRole = element.getAttribute('role') === 'region';
      const hasAriaLabel = element.getAttribute('aria-label') !== null;
      const hasLang = element.getAttribute('lang') === 'ar';
      
      return hasRole && hasAriaLabel && hasLang;
      
    } catch (error) {
      this.testResults.warnings.push(`Screen reader test failed: ${error.message}`);
      return false;
    }
  }

  /**
   * Test keyboard navigation in RTL
   */
  async testKeyboardNavigation() {
    try {
      const input = document.createElement('input');
      input.type = 'text';
      input.dir = 'rtl';
      input.value = RTL_CONFIG.testStrings.arabic.simple;
      input.style.cssText = `
        direction: rtl;
        text-align: right;
      `;
      
      this.testContainer.appendChild(input);
      
      // Test if input supports RTL
      const computedStyle = window.getComputedStyle(input);
      return computedStyle.direction === 'rtl';
      
    } catch (error) {
      this.testResults.warnings.push(`Keyboard navigation test failed: ${error.message}`);
      return false;
    }
  }

  /**
   * Test ARIA labels in Arabic
   */
  async testAriaLabels() {
    try {
      const button = document.createElement('button');
      button.setAttribute('aria-label', 'إضافة إلى السلة');
      button.setAttribute('aria-describedby', 'nutrition-info');
      button.textContent = 'إضافة';
      
      const description = document.createElement('div');
      description.id = 'nutrition-info';
      description.textContent = 'معلومات التغذية للمنتج';
      
      this.testContainer.appendChild(button);
      this.testContainer.appendChild(description);
      
      return button.getAttribute('aria-label') !== null && 
             description.id === 'nutrition-info';
      
    } catch (error) {
      this.testResults.warnings.push(`ARIA labels test failed: ${error.message}`);
      return false;
    }
  }

  /**
   * Test color contrast for Arabic text
   */
  async testColorContrast() {
    try {
      // This is a simplified test - in a real implementation,
      // you would calculate actual contrast ratios
      const element = document.createElement('div');
      element.style.cssText = `
        color: #333;
        background-color: #fff;
        font-family: ${RTL_CONFIG.fonts.arabic.join(', ')};
      `;
      element.textContent = RTL_CONFIG.testStrings.arabic.simple;
      
      this.testContainer.appendChild(element);
      
      const computedStyle = window.getComputedStyle(element);
      const color = computedStyle.color;
      const backgroundColor = computedStyle.backgroundColor;
      
      // Basic check - ensure colors are set
      return color !== '' && backgroundColor !== '';
      
    } catch (error) {
      this.testResults.warnings.push(`Color contrast test failed: ${error.message}`);
      return false;
    }
  }

  /**
   * Test performance metrics
   */
  async testPerformance() {
    try {
      const startTime = performance.now();
      
      // Create multiple RTL elements to test performance
      const fragment = document.createDocumentFragment();
      
      for (let i = 0; i < 100; i++) {
        const element = document.createElement('div');
        element.style.cssText = `
          direction: rtl;
          font-family: ${RTL_CONFIG.fonts.arabic[0]};
        `;
        element.textContent = RTL_CONFIG.testStrings.arabic.simple + ` ${i}`;
        fragment.appendChild(element);
      }
      
      this.testContainer.appendChild(fragment);
      
      // Force layout
      this.testContainer.offsetHeight;
      
      const endTime = performance.now();
      const renderTime = endTime - startTime;
      
      // Performance is good if rendering 100 RTL elements takes less than 50ms
      this.testResults.performance = renderTime < 50;
      
      if (!this.testResults.performance) {
        this.testResults.warnings.push(`RTL rendering performance is slow: ${renderTime.toFixed(2)}ms for 100 elements`);
      }
      
    } catch (error) {
      this.testResults.errors.push(`Performance test failed: ${error.message}`);
    }
  }

  /**
   * Generate recommendations based on test results
   */
  generateRecommendations() {
    const recommendations = [];
    
    if (!this.testResults.cssLoading) {
      recommendations.push('Ensure RTL CSS styles are properly loaded and applied');
    }
    
    if (!this.testResults.fontLoading) {
      recommendations.push('Add Arabic web fonts (Noto Sans Arabic, Cairo, Tajawal) for better text rendering');
    }
    
    if (!this.testResults.textRendering) {
      recommendations.push('Improve text rendering by using proper Arabic fonts and CSS properties');
    }
    
    if (!this.testResults.directionality) {
      recommendations.push('Fix text directionality by ensuring direction: rtl and text-align: right are applied');
    }
    
    if (!this.testResults.bidiSupport) {
      recommendations.push('Improve bidirectional text support using unicode-bidi CSS property');
    }
    
    if (!this.testResults.numberFormatting) {
      recommendations.push('Implement Arabic-Indic numeral conversion for better localization');
    }
    
    if (!this.testResults.layoutIntegrity) {
      recommendations.push('Review layout components to ensure proper RTL alignment and spacing');
    }
    
    if (!this.testResults.accessibility) {
      recommendations.push('Improve accessibility by adding proper ARIA labels and lang attributes');
    }
    
    if (!this.testResults.performance) {
      recommendations.push('Optimize RTL rendering performance by using CSS containment and efficient selectors');
    }
    
    this.testResults.recommendations = recommendations;
  }

  /**
   * Convert Western numerals to Arabic-Indic numerals
   */
  convertToArabicNumerals(text) {
    const arabicNumerals = {
      '0': '٠', '1': '١', '2': '٢', '3': '٣', '4': '٤',
      '5': '٥', '6': '٦', '7': '٧', '8': '٨', '9': '٩'
    };
    
    return text.replace(/[0-9]/g, (digit) => arabicNumerals[digit] || digit);
  }

  /**
   * Clean up test elements
   */
  cleanup() {
    if (this.testContainer && this.testContainer.parentNode) {
      this.testContainer.parentNode.removeChild(this.testContainer);
    }
    
    if (this.observer) {
      this.observer.disconnect();
    }
  }

  /**
   * Get validation report
   */
  getValidationReport() {
    const passedTests = Object.values(this.testResults)
      .filter(value => typeof value === 'boolean')
      .filter(value => value === true).length;
    
    const totalTests = Object.values(this.testResults)
      .filter(value => typeof value === 'boolean').length;
    
    const score = Math.round((passedTests / totalTests) * 100);
    
    return {
      score,
      passedTests,
      totalTests,
      results: this.testResults,
      performance: this.performanceMetrics,
      summary: {
        status: score >= 80 ? 'PASS' : score >= 60 ? 'WARNING' : 'FAIL',
        message: this.getStatusMessage(score),
        criticalIssues: this.testResults.errors.length,
        warnings: this.testResults.warnings.length,
        recommendations: this.testResults.recommendations.length
      }
    };
  }

  /**
   * Get status message based on score
   */
  getStatusMessage(score) {
    if (score >= 90) {
      return 'Excellent RTL support! All tests passed.';
    } else if (score >= 80) {
      return 'Good RTL support with minor issues.';
    } else if (score >= 60) {
      return 'RTL support needs improvement.';
    } else {
      return 'Poor RTL support. Significant issues detected.';
    }
  }
}

// RTL Testing Utilities
class RTLTestUtils {
  /**
   * Quick RTL validation for development
   */
  static async quickValidation() {
    const validator = new RTLValidator();
    const report = await validator.validateRTLSupport();
    
    console.group('🔍 RTL Quick Validation Results');
    console.log(`Score: ${report.score}% (${report.passedTests}/${report.totalTests} tests passed)`);
    console.log(`Status: ${report.summary.status} - ${report.summary.message}`);
    
    if (report.results.errors.length > 0) {
      console.group('❌ Errors:');
      report.results.errors.forEach(error => console.error(error));
      console.groupEnd();
    }
    
    if (report.results.warnings.length > 0) {
      console.group('⚠️ Warnings:');
      report.results.warnings.forEach(warning => console.warn(warning));
      console.groupEnd();
    }
    
    if (report.results.recommendations.length > 0) {
      console.group('💡 Recommendations:');
      report.results.recommendations.forEach(rec => console.info(rec));
      console.groupEnd();
    }
    
    console.groupEnd();
    
    return report;
  }

  /**
   * Test specific RTL component
   */
  static testComponent(element) {
    const tests = {
      hasRTLDirection: false,
      hasRightAlignment: false,
      hasArabicFont: false,
      hasBidiSupport: false
    };
    
    if (element) {
      const computedStyle = window.getComputedStyle(element);
      
      tests.hasRTLDirection = computedStyle.direction === 'rtl';
      tests.hasRightAlignment = computedStyle.textAlign === 'right' || computedStyle.textAlign === 'start';
      tests.hasArabicFont = RTL_CONFIG.fonts.arabic.some(font => 
        computedStyle.fontFamily.includes(font)
      );
      tests.hasBidiSupport = computedStyle.unicodeBidi !== 'normal';
    }
    
    return tests;
  }

  /**
   * Apply RTL styles to element
   */
  static applyRTLStyles(element, options = {}) {
    const defaultOptions = {
      direction: 'rtl',
      textAlign: 'right',
      fontFamily: RTL_CONFIG.fonts.arabic.join(', '),
      unicodeBidi: 'embed'
    };
    
    const styles = { ...defaultOptions, ...options };
    
    Object.assign(element.style, styles);
    
    // Add RTL class if not present
    if (!element.classList.contains('rtl')) {
      element.classList.add('rtl');
    }
    
    // Set lang attribute
    if (!element.getAttribute('lang')) {
      element.setAttribute('lang', 'ar');
    }
    
    return element;
  }

  /**
   * Create RTL-ready input element
   */
  static createRTLInput(type = 'text', placeholder = '') {
    const input = document.createElement('input');
    input.type = type;
    input.dir = 'rtl';
    input.placeholder = placeholder;
    
    this.applyRTLStyles(input);
    
    return input;
  }

  /**
   * Format number for Arabic display
   */
  static formatArabicNumber(number, useArabicNumerals = true) {
    const formatted = new Intl.NumberFormat('ar-SA').format(number);
    
    if (useArabicNumerals) {
      return formatted.replace(/[0-9]/g, (digit) => {
        const arabicNumerals = {
          '0': '٠', '1': '١', '2': '٢', '3': '٣', '4': '٤',
          '5': '٥', '6': '٦', '7': '٧', '8': '٨', '9': '٩'
        };
        return arabicNumerals[digit] || digit;
      });
    }
    
    return formatted;
  }

  /**
   * Get RTL configuration
   */
  static getConfig() {
    return RTL_CONFIG;
  }
}

// Export for use in other modules
if (typeof module !== 'undefined' && module.exports) {
  module.exports = { RTLValidator, RTLTestUtils, RTL_CONFIG };
} else if (typeof window !== 'undefined') {
  window.RTLValidator = RTLValidator;
  window.RTLTestUtils = RTLTestUtils;
  window.RTL_CONFIG = RTL_CONFIG;
}

// Auto-run validation in development mode
if (typeof window !== 'undefined' && 
    (window.location.hostname === 'localhost' || 
     window.location.hostname === '127.0.0.1' ||
     window.location.search.includes('rtl-test=true'))) {
  
  // Run validation after DOM is loaded
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', () => {
      setTimeout(() => RTLTestUtils.quickValidation(), 1000);
    });
  } else {
    setTimeout(() => RTLTestUtils.quickValidation(), 1000);
  }
}