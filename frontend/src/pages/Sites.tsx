import { Link } from "react-router-dom";

export default function Sites() {
  const sites = [
    { id: 'a', name: 'Site A' },
    { id: 'b', name: 'Site B' },
    { id: 'c', name: 'Site C' },
    { id: 'd', name: 'Site D' },
    { id: 'e', name: 'Site E' },
  ]

  return (
    <div className="bg-slate-100 min-h-screen">
      <div className="max-w-2xl py-8 mx-auto text-center flex flex-col items-center text-slate-800">
        <p className="text-xl">SulatCMS</p>
        <h1 className="text-3xl font-bold">Choose a Site</h1>

        <div className="flex flex-col pt-8 space-y-2 w-3/4">
          {sites.map(s => (
            <Link to={`/sites/${s.id}`} key={`site_${s.id}`}
              className="bg-white border shadow rounded flex space-x-4 items-center py-3 px-6 hover:bg-slate-600/[.01]">
              <div className="h-8 w-8 rounded bg-slate-300"></div>
              <div className="flex flex-col items-start">
                <span className="text-sm">{s.name}</span>
                <p className="text-xs text-slate-500">Site type goes here</p>
              </div>
            </Link>
          )
          )}
          <button className="bg-white border text-sm font-semibold text-center w-full justify-center shadow rounded flex space-x-4 items-center py-3 px-6 hover:bg-slate-600/[.01]">
            Add site
          </button>
        </div>
      </div>
    </div>
  )
}
