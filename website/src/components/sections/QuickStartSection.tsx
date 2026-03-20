'use client';

import React, { useState } from 'react';
import { CodeBlock } from '@/components/CodeBlock';

export function QuickStartSection() {
  const [installMethod, setInstallMethod] = useState<'go' | 'brew' | 'release'>('go');

  const installExamples = {
    go: `go install github.com/dotbrains/distill@latest`,
    brew: `brew tap dotbrains/tap\nbrew install --cask distill`,
    release: `gh release download --repo dotbrains/distill \\\n  --pattern 'distill_darwin_arm64.tar.gz' --dir /tmp\ntar -xzf /tmp/distill_darwin_arm64.tar.gz -C /usr/local/bin`,
  };

  return (
    <section id="quick-start" className="py-12 sm:py-16 lg:py-20 bg-dark-slate overflow-hidden">
      <div className="max-w-7xl mx-auto px-4 sm:px-6">
        <div className="text-center mb-10 sm:mb-16">
          <h2 className="text-3xl sm:text-4xl lg:text-5xl font-bold text-cream mb-3 sm:mb-4">Quick Start</h2>
          <p className="text-slate-gray text-base sm:text-lg lg:text-xl max-w-3xl mx-auto">Install distill and compact your first source in under a minute</p>
        </div>
        <div className="grid lg:grid-cols-2 gap-8 lg:gap-12 items-start">
          <div className="bg-dark-gray/50 rounded-xl p-6 sm:p-8 border border-distill-purple/20 min-w-0">
            <h3 className="text-xl sm:text-2xl font-bold text-cream mb-4 sm:mb-6">1. Install</h3>
            <div className="flex gap-2 sm:gap-3 mb-6">
              {[
                { key: 'go' as const, label: 'go install' },
                { key: 'brew' as const, label: 'Homebrew' },
                { key: 'release' as const, label: 'Release' },
              ].map((method) => (
                <button
                  key={method.key}
                  onClick={() => setInstallMethod(method.key)}
                  className={`flex-1 px-3 sm:px-4 py-2.5 rounded-lg text-sm font-semibold transition-all ${
                    installMethod === method.key
                      ? 'bg-gradient-to-r from-distill-purple to-distill-violet text-white shadow-lg shadow-distill-purple/30'
                      : 'bg-dark-slate text-slate-gray hover:text-cream hover:border-distill-purple/50 border border-distill-purple/30'
                  }`}
                >
                  {method.label}
                </button>
              ))}
            </div>
            <CodeBlock code={installExamples[installMethod]} language="bash" />
          </div>
          <div className="bg-dark-gray/50 rounded-xl p-6 sm:p-8 border border-distill-violet/20 min-w-0">
            <h3 className="text-xl sm:text-2xl font-bold text-cream mb-4 sm:mb-6">2. Compact</h3>
            <CodeBlock
              code={`# Initialize config
distill config init

# Add a source
distill add pdf ~/Books/tao-of-react.pdf \\
  --name tao-of-react --template rules

# Compact it
distill tao-of-react

# Check agents
distill agents`}
              language="bash"
            />
            <div className="mt-6 bg-distill-purple/10 border border-distill-purple/30 rounded-lg p-4 sm:p-5">
              <p className="text-cream text-sm leading-relaxed">
                <span className="text-distill-purple font-semibold">Tip:</span> The default agent is <code className="bg-dark-slate/80 px-1.5 py-0.5 rounded text-distill-lavender font-mono text-xs">claude-cli</code> — no API key needed if you have Claude Code installed.
              </p>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
