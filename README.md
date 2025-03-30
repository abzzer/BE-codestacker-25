# 🕵️‍♂️ BE-codestacker-25

> Crime Case Management System – Backend API

---

## 📦 Pre-requisites

Ensure the following are installed:

- Docker + Docker Compose
- Make (for `make docker-up`)
- Postman (for testing endpoints manually)
- Available local ports:
  - `5432` (PostgreSQL)
  - `9000`, `9001` (MinIO)
  - `8080` (Application)

**Note on Volumes:**  
PostgreSQL data is currently _not persisted_ due to disabled volumes in `docker-compose.yml`. To enable persistence, uncomment the following:

```yaml
volumes:
  - pg-data:/var/lib/postgresql/data
```

---

## 🏁 Starting the API

```bash
make docker-up     # Start the full system
make docker-down   # Stop and clean containers
```

Access the API at:  
`http://localhost:8080`

When you visit the root endpoint `/`, you should receive the following confirmation:

```json
"We have a working API with databases and RBAC!!"
```

---

## 🛠️ Database Schema & System Structure

![DB Schema](/utils/Schema.png)


---

## 📚 Table of Contents

1. [User Management API](#1-user-management-api)
2. [Case Management APIs](#2-case-management-apis)
3. [Case Listing API](#3-case-listing-api)
4. [Case Details API](#4-case-details-api)
5. [Additional Case APIs](#5-additional-case-apis)
6. [Evidence Management APIs](#6-evidence-management-apis)
7. [Evidence Retrieval API](#7-evidence-retrieval-api)
8. [Evidence Image Retrieval API](#8-evidence-image-retrieval-api)
9. [Evidence Update API](#9-evidence-update-api)
10. [Soft Delete API](#10-soft-delete-api)
11. [Hard Delete API](#11-hard-delete-api)
12. [Text Analysis API](#12-text-analysis-api)
13. [Link Extraction API](#13-link-extraction-api)
14. [Audit Log API](#14-audit-log-api)
15. [Generate Report API](#15-generate-report-api)
16. [Report Tracking API (Public)](#16-report-tracking-api-public)
17. [Long Polling for Evidence Hard Delete](#17-long-polling-for-evidence-hard-delete)
18. [Deployment](#18-deployment)
19. [Appendix 1: How to upload photo evidnece with post-man ](#Appendix-1-How-to-Upload-Image-Evidence-with-Postman)

---

## 1. User Management API

> Admin can add, update, and delete users. Also manage user roles and clearance levels.

---

### Login to Get Token

Send a `POST` request to:

```
http://localhost:8080/login
```

#### Body (raw → JSON):

```json
{
  "user_id": "A001",
  "password": "123"
}
```

#### Successful Response:

```json
{
  "role": "admin",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjbGVhcmFuY2VfbGV2ZWwiOiJjcml0aWNhbCIsImV4cCI6MTc0MzM2MDE3OCwicm9sZSI6ImFkbWluIiwidXNlcl9pZCI6IkEwMDEifQ.WOwE6uTo6Vms6HI7SApele0DMc24DXvkkPgEOTvExG0",
  "userID": "A001"
}
```

---

### Using the Token in Requests

1. In Postman, after login, go to your next request (e.g. `/admin/create-user`)
2. Click the **Headers** tab
3. Add the following header:

```
Key: Authorization
Value: Bearer <your_token>
```

> Replace `<your_token>` with the actual token from your login response.

---

### Create New User (Admin Only)

Send a `POST` request to:

```
http://localhost:8080/admin/create-user
```

#### Headers:

```
Authorization: Bearer <your_token>
```

#### Body (raw → JSON):

```json
{
  "name": "Officer John",
  "password": "secure123",
  "role": "officer",
  "clearance_level": "medium"
}
```

#### Expected Response:

```json
{
  "message": "New user successfully created by admin",
  "created_id": "A103",
  "role": "officer",
  "clearance": "medium"
}
```


---

## 2. Case Management APIs

- Submit public crime report
- Create new case (investigator)
- Update existing case
- Link case to crime reports

More info soon

---

## 3. Case Listing API

> Search and list cases by name or description.

More info soon

---

## 4. Case Details API

> Returns full metadata and stats about a case.

More info soon

---

## 5. Additional Case APIs

- Get all assignees of a case
- Get all evidence of a case
- Get all suspects, victims, and witnesses

More info soon

---

## 6. Evidence Management APIs

> Submit image/text evidence with remarks. Includes upload to MinIO if image.

More info soon

---

## 7. Evidence Retrieval API

> Retrieve evidence by ID, returning content and metadata.

More info soon

---

## 8. Evidence Image Retrieval API

> Stream the image directly from MinIO if content is type "image".

More info soon

---

## 9. Evidence Update API

> Update the content or remarks (not the type).

More info soon

---

## 10. Soft Delete API

> Marks evidence as deleted and logs the action.

More info soon

---

## 11. Hard Delete API

> Multi-step confirmation before full deletion. Deletes image from MinIO if required. Logs the action in audit table.

More info soon

---

## 12. Text Analysis API

> Extracts top 10 used words across all **text-based evidence**, ignoring stopwords.

More info soon

---

## 13. Link Extraction API

> Retrieves all URLs found in evidence content for a given case.

More info soon

---

## 14. Audit Log API

> Admin-only API to view all audit logs of evidence actions.

More info soon

---

## 15. Generate Report API

> Download a **PDF** report for a given case. Includes:
- Case Details
- Evidence (Text + Images)
- Suspects / Victims / Witnesses

More info soon

---

## 16. Report Tracking API (Public)

> Public users can track their crime report status using the report ID they receive upon submission.

_Note: Public Endpoint - no auth needed_

Send a `GET` request to:

```
http://localhost:8080/check-report/1
```

> You can also replace `1` with the actual `report_id` you received after submitting a report.


#### Expected reponse - Default is pending when a crime has been reported but a case has not been opened yet. Hence the design decision to allow for a NULL field in "case_id" in report

```json
{
  "status": "ongoing"
}
```




---

## 17. Long Polling for Evidence Hard Delete

> [Bonus Challenge]

- Initiate long-poll delete
- Poll status: "In Progress", "Completed", or "Failed"
- Enables real-time feedback loop for deletion events

More info soon

---

## 18. Deployment

> This project is fully Dockerised. You can deploy locally via `make docker-up`.

More info soon


---

## Appendix 1: How to Upload Image Evidence with Postman

Use this guide to upload image-based evidence to a case using Postman.

### Step-by-step:

1. Open **Postman** and select `POST`.

2. Use this URL:  
   `http://localhost:8080/add-evidence/image`

3. Click the **Body** tab → Select **form-data**.

4. Add the following fields:

   - `case_number` → Type: Text → Value: `C12345`  
   - `remarks` → Type: Text → Value: `Attached image from scene`  
   - `image` → Type: File → Choose a `.jpg`, `.png`, etc.

5. Ensure you header contains the required authorisation token.
---

