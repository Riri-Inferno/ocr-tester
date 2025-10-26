import type { PromptInfo, Prompt } from "../types/prompt";

// === モックデータ（一覧用） ===
export const mockPrompts: PromptInfo[] = [
  {
    id: "1",
    name: "請求書OCR",
    createdAt: "2024-01-01T10:15:00Z",
    updatedAt: "2024-03-12T08:30:00Z",
  },
  {
    id: "2",
    name: "領収書OCR",
    createdAt: "2024-01-02T09:00:00Z",
    updatedAt: "2024-03-15T11:45:00Z",
  },
  {
    id: "3",
    name: "契約書OCR",
    createdAt: "2024-01-03T13:20:00Z",
    updatedAt: "2024-03-18T16:10:00Z",
  },
  {
    id: "4",
    name: "見積書OCR",
    createdAt: "2024-01-05T14:10:00Z",
    updatedAt: "2024-03-22T09:00:00Z",
  },
];

// === モックデータ（詳細用） ===
export const mockPromptDetails: Record<string, Prompt> = {
  "1": {
    id: "1",
    name: "請求書OCR",
    promptContent: `請求書の内容を構造化データとして抽出してください。
- 会社名
- 請求日
- 合計金額
- 明細項目（品目名、単価、数量）`,
    isDeleted: false,
    createdAt: "2024-01-01T10:15:00Z",
    updatedAt: "2024-03-12T08:30:00Z",
  },
  "2": {
    id: "2",
    name: "領収書OCR",
    promptContent: `領収書から以下の項目を抽出してください。
- 支払先名
- 支払金額（税込）
- 日付
- 支払方法`,
    isDeleted: false,
    createdAt: "2024-01-02T09:00:00Z",
    updatedAt: "2024-03-15T11:45:00Z",
  },
  "3": {
    id: "3",
    name: "契約書OCR",
    promptContent: `契約書から以下の項目を抽出してください。
- 契約当事者
- 契約日
- 契約期間
- 契約金額
- 契約条件の要約`,
    isDeleted: false,
    createdAt: "2024-01-03T13:20:00Z",
    updatedAt: "2024-03-18T16:10:00Z",
  },
  "4": {
    id: "4",
    name: "見積書OCR",
    promptContent: `見積書から以下を抽出してください。
- 発行元会社名
- 見積日
- 見積金額
- 有効期限
- 明細項目`,
    isDeleted: false,
    createdAt: "2024-01-05T14:10:00Z",
    updatedAt: "2024-03-22T09:00:00Z",
  },
};
