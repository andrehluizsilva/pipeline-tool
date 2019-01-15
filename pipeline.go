package main

import (
    "fmt"
    "gopkg.in/yaml.v2"
    "gopkg.in/src-d/go-git.v4"
    "gopkg.in/src-d/go-git.v4/plumbing"
    "io/ioutil"
    "os"
    "os/exec"
    "strings"
    "reflect"
)

// checkArgs should be used to ensure the right command line arguments are
// passed before executing script.
func checkArgs(arg ...string) {
    if len(os.Args) != len(arg)+1 {
        showWarning("Usage: %s %s", os.Args[0], strings.Join(arg, " "))
        os.Exit(1)
    }
}

// checkIfError should be used to naively panics if an error is not nil.
func checkIfError(err error) {
    if err == nil {
        return
    }

    fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
    os.Exit(1)
}

// showInfo should be used to describe the example commands that are about to run.
func showInfo(format string, args ...interface{}) {
    //fmt.Printf("\x1b[34;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
    fmt.Printf("%s\n\n", fmt.Sprintf(format, args...))
}

// showWarning should be used to display a warning
func showWarning(format string, args ...interface{}) {
    //fmt.Printf("\x1b[36;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
    fmt.Printf("%s\n\n", fmt.Sprintf(format, args...))
}

func clone(url string) string {

    url_path := strings.Split(url, "/")
    repo_full_name := url_path[len(url_path)-1]
    repo_name := strings.Split(repo_full_name, ".git")[0]

    errDir := os.RemoveAll(repo_name)
    checkIfError(errDir)

    // Clone the given repository to the given directory
    showInfo("Cloning Repository: %s", url)

    _, err := git.PlainClone(repo_name, false, &git.CloneOptions{
        URL:      url,
        RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
        Progress: os.Stdout,
    })

    checkIfError(err)

    return repo_name
}

func checkout(branch string, path string) {
    // We instance a new repository targeting the given path (the .git folder)
    repo, err := git.PlainOpen(path)
    checkIfError(err)

    w, err := repo.Worktree()
    checkIfError(err)

    showInfo("Checkout branch: %s", branch)
    err = w.Checkout(&git.CheckoutOptions{
        Branch: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)),
    })
    checkIfError(err)
}

func getBranch(pipeline map[string]interface {}) string {
    for k, v := range pipeline {
        if k == "branch" {
            return reflect.ValueOf(v).String()
        }
    }

    return ""
}

func getTasks(pipeline map[string]interface {}) []interface{}{
    for k, v := range pipeline {
        if k == "tasks" {
            return v.([]interface{})
        }
    }

    return nil
}

func getTaskCommand(key string, tasks []interface{}) (result string, err error) {
    for _, task := range tasks {
        for k, v := range task.(map[interface{}]interface{}) {
            if k == key {
                for k1, v1 := range v.(map[interface{}]interface{}) {
                    if k1 == "cmd" {
                        return reflect.ValueOf(v1).String(), nil
                    }
                }
            }
        }
    }

    return "", fmt.Errorf("Task [%s] not found!!", key)
}

func getPipelineCommands(key string, pipelines []interface{}) (result []interface{}, err error) {
    for _, step := range pipelines {
        for k, v := range step.(map[interface{}]interface{}) {
            if k == key {
                return v.([]interface{}), nil
            }
        }
    }

    return nil, fmt.Errorf("Pipeline [%s] not found!!", key)
}

func getPipelines(pipeline map[string]interface {}) []interface{}{
    for k, v := range pipeline {
        if k == "pipelines" {
            return v.([]interface{})
        }
    }

    return nil
}

func runCommand(command string, dir string) {

    cmd := exec.Command("bash", "-c", command)
    cmd.Dir = dir

    stdoutStderr, err := cmd.CombinedOutput()
    checkIfError(err)
    fmt.Printf("%s\n", stdoutStderr)
}

func main() {

    checkArgs("<action>", "<git_repo>")

    action   := os.Args[1]
    git_repo := os.Args[2]

    repository := clone(git_repo)
    repo_dir := "./" + repository

    filename := repo_dir + "/pipeline.yml"
    showInfo("\nRunning pipeline: %s", filename)
    source, err := ioutil.ReadFile(filename)
    checkIfError(err)

    pipeline := make(map[string]interface{})

    erro := yaml.Unmarshal(source, &pipeline)
    checkIfError(erro)

    branch := getBranch(pipeline)
    checkout(branch, repo_dir)

    tasks := getTasks(pipeline)
    pipelines := getPipelines(pipeline)

    steps, err := getPipelineCommands(action, pipelines)
    checkIfError(err)
    showInfo("Pipeline Steps: %v", steps)

    for _, step := range steps {
        command, err := getTaskCommand(reflect.ValueOf(step).String(), tasks)
        checkIfError(err)
        showInfo("Running [%s] step: %s", step,  command)
        runCommand(command, repository)
    }
}