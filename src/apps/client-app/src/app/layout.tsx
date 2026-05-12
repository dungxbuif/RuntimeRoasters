import Providers from "@/components/common/Providers";
import type { Metadata } from "next";
import { Gaegu, Indie_Flower, JetBrains_Mono, Roboto_Mono } from "next/font/google";
import "./globals.css";

const jetbrainsMono = JetBrains_Mono({ subsets: ["latin"], variable: "--font-headline" });
const robotoMono = Roboto_Mono({ subsets: ["latin"], variable: "--font-body" });
const gaegu = Gaegu({ weight: ["400", "700"], subsets: ["latin"], variable: "--font-marker" });
const indieFlower = Indie_Flower({ weight: ["400"], subsets: ["latin"], variable: "--font-annotation" });

export const metadata: Metadata = {
  title: "Runtime Roasters",
  description: "Advanced Supply Chain Management System",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className={`h-full antialiased ${jetbrainsMono.variable} ${robotoMono.variable} ${gaegu.variable} ${indieFlower.variable}`}>
      <head>
        {/* eslint-disable-next-line @next/next/no-page-custom-font */}
        <link href="https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined:wght,FILL@100..700,0..1&display=swap" rel="stylesheet" />
      </head>
      <body 
        className="bg-surface text-on-surface font-body selection:bg-primary/20 selection:text-primary min-h-screen flex flex-col"
        suppressHydrationWarning
      >
        <Providers>
          {children}
        </Providers>
      </body>
    </html>
  );
}
