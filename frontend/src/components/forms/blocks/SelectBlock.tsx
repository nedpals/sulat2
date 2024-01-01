import { Button, Label, ListBox, ListBoxItem, Popover, Select, SelectValue } from "react-aria-components";
import { cn } from "../../../utils";
import { FormBlockRenderProps } from "../render";
import { useFormContext } from "../context";

interface SelectBlockProps {
  label: string
  size: 'small' | 'medium' | 'large'
  options: string[]
  default_value: string
  placeholder?: string
}

export default function SelectBlock({ block, className }: FormBlockRenderProps<SelectBlockProps>) {
  const { getField, updateField } = useFormContext();

  // return (
  //   <Select
  //     className={cn(className)}
  //     defaultSelectedKey={getField(block.fieldKey, block.properties.default_value)}>
  //     <Label>{block.properties.label}</Label>
  //     <Button>
  //       <SelectValue />
  //       <span aria-hidden="true">▼</span>
  //     </Button>
  //     <Popover>
  //       <ListBox>
  //         {block.properties.options.map((option) => (
  //           <ListBoxItem key={option} textValue={option}>
  //             {option}
  //           </ListBoxItem>
  //         ))}
  //       </ListBox>
  //     </Popover>
  //   </Select>
  // );

  return (
    <select
      className={cn('sulat-input w-full', className, {
        'is-small': block.properties.size === 'small',
        'is-large': block.properties.size === 'large',
      })}
      defaultValue={getField(block.fieldKey, block.properties.default_value)}
      onChange={(ev) => updateField(block.fieldKey, ev.target.value)}>
      {block.properties.placeholder && <option>{block.properties.placeholder}</option>}
      {block.properties.options.map((option) => (
        <option key={option} value={option}>
          {option}
        </option>
      ))}
    </select>
  )
}

SelectBlock.properties = {
  id: 'select',
  name: 'Select',
  description: 'A select input',
  propertiesSchema: {
    label: { type: 'string', default: '' },
    size: { type: 'string', enum: ['small', 'medium', 'large'], default: 'medium' },
    options: { type: 'array', default: [] },
    default_value: { type: 'string', default: '' },
    placeholder: { type: 'string', default: 'Select' },
  }
}
