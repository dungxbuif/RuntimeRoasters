import { describe, expect, it } from 'vitest';
import { createProductId, expandSaleItems } from './sale';

describe('retail sale identity', () => {
  it('creates a stable valid demo product id format', () => {
    expect(createProductId('STORE-HN-HK', new Date('2026-06-06T00:00:00Z'), 'abc-123')).toBe('RR-CUP-HNHK-20260606-ABC123');
  });

  it('expands quantity to one unique product id per cup', () => {
    const items = expandSaleItems([{ menuItemId: 'menu-1', quantity: 3 }], 'HK');
    expect(items).toHaveLength(3);
    expect(new Set(items.map(item => item.product_id)).size).toBe(3);
    expect(items.every(item => item.menu_item_id === 'menu-1')).toBe(true);
  });
});
