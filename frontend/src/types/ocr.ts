export interface TestOcrRequest {
  pdfFile: File
  usePdfDirectly: boolean
  customPrompt?: string
}

export interface TestOcrResponse {
  ocrResult: string
  promptVersion: number
  processingMethod: string
}
