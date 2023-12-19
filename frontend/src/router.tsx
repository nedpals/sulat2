import { Navigate, createBrowserRouter } from 'react-router-dom'

// Pages
import App from './App'
import Collection from './pages/Collection'
import CollectionEditor from './pages/CollectionEditor'
import CollectionIndex from './pages/CollectionIndex'
import RecordEditor from './pages/RecordEditor'
import Sites from './pages/Sites'
import SiteIndex from './pages/SiteIndex'

// Layout pages
import SiteView from './layout_pages/SiteView'

export const router = createBrowserRouter([
  {
    path: '/',
    element: <App />,
    children: [
      {
        path: '/',
        element: <Navigate to="/sites" />
      },
      {
        path: 'sites',
        children: [
          {
            index: true,
            element: <Sites />,
          },
          {
            path: ':siteId',
            element: <SiteView />,
            children: [
              {
                index: true,
                element: <SiteIndex />
              },
              {
                path: 'collections',
                element: <CollectionIndex />
              },
              {
                path: 'collections/:collectionId',
                children: [
                  {
                    index: true,
                    element: <Collection />
                  },
                  {
                    path: ':recordId',
                    element: <RecordEditor />
                  }
                ]
              },
            ]
          },
        ]
      },
      {
        path: '/sites/:siteId/collections/:collectionId/edit',
        element: <CollectionEditor />
      },
    ]
  },
])