Pipeline CI GO
======

I wrote a simple dummy Java app but and decided that I didn't want to use Jenkins nor CircleCI for our Continous Integration, instead, I've decided to write a simple Continous Integration system :)

Sample Application:
https://github.com/andrehluizsilva/pipeline-sample

Inside the repository, I've added a `pipeline.yml` file. The CI tool will lookup this file to execute the requested pipeline that will be implemented here.

The will script do:

1. Fetch git code (https://github.com/andrehluizsilva/pipeline-sample)
2. Parse the `pipeline.yml`
3. Run selected pipeline


Pipeline.yml Format
======
The pipeline file has only three main hashes.

Branch: Repository branch that should be used to execute the code.

Tasks: Saved commands that can be used to create a pipeline

Pipelines: Group of ordered tasks

Script Arguments
======

1. Pipeline name, E.g: build
2. Git repository URL


Examples
=====
Should fetch git code and execute `build` pipeline:
```shell
go build ./pipeline.go
./pipeline build https://github.com/andrehluizsilva/pipeline-sample.git
```

Testing
=====
To test the script on your personal computer you must have `git`, `maven`, `unzip` and `java` installed.

If you don't have it installed, there is `Dockerfile` in this repo that you can use to run your script.
