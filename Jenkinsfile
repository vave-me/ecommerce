// groovyfunctions
def applyManifests(String path, String kubectlCmd, boolean required = true, String namespace = null) {
    def cmd = "${kubectlCmd}"
    def applyType = "-f"
    def originalPath = path
    if (namespace) {
        cmd += " -n ${namespace}"
    }

    if (path.endsWith('/')) {
        applyType = "-R -f"
        path = path.substring(0, path.length() - 1)
    }

    cmd += " apply ${applyType} ${originalPath}"
    def label = "Apply ${originalPath}"

    try {
        sh(script: cmd, label: label)
    } catch (e) {
        if (required) {
            error("Failed to apply manifests at ${originalPath}: ${e.getMessage()}")
        } else {
            echo "WARN: Manifests at ${originalPath} not found or failed to apply (optional). Error: ${e.getMessage()}"
        }
    }
}

def updateImageTag(String resourceType, String resourceName, String containerName, String imageName, String kubectlCmd, String namespace = env.K8S_NAMESPACE) {
    try {
        echo "Updating image for ${resourceType}/${resourceName} container ${containerName} to ${imageName} in namespace ${namespace}"
        sh(script: "${kubectlCmd} set image ${resourceType}/${resourceName} ${containerName}=${imageName} -n ${namespace}", label: "Set image ${resourceName}")
    } catch (e) {
        error("Failed to set image for ${resourceType}/${resourceName}: ${e.getMessage()}")
    }
}

def ensureServiceDb(String serviceName, String needsPostgis, String kubectlCmd, String dbNamespace = "database", String jobNamespace = env.K8S_NAMESPACE) {
    // Retrieve password directly
    def getPasswordCmd = "${kubectlCmd} get secret -n ${jobNamespace} service-db-passwords -o jsonpath='{.data.${serviceName.toUpperCase()}_PASS}' | base64 -d"
    def servicePassword = sh(script: getPasswordCmd, returnStdout: true).trim()

    if (!servicePassword) {
        error("Failed to retrieve password for service ${serviceName} from secret service-db-passwords key ${serviceName.toUpperCase()}_PASS in namespace ${jobNamespace}")
    }

    // Helper function for SQL single-quote escaping (replace ' with '')
    def sqlQuoteLiteral = { str -> str.replaceAll("'", "''") }

    // Write SQL to a temporary file instead of trying to pass it through shell heredoc
    def tempSqlFile = "db_setup_${serviceName}.sql"
    def dbUser = "${serviceName}_user" // Define user variable for clarity

    // Construct SQL script without complex escaping - use literal quotes for identifiers
    // MODIFIED SQL to update password if user exists
    def sqlScript = """
    -- Ensure database exists
    DO \$\$
    BEGIN
       IF NOT EXISTS (SELECT FROM pg_database WHERE datname = '${sqlQuoteLiteral(serviceName)}') THEN
          CREATE DATABASE "${sqlQuoteLiteral(serviceName)}" TEMPLATE commondb; -- Quoting database name for safety
          RAISE NOTICE 'Database created: ${sqlQuoteLiteral(serviceName)}';
       ELSE
          RAISE NOTICE 'Database already exists: ${sqlQuoteLiteral(serviceName)}';
       END IF;
    END
    \$\$;

    -- Ensure user exists and password is synchronized
    DO \$\$
    BEGIN
       IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = '${sqlQuoteLiteral(dbUser)}') THEN
          CREATE USER "${sqlQuoteLiteral(dbUser)}" WITH ENCRYPTED PASSWORD '${sqlQuoteLiteral(servicePassword)}'; -- Quoting user name for safety
          RAISE NOTICE 'User created: ${sqlQuoteLiteral(dbUser)}';
       ELSE
          RAISE NOTICE 'User already exists: ${sqlQuoteLiteral(dbUser)}. Ensuring password is synchronized.';
          ALTER USER "${sqlQuoteLiteral(dbUser)}" WITH ENCRYPTED PASSWORD '${sqlQuoteLiteral(servicePassword)}'; -- Quoting user name for safety
          RAISE NOTICE 'Password updated for existing user: ${sqlQuoteLiteral(dbUser)}';
       END IF;
    END
    \$\$;

    -- Grant basic connect rights
    GRANT CREATE, CONNECT ON DATABASE "${sqlQuoteLiteral(serviceName)}" TO "${sqlQuoteLiteral(dbUser)}"; -- Quoting names

    -- Switch connection to the target DB to apply schema changes
    \\c "${sqlQuoteLiteral(serviceName)}"

    GRANT CONNECT ON DATABASE "${sqlQuoteLiteral(serviceName)}" TO postgres; -- Quoting name
    CREATE SCHEMA IF NOT EXISTS "${sqlQuoteLiteral(serviceName)}"; -- Quoting name
    GRANT CREATE, USAGE ON SCHEMA "${sqlQuoteLiteral(serviceName)}" TO "${sqlQuoteLiteral(dbUser)}"; -- Quoting names

    ${ needsPostgis == 'true' ? """
    -- Ensuring PostGIS extension
    CREATE EXTENSION IF NOT EXISTS postgis SCHEMA public; -- Typically 'public' or a dedicated extensions schema
    GRANT SELECT, INSERT, UPDATE, DELETE ON public.geometry_columns TO "${sqlQuoteLiteral(dbUser)}";
    GRANT SELECT, INSERT, UPDATE, DELETE ON public.spatial_ref_sys TO "${sqlQuoteLiteral(dbUser)}";
    GRANT SELECT, INSERT, UPDATE, DELETE ON public.geography_columns TO "${sqlQuoteLiteral(dbUser)}";
    """ : "-- PostGIS not required" }

    -- Grant default privileges for future objects
    ALTER DEFAULT PRIVILEGES IN SCHEMA "${sqlQuoteLiteral(serviceName)}" GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO "${sqlQuoteLiteral(dbUser)}";
    ALTER DEFAULT PRIVILEGES IN SCHEMA "${sqlQuoteLiteral(serviceName)}" GRANT USAGE, SELECT ON SEQUENCES TO "${sqlQuoteLiteral(dbUser)}";
    ALTER DEFAULT PRIVILEGES IN SCHEMA "${sqlQuoteLiteral(serviceName)}" GRANT EXECUTE ON FUNCTIONS TO "${sqlQuoteLiteral(dbUser)}";

    -- Grant privileges on existing objects
    GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA "${sqlQuoteLiteral(serviceName)}" TO "${sqlQuoteLiteral(dbUser)}";
    GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA "${sqlQuoteLiteral(serviceName)}" TO "${sqlQuoteLiteral(dbUser)}";
    GRANT EXECUTE ON ALL FUNCTIONS IN SCHEMA "${sqlQuoteLiteral(serviceName)}" TO "${sqlQuoteLiteral(dbUser)}";

    -- Set schema owner
    ALTER SCHEMA "${sqlQuoteLiteral(serviceName)}" OWNER TO "${sqlQuoteLiteral(dbUser)}";
    """ // End of sqlScript

    echo "Ensuring Database and User exist for service: ${serviceName}"
    try {
        // Write SQL to a temporary file instead of using shell heredoc
        writeFile file: tempSqlFile, text: sqlScript

        // Execute the SQL using the temporary file
        sh """
           set -e
           ${kubectlCmd} cp ${tempSqlFile} ${dbNamespace}/postgres-db-0:/tmp/${tempSqlFile} -c postgres
           ${kubectlCmd} exec -n ${dbNamespace} postgres-db-0 -c postgres -- psql -v ON_ERROR_STOP=1 -U postgres -d postgres -f /tmp/${tempSqlFile}
           ${kubectlCmd} exec -n ${dbNamespace} postgres-db-0 -c postgres -- rm /tmp/${tempSqlFile} || true
        """

        // Clean up the local temporary file
        sh "rm ${tempSqlFile}"

        echo "Database/User/Schema/Permissions ensured for ${serviceName}."
    } catch (e) {
         error("Failed to ensure database/user for service ${serviceName} in namespace ${dbNamespace}. Error: ${e.getMessage()}")
    }
}

def getPostgisRequirement(String serviceName) {
    def postgisMap = [
        users: 'true', products: 'true', geocoding: 'true', posts: 'true',
        services: 'true', properties: 'true', vehicles: 'true', deals: 'true', jobs: 'true'
    ]
    return postgisMap.getOrDefault(serviceName, 'false')
}

def rolloutRestart(String resourceType, String kubectlCmd, String resourceName = "--all", String namespace = env.K8S_NAMESPACE) {
    try {
        sh(script: "${kubectlCmd} rollout restart ${resourceType}/${resourceName} -n ${namespace}", label: "Restart ${resourceType} ${resourceName}")
    } catch (e) {
        echo "WARN: Failed to trigger rollout restart for ${resourceType} ${resourceName}: ${e.getMessage()}"
    }
}

pipeline {
    agent any

    parameters {
        choice(name: 'SERVICE', choices: ['all', 'activity', 'users', 'products', 'baskets', 'ordering', 'payments', 'search', 'notifications', 'comments', 'media', 'messages', 'categories', 'mailer', 'newsletters', 'offers', 'posts','shipping', 'support', 'wishlists','reviews','following','metrics','geocoding', 'assistants', 'vectors','services', 'qdrant','frontend', 'cosec','erp', 'scheduler', 'merchant','managers','reverse-proxy'], description: 'Service to build/deploy (or "all")')
        choice(name: 'ENVIRONMENT', choices: ['development', 'testing', 'staging', 'production'], description: 'Target environment')
        string(name: 'TAG', defaultValue: "${BUILD_NUMBER}", description: 'Docker image tag')
        booleanParam(name: 'RUN_TESTS', defaultValue: true, description: 'Run tests')
        booleanParam(name: 'PUSH_IMAGE', defaultValue: true, description: 'Push Docker image(s)')
        booleanParam(name: 'DEPLOY', defaultValue: true, description: 'Deploy the service(s)')
        string(name: 'KUBECONFIG_CREDENTIAL_ID', defaultValue: '', description: 'Optional: Jenkins credential ID for kubeconfig file')
        string(name: 'REGISTRY_CREDENTIAL_ID', defaultValue: 'registry-sfx-markt-de-creds', description: 'Jenkins credential ID for registry')
    }

    environment {
        REGISTRY = "registry.sfx-markt.de"
        GO_VERSION = "1.23.1"
        DOCKER_BUILDKIT = "1"
        // ADDED reverse-proxy to BUILDABLE_SERVICES_LIST
        BUILDABLE_SERVICES_LIST = 'activity users products baskets ordering payments search notifications comments media messages categories geocoding mailer newsletters offers posts shipping support wishlists frontend reviews following metrics assistants vectors cosec services erp merchant scheduler managers reverse-proxy'
        MICROSERVICES_LIST = 'activity users products baskets ordering payments search notifications comments media messages categories geocoding mailer newsletters offers posts shipping support wishlists reviews metrics assistants vectors following cosec erp services merchant managers scheduler'
        K8S_MANIFEST_DIR = "k8s"
        K8S_NAMESPACE = "my-project"
        KUBECONFIG_PATH = ""
    }

    stages {

        stage('Setup Environment') {
            steps {
                script {
                    if (params.KUBECONFIG_CREDENTIAL_ID) {
                        withCredentials([file(credentialsId: params.KUBECONFIG_CREDENTIAL_ID, variable: 'KUBECONFIG_SECRET_FILE')]) {
                            env.KUBECONFIG_PATH = KUBECONFIG_SECRET_FILE
                        }
                        echo "Using Kubeconfig from credentials: ${env.KUBECONFIG_PATH}"
                    } else {
                        echo "Using default kubectl context."
                    }
                }
            }
        }

        stage('Validate') {
             steps {
                 echo "Validation..."
                 sh 'docker --version'
                 script {
                     def kubectlCmdValidate = "kubectl"
                     if (env.KUBECONFIG_PATH) {
                         kubectlCmdValidate = "kubectl --kubeconfig ${env.KUBECONFIG_PATH}"
                     }
                     sh "${kubectlCmdValidate} version --client"
                 }
             }
        }

        stage('Test') {
            when {
                expression {
                    params.RUN_TESTS &&
                    params.SERVICE != 'all' &&
                    env.MICROSERVICES_LIST.split(' ').contains(params.SERVICE) // reverse-proxy is not a microservice, so won't run tests by default
                }
            }
            steps {
                 echo "Running tests for Go service: ${params.SERVICE}..."
                 script {
                     def serviceDir = "${WORKSPACE}/${params.SERVICE}"
                     if (fileExists(serviceDir)) {
                         sh """
                             set -e
                             echo "Executing tests in a Go container..."
                             docker run --rm -v "${WORKSPACE}:/middleman" -w "/middleman/${params.SERVICE}" \\
                                 golang:${env.GO_VERSION}-alpine \\
                                 sh -c "go test ./... -v || exit 1"
                             echo "Tests for ${params.SERVICE} completed."
                         """
                     } else {
                          echo "WARN: Test directory '${serviceDir}' not found. Skipping tests."
                     }
                 }
            }
        }

        stage('Build Docker Image') {
            // REMOVED when condition to allow reverse-proxy build if selected
            steps {
                script {
                    def servicesToBuild = []
                    def allBuildableServices = env.BUILDABLE_SERVICES_LIST.split(' ')
                    if (params.SERVICE == 'all') {
                        servicesToBuild.addAll(allBuildableServices)
                    } else if (allBuildableServices.contains(params.SERVICE)) {
                        servicesToBuild.add(params.SERVICE)
                    } else {
                        // This condition might now be hit if an unbuildable service is chosen directly
                        // and not 'all'. Consider if 'reverse-proxy' can be chosen directly if not buildable.
                        // For now, this assumes if it's in BUILDABLE_SERVICES_LIST, it's intended to be built.
                        echo "WARN: Service ${params.SERVICE} not in BUILDABLE_SERVICES_LIST, or unknown. Skipping build."
                        // error("Unknown service specified or service not in BUILDABLE_SERVICES_LIST: ${params.SERVICE}")
                    }

                    def parallelBuilds = [:]
                    servicesToBuild.each { svc ->
                        def dockerfilePath = ""
                        def buildArgs = ""
                        def dockerBuildFlags = "--pull"

                        if (svc == 'mailer') {
                            dockerfilePath = "docker/Dockerfile.mailer"
                            buildArgs = "--build-arg service=${svc}"
                            echo "[Build mailer] Using dedicated Dockerfile."
                        } else if (env.MICROSERVICES_LIST.split(' ').contains(svc)) {
                            dockerfilePath = "docker/Dockerfile.microservices"
                            buildArgs = "--build-arg service=${svc}"
                        } else if (svc == 'frontend') {
                            if (fileExists("${WORKSPACE}/Dockerfile.frontend")) {
                                dockerfilePath = "Dockerfile.frontend"
                            } else if (fileExists("${WORKSPACE}/docker/Dockerfile.frontend")) {
                                dockerfilePath = "docker/Dockerfile.frontend"
                            } else {
                                error("Dockerfile.frontend not found in either workspace root or docker/ directory")
                            }
                            def timestamp = new Date().getTime()
                            dockerBuildFlags += " --no-cache --pull=true --force-rm=true --build-arg CACHE_BUST=${timestamp}"
                            echo "[Build frontend] Adding full rebuild flags with timestamp ${timestamp} to force clean build with no caching."
                        } else if (svc == 'reverse-proxy') { // ADDED condition for reverse-proxy
                            dockerfilePath = "docker/Dockerfile.nginx" // Assuming this is the correct path from your provided files
                            // No specific buildArgs for the simple Nginx Dockerfile
                            echo "[Build reverse-proxy] Using Dockerfile: ${dockerfilePath}"
                        } else {
                            echo "WARN: Unknown service type logic for '${svc}'. Skipping build determination."
                            return // Skips this iteration of the loop
                        }

                        if (dockerfilePath && fileExists("${WORKSPACE}/${dockerfilePath}")) {
                            parallelBuilds[svc] = {
                                stage("Build ${svc}") {
                                    try {
                                        def imageName = "${env.REGISTRY}/${svc}:${params.TAG}"
                                        echo "Building ${svc} image: ${imageName} using ${dockerfilePath}..."
                                        if (svc == 'frontend') {
                                            sh """
                                                echo "Removing any existing frontend images to ensure clean build..."
                                                docker image rm -f ${imageName} || true
                                                docker builder prune -f
                                            """
                                        }
                                        def fullBuildCommand = """
                                            set -e
                                            docker build ${dockerBuildFlags} -t ${imageName} \\
                                                -f ${WORKSPACE}/${dockerfilePath} ${buildArgs} ${WORKSPACE}
                                        """
                                        echo "Executing build command for ${svc}: ${fullBuildCommand}"
                                        sh (script: fullBuildCommand, label: "Build Docker ${svc}")
                                    } catch (e) {
                                        error("Failed to build ${svc}: ${e.getMessage()}")
                                    }
                                }
                            }
                        } else {
                            if (dockerfilePath) { // Only warn if a dockerfilePath was determined but not found
                                echo "WARN: Dockerfile '${WORKSPACE}/${dockerfilePath}' not found for service '${svc}'. Skipping build."
                            }
                            // If dockerfilePath was empty (due to unknown service type), it was already logged.
                        }
                    }

                    if (!parallelBuilds.isEmpty()) {
                        parallel parallelBuilds
                    } else {
                        if (!servicesToBuild.isEmpty()) { // Only error if services were selected but no builds were possible
                           error("No Dockerfiles found or build logic defined for the selected, buildable service(s): ${servicesToBuild.join(', ')}")
                        } else {
                           echo "No services selected or eligible for build."
                        }
                    }
                }
            }
        }

        stage('Push to Registry') {
            // REMOVED when condition to allow reverse-proxy push if built
            steps {
                script {
                    def servicesToAttemptPush = []
                    def allBuildableServices = env.BUILDABLE_SERVICES_LIST.split(' ')

                    if (params.SERVICE == 'all') {
                        servicesToAttemptPush.addAll(allBuildableServices)
                    } else if (allBuildableServices.contains(params.SERVICE)) {
                        // Check if the service was intended to be built (has Dockerfile logic)
                        // This logic could be more robust by checking if the build stage actually produced an image
                        // For now, assume if it's buildable and selected, an attempt should be made.
                        servicesToAttemptPush.add(params.SERVICE)
                    }


                    if (!servicesToAttemptPush.isEmpty()) {
                        echo "Logging in to registry ${env.REGISTRY}..."
                        withCredentials([usernamePassword(credentialsId: params.REGISTRY_CREDENTIAL_ID, usernameVariable: 'REG_USER', passwordVariable: 'REG_PASS')]) {
                            try {
                                sh(script: "echo \$REG_PASS | docker login ${env.REGISTRY} -u \$REG_USER --password-stdin", label: "Docker Login")
                                def parallelPushes = [:]
                                servicesToAttemptPush.each { svc ->
                                     parallelPushes[svc] = {
                                         stage("Push ${svc}") {
                                             try {
                                                 def imageName = "${env.REGISTRY}/${svc}:${params.TAG}"
                                                 // Check if image exists locally before pushing
                                                 def imageExists = sh(script: "docker image inspect ${imageName} > /dev/null 2>&1", returnStatus: true) == 0
                                                 if (imageExists) {
                                                     echo "Pushing ${imageName}..."
                                                     sh "docker push ${imageName}"
                                                 } else {
                                                     echo "Skipping push for ${imageName}, image not found locally (likely build failed or skipped)."
                                                 }
                                             } catch (e) {
                                                 error("Failed to push ${svc}: ${e.getMessage()}")
                                             }
                                         }
                                     }
                                }
                                if (!parallelPushes.isEmpty()) {
                                   parallel parallelPushes
                                }
                            } catch (e) {
                                error("Docker login failed: ${e.getMessage()}")
                            }
                        }
                    } else {
                        echo "No services identified for push."
                    }
                }
            }
            post {
                always {
                    echo "Logging out from registry ${env.REGISTRY}..."
                    sh "docker logout ${env.REGISTRY} || true"
                }
            }
        }

        stage('Deploy to Kubernetes') {
            when { expression { params.DEPLOY } }
            steps {
                script {
                    def kubectlCmd = "kubectl"
                    if (env.KUBECONFIG_PATH) {
                        kubectlCmd = "kubectl --kubeconfig ${env.KUBECONFIG_PATH}"
                    }

                    echo "Deploying to Kubernetes namespace: ${env.K8S_NAMESPACE} using tag: ${params.TAG}"
                    def manifestBaseDir = "${WORKSPACE}/${env.K8S_MANIFEST_DIR}"
                    if (!fileExists(manifestBaseDir)) {
                        error("Kubernetes manifest directory '${manifestBaseDir}' not found!")
                    }

                    echo "Step 1: Applying Namespaces..."
                    applyManifests("${manifestBaseDir}/00-namespaces.yaml", kubectlCmd, false)

                    echo "Step 2: Applying Infrastructure Services (DB, NATS, Redis, Minio, Qdrant)..."
                    applyManifests("${manifestBaseDir}/02-database/", kubectlCmd)
                    applyManifests("${manifestBaseDir}/03-messaging/", kubectlCmd)
                    applyManifests("${manifestBaseDir}/04-redis/", kubectlCmd)
                    applyManifests("${manifestBaseDir}/05-minio/", kubectlCmd, false)
                    applyManifests("${manifestBaseDir}/06-common/", kubectlCmd)
                    
                    // Deploy Qdrant only if qdrant service is selected, vectors service is selected, or 'all'
                    if (params.SERVICE == 'all' || params.SERVICE == 'qdrant' || params.SERVICE == 'vectors') {
                        echo "Applying Qdrant Vector Database..."
                        applyManifests("${manifestBaseDir}/14-vectors/", kubectlCmd)
                    }

                    echo "Step 3: Applying Base Reverse Proxy Manifest and Ingress..." // Changed "Reverse Proxy" to "Base Reverse Proxy Manifest"
                    // Apply reverse-proxy manifest first, as other services might depend on its existence conceptually
                    // or if it were to be an actual proxy layer for other services.
                    // Its image will be updated later in the loop if it's built by this pipeline.
                    applyManifests("${manifestBaseDir}/07-reverse-proxy/nginx-deployment.yaml", kubectlCmd, true, env.K8S_NAMESPACE)
                    applyManifests("${manifestBaseDir}/08-middleware/", kubectlCmd)


                    echo "Step 4: Ensure DB Structure and Deploy Application Services..."
                    def servicesToDeploy = []
                    def allBuildableServices = env.BUILDABLE_SERVICES_LIST.split(' ') // Now includes reverse-proxy
                    def allMicroservices = env.MICROSERVICES_LIST.split(' ')


                    if (params.SERVICE == 'all') {
                        echo "Ensuring DB structures for all microservices..."
                        allMicroservices.each { svc ->
                            // Skip database setup for vectors service (uses Qdrant, not PostgreSQL)
                            if (svc != 'vectors') {
                                def needsPostgis = getPostgisRequirement(svc)
                                ensureServiceDb(svc, needsPostgis, kubectlCmd)
                            }
                        }
                        servicesToDeploy.addAll(allBuildableServices) // This will include reverse-proxy now
                    } else if (allBuildableServices.contains(params.SERVICE)) { // Handles individual service deployment, including reverse-proxy
                        if (allMicroservices.contains(params.SERVICE) && params.SERVICE != 'vectors') {
                             echo "Ensuring DB structure for single service: ${params.SERVICE}"
                             def needsPostgis = getPostgisRequirement(params.SERVICE)
                             ensureServiceDb(params.SERVICE, needsPostgis, kubectlCmd)
                        }
                        servicesToDeploy.add(params.SERVICE)
                    } else {
                        error("Service parameter '${params.SERVICE}' is not in the known buildable services list.")
                    }

                    servicesToDeploy.each { svc ->
                        echo "--- Processing deployment for ${svc} ---"
                        def manifestPath = ""
                        def resourceType = "deployment"
                        def containerName = svc // Default container name to service name

                        if (allMicroservices.contains(svc)) {
                            manifestPath = "${manifestBaseDir}/09-services/${svc}.yaml"
                        } else if (svc == 'frontend') {
                            manifestPath = "${manifestBaseDir}/10-frontend/frontend.yaml"
                            containerName = "frontend" // Specific container name for frontend
                            if (fileExists(manifestPath)) {
                                echo "Patching frontend deployment to use imagePullPolicy: Always"
                                sh """
                                    ${kubectlCmd} patch deployment frontend -n ${env.K8S_NAMESPACE} \
                                    --patch '{"spec":{"template":{"spec":{"containers":[{"name":"frontend","imagePullPolicy":"Always"}]}}}}' || true
                                """
                            }
                        } else if (svc == 'reverse-proxy') {
                            // Manifest already applied. Image update will happen below.
                            // We still need to define containerName for updateImageTag.
                            containerName = "proxy" // Container name in reverse-proxy deployment
                            manifestPath = "${manifestBaseDir}/07-reverse-proxy/nginx-deployment.yaml" // Needed for fileExists check
                        } else {
                            echo "WARN: Unknown service type '${svc}' in deployment loop. Skipping full processing."
                            // return // or continue if in Groovy 2.x+
                        }

                        if (manifestPath) { // If a manifest path was determined (even if already applied for reverse-proxy)
                            if (fileExists(manifestPath)) {
                                 // For services other than reverse-proxy, apply their specific main manifest here.
                                 // reverse-proxy's main deployment manifest was applied earlier.
                                 if (svc != 'reverse-proxy') {
                                     echo "Applying base manifest ${manifestPath} for ${svc}..."
                                     applyManifests(manifestPath, kubectlCmd, true, env.K8S_NAMESPACE)
                                 }

                                 // Now, update the image for all services in servicesToDeploy (including reverse-proxy)
                                 def fullImageName = "${env.REGISTRY}/${svc}:${params.TAG}"
                                 echo "Setting image for ${resourceType}/${svc} container ${containerName} to ${fullImageName}..."
                                 updateImageTag(resourceType, svc, containerName, fullImageName, kubectlCmd, env.K8S_NAMESPACE)

                                 if (svc == 'frontend') {
                                     echo "Forcing restart of frontend deployment to ensure latest image is used"
                                     rolloutRestart(resourceType, kubectlCmd, svc, env.K8S_NAMESPACE)
                                 }
                            } else {
                                 error("Manifest file '${manifestPath}' not found for service '${svc}'.")
                            }
                        } else if (!allBuildableServices.contains(svc)) {
                            // This case handles if a service name was in servicesToDeploy but didn't match any known type for manifestPath
                            echo "WARN: No manifest path logic for service '${svc}'. Skipping manifest application and image update."
                        }
                        echo "--- Finished deployment processing for ${svc} ---"
                    }

                    // Decide if reverse-proxy restart is needed.
                    // If 'all' services were deployed, or if a microservice was deployed (implying potential backend changes)
                    // or if reverse-proxy itself was deployed.
                    if (params.SERVICE == 'all' || allMicroservices.contains(params.SERVICE) || params.SERVICE == 'reverse-proxy') {
                         echo "Triggering rollout restart for reverse-proxy to pick up any upstream changes or its own update..."
                         rolloutRestart("deployment", kubectlCmd, "reverse-proxy", env.K8S_NAMESPACE)
                    }

                    echo "Kubernetes deployment steps completed."
                }
            }
        }
    }

    post {
        always {
            echo "Pipeline finished."
            cleanWs()
        }
        success {
            echo "Pipeline completed successfully!"
        }
        failure {
            echo "Pipeline failed. Check logs for details."
        }
    }
}