# Overview

A production-grade expense tracking backend API built with Go, following [Clean Architecture principles](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html). This service powers the [Expense Tracker frontend application](https://github.com/eyo-chen/expense-tracker).

## Technology Stack

- **Go** - Core backend language
- **MySQL** - Primary database
- **Redis** - Caching layer
- **RabbitMQ** - Message queue
- **Docker** - Containerization
- **Amazon ECS** - Container orchestration
- **Amazon Lambda** - Cron job
- **Amazon S3** - File storage
- **API Gateway** - API management
- **GitHub Actions** - CI/CD pipeline

## Features

- RESTful API endpoints for:
  - User authentication and authorization
  - Transaction management (create, read, update, delete)
  - Category management
  - Upload custom icons for categories
  - Financial reporting and analytics
- Comprehensive test coverage
- Clean Architecture implementation
- Scalable infrastructure design
- Secure API endpoints
- Performance optimized with caching

## Architecture

This project follows [Clean Architecture principles](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html), ensuring:

- Independence of layers: Changes in one layer don't affect others, ensuring modular development.

- Testability: Business logic can be tested in isolation without external dependencies.

- Maintainability through separation of concerns: Each component has a single responsibility, simplifying maintenance.

- Database independence: Business logic is decoupled from storage implementation, allowing flexible database choices.

- Easier feature implementation: Clear boundaries enable safe addition of new features without side effects.

- Long-term scalability: Architecture supports growing complexity while maintaining code quality.


## How To Run Locally

There are two ways to run this project locally:

### Prerequisites
You need to install docker to run this project locally.

### Option 1: Only Backend

If you only want to run the backend API without cloning the frontend repository:

1. Clone this repository:
   ```bash
   git clone git@github.com:eyo-chen/expense-tracker-go.git
   cd expense-tracker-go
   ```

2. Copy the environment file:
   ```bash
   cp .env.example .env
   ```

3. Run the application:
   ```bash
   make run
   ```

The API will be available at `http://localhost:8000` and the frontend will be available at `http://localhost:3000`

### Option 2: Full Stack Development (Frontend + Backend)

To run both frontend and backend together, you'll need to set up the following directory structure:

```bash
expense-tracker-app/
├── expense-tracker-go/     # Backend
└── expense-tracker/        # Frontend
```

Follow these steps:

1. Create and enter the main project directory:
   ```bash
   mkdir expense-tracker-app
   cd expense-tracker-app
   ```

2. Clone both repositories:
   ```bash
   # Clone backend
   git clone git@github.com:eyo-chen/expense-tracker-go.git
   
   # Clone frontend
   git clone git@github.com:eyo-chen/expense-tracker.git
   ```

3. Enter the backend directory and set up environment:
   ```bash
   cd expense-tracker-go
   cp .env.example .env
   ```

4. Update the `.env` file with your configurations

5. Run both services:
   ```bash
   # inside expense-tracker-go directory
   make run-with-frontend
   ```

The API will be available at `http://localhost:8000` and the frontend will be available at `http://localhost:3000`

