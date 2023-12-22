import { cn } from "../../../utils";
import { FormBlockRendererProps } from "../FormBlockRenderer";
import FormBlockZone from "../FormBlockZone";
import { useFormContext } from "../FormContext";
import { FormBlock } from "../types";

export interface StackBlockProps {
  direction: 'horizontal' | 'vertical'
  children?: FormBlock[]
}

export default function StackBlock({ block }: FormBlockRendererProps<StackBlockProps>) {
  const context = useFormContext();

  return (
    <FormBlockZone
      containerClassName={cn('flex-1 min-h-8 h-full w-full flex', {
        'flex-row': block.properties.direction === 'horizontal',
        'flex-col space-y-4': block.properties.direction === 'vertical',
      })}
      editable={context.isEditable}
      zoneKey={[context.parentKey, block.key].filter(Boolean).join('.')}
      children={block.properties.children?.filter(b => b.type !== 'stack')} />
  );
}

StackBlock.properties = {
  id: 'stack',
  name: 'Stack',
  description: 'A stack of blocks. Can be horizontal or vertical.',
  propertiesSchema: {
    direction: { type: 'string', enum: ['primary', 'secondary'], default: 'vertical' },
    children: { type: 'blocks', default: [] }
  }
}
