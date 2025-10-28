const API_URL = window.location.origin + '/api/v1';

// Load user data on page load
window.onload = function() {
    loadUserProfile();
    loadUserAppointments();
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

function editField(fieldName) {
    const fieldMap = {
        'name': { element: 'userName', label: 'Имя', type: 'text' },
        'email': { element: 'userEmail', label: 'Email', type: 'email' },
        'phone': { element: 'userPhone', label: 'Телефон', type: 'tel' }
    };
    
    const field = fieldMap[fieldName];
    if (!field) return;
    
    const spanElement = document.getElementById(field.element);
    const currentValue = spanElement.textContent;
    
    // Create input field
    const input = document.createElement('input');
    input.type = field.type;
    input.value = currentValue;
    input.style.cssText = `
        flex: 1;
        padding: 8px 12px;
        border: 2px solid var(--primary);
        border-radius: 6px;
        font-size: 14px;
        outline: none;
    `;
    
    // Create save and cancel buttons
    const saveBtn = document.createElement('button');
    saveBtn.innerHTML = '✓';
    saveBtn.style.cssText = `
        width: 32px;
        height: 32px;
        background: #10b981;
        color: white;
        border: none;
        border-radius: 6px;
        cursor: pointer;
        font-size: 14px;
        margin-left: 5px;
    `;
    
    const cancelBtn = document.createElement('button');
    cancelBtn.innerHTML = '✕';
    cancelBtn.style.cssText = `
        width: 32px;
        height: 32px;
        background: #ef4444;
        color: white;
        border: none;
        border-radius: 6px;
        cursor: pointer;
        font-size: 14px;
        margin-left: 5px;
    `;
    
    // Replace span with input and buttons
    const parentDiv = spanElement.parentElement;
    parentDiv.replaceChild(input, spanElement);
    parentDiv.appendChild(saveBtn);
    parentDiv.appendChild(cancelBtn);
    
    // Focus input
    input.focus();
    input.select();
    
    // Save function
    const saveField = async () => {
        const newValue = input.value.trim();
        if (newValue === currentValue) {
            cancelEdit();
            return;
        }
        
        try {
            const token = localStorage.getItem('token');
            const updateData = { [fieldName]: newValue };
            
            const response = await fetch(`${API_URL}/user/profile`, {
                method: 'PUT',
                headers: { 
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify(updateData)
            });

            if (!response.ok) {
                throw new Error('Failed to update field');
            }

            // Update localStorage
            const updatedUser = { ...window.currentUser, [fieldName]: newValue };
            localStorage.setItem('user', JSON.stringify(updatedUser));
            window.currentUser = updatedUser;
            
            // Update display
            spanElement.textContent = newValue;
            parentDiv.replaceChild(spanElement, input);
            parentDiv.removeChild(saveBtn);
            parentDiv.removeChild(cancelBtn);
            
            showMessage(`${field.label} успешно обновлен!`, 'success');
            
        } catch (error) {
            console.error('Error updating field:', error);
            showMessage(`Ошибка при обновлении ${field.label.toLowerCase()}: ` + error.message, 'error');
        }
    };
    
    // Cancel function
    const cancelEdit = () => {
        spanElement.textContent = currentValue;
        parentDiv.replaceChild(spanElement, input);
        parentDiv.removeChild(saveBtn);
        parentDiv.removeChild(cancelBtn);
    };
    
    // Event listeners
    saveBtn.onclick = saveField;
    cancelBtn.onclick = cancelEdit;
    
    input.addEventListener('keydown', (e) => {
        if (e.key === 'Enter') {
            saveField();
        } else if (e.key === 'Escape') {
            cancelEdit();
        }
    });
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
function uploadProfilePhoto(event) {
    const file = event.target.files[0];
    if (!file) return;
    
    if (!file.type.startsWith('image/')) {
        showMessage('Пожалуйста, выберите изображение', 'error');
        return;
    }
    
    const reader = new FileReader();
    reader.onload = function(e) {
        document.getElementById('profilePicture').style.background = `url(${e.target.result}) center/cover`;
        document.getElementById('profilePicture').style.color = 'transparent';
        document.getElementById('profileInitial').style.display = 'none';
        showMessage('Фото загружено!', 'success');
    };
    reader.readAsDataURL(file);
}

// Edit appointment
function editAppointment(appointmentId) {
    const modal = document.getElementById('editAppointmentModal');
    const form = document.getElementById('editAppointmentForm');
    
    // Load appointment data
    fetch(`${API_URL}/appointments/${appointmentId}`)
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
        closeEditModal();
        loadUserAppointments(); // Reload appointments
        
    } catch (error) {
        console.error('Error updating appointment:', error);
        showMessage('Ошибка при обновлении записи: ' + error.message, 'error');
    }
}

async function loadUserAppointments() {
    try {
        const response = await fetch(`${API_URL}/appointments`);
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
                </div>
            `;
            container.appendChild(div);
        });
    } catch (error) {
        console.error('Error loading appointments:', error);
        document.getElementById('userAppointments').innerHTML = '<p style="color: #e74c3c;">Ошибка загрузки заказов</p>';
    }
}

