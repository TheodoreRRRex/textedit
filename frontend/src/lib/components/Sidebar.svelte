<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { PickDirectoryDialog, ListDirectory } from '../../../wailsjs/go/main/App.js';
  import { sidebarRootPath, currentFilePath } from '../stores/editor';
  import FileTree from './FileTree.svelte';
  import type { DirEntry } from '../types';

  const dispatch = createEventDispatcher<{ fileSelect: string }>();

  let entries: DirEntry[] = [];

  async function openFolder() {
    try {
      const path = await PickDirectoryDialog();
      if (path) {
        $sidebarRootPath = path;
        entries = await ListDirectory(path);
      }
    } catch (e) {
      alert(`Failed to open folder: ${e}`);
    }
  }

  function handleFileClick(event: CustomEvent<string>) {
    dispatch('fileSelect', event.detail);
  }
</script>

<aside class="sidebar">
  <div class="sidebar-header">
    <button class="open-folder-btn" on:click={openFolder}>
      Open Folder
    </button>
  </div>
  <div class="file-tree-container">
    {#if $sidebarRootPath}
      <div class="root-label" title={$sidebarRootPath}>
        {$sidebarRootPath.split('/').pop() || $sidebarRootPath}
      </div>
      <FileTree {entries} activePath={$currentFilePath} on:fileClick={handleFileClick} />
    {:else}
      <div class="empty-state">
        No folder open
      </div>
    {/if}
  </div>
</aside>

<style>
  .sidebar {
    background: #252526;
    display: flex;
    flex-direction: column;
    overflow: hidden;
    border-right: 1px solid #3c3c3c;
    grid-row: 1 / 3;
  }

  .sidebar-header {
    padding: 8px;
    border-bottom: 1px solid #3c3c3c;
  }

  .open-folder-btn {
    width: 100%;
    padding: 6px 12px;
    background: #0e639c;
    color: #fff;
    border: none;
    border-radius: 4px;
    font-size: 12px;
  }

  .open-folder-btn:hover {
    background: #1177bb;
  }

  .file-tree-container {
    flex: 1;
    overflow-y: auto;
    padding: 4px 0;
  }

  .root-label {
    padding: 4px 12px;
    font-size: 11px;
    font-weight: bold;
    text-transform: uppercase;
    color: #bbbbbb;
    letter-spacing: 0.5px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .empty-state {
    padding: 20px 12px;
    color: #888;
    font-size: 12px;
    text-align: center;
  }
</style>
