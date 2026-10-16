# orders-api source repository

Change `REPLACE_GITHUB_OWNER` in Jenkinsfile. Create a Jenkins **multibranch pipeline** and a build agent labeled `go-docker-git-yq` with Go, Docker client/daemon access, Git, `yq` v4 and SSH available. Install Jenkins Git, GitHub, Credentials Binding, SSH Agent, Pipeline plugins. Create `ghcr-push` username/token credential with registry package-write permissions and `gitops-deploy-key` SSH private-key credential with GitOps repository write access. Use a dedicated agent, not the Jenkins controller, for builds.

For a first manual check: `go test ./...`, then `docker build -t orders-api:dev .` and `docker run --rm -p 8080:8080 orders-api:dev`. `curl localhost:8080/` and `curl localhost:8080/metrics`.

Set GHCR package public for this demo if you don't want Kubernetes imagePullSecrets. Do not put tokens in Git. This repo only does CI; Argo CD deploys from the other repo.
