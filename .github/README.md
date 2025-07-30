# 🤖 AdvanceGG GitHub Workflows

Automated CI/CD pipeline for AdvanceGG with comprehensive testing, security, and release management.

## 🚀 Workflows Overview

### 1. Release Workflow (`release.yml`)
**Trigger:** Manual dispatch with version input
**Purpose:** Automated release creation with multi-platform binaries

**Features:**
- ✅ Cross-platform binary building (Windows, macOS, Linux)
- ✅ Python package building and PyPI publishing
- ✅ Node.js package building and npm publishing
- ✅ GitHub release creation with assets
- ✅ Automatic version tagging
- ✅ Comprehensive release notes generation
- ✅ SHA256 checksums for all binaries

**Usage:**
1. Go to Actions → Release AdvanceGG
2. Click "Run workflow"
3. Enter version (e.g., `1.0.0`)
4. Select release type
5. Choose if pre-release or draft
6. Click "Run workflow"

### 2. CI Workflow (`ci.yml`)
**Trigger:** Push to main/develop, Pull requests
**Purpose:** Continuous integration testing

**Features:**
- ✅ Go core library testing (multiple versions)
- ✅ Cross-platform native library building
- ✅ Python ecosystem testing (multiple Python versions)
- ✅ Node.js ecosystem testing (multiple Node versions)
- ✅ Security scanning and vulnerability checks
- ✅ Code coverage reporting
- ✅ Performance benchmarking

### 3. Security Workflow (`security.yml`)
**Trigger:** Daily schedule, dependency changes
**Purpose:** Security monitoring and maintenance

**Features:**
- ✅ Vulnerability scanning (Go, Python, Node.js)
- ✅ Automated dependency updates
- ✅ License compliance checking
- ✅ Code quality metrics
- ✅ Security issue creation
- ✅ Automated PR creation for updates

### 4. Documentation Workflow (`docs.yml`)
**Trigger:** Documentation changes, manual dispatch
**Purpose:** Documentation building and deployment

**Features:**
- ✅ Documentation validation and building
- ✅ GitHub Pages deployment
- ✅ Documentation quality checks
- ✅ Automatic index generation
- ✅ Markdown linting and optimization

## 🔧 Setup Requirements

### Repository Secrets
Add these secrets in GitHub repository settings:

```bash
# Required for package publishing
PYPI_API_TOKEN=pypi-...          # PyPI API token for Python packages
NPM_TOKEN=npm_...                # npm token for Node.js packages

# Optional for enhanced features
PYPI_TEST_TOKEN=pypi-...         # PyPI Test API token for pre-releases
CODECOV_TOKEN=...                # Codecov token for coverage reports
```

### Repository Settings
Configure these settings for optimal security:

1. **Branch Protection Rules** (main branch):
   - ✅ Require status checks to pass
   - ✅ Require branches to be up to date
   - ✅ Require review from code owners
   - ✅ Dismiss stale reviews
   - ✅ Restrict pushes to admins only

2. **Security Settings**:
   - ✅ Enable Dependabot alerts
   - ✅ Enable Dependabot security updates
   - ✅ Enable secret scanning
   - ✅ Enable code scanning (CodeQL)

3. **Actions Settings**:
   - ✅ Allow actions and reusable workflows
   - ✅ Allow actions created by GitHub
   - ✅ Allow specified actions (add trusted actions)

## 📋 Release Process

### Automated Release (Recommended)
1. **Prepare Release**:
   - Update CHANGELOG.md
   - Ensure all tests pass
   - Update version in documentation

2. **Trigger Release**:
   - Go to Actions → Release AdvanceGG
   - Enter version number (semantic versioning)
   - Select release type
   - Add changelog notes
   - Run workflow

3. **Post-Release**:
   - Verify packages on PyPI and npm
   - Test installation on different platforms
   - Update documentation if needed

### Manual Release (Fallback)
If automated release fails:

```bash
# 1. Build binaries locally
make build-all

# 2. Create release manually
gh release create v1.0.0 \
  --title "AdvanceGG v1.0.0" \
  --notes-file RELEASE_NOTES.md \
  dist/binaries/*

# 3. Publish packages
cd ecosystem/python && twine upload dist/*
cd ecosystem/nodejs && npm publish
```

## 🛡️ Security Features

### Automated Security Scanning
- **Daily vulnerability scans** for all dependencies
- **License compliance checking** for legal requirements
- **Code quality metrics** tracking over time
- **Automated security issue creation** for critical findings

### Dependency Management
- **Dependabot integration** for automated updates
- **Security-first update strategy** prioritizing patches
- **Automated testing** of dependency updates
- **Pull request creation** for review and approval

### Safe Release Process
- **Multi-stage validation** before release
- **Checksum verification** for all binaries
- **Signed releases** with GPG (when configured)
- **Rollback capability** through GitHub releases

## 📊 Monitoring & Metrics

### Build Status
Monitor workflow status at:
- **CI Status**: ![CI](https://github.com/GrandpaEJ/advancegg/workflows/CI/badge.svg)
- **Security**: ![Security](https://github.com/GrandpaEJ/advancegg/workflows/Security/badge.svg)
- **Docs**: ![Docs](https://github.com/GrandpaEJ/advancegg/workflows/Documentation/badge.svg)

### Performance Tracking
- **Benchmark results** stored in GitHub Pages
- **Performance regression detection** in CI
- **Memory usage monitoring** across platforms
- **Build time optimization** tracking

### Quality Metrics
- **Code coverage** reporting via Codecov
- **Security score** from vulnerability scans
- **Documentation completeness** percentage
- **Dependency freshness** tracking

## 🚨 Troubleshooting

### Common Issues

**Release workflow fails:**
```bash
# Check workflow logs
gh run list --workflow=release.yml
gh run view <run-id>

# Common fixes:
# 1. Verify version format (semantic versioning)
# 2. Check if tag already exists
# 3. Ensure secrets are configured
# 4. Verify Go/Python/Node.js versions
```

**Package publishing fails:**
```bash
# Check API tokens
# 1. PyPI token has correct permissions
# 2. npm token is not expired
# 3. Package name is available
# 4. Version doesn't already exist
```

**Security scan alerts:**
```bash
# Review security reports
# 1. Check workflow artifacts
# 2. Review Dependabot PRs
# 3. Update dependencies manually if needed
# 4. Create security issues for tracking
```

### Emergency Procedures

**Critical Security Issue:**
1. Create hotfix branch
2. Apply security patch
3. Run security scan locally
4. Create emergency release
5. Notify users immediately

**Build System Failure:**
1. Check GitHub status page
2. Verify runner availability
3. Use manual build process
4. Contact GitHub support if needed

## 🤝 Contributing to Workflows

### Adding New Workflows
1. Create workflow file in `.github/workflows/`
2. Follow naming convention: `purpose.yml`
3. Add comprehensive documentation
4. Test thoroughly before merging
5. Update this README

### Modifying Existing Workflows
1. Create feature branch
2. Make changes with clear commit messages
3. Test in fork repository first
4. Create pull request with detailed description
5. Request review from maintainers

### Best Practices
- ✅ Use semantic versioning for all releases
- ✅ Include comprehensive error handling
- ✅ Add timeout limits to prevent hanging
- ✅ Use caching to improve performance
- ✅ Follow security best practices
- ✅ Document all configuration options

## 📞 Support

For workflow-related issues:
1. Check existing GitHub Issues
2. Review workflow logs and artifacts
3. Consult this documentation
4. Create new issue with detailed information

**Workflow Maintainers:**
- @GrandpaEJ (Primary maintainer)

---

*This documentation is automatically updated by the documentation workflow.*
