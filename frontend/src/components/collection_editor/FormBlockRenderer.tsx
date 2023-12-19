import { useState } from "react"
import { Button, Dialog, DialogTrigger, OverlayArrow, Popover } from "react-aria-components"
import { cn } from "../../utils"
import { FormContext, defaultFormContext, useFormContext } from "./FormContext"

import ButtonBlock, { buttonBlockInfo } from "./blocks/ButtonBlock"
import StackBlock, { stackBlockInfo } from "./blocks/StackBlock"
import TextBlock, { textBlockInfo } from "./blocks/TextBlock"
import SelectBlock, { selectBlockInfo } from "./blocks/SelectBlock"
import TextareaBlock, { textareaBlockInfo } from "./blocks/TextareaBlock"
import FormBlockZone from "./FormBlockZone"
import { FormBlock } from "./types"

export interface FormBlockRendererProps<T = Record<string, any>> { 
  block: FormBlock<T>
  className?: string 
}

function FallbackBlockRenderer({ block, className }: FormBlockRendererProps) {
  const { parentKey } = useFormContext();

  return (
    <div className={cn(
      className, 
      "h-full w-full py-8 border bg-white text-center flex items-center justify-center"
    )}>
      {parentKey}.{block.key}
    </div>
  )
}

const blocks = {
  [stackBlockInfo.id]: StackBlock,
  [buttonBlockInfo.id]: ButtonBlock,
  [textBlockInfo.id]: TextBlock,
  [textareaBlockInfo.id]: TextareaBlock,
  [selectBlockInfo.id]: SelectBlock,
};

const propTypesToExclude = ['blocks'];

function propTypeToBlockType(schema: Record<string, any>) {
  switch (schema.type) {
    case 'string':
      if (schema.enum) {
        return 'select';
      }
      return 'text';
    default:
      return schema.type;
  }
}

function propEntryToBlockPropValue(schema: Record<string, any>): Record<string, any> {
  const blockType = propTypeToBlockType(schema);

  switch (blockType) {
    case 'select':
      return { options: schema.enum, default_value: schema.default };
    case 'text':
    case 'textarea':
      return { default_value: schema.default };
    default:
      return {};
  }
}

export default function FormBlockRenderer(props: FormBlockRendererProps) {
  const { block } = props;
  const { isEditable, removeField, getPropertiesSchema } = useFormContext();
  const [isOptionsOpen, setIsOptionsOpen] = useState(false);
  const schema = getPropertiesSchema(block.type);

  let BlockComponent: React.FC<FormBlockRendererProps<any>> = FallbackBlockRenderer;
  if (block.type in blocks) {
    BlockComponent = blocks[block.type];
  }

  if (isEditable) {
    return (
      <div className="sulat-editable-block">
        <div className={cn('sulat-editable-block-options', { '!block': isOptionsOpen })}>
          <span className="bg-violet-400 text-white px-4 py-1 rounded font-medium text-sm inline-block">
            {block.key}
          </span>
          {schema && <DialogTrigger onOpenChange={setIsOptionsOpen}>
            <Button className="sulat-btn sulat-block-edit-btn is-small">Edit</Button>
            <Popover className="sulat-edit-block-dialog">
              <OverlayArrow>
                <svg width={12} height={12} className="stroke-[1px] fill-white block" viewBox="0 0 12 12">
                  <path d="M0 0 L6 6 L12 0" />
                </svg>
              </OverlayArrow>

              <Dialog className="outline-none">
                <div className="flex-col">
                  <p>Block settings</p>
                  
                  <FormContext.Provider value={defaultFormContext.copyWith({
                    isEditable: false,
                    parentKey: block.key,
                    values: block.properties,
                  })}>
                    <FormBlockZone 
                      max={1}
                      zoneKey={block.key + "_options"} 
                      children={Object.entries(schema)
                        .filter(([_, s]) => !propTypesToExclude.includes(s.type))
                        .map(([property, s]) => ({
                          key: property,
                          label: property,
                          type: propTypeToBlockType(s),
                          properties: propEntryToBlockPropValue(s),
                        }))} />
                  </FormContext.Provider>

                  <button 
                    onClick={() => removeField(block.key)} 
                    className="sulat-btn is-danger is-small mt-4">
                    Remove
                  </button>
                </div>
              </Dialog>
            </Popover>
          </DialogTrigger>}
        </div>

        <BlockComponent {...props} className={cn('pointer-events-none', props.className)} />
      </div>
    );
  }

  return <BlockComponent {...props} />
}