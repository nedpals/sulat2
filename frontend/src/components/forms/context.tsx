import { deleteProperty, getProperty, setProperty } from "dot-prop";
import { produce } from "immer";
import { createContext, useContext } from "react";

export interface FormContextData {
  state: Record<string, any>
  getField: (key: string, defaultValue: any) => any
  updateField: (key: string, value: any) => void
  removeField: (key: string) => void
}

export const FormContext = createContext<FormContextData>({
  state: {},
  getField: () => null,
  updateField: () => {},
  removeField: () => {},
});
FormContext.displayName = 'FormContext';

export const useFormContext = () => useContext(FormContext);
export const FormProvider = ({ state, onChange, children }: {
  state: Record<string, any>,
  onChange:(s: Record<string,any>) => void,
  children: React.ReactNode
}) => {
  return (
    <FormContext.Provider value={{
      state,
      getField(key, defaultValue) {
        return getProperty(state, key, defaultValue);
      },
      updateField(key, value) {
        onChange(produce(state, (draft) => {
          setProperty(draft, key, value);
        }));
      },
      removeField(key) {
        onChange(produce(state, (draft) => {
          deleteProperty(draft, key);
        }));
      },
    }}>
      {children}
    </FormContext.Provider>
  );
}
