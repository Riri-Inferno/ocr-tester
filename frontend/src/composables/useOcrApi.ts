import { ref } from 'vue'
import type{ Ref } from 'vue'
import { ocrRepository } from '../repositories/ocrRepository'
import type { TestOcrRequest, TestOcrResponse } from '../types/ocr'

export function useOcrApi() {
  const loading: Ref<boolean> = ref(false)
  const error: Ref<Error | null> = ref(null)
  const ocrResult: Ref<TestOcrResponse | null> = ref(null)

  /**
   * PDFファイルをOCRテスト処理
   */
  const testOcr = async (request: TestOcrRequest) => {
    loading.value = true
    error.value = null
    ocrResult.value = null
    
    try {
      const data = await ocrRepository.testOcr(request)
      ocrResult.value = data
      return data
    } catch (err) {
      error.value = err instanceof Error ? err : new Error('OCRテストに失敗しました')
      throw err
    } finally {
      loading.value = false
    }
  }

  return {
    loading,
    error,
    ocrResult,
    testOcr,
  }
}
