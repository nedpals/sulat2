import { createContext, useContext } from "react"
import { FormBlockInstance } from "./blocks";

export interface FormEditorContextData {
  // Editing blocks
  updateBlock: (parentKey: string | null, key: string, properties: Record<string, any>) => void
  removeBlock: (parentKey: string | null, key: string) => void
}

export const FormEditorContext = createContext<FormEditorContextData>({
  updateBlock: () => {},
  removeBlock: () => {},
});

FormEditorContext.displayName = 'FormEditorContext';

export const FormEditorProvider = FormEditorContext.Provider;
export const useFormEditorContext = () => useContext(FormEditorContext);

export function findAndInsertFormBlock(blocks: FormBlockInstance[], location: string, block: FormBlockInstance) {
  if (location.length === 0) {
    const existingBlock = blocks.find(b => b.key === block.key);
    if (existingBlock) {
      for (let i = 1;; i++) {
        const newBlockKey = `${block.key}_${i}`;
        const existingBlock = blocks.find(b => b.key === newBlockKey);
        if (existingBlock) {
          continue;
        }

        block.key = newBlockKey;
        break;
      }
    }

    blocks.push(block);
    return;
  }

  const locationArr = location.split('.').filter(Boolean);
  const loc = locationArr.shift();
  if (typeof loc === 'undefined') {
    return;
  }

  for (let i = 0; i < blocks.length; i++) {
    const _block = blocks[i];
    if (_block.key !== loc) {
      continue;
    }

    // If there are still keys left, we need to go deeper
    if ('children' in _block.properties) {
      const otherKeys = locationArr.join('.');
      findAndInsertFormBlock(blocks[i].properties.children, otherKeys, block);
    }

    break;
  }
}
