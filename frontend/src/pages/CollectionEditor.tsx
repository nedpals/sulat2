import { DndContext, DragEndEvent, DragOverlay, DragStartEvent, UniqueIdentifier } from "@dnd-kit/core";
import MainLayout from "../components/MainLayout";
import { useEffect, useMemo, useState } from "react";
import { Link } from "react-router-dom";
import FormBlockZone from "../components/collection_editor/FormBlockZone";
import FormBlockButton from "../components/collection_editor/FormBlockButton";
import { FormContext, defaultFormContext, removeField } from "../components/collection_editor/FormContext";
import { stackBlockInfo } from "../components/collection_editor/blocks/StackBlock";
import { buttonBlockInfo } from "../components/collection_editor/blocks/ButtonBlock";
import { textBlockInfo } from "../components/collection_editor/blocks/TextBlock";
import { textareaBlockInfo } from "../components/collection_editor/blocks/TextareaBlock";
import { selectBlockInfo } from "../components/collection_editor/blocks/SelectBlock";
import { FormBlock } from "../components/collection_editor/types";

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

function createBlockFromInfo(blockInfo: BlockCategory['blocks'][0], location: string): FormBlock<Record<string, any>> {
  let properties: Record<string, any> = {};

  for (const [key, value] of Object.entries(blockInfo.propertiesSchema)) {
    properties[key] = (value as any).default ?? null;
  }
  
  return {
    label: `New ${blockInfo.name}`,
    key: blockInfo.id,
    type: blockInfo.id,
    location,
    properties
  }
}

const blockList: BlockCategory[] = [
  {
    id: 'layout',
    title: 'Layout',
    blocks: [
      stackBlockInfo
    ]
  },
  {
    id: 'buttons',
    title: 'Buttons',
    blocks: [
      buttonBlockInfo,
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
      textBlockInfo,
      textareaBlockInfo,
      selectBlockInfo,
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
    const locationArr = location.split('.').filter(Boolean);
    if (locationArr.length > 0) {
      const loc = locationArr[0];
  
      for (let i = 0; i < blocks.length; i++) {
        const _block = blocks[i];
        if (_block.key !== loc) {
          continue;
        }
  
        // If there are still keys left, we need to go deeper
        const otherKeys = locationArr.slice(1);
        if ('children' in _block.properties) {
          findAndInsertFormBlock(blocks[i].properties.children, otherKeys.join('.'), block);
        }

        break;
      }
      return;
    }
     
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

    blocks.push({
      ...block,
      location
    });

    // console.log({[location]: blocks});
  }

  const handleDragEnd = (evt: DragEndEvent) => {
    // console.log(evt.collisions);
    if (evt.collisions && evt.collisions.length !== 0) {
      let location = (evt.collisions[0].id as string)
      location = location.substring('main.'.length + Math.min(location.indexOf('.'), 0));

      setSchema(s => {
        findAndInsertFormBlock(
          s, location,
          createBlockFromInfo(activeDraggedBlock!, location)
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
        headerDisabled
        navClassName="pt-6 pb-8"
        containerClassName="px-12"
        navigationSlot={() => (<>
          <Link to="/sites" className="px-6 py-2 hover:bg-slate-200 text-sm text-slate-800">
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
        <header className="flex justify-between items-center pt-6 pb-4 border-b">
          <div>
            <p className="text-sm text-slate-600 block">Edit Collections</p>
            <h1 className="text-lg font-bold text-slate-800">Collection Name</h1>
          </div>
        
          <div className="flex space-x-4">
            <button className="sulat-btn is-primary">Save</button>
          </div>
        </header>

        <div className="pt-4">
          <DragOverlay dropAnimation={null}>
            {activeDraggedBlock && 
              <FormBlockButton 
                id={activeDraggedBlock.id} 
                title={activeDraggedBlock.name} 
                className="w-[17rem]"
                description={activeDraggedBlock.description} />}
          </DragOverlay>

          <FormContext.Provider value={defaultFormContext.copyWith({ 
            blocks: schema,
            getPropertiesSchema(key) {
              return getPropertiesSchema(key);
            },
            removeField(key) {
              setSchema(s => {
                return removeField(s, key);
              });
            },
          })}>
            <FormBlockZone max={1} editable zoneKey="main" children={schema} />
          </FormContext.Provider>
        </div>
      </MainLayout>
    </DndContext>
  )
}