const API_BASE_URL = '/api';

function formatDate(date) {
    return new Date(date).toLocaleString('es-ES', {
        day: '2-digit',
        month: '2-digit',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    });
}

async function loadPosts() {
    try {
        const response = await axios.get(`${API_BASE_URL}/posts`);
        const posts = response.data;
        const postsList = document.getElementById('posts-list');
        postsList.innerHTML = posts.map(post => `
            <div class="post">
                <div class="post-header">
                    <h3>${post.title}</h3>
                    <div class="likes">
                        <button onclick="likePost('${post.id}')">👍 ${post.likes}</button>
                    </div>
                </div>
                <p>${post.description}</p>
                <a href="${post.url}" class="post-link" target="_blank">Ver descuento</a>
                <div class="post-meta">
                    <span>Por: ${post.author.name || 'Anónimo'}</span>
                    <span class="timestamp">Creado: ${formatDate(post.creation_time)}</span>
                    <span class="timestamp">Expira: ${formatDate(post.expire_time)}</span>
                </div>
                <button onclick="deletePost('${post.id}')">Eliminar</button>
            </div>
        `).join('');
    } catch (error) {
        console.error('Error:', error);
    }
}

function clearErrors() {
    document.querySelectorAll('.error').forEach(el => el.textContent = '');
}

function showErrors(errors) {
    errors.forEach(error => {
        const errorElement = document.getElementById(`${error.field}-error`);
        if (errorElement) {
            errorElement.textContent = error.message;
        }
    });
}

document.getElementById('create-post-form').onsubmit = async (e) => {
    e.preventDefault();
    clearErrors();

    const title = document.getElementById('title').value;
    const description = document.getElementById('description').value;
    const url = document.getElementById('url').value;
    const expireTime = document.getElementById('expire-time').value;

    try {
        await axios.post(`${API_BASE_URL}/posts`, {
            title,
            description,
            url,
            expire_time: new Date(expireTime).toISOString()
        });
        loadPosts();
        e.target.reset();
    } catch (error) {
        if (error.response?.data?.errors) {
            showErrors(error.response.data.errors);
        } else {
            console.error('Error:', error);
        }
    }
};

async function deletePost(id) {
    if (confirm('¿Estás seguro de que quieres eliminar este descuento?')) {
        try {
            await axios.delete(`${API_BASE_URL}/posts/${id}`);
            loadPosts();
        } catch (error) {
            console.error('Error:', error);
        }
    }
}

async function likePost(id) {
    try {
        await axios.post(`${API_BASE_URL}/posts/${id}/like`);
        loadPosts();
    } catch (error) {
        console.error('Error:', error);
    }
}

const expireTimeInput = document.getElementById('expire-time');
const now = new Date();
now.setMinutes(now.getMinutes() - now.getTimezoneOffset());
expireTimeInput.min = now.toISOString().slice(0, 16);

// Cargar posts al iniciar
loadPosts();