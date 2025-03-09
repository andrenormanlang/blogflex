/**
 * BlogFlex Mobile-First JavaScript
 * Enhances the mobile experience with responsive behaviors
 */

document.addEventListener("DOMContentLoaded", function () {
  // Handle mobile navigation toggle
  const navbarToggler = document.querySelector(".navbar-toggler");
  if (navbarToggler) {
    navbarToggler.addEventListener("click", function () {
      const navbarCollapse = document.querySelector(".navbar-collapse");
      if (navbarCollapse) {
        navbarCollapse.classList.toggle("show");
      }
    });
  }

  // Add touch-friendly hover effects for cards on mobile
  const cards = document.querySelectorAll(".card");
  cards.forEach((card) => {
    card.addEventListener(
      "touchstart",
      function () {
        this.classList.add("card-touch-active");
      },
      { passive: true }
    );

    card.addEventListener(
      "touchend",
      function () {
        setTimeout(() => {
          this.classList.remove("card-touch-active");
        }, 150);
      },
      { passive: true }
    );
  });

  // Responsive image handling
  const contentImages = document.querySelectorAll(
    ".post-content img, .card-content img"
  );
  contentImages.forEach((img) => {
    // Remove any inline dimensions that might break responsive layout
    img.removeAttribute("width");
    img.removeAttribute("height");
    img.style.maxWidth = "100%";
    img.style.height = "auto";

    // Add loading="lazy" for better performance
    img.setAttribute("loading", "lazy");
  });

  // Enhance form elements for mobile
  const formControls = document.querySelectorAll(".form-control");
  formControls.forEach((control) => {
    // Increase touch target size on mobile
    if (window.innerWidth < 768) {
      control.classList.add("form-control-lg");
    }
  });

  // Handle viewport height issues on mobile browsers
  function setMobileViewportHeight() {
    const vh = window.innerHeight * 0.01;
    document.documentElement.style.setProperty("--vh", `${vh}px`);
  }

  // Set initial viewport height
  setMobileViewportHeight();

  // Update on resize and orientation change
  window.addEventListener("resize", setMobileViewportHeight);
  window.addEventListener("orientationchange", setMobileViewportHeight);
});

// Add CSS class for touch devices
if ("ontouchstart" in document.documentElement) {
  document.body.classList.add("touch-device");
}

// Improve scrolling performance on mobile
let scrollTimeout;
window.addEventListener(
  "scroll",
  function () {
    if (!document.body.classList.contains("is-scrolling")) {
      document.body.classList.add("is-scrolling");
    }

    clearTimeout(scrollTimeout);
    scrollTimeout = setTimeout(function () {
      document.body.classList.remove("is-scrolling");
    }, 200);
  },
  { passive: true }
);
