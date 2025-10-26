<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { usePromptApi } from '../composables/usePromptApi'

const { 
  versionList, 
  fetchVersionList, 
  fetchPromptByVersion,
  activePrompt 
} = usePromptApi()

const emit = defineEmits<{
  promptSelected: [prompt: any]
}>()

const selectedVersion = ref<number | null>(null)

// バージョンリストを取得
const loadVersionList = async () => {
  try {
    await fetchVersionList()
    // 有効なバージョンを自動選択
    const activeVersion = versionList.value.find(v => v.isActive)
    if (activeVersion) {
      selectedVersion.value = activeVersion.version
      await fetchPromptByVersion(activeVersion.version)
      // 親コンポーネントに通知
      if (activePrompt.value) {
        emit('promptSelected', {
          ...activePrompt.value,
          isActive: activeVersion.isActive
        })
      }
    }
  } catch (error) {
    console.error('Failed to load version list:', error)
  }
}

// バージョンクリック時の処理
const handleVersionClick = async (version: number) => {
  selectedVersion.value = version
  try {
    // 選択されたバージョンの完全な情報を取得
    const selectedVersionInfo = versionList.value.find(v => v.version === version)
    
    await fetchPromptByVersion(version)
    
    // デバッグ用
    console.log('selectedVersionInfo:', selectedVersionInfo)
    console.log('activePrompt.value:', activePrompt.value)
    
    // 親コンポーネントに通知（versionListからisActive情報を取得）
    if (activePrompt.value && selectedVersionInfo) {
      const promptToEmit = {
        ...activePrompt.value,
        isActive: selectedVersionInfo.isActive
      }
      console.log('Emitting prompt:', promptToEmit)
      emit('promptSelected', promptToEmit)
    }
  } catch (error) {
    console.error('Failed to load prompt:', error)
  }
}

// 再読み込みボタンの処理（親から呼ばれる想定）
const reload = async () => {
  await loadVersionList()
}

// 初期読み込み
onMounted(() => {
  loadVersionList()
})

// 親コンポーネントに公開
defineExpose({
  reload
})
</script>

<template>
  <div class="prompt-list">
    <div class="list-header">
      <h2 class="list-title">プロンプト変更履歴</h2>
      <!-- <button class="reload-button" @click="loadVersionList">
        <span>🔄</span> 更新
      </button> -->
    </div>
    
    <ul class="prompt-items">
      <li 
        v-for="prompt in versionList" 
        :key="prompt.id"
        :class="['prompt-item', { active: selectedVersion === prompt.version }]"
        @click="handleVersionClick(prompt.version)"
      >
        <div class="prompt-info">
          <span class="prompt-name">Version {{ prompt.version }}</span>
          <span class="prompt-date">{{ new Date(prompt.createdAt).toLocaleDateString('ja-JP') }}</span>
        </div>
        <span :class="['status-badge', { active: prompt.isActive }]">
          {{ prompt.isActive ? '有効' : '無効' }}
        </span>
      </li>
    </ul>

    <div v-if="versionList.length === 0" class="empty-state">
      履歴がありません
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
</style>
