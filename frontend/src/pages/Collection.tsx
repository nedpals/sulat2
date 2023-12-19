import { useState } from 'react'
import { Link, useParams } from 'react-router-dom';

function Collection() {
  const params = useParams();

  const records = [
    {
      id: '1.md',
      data: {
        title: 'Record 1',
        lastModified: '2021-04-01T00:00:00Z',
      },
    },
    {
      id: '2.md',
      data: {
        title: 'Record 2',
        lastModified: '2021-04-02T00:00:00Z',
      },
    },
    {
      id: '3.md',
      data: {
        title: 'Record 3',
        lastModified: '2021-04-03T00:00:00Z',
      },
    },
    {
      id: '4.md',
      data: {
        title: 'Record 4',
        lastModified: '2021-04-04T00:00:00Z',
      },
    },
    {
      id: '5.md',
      data: {
        title: 'Record 5',
        lastModified: '2021-04-05T00:00:00Z',
      },
    },
  ]

  return (
    <div className="px-12">
      <header className="flex justify-between items-center py-6 border-b">
        <h1 className="text-lg font-bold text-slate-800">Collection Name</h1>
      
        <div className="flex space-x-4">
          <Link to={`/sites/${params.siteId}/collections/a/edit`} className="sulat-btn">Edit collection</Link>
          <button className="sulat-btn is-primary">Add new</button>
        </div>
      </header>

      <section className="py-3">
        <div className="rounded border pl-8">
          <input type="text" className="bg-none rounded-r border-none w-full py-3 outline-none" placeholder="Search for records..." />
        </div>
      </section>

      <section>
        <div className="flex border-b pt-4 pb-2 text-sm uppercase font-medium text-slate-600">
          <div className="w-2/3 pl-20 pr-2">Name / ID</div>
          <div className="w-1/3 px-2">Last modified</div>
        </div>

        <div className="flex flex-col divide-y">
          {records.map(r => (
            <Link to={`/sites/${params.siteId}/collections/a/${r.id}`} key={`record_${r.id}`} className="group flex py-4 text-sm text-slate-800 hover:bg-slate-100">
              <div className="flex space-x-4 items-center w-2/3 pl-8 pr-2">
                <div className="h-8 w-8 bg-violet-500 rounded"></div>
                <div className="flex flex-col font-bold group-hover:underline">
                  <span>{r.data.title}</span>
                  <span className="text-xs text-slate-500 font-normal">{r.id}</span>
                </div>
              </div>
              <div className="w-1/3 self-center px-2">{r.data.lastModified}</div>
            </Link>
          ))}
        </div>
      </section>
    </div>
  );
}

export default Collection
