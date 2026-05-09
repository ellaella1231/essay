"use client";

import { useState } from "react";
import { Camera, Upload, BookOpen } from "lucide-react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useResultStore } from "@/store/useResultStore";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://127.0.0.1:8080";

export default function HomePage() {
  const [grade, setGrade] = useState("Grade 7");
  const [promptText, setPromptText] = useState("");
  const [file, setFile] = useState<File | null>(null);
  const [isUploading, setIsUploading] = useState(false);
  const [error, setError] = useState("");
  const setResult = useResultStore((state) => state.setResult);
  const router = useRouter();

  const handleUpload = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!file) {
      return;
    }

    setIsUploading(true);
    setError("");

    const formData = new FormData();
    formData.append("image", file);
    formData.append("prompt_text", promptText);
    formData.append("grade", grade);

    try {
      const res = await fetch(`${API_BASE_URL}/essays/grade`, {
        method: "POST",
        body: formData,
      });

      const data = await res.json();
      if (!res.ok) {
        setError(data.error ?? "Upload failed");
        return;
      }

      setResult(data);
      router.push("/result/latest");
    } catch {
      setError("Unable to reach the server");
    } finally {
      setIsUploading(false);
    }
  };

  return (
    <div className="min-h-screen bg-slate-50 pb-20">
      <header className="bg-white px-6 py-4 shadow-sm sticky top-0 z-10">
        <h1 className="text-xl font-bold text-slate-900 flex items-center gap-2">
          <BookOpen className="text-blue-600" />
          AI Essay Master
        </h1>
      </header>

      <main className="px-4 py-6 max-w-lg mx-auto space-y-8">
        <section className="bg-white rounded-3xl p-6 shadow-sm border border-slate-100">
          <h2 className="text-lg font-semibold text-slate-900 mb-4">Grade Level</h2>
          <select 
            value={grade} 
            onChange={(e) => setGrade(e.target.value)}
            className="w-full bg-slate-50 border border-slate-200 text-slate-900 rounded-xl px-4 py-3 focus:ring-2 focus:ring-blue-600 outline-none"
          >
            <option>Grade 7</option>
            <option>Grade 8</option>
            <option>Grade 9</option>
            <option>High School</option>
          </select>
        </section>

        <section className="bg-white rounded-3xl p-6 shadow-sm border border-slate-100">
          <h2 className="text-lg font-semibold text-slate-900 mb-4">Essay Prompt (Optional)</h2>
          <textarea
            value={promptText}
            onChange={(e) => setPromptText(e.target.value)}
            placeholder="What is the topic of your essay?"
            className="w-full bg-slate-50 border border-slate-200 text-slate-900 rounded-xl px-4 py-3 h-24 resize-none focus:ring-2 focus:ring-blue-600 outline-none"
          />
        </section>

        <form onSubmit={handleUpload} className="space-y-6">
          <div className="relative group cursor-pointer">
            <input 
              type="file" 
              accept="image/*" 
              onChange={(e) => setFile(e.target.files?.[0] || null)}
              className="absolute inset-0 w-full h-full opacity-0 cursor-pointer z-10"
            />
            <div className={`border-2 border-dashed rounded-3xl p-10 flex flex-col items-center justify-center transition-colors ${file ? 'border-blue-500 bg-blue-50' : 'border-slate-300 bg-white hover:bg-slate-50'}`}>
              <Camera className={`h-12 w-12 mb-4 ${file ? 'text-blue-500' : 'text-slate-400'}`} />
              <span className={`text-sm font-medium ${file ? 'text-blue-700' : 'text-slate-600'}`}>
                {file ? file.name : "Tap to take a photo of your essay"}
              </span>
            </div>
          </div>

          <button
            type="submit"
            disabled={!file || isUploading}
            className="w-full flex items-center justify-center gap-2 bg-blue-600 text-white rounded-2xl py-4 font-semibold text-lg shadow-lg shadow-blue-200 disabled:opacity-50 disabled:cursor-not-allowed hover:bg-blue-700 transition-all active:scale-95"
          >
            {isUploading ? (
              <span className="animate-pulse">AI is correcting your essay...</span>
            ) : (
              <>
                <Upload className="h-5 w-5" />
                Upload & Grade
              </>
            )}
          </button>
          {error ? <p className="text-sm text-red-500 text-center">{error}</p> : null}
        </form>
      </main>

      {/* Bottom Navigation */}
      <nav className="fixed bottom-0 w-full bg-white border-t border-slate-200 px-6 py-3 flex justify-between items-center z-20">
        <Link href="/" className="flex flex-col items-center text-blue-600">
          <Camera className="h-6 w-6" />
          <span className="text-xs mt-1 font-medium">Grade</span>
        </Link>
        <Link href="/notebook" className="flex flex-col items-center text-slate-400 hover:text-slate-600">
          <BookOpen className="h-6 w-6" />
          <span className="text-xs mt-1 font-medium">Notebook</span>
        </Link>
      </nav>
    </div>
  );
}
