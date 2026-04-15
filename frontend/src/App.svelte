<script lang="ts">
  import Sidebar from './lib/components/Sidebar.svelte';
  import Editor from './lib/components/Editor.svelte';
  import StatusBar from './lib/components/StatusBar.svelte';
  import { onMount } from 'svelte';
  import { EventsOn } from '../wailsjs/runtime/runtime.js';
  import { OpenFile, SaveFile, SaveFileDialog, PickFileDialog, SetDirty, SetWindowTitle, GetStartupFile } from '../wailsjs/go/main/App.js';
  import {
    currentFilePath,
    currentFileName,
    isDirty,
    currentContent,
    savedContent,
    lineEnding,
  } from './lib/stores/editor';

  let editorComponent: Editor;

  function updateTitle() {
    const dirty = $isDirty ? '● ' : '';
    const name = $currentFileName;
    SetWindowTitle(`${dirty}${name} — TextEdit`);
  }

  $: $isDirty, $currentFileName, updateTitle();
  $: SetDirty($isDirty);

  async function openFileByPath(path: string) {
    try {
      const result = await OpenFile(path);
      $currentFilePath = result.path;
      $currentFileName = result.name;
      $savedContent = result.content;
      $currentContent = result.content;
      $isDirty = false;
      $lineEnding = result.hasCRLF ? 'CRLF' : 'LF';
      editorComponent?.loadFile(result.content, result.path);
    } catch (e) {
      alert(`Failed to open file: ${e}`);
    }
  }

  async function handleOpenFile() {
    try {
      const path = await PickFileDialog();
      if (path) {
        await openFileByPath(path);
      }
    } catch (e) {
      alert(`Failed to open file: ${e}`);
    }
  }

  async function handleSave() {
    if (!$currentFilePath) {
      handleSaveAs();
      return;
    }
    try {
      const content = editorComponent?.getContent() ?? $currentContent;
      await SaveFile($currentFilePath, content);
      $savedContent = content;
      $isDirty = false;
    } catch (e) {
      alert(`Failed to save file: ${e}`);
    }
  }

  async function handleSaveAs() {
    try {
      const defaultName = $currentFileName !== 'Untitled' ? $currentFileName : '';
      const path = await SaveFileDialog(defaultName);
      if (path) {
        const content = editorComponent?.getContent() ?? $currentContent;
        await SaveFile(path, content);
        $currentFilePath = path;
        $currentFileName = path.split('/').pop() || path.split('\\').pop() || 'Untitled';
        $savedContent = content;
        $isDirty = false;
      }
    } catch (e) {
      alert(`Failed to save file: ${e}`);
    }
  }

  function handleNewFile() {
    $currentFilePath = null;
    $currentFileName = 'Untitled';
    $currentContent = '';
    $savedContent = '';
    $isDirty = false;
    editorComponent?.loadFile('', null);
  }

  function handleFileSelect(event: CustomEvent<string>) {
    openFileByPath(event.detail);
  }

  // On startup, check if a file was passed via CLI (e.g. "Open with")
  onMount(async () => {
    const startupFile = await GetStartupFile();
    if (startupFile) {
      await openFileByPath(startupFile);
    }
  });

  // Listen for menu events from Go
  EventsOn('menu:open-file', handleOpenFile);
  EventsOn('menu:save', handleSave);
  EventsOn('menu:save-as', handleSaveAs);
  EventsOn('menu:new-file', handleNewFile);
</script>

<div id="layout">
  <Sidebar on:fileSelect={handleFileSelect} />
  <Editor bind:this={editorComponent} />
  <StatusBar />
</div>

<style>
  #layout {
    display: grid;
    grid-template-columns: 250px 1fr;
    grid-template-rows: 1fr 28px;
    height: 100vh;
    overflow: hidden;
  }
</style>
