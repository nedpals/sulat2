import { FormBlockRendererProps } from "../FormBlockRenderer";
import { useFormContext } from "../FormContext";

export interface TextBlockProps {
  label: string
  default_value?: string
  placeholder: string
}

export default function TextBlock({ block }: FormBlockRendererProps<TextBlockProps>) {
  const { getValue, setValue } = useFormContext();

  return (
    <div>
      <input
        type="text"
        defaultValue={getValue(block.key, block.properties.default_value)}
        onChange={() => setValue(block.key, block.properties.default_value)}
        className="sulat-input text-xl w-full"
        placeholder={block.properties.placeholder ?? 'Title'} />
    </div>
  );
}

TextBlock.properties = {
  id: 'text',
  name: 'Text',
  description: 'A text input',
  propertiesSchema: {
    label: { type: 'string', default: 'Text' },
    default_value: { type: 'string', default: '' },
    placeholder: { type: 'string', default: 'Text' },
  }
}
