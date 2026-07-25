# .github/PULL_REQUEST_TEMPLATE.md

##  Summary
<!-- Provide a high-level overview of the feature, refactoring, or fix introduced in this PR. -->


##  Key Changes
<!-- Bullet points highlighting the significant technical updates. -->
- 

##  How to Test
<!-- Step-by-step instructions for testing or verifying the changes locally. -->
```bash
# Example commands to verify this PR
```

##  Verification Results / Logs
<!-- Paste relevant command outputs, terminal logs, or metrics JSON to prove functionality. -->
```json

```

##  Checklist
- [ ] Code follows the style guidelines and architectural patterns of this project.
- [ ] Self-review has been conducted before requesting review.
- [ ] Unit tests or integration tests have been added/updated.
- [ ] All tests pass locally (`go test ./...` or `docker compose`).


---
#  English Action Verbs for PR Titles & Commits

| Verb | Usage | Example |
|---|---|---|
| **Add / Implement** | Adding new features or modules | `Implement AnalyticsClient for parallel execution` |
| **Fix / Resolve** | Bug fixes, type errors, parse fixes | `Fix JSON3 deserialization type mismatch in Julia service` |
| **Refactor** | Restructuring code without changing functionality | `Refactor HTTP request handling in gateway-go` |
| **Optimize** | Performance or resource improvements | `Optimize parallel goroutine response aggregation` |
| **Integrate** | Combining multiple services or libraries | `Integrate R and Julia microservices into Go gateway` |
