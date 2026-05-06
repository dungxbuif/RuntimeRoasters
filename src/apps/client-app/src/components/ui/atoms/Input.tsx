import React from "react";
import { cn } from "@/lib/utils/cn";

export const Input = React.forwardRef<HTMLInputElement, React.InputHTMLAttributes<HTMLInputElement>>(
  ({ className, ...props }, ref) => (
    <input
      ref={ref}
      className={cn(
        "w-full bg-white/50 backdrop-blur-sm border border-outline-variant/30 text-on-surface font-headline text-sm px-5 py-4 rounded-2xl focus:ring-2 focus:ring-primary/20 focus:border-primary transition-all placeholder:text-on-surface-variant/30 outline-none",
        className
      )}
      {...props}
    />
  )
);
Input.displayName = "Input";

interface ToggleProps extends React.InputHTMLAttributes<HTMLInputElement> {
  activeVariant?: 'primary' | 'error';
}

export const ToggleSwitch = ({ className, activeVariant = 'primary', ...props }: ToggleProps) => {
  const activeColors = {
    primary: "peer-checked:bg-primary",
    error: "peer-checked:bg-error"
  };

  return (
    <label className={cn("relative inline-flex items-center cursor-pointer", className)}>
      <input type="checkbox" className="sr-only peer" {...props} />
      <div className={cn(
        "w-11 h-6 bg-surface-container-high rounded-full peer shadow-inner transition-all after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all after:shadow-sm peer-checked:after:translate-x-full peer-checked:after:border-white",
        activeColors[activeVariant]
      )}></div>
    </label>
  );
};
