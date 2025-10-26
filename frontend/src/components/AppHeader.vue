<template>
  <header class="app-header">
    <div class="header-content">
      <div class="header-main">
        <img src="/icons/logo.svg" alt="ロゴ" class="logo" />
        <div>
          <h1 class="header-title">OCR統合開発者向けツール</h1>
          <p class="header-subtitle">OCR用プロンプト編集ツール</p>
        </div>
      </div>
    </div>
  </header>
</template>

<style scoped>
.app-header {
  background: linear-gradient(135deg, #a8e063 10%, #56ab2f 100%);
  color: white;
  padding: 1.5rem 2rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.header-content {
  max-width: 1400px;
  margin: 0;
}

.header-main {
  display: flex;
  align-items: center;
  gap: 1.5rem;
}

.logo {
  height: 3rem;
  width: auto;
  /* filter: brightness(0) invert(1); */
}

.header-title {
  font-size: 1.75rem;
  font-weight: 700;
  margin: 0;
}

.header-subtitle {
  margin: 0.25rem 0 0;
  opacity: 0.9;
  font-size: 0.875rem;
}
</style>
