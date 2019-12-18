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
    docker.image('golang:1.13').inside("--user root --privileged -v /home/ubuntu/test-${ghprbPullId}:/home/test-${ghprbPullId}") {
        stage('Build') {
            sh "cd /home/test-${ghprbPullId}/cmd/manager/ && go build main.go"
        }
    }
    stage('Cleanup'){
        dir("/home/ubuntu/test-${ghprbPullId}"){
            cleanWs()
        }
    }
}