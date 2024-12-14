# Personal Blog Project

A full-stack web application for creating and managing a personal blog with guest and admin sections. Built with Go backend and React frontend.

## Features

### Guest Section
- View list of published articles on the home page
- Read individual articles with publication dates
- Responsive design for optimal viewing on all devices

### Admin Section
- Secure authentication system
- Dashboard displaying all articles (including drafts)
- Create, edit, and delete articles
- Restore deleted articles
- Basic user management

## Project Submission

This project is a submission for the Personal Blog project on [Roadmap.sh](https://roadmap.sh/projects/personal-blog). It aims to fulfill the requirements of building a personal blog with both guest and admin sections, demonstrating key web development skills including:

- Filesystem-based storage
- Server-side rendering
- Basic authentication
- CRUD operations for articles
- Responsive design

The implementation follows the guidelines provided on the Roadmap.sh project page, focusing on creating a practical, functional personal blog application.

## Tech Stack

### Backend
- Go (Golang)
- Chi router for HTTP routing
- CORS middleware
- Swagger for API documentation
- File-based storage system
- Basic authentication and session management

### Frontend
- React 18
- Vite for build tooling
- React Router for navigation
- Tailwind CSS for styling
- Axios for API communication
- DOMPurify for content sanitization

## Project Structure

```
personal-blog/
├── cmd/
│   ├── api/              # Backend application code
│   └── docs/             # Swagger documentation
├── internal/             # Internal packages
├── web/                 # Frontend React application
│   ├── src/
│   ├── public/
│   └── package.json
└── README.md
```

## Getting Started

### Prerequisites
- Go 1.21 or higher
- Node.js 18 or higher
- npm or yarn

### Backend Setup

1. Navigate to the project root:
```bash
cd personal-blog
```

2. Install Go dependencies:
```bash
go mod tidy
```

3. Run the backend server:
```bash
# Standard run
go run ./cmd/api

# Or with live reload
make run/live
```

The API server will start at `http://localhost:4444`

### Frontend Setup

1. Navigate to the web directory:
```bash
cd web
```

2. Install dependencies:
```bash
npm install
# or
yarn
```

3. Start the development server:
```bash
npm run dev
# or
yarn dev
```

The frontend application will be available at `http://localhost:5173`

## API Documentation

### APIdog Documentation
Comprehensive API documentation is available on APIdog.com:
[Personal Blog API Documentation](https://www.apidog.com/apidoc/shared-a3bc999d-9f57-4565-bc95-af4b2d99c2d3)

The APIdog documentation includes:
- Detailed endpoint descriptions
- Request/response examples
- Authentication requirements
- Error responses
- API schemas

### API Endpoints Summary

#### Public Routes
- `GET /api/v1/articles` - Get all published articles
- `GET /api/v1/article/{id}` - Get specific article by ID

#### Admin Routes (Requires Authentication)
- `POST /api/v1/sign-in` - Admin sign-in
- `POST /api/v1/sign-out` - Admin sign-out
- `GET /api/v1/admin` - Get all articles (including drafts)
- `POST /api/v1/create` - Create new article
- `PATCH /api/v1/edit/{id}` - Update existing article
- `DELETE /api/v1/delete/{id}` - Delete article
- `POST /api/v1/restore/{id}` - Restore deleted article

### Local API Documentation
- Swagger UI is accessible at `/swagger/*`
- Raw API specification available at `/swagger/doc.json`

## Development Tools

### Backend
- `make tidy` - Format code and tidy go.mod
- `make audit` - Run security and code quality checks
- `make test` - Run all tests
- `make build` - Build the application
- `make run/live` - Run with live reload

### Frontend
- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run check` - Run ESLint checks
- `npm run preview` - Preview production build


## License

This project is licensed under the MIT License - see the LICENSE file for details.
