<script>
    import { onMount } from 'svelte';
    import { goto } from '$app/navigation';
    import { isAuthenticated } from '$lib/auth.js';

    onMount(() => {
        // 這是一個簡單的客戶端路由保護
        // 如果使用者未登入，將他們導向到登入頁面
        const unsub = isAuthenticated.subscribe(isAuth => {
            if (!isAuth) {
                goto('/login');
            }
        });

        // 組件銷毀時取消訂閱
        return unsub;
    });
</script>

<svelte:head>
    <title>後台管理 - 我的部落格</title>
</svelte:head>

<div class="max-w-4xl mx-auto">
    <h1 class="text-4xl font-bold mb-8">後台管理</h1>
    <p class="text-lg text-gray-700">歡迎來到管理後台！</p>
    
    <div class="mt-8 grid grid-cols-1 md:grid-cols-2 gap-6">
        <a href="/admin/new-post" class="bg-white p-6 rounded-lg shadow-md hover:shadow-xl transition-shadow duration-300">
            <h2 class="text-2xl font-bold text-indigo-600">撰寫新文章</h2>
            <p class="text-gray-600 mt-2">建立一篇全新的部落格文章。</p>
        </a>
        <div class="bg-white p-6 rounded-lg shadow-md">
            <h2 class="text-2xl font-bold text-gray-500">管理現有文章</h2>
            <p class="text-gray-600 mt-2">（功能開發中）編輯或刪除已發表的文章。</p>
        </div>
    </div>
</div>
