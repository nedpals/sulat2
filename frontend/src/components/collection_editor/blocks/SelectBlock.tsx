import { Button, Label, ListBox, ListBoxItem, Popover, Select, SelectValue } from "react-aria-components";
import { FormBlockRendererProps } from "../FormBlockRenderer";
import { cn } from "../../../utils";
import { useFormContext } from "../FormContext";

interface SelectBlockProps {
  label: string
  options: string[]
  default_value: string
  placeholder?: string
}

export const selectBlockInfo = {
  id: 'select',
  name: 'Select',
  description: 'A select input',
  propertiesSchema: {
    label: { type: 'string', default: 'Select' },
    options: { type: 'array', default: [] },
    default_value: { type: 'string', default: '' },
    placeholder: { type: 'string', default: 'Select' },
  }
}

export default function SelectBlock({ block, className }: FormBlockRendererProps<SelectBlockProps>) {
  const { getFieldValue } = useFormContext();
  
  return (
    <Select 
      className={cn(className)}
      defaultSelectedKey={getFieldValue(block.key, block.properties.default_value)}>
      <Label>{block.properties.label}</Label>
      <Button>
        <SelectValue />
        <span aria-hidden="true">▼</span>
      </Button>
      <Popover>
        <ListBox>
          {block.properties.options.map((option) => (
            <ListBoxItem key={option} textValue={option}>
              {option}
            </ListBoxItem>
          ))}
        </ListBox>
      </Popover>
    </Select>
  );
}