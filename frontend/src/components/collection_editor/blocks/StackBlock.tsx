import { cn } from "../../../utils";
import { FormBlock, FormBlockRendererProps } from "../FormBlockRenderer";
import FormBlockZone from "../FormBlockZone";
import { useFormContext } from "../FormContext";

export interface StackBlockProps {
  direction: 'horizontal' | 'vertical'
  children?: FormBlock[]
}

export const stackBlockInfo = {
  id: 'stack',
  name: 'Stack',
  description: 'A stack of blocks. Can be horizontal or vertical.',
  propertiesSchema: {
    direction: { type: 'string', enum: ['primary', 'secondary'], default: 'vertical' },
    children: { type: 'blocks', default: [] }
  }
}

export default function StackBlock({ block }: FormBlockRendererProps<StackBlockProps>) {
  const context = useFormContext();
  const parentKey = context.parentKey;
  const zoneKey = [parentKey, block.key].filter(Boolean).join('.');

  return (
    <div className={cn('flex-1 min-h-8 h-full w-full flex', {
      'flex-row': block.properties.direction === 'horizontal',
      'flex-col space-y-4': block.properties.direction === 'vertical',
    })}>
      <FormBlockZone 
        zoneKey={zoneKey}
        children={block.properties.children} />
    </div>
  );
}