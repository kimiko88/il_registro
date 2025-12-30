# Branch Protection Rules Configuration

This document outlines the recommended branch protection rules for the `main` and `develop` branches.

## Rules for `main`

1.  **Require a pull request before merging**
    *   [x] Require approvals: **2**
    *   [x] Dismiss stale pull request approvals when new commits are pushed
    *   [x] Require review from Code Owners

2.  **Require status checks to pass before merging**
    *   **Backend Checks:**
        *   `test-backend`
        *   `lint-backend`
        *   `security-backend`
    *   **Frontend Checks:**
        *   `test-frontend`
        *   `lint-frontend`
        *   `security-frontend`
    *   **Common Checks:**
        *   `trivy-scan`

3.  **Require conversation resolution before merging**
    *   [x] All threads must be resolved

4.  **Include administrators**
    *   [x] Enforce all configured restrictions for administrators

5.  **Restrict who can push to matching branches**
    *   No one (force PRs)

## Rules for `develop`

1.  **Require a pull request before merging**
    *   [x] Require approvals: **1**

2.  **Require status checks to pass before merging**
    *   Similar to `main`, ensure tests and linting pass.

3.  **Allow force pushes**
    *   [ ] Disabled (Recommended)

---

**Apply via:** GitHub Repository Settings -> Branches -> Add rule
