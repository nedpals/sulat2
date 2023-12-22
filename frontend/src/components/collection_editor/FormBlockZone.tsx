import FormBlockRenderer from "./FormBlockRenderer";
import { FormBlock } from "./types";
import EditableFormBlockZone from "./EditableFormBlockZone";
import { useMemo } from "react";
import { nanoid } from "nanoid";
import { cn } from "../../utils";

export default function FormBlockZone({ zoneKey, max, editable, children = [], containerClassName }: {
  zoneKey: string,
  max?: number,
  editable?: boolean,
  children?: FormBlock[],
  containerClassName?: string
}) {
  const uniqueKey = useMemo(() => nanoid(), []);

  if (editable) {
    return <EditableFormBlockZone
              max={max}
              zoneKey={zoneKey}
              children={children}
              editable={editable}
              uniqueKey={uniqueKey} />;
  }

  return (
    <div className={cn("flex flex-col", containerClassName)}>
      {children.map((block) => (
        <FormBlockRenderer
          key={`editable_block_${block.type}_${zoneKey}.${block.key}_${uniqueKey}`}
          block={block} />
      ))}
    </div>
  );
}
