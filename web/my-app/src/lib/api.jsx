const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

// Generic fetch wrapper - update fungsi ini
async function fetchAPI(endpoint, options = {}) {
  const url = `${API_BASE_URL}${endpoint}`;
  
  const defaultHeaders = {
    'Content-Type': 'application/json',
  };

  if (typeof window !== 'undefined') {
    const token = localStorage.getItem('access_token');
    if (token) {
      defaultHeaders['Authorization'] = `Bearer ${token}`;
    }
  }

  try {
    const response = await fetch(url, {
      ...options,
      headers: {
        ...defaultHeaders,
        ...options.headers,
      },
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({ message: 'Request failed' }));
      throw new Error(error.message || `HTTP ${response.status}`);
    }

    const data = await response.json();
    return data;
  } catch (error) {
    console.error(`API Error [${endpoint}]:`, error.message);
    throw error;
  }
}

// ========================================
// ARTICLES API
// ========================================
export const articlesAPI = {
  getAll: (params = {}) => {
    const searchParams = new URLSearchParams();
    if (params.kategori) searchParams.set('kategori', params.kategori);
    if (params.status) searchParams.set('status', params.status);
    if (params.limit) searchParams.set('limit', params.limit.toString());
    
    const query = searchParams.toString();
    return fetchAPI(`/articles${query ? `?${query}` : ''}`);
  },

  getBySlug: (slug) => fetchAPI(`/articles/slug/${slug}`),

  getById: (id) => fetchAPI(`/articles/${id}`),
};

// ========================================
// CATEGORIES API
// ========================================
export const categoriesAPI = {
  getAll: () => fetchAPI('/categories'),
  getById: (id) => fetchAPI(`/categories/${id}`),
};

// ========================================
// TAGS API
// =======================================
export const tagsAPI = {
  getAll: () => fetchAPI('/tags'),
};

// ========================================
// AUTH API
// ========================================
export const authAPI = {
  login: (email, password) =>
    fetchAPI('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    }),

  register: (username, email, password) =>
    fetchAPI('/auth/register', {
      method: 'POST',
      body: JSON.stringify({ username, email, password }),
    }),

  logout: () =>
    fetchAPI('/auth/logout', {
      method: 'POST',
    }),

  refresh: (refreshToken) =>
    fetchAPI('/auth/refresh', {
      method: 'POST',
      body: JSON.stringify({ refresh_token: refreshToken }),
    }),
};

// ========================================
// COMMENTS API
// ========================================
export const commentsAPI = {
  getByArticle: (articleId) => fetchAPI(`/articles/${articleId}/comments`),

  create: (articleId, konten) =>
    fetchAPI(`/articles/${articleId}/comments`, {
      method: 'POST',
      body: JSON.stringify({ konten }),
    }),
};

// ========================================
// USER API
// ========================================
export const userAPI = {
  getProfile: () => fetchAPI('/users/me'),
  
  updateProfile: (username, email) =>
    fetchAPI('/users/me', {
      method: 'PUT',
      body: JSON.stringify({ username, email }),
    }),

  changePassword: (currentPassword, newPassword) =>
    fetchAPI('/users/me/password', {
      method: 'PUT',
      body: JSON.stringify({ 
        current_password: currentPassword, 
        new_password: newPassword 
      }),
    }),
};