# App accelerator Tanzu CLI plugin

This plugin lets you manage your accelerator resources using the tanzu CLI

# Commands

## Create

```
tanzu accelerator create my-accelerator-name --gitRepoUrl http://www.repourl.com --gitBranch main
```

### Update

```
tanzu accelerator update existing-accelerator-name --description "Lorem Ipsum"
```

## Delete

```
tanzu accelerator delete existing-accelerator-name
```

## Get

```
tanzu accelerator get existing-accelerator-name
```

## List

```
tanzu accelerator list
```