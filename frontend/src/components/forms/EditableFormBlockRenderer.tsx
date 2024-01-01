import { useEffect, useState } from "react"
import { Button, Dialog, DialogTrigger, OverlayArrow, Popover } from "react-aria-components"
import { cn } from "../../utils"

import BlockRender from "./blocks/BlockRender";
import { FormBlockInstance, getBlockByType } from "./blocks"
import { FormBlockRenderProps } from "./render"
import { useFormEditorContext } from "./editor"
import { FormProvider } from "./context";
import FormZone from "./FormSection";
import { useDraggable } from "@dnd-kit/core";
import { setProperty } from "dot-prop";
import { useImmer } from "use-immer";

function propTypeToBlockType(schema: Record<string, any>) {
  switch (schema.type) {
    case 'string':
    case 'number':
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
      return { options: schema.enum, defaultValue: schema.default };
    case 'text':
      return {
        type: schema.type !== 'number' ? 'text' : 'number',
        defaultValue: schema.default
      };
    case 'textarea':
      return { defaultValue: schema.default };
    default:
      return {};
  }
}

function propertiesToFormBlockInstance(
  schema: Record<string, { type: string, [k: string]: any }>,
  fieldKey?: boolean,
) {
  const blocks: FormBlockInstance[] = [];
  const excludedFormTypes = ['blocks', 'array'];

  for (const propKey in schema) {
    const propSchema = schema[propKey];
    if (excludedFormTypes.includes(propSchema.type)) continue;

    const blockType = propTypeToBlockType(propSchema);
    blocks.push({
      key: propKey,
      type: blockType,
      fieldKey: 'properties.' + propKey,
      properties: {
        label: propKey,
        ...propEntryToBlockPropValue(propSchema),
      },
    });
  }

  return [
    {
      key: 'properties',
      label: '',
      type: 'stack',
      fieldKey: '',
      properties: {
        children: [
          ...(typeof fieldKey === 'undefined' || fieldKey ? [
            {
              key: 'fieldKey',
              type: 'text',
              fieldKey: 'fieldKey',
              properties: {
                label: 'Field key',
                defaultValue: '',
              }
            }
          ] : []),
          ...blocks
        ],
      }
    }
  ];
}

function filterEditableProperties(
  schema: Record<string, { type: string, [k: string]: any }>,
  properties: Record<string, any>
) {
  const excludedFormTypes = ['blocks', 'array'];
  const filteredProperties: Record<string, any> = {};

  for (const propKey in schema) {
    if (excludedFormTypes.includes(schema[propKey].type)) continue;
    filteredProperties[propKey] = properties[propKey];
  }

  return filteredProperties;
}

export default function EditableFormBlockRenderer(props: FormBlockRenderProps) {
  const { block } = props;
  const blockProperties = getBlockByType(block.type).properties;
  const schema = blockProperties.propertiesSchema;
  const { attributes, listeners, setNodeRef, transform } = useDraggable({ id: props.block.key });

  const { removeBlock, updateBlock } = useFormEditorContext();
  const [isOptionsOpen, setIsOptionsOpen] = useState(false);
  const [properties, setProperties] = useImmer({
    fieldKey: block.fieldKey,
    properties: filterEditableProperties(schema, block.properties),
  });

  // useEffect(() => {
  //   console.log('INITIALIZING', block.key, properties);
  // }, []);

  useEffect(() => {
    // console.log('UPDATING', block.key);
    updateBlock(null, block.key, properties);
  }, [properties]);

  // useEffect(() => {
  //   console.log(block);
  // }, [block]);

  return (
    <div
      style={transform ? {
        transform: `translate3d(${transform.x}px, ${transform.y}px, 0)`,
        backgroundColor: 'white',
      } : undefined}
      className="sulat-editable-block">
      <div className={cn('sulat-editable-block-options', { '!block': isOptionsOpen })}>
        <button
          ref={setNodeRef}
          {...listeners} {...attributes}
          className="bg-violet-400 text-white cursor-move px-4 py-1 rounded font-medium text-sm inline-block">
          {block.key}
        </button>

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
                <div className="mb-2 flex items-center justify-between">
                  <p className="text-sm font-semibold">Block settings</p>

                  <button
                    onClick={() => removeBlock(null, block.key)}
                    className="sulat-btn is-danger is-small">
                    Remove
                  </button>
                </div>

                <FormProvider state={properties ?? {}} onChange={(p) => {
                  setProperties(draft => {
                    for (const propKey in p) {
                      setProperty(draft, propKey, p[propKey]);
                    }
                  });
                }}>
                  <FormZone
                    name="properties"
                    editable={false}
                    children={propertiesToFormBlockInstance(
                      schema,
                      //@ts-expect-error
                      blockProperties.fieldKey,
                    )} />
                </FormProvider>
              </div>
            </Dialog>
          </Popover>
        </DialogTrigger>}
      </div>

      <BlockRender {...props} className={cn('pointer-events-none', props.className)} />
    </div>
  );
}
