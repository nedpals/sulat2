import { createContext, useContext } from "react"
import { FormBlock } from "./types"

interface FormContextDataProperties {
  level: number
  isEditable: boolean
  parentKey?: string
  blocks?: FormBlock[]
  values: Record<string, any>
}

export interface FormContextData extends FormContextDataProperties {
  getBlockSchema: (blockKey: string) => Record<string, any> | null
  getValue(key: string, defaultValue?: any): any
  setValue(key: string, value: any): void
  setField(key: string, property: string, value: any): void
  removeField(key: string): void
}

type ModifiableFormContextFields = Partial<FormContextData>;

export function copyContextWith(oldContext: FormContextData, context: ModifiableFormContextFields): FormContextData {
  console.log(oldContext);

  return {
    ...oldContext,
    ...context,
    level: oldContext.level + 1,
  }
}

export function setFormField(blocks: FormBlock[], key: string, property: string, value: any) {
  // Find block first
  const keyArr = key.split('.');
  const blockKey = keyArr.shift();
  if (!blockKey) {
    return;
  }

  const blockIdx = blocks.findIndex((block) => block.key === blockKey);
  if (blockIdx === -1) {
    return;
  }

  const block = blocks[blockIdx];
  if (keyArr.length === 0) {
    block.properties[property] = value;
    return;
  }
}

export function getFieldValue(values: Record<string, any>, key: string, defaultValue: any = null): any {
  const keyArr = key.split('.');
  const fieldKey = keyArr.shift();
  if (!fieldKey || !values[fieldKey]) {
    return defaultValue;
  }

  if (keyArr.length === 0) {
    return values[fieldKey];
  }

  if (Array.isArray(values[fieldKey])) {
    const nextKey = keyArr.shift();
    if (!nextKey) {
      return defaultValue;
    }

    const index = parseInt(nextKey);
    if (isNaN(index)) {
      return defaultValue;
    }

    return getFieldValue(values[fieldKey][index], keyArr.join('.'), defaultValue);
  }

  return getFieldValue(values[fieldKey], keyArr.join('.'), defaultValue);
}

export function setFormFieldValue(values: Record<string, any>, key: string, value: any) {
  const keyArr = key.split('.');
  let firstKey = keyArr.shift();
  if (!firstKey) {
    return;
  }

  const isArrayPush = firstKey.match(/\[\]$/);
  if (isArrayPush) {
    firstKey = firstKey.replace(/\[\]$/, '');
  }

  if (!values[firstKey] && keyArr.length !== 0) {
    const nextKey = keyArr.shift();
    if (!nextKey) {
      return;
    }

    if (!isNaN(parseInt(nextKey))) {
      values[firstKey] = [];
    } else {
      values[firstKey] = {};
    }
  }

  if (keyArr.length !== 0) {
    return setFormFieldValue(values[firstKey], keyArr.join('.'), value);
  } else if (isArrayPush && Array.isArray(values[firstKey])) {
    values[firstKey].push(value);
    return;
  }

  values[firstKey] = value;
}

export function removeField(blocks: FormBlock[], key: string): FormBlock[] {
  console.log('REMOVE FIELD', key);

  const keyArr = key.split('.');
  const blockKey = keyArr.shift();

  if (!blockKey) {
    return blocks;
  }

  for (let i = 0; i < blocks.length; i++) {
    const block = blocks[i];
    if (block.key !== blockKey) {
      continue;
    }

    if (keyArr.length === 0) {
      let newBlocks = [...blocks];
      newBlocks.splice(i, 1);
      return newBlocks;
    }

    if ('children' in block.properties && Array.isArray(block.properties.children)) {
      blocks[i].properties.children = removeField(block.properties.children, keyArr.join('.'));
      return blocks;
    }
  }

  return blocks;
}

export function makeContextValue(data: Partial<FormContextDataProperties>): FormContextData {
  return {
    ...defaultFormContext,
    ...data,
  }
}

export const defaultFormContext: FormContextData = {
  level: 0,
  isEditable: false,
  values: {},
  getBlockSchema() {
    return null;
  },
  getValue: function(key, defaultValue = null) {
    console.log('GET FIELD VALUE', key, defaultValue, this);
    if (!this.values) {
      return defaultValue;
    }
    return getFieldValue(this.values, key, defaultValue);
  },
  setValue: function(key: string, val: any) {
    setFormFieldValue(this.values, key, val);
  },
  setField: function(key: string, property: string, value: any) {
    if (!this.blocks || !this.isEditable) return;
    setFormField(this.blocks, (this.parentKey ?? '') + key, property, value);
  },
  removeField: function(key) {
    if (!this.blocks || !this.isEditable) return;
    this.blocks = removeField(this.blocks, (this.parentKey ?? '') + key);
  },
};

export const FormContext = createContext<FormContextData>(defaultFormContext);
FormContext.displayName = "FormContext";

export const useFormContext = () => useContext(FormContext);

export function FormContextProvider({ children }: { children: React.ReactNode }) {
  return <FormContext.Provider value={defaultFormContext}>
    {children}
  </FormContext.Provider>
}
