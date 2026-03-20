'use client';

import { BookOpen, Code, Users, RefreshCw, Layers, FileText } from 'lucide-react';

export function UseCasesSection() {
  const useCases = [
    { icon: <BookOpen className="w-6 h-6" />, title: 'Framework Guides', description: 'Compact Tao of React, Tao of Node, or any framework guide into numbered rules your agents follow when writing code.' },
    { icon: <Layers className="w-6 h-6" />, title: 'Architecture Books', description: 'Distill DDIA, Clean Architecture, or domain-driven design books into chapter-based principles for data and system design decisions.' },
    { icon: <Code className="w-6 h-6" />, title: 'Internal Style Guides', description: 'Turn your team\'s coding conventions doc into agent-readable rules. Agents produce code that matches your patterns.' },
    { icon: <Users className="w-6 h-6" />, title: 'Shared Team Context', description: 'Publish compacted docs to a shared git repo. Every teammate\'s agent loads the same rules — consistent AI behavior across the team.' },
    { icon: <RefreshCw className="w-6 h-6" />, title: 'Living Documentation', description: 'When source material updates or you improve a template, distill update re-compacts automatically. Context stays current.' },
    { icon: <FileText className="w-6 h-6" />, title: 'Design Pattern References', description: 'Extract named patterns (Problem/Solution/Rationale) from pattern catalogs. Agents apply the right pattern for the right situation.' },
  ];

  return (
    <section id="use-cases" className="py-12 sm:py-16 lg:py-20 bg-dark-slate">
      <div className="max-w-7xl mx-auto px-4 sm:px-6">
        <div className="text-center mb-10 sm:mb-16">
          <h2 className="text-3xl sm:text-4xl lg:text-5xl font-bold text-cream mb-3 sm:mb-4">Use Cases</h2>
          <p className="text-cream/70 text-base sm:text-lg lg:text-xl max-w-3xl mx-auto">distill turns any technical knowledge into agent-ready context</p>
        </div>
        <div className="grid sm:grid-cols-2 lg:grid-cols-3 gap-4 sm:gap-6 lg:gap-8">
          {useCases.map((useCase, index) => (
            <div key={index} className="bg-dark-gray/50 border border-distill-purple/20 rounded-xl p-5 sm:p-6 hover:border-distill-lavender/40 transition-all">
              <div className="w-10 h-10 sm:w-12 sm:h-12 bg-gradient-to-br from-distill-purple to-distill-violet rounded-lg flex items-center justify-center text-white mb-3 sm:mb-4">{useCase.icon}</div>
              <h3 className="text-lg sm:text-xl font-semibold text-cream mb-2">{useCase.title}</h3>
              <p className="text-cream/60 text-sm sm:text-base leading-relaxed">{useCase.description}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
