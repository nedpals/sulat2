import { createContext, useContext } from "react"
import { FormBlock } from "./types"

export interface FormContext {
  isEditable: boolean
  parentKey?: string
  blocks?: FormBlock[]
  values: Record<string, any>
  getPropertiesSchema(key: string): Record<string, any>
  getFieldValue(key: string, defaultValue: any): any
  setFieldValue(key: string, value: any): void
  setField(key: string, property: string, value: any): void
  removeField(key: string): void
  copyWith(f: ModifiableFormContextFields): FormContext
}

type ModifiableFormContextFields = Partial<FormContext>;

function copyContextWith(oldContext: FormContext, context: ModifiableFormContextFields): FormContext {
  return {
    ...oldContext,
    ...context,
  }
}

function setFormField(blocks: FormBlock[], key: string, property: string, value: any) {
  // Find block first
  const keyArr = key.split('.');
  const blockKey = keyArr.shift();

  if (!blockKey) {
    return;
  }

  for (let i = 0; i < blocks.length; i++) {
    const block = blocks[i];
    if (block.key !== blockKey) {
      continue;
    }

    if (keyArr.length === 0) {
      block.properties[property] = value;
      return;
    }

    if ('children' in block.properties && Array.isArray(block.properties.children)) {
      return setFormField(block.properties.children, keyArr.join('.'), property, value);
    }
  }
}

function getFieldValue(values: Record<string, any>, key: string, defaultValue: any = null): any {
  const keyArr = key.split('.');
  const blockKey = keyArr.shift();

  if (!blockKey || !values[blockKey]) {
    return defaultValue;
  }

  if (keyArr.length === 0) {
    return values[blockKey];
  }

  if (Array.isArray(values[blockKey])) {
    const nextKey = keyArr.shift();
    if (!nextKey) {
      return defaultValue;
    }

    const index = parseInt(nextKey);
    if (isNaN(index)) {
      return defaultValue;
    }

    return getFieldValue(values[blockKey][index], keyArr.join('.'), defaultValue);
  }

  return getFieldValue(values[blockKey], keyArr.join('.'), defaultValue);
}

function setFormFieldValue(values: Record<string, any>, key: string, value: any) {
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

export const defaultFormContext: FormContext = {
  isEditable: false,
  values: {},
  getPropertiesSchema: function(_: string) {
    return {};
  },
  getFieldValue: function(key, defaultValue = null) {
    console.log('GET FIELD VALUE', key, defaultValue, this);
    if (!this.values) {
      return defaultValue;
    }
    return getFieldValue(this.values, key, defaultValue);
  },
  setFieldValue: function(key: string, val: any) { 
    setFormFieldValue(this.values, key, val); 
  },
  setField: function(key, property, value) {
    if (!this.blocks || !this.isEditable) return;
    setFormField(this.blocks, (this.parentKey ?? '') + key, property, value);
  },
  removeField: function(key) {
    if (!this.blocks || !this.isEditable) return;
    this.blocks = removeField(this.blocks, (this.parentKey ?? '') + key);
  },
  copyWith: function(context: ModifiableFormContextFields) { 
    return copyContextWith(this, context);
  },
};

export const FormContext = createContext<FormContext>(defaultFormContext);

export const useFormContext = () => useContext(FormContext);