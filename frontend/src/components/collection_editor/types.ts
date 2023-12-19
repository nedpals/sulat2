export interface FormBlock<T = Record<string, any>> {
  key: string
  label: string
  type: string
  location?: string
  properties: T
}