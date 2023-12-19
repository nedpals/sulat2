import { useDraggable } from "@dnd-kit/core";
import { cn } from "../../utils";

export default function FormBlockButton({ id, title, description, className }: {
  id: string
  title: string
  description: string
  className?: string
}) {
  const { attributes, listeners, setNodeRef, transform } = useDraggable({ id });

  return (
    <button
      ref={setNodeRef} {...listeners} {...attributes}
      style={transform ? {
        transform: `translate3d(${transform.x}px, ${transform.y}px, 0)`
      } : undefined}
      className={cn(
        className,
        "flex space-x-2 bg-white border hover:border-violet-500 transition-colors shadow rounded px-3 py-2 cursor-grab",
        { 'cursor-grabbing': transform },
      )}>
      <div className="h-6 w-6 rounded bg-slate-400"></div>
      <div className="flex flex-col items-start text-left">
        <span className="text-sm">{title}</span>
        <span className="text-xs text-ellipsis w-full">{description}</span>
      </div>
    </button>
  )
}