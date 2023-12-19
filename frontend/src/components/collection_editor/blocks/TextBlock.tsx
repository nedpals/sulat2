import { FormBlockRendererProps } from "../FormBlockRenderer";
import { useFormContext } from "../FormContext";

export interface TextBlockProps {
  label: string
  default_value?: string
  placeholder: string
}

export const textBlockInfo = {
  id: 'text',
  name: 'Text',
  description: 'A text input',
  propertiesSchema: {
    label: { type: 'string', default: 'Text' },
    default_value: { type: 'string', default: '' },
    placeholder: { type: 'string', default: 'Text' },
  }
}

export default function TextBlock({ block }: FormBlockRendererProps<TextBlockProps>) {
  const { getFieldValue, setFieldValue } = useFormContext();

  return (
    <div>
      <input 
        type="text" 
        defaultValue={getFieldValue(block.key, block.properties.default_value)}
        onChange={() => setFieldValue(block.key, block.properties.default_value)}
        className="sulat-input text-xl w-full" 
        placeholder={block.properties.placeholder ?? 'Title'} />
    </div>
  );
}