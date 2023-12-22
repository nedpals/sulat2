import { Link, useParams } from 'react-router-dom';
import Scaffold from '../components/Scaffold';

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
    <Scaffold
      className="pt-0"
      leftHeader={() => (<div className="flex space-x-2 items-center self-start">
        <h1 className="text-lg font-bold text-slate-800">Collection Name</h1>
        <Link to={`/sites/${params.siteId}/collections/a/edit`} className="sulat-btn is-small">Edit</Link>
      </div>)}
      actions={() => (<>
        <button className="sulat-btn is-primary">Add new</button>
      </>)}>
      <div className="flex border-b lg:pl-[4.5rem] text-sm uppercase font-medium text-slate-600">
        <div className="w-2/3 pt-4 pb-2 px-2 mt-auto hover:bg-slate-600/10">Name / ID</div>
        <div className="w-1/3 pt-4 pb-2 px-2 mt-auto hover:bg-slate-600/10">Last modified</div>
      </div>

      <div className="flex flex-col divide-y">
        {records.map(r => (
          <Link
            to={`/sites/${params.siteId}/collections/${params.collectionId}/${r.id}`} key={`record_${r.id}`}
            className="group flex flex-nowrap py-4 text-sm text-slate-800 hover:bg-slate-100">
            <div className="ml-2 lg:ml-10 h-8 w-8 bg-violet-500 rounded"></div>
            <div className="flex space-x-4 items-center w-2/3 px-3">
              <div className="flex flex-col font-bold group-hover:underline">
                <span>{r.data.title}</span>
                <span className="text-xs text-slate-500 font-normal">{r.id}</span>
              </div>
            </div>
            <div className="w-1/3 self-center px-2">{r.data.lastModified}</div>
          </Link>
        ))}
      </div>
    </Scaffold>
  );
}

export default Collection
