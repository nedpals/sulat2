import { useFormSectionContext } from "."
import { FormBlockRenderProps } from "./render"
import BlockRender from "./blocks/BlockRender"
import EditableFormBlockRenderer from "./EditableFormBlockRenderer"

export default function FormBlockRenderer(props: FormBlockRenderProps) {
  const { isEditable } = useFormSectionContext();

  if (isEditable) {
    return <EditableFormBlockRenderer {...props} />
  }

  return <BlockRender {...props} className={props.className} />
}
