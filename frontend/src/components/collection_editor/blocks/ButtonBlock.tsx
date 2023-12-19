import { emitter } from "../../../eventbus";
import { cn } from "../../../utils";
import { FormBlockRendererProps } from "../FormBlockRenderer";

// NOTE to self: when building extension feature, create "actions" feature
// which can be registered by attaching a .on in the eventbus

interface ButtonBlockProps {
  text: string
  type: 'primary' | 'secondary'
  task: string
  arguments: any
}

export const buttonBlockInfo = {
  id: 'button',
  name: 'Button',
  description: 'A button',
  propertiesSchema: {
    text: { type: 'string', default: 'Button' },
    type: { type: 'string', enum: ['primary', 'secondary'], default: 'primary' },
    task: { type: 'string', default: 'log' },
    arguments: { type: 'object', default: {} },
  }
}

export default function ButtonBlock({ block, className }: FormBlockRendererProps<ButtonBlockProps>) {
  return (
    <button
      onClick={() => emitter.emit(block.properties.task, block.properties.arguments)} 
      className={cn('sulat-btn is-primary self-stretch w-full', className)}>
      {block.properties.text}
    </button>
  );
}