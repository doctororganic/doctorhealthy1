// Enhanced Validation System using Validator.js v10
class ValidationSystem {
    constructor() {
        this.errors = [];
        this.currentLanguage = localStorage.getItem('language') || 'en';
    }

    // Clear previous errors
    clearErrors() {
        this.errors = [];
        this.removeErrorMessages();
    }

    // Add error to the list
    addError(field, message) {
        this.errors.push({ field, message });
    }

    // Check if validation passed
    isValid() {
        return this.errors.length === 0;
    }

    // Get all errors
    getErrors() {
        return this.errors;
    }

    // Display errors in UI
    displayErrors() {
        this.errors.forEach(error => {
            this.showFieldError(error.field, error.message);
        });
    }

    // Show error for specific field
    showFieldError(fieldId, message) {
        const field = document.getElementById(fieldId);
        if (field) {
            field.classList.add('is-invalid');
            
            // Remove existing error message
            const existingError = field.parentNode.querySelector('.invalid-feedback');
            if (existingError) {
                existingError.remove();
            }

            // Add new error message
            const errorDiv = document.createElement('div');
            errorDiv.className = 'invalid-feedback';
            errorDiv.textContent = message;
            field.parentNode.appendChild(errorDiv);
        }
    }

    // Remove all error messages from UI
    removeErrorMessages() {
        document.querySelectorAll('.is-invalid').forEach(field => {
            field.classList.remove('is-invalid');
        });
        document.querySelectorAll('.invalid-feedback').forEach(error => {
            error.remove();
        });
    }

    // Email validation
    validateEmail(email, fieldId = 'email') {
        if (!email || email.trim() === '') {
            this.addError(fieldId, this.getErrorMessage('email_required'));
            return false;
        }
        
        if (!validator.isEmail(email)) {
            this.addError(fieldId, this.getErrorMessage('email_invalid'));
            return false;
        }
        
        return true;
    }

    // Password validation
    validatePassword(password, fieldId = 'password') {
        if (!password || password.trim() === '') {
            this.addError(fieldId, this.getErrorMessage('password_required'));
            return false;
        }
        
        if (!validator.isLength(password, { min: 8 })) {
            this.addError(fieldId, this.getErrorMessage('password_min_length'));
            return false;
        }
        
        if (!validator.matches(password, /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)/)) {
            this.addError(fieldId, this.getErrorMessage('password_complexity'));
            return false;
        }
        
        return true;
    }

    // Confirm password validation
    validateConfirmPassword(password, confirmPassword, fieldId = 'confirmPassword') {
        if (!confirmPassword || confirmPassword.trim() === '') {
            this.addError(fieldId, this.getErrorMessage('confirm_password_required'));
            return false;
        }
        
        if (!validator.equals(password, confirmPassword)) {
            this.addError(fieldId, this.getErrorMessage('passwords_not_match'));
            return false;
        }
        
        return true;
    }

    // Name validation
    validateName(name, fieldId = 'name') {
        if (!name || name.trim() === '') {
            this.addError(fieldId, this.getErrorMessage('name_required'));
            return false;
        }
        
        if (!validator.isLength(name.trim(), { min: 2, max: 50 })) {
            this.addError(fieldId, this.getErrorMessage('name_length'));
            return false;
        }
        
        if (!validator.matches(name.trim(), /^[a-zA-Z\u0600-\u06FF\s]+$/)) {
            this.addError(fieldId, this.getErrorMessage('name_invalid_chars'));
            return false;
        }
        
        return true;
    }

    // Age validation
    validateAge(age, fieldId = 'age', min = 13, max = 100) {
        if (!age || age === '') {
            this.addError(fieldId, this.getErrorMessage('age_required'));
            return false;
        }
        
        const ageNum = parseInt(age);
        if (!validator.isInt(age.toString(), { min, max })) {
            this.addError(fieldId, this.getErrorMessage('age_range').replace('{min}', min).replace('{max}', max));
            return false;
        }
        
        return true;
    }

    // Weight validation
    validateWeight(weight, fieldId = 'weight', min = 30, max = 300) {
        if (!weight || weight === '') {
            this.addError(fieldId, this.getErrorMessage('weight_required'));
            return false;
        }
        
        if (!validator.isFloat(weight.toString(), { min, max })) {
            this.addError(fieldId, this.getErrorMessage('weight_range').replace('{min}', min).replace('{max}', max));
            return false;
        }
        
        return true;
    }

    // Height validation
    validateHeight(height, fieldId = 'height', min = 100, max = 250) {
        if (!height || height === '') {
            this.addError(fieldId, this.getErrorMessage('height_required'));
            return false;
        }
        
        if (!validator.isFloat(height.toString(), { min, max })) {
            this.addError(fieldId, this.getErrorMessage('height_range').replace('{min}', min).replace('{max}', max));
            return false;
        }
        
        return true;
    }

    // Phone number validation
    validatePhone(phone, fieldId = 'phone') {
        if (!phone || phone.trim() === '') {
            this.addError(fieldId, this.getErrorMessage('phone_required'));
            return false;
        }
        
        // Remove spaces and special characters for validation
        const cleanPhone = phone.replace(/[\s\-\(\)\+]/g, '');
        
        if (!validator.isMobilePhone(cleanPhone, 'any', { strictMode: false })) {
            this.addError(fieldId, this.getErrorMessage('phone_invalid'));
            return false;
        }
        
        return true;
    }

    // URL validation
    validateURL(url, fieldId = 'url') {
        if (!url || url.trim() === '') {
            this.addError(fieldId, this.getErrorMessage('url_required'));
            return false;
        }
        
        if (!validator.isURL(url, { require_protocol: true })) {
            this.addError(fieldId, this.getErrorMessage('url_invalid'));
            return false;
        }
        
        return true;
    }

    // Date validation
    validateDate(date, fieldId = 'date') {
        if (!date || date.trim() === '') {
            this.addError(fieldId, this.getErrorMessage('date_required'));
            return false;
        }
        
        if (!validator.isISO8601(date)) {
            this.addError(fieldId, this.getErrorMessage('date_invalid'));
            return false;
        }
        
        return true;
    }

    // Required field validation
    validateRequired(value, fieldId, fieldName) {
        if (!value || value.toString().trim() === '') {
            this.addError(fieldId, this.getErrorMessage('field_required').replace('{field}', fieldName));
            return false;
        }
        return true;
    }

    // Numeric validation
    validateNumeric(value, fieldId, min = null, max = null) {
        if (!value || value === '') {
            this.addError(fieldId, this.getErrorMessage('numeric_required'));
            return false;
        }
        
        if (!validator.isNumeric(value.toString())) {
            this.addError(fieldId, this.getErrorMessage('numeric_invalid'));
            return false;
        }
        
        if (min !== null && parseFloat(value) < min) {
            this.addError(fieldId, this.getErrorMessage('numeric_min').replace('{min}', min));
            return false;
        }
        
        if (max !== null && parseFloat(value) > max) {
            this.addError(fieldId, this.getErrorMessage('numeric_max').replace('{max}', max));
            return false;
        }
        
        return true;
    }

    // File validation
    validateFile(file, fieldId = 'file', allowedTypes = ['image/jpeg', 'image/png', 'image/gif'], maxSize = 5 * 1024 * 1024) {
        if (!file) {
            this.addError(fieldId, this.getErrorMessage('file_required'));
            return false;
        }
        
        if (!allowedTypes.includes(file.type)) {
            this.addError(fieldId, this.getErrorMessage('file_type_invalid'));
            return false;
        }
        
        if (file.size > maxSize) {
            this.addError(fieldId, this.getErrorMessage('file_size_invalid').replace('{size}', Math.round(maxSize / 1024 / 1024)));
            return false;
        }
        
        return true;
    }

    // Get localized error messages
    getErrorMessage(key) {
        const messages = {
            en: {
                email_required: 'Email is required',
                email_invalid: 'Please enter a valid email address',
                password_required: 'Password is required',
                password_min_length: 'Password must be at least 8 characters long',
                password_complexity: 'Password must contain at least one uppercase letter, one lowercase letter, and one number',
                confirm_password_required: 'Please confirm your password',
                passwords_not_match: 'Passwords do not match',
                name_required: 'Name is required',
                name_length: 'Name must be between 2 and 50 characters',
                name_invalid_chars: 'Name can only contain letters and spaces',
                age_required: 'Age is required',
                age_range: 'Age must be between {min} and {max} years',
                weight_required: 'Weight is required',
                weight_range: 'Weight must be between {min} and {max} kg',
                height_required: 'Height is required',
                height_range: 'Height must be between {min} and {max} cm',
                phone_required: 'Phone number is required',
                phone_invalid: 'Please enter a valid phone number',
                url_required: 'URL is required',
                url_invalid: 'Please enter a valid URL',
                date_required: 'Date is required',
                date_invalid: 'Please enter a valid date',
                field_required: '{field} is required',
                numeric_required: 'This field is required',
                numeric_invalid: 'Please enter a valid number',
                numeric_min: 'Value must be at least {min}',
                numeric_max: 'Value must not exceed {max}',
                file_required: 'Please select a file',
                file_type_invalid: 'Invalid file type. Please select a valid image file',
                file_size_invalid: 'File size must not exceed {size} MB'
            },
            ar: {
                email_required: 'البريد الإلكتروني مطلوب',
                email_invalid: 'يرجى إدخال عنوان بريد إلكتروني صحيح',
                password_required: 'كلمة المرور مطلوبة',
                password_min_length: 'يجب أن تكون كلمة المرور 8 أحرف على الأقل',
                password_complexity: 'يجب أن تحتوي كلمة المرور على حرف كبير وحرف صغير ورقم واحد على الأقل',
                confirm_password_required: 'يرجى تأكيد كلمة المرور',
                passwords_not_match: 'كلمات المرور غير متطابقة',
                name_required: 'الاسم مطلوب',
                name_length: 'يجب أن يكون الاسم بين 2 و 50 حرف',
                name_invalid_chars: 'يمكن أن يحتوي الاسم على أحرف ومسافات فقط',
                age_required: 'العمر مطلوب',
                age_range: 'يجب أن يكون العمر بين {min} و {max} سنة',
                weight_required: 'الوزن مطلوب',
                weight_range: 'يجب أن يكون الوزن بين {min} و {max} كيلو',
                height_required: 'الطول مطلوب',
                height_range: 'يجب أن يكون الطول بين {min} و {max} سم',
                phone_required: 'رقم الهاتف مطلوب',
                phone_invalid: 'يرجى إدخال رقم هاتف صحيح',
                url_required: 'الرابط مطلوب',
                url_invalid: 'يرجى إدخال رابط صحيح',
                date_required: 'التاريخ مطلوب',
                date_invalid: 'يرجى إدخال تاريخ صحيح',
                field_required: '{field} مطلوب',
                numeric_required: 'هذا الحقل مطلوب',
                numeric_invalid: 'يرجى إدخال رقم صحيح',
                numeric_min: 'يجب أن تكون القيمة {min} على الأقل',
                numeric_max: 'يجب ألا تتجاوز القيمة {max}',
                file_required: 'يرجى اختيار ملف',
                file_type_invalid: 'نوع ملف غير صحيح. يرجى اختيار ملف صورة صحيح',
                file_size_invalid: 'يجب ألا يتجاوز حجم الملف {size} ميجابايت'
            }
        };
        
        return messages[this.currentLanguage]?.[key] || messages.en[key] || key;
    }

    // Update language
    updateLanguage(language) {
        this.currentLanguage = language;
    }

    // Validate entire form
    validateForm(formId, validationRules) {
        this.clearErrors();
        const form = document.getElementById(formId);
        
        if (!form) {
            console.error(`Form with ID '${formId}' not found`);
            return false;
        }
        
        const formData = new FormData(form);
        
        validationRules.forEach(rule => {
            const value = formData.get(rule.field) || '';
            
            if (rule.required && !this.validateRequired(value, rule.field, rule.name || rule.field)) {
                return;
            }
            
            if (value && rule.type) {
                switch (rule.type) {
                    case 'email':
                        this.validateEmail(value, rule.field);
                        break;
                    case 'password':
                        this.validatePassword(value, rule.field);
                        break;
                    case 'name':
                        this.validateName(value, rule.field);
                        break;
                    case 'age':
                        this.validateAge(value, rule.field, rule.min, rule.max);
                        break;
                    case 'weight':
                        this.validateWeight(value, rule.field, rule.min, rule.max);
                        break;
                    case 'height':
                        this.validateHeight(value, rule.field, rule.min, rule.max);
                        break;
                    case 'phone':
                        this.validatePhone(value, rule.field);
                        break;
                    case 'url':
                        this.validateURL(value, rule.field);
                        break;
                    case 'date':
                        this.validateDate(value, rule.field);
                        break;
                    case 'numeric':
                        this.validateNumeric(value, rule.field, rule.min, rule.max);
                        break;
                }
            }
        });
        
        if (!this.isValid()) {
            this.displayErrors();
            return false;
        }
        
        return true;
    }
}

// Create global validation instance
const validation = new ValidationSystem();

// Export for use in other modules
if (typeof module !== 'undefined' && module.exports) {
    module.exports = ValidationSystem;
}

// Make available globally
window.validation = validation;
window.ValidationSystem = ValidationSystem;