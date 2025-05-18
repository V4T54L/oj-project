import { diffWords } from "diff";

export function getDiffText(expected: string, actual: string) {
    const diff = diffWords(expected, actual);
    return diff.map((part, i) => {
        const className = part.added
            ? "bg-green-800 text-green-300"
            : part.removed
                ? "bg-red-800 text-red-300 line-through"
                : "";

        return (
            <span key={i} className={`${className} px-1`
            }>
                {part.value}
            </span>
        );
    });
}
