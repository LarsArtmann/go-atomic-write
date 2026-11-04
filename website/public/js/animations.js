// Scroll-reveal for [data-animate] elements. Runs on astro:page-load so it
// re-binds after every ClientRouter navigation (the elements come fresh from
// the swapped DOM); without ClientRouter the event still fires once on load.
function bindAnimations() {
  const targets = document.querySelectorAll("[data-animate]:not([data-animate-bound])");

  if (window.matchMedia("(prefers-reduced-motion: reduce)").matches) {
    targets.forEach((el) => {
      el.classList.add("animate-fade-in");
      el.dataset.animateBound = "";
    });

    return;
  }

  const observer = new IntersectionObserver(
    (entries) => {
      entries.forEach((entry) => {
        if (entry.isIntersecting) {
          entry.target.classList.add("animate-fade-in");
          observer.unobserve(entry.target);
        }
      });
    },
    { threshold: 0.1 },
  );

  targets.forEach((el) => {
    el.dataset.animateBound = "";
    observer.observe(el);
  });
}

document.addEventListener("astro:page-load", bindAnimations);
