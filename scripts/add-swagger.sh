#!/bin/bash

# Script untuk menambahkan Swagger annotations ke semua controllers

echo "Adding Swagger annotations to controllers..."

# Function to add swagger annotations to a file
add_swagger_to_controller() {
    local file=$1
    local controller_name=$2
    
    echo "Processing $file for $controller_name..."
    
    # Add swagger annotations for each function
    # This is a template - you'll need to customize for each controller
    
    echo "Manual annotation needed for $file"
    echo "Please add @Summary, @Description, @Tags, @Param, @Success, @Failure annotations"
    echo "Example template saved in docs/swagger-template.md"
}

# Create swagger template
cat > docs/swagger-template.md << 'EOF'
# Swagger Annotation Template

## Basic Structure
```go
// FunctionName godoc
// @Summary Brief description
// @Description Detailed description
// @Tags TagName
// @Accept json
// @Produce json
// @Param id path string true "ID parameter"
// @Param body body models.ModelName true "Request body"
// @Success 200 {object} models.ModelName
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/endpoint [method]
func FunctionName(c *fiber.Ctx) error {
    // function body
}
```

## Tags to Use
- Films (for filmController)
- Schedules (for jadwalController)
- Tickets (for tiketController)
- Users (for userController)
- Payments (for pembayaranController)

## Common Responses
- 200: Success with data
- 201: Created successfully
- 400: Bad request / Invalid input
- 404: Not found
- 500: Internal server error
EOF

# Process each controller
add_swagger_to_controller "controllers/jadwalController.go" "Schedules"
add_swagger_to_controller "controllers/tiketController.go" "Tickets"
add_swagger_to_controller "controllers/userController.go" "Users"
add_swagger_to_controller "controllers/pembayaranController.go" "Payments"

echo "Swagger annotations template created!"
echo "Please manually add annotations to remaining controllers using the template in docs/swagger-template.md"
echo "After adding annotations, run: make docs"
