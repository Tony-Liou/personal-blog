<script>
  import { page } from '$app/stores';
  import { onMount } from 'svelte';
  import api from '$lib/api.js';

  let authorId = null;
  let author = null;
  let authorPosts = [];
  let isLoading = true;
  let error = null;

  $: authorId = $page.params.id;

  onMount(async () => {
    if (authorId) {
      await loadAuthorData();
    }
  });

  async function loadAuthorData() {
    try {
      // 載入作者資訊
      const authorResponse = await api.getAuthor(authorId);
      author = authorResponse.data;

      // 載入作者的文章
      const postsResponse = await api.getPosts(1, 50);
      authorPosts = postsResponse.data.filter(post => post.author_id === parseInt(authorId));
    } catch (err) {
      console.error('載入作者資訊失敗:', err);
      error = '無法載入作者資訊，請檢查作者ID是否正確。';
    } finally {
      isLoading = false;
    }
  }

  function formatDate(dateString) {
    return new Date(dateString).toLocaleDateString('zh-TW', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });
  }
</script>

<svelte:head>
  <title>{author ? `${author.username} - 作者頁面` : '作者頁面'} - 我的部落格</title>
</svelte:head>

<div class="max-w-6xl mx-auto">
  {#if isLoading}
    <div class="text-center py-12">
      <p class="text-gray-500">正在載入作者資訊...</p>
    </div>
  {:else if error}
    <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded-lg">
      <span class="block sm:inline">{error}</span>
      <div class="mt-4">
        <a href="/" class="text-indigo-600 hover:text-indigo-800 underline">返回首頁</a>
      </div>
    </div>
  {:else if author}
    <!-- 作者資訊區塊 -->
    <div class="bg-white p-8 rounded-lg shadow-md mb-8">
      <div class="flex flex-col md:flex-row items-center md:items-start gap-6">
        <div class="flex-shrink-0">
          {#if author.avatar_url}
            <img 
              src={author.avatar_url} 
              alt={author.username}
              class="w-32 h-32 rounded-full object-cover border-4 border-gray-200"
            />
          {:else}
            <div class="w-32 h-32 rounded-full bg-gray-200 flex items-center justify-center border-4 border-gray-300">
              <span class="text-4xl text-gray-500">{author.username.charAt(0).toUpperCase()}</span>
            </div>
          {/if}
        </div>
        
        <div class="flex-1 text-center md:text-left">
          <h1 class="text-4xl font-bold text-gray-900 mb-2">{author.username}</h1>
          <p class="text-gray-600 mb-4">部落格作者</p>
          
          {#if author.bio}
            <div class="text-gray-700 leading-relaxed">
              {#each author.bio.split('\n') as paragraph}
                {#if paragraph.trim()}
                  <p class="mb-2">{paragraph}</p>
                {/if}
              {/each}
            </div>
          {:else}
            <p class="text-gray-500">作者尚未填寫個人簡介。</p>
          {/if}
          
          <div class="mt-4 text-sm text-gray-500">
            <p>加入時間: {formatDate(author.created_at)}</p>
            <p>文章數量: {authorPosts.length} 篇</p>
          </div>
        </div>
      </div>
    </div>

    <!-- 作者文章列表 -->
    <div class="mb-8">
      <h2 class="text-3xl font-bold mb-6">{author.username} 的文章</h2>
      
      {#if authorPosts.length === 0}
        <div class="text-center py-8">
          <p class="text-gray-500">此作者尚未發表任何文章。</p>
        </div>
      {:else}
        <div class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          {#each authorPosts as post (post.id)}
            <div class="bg-white p-6 rounded-lg shadow-md hover:shadow-xl transition-shadow duration-300">
              {#if post.cover_image_url}
                <img 
                  src={post.cover_image_url} 
                  alt={post.title}
                  class="w-full h-48 object-cover rounded-lg mb-4"
                />
              {/if}
              
              <h3 class="text-xl font-bold mb-2">
                <a href={`/posts/${post.id}`} class="text-gray-900 hover:text-indigo-600 transition-colors">
                  {post.title}
                </a>
              </h3>
              
              <p class="text-gray-600 text-sm mb-3">
                {formatDate(post.created_at)}
              </p>
              
              <p class="text-gray-700 leading-relaxed mb-4">
                {post.content.substring(0, 100)}{post.content.length > 100 ? '...' : ''}
              </p>
              
              <a 
                href={`/posts/${post.id}`} 
                class="text-indigo-600 hover:text-indigo-800 font-semibold"
              >
                繼續閱讀 &rarr;
              </a>
            </div>
          {/each}
        </div>
      {/if}
    </div>

    <div class="text-center">
      <a 
        href="/"
        class="inline-block bg-indigo-600 text-white font-bold py-3 px-6 rounded-lg hover:bg-indigo-700 transition-colors"
      >
        瀏覽所有文章
      </a>
    </div>
  {/if}
</div>