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
            const dateStr = apt.date ? (typeof apt.date === 'string' ? apt.date : apt.date.split('T')[0]) : 'Не указана';
            div.innerHTML = `
                <div class="appointment-header">
                    <h3>Заказ #${apt.id}</h3>
                    <span class="status status-${apt.status || 'pending'}">${(apt.status || 'pending')}</span>
                </div>
                <div class="appointment-body">
                    <p><strong>Дата:</strong> ${dateStr} в ${apt.time || 'Не указано'}</p>
                    <p><strong>Услуга:</strong> ${apt.service_name || 'Не указана'}</p>
                    <p><strong>Мастер:</strong> ${apt.master_name || 'Не назначен'}</p>
                    ${apt.comment ? `<p><strong>Комментарий:</strong> ${apt.comment}</p>` : ''}
                </div>
            `;
            container.appendChild(div);
        });
    } catch (error) {
        console.error('Error loading appointments:', error);
        document.getElementById('userAppointments').innerHTML = '<p style="color: #e74c3c;">Ошибка загрузки заказов</p>';
    }
}

