# Security remediation checklist

This checklist tracks authorization fixes required before production use with real school data.

## Required controls

- Verbali: verify that every class belongs to the requested school before creating meetings or minutes.
- Rubrics: verify that the teacher is assigned to the student's class before recording an assessment.
- PCTO: verify the complete tutor-to-project ownership chain before approving hour logs.
- Didactic materials: verify that `class_id` belongs to `school_id` at write time.
- Agenda: verify that the student is enrolled in the class before accepting completion data.

## Required tests

Each control needs a negative test proving that a user from another school, class, teacher assignment, or PCTO project receives an authorization error.

This branch intentionally does not weaken existing authorization checks or claim that these controls are complete. The application must not be considered ready for real school data until the checks and negative tests pass.
