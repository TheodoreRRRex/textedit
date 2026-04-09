<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { EditorView, basicSetup } from 'codemirror';
  import { highlightWhitespace } from '@codemirror/view';
  import { EditorState, Compartment } from '@codemirror/state';
  import { json } from '@codemirror/lang-json';
  import { oneDark } from '@codemirror/theme-one-dark';
  import {
    currentContent,
    savedContent,
    isDirty,
    cursorPosition,
    showWhitespace,
  } from '../stores/editor';

  const whitespaceCompartment = new Compartment();

  let container: HTMLDivElement;
  let view: EditorView;

  export function getContent(): string {
    if (view) {
      return view.state.doc.toString();
    }
    return $currentContent;
  }

  export function loadFile(content: string, filePath: string | null) {
    if (view) {
      view.setState(createState(content, filePath));
    }
  }

  function getLanguageExtension(path: string | null) {
    if (!path) return [];
    const ext = path.split('.').pop()?.toLowerCase();
    switch (ext) {
      case 'json':
        return [json()];
      default:
        return [];
    }
  }

  function createState(content: string, filePath: string | null): EditorState {
    return EditorState.create({
      doc: content,
      extensions: [
        basicSetup,
        oneDark,
        EditorView.lineWrapping,
        whitespaceCompartment.of($showWhitespace ? highlightWhitespace() : []),
        ...getLanguageExtension(filePath),
        EditorView.updateListener.of((update) => {
          if (update.docChanged) {
            const newContent = update.state.doc.toString();
            $currentContent = newContent;
            $isDirty = newContent !== $savedContent;
          }
          if (update.selectionSet) {
            const pos = update.state.selection.main.head;
            const line = update.state.doc.lineAt(pos);
            $cursorPosition = {
              line: line.number,
              col: pos - line.from + 1,
            };
          }
        }),
      ],
    });
  }

  onMount(() => {
    view = new EditorView({
      state: createState('', null),
      parent: container,
    });

    const handleKeydown = (e: KeyboardEvent) => {
      if (e.ctrlKey && e.key === 'h') {
        e.preventDefault();
        $showWhitespace = !$showWhitespace;
      }
    };
    window.addEventListener('keydown', handleKeydown);

    return () => {
      window.removeEventListener('keydown', handleKeydown);
    };
  });

  $: if (view) {
    view.dispatch({
      effects: whitespaceCompartment.reconfigure(
        $showWhitespace ? highlightWhitespace() : []
      ),
    });
  }

  onDestroy(() => {
    if (view) view.destroy();
  });
</script>

<div class="editor-container" bind:this={container}></div>

<style>
  .editor-container {
    overflow: hidden;
    background: #1e1e1e;
  }

  .editor-container :global(.cm-editor) {
    height: 100%;
  }

  .editor-container :global(.cm-scroller) {
    font-family: "Cascadia Code", "Consolas", "Courier New", monospace;
    font-size: 14px;
  }
</style>
