const API_URL = 'http://localhost:8080';
const content = document.getElementById('content');

function getToken() {
  return localStorage.getItem('token');
}

function setToken(token) {
  localStorage.setItem('token', token);
  updateAuthLinks();
}

function removeToken() {
  localStorage.removeItem('token');
  updateAuthLinks();
}

function updateAuthLinks() {
  if (getToken()) {
    document.getElementById('authLinks').style.display = 'none';
    document.getElementById('logoutLink').style.display = 'inline';
  } else {
    document.getElementById('authLinks').style.display = 'inline';
    document.getElementById('logoutLink').style.display = 'none';
  }
}

async function apiFetch(path, options = {}) {
  const headers = options.headers || {};
  headers['Content-Type'] = 'application/json';
  if (getToken()) {
    headers['Authorization'] = 'Bearer ' + getToken();
  }
  const res = await fetch(API_URL + path, {...options, headers});
  if (!res.ok) throw new Error(await res.text());
  return res.json();
}

function showHome() {
  content.innerHTML = 'Loading products...';
  fetch(API_URL + '/products')
    .then(res => res.json())
    .then(products => {
      content.innerHTML = '<h2>Products</h2>';
      products.forEach(p => {
        const card = document.createElement('div');
        card.className = 'card';
        card.innerHTML = `
          <h3>${p.name}</h3>
          <img src="/uploads/${p.file_name}" style="max-width:100%;">
          <p>${p.description}</p>
          <p>Price: $${p.price}</p>
          <button onclick="showProduct(${p.id})">View</button>
        `;
        content.appendChild(card);
      });
    });
}

function showProduct(id) {
  content.innerHTML = 'Loading product...';
  fetch(API_URL + '/products/' + id)
    .then(res => res.json())
    .then(p => {
      content.innerHTML = `
        <div class="card">
          <h2>${p.name}</h2>
          <img src="/uploads/${p.file_name}" style="max-width:100%;">
          <p>${p.description}</p>
          <p>Price: $${p.price}</p>
          <button onclick="addToCart(${p.id})">Add to Cart</button>
          <button onclick="showHome()">Back</button>
        </div>
      `;
    });
}

async function addToCart(productId) {
  if (!getToken()) {
    alert('You must login first!');
    showLogin();
    return;
  }
  await apiFetch('/cart', {
    method: 'POST',
    body: JSON.stringify({ user_id: 1, product_id: productId, quantity: 1 })
  });
  alert('Added to cart!');
}

function showCart() {
  if (!getToken()) {
    alert('You must login first!');
    showLogin();
    return;
  }
  content.innerHTML = 'Loading cart...';
  apiFetch('/cart?user_id=1')
    .then(items => {
      content.innerHTML = '<h2>Your Cart</h2>';
      if (!items.length) {
        content.innerHTML += '<p>Cart is empty.</p>';
      }
      items.forEach(item => {
        const card = document.createElement('div');
        card.className = 'card';
        card.innerHTML = `
          <p>Product ID: ${item.product_id}</p>
          <p>Quantity: ${item.quantity}</p>
          <button class="danger" onclick="deleteCartItem(${item.id})">Delete</button>
        `;
        content.appendChild(card);
      });
      content.innerHTML += `<button onclick="showHome()">Back to Shop</button>`;
    });
}

async function deleteCartItem(id) {
  await apiFetch('/cart?id=' + id, { method: 'DELETE' });
  showCart();
}

function showLogin() {
  content.innerHTML = `
    <h2>Login</h2>
    <form onsubmit="login(event)">
      <input id="loginUser" placeholder="Username" required><br>
      <input id="loginPass" type="password" placeholder="Password" required><br>
      <button type="submit">Login</button>
    </form>
  `;
}

async function login(e) {
  e.preventDefault();
  const username = document.getElementById('loginUser').value;
  const password = document.getElementById('loginPass').value;
  try {
    const res = await fetch(API_URL + '/login', {
      method: 'POST',
      headers: {'Content-Type': 'application/json'},
      body: JSON.stringify({username, password})
    });
    if (!res.ok) throw new Error(await res.text());
    const data = await res.json();
    setToken(data.token);
    showHome();
  } catch (err) {
    alert('Login failed: ' + err.message);
  }
}

function showRegister() {
  content.innerHTML = `
    <h2>Register</h2>
    <form onsubmit="register(event)">
      <input id="regUser" placeholder="Username" required><br>
      <input id="regEmail" placeholder="Email" required><br>
      <input id="regPass" type="password" placeholder="Password" required><br>
      <button type="submit">Register</button>
    </form>
  `;
}

async function register(e) {
  e.preventDefault();
  const username = document.getElementById('regUser').value;
  const email = document.getElementById('regEmail').value;
  const password = document.getElementById('regPass').value;
  try {
    const res = await fetch(API_URL + '/users', {
      method: 'POST',
      headers: {'Content-Type': 'application/json'},
      body: JSON.stringify({username, email, password})
    });
    if (!res.ok) throw new Error(await res.text());
    alert('Registered! Please login.');
    showLogin();
  } catch (err) {
    alert('Register failed: ' + err.message);
  }
}

function logout() {
  removeToken();
  showHome();
}

// Init
updateAuthLinks();
showHome();
