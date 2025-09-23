<script>
  import { onMount } from 'svelte';
  import api from '$lib/api.js';

  let posts = [];
  let isLoading = true;
  let error = null;

  onMount(async () => {
    try {
      const response = await api.getPosts();
      posts = response.data;
    } catch (err) {
      console.error("讀取文章失敗:", err);
      error = "無法載入文章，請稍後再試。";
    } finally {
      isLoading = false;
    }
  });

  function formatDate(dateString) {
    return new Date(dateString).toLocaleDateString('zh-TW', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });
  }
</script>

<svelte:head>
  <title>首頁 - 我的部落格</title>
</svelte:head>

<div class="max-w-4xl mx-auto">
  <h1 class="text-4xl font-bold mb-8 text-center">最新文章</h1>

  {#if isLoading}
    <p class="text-center text-gray-500">正在載入文章...</p>
  {:else if error}
    <p class="text-center text-red-500">{error}</p>
  {:else if posts.length === 0}
    <p class="text-center text-gray-500">目前沒有任何文章。</p>
  {:else}
    <div class="space-y-8">
      {#each posts as post (post.id)}
        <div class="bg-white p-6 rounded-lg shadow-md hover:shadow-xl transition-shadow duration-300">
          <h2 class="text-2xl font-bold mb-2">
            <!-- TODO: 連結到單篇文章頁面 -->
            <a href={`/posts/${post.id}`} class="text-gray-900 hover:text-indigo-600">{post.title}</a>
          </h2>
          <div class="text-gray-500 text-sm mb-4">
            <span>由 {post.author.username} 發表於 {formatDate(post.created_at)}</span>
          </div>
          <p class="text-gray-700 leading-relaxed">
            <!-- 顯示部分內容作為預覽 -->
            {post.content.substring(0, 150)}{post.content.length > 150 ? '...' : ''}
          </p>
           <a href={`/posts/${post.id}`} class="text-indigo-600 hover:text-indigo-800 font-semibold mt-4 inline-block">
            繼續閱讀 &rarr;
          </a>
        </div>
      {/each}
    </div>
  {/if}
</div>
