node('contrail-operator-node') {
    docker.image('kaweue/testrunner:2').inside("--user root -v /var/run/docker.sock:/var/run/docker.sock --net host") {
        stage('Build and test') {
            checkout([$class: 'GitSCM', branches: [[name: "*/${ghprbSourceBranch}"]], 
                    doGenerateSubmoduleConfigurations: false, extensions: [[$class: 'PreBuildMerge', 
                    options: [mergeRemote: 'origin', mergeTarget: "${ghprbTargetBranch}"]]], 
                    submoduleCfg: [], 
                    userRemoteConfigs: [[credentialsId: 'acitalkey', 
                    url: 'git@github.com:Juniper/contrail-operator.git']]])

            try {
                sh "kind delete cluster --name=${testEnvPrefix}${ghprbPullId} || true"
                sh "./test/env/create_k8s_cluster.sh ${testEnvPrefix}${ghprbPullId} ${registry}"
                sh "kubectl create namespace contrail"
                sh 'BUILD_SCM_REVISION=`echo "${ghprbActualCommit}" | head -c 7` BUILD_SCM_BRANCH=${GIT_BRANCH} ./test/env/run_e2e_tests.sh'
            } finally {
                sh "kubectl delete namespace contrail"
                sh "kind delete cluster --name=${testEnvPrefix}${ghprbPullId}"
            }
        }
    }
}