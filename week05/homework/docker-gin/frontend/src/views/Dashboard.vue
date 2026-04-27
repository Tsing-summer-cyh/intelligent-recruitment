<template>
  <div class="dashboard">
    <header>
      <h2>全栈智能单词本 (AI 驱动)</h2>
      <button @click="handleLogout" class="logout-btn">退出登录</button>
    </header>

    <main class="content">
      <section class="search-section">
        <h3>查询单词</h3>
        <div class="search-box">
          <input v-model="searchWord" placeholder="输入单词，如: apple" />
          <select v-model="aiProvider">
            <option value="DeepSeek">DeepSeek</option>
            <option value="通义千问">通义千问</option>
          </select>
          <button @click="handleSearch">智能查询</button>
        </div>

        <div v-if="searchResult" class="result-card">
          <div class="word-title">{{ searchResult.word }}</div>
          <div class="meaning"><strong>释义：</strong>{{ searchResult.meaning }}</div>
          <div class="sentences">
            <strong>例句：</strong>
            <ul>
              <li v-for="(sentence, index) in searchResult.sentences" :key="index">
                {{ sentence }}
              </li>
            </ul>
          </div>
          <div class="source-tag">数据来源: {{ resultSource }} ({{ searchResult.ai_provider }})</div>
          
          <button v-if="resultSource === 'ai'" @click="handleSave" class="save-btn">
            手动保存到生词本
          </button>
        </div>
      </section>

      <section class="list-section">
        <h3>我的单词本</h3>
        <div class="word-list">
          <div v-for="item in wordList" :key="item.id" class="list-item">
            <div class="item-header">
              <strong>{{ item.word }}</strong>
              <button @click="handleDelete(item.id)" class="delete-btn">删除</button>
            </div>
            <div class="item-meaning">{{ item.meaning }}</div>
          </div>
        </div>
        
        <div class="pagination">
          <button :disabled="page <= 1" @click="changePage(page - 1)">上一页</button>
          <span>第 {{ page }} 页 / 共 {{ Math.ceil(total / pageSize) }} 页</span>
          <button :disabled="page * pageSize >= total" @click="changePage(page + 1)">下一页</button>
        </div>
      </section>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import request from '../api/request';

const router = useRouter();

// 搜索相关的状态
const searchWord = ref('');
const aiProvider = ref('DeepSeek');
const searchResult = ref(null);
const resultSource = ref('');

// 单词列表相关的状态
const wordList = ref([]);
const page = ref(1);
const pageSize = ref(10);
const total = ref(0);

// 获取单词列表
const fetchWordList = async () => {
  try {
    const res = await request.get(`/words?page=${page.value}&page_size=${pageSize.value}`);
    wordList.value = res.data || [];
    total.value = res.total || 0;
  } catch (error) {
    console.error("获取列表失败", error);
  }
};

// 页面加载时拉取一次列表
onMounted(() => {
  fetchWordList();
});

// 智能查询
const handleSearch = async () => {
  if (!searchWord.value.trim()) return alert("请输入单词");
  try {
    searchResult.value = null; // 清空旧数据
    const res = await request.get(`/query?word=${searchWord.value}&ai_provider=${aiProvider.value}`);
    searchResult.value = res.data;
    resultSource.value = res.source; // 'ai' 或者 'database'
  } catch (error) {
    alert("查询失败，请检查后端或AI接口配置");
  }
};

// 手动保存单词
const handleSave = async () => {
  if (!searchResult.value) return;
  try {
    await request.post('/save', searchResult.value);
    alert("保存成功！");
    resultSource.value = 'database'; // 保存后变为已入库状态，隐藏按钮
    fetchWordList(); // 刷新右侧列表
  } catch (error) {
    alert("保存失败");
  }
};

// 删除单词
const handleDelete = async (id) => {
  if (!confirm("确定要删除这个单词吗？")) return;
  try {
    await request.delete(`/words/${id}`);
    fetchWordList(); // 刷新列表
  } catch (error) {
    alert("删除失败");
  }
};

// 分页切换
const changePage = (newPage) => {
  page.value = newPage;
  fetchWordList();
};

// 退出登录
const handleLogout = () => {
  localStorage.removeItem('jwt_token');
  router.push('/login');
};
</script>

<style scoped>
.dashboard { padding: 30px; max-width: 1100px; margin: 0 auto; }

/* 头部样式 */
header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 30px; }
h2 { margin: 0; color: #111827; font-size: 28px; font-weight: bold; letter-spacing: -0.5px; }
.logout-btn { background: white; color: #ef4444; border: 1px solid #f87171; padding: 8px 16px; border-radius: 8px; cursor: pointer; font-weight: 600; transition: all 0.2s; }
.logout-btn:hover { background: #fee2e2; border-color: #ef4444; }

/* 主体双栏布局 */
.content { display: grid; grid-template-columns: 1fr 1fr; gap: 30px; align-items: start; }
@media (max-width: 768px) { .content { grid-template-columns: 1fr; } } /* 屏幕变小时自动变成上下排布 */

/* 卡片通用样式 */
/* === 新增交互：左侧固定，右侧独立滚动 === */

/* 1. 让左侧面板吸顶固定 */
.search-section {
  position: sticky;
  top: 30px; /* 往下滚动时，距离浏览器顶部 30px 的位置停住 */
  height: fit-content; /* 确保它只占据自己内容的高度，不会被右侧拉长 */
}

/* 2. 让右侧面板变成内部独立滚动 */
.list-section {
  max-height: calc(100vh - 110px); /* 动态计算高度：屏幕总高度减去头部和间距 */
  overflow-y: auto; /* 内容超过高度时，出现滚动条 */
}

/* 3. 美化右侧的滚动条（苹果风格的细线滚动条） */
.list-section::-webkit-scrollbar { 
  width: 6px; 
}
.list-section::-webkit-scrollbar-thumb { 
  background: #cbd5e1; 
  border-radius: 10px; 
}
.list-section::-webkit-scrollbar-track { 
  background: transparent; 
}

/* 搜索栏区 */
.search-box { display: flex; gap: 10px; margin-bottom: 25px; }
.search-box input { flex: 1; padding: 10px 14px; border: 1px solid #d1d5db; border-radius: 8px; font-size: 15px; transition: all 0.2s;}
.search-box input:focus { outline: none; border-color: #4f46e5; box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.2); }
.search-box select { padding: 10px; border: 1px solid #d1d5db; border-radius: 8px; background: white; cursor: pointer; outline: none; color: #4b5563;}
.search-box button { background: #4f46e5; color: white; border: none; padding: 0 20px; border-radius: 8px; cursor: pointer; font-weight: 600; transition: all 0.2s; }
.search-box button:hover { background: #4338ca; box-shadow: 0 4px 12px rgba(79, 70, 229, 0.2); }

/* 查询结果卡片 */
.result-card { background: #f8fafc; padding: 24px; border-radius: 12px; border: 1px solid #e2e8f0; }
.word-title { font-size: 32px; font-weight: 800; color: #1e293b; margin-bottom: 12px; }
.meaning { font-size: 16px; color: #334155; margin-bottom: 15px; line-height: 1.6; }
.sentences strong { color: #475569; }
.sentences ul { padding-left: 20px; margin: 8px 0; color: #475569; line-height: 1.6; }
.sentences li { margin-bottom: 6px; }
.source-tag { font-size: 12px; color: #94a3b8; margin-top: 15px; text-align: right; font-style: italic; }

/* 保存按钮 */
.save-btn { background: #10b981; color: white; border: none; padding: 12px 16px; border-radius: 8px; cursor: pointer; width: 100%; margin-top: 15px; font-weight: bold; font-size: 15px; transition: all 0.2s; }
.save-btn:hover { background: #059669; transform: translateY(-1px); box-shadow: 0 4px 12px rgba(16, 185, 129, 0.25); }

/* 单词列表区 */
.word-list { display: flex; flex-direction: column; gap: 12px; }
.list-item { background: #f8fafc; border: 1px solid #e2e8f0; padding: 16px; border-radius: 12px; transition: all 0.2s; }
.list-item:hover { transform: translateY(-2px); box-shadow: 0 6px 12px rgba(0,0,0,0.05); border-color: #cbd5e1; }
.item-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 8px; }
.item-header strong { font-size: 20px; color: #1e293b; font-weight: 700;}
.delete-btn { background: #fee2e2; color: #ef4444; border: none; border-radius: 6px; cursor: pointer; font-size: 13px; padding: 6px 12px; font-weight: 600; transition: all 0.2s; }
.delete-btn:hover { background: #ef4444; color: white; }
.item-meaning { color: #475569; font-size: 14px; line-height: 1.6; }

/* 分页控件 */
.pagination { margin-top: 25px; display: flex; justify-content: space-between; align-items: center; background: #f8fafc; padding: 12px 15px; border-radius: 10px; border: 1px solid #e2e8f0;}
.pagination button { background: white; border: 1px solid #d1d5db; padding: 6px 16px; border-radius: 6px; cursor: pointer; font-weight: 600; color: #374151; transition: all 0.2s; }
.pagination button:hover:not(:disabled) { border-color: #4f46e5; color: #4f46e5; }
.pagination button:disabled { background: #f3f4f6; color: #9ca3af; cursor: not-allowed; border-color: #e5e7eb;}
.pagination span { font-size: 14px; color: #6b7280; font-weight: 500; }
</style>