import { FormBlockRenderProps } from "../render";
import { emitter } from "../../../eventbus";
import { cn } from "../../../utils";

// NOTE to self: when building extension feature, create "actions" feature
// which can be registered by attaching a .on in the eventbus

interface ButtonBlockProps {
  text: string
  type: 'primary' | 'secondary'
  action: string
  arguments: any
}

export default function ButtonBlock({ block, className }: FormBlockRenderProps<ButtonBlockProps>) {
  return (
    <button
      onClick={() => emitter.emit('triggerAction', {
        action: block.properties.action,
        arguments: block.properties.arguments
      })}
      className={cn('sulat-btn self-stretch w-full', className, {
        'is-primary': block.properties.type === 'primary',
        'is-secondary': block.properties.type === 'secondary',
      })}>
      {block.properties.text}
    </button>
  );
}

ButtonBlock.properties = {
  id: 'button',
  name: 'Button',
  description: 'A button',
  propertiesSchema: {
    text: { type: 'string', default: 'Button' },
    type: { type: 'string', enum: ['primary', 'secondary'], default: 'primary' },
    action: { type: 'string', default: 'log' },
    arguments: { type: 'object', default: {} },
  }
}
