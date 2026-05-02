import React from "react";
import { cn } from "@/lib/utils/cn";

export const FlatInput = React.forwardRef<HTMLInputElement, React.InputHTMLAttributes<HTMLInputElement>>(
  ({ className, ...props }, ref) => (
    <input
      ref={ref}
      className={cn(
        "bg-transparent border-none border-b-2 border-outline-variant/30 text-on-surface font-mono text-xs py-2 focus:ring-0 focus:border-primary transition-all placeholder:text-on-surface-variant/30",
        className
      )}
      {...props}
    />
  )
);
FlatInput.displayName = "FlatInput";

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
