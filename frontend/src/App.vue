<script setup lang="ts">
import { ref } from "vue";
import AppHeader from "./components/AppHeader.vue";
import PromptEditor from "./components/PromptEditor.vue";
import PromptList from "./components/PromptList.vue";
import PdfViewer from "./components/PdfViewer.vue";

const pdfViewerRef = ref<InstanceType<typeof PdfViewer> | null>(null);
const selectedPrompt = ref<any>(null);
const promptListRef = ref<any>(null);
const hasPdfFile = ref(false);

const handlePromptSelected = (prompt: any) => {
  selectedPrompt.value = prompt;
};

const handleReloadRequested = () => {
  // PromptListの再読み込みを実行
  if (promptListRef.value) {
    promptListRef.value.reload();
  }
};

// PDFアップロード時の処理
const handlePdfUploaded = (file: File) => {
  console.log("PDF uploaded:", file.name);
  hasPdfFile.value = true;
};

// PDFクリア時の処理
const handlePdfCleared = () => {
  console.log("PDF cleared");
  hasPdfFile.value = false;
};
</script>

<template>
  <div class="app-container">
    <AppHeader />

    <main class="main-content">
      <aside class="sidebar">
        <PromptList
          ref="promptListRef"
          @prompt-selected="handlePromptSelected"
        />
      </aside>

      <div class="content-area">
        <section class="editor-section">
          <PromptEditor
            :selected-prompt="selectedPrompt"
            :pdf-viewer-ref="pdfViewerRef"
            :has-pdf-file="hasPdfFile"
            @reload-requested="handleReloadRequested"
          />
        </section>

        <section class="viewer-section">
          <PdfViewer
            ref="pdfViewerRef"
            @pdf-uploaded="handlePdfUploaded"
            @pdf-cleared="handlePdfCleared"
          />
        </section>
      </div>
    </main>
  </div>
</template>

<style scoped>
.app-container {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background-color: #f5f5f5;
  margin: 0;
  padding: 0;
}

.main-content {
  flex: 1;
  display: flex;
  gap: 1.5rem;
  padding: 1.5rem;
  margin: 0;
  overflow: hidden;
}

.sidebar {
  width: 300px;
  flex-shrink: 0;
  overflow-y: auto;
}

.content-area {
  flex: 1;
  display: flex;
  gap: 1.5rem;
  overflow: hidden;
}

.editor-section {
  flex: 1;
  overflow-y: auto;
  min-width: 0;
}

.viewer-section {
  flex: 1;
  overflow-y: auto;
  min-width: 0;
}
</style>
