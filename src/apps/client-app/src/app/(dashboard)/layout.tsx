import type { Metadata } from 'next';

export const metadata: Metadata = {
  title: 'System Topography — Runtime Roasters',
  description: 'Visualize how user intent flows through the RuntimeRoasters microservices infrastructure stack.',
};

export default function DashboardLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="flex flex-col min-h-screen bg-[#f7f9fb]">
      <main className="flex-1 relative">
        {children}
      </main>
    </div>
  );
}
