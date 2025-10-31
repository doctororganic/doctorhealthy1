// Contact and subscription functions

// WhatsApp contact function
function contactWhatsApp() {
    const phoneNumber = '+201234567890'; // Replace with actual WhatsApp number
    const message = encodeURIComponent('Hello! I would like to know more about your nutrition services.');
    const whatsappUrl = `https://wa.me/${phoneNumber}?text=${message}`;
    window.open(whatsappUrl, '_blank');
}

// Medical subscription form function
function openMedicalSubscriptionForm() {
    // Create modal for medical subscription
    const modal = document.createElement('div');
    modal.className = 'modal fade';
    modal.id = 'medicalSubscriptionModal';
    modal.innerHTML = `
        <div class="modal-dialog modal-lg">
            <div class="modal-content">
                <div class="modal-header">
                    <h5 class="modal-title">Medical Subscription Form</h5>
                    <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
                </div>
                <div class="modal-body">
                    <form id="medicalSubscriptionForm">
                        <div class="row">
                            <div class="col-md-6 mb-3">
                                <label for="fullName" class="form-label">Full Name</label>
                                <input type="text" class="form-control" id="fullName" required>
                            </div>
                            <div class="col-md-6 mb-3">
                                <label for="email" class="form-label">Email</label>
                                <input type="email" class="form-control" id="email" required>
                            </div>
                        </div>
                        <div class="row">
                            <div class="col-md-6 mb-3">
                                <label for="phone" class="form-label">Phone Number</label>
                                <input type="tel" class="form-control" id="phone" required>
                            </div>
                            <div class="col-md-6 mb-3">
                                <label for="age" class="form-label">Age</label>
                                <input type="number" class="form-control" id="age" min="1" max="120" required>
                            </div>
                        </div>
                        <div class="mb-3">
                            <label for="medicalCondition" class="form-label">Medical Condition</label>
                            <select class="form-select" id="medicalCondition" required>
                                <option value="">Select Condition</option>
                                <option value="diabetes_type1">Type 1 Diabetes</option>
                                <option value="diabetes_type2">Type 2 Diabetes</option>
                                <option value="hypertension">Hypertension</option>
                                <option value="heart_disease">Heart Disease</option>
                                <option value="kidney_disease">Kidney Disease</option>
                                <option value="liver_disease">Liver Disease</option>
                                <option value="obesity">Obesity</option>
                                <option value="celiac">Celiac Disease</option>
                                <option value="ibs">IBS</option>
                                <option value="other">Other</option>
                            </select>
                        </div>
                        <div class="mb-3">
                            <label for="medications" class="form-label">Current Medications</label>
                            <textarea class="form-control" id="medications" rows="3" placeholder="List your current medications..."></textarea>
                        </div>
                        <div class="mb-3">
                            <label for="subscriptionType" class="form-label">Subscription Type</label>
                            <select class="form-select" id="subscriptionType" required>
                                <option value="">Select Subscription</option>
                                <option value="basic">Basic Plan - $29/month</option>
                                <option value="premium">Premium Plan - $49/month</option>
                                <option value="medical">Medical Plan - $79/month</option>
                            </select>
                        </div>
                    </form>
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Cancel</button>
                    <button type="button" class="btn btn-primary" onclick="submitMedicalSubscription()">Subscribe</button>
                </div>
            </div>
        </div>
    `;
    
    document.body.appendChild(modal);
    const bootstrapModal = new bootstrap.Modal(modal);
    bootstrapModal.show();
    
    // Clean up modal when closed
    modal.addEventListener('hidden.bs.modal', () => {
        document.body.removeChild(modal);
    });
}

// Submit medical subscription form
function submitMedicalSubscription() {
    const form = document.getElementById('medicalSubscriptionForm');
    const formData = new FormData(form);
    
    // Basic validation
    if (!form.checkValidity()) {
        form.reportValidity();
        return;
    }
    
    // Collect form data
    const subscriptionData = {
        fullName: document.getElementById('fullName').value,
        email: document.getElementById('email').value,
        phone: document.getElementById('phone').value,
        age: document.getElementById('age').value,
        medicalCondition: document.getElementById('medicalCondition').value,
        medications: document.getElementById('medications').value,
        subscriptionType: document.getElementById('subscriptionType').value
    };
    
    // Show success message
    alert('Thank you for your subscription request! We will contact you within 24 hours to confirm your medical plan.');
    
    // Close modal
    const modal = bootstrap.Modal.getInstance(document.getElementById('medicalSubscriptionModal'));
    modal.hide();
    
    // In a real application, you would send this data to your backend
    console.log('Medical subscription data:', subscriptionData);
}

// Make functions globally available
window.contactWhatsApp = contactWhatsApp;
window.openMedicalSubscriptionForm = openMedicalSubscriptionForm;
window.submitMedicalSubscription = submitMedicalSubscription;