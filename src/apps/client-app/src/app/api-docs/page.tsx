'use client';

import dynamic from 'next/dynamic';
import 'swagger-ui-react/swagger-ui.css';

// Swagger UI needs to be rendered on the client side only
const SwaggerUI = dynamic(() => import('swagger-ui-react'), { 
  ssr: false,
  loading: () => <div className="p-8 text-center">Loading API Documentation...</div>
});

export default function ApiDocsPage() {
  // Only show in development
  if (process.env.NODE_ENV === 'production') {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <h1 className="text-2xl font-bold">444 - Not Found</h1>
      </div>
    );
  }

  return (
    <div className="bg-white min-h-screen">
      <div className="p-4 bg-gray-100 border-b border-gray-200">
        <h1 className="text-xl font-bold text-gray-800">RuntimeRoasters - Public API Portal (KrakenD Export)</h1>
        <p className="text-sm text-gray-500">Generated directly from Gateway configuration</p>
      </div>
      <SwaggerUI url="/api-spec.json" />
    </div>
  );
}
