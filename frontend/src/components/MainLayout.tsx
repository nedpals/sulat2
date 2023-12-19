import { Link, To } from "react-router-dom";
import { cn } from "../utils";

export default function MainLayout({ 
  defaultLink, 
  currentSiteSlot: CurrentSiteSlot, 
  navigationSlot: NavigationSlot, 
  headerDisabled = false,
  children, 
  navClassName,
  containerClassName 
}: {
  defaultLink?: To
  headerDisabled?: boolean
  currentSiteSlot?: React.FC
  navigationSlot?: React.FC
  children: React.ReactNode
  navClassName?: string
  containerClassName?: string
}) {
  return (
    <div className="flex">
      <div className={cn('w-[20rem] max-w-[20rem] fixed left-0 inset-y-0 overscroll-x-contain h-screen max-h-screen overflow-x-hidden overflow-y-auto flex flex-col bg-slate-100', navClassName)}>
        {!headerDisabled && <div className="px-6 pt-8 flex flex-col pb-3">
          <Link to={defaultLink ?? '/'} className="text-xl">SulatCMS</Link>
          {CurrentSiteSlot && <CurrentSiteSlot />}
        </div>}

        {NavigationSlot && <NavigationSlot />}
      </div>

      <div className={cn('ml-[20rem] flex-1', containerClassName)}>
        {children}
      </div>
    </div>
  );
}