import type { TestOcrRequest, TestOcrResponse } from '../types/ocr'

const API_BASE_URL = import.meta.env.DEV 
  ? '/api' 
  : import.meta.env.VITE_API_TARGET_URL

export const ocrRepository = {
  /**
   * PDFファイルをOCRテスト処理
   */
  async testOcr(request: TestOcrRequest): Promise<TestOcrResponse> {
    const formData = new FormData()
    formData.append('PdfFile', request.pdfFile)
    formData.append('UsePdfDirectly', request.usePdfDirectly.toString())
    if (request.customPrompt) {
      formData.append('CustomPrompt', request.customPrompt)
    }

    try {
      const response = await fetch(`${API_BASE_URL}/Document/test-ocr`, {
        method: 'POST',
        headers: {
          'Accept': 'application/json',
        },
        body: formData,
      })

      if (!response.ok) {
        throw new Error(`API Error: ${response.status} ${response.statusText}`)
      }

      return await response.json()
    } catch (error) {
      console.error('OCR test error:', error)
      throw error
    }
  },
}
