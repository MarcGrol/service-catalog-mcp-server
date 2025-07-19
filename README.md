# Claude Desktop Configuration

## 1. Build and Prepare Your Server
 
First, build your advanced MCP server

    go install ./...

Make sure it works with stdio (not HTTP for Claude Desktop)

    ./learnmcp

## 2. Configure Claude Desktop

    # Edit your Claude Desktop configuration file:
    # macOS: ~/.config/claude-desktop/claude_desktop_config.json
    # Windows: %APPDATA%\Claude\claude_desktop_config.json

    {
        "mcpServers": {
            "aproject-server": {
            "command": "/Users/marcgrol/go/bin//learnmcp",
            "args": [],
            "env": {}
            }
        }
    }

IMPORTANT: Use the full absolute path to your compiled binary!
Example: "/Users/yourusername/projects/mcp-demo/advanced-mcp-server"

## 3. Restart Claude Desktop

After saving the config, completely quit and restart Claude Desktop
heck the Developer tab in settings to see if your server connected

## 4. Test Prompts in Claude Desktop

### TOOLS EXPLORATION PROMPTS

"What tools are available in the advanced project manager?"

"Create a new project called 'AI Assistant Platform' for building intelligent chatbots. The project should be described as 'A comprehensive platform for creating and deploying AI-powered conversational agents' and list the authors as 'Alice Johnson, Bob Smith, Carol Davis'."

"Now create a high-priority task called 'Implement user authentication' with the description 'Add OAuth 2.0 login functionality with Google and GitHub providers' and set the due date to 2025-08-15."

"Create another task for 'Design system architecture' with medium priority and description 'Create system design for scalable microservices architecture'."

"Generate a project summary analytics report to see the current status."

"Generate a tasks analytics report to understand the task distribution."

"Search for anything related to 'authentication' in the project."

"What happens if I try to create a project without providing a name?"

### RESOURCES EXPLORATION PROMPTS  

"What resources are available in the project manager?"

"Show me the current project configuration."

"Display the list of all tasks in the project."

"What are the current project statistics and metrics?"

"Read the project documentation for me."

"Compare the project configuration with the task list - what insights can you gather?"

"Based on the project statistics, what recommendations do you have for improving productivity?"

### PROMPTS EXPLORATION PROMPTS

"What prompts are available in the project manager?"

"I want to plan a new web application project that will take 8 weeks to complete using Go, React, PostgreSQL, and Docker. Give me detailed project planning guidance."

"Help me plan a mobile app project using React Native and Firebase for a 6-week timeline."

"I need a code review checklist for Go code with a focus on security best practices."

"Generate a code review template for Python code focusing on performance optimization."

"Help me plan a 2-week sprint for a team of 5 developers."

"Create a sprint plan for a 10-day sprint with 3 team members."

"What's the difference between the project planning prompt for a web project vs an API project?"

### INTEGRATION AND WORKFLOW PROMPTS

"Walk me through creating a complete project workflow: create a project, add some tasks, check the resources, and then use prompts to plan the next steps."

"How would you use the available tools and resources to manage a real software development project?"

"Based on the current project data, use the sprint planning prompt to help organize the next development cycle."

"Create a task, then immediately check the task list resource to see how it's stored."

"Use the project planning prompt to suggest improvements based on the current project configuration."

"What's missing from this project management system that you would need for real-world use?"

### UNDERSTANDING MCP CONCEPTS

"Explain the difference between the tools, resources, and prompts in this system. Give me specific examples from what's available."

"How do the three MCP capabilities (tools, resources, prompts) work together in this project manager?"

"Which of these would you use to: 1) Add a new feature request, 2) Check project status, 3) Get planning advice?"

"Show me how the same information appears differently when accessed as a tool result vs a resource vs incorporated into a prompt."

"If you were building an AI assistant for project management, how would you use these MCP capabilities?"

### TESTING EDGE CASES AND PARAMETERS

"Try to create a project with comma-separated authors and see how the system handles it."

"Create a task with an invalid due date format and see what happens."

"Use the project planning prompt with unusual parameters like a 1-day timeline or 20 team members."

"Test the search functionality with terms that don't exist in the project."

"Generate an analytics report with an invalid report type."

## 5. Observing MCP in Action

Watch for these behaviors in Claude Desktop:

✅ **Tool Calls**: Claude will show "Using tool: create_project" before executing
✅ **Resource Access**: Claude will read resources automatically for context  
✅ **Prompt Usage**: Claude will use prompts to generate better responses
✅ **Error Handling**: See how Claude handles invalid parameters
✅ **Data Flow**: Notice how tools create data that resources then expose

## 6. Advanced Exploration

Once you're comfortable with the basics:

"Help me design a workflow that uses all three MCP capabilities together to manage a software project from inception to deployment."

"Based on this MCP server example, how would you extend it to integrate with real systems like GitHub, Jira, and Slack?"

"What are the security considerations when building MCP servers for production use?"

"How would you modify this to support multiple projects and user permissions?"