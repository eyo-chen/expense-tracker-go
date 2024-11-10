# Overview

A production-grade expense tracking backend API built with Go, following [Clean Architecture principles](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html). This service powers the [Expense Tracker frontend application](https://github.com/eyo-chen/expense-tracker).

## Technology Stack

- **Go** - Core backend language
- **MySQL** - Primary database
- **Redis** - Caching layer
- **Docker** - Containerization
- **Amazon ECS** - Container orchestration
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

- Independence of frameworks
- Independence of Database
- Independence of any external services
- Testability


## How To Run Locally

There are two ways to run this project locally:

### Prerequisites
You need to install docker to run this project locally.

### Option 1: Without Cloning Frontend Repository

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

3. Update the `.env` file with your configurations

4. Run the application:
   ```bash
   make run
   ```

The API will be available at `http://localhost:8000` and the frontend will be available at `http://localhost:3000`

### Option 2: Full Stack Development (Cloning Frontend Repository)

To run both frontend and backend together:

1. Clone both repositories:
   ```bash
   # Clone backend
   git clone git@github.com:eyo-chen/expense-tracker-go.git
   cd expense-tracker-go

   # Clone frontend
   git clone git@github.com:eyo-chen/expense-tracker.git ../expense-tracker
   ```

2. Copy the environment file:
   ```bash
   cp .env.example .env
   ```

3. Update the `.env` file with your configurations

4. Run both services:
   ```bash
   make run-with-frontend
   ```

The API will be available at `http://localhost:8000` and the frontend will be available at `http://localhost:3000`

