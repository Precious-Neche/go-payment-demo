// Payment page JavaScript
document.addEventListener('DOMContentLoaded', function() {
    // Auto-start countdown for redirect
    let countdown = 5;
    const countdownElement = document.getElementById('countdown');
    const countdownContainer = document.getElementById('countdownContainer');
    
    if (countdownElement && countdownContainer) {
        const countdownInterval = setInterval(() => {
            countdown--;
            countdownElement.textContent = countdown;
            
            if (countdown <= 0) {
                clearInterval(countdownInterval);
                // Auto-click the payment link
                const paystackLink = document.getElementById('paystackLink');
                if (paystackLink) {
                    paystackLink.click();
                }
                countdownContainer.innerHTML = '<p class="text-success">Redirecting to Paystack...</p>';
            }
        }, 1000);
    }
    
    // Copy reference to clipboard
    window.copyReference = function() {
        const reference = document.getElementById('reference').textContent;
        navigator.clipboard.writeText(reference).then(() => {
            alert('Reference copied to clipboard: ' + reference);
        });
    };
    
    // Cancel payment button
    const cancelBtn = document.getElementById('cancelBtn');
    if (cancelBtn) {
        cancelBtn.addEventListener('click', function() {
            if (confirm('Are you sure you want to cancel this payment?')) {
                window.location.href = '/';
            }
        });
    }
    
    // Update status when payment link is clicked
    const paystackLink = document.getElementById('paystackLink');
    if (paystackLink) {
        paystackLink.addEventListener('click', function() {
            const paymentStatus = document.getElementById('paymentStatus');
            if (paymentStatus) {
                paymentStatus.innerHTML = `
                    <div class="spinner-border spinner-border-sm me-2" role="status"></div>
                    Redirected to Paystack. Complete your payment on the new page...
                `;
                paymentStatus.className = 'alert alert-warning';
            }
        });
    }
    
    // Check payment status every 5 seconds (for async verification)
    const reference = document.getElementById('reference').textContent;
    if (reference) {
        setInterval(() => {
            fetch(`/api/payments/verify/${reference}`)
                .then(response => response.json())
                .then(data => {
                    if (data.status === 'success') {
                        // Redirect to success page
                        window.location.href = `/success?reference=${reference}`;
                    } else if (data.status === 'failed') {
                        // Redirect to failed page
                        window.location.href = `/failed?reference=${reference}`;
                    }
                })
                .catch(error => console.error('Verification error:', error));
        }, 5000);
    }
});
