# File: docker/jenkins/Dockerfile
# --------------------------------------------------------------------
# Dockerfile for Custom Jenkins Controller:
#   - Based on a specific Jenkins release using JDK 21. (User specified 2.505)
#   - Installs Docker CLI tools (docker, buildx, compose) inside the container.
#   - Installs kubectl CLI tool. (RE-ADDED)
#   - Installs Jenkins plugins specified in 'plugins.txt'.
#   - Configures Jenkins via Configuration-as-Code (JCasC) files.
#   - Grants the Jenkins user access to the Docker socket via group membership.
# --------------------------------------------------------------------

# Using user-specified Jenkins weekly release with JDK 21
# Note: Weekly releases may be less stable than LTS (e.g., 2.452.3-lts-jdk21)
FROM jenkins/jenkins:2.505-jdk21

# --- Arguments ---
# Build argument for the host's Docker group GID.
# Default value used if not provided during 'docker build' or via 'docker compose build args'.
# Find on host via: stat -c '%g' /var/run/docker.sock
ARG HOST_DOCKER_GROUP_ID=999

# --- System Setup (as root) ---
USER root

# Install prerequisites (curl, gpg) and Docker CLI tools (docker, buildx, compose plugin)
# Uses official Docker repository for Debian (base OS of jenkins/jenkins images)
# This part is independent of the Java version used by Jenkins itself.
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates curl gnupg make \
 && install -m 0755 -d /etc/apt/keyrings \
 && curl -fsSL https://download.docker.com/linux/debian/gpg -o /etc/apt/keyrings/docker.asc \
 && chmod a+r /etc/apt/keyrings/docker.asc \
 && echo \
      "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/debian \
      $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
      tee /etc/apt/sources.list.d/docker.list > /dev/null \
 && apt-get update && apt-get install -y --no-install-recommends \
      docker-ce-cli \
      docker-buildx-plugin \
      docker-compose-plugin \
 # Clean up APT cache to reduce image size
 && apt-get clean && rm -rf /var/lib/apt/lists/*

# *** RE-ADDED: Install kubectl (Essential for the Jenkins Pipeline) ***
RUN echo "Installing kubectl..." \
 && KUBECTL_VERSION=$(curl -L -s https://dl.k8s.io/release/stable.txt) \
 && echo "Downloading kubectl version ${KUBECTL_VERSION}..." \
 && curl -LO "https://dl.k8s.io/release/${KUBECTL_VERSION}/bin/linux/amd64/kubectl" \
 # Optional: Validate checksum (get checksum from release page)
 # && echo "<checksum> kubectl" | sha256sum --check \
 && install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl \
 && rm kubectl \
 && echo "kubectl installation complete." \
 && /usr/local/bin/kubectl version --client

# Configure Docker Group Access for the 'jenkins' user
# Ensures the 'docker' group exists with the specified HOST_DOCKER_GROUP_ID
# and adds the 'jenkins' user to it, allowing access to the mounted docker.sock
RUN echo "Preparing docker group with GID=${HOST_DOCKER_GROUP_ID}..." \
 && DOCKER_GROUP_NAME="docker" \
 && if ! getent group ${HOST_DOCKER_GROUP_ID} >/dev/null 2>&1; then \
      echo "Group with GID ${HOST_DOCKER_GROUP_ID} doesn't exist. Creating '${DOCKER_GROUP_NAME}' group..." ; \
      groupadd -g ${HOST_DOCKER_GROUP_ID} ${DOCKER_GROUP_NAME} ; \
    else \
      EXISTING_GROUP_NAME=$(getent group ${HOST_DOCKER_GROUP_ID} | cut -d: -f1); \
      if [ "${EXISTING_GROUP_NAME}" != "${DOCKER_GROUP_NAME}" ]; then \
        # GID exists but has the wrong name - this is usually an error state.
        echo "ERROR: GID ${HOST_DOCKER_GROUP_ID} is already used by group '${EXISTING_GROUP_NAME}' but expected '${DOCKER_GROUP_NAME}'." >&2; \
        exit 1; \
      else \
        echo "Group '${DOCKER_GROUP_NAME}' with GID ${HOST_DOCKER_GROUP_ID} already exists."; \
      fi; \
    fi \
 && echo "Adding user 'jenkins' to '${DOCKER_GROUP_NAME}' (GID: ${HOST_DOCKER_GROUP_ID}) group." \
 # Add jenkins user to the potentially newly created/existing docker group
 && usermod -aG ${DOCKER_GROUP_NAME} jenkins

# --- Jenkins Configuration ---

# Install Plugins defined in plugins.txt
# Assumes plugins.txt is located at docker/jenkins/plugins.txt relative to the build context
# Ensure plugin versions in plugins.txt are compatible with the Jenkins core version being used
COPY --chown=jenkins:jenkins docker/jenkins/plugins.txt /usr/share/jenkins/ref/plugins.txt
RUN echo "Installing plugins specified in /usr/share/jenkins/ref/plugins.txt..." \
 && jenkins-plugin-cli --verbose --plugin-file /usr/share/jenkins/ref/plugins.txt

# Setup Jenkins Configuration as Code (JCasC)
# Copy the directory containing JCasC *.yaml files
# Assumes casc_configs dir is located at docker/jenkins/casc_configs relative to the build context
COPY --chown=jenkins:jenkins docker/jenkins/casc_configs /var/jenkins_home/casc_configs
# Optional: Explicitly set the JCasC config path environment variable
ENV CASC_JENKINS_CONFIG /var/jenkins_home/casc_configs

# --- Final Permissions and User ---

# Ensure jenkins user owns its home directory contents
RUN chown -R jenkins:jenkins /var/jenkins_home

# Switch back to the non-root 'jenkins' user to run the Jenkins process
USER jenkins

# Jenkins entrypoint is inherited from the base image
