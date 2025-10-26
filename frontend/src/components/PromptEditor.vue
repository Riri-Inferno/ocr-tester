<script setup lang="ts">
import { ref, onMounted, watch } from "vue";
import { usePromptApi } from "../composables/usePromptApi";
import { useOcrApi } from "../composables/useOcrApi";

const isInitialized = ref(false);
const props = defineProps<{
  selectedPrompt?: any;
  pdfViewerRef?: any;
  hasPdfFile?: boolean;
}>();

const emit = defineEmits<{
  reloadRequested: [];
}>();

const { loading, error, savePrompt, fetchPromptById, deletePrompt } =
  usePromptApi();
const { testOcr, ocrResult, loading: ocrLoading } = useOcrApi();

const promptName = ref("");
const promptContent = ref("");
const isSaving = ref(false);
const isDeleting = ref(false);
const displayPrompt = ref<any>(null);
const isCopying = ref(false);

// PDFダイレクト送信フラグ
const usePdfDirectly = ref(true);

// 選択されたプロンプトを監視
watch(
  () => props.selectedPrompt,
  async (newPrompt) => {
    if (newPrompt === null) {
      // プロンプトが0件の場合
      displayPrompt.value = null;
      promptName.value = "";
      promptContent.value = "";
    } else if (newPrompt) {
      if (newPrompt.id) {
        // 既存のプロンプト：IDで詳細を取得
        try {
          const fullPrompt = await fetchPromptById(newPrompt.id);
          displayPrompt.value = fullPrompt;
          promptName.value = fullPrompt.name;
          promptContent.value = fullPrompt.promptContent;
        } catch (error) {
          console.error("Failed to fetch prompt details:", error);
        }
      } else {
        // 新規作成：空のプロンプトをそのまま使用
        displayPrompt.value = null; // 新規なのでnull
        promptName.value = "";
        promptContent.value = "";
      }
    }
  }
);

// コンポーネントマウント時の処理
onMounted(() => {
  isInitialized.value = true;
  // 初回マウント時、既に選択されたプロンプトがあれば反映
  if (props.selectedPrompt) {
    displayPrompt.value = props.selectedPrompt;
    promptName.value = props.selectedPrompt.name;
    promptContent.value = props.selectedPrompt.promptContent;
  }
});

// 再読み込み
const reloadPrompt = async () => {
  if (displayPrompt.value?.id) {
    try {
      const prompt = await fetchPromptById(displayPrompt.value.id);
      displayPrompt.value = prompt;
      promptName.value = prompt.name;
      promptContent.value = prompt.promptContent;
    } catch {
      // エラーはcomposable内で処理
    }
  }
  emit("reloadRequested");
};

// 更新（既存のプロンプトを上書き保存）
const handleUpdatePrompt = async () => {
  if (!displayPrompt.value) return;

  isSaving.value = true;
  try {
    const result = await savePrompt(
      promptName.value,
      promptContent.value,
      displayPrompt.value.id
    );
    if (result) {
      alert("プロンプトを更新しました");
      emit("reloadRequested");
    }
  } catch {
    alert("プロンプトの更新に失敗しました");
  } finally {
    isSaving.value = false;
  }
};

// 新しいプロンプトとして保存
const handleSaveAsNewPrompt = async () => {
  isSaving.value = true;
  try {
    const result = await savePrompt(promptName.value, promptContent.value);
    if (result) {
      alert("新しいプロンプトとして保存しました");
      emit("reloadRequested");
    }
  } catch {
    alert("プロンプトの保存に失敗しました");
  } finally {
    isSaving.value = false;
  }
};

// プロンプトを削除
const handleDeletePrompt = async () => {
  if (!displayPrompt.value) return;

  if (confirm(`「${displayPrompt.value.name}」を削除しますか？`)) {
    isDeleting.value = true;
    try {
      await deletePrompt(displayPrompt.value.id);
      alert("プロンプトを削除しました");
      displayPrompt.value = null;
      promptName.value = "";
      promptContent.value = "";
      emit("reloadRequested");
    } catch {
      alert("プロンプトの削除に失敗しました");
    } finally {
      isDeleting.value = false;
    }
  }
};

// プロンプトをコピー
const copyPrompt = async () => {
  isCopying.value = true;
  try {
    await navigator.clipboard.writeText(promptContent.value);
    const button = document.querySelector(".copy-button");
    if (button) {
      const originalHTML = button.innerHTML;
      button.innerHTML = "<span>✅</span> コピーしました";
      setTimeout(() => {
        button.innerHTML = originalHTML;
      }, 2000);
    }
  } catch (err) {
    alert("コピーに失敗しました");
  } finally {
    isCopying.value = false;
  }
};

// OCRテスト実行
const runOcrTest = async () => {
  if (!props.pdfViewerRef?.hasPdf || !promptContent.value) {
    alert("PDFファイルとプロンプトを入力してください");
    return;
  }

  const pdfFile = props.pdfViewerRef.pdfFile;

  try {
    await testOcr({
      pdfFile,
      usePdfDirectly: usePdfDirectly.value,
      customPrompt: promptContent.value,
    });
  } catch (error) {
    console.error("OCRテストエラー:", error);
  }
};

// OCR結果をクリア
const clearOcrResult = () => {
  ocrResult.value = null;
};

// プロンプトが変更されたらOCR結果をクリア
watch(
  () => promptContent.value,
  () => {
    clearOcrResult();
  }
);

// OCR結果をコピー
const copyOcrResult = async () => {
  if (!ocrResult.value?.ocrResult) return;

  try {
    await navigator.clipboard.writeText(ocrResult.value.ocrResult);
    const button = document.querySelector(".ocr-copy-button");
    if (button) {
      const originalHTML = button.innerHTML;
      button.innerHTML = "<span>✅</span> コピーしました";
      setTimeout(() => {
        button.innerHTML = originalHTML;
      }, 2000);
    }
  } catch (err) {
    alert("コピーに失敗しました");
  }
};
</script>

<template>
  <div class="editor-container">
    <div class="editor-header">
      <h2 class="editor-title">プロンプト編集</h2>
      <div class="action-buttons">
        <button
          class="button button-secondary"
          @click="reloadPrompt"
          :disabled="loading"
        >
          <span>🔄</span> {{ loading ? "読み込み中..." : "再読み込み" }}
        </button>
        <button
          class="button button-danger"
          @click="handleDeletePrompt"
          :disabled="loading || isDeleting || !displayPrompt"
          v-if="displayPrompt"
        >
          <span>🗑️</span> {{ isDeleting ? "削除中..." : "削除" }}
        </button>
        <button
          class="button button-primary"
          @click="handleUpdatePrompt"
          :disabled="loading || isSaving || !displayPrompt"
          v-if="displayPrompt"
        >
          <span>💾</span>
          {{ isSaving ? "更新中..." : "更新" }}
        </button>
        <button
          class="button button-secondary"
          @click="handleSaveAsNewPrompt"
          :disabled="loading || isSaving || !promptName || !promptContent"
        >
          <span>➕</span>
          {{ isSaving ? "保存中..." : "新しいプロンプトとして保存" }}
        </button>
      </div>
    </div>

    <div class="editor-content">
      <!-- エラー表示 -->
      <div v-if="error" class="error-message">⚠️ {{ error.message }}</div>

      <!-- ローディング表示 -->
      <div v-if="loading" class="loading">読み込み中...</div>

      <!-- プロンプト編集エリア -->
      <div v-else>
        <!-- プロンプト情報 -->
        <div v-if="displayPrompt">
          <div class="form-group">
            <label class="form-label">プロンプトID</label>
            <input
              type="text"
              class="form-input"
              :value="displayPrompt.id"
              readonly
            />
          </div>
        </div>

        <!-- プロンプト名 -->
        <div class="form-group">
          <label class="form-label">プロンプト名</label>
          <input
            type="text"
            class="form-input"
            v-model="promptName"
            placeholder="プロンプト名を入力してください..."
          />
        </div>

        <!-- プロンプト内容 -->
        <div class="form-group">
          <label class="form-label">プロンプト内容</label>
          <textarea
            v-model="promptContent"
            class="form-textarea"
            rows="20"
            placeholder="プロンプトを入力してください..."
          ></textarea>

          <!-- 処理方式トグル -->
          <div class="processing-options">
            <label class="toggle-label">
              処理方式:
              <div class="toggle-container">
                <span :class="['toggle-option', { active: !usePdfDirectly }]"
                  >画像変換</span
                >
                <label class="toggle-switch">
                  <input
                    type="checkbox"
                    v-model="usePdfDirectly"
                    class="toggle-input"
                  />
                  <span class="toggle-slider"></span>
                </label>
                <span :class="['toggle-option', { active: usePdfDirectly }]"
                  >PDF直送</span
                >
              </div>
            </label>
          </div>

          <div class="textarea-actions">
            <button
              class="button button-secondary copy-button"
              @click="copyPrompt"
              :disabled="loading || isCopying || !promptContent"
            >
              <span>📋</span> コピー
            </button>
            <button
              class="button button-ocr"
              @click="runOcrTest"
              :disabled="ocrLoading || !props.pdfViewerRef?.hasPdf"
              :title="
                !props.pdfViewerRef?.hasPdf
                  ? 'PDFファイルを選択してください'
                  : ''
              "
            >
              <span>🚀</span>
              {{ ocrLoading ? "OCR実行中..." : "OCRテスト実行" }}
            </button>
          </div>
        </div>

        <!-- メタデータ -->
        <div v-if="displayPrompt" class="metadata">
          <p>
            作成日時:
            {{ new Date(displayPrompt.createdAt).toLocaleString("ja-JP") }}
          </p>
          <p>
            更新日時:
            {{ new Date(displayPrompt.updatedAt).toLocaleString("ja-JP") }}
          </p>
        </div>

        <!-- OCR結果表示エリア（変更なし） -->
        <div v-if="ocrResult" class="ocr-result-container">
          <div class="ocr-result-header">
            <h4>OCR結果</h4>
            <button class="close-button" @click="clearOcrResult" title="閉じる">
              <span>✕</span>
            </button>
          </div>
          <div class="ocr-result-info">
            <span class="info-item">
              <strong>処理方式:</strong>
              {{
                ocrResult.processingMethod === "PDFDirectly"
                  ? "PDF直送"
                  : "画像変換"
              }}
            </span>
          </div>
          <pre class="ocr-result-content">{{ ocrResult.ocrResult }}</pre>
          <div class="ocr-result-actions">
            <button
              class="button button-secondary ocr-copy-button"
              @click="copyOcrResult"
            >
              <span>📋</span> コピー
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.editor-container {
  background: white;
  border-radius: 0.5rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  padding: 1.5rem;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.editor-header {
  display: flex;
  justify-content: flex-start;
  align-items: center;
  margin-bottom: 1.5rem;
  gap: 1rem;
}

.editor-title {
  font-size: 1.25rem;
  font-weight: 600;
  margin: 0;
  flex: 1;
}

.action-buttons {
  display: flex;
  gap: 0.5rem;
  align-items: flex-start;
  white-space: nowrap;
}

.save-button-wrapper {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  position: relative;
}

.button-caption {
  font-size: 0.7rem;
  color: #6c757d;
  margin-top: 0.25rem;
  position: absolute;
  top: 100%;
  right: 0;
  white-space: nowrap;
}

.editor-content {
  flex: 1;
  overflow-y: auto;
}

.editor-content::-webkit-scrollbar {
  display: none;
}

.button {
  padding: 0.5rem 1rem;
  border-radius: 0.375rem;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  border: none;
  display: flex;
  align-items: center;
  gap: 0.375rem;
  white-space: nowrap;
}

.button-primary {
  background: #667eea;
  color: white;
}

.button-primary:hover:not(:disabled) {
  background: #5a67d8;
}

.button-secondary {
  background: #e9ecef;
  color: #495057;
}

.button-secondary:hover:not(:disabled) {
  background: #dee2e6;
}

.button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.copy-button {
  width: auto;
}

.textarea-actions {
  display: flex;
  gap: 0.5rem;
  margin-top: 0.5rem;
  margin-left: 0.5rem;
}

.button-ocr {
  background: #17a2b8;
  color: white;
}

.button-ocr:hover:not(:disabled) {
  background: #138496;
}

.ocr-result-container {
  margin-top: 1rem;
  border: 1px solid #dee2e6;
  border-radius: 0.375rem;
  background: #f8f9fa;
  overflow: hidden;
}

.ocr-result-actions {
  padding: 0.5rem 1rem;
  background: white;
  border-top: 1px solid #dee2e6;
}

.ocr-result-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.75rem 1rem;
  border-bottom: 1px solid #dee2e6;
  background: white;
}

.ocr-result-header h4 {
  margin: 0;
  font-size: 1rem;
  font-weight: 600;
}

.close-button {
  background: none;
  border: none;
  font-size: 1.25rem;
  color: #6c757d;
  cursor: pointer;
  padding: 0;
  width: 1.5rem;
  height: 1.5rem;
  display: flex;
  align-items: center;
  justify-content: center;
}

.close-button:hover {
  color: #495057;
}

.ocr-result-info {
  padding: 0.75rem 1rem;
  display: flex;
  gap: 1.5rem;
  font-size: 0.875rem;
  color: #6c757d;
  background: #e9ecef;
}

.info-item {
  display: flex;
  gap: 0.5rem;
}

.ocr-result-content {
  padding: 1rem;
  margin: 0;
  font-size: 0.875rem;
  line-height: 1.5;
  white-space: pre-wrap;
  word-wrap: break-word;
  font-family: "Consolas", "Monaco", monospace;
  max-height: 400px;
  overflow-y: auto;
  background: white;
}

.form-group {
  margin-bottom: 1.5rem;
}

.form-label {
  display: block;
  font-size: 0.875rem;
  font-weight: 500;
  margin-bottom: 0.5rem;
  color: #495057;
}

.form-input,
.form-textarea {
  width: 100%;
  padding: 0.5rem 0.75rem;
  border: 1px solid #ced4da;
  border-radius: 0.375rem;
  font-size: 0.875rem;
  font-family: inherit;
  transition: border-color 0.2s;
}

.form-input:focus,
.form-textarea:focus:not([readonly]) {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.form-textarea {
  resize: vertical;
  line-height: 1.5;
}

.form-input[readonly] {
  background-color: #f8f9fa;
  cursor: not-allowed;
}

.version-info {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.version-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 0.75rem;
  font-weight: 600;
  border-radius: 9999px;
  padding: 0.25rem 0.75rem;
  min-width: 2.5rem;
  min-height: 2rem;
  line-height: 1;
  white-space: nowrap;
}

.version-badge.active {
  background-color: #28a745;
  color: #fff;
}

.version-badge.inactive {
  background-color: #6c757d;
  color: #fff;
}

.error-message {
  background-color: #fee;
  color: #c33;
  padding: 0.75rem;
  border-radius: 0.375rem;
  margin-bottom: 1rem;
  font-size: 0.875rem;
}

.loading {
  text-align: center;
  padding: 2rem;
  color: #666;
}

.metadata {
  margin-top: 1.5rem;
  padding-top: 0rem;
  border-bottom: 1px solid #e9ecef;
  font-size: 0.75rem;
  color: #6c757d;
}

.metadata p {
  margin: 0.25rem 0;
}

.processing-options {
  margin-top: 0.75rem;
  margin-bottom: 0.75rem;
  padding: 0.75rem;
  margin-left: 0.5rem;
  background: #f8f9fa;
  border-radius: 0.375rem;
}

.toggle-label {
  display: flex;
  align-items: center;
  gap: 1rem;
  font-size: 0.875rem;
  font-weight: 500;
  color: #495057;
}

.toggle-container {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.toggle-option {
  font-size: 0.875rem;
  color: #6c757d;
  transition: color 0.2s;
}

.toggle-option.active {
  color: #495057;
  font-weight: 600;
}

.toggle-switch {
  position: relative;
  display: inline-block;
  width: 48px;
  height: 24px;
}

.toggle-input {
  opacity: 0;
  width: 0;
  height: 0;
}

.toggle-slider {
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: #ccc;
  transition: 0.3s;
  border-radius: 24px;
}

.toggle-slider:before {
  position: absolute;
  content: "";
  height: 18px;
  width: 18px;
  left: 3px;
  bottom: 3px;
  background-color: white;
  transition: 0.3s;
  border-radius: 50%;
}

.toggle-input:checked + .toggle-slider {
  background-color: #17a2b8;
}

.toggle-input:checked + .toggle-slider:before {
  transform: translateX(24px);
}

.toggle-slider:hover {
  background-color: #a0a0a0;
}

.toggle-input:checked + .toggle-slider:hover {
  background-color: #138496;
}

.button-danger {
  background: #dc3545;
  color: white;
}

.button-danger:hover:not(:disabled) {
  background: #c82333;
}
</style>
