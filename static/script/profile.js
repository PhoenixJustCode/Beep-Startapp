const API_URL = window.location.origin + '/api/v1';

// Load user data on page load
window.onload = function() {
    loadUserProfile();
    loadUserAppointments();
    loadMasterProfile(); // Add this line
};

function openTab(evt, tabName) {
    // Скрываем все вкладки
    const tabContents = document.getElementsByClassName('tab-content');
    for (let i = 0; i < tabContents.length; i++) {
        tabContents[i].style.display = 'none';
    }
    
    // Убираем активный класс у всех кнопок
    const tabBtns = document.getElementsByClassName('tab-btn');
    for (let i = 0; i < tabBtns.length; i++) {
        tabBtns[i].classList.remove('active');
        tabBtns[i].style.borderBottom = '3px solid transparent';
        tabBtns[i].style.color = '#64748b';
    }
    
    // Показываем выбранную вкладку
    document.getElementById(tabName).style.display = 'block';
    
    // Добавляем активный класс к кнопке
    evt.currentTarget.classList.add('active');
    evt.currentTarget.style.borderBottom = '3px solid var(--primary)';
    evt.currentTarget.style.color = 'var(--primary)';
}


async function loadUserProfile() {
    const user = localStorage.getItem('user');
    if (!user) {
        window.location.href = '/';
        return;
    }

    const userObj = JSON.parse(user);
    
    document.getElementById('userName').textContent = userObj.name || '-';
    document.getElementById('userEmail').textContent = userObj.email || '-';
    document.getElementById('userPhone').textContent = userObj.phone || '-';
    document.getElementById('userCreated').textContent = userObj.created_at ? new Date(userObj.created_at).toLocaleDateString('ru-RU') : '-';
    
        // Set profile initial
        const firstLetter = userObj.name ? userObj.name.charAt(0).toUpperCase() : 'A';
        document.getElementById('profileInitial').textContent = firstLetter;
        
        // Set profile photo if exists
        if (userObj.photo_url) {
            document.getElementById('profilePicture').style.background = `url(${userObj.photo_url}) center/cover`;
            document.getElementById('profilePicture').style.color = 'transparent';
            document.getElementById('profileInitial').style.display = 'none';
        } else {
            document.getElementById('profilePicture').style.background = 'linear-gradient(135deg, var(--primary), var(--secondary))';
            document.getElementById('profilePicture').style.color = 'white';
            document.getElementById('profileInitial').style.display = 'flex';
        }
        
        // Store current user data for editing
        window.currentUser = userObj;
}

async function saveProfile(event) {
    event.preventDefault();
    
    const name = document.getElementById('editName').value;
    const email = document.getElementById('editEmail').value;
    const phone = document.getElementById('editPhone').value;
    
    try {
        const token = localStorage.getItem('token');
        const response = await fetch(`${API_URL}/user/profile`, {
            method: 'PUT',
            headers: { 
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({ name, email, phone })
        });

        if (!response.ok) {
            throw new Error('Failed to update profile');
        }

        const data = await response.json();
        
        // Update localStorage
        const updatedUser = { ...window.currentUser, name, email, phone };
        localStorage.setItem('user', JSON.stringify(updatedUser));
        window.currentUser = updatedUser;
        
        // Update display
        document.getElementById('userName').textContent = name;
        document.getElementById('userEmail').textContent = email;
        document.getElementById('userPhone').textContent = phone;
        
        // Update profile initial
        const firstLetter = name ? name.charAt(0).toUpperCase() : 'A';
        document.getElementById('profileInitial').textContent = firstLetter;
        
        // Show success message
        showMessage('Профиль успешно обновлен!', 'success');
        
        // Hide edit mode
        document.getElementById('profileViewMode').style.display = 'block';
        document.getElementById('profileEditMode').style.display = 'none';
        
    } catch (error) {
        console.error('Error updating profile:', error);
        showMessage('Ошибка при обновлении профиля: ' + error.message, 'error');
    }
}

// Open edit modal for field
function openEditModal(fieldName) {
    const fieldMap = {
        'name': { label: 'Имя', type: 'text' },
        'email': { label: 'Email', type: 'email' },
        'phone': { label: 'Телефон', type: 'tel' }
    };
    
    const field = fieldMap[fieldName];
    if (!field) return;
    
    const currentValue = document.getElementById(`user${fieldName.charAt(0).toUpperCase() + fieldName.slice(1)}`).textContent;
    
    document.getElementById('editFieldName').value = fieldName;
    document.getElementById('editFieldTitle').textContent = `Редактировать ${field.label.toLowerCase()}`;
    document.getElementById('editFieldLabel').textContent = field.label + ':';
    document.getElementById('editFieldValue').type = field.type;
    document.getElementById('editFieldValue').value = currentValue;
    
    document.getElementById('editFieldModal').style.display = 'flex';
    document.getElementById('editFieldValue').focus();
    document.getElementById('editFieldValue').select();
}

// Close edit modal
function closeEditModal() {
    document.getElementById('editFieldModal').style.display = 'none';
}

// Close appointment edit modal
function closeAppointmentModal() {
    document.getElementById('editAppointmentModal').style.display = 'none';
}

// Save field edit
async function saveFieldEdit(event) {
    event.preventDefault();
    
    const fieldName = document.getElementById('editFieldName').value;
    const newValue = document.getElementById('editFieldValue').value.trim();
    
    if (!newValue) {
        showMessage('Поле не может быть пустым', 'error');
        return;
    }
    
    try {
        const token = localStorage.getItem('token');
        
        // Create update object with only the changed field
        const updateData = {};
        updateData[fieldName] = newValue;
        
        const response = await fetch(`${API_URL}/user/profile`, {
            method: 'PUT',
            headers: { 
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify(updateData)
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || 'Failed to update field');
        }

        // Update display
        const elementId = fieldName === 'name' ? 'userName' : 
                          fieldName === 'email' ? 'userEmail' : 
                          'userPhone';
        
        document.getElementById(elementId).textContent = newValue;
        
        // Update localStorage
        const user = JSON.parse(localStorage.getItem('user'));
        user[fieldName] = newValue;
        localStorage.setItem('user', JSON.stringify(user));
        window.currentUser = user;
        
        // Update profile initial if name changed
        if (fieldName === 'name') {
            const firstLetter = newValue ? newValue.charAt(0).toUpperCase() : 'A';
            document.getElementById('profileInitial').textContent = firstLetter;
        }
        
        // Close only the field edit modal
        document.getElementById('editFieldModal').style.display = 'none';
        showMessage('Успешно обновлено!', 'success');
        
    } catch (error) {
        console.error('Error updating field:', error);
        showMessage(`Ошибка при обновлении: ${error.message}`, 'error');
    }
}

function showMessage(message, type) {
    // Create message element
    const messageDiv = document.createElement('div');
    messageDiv.className = `message ${type}`;
    messageDiv.textContent = message;
    messageDiv.style.cssText = `
        position: fixed;
        top: 20px;
        right: 20px;
        padding: 15px 20px;
        border-radius: 8px;
        color: white;
        font-weight: 600;
        z-index: 1000;
        animation: slideIn 0.3s ease;
        ${type === 'success' ? 'background: #10b981;' : 'background: #ef4444;'}
    `;
    
    document.body.appendChild(messageDiv);
    
    // Remove message after 3 seconds
    setTimeout(() => {
        messageDiv.remove();
    }, 3000);
}

// Profile photo upload
async function uploadProfilePhoto(event) {
    const file = event.target.files[0];
    if (!file) return;
    
    if (!file.type.startsWith('image/')) {
        showMessage('Пожалуйста, выберите изображение', 'error');
        return;
    }
    
    // Check file size (max 5MB)
    if (file.size > 5 * 1024 * 1024) {
        showMessage('Файл слишком большой (максимум 5MB)', 'error');
        return;
    }
    
    try {
        const formData = new FormData();
        formData.append('photo', file);
        
        const token = localStorage.getItem('token');
        const response = await fetch(`${API_URL}/user/photo`, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${token}`
            },
            body: formData
        });
        
        // Check if response is ok
        if (!response.ok) {
            let errorMessage = 'Failed to upload photo';
            try {
                const errorData = await response.json();
                errorMessage = errorData.error || errorMessage;
            } catch (e) {
                // If response is not JSON, get text
                const text = await response.text();
                errorMessage = text || errorMessage;
            }
            throw new Error(errorMessage);
        }
        
        const result = await response.json();
        
        // Update profile picture display
        document.getElementById('profilePicture').style.background = `url(${result.photo_url}) center/cover`;
        document.getElementById('profilePicture').style.color = 'transparent';
        document.getElementById('profileInitial').style.display = 'none';
        
        // Update localStorage
        const user = JSON.parse(localStorage.getItem('user'));
        user.photo_url = result.photo_url;
        localStorage.setItem('user', JSON.stringify(user));
        
        showMessage('Фото успешно загружено!', 'success');
        
    } catch (error) {
        console.error('Error uploading photo:', error);
        showMessage(`Ошибка загрузки фото: ${error.message}`, 'error');
    }
}

// Edit appointment
function editAppointment(appointmentId) {
    const modal = document.getElementById('editAppointmentModal');
    
    // Close profile edit modal if open
    document.getElementById('editFieldModal').style.display = 'none';
    
    // Load appointment data
    const token = localStorage.getItem('token');
    const headers = {};
    if (token) {
        headers['Authorization'] = 'Bearer ' + token;
    }
    
    fetch(`${API_URL}/appointments/${appointmentId}`, { headers })
        .then(response => response.json())
        .then(data => {
            document.getElementById('editAppointmentId').value = data.id;
            document.getElementById('editDate').value = data.date ? data.date.split('T')[0] : '';
            document.getElementById('editTime').value = data.time || '';
            document.getElementById('editComment').value = data.comment || '';
            
            modal.style.display = 'flex';
        })
        .catch(error => {
            console.error('Error loading appointment:', error);
            showMessage('Ошибка загрузки записи', 'error');
        });
}

// Close edit modal
function closeEditModal() {
    document.getElementById('editAppointmentModal').style.display = 'none';
}

// Save appointment changes
async function cancelAppointment(appointmentId) {
    if (!confirm('Вы точно хотите отменить эту запись?')) {
        return;
    }
    
    try {
        const token = localStorage.getItem('token');
        const response = await fetch(`${API_URL}/appointments/${appointmentId}/cancel`, {
            method: 'PUT',
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        
        if (!response.ok) {
            let errorMessage = 'Failed to cancel appointment';
            try {
                const error = await response.json();
                errorMessage = error.error || errorMessage;
            } catch (e) {
                const text = await response.text();
                console.error('Cancel appointment error (non-JSON):', text);
                errorMessage = text || errorMessage;
            }
            throw new Error(errorMessage);
        }
        
        showMessage('Запись успешно отменена!', 'success');
        loadUserAppointments(); // Reload appointments
        
    } catch (error) {
        console.error('Error cancelling appointment:', error);
        showMessage(`Ошибка отмены записи: ${error.message}`, 'error');
    }
}

async function deleteAppointment(appointmentId) {
    if (!confirm('Вы точно хотите удалить эту запись?')) {
        return;
    }
    
    try {
        const token = localStorage.getItem('token');
        const response = await fetch(`${API_URL}/appointments/${appointmentId}`, {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        
        if (!response.ok) {
            let errorMessage = 'Failed to delete appointment';
            try {
                const error = await response.json();
                errorMessage = error.error || errorMessage;
            } catch (e) {
                const text = await response.text();
                console.error('Delete appointment error (non-JSON):', text);
                errorMessage = text || errorMessage;
            }
            throw new Error(errorMessage);
        }
        
        showMessage('Запись успешно удалена!', 'success');
        loadUserAppointments(); // Reload appointments
        
    } catch (error) {
        console.error('Error deleting appointment:', error);
        showMessage(`Ошибка удаления записи: ${error.message}`, 'error');
    }
}

async function saveAppointmentChanges(event) {
    event.preventDefault();
    
    const appointmentId = document.getElementById('editAppointmentId').value;
    const date = document.getElementById('editDate').value;
    const time = document.getElementById('editTime').value;
    const comment = document.getElementById('editComment').value;
    
    try {
        const token = localStorage.getItem('token');
        const response = await fetch(`${API_URL}/appointments/${appointmentId}`, {
            method: 'PUT',
            headers: { 
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({ date, time, comment })
        });

        if (!response.ok) {
            throw new Error('Failed to update appointment');
        }

        showMessage('Запись успешно обновлена!', 'success');
        closeAppointmentModal();
        loadUserAppointments(); // Reload appointments
        
    } catch (error) {
        console.error('Error updating appointment:', error);
        showMessage('Ошибка при обновлении записи: ' + error.message, 'error');
    }
}

async function loadUserAppointments() {
    try {
        const token = localStorage.getItem('token');
        const headers = {};
        if (token) {
            headers['Authorization'] = 'Bearer ' + token;
        }
        
        const response = await fetch(`${API_URL}/appointments`, { headers });
        const data = await response.json();
        const container = document.getElementById('userAppointments');

        if (data.length === 0) {
            container.innerHTML = '<p style="color: #999;">У вас пока нет заказов</p>';
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
                <div class="appointment-actions" style="margin-top: 15px; display: flex; gap: 10px;">
                    <button onclick="editAppointment(${apt.id})" style="flex: 1; padding: 10px; background: var(--primary); color: white; border: none; border-radius: 8px; cursor: pointer; font-weight: 600;">✏️ Редактировать</button>
                    ${apt.status === 'cancelled' ? 
                        `<button onclick="deleteAppointment(${apt.id})" style="flex: 1; padding: 10px; background: #e74c3c; color: white; border: none; border-radius: 8px; cursor: pointer; font-weight: 600;">🗑️ Удалить</button>` :
                        `<button onclick="cancelAppointment(${apt.id})" style="flex: 1; padding: 10px; background: #f39c12; color: white; border: none; border-radius: 8px; cursor: pointer; font-weight: 600;">❌ Отменить</button>`
                    }
                </div>
            `;
            container.appendChild(div);
        });
    } catch (error) {
        console.error('Error loading appointments:', error);
        document.getElementById('userAppointments').innerHTML = '<p style="color: #e74c3c;">Ошибка загрузки заказов</p>';
    }
}

// Master Profile Functions

// Load master profile
async function loadMasterProfile() {
    try {
        const token = localStorage.getItem('token');
        const headers = {};
        if (token) {
            headers['Authorization'] = 'Bearer ' + token;
        }
        
        const response = await fetch(`${API_URL}/master/profile`, { headers });
        if (response.status === 404) {
            // No master profile exists
            document.getElementById('noMasterProfile').style.display = 'block';
            document.getElementById('createMasterForm').style.display = 'none';
            document.getElementById('existingMasterProfile').style.display = 'none';
            return;
        }
        
        if (!response.ok) {
            throw new Error('Failed to load master profile');
        }
        
        const master = await response.json();
        
        // Show existing master profile
        document.getElementById('noMasterProfile').style.display = 'none';
        document.getElementById('createMasterForm').style.display = 'none';
        document.getElementById('existingMasterProfile').style.display = 'block';
        
        // Update display
        document.getElementById('displayMasterName').textContent = master.name || '-';
        document.getElementById('displayMasterEmail').textContent = master.email || '-';
        document.getElementById('displayMasterPhone').textContent = master.phone || '-';
        document.getElementById('displayMasterSpecialization').textContent = master.specialization || '-';
        document.getElementById('displayMasterAddress').textContent = master.address || '-';
        
        // Update rating display
        const rating = master.rating || 0;
        const stars = '★'.repeat(Math.floor(rating)) + '☆'.repeat(5 - Math.floor(rating));
        document.getElementById('masterRating').textContent = stars;
        
        // Update master profile picture - use user photo if available
        const firstLetter = master.name ? master.name.charAt(0).toUpperCase() : 'М';
        document.getElementById('masterInitial').textContent = firstLetter;
        
        // Get user photo from localStorage
        const user = JSON.parse(localStorage.getItem('user'));
        if (user && user.photo_url) {
            document.getElementById('masterProfilePicture').style.background = `url(${user.photo_url}) center/cover`;
            document.getElementById('masterProfilePicture').style.color = 'transparent';
            document.getElementById('masterInitial').style.display = 'none';
        } else {
            document.getElementById('masterProfilePicture').style.background = 'linear-gradient(135deg, var(--primary), var(--secondary))';
            document.getElementById('masterProfilePicture').style.color = 'white';
            document.getElementById('masterInitial').style.display = 'flex';
        }
        
        // Load additional data
        loadMasterWorks();
        loadMasterPaymentInfo();
        loadMasterReviews();
        
    } catch (error) {
        console.error('Error loading master profile:', error);
        // Show no profile state
        document.getElementById('noMasterProfile').style.display = 'block';
        document.getElementById('createMasterForm').style.display = 'none';
        document.getElementById('existingMasterProfile').style.display = 'none';
    }
}

// Show create master form
function showCreateMasterForm() {
    document.getElementById('noMasterProfile').style.display = 'none';
    document.getElementById('createMasterForm').style.display = 'block';
    
    // Pre-fill with user data from profile
    const user = JSON.parse(localStorage.getItem('user'));
    if (user) {
        document.getElementById('masterName').value = user.name || '';
        document.getElementById('masterEmail').value = user.email || '';
        document.getElementById('masterPhone').value = user.phone || '';
        document.getElementById('masterSpecialization').value = 'Автосервис';
        document.getElementById('masterAddress').value = 'Алматы, ул. Тестовая 1';
    }
}

// Hide create master form
function hideCreateMasterForm() {
    document.getElementById('noMasterProfile').style.display = 'block';
    document.getElementById('createMasterForm').style.display = 'none';
}

// Create master profile
async function createMasterProfile(event) {
    event.preventDefault();
    
    const formData = {
        name: document.getElementById('masterName').value,
        email: document.getElementById('masterEmail').value,
        phone: document.getElementById('masterPhone').value,
        specialization: document.getElementById('masterSpecialization').value,
        address: document.getElementById('masterAddress').value
    };
    
    try {
        const token = localStorage.getItem('token');
        const headers = {
            'Content-Type': 'application/json'
        };
        if (token) {
            headers['Authorization'] = 'Bearer ' + token;
        }
        
        const response = await fetch(`${API_URL}/master/profile`, {
            method: 'POST',
            headers: headers,
            body: JSON.stringify(formData)
        });
        
        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || 'Failed to create master profile');
        }
        
        showMessage('Профиль мастера успешно создан!', 'success');
        loadMasterProfile(); // Reload to show the new profile
        
    } catch (error) {
        console.error('Error creating master profile:', error);
        showMessage(`Ошибка создания профиля: ${error.message}`, 'error');
    }
}

// Load master works
async function loadMasterWorks() {
    try {
        const token = localStorage.getItem('token');
        const headers = {};
        if (token) {
            headers['Authorization'] = 'Bearer ' + token;
        }
        
        const response = await fetch(`${API_URL}/master/works`, { headers });
        const works = await response.json();
        
        const container = document.getElementById('worksGallery');
        if (works.length === 0) {
            container.innerHTML = '<p style="color: #999; text-align: center; padding: 20px;">Пока нет работ</p>';
            return;
        }
        
        container.innerHTML = '';
        works.forEach(work => {
            const workCard = document.createElement('div');
            workCard.style.cssText = 'background: white; border-radius: 10px; padding: 15px; box-shadow: 0 2px 10px rgba(0,0,0,0.1);';
            
            workCard.innerHTML = `
                <div style="display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 10px;">
                    <h4 style="margin: 0; color: var(--primary); flex: 1;">${work.title}</h4>
                    <div style="display: flex; gap: 5px;">
                        <button onclick="editWork(${work.id})" style="background: var(--primary); color: white; border: none; border-radius: 50%; width: 30px; height: 30px; cursor: pointer; font-size: 14px; display: flex; align-items: center; justify-content: center;" title="Редактировать работу">✏️</button>
                        <button onclick="deleteWork(${work.id})" style="background: #e74c3c; color: white; border: none; border-radius: 50%; width: 30px; height: 30px; cursor: pointer; font-size: 16px; display: flex; align-items: center; justify-content: center;" title="Удалить работу">×</button>
                    </div>
                </div>
                <p style="margin: 5px 0; color: #666;"><strong>Дата:</strong> ${new Date(work.work_date).toLocaleDateString('ru-RU')}</p>
                <p style="margin: 5px 0; color: #666;"><strong>Заказчик:</strong> ${work.customer_name}</p>
                <p style="margin: 5px 0; color: #666;"><strong>Сумма:</strong> ${work.amount} KZT</p>
                ${work.photo_urls && work.photo_urls.length > 0 ? 
                    `<img src="${work.photo_urls[0]}" style="width: 100%; height: 150px; object-fit: cover; border-radius: 8px; margin-top: 10px;">` : 
                    '<div style="width: 100%; height: 150px; background: #f0f0f0; border-radius: 8px; margin-top: 10px; display: flex; align-items: center; justify-content: center; color: #999;">Нет фото</div>'
                }
            `;
            
            container.appendChild(workCard);
        });
        
    } catch (error) {
        console.error('Error loading master works:', error);
        document.getElementById('worksGallery').innerHTML = '<p style="color: #e74c3c;">Ошибка загрузки работ</p>';
    }
}

// Load master payment info
async function loadMasterPaymentInfo() {
    try {
        const token = localStorage.getItem('token');
        const headers = {};
        if (token) {
            headers['Authorization'] = 'Bearer ' + token;
        }
        
        const response = await fetch(`${API_URL}/master/payment-info`, { headers });
        if (response.status === 404) {
            // No payment info exists
            document.getElementById('kaspiCard').textContent = '-';
            document.getElementById('freedomCard').textContent = '-';
            document.getElementById('halykCard').textContent = '-';
            return;
        }
        
        const info = await response.json();
        
        document.getElementById('kaspiCard').textContent = info.kaspi_card || '-';
        document.getElementById('freedomCard').textContent = info.freedom_card || '-';
        document.getElementById('halykCard').textContent = info.halyk_card || '-';
        
    } catch (error) {
        console.error('Error loading payment info:', error);
    }
}

// Load master reviews
async function loadMasterReviews() {
    try {
        const token = localStorage.getItem('token');
        const headers = {};
        if (token) {
            headers['Authorization'] = 'Bearer ' + token;
        }
        
        const response = await fetch(`${API_URL}/master/reviews`, { headers });
        const reviews = await response.json();
        
        const container = document.getElementById('reviewsPreview');
        if (reviews.length === 0) {
            container.innerHTML = '<p style="color: #999; text-align: center; padding: 20px;">Пока нет отзывов</p>';
            return;
        }
        
        // Show preview of recent reviews
        const recentReviews = reviews.slice(0, 3);
        container.innerHTML = '';
        
        recentReviews.forEach(review => {
            const reviewDiv = document.createElement('div');
            reviewDiv.style.cssText = 'background: #f8f9fa; border-radius: 8px; padding: 15px; margin-bottom: 10px;';
            
            const stars = '★'.repeat(review.rating) + '☆'.repeat(5 - review.rating);
            
            reviewDiv.innerHTML = `
                <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 8px;">
                    <strong>${review.user_name}</strong>
                    <span style="color: #f59e0b;">${stars}</span>
                </div>
                <p style="margin: 0; color: #666; font-size: 14px;">${review.comment || 'Без комментария'}</p>
                <p style="margin: 5px 0 0 0; color: #999; font-size: 12px;">${new Date(review.created_at).toLocaleDateString('ru-RU')}</p>
            `;
            
            container.appendChild(reviewDiv);
        });
        
        if (reviews.length > 3) {
            const moreDiv = document.createElement('div');
            moreDiv.style.cssText = 'text-align: center; margin-top: 10px;';
            moreDiv.innerHTML = `<p style="color: var(--primary); font-size: 14px;">И еще ${reviews.length - 3} отзывов...</p>`;
            container.appendChild(moreDiv);
        }
        
    } catch (error) {
        console.error('Error loading reviews:', error);
        document.getElementById('reviewsPreview').innerHTML = '<p style="color: #e74c3c;">Ошибка загрузки отзывов</p>';
    }
}

// Modal functions

// Show add work modal
function showAddWorkModal() {
    document.getElementById('addWorkModal').style.display = 'flex';
}

// Close add work modal
function closeAddWorkModal() {
    document.getElementById('addWorkModal').style.display = 'none';
    document.getElementById('addWorkForm').reset();
    
    // Reset form title and button
    const modalTitle = document.querySelector('#addWorkModal h2');
    modalTitle.textContent = 'Добавить работу';
    
    const submitButton = document.querySelector('#addWorkForm button[type="submit"]');
    submitButton.textContent = 'Сохранить';
    
    // Clear work ID
    document.getElementById('workId').value = '';
}

// Save work
async function saveWork(event) {
    event.preventDefault();
    
    // Handle photo uploads first
    const photoFiles = document.getElementById('workPhotos').files;
    const photoUrls = [];
    
    // Upload photos if any
    if (photoFiles.length > 0) {
        try {
            const token = localStorage.getItem('token');
            
            for (let i = 0; i < photoFiles.length; i++) {
                const file = photoFiles[i];
                
                if (!file.type.startsWith('image/')) {
                    showMessage('Пожалуйста, выберите только изображения', 'error');
                    return;
                }
                
                if (file.size > 5 * 1024 * 1024) {
                    showMessage('Файл слишком большой (максимум 5MB)', 'error');
                    return;
                }
                
                const formData = new FormData();
                formData.append('photo', file);
                
                console.log('Uploading photo:', file.name);
                
                const response = await fetch(`${API_URL}/master/work-photo`, {
                    method: 'POST',
                    headers: {
                        'Authorization': `Bearer ${token}`
                    },
                    body: formData
                });
                
                console.log('Photo upload response status:', response.status);
                
                if (!response.ok) {
                    let errorMessage = 'Failed to upload photo';
                    try {
                        const error = await response.json();
                        errorMessage = error.error || errorMessage;
                    } catch (e) {
                        const text = await response.text();
                        console.error('Photo upload error (non-JSON):', text);
                        errorMessage = text || errorMessage;
                    }
                    throw new Error(errorMessage);
                }
                
                const result = await response.json();
                console.log('Photo upload result:', result);
                photoUrls.push(result.photo_url);
            }
        } catch (error) {
            console.error('Error uploading photos:', error);
            showMessage(`Ошибка загрузки фото: ${error.message}`, 'error');
            return;
        }
    }
    
    const formData = {
        title: document.getElementById('workTitle').value,
        work_date: document.getElementById('workDate').value,
        customer_name: document.getElementById('workCustomer').value,
        amount: parseFloat(document.getElementById('workAmount').value),
        photo_urls: photoUrls
    };
    
    try {
        const token = localStorage.getItem('token');
        const headers = {
            'Content-Type': 'application/json'
        };
        if (token) {
            headers['Authorization'] = 'Bearer ' + token;
        }
        
        const workId = document.getElementById('workId').value;
        let response;
        
        if (workId) {
            // Update existing work
            response = await fetch(`${API_URL}/master/works/${workId}`, {
                method: 'PUT',
                headers: headers,
                body: JSON.stringify(formData)
            });
        } else {
            // Create new work
            response = await fetch(`${API_URL}/master/works`, {
                method: 'POST',
                headers: headers,
                body: JSON.stringify(formData)
            });
        }
        
        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || 'Failed to save work');
        }
        
        showMessage(workId ? 'Работа успешно обновлена!' : 'Работа успешно добавлена!', 'success');
        closeAddWorkModal();
        loadMasterWorks(); // Reload works
        
    } catch (error) {
        console.error('Error saving work:', error);
        showMessage(`Ошибка сохранения работы: ${error.message}`, 'error');
    }
}

// Show reviews modal
function showReviewsModal() {
    document.getElementById('reviewsModal').style.display = 'flex';
    loadAllReviews();
}

// Close reviews modal
function closeReviewsModal() {
    document.getElementById('reviewsModal').style.display = 'none';
}

// Load all reviews for modal
async function loadAllReviews() {
    try {
        const token = localStorage.getItem('token');
        const headers = {};
        if (token) {
            headers['Authorization'] = 'Bearer ' + token;
        }
        
        const response = await fetch(`${API_URL}/master/reviews`, { headers });
        const reviews = await response.json();
        
        const container = document.getElementById('reviewsList');
        if (reviews.length === 0) {
            container.innerHTML = '<p style="color: #999; text-align: center; padding: 20px;">Пока нет отзывов</p>';
            return;
        }
        
        container.innerHTML = '';
        reviews.forEach(review => {
            const reviewDiv = document.createElement('div');
            reviewDiv.style.cssText = 'background: #f8f9fa; border-radius: 8px; padding: 15px; margin-bottom: 15px;';
            
            const stars = '★'.repeat(review.rating) + '☆'.repeat(5 - review.rating);
            
            reviewDiv.innerHTML = `
                <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 10px;">
                    <strong>${review.user_name}</strong>
                    <span style="color: #f59e0b; font-size: 18px;">${stars}</span>
                </div>
                <p style="margin: 0; color: #666;">${review.comment || 'Без комментария'}</p>
                <p style="margin: 8px 0 0 0; color: #999; font-size: 12px;">${new Date(review.created_at).toLocaleDateString('ru-RU')}</p>
            `;
            
            container.appendChild(reviewDiv);
        });
        
    } catch (error) {
        console.error('Error loading all reviews:', error);
        document.getElementById('reviewsList').innerHTML = '<p style="color: #e74c3c;">Ошибка загрузки отзывов</p>';
    }
}

// Show edit payment modal
function showEditPaymentModal() {
    document.getElementById('editPaymentModal').style.display = 'flex';
}

// Close edit payment modal
function closeEditPaymentModal() {
    document.getElementById('editPaymentModal').style.display = 'none';
}

// Delete master profile
async function deleteMasterProfile() {
    if (!confirm('Вы уверены, что хотите удалить профиль мастера? Это действие нельзя отменить.')) {
        return;
    }
    
    try {
        const token = localStorage.getItem('token');
        const headers = {
            'Content-Type': 'application/json'
        };
        if (token) {
            headers['Authorization'] = 'Bearer ' + token;
        }
        
        const response = await fetch(`${API_URL}/master/profile`, {
            method: 'DELETE',
            headers: headers
        });
        
        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || 'Failed to delete master profile');
        }
        
        showMessage('Профиль мастера успешно удален!', 'success');
        loadMasterProfile(); // Reload to show the no profile state
        
    } catch (error) {
        console.error('Error deleting master profile:', error);
        showMessage(`Ошибка удаления профиля: ${error.message}`, 'error');
    }
}

// Save payment info
async function savePaymentInfo(event) {
    event.preventDefault();
    
    const formData = {
        kaspi_card: document.getElementById('editKaspiCard').value,
        freedom_card: document.getElementById('editFreedomCard').value,
        halyk_card: document.getElementById('editHalykCard').value
    };
    
    try {
        const token = localStorage.getItem('token');
        const headers = {
            'Content-Type': 'application/json'
        };
        if (token) {
            headers['Authorization'] = 'Bearer ' + token;
        }
        
        const response = await fetch(`${API_URL}/master/payment-info`, {
            method: 'PUT',
            headers: headers,
            body: JSON.stringify(formData)
        });
        
        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || 'Failed to save payment info');
        }
        
        showMessage('Реквизиты успешно обновлены!', 'success');
        closeEditPaymentModal();
        loadMasterPaymentInfo(); // Reload payment info
        
    } catch (error) {
        console.error('Error saving payment info:', error);
        showMessage(`Ошибка сохранения реквизитов: ${error.message}`, 'error');
    }
}

// Master photo upload
async function uploadMasterPhoto(event) {
    const file = event.target.files[0];
    if (!file) return;
    
    if (!file.type.startsWith('image/')) {
        showMessage('Пожалуйста, выберите изображение', 'error');
        return;
    }
    
    // Check file size (max 5MB)
    if (file.size > 5 * 1024 * 1024) {
        showMessage('Файл слишком большой (максимум 5MB)', 'error');
        return;
    }
    
    try {
        const formData = new FormData();
        formData.append('photo', file);
        
        const token = localStorage.getItem('token');
        const response = await fetch(`${API_URL}/master/photo`, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${token}`
            },
            body: formData
        });
        
        if (!response.ok) {
            let errorMessage = 'Failed to upload photo';
            try {
                const errorData = await response.json();
                errorMessage = errorData.error || errorMessage;
            } catch (e) {
                const text = await response.text();
                errorMessage = text || errorMessage;
            }
            throw new Error(errorMessage);
        }
        
        const result = await response.json();
        
        // Update both user and master profile pictures
        document.getElementById('profilePicture').style.background = `url(${result.photo_url}) center/cover`;
        document.getElementById('profilePicture').style.color = 'transparent';
        document.getElementById('profileInitial').style.display = 'none';
        
        document.getElementById('masterProfilePicture').style.background = `url(${result.photo_url}) center/cover`;
        document.getElementById('masterProfilePicture').style.color = 'transparent';
        document.getElementById('masterInitial').style.display = 'none';
        
        // Update localStorage
        const user = JSON.parse(localStorage.getItem('user'));
        user.photo_url = result.photo_url;
        localStorage.setItem('user', JSON.stringify(user));
        
        showMessage('Фото профиля успешно загружено!', 'success');
        
    } catch (error) {
        console.error('Error uploading master photo:', error);
        showMessage(`Ошибка загрузки фото: ${error.message}`, 'error');
    }
}

// Review functions

// Show add review modal
function showAddReviewModal() {
    document.getElementById('addReviewModal').style.display = 'flex';
    loadMastersForReview();
}

// Close add review modal
function closeAddReviewModal() {
    document.getElementById('addReviewModal').style.display = 'none';
    document.getElementById('addReviewForm').reset();
    resetRating();
}

// Load masters for review dropdown
async function loadMastersForReview() {
    try {
        const token = localStorage.getItem('token');
        const headers = {};
        if (token) {
            headers['Authorization'] = 'Bearer ' + token;
        }
        
        const response = await fetch(`${API_URL}/masters`, { headers });
        const masters = await response.json();
        
        const select = document.getElementById('reviewMasterId');
        select.innerHTML = '<option value="">Выберите мастера...</option>';
        
        masters.forEach(master => {
            const option = document.createElement('option');
            option.value = master.id;
            option.textContent = master.name;
            select.appendChild(option);
        });
        
    } catch (error) {
        console.error('Error loading masters:', error);
        showMessage('Ошибка загрузки списка мастеров', 'error');
    }
}

// Set rating
function setRating(rating) {
    document.getElementById('reviewRating').value = rating;
    
    // Update star display
    for (let i = 1; i <= 5; i++) {
        const star = document.getElementById(`star${i}`);
        if (i <= rating) {
            star.style.color = '#f59e0b';
        } else {
            star.style.color = '#ddd';
        }
    }
}

// Reset rating
function resetRating() {
    document.getElementById('reviewRating').value = '0';
    for (let i = 1; i <= 5; i++) {
        document.getElementById(`star${i}`).style.color = '#ddd';
    }
}

// Save review
async function saveReview(event) {
    event.preventDefault();
    
    const masterId = document.getElementById('reviewMasterId').value;
    const rating = parseInt(document.getElementById('reviewRating').value);
    const comment = document.getElementById('reviewComment').value.trim();
    
    if (!masterId) {
        showMessage('Выберите мастера', 'error');
        return;
    }
    
    if (rating === 0) {
        showMessage('Поставьте оценку', 'error');
        return;
    }
    
    try {
        const token = localStorage.getItem('token');
        const response = await fetch(`${API_URL}/reviews`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({
                master_id: parseInt(masterId),
                rating: rating,
                comment: comment
            })
        });
        
        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || 'Failed to save review');
        }
        
        showMessage('Отзыв успешно добавлен!', 'success');
        closeAddReviewModal();
        
        // Reload master reviews if we're viewing master profile
        if (document.getElementById('existingMasterProfile').style.display !== 'none') {
            loadMasterReviews();
        }
        
    } catch (error) {
        console.error('Error saving review:', error);
        showMessage(`Ошибка сохранения отзыва: ${error.message}`, 'error');
    }
}

// Master profile editing functions

// Show master edit form
function showMasterEditForm() {
    document.getElementById('editMasterModal').style.display = 'flex';
    
    // Pre-fill form with current master data
    document.getElementById('editMasterName').value = document.getElementById('displayMasterName').textContent;
    document.getElementById('editMasterEmail').value = document.getElementById('displayMasterEmail').textContent;
    document.getElementById('editMasterPhone').value = document.getElementById('displayMasterPhone').textContent;
    document.getElementById('editMasterSpecialization').value = document.getElementById('displayMasterSpecialization').textContent;
    document.getElementById('editMasterAddress').value = document.getElementById('displayMasterAddress').textContent;
}

// Close master edit modal
function closeMasterEditModal() {
    document.getElementById('editMasterModal').style.display = 'none';
}

// Save master profile
async function saveMasterProfile(event) {
    event.preventDefault();
    
    const formData = {
        name: document.getElementById('editMasterName').value,
        email: document.getElementById('editMasterEmail').value,
        phone: document.getElementById('editMasterPhone').value,
        specialization: document.getElementById('editMasterSpecialization').value,
        address: document.getElementById('editMasterAddress').value
    };
    
    try {
        const token = localStorage.getItem('token');
        const response = await fetch(`${API_URL}/master/profile`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify(formData)
        });
        
        console.log('Master profile update response status:', response.status);
        
        if (!response.ok) {
            let errorMessage = 'Failed to update master profile';
            try {
                const error = await response.json();
                console.error('Master profile update error:', error);
                errorMessage = error.error || errorMessage;
            } catch (e) {
                const text = await response.text();
                console.error('Master profile update error (non-JSON):', text);
                errorMessage = text || errorMessage;
            }
            throw new Error(errorMessage);
        }
        
        const result = await response.json();
        console.log('Master profile update result:', result);
        
        showMessage('Профиль мастера успешно обновлен!', 'success');
        closeMasterEditModal();
        loadMasterProfile(); // Reload master profile
        
    } catch (error) {
        console.error('Error updating master profile:', error);
        showMessage(`Ошибка обновления профиля: ${error.message}`, 'error');
    }
}

// Edit work function
async function editWork(workId) {
    try {
        console.log('Loading work with ID:', workId);
        const token = localStorage.getItem('token');
        const response = await fetch(`${API_URL}/master/works/${workId}`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        
        console.log('Load work response status:', response.status);
        
        if (!response.ok) {
            let errorMessage = 'Failed to load work data';
            try {
                const error = await response.json();
                errorMessage = error.error || errorMessage;
            } catch (e) {
                const text = await response.text();
                console.error('Load work error (non-JSON):', text);
                errorMessage = text || errorMessage;
            }
            throw new Error(errorMessage);
        }
        
        const work = await response.json();
        console.log('Loaded work data:', work);
        
        // Fill the form with existing data
        document.getElementById('workId').value = work.id;
        document.getElementById('workTitle').value = work.title;
        document.getElementById('workDate').value = work.work_date.split('T')[0];
        document.getElementById('workCustomer').value = work.customer_name;
        document.getElementById('workAmount').value = work.amount;
        
        // Show the modal
        document.getElementById('addWorkModal').style.display = 'flex';
        
        // Change the form title
        const modalTitle = document.querySelector('#addWorkModal h2');
        modalTitle.textContent = 'Редактировать работу';
        
        // Change the submit button text
        const submitButton = document.querySelector('#addWorkForm button[type="submit"]');
        submitButton.textContent = 'Сохранить изменения';
        
    } catch (error) {
        console.error('Error loading work:', error);
        showMessage(`Ошибка загрузки работы: ${error.message}`, 'error');
    }
}

// Delete work function
async function deleteWork(workId) {
    if (!confirm('Вы точно хотите удалить эту работу?')) {
        return;
    }
    
    try {
        const token = localStorage.getItem('token');
        const response = await fetch(`${API_URL}/master/works/${workId}`, {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        
        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || 'Failed to delete work');
        }
        
        showMessage('Работа успешно удалена!', 'success');
        loadMasterWorks(); // Reload works
        
    } catch (error) {
        console.error('Error deleting work:', error);
        showMessage(`Ошибка удаления работы: ${error.message}`, 'error');
    }
}

