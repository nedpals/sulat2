import { useFormSectionContext } from "..";
import { useFormContext } from "../context";
import { FormBlockRenderProps } from "../render";

export interface TextareaBlockProps {
  label: string
  default_value?: string
  placeholder: string
  columns?: number
  rows?: number
}

export default function TextareaBlock({ block }: FormBlockRenderProps<TextareaBlockProps>) {
  const { getField, updateField } = useFormContext();
  const { isEditable } = useFormSectionContext();

  return (
    <div>
      <textarea
        cols={block.properties.columns}
        rows={block.properties.rows}
        defaultValue={getField(block.fieldKey, block.properties.default_value)}
        disabled={isEditable}
        onChange={() => updateField(block.fieldKey, block.properties.default_value)}
        className="sulat-input text-xl w-full"
        placeholder={block.properties.placeholder ?? 'Title'} />
    </div>
  );
}

TextareaBlock.properties = {
  id: 'textarea',
  name: 'Textarea',
  description: 'A textarea input',
  propertiesSchema: {
    label: { type: 'string', default: '' },
    default_value: { type: 'string', default: '' },
    placeholder: { type: 'string', default: 'Textarea' },
    columns: { type: 'number', default: 10 },
    rows: { type: 'number', default: 10 },
  }
}
