import { DndContext, DragEndEvent, DragOverlay, DragStartEvent, UniqueIdentifier } from "@dnd-kit/core";
import MainLayout from "../components/MainLayout";
import { useEffect, useMemo, useState } from "react";
import { Link, useParams } from "react-router-dom";
import FormBlockZone from "../components/collection_editor/FormBlockZone";
import FormBlockButton from "../components/collection_editor/FormBlockButton";
import { FormContext, copyContextWith, defaultFormContext } from "../components/collection_editor/FormContext";
import StackBlock from "../components/collection_editor/blocks/StackBlock";
import ButtonBlock from "../components/collection_editor/blocks/ButtonBlock";
import TextBlock from "../components/collection_editor/blocks/TextBlock";
import SelectBlock from "../components/collection_editor/blocks/SelectBlock";
import { FormBlock } from "../components/collection_editor/types";
import Scaffold from "../components/Scaffold";
import TextareaBlock from "../components/collection_editor/blocks/TextareaBlock";

interface BlockCategory {
  id: string
  title: string
  blocks: {
    id: string
    name: string
    description: string
    propertiesSchema: Record<string, any>
  }[]
}

function createBlockFromInfo(blockInfo: BlockCategory['blocks'][0]): FormBlock {
  const properties: Record<string, any> = {};

  for (const [key, value] of Object.entries(blockInfo.propertiesSchema)) {
    properties[key] = value.default ?? null;
  }

  return {
    label: `New ${blockInfo.name}`,
    key: blockInfo.id,
    type: blockInfo.id,
    properties
  }
}

const blockList: BlockCategory[] = [
  {
    id: 'layout',
    title: 'Layout',
    blocks: [
      StackBlock.properties
    ]
  },
  {
    id: 'buttons',
    title: 'Buttons',
    blocks: [
      ButtonBlock.properties,
      // {
      //   id: 'link',
      //   label: 'Link',
      //   description: 'A link',
      //   payload: {
      //     key: 'link',
      //     type: 'link',
      //     properties: {
      //       text: 'Link',
      //       href: '#'
      //     }
      //   }
      // }
    ]
  },
  {
    id: 'input',
    title: 'Input',
    blocks: [
      TextBlock.properties,
      TextareaBlock.properties,
      SelectBlock.properties,
      // {
      //   id: 'checkbox',
      //   label: 'Checkbox',
      //   description: 'A checkbox input',
      //   payload: {
      //     key: 'checkbox',
      //     type: 'checkbox',
      //     properties: {
      //       label: 'Checkbox'
      //     }
      //   }
      // },
      // {
      //   id: 'radio',
      //   label: 'Radio',
      //   description: 'A radio input',
      //   payload: {
      //     key: 'radio',
      //     type: 'radio',
      //     properties: {
      //       label: 'Radio',
      //       options: [
      //         { label: 'Option 1', value: 'option_1' },
      //         { label: 'Option 2', value: 'option_2' },
      //         { label: 'Option 3', value: 'option_3' },
      //       ]
      //     }
      //   }
      // },
      // {
      //   id: 'file',
      //   label: 'File',
      //   description: 'A file input',
      //   payload: {
      //     key: 'file',
      //     type: 'file',
      //     properties: {
      //       label: 'File',
      //       placeholder: 'File'
      //     }
      //   }
      // },
      // {
      //   id: 'image',
      //   label: 'Image',
      //   description: 'An image input',
      //   payload: {
      //     key: 'image',
      //     type: 'image',
      //     properties: {
      //       label: 'Image',
      //       placeholder: 'Image'
      //     }
      //   }
      // },
    ]
  }
];

function getPropertiesSchema(blockId: string): Record<string, any> {
  for (const blockCategory of blockList) {
    for (const block of blockCategory.blocks) {
      if (block.id !== blockId) {
        continue;
      }
      return block.propertiesSchema;
    }
  }

  return {};
}

export default function CollectionEditor() {
  const params = useParams();
  const [schema, setSchema] = useState<FormBlock[]>([]);
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
    setActiveDraggedBlockId(event.active.id);
  }

  const findAndInsertFormBlock = (blocks: FormBlock[], location: string, block: FormBlock): void => {
    if (location.length === 0) {
      const existingBlock = blocks.find(b => b.key === block.key);
      if (existingBlock) {
        for (let i = 1;; i++) {
          const newBlockKey = `${block.key}_${i}`;
          const existingBlock = blocks.find(b => b.key === newBlockKey);
          if (existingBlock) {
            continue;
          }

          block.key = newBlockKey;
          break;
        }
      }

      blocks.push(block);
      return;
    }

    const locationArr = location.split('.').filter(Boolean);
    const loc = locationArr.shift();
    if (typeof loc === 'undefined') {
      return;
    }

    for (let i = 0; i < blocks.length; i++) {
      const _block = blocks[i];
      if (_block.key !== loc) {
        continue;
      }

      // If there are still keys left, we need to go deeper
      if ('children' in _block.properties) {
        const otherKeys = locationArr.join('.');
        findAndInsertFormBlock(blocks[i].properties.children, otherKeys, block);
      }

      break;
    }
  }

  const handleDragEnd = (evt: DragEndEvent) => {
    // console.log(evt.collisions);
    if (evt.collisions && evt.collisions.length !== 0) {
      let location = (evt.collisions[0].id as string)
      location = location.substring('main.'.length + Math.min(location.indexOf('.'), 0));

      setSchema(s => {
        findAndInsertFormBlock(
          s, location,
          createBlockFromInfo(activeDraggedBlock!)
        );
        return s;
      });
    }

    setActiveDraggedBlockId(null);
  }

  useEffect(() => {
    console.log(schema);
  }, [schema]);

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

          <FormContext.Provider value={copyContextWith(defaultFormContext, { blocks: schema, getBlockSchema: getPropertiesSchema })}>
            <FormBlockZone max={1} editable zoneKey="main" children={schema} />
          </FormContext.Provider>
        </Scaffold>
      </MainLayout>
    </DndContext>
  )
}
