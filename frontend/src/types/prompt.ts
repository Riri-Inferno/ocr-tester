export interface Prompt {
  id: string;
  name: string;
  promptContent: string;
  version?: number;
  createdAt: string;
  updatedAt: string;
  userId?: string;
}

export interface CreatePromptRequest {
  name: string;
  promptContent: string;
  userId?: string;
}

export interface ApiError {
  message: string;
  code?: string;
}

export interface PromptInfo {
  id: string;
  name: string;
  version?: number;
  createdAt: string;
  updatedAt: string;
  userId?: string;
}
