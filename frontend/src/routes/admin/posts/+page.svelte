<script>
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import api from '$lib/api.js';
  import { authToken } from '$lib/auth.js';

  let posts = [];
  let isLoading = true;
  let error = null;
  let token = null;

  onMount(async () => {
    const unsubToken = authToken.subscribe(value => {
      token = value;
    });

    await loadPosts();

    return unsubToken;
  });

  async function loadPosts() {
    try {
      const response = await api.getPosts(1, 50); // Load more posts for admin
      posts = response.data;
    } catch (err) {
      console.error("讀取文章失敗:", err);
      error = "無法載入文章，請稍後再試。";
    } finally {
      isLoading = false;
    }
  }

  async function deletePost(postId, postTitle) {
    if (!confirm(`確定要刪除文章「${postTitle}」嗎？此操作無法復原。`)) {
      return;
    }

    try {
      await api.deletePost(postId, token);
      posts = posts.filter(post => post.id !== postId);
    } catch (err) {
      console.error('刪除文章失敗:', err);
      if (err.response && err.response.data && err.response.data.error) {
        alert(`刪除失敗: ${err.response.data.error}`);
      } else {
        alert('刪除文章時發生錯誤，請稍後再試。');
      }
    }
  }

  function editPost(postId) {
    goto(`/admin/posts/edit/${postId}`);
  }

  function formatDate(dateString) {
    return new Date(dateString).toLocaleDateString('zh-TW', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  }
</script>

<svelte:head>
  <title>文章管理 - 我的部落格</title>
</svelte:head>

<div class="max-w-6xl mx-auto">
  <div class="flex items-center justify-between mb-8">
    <h1 class="text-4xl font-bold">文章管理</h1>
    <div class="flex gap-4">
      <a 
        href="/admin/new" 
        class="bg-indigo-600 text-white font-bold py-2 px-4 rounded-lg hover:bg-indigo-700 transition-colors"
      >
        撰寫新文章
      </a>
      <a 
        href="/admin" 
        class="px-4 py-2 text-gray-600 hover:text-gray-800 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
      >
        返回管理頁面
      </a>
    </div>
  </div>

  {#if isLoading}
    <div class="text-center py-12">
      <p class="text-gray-500">正在載入文章...</p>
    </div>
  {:else if error}
    <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded-lg">
      <span class="block sm:inline">{error}</span>
    </div>
  {:else if posts.length === 0}
    <div class="text-center py-12">
      <p class="text-gray-500 mb-4">目前沒有任何文章。</p>
      <a 
        href="/admin/new" 
        class="bg-indigo-600 text-white font-bold py-2 px-4 rounded-lg hover:bg-indigo-700 transition-colors"
      >
        撰寫第一篇文章
      </a>
    </div>
  {:else}
    <div class="bg-white rounded-lg shadow-md overflow-hidden">
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                文章標題
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                作者
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                發布時間
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                更新時間
              </th>
              <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                操作
              </th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            {#each posts as post (post.id)}
              <tr class="hover:bg-gray-50">
                <td class="px-6 py-4 whitespace-nowrap">
                  <div class="flex items-center">
                    <div>
                      <div class="text-sm font-medium text-gray-900">
                        <a 
                          href={`/posts/${post.id}`} 
                          class="hover:text-indigo-600 transition-colors"
                          target="_blank"
                        >
                          {post.title}
                        </a>
                      </div>
                      <div class="text-sm text-gray-500 truncate max-w-xs">
                        {post.content.substring(0, 100)}{post.content.length > 100 ? '...' : ''}
                      </div>
                    </div>
                  </div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {post.author.username}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {formatDate(post.created_at)}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {formatDate(post.updated_at)}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                  <div class="flex justify-end gap-2">
                    <button
                      on:click={() => editPost(post.id)}
                      class="text-indigo-600 hover:text-indigo-900 px-3 py-1 rounded border border-indigo-200 hover:bg-indigo-50 transition-colors"
                    >
                      編輯
                    </button>
                    <button
                      on:click={() => deletePost(post.id, post.title)}
                      class="text-red-600 hover:text-red-900 px-3 py-1 rounded border border-red-200 hover:bg-red-50 transition-colors"
                    >
                      刪除
                    </button>
                  </div>
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    </div>
  {/if}
</div>