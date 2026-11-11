/* ============================================================
   Todo List App — app.js
   Pure client-side, no build step, no external dependencies.
   All state lives in the `todos` array for the browser session.
   ============================================================ */

// ---------------------------------------------------------------------------
// In-memory store
// ---------------------------------------------------------------------------

/** @type {Array<{id: string, title: string, description: string, completed: boolean, createdAt: number}>} */
let todos = [];

// ---------------------------------------------------------------------------
// TodoStore
// ---------------------------------------------------------------------------

/**
 * Returns a shallow copy of all todos, ordered ascending by createdAt.
 * @returns {Array<{id: string, title: string, description: string, completed: boolean, createdAt: number}>}
 */
function getTodos() {
  return [...todos].sort((a, b) => a.createdAt - b.createdAt);
}

/**
 * Generates a unique ID using crypto.randomUUID() when available,
 * falling back to a timestamp+random string.
 * @returns {string}
 */
function generateId() {
  if (typeof crypto !== 'undefined' && typeof crypto.randomUUID === 'function') {
    return crypto.randomUUID();
  }
  return Date.now().toString(36) + Math.random().toString(36).slice(2);
}

/**
 * Creates a new Todo, appends it to todos[], and returns the created Todo.
 * Callers should validate the title before calling this function.
 * @param {string} title
 * @param {string} [description='']
 * @returns {{id: string, title: string, description: string, completed: boolean, createdAt: number}}
 */
function addTodo(title, description = '') {
  const todo = {
    id: generateId(),
    title: title.trim(),
    description: description,
    completed: false,
    createdAt: Date.now(),
  };
  todos.push(todo);
  return todo;
}

/**
 * Updates the todo with the given id by merging `fields` into it.
 * Only the fields present in `fields` are changed; all other fields are preserved.
 * @param {string} id
 * @param {{ title?: string, description?: string, completed?: boolean }} fields
 * @returns {{id: string, title: string, description: string, completed: boolean, createdAt: number} | null}
 */
function updateTodo(id, fields) {
  const index = todos.findIndex((t) => t.id === id);
  if (index === -1) return null;
  todos[index] = { ...todos[index], ...fields };
  return todos[index];
}

/**
 * Removes the todo with the given id from todos[].
 * @param {string} id
 * @returns {boolean} true if removed, false if id not found
 */
function deleteTodo(id) {
  const index = todos.findIndex((t) => t.id === id);
  if (index === -1) return false;
  todos.splice(index, 1);
  return true;
}

// ---------------------------------------------------------------------------
// Validator
// ---------------------------------------------------------------------------

/**
 * Validates a todo title string.
 * Returns an array of human-readable error message strings.
 * An empty array means the title is valid.
 *
 * Rules:
 *   - After trimming whitespace, length must be 1–255 characters.
 *
 * @param {string} title
 * @returns {string[]}
 */
function validateTitle(title) {
  const trimmed = title.trim();
  if (trimmed.length === 0) {
    return ['Title is required.'];
  }
  if (trimmed.length > 255) {
    return ['Title must be 255 characters or fewer.'];
  }
  return [];
}

// ---------------------------------------------------------------------------
// Renderer
// ---------------------------------------------------------------------------

/**
 * Clears and re-renders the entire todo list from the current store state.
 * Shows an empty-state message when the store is empty.
 */
function renderTodoList() {
  const listEl = document.getElementById('todo-list');
  if (!listEl) return;

  // Clear existing content
  listEl.innerHTML = '';

  const items = getTodos();

  if (items.length === 0) {
    const empty = document.createElement('p');
    empty.id = 'empty-state';
    empty.textContent = 'No tasks yet. Add one above!';
    listEl.appendChild(empty);
    return;
  }

  items.forEach((todo) => {
    listEl.appendChild(renderTodoItem(todo));
  });
}

/**
 * Displays an inline error message adjacent to `inputEl`.
 * If an error element for this input already exists it is updated in place;
 * otherwise a new <span> is inserted immediately after `inputEl`.
 * The error element carries `data-error-for` set to the input's id (or a
 * generated key) so that clearError can locate it without ambiguity.
 *
 * Requirements: 5.4, 5.5
 *
 * @param {HTMLElement} inputEl - The input or textarea the error belongs to.
 * @param {string} message - Human-readable error text to display.
 */
function showError(inputEl, message) {
  // Use the element's id as a stable key; fall back to a data attribute we set.
  let key = inputEl.id;
  if (!key) {
    key = inputEl.dataset.errorKey;
    if (!key) {
      key = 'err-' + Date.now().toString(36);
      inputEl.dataset.errorKey = key;
    }
  }

  // Look for an existing sibling error element tied to this input.
  let errorEl = inputEl.parentNode
    ? inputEl.parentNode.querySelector(`[data-error-for="${key}"]`)
    : null;

  if (!errorEl) {
    errorEl = document.createElement('span');
    errorEl.dataset.errorFor = key;
    errorEl.className = 'field-error';
    errorEl.setAttribute('role', 'alert');
    // Insert immediately after inputEl.
    inputEl.insertAdjacentElement('afterend', errorEl);
  }

  errorEl.textContent = message;
}

/**
 * Removes any inline error message previously created by `showError` for
 * `inputEl`. If no error element exists this is a no-op.
 *
 * Requirements: 5.5
 *
 * @param {HTMLElement} inputEl - The input or textarea whose error should be removed.
 */
function clearError(inputEl) {
  const key = inputEl.id || inputEl.dataset.errorKey;
  if (!key) return;

  const errorEl = inputEl.parentNode
    ? inputEl.parentNode.querySelector(`[data-error-for="${key}"]`)
    : null;

  if (errorEl) {
    errorEl.remove();
  }
}

/**
 * Creates and returns a <li> DOM element representing a single todo item.
 * Includes the title, description, completion indicator, and placeholder
 * controls for toggle, edit, and delete (wired up in later tasks).
 * @param {{id: string, title: string, description: string, completed: boolean, createdAt: number}} todo
 * @returns {HTMLLIElement}
 */
function renderTodoItem(todo) {
  const li = document.createElement('li');
  li.dataset.id = todo.id;
  li.classList.toggle('completed', todo.completed);

  // --- Content area ---
  const content = document.createElement('div');
  content.className = 'todo-content';

  const titleEl = document.createElement('span');
  titleEl.className = 'todo-title';
  titleEl.textContent = todo.title;

  content.appendChild(titleEl);

  if (todo.description) {
    const descEl = document.createElement('p');
    descEl.className = 'todo-description';
    descEl.textContent = todo.description;
    content.appendChild(descEl);
  }

  // --- Controls area ---
  const controls = document.createElement('div');
  controls.className = 'todo-controls';

  // Toggle button — wired in task 6.1
  const toggleBtn = document.createElement('button');
  toggleBtn.type = 'button';
  toggleBtn.className = 'btn-toggle';
  toggleBtn.setAttribute('aria-label', todo.completed ? 'Mark incomplete' : 'Mark complete');
  toggleBtn.textContent = todo.completed ? '↩ Undo' : '✓ Done';
  toggleBtn.addEventListener('click', () => handleToggle(todo.id));

  // Edit button — switches the item into edit mode
  const editBtn = document.createElement('button');
  editBtn.type = 'button';
  editBtn.className = 'btn-edit';
  editBtn.setAttribute('aria-label', 'Edit todo');
  editBtn.textContent = '✎ Edit';
  editBtn.addEventListener('click', () => {
    // Replace view mode with inline edit mode on this <li>.
    // Keep a reference to the li so we can swap content in place.
    const liEl = li;

    // --- Build edit form ---
    const editForm = document.createElement('div');
    editForm.className = 'todo-edit-form';

    // Title input
    const titleInput = document.createElement('input');
    titleInput.type = 'text';
    titleInput.id = `edit-title-${todo.id}`;
    titleInput.className = 'edit-title-input';
    titleInput.value = todo.title;
    titleInput.setAttribute('aria-label', 'Edit title');

    // Description textarea
    const descInput = document.createElement('textarea');
    descInput.id = `edit-desc-${todo.id}`;
    descInput.className = 'edit-desc-input';
    descInput.value = todo.description;
    descInput.setAttribute('aria-label', 'Edit description');
    descInput.rows = 2;

    // Save button
    const saveBtn = document.createElement('button');
    saveBtn.type = 'button';
    saveBtn.className = 'btn-save';
    saveBtn.textContent = '✔ Save';
    saveBtn.addEventListener('click', () => handleEditSave(todo.id));

    // Cancel button
    const cancelBtn = document.createElement('button');
    cancelBtn.type = 'button';
    cancelBtn.className = 'btn-cancel';
    cancelBtn.textContent = '✕ Cancel';
    cancelBtn.addEventListener('click', () => handleEditCancel(todo.id));

    // Allow saving with Enter key on the title input
    titleInput.addEventListener('keydown', (e) => {
      if (e.key === 'Enter') {
        e.preventDefault();
        handleEditSave(todo.id);
      } else if (e.key === 'Escape') {
        handleEditCancel(todo.id);
      }
    });

    // Clear title error when user types
    titleInput.addEventListener('input', () => clearError(titleInput));

    editForm.appendChild(titleInput);
    editForm.appendChild(descInput);

    const editControls = document.createElement('div');
    editControls.className = 'edit-controls';
    editControls.appendChild(saveBtn);
    editControls.appendChild(cancelBtn);
    editForm.appendChild(editControls);

    // Swap: remove existing children and insert the edit form
    while (liEl.firstChild) {
      liEl.removeChild(liEl.firstChild);
    }
    liEl.classList.add('editing');
    liEl.appendChild(editForm);

    // Focus the title input so the user can start typing immediately
    titleInput.focus();
  });

  // Delete button — wired in task 6.1
  const deleteBtn = document.createElement('button');
  deleteBtn.type = 'button';
  deleteBtn.className = 'btn-delete';
  deleteBtn.setAttribute('aria-label', 'Delete todo');
  deleteBtn.textContent = '✕ Delete';
  deleteBtn.addEventListener('click', () => handleDelete(todo.id));

  controls.appendChild(toggleBtn);
  controls.appendChild(editBtn);
  controls.appendChild(deleteBtn);

  li.appendChild(content);
  li.appendChild(controls);

  return li;
}

// ---------------------------------------------------------------------------
// EventHandlers
// ---------------------------------------------------------------------------

/**
 * Handles toggling the completed status of a todo.
 * Flips the current completed boolean and re-renders the list.
 *
 * Requirements: 3.6, 3.7
 *
 * @param {string} todoId
 */
function handleToggle(todoId) {
  const current = getTodos().find((t) => t.id === todoId);
  if (!current) return;
  updateTodo(todoId, { completed: !current.completed });
  renderTodoList();
}

/**
 * Handles clicking the delete button for a todo.
 * Removes the todo from the store and re-renders the list.
 * When the last todo is deleted the empty-state message is shown automatically
 * by renderTodoList().
 *
 * Requirements: 4.1, 4.2, 4.3
 *
 * @param {string} todoId
 */
function handleDelete(todoId) {
  deleteTodo(todoId);
  renderTodoList();
}

/**
 * Handles saving an inline edit.
 * Reads the title/description inputs rendered by the edit mode in renderTodoItem,
 * validates the title, shows/clears inline errors, and on success updates the
 * store and re-renders.
 *
 * Requirements: 3.1, 3.2, 3.3, 3.4, 3.5
 *
 * @param {string} todoId
 */
function handleEditSave(todoId) {
  const titleInput = document.getElementById(`edit-title-${todoId}`);
  const descInput = document.getElementById(`edit-desc-${todoId}`);

  if (!titleInput) return; // edit mode not active — nothing to do

  const title = titleInput.value;
  const description = descInput ? descInput.value : '';

  const errors = validateTitle(title);

  if (errors.length > 0) {
    showError(titleInput, errors[0]);
    return;
  }

  clearError(titleInput);

  updateTodo(todoId, { title: title.trim(), description });
  renderTodoList();
}

/**
 * Handles cancelling an inline edit without saving.
 * Simply re-renders the list to restore the view-mode display.
 *
 * Requirements: 3.1
 *
 * @param {string} todoId
 */
function handleEditCancel(todoId) { // eslint-disable-line no-unused-vars
  renderTodoList();
}

/**
 * Handles create-form submission.
 * - Prevents the default page reload.
 * - Reads title and description from the form inputs.
 * - Validates the title; shows/clears inline errors.
 * - On success: adds the todo, clears the form, and re-renders the list.
 *
 * Requirements: 2.1, 2.2, 2.3, 2.4, 2.5, 2.6
 *
 * @param {SubmitEvent} event
 */
function handleCreateSubmit(event) {
  event.preventDefault();

  const titleInput = document.getElementById('create-title');
  const descInput = document.getElementById('create-description');

  const title = titleInput.value;
  const description = descInput ? descInput.value : '';

  const errors = validateTitle(title);

  if (errors.length > 0) {
    // Show the first error message adjacent to the title input.
    showError(titleInput, errors[0]);
    return;
  }

  // Input is valid — clear any previous error.
  clearError(titleInput);

  addTodo(title.trim(), description);

  // Clear the form inputs (Requirement 2.6).
  titleInput.value = '';
  if (descInput) descInput.value = '';

  renderTodoList();
}

// ---------------------------------------------------------------------------
// Bootstrap
// ---------------------------------------------------------------------------

document.addEventListener('DOMContentLoaded', () => {
  renderTodoList();

  // Wire up create form (task 5.1).
  const createForm = document.getElementById('create-form');
  if (createForm) {
    createForm.addEventListener('submit', handleCreateSubmit);
  }

  // Clear the title error as soon as the user starts correcting it (Requirement 5.5).
  const titleInput = document.getElementById('create-title');
  if (titleInput) {
    titleInput.addEventListener('input', () => clearError(titleInput));
  }
});
