<script>
  import { onMount } from 'svelte';
  import api from '$lib/api.js';

  let author = null;
  let isLoading = true;
  let error = null;

  onMount(async () => {
    try {
      // 假設作者 ID 為 1，在實際應用中可能需要從配置中取得
      const response = await api.getAuthor(1);
      author = response.data;
    } catch (err) {
      console.error('載入作者資訊失敗:', err);
      error = '無法載入作者資訊';
    } finally {
      isLoading = false;
    }
  });
</script>

<svelte:head>
  <title>關於作者 - 我的部落格</title>
</svelte:head>

<div class="max-w-4xl mx-auto">
  {#if isLoading}
    <div class="text-center py-12">
      <p class="text-gray-500">正在載入作者資訊...</p>
    </div>
  {:else if error}
    <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded-lg">
      <span class="block sm:inline">{error}</span>
    </div>
  {:else if author}
    <div class="bg-white p-8 rounded-lg shadow-md">
      <div class="text-center mb-8">
        {#if author.avatar_url}
          <img 
            src={author.avatar_url} 
            alt={author.username}
            class="w-32 h-32 rounded-full mx-auto mb-4 object-cover border-4 border-gray-200"
          />
        {:else}
          <div class="w-32 h-32 rounded-full mx-auto mb-4 bg-gray-200 flex items-center justify-center border-4 border-gray-300">
            <span class="text-4xl text-gray-500">{author.username.charAt(0).toUpperCase()}</span>
          </div>
        {/if}
        <h1 class="text-4xl font-bold text-gray-900 mb-2">{author.username}</h1>
        <p class="text-gray-600">部落格作者</p>
      </div>

      <div class="prose max-w-none">
        {#if author.bio}
          <div class="text-gray-700 leading-relaxed text-lg">
            {#each author.bio.split('\n') as paragraph}
              {#if paragraph.trim()}
                <p class="mb-4">{paragraph}</p>
              {/if}
            {/each}
          </div>
        {:else}
          <div class="text-center py-8">
            <p class="text-gray-500 text-lg">作者尚未填寫個人簡介。</p>
          </div>
        {/if}
      </div>

      <div class="mt-8 pt-8 border-t border-gray-200">
        <div class="flex flex-wrap gap-4 justify-center">
          <div class="text-center">
            <p class="text-sm text-gray-500">加入時間</p>
            <p class="font-medium">{new Date(author.created_at).toLocaleDateString('zh-TW', {
              year: 'numeric',
              month: 'long',
              day: 'numeric'
            })}</p>
          </div>
        </div>
      </div>

      <div class="mt-8 text-center">
        <a 
          href="/"
          class="inline-block bg-indigo-600 text-white font-bold py-3 px-6 rounded-lg hover:bg-indigo-700 transition-colors"
        >
          瀏覽文章
        </a>
      </div>
    </div>
  {/if}
</div>