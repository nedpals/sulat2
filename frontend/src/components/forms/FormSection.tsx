import EditableFormSection from "./EditableFormSection";
import { useCallback, useMemo } from "react";
import { cn } from "../../utils";
import { FormBlockInstance } from "./blocks";
import FormBlockRenderer from "./FormBlockRenderer";
import { FormSectionProvider, useFormSectionContext } from ".";

export default function FormZone({
  name,
  parent,
  children = [],
  containerClassName,
  editable,
  childWrapper,
  onChange
}: {
  name: string,
  parent?: string,
  children?: FormBlockInstance[],
  childWrapper?: React.FC<{children: React.ReactNode, idx: number}>,
  containerClassName?: string,
  editable?: boolean,
  onChange?: (blocks: FormBlockInstance[]) => void
}) {
  const sectionContext = useFormSectionContext();
  // const memoChildren = useDeepCompareMemoize(children, []);
  const memoChildren = children;

  const combinedParent = useMemo(() => {
    if (sectionContext.parentKey) {
      return [sectionContext.parentKey, parent].filter(Boolean).join('.');
    }
    return parent ?? '';
  }, [sectionContext, parent]);

  const combinedName = useMemo(() => {
    if (sectionContext.name) {
      return [sectionContext.name, name].join('.');
    }
    return name;
  }, [sectionContext, name]);

  const isEditable = useMemo(() => {
    return editable ?? sectionContext.isEditable;
  }, [sectionContext, editable]);

  const ChildWrapper = useCallback(
    childWrapper ?? (({children}: {children: React.ReactNode}) => <>{children}</>
  ), []);

  return (
    <FormSectionProvider value={{
      name: combinedName,
      parentKey: combinedParent,
      isEditable,
      onChange,
      children: memoChildren,
    }}>
      {isEditable ? (
        <EditableFormSection
          containerClassName={containerClassName}
          childWrapper={ChildWrapper} />
      ) : (
        <div className={cn("flex flex-col", containerClassName)}>
          {memoChildren?.map((block, idx) => {
            const key = `block\$\$${combinedParent}_${block.type}_${combinedName}.${block.key}`;
            return <ChildWrapper idx={idx} key={key}>
              <FormBlockRenderer block={block} />
            </ChildWrapper>;
          })}
        </div>
      )}
    </FormSectionProvider>
  )
}
