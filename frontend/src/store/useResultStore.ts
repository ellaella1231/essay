import { create } from "zustand";
import { createJSONStorage, persist } from "zustand/middleware";

export interface EssayErrorDetail {
  original_segment: string;
  suggested_segment: string;
  explanation: string;
}

export interface EssayResult {
  score: number;
  perfect_essay: string;
  original_text?: string;
  errors: EssayErrorDetail[];
}

interface ResultState {
  result: EssayResult | null;
  setResult: (result: EssayResult) => void;
  clearResult: () => void;
}

export const useResultStore = create<ResultState>()(
  persist(
    (set) => ({
      result: null,
      setResult: (result) => set({ result }),
      clearResult: () => set({ result: null }),
    }),
    {
      name: "essay-result-store",
      storage: createJSONStorage(() => sessionStorage),
    }
  )
);
