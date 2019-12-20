node('multicloud-node') {
    docker.image('golang:1.13').inside("--user root") {
        stage('Build and test') {
            checkout([$class: 'GitSCM', branches: [[name: "*/${ghprbSourceBranch}"]], 
                    doGenerateSubmoduleConfigurations: false, extensions: [[$class: 'PreBuildMerge', 
                    options: [mergeRemote: 'origin', mergeTarget: "${ghprbTargetBranch}"]]], 
                    submoduleCfg: [], 
                    userRemoteConfigs: [[credentialsId: 'acitalkey', 
                    url: 'git@github.com:Juniper/contrail-operator.git']]])

            sh "go build cmd/manager/main.go"
            sh "go test -race -v ./pkg/..."
        }
    }
}