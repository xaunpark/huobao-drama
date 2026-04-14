---
date: "2026-04-14"
problem_type: "integration_issue"
component: "api_integration"
severity: "high"
symptoms:
  - "Video upscaling tasks get stuck indefinitely in 'Upscaling' state in the UI."
  - "Backend logs show jobs stuck polling in 'PENDING' state up to the maximum timeout (10 minutes/200 attempts), even when the Flow-Tool API explicitly returned failed statuses or re-queue errors."
root_cause: "logic_error"
tags:
  - polling
  - flow-tool
  - transient-errors
  - stalled-tasks
related_issues: []
---

## 1. Problem Statement

Video upscaling tasks triggered from the backend were getting permanently "stuck" in the `Upscaling` state. Both the Batch Action Studio and the Professional Editor UI reflected this frozen status, blocking further progress for 10 minutes until the backend's hard polling timeout eventually failed the tasks.

## 2. Symptoms Observed

- `POST /v1/jobs/.../upscale` queued successfully, but jobs either failed with `403 reCAPTCHA rejected` or `500 Internal error` immediately without the UI state updating.
- Backend continuously logged: `Video generation in progress {"flow_status": "PENDING"}` up to 200 attempts despite Flow-Tool failing the job.
- When jobs were genuinely stalled in Flow-Tool's queue (nobody picking them up), backend also endlessly waited in `PENDING` state.
- Even when the backend successfully caught a terminal error and saved the status as `upscale_failed`, the frontend remained frozen at "Upscaling...".

## 3. Investigation Steps

| Step | Discovery |
|---|---|
| Checked Flow-Tool API spec | Confirmed Flow-Tool returns status in uppercase (`SUCCESS`, `FAILED`), and uses `PENDING` alongside an `error` field for transient issues (e.g., `401 token expired` triggering an auto-requeue). |
| Reviewed Backend Polling Logic (`flowtool_video_client.go`) | Found that the `default` case in the status `switch` was silently swallowing the `error` field if the status wasn't exactly `"failed"` or `"error"`. |
| Analyzed Backend State Machine (`video_generation_service.go`) | Discovered polling loop failed to distinguish between terminal errors (fail immediately), transient errors (retry a few times), and stale queue states (fail if never dispatched). |
| Traced Frontend Status Updates | Found that frontend polling dynamically stopped when status left `processing`, omitting `upscaling`. Furthermore, it only checked for the `failed` status, missing the backend's specific `upscale_failed` state. |

## 4. Root Cause Analysis

There were multiple overlapping logic flaws causing desynchronization between Flow-Tool, the Backend, and the Frontend:
1. **Error Swallowing:** The backend ignored any error messages from the Flow-Tool API if the status wasn't exactly `FAILED` (e.g., when it was `PENDING` because of a transient 401 re-queue).
2. **Missing Terminal Identification:** The backend treated explicit `"FAILED"` statuses the same as transient errors, resulting in unnecessary retries.
3. **No Stale Queue Detection:** Jobs queued but never dispatched by Flow-Tool would stay `PENDING` indefinitely, burning up the full 10-minute timeout.
4. **UI Polling Conditions:** The frontend's polling logic (`BatchGenerationDialog`, `BatchReferenceVideoStudio.vue`, `ProfessionalEditor.vue`) did not activate or continue correctly when a job was in the `upscaling` state, and completely ignored the `upscale_failed` status.

## 5. Working Solution

**Backend Changes (`pkg/video/flowtool_video_client.go` & `application/services/video_generation_service.go`):**
- Updated `flowtool_video_client.go` to always propagate the `error` field in the API response, regardless of the accompanying `status`.
- Introduced a `consecutivePending` counter (15 polls ≈ 45s) to fail jobs that are stuck queued without being dispatched.
- Introduced a `consecutiveErrors` counter (3 polls ≈ 9s) to tolerate transient `PENDING + error` states (like 401 re-queues).
- Automatically fail immediately with no retries if Flow-Tool explicitly returns `status: "FAILED"` or `status: "ERROR"`.

**Frontend Changes (`BatchGenerationDialog.vue`, `BatchReferenceVideoStudio.vue`):**
```javascript
// Changed:
} else if (check.status === 'failed') {
// To:
} else if (check.status === 'failed' || check.status === 'upscale_failed') {
```

**Frontend Changes (`ProfessionalEditor.vue`):**
- Added `'upscaling'` status to both `hasPendingOrProcessing` conditions to ensure polling activates and sustains while an upscale task is ongoing.
- Added localization maps for `upscaled` and `upscale_failed` to `getStatusText`.

## 6. Prevention Strategies

- **Strict Status Mapping Definition:** Maintain a strict shared definition of all terminal and non-terminal states between the UI and Backend (`pending`, `processing`, `completed`, `failed`, `upscaling`, `upscaled`, `upscale_failed`).
- **Always Validate `status` against `error_msg`:** Treat an explicit `error_msg` payload from any third-party API as a high-priority signal, even if the primary `status` enum suggests the job is still active.
- **Fail Fast over Timeouts:** Implement active stall detection (`consecutivePending`) instead of relying solely on a hard maximum polling limit, to dramatically improve UX.

## 7. Cross-References

- `flowtool_video_client.go:GetTaskStatus`
- `video_generation_service.go:pollTaskStatus`
- `docs/upscale_api_response_spec.md.resolved`
