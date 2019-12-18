node('multicloud-node') {
    stage('Get repository'){
        dir("/home/ubuntu/test-${ghprbPullId}"){
            checkout([$class: 'GitSCM', branches: [[name: "*/${ghprbSourceBranch}"]], 
                    doGenerateSubmoduleConfigurations: false, extensions: [[$class: 'PreBuildMerge', 
                    options: [mergeRemote: 'origin', mergeTarget: "${ghprbTargetBranch}"]]], 
                    submoduleCfg: [], 
                    userRemoteConfigs: [[credentialsId: 'acitalkey', 
                    url: 'git@github.com:Juniper/contrail-operator.git']]])
        }
    }
    docker.image('golang:1.13').inside("--user root -v /home/ubuntu/test-${ghprbPullId}:/home/test-${ghprbPullId}") {
        stage('Build and test') {
            try {
                sh "cd /home/test-${ghprbPullId} && go build cmd/manager/main.go"
                sh "cd /home/test-${ghprbPullId} && go test -race -v ./pkg/..."
            } finally {
                sh "cd /home/test-${ghprbPullId} && rm -rf *"
            }
        }
    }
}