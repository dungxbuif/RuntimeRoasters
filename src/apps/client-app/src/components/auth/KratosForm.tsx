'use client';

import { Button } from '@/components/ui/atoms/Button';
import { Input } from '@/components/ui/atoms/Input';
import { testId, e2eSelectors } from '@/lib/utils/test-id';
import { LoginFlow, UiNode, UpdateLoginFlowBody } from '@ory/client';
import { motion } from 'framer-motion';
import React from 'react';

interface KratosFormProps {
  flow: LoginFlow;
  onSubmit: (values: UpdateLoginFlowBody) => Promise<void>;
  isLoading?: boolean;
}

export const KratosForm: React.FC<KratosFormProps> = ({ flow, onSubmit, isLoading }) => {
  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const formData = new FormData(event.currentTarget);
    const body: Record<string, string | number | boolean> = {};
    
    formData.forEach((value, key) => {
      body[key] = value as string;
    });

    body['method'] = 'password';

    await onSubmit(body as unknown as UpdateLoginFlowBody);
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div className="space-y-4">
        {flow.ui.nodes.map((node, index) => (
          <FormNode key={`${node.group}-${index}`} node={node} />
        ))}
      </div>

      {flow.ui.messages?.map((msg) => (
        <div 
          key={msg.id} 
          className={`p-3 rounded-lg text-xs font-bold uppercase tracking-wider ${
            msg.type === 'error' ? 'bg-error/10 text-error border border-error/20' : 'bg-primary/10 text-primary border border-primary/20'
          }`}
        >
          {msg.text}
        </div>
      ))}

      <motion.div
        whileHover={{ scale: 1.02 }}
        whileTap={{ scale: 0.98 }}
      >
        <Button 
          type="submit" 
          className="w-full py-6 text-lg font-black italic uppercase tracking-[0.2em]"
          disabled={isLoading}
          {...testId(e2eSelectors.LOGIN_SUBMIT)}
        >
          {isLoading ? 'Verifying Identity...' : 'Authenticate'}
        </Button>
      </motion.div>
    </form>
  );
};

const FormNode: React.FC<{ node: UiNode }> = ({ node }) => {
  const { attributes, messages, meta } = node;

  if (attributes.node_type === 'input') {
    const inputAttr = attributes;
    
    // Skip hidden inputs like CSRF (handled by browser usually, or we can include them)
    if (inputAttr.type === 'hidden') {
      return <input type="hidden" name={inputAttr.name} value={inputAttr.value as string} />;
    }

    if (inputAttr.type === 'submit') {
      return null; // We use our own submit button
    }

    return (
      <div className="space-y-3" data-e2e={`form-field-${inputAttr.name}`}>
        <label className="text-[10px] font-black uppercase tracking-[0.2em] text-on-surface-variant opacity-80 ml-4">
          {meta.label?.text || inputAttr.name}
        </label>
        <Input
          name={inputAttr.name}
          type={inputAttr.type}
          defaultValue={inputAttr.value as string}
          placeholder={`Enter your ${meta.label?.text?.toLowerCase() || inputAttr.name}...`}
          required={inputAttr.required}
          className="h-16 shadow-inner"
          data-e2e={`input-${inputAttr.name}`}
        />
        {messages.map((msg) => (
          <p key={msg.id} className="text-[10px] text-error font-bold ml-1 italic">
            {msg.text}
          </p>
        ))}
      </div>
    );
  }

  return null;
};
