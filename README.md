# AOspace-Upgrade

English | [简体中文](./README_cn.md)

## Introduce

On-demand startup, mainly responsible for server-side upgrades

## Build

### Prepare environment

- docker (>=18.09)
- git
- golang 1.18 +

### Code Download

```shell
git clone git@github.com:ao-space/space-upgrade.git
```

### Build image

go into project root path and run:

```shell
docker build -t local/space-upgrade:{tag} . 
````

The tag parameter can be modified to be consistent with the docker-compose.yml that is running on the server as a whole.

## Run

Generally, it does not run permanently, but only starts when the upgrade is performed, and exits when the upgrade is complete.

To run it manually, run

```shell
docker exec -it aospace-all-in-one docker-compose -f /etc/ao-space/aospace-upgrade.yml up -d
```

## Contribution Guidelines

Contributions to this project are very welcome. Here are some guidelines and suggestions to help you get involved in the project.

[Contribution Guidelines](https://github.com/ao-space/ao.space/blob/dev/docs/en/contribution-guidelines.md)

## Contact us

- Email: <developer@ao.space>
- [Official Website](https://ao.space)
- [Discussion group](https://slack.ao.space)

## Thanks for your contribution

Finally, thank you for your contribution to this project. We welcome contributions in all forms, including but not limited to code contributions, issue reports, feature requests, documentation writing, etc. We believe that with your help, this project will become more perfect and stronger.
