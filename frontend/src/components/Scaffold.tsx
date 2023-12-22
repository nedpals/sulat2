import { cn } from "../utils";

export default function Scaffold({ leftHeader: LeftHeader, actions: Actions, children, className }: {
  leftHeader?: React.FC
  actions?: React.FC
  children: React.ReactNode
  className?: string
}) {
  return (
    <div className="px-6 md:px-8 lg:px-12">
      {(LeftHeader || Actions) && <header className="flex justify-between items-center py-3 border-b">
        {LeftHeader && <div className="justify-self-start">
          <LeftHeader />
        </div>}
        {Actions && <div className="justify-end items-center flex space-x-4">
          <Actions />
        </div>}
      </header>}

      <section className={cn('pt-4', className)}>
        {children}
      </section>
    </div>
  );
}
