import { useFormSectionContext } from "..";
import { cn } from "../../../utils";
import { useFormContext } from "../context";
import { FormBlockRenderProps } from "../render";

export interface TextBlockProps {
  label: string
  size: 'small' | 'medium' | 'large'
  type: 'text' | 'password' | 'email' | 'number'
  default_value?: string
  placeholder: string
}

export default function TextBlock({ block }: FormBlockRenderProps<TextBlockProps>) {
  const { getField, updateField } = useFormContext();
  const { isEditable } = useFormSectionContext();

  return (
    <div>
      <input
        type={block.properties.type}
        disabled={isEditable}
        defaultValue={getField(block.fieldKey, block.properties.default_value)}
        onChange={(evt) => updateField(block.fieldKey, evt.currentTarget.value)}
        className={cn('sulat-input w-full', {
          'is-small': block.properties.size === 'small',
          'is-large': block.properties.size === 'large',
        })}
        placeholder={block.properties.placeholder} />
    </div>
  );
}

TextBlock.properties = {
  id: 'text',
  name: 'Text',
  description: 'A text input',
  propertiesSchema: {
    label: { type: 'string', default: '' },
    type: { type: 'string', enum: ['text', 'password', 'email', 'number'], default: 'text' },
    size: { type: 'string', enum: ['small', 'medium', 'large'], default: 'medium' },
    default_value: { type: 'string', default: '' },
    placeholder: { type: 'string', default: 'Text' },
  }
}
