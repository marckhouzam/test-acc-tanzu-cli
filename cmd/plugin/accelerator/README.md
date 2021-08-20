# App accelerator Tanzu CLI plugin

This plugin lets you manage your accelerator resources using the tanzu CLI

# Commands

## Create

```
tanzu accelerator create my-accelerator-name --git-repository http://www.repourl.com --git-branch main

created accelerator my-accelerator-name in namespace default
```

### Update

```
tanzu accelerator update existing-accelerator-name --description "Lorem Ipsum"

updated accelerator existing-accelerator-name in namespace default
```

## Delete

```
tanzu accelerator delete existing-accelerator-name

deleted accelerator existing-accelerator-name in namespace default
```

## Get

```
tanzu accelerator get existing-accelerator-name

NAME                        GIT REPOSITORY                                    BRANCH
existing-accelerator-name   https://github.com/example/existing-accelerator   main
```

## List

```
tanzu accelerator list

NAME               GIT REPOSITORY                                            BRANCH
engine-features    https://github.com/simple-starters/e2e-engine-features    main
new-accelerator    https://github.com/simple-starters/e2e-new-accelerator    main
podinfo            https://github.com/simple-starters/e2e-podinfo            main
spring-petclinic   https://github.com/simple-starters/e2e-spring-petclinic   main
```