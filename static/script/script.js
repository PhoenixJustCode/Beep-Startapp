
const API_URL = window.location.origin + '/api/v1';
let selectedMasterId = null;
let selectedServiceId = null;

// Load categories on page load
window.onload = function() {
    loadCategories();
    loadMasters();
    loadCars();
    setDefaultDate();
};

function setDefaultDate() {
    const tomorrow = new Date();
    tomorrow.setDate(tomorrow.getDate() + 1);
    document.getElementById('appointment-date').value = tomorrow.toISOString().split('T')[0];
}

async function loadCategories() {
    try {
const response = await fetch(`${API_URL}/categories`);
const data = await response.json();
const select = document.getElementById('category');
select.innerHTML = '<option value="">Select a category...</option>';
data.forEach(cat => {
    const option = document.createElement('option');
    option.value = cat.id;
    option.textContent = cat.name;
    select.appendChild(option);
});
    } catch (error) {
console.error('Error loading categories:', error);
    }
}

async function loadServices() {
    const categoryId = document.getElementById('category').value;
    if (!categoryId) return;

    try {
const response = await fetch(`${API_URL}/services?category_id=${categoryId}`);
const data = await response.json();
const select = document.getElementById('service');
select.innerHTML = '<option value="">Select a service...</option>';
data.forEach(service => {
    const option = document.createElement('option');
    option.value = service.id;
    option.textContent = service.name;
    select.appendChild(option);
});
    } catch (error) {
console.error('Error loading services:', error);
    }
}

async function loadCars() {
    selectedServiceId = document.getElementById('service').value;
    
    try {
const response = await fetch(`${API_URL}/cars`);
const data = await response.json();
const select = document.getElementById('car');
select.innerHTML = '<option value="">Select a car...</option>';
data.forEach(car => {
    const option = document.createElement('option');
    option.value = car.id;
    option.textContent = `${car.brand} ${car.model}`;
    select.appendChild(option);
});
    } catch (error) {
console.error('Error loading cars:', error);
    }
}

async function loadMasters() {
    try {
const response = await fetch(`${API_URL}/masters`);
const data = await response.json();
const container = document.getElementById('masters-container');

if (data.length === 0) {
    container.innerHTML = '<div class="loading">No masters available</div>';
    return;
}

container.innerHTML = '';
data.forEach(master => {
    const card = document.createElement('div');
    card.className = 'master-card';
    card.onclick = () => selectMaster(master.id);
    card.innerHTML = `
<div class="master-avatar">${master.name.charAt(0)}</div>
<h3>${master.name}</h3>
<div class="rating">⭐ ${master.rating || 'N/A'}</div>
<p>${master.specialization || 'Expert'}</p>
    `;
    container.appendChild(card);
});
    } catch (error) {
console.error('Error loading masters:', error);
document.getElementById('masters-container').innerHTML = '<div class="error">Error loading masters</div>';
    }
}

function selectMaster(masterId) {
    selectedMasterId = masterId;
    document.querySelectorAll('.master-card').forEach(card => {
card.classList.remove('selected');
    });
    event.currentTarget.classList.add('selected');
    updateTimeSlots();
}

async function updateTimeSlots() {
    if (!selectedMasterId) return;
    
    const date = document.getElementById('appointment-date').value;
    if (!date) return;

    try {
const response = await fetch(`${API_URL}/masters/${selectedMasterId}/available-slots?date=${date}`);
const data = await response.json();
const select = document.getElementById('appointment-time');
select.innerHTML = '<option value="">Select time slot...</option>';

if (data.slots && data.slots.length > 0) {
    data.slots.forEach(slot => {
const option = document.createElement('option');
option.value = slot;
option.textContent = slot;
select.appendChild(option);
    });
} else {
    const option = document.createElement('option');
    option.textContent = 'No slots available';
    select.appendChild(option);
}
    } catch (error) {
console.error('Error loading time slots:', error);
    }
}

document.getElementById('appointment-date').addEventListener('change', updateTimeSlots);

async function calculatePrice() {
    const serviceId = document.getElementById('service').value;
    const carId = document.getElementById('car').value;

    if (!serviceId || !carId) {
showResults('Please select both service and car', 'error');
return;
    }

    try {
const response = await fetch(`${API_URL}/pricing/calculate`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ service_id: parseInt(serviceId), car_id: parseInt(carId) })
});

const data = await response.json();
showResults(`Estimated Price: $${data.price}`, 'success');
    } catch (error) {
showResults('Error calculating price', 'error');
console.error(error);
    }
}

async function bookAppointment() {
    if (!selectedMasterId || !selectedServiceId) {
showResults('Please select a master and service', 'error');
return;
    }

    const date = document.getElementById('appointment-date').value;
    const time = document.getElementById('appointment-time').value;
    const comment = document.getElementById('comment').value;

    if (!date || !time) {
showResults('Please select date and time', 'error');
return;
    }

    try {
const response = await fetch(`${API_URL}/appointments`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
master_id: selectedMasterId,
service_id: parseInt(selectedServiceId),
date: date,
time: time,
comment: comment
    })
});

const data = await response.json();
showResults(`Appointment booked successfully! ID: ${data.id}`, 'success');
loadMyAppointments();
    } catch (error) {
showResults('Error booking appointment', 'error');
console.error(error);
    }
}

async function loadMyAppointments() {
    try {
const response = await fetch(`${API_URL}/appointments`);
const data = await response.json();
const container = document.getElementById('appointments-list');

if (data.length === 0) {
    container.innerHTML = '<div class="loading">No appointments found</div>';
    return;
}

container.innerHTML = '';
data.forEach(apt => {
    const div = document.createElement('div');
    div.className = 'result-item';
    div.innerHTML = `
<h3>Appointment #${apt.id}</h3>
<p><strong>Date:</strong> ${apt.date} at ${apt.time}</p>
<p><strong>Status:</strong> ${apt.status}</p>
<p><strong>Comment:</strong> ${apt.comment || 'None'}</p>
    `;
    container.appendChild(div);
});
    } catch (error) {
console.error('Error loading appointments:', error);
    }
}

function showResults(message, type) {
    const results = document.getElementById('results');
    const content = document.getElementById('results-content');
    
    results.style.display = 'block';
    content.innerHTML = `<div class="${type}">${message}</div>`;
}

// Login Modal Functions
function openLoginModal() {
    document.getElementById('loginModal').style.display = 'block';
}

function closeLoginModal() {
    document.getElementById('loginModal').style.display = 'none';
}

window.onclick = function(event) {
    const modal = document.getElementById('loginModal');
    if (event.target == modal) {
closeLoginModal();
    }
}

function showRegisterForm() {
    document.getElementById('registerForm').style.display = 'block';
    document.getElementById('loginForm').style.display = 'none';
    document.getElementById('dividerOr').style.display = 'none';
    document.getElementById('btnCreateAccount').style.display = 'none';
}

function hideRegisterForm() {
    document.getElementById('registerForm').style.display = 'none';
    document.getElementById('loginForm').style.display = 'block';
    document.getElementById('dividerOr').style.display = 'flex';
    document.getElementById('btnCreateAccount').style.display = 'block';
}

async function handleLogin(event) {
    event.preventDefault();
    const email = document.getElementById('modalEmail').value;
    const password = document.getElementById('modalPassword').value;
    
    try {
const response = await fetch(`${API_URL}/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password })
});

const data = await response.json();

if (response.ok) {
    localStorage.setItem('token', data.token);
    localStorage.setItem('user', JSON.stringify(data.user));
    closeLoginModal();
    window.location.href = '/profile';
} else {
    document.getElementById('loginError').textContent = data.error || 'Ошибка входа';
    document.getElementById('loginError').style.display = 'block';
}
    } catch (error) {
document.getElementById('loginError').textContent = 'Ошибка подключения';
document.getElementById('loginError').style.display = 'block';
    }
}

async function handleRegister(event) {
    event.preventDefault();
    const name = document.getElementById('regName').value;
    const email = document.getElementById('regEmail').value;
    const phone = document.getElementById('regPhone').value;
    const password = document.getElementById('regPassword').value;
    
    try {
const response = await fetch(`${API_URL}/auth/register`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ name, email, phone, password })
});

const data = await response.json();

if (response.ok) {
    localStorage.setItem('token', data.token);
    localStorage.setItem('user', JSON.stringify(data.user));
    closeLoginModal();
    window.location.href = '/profile';
} else {
    document.getElementById('loginError').textContent = data.error || 'Ошибка регистрации';
    document.getElementById('loginError').style.display = 'block';
}
    } catch (error) {
document.getElementById('loginError').textContent = 'Ошибка подключения';
document.getElementById('loginError').style.display = 'block';
    }
}

// Check if user is logged in on page load
const user = localStorage.getItem('user');
if (user) {
    const userObj = JSON.parse(user);
    const header = document.querySelector('header > div');
    header.innerHTML = `
<div>
    <h1>✨ BEEP</h1>
    <p class="tagline">Auto service Booking System</p>
</div>
<div class="profile-menu">
    <div class="profile-dropdown">
<button class="avatar" onclick="toggleProfile()">${userObj.name.charAt(0)}</button>
<div id="profileDropdown" class="dropdown-content">
    <a href="#" onclick="viewProfile()">Личный кабинет</a>
    <a href="#" onclick="logout()">Выйти</a>
</div>
    </div>
</div>
    `;
}

function toggleProfile() {
    const dropdown = document.getElementById('profileDropdown');
    dropdown.style.display = dropdown.style.display === 'block' ? 'none' : 'block';
}

function logout() {
    localStorage.removeItem('user');
    localStorage.removeItem('token');
    location.reload();
}

function viewProfile() {
    window.location.href = '/profile';
}
    