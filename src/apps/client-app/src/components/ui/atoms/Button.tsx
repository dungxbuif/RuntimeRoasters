import React from "react";
import { cn } from "@/lib/utils/cn";

export interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'ghost' | 'outline' | 'danger' | 'surface';
  size?: 'sm' | 'md' | 'lg' | 'xl';
  isLoading?: boolean;
}

export const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant = 'primary', size = 'md', isLoading, children, ...props }, ref) => {
    const variants = {
      primary: "bg-primary text-on-primary shadow-lg shadow-primary/20 hover:scale-[0.98] active:scale-95 italic",
      ghost: "bg-transparent text-on-surface-variant hover:bg-surface-container-high font-bold",
      outline: "bg-transparent border-2 border-outline-variant/30 text-on-surface hover:border-primary/50 font-bold",
      danger: "bg-error text-on-error shadow-lg shadow-error/20 hover:scale-[0.98] italic",
      surface: "bg-surface-container-lowest text-on-surface border border-outline-variant/10 shadow-sm hover:shadow-md font-bold"
    };

    const sizes = {
      sm: "px-3 py-1.5 text-[10px] tracking-widest",
      md: "px-6 py-3 text-xs tracking-widest",
      lg: "px-8 py-4 text-sm tracking-widest",
      xl: "px-10 py-5 text-base tracking-widest"
    };

    return (
      <button
        ref={ref}
        className={cn(
          "inline-flex items-center justify-center rounded-xl font-headline uppercase transition-all duration-200 disabled:opacity-50 disabled:pointer-events-none",
          variants[variant],
          sizes[size],
          className
        )}
        {...props}
      >
        {isLoading ? (
          <span className="mr-2 h-4 w-4 animate-spin material-symbols-outlined !text-xs">autorenew</span>
        ) : null}
        {children}
      </button>
    );
  }
);

Button.displayName = "Button";
