<script setup lang="ts">
import { ref, computed } from 'vue'

const pdfFile = ref<File | null>(null)
const pdfUrl = ref<string>('')
const isLoading = ref(false)
const currentPage = ref(1)
const totalPages = ref(1)
const scale = ref(1.0)

const emit = defineEmits<{
  pdfUploaded: [file: File]
  pdfCleared: []
}>()

// PDFファイルの選択
const handleFileSelect = (event: Event) => {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  
  if (file && file.type === 'application/pdf') {
    pdfFile.value = file
    displayPdf(file)
  } else {
    alert('PDFファイルを選択してください')
  }
}

// PDFの表示
const displayPdf = (file: File) => {
  isLoading.value = true
  
  // FileReaderでPDFをDataURLに変換
  const reader = new FileReader()
  reader.onload = (e) => {
    pdfUrl.value = e.target?.result as string
    isLoading.value = false
    // PDFアップロード完了を親に通知
    emit('pdfUploaded', file)
  }
  reader.readAsDataURL(file)
}

// ファイルをドロップ
const handleDrop = (event: DragEvent) => {
  event.preventDefault()
  const file = event.dataTransfer?.files[0]
  
  if (file && file.type === 'application/pdf') {
    pdfFile.value = file
    displayPdf(file)
  }
}

const handleDragOver = (event: DragEvent) => {
  event.preventDefault()
}

// PDFをクリア
const clearPdf = () => {
  pdfFile.value = null
  pdfUrl.value = ''
  currentPage.value = 1
  totalPages.value = 1
  scale.value = 1.0
  // PDFクリアを親に通知
  emit('pdfCleared')
}

// PDFをBase64形式で取得
const getPdfAsBase64 = async (): Promise<string | null> => {
  if (!pdfFile.value) return null
  
  return new Promise((resolve) => {
    const reader = new FileReader()
    reader.onload = (e) => {
      const base64 = e.target?.result as string
      const base64Data = base64.split(',')[1]
      resolve(base64Data || null)
    }
    // 型アサーションでFileであることを明示
    reader.readAsDataURL(pdfFile.value as File)
  })
}

// PDFをFormDataで取得（ファイルアップロード用）
const getPdfAsFormData = (): FormData | null => {
  if (!pdfFile.value) return null
  
  const formData = new FormData()
  formData.append('pdfFile', pdfFile.value)
  formData.append('fileName', pdfFile.value.name)
  formData.append('fileSize', pdfFile.value.size.toString())
  
  return formData
}

// PDFファイルの存在確認
const hasPdf = computed(() => !!pdfFile.value)

// 親コンポーネントに公開
defineExpose({
  pdfFile,
  getPdfAsBase64,
  getPdfAsFormData,
  hasPdf,
})
</script>

<template>
  <div class="pdf-viewer-container">
    <div class="pdf-viewer-header">
      <h3 class="viewer-title">PDFテストビューア</h3>
      <div class="viewer-controls" v-if="pdfFile">
        <button class="control-button" @click="clearPdf">
          <span>🗑️</span> クリア
        </button>
      </div>
    </div>

    <div class="pdf-viewer-content">
      <!-- PDFアップロードエリア -->
      <div 
        v-if="!pdfFile" 
        class="upload-area"
        @drop="handleDrop"
        @dragover="handleDragOver"
      >
        <div class="upload-prompt">
          <span class="upload-icon">📄</span>
          <p class="upload-text">PDFファイルをドラッグ＆ドロップ</p>
          <p class="upload-subtext">または</p>
          <label class="upload-button">
            <input 
              type="file" 
              accept="application/pdf"
              @change="handleFileSelect"
              class="file-input"
            />
            ファイルを選択
          </label>
        </div>
      </div>

      <!-- PDF表示エリア -->
      <div v-else class="pdf-display">
        <div v-if="isLoading" class="loading">
          PDFを読み込み中...
        </div>
        
        <div v-else class="pdf-frame">
          <iframe 
            :src="pdfUrl" 
            class="pdf-iframe"
            title="PDF Viewer"
          />
        </div>

        <div class="pdf-info">
          <span class="file-name">{{ pdfFile.name }}</span>
          <span class="file-size">{{ (pdfFile.size / 1024 / 1024).toFixed(2) }} MB</span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.pdf-viewer-container {
  background: white;
  border-radius: 0.5rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  padding: 1.5rem;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.pdf-viewer-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.viewer-title {
  font-size: 1.125rem;
  font-weight: 600;
  margin: 0;
}

.viewer-controls {
  display: flex;
  gap: 0.5rem;
}

.control-button {
  background: #e9ecef;
  color: #495057;
  border: none;
  padding: 0.375rem 0.75rem;
  border-radius: 0.375rem;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  gap: 0.25rem;
}

.control-button:hover {
  background: #dee2e6;
}

.pdf-viewer-content {
  flex: 1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.upload-area {
  flex: 1;
  border: 2px dashed #dee2e6;
  border-radius: 0.5rem;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s;
  cursor: pointer;
}

.upload-area:hover {
  border-color: #667eea;
  background: #f8f9fa;
}

.upload-prompt {
  text-align: center;
}

.upload-icon {
 font-size: 3rem;
 display: block;
 margin-bottom: 1rem;
}

.upload-text {
  font-size: 1rem;
  font-weight: 500;
  margin: 0 0 0.5rem;
  color: #495057;
}

.upload-subtext {
  font-size: 0.875rem;
  color: #6c757d;
  margin: 0 0 1rem;
}

.upload-button {
  background: #667eea;
  color: white;
  padding: 0.5rem 1rem;
  border-radius: 0.375rem;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  display: inline-block;
}

.upload-button:hover {
  background: #5a67d8;
}

.file-input {
  display: none;
}

.pdf-display {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.loading {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #6c757d;
}

.pdf-frame {
  flex: 1;
  border: 1px solid #dee2e6;
  border-radius: 0.375rem;
  overflow: hidden;
}

.pdf-iframe {
  width: 100%;
  height: 100%;
  border: none;
}

.pdf-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 0.75rem;
  padding: 0.5rem;
  background: #f8f9fa;
  border-radius: 0.375rem;
  font-size: 0.75rem;
  color: #6c757d;
}

.file-name {
  font-weight: 500;
  max-width: 70%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.file-size {
  flex-shrink: 0;
}
</style>
