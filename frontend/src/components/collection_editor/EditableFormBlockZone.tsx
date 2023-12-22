import { useDroppable } from "@dnd-kit/core";
import { cn } from "../../utils";
import FormBlockRenderer from "./FormBlockRenderer";
import { FormContext, copyContextWith, useFormContext } from "./FormContext";
import { FormBlock } from "./types";
import { useMemo } from "react";

export default function EditableFormBlockZone({ zoneKey, max, editable, uniqueKey, children = [] }: {
  zoneKey: string,
  max?: number,
  editable?: boolean,
  uniqueKey: string,
  children?: FormBlock[]
}) {
  const context = useFormContext();
  const disabled = useMemo(() => (max && children.length >= max) || !(editable ?? context.isEditable), [max, children, editable, context.isEditable]);
  const {isOver, setNodeRef} = useDroppable({
    id: zoneKey,
    disabled: disabled
  });

  return (
    <FormContext.Provider value={copyContextWith(context, {
      parentKey: zoneKey,
      isEditable: editable ?? context.isEditable
    })}>
      <div ref={setNodeRef} className={cn('border rounded flex flex-col min-h-[24rem] h-full w-full', {
        'border-violet-500 rounded': isOver,
        'border-opacity-5': children.length !== 0,
        'border-none': disabled
      })}>
        {children.map((block) => (
          <FormBlockRenderer
            key={`editable_block_${block.type}_${context.parentKey}.${zoneKey}.${block.key}_${uniqueKey}`}
            block={block} />
        ))}

        {children.length === 0 && (
          <div className="rounded py-8 flex flex-col items-center justify-center h-full w-full flex-1">
            <div className="text-sm text-gray-600">Drop a block here</div>
          </div>
        )}
      </div>
    </FormContext.Provider>
  );
}
