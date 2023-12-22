import { Link, To } from "react-router-dom";
import { cn } from "../utils";

export default function MainLayout({
  defaultLink,
  currentSiteSlot: CurrentSiteSlot,
  navigationSlot: NavigationSlot,
  children,
  navClassName,
  containerClassName
}: {
  defaultLink?: To
  currentSiteSlot?: React.FC
  navigationSlot?: React.FC
  children: React.ReactNode
  navClassName?: string
  containerClassName?: string
}) {
  return (
    <div>
      <header className="flex items-center bg-slate-900 text-white pr-4 shadow">
        <Link to={defaultLink ?? '/'} className="py-4 px-4 hover:bg-white hover:text-slate-900 transition-colors">
          Sulat<span className="font-bold">CMS</span>
        </Link>

        <div className="mx-auto w-full flex">
          <input
            type="text"
            className="border border-white/20 bg-white/20 hover:border-white/30 hover:bg-white/30 focus:bg-white focus:text-slate-900 focus:border-white w-1/2 mx-auto transition-colors px-6 py-2 rounded outline-none"
            placeholder="Search for collections, records, and sites..." />
        </div>
      </header>

      <div className="flex">
        <div className={cn('w-[15rem] max-w-[15rem] lg:w-[20rem] lg:max-w-[20rem] fixed left-0 inset-b-0 overscroll-x-contain h-screen max-h-screen overflow-x-hidden overflow-y-auto flex flex-col bg-slate-100', navClassName)}>
          {/* {CurrentSiteSlot && <div className="px-6 pt-8 flex flex-col pb-3">
            <CurrentSiteSlot />
          </div>} */}
          {NavigationSlot && <NavigationSlot />}
        </div>
        <div className={cn('ml-[15rem] lg:ml-[20rem] flex-1', containerClassName)}>
          {children}
        </div>
      </div>
    </div>
  );
}
