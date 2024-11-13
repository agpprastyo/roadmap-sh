function toggleResetButton(form) {
  const inputs = form.querySelectorAll('input[type="text"], input[type="number"]');
  const resetButton = form.querySelector('button[type="reset"]');

  if (resetButton) { // Check if reset button is found
    let isEmpty = true;

    inputs.forEach(input => {
      if (input.value.trim()) {
        isEmpty = false; // If any input has a value, set to false
      }
    });

    // Show or hide the reset button based on input values
    resetButton.style.display = isEmpty ? 'none' : 'inline-block';
  }
}

function handleReset(form) {
  const inputs = form.querySelectorAll('input[type="text"], input[type="number"]');
  inputs.forEach(input => {
    input.value = ''; // Clear the input
  });
  toggleResetButton(form); // Update reset button visibility after reset
}
