import { Navigate, useParams } from "react-router-dom";

export default function SiteIndex() {
  const params = useParams();

  return <Navigate to={`/sites/${params.siteId}/collections/a`} />
}