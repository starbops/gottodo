package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/starbops/gottodo/internal/services"
	"github.com/starbops/gottodo/pkg/auth"
)

// PageHandler handles HTTP requests for HTML pages
type PageHandler struct {
	todoService *services.TodoService
	authService *auth.AuthService
}

// NewPageHandler creates a new PageHandler
func NewPageHandler(todoService *services.TodoService, authService *auth.AuthService) *PageHandler {
	return &PageHandler{
		todoService: todoService,
		authService: authService,
	}
}

// Home handles GET /
func (h *PageHandler) Home(c echo.Context) error {
	// In a real implementation, this would render the home page template
	return c.HTML(http.StatusOK, `
		<html>
			<head>
				<title>GoTToDo - A Simple Todo App</title>
				<meta name="viewport" content="width=device-width, initial-scale=1">
				<script src="https://cdn.tailwindcss.com"></script>
				<script src="https://unpkg.com/htmx.org@1.9.10"></script>
				<style>
					.htmx-indicator {
						display: none;
					}
					.htmx-request .htmx-indicator {
						display: inline;
					}
					.htmx-request.htmx-indicator {
						display: inline;
					}
				</style>
			</head>
			<body class="bg-gray-100 min-h-screen">
				<div class="container mx-auto px-4 py-8">
					<h1 class="text-3xl font-bold text-center mb-8">GoTToDo</h1>
					<div class="max-w-md mx-auto bg-white rounded-lg shadow-md p-6">
						<p class="text-gray-700 mb-4">A simple todo app built with Go, Templ, Tailwind CSS, and HTMX.</p>
						<div class="flex flex-col space-y-4">
							<a href="/auth/github" class="bg-gray-900 hover:bg-gray-800 text-white font-semibold py-2 px-4 rounded flex items-center justify-center">
								<svg class="w-5 h-5 mr-2" fill="currentColor" viewBox="0 0 24 24" aria-hidden="true">
									<path fill-rule="evenodd" d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z" clip-rule="evenodd"></path>
								</svg>
								Login with GitHub
							</a>
							<div class="flex justify-between">
								<a href="/login" class="bg-blue-500 hover:bg-blue-600 text-white font-semibold py-2 px-4 rounded w-[48%] text-center">Login</a>
								<a href="/register" class="bg-green-500 hover:bg-green-600 text-white font-semibold py-2 px-4 rounded w-[48%] text-center">Register</a>
							</div>
						</div>
					</div>
				</div>
			</body>
		</html>
	`)
}

// Login handles GET /login
func (h *PageHandler) Login(c echo.Context) error {
	// In a real implementation, this would render the login page template
	return c.HTML(http.StatusOK, `
		<html>
			<head>
				<title>Login - GoTToDo</title>
				<meta name="viewport" content="width=device-width, initial-scale=1">
				<script src="https://cdn.tailwindcss.com"></script>
				<script src="https://unpkg.com/htmx.org@1.9.10"></script>
				<style>
					.htmx-indicator {
						display: none;
					}
					.htmx-request .htmx-indicator {
						display: inline;
					}
					.htmx-request.htmx-indicator {
						display: inline;
					}
				</style>
			</head>
			<body class="bg-gray-100 min-h-screen">
				<div class="container mx-auto px-4 py-8">
					<h1 class="text-3xl font-bold text-center mb-8">Login</h1>
					<div class="max-w-md mx-auto bg-white rounded-lg shadow-md p-6">
						<a href="/auth/github" class="bg-gray-900 hover:bg-gray-800 text-white font-semibold py-2 px-4 rounded flex items-center justify-center mb-4">
							<svg class="w-5 h-5 mr-2" fill="currentColor" viewBox="0 0 24 24" aria-hidden="true">
								<path fill-rule="evenodd" d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z" clip-rule="evenodd"></path>
							</svg>
							Login with GitHub
						</a>
						<div class="text-center mb-4">
							<span class="text-gray-500">Or login with email</span>
						</div>
						<form hx-post="/auth/login" hx-swap="outerHTML">
							<div class="mb-4">
								<label class="block text-gray-700 text-sm font-bold mb-2" for="email">Email</label>
								<input class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" id="email" name="email" type="email" placeholder="Email">
							</div>
							<div class="mb-6">
								<label class="block text-gray-700 text-sm font-bold mb-2" for="password">Password</label>
								<input class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" id="password" name="password" type="password" placeholder="Password">
							</div>
							<div class="flex items-center justify-between">
								<button class="bg-blue-500 hover:bg-blue-600 text-white font-semibold py-2 px-4 rounded focus:outline-none focus:shadow-outline" type="submit">Sign In</button>
								<a class="inline-block align-baseline font-bold text-sm text-blue-500 hover:text-blue-800" href="/register">Don't have an account?</a>
							</div>
						</form>
					</div>
				</div>
			</body>
		</html>
	`)
}

// Register handles GET /register
func (h *PageHandler) Register(c echo.Context) error {
	// In a real implementation, this would render the registration page template
	return c.HTML(http.StatusOK, `
		<html>
			<head>
				<title>Register - GoTToDo</title>
				<meta name="viewport" content="width=device-width, initial-scale=1">
				<script src="https://cdn.tailwindcss.com"></script>
				<script src="https://unpkg.com/htmx.org@1.9.10"></script>
				<style>
					.htmx-indicator {
						display: none;
					}
					.htmx-request .htmx-indicator {
						display: inline;
					}
					.htmx-request.htmx-indicator {
						display: inline;
					}
				</style>
			</head>
			<body class="bg-gray-100 min-h-screen">
				<div class="container mx-auto px-4 py-8">
					<h1 class="text-3xl font-bold text-center mb-8">Register</h1>
					<div class="max-w-md mx-auto bg-white rounded-lg shadow-md p-6">
						<a href="/auth/github" class="bg-gray-900 hover:bg-gray-800 text-white font-semibold py-2 px-4 rounded flex items-center justify-center mb-4">
							<svg class="w-5 h-5 mr-2" fill="currentColor" viewBox="0 0 24 24" aria-hidden="true">
								<path fill-rule="evenodd" d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z" clip-rule="evenodd"></path>
							</svg>
							Register with GitHub
						</a>
						<div class="text-center mb-4">
							<span class="text-gray-500">Or register with email</span>
						</div>
						<form hx-post="/auth/register" hx-swap="outerHTML">
							<div class="mb-4">
								<label class="block text-gray-700 text-sm font-bold mb-2" for="email">Email</label>
								<input class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" id="email" name="email" type="email" placeholder="Email">
							</div>
							<div class="mb-6">
								<label class="block text-gray-700 text-sm font-bold mb-2" for="password">Password</label>
								<input class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" id="password" name="password" type="password" placeholder="Password">
							</div>
							<div class="flex items-center justify-between">
								<button class="bg-green-500 hover:bg-green-600 text-white font-semibold py-2 px-4 rounded focus:outline-none focus:shadow-outline" type="submit">Register</button>
								<a class="inline-block align-baseline font-bold text-sm text-blue-500 hover:text-blue-800" href="/login">Already have an account?</a>
							</div>
						</form>
					</div>
				</div>
			</body>
		</html>
	`)
}

// Dashboard handles GET /dashboard
func (h *PageHandler) Dashboard(c echo.Context) error {
	// Get user from context
	userID := c.Get("user_id").(string)

	// Get todos for the user - in a real implementation, we would use these to render the page
	todos, err := h.todoService.GetUserTodos(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	// Generate HTML for todos
	todoHTML := ""
	if len(todos) == 0 {
		todoHTML = `<p class="text-gray-500 text-center">No todos yet. Add one above!</p>`
	} else {
		for _, todo := range todos {
			completedClass := ""
			titleClass := ""
			buttonIcon := ""

			if todo.Completed {
				completedClass = "bg-gray-100"
				titleClass = "line-through text-gray-500"
				buttonIcon = `
					<button class="text-yellow-500 hover:text-yellow-700 mr-2" hx-put="/todos/` + todo.ID + `/incomplete" hx-swap="outerHTML" hx-target="#todo-` + todo.ID + `">
						<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
							<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.28 7.22a.75.75 0 00-1.06 1.06L8.94 10l-1.72 1.72a.75.75 0 101.06 1.06L10 11.06l1.72 1.72a.75.75 0 101.06-1.06L11.06 10l1.72-1.72a.75.75 0 00-1.06-1.06L10 8.94 8.28 7.22z" clip-rule="evenodd" />
						</svg>
					</button>
				`
			} else {
				buttonIcon = `
					<button class="text-green-500 hover:text-green-700 mr-2" hx-put="/todos/` + todo.ID + `/complete" hx-swap="outerHTML" hx-target="#todo-` + todo.ID + `">
						<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
							<path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
						</svg>
					</button>
				`
			}

			todoHTML += `
			<div class="border rounded-lg p-4 bg-white shadow-sm mb-4 ` + completedClass + `" id="todo-` + todo.ID + `">
				<div class="flex justify-between items-start">
					<div>
						<h3 class="font-semibold text-lg ` + titleClass + `">` + todo.Title + `</h3>
						<p class="text-gray-600 mt-1">` + todo.Description + `</p>
					</div>
					<div class="flex">
						` + buttonIcon + `
						<button class="text-red-500 hover:text-red-700" hx-delete="/todos/` + todo.ID + `" hx-swap="outerHTML" hx-target="#todo-` + todo.ID + `" hx-confirm="Are you sure you want to delete this todo?">
							<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
								<path fill-rule="evenodd" d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z" clip-rule="evenodd" />
							</svg>
						</button>
					</div>
				</div>
			</div>
			`
		}
	}

	// In a real implementation, this would render the dashboard template with the todos
	return c.HTML(http.StatusOK, `
		<html>
			<head>
				<title>Dashboard - GoTToDo</title>
				<meta name="viewport" content="width=device-width, initial-scale=1">
				<script src="https://cdn.tailwindcss.com"></script>
				<script src="https://unpkg.com/htmx.org@1.9.10"></script>
				<style>
					.htmx-indicator {
						display: none;
					}
					.htmx-request .htmx-indicator {
						display: inline;
					}
					.htmx-request.htmx-indicator {
						display: inline;
					}
				</style>
				<script>
					// Listen for successful form submission
					document.addEventListener('DOMContentLoaded', function() {
						// Add HTMX event listener for after the swap completes
						document.body.addEventListener('htmx:afterSwap', function(event) {
							// Check if the swap target was the todo list
							if (event.detail.target.id === 'todo-list') {
								// Clear the form AFTER data has been sent successfully
								const form = document.getElementById('todo-form');
								if (form) {
									// Reset form
									form.reset();
									
									// Show success message
									const message = document.getElementById('form-message');
									if (message) {
										message.classList.remove('hidden');
										message.textContent = "Todo added successfully!";
										
										// Hide the message after 2 seconds
										setTimeout(function() {
											message.classList.add('hidden');
										}, 2000);
									}
								}
							}
						});
					});
				</script>
			</head>
			<body class="bg-gray-100 min-h-screen">
				<div class="container mx-auto px-4 py-8">
					<div class="flex justify-between items-center mb-8">
						<h1 class="text-3xl font-bold">Your Todos</h1>
						<form hx-post="/auth/logout" hx-swap="none">
							<button class="bg-red-500 hover:bg-red-600 text-white font-semibold py-2 px-4 rounded">Logout</button>
						</form>
					</div>
					<div class="bg-white rounded-lg shadow-md p-6 mb-6">
						<h2 class="text-xl font-semibold mb-4">Add New Todo</h2>
						<form id="todo-form" hx-post="/todos" hx-swap="outerHTML" hx-target="#todo-list" hx-headers='{"Content-Type": "application/x-www-form-urlencoded"}' hx-indicator="#form-indicator">
							<div class="mb-4">
								<label class="block text-gray-700 text-sm font-bold mb-2" for="title">Title</label>
								<input class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" id="title" name="title" type="text" placeholder="Todo title" required>
							</div>
							<div class="mb-4">
								<label class="block text-gray-700 text-sm font-bold mb-2" for="description">Description</label>
								<textarea class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" id="description" name="description" placeholder="Todo description" required></textarea>
							</div>
							<div class="flex items-center">
								<button class="bg-blue-500 hover:bg-blue-600 text-white font-semibold py-2 px-4 rounded focus:outline-none focus:shadow-outline" type="submit">
									Add Todo
									<span id="form-indicator" class="htmx-indicator ml-2">
										<svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white inline" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
											<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
											<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
										</svg>
									</span>
								</button>
								<span id="form-message" class="ml-4 text-green-600 hidden">Todo added successfully!</span>
							</div>
						</form>
					</div>
					<div id="todo-list" class="bg-white rounded-lg shadow-md p-6">
						<h2 class="text-xl font-semibold mb-4">Your Todos</h2>
						<div class="space-y-4">
							`+todoHTML+`
						</div>
					</div>
				</div>
			</body>
		</html>
	`)
}
