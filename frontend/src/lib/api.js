import axios from 'axios';

// 建立一個 axios 實例，設定後端 API 的基礎 URL
// 在開發環境中，這通常指向你的 Go 伺服器
const apiClient = axios.create({
  baseURL: 'http://localhost:8080/api/v1', // 確保這與你的 Go 後端 Port 一致
  headers: {
    'Content-Type': 'application/json',
  },
});

export default {
  // 取得所有文章
  getPosts(page = 1, limit = 10) {
    return apiClient.get(`/posts/?page=${page}&limit=${limit}`);
  },

  // 根據 ID 取得單篇文章
  getPost(id) {
    return apiClient.get(`/posts/${id}`);
  },

  // --- 需要驗證的 API ---

  // 登入
  login(username, password) {
    return apiClient.post('/auth/login', { username, password });
  },

  // 建立文章
  createPost(postData, token) {
    return apiClient.post('/posts', postData, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
  },

  // 更新文章
  updatePost(id, postData, token) {
    return apiClient.put(`/posts/${id}`, postData, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
  },

  // 刪除文章
  deletePost(id, token) {
    return apiClient.delete(`/posts/${id}`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
  },
};
