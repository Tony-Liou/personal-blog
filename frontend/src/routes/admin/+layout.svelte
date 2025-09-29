<script>
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { page } from '$app/stores';
  import { isAuthenticated, logout } from '$lib/auth.js';

  let isAuth = false;

  onMount(() => {
    const unsub = isAuthenticated.subscribe(value => {
      isAuth = value;
      if (!value) {
        goto('/login');
      }
    });

    return unsub;
  });

  function handleLogout() {
    logout();
    goto('/');
  }

  $: currentPath = $page.url.pathname;
</script>

<svelte:head>
  <title>後台管理 - 我的部落格</title>
</svelte:head>

{#if isAuth}
  <div class="min-h-screen bg-gray-50">
    <!-- Admin Navigation Bar -->
    <nav class="bg-white shadow-sm border-b border-gray-200">
      <div class="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <div class="flex items-center">
            <a href="/" class="text-xl font-bold text-gray-900 hover:text-indigo-600 transition-colors">
              我的部落格
            </a>
            <span class="ml-4 px-3 py-1 bg-indigo-100 text-indigo-800 text-sm font-medium rounded-full">
              管理模式
            </span>
          </div>
          
          <div class="flex items-center space-x-4">
            <a 
              href="/admin"
              class="px-3 py-2 rounded-md text-sm font-medium transition-colors {currentPath === '/admin' ? 'bg-indigo-100 text-indigo-700' : 'text-gray-700 hover:text-indigo-600 hover:bg-gray-100'}"
            >
              管理首頁
            </a>
            <a 
              href="/admin/new"
              class="px-3 py-2 rounded-md text-sm font-medium transition-colors {currentPath === '/admin/new' ? 'bg-indigo-100 text-indigo-700' : 'text-gray-700 hover:text-indigo-600 hover:bg-gray-100'}"
            >
              撰寫文章
            </a>
            <a 
              href="/admin/posts"
              class="px-3 py-2 rounded-md text-sm font-medium transition-colors {currentPath.startsWith('/admin/posts') ? 'bg-indigo-100 text-indigo-700' : 'text-gray-700 hover:text-indigo-600 hover:bg-gray-100'}"
            >
              管理文章
            </a>
            <a 
              href="/"
              class="px-3 py-2 rounded-md text-sm font-medium text-gray-700 hover:text-indigo-600 hover:bg-gray-100 transition-colors"
              target="_blank"
            >
              檢視網站
            </a>
            <button
              on:click={handleLogout}
              class="px-3 py-2 rounded-md text-sm font-medium text-red-700 hover:text-red-900 hover:bg-red-50 transition-colors"
            >
              登出
            </button>
          </div>
        </div>
      </div>
    </nav>

    <!-- Main Content -->
    <main class="py-8 px-4 sm:px-6 lg:px-8">
      <slot />
    </main>
  </div>
{:else}
  <div class="min-h-screen bg-gray-50 flex items-center justify-center">
    <div class="text-center">
      <p class="text-gray-500">正在驗證身份...</p>
    </div>
  </div>
{/if}