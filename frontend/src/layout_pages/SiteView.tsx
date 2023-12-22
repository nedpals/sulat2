import { Link, Outlet, useParams } from "react-router-dom"
import { cn } from "../utils"
import MainLayout from "../components/MainLayout";

export default function SiteView() {
  const params = useParams();

  const collections = [
    { id: 'a', name: 'Collection A' },
    { id: 'b', name: 'Collection B' },
    { id: 'c', name: 'Collection C' },
    { id: 'd', name: 'Collection D' },
    { id: 'e', name: 'Collection E' },
  ]

  return (
    <MainLayout
      defaultLink="/sites"
      currentSiteSlot={() => (
        <a href="#" className="text-lg font-bold hover:underline">Site Name</a>
      )}
      navigationSlot={() => (<>
        <div className="flex flex-col pt-6">
          <span className="px-6 pb-1 text-sm font-bold text-gray-600 uppercase block">Collections</span>
          {collections.map(c => (
            <Link
              to={`/sites/${params.siteId}/collections/${c.id}`}
              key={`collection_${c.id}`}
              className={cn(
                'flex space-x-4 items-center py-3 px-6 hover:bg-slate-600/5',
                {
                  'bg-slate-600/5 hover:bg-slate-600/10': params.collectionId && params.collectionId === c.id
                }
              )}>
              <div className="h-8 w-8 rounded bg-slate-300"></div>
              <span className="text-sm">{c.name}</span>
            </Link>
          ))}

          <button className="flex space-x-4 items-center mt-4 py-3 px-6 hover:bg-slate-200">
            <div className="h-8 w-8 rounded bg-slate-300"></div>
            <span className="text-sm">Add new collection</span>
          </button>
        </div>
      </>)}>
      <Outlet />
    </MainLayout>
  )
}
