import * as React from "react";
import { cva, type VariantProps } from "class-variance-authority";
import { cn } from "../../lib/utils";

const buttonVariants = cva(
  "inline-flex items-center justify-center gap-1.5 whitespace-nowrap rounded-md text-[13.5px] font-medium transition-colors disabled:pointer-events-none disabled:opacity-50 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary-500 focus-visible:ring-offset-1",
  {
    variants: {
      variant: {
        default:
          "bg-surface border border-border-strong text-ink-primary hover:bg-surface-secondary",
        primary:
          "bg-primary-600 border border-primary-600 text-white hover:bg-primary-700 hover:border-primary-700",
        danger:
          "bg-danger-600 border border-danger-600 text-white hover:bg-danger-700 hover:border-danger-700",
        ghost:
          "border border-transparent text-ink-secondary hover:bg-surface-secondary hover:text-ink-primary",
        outline:
          "bg-transparent border border-border-strong text-ink-primary hover:bg-surface-secondary",
      },
      size: {
        default: "h-9 px-3.5 py-2",
        sm: "h-7 px-2.5 text-[12.5px]",
        lg: "h-10 px-5",
        icon: "h-8 w-8 p-0",
      },
    },
    defaultVariants: {
      variant: "default",
      size: "default",
    },
  }
);

export interface ButtonProps
  extends React.ButtonHTMLAttributes<HTMLButtonElement>,
    VariantProps<typeof buttonVariants> {}

const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant, size, ...props }, ref) => {
    return (
      <button className={cn(buttonVariants({ variant, size, className }))} ref={ref} {...props} />
    );
  }
);
Button.displayName = "Button";

export { Button, buttonVariants };