import { RouterProvider as AriaRouterProvider } from 'react-aria';
import { Outlet, useNavigate } from 'react-router-dom';

export default function App() {
  const navigate = useNavigate();
    
  return (
    <AriaRouterProvider navigate={navigate}>
      <Outlet />
    </AriaRouterProvider>
  );
}