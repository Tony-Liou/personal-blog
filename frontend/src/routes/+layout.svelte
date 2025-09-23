<script>
  import '../app.css';
</script>

<div class="min-h-screen bg-gray-50 font-sans text-gray-800">
  <!-- 導覽列 -->
  <nav class="bg-white shadow-md">
    <div class="container mx-auto px-6 py-4 flex justify-between items-center">
      <a href="/" class="text-2xl font-bold text-gray-800 hover:text-indigo-600">
        我的部落格
      </a>
      <div class="space-x-4">
        <a href="/" class="text-gray-600 hover:text-indigo-600">首頁</a>
        <!-- TODO: 新增 "關於作者" 和 "登入" 連結 -->
        <a href="/login" class="text-gray-600 hover:text-indigo-600">登入</a>
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
