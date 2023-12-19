import { FormBlockRendererProps } from "../FormBlockRenderer";
import { useFormContext } from "../FormContext";

export interface TextareaBlockProps {
  label: string
  default_value?: string
  placeholder: string
}

export const textareaBlockInfo = {
  id: 'textarea',
  name: 'Textarea',
  description: 'A textarea input',
  propertiesSchema: {
    label: { type: 'string', default: 'Textarea' },
    default_value: { type: 'string', default: '' },
    placeholder: { type: 'string', default: 'Textarea' },
  }
}

export default function TextareaBlock({ block }: FormBlockRendererProps<TextareaBlockProps>) {
  const { getFieldValue, setFieldValue } = useFormContext();

  return (
    <div>
      <textarea
        defaultValue={getFieldValue(block.key, block.properties.default_value)}
        onChange={() => setFieldValue(block.key, block.properties.default_value)}
        className="sulat-input text-xl w-full"
        placeholder={block.properties.placeholder ?? 'Title'} />
    </div>
  );
}