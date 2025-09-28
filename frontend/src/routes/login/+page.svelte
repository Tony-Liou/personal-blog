<script>
  import { goto } from '$app/navigation';
  import api from '$lib/api.js';
  import { authToken } from '$lib/auth.js';

  let username = '';
  let password = '';
  let error = null;
  let isLoading = false;

  async function handleLogin() {
    error = null;
    isLoading = true;
    try {
      const response = await api.login(username, password);
      const token = response.data.token;

      // 登入成功，將 token 存入 store
      authToken.set(token);

      // 導向到後台管理頁面
      goto('/admin');

    } catch (err) {
      if (err.response && err.response.data && err.response.data.error) {
        error = err.response.data.error;
      } else {
        error = '登入時發生未知的錯誤，請稍後再試。';
      }
      console.error('登入失敗:', err);
    } finally {
      isLoading = false;
    }
  }
</script>

<svelte:head>
  <title>管理員登入 - 我的部落格</title>
</svelte:head>

<div class="max-w-md mx-auto mt-10">
  <div class="bg-white p-8 rounded-lg shadow-md">
    <h1 class="text-3xl font-bold text-center mb-6">管理員登入</h1>
    <form on:submit|preventDefault={handleLogin}>
      <div class="mb-4">
        <label for="username" class="block text-gray-700 font-medium mb-2">使用者名稱</label>
        <input
          type="text"
          id="username"
          bind:value={username}
          class="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500"
          required
        />
      </div>
      <div class="mb-6">
        <label for="password" class="block text-gray-700 font-medium mb-2">密碼</label>
        <input
          type="password"
          id="password"
          bind:value={password}
          class="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500"
          required
        />
      </div>

      {#if error}
        <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded-lg relative mb-4" role="alert">
          <span class="block sm:inline">{error}</span>
        </div>
      {/if}

      <button
        type="submit"
        disabled={isLoading}
        class="w-full bg-indigo-600 text-white font-bold py-2 px-4 rounded-lg hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 disabled:bg-indigo-400"
      >
        {isLoading ? '登入中...' : '登入'}
      </button>
    </form>
  </div>
</div>
