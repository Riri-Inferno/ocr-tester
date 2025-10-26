export interface ActivePrompt {
  id: string
  promptContent: string
  version: number
  isDeleted: boolean
  createdAt: string
  updatedAt: string
}

export interface CreatePromptRequest {
  promptContent: string
}

export interface ApiError {
  message: string
  code?: string
}

export interface OcrPromptVersionInfo {
  id: string
  version: number
  isActive: boolean
  createdAt: string
  updatedAt: string
}
