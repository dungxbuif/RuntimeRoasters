import Sidebar from "@/components/common/Sidebar";

export default function AdminLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <>
      <header className="flex justify-between items-center px-8 py-4 w-full sticky top-0 z-50 bg-slate-900 backdrop-blur-md border-b border-white/10">
        <div className="flex items-center gap-8">
          <span className="text-2xl font-black text-white tracking-tighter font-headline uppercase italic">
            Runtime Roasters <span className="text-primary-fixed ml-2 font-black not-italic">[Admin Portal]</span>
          </span>
        </div>
      </header>

      <div className="flex flex-1">
        <Sidebar />
        <main className="ml-64 flex-1 bg-surface min-h-screen relative overflow-x-hidden">
          <div className="p-8">
            {children}
          </div>
          <div className="absolute top-0 right-0 w-1/3 h-1/3 bg-primary/10 rounded-full blur-[120px] pointer-events-none z-0 mix-blend-multiply opacity-30"></div>
        </main>
      </div>
    </>
  );
}
