import { Navigate, useParams } from "react-router-dom";

export default function CollectionIndex() {
  const params = useParams();

  if (!params.siteId) {
    return <Navigate to="/sites" />
  }

  return <Navigate to={`/sites/${params.siteId}/collections/a`} />
}