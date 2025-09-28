import api from '$lib/api.js';
import { error } from '@sveltejs/kit';

/** @type {import('./$types').PageServerLoad} */
export async function load({ params }) {
  try {
    const response = await api.getPost(params.id);
    return {
      post: response.data
    };
  } catch (err) {
    // 如果 API 回傳 404 (找不到文章) 或其他錯誤
    if (err.response && err.response.status) {
        throw error(err.response.status, '找不到這篇文章');
    }
    // 其他伺服器錯誤
    throw error(500, '無法從伺服器載入文章，請稍後再試。');
  }
}
