node('contrail-operator-node') {
    docker.image('kaweue/testrunner:2').inside("--user root -v /var/run/docker.sock:/var/run/docker.sock --net host") {
        stage('Build and test') {
            checkout([$class: 'GitSCM', branches: [[name: "*/${ghprbSourceBranch}"]], 
                    doGenerateSubmoduleConfigurations: false, extensions: [[$class: 'PreBuildMerge', 
                    options: [mergeRemote: 'origin', mergeTarget: "${ghprbTargetBranch}"]]], 
                    submoduleCfg: [], 
                    userRemoteConfigs: [[credentialsId: 'acitalkey', 
                    url: 'git@github.com:Juniper/contrail-operator.git']]])

            sh "go build cmd/manager/main.go"
            sh "go test -race -v ./pkg/..."

            try {
                sh "./test/env/create_k8s_cluster.sh ${ghprbPullId} ${registry}"
                sh "kubectl create namespace contrail"
                sh 'operator-sdk test local ./test/e2e/ --namespace contrail --go-test-flags "-failfast -v -timeout=30m" --up-local '
            } finally {
                sh "kubectl delete namespace contrail"
                sh "kind delete cluster --name=${ghprbPullId}"
            }
        }
    }
}