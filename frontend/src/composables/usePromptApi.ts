import { ref } from 'vue'
import type{ Ref } from 'vue'
import { promptRepository } from '../repositories/promptRepository'
import type { ActivePrompt, ApiError, OcrPromptVersionInfo } from '../types/prompt'

export function usePromptApi() {
  const loading: Ref<boolean> = ref(false)
  const error: Ref<ApiError | null> = ref(null)
  const activePrompt: Ref<ActivePrompt | null> = ref(null)
  const versionList: Ref<OcrPromptVersionInfo[]> = ref([])

  /**
   * 現在有効なプロンプトを取得
   */
  const fetchActivePrompt = async () => {
    loading.value = true
    error.value = null
    
    try {
      const data = await promptRepository.getActivePrompt()
      activePrompt.value = data
      return data
    } catch (err) {
      error.value = {
        message: err instanceof Error ? err.message : '不明なエラーが発生しました',
      }
      throw err
    } finally {
      loading.value = false
    }
  }

  /**
   * プロンプトを新バージョンとして保存
   */
  const savePrompt = async (promptContent: string) => {
    loading.value = true
    error.value = null
    
    try {
      const data = await promptRepository.upsertPrompt({ promptContent })
      activePrompt.value = data
      return data
    } catch (err) {
      error.value = {
        message: err instanceof Error ? err.message : '不明なエラーが発生しました',
      }
      throw err
    } finally {
      loading.value = false
    }
  }

  /**
   * バージョンリストを取得
   */
  const fetchVersionList = async () => {
    loading.value = true
    error.value = null
    
    try {
      const data = await promptRepository.getVersionList()
      versionList.value = data
      return data
    } catch (err) {
      error.value = {
        message: err instanceof Error ? err.message : '不明なエラーが発生しました',
      }
      throw err
    } finally {
      loading.value = false
    }
  }

  /**
   * 特定バージョンのプロンプトを取得
   */
  const fetchPromptByVersion = async (version: number) => {
    loading.value = true
    error.value = null
    
    try {
      const data = await promptRepository.getPromptByVersion(version)
      activePrompt.value = data
      return data
    } catch (err) {
      error.value = {
        message: err instanceof Error ? err.message : '不明なエラーが発生しました',
      }
      throw err
    } finally {
      loading.value = false
    }
  }

  return {
    loading,
    error,
    activePrompt,
    versionList,
    fetchActivePrompt,
    savePrompt,
    fetchVersionList,
    fetchPromptByVersion,
  }
}
