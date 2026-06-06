export interface SelectedMenuItem {
  menuItemId: string;
  quantity: number;
}

export function createProductId(storeCode: string, now = new Date(), random = crypto.randomUUID()): string {
  const store = storeCode.replace(/^STORE-/i, '').replace(/[^A-Za-z0-9]/g, '').slice(0, 8).toUpperCase() || 'STORE';
  const date = now.toISOString().slice(0, 10).replace(/-/g, '');
  const suffix = random.replace(/[^A-Za-z0-9]/g, '').slice(0, 10).toUpperCase();
  return `RR-CUP-${store}-${date}-${suffix}`;
}

export function expandSaleItems(items: SelectedMenuItem[], storeCode: string) {
  return items.flatMap(item =>
    Array.from({ length: item.quantity }, () => ({
      menu_item_id: item.menuItemId,
      product_id: createProductId(storeCode),
    }))
  );
}
