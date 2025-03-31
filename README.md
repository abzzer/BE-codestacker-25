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
16. [Report Tracking API](#16-report-tracking-api)
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
### Logout
> This endpoint simulates logout. It does not invalidate the token server-side.  

**POST** `http://localhost:8080/logout`

**Headers:**

```
Authorization: Bearer <your_token>
```

**Response:**

```json
{
  "message": "Successfully logged out. Please discard your token on the client side."
}
```
---

### Update Existing User

**PATCH** `http://localhost:8080/admin/update-user/A104`

Update user details such as name, password, role, or clearance level.

**Headers:**

```
Authorization: Bearer <admin_token>
```

#### Payload (raw → JSON):

```json
{
  "password": "newpass123",
  "role": "investigator",
  "clearance_level": "medium"
}
```

#### Response:

```json
{
  "message": "User A104 successfully updated"
}
```

> **Note:** If the updated user is currently logged in, they must **log in again** to refresh their token and apply changes.

---

### Delete a User

**DELETE** `http://localhost:8080/admin/delete-user/A104`

Soft deletes a user (marks them as deleted). Cannot be undone.

**Headers:**

```
Authorization: Bearer <admin_token>
```

#### Response:

```json
{
  "message": "User A104 deleted successfully"
}
```
> Note this wont work with A0001 - hard coded the Original ADMIN user to stay forever


---

## 2. Case Management APIs

- Submit public crime report
- Create new case (investigator)
- Update existing case
- Link case to crime reports

--- 

### Submit Crime Report (Public)

**POST** `http://localhost:8080/submit-report`

No authentication required.

**Body (raw → JSON):**

```json
{
  "email": "jane.doe@example.com",
  "civil_id": "1234567890",
  "name": "Jane Doe",
  "description": "Suspicious activity noticed near the alley.",
  "area": "Downtown",
  "city": "Muscat"
}
```

**Response:**

```json
{
  "message": "Report submitted successfully. Please keep your report ID to check status.",
  "report_id": 2
}
```

---

### Create a New Case (Investigator/Admin)

**POST** `http://localhost:8080/add-case`

Requires a valid token from a user with role `investigator` or `admin`.

**Headers:**

```
Authorization: Bearer <your_token>
```

**Body (raw → JSON):**

```json
{
  "case_name": "Alley Incident Investigation",
  "description": "Initial investigation into the alley disturbance",
  "area": "Downtown",
  "city": "Riyadh",
  "level": "medium"
}
```

**Response:**

```json
{
  "case_number": "C10000",
  "message": "Case was created successfully here is your case number. No one is assigned this case yet."
}
```
---

## Update Case

### ✏️ Update General Case Info

**PUT** `http://localhost:8080/update-case/C12345`

> Roles allowed: `admin`, `investigator`  
> Requires token

**Headers:**

```
Authorization: Bearer <your_token>
```

**Body (raw → JSON):**

```json
{
  "case_name": "C12345 Renamed",
  "description": "Updated case details for testing",
  "area": "Al Khobar",
  "city": "Dammam",
  "level": "high"
}
```

**Expected Response:**

```json
{
  "message": "Case updated successfully"
}
```

---

### 🧍 Add Person (Victim / Suspect / Witness)

**POST** `http://localhost:8080/update-case/C12345/add-person`

> Roles allowed: `admin`, `investigator`  
> Requires token

**Body (raw → JSON):**

```json
{
  "type": "victim",
  "name": "Fatima Al-Hassan",
  "age": 32,
  "gender": "female",
  "role": "Eyewitness"
}
```

**Expected Response:**

```json
{
  "message": "Person added successfully",
  "person_id": 1
}
```

---


### Add Officer or Investigator to Case

**POST** `http://localhost:8080/update-case/C12345/add-officer`

> Roles allowed: `admin`, `investigator`  
> Requires token

**Headers:**

```
Authorization: Bearer <your_token>
```

#### Payload (raw → JSON):

```json
{
  "user_id": "A102"
}
```
> `A102` is an **officer** with `medium` clearance.  
> Case `C12345` has `high` level → Officer cannot be assigned.

**Expected Response:**

```json
{
  "error": "officer's clearance level is insufficient for this case"
}
```

---

#### Example with Investigator (Can Always Be Assigned)

```json
{
  "user_id": "A101"
}
```

> `A101` is an investigator → clearance check skipped.

**Expected Response:**

```json
{
  "message": "Officer assigned successfully"
}
```

### Update Case Status Only

**PATCH** `http://localhost:8080/update-case/C12345/status`

> Roles allowed: `admin`, `investigator`, `officer`  
> Requires token

**Body (raw → JSON):**

```json
{
  "status": "ongoing"
}
```

**Expected Response:**

```json
{
  "message": "Case status updated successfully"
}
```
---


---

## 3. Case Listing API

> Search and list cases by name or description.

More info soon

---

## 4. Case Details API

**URL:**  
```
GET http://localhost:8080/case/:caseid
```

**Authorization Required:**  
- Admins and Investigators: Always allowed  
- Officers: Only allowed if their clearance level is equal to or higher than the case level

**Headers:**  
```
Authorization: Bearer <your_token>
```

**Example Request:**  
```
GET http://localhost:8080/case/partial/C12345
```

**Example Successful Response:**  
```json
{
  "case_number": "C12345",
  "case_name": "Theft Investigation",
  "description": "Investigation of a reported theft at a local store.",
  "area": "Downtown",
  "city": "New York",
  "created_by": "A001",
  "created_at": "2025-03-10T14:30:00Z",
  "case_type": "criminal",
  "level": "high",
  "status": "ongoing",
  "reported_by": 1,
  "num_assignees": 3,
  "num_evidences": 2,
  "num_suspects": 1,
  "num_victims": 3,
  "num_witnesses": 0
}
```

---

---

## 5. Additional Case APIs

> Returns full details of a case including:
- All base case metadata
- All assignees (users assigned to the case)
- All evidence (text and images)
- All involved persons (suspects, victims, witnesses)

**Protected route** — Requires `admin`, `investigator`, or `officer` with **appropriate clearance level**.

### Endpoint:

```
GET /case/full/:caseid
```

### Example:

```
GET http://localhost:8080/case/full/C12345
Authorization: Bearer <your_token>
```

### Response:

```json
{
  "case_number": "C12345",
  "case_name": "Theft Investigation",
  "description": "Investigation of a reported theft at a local store.",
  "area": "Downtown",
  "city": "New York",
  "created_by": "A001",
  "created_at": "2025-03-10T14:30:00Z",
  "case_type": "criminal",
  "level": "high",
  "status": "ongoing",
  "reported_by": 1,
  "num_assignees": 3,
  "num_evidences": 2,
  "num_suspects": 1,
  "num_victims": 3,
  "num_witnesses": 0,
  "assignees": [
    ....
  ],
  "evidence": [
    ...
  ],
  "people": [
    .....
    }
  ]
}
```

---

## 6. Evidence Management APIs

> Officers, Investigators, and Admins can add evidence to a case.

---

### Add Text-Based Evidence

**POST** `http://localhost:8080/add-evidence/text`

Adds plain text evidence to a case.

**Headers:**

```
Authorization: Bearer <your_token>
Content-Type: application/json
```

#### Body (raw → JSON):

```json
{
  "case_number": "C12345",
  "type": "text",
  "content": "Note recovered JASDF ASDF ASDF ASDF ASDF ASD FASDF ASDF ASDF ASDF ASDF  at the scene with suspicious handwriting.",
  "remarks": "Check for handwriting match"
}
```

#### Response:

```json
{
  "evidence_id": 3
  "message": "Text evidence added",
}
```

---

### Add Image-Based Evidence

**POST** `http://localhost:8080/add-evidence/image`

Uploads an image file to a case. Image is stored in MinIO and referenced by URL.

**Headers:**

```
Authorization: Bearer <your_token>
Content-Type: multipart/form-data
```

#### Body (form-data → KEY/VALUE):

| Key          | Type | Value                         |
|--------------|------|-------------------------------|
| `case_number`| Text | `C12345`                      |
| `remarks`    | Text | `Photo of the broken window`  |
| `image`      | File | Select a `.jpg`, `.png`, etc. |

#### Response:

```json
{
  "contentSize": "124512 bytes"
  "evidence_id": 4,
  "message": "Image evidence added successfully",
  "minio_url": "http://minio:9000/evidence-bucket/evidence/abc123.jpg",
}
```
> Also You wont be able to access the link directly even if you change minio -> localhost due to Minio Auth
> Ensure the image is valid. Invalid file formats will be rejected.

---

## 7. Evidence Retrieval API

> Must be authetnicated to access this 
> In the belwo examples Evidence ID `1` is TEXT and Evidence ID `2` is an IMAGE
> Note also that police officers can access all evidence I didn't see anything that opposes this in requiremnts:
> I think a real life implementation would restrict their evidence retrieval access

### Get Evidence Details by ID

**GET** `http://localhost:8080/evidence/details/:evidenceid`

Returns either text content or image metadata based on the type.

#### Example: Text Evidence (ID = 1)

```json
{
  "content": "Note recovered JASDF ASDF ASDF ASDF ASDF ASD FASDF ASDF ASDF ASDF ASDF  at the scene with suspicious handwriting.",
  "remarks": "Check for handwriting match",
  "type": "text"
}
```

#### Example: Image Evidence (ID = 2)

```json
{
  "remarks": "Photo of the broken window",
  "size": "113420 bytes",
  "type": "image"
}
```

---

## 8. Evidence Image Retrieval API

> Note here that Evidence with ID 2 is an image and ID 1 is a text

### View Full Image Evidence (ID = 2)

**GET** `http://localhost:8080/evidence/get-image/2`

Returns raw image file from MinIO storage.

- Will trigger file download or display in postman depending on where you access from

---

### Invalid Image Request (Text ID = 1)

**GET** `http://localhost:8080/evidence/get-image/1`

If the evidence ID does not point to an image, this will return an error:

```json
{
  "error": "evidence is not an image"
}
```

---

## 9. Evidence Update API

> Update the content of existing evidence (either text or image).  
> **The evidence type cannot be changed. Only the `content` or image file can be replaced.**

### Roles: `admin`, `investigator`

---

### Update Text Evidence

**Endpoint:**

```
PUT /evidence/update/:evidenceid
```

**Request Example (Body → raw → JSON):**

```json
{
  "content": "Updated text content with valid information and links like https://example.com."
}
```

**Example:**

```
PUT http://localhost:8080/evidence/update/1
Authorization: Bearer <your_token>
```

**Success Response:**

```json
{
  "message": "Text evidence updated successfully"
}
```

---

### Update Image Evidence

**Endpoint:**

```
PUT /evidence/update/:evidenceid
```

**Required Form Field:**

```
Key: Authorization
Value: Bearer <your_token>
```
**Choose FORM and not RAW json**

- `image` → Upload new image file (e.g. `.jpg`, `.png`)

**Success Response:**

```json
{
  "message": "Image evidence updated successfully"
}
```

---

> Notes
- This endpoint is type-specific - types do not change you will get error if you update TEXT with an image and vice-versa.
- Only the evidence `content` field is updated.

---


## 10. Soft Delete API

> Also writes a log entry in the `audit_logs` table with the action `soft_deleted`.
> valid Roles: `admin`, `investigator`


**Endpoint:**

```
GET /evidence/soft-delete/:evidenceid
```

**Example Request:**

```
GET http://localhost:8080/evidence/soft-delete/1
Authorization: Bearer <your_token>
```

**Response Example:**

```json
{
  "message": "Evidence soft-deleted successfully and audit log written"
}
```

---

**Notes:**

- The evidence remains in the database, but its `deleted` flag is set to `true`.
- Soft-deleted is treated as deleted. Currently no endpoint exists to UNDO this.

---

## 11. Hard Delete API

> Multi-step confirmation before full deletion. Deletes image from MinIO if required. Logs the action in audit table.

More info soon

---

## 12. Text Analysis API

**GET** `http://localhost:8080/evidence/top-ten`

Returns most frequent words across all text evidence in the system (excluding common stop words).

#### ✅ Response:

```json
{
  "top_words": [
    "asdf",
    "jasdf",
    "fasdf",
    "suspicious",
    "note",
    "recovered",
    "asd",
    "scene",
    "handwriting"
  ]
}
```


---

## 13. Link Extraction API

> Retrieves all URLs found in text-based evidence content for a given case.


**URL:**  
```
GET http://localhost:8080/evidence/get-urls/:caseid
```

**Authorization Required:**  
- Admin  
- Investigator  
- Officer (any clearence)

**Headers:**  
```
Authorization: Bearer <your_token>
```

**Example Request:**  
```
GET http://localhost:8080/evidence/get-urls/C12345
```

**Example Pre-requisite (Text Evidence to Insert First):**  
```json
{
  "case_number": "C12345",
  "type": "text",
  "content": "Note recovered JASDF https://instagram.com/?hl=en ASDF ASDF ASDF https://www.youtube.com/ ASDF ASD FASDF ASDF ASDF ASDF ASDF  at the scene with https://www.youtube.com/ suspicious handwriting https://www.youtube.com/.",
  "remarks": "Check for handwriting match"
}
```

**Example Response:**  
```json
{
  "urls": [
    "https://www.instagram.com/?hl=en",
    "https://www.youtube.com/",
    "https://www.youtube.com/",
    "https://www.youtube.com/."
  ]
}
```

---

## 14. Audit Log API - View Audit Logs

> Admin-only API to view all audit logs of evidence actions.


**GET** `http://localhost:8080/admin/audit-log`

Returns a list of all evidence-related admin actions (add, update, delete).

**Headers:**

```
Authorization: Bearer <admin_token>
```

#### Sample Response:

```json
[
  {
    "id": 1,
    "action": "added",
    "evidence_id": 10,
    "user_id": "A001",
    "timestamp": "2024-03-30T14:01:25Z"
  },
  {
    "id": 2,
    "action": "soft_deleted",
    "evidence_id": 15,
    "user_id": "A001",
    "timestamp": "2024-03-30T14:05:00Z"
  }
]
```
> NOTE: If you try this without doing anything to evidence you well get an empty json back

---

## 15. Generate Report API

> Generate and download a detailed PDF report for any case. The report includes:
- Case metadata (name, description, level, status, etc.)
- Case assignees and their roles
- People involved (suspects, victims, witnesses)
- Linked citizen reports
- All text evidence with remarks
- All image evidence (JPEG only) with preview

### Generate Case PDF

**URL:**  
```
GET http://localhost:8080/case/pdf/:caseid
```

**Example Request:**  
```
GET http://localhost:8080/case/pdf/C12345
```

**Authorization Required:**  
- Admins and Investigators only

**Headers:**  
```
Authorization: Bearer <your_token>
```

**Output:**  
- PDF file will be returned directly as a file download or preview.
- Only `.jpg` or `.jpeg` images will be embedded. Other image types will be skipped silently.

---

## 16. Report Tracking API

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

### View ALL reports

> You require access token to get ot this endpoint

**Endpoint:**

```
GET /reports/all
```

**Example Request:**

```
GET http://localhost:8080/reports/all
Authorization: Bearer <your_token>
```

**Response Example:**

```json
[
  {
    "report_id": 1,
    "email": "jane.doe@example.com",
    "civil_id": "1234567890",
    "name": "Jane Doe",
    "role": "Citizen",
    "case_number": "C12345",
    "description": "Suspicious activity noticed near the alley.",
    "area": "Downtown",
    "city": "Muscat"
  },
  ...
]
```

---

### Link Report to Case (Not Clarified in outline but practical)

Use this to associate an existing report with a case.

**Endpoint:**

```
POST /reports/case/:reportID
```

**Example Request:**

```
POST http://localhost:8080/reports/case/1
Authorization: Bearer <your_token>
```

**Body (raw JSON):**

```json
{
  "case_number": "C12345"
}
```

**Response Example:**

```json
{
  "message": "Report successfully linked to case C12345"
}
```

---

## 17. Long Polling for Evidence Hard Delete

This section describes how the system tracks and reports the status of a hard delete operation on a piece of evidence using long polling. This mechanism is useful for admins to monitor the progress of deletions in near real-time.

---

### Endpoint

```
GET /evidence/hard-delete-status/:evidenceid
```

Requires a valid **admin token** in the `Authorization` header.

---

### How It Works

When an admin initiates a delete request using the multi-step hard delete process (`POST`, `PATCH`, then `DELETE`), the system tracks its internal state using an in-memory state manager.

This polling endpoint checks the current state of the delete request every second for up to 30 seconds. 

Once the deletion is complete (or fails), the response returns immediately. If the process hasn’t resolved within 30 seconds, it returns the last known state and exits.

---

### Deletion States

Each delete request progresses through the following possible statuses:

#### `initiated`

This status is set when the admin sends a `POST` request to:

```
POST /evidence/hard-delete/:evidenceid
```

It signals the start of the deletion process, prompting the admin to confirm their action.

#### `confirmed`

This status is set when the admin sends a `PATCH` request with a body:

```json
{
  "confirm": "yes"
}
```

It confirms that the admin does intend to delete the evidence.

#### `deleting`

This status is set internally once the actual deletion process begins (usually immediately after the `DELETE` request is received).

#### `done`

This is the final success status. Delete from DB OR minio and Audited

#### `failed`

There was an issue with DB, Minio, or Auditing.

---

### Example Response When Completed

```json
{
  "status": "done",
  "message": "Deletion status resolved"
}
```

---

### Example Response When Timeout is Reached

```json
{
  "status": "confirmed",
  "message": "Timeout reached, no final status yet"
}
```

---

### Notes

- This endpoint only supports the `GET` method.
- It must be called by an authenticated admin.
- State data is retained for 5 minutes from the last status change.
- If no action is taken within that window, the state will expire and be cleared automatically.

---



## 18. Deployment

> This project is fully Dockerised. You can deploy locally via `make docker-up`. See the Pre-requisites above. I have not yet hosted it.

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

