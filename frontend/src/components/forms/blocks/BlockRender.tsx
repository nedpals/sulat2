import { getBlockByType } from ".";
import { FormBlockRenderProps } from "../render";

export default function BlockRender(props: FormBlockRenderProps) {
  const { block } = props;
  const BlockComponent: React.FC<FormBlockRenderProps<any>> = getBlockByType(block.type);

  if (block.properties.label) {
    return (
      <div className="flex flex-col">
        <label className="text-sm text-gray-600 mb-1">{block.properties.label}</label>
        <BlockComponent {...props} />
      </div>
    )
  }

  return <BlockComponent {...props} />
}
