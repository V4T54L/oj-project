import React, { useEffect } from "react";
import { AnimatePresence, motion, type MotionProps } from "framer-motion";
import clsx from "clsx";

interface ModalProps {
    children: React.ReactNode;
    isOpen: boolean;
    onOpenChange: (isOpen: boolean) => void;
    size?: "xs" | "sm" | "md" | "lg" | "xl" | "2xl" | "3xl" | "4xl" | "5xl" | "full";
    radius?: "none" | "sm" | "md" | "lg";
    placement?:
    | "auto"
    | "top"
    | "bottom"
    | "center"
    | "top-center"
    | "bottom-center";
    backdrop?: "transparent" | "opaque" | "blur";
    scrollBehavior?: "normal" | "inside" | "outside";
    isDismissable?: boolean;
    isKeyboardDismissDisabled?: boolean;
    hideCloseButton?: boolean;
    motionProps?: MotionProps;
}

const sizeMap: Record<string, string> = {
    xs: "max-w-xs",
    sm: "max-w-sm",
    md: "max-w-md",
    lg: "max-w-lg",
    xl: "max-w-xl",
    "2xl": "max-w-2xl",
    "3xl": "max-w-3xl",
    "4xl": "max-w-4xl",
    "5xl": "max-w-5xl",
    full: "w-full h-full",
};

const radiusMap: Record<string, string> = {
    none: "rounded-none",
    sm: "rounded-sm",
    md: "rounded-md",
    lg: "rounded-lg",
};

const placementMap: Record<string, string> = {
    center: "items-center justify-center",
    top: "items-start justify-center",
    bottom: "items-end justify-center",
    "top-center": "items-start justify-center",
    "bottom-center": "items-end justify-center",
    auto: "items-center justify-center",
};

const backdropMap: Record<string, string> = {
    transparent: "bg-transparent",
    opaque: "bg-black bg-opacity-50",
    blur: "bg-black/40 backdrop-blur-md",
};

export const Modal: React.FC<ModalProps> = ({
    children,
    isOpen,
    onOpenChange,
    size = "md",
    radius = "md",
    placement = "center",
    backdrop = "opaque",
    scrollBehavior = "normal",
    isDismissable = true,
    isKeyboardDismissDisabled = false,
    hideCloseButton = false,
    motionProps = {},
}) => {
    useEffect(() => {
        const handleKeyDown = (e: KeyboardEvent) => {
            if (e.key === "Escape" && !isKeyboardDismissDisabled) {
                onOpenChange(false);
            }
        };
        if (isOpen) {
            document.addEventListener("keydown", handleKeyDown);
        }
        return () => document.removeEventListener("keydown", handleKeyDown);
    }, [isOpen, isKeyboardDismissDisabled, onOpenChange]);

    if (typeof window === "undefined") return null;

    return (
        <AnimatePresence>
            {isOpen && (
                <div
                    className={clsx(
                        "fixed inset-0 z-50 flex",
                        placementMap[placement],
                        backdropMap[backdrop]
                    )}
                    onClick={() => isDismissable && onOpenChange(false)}
                >
                    <motion.div
                        onClick={(e: any) => e.stopPropagation()}
                        initial={{ opacity: 0, scale: 0.95 }}
                        animate={{ opacity: 1, scale: 1 }}
                        exit={{ opacity: 0, scale: 0.95 }}
                        transition={{ duration: 0.2 }}
                        {...motionProps}
                        className={clsx(
                            "w-full relative",
                            "bg-gray-800 text-white",
                            sizeMap[size],
                            radiusMap[radius],
                            scrollBehavior === "inside" ? "overflow-y-auto max-h-[90vh]" : "",
                            "shadow-xl p-6"
                        )}
                    >
                        {!hideCloseButton && (
                            <button
                                onClick={() => onOpenChange(false)}
                                className="absolute top-4 right-4 text-gray-400 hover:text-gray-200 text-xl"
                                aria-label="Close"
                            >
                                ×
                            </button>
                        )}
                        {children}
                    </motion.div>
                </div>
            )}
        </AnimatePresence>
    );
};

export const ModalHeader: React.FC<{ children: React.ReactNode }> = ({ children }) => (
    <div className="text-xl font-semibold text-white mb-4 border-b border-gray-700 pb-2">
        {children}
    </div>
);

export const ModalBody: React.FC<{ children: React.ReactNode }> = ({ children }) => (
    <div className="text-gray-200 mb-4">{children}</div>
);

export const ModalFooter: React.FC<{ children: React.ReactNode }> = ({ children }) => (
    <div className="flex justify-end gap-2 border-t border-gray-700 pt-4">{children}</div>
);
