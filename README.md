# Releasetracker (for Kubernetes)

This is an API tool for tracking Kubernetes deployed releases.

It monitors Kubernetes objects to track deployed release and will compare
to upstream release to determine if the deployed release is up to date.

## Components

- Release Syncer: a binary intended to sync tracked releases with upstream
- Release Tracker API: a REST API to register syncer and query release status