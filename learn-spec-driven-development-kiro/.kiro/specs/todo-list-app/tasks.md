# Implementation Plan: Todo List App

## Overview

Implement a pure client-side todo list application as exactly three files: `index.html`, `style.css`, and `app.js`. No server, no build step, no external runtime dependencies. All state lives in an in-memory JavaScript array for the duration of the browser session. Every CRUD operation mutates the array and immediately re-renders the DOM via a single `renderTodoList()` entry point.

## Tasks

- [x] 1. Set up project files
  - [x] 1.1 Create `index.html` with semantic HTML boilerplate
    - Add `<!DOCTYPE html>`, `<html lang="en">`, `<head>`, and `<body>` structure
    - Include a visible `<h1>` heading that identifies the application
    - Add a `<link>` tag referencing `style.css` and a `<script defer src="app.js">` tag
    - Add placeholder `<main>` with an empty todo list container (`<ul id="todo-list">`) and a create-form placeholder
    - _Requirements: 5.1, 5.2_

  - [x] 1.2 Create `style.css` and `app.js` with empty boilerplate
    - `style.css`: add a CSS reset and basic body/font variables
    - `app.js`: add a `DOMContentLoaded` listener stub and a module-level `let todos = [];` declaration
    - _Requirements: 5.1_

- [x] 2. Implement TodoStore
  - [x] 2.1 Implement `getTodos`, `addTodo`, `updateTodo`, and `deleteTodo` in `app.js`
    - `getTodos()` — returns a shallow copy of `todos`, ordered ascending by `createdAt`
    - `addTodo(title, description = '')` — creates a Todo object (id via `crypto.randomUUID()` with `Date.now().toString(36) + Math.random().toString(36).slice(2)` fallback, `completed: false`, `createdAt: Date.now()`), appends to `todos[]`, returns the new todo
    - `updateTodo(id, fields)` — merges `fields` into the matching todo; returns updated todo or `null` if not found
    - `deleteTodo(id)` — removes todo with matching id; returns `true` if removed, `false` if not found
    - _Requirements: 2.2, 3.2, 3.7, 4.2_

  - [x]* 2.2 Write property test for TodoStore — Property 1 (Create round-trip)
    - **Property 1: Create round-trip** — `addTodo(title, desc)` then `getTodos()` contains exactly one todo matching title, desc, `completed: false`, and positive `createdAt`
    - **Validates: Requirements 2.2, 2.3**
    - Use `fc.string({ minLength: 1, maxLength: 255 })` for title; `fc.string()` for description
    - Tag: `// Feature: todo-list-app, Property 1: Create round-trip`

  - [x]* 2.3 Write property test for TodoStore — Property 2 (Unique IDs)
    - **Property 2: All todo IDs are unique** — N calls to `addTodo` produce N pairwise-distinct IDs
    - **Validates: Requirements 2.2**
    - Use `fc.integer({ min: 2, max: 50 })` for N
    - Tag: `// Feature: todo-list-app, Property 2: Unique IDs`

  - [x]* 2.4 Write property test for TodoStore — Property 3 (Delete removes todo)
    - **Property 3: Delete permanently removes the todo** — after `deleteTodo(id)`, `getTodos()` contains no todo with that id
    - **Validates: Requirements 4.2**
    - Use an arbitrary store of 1–20 todos; randomly pick one id to delete
    - Tag: `// Feature: todo-list-app, Property 3: Delete removes todo`

  - [x]* 2.5 Write property test for TodoStore — Property 4 (Toggle preserves fields)
    - **Property 4: Toggle flips completed and preserves all other fields** — `updateTodo(id, { completed: !todo.completed })` flips only `completed`; `id`, `title`, `description`, `createdAt` remain identical
    - **Validates: Requirements 3.6, 3.7**
    - Tag: `// Feature: todo-list-app, Property 4: Toggle preserves fields`

- [x] 3. Implement Validator
  - [x] 3.1 Implement `validateTitle(title)` in `app.js`
    - Trim whitespace before measuring length
    - Return `[]` (empty array) when trimmed length is 1–255
    - Return a non-empty array with a human-readable error message when trimmed length is 0 (title required) or > 255 (title too long)
    - _Requirements: 2.4, 2.5, 3.4, 3.5_

  - [x]* 3.2 Write property test for Validator — Property 5 (Validation rejects bad titles)
    - **Property 5: Validation rejects every non-conforming title** — empty strings, whitespace-only strings, and strings with trimmed length > 255 always produce a non-empty errors array
    - **Validates: Requirements 2.4, 2.5, 3.4, 3.5**
    - Use `fc.oneof(fc.constant(""), fc.string().map(s => " ".repeat(s.length)), fc.string({ minLength: 256 }))`
    - Tag: `// Feature: todo-list-app, Property 5: Validation rejects bad titles`

- [x] 4. Implement Renderer
  - [x] 4.1 Implement `renderTodoList()` and `renderTodoItem(todo)` in `app.js`
    - `renderTodoList()` — clears `#todo-list`, iterates `getTodos()`, appends `renderTodoItem(todo)` for each; shows an empty-state `<p>` when store is empty
    - `renderTodoItem(todo)` — returns a `<li>` element containing the todo's title, description, completion status indicator, and placeholder controls for toggle, edit, and delete (can be stubs at this stage)
    - Apply a CSS class (e.g. `completed`) to finished todos so styles can target them
    - _Requirements: 1.1, 1.2, 1.3, 1.4, 5.3_

  - [x] 4.2 Implement `showError(inputEl, message)` and `clearError(inputEl)` in `app.js`
    - `showError` — inserts or updates a sibling error element adjacent to `inputEl` with the given message
    - `clearError` — removes the sibling error element if present
    - _Requirements: 5.4, 5.5_

- [x] 5. Implement Create Todo
  - [x] 5.1 Add create form to `index.html` and wire `handleCreateSubmit` in `app.js`
    - Add `<form id="create-form">` with a title `<input>`, description `<textarea>`, and submit `<button>` to `index.html`
    - Implement `handleCreateSubmit(event)`: call `event.preventDefault()`, read inputs, call `validateTitle`, show/clear errors, call `addTodo` on success, clear form inputs, call `renderTodoList()`
    - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5, 2.6_

- [x] 6. Implement Toggle and Delete
  - [x] 6.1 Add toggle and delete controls to `renderTodoItem` and implement `handleToggle` / `handleDelete`
    - `renderTodoItem` — add a checkbox or toggle button (checked state reflects `todo.completed`) and a delete button to each `<li>`; attach click listeners that call `handleToggle(todo.id)` / `handleDelete(todo.id)`
    - `handleToggle(todoId)` — calls `updateTodo(id, { completed: !current.completed })`, then `renderTodoList()`
    - `handleDelete(todoId)` — calls `deleteTodo(id)`, then `renderTodoList()`
    - _Requirements: 3.6, 3.7, 4.1, 4.2, 4.3_

- [x] 7. Implement Inline Edit
  - [x] 7.1 Add inline edit controls to `renderTodoItem` and implement `handleEditSave` / `handleEditCancel`
    - `renderTodoItem` — add an edit button that, when clicked, replaces the title/description display with editable `<input>` and `<textarea>` pre-filled with current values, plus Save and Cancel buttons
    - `handleEditSave(todoId)` — reads inputs, calls `validateTitle`, shows/clears errors, calls `updateTodo` on success, then `renderTodoList()`
    - `handleEditCancel(todoId)` — calls `renderTodoList()` to restore the original view without saving
    - _Requirements: 3.1, 3.2, 3.3, 3.4, 3.5_

- [x] 8. Checkpoint — Core functionality complete
  - Verify all CRUD operations work in the browser: create, read, toggle, delete, inline edit, inline edit cancel, and validation errors all function without page reloads
  - Ensure all tests pass, ask the user if questions arise.

- [x] 9. Style the app (`style.css`)
  - [x] 9.1 Add layout and component styles
    - Center the page content; style the `<h1>` heading and overall page layout
    - Style the create form: input, textarea, and submit button with clear focus states
    - Style the todo list and individual `<li>` items with readable spacing
    - _Requirements: 5.2, 5.6, 5.7_

  - [x] 9.2 Add completion and error styles
    - `.completed` class — apply `text-decoration: line-through` and a muted color to completed todo titles
    - Inline error message element — style with a visible error color (e.g. red) in small text adjacent to the input
    - Visually distinguish complete vs incomplete todos so the status is immediately apparent
    - _Requirements: 1.4, 5.4, 5.5, 5.6_

- [x] 10. Final checkpoint — Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 11. (Optional) Write unit tests (`test.js`)
  - [ ]* 11.1 Write unit tests for TodoStore
    - Test `addTodo` with valid title: store length increases by 1, returned todo has correct fields, `completed` is `false`
    - Test N `addTodo` calls: all N resulting IDs are distinct
    - Test `getTodos` returns a copy: mutating the returned array does not affect the internal store
    - Test `getTodos` ordering: returns todos in ascending `createdAt` order
    - Test `updateTodo` partial update: only supplied fields change; others preserved
    - Test `updateTodo` unknown id: returns `null`; store unchanged
    - Test `deleteTodo` owned id: store no longer contains that id; returns `true`
    - Test `deleteTodo` unknown id: returns `false`; store unchanged
    - _Requirements: 2.2, 3.2, 3.7, 4.2_

  - [ ]* 11.2 Write unit tests for Validator
    - Empty string `""` → errors non-empty
    - Whitespace only `"   "` → errors non-empty
    - Single character `"a"` → errors empty
    - Exactly 255 characters → errors empty
    - 256 characters → errors non-empty
    - Untrimmed valid `"  hello  "` → errors empty (trimmed to 5 chars)
    - _Requirements: 2.4, 2.5, 3.4, 3.5_

- [ ] 12. (Optional) Write property-based tests (`test.pbt.js`)
  - [ ]* 12.1 Write property-based tests for Properties 1–6 using fast-check
    - **Property 1: Create round-trip** — `fc.string({ minLength: 1, maxLength: 255 })` for title; `fc.string()` for description; **Validates: Requirements 2.2, 2.3**
    - **Property 2: Unique IDs** — `fc.integer({ min: 2, max: 50 })` for N; **Validates: Requirements 2.2**
    - **Property 3: Delete removes todo** — arbitrary store of 1–20 todos, random pick and delete; **Validates: Requirements 4.2**
    - **Property 4: Toggle preserves fields** — arbitrary store, random pick, toggle; **Validates: Requirements 3.6, 3.7**
    - **Property 5: Validation rejects bad titles** — `fc.oneof(fc.constant(""), fc.string().map(s => " ".repeat(s.length)), fc.string({ minLength: 256 }))`; **Validates: Requirements 2.4, 2.5, 3.4, 3.5**
    - **Property 6: UI reflects store state** — arbitrary sequence of add/update/delete; assert DOM item count equals `getTodos().length` and items appear in `createdAt` order after each mutation; **Validates: Requirements 1.1, 1.3, 2.3, 3.3, 4.2, 5.3**
    - Minimum 100 iterations per property; tag each: `// Feature: todo-list-app, Property N: ...`
    - Run with `node test.pbt.js` (no browser required)
    - _Requirements: 2.2, 2.3, 2.4, 2.5, 3.4, 3.5, 3.6, 3.7, 4.2_

## Notes

- Tasks marked with `*` are optional and can be skipped for a faster MVP
- Each task references specific requirements for traceability
- Checkpoints ensure incremental browser validation before moving on
- Property tests (tasks 2.2–2.5, 3.2, 12.1) validate universal correctness properties; unit tests (tasks 11.1–11.2) validate specific examples and edge cases
- Run `node test.js` and `node test.pbt.js` from the project root; no build step or server required

## Task Dependency Graph

```json
{
  "waves": [
    { "id": 0, "tasks": ["1.1", "1.2"] },
    { "id": 1, "tasks": ["2.1"] },
    { "id": 2, "tasks": ["3.1", "2.2", "2.3", "2.4", "2.5"] },
    { "id": 3, "tasks": ["4.1", "3.2"] },
    { "id": 4, "tasks": ["4.2"] },
    { "id": 5, "tasks": ["5.1"] },
    { "id": 6, "tasks": ["6.1"] },
    { "id": 7, "tasks": ["7.1"] },
    { "id": 8, "tasks": ["9.1"] },
    { "id": 9, "tasks": ["9.2", "11.1", "11.2"] },
    { "id": 10, "tasks": ["12.1"] }
  ]
}
```
