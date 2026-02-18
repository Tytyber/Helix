document.querySelector(".btn-primary").addEventListener("click", () => {
    alert("Добро пожаловать в Helix 🚀");
});

document.querySelector(".btn-secondary").addEventListener("click", () => {
    document.getElementById("features").scrollIntoView({
        behavior: "smooth"
    });
});

// Анимация появления
const elements = document.querySelectorAll(".fade-in");

const observer = new IntersectionObserver((entries) => {
    entries.forEach(entry => {
        if (entry.isIntersecting) {
            entry.target.classList.add("visible");
        }
    });
}, { threshold: 0.2 });

elements.forEach(el => observer.observe(el));


// Плавный скролл
document.querySelector(".btn-secondary").addEventListener("click", () => {
    document.getElementById("features").scrollIntoView({
        behavior: "smooth"
    });
});


// Параллакс фона
document.addEventListener("mousemove", (e) => {
    const glow1 = document.querySelector(".glow1");
    const glow2 = document.querySelector(".glow2");

    const x = (e.clientX / window.innerWidth - 0.5) * 30;
    const y = (e.clientY / window.innerHeight - 0.5) * 30;

    glow1.style.transform = `translate(${x}px, ${y}px)`;
    glow2.style.transform = `translate(${-x}px, ${-y}px)`;
});
