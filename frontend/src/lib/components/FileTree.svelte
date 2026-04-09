<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import FileTreeItem from './FileTreeItem.svelte';
  import type { DirEntry } from '../types';

  export let entries: DirEntry[] = [];
  export let depth: number = 0;
  export let activePath: string | null = null;

  const dispatch = createEventDispatcher<{ fileClick: string }>();

  function handleFileClick(event: CustomEvent<string>) {
    dispatch('fileClick', event.detail);
  }
</script>

<ul class="tree">
  {#each entries as entry (entry.path)}
    <FileTreeItem {entry} {depth} {activePath} on:fileClick={handleFileClick} />
  {/each}
</ul>

<style>
  .tree {
    list-style: none;
    padding: 0;
    margin: 0;
  }
</style>
