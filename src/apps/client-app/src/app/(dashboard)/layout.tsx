import Sidebar from "@/components/common/Sidebar";

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <>
      <header className="flex justify-between items-center px-8 py-4 w-full sticky top-0 z-50 bg-slate-100/50 backdrop-blur-md border-b border-outline-variant/10">
        <div className="flex items-center gap-8">
          <span className="text-2xl font-black text-blue-600 tracking-tighter font-headline uppercase italic text-glow">Runtime Roasters</span>
        </div>
      </header>

      <div className="flex flex-1">
        <Sidebar />
        <main className="ml-64 flex-1 bg-surface min-h-screen relative overflow-x-hidden">
          {children}
          <div className="absolute top-0 right-0 w-1/3 h-1/3 bg-primary-fixed/20 rounded-full blur-[120px] pointer-events-none z-0 mix-blend-multiply opacity-30"></div>
        </main>
      </div>
    </>
  );
}
