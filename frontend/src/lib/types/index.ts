export interface FileResult {
  path: string;
  name: string;
  content: string;
  size: number;
  readOnly: boolean;
}

export interface DirEntry {
  name: string;
  path: string;
  isDir: boolean;
  ext: string;
  size: number;
}
