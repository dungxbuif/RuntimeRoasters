import type { Metadata } from "next";
import "./globals.css";
import Sidebar from "@/components/common/Sidebar";
import Providers from "@/components/common/Providers";

export const metadata: Metadata = {
  title: "ArchitectureNarrator - Runtime Roasters",
  description: "Advanced Supply Chain Management System",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className="h-full antialiased">
      <head>
        <link rel="preconnect" href="https://fonts.googleapis.com" />
        <link rel="preconnect" href="https://fonts.gstatic.com" crossOrigin="anonymous" />
        <link href="https://fonts.googleapis.com/css2?family=Space+Grotesk:wght@300;400;500;600;700;900&family=Inter:wght@300;400;500;600;700&family=Gaegu:wght@400;700&family=Indie+Flower&display=swap" rel="stylesheet" />
        <link href="https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined:wght,FILL@100..700,0..1&display=swap" rel="stylesheet" />
      </head>
      <body className="bg-surface text-on-surface font-body selection:bg-tertiary-fixed min-h-screen flex flex-col">
        <header className="flex justify-between items-center px-8 py-4 w-full sticky top-0 z-50 bg-slate-100/50 backdrop-blur-md border-b border-outline-variant/10">
           {/* ... Giữ nguyên Header UI ... */}
           <div className="flex items-center gap-8">
            <span className="text-2xl font-black text-blue-600 tracking-tighter font-headline uppercase italic text-glow">ArchitectureNarrator</span>
          </div>
        </header>

        <div className="flex flex-1">
          <Sidebar />
          <main className="ml-64 flex-1 bg-surface min-h-screen relative overflow-x-hidden">
            <Providers>
              {children}
            </Providers>
            <div className="absolute top-0 right-0 w-1/3 h-1/3 bg-primary-fixed/20 rounded-full blur-[120px] pointer-events-none z-0 mix-blend-multiply opacity-30"></div>
          </main>
        </div>
      </body>
    </html>
  );
}
