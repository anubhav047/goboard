import './styles/main.css';
import { authAPI, boardsAPI, User } from './api/client';
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

// Global functions for onclick handlers
(window as any).showLogin = showLogin;
(window as any).showRegister = showRegister;
(window as any).logout = logout;
(window as any).showCreateBoardModal = showCreateBoardModal;
(window as any).hideCreateBoardModal = hideCreateBoardModal;
(window as any).showBoard = (boardId: number) => {
  console.log('Show board:', boardId);
  // TODO: Implement board view
};

// Start the app
init();
