'use client';

import React from 'react';

export interface PipelineStep {
  key: string;
  label: string;
}

interface StatusPipelineProps {
  steps: PipelineStep[];
  currentStep: string;
  size?: 'sm' | 'md';
}

export default function StatusPipeline({ steps, currentStep, size = 'md' }: StatusPipelineProps) {
  const currentIndex = steps.findIndex(s => s.key === currentStep);

  return (
    <div className="flex items-center gap-0 w-full overflow-x-auto py-2">
      {steps.map((step, i) => {
        const isCompleted = i < currentIndex;
        const isCurrent = i === currentIndex;
        const isFuture = i > currentIndex;
        const dotSize = size === 'sm' ? 'w-3 h-3' : 'w-4 h-4';
        const textSize = size === 'sm' ? 'text-[7px]' : 'text-[8px]';

        return (
          <React.Fragment key={step.key}>
            <div className="flex flex-col items-center min-w-0 flex-shrink-0">
              <div className={`${dotSize} rounded-full border-2 flex items-center justify-center transition-all ${
                isCompleted ? 'bg-green-500 border-green-500' :
                isCurrent ? 'bg-primary border-primary animate-pulse shadow-[0_0_8px_rgba(0,74,198,0.4)]' :
                'bg-transparent border-slate-300'
              }`}>
                {isCompleted && (
                  <svg className="w-2 h-2 text-white" fill="currentColor" viewBox="0 0 20 20">
                    <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                  </svg>
                )}
              </div>
              <span className={`${textSize} font-black uppercase tracking-tight mt-1 text-center leading-tight max-w-[60px] ${
                isCompleted ? 'text-green-600' :
                isCurrent ? 'text-primary' :
                'text-slate-400'
              }`}>
                {step.label}
              </span>
            </div>
            {i < steps.length - 1 && (
              <div className={`flex-1 h-0.5 min-w-[12px] mx-0.5 rounded-full ${
                i < currentIndex ? 'bg-green-400' : 'bg-slate-200'
              }`} />
            )}
          </React.Fragment>
        );
      })}
    </div>
  );
}
