# Verv CLI
## verv - is a simple cli tool for managing Go projects

### Installation
```shell
wget -qO- https://github.com/vervstack/verv/releases/latest/download/install.sh | bash
```
Detects your OS/arch (Linux/macOS, amd64/arm64) and installs the matching binary to `/usr/local/bin/verv`.

Or, with a Go toolchain:
```shell
go install go.vervstack.ru/verv@latest
```


### Features: 
  - create and manage a Golang projects

### TODO
 - manage cloud infrastructure

### Configuration
Can be used in order to override default settings. All fields are optional to specify

#### Via conf file
You can create verv.yaml file and specify options you'd like to override 
- Either put config next to verv binary  
```
  $GOPATH/bin/verv.yaml
```
- or pass config via flag 
```
  verv [COMMAND] [ARGUMENTS] --cfg ./verv.yaml
```

#### Via environment variables
Alternatively (or additionally) you can specify environment variable(s)
that start with "VERV_" and followed by field name.

##### Configuration structure (all fields are optional)
- **default_project_git_path** - URL path to git system where this package will be available to fetch
```shell
  export VERV_DEFAULT_PROJECT_GIT_PATH=github.com/vervstack
```
- env - object. defines environment handling
  - **path_to_main** - path to main file. Used for project scan. Entrypoint when starts
  ```shell
    export VERV_PATH_TO_MAIN=cmd/proj_name/main.go
  ```
  - **path_to_config** - path to project config file. Used for project scan 
  ```shell
    export VERV_PATH_TO_CONFIG=config/dev.yaml
  ```

#### Config example can be found in internal/config/verv.yaml
```yaml
env:
  path_to_main: 'cmd/proj_name/main.go'
  path_to_config: 'config/dev.yaml'

default_project_git_path: github.com/vervstack
```
#### SPECIAL VARS
Reading examples you might notice that there is a word **proj_name** used to define a project name. 
This is not just an example. In order to make tool more useful some sequences are predefined for internal use.
- proj_name - can be used in order to put actual name of project somewere
```text
example: 
    When you set variable 
    env.path_to_main to "cmd/proj_name/main.go" 
    and create a project named Verv via 
    
    verv init
    
    it creates project with main file at
    "./github.com/vervstack/verv/cmd/verv/main.go"
    
    instead of "./github.com/vervstack/verv/cmd/proj_name/main.go"
```
### THE LAST, NOT THE LEAST on project creation
- ALL project created via verv tool require some url at the begging.
- Bear in mind that project is being created withing the whole path, meaning github.com/vervstack/ folders being created
