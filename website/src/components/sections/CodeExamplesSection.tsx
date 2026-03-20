'use client';

import React, { useState } from 'react';
import { CodeBlock } from '@/components/CodeBlock';

export function CodeExamplesSection() {
  const [activeTab, setActiveTab] = useState<'cli' | 'config' | 'output' | 'templates'>('cli');

  const examples = {
    cli: `# Add a source and compact it
$ distill add pdf ~/Books/tao-of-react.pdf \\
    --name tao-of-react --template rules --output-dir tao
✓ Added source "tao-of-react" (pdf)

$ distill tao-of-react
→ source:   tao-of-react (pdf)
→ template: rules
→ agent:    claude-cli (sonnet)
→ Compacting (1 chunk)...

✓ Compaction complete.
→ tokens:  2847 (budget: 4000)
→ output:  ./output/tao/tao-of-react-minified.md

# Re-compact everything after improving a template
$ distill update --force
→ Re-compacting 4 sources...
✓ 4 sources updated.

# List tracked sources
$ distill list
  tao-of-react    pdf     rules     tao/     2847 tok   ✓ current
  ddia            pdf     principles ddia/   -          ✗ not yet

# Install a shared context repo for agents
$ distill install https://github.com/myteam/context.git
✓ Context repo installed at ~/.claude/docs/`,
    config: `# distill.yaml
default_agent: claude-cli

agents:
  claude-cli:
    provider: claude-cli
    model: sonnet
  claude-api:
    provider: anthropic
    model: claude-sonnet-4-20250514
    api_key_env: ANTHROPIC_API_KEY
  gpt-api:
    provider: openai
    model: gpt-4o
    api_key_env: OPENAI_API_KEY

output:
  dir: ./output
  generate_indexes: true
  token_budget: 4000

sources:
  tao-of-react:
    type: pdf
    path: ~/Books/tao-of-react.pdf
    template: rules
    output_dir: tao
    output_file: tao-of-react-minified.md`,
    output: `# output directory structure
output/
├── index.md                      # root index
├── tao/
│   ├── index.md                  # "Load when working with React..."
│   ├── tao-of-react-minified.md
│   └── tao-of-node-minified.md
└── ddia/
    ├── index.md                  # "Load when making data decisions..."
    ├── ddia_01_minified.md
    ├── ddia_02_minified.md
    └── ...

# Example: tao-of-react-minified.md
# React Best Practices - Condensed Guide

## 1. Architecture

**1.1 Common Module** - Create shared module for reusable
    components, hooks, utils. Avoid bloat; split when large.

**1.2 Absolute Paths** - Use @modules/common vs ../../../.
    Configure bundler, IDE, eslint, Jest.

## Quick Decision Tree

- **Extract component?** 3+ uses or 50+ lines
- **Use reducer?** 3+ related state values`,
    templates: `# Available templates

$ distill templates
  rules        Numbered imperative rules grouped by section
  principles   Chapter-based core principles with loading guidance
  patterns     Named patterns with rationale and code examples
  raw          Minimal compaction, preserves original structure

# Custom template example (./templates/checklist.md)
---
name: checklist
description: Actionable checklist for code review
---

You are a technical knowledge compactor.
Produce a checklist of items an AI code reviewer
should verify when reviewing code in this domain.

OUTPUT FORMAT:
## {Category}
- [ ] {Check item} — {Why it matters in 1 sentence}`,
  };

  const tabs = [
    { key: 'cli' as const, label: 'CLI', language: 'bash' },
    { key: 'config' as const, label: 'Config', language: 'yaml' },
    { key: 'output' as const, label: 'Output', language: 'markdown' },
    { key: 'templates' as const, label: 'Templates', language: 'bash' },
  ];

  return (
    <section id="code-examples" className="py-12 sm:py-16 lg:py-20 bg-dark-gray/50">
      <div className="max-w-6xl mx-auto px-4 sm:px-6">
        <div className="text-center mb-10 sm:mb-16">
          <h2 className="text-3xl sm:text-4xl lg:text-5xl font-bold text-cream mb-3 sm:mb-4">Code Examples</h2>
          <p className="text-cream/70 text-base sm:text-lg lg:text-xl max-w-3xl mx-auto">See distill in action — CLI commands, configuration, and output format</p>
        </div>
        <div className="bg-dark-slate border border-distill-purple/30 rounded-xl overflow-hidden">
          <div className="flex border-b border-distill-purple/30 overflow-x-auto">
            {tabs.map((tab) => (
              <button
                key={tab.key}
                onClick={() => setActiveTab(tab.key)}
                className={`flex-1 px-3 sm:px-6 py-3 sm:py-4 text-xs sm:text-sm font-semibold transition-colors whitespace-nowrap ${
                  activeTab === tab.key
                    ? 'bg-dark-gray/50 text-distill-purple border-b-2 border-distill-purple'
                    : 'text-cream/70 hover:text-cream hover:bg-dark-gray/30'
                }`}
              >
                {tab.label}
              </button>
            ))}
          </div>
          <div className="p-4 sm:p-6 overflow-x-auto">
            <CodeBlock code={examples[activeTab]} language={tabs.find((t) => t.key === activeTab)?.language} />
          </div>
        </div>
      </div>
    </section>
  );
}
