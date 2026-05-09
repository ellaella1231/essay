"use client";

import { useState } from "react";
import Link from "next/link";
import { ChevronLeft, Sparkles, BookOpen } from "lucide-react";

const MOCK_RESULT = {
  score: 85,
  original: "I very like play football. Every day I playing with my friends. It make me happy.",
  perfect: "I really enjoy playing football. I play with my friends every day. It makes me very happy.",
  errors: [
    {
      id: 1,
      original_segment: "I very like play",
      suggested_segment: "I really enjoy playing",
      explanation: "Use 'really enjoy' or 'like' instead of 'very like'. 'Enjoy' is followed by a gerund (-ing).",
    },
    {
      id: 2,
      original_segment: "Every day I playing",
      suggested_segment: "I play with my friends every day",
      explanation: "For habitual actions, use the present simple tense 'I play'.",
    },
    {
      id: 3,
      original_segment: "It make me happy",
      suggested_segment: "It makes me very happy",
      explanation: "'It' is a third-person singular subject, so the verb needs an 's'.",
    }
  ]
};

export default function ResultPage() {
  const [showPerfect, setShowPerfect] = useState(true);

  return (
    <div className="min-h-screen bg-slate-50 pb-10">
      <header className="bg-white px-4 py-4 shadow-sm sticky top-0 z-10 flex items-center gap-4">
        <Link href="/" className="p-2 -ml-2 rounded-full hover:bg-slate-100">
          <ChevronLeft className="h-6 w-6 text-slate-700" />
        </Link>
        <h1 className="text-xl font-bold text-slate-900">Grading Result</h1>
      </header>

      <main className="px-4 py-6 max-w-lg mx-auto space-y-6">
        {/* Score Card */}
        <div className="bg-gradient-to-br from-blue-600 to-indigo-700 rounded-3xl p-8 text-center shadow-lg text-white">
          <h2 className="text-blue-100 font-medium mb-2">Final Score</h2>
          <div className="text-6xl font-black">{MOCK_RESULT.score}</div>
          <p className="mt-4 text-blue-50 flex items-center justify-center gap-2">
            <Sparkles className="h-4 w-4 text-yellow-300" />
            Great job! Keep practicing.
          </p>
        </div>

        {/* Text View Toggle */}
        <div className="bg-white rounded-3xl p-6 shadow-sm border border-slate-100">
          <div className="flex bg-slate-100 rounded-xl p-1 mb-6">
            <button
              onClick={() => setShowPerfect(false)}
              className={`flex-1 py-2 rounded-lg text-sm font-semibold transition-colors ${!showPerfect ? 'bg-white shadow text-slate-900' : 'text-slate-500'}`}
            >
              Original Text
            </button>
            <button
              onClick={() => setShowPerfect(true)}
              className={`flex-1 py-2 rounded-lg text-sm font-semibold transition-colors ${showPerfect ? 'bg-blue-600 shadow text-white' : 'text-slate-500'}`}
            >
              Perfect Version
            </button>
          </div>
          <div className="text-slate-700 leading-relaxed font-serif text-lg">
            {showPerfect ? MOCK_RESULT.perfect : MOCK_RESULT.original}
          </div>
        </div>

        {/* Error Breakdown */}
        <div>
          <h3 className="text-lg font-bold text-slate-900 mb-4 px-2 flex items-center gap-2">
            <BookOpen className="h-5 w-5 text-blue-600" />
            Detailed Corrections
          </h3>
          <div className="space-y-4">
            {MOCK_RESULT.errors.map((err) => (
              <div key={err.id} className="bg-white rounded-2xl p-5 shadow-sm border border-red-100">
                <div className="mb-3 flex items-start gap-3">
                  <span className="bg-red-100 text-red-700 text-xs font-bold px-2 py-1 rounded">Original</span>
                  <span className="text-slate-800 font-medium line-through decoration-red-400 decoration-2">{err.original_segment}</span>
                </div>
                <div className="mb-3 flex items-start gap-3">
                  <span className="bg-green-100 text-green-700 text-xs font-bold px-2 py-1 rounded">Corrected</span>
                  <span className="text-slate-900 font-semibold">{err.suggested_segment}</span>
                </div>
                <div className="mt-4 pt-4 border-t border-slate-100 text-sm text-slate-600 flex gap-2">
                  <span className="shrink-0 text-amber-500">💡</span>
                  <p>{err.explanation}</p>
                </div>
              </div>
            ))}
          </div>
        </div>
      </main>
    </div>
  );
}
