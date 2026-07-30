## Setting Up your Development Environment

If you plan to modify the source code and build your own version of the plugin, follow the steps in this section to set up the development environment.

This fork's locally built and packaged artifacts are unsigned. The signature is
based on a hash of the distribution, so Grafana must be configured to allow this
plugin to load unsigned.

For the plugin to load, add the following configuration parameter to the  `/etc/grafana/grafana.ini` file in the [plugins] section:

`allow_loading_unsigned_plugins = vertica-grafana-datasource`

![Allow Unsigned Plugin](https://raw.githubusercontent.com/vertica/vertica-grafana-datasource/main/src/img/allow-unsigned-plugin.png)

### Prerequisites
* [Grafana](https://grafana.com/docs/grafana/latest/installation/) (version 10.x)
* [Vertica](https://www.vertica.com/download/vertica/)
* [NodeJS](https://nodejs.org/en/download/package-manager/) 20
* [Yarn Classic](https://classic.yarnpkg.com/en/docs/install/) 1.22.19
* [Go](https://golang.org/doc/install) 1.26.3
* [Mage](https://magefile.org/)

**Note:** This plugin is tested on **Linux (Ubuntu)** with the following versions:

* Grafana - v13.0.0
* NodeJS -  v20.20.2
* Yarn - v1.22.19
* Npm - v10.8.0
* Go - v1.26.3
* React - v19

 
### Installing Vertica Grafana Data Source Plugin 
**Note:** If you have the older version of the Vertica Grafana plugin, remove it using grafana-cli:
     
      grafana-cli plugins remove vertica-grafana-datasource
      
      
1. Clone the new Vertica Grafana data source repository into the plugins directory. 
2. Go to the Vertica plugin folder. 
3. To build the plugin 
   1. Install dependencies: 
   
       ```bash 
          yarn install 
       ``` 
   2. Build plugin in production mode: 
   
       ```bash 
          yarn build 
       ``` 
   3. Build the Linux AMD64 backend plugin binary used by this fork's release:
   
       ```bash
          go run github.com/magefile/mage@v1.17.2 -v build:linux
       ```
   4. Package the reproducible unsigned Linux AMD64 release:

       ```bash
          rtk bash scripts/package-release.sh v3.2.2
       ```

      The ZIP is written to
      `release/vertica-grafana-datasource-3.2.2.zip`.
4. Restart the Grafana server.
