const API_URL = window.location.origin + '/api/v1';

// Get master ID from URL
function getMasterIdFromURL() {
    const path = window.location.pathname;
    const match = path.match(/\/reviews\/(\d+)/);
    return match ? parseInt(match[1]) : null;
}

let currentMasterId = null;
let currentMaster = null;

// Load page data
window.onload = function() {
    currentMasterId = getMasterIdFromURL();
    if (!currentMasterId) {
        document.getElementById('masterInfoContent').innerHTML = '<p style="color: #e74c3c;">Ошибка: ID мастера не найден</p>';
        return;
    }
    
    loadMasterInfo();
    loadMasterWorks();
    loadMasterCertificates();
    loadReviews();
    checkAuth();
};

// Check if user is authenticated
function checkAuth() {
    const token = localStorage.getItem('token');
    if (token) {
        document.getElementById('addReviewCard').style.display = 'block';
    }
}

// Load master info
async function loadMasterInfo() {
    try {
        const response = await fetch(`${API_URL}/masters/${currentMasterId}`);
        if (!response.ok) {
            throw new Error('Не удалось загрузить информацию о мастере');
        }
        
        currentMaster = await response.json();
        
        // Load verification status
        const statusResponse = await fetch(`${API_URL}/masters/${currentMasterId}/verification-status`);
        let isVerified = false;
        let reviewCount = 0;
        let workCount = 0;
        
        if (statusResponse.ok) {
            const status = await statusResponse.json();
            isVerified = status.is_verified;
            reviewCount = status.review_count;
            workCount = status.work_count;
        }
        
        const rating = currentMaster.rating || 0;
        const stars = '★'.repeat(Math.floor(rating)) + '☆'.repeat(5 - Math.floor(rating));
        
        const verifiedBadge = isVerified ? 
            '<span style="color: #10b981; font-weight: 600; font-size: 16px; margin-left: 10px;">✓ Проверенный мастер</span>' : 
            '<span style="color: #64748b; font-weight: 600; font-size: 16px; margin-left: 10px;">Обычный мастер</span>';
        
        document.getElementById('masterInfoContent').innerHTML = `
            <h2 style="margin-bottom: 20px;">
                👨‍🔧 ${currentMaster.name || '-'}
                ${verifiedBadge}
            </h2>
            <div style="margin-bottom: 15px;">
                <span style="color: #f59e0b; font-size: 24px; margin-right: 10px;">${stars}</span>
                <span style="font-size: 18px; font-weight: 600;">${rating.toFixed(1)}</span>
            </div>
            <div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 20px; margin-top: 20px; text-align: left;">
                <div>
                    <p style="margin: 5px 0; color: #64748b;"><strong>Email:</strong></p>
                    <p style="margin: 0; font-size: 16px;">${currentMaster.email || '-'}</p>
                </div>
                <div>
                    <p style="margin: 5px 0; color: #64748b;"><strong>Телефон:</strong></p>
                    <p style="margin: 0; font-size: 16px;">${currentMaster.phone || '-'}</p>
                </div>
                <div>
                    <p style="margin: 5px 0; color: #64748b;"><strong>Специализация:</strong></p>
                    <p style="margin: 0; font-size: 16px;">${currentMaster.specialization || '-'}</p>
                </div>
                <div>
                    <p style="margin: 5px 0; color: #64748b;"><strong>Адрес:</strong></p>
                    <p style="margin: 0; font-size: 16px;">${currentMaster.address || '-'}</p>
                </div>
            </div>
            <div style="margin-top: 20px; padding: 15px; background: #f0f9ff; border-radius: 8px;">
                <p style="margin: 0; color: #334155;">
                    <strong>Отзывов:</strong> ${reviewCount} | 
                    <strong>Работ:</strong> ${workCount}
                </p>
            </div>
        `;
        
    } catch (error) {
        console.error('Ошибка загрузки информации о мастере:', error);
        document.getElementById('masterInfoContent').innerHTML = '<p style="color: #e74c3c;">Ошибка загрузки информации о мастере</p>';
    }
}

// Load master works
async function loadMasterWorks() {
    try {
        // Note: This endpoint should be created in backend
        const response = await fetch(`${API_URL}/masters/${currentMasterId}/works`);
        if (!response.ok) {
            throw new Error('Не удалось загрузить работы мастера');
        }
        
        const works = await response.json();
        const container = document.getElementById('masterWorksList');
        
        if (works.length === 0) {
            container.innerHTML = '<p style="color: #999; grid-column: 1/-1; text-align: center; padding: 20px;">Пока нет работ</p>';
            return;
        }
        
        container.innerHTML = '';
        works.forEach(work => {
            const workCard = document.createElement('div');
            workCard.style.cssText = 'background: #f8f9fa; border-radius: 10px; padding: 15px; box-shadow: 0 2px 10px rgba(0,0,0,0.1);';
            
            workCard.innerHTML = `
                <h4 style="margin: 0 0 10px 0; color: var(--primary);">${work.title}</h4>
                <p style="margin: 5px 0; color: #666;"><strong>Дата:</strong> ${new Date(work.work_date).toLocaleDateString('ru-RU')}</p>
                <p style="margin: 5px 0; color: #666;"><strong>Заказчик:</strong> ${work.customer_name}</p>
                <p style="margin: 5px 0; color: #666;"><strong>Сумма:</strong> ${work.amount} KZT</p>
                ${work.photo_urls && work.photo_urls.length > 0 ? 
                    `<img src="${work.photo_urls[0]}" style="width: 100%; height: 150px; object-fit: cover; border-radius: 8px; margin-top: 10px;">` : 
                    '<div style="width: 100%; height: 150px; background: #e2e8f0; border-radius: 8px; margin-top: 10px; display: flex; align-items: center; justify-content: center; color: #999;">Нет фото</div>'
                }
            `;
            
            container.appendChild(workCard);
        });
        
    } catch (error) {
        console.error('Ошибка загрузки работ:', error);
        document.getElementById('masterWorksList').innerHTML = '<p style="color: #e74c3c;">Ошибка загрузки работ</p>';
    }
}

// Load master certificates
async function loadMasterCertificates() {
    try {
        const response = await fetch(`${API_URL}/masters/${currentMasterId}/certificates`);
        if (!response.ok) {
            throw new Error('Не удалось загрузить сертификаты');
        }
        
        const certificates = await response.json();
        const container = document.getElementById('masterCertificatesList');
        
        if (certificates.length === 0) {
            container.innerHTML = '<p style="color: #999; grid-column: 1/-1; text-align: center; padding: 20px;">Пока нет сертификатов</p>';
            return;
        }
        
        container.innerHTML = '';
        certificates.forEach(cert => {
            const certCard = document.createElement('div');
            certCard.style.cssText = 'background: #f8f9fa; border-radius: 10px; padding: 15px; box-shadow: 0 2px 10px rgba(0,0,0,0.1);';
            
            certCard.innerHTML = `
                <h4 style="margin: 0 0 10px 0; color: var(--primary);">${cert.name}</h4>
                <img src="${cert.photo_url}" style="width: 100%; height: 200px; object-fit: cover; border-radius: 8px; margin-top: 10px;" alt="${cert.name}">
            `;
            
            container.appendChild(certCard);
        });
        
    } catch (error) {
        console.error('Ошибка загрузки сертификатов:', error);
        document.getElementById('masterCertificatesList').innerHTML = '<p style="color: #999; grid-column: 1/-1; text-align: center; padding: 20px;">Пока нет сертификатов</p>';
    }
}

// Load reviews
async function loadReviews() {
    try {
        const response = await fetch(`${API_URL}/masters/${currentMasterId}/reviews`);
        if (!response.ok) {
            throw new Error('Не удалось загрузить отзывы');
        }
        
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
                    <strong style="font-size: 16px;">${review.user_name}</strong>
                    <span style="color: #f59e0b; font-size: 18px;">${stars}</span>
                </div>
                <p style="margin: 0; color: #666; line-height: 1.6;">${review.comment || 'Без комментария'}</p>
                <p style="margin: 8px 0 0 0; color: #999; font-size: 12px;">${new Date(review.created_at).toLocaleDateString('ru-RU')}</p>
            `;
            
            container.appendChild(reviewDiv);
        });
        
    } catch (error) {
        console.error('Ошибка загрузки отзывов:', error);
        document.getElementById('reviewsList').innerHTML = '<p style="color: #e74c3c;">Ошибка загрузки отзывов</p>';
    }
}

// Set review rating
function setReviewRating(rating) {
    document.getElementById('reviewRating').value = rating;
    
    for (let i = 1; i <= 5; i++) {
        const star = document.getElementById(`star${i}`);
        if (i <= rating) {
            star.style.color = '#f59e0b';
        } else {
            star.style.color = '#ddd';
        }
    }
}

// Reset review form
function resetReviewForm() {
    document.getElementById('reviewRating').value = '0';
    document.getElementById('reviewComment').value = '';
    for (let i = 1; i <= 5; i++) {
        document.getElementById(`star${i}`).style.color = '#ddd';
    }
}

// Save review
async function saveReview(event) {
    event.preventDefault();
    
    const token = localStorage.getItem('token');
    if (!token) {
        alert('Необходима авторизация для добавления отзыва');
        return;
    }
    
    const rating = parseInt(document.getElementById('reviewRating').value);
    const comment = document.getElementById('reviewComment').value.trim();
    
    if (rating === 0) {
        alert('Поставьте оценку');
        return;
    }
    
    try {
        const response = await fetch(`${API_URL}/reviews`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({
                master_id: currentMasterId,
                rating: rating,
                comment: comment
            })
        });
        
        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || 'Не удалось сохранить отзыв');
        }
        
        alert('Отзыв успешно добавлен!');
        resetReviewForm();
        loadReviews();
        
        // Reload master info to update rating
        loadMasterInfo();
        
    } catch (error) {
        console.error('Ошибка сохранения отзыва:', error);
        alert(`Ошибка сохранения отзыва: ${error.message}`);
    }
}

