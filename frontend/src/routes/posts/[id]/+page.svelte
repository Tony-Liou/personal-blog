<script>
  /** @type {import('./$types').PageData} */
  export let data;

  const { post } = data;

  function formatDate(dateString) {
    return new Date(dateString).toLocaleString('zh-TW', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  }
</script>

<svelte:head>
  <title>{post.title} - 我的部落格</title>
  <meta name="description" content={post.content.substring(0, 150)} />
</svelte:head>

<div class="max-w-3xl mx-auto bg-white p-6 sm:p-8 rounded-lg shadow-lg">
  <!-- 文章標題 -->
  <h1 class="text-3xl sm:text-4xl font-extrabold text-gray-900 mb-4">{post.title}</h1>

  <!-- 作者與日期資訊 -->
  <div class="flex items-center text-gray-500 text-sm mb-6">
    <!-- TODO: 連結到作者頁面 -->
    <a href="#" class="font-medium hover:text-indigo-600">{post.author.username}</a>
    <span class="mx-2">&middot;</span>
    <time datetime={post.created_at}>{formatDate(post.created_at)}</time>
  </div>

  <!-- 封面圖片 (如果有的話) -->
  {#if post.cover_image_url}
    <img
      src={post.cover_image_url}
      alt="文章封面圖片"
      class="w-full h-auto rounded-lg mb-8 object-cover"
    />
  {/if}

  <!-- 文章內容 -->
  <div class="prose prose-lg max-w-none text-gray-800 leading-relaxed">
    <!-- 使用 <pre> 標籤來保留換行和空格，適合顯示純文字內容 -->
    <pre class="whitespace-pre-wrap font-sans">{post.content}</pre>
  </div>

  <div class="mt-12 border-t pt-6">
    <a href="/" class="text-indigo-600 hover:text-indigo-800 font-semibold">
      &larr; 返回文章列表
    </a>
  </div>
</div>
