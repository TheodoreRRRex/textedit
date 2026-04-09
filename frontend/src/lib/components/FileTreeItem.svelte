<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { ListDirectory } from '../../../wailsjs/go/main/App.js';
  import FileTree from './FileTree.svelte';
  import type { DirEntry } from '../types';

  export let entry: DirEntry;
  export let depth: number = 0;
  export let activePath: string | null = null;

  const dispatch = createEventDispatcher<{ fileClick: string }>();

  let expanded = false;
  let children: DirEntry[] = [];
  let loaded = false;

  async function toggle() {
    if (!entry.isDir) {
      dispatch('fileClick', entry.path);
      return;
    }
    expanded = !expanded;
    if (expanded && !loaded) {
      try {
        children = await ListDirectory(entry.path);
        loaded = true;
      } catch (e) {
        console.error('Failed to list directory:', e);
      }
    }
  }

  function handleChildFileClick(event: CustomEvent<string>) {
    dispatch('fileClick', event.detail);
  }

  function getIcon(e: DirEntry): string {
    if (e.isDir) return expanded ? '📂' : '📁';
    const ext = e.ext.toLowerCase();
    switch (ext) {
      case 'json': return '{ }';
      case 'env': return '🔑';
      case 'conf':
      case 'cfg':
      case 'ini': return '⚙';
      case 'yaml':
      case 'yml': return '📋';
      case 'toml': return '📋';
      case 'txt': return '📄';
      case 'md': return '📝';
      default: return '📄';
    }
  }
</script>

<li>
  <div
    class="tree-item"
    class:active={activePath === entry.path}
    class:directory={entry.isDir}
    style="padding-left: {12 + depth * 16}px"
    on:click={toggle}
    on:keydown={(e) => e.key === 'Enter' && toggle()}
    tabindex="0"
    role="treeitem"
    aria-selected={activePath === entry.path}
  >
    {#if entry.isDir}
      <span class="chevron" class:expanded>{expanded ? '▾' : '▸'}</span>
    {:else}
      <span class="chevron-spacer"></span>
    {/if}
    <span class="icon">{getIcon(entry)}</span>
    <span class="name">{entry.name}</span>
  </div>
  {#if entry.isDir && expanded && children.length > 0}
    <FileTree entries={children} depth={depth + 1} {activePath} on:fileClick={handleChildFileClick} />
  {/if}
</li>

<style>
  li {
    list-style: none;
  }

  .tree-item {
    display: flex;
    align-items: center;
    padding: 3px 8px;
    cursor: pointer;
    white-space: nowrap;
    user-select: none;
    font-size: 13px;
    color: #cccccc;
    border-radius: 0;
  }

  .tree-item:hover {
    background: #2a2d2e;
  }

  .tree-item.active {
    background: #37373d;
    color: #ffffff;
  }

  .chevron {
    width: 16px;
    font-size: 10px;
    flex-shrink: 0;
    text-align: center;
  }

  .chevron-spacer {
    width: 16px;
    flex-shrink: 0;
  }

  .icon {
    margin-right: 6px;
    font-size: 12px;
    flex-shrink: 0;
  }

  .name {
    overflow: hidden;
    text-overflow: ellipsis;
  }
</style>
