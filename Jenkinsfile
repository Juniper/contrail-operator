node('contrail-operator-node') {
    docker.image('kaweue/testrunner:6').inside("--user root -v /var/run/docker.sock:/var/run/docker.sock --net host") {
        stage('e2e test') {
            checkout([$class: 'GitSCM', branches: [[name: "*/${ghprbSourceBranch}"]], 
                    doGenerateSubmoduleConfigurations: false, extensions: [[$class: 'PreBuildMerge', 
                    options: [mergeRemote: 'origin', mergeTarget: "${ghprbTargetBranch}"]]], 
                    submoduleCfg: [], 
                    userRemoteConfigs: [[credentialsId: 'acitalkey', 
                    url: 'git@github.com:Juniper/contrail-operator.git']]])

            withCredentials([string(credentialsId: 'github-api-read-token', variable: 'TOKEN')]) {
                sh '''#!/bin/bash

                    while : ; do
                        check_suites=$(curl -u $TOKEN -s -H "Accept: application/vnd.github.antiope-preview+json" https://api.github.com/repos/Juniper/contrail-operator/commits/${ghprbActualCommit}/check-suites)
                        status=$(echo $check_suites | jq -cr '.check_suites[]|select(.app.slug == "google-cloud-build").status')
                        [[ $status != "completed" ]] && echo "Waiting for upstream job. Current status: $status" && sleep 5 && continue
                        conclusion=$(echo "$check_suites" | jq -cr '.check_suites[]|select(.app.slug == "google-cloud-build").conclusion')
                        [[ $conclusion == "success" ]] && break
                        echo "Upstream job failed with conclusion: $conclusion"
                        exit 1
                    done
                '''
            }

            try {
                sh "kind delete cluster --name=${testEnvPrefix}${ghprbPullId} || true"
                sh "./test/env/create_k8s_cluster.sh ${testEnvPrefix}${ghprbPullId} ${registry} ${numberOfNodes}"
                sh "kubectl create namespace contrail"
                sh 'BUILD_SCM_REVISION=`echo "${ghprbActualCommit}" | head -c 7` BUILD_SCM_BRANCH=${GIT_BRANCH} E2E_TEST_SUITE=${e2eTestSuite} ./test/env/run_e2e_tests.sh'
            } finally {
                sh "kubectl delete --ignore-not-found namespace contrail"
                sh "kind delete cluster --name=${testEnvPrefix}${ghprbPullId}"
            }
        }
    }
}
