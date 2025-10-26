<script setup lang="ts">
import { ref, onMounted } from "vue";
import { usePromptApi } from "../composables/usePromptApi";
import type { PromptInfo } from "../types/prompt";

const { promptList, fetchPromptList, fetchPromptById, currentPrompt } =
  usePromptApi();

const emit = defineEmits<{
  promptSelected: [prompt: any];
}>();

const selectedId = ref<string | null>(null);
const firstPrompt = ref<PromptInfo | undefined>();

// プロンプト一覧を取得
const loadPromptList = async () => {
  try {
    await fetchPromptList();
    // 初回は最初のプロンプトを選択
    if (promptList.value.length > 0 && !selectedId.value) {
      firstPrompt.value = promptList.value[0];
      if (firstPrompt.value) {
        selectedId.value = firstPrompt.value.id;
        await fetchPromptById(firstPrompt.value.id);
        if (currentPrompt.value) {
          emit("promptSelected", currentPrompt.value);
        }
      }
    }
  } catch (error) {
    console.error("Failed to load prompt list:", error);
  }
};

// プロンプトクリック時の処理
const handlePromptClick = async (id: string) => {
  selectedId.value = id;
  try {
    await fetchPromptById(id);
    if (currentPrompt.value) {
      emit("promptSelected", currentPrompt.value);
    }
  } catch (error) {
    console.error("Failed to load prompt:", error);
  }
};

// 新規作成
const handleCreateNew = () => {
  // 選択を解除
  selectedId.value = null;

  // 空のプロンプトを作成
  const newPrompt = {
    id: null,
    name: "",
    promptContent: "",
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
    userId: null,
  };

  // 親コンポーネントに通知
  emit("promptSelected", newPrompt);
};

// 再読み込みボタンの処理（親から呼ばれる想定）
const reload = async () => {
  await loadPromptList();
};

// 初期読み込み
onMounted(() => {
  loadPromptList();
});

// 親コンポーネントに公開
defineExpose({
  reload,
});
</script>

<template>
  <div class="prompt-list">
    <div class="list-header">
      <h2 class="list-title">プロンプト一覧</h2>
      <button class="button button-primary" @click="handleCreateNew">
        <span>➕</span> 新規作成
      </button>
    </div>

    <ul class="prompt-items">
      <li
        v-for="prompt in promptList"
        :key="prompt.id"
        :class="['prompt-item', { active: selectedId === prompt.id }]"
        @click="handlePromptClick(prompt.id)"
      >
        <div class="prompt-info">
          <span class="prompt-name">{{ prompt.name }}</span>
          <span class="prompt-date">{{
            new Date(prompt.createdAt).toLocaleDateString("ja-JP")
          }}</span>
        </div>
      </li>
    </ul>

    <div v-if="promptList.length === 0" class="empty-state">
      プロンプトがありません
      <button
        class="button button-primary"
        @click="handleCreateNew"
        style="margin-top: 1rem"
      >
        <span>➕</span> 最初のプロンプトを作成
      </button>
    </div>
  </div>
</template>

<style scoped>
.prompt-list {
  background: white;
  border-radius: 0.5rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  padding: 1.5rem;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.list-title {
  font-size: 1.125rem;
  font-weight: 600;
  margin: 0;
}

.reload-button {
  background: #e9ecef;
  color: #495057;
  border: none;
  padding: 0.375rem 0.75rem;
  border-radius: 0.375rem;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  gap: 0.25rem;
}

.reload-button:hover {
  background: #dee2e6;
}

.prompt-items {
  list-style: none;
  padding: 0;
  margin: 0;
  overflow-y: auto;
  flex: 1;
}

.prompt-items::-webkit-scrollbar {
  display: none;
}

.prompt-item {
  padding: 0.75rem 1rem;
  margin-bottom: 0.5rem;
  background: #f8f9fa;
  border-radius: 0.375rem;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.prompt-item:hover {
  background: #e9ecef;
}

.prompt-item.active {
  background: #5eb0f9;
  color: white;
}

.prompt-info {
  display: flex;
  flex-direction: column;
  gap: 0.125rem;
}

.prompt-name {
  font-size: 0.875rem;
  font-weight: 500;
}

.prompt-date {
  font-size: 0.75rem;
  opacity: 0.7;
}

.prompt-item.active .prompt-date {
  opacity: 0.9;
}

.status-badge {
  font-size: 0.75rem;
  padding: 0.125rem 0.5rem;
  border-radius: 9999px;
  background: #e9ecef;
  color: #6c757d;
  white-space: nowrap;
}

.status-badge.active {
  background: #28a745;
  color: white;
}

.prompt-item.active .status-badge {
  background: rgba(255, 255, 255, 0.3);
  color: white;
}

.empty-state {
  text-align: center;
  color: #6c757d;
  padding: 2rem;
  font-size: 0.875rem;
}

.button {
  padding: 0.5rem 1rem;
  border-radius: 0.375rem;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  border: none;
  display: flex;
  align-items: center;
  gap: 0.375rem;
  white-space: nowrap;
}

.button-primary {
  background: #667eea;
  color: white;
}

.button-primary:hover:not(:disabled) {
  background: #5a67d8;
}

/* その他既存のスタイル */
.prompt-list {
  background: white;
  border-radius: 0.5rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  padding: 1.5rem;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.list-title {
  font-size: 1.125rem;
  font-weight: 600;
  margin: 0;
}

.prompt-items {
  list-style: none;
  padding: 0;
  margin: 0;
  overflow-y: auto;
  flex: 1;
}

.prompt-items::-webkit-scrollbar {
  display: none;
}

.prompt-item {
  padding: 0.75rem 1rem;
  margin-bottom: 0.5rem;
  background: #f8f9fa;
  border-radius: 0.375rem;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.prompt-item:hover {
  background: #e9ecef;
}

.prompt-item.active {
  background: #5eb0f9;
  color: white;
}

.prompt-info {
  display: flex;
  flex-direction: column;
  gap: 0.125rem;
}

.prompt-name {
  font-size: 0.875rem;
  font-weight: 500;
}

.prompt-date {
  font-size: 0.75rem;
  opacity: 0.7;
}

.prompt-item.active .prompt-date {
  opacity: 0.9;
}

.empty-state {
  text-align: center;
  color: #6c757d;
  padding: 2rem;
  font-size: 0.875rem;
}
</style>
