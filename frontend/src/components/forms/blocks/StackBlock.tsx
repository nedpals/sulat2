import { FormBlockInstance } from ".";
import { cn } from "../../../utils";
import FormZone from "../FormSection";
import { FormBlockRenderProps } from "../render";

export interface StackBlockProps {
  direction: 'horizontal' | 'vertical'
  children: FormBlockInstance[]
}

export default function StackBlock({ block }: FormBlockRenderProps<StackBlockProps>) {
  return (
    <FormZone
      containerClassName={cn('flex-1 h-full w-full flex', {
        'flex-row': block.properties.direction === 'horizontal',
        'flex-col': block.properties.direction === 'vertical',
      })}
      childWrapper={({ children }) => <div className="flex-shrink-1">{children}</div>}
      name={block.key}
      parent={block.key}
      children={block.properties.children} />
  );
}

StackBlock.properties = {
  id: 'stack',
  name: 'Stack',
  description: 'A stack of blocks. Can be horizontal or vertical.',
  fieldKey: false,
  propertiesSchema: {
    direction: { type: 'string', enum: ['vertical', 'horizontal'], default: 'vertical' },
    children: { type: 'blocks', default: [] }
  }
}
