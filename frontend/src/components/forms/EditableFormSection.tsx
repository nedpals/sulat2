import { useDroppable } from "@dnd-kit/core";
import { cn } from "../../utils";
import FormBlockRenderer from "./FormBlockRenderer";
import { FormEditorProvider, useFormEditorContext } from "./editor";
import { FormBlockInstance } from "./blocks";
import { produce } from "immer";
import { getProperty } from "dot-prop";
import { useFormSectionContext } from ".";

function getParent(pKey: string, draft: FormBlockInstance[]) {
  if (!pKey) return draft;

  let _pKeyArr = pKey.split('.');

  // the first key when located is in an array so we find the index of it instead
  const parentIdx = draft.findIndex(c => c.properties.children && c.key === _pKeyArr[0]);
  if (parentIdx === -1) return;

  // replace the first key with the index
  _pKeyArr[0] = `[${parentIdx}]`;

  const _pKey = _pKeyArr.map((k, i) => i === 0 ? k : k + '.properties.children').join('.');
  const parentBlock = getProperty(draft, _pKey) as FormBlockInstance | undefined;
  if (!parentBlock || !parentBlock.properties.children) return;

  return parentBlock.properties.children;
}

export default function EditableFormSection({ containerClassName, childWrapper: ChildWrapper }: {
  containerClassName?: string,
  childWrapper: React.FC<{children: React.ReactNode, idx: number}>,
}) {
  const { name, parentKey, onChange, children } = useFormSectionContext();
  const context = useFormEditorContext();
  const {isOver, setNodeRef} = useDroppable({
    id: name
  });

  return (<FormEditorProvider value={{
    updateBlock: (pKey, key, rawProperties) => {
      if (onChange) {
        return onChange(produce(children, draft => {
          let properties = rawProperties;
          let hasFieldKey = false;

          if (rawProperties && 'fieldKey' in rawProperties) {
            hasFieldKey = true;
            properties = rawProperties.properties;
          }

          const parent = getParent(pKey ?? '', draft);
          if (!parent) return;

          const block = parent.find((b: FormBlockInstance) => b.key === key);
          if (!block) return;

          if (hasFieldKey) {
            block.fieldKey = rawProperties.fieldKey;
          }

          if (properties) {
            for (const propKey in properties) {
              block.properties[propKey] = properties[propKey];
            }
          }
        }));
      }

      return context.updateBlock(pKey ?? parentKey, key, rawProperties);
    },
    removeBlock: (pKey, key) => {
      if (onChange) {
        return onChange(produce(children, draft => {
          const parent = getParent(pKey ?? '', draft);
          const blockIndex = parent.findIndex((b: FormBlockInstance) => b.key === key);

          if (blockIndex === -1) return;
          parent.splice(blockIndex, 1);
        }));
      }
      return context.removeBlock(pKey ?? parentKey, key);
    },
  }}>
    <div ref={setNodeRef} className={cn('flex flex-col', 'border rounded min-h-[24rem] h-full w-full', containerClassName, {
      'border-violet-500 rounded': isOver,
      'border-opacity-5': children.length !== 0,
    })}>
      {children?.map((block, idx) => {
        const key = `block\$\$${parentKey}_${block.type}_${name}.${block.key}`;
        return <ChildWrapper idx={idx} key={key}>
          <FormBlockRenderer block={block} />
        </ChildWrapper>;
      })}

      {children.length === 0 && (
        <div className="rounded py-8 flex flex-col items-center justify-center h-full w-full flex-1">
          <div className="text-sm text-gray-600">Drop a block here</div>
        </div>
      )}
    </div>
  </FormEditorProvider>);
}
