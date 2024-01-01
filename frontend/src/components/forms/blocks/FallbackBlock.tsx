import { useFormSectionContext } from "..";
import { cn } from "../../../utils";
import { FormBlockRenderProps } from "../render";

export default function FallbackBlock({ block, className }: FormBlockRenderProps) {
  const { parentKey } = useFormSectionContext();

  return (
    <div className={cn(
      className,
      "h-full w-full py-2 border bg-white text-center flex items-center justify-center"
    )}>
      {parentKey}.{block.key}
    </div>
  )
}

FallbackBlock.properties = {
  id: 'fallback',
  name: 'Fallback',
  description: 'A fallback block',
  propertiesSchema: {}
}
