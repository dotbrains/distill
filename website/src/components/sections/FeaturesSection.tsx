'use client';

import { FileText, RefreshCw, Layers, Hash, BookOpen, Terminal, Cpu, FolderTree, Settings } from 'lucide-react';

export function FeaturesSection() {
  const features = [
    {
      icon: <Layers className="w-6 h-6" />,
      title: 'Template-Driven Compaction',
      description: 'Four built-in templates (rules, principles, patterns, raw) enforce consistent output structure. Custom templates supported via markdown files.',
    },
    {
      icon: <FileText className="w-6 h-6" />,
      title: 'Multi-Source Ingestion',
      description: 'Read from PDFs, markdown files, Notion pages, web URLs, EPUBs, and GitHub repos. All normalized into uniform chunks.',
    },
    {
      icon: <Hash className="w-6 h-6" />,
      title: 'Content Hashing',
      description: 'SHA-256 tracks source content, template, and agent identity. Changed sources are re-compacted automatically. Unchanged sources are skipped.',
    },
    {
      icon: <Cpu className="w-6 h-6" />,
      title: 'Pluggable Agents',
      description: 'Claude CLI, Codex CLI, Anthropic API, and OpenAI API. Same interface, swap with --agent. Claude CLI is zero-config default.',
    },
    {
      icon: <FolderTree className="w-6 h-6" />,
      title: 'Hierarchical Indexes',
      description: 'Auto-generated index.md files at every level. Agents know what to load and when — no token waste from loading everything.',
    },
    {
      icon: <RefreshCw className="w-6 h-6" />,
      title: 'Incremental Updates',
      description: 'distill update skips clean sources. Improve a template? All sources using it are automatically re-compacted.',
    },
    {
      icon: <BookOpen className="w-6 h-6" />,
      title: 'Chapter Splitting',
      description: 'Split books by chapter with split_by: chapter. Each chapter becomes a separate output file with its own entry in the index.',
    },
    {
      icon: <Terminal className="w-6 h-6" />,
      title: 'Context Repo Publishing',
      description: 'distill publish copies output to a git repo and commits. Teams clone into ~/.claude/docs/ for shared agent context.',
    },
    {
      icon: <Settings className="w-6 h-6" />,
      title: 'Single Binary CLI',
      description: 'Written in Go with Cobra. Cross-compiles to macOS and Linux. Install via go install, Homebrew, or GitHub Release.',
    },
  ];

  return (
    <section id="features" className="py-12 sm:py-16 lg:py-20 bg-dark-slate">
      <div className="max-w-7xl mx-auto px-4 sm:px-6">
        <div className="text-center mb-10 sm:mb-16">
          <h2 className="text-3xl sm:text-4xl lg:text-5xl font-bold text-cream mb-3 sm:mb-4">
            Built for Agent-Driven Development
          </h2>
          <p className="text-cream/70 text-base sm:text-lg lg:text-xl max-w-3xl mx-auto">
            Every feature exists to make agent context production repeatable and maintainable
          </p>
        </div>
        <div className="grid sm:grid-cols-2 lg:grid-cols-3 gap-4 sm:gap-6 lg:gap-8">
          {features.map((feature, index) => (
            <div
              key={index}
              className="group bg-dark-gray/50 border border-distill-purple/20 hover:border-distill-lavender/40 rounded-xl p-5 sm:p-6 transition-all hover:shadow-lg hover:shadow-distill-purple/10"
            >
              <div className="w-10 h-10 sm:w-12 sm:h-12 bg-gradient-to-br from-distill-purple to-distill-violet rounded-lg flex items-center justify-center text-white mb-3 sm:mb-4 group-hover:scale-110 transition-transform">
                {feature.icon}
              </div>
              <h3 className="text-lg sm:text-xl font-semibold text-cream mb-2">{feature.title}</h3>
              <p className="text-cream/60 text-sm sm:text-base leading-relaxed">{feature.description}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
