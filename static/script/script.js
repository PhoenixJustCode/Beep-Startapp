
const API_URL = window.location.origin + '/api/v1';
let selectedMasterId = null;
let selectedServiceId = null;
let allMasters = []; // Store all masters for filtering

// Utility function to log duplicates
function logDuplicates(data, type, uniqueData) {
    if (data.length !== uniqueData.length) {
        console.warn(`⚠️ Found ${data.length - uniqueData.length} duplicate ${type}:`, 
            data.filter((item, index) => 
                index !== data.findIndex(d => 
                    type === 'categories' ? d.name === item.name :
                    type === 'services' ? d.name === item.name :
                    type === 'cars' ? (d.brand === item.brand && d.model === item.model && d.year === item.year) :
                    type === 'masters' ? (d.name === item.name && d.email === item.email) : false
                )
            )
        );
    }
}

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
        
        // Remove duplicates by name
        const uniqueCategories = data.filter((category, index, self) => 
            index === self.findIndex(c => c.name === category.name)
        );
        
        logDuplicates(data, 'categories', uniqueCategories);
        console.log(`Loaded ${data.length} categories, ${uniqueCategories.length} unique`);
        
        const select = document.getElementById('category');
        select.innerHTML = '<option value="">Выберите категорию...</option>';
        
        uniqueCategories.forEach(cat => {
            const option = document.createElement('option');
            option.value = cat.id;
            option.textContent = cat.name;
            select.appendChild(option);
        });
    } catch (error) {
        console.error('Ошибка загрузки категорий:', error);
    }
}

async function loadServices() {
    const categoryId = document.getElementById('category').value;
    if (!categoryId) return;

    try {
        const response = await fetch(`${API_URL}/services?category_id=${categoryId}`);
        const data = await response.json();
        
        // Remove duplicates by name
        const uniqueServices = data.filter((service, index, self) => 
            index === self.findIndex(s => s.name === service.name)
        );
        
        logDuplicates(data, 'services', uniqueServices);
        console.log(`Loaded ${data.length} services, ${uniqueServices.length} unique`);
        
        const select = document.getElementById('service');
        select.innerHTML = '<option value="">Выберите услугу...</option>';
        
        uniqueServices.forEach(service => {
            const option = document.createElement('option');
            option.value = service.id;
            option.textContent = service.name;
            select.appendChild(option);
        });
    } catch (error) {
        console.error('Ошибка загрузки услуг:', error);
    }
}

async function loadCars() {
    selectedServiceId = document.getElementById('service').value;
    
    try {
        const response = await fetch(`${API_URL}/cars`);
        const data = await response.json();
        
        // Remove duplicates by brand + model + year
        const uniqueCars = data.filter((car, index, self) => 
            index === self.findIndex(c => 
                c.brand === car.brand && 
                c.model === car.model && 
                c.year === car.year
            )
        );
        
        logDuplicates(data, 'cars', uniqueCars);
        console.log(`Loaded ${data.length} cars, ${uniqueCars.length} unique`);
        
        const select = document.getElementById('car');
        select.innerHTML = '<option value="">Выберите автомобиль...</option>';
        
        uniqueCars.forEach(car => {
            const option = document.createElement('option');
            option.value = car.id;
            option.textContent = `${car.brand} ${car.model}`;
            select.appendChild(option);
        });
    } catch (error) {
        console.error('Ошибка загрузки автомобилей:', error);
    }
}

async function loadMasters() {
    try {
        const response = await fetch(`${API_URL}/masters`);
        const data = await response.json();
        
        // Remove duplicates by name + email
        const uniqueMasters = data.filter((master, index, self) => 
            index === self.findIndex(m => 
                m.name === master.name && 
                m.email === master.email
            )
        );
        
        allMasters = uniqueMasters; // Store unique masters
        logDuplicates(data, 'masters', uniqueMasters);
        console.log(`Loaded ${data.length} masters, ${uniqueMasters.length} unique`);
        
        const container = document.getElementById('masters-container');

        if (uniqueMasters.length === 0) {
            container.innerHTML = '<div class="loading">Мастера недоступны</div>';
            return;
        }

        renderMasters(uniqueMasters);
    } catch (error) {
        console.error('Ошибка загрузки мастеров:', error);
        document.getElementById('masters-container').innerHTML = '<div class="error">Ошибка загрузки мастеров</div>';
    }
}

function renderMasters(masters) {
    const container = document.getElementById('masters-container');
    container.innerHTML = '';

    masters.forEach(master => {
        const card = document.createElement('div');
        card.className = 'master-card';
        card.onclick = (e) => {
            if (e.target.classList.contains('favorite-star') || e.target.parentElement.classList.contains('favorite-star')) {
                return; // Don't select master when clicking star
            }
            selectMaster(master.id);
        };
        
        const isVerified = master.is_verified ? '<span style="color: #10b981; font-size: 14px; font-weight: 600;">✓ Проверенный</span>' : '';
        const favoriteClass = master.is_favorite ? 'favorite-active' : '';
        const favoriteIcon = master.is_favorite ? '⭐' : '☆';
        
        card.innerHTML = `
            <div style="position: relative;">
                <div class="master-avatar">${master.name.charAt(0)}</div>
                <span class="favorite-star ${favoriteClass}" onclick="toggleFavorite(${master.id}, event)" style="position: absolute; top: 5px; right: 5px; font-size: 16px; cursor: pointer; z-index: 10; padding: 2px; background: rgba(255,255,255,0.8); border-radius: 50%; width: 24px; height: 24px; display: flex; align-items: center; justify-content: center; line-height: 1;" title="${master.is_favorite ? 'Удалить из избранного' : 'Добавить в избранное'}">${favoriteIcon}</span>
            </div>
            <h3>${master.name}</h3>
            ${isVerified}
            <div class="rating">⭐ ${master.rating || 'N/A'}</div>
            <p class="specialization">${master.specialization || 'Expert'}</p>
            <div class="master-location">📍 ${master.address || 'Location not specified'}</div>
        `;
        container.appendChild(card);
    });
}

// Toggle favorite master
async function toggleFavorite(masterId, event) {
    event.stopPropagation();
    
    try {
        const token = localStorage.getItem('token');
        if (!token) {
            showResults('Необходима авторизация для добавления в избранное', 'error');
            return;
        }
        
        // Find master in allMasters
        const master = allMasters.find(m => m.id === masterId);
        if (!master) return;
        
        const isFavorite = master.is_favorite;
        const url = isFavorite ? `${API_URL}/favorites/${masterId}` : `${API_URL}/favorites`;
        const method = isFavorite ? 'DELETE' : 'POST';
        
        const requestBody = isFavorite ? undefined : JSON.stringify({ master_id: masterId });
        
        // Ensure token format is correct for the backend
        // Backend expects token with "mock-jwt-token-" prefix
        // If token is just email, add prefix; if already has prefix, use as is
        let authToken = token;
        if (token && !token.includes('mock-jwt-token-') && token.includes('@')) {
            // Token is email, add prefix
            authToken = `mock-jwt-token-${token}`;
        }
        // Remove "Bearer " prefix if present, backend will add it
        authToken = authToken.replace(/^Bearer\s+/i, '');
        
        const response = await fetch(url, {
            method,
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${authToken}`
            },
            body: requestBody
        });
        
        if (!response.ok) {
            let errorMessage = 'Не удалось обновить избранное';
            try {
                const errorData = await response.json();
                errorMessage = errorData.error || errorMessage;
                console.error('API Error:', errorData);
                
                // If user not found, suggest to re-login or register
                if (response.status === 401) {
                    if (errorData.error && (errorData.error.includes('не найден') || errorData.error.includes('not found'))) {
                        errorMessage = 'Пользователь не найден. Пожалуйста, войдите в систему или зарегистрируйтесь.';
                        // Clear token and redirect to login
                        localStorage.removeItem('token');
                        localStorage.removeItem('user');
                        
                        // Show login modal after a short delay
                        setTimeout(() => {
                            if (typeof showLoginModal === 'function') {
                                showLoginModal();
                            } else {
                                // If modal function not available, reload page
                                window.location.reload();
                            }
                        }, 1000);
                    } else {
                        errorMessage = 'Необходима авторизация. Пожалуйста, войдите в систему.';
                        localStorage.removeItem('token');
                        localStorage.removeItem('user');
                    }
                }
            } catch (e) {
                console.error('Response status:', response.status, response.statusText);
            }
            throw new Error(errorMessage);
        }
        
        // Update local state immediately for UI responsiveness
        master.is_favorite = !isFavorite;
        
        console.log(`Favorite ${!isFavorite ? 'added' : 'removed'} for master ${masterId}, new state: ${master.is_favorite}`);
        
        // Re-render masters with updated state immediately
        const searchTerm = document.getElementById('master-search').value.toLowerCase();
        if (searchTerm) {
            filterMasters();
        } else {
            renderMasters(allMasters);
        }
        
        showResults(isFavorite ? 'Мастер удален из избранного' : 'Мастер добавлен в избранное', 'success');
        
        // Don't reload immediately - let user see the change
        // The state will be synced on next page load or manual refresh
        // This prevents the issue where favorite gets added then immediately removed
    } catch (error) {
        console.error('Ошибка обновления избранного:', error);
        showResults(`Ошибка: ${error.message}`, 'error');
    }
}

// Filter by favorites
let showOnlyFavorites = false;
function toggleFavoritesFilter() {
    showOnlyFavorites = !showOnlyFavorites;
    const button = document.getElementById('favoritesFilterBtn');
    if (button) {
        button.style.background = showOnlyFavorites ? 'var(--primary)' : '#6b7280';
    }
    
    if (showOnlyFavorites) {
        const favoriteMasters = allMasters.filter(m => m.is_favorite);
        renderMasters(favoriteMasters);
    } else {
        renderMasters(allMasters);
    }
}

function filterMasters() {
    const searchTerm = document.getElementById('master-search').value.toLowerCase();
    let filteredMasters = allMasters.filter(master => 
        master.name.toLowerCase().includes(searchTerm) ||
        master.specialization.toLowerCase().includes(searchTerm) ||
        (master.address && master.address.toLowerCase().includes(searchTerm))
    );
    
    // Apply favorites filter if active
    if (showOnlyFavorites) {
        filteredMasters = filteredMasters.filter(m => m.is_favorite);
    }
    
    console.log(`Filtered ${allMasters.length} masters to ${filteredMasters.length} results`);
    renderMasters(filteredMasters);
}

function toggleMasterDropdown() {
    const dropdown = document.getElementById('master-dropdown-body');
    const header = document.querySelector('.dropdown-header');
    const isOpen = dropdown.style.display === 'block';
    
    if (isOpen) {
        dropdown.style.display = 'none';
        header.classList.remove('open');
    } else {
        dropdown.style.display = 'block';
        header.classList.add('open');
        // Focus on search input when opening
        setTimeout(() => {
            document.getElementById('master-search').focus();
        }, 100);
    }
}

// Close dropdown when clicking outside
document.addEventListener('click', function(event) {
    const dropdown = document.getElementById('master-dropdown');
    if (!dropdown.contains(event.target)) {
        document.getElementById('master-dropdown-body').style.display = 'none';
    }
});

function selectMaster(masterId) {
    selectedMasterId = masterId;
    const master = allMasters.find(m => m.id === masterId);
    
    document.querySelectorAll('.master-card').forEach(card => {
        card.classList.remove('selected');
    });
    event.currentTarget.classList.add('selected');
    
    // Update selected master display
    document.getElementById('selected-master').innerHTML = `
        <div style="display: flex; align-items: center; gap: 10px;">
            <div class="master-avatar-small">${master.name.charAt(0)}</div>
            <div>
                <div style="font-weight: 600; color: var(--primary);">${master.name}</div>
                <div style="font-size: 12px; color: #666;">
                    ⭐ ${master.rating || 'N/A'} • ${master.specialization || 'Expert'}
                </div>
            </div>
        </div>
    `;
    
    // Close dropdown
    document.getElementById('master-dropdown-body').style.display = 'none';
    
    // Clear search
    document.getElementById('master-search').value = '';
    
    updateTimeSlots();
}

async function updateTimeSlots() {
    if (!selectedMasterId) return;
    
    const date = document.getElementById('appointment-date').value;
    if (!date) return;

    const select = document.getElementById('appointment-time');
    select.innerHTML = '<option value="">Загрузка...</option>';
    select.disabled = true;

    try {
        const response = await fetch(`${API_URL}/masters/${selectedMasterId}/available-slots?date=${date}`);
        
        if (!response.ok) {
            select.innerHTML = '<option value="">Ошибка загрузки времени</option>';
            select.disabled = false;
            return;
        }
        
        const data = await response.json();
        select.innerHTML = '<option value="">Выберите время...</option>';

        if (data.slots && data.slots.length > 0) {
            data.slots.forEach(slot => {
                const option = document.createElement('option');
                option.value = slot;
                option.textContent = slot;
                select.appendChild(option);
            });
        } else {
            const option = document.createElement('option');
            option.value = '';
            option.textContent = 'Нет доступного времени';
            option.disabled = true;
            select.appendChild(option);
        }
        select.disabled = false;
    } catch (error) {
        select.innerHTML = '<option value="">Ошибка загрузки времени</option>';
        select.disabled = false;
    }
}

document.getElementById('appointment-date').addEventListener('change', updateTimeSlots);

async function calculatePrice() {
    const serviceId = document.getElementById('service').value;
    const carId = document.getElementById('car').value;

    if (!serviceId || !carId) {
        showResults('Пожалуйста, выберите услугу и автомобиль', 'error');
        return;
    }

    try {
        const response = await fetch(`${API_URL}/pricing/calculate`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ service_id: parseInt(serviceId), car_id: parseInt(carId) })
        });

        if (!response.ok) {
            throw new Error('Не удалось рассчитать цену');
        }

        const data = await response.json();
        displayPriceCalculation(data);
    } catch (error) {
        showResults('Ошибка расчета цены: ' + error.message, 'error');
        console.error(error);
    }
}

function displayPriceCalculation(data) {
    const results = document.getElementById('results');
    const content = document.getElementById('results-content');
    
    results.style.display = 'block';
    
    let html = `
        <div class="price-calculation">
            <div class="price-header">
                <h3>💰 Расчет стоимости</h3>
                <div class="final-price">${Math.round(data.final_price)} KZT</div>
            </div>
            
            <div class="calculation-details">
                <div class="service-info">
                    <h4>${data.service_name}</h4>
                    <div class="car-info">
                        <span class="car-brand">${data.car_brand} ${data.car_model}</span>
                        <span class="car-year">(${data.car_year} год, ${data.car_age} лет)</span>
                        <span class="car-type">${data.car_type}</span>
                    </div>
                </div>
                
                <div class="price-breakdown">
                    <h4>Детализация расчета:</h4>
                    <div class="price-items">
    `;
    
    data.price_details.forEach(detail => {
        const amountClass = detail.is_addition ? 
            (detail.amount > 0 ? 'price-addition' : 'price-discount') : 
            'price-base';
        
        html += `
            <div class="price-item ${amountClass}">
                <span class="price-description">${detail.description}</span>
                <span class="price-amount">
                    ${detail.amount > 0 ? '+' : ''}${Math.round(detail.amount)} KZT
                    ${detail.multiplier ? `(${detail.multiplier}x)` : ''}
                </span>
            </div>
        `;
    });
    
    html += `
                    </div>
                    
                    <div class="price-summary">
                        <div class="price-item price-total">
                            <span class="price-description">Итого:</span>
                            <span class="price-amount">${Math.round(data.final_price)} KZT</span>
                        </div>
                        
                        <div class="price-range">
                            <small>Диапазон: ${Math.round(data.min_price)} - ${Math.round(data.max_price)} KZT</small>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    `;
    
    content.innerHTML = html;
}

async function bookAppointment() {
    if (!selectedMasterId || !selectedServiceId) {
showResults('Пожалуйста, выберите мастера и услугу', 'error');
return;
    }

    const date = document.getElementById('appointment-date').value;
    const time = document.getElementById('appointment-time').value;
    const comment = document.getElementById('comment').value;

    if (!date || !time) {
showResults('Пожалуйста, выберите дату и время', 'error');
return;
    }

    try {
        const token = localStorage.getItem('token');
        const headers = { 'Content-Type': 'application/json' };
        if (token) {
            headers['Authorization'] = 'Bearer ' + token;
        }
        
        const response = await fetch(`${API_URL}/appointments`, {
    method: 'POST',
    headers: headers,
    body: JSON.stringify({
master_id: selectedMasterId,
service_id: parseInt(selectedServiceId),
date: date,
time: time,
comment: comment
    })
});

const data = await response.json();
if (response.ok) {
    showResults(`Запись успешно создана! ID: ${data.id}`, 'success');
} else {
    showResults(data.error || 'Ошибка создания записи', 'error');
}
loadMyAppointments();
    } catch (error) {
showResults('Ошибка создания записи', 'error');
console.error(error);
    }
}

async function loadMyAppointments() {
    try {
        const token = localStorage.getItem('token');
        const headers = { 'Content-Type': 'application/json' };
        if (token) {
            headers['Authorization'] = 'Bearer ' + token;
        }
        
        const response = await fetch(`${API_URL}/appointments`, { headers });
        const data = await response.json();
        const container = document.getElementById('appointments-list');

        if (data.length === 0) {
            container.innerHTML = '<div class="loading">У вас пока нет заказов</div>';
            return;
        }

        container.innerHTML = '';
        data.forEach(apt => {
            const div = document.createElement('div');
            div.className = 'appointment-card';
            
            // Format date beautifully
            let dateStr = 'Не указана';
            let timeStr = 'Не указано';
            
            if (apt.date) {
                try {
                    const date = new Date(apt.date);
                    const months = ['Января', 'Февраля', 'Марта', 'Апреля', 'Мая', 'Июня', 
                                  'Июля', 'Августа', 'Сентября', 'Октября', 'Ноября', 'Декабря'];
                    dateStr = `${date.getDate()} ${months[date.getMonth()]} ${date.getFullYear()}`;
                } catch (e) {
                    dateStr = apt.date;
                }
            }
            
            if (apt.time) {
                try {
                    const time = new Date(apt.time);
                    timeStr = `${time.getHours().toString().padStart(2, '0')}:${time.getMinutes().toString().padStart(2, '0')}`;
                } catch (e) {
                    timeStr = apt.time;
                }
            }
            
            div.innerHTML = `
                <div class="appointment-header">
                    <h3>Заказ #${apt.id}</h3>
                    <span class="status status-${apt.status || 'pending'}">${(apt.status || 'pending')}</span>
                </div>
                <div class="appointment-body">
                    <p><strong>📅 Дата:</strong> ${dateStr}</p>
                    <p><strong>🕐 Время:</strong> ${timeStr}</p>
                    <p><strong>🔧 Услуга:</strong> ${apt.service_name || 'Не указана'}</p>
                    <p><strong>👨‍🔧 Мастер:</strong> ${apt.master_name || 'Не назначен'}</p>
                    ${apt.comment ? `<p><strong>💬 Комментарий:</strong> ${apt.comment}</p>` : ''}
                </div>
            `;
            container.appendChild(div);
        });
    } catch (error) {
        console.error('Ошибка загрузки записей:', error);
        document.getElementById('appointments-list').innerHTML = '<p style="color: #e74c3c;">Ошибка загрузки записей</p>';
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
<div style="flex: 1; display: flex; align-items: center; gap: 15px;">
    <a href="/"><img src="/static/logo.png?v=3&t=1730123456" alt="BEEP" style="width: 60px; height: 60px;"></a>
    <div>
        <h1 style="margin: 0;">BEEP</h1>
        <p class="tagline" style="margin: 0;">Система записи в автосервис</p>
    </div>
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
    