import { DndContext, DragEndEvent, DragOverlay, DragStartEvent, UniqueIdentifier } from "@dnd-kit/core";
import MainLayout from "../components/MainLayout";
import { useMemo, useState } from "react";
import { Link, useParams } from "react-router-dom";
import FormZone from "../components/forms/FormSection";
import FormBlockButton from "../components/collection_editor/FormBlockButton";
import Scaffold from "../components/Scaffold";
import { blockList } from "../components/forms/blocks";
import { FormSchema, createBlockInstance } from "../components/forms";
import { produce } from "immer";
import { findAndInsertFormBlock } from "../components/forms/editor";

export default function CollectionEditor() {
  const params = useParams();
  const [schema, setSchema] = useState<FormSchema>({ main: [] });
  const [activeDraggedBlockId, setActiveDraggedBlockId] = useState<UniqueIdentifier | null>(null);

  const activeDraggedBlock = useMemo(() => {
    if (activeDraggedBlockId) {
      for (const blockCategory of blockList) {
        for (const block of blockCategory.blocks) {
          if (block.id !== activeDraggedBlockId) {
            continue;
          }
          return block;
        }
      }
    }
    return null;
  }, [activeDraggedBlockId])

  const handleDragStart = (event: DragStartEvent) => {
    if (event.active.data.current && event.active.data.current.type === 'block-type') {
      setActiveDraggedBlockId(event.active.data.current.id);
    }
  }

  const handleDragEnd = (evt: DragEndEvent) => {
    if (!activeDraggedBlock) return;
    // console.log(evt.collisions);
    if (evt.collisions && evt.collisions.length !== 0) {
      let location = (evt.collisions[0].id as string);
      if (location === 'main' && schema[location].length === 1) {
        return;
      }

      const locationArr = location.split('.');
      locationArr.shift();
      location = locationArr.join('.');

      setSchema(
        produce(schema, draft => {
          findAndInsertFormBlock(
            draft.main,
            location,
            createBlockInstance(activeDraggedBlock));
        })
      );
    }
    setActiveDraggedBlockId(null);
  }

  return (
    <DndContext
      autoScroll={{threshold: { x: 0, y: 0.2 }, layoutShiftCompensation: false}}
      onDragStart={handleDragStart}
      onDragEnd={handleDragEnd}>
      <MainLayout
        navClassName="pt-6 pb-8"
        navigationSlot={() => (<>
          <Link to={`/sites/${params.siteId}/collections/${params.collectionId}`} className="px-6 py-2 hover:bg-slate-200 text-sm text-slate-800">
            &lt; Exit collection editor
          </Link>

          <div className="flex flex-col space-y-6 px-6 pt-4">
            {blockList.map(b => (
              <div key={`block_category_${b.id}`}>
                <span className="pb-1 text-sm font-bold text-gray-600 uppercase block">{b.title}</span>

                <div className="flex flex-col space-y-2 pt-2">
                  {b.blocks.map(block => (
                    <FormBlockButton
                      key={`block_${block.id}`}
                      id={block.id}
                      title={block.name}
                      description={block.description} />
                  ))}
                </div>
              </div>
            ))}
          </div>
          </>)}>
        <Scaffold
          leftHeader={() => <div className="flex flex-col">
            <p className="text-sm text-slate-600 block">Edit Collections</p>
            <h1 className="text-lg font-bold text-slate-800">Collection Name</h1>
          </div>}
          actions={() => (<>
            <button className="sulat-btn is-primary">Save</button>
          </>)}>
          <DragOverlay dropAnimation={null}>
            {activeDraggedBlock &&
              <FormBlockButton
                id={activeDraggedBlock.id}
                title={activeDraggedBlock.name}
                className="w-[17rem]"
                description={activeDraggedBlock.description} />}
          </DragOverlay>

          <FormZone
            name="main" children={schema.main}
            editable
            onChange={(children) => {
              setSchema(
                produce(schema, draft => {
                  draft.main = children;
                })
              );
            }} />
        </Scaffold>
      </MainLayout>
    </DndContext>
  )
}
