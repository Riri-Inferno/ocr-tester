import type { Prompt, CreatePromptRequest, PromptInfo } from "../types/prompt";
import { mockPrompts, mockPromptDetails } from "../mocks/promptMock";

// 開発環境ではプロキシ経由、本番環境では直接API
const API_BASE_URL = import.meta.env.DEV
  ? "/api"
  : import.meta.env.VITE_API_TARGET_URL;

// モックモードかどうか
const USE_MOCK = import.meta.env.VITE_USE_MOCK === "true";

export const promptRepository = {
  /**
   * プロンプト情報一覧を取得
   */
  async getPromptList(): Promise<PromptInfo[]> {
    if (USE_MOCK) {
      // モックデータを返す（少し遅延を入れて本物っぽく）
      await new Promise((resolve) => setTimeout(resolve, 300));
      return mockPrompts;
    }

    try {
      const response = await fetch(`${API_BASE_URL}/prompts`, {
        method: "GET",
        headers: {
          Accept: "application/json",
        },
      });

      if (!response.ok) {
        throw new Error(`API Error: ${response.status} ${response.statusText}`);
      }

      return await response.json();
    } catch (error) {
      console.error("Fetch error:", error);
      throw error;
    }
  },

  /**
   * 特定のプロンプトを取得
   */
  async getPromptById(id: string): Promise<Prompt> {
    if (USE_MOCK) {
      await new Promise((resolve) => setTimeout(resolve, 300));
      const prompt = mockPromptDetails[id];
      if (!prompt) {
        throw new Error(`Prompt not found: ${id}`);
      }
      return prompt;
    }

    try {
      const response = await fetch(`${API_BASE_URL}/prompts/${id}`, {
        method: "GET",
        headers: {
          Accept: "application/json",
        },
      });

      if (!response.ok) {
        throw new Error(`API Error: ${response.status} ${response.statusText}`);
      }

      return await response.json();
    } catch (error) {
      console.error("Fetch error:", error);
      throw error;
    }
  },

  /**
   * プロンプトを更新/作成（upsert）
   */
  async upsertPrompt(data: CreatePromptRequest): Promise<Prompt> {
    if (USE_MOCK) {
      await new Promise((resolve) => setTimeout(resolve, 300));
      // 新規作成のモック
      const newId = String(Date.now());
      const newPrompt: Prompt = {
        id: newId,
        name: data.name,
        promptContent: data.promptContent,
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
        userId: data.userId,
      };
      // モックデータに追加（実際には永続化されない）
      mockPromptDetails[newId] = newPrompt;
      mockPrompts.push({
        id: newId,
        name: data.name,
        createdAt: newPrompt.createdAt,
        updatedAt: newPrompt.updatedAt,
        userId: data.userId,
      });
      return newPrompt;
    }

    try {
      const response = await fetch(`${API_BASE_URL}/prompts/upsert`, {
        method: "POST",
        headers: {
          Accept: "application/json",
          "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
      });

      if (!response.ok) {
        throw new Error(`API Error: ${response.status} ${response.statusText}`);
      }

      return await response.json();
    } catch (error) {
      console.error("Fetch error:", error);
      throw error;
    }
  },

  /**
   * プロンプトを削除
   */
  async deletePrompt(id: string): Promise<void> {
    if (USE_MOCK) {
      await new Promise((resolve) => setTimeout(resolve, 300));
      // モックから削除
      const index = mockPrompts.findIndex((p) => p.id === id);
      if (index !== -1) {
        mockPrompts.splice(index, 1);
      }
      delete mockPromptDetails[id];
      return;
    }

    try {
      const response = await fetch(`${API_BASE_URL}/prompts/${id}`, {
        method: "DELETE",
        headers: {
          Accept: "application/json",
        },
      });

      if (!response.ok) {
        throw new Error(`API Error: ${response.status} ${response.statusText}`);
      }
    } catch (error) {
      console.error("Fetch error:", error);
      throw error;
    }
  },
};
