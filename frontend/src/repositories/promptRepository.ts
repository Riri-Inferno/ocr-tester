import type { ActivePrompt, CreatePromptRequest, OcrPromptVersionInfo } from '../types/prompt'

// 開発環境ではプロキシ経由、本番環境では直接API
const API_BASE_URL = import.meta.env.DEV 
  ? '/api' 
  : import.meta.env.VITE_API_TARGET_URL

export const promptRepository = {
  /**
   * 現在有効なプロンプトを取得
   */
  async getActivePrompt(): Promise<ActivePrompt> {
    try {
      const response = await fetch(`${API_BASE_URL}/OcrPrompt/active`, {
        method: 'GET',
        headers: {
          'Accept': 'application/json',
        },
      })

      if (!response.ok) {
        throw new Error(`API Error: ${response.status} ${response.statusText}`)
      }

      return await response.json()
    } catch (error) {
      console.error('Fetch error:', error)
      throw error
    }
  },

  /**
   * プロンプトを更新/作成（upsert）
   */
  async upsertPrompt(data: CreatePromptRequest): Promise<ActivePrompt> {
    try {
      const response = await fetch(`${API_BASE_URL}/OcrPrompt/upsert`, {
        method: 'POST',
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(data),
      })

      if (!response.ok) {
        throw new Error(`API Error: ${response.status} ${response.statusText}`)
      }

      return await response.json()
    } catch (error) {
      console.error('Fetch error:', error)
      throw error
    }
  },
  /**
   * プロンプトのバージョンリストを取得
   */
  async getVersionList(): Promise<OcrPromptVersionInfo[]> {
    try {
      const response = await fetch(`${API_BASE_URL}/OcrPrompt/versions`, {
        method: 'GET',
        headers: {
          'Accept': 'application/json',
        },
      })

      if (!response.ok) {
        throw new Error(`API Error: ${response.status} ${response.statusText}`)
      }

      return await response.json()
    } catch (error) {
      console.error('Fetch error:', error)
      throw error
    }
  },

  /**
   * 特定バージョンのプロンプトを取得
   */
  async getPromptByVersion(version: number): Promise<ActivePrompt> {
    try {
      const response = await fetch(`${API_BASE_URL}/OcrPrompt/version/${version}`, {
        method: 'GET',
        headers: {
          'Accept': 'application/json',
        },
      })

      if (!response.ok) {
        throw new Error(`API Error: ${response.status} ${response.statusText}`)
      }

      return await response.json()
    } catch (error) {
      console.error('Fetch error:', error)
      throw error
    }
  },
}
