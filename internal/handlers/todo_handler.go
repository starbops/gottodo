package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/starbops/gottodo/internal/services"
)

// TodoHandler handles HTTP requests for todos
type TodoHandler struct {
	service *services.TodoService
}

// NewTodoHandler creates a new TodoHandler
func NewTodoHandler(service *services.TodoService) *TodoHandler {
	return &TodoHandler{
		service: service,
	}
}

// CreateTodoRequest represents the request body for creating a todo
type CreateTodoRequest struct {
	Title       string `json:"title" form:"title"`
	Description string `json:"description" form:"description"`
}

// UpdateTodoRequest represents the request body for updating a todo
type UpdateTodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// GetAllTodos handles GET /todos
func (h *TodoHandler) GetAllTodos(c echo.Context) error {
	userID := c.Get("user_id").(string)
	todos, err := h.service.GetUserTodos(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, todos)
}

// GetTodo handles GET /todos/:id
func (h *TodoHandler) GetTodo(c echo.Context) error {
	id := c.Param("id")
	todo, err := h.service.GetTodo(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	// Check if the todo belongs to the authenticated user
	userID := c.Get("user_id").(string)
	if todo.UserID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "You don't have permission to access this todo",
		})
	}

	return c.JSON(http.StatusOK, todo)
}

// CreateTodo handles POST /todos
func (h *TodoHandler) CreateTodo(c echo.Context) error {
	// Get form values directly
	title := c.FormValue("title")
	description := c.FormValue("description")

	// Validate required fields
	if title == "" || description == "" {
		// Return a user-friendly error with the same structure as success
		errorHTML := `
		<div id="todo-list" class="bg-white rounded-lg shadow-md p-6">
			<h2 class="text-xl font-semibold mb-4">Your Todos</h2>
			<div class="space-y-4">
				<div class="bg-red-100 text-red-800 p-4 rounded-lg mb-4">
					<p>Error: Title and description are required.</p>
				</div>
				<!-- Re-fetch existing todos to keep the list current -->
				`

		// Get todos to display (ensure list is still populated if there are existing todos)
		userID := c.Get("user_id").(string)
		todos, err := h.service.GetUserTodos(c.Request().Context(), userID)
		if err == nil && len(todos) > 0 {
			for _, t := range todos {
				completedClass := ""
				titleClass := ""
				buttonIcon := ""

				if t.Completed {
					completedClass = "bg-gray-100"
					titleClass = "line-through text-gray-500"
					buttonIcon = `
						<button class="text-yellow-500 hover:text-yellow-700 mr-2" hx-put="/todos/` + t.ID + `/incomplete" hx-swap="outerHTML" hx-target="#todo-` + t.ID + `">
							<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
								<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.28 7.22a.75.75 0 00-1.06 1.06L8.94 10l-1.72 1.72a.75.75 0 101.06 1.06L10 11.06l1.72 1.72a.75.75 0 101.06-1.06L11.06 10l1.72-1.72a.75.75 0 00-1.06-1.06L10 8.94 8.28 7.22z" clip-rule="evenodd" />
							</svg>
						</button>
					`
				} else {
					buttonIcon = `
						<button class="text-green-500 hover:text-green-700 mr-2" hx-put="/todos/` + t.ID + `/complete" hx-swap="outerHTML" hx-target="#todo-` + t.ID + `">
							<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
								<path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
							</svg>
						</button>
					`
				}

				errorHTML += `
				<div class="border rounded-lg p-4 bg-white shadow-sm mb-4 ` + completedClass + `" id="todo-` + t.ID + `">
					<div class="flex justify-between items-start">
						<div>
							<h3 class="font-semibold text-lg ` + titleClass + `">` + t.Title + `</h3>
							<p class="text-gray-600 mt-1">` + t.Description + `</p>
						</div>
						<div class="flex">
							` + buttonIcon + `
							<button class="text-red-500 hover:text-red-700" hx-delete="/todos/` + t.ID + `" hx-swap="outerHTML" hx-target="#todo-` + t.ID + `" hx-confirm="Are you sure you want to delete this todo?">
								<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
									<path fill-rule="evenodd" d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z" clip-rule="evenodd" />
								</svg>
							</button>
						</div>
					</div>
				</div>
				`
			}
		} else {
			errorHTML += `<p class="text-gray-500 text-center">No todos yet. Add one above!</p>`
		}

		errorHTML += `
			</div>
		</div>
		`

		// Return the error with the todo list HTML
		return c.HTML(http.StatusBadRequest, errorHTML)
	}

	userID := c.Get("user_id").(string)
	_, err := h.service.CreateTodo(c.Request().Context(), userID, title, description)
	if err != nil {
		// Handle service errors with user-friendly HTML
		return c.HTML(http.StatusInternalServerError, `
		<div id="todo-list" class="bg-white rounded-lg shadow-md p-6">
			<h2 class="text-xl font-semibold mb-4">Your Todos</h2>
			<div class="space-y-4">
				<div class="bg-red-100 text-red-800 p-4 rounded-lg mb-4">
					<p>Error: Unable to create todo. Please try again later.</p>
				</div>
				<p class="text-gray-500 text-center">Refresh the page to see your current todos.</p>
			</div>
		</div>
		`)
	}

	// Return a refresh of the todo list - this allows the frontend to update
	todos, err := h.service.GetUserTodos(c.Request().Context(), userID)
	if err != nil {
		return c.HTML(http.StatusInternalServerError, `
		<div id="todo-list" class="bg-white rounded-lg shadow-md p-6">
			<h2 class="text-xl font-semibold mb-4">Your Todos</h2>
			<div class="space-y-4">
				<div class="bg-red-100 text-red-800 p-4 rounded-lg mb-4">
					<p>Error: Unable to retrieve todos. Please refresh the page.</p>
				</div>
			</div>
		</div>
		`)
	}

	// Generate HTML for todos
	todoHTML := ""
	if len(todos) == 0 {
		todoHTML = `<p class="text-gray-500 text-center">No todos yet. Add one above!</p>`
	} else {
		for _, t := range todos {
			completedClass := ""
			titleClass := ""
			buttonIcon := ""

			if t.Completed {
				completedClass = "bg-gray-100"
				titleClass = "line-through text-gray-500"
				buttonIcon = `
					<button class="text-yellow-500 hover:text-yellow-700 mr-2" hx-put="/todos/` + t.ID + `/incomplete" hx-swap="outerHTML" hx-target="#todo-` + t.ID + `">
						<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
							<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.28 7.22a.75.75 0 00-1.06 1.06L8.94 10l-1.72 1.72a.75.75 0 101.06 1.06L10 11.06l1.72 1.72a.75.75 0 101.06-1.06L11.06 10l1.72-1.72a.75.75 0 00-1.06-1.06L10 8.94 8.28 7.22z" clip-rule="evenodd" />
						</svg>
					</button>
				`
			} else {
				buttonIcon = `
					<button class="text-green-500 hover:text-green-700 mr-2" hx-put="/todos/` + t.ID + `/complete" hx-swap="outerHTML" hx-target="#todo-` + t.ID + `">
						<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
							<path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
						</svg>
					</button>
				`
			}

			todoHTML += `
			<div class="border rounded-lg p-4 bg-white shadow-sm mb-4 ` + completedClass + `" id="todo-` + t.ID + `">
				<div class="flex justify-between items-start">
					<div>
						<h3 class="font-semibold text-lg ` + titleClass + `">` + t.Title + `</h3>
						<p class="text-gray-600 mt-1">` + t.Description + `</p>
					</div>
					<div class="flex">
						` + buttonIcon + `
						<button class="text-red-500 hover:text-red-700" hx-delete="/todos/` + t.ID + `" hx-swap="outerHTML" hx-target="#todo-` + t.ID + `" hx-confirm="Are you sure you want to delete this todo?">
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

	// Important: we're returning the entire todo-list div with its header to preserve the structure
	htmlResponse := `
		<div id="todo-list" class="bg-white rounded-lg shadow-md p-6">
			<h2 class="text-xl font-semibold mb-4">Your Todos</h2>
			<div class="space-y-4">
				` + todoHTML + `
			</div>
		</div>
	`

	// Add HX-Trigger header to ensure client-side events fire
	c.Response().Header().Set("HX-Trigger", "todoCreated")

	return c.HTML(http.StatusOK, htmlResponse)
}

// UpdateTodo handles PUT /todos/:id
func (h *TodoHandler) UpdateTodo(c echo.Context) error {
	id := c.Param("id")
	var req UpdateTodoRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Check if the todo exists and belongs to the authenticated user
	todo, err := h.service.GetTodo(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	userID := c.Get("user_id").(string)
	if todo.UserID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "You don't have permission to update this todo",
		})
	}

	updatedTodo, err := h.service.UpdateTodo(c.Request().Context(), id, req.Title, req.Description)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, updatedTodo)
}

// DeleteTodo handles DELETE /todos/:id
func (h *TodoHandler) DeleteTodo(c echo.Context) error {
	id := c.Param("id")

	// Check if the todo exists and belongs to the authenticated user
	todo, err := h.service.GetTodo(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	userID := c.Get("user_id").(string)
	if todo.UserID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "You don't have permission to delete this todo",
		})
	}

	err = h.service.DeleteTodo(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	// Return an empty div with a success message that will fade out
	return c.HTML(http.StatusOK, `
		<div id="todo-`+id+`" class="bg-green-100 text-green-800 p-4 rounded-lg mb-4 fade-out">
			<p>Todo successfully deleted!</p>
			<script>
				setTimeout(function() {
					const element = document.getElementById("todo-`+id+`");
					element.style.opacity = "0";
					element.style.transition = "opacity 0.5s";
					setTimeout(function() {
						element.remove();
					}, 500);
				}, 1000);
			</script>
		</div>
	`)
}

// CompleteTodo handles PUT /todos/:id/complete
func (h *TodoHandler) CompleteTodo(c echo.Context) error {
	id := c.Param("id")

	// Check if the todo exists and belongs to the authenticated user
	todo, err := h.service.GetTodo(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	userID := c.Get("user_id").(string)
	if todo.UserID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "You don't have permission to update this todo",
		})
	}

	updatedTodo, err := h.service.CompleteTodo(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	// Return HTML representation of the updated todo
	completedClass := ""
	buttonIcon := ""

	if updatedTodo.Completed {
		completedClass = "bg-gray-100"
		buttonIcon = `
			<button class="text-yellow-500 hover:text-yellow-700 mr-2" hx-put="/todos/` + updatedTodo.ID + `/incomplete" hx-swap="outerHTML" hx-target="#todo-` + updatedTodo.ID + `">
				<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
					<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.28 7.22a.75.75 0 00-1.06 1.06L8.94 10l-1.72 1.72a.75.75 0 101.06 1.06L10 11.06l1.72 1.72a.75.75 0 101.06-1.06L11.06 10l1.72-1.72a.75.75 0 00-1.06-1.06L10 8.94 8.28 7.22z" clip-rule="evenodd" />
				</svg>
			</button>
		`
	} else {
		buttonIcon = `
			<button class="text-green-500 hover:text-green-700 mr-2" hx-put="/todos/` + updatedTodo.ID + `/complete" hx-swap="outerHTML" hx-target="#todo-` + updatedTodo.ID + `">
				<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
					<path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
				</svg>
			</button>
		`
	}

	titleClass := ""
	if updatedTodo.Completed {
		titleClass = "line-through text-gray-500"
	}

	return c.HTML(http.StatusOK, `
		<div class="border rounded-lg p-4 bg-white shadow-sm mb-4 `+completedClass+`" id="todo-`+updatedTodo.ID+`">
			<div class="flex justify-between items-start">
				<div>
					<h3 class="font-semibold text-lg `+titleClass+`">`+updatedTodo.Title+`</h3>
					<p class="text-gray-600 mt-1">`+updatedTodo.Description+`</p>
				</div>
				<div class="flex">
					`+buttonIcon+`
					<button class="text-red-500 hover:text-red-700" hx-delete="/todos/`+updatedTodo.ID+`" hx-swap="outerHTML" hx-target="#todo-`+updatedTodo.ID+`" hx-confirm="Are you sure you want to delete this todo?">
						<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
							<path fill-rule="evenodd" d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z" clip-rule="evenodd" />
						</svg>
					</button>
				</div>
			</div>
		</div>
	`)
}

// IncompleteTodo marks a todo as not completed
func (h *TodoHandler) IncompleteTodo(c echo.Context) error {
	id := c.Param("id")

	// Check if the todo exists and belongs to the authenticated user
	todo, err := h.service.GetTodo(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	userID := c.Get("user_id").(string)
	if todo.UserID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "You don't have permission to update this todo",
		})
	}

	updatedTodo, err := h.service.IncompleteTodo(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	// Return HTML representation of the updated todo
	buttonIcon := `
		<button class="text-green-500 hover:text-green-700 mr-2" hx-put="/todos/` + updatedTodo.ID + `/complete" hx-swap="outerHTML" hx-target="#todo-` + updatedTodo.ID + `">
			<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
				<path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
			</svg>
		</button>
	`

	return c.HTML(http.StatusOK, `
		<div class="border rounded-lg p-4 bg-white shadow-sm mb-4" id="todo-`+updatedTodo.ID+`">
			<div class="flex justify-between items-start">
				<div>
					<h3 class="font-semibold text-lg">`+updatedTodo.Title+`</h3>
					<p class="text-gray-600 mt-1">`+updatedTodo.Description+`</p>
				</div>
				<div class="flex">
					`+buttonIcon+`
					<button class="text-red-500 hover:text-red-700" hx-delete="/todos/`+updatedTodo.ID+`" hx-swap="outerHTML" hx-target="#todo-`+updatedTodo.ID+`" hx-confirm="Are you sure you want to delete this todo?">
						<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
							<path fill-rule="evenodd" d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z" clip-rule="evenodd" />
						</svg>
					</button>
				</div>
			</div>
		</div>
	`)
}
