# GitLab CI/CD Variables Setup Guide

This document outlines all required GitLab CI/CD variables for the pipeline with SonarQube and ArgoCD integration.

## Required Variables

Configure these variables in **Settings → CI/CD → Variables** in your GitLab project.

### Container Registry

| Variable | Type | Protected | Masked | Description |
|----------|------|-----------|--------|-------------|
| `CI_REGISTRY` | Variable | No | No | GitLab container registry URL (auto-provided) |
| `CI_REGISTRY_USER` | Variable | No | No | Registry username (auto-provided) |
| `CI_REGISTRY_PASSWORD` | Variable | No | Yes | Registry password (auto-provided) |
| `CI_REGISTRY_IMAGE` | Variable | No | No | Full image path (auto-provided) |

### SonarQube Configuration

| Variable | Type | Protected | Masked | Description |
|----------|------|-----------|--------|-------------|
| `SONAR_HOST_URL` | Variable | No | No | SonarQube server URL (e.g., https://sonarqube.example.com) |
| `SONAR_TOKEN` | Variable | Yes | Yes | SonarQube project token (format: sqp_xxxxxxxxxxxxx) |

### ArgoCD Configuration

| Variable | Type | Protected | Masked | Description |
|----------|------|-----------|--------|-------------|
| `ARGOCD_SERVER` | Variable | No | No | ArgoCD server URL (e.g., argocd.example.com) |
| `ARGOCD_USERNAME` | Variable | Yes | No | ArgoCD username (e.g., admin or service account) |
| `ARGOCD_PASSWORD` | Variable | Yes | Yes | ArgoCD password or API token |

### GitOps Repository

| Variable | Type | Protected | Masked | Description |
|----------|------|-----------|--------|-------------|
| `GITOPS_SSH_PRIVATE_KEY` | File | Yes | Yes | SSH private key for GitOps repository access |
| `GITOPS_REPO` | Variable | No | No | GitOps repository URL (e.g., git@gitlab.com:yourorg/gitops-manifests.git) |
| `GITOPS_BRANCH` | Variable | No | No | GitOps repository branch (default: main) |

### Optional Configuration

| Variable | Type | Protected | Masked | Description |
|----------|------|-----------|--------|-------------|
| `GITLAB_USER_EMAIL` | Variable | No | No | Git commit email for manifest updates (default: uses CI user) |
| `GITLAB_USER_NAME` | Variable | No | No | Git commit name for manifest updates (default: uses CI user) |

## Setup Instructions

### 1. Generate SonarQube Token

```bash
# Login to SonarQube UI
# Navigate to: My Account → Security → Generate Tokens

# Token settings:
Name: gitlab-ci
Type: Project Analysis Token
Expires: Never (or 90 days for security)

# Copy token (format: sqp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxx)
```

Add to GitLab:
```
Variable: SONAR_TOKEN
Value: sqp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxx
Protected: Yes
Masked: Yes
Environment scope: *
```

### 2. Generate ArgoCD Token

**Option A: Use API Token (Recommended)**

```bash
# Login to ArgoCD CLI
argocd login argocd.example.com --username admin

# Generate token for CI/CD
argocd account generate-token \
  --account gitlab-ci \
  --id gitlab-ci-token

# Copy the generated token
```

**Option B: Use Password**

Use the admin password or create a dedicated service account.

Add to GitLab:
```
Variable: ARGOCD_USERNAME
Value: admin
Protected: Yes
Masked: No

Variable: ARGOCD_PASSWORD  
Value: <token or password>
Protected: Yes
Masked: Yes
```

### 3. Generate SSH Key for GitOps Repository

```bash
# Generate SSH key pair
ssh-keygen -t ed25519 -C "gitlab-ci@example.com" -f gitlab-ci-gitops

# This creates:
# - gitlab-ci-gitops (private key)
# - gitlab-ci-gitops.pub (public key)
```

**Add Public Key to GitLab GitOps Repository:**

1. Navigate to GitOps repository
2. Go to: Settings → Repository → Deploy Keys
3. Click "Add new key"
4. Paste contents of `gitlab-ci-gitops.pub`
5. Title: `GitLab CI/CD`
6. **Important**: Check "Write access allowed"
7. Click "Add key"

**Add Private Key to GitLab CI/CD Variables:**

```bash
# Display private key
cat gitlab-ci-gitops
```

Add to GitLab:
```
Variable: GITOPS_SSH_PRIVATE_KEY
Type: File
Value: <paste entire private key including BEGIN/END lines>
Protected: Yes
Masked: Yes
Environment scope: *
```

**Important**: The private key must include:
```
-----BEGIN OPENSSH PRIVATE KEY-----
<key content>
-----END OPENSSH PRIVATE KEY-----
```

### 4. Configure Environment-Specific Variables (Optional)

For different configurations per environment:

#### Development Environment

Create variables with **Environment scope: develop**:

```
Variable: ARGOCD_SERVER
Value: argocd-dev.example.com
Environment scope: develop

Variable: SONAR_HOST_URL
Value: https://sonarqube-dev.example.com
Environment scope: develop
```

#### Staging Environment

Create variables with **Environment scope: staging** or **Environment scope: main**:

```
Variable: ARGOCD_SERVER
Value: argocd-staging.example.com
Environment scope: staging
```

#### Production Environment

Create variables with **Environment scope: production** and **Protected: Yes**:

```
Variable: ARGOCD_SERVER
Value: argocd.example.com
Environment scope: production
Protected: Yes
```

## Verification

### Test SonarQube Connection

```bash
# Test from your local machine or CI runner
curl -u "${SONAR_TOKEN}:" "${SONAR_HOST_URL}/api/system/status"

# Should return: {"status":"UP",...}
```

### Test ArgoCD Connection

```bash
# Test ArgoCD authentication
argocd login ${ARGOCD_SERVER} \
  --username ${ARGOCD_USERNAME} \
  --password ${ARGOCD_PASSWORD} \
  --grpc-web

# List applications
argocd app list
```

### Test GitOps Repository Access

```bash
# Test SSH key
ssh-add gitlab-ci-gitops
ssh -T git@gitlab.com

# Should return: Welcome to GitLab, @username!

# Test clone
git clone ${GITOPS_REPO} test-clone
cd test-clone
git log --oneline -n 5
```

## Variable Security Best Practices

### 1. Use Protected Variables

Mark variables as **Protected** if they should only be available on:
- Protected branches (main, develop)
- Protected tags (v1.0.0, v2.0.0)

**Recommended for protection:**
- `SONAR_TOKEN`
- `ARGOCD_USERNAME`
- `ARGOCD_PASSWORD`
- `GITOPS_SSH_PRIVATE_KEY`

### 2. Use Masked Variables

Mark variables as **Masked** to hide values in job logs.

**Recommended for masking:**
- All tokens and passwords
- SSH private keys
- API keys

### 3. Use Environment Scopes

Limit variable availability to specific branches/environments:

```
Production secrets → Environment scope: production
Staging secrets → Environment scope: staging
Dev secrets → Environment scope: develop
```

### 4. Rotate Secrets Regularly

**Recommended rotation schedule:**
- SonarQube tokens: Every 90 days
- ArgoCD passwords: Every 90 days
- SSH keys: Every 180 days
- API tokens: Every 90 days

### 5. Audit Variable Access

Regularly review:
- Who has access to variable settings
- Which variables are in use
- When variables were last updated

GitLab audit logs: **Settings → Audit Events**

## Troubleshooting

### Issue: "SonarQube token is invalid"

**Solution:**
1. Verify token format starts with `sqp_`
2. Check token hasn't expired in SonarQube
3. Regenerate token if needed
4. Ensure variable is not accidentally protected when used on feature branches

### Issue: "ArgoCD authentication failed"

**Solution:**
1. Verify ArgoCD server URL (no https:// prefix in variable)
2. Check username/password are correct
3. Test login manually with argocd CLI
4. Check if password contains special characters that need escaping

### Issue: "Permission denied (publickey)" for GitOps repo

**Solution:**
1. Verify SSH private key format (includes BEGIN/END lines)
2. Check deploy key is added to GitOps repository
3. Ensure "Write access allowed" is checked on deploy key
4. Test SSH key locally: `ssh -i gitlab-ci-gitops -T git@gitlab.com`
5. Check if key is passphrase-protected (must be passphrase-free for CI)

### Issue: "Variable not found" in pipeline

**Solution:**
1. Check variable name matches exactly (case-sensitive)
2. Verify environment scope matches your branch/tag
3. Check if variable is protected but branch isn't protected
4. Look for typos in variable name

### Issue: Pipeline shows variable value in logs (security issue)

**Solution:**
1. Mark variable as **Masked** immediately
2. Rotate the exposed secret
3. Check if you're echo'ing variables in scripts (remove)
4. Review job logs and clear sensitive data

## Environment-Specific Configuration Matrix

| Environment | Branch/Tag | Variables Scope | Protected | Auto-Deploy |
|-------------|-----------|-----------------|-----------|-------------|
| Development | `develop` | develop | No | Yes |
| Staging | `main` | staging or main | Recommended | Yes |
| Production | `v*` tags | production | Yes | Manual |

## Complete Variable List Template

Copy this checklist when setting up a new project:

```
Container Registry (Auto-provided)
[ ] CI_REGISTRY
[ ] CI_REGISTRY_USER
[ ] CI_REGISTRY_PASSWORD
[ ] CI_REGISTRY_IMAGE

SonarQube
[ ] SONAR_HOST_URL = https://sonarqube.example.com
[ ] SONAR_TOKEN = sqp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxx

ArgoCD
[ ] ARGOCD_SERVER = argocd.example.com
[ ] ARGOCD_USERNAME = admin
[ ] ARGOCD_PASSWORD = ********

GitOps Repository
[ ] GITOPS_SSH_PRIVATE_KEY = (SSH private key content)
[ ] GITOPS_REPO = git@gitlab.com:yourorg/gitops-manifests.git
[ ] GITOPS_BRANCH = main

Optional
[ ] GITLAB_USER_EMAIL = ci@example.com
[ ] GITLAB_USER_NAME = GitLab CI
```

## Additional Resources

- [GitLab CI/CD Variables Documentation](https://docs.gitlab.com/ee/ci/variables/)
- [Protected Variables](https://docs.gitlab.com/ee/ci/variables/#protected-cicd-variables)
- [Masked Variables](https://docs.gitlab.com/ee/ci/variables/#mask-a-cicd-variable)
- [Environment Scopes](https://docs.gitlab.com/ee/ci/environments/)
- [SonarQube Authentication](https://docs.sonarqube.org/latest/user-guide/user-token/)
- [ArgoCD Tokens](https://argo-cd.readthedocs.io/en/stable/user-guide/commands/argocd_account_generate-token/)