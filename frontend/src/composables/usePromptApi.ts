import { ref } from "vue";
import type { Ref } from "vue";
import { promptRepository } from "../repositories/promptRepository";
import type { Prompt, ApiError, PromptInfo } from "../types/prompt";

export function usePromptApi() {
  const loading: Ref<boolean> = ref(false);
  const error: Ref<ApiError | null> = ref(null);
  const currentPrompt: Ref<Prompt | null> = ref(null);
  const promptList: Ref<PromptInfo[]> = ref([]);

  /**
   * プロンプトを保存
   */
  const savePrompt = async (
    name: string,
    promptContent: string,
    id?: string
  ) => {
    loading.value = true;
    error.value = null;

    try {
      const data = await promptRepository.upsertPrompt({
        id,
        name,
        promptContent,
      });
      currentPrompt.value = data;
      return data;
    } catch (err) {
      error.value = {
        message:
          err instanceof Error ? err.message : "不明なエラーが発生しました",
      };
      throw err;
    } finally {
      loading.value = false;
    }
  };

  /**
   * プロンプト一覧を取得
   */
  const fetchPromptList = async () => {
    loading.value = true;
    error.value = null;

    try {
      const data = await promptRepository.getPromptList();
      promptList.value = data;
      return data;
    } catch (err) {
      error.value = {
        message:
          err instanceof Error ? err.message : "不明なエラーが発生しました",
      };
      throw err;
    } finally {
      loading.value = false;
    }
  };

  /**
   * 特定のプロンプトを取得
   */
  const fetchPromptById = async (id: string) => {
    loading.value = true;
    error.value = null;

    try {
      const data = await promptRepository.getPromptById(id);
      currentPrompt.value = data;
      return data;
    } catch (err) {
      error.value = {
        message:
          err instanceof Error ? err.message : "不明なエラーが発生しました",
      };
      throw err;
    } finally {
      loading.value = false;
    }
  };

  /**
   * プロンプトを削除
   */
  const deletePrompt = async (id: string) => {
    loading.value = true;
    error.value = null;

    try {
      await promptRepository.deletePrompt(id);
      // 削除後、リストから該当項目を除外
      promptList.value = promptList.value.filter((p) => p.id !== id);
      // 現在のプロンプトが削除された場合はクリア
      if (currentPrompt.value?.id === id) {
        currentPrompt.value = null;
      }
    } catch (err) {
      error.value = {
        message:
          err instanceof Error ? err.message : "不明なエラーが発生しました",
      };
      throw err;
    } finally {
      loading.value = false;
    }
  };

  return {
    loading,
    error,
    currentPrompt,
    promptList,
    savePrompt,
    fetchPromptList,
    fetchPromptById,
    deletePrompt,
  };
}
