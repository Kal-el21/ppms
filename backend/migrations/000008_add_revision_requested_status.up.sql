ALTER TABLE project_requests
    DROP CONSTRAINT IF EXISTS project_requests_status_check;

ALTER TABLE project_requests
    ADD CONSTRAINT project_requests_status_check CHECK (status IN (
        'DRAFT',
        'SUBMITTED',
        'UNDER_REVIEW',
        'APPROVED',
        'REJECTED',
        'REVISION_REQUESTED',
        'REVISED'
    ));
