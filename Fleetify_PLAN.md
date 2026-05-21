# Fleetify - Fleet Maintenance Project Plan

## 1. Project Summary

Fleetify is an internal fleet maintenance system for managing vehicle maintenance reports. The application must feel like a single-page app, with a REST API backend and a Vanilla JavaScript frontend.

**Required stack:**

- Backend: Go, GoFiber, GORM
- Database: MySQL 8.0 with InnoDB tables
- Frontend: Vanilla JavaScript and Bootstrap 5
- DevOps: Docker Compose
- Deadline target: 5 working days

The main workflow involves two roles:

1. **Service Advisor (SA)** creates an initial maintenance report, adds estimated parts/services, and later completes approved reports.
2. **Approval / Management** reviews pending reports and approves them.

---

## 2. Main Objectives

The implementation should prioritize:

1. Correct role-based workflow.
2. Reliable database integrity using foreign keys and transactions.
3. Clean Go code structure with repository/service separation.
4. Safe frontend rendering without `.innerHTML`.
5. Easy setup using `docker-compose up --build`.
6. Clear documentation and at least 5 Git commits.

---

## 3. Functional Scope

### 3.1 Required Features

#### F-01 - Create Initial Maintenance Report

**Role:** SA

SA can create a maintenance report by selecting a vehicle from master data and entering:

- Vehicle
- Odometer
- Complaint
- Initial photo or simulated photo
- Estimated parts/services list
- Quantity for each item

System behavior:

- Initial status must automatically become `PENDING_APPROVAL`.
- The item price must be copied from `master_items.price` into `report_items.price_snapshot`.
- Header and detail insertions must run inside one atomic transaction.
- If one detail insert fails, the report header must not be saved.

#### F-02 - Approve Maintenance Report

**Role:** Approval

Approval users can review incoming reports and approve reports with status `PENDING_APPROVAL`.

System behavior:

- Status changes from `PENDING_APPROVAL` to `APPROVED`.
- Invalid transitions must be rejected.
- Optional bonus: trigger webhook asynchronously after approval.

#### F-03 - Complete Approved Report

**Role:** SA

SA can complete approved reports by uploading or simulating a proof photo.

System behavior:

- Only reports with status `APPROVED` can be completed.
- Status changes to `COMPLETED`.
- `proof_photo` is saved.
- Optional bonus: trigger webhook asynchronously after completion.

#### F-04 - Maintenance Report History

**Role:** SA and Approval

Display all maintenance reports with complete key information:

- Report ID
- SA name
- Vehicle license plate
- Vehicle model
- Odometer
- Complaint
- Status
- Initial photo
- Proof photo
- Report items
- Total estimated cost
- Created date

---

## 4. Bonus Features

### B-01 - Export CSV with Native JavaScript

Add a button on the report history page to export reports into CSV.

Rules:

- Use only Native JavaScript.
- Do not use third-party export libraries.
- Generate CSV using `Blob`, `URL.createObjectURL()`, and a temporary `<a>` download element.

### B-02 - Backend Webhook with Goroutine

Trigger an asynchronous HTTP POST request when a report status changes to:

- `APPROVED`
- `COMPLETED`

Implementation notes:

- Use Goroutine after the database transaction/update succeeds.
- Do not block the main API response.
- Read target URL from environment variable, for example `WEBHOOK_URL`.
- Add timeout to the HTTP client.
- Log webhook failures without failing the main status update.

Example webhook payload:

```json
{
  "event": "REPORT_APPROVED",
  "report_id": 1,
  "status": "APPROVED",
  "vehicle_license_plate": "B 1234 FTY",
  "updated_at": "2026-05-20T10:00:00Z"
}
```

---

## 5. Database Design

Use MySQL 8.0 and InnoDB engine for all tables.

### 5.1 Tables

#### users

Stores application users and roles.

| Column | Type | Notes |
|---|---|---|
| id | BIGINT UNSIGNED PK | Auto increment |
| username | VARCHAR(100) | Unique |
| role | ENUM('SA', 'APPROVAL') | Required |
| created_at | DATETIME | Optional GORM timestamp |
| updated_at | DATETIME | Optional GORM timestamp |

#### vehicles

Stores vehicle master data.

| Column | Type | Notes |
|---|---|---|
| id | BIGINT UNSIGNED PK | Auto increment |
| license_plate | VARCHAR(30) | Unique, indexed |
| model | VARCHAR(100) | Required |
| created_at | DATETIME | Optional GORM timestamp |
| updated_at | DATETIME | Optional GORM timestamp |

#### master_items

Stores parts and services master data.

| Column | Type | Notes |
|---|---|---|
| id | BIGINT UNSIGNED PK | Auto increment |
| item_name | VARCHAR(150) | Required |
| type | ENUM('PART', 'SERVICE') | Required |
| price | DECIMAL(15,2) | Required |
| created_at | DATETIME | Optional GORM timestamp |
| updated_at | DATETIME | Optional GORM timestamp |

#### maintenance_reports

Stores maintenance report header data.

| Column | Type | Notes |
|---|---|---|
| id | BIGINT UNSIGNED PK | Auto increment |
| vehicle_id | BIGINT UNSIGNED FK | References vehicles.id |
| created_by | BIGINT UNSIGNED FK | References users.id |
| odometer | INT UNSIGNED | Required |
| complaint | TEXT | Required |
| status | ENUM('PENDING_APPROVAL', 'APPROVED', 'COMPLETED') | Required |
| initial_photo | VARCHAR(255) | Initial photo path or simulated filename |
| proof_photo | VARCHAR(255) | Proof photo path or simulated filename |
| created_at | DATETIME | Required |
| updated_at | DATETIME | Optional but recommended |

#### report_items

Stores maintenance report detail items.

| Column | Type | Notes |
|---|---|---|
| id | BIGINT UNSIGNED PK | Auto increment |
| report_id | BIGINT UNSIGNED FK | References maintenance_reports.id |
| item_id | BIGINT UNSIGNED FK | References master_items.id |
| quantity | INT UNSIGNED | Required, minimum 1 |
| price_snapshot | DECIMAL(15,2) | Copied from master_items.price |
| created_at | DATETIME | Optional GORM timestamp |
| updated_at | DATETIME | Optional GORM timestamp |

### 5.2 Relationship Summary

- One `user` can create many `maintenance_reports`.
- One `vehicle` can have many `maintenance_reports`.
- One `maintenance_report` has many `report_items`.
- One `master_item` can be used in many `report_items`.

### 5.3 Data Integrity Rules

- Use foreign keys for all relationships.
- Use indexes on frequently queried columns:
  - `maintenance_reports.status`
  - `maintenance_reports.created_by`
  - `maintenance_reports.vehicle_id`
  - `vehicles.license_plate`
- Use transaction for report creation.
- Use status validation to prevent invalid workflow transitions.

---

## 6. Seeder Plan

Seeder must run automatically through Docker init SQL or application startup.

Minimum seed data:

### Users

| Username | Role |
|---|---|
| sa_user | SA |
| approval_user | APPROVAL |

### Vehicles

| License Plate | Model |
|---|---|
| B 1234 FTY | Toyota Avanza |
| B 5678 FTY | Daihatsu Gran Max |
| B 9012 FTY | Mitsubishi L300 |

### Master Items

| Item Name | Type | Price |
|---|---|---:|
| Engine Oil | PART | 350000 |
| Oil Filter | PART | 85000 |
| Brake Pad | PART | 450000 |
| General Service | SERVICE | 250000 |
| Brake Inspection | SERVICE | 150000 |

Recommended approach:

- Use `init.sql` under `docker/mysql/init/` for predictable initial data.
- Keep `schema.sql` available in repository for documentation and manual review.
- Make seeder idempotent using `INSERT IGNORE` or `ON DUPLICATE KEY UPDATE`.

---

## 7. Backend Architecture Plan

### 7.1 Recommended Folder Structure

```text
backend/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── database/
│   │   └── mysql.go
│   ├── models/
│   │   ├── user.go
│   │   ├── vehicle.go
│   │   ├── master_item.go
│   │   ├── maintenance_report.go
│   │   └── report_item.go
│   ├── repositories/
│   │   ├── user_repository.go
│   │   ├── vehicle_repository.go
│   │   ├── item_repository.go
│   │   └── report_repository.go
│   ├── services/
│   │   ├── report_service.go
│   │   └── webhook_service.go
│   ├── handlers/
│   │   ├── vehicle_handler.go
│   │   ├── item_handler.go
│   │   └── report_handler.go
│   ├── middlewares/
│   │   ├── user_context.go
│   │   └── rbac.go
│   └── routes/
│       └── routes.go
├── uploads/
├── Dockerfile
├── go.mod
└── go.sum
```

### 7.2 Backend Layers

#### Handler Layer

Responsibilities:

- Parse request body/form-data.
- Read authenticated user from Fiber context.
- Return consistent JSON responses.
- Avoid business logic inside handlers.

#### Service Layer

Responsibilities:

- Validate workflow status.
- Run transaction for report creation.
- Snapshot item prices.
- Trigger webhook after successful status changes.

#### Repository Layer

Responsibilities:

- Encapsulate GORM queries.
- Handle preloading relationships.
- Keep database access clean and reusable.

---

## 8. Backend API Plan

### 8.1 Header-Based User Context

All protected requests must include:

```http
X-User-ID: 1
```

Middleware flow:

1. Read `X-User-ID`.
2. Validate that it exists and is numeric.
3. Find user in database.
4. Store user object in Fiber context.
5. RBAC middleware checks required role.

### 8.2 API Endpoints

#### Health Check

```http
GET /health
```

Response:

```json
{
  "status": "ok"
}
```

#### Get Testing Users

```http
GET /api/users
```

Purpose:

- Used by frontend user switcher for testing SA and Approval workflows.

#### Get Vehicles

```http
GET /api/vehicles
```

Role:

- SA
- Approval

#### Get Master Items

```http
GET /api/master-items
```

Role:

- SA
- Approval

#### Create Maintenance Report

```http
POST /api/reports
```

Role:

- SA only

Request example:

```json
{
  "vehicle_id": 1,
  "odometer": 120000,
  "complaint": "Brake noise and engine vibration",
  "initial_photo": "initial-photo.jpg",
  "items": [
    {
      "item_id": 1,
      "quantity": 1
    },
    {
      "item_id": 4,
      "quantity": 1
    }
  ]
}
```

Business rules:

- `status` must be set by backend as `PENDING_APPROVAL`.
- Do not trust status from frontend.
- Validate vehicle exists.
- Validate all selected master items exist.
- Snapshot price from `master_items.price`.
- Use one database transaction for header and detail.

#### Get All Reports

```http
GET /api/reports
```

Role:

- SA
- Approval

Query options:

```http
GET /api/reports?status=PENDING_APPROVAL
```

Response should include:

- Report header
- User/SA data
- Vehicle data
- Report items with master item data
- Computed total estimate

#### Get Report Detail

```http
GET /api/reports/:id
```

Role:

- SA
- Approval

#### Approve Report

```http
PATCH /api/reports/:id/approve
```

Role:

- Approval only

Business rules:

- Only `PENDING_APPROVAL` reports can be approved.
- Change status to `APPROVED`.
- Trigger webhook asynchronously if enabled.

#### Complete Report

```http
PATCH /api/reports/:id/complete
```

Role:

- SA only

Request example:

```json
{
  "proof_photo": "proof-photo.jpg"
}
```

Business rules:

- Only `APPROVED` reports can be completed.
- Save proof photo.
- Change status to `COMPLETED`.
- Trigger webhook asynchronously if enabled.

---

## 9. Critical Backend Logic

### 9.1 Report Creation Transaction

Pseudo-flow:

```go
err := db.Transaction(func(tx *gorm.DB) error {
    // 1. Validate vehicle exists
    // 2. Create maintenance_reports header with status PENDING_APPROVAL
    // 3. Loop through request items
    // 4. Load each master_item from DB
    // 5. Create report_items using master_item.Price as price_snapshot
    // 6. Return nil to commit transaction
})
```

Important points:

- Never use price sent from frontend.
- Always read price from `master_items` inside backend logic.
- Return error from transaction callback to rollback.
- Keep transaction only for database operations, not webhook calls.

### 9.2 Status Transition Validation

Allowed transitions:

| Current Status | Action | New Status | Role |
|---|---|---|---|
| PENDING_APPROVAL | Approve | APPROVED | APPROVAL |
| APPROVED | Complete | COMPLETED | SA |

Rejected transitions:

- `PENDING_APPROVAL` directly to `COMPLETED`.
- `APPROVED` back to `PENDING_APPROVAL`.
- `COMPLETED` to any other status.
- Approval by SA.
- Completion by Approval.

### 9.3 Error Response Format

Use a consistent JSON structure:

```json
{
  "success": false,
  "message": "Only APPROVAL role can approve reports",
  "error": "FORBIDDEN"
}
```

Success response example:

```json
{
  "success": true,
  "message": "Report approved successfully",
  "data": {}
}
```

---

## 10. Frontend Architecture Plan

### 10.1 Frontend Folder Structure

```text
frontend/
├── index.html
├── assets/
│   ├── css/
│   │   └── style.css
│   └── js/
│       ├── api.js
│       ├── state.js
│       ├── dom.js
│       ├── components.js
│       ├── pages/
│       │   ├── create-report.js
│       │   ├── approval.js
│       │   └── history.js
│       └── app.js
```

### 10.2 Frontend Pages / Sections

The frontend should behave like a single-page app using section switching.

Recommended sections:

1. **User Switcher**
   - Select seeded user.
   - Store selected user ID in frontend state/localStorage.
   - Send selected user ID through `X-User-ID` header for every API call.

2. **Create Report**
   - Visible for SA.
   - Vehicle dropdown.
   - Odometer input.
   - Complaint textarea.
   - Initial photo field.
   - Dynamic item rows.
   - Submit button.

3. **Approval Queue**
   - Visible for Approval.
   - Show reports with status `PENDING_APPROVAL`.
   - Report detail preview.
   - Approve button.

4. **Complete Report**
   - Visible for SA.
   - Show approved reports.
   - Proof photo input.
   - Complete button.

5. **Report History**
   - Visible for SA and Approval.
   - Table of all reports.
   - Status badge.
   - Detail button/modal.
   - Optional CSV export button.

### 10.3 DOM Rendering Rules

The document explicitly prohibits using `.innerHTML` to render data.

Use:

- `document.createElement()`
- `document.createTextNode()`
- `DocumentFragment`
- `element.appendChild()`
- `element.replaceChildren()`

Example pattern:

```js
const fragment = document.createDocumentFragment();

reports.forEach((report) => {
  const row = document.createElement("tr");

  const statusCell = document.createElement("td");
  statusCell.textContent = report.status;

  row.appendChild(statusCell);
  fragment.appendChild(row);
});

tableBody.replaceChildren(fragment);
```

### 10.4 Bootstrap 5 Usage

Use Bootstrap for:

- Navbar or tab navigation
- Responsive forms
- Tables
- Cards
- Badges for status
- Modal for report detail
- Alert messages
- Loading buttons/spinners

Status badge suggestion:

| Status | Bootstrap Badge |
|---|---|
| PENDING_APPROVAL | `bg-warning text-dark` |
| APPROVED | `bg-primary` |
| COMPLETED | `bg-success` |

---

## 11. Docker Plan

### 11.1 Required Services

`docker-compose.yml` should include:

1. `app`
   - Go application container
   - Exposes backend API port, for example `8080`
   - Depends on MySQL

2. `mysql`
   - MySQL 8.0
   - Uses persistent volume
   - Runs init SQL/seeder automatically

Optional:

3. `phpmyadmin` or `adminer`
   - Helpful for checking database during development
   - Keep optional to avoid unnecessary complexity

### 11.2 Expected Command

The whole application must run with:

```bash
docker-compose up --build
```

### 11.3 Environment Variables

Recommended `.env.example`:

```env
APP_PORT=8080
DB_HOST=mysql
DB_PORT=3306
DB_USER=fleetify
DB_PASSWORD=fleetify_password
DB_NAME=fleetify_db
DB_ROOT_PASSWORD=root_password
WEBHOOK_URL=https://webhook.site/example
UPLOAD_DIR=uploads
```

### 11.4 Startup Reliability

Because MySQL may not be ready instantly, the Go application should:

- Retry DB connection several times, or
- Use a wait script, or
- Implement simple retry logic in `database/mysql.go`.

---

## 12. Development Timeline - 5 Working Days

### Day 1 - Project Setup and Database Foundation

Goals:

- Initialize Git repository.
- Create GoFiber project structure.
- Create Dockerfile and Docker Compose.
- Configure MySQL connection.
- Create models and migrations/schema.
- Add seed data.

Deliverables:

- Backend starts successfully.
- MySQL starts successfully.
- Seed data available.
- `/health` endpoint works.

Suggested commit:

```text
chore: initialize project structure with docker and database setup
```

### Day 2 - Backend Core Workflow

Goals:

- Implement user context middleware using `X-User-ID`.
- Implement RBAC middleware.
- Implement vehicles and master items endpoints.
- Implement create report endpoint with transaction.
- Implement report listing with preloaded relations.

Deliverables:

- SA can create report.
- Report item prices are snapshotted.
- Report list includes SA, vehicle, and items.

Suggested commit:

```text
feat: implement report creation with transaction and price snapshot
```

### Day 3 - Approval and Completion Workflow

Goals:

- Implement approve endpoint.
- Implement complete endpoint.
- Add status transition validation.
- Add consistent error handling.
- Add optional file/photo simulation handling.

Deliverables:

- Approval can approve pending reports.
- SA can complete approved reports.
- Invalid role/action combinations are rejected.

Suggested commit:

```text
feat: implement approval and completion workflow with role validation
```

### Day 4 - Frontend Implementation

Goals:

- Build Bootstrap layout.
- Build user switcher.
- Build create report form.
- Build approval queue.
- Build completion section.
- Build report history table.
- Integrate Fetch API with `X-User-ID` header.
- Render data using `createElement` and `DocumentFragment` only.

Deliverables:

- Full workflow can be tested from browser.
- UI is responsive.
- No `.innerHTML` is used to render API data.

Suggested commit:

```text
feat: build vanilla js frontend for maintenance workflow
```

### Day 5 - Bonus, Testing, Documentation, Cleanup

Goals:

- Add CSV export using Native JS.
- Add webhook using Goroutine if time allows.
- Improve README.
- Add API documentation.
- Add testing accounts to README.
- Perform full manual testing.
- Ensure repository has at least 5 commits.

Deliverables:

- `README.md` is complete.
- Docker setup works from clean environment.
- Bonus features are implemented if possible.
- Final repository is ready to submit.

Suggested commits:

```text
feat: add native csv export for report history
feat: add async webhook notification for status changes
```

```text
docs: add setup guide api documentation and testing accounts
```

---

## 13. Manual Testing Scenarios

### 13.1 Seeder Validation

- Run `docker-compose up --build`.
- Confirm database is created.
- Confirm 2 users exist.
- Confirm 3 vehicles exist.
- Confirm 5 master items exist.

### 13.2 SA Creates Report

Steps:

1. Select `sa_user`.
2. Open Create Report section.
3. Select a vehicle.
4. Fill odometer and complaint.
5. Add at least 2 items.
6. Submit report.

Expected result:

- Report is created.
- Status is `PENDING_APPROVAL`.
- Report items contain `price_snapshot`.

### 13.3 SA Cannot Approve Report

Steps:

1. Send approve request using `X-User-ID` of SA.

Expected result:

- API returns forbidden error.
- Report status remains `PENDING_APPROVAL`.

### 13.4 Approval Approves Report

Steps:

1. Select `approval_user`.
2. Open Approval Queue.
3. Click Approve.

Expected result:

- Report status changes to `APPROVED`.
- Optional webhook is triggered.

### 13.5 Approval Cannot Complete Report

Steps:

1. Send complete request using `X-User-ID` of Approval.

Expected result:

- API returns forbidden error.
- Report status remains `APPROVED`.

### 13.6 SA Completes Report

Steps:

1. Select `sa_user`.
2. Open Complete Report section.
3. Add proof photo.
4. Submit completion.

Expected result:

- Report status changes to `COMPLETED`.
- `proof_photo` is saved.
- Optional webhook is triggered.

### 13.7 CSV Export

Steps:

1. Open Report History.
2. Click Export CSV.

Expected result:

- Browser downloads CSV file.
- CSV contains report rows.

---

## 14. README.md Content Plan

The repository README should include the following sections:

1. **Project Overview**
   - Short explanation of Fleetify and its workflow.

2. **Tech Stack**
   - Go, GoFiber, GORM, MySQL 8.0, Vanilla JS, Bootstrap 5, Docker.

3. **How to Run**

```bash
git clone <repository-url>
cd <repository-folder>
cp .env.example .env
docker-compose up --build
```

4. **Seeder Information**
   - Explain how seed data runs.
   - Mention seeded users, vehicles, and master items.

5. **Environment Variables**
   - Explain each variable in `.env.example`.

6. **Testing Accounts**

| Username | Role | X-User-ID |
|---|---|---:|
| sa_user | SA | 1 |
| approval_user | APPROVAL | 2 |

7. **API Documentation**
   - List all endpoints, methods, roles, and example request bodies.

8. **Workflow Explanation**
   - `PENDING_APPROVAL -> APPROVED -> COMPLETED`.

9. **Technical Decisions**
   - Why transaction is used.
   - Why price snapshot is used.
   - Why RBAC uses `X-User-ID`.
   - Why frontend avoids `.innerHTML`.

10. **Bonus Features**
   - CSV Export.
   - Webhook.

---

## 15. Git Commit Strategy

The document requires at least 5 commit/push histories. Recommended commits:

```text
1. chore: initialize fleetify project with docker compose
2. feat: add database schema and seed data
3. feat: implement backend report creation with transaction
4. feat: implement approval and completion workflow
5. feat: build vanilla js bootstrap frontend
6. feat: add csv export and async webhook bonus
7. docs: add complete readme and api documentation
```

Keep commits meaningful and avoid pushing everything in one final commit.

---

## 16. Final Submission Checklist

### Repository

- [ ] Repository is public.
- [ ] Repository has at least 5 commits.
- [ ] `README.md` is complete.
- [ ] `.env.example` is included.
- [ ] `docker-compose.yml` is included.
- [ ] Seeder is included and runs automatically.
- [ ] `schema.sql`, `init.sql`, or Go seeder is included.

### Backend

- [ ] Uses GoFiber.
- [ ] Uses GORM.
- [ ] Connects to MySQL 8.0.
- [ ] Uses InnoDB tables.
- [ ] Has role validation using `X-User-ID`.
- [ ] Report creation uses database transaction.
- [ ] Report item price is snapshotted from master item price.
- [ ] Status transitions are validated.
- [ ] Error responses are consistent.
- [ ] Webhook bonus is implemented if possible.

### Frontend

- [ ] Uses Vanilla JavaScript.
- [ ] Uses Bootstrap 5.
- [ ] Uses Fetch API.
- [ ] Does not use `.innerHTML` to render data.
- [ ] Uses `document.createElement()` or `DocumentFragment`.
- [ ] SA can create reports.
- [ ] Approval can approve reports.
- [ ] SA can complete approved reports.
- [ ] Report history is displayed.
- [ ] CSV export bonus is implemented if possible.

### Docker

- [ ] App and MySQL run through Docker Compose.
- [ ] Command `docker-compose up --build` works from clean setup.
- [ ] MySQL data is seeded.
- [ ] App can connect to database reliably.

---

## 17. Recommended Implementation Priority

If time is limited, prioritize in this order:

1. Docker and database setup.
2. Seeder.
3. Backend models and relationships.
4. RBAC middleware using `X-User-ID`.
5. Create report with transaction and price snapshot.
6. Approve report.
7. Complete report.
8. Report history API.
9. Frontend workflow.
10. README and commit history.
11. CSV export bonus.
12. Webhook bonus.

The core scoring is mainly affected by backend workflow correctness, Go code quality, database integrity, frontend integration, Docker usability, and documentation. Bonus features should be added only after the required workflow is stable.
