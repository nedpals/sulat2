import { useDroppable } from "@dnd-kit/core";
import { cn } from "../../utils";
import FormBlockRenderer from "./FormBlockRenderer";
import { FormContext, useFormContext } from "./FormContext";
import { FormBlock } from "./types";

export default function FormBlockZone({ zoneKey, max, editable, children = [] }: { 
  zoneKey: string, 
  max?: number,
  editable?: boolean,
  children?: FormBlock[] 
}) {
  const context = useFormContext();
  const {isOver, setNodeRef} = useDroppable({ 
    id: zoneKey, 
    disabled: (max && children.length >= max) || !(editable ?? context.isEditable)
  });

  return (
    <FormContext.Provider value={context.copyWith({ 
      parentKey: zoneKey, 
      isEditable: editable ?? context.isEditable 
    })}>
      <div ref={setNodeRef} className={cn('flex flex-col min-h-[24rem] h-full w-full', {
        'border border-violet-500 rounded': isOver
      })}>
        {children.map((block) => (
          <FormBlockRenderer 
            key={`editable_block_${block.type}_${zoneKey}.${block.key}`}
            block={block} />
        ))}

        {children.length === 0 && (
          <div className="border rounded py-8 flex flex-col items-center justify-center h-full w-full flex-1">
            <div className="text-sm text-gray-600">Drop a block here</div>
          </div>
        )}
      </div>
    </FormContext.Provider>
  );
}