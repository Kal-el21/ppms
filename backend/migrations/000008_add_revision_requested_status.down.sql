UPDATE project_requests
SET status = 'REJECTED'
WHERE status = 'REVISION_REQUESTED';

ALTER TABLE project_requests
    DROP CONSTRAINT IF EXISTS project_requests_status_check;

ALTER TABLE project_requests
    ADD CONSTRAINT project_requests_status_check CHECK (status IN (
        'DRAFT',
        'SUBMITTED',
        'UNDER_REVIEW',
        'APPROVED',
        'REJECTED',
        'REVISED'
    ));
