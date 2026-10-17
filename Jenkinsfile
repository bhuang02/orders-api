
pipeline {
    agent {
        kubernetes {
            cloud 'kubernetes'
            defaultContainer 'go'
            yaml '''
apiVersion: v1
kind: Pod
spec:
  securityContext:
    runAsUser: 1000
    runAsGroup: 1000
    fsGroup: 1000
  containers:
    - name: go
      image: golang:1.25-bookworm
      command: ['cat']
      tty: true
      env:
        - name: GOPATH
          value: /tmp/gopath
        - name: GOCACHE
          value: /tmp/gocache
      resources:
        requests:
          cpu: 100m
          memory: 256Mi
        limits:
          cpu: '1000m'
          memory: 1Gi

    - name: buildkit
      image: moby/buildkit:rootless
      command: ['cat']
      tty: true
      env:
        - name: BUILDKITD_FLAGS
          value: '--oci-worker-no-process-sandbox --oci-worker-snapshotter=native'
      securityContext:
        runAsUser: 1000
        runAsGroup: 1000
        seccompProfile:
          type: Unconfined
        appArmorProfile:
          type: Unconfined
      resources:
        requests:
          cpu: 250m
          memory: 512Mi
        limits:
          cpu: '2000m'
          memory: 4Gi
      volumeMounts:
        - name: buildkit-state
          mountPath: /home/user/.local/share/buildkit

  volumes:
    - name: buildkit-state
      emptyDir: {}
'''
        }
    }

    options {
        skipDefaultCheckout(true)
        disableConcurrentBuilds()
    }

    environment {
        IMAGE_REPO = 'ghcr.io/bhuang02/orders-api'
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm

                script {
                    env.IMAGE_TAG = sh(
                        script: 'git rev-parse HEAD',
                        returnStdout: true
                    ).trim()
                }

                echo "Image to publish: ${env.IMAGE_REPO}:${env.IMAGE_TAG}"
            }
        }

        stage('Unit tests') {
            steps {
                sh 'go test ./...'
            }
        }

        stage('Build and push image') {
            steps {
                container('buildkit') {
                    withCredentials([
                        usernamePassword(
                            credentialsId: 'ghcr-push',
                            usernameVariable: 'GHCR_USER',
                            passwordVariable: 'GHCR_TOKEN'
                        )
                    ]) {
                        sh '''#!/bin/sh
set -eu
set +x

# Keep registry credentials out of Git and the Jenkins workspace.
export DOCKER_CONFIG="$(mktemp -d)"
trap 'rm -rf "$DOCKER_CONFIG"' EXIT

AUTH_B64="$(printf '%s' "$GHCR_USER:$GHCR_TOKEN" | base64 | tr -d '\\n')"

printf '{"auths":{"ghcr.io":{"auth":"%s"}}}\\n' "$AUTH_B64" > "$DOCKER_CONFIG/config.json"
chmod 600 "$DOCKER_CONFIG/config.json"
unset AUTH_B64

# Build directly from the checked-out Dockerfile, then push to GHCR.
buildctl-daemonless.sh build \\
  --progress=plain \\
  --frontend dockerfile.v0 \\
  --local context=. \\
  --local dockerfile=. \\
  --output "type=image,name=$IMAGE_REPO:$IMAGE_TAG,push=true"
'''
                    }
                }
            }
        }

        stage('Update GitOps staging') {
            steps {
                container('go') {
                    withCredentials([
                        usernamePassword(
                            credentialsId: 'gitops-write',
                            usernameVariable: 'GITOPS_USER',
                            passwordVariable: 'GITOPS_TOKEN'
                        )
                    ]) {
                        sh '''#!/bin/sh
set -eu
set +x

# Git authentication without putting the token in the repository URL.
ASKPASS="$(mktemp)"
trap 'rm -f "$ASKPASS"' EXIT

cat > "$ASKPASS" <<'EOF'
#!/bin/sh
case "$1" in
  *Username*) printf '%s\\n' "$GITOPS_USER" ;;
  *Password*) printf '%s\\n' "$GITOPS_TOKEN" ;;
  *) exit 1 ;;
esac
EOF

chmod 700 "$ASKPASS"
export GIT_ASKPASS="$ASKPASS"
export GIT_TERMINAL_PROMPT=0

# Clone the separate GitOps repository.
git clone https://github.com/bhuang02/orders-gitops.git gitops-workdir

cd gitops-workdir

# Update only the staging image tag.
VALUES=charts/orders-api/values-staging.yaml

sed -i "s|^  tag:.*|  tag: '$IMAGE_TAG'|" "$VALUES"

# Verify that the image tag was updated.
grep -Fxq "  tag: '$IMAGE_TAG'" "$VALUES"

git config user.name "Jenkins CI"
git config user.email "jenkins-ci@users.noreply.github.com"

git add "$VALUES"

if git diff --cached --quiet; then
    echo "Staging already uses this image."
else
    git commit -m "ci: deploy orders-api $IMAGE_TAG"
    git push origin HEAD:main
fi
'''
                    }
                }
            }
        }
    }

    post {
        success {
            echo "Published ${env.IMAGE_REPO}:${env.IMAGE_TAG}"
        }
    }
}

