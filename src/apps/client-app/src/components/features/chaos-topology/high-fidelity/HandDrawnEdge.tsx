'use client';

import React from 'react';

interface HandDrawnEdgeProps {
  d: string;
  isActive?: boolean;
  isError?: boolean;
  label?: string;
  markerEnd?: string;
}

export const HandDrawnEdge = ({ d, isActive, isError, label, markerEnd = 'url(#arrowhead)' }: HandDrawnEdgeProps) => {
  return (
    <g className="z-10">
      <path
        d={d}
        fill="none"
        stroke={isError ? '#ba1a1a' : (isActive ? '#004ac6' : '#737686')}
        strokeWidth={isActive ? "3" : "2"}
        className={`rough-stroke transition-all duration-500 ${isActive ? 'animate-flow-dash' : ''}`}
        markerEnd={isError ? 'url(#arrowhead-error)' : (isActive ? 'url(#arrowhead-primary)' : markerEnd)}
      />
      {label && (
        <text className="font-annotation text-xs fill-slate-500 italic">
          <textPath href={`#${label.replace(/\s+/g, '-')}`} startOffset="50%" textAnchor="middle">
            {label}
          </textPath>
        </text>
      )}
    </g>
  );
};
