import type { Prompt, CreatePromptRequest, PromptInfo } from "../types/prompt";

// 開発環境ではプロキシ経由、本番環境では直接API
const API_BASE_URL = import.meta.env.DEV
  ? "/api"
  : import.meta.env.VITE_API_TARGET_URL;

export const promptRepository = {
  /**
   * プロンプト情報一覧を取得
   */
  async getPromptList(): Promise<PromptInfo[]> {
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
