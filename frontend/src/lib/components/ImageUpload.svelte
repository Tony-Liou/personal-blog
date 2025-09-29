<script>
  import { createEventDispatcher } from 'svelte';
  import api from '$lib/api.js';
  import { authToken } from '$lib/auth.js';

  export let value = '';
  export let placeholder = '上傳圖片或輸入圖片網址...';
  export let label = '圖片';

  const dispatch = createEventDispatcher();

  let fileInput;
  let isUploading = false;
  let uploadError = null;
  let token = null;

  authToken.subscribe(t => token = t);

  async function handleFileSelect(event) {
    const file = event.target.files[0];
    if (!file) return;

    // 檢查檔案類型
    const allowedTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/gif', 'image/webp'];
    if (!allowedTypes.includes(file.type)) {
      uploadError = '請選擇有效的圖片檔案 (JPEG, PNG, GIF, WebP)';
      return;
    }

    // 檢查檔案大小 (10MB)
    const maxSize = 10 * 1024 * 1024;
    if (file.size > maxSize) {
      uploadError = '檔案大小不能超過 10MB';
      return;
    }

    uploadError = null;
    isUploading = true;

    try {
      const response = await api.uploadFile(file, token);
      const imageUrl = `http://localhost:8080${response.data.url}`;
      value = imageUrl;
      dispatch('upload', { url: imageUrl });
      
      // 清空 file input
      fileInput.value = '';
    } catch (err) {
      console.error('上傳失敗:', err);
      if (err.response && err.response.data && err.response.data.error) {
        uploadError = err.response.data.error;
      } else {
        uploadError = '上傳失敗，請稍後再試';
      }
    } finally {
      isUploading = false;
    }
  }

  function triggerFileInput() {
    fileInput.click();
  }

  function handleInput(event) {
    value = event.target.value;
    dispatch('input', { value });
  }
</script>

<div class="space-y-3">
  <label class="block text-gray-700 font-medium text-sm">{label}</label>
  
  <div class="flex gap-3">
    <input
      type="url"
      bind:value
      on:input={handleInput}
      {placeholder}
      class="flex-1 px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 text-sm"
    />
    
    <button
      type="button"
      on:click={triggerFileInput}
      disabled={isUploading}
      class="px-4 py-2 bg-gray-100 text-gray-700 border border-gray-300 rounded-lg hover:bg-gray-200 focus:outline-none focus:ring-2 focus:ring-gray-300 disabled:opacity-50 disabled:cursor-not-allowed transition-colors text-sm font-medium"
    >
      {isUploading ? '上傳中...' : '選擇檔案'}
    </button>
  </div>

  <input
    bind:this={fileInput}
    type="file"
    accept="image/*"
    on:change={handleFileSelect}
    class="hidden"
  />

  {#if uploadError}
    <p class="text-red-600 text-sm">{uploadError}</p>
  {/if}

  {#if value}
    <div class="mt-3">
      <p class="text-sm text-gray-600 mb-2">預覽：</p>
      <img 
        src={value} 
        alt="Preview" 
        class="max-w-full h-auto max-h-48 rounded-lg border border-gray-200"
        on:error={() => uploadError = '無法載入圖片預覽'}
      />
    </div>
  {/if}
</div>