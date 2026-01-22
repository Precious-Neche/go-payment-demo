// Wait for DOM to load
document.addEventListener("DOMContentLoaded", function () {
  // Get form elements
  const paymentForm = document.getElementById("paymentForm");
  const emailInput = document.getElementById("email");
  const amountInput = document.getElementById("amount");
  const productSelect = document.getElementById("product");
  const payButton = document.getElementById("payButton");
  const buttonText = document.getElementById("buttonText");
  const spinner = document.getElementById("spinner");
  const paymentStatus = document.getElementById("paymentStatus");
  const paymentLink = document.getElementById("paymentLink");
  const amountButtons = document.querySelectorAll(".amount-btn");

  // Define product amounts mapping
  const productAmounts = {
    'basic': 5000,
    'pro': 10000,
    'premium': 25000
  };

  // Define quick select amounts
  const quickSelectAmounts = [5000, 10000, 25000];

  // Function to sync product select based on amount
  function syncProductSelectFromAmount(amount) {
    // Remove active class from all amount buttons first
    amountButtons.forEach((btn) => btn.classList.remove("active"));
    
    // Check if amount matches any quick select button
    const matchingButton = Array.from(amountButtons).find(btn => 
      parseInt(btn.dataset.amount) === amount
    );
    
    if (matchingButton) {
      matchingButton.classList.add("active");
    }
    
    // Check if amount matches any product plan
    let matchedProduct = 'custom';
    
    if (amount === productAmounts.basic) {
      matchedProduct = 'basic';
    } else if (amount === productAmounts.pro) {
      matchedProduct = 'pro';
    } else if (amount === productAmounts.premium) {
      matchedProduct = 'premium';
    }
    
    // Update product select if different
    if (productSelect.value !== matchedProduct) {
      productSelect.value = matchedProduct;
    }
  }

  // Function to sync amount from product select
  function syncAmountFromProductSelect() {
    const product = productSelect.value;
    const currentAmount = parseInt(amountInput.value) || 0;
    
    // If selecting a predefined product, set the amount
    if (product !== 'custom' && productAmounts[product]) {
      const newAmount = productAmounts[product];
      amountInput.value = newAmount;
      
      // Sync quick select buttons
      syncProductSelectFromAmount(newAmount);
    } 
    // If selecting custom, don't change amount but sync quick select buttons
    else if (product === 'custom' && currentAmount > 0) {
      syncProductSelectFromAmount(currentAmount);
    }
  }

  // Quick amount buttons
  amountButtons.forEach((button) => {
    button.addEventListener("click", function () {
      const selectedAmount = parseInt(this.dataset.amount);
      
      // Set amount input value
      amountInput.value = selectedAmount;
      
      // Sync product select
      syncProductSelectFromAmount(selectedAmount);
    });
  });

  // Product selection change
  productSelect.addEventListener("change", function () {
    syncAmountFromProductSelect();
  });

  // Amount input change (for manual typing)
  amountInput.addEventListener("input", function () {
    const amount = parseInt(this.value) || 0;
    
    // Debounce to prevent too many updates
    clearTimeout(amountInput.debounceTimer);
    amountInput.debounceTimer = setTimeout(() => {
      syncProductSelectFromAmount(amount);
    }, 300); // 300ms delay
  });

  // Amount input blur (when user leaves the field)
  amountInput.addEventListener("blur", function () {
    const amount = parseInt(this.value) || 0;
    syncProductSelectFromAmount(amount);
  });

  // Form submission
  paymentForm.addEventListener("submit", async function (e) {
    e.preventDefault();

    // Get form values
    const email = emailInput.value.trim();
    const amount = parseInt(amountInput.value);

    // Validation
    if (!email || !amount) {
      showAlert("Please fill in all fields", "danger");
      return;
    }

    if (amount < 100) {
      showAlert("Minimum amount is ₦100", "danger");
      return;
    }

    // Show loading state
    setLoadingState(true);

    try {
      // Send request to backend
      const response = await fetch("/api/initialize", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          email: email,
          amount: amount,
          product: productSelect.value,
        }),
      });

      const data = await response.json();
      console.log("API Response:", data);

      if (!response.ok) {
        throw new Error(data.error || "Payment failed");
      }
      if (!data.status) {
        throw new Error(data.message || "Payment initialization failed");
      }

      // PAYSTACK RESPONSE STRUCTURE:
      // data.data.authorization_url - the payment link
      // data.data.reference - the reference code
      const paymentUrl = data.data?.authorization_url;
      const reference = data.data?.reference;

      if (!paymentUrl) {
        throw new Error("No payment URL received from server");
      }

      // Show success and payment link
      showAlert("Payment initialized successfully!", "success");

      paymentLink.innerHTML = `
            <div class="payment-link">
                <p class="mb-2"><strong>Payment Link:</strong></p>
                <a href="${paymentUrl}" target="_blank" class="btn btn-sm btn-success">
                    Click here to complete payment
                </a>
                <p class="mt-2 mb-0 small text-muted">Reference: ${reference}</p>
            </div>
        `;

      // Automatically redirect after 3 seconds (optional)
      setTimeout(() => {
        window.open(paymentUrl, "_blank");
      }, 3000);
    } catch (error) {
      showAlert("Error: " + error.message, "danger");
    } finally {
      // Reset loading state
      setLoadingState(false);
    }
  });

  // Helper functions
  function setLoadingState(isLoading) {
    if (isLoading) {
      buttonText.textContent = "Processing...";
      spinner.classList.remove("d-none");
      payButton.disabled = true;
    } else {
      buttonText.textContent = "Pay Now";
      spinner.classList.add("d-none");
      payButton.disabled = false;
    }
  }

  function showAlert(message, type) {
    paymentStatus.textContent = message;
    paymentStatus.className = `alert alert-${type}`;
    paymentStatus.classList.remove("d-none");

    // Auto-hide after 5 seconds
    setTimeout(() => {
      paymentStatus.classList.add("d-none");
    }, 5000);
  }

  // Check if there's a payment reference in URL
  const urlParams = new URLSearchParams(window.location.search);
  const reference = urlParams.get("reference");

  if (reference) {
    // Verify payment
    verifyPayment(reference);
  }
});

// Verify payment function
async function verifyPayment(reference) {
    try {
        // Make sure reference is valid
        if (!reference || reference === 'undefined') {
            console.error('Invalid reference:', reference);
            return;
        }
        
        const response = await fetch(`/api/verify/${reference}`);
        
        if (!response.ok) {
            throw new Error(`Verification failed: ${response.status}`);
        }
        
        const data = await response.json();
        console.log("Verification response:", data);
        
        // Handle both response structures
        const status = data.status || data.data?.status;
        
        if (status === 'success') {
            showPaymentResult('success', {
                reference: reference,
                amount: data.amount || data.data?.amount,
                status: status,
                message: data.message || 'Payment successful'
            });
        } else {
            showPaymentResult('failed', {
                reference: reference,
                amount: data.amount || data.data?.amount,
                status: status || 'failed',
                message: data.message || 'Payment could not be completed'
            });
        }
    } catch (error) {
        console.error('Verification error:', error);
    }
}

function showPaymentResult(status, data) {
  // Create result modal
  const modalHtml = `
        <div class="modal fade" id="paymentResultModal" tabindex="-1">
            <div class="modal-dialog modal-dialog-centered">
                <div class="modal-content">
                    <div class="modal-header border-0">
                        <h5 class="modal-title">
                            ${status === "success" ? "🎉 Payment Successful!" : "❌ Payment Failed"}
                        </h5>
                        <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
                    </div>
                    <div class="modal-body text-center">
                        ${
                          status === "success"
                            ? '<div class="text-success mb-3"><i class="fas fa-check-circle fa-4x"></i></div>'
                            : '<div class="text-danger mb-3"><i class="fas fa-times-circle fa-4x"></i></div>'
                        }
                        <h4>${data.message || (status === "success" ? "Thank you for your payment!" : "Payment could not be completed")}</h4>
                        <div class="mt-4 text-start">
                            <p><strong>Reference:</strong> ${data.reference || "N/A"}</p>
                            <p><strong>Amount:</strong> ₦${(data.amount / 100).toLocaleString()}</p>
                            <p><strong>Status:</strong> ${data.status || "unknown"}</p>
                            <p><strong>Date:</strong> ${new Date().toLocaleString()}</p>
                        </div>
                    </div>
                    <div class="modal-footer border-0">
                        <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
                        ${
                          status === "success"
                            ? '<a href="/" class="btn btn-primary">Back to Home</a>'
                            : '<button type="button" class="btn btn-primary" onclick="window.location.reload()">Try Again</button>'
                        }
                    </div>
                </div>
            </div>
        </div>
    `;

  // Add modal to body
  document.body.insertAdjacentHTML("beforeend", modalHtml);

  // Show modal
  const modal = new bootstrap.Modal(
    document.getElementById("paymentResultModal"),
  );
  modal.show();
}