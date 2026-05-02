import React from "react";
import { cn } from "@/lib/utils/cn";

interface TypographyProps extends React.HTMLAttributes<HTMLHeadingElement> {
  children: React.ReactNode;
}

export const HeroTitle = ({ children, className, ...props }: TypographyProps) => (
  <h1 
    className={cn(
      "text-5xl md:text-6xl font-black font-headline tracking-tighter uppercase italic text-on-surface",
      className
    )} 
    {...props}
  >
    {children}
  </h1>
);

export const SectionHeader = ({ children, className, ...props }: TypographyProps) => (
  <h2 
    className={cn(
      "text-3xl font-black font-headline tracking-tighter uppercase italic text-on-surface",
      className
    )} 
    {...props}
  >
    {children}
  </h2>
);

export const MicroLabel = ({ children, className, ...props }: React.HTMLAttributes<HTMLParagraphElement>) => (
  <p 
    className={cn(
      "text-[10px] font-black text-on-surface-variant uppercase tracking-[0.3em] italic opacity-60",
      className
    )} 
    {...props}
  >
    {children}
  </p>
);

export const MarkerText = ({ children, className, ...props }: React.HTMLAttributes<HTMLSpanElement>) => (
  <span 
    className={cn(
      "font-marker text-3xl text-primary rotate-[-1deg]",
      className
    )} 
    {...props}
  >
    {children}
  </span>
);
