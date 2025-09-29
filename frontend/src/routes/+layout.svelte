<script>
  import '../app.css';
  import { goto } from '$app/navigation';
  import { isAuthenticated, logout } from '$lib/auth.js';

  function handleLogout() {
    logout();
    goto('/'); // 登出後導向回首頁
  }
</script>

<div class="min-h-screen bg-gray-50 font-sans text-gray-800">
  <!-- 導覽列 -->
  <nav class="bg-white shadow-md">
    <div class="container mx-auto px-6 py-4 flex justify-between items-center">
      <a href="/" class="text-2xl font-bold text-gray-800 hover:text-indigo-600">
        我的部落格
      </a>
      <div class="space-x-4 flex items-center">
        <a href="/" class="text-gray-600 hover:text-indigo-600">首頁</a>
        <a href="/about" class="text-gray-600 hover:text-indigo-600">關於作者</a>
        
        {#if $isAuthenticated}
          <!-- 如果已登入 -->
          <a href="/admin" class="text-gray-600 hover:text-indigo-600">後台管理</a>
          <button on:click={handleLogout} class="text-gray-600 hover:text-indigo-600">登出</button>
        {:else}
          <!-- 如果未登入 -->
          <a href="/login" class="bg-indigo-600 text-white font-semibold py-2 px-4 rounded-lg hover:bg-indigo-700">
            登入
          </a>
        {/if}
      </div>
    </div>
  </nav>

  <!-- 頁面內容 -->
  <main class="container mx-auto px-6 py-8">
    <slot />
  </main>

  <!-- 頁腳 -->
  <footer class="bg-white mt-8 py-4">
    <div class="container mx-auto px-6 text-center text-gray-500">
      &copy; {new Date().getFullYear()} Tony Liou's Blog. All rights reserved.
    </div>
  </footer>
</div>
