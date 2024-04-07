#!/bin/bash

# get kubectl pod status
# Usage: ./status.sh <kind>
# kind: top, node, pod

function top(){
    all_pods=$(kubectl get pods | grep -v NAME | grep Running | cut -d ' ' -f 1)
    for pod in $all_pods
    do
        echo "Pod: $pod"
        kubectl top pod $pod
    done
}

function pod(){
    kubectl get pod
}

function node(){
    kubectl get node
    # sort by node name
    kubectl get pod -o=custom-columns=NAME:.metadata.name,STATUS:.status.phase,NODE:.spec.nodeName | sort -k3
}

function help(){
cat <<EOF
  Usage: ./status.sh <kind>

    kind    : top, node, pod
EOF
}

function main(){
    if [ $# -lt 1 ]; then
        help
        exit 1
    fi

    option=$1
    case $option in
        top)
            top
            ;;
        node)
            node
            ;;
        pod)
            pod
            ;;
        *)
            echo "Invalid option: $option"
            help
            ;;
    esac
}

main $@