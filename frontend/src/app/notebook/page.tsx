"use client";

import { useState } from "react";
import Link from "next/link";
import { Camera, BookOpen, Check } from "lucide-react";

const INITIAL_MISTAKES = [
  {
    id: 1,
    original_segment: "I very like play",
    suggested_segment: "I really enjoy playing",
    explanation: "Use 'really enjoy' or 'like' instead of 'very like'. 'Enjoy' is followed by a gerund (-ing).",
  },
  {
    id: 2,
    original_segment: "She don't know",
    suggested_segment: "She doesn't know",
    explanation: "Use 'doesn't' for third-person singular subjects (he, she, it).",
  }
];

export default function NotebookPage() {
  const [mistakes, setMistakes] = useState(INITIAL_MISTAKES);

  const markAsLearned = (id: number) => {
    // In a real app, this would be an API call to soft-delete
    setMistakes(mistakes.filter(m => m.id !== id));
  };

  return (
    <div className="min-h-screen bg-slate-50 pb-20">
      <header className="bg-white px-6 py-4 shadow-sm sticky top-0 z-10">
        <h1 className="text-xl font-bold text-slate-900 flex items-center gap-2">
          <BookOpen className="text-blue-600" />
          Mistake Notebook
        </h1>
      </header>

      <main className="px-4 py-6 max-w-lg mx-auto space-y-4">
        {mistakes.length === 0 ? (
          <div className="text-center py-20 text-slate-500">
            <div className="bg-blue-50 w-20 h-20 rounded-full flex items-center justify-center mx-auto mb-4">
              <Check className="h-10 w-10 text-blue-500" />
            </div>
            <p className="font-medium text-lg">All caught up!</p>
            <p className="text-sm mt-1">You have reviewed all your past mistakes.</p>
          </div>
        ) : (
          mistakes.map((mistake) => (
            <div key={mistake.id} className="bg-white rounded-2xl p-5 shadow-sm border border-slate-100 flex flex-col gap-4">
              <div>
                <div className="mb-2">
                  <span className="text-xs font-bold text-slate-400 uppercase tracking-wider">You wrote</span>
                  <p className="text-slate-800 font-medium line-through decoration-red-400 decoration-2">{mistake.original_segment}</p>
                </div>
                <div>
                  <span className="text-xs font-bold text-blue-600 uppercase tracking-wider">Correction</span>
                  <p className="text-slate-900 font-bold">{mistake.suggested_segment}</p>
                </div>
                <div className="mt-3 text-sm text-slate-600 bg-slate-50 p-3 rounded-xl border border-slate-100">
                  {mistake.explanation}
                </div>
              </div>
              
              <div className="flex justify-end pt-2 border-t border-slate-100">
                <button
                  onClick={() => markAsLearned(mistake.id)}
                  className="flex items-center gap-2 text-sm font-semibold text-emerald-600 hover:text-emerald-700 bg-emerald-50 px-4 py-2 rounded-xl transition-colors"
                >
                  <Check className="h-4 w-4" />
                  Got it!
                </button>
              </div>
            </div>
          ))
        )}
      </main>

      {/* Bottom Navigation */}
      <nav className="fixed bottom-0 w-full bg-white border-t border-slate-200 px-6 py-3 flex justify-between items-center z-20">
        <Link href="/" className="flex flex-col items-center text-slate-400 hover:text-slate-600">
          <Camera className="h-6 w-6" />
          <span className="text-xs mt-1 font-medium">Grade</span>
        </Link>
        <Link href="/notebook" className="flex flex-col items-center text-blue-600">
          <BookOpen className="h-6 w-6" />
          <span className="text-xs mt-1 font-medium">Notebook</span>
        </Link>
      </nav>
    </div>
  );
}
