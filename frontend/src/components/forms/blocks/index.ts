import ButtonBlock from "./ButtonBlock"
import FallbackBlock from "./FallbackBlock"
import SelectBlock from "./SelectBlock"
import StackBlock from "./StackBlock"
import TextBlock from "./TextBlock"
import TextareaBlock from "./TextareaBlock"

export interface FormBlockInstance<T = Record<string, any>> {
  key: string
  type: string
  fieldKey: string
  properties: T
}

export interface FormBlock<T = Record<string, any>> {
  id: string
  name: string
  description: string
  propertiesSchema: T
}

export interface BlockGroup {
  id: string
  title: string
  blocks: FormBlock[]
}

export const blockList: BlockGroup[] = [
  {
    id: 'layout',
    title: 'Layout',
    blocks: [
      StackBlock.properties
    ]
  },
  {
    id: 'buttons',
    title: 'Buttons',
    blocks: [
      ButtonBlock.properties,
    ]
  },
  {
    id: 'input',
    title: 'Input',
    blocks: [
      TextBlock.properties,
      TextareaBlock.properties,
      SelectBlock.properties,
    ]
  }
];

export const blockTypeToBlocks = {
  [StackBlock.properties.id]: StackBlock,
  [ButtonBlock.properties.id]: ButtonBlock,
  [TextBlock.properties.id]: TextBlock,
  [TextareaBlock.properties.id]: TextareaBlock,
  [SelectBlock.properties.id]: SelectBlock,
};

export function getBlockByType(type: string) {
  if (!(type in blockTypeToBlocks)) {
    return FallbackBlock;
  }
  return blockTypeToBlocks[type];
}

export {
  FallbackBlock,
  ButtonBlock,
  SelectBlock,
  StackBlock,
  TextBlock,
  TextareaBlock,
}
