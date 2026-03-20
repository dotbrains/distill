'use client';

import { ExternalLink, BookOpen, Layers, Hash, Cpu, FolderTree } from 'lucide-react';

const contributions = [
  { icon: <Layers className="w-5 h-5" />, title: 'Template-Driven Compaction', description: 'Four built-in output formats enforce consistent structure across heterogeneous source material. Custom templates supported.' },
  { icon: <FolderTree className="w-5 h-5" />, title: 'Multi-Source Ingestion', description: 'Reads from PDFs, markdown, Notion, web URLs, EPUBs, and GitHub repos — all normalized into uniform text chunks.' },
  { icon: <Hash className="w-5 h-5" />, title: 'Content-Hashing State Tracker', description: 'SHA-256 tracks source, template, and agent identity. Changed inputs trigger re-compaction; unchanged sources are skipped.' },
  { icon: <Cpu className="w-5 h-5" />, title: 'Pluggable Agent Architecture', description: 'CLI-based (Claude, Codex) and API-based (Anthropic, OpenAI) providers with the same registry pattern as prr.' },
];

export function PaperSection() {
  return (
    <section id="paper" className="py-12 sm:py-16 lg:py-20 bg-dark-gray/50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6">
        <div className="text-center mb-10 sm:mb-16">
          <div className="inline-flex items-center gap-2 px-3 py-1 bg-distill-purple/10 border border-distill-purple/30 rounded-full text-distill-purple text-xs font-medium mb-4">Technical Paper</div>
          <h2 className="text-3xl sm:text-4xl lg:text-5xl font-bold text-cream mb-3 sm:mb-4">The Design Behind distill</h2>
          <p className="text-cream/70 text-base sm:text-lg lg:text-xl max-w-3xl mx-auto">A deep-dive into the four technical contributions that make repeatable knowledge compaction possible</p>
        </div>

        <div className="grid lg:grid-cols-2 gap-8 lg:gap-12 items-start">
          <div className="bg-dark-slate border border-distill-purple/20 rounded-2xl p-8 sm:p-10">
            <div className="flex items-start gap-4 mb-6">
              <div className="w-12 h-12 bg-gradient-to-br from-distill-purple to-distill-violet rounded-lg flex items-center justify-center text-white flex-shrink-0">
                <BookOpen className="w-6 h-6" />
              </div>
              <div>
                <p className="text-distill-purple text-xs font-medium uppercase tracking-wider mb-1">Technical Paper</p>
                <h3 className="text-cream font-bold text-lg sm:text-xl leading-snug">distill: A Template-Driven Knowledge Compaction Pipeline for AI Agents</h3>
              </div>
            </div>
            <p className="text-cream/50 text-xs mb-5">Nicholas Adamou — dotbrains</p>
            <p className="text-cream/70 text-sm leading-relaxed mb-8">
              AI agents are constrained by their context windows. Technical books contain critical guidance but at 200-600 pages they are too verbose to load. This paper presents distill and the four design decisions that make knowledge compaction repeatable, consistent, and cheap to maintain.
            </p>
            <a
              href="https://github.com/dotbrains/distill/blob/main/PAPER.md"
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center gap-2 bg-gradient-to-r from-distill-purple to-distill-violet hover:from-distill-violet hover:to-distill-lavender text-white px-6 py-3 rounded-lg shadow-lg shadow-distill-purple/30 text-sm font-semibold transition-all"
            >
              Read the Paper
              <ExternalLink className="w-4 h-4" />
            </a>
          </div>

          <div className="space-y-4">
            {contributions.map((c, i) => (
              <div key={i} className="flex gap-4 bg-dark-slate border border-distill-purple/10 hover:border-distill-lavender/30 rounded-xl p-4 sm:p-5 transition-all group">
                <div className="w-9 h-9 bg-gradient-to-br from-distill-purple to-distill-violet rounded-lg flex items-center justify-center text-white flex-shrink-0 group-hover:scale-110 transition-transform">{c.icon}</div>
                <div>
                  <div className="flex items-center gap-2 mb-1">
                    <span className="text-distill-purple text-xs font-medium">({i + 1})</span>
                    <h4 className="text-cream font-semibold text-sm sm:text-base">{c.title}</h4>
                  </div>
                  <p className="text-cream/60 text-xs sm:text-sm leading-relaxed">{c.description}</p>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </section>
  );
}
