<script>
  import { goto } from '$app/navigation';
  import { onMount } from 'svelte';
  import api from '$lib/api.js';
  import { authToken } from '$lib/auth.js';
  import ImageUpload from '$lib/components/ImageUpload.svelte';

  let title = '';
  let content = '';
  let coverImageUrl = '';
  let isLoading = false;
  let error = null;
  let token = null;

  onMount(() => {
    const unsubToken = authToken.subscribe(value => {
      token = value;
    });

    return unsubToken;
  });

  async function handleSubmit() {
    if (!title.trim() || !content.trim()) {
      error = '標題和內容都是必填的';
      return;
    }

    error = null;
    isLoading = true;

    try {
      const postData = {
        title: title.trim(),
        content: content.trim(),
        cover_image_url: coverImageUrl.trim() || ''
      };

      await api.createPost(postData, token);
      
      goto('/admin');
    } catch (err) {
      if (err.response && err.response.data && err.response.data.error) {
        error = err.response.data.error;
      } else {
        error = '建立文章時發生錯誤，請稍後再試。';
      }
      console.error('建立文章失敗:', err);
    } finally {
      isLoading = false;
    }
  }

  function handleCancel() {
    goto('/admin');
  }
</script>

<svelte:head>
  <title>撰寫新文章 - 我的部落格</title>
</svelte:head>

<div class="max-w-4xl mx-auto">
  <div class="flex items-center justify-between mb-8">
    <h1 class="text-4xl font-bold">撰寫新文章</h1>
    <button
      type="button"
      on:click={handleCancel}
      class="px-4 py-2 text-gray-600 hover:text-gray-800 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
    >
      返回管理頁面
    </button>
  </div>

  <form on:submit|preventDefault={handleSubmit} class="bg-white p-8 rounded-lg shadow-md">
    <div class="mb-6">
      <label for="title" class="block text-gray-700 font-medium mb-2">文章標題 *</label>
      <input
        type="text"
        id="title"
        bind:value={title}
        placeholder="輸入文章標題..."
        class="w-full px-4 py-3 border rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
        required
      />
    </div>

    <div class="mb-6">
      <ImageUpload 
        bind:value={coverImageUrl}
        label="封面圖片（選填）"
        placeholder="上傳圖片或輸入圖片網址..."
      />
    </div>

    <div class="mb-8">
      <label for="content" class="block text-gray-700 font-medium mb-2">文章內容 *</label>
      <textarea
        id="content"
        bind:value={content}
        placeholder="開始撰寫您的文章內容..."
        rows="20"
        class="w-full px-4 py-3 border rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 resize-y min-h-96"
        required
      ></textarea>
      <p class="text-sm text-gray-500 mt-2">
        支援 Markdown 格式。使用空行來分段。
      </p>
    </div>

    {#if error}
      <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded-lg mb-6" role="alert">
        <span class="block sm:inline">{error}</span>
      </div>
    {/if}

    <div class="flex gap-4">
      <button
        type="submit"
        disabled={isLoading}
        class="flex-1 bg-indigo-600 text-white font-bold py-3 px-6 rounded-lg hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 disabled:bg-indigo-400 transition-colors"
      >
        {isLoading ? '發布中...' : '發布文章'}
      </button>
      
      <button
        type="button"
        on:click={handleCancel}
        disabled={isLoading}
        class="px-6 py-3 text-gray-700 border border-gray-300 rounded-lg hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-gray-300 focus:ring-offset-2 disabled:opacity-50 transition-colors"
      >
        取消
      </button>
    </div>
  </form>
</div>