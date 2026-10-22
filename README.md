# Omarix

**Omarix** is a full-stack image-processing web application built with Go and React. It allows users to upload images, apply transformations, preview their changes, and save edited copies to a personal gallery.

The project focuses on integrating a REST API, authenticated image management, Cloudinary-based transformations, and an interactive editing experience.

## Features

### Image editing

- Upload JPG and PNG images (the backend currently enforces a 5 MB upload limit).
- Resize, crop, rotate, and flip images.
- Adjust brightness, contrast, and saturation.
- Apply grayscale, sepia, blur, pixelation, and cartoon-style effects.
- Access Cloudinary AI transformation options such as background removal, generative restoration, and upscaling (subject to Cloudinary account support).
- Preview transformations using Cloudinary-rendered images rather than CSS approximations.
- Compare an edited preview against the original.
- Save transformed images as separate gallery entries, preserving the original.
- Download completed images.

### Accounts and gallery

- Register and sign in with an application account, or use Google Sign-In through Firebase Authentication.
- Authenticate protected API requests using JWT bearer tokens.
- Access a personal gallery with paginated image listings.
- Open, edit, and delete images belonging to the signed-in user.
- Navigate directly to the completed image after saving an edit.
- Use a responsive interface with light and dark themes.

## Tech stack

| Layer | Technology |
| --- | --- |
| Frontend | React, Vite, JavaScript |
| Backend | Go, Gin |
| Database | PostgreSQL, GORM |
| Image storage and transformations | Cloudinary |
| Google Sign-In | Firebase Authentication |
| API authentication | JWT |

## How it works

1. The user signs in and receives an application access token.
2. The frontend sends an image upload request to the Go API.
3. The backend validates the upload, stores the image in Cloudinary, and saves its metadata in PostgreSQL.
4. The workspace builds transformation parameters from the user's editing choices.
5. The frontend requests Cloudinary-rendered previews to show the effect of those parameters.
6. When the user saves, the backend checks image ownership, validates the transformation request, processes the image through Cloudinary, and creates a new image record.
7. The user can view the result and access it later in their gallery.

The original image and the transformed copy are stored as separate records.

## Architecture

```text
React + Vite frontend
    |
    | REST requests + JWT bearer token
    v
Go / Gin API
    |-- Authentication and authorization middleware
    |-- Controllers (HTTP handlers)
    |-- Services (application logic)
    |-- Repositories (database access)
    |       |
    |       v
    |   PostgreSQL / GORM
    |
    |-- Firebase Admin SDK (Google ID-token verification)
    |
    `-- Cloudinary (image upload and transformations)
```

The frontend uses Firebase Authentication for Google Sign-In. The Go backend verifies the Firebase ID token before issuing its own application JWT. Image endpoints use the application JWT and check that the requested image belongs to the authenticated user.

## API overview

The image routes are registered under `/api` and require a valid bearer token.

| Method | Endpoint | Purpose |
| --- | --- | --- |
| `POST` | `/api/images` | Upload an image (`multipart/form-data`, field `image`) |
| `GET` | `/api/images` | List the user's images (`page`, `limit`) |
| `GET` | `/api/images/:id` | Retrieve an owned image |
| `POST` | `/api/images/:id/transform` | Apply transformations and save a new image |
| `DELETE` | `/api/images/:id` | Delete an owned image |

### Example transformation request

```json
{
  "effects": {
    "brightness": -15,
    "contrast": 25,
    "grayscale": true
  },
  "geometry": {
    "flip": "horizontal"
  }
}
```

Example request:

```bash
curl -X POST http://localhost:8080/api/images/1/transform \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"effects":{"brightness":-15,"contrast":25,"grayscale":true}}'
```

The exact result depends on Cloudinary transformation support and configuration.

## Local development

### Requirements

- Go (version compatible with the repository's `go.mod`)
- Node.js and npm (versions compatible with the frontend's `package.json`)
- PostgreSQL
- Cloudinary account
- Firebase project configured for Google Sign-In

### Frontend configuration

In the frontend directory, create a `.env` file with the settings used by the Vite application:

```dotenv
VITE_API_BASE_URL=http://localhost:8080
VITE_FIREBASE_API_KEY=your_firebase_web_api_key
VITE_FIREBASE_AUTH_DOMAIN=your-project.firebaseapp.com
VITE_FIREBASE_PROJECT_ID=your-project-id
VITE_FIREBASE_APP_ID=your-firebase-web-app-id
```

Firebase Web App configuration values are intended for client-side use. **Never put Firebase service-account credentials, JWT signing secrets, database passwords, or Cloudinary API secrets in `VITE_` variables.**

Install and run the frontend from its directory:

```bash
npm install
npm run dev
```

### Backend configuration

The Go service requires a PostgreSQL connection, JWT signing configuration, a Cloudinary URL, and Firebase Admin SDK credentials. Set these according to the repository's existing configuration loader; the exact environment variable names and database setup depend on that code.

The backend initializes Firebase Admin SDK using a service-account credentials file (or the configured credentials path). Keep that file out of version control.

Run the backend from the Go module directory after configuring its dependencies:

```bash
go mod download
go run .
```

The frontend development server is typically available at `http://localhost:5173`. The Go API address is controlled by the backend configuration; the frontend examples above assume `http://localhost:8080`.

## Engineering challenges

### Consistent previews

An early version used CSS filters to approximate image effects in the browser. The saved results, however, were generated by Cloudinary, so users could see differences between the preview and final output. The workspace was updated to use Cloudinary-rendered previews and align transformation parameters across the frontend and backend.

### Cropping and editing workflows

The crop interface was iterated to support selecting a region of the image and saving a crop without requiring another effect. The result flow was also improved so users are taken directly to the saved image instead of having to search through gallery pages.

### Access control

Protected image routes validate JWTs and scope image queries to the authenticated user. The backend also validates transformation inputs before sending them to Cloudinary.

## Known limitations and future improvements

- Introduce secure refresh-token/session renewal so users do not need to sign in again when short-lived access tokens expire.
- Expand automated tests for transformations, ownership checks, and frontend workflows.
- Optimize Cloudinary preview responsiveness and error handling under slow networks.
- Verify all advanced AI transformations against the Cloudinary plan and configuration used for deployment.
- Improve observability, production configuration, and deployment documentation.

## Repository

[GitHub — image_processing_service](https://github.com/omarelmadboly54-sketch/image_processing_service)

## Author

Developed as a personal full-stack project to practice Go backend development, React integration, authentication, image processing, and user-focused debugging.
