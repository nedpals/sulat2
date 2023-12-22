export interface FormBlock<T = Record<string, any>> {
  key: string
  label: string
  type: string
  properties: T
}
