export interface SearchResult {
  entity_type: string;
  entity_id: number;
  title: string;
  snippet: string;
}

export interface SearchResponse {
  projects: SearchResult[];
  tasks: SearchResult[];
  requests: SearchResult[];
}