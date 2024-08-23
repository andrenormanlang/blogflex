
  // Show spinner when an HTMX request is being sent
  document.addEventListener("htmx:configRequest", function(evt) {
    document.getElementById("loading-spinner").style.display = "block";
    document.getElementById("loading-message").style.display = "block";
  });

  // Hide spinner when the HTMX request is finished
  document.addEventListener("htmx:afterRequest", function(evt) {
    document.getElementById("loading-spinner").style.display = "none";
    document.getElementById("loading-message").style.display = "none";
  });

  // Handle errors by hiding the spinner and showing an error message
  document.addEventListener("htmx:responseError", function(evt) {
    document.getElementById("loading-spinner").style.display = "none";
    document.getElementById("loading-message").style.display = "none";
    alert("An error occurred while processing your request.");
  });

