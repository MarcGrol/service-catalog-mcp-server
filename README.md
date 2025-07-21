# Service-catalog mcp-server

## Context
This experimental "mcp-server" that makes the Adyen "service-catalog"accessible via normal text.
The service-catalog itself is harvested from the adyen-main codebase and describes all our software modules, with their owners, databases, exposed interfaces, consumed interfaces and jobs.

## What it exposes?

### Module Management Tools

#### 1. ```list_modules```

Lists all modules in your service catalog (optionally filtered onm keyword) 
Shows module names and short descriptions

Usage: ```list_module <keyword>``` (e.g. "kyc")

#### 2. ```get_module```

Get detailed information about a specific module
Shows: lines of code, file count, teams, exposed/consumed interfaces, databases, jobs

Usage: ```get_module <module_id>``` (e.g., "psp", "partner", "adyen")

#### 3. ```list_modules_of_teams```

Find all modules owned by a specific team

Usage: ```list_modules_of_teams <team_id>``` (e.g., "PartnerExperience")

### Interface Management Tools

#### 4. ```list_interfaces```

Lists all interfaces/APIs in your catalog (1,649 interfaces)
Shows interface IDs and descriptions

Usage: ```list_interfaces```

#### 5. ```get_interface```

Get detailed information about a specific interface
Shows: description, type, methods, specifications

Usage: ```get_interface <interface_id>``` (e.g., "com.adyen.services.acm.AcmService")

#### 6. ```list_interface_consumers```

Find all modules that consume/depend on a specific interface

Usage: ```list_interface_consumers``` <interface_id>

### Database Dependency Tools

#### 7. ```list_database_consumers```

Find all modules that use a specific database

Usage: ```list_database_consumers <database_id>``` (e.g., "partner", "config")

## Analysis Capabilities

### What You Can Discover:

- Architecture Overview: Complete module and interface catalog
- Code Metrics: Lines of code, file counts per module
- Team Ownership: Which teams own which modules
- Dependencies: Module-to-interface and module-to-database relationships
- API Catalog: Complete list of 1,649 available APIs/services
- Service Details: Method signatures, specifications for any interface

### Example Analysis Workflows:

- Find largest modules by code size
- Map team responsibilities across modules
- Trace interface dependencies across the system
- Identify database usage patterns
- Explore API capabilities and versions
- These tools give you comprehensive visibility into your Adyen service catalog architecture, dependencies, and ownership patterns. You can use them to understand system 
complexity, plan refactoring, analyze team boundaries, or explore available APIs for integration work.
