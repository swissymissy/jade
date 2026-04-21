(function () {
  const state = {
    products: [],
    search: '',
  };

  const els = {
    createForm: document.getElementById('createProductForm'),
    feedback: document.getElementById('dashboardFeedback'),
    products: document.getElementById('adminProducts'),
    search: document.getElementById('productSearch'),
    statTotal: document.getElementById('statTotal'),
    statVisible: document.getElementById('statVisible'),
    statHidden: document.getElementById('statHidden'),
  };

  function escapeHTML(value) {
    const div = document.createElement('div');
    div.textContent = value ?? '';
    return div.innerHTML;
  }

  function readNullString(value) {
    if (!value) return '';
    if (typeof value === 'string') return value;
    if (typeof value === 'object' && value.Valid) return value.String || '';
    return '';
  }

  function slugify(value) {
    return String(value || '')
      .trim()
      .toLowerCase()
      .split(/\s+/)
      .filter(Boolean)
      .join('-');
  }

  function formatMoney(value) {
    const amount = Number(value || 0);
    return Number.isFinite(amount) ? amount.toFixed(2) : '0.00';
  }

  function formatDate(value) {
    if (!value) return '--';
    const parsed = new Date(value.replace(' ', 'T') + 'Z');
    if (Number.isNaN(parsed.getTime())) return value;
    return parsed.toLocaleString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: 'numeric',
      minute: '2-digit',
    });
  }

  function setFeedback(message, kind = 'info') {
    if (!message) {
      els.feedback.textContent = '';
      els.feedback.className = 'admin-feedback hidden';
      return;
    }
    els.feedback.textContent = message;
    els.feedback.className = `admin-feedback admin-feedback-${kind}`;
  }

  function setBusy(container, busy) {
    if (!container) return;
    container.querySelectorAll('button, input, textarea, select').forEach((el) => {
      el.disabled = busy;
    });
  }

  async function requestJSON(url, options = {}) {
    const config = {
      credentials: 'same-origin',
      ...options,
    };
    const response = await fetch(url, config);
    const data = await response.json().catch(() => ({}));

    if (response.status === 401) {
      window.location.href = '/admin/login';
      throw new Error('Your admin session expired. Please sign in again.');
    }

    if (!response.ok) {
      throw new Error(data.error || 'Request failed.');
    }
    return data;
  }

  function filteredProducts() {
    const query = state.search.trim().toLowerCase();
    if (!query) return state.products;
    return state.products.filter((product) => {
      return [product.name, product.type, product.slug]
        .map((value) => String(value || '').toLowerCase())
        .some((value) => value.includes(query));
    });
  }

  function productImageHTML(product) {
    const url = product.cover_image && product.cover_image.image_url;
    if (url) {
      return `<img src="${url}" alt="${escapeHTML(product.name || 'Product image')}">`;
    }

    return `
      <div class="admin-image-placeholder">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.4" aria-hidden="true">
          <path d="M4 7a2 2 0 0 1 2-2h12a2 2 0 0 1 2 2v10a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2z"/>
          <path d="m8 14 2.5-2.5L15 16"/>
          <circle cx="9" cy="10" r="1.5"/>
        </svg>
        <span>No image yet</span>
      </div>
    `;
  }

  function productGalleryHTML(product) {
    const images = Array.isArray(product.images) ? product.images : [];
    if (!images.length) {
      return `
        <div class="admin-gallery-empty">
          <p>No images uploaded yet. Use the upload card to add photos.</p>
        </div>
      `;
    }

    return `
      <ul class="admin-gallery-grid">
        ${images.map((image) => {
          const isCover = Number(image.cover) === 1;
          return `
            <li class="admin-gallery-item${isCover ? ' is-cover' : ''}">
              <img src="${image.image_url}" alt="${escapeHTML(product.name || 'Product image')}" loading="lazy">
              ${isCover ? '<span class="admin-gallery-tag">Cover</span>' : ''}
              <button type="button" class="admin-gallery-delete" data-delete-image="${image.id}" aria-label="Delete image">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6" aria-hidden="true">
                  <path d="M4 7h16"/>
                  <path d="M10 11v6"/>
                  <path d="M14 11v6"/>
                  <path d="M6 7l1 13a2 2 0 0 0 2 2h6a2 2 0 0 0 2-2l1-13"/>
                  <path d="M9 7V4a1 1 0 0 1 1-1h4a1 1 0 0 1 1 1v3"/>
                </svg>
              </button>
            </li>
          `;
        }).join('')}
      </ul>
    `;
  }

  function renderProducts() {
    const products = filteredProducts();
    const total = state.products.length;
    const visible = state.products.filter((product) => Number(product.is_available) === 1).length;
    const hidden = total - visible;

    els.statTotal.textContent = String(total);
    els.statVisible.textContent = String(visible);
    els.statHidden.textContent = String(hidden);

    if (!products.length) {
      els.products.innerHTML = `
        <article class="admin-empty-state admin-panel">
          <p class="eyebrow">Catalog</p>
          <h2>No products match this view</h2>
          <p>Try a different search term, or create a new product above.</p>
        </article>
      `;
      return;
    }

    els.products.innerHTML = products.map((product) => {
      const description = readNullString(product.description);
      const about = readNullString(product.about);
      const videoKey = readNullString(product.video_url);
      const live = Number(product.is_available) === 1;

      return `
        <article class="admin-product-card admin-panel" data-product-id="${product.id}">
          <div class="admin-product-top">
            <div class="admin-product-cover">${productImageHTML(product)}</div>
            <div class="admin-product-summary">
              <div class="admin-product-meta">
                <p class="admin-product-code">Product #${product.id}</p>
                <h3>${escapeHTML(product.name || 'Untitled product')}</h3>
                <p class="admin-product-slug">/${escapeHTML(product.slug || slugify(product.name))}</p>
              </div>
              <div class="admin-product-badges">
                <span class="admin-badge ${live ? 'is-live' : 'is-hidden'}">${live ? 'Visible' : 'Hidden'}</span>
                <span class="admin-badge">Qty ${Number(product.quantity || 0)}</span>
                <span class="admin-badge">USD ${formatMoney(product.price)}</span>
                <span class="admin-badge">${videoKey ? 'Video attached' : 'No video'}</span>
              </div>
              <p class="admin-product-note">Last updated ${escapeHTML(formatDate(product.updated_at || product.created_at))}</p>
            </div>
          </div>

          <form class="admin-form admin-form-grid admin-product-form" data-update-product="${product.id}" novalidate>
            <label>
              <span>Name</span>
              <input type="text" name="name" value="${escapeHTML(product.name || '')}" required maxlength="120" data-name-input>
            </label>
            <label>
              <span>Type</span>
              <input type="text" name="type" value="${escapeHTML(product.type || '')}" required maxlength="80">
            </label>
            <label>
              <span>Price (USD)</span>
              <input type="number" name="price" value="${formatMoney(product.price)}" min="0.01" step="0.01" required>
            </label>
            <label>
              <span>Quantity</span>
              <input type="number" name="quantity" value="${Number(product.quantity || 0)}" min="0" step="1" required>
            </label>
            <label>
              <span>Availability</span>
              <select name="is_available">
                <option value="1"${live ? ' selected' : ''}>Visible on storefront</option>
                <option value="0"${live ? '' : ' selected'}>Hidden from storefront</option>
              </select>
            </label>
            <label>
              <span>Slug</span>
              <input type="text" value="${escapeHTML(product.slug || slugify(product.name))}" readonly data-slug-preview>
            </label>
            <label class="admin-field-wide">
              <span>Description (tech specs)</span>
              <textarea name="description" rows="4" placeholder="Gemstone, color, size, certificate — the factual details.">${escapeHTML(description)}</textarea>
            </label>
            <label class="admin-field-wide">
              <span>About this piece (story / hook)</span>
              <textarea name="about" rows="4" placeholder="The sensory hook — how the piece feels, the mood, what makes it special.">${escapeHTML(about)}</textarea>
            </label>
            <div class="admin-form-actions">
              <button type="submit" class="btn btn-dark">Save changes</button>
              <button type="button" class="btn btn-outline-dark" data-delete-product="${product.id}">Delete product</button>
            </div>
          </form>

          <div class="admin-gallery" data-product-gallery="${product.id}">
            <div class="admin-gallery-head">
              <p class="eyebrow">Current images</p>
              <h4>Uploaded photos</h4>
            </div>
            ${productGalleryHTML(product)}
          </div>

          <div class="admin-media-grid">
            <form class="admin-upload-card" data-image-upload="${product.id}">
              <div class="admin-upload-copy">
                <p class="eyebrow">Images</p>
                <h4>Upload product photos</h4>
                <p>Add up to five images per product. The first image becomes the cover photo automatically.</p>
              </div>
              <label>
                <span>Select image files</span>
                <input type="file" name="images" accept=".jpg,.jpeg,.png,.webp" multiple>
              </label>
              <button type="submit" class="btn btn-outline-dark">Upload images</button>
            </form>

            <form class="admin-upload-card" data-video-upload="${product.id}">
              <div class="admin-upload-copy">
                <p class="eyebrow">Video</p>
                <h4>Replace product video</h4>
                <p>${escapeHTML(videoKey ? `Current video key: ${videoKey}` : 'No product video uploaded yet.')}</p>
              </div>
              <label>
                <span>Select video file</span>
                <input type="file" name="video" accept=".mp4,.mov,.webm">
              </label>
              <button type="submit" class="btn btn-outline-dark">Upload video</button>
            </form>
          </div>
        </article>
      `;
    }).join('');
  }

  async function loadProducts(message = '') {
    els.products.innerHTML = `
      <article class="admin-panel admin-empty-state">
        <p class="eyebrow">Catalog</p>
        <h2>Loading products...</h2>
        <p>The dashboard is pulling the latest catalog from the admin API.</p>
      </article>
    `;

    try {
      const products = await requestJSON('/api/admin/products?limit=200');
      state.products = Array.isArray(products) ? products : [];
      renderProducts();
      if (message) setFeedback(message, 'success');
    } catch (error) {
      setFeedback(error.message, 'error');
      els.products.innerHTML = `
        <article class="admin-panel admin-empty-state">
          <p class="eyebrow">Catalog</p>
          <h2>Unable to load the dashboard</h2>
          <p>${escapeHTML(error.message || 'Please try again in a moment.')}</p>
        </article>
      `;
    }
  }

  async function handleCreateProduct(event) {
    event.preventDefault();
    setFeedback('');

    const form = els.createForm;
    const formData = new FormData(form);
    setBusy(form, true);

    try {
      await requestJSON('/api/admin/products', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          name: String(formData.get('name') || '').trim(),
          type: String(formData.get('type') || '').trim(),
          price: Number(formData.get('price') || 0),
          quantity: Number(formData.get('quantity') || 0),
          description: String(formData.get('description') || '').trim(),
          about: String(formData.get('about') || '').trim(),
        }),
      });

      form.reset();
      await loadProducts('New product created successfully.');
    } catch (error) {
      setFeedback(error.message, 'error');
    } finally {
      setBusy(els.createForm, false);
    }
  }

  async function handleUpdateProduct(form) {
    setFeedback('');
    const formData = new FormData(form);
    const productID = form.getAttribute('data-update-product');
    setBusy(form, true);

    try {
      await requestJSON(`/api/admin/products/${productID}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          name: String(formData.get('name') || '').trim(),
          type: String(formData.get('type') || '').trim(),
          price: Number(formData.get('price') || 0),
          quantity: Number(formData.get('quantity') || 0),
          description: String(formData.get('description') || '').trim(),
          about: String(formData.get('about') || '').trim(),
          is_available: Number(formData.get('is_available') || 0),
        }),
      });

      await loadProducts('Product details updated.');
    } catch (error) {
      setFeedback(error.message, 'error');
    } finally {
      setBusy(form, false);
    }
  }

  async function handleDeleteProduct(button) {
    const productID = button.getAttribute('data-delete-product');
    const card = button.closest('[data-product-id]');
    const productName = card?.querySelector('h3')?.textContent || 'this product';

    if (!window.confirm(`Delete ${productName}? This also removes its images and video.`)) {
      return;
    }

    setFeedback('');
    setBusy(card, true);

    try {
      await requestJSON(`/api/admin/products/${productID}`, { method: 'DELETE' });
      await loadProducts('Product deleted successfully.');
    } catch (error) {
      setFeedback(error.message, 'error');
      setBusy(card, false);
    }
  }

  async function handleDeleteImage(button) {
    const imageID = button.getAttribute('data-delete-image');
    const item = button.closest('.admin-gallery-item');
    const isCover = item?.classList.contains('is-cover');

    const message = isCover
      ? 'Delete this cover image? Another image will need to be set as the cover.'
      : 'Delete this image?';
    if (!window.confirm(message)) return;

    setFeedback('');
    button.disabled = true;

    try {
      await requestJSON(`/api/admin/images/${imageID}`, { method: 'DELETE' });
      await loadProducts('Image deleted successfully.');
    } catch (error) {
      setFeedback(error.message, 'error');
      button.disabled = false;
    }
  }

  async function handleImageUpload(form) {
    const productID = form.getAttribute('data-image-upload');
    const input = form.querySelector('input[name="images"]');
    const files = Array.from(input?.files || []);

    if (!files.length) {
      setFeedback('Choose at least one image to upload.', 'error');
      return;
    }

    setFeedback('');
    setBusy(form, true);

    try {
      for (const file of files) {
        const data = new FormData();
        data.append('image', file);
        await requestJSON(`/api/admin/products/${productID}/images`, {
          method: 'POST',
          body: data,
        });
      }

      form.reset();
      await loadProducts(files.length === 1 ? 'Image uploaded successfully.' : `${files.length} images uploaded successfully.`);
    } catch (error) {
      setFeedback(error.message, 'error');
      setBusy(form, false);
    }
  }

  async function handleVideoUpload(form) {
    const productID = form.getAttribute('data-video-upload');
    const input = form.querySelector('input[name="video"]');
    const file = input?.files?.[0];

    if (!file) {
      setFeedback('Choose a video file before uploading.', 'error');
      return;
    }

    setFeedback('');
    setBusy(form, true);

    try {
      const data = new FormData();
      data.append('video', file);
      await requestJSON(`/api/admin/products/${productID}/video`, {
        method: 'POST',
        body: data,
      });

      form.reset();
      await loadProducts('Video uploaded successfully.');
    } catch (error) {
      setFeedback(error.message, 'error');
      setBusy(form, false);
    }
  }

  els.createForm.addEventListener('submit', handleCreateProduct);
  els.search.addEventListener('input', (event) => {
    state.search = event.currentTarget.value || '';
    renderProducts();
  });

  els.products.addEventListener('submit', async (event) => {
    const updateForm = event.target.closest('[data-update-product]');
    const imageForm = event.target.closest('[data-image-upload]');
    const videoForm = event.target.closest('[data-video-upload]');

    if (!updateForm && !imageForm && !videoForm) return;

    event.preventDefault();

    if (updateForm) {
      await handleUpdateProduct(updateForm);
      return;
    }
    if (imageForm) {
      await handleImageUpload(imageForm);
      return;
    }
    if (videoForm) {
      await handleVideoUpload(videoForm);
    }
  });

  els.products.addEventListener('click', async (event) => {
    const deleteImageButton = event.target.closest('[data-delete-image]');
    if (deleteImageButton) {
      await handleDeleteImage(deleteImageButton);
      return;
    }

    const deleteButton = event.target.closest('[data-delete-product]');
    if (!deleteButton) return;
    await handleDeleteProduct(deleteButton);
  });

  els.products.addEventListener('input', (event) => {
    const nameInput = event.target.closest('[data-name-input]');
    if (!nameInput) return;
    const form = nameInput.closest('[data-update-product]');
    const slugInput = form?.querySelector('[data-slug-preview]');
    if (slugInput) {
      slugInput.value = slugify(nameInput.value);
    }
  });

  loadProducts();
})();
