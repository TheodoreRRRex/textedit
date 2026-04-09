<script lang="ts">
  import {
    currentFilePath,
    currentFileName,
    isDirty,
    cursorPosition,
    lineEnding,
    showWhitespace,
  } from '../stores/editor';
</script>

<footer class="status-bar">
  <div class="left">
    {#if $isDirty}
      <span class="dirty-indicator">●</span>
    {/if}
    <span class="file-path" title={$currentFilePath || 'Untitled'}>
      {$currentFilePath || 'Untitled'}
    </span>
  </div>
  <div class="right">
    {#if $showWhitespace}
      <span class="ws-indicator">[WS]</span>
      <span class="separator">|</span>
    {/if}
    <span>Ln {$cursorPosition.line}, Col {$cursorPosition.col}</span>
    <span class="separator">|</span>
    <span>{$lineEnding}</span>
    <span class="separator">|</span>
    <span>UTF-8</span>
  </div>
</footer>

<style>
  .status-bar {
    grid-column: 1 / -1;
    background: #007acc;
    color: #ffffff;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 12px;
    font-size: 12px;
    white-space: nowrap;
    user-select: none;
  }

  .left {
    display: flex;
    align-items: center;
    gap: 6px;
    overflow: hidden;
    min-width: 0;
  }

  .file-path {
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .dirty-indicator {
    color: #ffffff;
    font-size: 14px;
  }

  .right {
    display: flex;
    align-items: center;
    gap: 8px;
    flex-shrink: 0;
  }

  .separator {
    opacity: 0.6;
  }

  .ws-indicator {
    font-weight: bold;
  }
</style>
