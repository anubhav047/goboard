import './styles/main.css';
import { authAPI, boardsAPI, listsAPI, cardsAPI, User, Board, List, Card } from './api/client';
import logoUrl from './assets/logo.png';

// Global state
let currentUser: User | null = null;

// DOM elements
const navBrand = document.getElementById('nav-brand')!;
const navLinks = document.getElementById('nav-links')!;
const mainContent = document.getElementById('main-content')!;

// Set up navbar logo
function setupNavbar() {
  navBrand.innerHTML = `
    <img src="${logoUrl}" alt="GoBoard Logo" />
  `;
}

// Initialize app
async function init() {
  // Set up the navbar logo
  setupNavbar();
  
  try {
    // Check if user is already logged in
    currentUser = await authAPI.me();
    showDashboard();
  } catch (error) {
    // User not logged in, show login form
    showLogin();
  }
}

// Show login form
function showLogin() {
  navLinks.innerHTML = `
    <button class="btn" onclick="showRegister()">Sign Up</button>
  `;

  mainContent.innerHTML = `
    <div class="container">
      <div class="card" style="max-width: 400px; margin: 4rem auto;">
        <h1 style="text-align: center; margin-bottom: 2rem; color: #ffffff; font-weight: 600; font-size: 1.5rem;">Welcome Back</h1>
        
        <div id="login-error" class="error hidden"></div>

        <form id="login-form">
          <div class="form-group">
            <label for="email">Email</label>
            <input type="email" id="email" name="email" required>
          </div>

          <div class="form-group">
            <label for="password">Password</label>
            <input type="password" id="password" name="password" required>
          </div>

          <button type="submit" class="btn" style="width: 100%;">Sign In</button>
        </form>

        <p style="text-align: center; margin-top: 2rem; color: #94a3b8; font-size: 0.9rem;">
          Don't have an account? <a href="#" onclick="showRegister()" style="color: #ffffff; text-decoration: underline;">Sign up here</a>
        </p>
      </div>
    </div>
  `;

  // Handle login form submission
  const loginForm = document.getElementById('login-form') as HTMLFormElement;
  loginForm.addEventListener('submit', handleLogin);
}

// Show register form
function showRegister() {
  navLinks.innerHTML = `
    <button class="btn btn-secondary" onclick="showLogin()">Sign In</button>
  `;

  mainContent.innerHTML = `
    <div class="container">
      <div class="card" style="max-width: 400px; margin: 4rem auto;">
        <h1 style="text-align: center; margin-bottom: 2rem; color: #ffffff; font-weight: 600; font-size: 1.5rem;">Create Account</h1>
        
        <div id="register-error" class="error hidden"></div>

        <form id="register-form">
          <div class="form-group">
            <label for="name">Full Name</label>
            <input type="text" id="name" name="name" required>
          </div>

          <div class="form-group">
            <label for="email">Email</label>
            <input type="email" id="email" name="email" required>
          </div>

          <div class="form-group">
            <label for="password">Password</label>
            <input type="password" id="password" name="password" required minlength="6">
          </div>

          <button type="submit" class="btn" style="width: 100%;">Create Account</button>
        </form>

        <p style="text-align: center; margin-top: 2rem; color: #94a3b8; font-size: 0.9rem;">
          Already have an account? <a href="#" onclick="showLogin()" style="color: #ffffff; text-decoration: underline;">Sign in here</a>
        </p>
      </div>
    </div>
  `;

  // Handle register form submission
  const registerForm = document.getElementById('register-form') as HTMLFormElement;
  registerForm.addEventListener('submit', handleRegister);
}

// Show board view with lists and cards
async function showBoard(boardId: number) {
  navLinks.innerHTML = `
    <button class="btn btn-secondary" onclick="showDashboard()">← Back to Dashboard</button>
    <span>Welcome, ${currentUser?.name}</span>
    <button class="btn btn-secondary" onclick="logout()">Logout</button>
  `;

  try {
    console.log('Loading board with ID:', boardId);
    
    // Get board details and lists
    const [board, lists] = await Promise.all([
      boardsAPI.getById(boardId),
      listsAPI.getByBoard(boardId)
    ]);
    
    console.log('Board loaded:', board);
    console.log('Lists loaded:', lists);

    // Get cards for each list
    const listsWithCards = await Promise.all(
      lists.map(async (list) => {
        const cards = await cardsAPI.getByList(list.ID);
        return { ...list, cards: cards || [] };
      })
    );

    mainContent.innerHTML = `
      <div class="container">
        <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 2rem; padding-bottom: 1rem; border-bottom: 1px solid #1a1a1a;">
          <div>
            <h1 style="color: #ffffff; font-size: 2rem; font-weight: 600; margin-bottom: 0.5rem;">${board.Name}</h1>
            <p style="color: #94a3b8; font-size: 0.9rem;">${board.Description}</p>
          </div>
          <button onclick="showCreateListModal()" class="btn">+ Add List</button>
        </div>

        <div class="board-container">
          ${listsWithCards.map(list => `
            <div class="list" data-list-id="${list.ID}">
              <div class="list-header">
                <h3 class="list-title">${list.Name}</h3>
                <button onclick="showCreateCardModal(${list.ID})" class="btn" style="padding: 0.5rem; font-size: 0.8rem;">+ Add Card</button>
              </div>
              <div class="cards-container">
                ${(list.cards || []).map(card => `
                  <div class="card-item" data-card-id="${card.ID}" draggable="true">
                    <div class="card-title">${card.Title}</div>
                    ${card.Description ? `<div class="card-description">${card.Description}</div>` : ''}
                  </div>
                `).join('')}
              </div>
            </div>
          `).join('')}
        </div>
      </div>

      <!-- Create List Modal -->
      <div id="createListModal" class="hidden" style="position: fixed; top: 0; left: 0; width: 100%; height: 100%; background: rgba(0,0,0,0.8); z-index: 1000;">
        <div style="position: absolute; top: 50%; left: 50%; transform: translate(-50%, -50%); background: #000000; border: 1px solid #1a1a1a; padding: 2rem; border-radius: 4px; width: 90%; max-width: 400px;">
          <h2 style="margin-bottom: 2rem; color: #ffffff; font-weight: 600;">Create New List</h2>
          
          <form id="create-list-form">
            <div class="form-group">
              <label for="list-name">List Name</label>
              <input type="text" id="list-name" name="name" required placeholder="e.g., To Do">
            </div>

            <div style="display: flex; gap: 1rem; justify-content: flex-end;">
              <button type="button" onclick="hideCreateListModal()" class="btn btn-secondary">Cancel</button>
              <button type="submit" class="btn">Create List</button>
            </div>
          </form>
        </div>
      </div>

      <!-- Create Card Modal -->
      <div id="createCardModal" class="hidden" style="position: fixed; top: 0; left: 0; width: 100%; height: 100%; background: rgba(0,0,0,0.8); z-index: 1000;">
        <div style="position: absolute; top: 50%; left: 50%; transform: translate(-50%, -50%); background: #000000; border: 1px solid #1a1a1a; padding: 2rem; border-radius: 4px; width: 90%; max-width: 500px;">
          <h2 style="margin-bottom: 2rem; color: #ffffff; font-weight: 600;">Create New Card</h2>
          
          <form id="create-card-form">
            <div class="form-group">
              <label for="card-title">Card Title</label>
              <input type="text" id="card-title" name="title" required placeholder="e.g., Design homepage">
            </div>

            <div class="form-group">
              <label for="card-description">Description (Optional)</label>
              <textarea id="card-description" name="description" rows="3" placeholder="Add more details..."></textarea>
            </div>

            <div style="display: flex; gap: 1rem; justify-content: flex-end;">
              <button type="button" onclick="hideCreateCardModal()" class="btn btn-secondary">Cancel</button>
              <button type="submit" class="btn">Create Card</button>
            </div>
          </form>
        </div>
      </div>
    `;

    // Store current board ID for modal handlers
    (window as any).currentBoardId = boardId;
    (window as any).currentListId = null;

    // Handle create list form
    const createListForm = document.getElementById('create-list-form') as HTMLFormElement;
    createListForm.addEventListener('submit', handleCreateList);

    // Handle create card form  
    const createCardForm = document.getElementById('create-card-form') as HTMLFormElement;
    createCardForm.addEventListener('submit', handleCreateCard);

    // Initialize drag and drop
    initializeDragAndDrop();

  } catch (error: any) {
    console.error('Error loading board:', error);
    console.error('Error details:', {
      message: error.message,
      response: error.response?.data,
      status: error.response?.status
    });
    showError(`Failed to load board: ${error.response?.data?.error || error.message || 'Unknown error'}`);
  }
}

// Show dashboard
async function showDashboard() {
  navLinks.innerHTML = `
    <span>Welcome, ${currentUser?.name}</span>
    <button class="btn btn-secondary" onclick="logout()">Logout</button>
  `;

  try {
    const boards = await boardsAPI.getAll();
    
    mainContent.innerHTML = `
      <div class="container">
        <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 3rem; padding-bottom: 1rem; border-bottom: 1px solid #1a1a1a;">
          <h1 style="color: #ffffff; font-size: 2rem; font-weight: 600; letter-spacing: -0.02em;">Your Boards</h1>
          <button onclick="showCreateBoardModal()" class="btn">+ New Board</button>
        </div>

        <div class="boards-grid" style="display: grid; grid-template-columns: repeat(auto-fill, minmax(300px, 1fr)); gap: 1.5rem;">
          ${boards.length > 0 ? boards.map(board => `
            <div class="card board-card" style="cursor: pointer; transition: all 0.2s ease;" onclick="showBoard(${board.ID})">
              <h3 style="margin-bottom: 1rem; color: #ffffff; font-weight: 600; font-size: 1.1rem;">${board.Name}</h3>
              <p style="color: #94a3b8; margin-bottom: 1rem; font-size: 0.9rem; line-height: 1.4;">${board.Description}</p>
              <small style="color: #64748b; font-size: 0.8rem;">Created ${new Date(board.CreatedAt).toLocaleDateString()}</small>
            </div>
          `).join('') : `
            <div class="card" style="grid-column: 1 / -1; text-align: center; padding: 4rem;">
              <h3 style="color: #94a3b8; margin-bottom: 1rem; font-weight: 500;">No boards yet</h3>
              <p style="color: #64748b; margin-bottom: 2rem; font-size: 0.9rem;">Create your first board to get started!</p>
              <button onclick="showCreateBoardModal()" class="btn">Create Your First Board</button>
            </div>
          `}
        </div>
      </div>

      <!-- Create Board Modal -->
      <div id="createBoardModal" class="hidden" style="position: fixed; top: 0; left: 0; width: 100%; height: 100%; background: rgba(0,0,0,0.8); z-index: 1000;">
        <div style="position: absolute; top: 50%; left: 50%; transform: translate(-50%, -50%); background: #000000; border: 1px solid #1a1a1a; padding: 2rem; border-radius: 4px; width: 90%; max-width: 500px;">
          <h2 style="margin-bottom: 2rem; color: #ffffff; font-weight: 600;">Create New Board</h2>
          
          <form id="create-board-form">
            <div class="form-group">
              <label for="board-name">Board Name</label>
              <input type="text" id="board-name" name="name" required placeholder="e.g., Website Redesign">
            </div>

            <div class="form-group">
              <label for="board-description">Description (Optional)</label>
              <textarea id="board-description" name="description" rows="3" placeholder="What's this board for?"></textarea>
            </div>

            <div style="display: flex; gap: 1rem; justify-content: flex-end;">
              <button type="button" onclick="hideCreateBoardModal()" class="btn btn-secondary">Cancel</button>
              <button type="submit" class="btn">Create Board</button>
            </div>
          </form>
        </div>
      </div>
    `;

    // Add board card hover effects
    const style = document.createElement('style');
    style.textContent = `
      .board-card:hover {
        border-color: #ffffff;
      }
    `;
    document.head.appendChild(style);

    // Handle create board form
    const createBoardForm = document.getElementById('create-board-form') as HTMLFormElement;
    createBoardForm.addEventListener('submit', handleCreateBoard);

  } catch (error) {
    console.error('Error loading dashboard:', error);
    showError('Failed to load dashboard');
  }
}

// Event handlers
async function handleLogin(e: Event) {
  e.preventDefault();
  const form = e.target as HTMLFormElement;
  const formData = new FormData(form);
  const email = formData.get('email') as string;
  const password = formData.get('password') as string;

  try {
    currentUser = await authAPI.login(email, password);
    showDashboard();
  } catch (error: any) {
    showError('Invalid credentials', 'login-error');
  }
}

async function handleRegister(e: Event) {
  e.preventDefault();
  const form = e.target as HTMLFormElement;
  const formData = new FormData(form);
  const name = formData.get('name') as string;
  const email = formData.get('email') as string;
  const password = formData.get('password') as string;

  try {
    currentUser = await authAPI.register(name, email, password);
    showDashboard();
  } catch (error: any) {
    showError('Registration failed', 'register-error');
  }
}

async function handleCreateBoard(e: Event) {
  e.preventDefault();
  const form = e.target as HTMLFormElement;
  const formData = new FormData(form);
  const name = formData.get('name') as string;
  const description = formData.get('description') as string;

  try {
    await boardsAPI.create(name, description);
    hideCreateBoardModal();
    showDashboard(); // Refresh dashboard
  } catch (error: any) {
    showError('Failed to create board');
  }
}

async function handleCreateList(e: Event) {
  e.preventDefault();
  const form = e.target as HTMLFormElement;
  const formData = new FormData(form);
  const name = formData.get('name') as string;
  const boardId = (window as any).currentBoardId;

  try {
    // Get current lists to determine position
    const lists = await listsAPI.getByBoard(boardId);
    const position = lists.length;

    await listsAPI.create(boardId, name, position);
    hideCreateListModal();
    showBoard(boardId); // Refresh board view
  } catch (error: any) {
    showError('Failed to create list');
  }
}

async function handleCreateCard(e: Event) {
  e.preventDefault();
  const form = e.target as HTMLFormElement;
  const formData = new FormData(form);
  const title = formData.get('title') as string;
  const description = formData.get('description') as string;
  const listId = (window as any).currentListId;

  try {
    // Get current cards to determine position
    const cards = await cardsAPI.getByList(listId);
    const position = cards.length;

    await cardsAPI.create(listId, title, description || '', position);
    hideCreateCardModal();
    showBoard((window as any).currentBoardId); // Refresh board view
  } catch (error: any) {
    showError('Failed to create card');
  }
}

// Helper function to find the element after which to insert the dragged card
function getDragAfterElement(container: HTMLElement, y: number): HTMLElement | null {
  const draggableElements = [...container.querySelectorAll('.card-item:not(.dragging)')] as HTMLElement[];
  
  return draggableElements.reduce((closest, child) => {
    const box = child.getBoundingClientRect();
    const offset = y - box.top - box.height / 2;
    
    if (offset < 0 && offset > closest.offset) {
      return { offset: offset, element: child };
    } else {
      return closest;
    }
  }, { offset: Number.NEGATIVE_INFINITY, element: null as HTMLElement | null }).element;
}

// Drag and Drop functionality
function initializeDragAndDrop() {
  const cards = document.querySelectorAll('.card-item');
  const containers = document.querySelectorAll('.cards-container');

  // Add drag event listeners to cards
  cards.forEach(card => {
    card.addEventListener('dragstart', handleDragStart as EventListener);
    card.addEventListener('dragend', handleDragEnd as EventListener);
  });

  // Add drop event listeners to containers
  containers.forEach(container => {
    container.addEventListener('dragover', handleDragOver as EventListener);
    container.addEventListener('dragenter', handleDragEnter as EventListener);
    container.addEventListener('dragleave', handleDragLeave as EventListener);
    container.addEventListener('drop', handleDrop as unknown as EventListener);
  });
}

// let draggedCard: HTMLElement | null = null; // Unused for now but may be needed for advanced features
let draggedCardId: string | null = null;
let draggedFromListId: string | null = null;

function handleDragStart(e: DragEvent) {
  const card = e.target as HTMLElement;
  // draggedCard = card; // Store for potential future use
  draggedCardId = (card as any).dataset.cardId || null;
  draggedFromListId = (card.closest('.list') as any)?.dataset.listId || null;
  
  card.classList.add('dragging');
  
  if (e.dataTransfer) {
    e.dataTransfer.effectAllowed = 'move';
    e.dataTransfer.setData('text/html', card.outerHTML);
  }
}

function handleDragEnd(e: DragEvent) {
  const card = e.target as HTMLElement;
  card.classList.remove('dragging');
  
  // Clean up
  // draggedCard = null;
  draggedCardId = null;
  draggedFromListId = null;
  
  // Remove drag-over classes from all containers and cards
  document.querySelectorAll('.cards-container').forEach(container => {
    container.classList.remove('drag-over');
  });
  document.querySelectorAll('.list').forEach(list => {
    list.classList.remove('drag-over');
  });
  document.querySelectorAll('.card-item').forEach(card => {
    card.classList.remove('drag-over-top', 'drag-over-bottom');
  });
}

function handleDragOver(e: DragEvent) {
  e.preventDefault();
  if (e.dataTransfer) {
    e.dataTransfer.dropEffect = 'move';
  }
  
  // Add visual feedback for intra-list reordering
  const container = (e.target as HTMLElement).closest('.cards-container') as HTMLElement;
  if (container) {
    const afterElement = getDragAfterElement(container, e.clientY);
    
    // Remove previous indicators
    container.querySelectorAll('.card-item').forEach(card => {
      card.classList.remove('drag-over-top', 'drag-over-bottom');
    });
    
    // Add new indicator
    if (afterElement) {
      afterElement.classList.add('drag-over-top');
    } else {
      // If no afterElement, we're at the end
      const lastCard = container.querySelector('.card-item:last-child');
      if (lastCard) {
        lastCard.classList.add('drag-over-bottom');
      }
    }
  }
}

function handleDragEnter(e: DragEvent) {
  e.preventDefault();
  const container = e.target as HTMLElement;
  if (container.classList.contains('cards-container')) {
    container.classList.add('drag-over');
    container.closest('.list')?.classList.add('drag-over');
  }
}

function handleDragLeave(e: DragEvent) {
  const container = e.target as HTMLElement;
  if (container.classList.contains('cards-container')) {
    // Only remove if we're actually leaving the container
    const rect = container.getBoundingClientRect();
    const x = e.clientX;
    const y = e.clientY;
    
    if (x < rect.left || x > rect.right || y < rect.top || y > rect.bottom) {
      container.classList.remove('drag-over');
      container.closest('.list')?.classList.remove('drag-over');
    }
  }
}

async function handleDrop(e: DragEvent) {
  e.preventDefault();
  
  let container = e.target as HTMLElement;
  
  // If dropped on a card, find the container
  if (container.classList.contains('card-item')) {
    container = container.closest('.cards-container') as HTMLElement;
  }
  
  // If dropped on card content (title/description), find the container
  if (!container || !container.classList.contains('cards-container')) {
    container = (e.target as HTMLElement).closest('.cards-container') as HTMLElement;
  }
  
  if (!container || !container.classList.contains('cards-container') || !draggedCardId) {
    return;
  }
  
  const targetList = container.closest('.list');
  const targetListId = (targetList as any)?.dataset.listId;
  
  if (!targetListId) {
    return;
  }
  
  // Remove drag-over classes and position indicators
  container.classList.remove('drag-over');
  targetList?.classList.remove('drag-over');
  document.querySelectorAll('.card-item').forEach(card => {
    card.classList.remove('drag-over-top', 'drag-over-bottom');
  });
  
  try {
    let newPosition: number;
    
    // If dropped in the same list, calculate position based on drop location
    if (targetListId === draggedFromListId) {
      // Find the drop position based on mouse position
      const cards = Array.from(container.querySelectorAll('.card-item'));
      const afterElement = getDragAfterElement(container, e.clientY);
      
      if (afterElement == null) {
        newPosition = cards.length - 1; // Drop at end
      } else {
        const afterIndex = cards.indexOf(afterElement);
        newPosition = afterIndex;
      }
    } else {
      // Different list - append to end
      const targetCards = container.querySelectorAll('.card-item');
      newPosition = targetCards.length;
    }
    
    // Call API to move the card
    await cardsAPI.move(parseInt(draggedCardId), parseInt(targetListId), newPosition);
    
    // Refresh the board view
    showBoard((window as any).currentBoardId);
    
  } catch (error) {
    console.error('Error moving card:', error);
    showError('Failed to move card');
  }
}

// Utility functions
function showError(message: string, elementId?: string) {
  const errorElement = document.getElementById(elementId || 'error');
  if (errorElement) {
    errorElement.textContent = message;
    errorElement.classList.remove('hidden');
  } else {
    alert(message); // Fallback
  }
}

function logout() {
  currentUser = null;
  showLogin();
}

// Utility functions
function showCreateBoardModal() {
  const modal = document.getElementById('createBoardModal');
  if (modal) modal.classList.remove('hidden');
}

function hideCreateBoardModal() {
  const modal = document.getElementById('createBoardModal');
  if (modal) modal.classList.add('hidden');
}

function showCreateListModal() {
  const modal = document.getElementById('createListModal');
  if (modal) modal.classList.remove('hidden');
}

function hideCreateListModal() {
  const modal = document.getElementById('createListModal');
  if (modal) modal.classList.add('hidden');
}

function showCreateCardModal(listId: number) {
  (window as any).currentListId = listId;
  const modal = document.getElementById('createCardModal');
  if (modal) modal.classList.remove('hidden');
}

function hideCreateCardModal() {
  const modal = document.getElementById('createCardModal');
  if (modal) modal.classList.add('hidden');
}

// Global functions for onclick handlers
(window as any).showLogin = showLogin;
(window as any).showRegister = showRegister;
(window as any).showDashboard = showDashboard;
(window as any).logout = logout;
(window as any).showCreateBoardModal = showCreateBoardModal;
(window as any).hideCreateBoardModal = hideCreateBoardModal;
(window as any).showCreateListModal = showCreateListModal;
(window as any).hideCreateListModal = hideCreateListModal;
(window as any).showCreateCardModal = showCreateCardModal;
(window as any).hideCreateCardModal = hideCreateCardModal;
(window as any).showBoard = showBoard;

// Start the app
init();
