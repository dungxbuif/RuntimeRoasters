'use client';

import React from 'react';

export default function ProfilePage() {
  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-4xl font-black font-headline text-on-surface uppercase italic tracking-tighter">
          Account <span className="text-primary">Settings</span>
        </h1>
        <p className="text-on-surface-variant font-medium">Update your security credentials and personal information.</p>
      </div>

      <div className="bg-surface-container-lowest rounded-[2rem] border border-outline-variant/10 p-12 text-center">
        <div className="w-24 h-24 bg-primary/10 rounded-full flex items-center justify-center mx-auto mb-6 text-primary">
          <span className="material-symbols-outlined text-4xl">construction</span>
        </div>
        <h2 className="text-xl font-bold uppercase tracking-tight mb-2">Under Construction</h2>
        <p className="text-on-surface-variant max-w-md mx-auto text-sm">
          The Profile Management module is being integrated with Ory Kratos Self-Service UI.
        </p>
      </div>
    </div>
  );
}
