export interface TraceEvent {
  id: string;
  message_id: string;
  topic: string;
  batch_id?: string;
  order_id?: string;
  shipment_id?: string;
  payload: string;
  occurred_at: string;
  created_at: string;
}

export interface TraceDocument {
  entity_id: string;
  entity_type: string;
  events: TraceEvent[];
  updated_at: string;
}

export interface TraceResponse {
  events: TraceEvent[];
}

export interface TraceDocumentResponse {
  document: TraceDocument;
}
