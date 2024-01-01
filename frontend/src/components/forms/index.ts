import { createContext, useContext } from "react";
import { FormBlock, FormBlockInstance } from "./blocks";

export type FormSchema = Record<string, FormBlockInstance[]>;

export interface FormSectionContextProps {
  name: string
  isEditable: boolean
  parentKey: string
  children: FormBlockInstance[]
  onChange?: (blocks: FormBlockInstance[]) => void
}

export const FormSectionContext = createContext<FormSectionContextProps>({ isEditable: false, name: '', children: [], parentKey: '' });
FormSectionContext.displayName = 'FormSectionContext';

export const FormSectionProvider = FormSectionContext.Provider;
export const useFormSectionContext = () => useContext(FormSectionContext);

export function makeKey(parentKey: string, name: string) {
  return [parentKey, name].filter(Boolean).join('.');
}

export function createBlockInstance(block: FormBlock): FormBlockInstance {
  const properties: Record<string, any> = {};

  for (const [key, value] of Object.entries(block.propertiesSchema)) {
    if (key === 'label') {
      properties[key] = value.default ? value.default : `New ${block.name}`;
    } else {
      properties[key] = value.default ?? null;
    }
  }

  return {
    fieldKey: '',
    key: block.id,
    type: block.id,
    properties: properties,
  }
}
