import { FormBlockInstance } from "./blocks"

export interface FormBlockRenderProps<T = Record<string, any>> {
  block: FormBlockInstance<T>
  className?: string
}
