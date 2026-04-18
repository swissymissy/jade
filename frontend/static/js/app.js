// === XSS SANITIZATION ===
function escapeHTML(str) {
  const div = document.createElement('div');
  div.textContent = str;
  return div.innerHTML;
}

// === LANGUAGE TOGGLE ===
let currentLang = 'en';

function toggleLang() {
  currentLang = currentLang === 'en' ? 'vi' : 'en';
  document.querySelector('.lang-toggle').textContent = currentLang === 'en' ? 'VI' : 'EN';

  document.querySelectorAll('[data-en]').forEach(el => {
    const text = el.getAttribute(`data-${currentLang}`);
    if (text) el.textContent = text;
  });
  document.querySelectorAll('[data-en-html]').forEach(el => {
    const html = el.getAttribute(`data-${currentLang}-html`);
    if (html) el.innerHTML = html;
  });
  document.querySelectorAll('[data-en-placeholder]').forEach(el => {
    const ph = el.getAttribute(`data-${currentLang}-placeholder`);
    if (ph) el.placeholder = ph;
  });
}

// === HERO CAROUSEL ===
const heroSlides = document.querySelectorAll('.hero-slide');
if (heroSlides.length > 1) {
  let heroIdx = 0;
  setInterval(() => {
    heroSlides[heroIdx].classList.remove('active');
    heroIdx = (heroIdx + 1) % heroSlides.length;
    heroSlides[heroIdx].classList.add('active');
  }, 4500);
}

// === NAV SCROLL ===
window.addEventListener('scroll', () => {
  document.getElementById('nav').classList.toggle('scrolled', window.scrollY > 20);
});

// === SCROLL REVEAL ===
const observer = new IntersectionObserver((entries) => {
  entries.forEach(entry => {
    if (entry.isIntersecting) {
      entry.target.classList.add('visible');
    }
  });
}, { threshold: 0.1 });

document.querySelectorAll('.reveal').forEach(el => observer.observe(el));

// === PRODUCT RENDERING ===
const API_BASE = '/api';

function renderProduct(product) {
  const coverSrc = product.cover_image
    ? product.cover_image.image_url
    : null;

  const safeName = escapeHTML(product.name);
  const safeType = escapeHTML(product.type);

  const imageHTML = coverSrc
    ? `<img src="${coverSrc}" alt="${safeName}" loading="lazy">`
    : `<div class="product-image-placeholder">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
          <path d="M21 16V8a2 2 0 00-1-1.73l-7-4a2 2 0 00-2 0l-7 4A2 2 0 002 8v8a2 2 0 001 1.73l7 4a2 2 0 002 0l7-4A2 2 0 0022 16z"/>
        </svg>
      </div>`;

  return `
    <div class="product-card" onclick="window.location.href='/products/${encodeURIComponent(product.slug)}'">
      <div class="product-image">${imageHTML}</div>
      <div class="product-info">
        <h3 class="product-name">${safeName}</h3>
        <p class="product-type">${safeType}</p>
        <p class="product-price">
          <span class="product-price-currency">$</span>${product.price.toFixed(2)}
        </p>
      </div>
    </div>
  `;
}

function renderSkeletons(count) {
  return Array(count).fill(`
    <div class="product-card">
      <div class="product-image"><div class="skeleton" style="width:100%;height:100%"></div></div>
      <div class="product-info">
        <div class="skeleton" style="height:20px;width:70%;margin-bottom:8px;border-radius:4px"></div>
        <div class="skeleton" style="height:14px;width:40%;margin-bottom:12px;border-radius:4px"></div>
        <div class="skeleton" style="height:22px;width:30%;border-radius:4px"></div>
      </div>
    </div>
  `).join('');
}

async function loadProducts() {
  const grid = document.getElementById('productsGrid');
  grid.innerHTML = renderSkeletons(6);

  try {
    const res = await fetch(`${API_BASE}/products`);
    const products = await res.json();

    if (products.length === 0) {
      grid.innerHTML = `<p style="text-align:center;grid-column:1/-1;color:var(--stone-400);font-size:15px;padding:60px 0;">No products found</p>`;
      return;
    }
    grid.innerHTML = products.map(renderProduct).join('');
  } catch (err) {
    console.error('Failed to load products:', err);
    grid.innerHTML = `<p style="text-align:center;grid-column:1/-1;color:var(--stone-400);font-size:15px;padding:60px 0;">Unable to load products</p>`;
  }
}

// === SEARCH ===
let searchTimeout;
document.getElementById('searchInput').addEventListener('input', (e) => {
  clearTimeout(searchTimeout);
  searchTimeout = setTimeout(() => {
    const query = e.target.value.trim();
    if (query.length === 0) {
      loadProducts();
      return;
    }
    searchProducts(query);
  }, 400);
});

async function searchProducts(query) {
  const grid = document.getElementById('productsGrid');
  grid.innerHTML = renderSkeletons(3);

  try {
    const res = await fetch(`${API_BASE}/products/search?q=${encodeURIComponent(query)}`);
    const products = await res.json();

    if (products.length === 0) {
      grid.innerHTML = `<p style="text-align:center;grid-column:1/-1;color:var(--stone-400);font-size:15px;padding:60px 0;">${currentLang === 'en' ? 'No results found' : 'Kh\u00f4ng t\u00ecm th\u1ea5y k\u1ebft qu\u1ea3'}</p>`;
      return;
    }
    grid.innerHTML = products.map(renderProduct).join('');
  } catch (err) {
    console.error('Search failed:', err);
  }
}

// === PRICE FILTER ===
async function applyFilter() {
  const min = document.getElementById('minPrice').value;
  const max = document.getElementById('maxPrice').value;
  const grid = document.getElementById('productsGrid');

  let url = `${API_BASE}/products/filter?`;
  if (min) url += `min=${min}&`;
  if (max) url += `max=${max}`;

  grid.innerHTML = renderSkeletons(3);

  try {
    const res = await fetch(url);
    const products = await res.json();

    if (products.length === 0) {
      grid.innerHTML = `<p style="text-align:center;grid-column:1/-1;color:var(--stone-400);font-size:15px;padding:60px 0;">${currentLang === 'en' ? 'No products in this price range' : 'Kh\u00f4ng c\u00f3 s\u1ea3n ph\u1ea9m trong kho\u1ea3ng gi\u00e1 n\u00e0y'}</p>`;
      return;
    }
    grid.innerHTML = products.map(renderProduct).join('');
  } catch (err) {
    console.error('Filter failed:', err);
  }
}

// === MOBILE MENU ===
function toggleMobileMenu() {
  const links = document.querySelector('.nav-links');
  links.style.display = links.style.display === 'flex' ? 'none' : 'flex';
}

// === INIT ===
loadProducts();
